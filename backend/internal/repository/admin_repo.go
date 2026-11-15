package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"strconv"
	"strings"
	"time"
)

func NewAdminRepo(db *pgxpool.Pool) *AdminRepo { return &AdminRepo{db} }

type AdminRepo struct{ db *pgxpool.Pool }

// appendCond writes " AND col=$N" to b and appends val to args, returning the new args slice.
func appendCond(b *strings.Builder, args []interface{}, col string, val interface{}) []interface{} {
	b.WriteString(col)
	b.WriteString(fmt.Sprintf("$%d", len(args)+1)) //nolint:gosec // G201: Building static SQL fragment with computed $N parameter placeholder
	return append(args, val)
}

func (r *AdminRepo) ListPending(ctx context.Context) ([]models.Product, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT p.id, p.seller_id, p.location_id, p.title, p.price,
		       COALESCE(p.description,''), p.category_id,
		       COALESCE(p.subcategory,''), COALESCE(p.item,''), p.condition,
		       p.status, COALESCE(p.admin_notes,''), p.created_at, p.updated_at
		FROM products p
		WHERE p.status = 'pending'
		ORDER BY p.created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products, err := scanProducts(rows)
	if err != nil {
		return nil, err
	}
	return r.attachImages(ctx, products)
}

func (r *AdminRepo) ListAllProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT p.id, p.seller_id, p.location_id, p.title, p.price,
		       COALESCE(p.description,''), p.category_id,
		       COALESCE(p.subcategory,''), COALESCE(p.item,''), p.condition,
		       p.status, COALESCE(p.admin_notes,''), p.created_at, p.updated_at
		FROM products p
		ORDER BY p.created_at DESC
		LIMIT 500`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products, err := scanProducts(rows)
	if err != nil {
		return nil, err
	}
	return r.attachImages(ctx, products)
}

// ListProductsByStatus returns products filtered by a specific status.
func (r *AdminRepo) ListProductsByStatus(ctx context.Context, status string) ([]models.Product, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT p.id, p.seller_id, p.location_id, p.title, p.price,
		       COALESCE(p.description,''), p.category_id,
		       COALESCE(p.subcategory,''), COALESCE(p.item,''), p.condition,
		       p.status, COALESCE(p.admin_notes,''), p.created_at, p.updated_at
		FROM products p
		WHERE p.status = $1
		ORDER BY p.created_at DESC
		LIMIT 500`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products, err := scanProducts(rows)
	if err != nil {
		return nil, err
	}
	return r.attachImages(ctx, products)
}

func (r *AdminRepo) ListAllProductsPaged(ctx context.Context, page, limit int) ([]models.Product, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	var total int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&total) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	rows, err :=                                                         //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT p.id, p.seller_id, p.location_id, p.title, p.price,
		       COALESCE(p.description,''), p.category_id,
		       COALESCE(p.subcategory,''), COALESCE(p.item,''), p.condition,
		       p.status, COALESCE(p.admin_notes,''), p.created_at, p.updated_at
		FROM products p
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	products, err := scanProducts(rows)
	if err != nil {
		return nil, 0, err
	}
	products, err = r.attachImages(ctx, products)
	return products, total, err
}

func (r *AdminRepo) ListProductsByStatusPaged(ctx context.Context, status string, page, limit int) ([]models.Product, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	var total int
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE status=$1`, status).Scan(&total) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	rows, err :=                                                                                 //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT p.id, p.seller_id, p.location_id, p.title, p.price,
		       COALESCE(p.description,''), p.category_id,
		       COALESCE(p.subcategory,''), COALESCE(p.item,''), p.condition,
		       p.status, COALESCE(p.admin_notes,''), p.created_at, p.updated_at
		FROM products p
		WHERE p.status = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	products, err := scanProducts(rows)
	if err != nil {
		return nil, 0, err
	}
	products, err = r.attachImages(ctx, products)
	return products, total, err
}

func (r *AdminRepo) ListSellersPaged(ctx context.Context, page, limit int) ([]models.User, int, error) {
	return r.ListSellersFilteredPaged(ctx, page, limit, "", "", "", "")
}

func (r *AdminRepo) ListSellersFilteredPaged(ctx context.Context, page, limit int, search, status, tier, strategy string) ([]models.User, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	whereClauses := []string{"(u.role='seller' OR u.seller_preference IN ('rent', 'both') OR EXISTS (SELECT 1 FROM user_capabilities uc WHERE uc.user_id = u.id AND uc.capability IN ('seller', 'landlord') AND uc.status = 'active'))"}
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(u.name ILIKE $%d OR u.phone ILIKE $%d OR u.business_name ILIKE $%d OR u.email ILIKE $%d)", argIdx, argIdx, argIdx, argIdx))
		args = append(args, "%"+search+"%")
		argIdx++
	}

	if status == "active" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.is_active = $%d", argIdx))
		args = append(args, true)
		argIdx++
	} else if status == "inactive" || status == "deactivated" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.is_active = $%d", argIdx))
		args = append(args, false)
		argIdx++
	}

	if tier != "" && tier != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.tier = $%d", argIdx))
		args = append(args, tier)
		argIdx++
	}

	if strategy == "sale" {
		whereClauses = append(whereClauses, `(
			EXISTS (SELECT 1 FROM products p WHERE p.seller_id = u.id AND p.listing_type = 'sale' AND p.status = 'approved')
			OR (
				NOT EXISTS (SELECT 1 FROM products p WHERE p.seller_id = u.id AND p.listing_type = 'rent' AND p.status = 'approved')
				AND COALESCE(u.seller_preference, 'sale') = 'sale'
			)
		)`)
	} else if strategy == "rent" {
		whereClauses = append(whereClauses, `(
			EXISTS (SELECT 1 FROM products p WHERE p.seller_id = u.id AND p.listing_type = 'rent' AND p.status = 'approved')
			OR (
				NOT EXISTS (SELECT 1 FROM products p WHERE p.seller_id = u.id AND p.listing_type = 'sale' AND p.status = 'approved')
				AND u.seller_preference = 'rent'
			)
		)`)
	} else if strategy == "both" {
		whereClauses = append(whereClauses, `(
			EXISTS (SELECT 1 FROM products p WHERE p.seller_id = u.id AND p.listing_type = 'sale' AND p.status = 'approved')
			AND EXISTS (SELECT 1 FROM products p WHERE p.seller_id = u.id AND p.listing_type = 'rent' AND p.status = 'approved')
		)`)
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users u WHERE %s", whereSQL)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf(`
		SELECT u.id, u.name, u.phone, COALESCE(u.business_name,''), COALESCE(u.email,''),
		       u.rating, u.role, u.is_active,
		       COALESCE(u.tier::text,'normal') AS tier,
		       COALESCE(u.total_value_sold,0) AS total_value_sold,
		       COALESCE(ps.net_profit, 0) AS net_profit,
		       COALESCE(u.deactivation_reason,'') AS deactivation_reason,
		       COALESCE(u.seller_preference,'sale') AS seller_preference,
		       COALESCE(spc.cnt, 0) AS sale_product_count,
		       COALESCE(rpc.cnt, 0) AS rent_product_count,
		       u.created_at, u.updated_at
		FROM users u
		LEFT JOIN seller_profit_summary ps ON ps.seller_id = u.id
		LEFT JOIN (SELECT seller_id, COUNT(*) AS cnt FROM products WHERE listing_type = 'sale' AND status = 'approved' GROUP BY seller_id) spc ON spc.seller_id = u.id
		LEFT JOIN (SELECT seller_id, COUNT(*) AS cnt FROM products WHERE listing_type = 'rent' AND status = 'approved' GROUP BY seller_id) rpc ON rpc.seller_id = u.id
		WHERE %s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d`, whereSQL, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.BusinessName, &u.Email,
			&u.Rating, &u.Role, &u.IsActive, &u.Tier, &u.TotalValueSold,
			&u.NetProfit, &u.DeactivationReason, &u.SellerPreference,
			&u.SaleProductCount, &u.RentProductCount,
			&u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if len(u.Capabilities) == 0 {
			u.Capabilities = []string{"buyer"}
			if u.Role == "seller" || u.SellerPreference == "sale" || u.SellerPreference == "both" {
				u.Capabilities = append(u.Capabilities, "seller")
			}
			if u.SellerPreference == "rent" || u.SellerPreference == "both" {
				u.Capabilities = append(u.Capabilities, "landlord")
			}
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *AdminRepo) ListAllTransactionsPaged(ctx context.Context, page, limit int) ([]models.Transaction, int, error) {
	return r.ListAllTransactionsFilteredPaged(ctx, page, limit, "", "", "", "", "", "")
}

func (r *AdminRepo) ListAllTransactionsFilteredPaged(ctx context.Context, page, limit int, search, status, paymentMethod, txType, dateFrom, dateTo string) ([]models.Transaction, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	whereClauses := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		if idVal, err := strconv.Atoi(search); err == nil && idVal > 0 {
			whereClauses = append(whereClauses, fmt.Sprintf("(t.id = $%d OR t.payment_reference ILIKE $%d OR t.buyer_name ILIKE $%d OR t.buyer_phone ILIKE $%d OR t.payer_phone ILIKE $%d)", argIdx, argIdx+1, argIdx+1, argIdx+1, argIdx+1))
			args = append(args, idVal, "%"+search+"%")
			argIdx += 2
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("(t.payment_reference ILIKE $%d OR t.buyer_name ILIKE $%d OR t.buyer_phone ILIKE $%d OR t.payer_phone ILIKE $%d)", argIdx, argIdx, argIdx, argIdx))
			args = append(args, "%"+search+"%")
			argIdx++
		}
	}

	if status != "" && status != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf("t.status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}

	if paymentMethod != "" && paymentMethod != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf("t.payment_method = $%d", argIdx))
		args = append(args, paymentMethod)
		argIdx++
	}

	if txType != "" && txType != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf("t.transaction_type = $%d", argIdx))
		args = append(args, txType)
		argIdx++
	}

	if dateFrom != "" {
		if tFrom, err := time.Parse("2006-01-02", dateFrom); err == nil {
			whereClauses = append(whereClauses, fmt.Sprintf("t.created_at >= $%d", argIdx))
			args = append(args, tFrom)
			argIdx++
		}
	}

	if dateTo != "" {
		if tTo, err := time.Parse("2006-01-02", dateTo); err == nil {
			tToEnd := tTo.Add(24*time.Hour - time.Nanosecond)
			whereClauses = append(whereClauses, fmt.Sprintf("t.created_at <= $%d", argIdx))
			args = append(args, tToEnd)
			argIdx++
		}
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions t WHERE %s", whereSQL)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf(`
		SELECT t.id, t.user_id, COALESCE(t.product_id, 0), COALESCE(t.seller_id, 0),
		       COALESCE(t.buyer_name, ''), COALESCE(t.buyer_phone, ''),
		       t.amount, COALESCE(t.fee_amount, 0), COALESCE(t.broker_fee, 0), COALESCE(t.seller_gets, 0),
		       COALESCE(t.transaction_type, 'listing_fee'), COALESCE(t.payment_method, ''),
		       COALESCE(t.payment_reference, ''), COALESCE(t.payer_phone, ''),
		       t.status, COALESCE(t.notes, ''), t.created_at, t.updated_at
		FROM transactions t
		WHERE %s
		ORDER BY t.created_at DESC
		LIMIT $%d OFFSET $%d`, whereSQL, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var txns []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.ProductID, &t.SellerID,
			&t.BuyerName, &t.BuyerPhone,
			&t.Amount, &t.FeeAmount, &t.BrokerFee, &t.SellerGets,
			&t.TransactionType, &t.PaymentMethod,
			&t.PaymentReference, &t.PayerPhone,
			&t.Status, &t.Notes, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		txns = append(txns, t)
	}
	return txns, total, rows.Err()
}

func (r *AdminRepo) attachImages(ctx context.Context, products []models.Product) ([]models.Product, error) {
	if len(products) == 0 {
		return products, nil
	}
	idIndex := make(map[int]int, len(products))
	productIDs := make([]int, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
		idIndex[p.ID] = i
	}
	imgRows, err := r.db.Query(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`SELECT id, product_id, image_url, order_number, created_at
		 FROM product_images WHERE product_id = ANY($1::int[])
		 ORDER BY product_id, order_number`, productIDs)
	if err != nil {
		return products, nil // images are non-fatal
	}
	defer imgRows.Close()
	for imgRows.Next() {
		var img models.ProductImage
		if err := imgRows.Scan(&img.ID, &img.ProductID, &img.ImageURL, &img.OrderNumber, &img.CreatedAt); err == nil {
			if i, ok := idIndex[img.ProductID]; ok {
				products[i].Images = append(products[i].Images, img)
			}
		}
	}
	return products, nil
}

// ReviewProduct updates product status to approved or rejected with admin notes.
// When approving, the listing clock starts NOW (not at creation time).
// Rejected products remain visible to the seller so they can see the notes and resubmit.
// Approving requires fee_paid=TRUE — admin must verify the listing fee payment first.
func (r *AdminRepo) ReviewProduct(ctx context.Context, productID, adminID int, status, notes string) error {
	// Gate approval on fee_paid — prevents approving a product whose fee hasn't been verified.
	if status == "approved" {
		var feePaid bool
		var feeAmount *float64
		_ = r.db.QueryRow(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
			`SELECT fee_paid, total_fee_amount FROM products WHERE id=$1`, productID,
		).Scan(&feePaid, &feeAmount)
		// Only block approval if there's actually a listing fee configured for the product.
		// Products with no fee (total_fee_amount IS NULL or 0) can be approved freely.
		if feeAmount != nil && *feeAmount > 0 && !feePaid {
			return fmt.Errorf("cannot approve: listing fee has not been verified yet — verify the payment in the Fee Payments panel first")
		}
	}

	var tag interface{ RowsAffected() int64 }
	var err error
	if status == "approved" {
		tag, err = r.db.Exec(ctx, `
			UPDATE products
			SET status=$1, admin_notes=$2, reviewed_by=$3, reviewed_at=NOW(), updated_at=NOW(),
			    expires_at = CASE
			      WHEN listing_duration_days IS NOT NULL
			      THEN NOW() + (listing_duration_days || ' days')::interval
			      ELSE NULL
			    END
			WHERE id=$4 AND status IN ('pending', 'rejected')`,
			status, notes, adminID, productID)
	} else {
		tag, err = r.db.Exec(ctx, `
			UPDATE products
			SET status=$1, admin_notes=$2, reviewed_by=$3, reviewed_at=NOW(), updated_at=NOW()
			WHERE id=$4 AND status IN ('pending', 'rejected')`,
			status, notes, adminID, productID)
	}
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product with id %d not found or is in terminal state ('sold', 'expired')", productID)
	}
	return nil
}

func (r *AdminRepo) DeleteProduct(ctx context.Context, productID int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM products WHERE id=$1`, productID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	return err
}

// GetProductSellerID returns the seller_id for a product (used for seller notifications).
func (r *AdminRepo) GetProductSellerID(ctx context.Context, productID int) (int, error) {
	var sellerID int
	err := r.db.QueryRow(ctx, `SELECT seller_id FROM products WHERE id=$1`, productID).Scan(&sellerID) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	return sellerID, err
}

func (r *AdminRepo) ListSellers(ctx context.Context) ([]models.User, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT u.id, u.name, u.phone, COALESCE(u.business_name,''), COALESCE(u.email,''),
		       u.rating, u.role, u.is_active,
		       COALESCE(u.tier::text,'normal') AS tier,
		       COALESCE(u.total_value_sold,0) AS total_value_sold,
		       COALESCE(ps.net_profit, 0) AS net_profit,
		       u.created_at, u.updated_at
		FROM users u
		LEFT JOIN seller_profit_summary ps ON ps.seller_id = u.id
		WHERE u.role='seller'
		ORDER BY u.created_at DESC
		LIMIT 500`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.BusinessName, &u.Email,
			&u.Rating, &u.Role, &u.IsActive, &u.Tier, &u.TotalValueSold,
			&u.NetProfit, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *AdminRepo) SetSellerActive(ctx context.Context, sellerID int, active bool, reason string) error {
	var tag interface{ RowsAffected() int64 }
	var err error
	if active {
		// Reactivating — clear the deactivation reason
		tag, err = r.db.Exec(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
			`UPDATE users SET is_active=$1, deactivation_reason=NULL, updated_at=NOW() WHERE id=$2 AND role='seller'`,
			active, sellerID)
	} else {
		// Deactivating — store the reason
		tag, err = r.db.Exec(ctx, //nolint:gosec // G201/G202: parameterized query, no user input in SQL structure
			`UPDATE users SET is_active=$1, deactivation_reason=$2, updated_at=NOW() WHERE id=$3 AND role='seller'`,
			active, reason, sellerID)
	}
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("seller with id %d not found", sellerID)
	}
	return nil
}

func (r *AdminRepo) ListAllTransactions(ctx context.Context) ([]models.Transaction, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT id, product_id, COALESCE(seller_id,0), buyer_name, buyer_phone, amount, broker_fee, seller_gets, status, COALESCE(notes,''), created_at, updated_at
		FROM transactions
		ORDER BY created_at DESC
		LIMIT 500`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var txns []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.ProductID, &t.SellerID, &t.BuyerName, &t.BuyerPhone,
			&t.Amount, &t.BrokerFee, &t.SellerGets, &t.Status, &t.Notes, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, rows.Err()
}

func (r *AdminRepo) AddLocation(ctx context.Context, l models.Location) (int, error) {
	var id int
	err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.QueryRow(ctx, `
		INSERT INTO locations (region, zone, subcity, woreda)
		VALUES ($1, $2, $3, $4) RETURNING id`,
			l.Region, l.Zone, l.Subcity, l.Woreda).Scan(&id)
	return id, err
}

func (r *AdminRepo) UpdateLocation(ctx context.Context, l models.Location) error {
	_, err := r.db.Exec(ctx, `
		UPDATE locations SET region=$1, zone=$2, subcity=$3, woreda=$4 WHERE id=$5`,
		l.Region, l.Zone, l.Subcity, l.Woreda, l.ID)
	return err
}

func (r *AdminRepo) DeleteLocation(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM locations WHERE id=$1`, id) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	return err
}

func (r *AdminRepo) ListLocations(ctx context.Context) ([]models.Location, error) {
	rows, err := //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		r.db.Query(ctx, `
		SELECT id, region, zone, subcity, woreda FROM locations ORDER BY region, zone, subcity, woreda`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var locs []models.Location
	for rows.Next() {
		var l models.Location
		if err := rows.Scan(&l.ID, &l.Region, &l.Zone, &l.Subcity, &l.Woreda); err != nil {
			return nil, err
		}
		locs = append(locs, l)
	}
	return locs, rows.Err()
}

func scanProducts(rows interface {
	Next() bool
	Scan(...interface{}) error
	Err() error
}) ([]models.Product, error) {
	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.SellerID, &p.LocationID, &p.Title, &p.Price,
			&p.Description, &p.CategoryID, &p.Subcategory, &p.Item, &p.Condition,
			&p.Status, &p.AdminNotes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// ── Analytics ─────────────────────────────────────────────────────────────────

// AnalyticsOverview returns aggregated platform stats for the admin dashboard.
