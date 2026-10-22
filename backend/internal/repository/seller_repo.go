package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type SellerRepoExport struct{ db *pgxpool.Pool }

func NewSellerRepo(db *pgxpool.Pool) *SellerRepoExport { return &SellerRepoExport{db} }

type SellerStats struct {
	TotalProducts  int     `json:"total_products"`
	Pending        int     `json:"pending"`
	Approved       int     `json:"approved"`
	Expired        int     `json:"expired"`
	Rejected       int     `json:"rejected"`
	TotalRevenue   float64 `json:"total_revenue"`
	TotalProfit    float64 `json:"total_profit"`
	PendingRevenue float64 `json:"pending_revenue"`
	UnpaidFees     float64 `json:"unpaid_fees"`
	PaidFees       float64 `json:"paid_fees"`
	PendingOrders  int     `json:"pending_orders"`
	Tier           string  `json:"tier"`
	TotalValueSold float64 `json:"total_value_sold"`
}

const sqlSellerProductStats = `
	SELECT
	  COUNT(*)                                                          AS total,
	  COUNT(*) FILTER (WHERE status='pending')                          AS pending,
	  COUNT(*) FILTER (WHERE status='approved')                         AS approved,
	  COUNT(*) FILTER (WHERE status='expired')                          AS expired,
	  COUNT(*) FILTER (WHERE status='rejected')                         AS rejected,
	  COALESCE(SUM(total_fee_amount) FILTER (WHERE fee_paid=FALSE AND status IN ('approved','expired')), 0) AS unpaid_fees
	 FROM products WHERE seller_id=$1`

const sqlSellerFinancialStats = `
	SELECT
	  (SELECT COALESCE(total_revenue, 0) FROM seller_profit_summary WHERE seller_id=$1) AS revenue,
	  (SELECT COALESCE(net_profit, 0) FROM seller_profit_summary WHERE seller_id=$1)    AS profit,
	  (SELECT COALESCE(total_fees_paid, 0) FROM seller_profit_summary WHERE seller_id=$1) AS paid_fees,
	  (SELECT COALESCE(SUM(t.seller_gets),0)
	   FROM transactions t
	   WHERE t.seller_id=$1 AND t.status='pending')                     AS pending_revenue,
	  (SELECT COUNT(*) FROM order_requests
	   WHERE seller_id=$1 AND status='pending')                         AS pending_orders,
	  COALESCE(tier::text,'normal')                                     AS tier,
	  COALESCE(total_value_sold,0)                                      AS total_value_sold
	 FROM users WHERE id=$1`

const sqlSellerTransactions = `
	SELECT t.id, t.product_id, COALESCE(t.seller_id,0), t.buyer_name, t.buyer_phone,
	        t.amount, t.broker_fee, t.seller_gets, t.status,
	        COALESCE(t.notes,''), t.created_at, t.updated_at
	 FROM transactions t
	 JOIN products p ON p.id = t.product_id
	 WHERE p.seller_id = $1
	 ORDER BY t.created_at DESC`

const sqlSellerProfile = `
	SELECT id, name, phone, additional_phones, COALESCE(business_name,''), COALESCE(email,''),
	        rating, role, is_active,
	        COALESCE(tier::text,'normal'), COALESCE(total_value_sold,0), created_at, updated_at
	 FROM users WHERE id=$1`

const sqlSellerUpdateProfile = `
	UPDATE users SET name=$1, business_name=$2, email=NULLIF($3,''), additional_phones=$4, updated_at=NOW()
	 WHERE id=$5 AND role='seller'`

const sqlSellerDeleteProduct = `
	DELETE FROM products
	 WHERE id=$1 AND seller_id=$2
	   AND status NOT IN ('approved', 'sold')
	   AND id NOT IN (SELECT product_id FROM transactions WHERE status='completed')
	   AND id NOT IN (SELECT product_id FROM order_requests WHERE status IN ('pending','accepted'))`

const sqlCreateOrderRequest = `
	INSERT INTO order_requests (product_id, seller_id, buyer_name, buyer_phone, buyer_email, message, offer_price, buyer_tier_at_order)
	 SELECT $1, $2, $3, $4, $5, $6, $7, $8
	 WHERE EXISTS (SELECT 1 FROM products WHERE id=$1 AND status='approved')
	 RETURNING id`

const sqlListOrdersBySeller = `
	SELECT o.id, o.product_id, o.seller_id, o.buyer_name, o.buyer_phone,
	        COALESCE(o.buyer_email,''), COALESCE(o.message,''), o.offer_price,
	        o.status, COALESCE(o.seller_note,''), COALESCE(o.buyer_tier_at_order,'normal'),
	        o.created_at, o.updated_at,
	        p.title, p.price
	 FROM order_requests o
	 JOIN products p ON p.id = o.product_id
	 WHERE o.seller_id=$1
	 ORDER BY o.created_at DESC`

const sqlUpdateOrderStatus = `
	UPDATE order_requests
	 SET status=$1, seller_note=$2, updated_at=NOW()
	 WHERE id=$3 AND seller_id=$4`

const sqlGetProductSellerIDAndTier = `
	SELECT seller_id, COALESCE(access_tier,'normal')
	 FROM products WHERE id=$1 AND status='approved'`

const sqlGetProductSellerID = `
	SELECT seller_id FROM products WHERE id=$1 AND status='approved'`

const sqlListProductsWithFee = `
	SELECT p.id, p.seller_id, p.location_id, p.title, p.price,
	        COALESCE(p.description,''), p.category_id,
	        COALESCE(p.subcategory,''), COALESCE(p.item,''), p.condition,
	        p.status, COALESCE(p.admin_notes,''), p.created_at, p.updated_at,
	        p.expires_at,
	        p.listing_duration_days, p.startup_fee_percent,
	        p.weekly_increment_percent, p.total_fee_percent, p.total_fee_amount,
	        p.fee_paid, COALESCE(p.fee_payment_pending, FALSE)
	 FROM products p
	 WHERE p.seller_id=$1
	 ORDER BY p.created_at DESC`

const sqlListProductImages = `
	SELECT product_id, id, image_url, order_number, created_at
	 FROM product_images WHERE product_id = ANY(
	   SELECT id FROM products WHERE seller_id=$1
	 ) ORDER BY product_id, order_number`

const sqlGetProductFeeStatus = `
	SELECT status, fee_paid, COALESCE(fee_payment_pending,FALSE) FROM products WHERE id=$1 AND seller_id=$2`

const sqlRecordFeePayment = `
	UPDATE products
	 SET fee_payment_pending=true, admin_notes=$1, updated_at=NOW()
	 WHERE id=$2 AND seller_id=$3`

const sqlSubmitUpgradeRequest = `
	INSERT INTO upgrade_requests
	   (user_id, role, requested_tier, payment_method, payer_phone, reference, amount, notes)
	 VALUES ($1, 'seller', $2, $3, $4, $5, $6, NULLIF($7,''))
	 RETURNING id`

const sqlGetUpgradeRequests = `
	SELECT id, requested_tier, payment_method, payer_phone, reference, amount,
	        COALESCE(notes,''), status, created_at
	 FROM upgrade_requests
	 WHERE user_id=$1 AND role='seller'
	 ORDER BY created_at DESC
	 LIMIT 10`

const sqlListNotifications = `
	SELECT id, title, message, type, is_read, created_at
	 FROM notifications
	 WHERE user_id = $1
	 ORDER BY created_at DESC
	 LIMIT 50`

const sqlMarkNotificationRead = `
	UPDATE notifications SET is_read=TRUE WHERE id=$1 AND user_id=$2`

const sqlCreateSellerNotification = `
	INSERT INTO notifications (user_id, title, message, type) VALUES ($1, $2, $3, $4)`

// GetStats returns aggregated dashboard stats for a seller.
func (r *SellerRepoExport) GetStats(ctx context.Context, sellerID int) (*SellerStats, error) {
	s := &SellerStats{}
	err := r.db.QueryRow(ctx, sqlSellerProductStats, sellerID).Scan(
		&s.TotalProducts, &s.Pending, &s.Approved, &s.Expired, &s.Rejected, &s.UnpaidFees)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(ctx, sqlSellerFinancialStats, sellerID).Scan(
		&s.TotalRevenue, &s.TotalProfit, &s.PaidFees, &s.PendingRevenue, &s.PendingOrders, &s.Tier, &s.TotalValueSold)
	if err != nil {
		s.Tier = "normal"
	}
	return s, nil
}

// ListTransactionsBySeller returns all transactions for products owned by a seller.
func (r *SellerRepoExport) ListTransactionsBySeller(ctx context.Context, sellerID int) ([]models.Transaction, error) {
	rows, err := r.db.Query(ctx, sqlSellerTransactions, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var txns []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.ProductID, &t.SellerID, &t.BuyerName, &t.BuyerPhone,
			&t.Amount, &t.BrokerFee, &t.SellerGets, &t.Status, &t.Notes,
			&t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, rows.Err()
}

// GetProfile returns a seller's profile by ID.
func (r *SellerRepoExport) GetProfile(ctx context.Context, sellerID int) (*models.User, error) {
	u := &models.User{}
	var extraPhones *[]string
	err := r.db.QueryRow(ctx, sqlSellerProfile, sellerID).Scan(
		&u.ID, &u.Name, &u.Phone, &extraPhones, &u.BusinessName, &u.Email,
		&u.Rating, &u.Role, &u.IsActive, &u.Tier, &u.TotalValueSold, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if extraPhones != nil {
		u.AdditionalPhones = *extraPhones
	}
	return u, nil
}

// UpdateProfile updates name, business_name, email, and additional_phones for a seller.
func (r *SellerRepoExport) UpdateProfile(ctx context.Context, u *models.User) error {
	maxPhones := u.MaxAdditionalPhones()
	if maxPhones >= 0 && len(u.AdditionalPhones) > maxPhones {
		return fmt.Errorf("tier %s allows max %d additional phones", u.Tier, maxPhones)
	}
	var phones interface{}
	if len(u.AdditionalPhones) > 0 {
		phones = u.AdditionalPhones
	}
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_, err := r.db.Exec(ctx, sqlSellerUpdateProfile, u.Name, u.BusinessName, u.Email, phones, u.ID)
	return err
}

// DeleteProduct hard-deletes a seller's own product.
func (r *SellerRepoExport) DeleteProduct(ctx context.Context, productID, sellerID int) error {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	tag, err := r.db.Exec(ctx, sqlSellerDeleteProduct, productID, sellerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cannot delete: product not found, is active/sold, or has open orders")
	}
	return nil
}

// CreateOrderRequest inserts a new buyer order request.
func (r *SellerRepoExport) CreateOrderRequest(ctx context.Context, o *models.OrderRequest) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, sqlCreateOrderRequest,
		o.ProductID, o.SellerID, o.BuyerName, o.BuyerPhone, o.BuyerEmail, o.Message, o.OfferPrice, o.BuyerTierAtOrder,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("product is no longer available")
		}
		return 0, err
	}
	return id, nil
}

// ListOrdersBySeller returns all order requests for a seller, newest first.
func (r *SellerRepoExport) ListOrdersBySeller(ctx context.Context, sellerID int) ([]models.OrderRequest, error) {
	rows, err := r.db.Query(ctx, sqlListOrdersBySeller, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.OrderRequest
	for rows.Next() {
		var o models.OrderRequest
		if err := rows.Scan(
			&o.ID, &o.ProductID, &o.SellerID, &o.BuyerName, &o.BuyerPhone,
			&o.BuyerEmail, &o.Message, &o.OfferPrice,
			&o.Status, &o.SellerNote, &o.BuyerTierAtOrder,
			&o.CreatedAt, &o.UpdatedAt,
			&o.ProductTitle, &o.ProductPrice,
		); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// UpdateOrderStatus updates the status and optional seller note on an order request.
func (r *SellerRepoExport) UpdateOrderStatus(ctx context.Context, orderID, sellerID int, status, note string) error {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	tag, err := r.db.Exec(ctx, sqlUpdateOrderStatus, status, note, orderID, sellerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

// GetProductSellerIDAndTier returns the seller_id and required_buyer_tier for a product.
func (r *SellerRepoExport) GetProductSellerIDAndTier(ctx context.Context, productID int) (int, string, error) {
	var sellerID int
	var requiredTier string
	err := r.db.QueryRow(ctx, sqlGetProductSellerIDAndTier, productID).Scan(&sellerID, &requiredTier)
	return sellerID, requiredTier, err
}

// GetProductSellerID returns the seller_id for a product.
func (r *SellerRepoExport) GetProductSellerID(ctx context.Context, productID int) (int, error) {
	var sellerID int
	err := r.db.QueryRow(ctx, sqlGetProductSellerID, productID).Scan(&sellerID)
	return sellerID, err
}

// ListProductsWithFee returns seller products with fee, expiry, and images.
func (r *SellerRepoExport) ListProductsWithFee(ctx context.Context, sellerID int) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, sqlListProductsWithFee, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	idxMap := map[int]int{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID, &p.SellerID, &p.LocationID, &p.Title, &p.Price,
			&p.Description, &p.CategoryID, &p.Subcategory, &p.Item, &p.Condition,
			&p.Status, &p.AdminNotes, &p.CreatedAt, &p.UpdatedAt,
			&p.ExpiresAt,
			&p.ListingDurationDays, &p.StartupFeePercent,
			&p.WeeklyIncrementPercent, &p.TotalFeePercent, &p.TotalFeeAmount,
			&p.FeePaid, &p.FeePaymentPending,
		); err != nil {
			return nil, err
		}
		idxMap[p.ID] = len(products)
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return products, nil
	}
	imgRows, err := r.db.Query(ctx, sqlListProductImages, sellerID)
	if err != nil {
		return products, nil
	}
	defer imgRows.Close()
	for imgRows.Next() {
		var img models.ProductImage
		var pid int
		if err := imgRows.Scan(&pid, &img.ID, &img.ImageURL, &img.OrderNumber, &img.CreatedAt); err != nil {
			continue
		}
		if i, ok := idxMap[pid]; ok {
			products[i].Images = append(products[i].Images, img)
		}
	}
	return products, nil
}

// RecordListingFeePayment stores the payment reference on the product row.
func (r *SellerRepoExport) RecordListingFeePayment(ctx context.Context, sellerID int, req dto.ListingFeePaymentRequest) error {
	var status string
	var feePaid bool
	var feePaymentPending bool
	err := r.db.QueryRow(ctx, sqlGetProductFeeStatus, req.ProductID, sellerID).Scan(&status, &feePaid, &feePaymentPending)
	if err != nil {
		return fmt.Errorf("product not found or does not belong to your account")
	}
	if feePaid {
		return fmt.Errorf("listing fee has already been verified and recorded for this product")
	}
	if feePaymentPending {
		return fmt.Errorf("a payment claim is already pending admin verification for this product")
	}
	if status == "sold" || status == "expired" || status == "rejected" {
		return fmt.Errorf("cannot pay listing fee for a %s product", status)
	}
	note := fmt.Sprintf("[FEE PAYMENT] method=%s phone=%s ref=%s amount=%.2f",
		req.PaymentMethod, req.PayerPhone, req.Reference, req.Amount)
	if req.Notes != "" {
		note += " | " + req.Notes
	}
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_, err = r.db.Exec(ctx, sqlRecordFeePayment, note, req.ProductID, sellerID)
	return err
}

// SubmitUpgradeRequest inserts a new tier upgrade request for a seller.
func (r *SellerRepoExport) SubmitUpgradeRequest(ctx context.Context, sellerID int, requestedTier, paymentMethod, payerPhone, reference, notes string, amount float64) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, sqlSubmitUpgradeRequest,
		sellerID, requestedTier, paymentMethod, payerPhone, reference, amount, notes,
	).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "uq_upgrade_requests_one_pending") {
			return 0, fmt.Errorf("you already have a pending upgrade request — wait for admin review before submitting another")
		}
		return 0, err
	}
	return id, nil
}

// GetUpgradeRequests returns all upgrade requests for a seller, newest first.
func (r *SellerRepoExport) GetUpgradeRequests(ctx context.Context, sellerID int) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, sqlGetUpgradeRequests, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var id int
		var tier, method, phone, ref, notes, status string
		var amount float64
		var createdAt interface{}
		if err := rows.Scan(&id, &tier, &method, &phone, &ref, &amount, &notes, &status, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"id": id, "requested_tier": tier, "payment_method": method,
			"payer_phone": phone, "reference": ref, "amount": amount,
			"notes": notes, "status": status, "created_at": createdAt,
		})
	}
	return out, rows.Err()
}

// ListNotifications returns unread notifications for a seller, newest first.
func (r *SellerRepoExport) ListNotifications(ctx context.Context, sellerID int) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, sqlListNotifications, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var id int
		var title, message, ntype string
		var isRead bool
		var createdAt interface{}
		if err := rows.Scan(&id, &title, &message, &ntype, &isRead, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"id": id, "title": title, "message": message,
			"type": ntype, "is_read": isRead, "created_at": createdAt,
		})
	}
	return out, rows.Err()
}

// MarkNotificationRead marks a seller notification as read.
func (r *SellerRepoExport) MarkNotificationRead(ctx context.Context, notifID, sellerID int) error {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_, err := r.db.Exec(ctx, sqlMarkNotificationRead, notifID, sellerID)
	return err
}

// CreateSellerNotification inserts a notification for a seller (fire-and-forget safe).
func (r *SellerRepoExport) CreateSellerNotification(ctx context.Context, sellerID int, title, message, ntype string) error {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_, err := r.db.Exec(ctx, sqlCreateSellerNotification, sellerID, title, message, ntype)
	return err
}

// UpdatePreference updates a seller's listing preference.
func (r *SellerRepoExport) UpdatePreference(ctx context.Context, sellerID int, preference string) error {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	_, err := r.db.Exec(ctx, `UPDATE users SET seller_preference = $1, updated_at = NOW() WHERE id = $2 AND role = 'seller'`, preference, sellerID)
	return err
}
