package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
)

type TelebirrRepo struct {
	db *pgxpool.Pool
}

func NewTelebirrRepo(db *pgxpool.Pool) *TelebirrRepo {
	return &TelebirrRepo{db: db}
}

// CreateOrder persists a new Telebirr order.
func (r *TelebirrRepo) CreateOrder(ctx context.Context, o *models.TelebirrOrder) error {
	query := `
		INSERT INTO telebirr_orders (
			out_trade_no, user_id, amount, currency, title, order_type, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		o.OutTradeNo, o.UserID, o.Amount, o.Currency, o.Title, o.OrderType, o.Status,
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

// GetOrderByOutTradeNo retrieves a Telebirr order by its unique merchant reference.
func (r *TelebirrRepo) GetOrderByOutTradeNo(ctx context.Context, outTradeNo string) (*models.TelebirrOrder, error) {
	var o models.TelebirrOrder
	var rawBytes []byte

	query := `
		SELECT id, out_trade_no, user_id, amount, currency, title, order_type, status,
		       prepay_id, transaction_no, payer_phone, raw_callback, verified_at, created_at, updated_at
		FROM telebirr_orders
		WHERE out_trade_no = $1
	`
	err := r.db.QueryRow(ctx, query, outTradeNo).Scan(
		&o.ID,
		&o.OutTradeNo,
		&o.UserID,
		&o.Amount,
		&o.Currency,
		&o.Title,
		&o.OrderType,
		&o.Status,
		&o.PrepayID,
		&o.TransactionNo,
		&o.PayerPhone,
		&rawBytes,
		&o.VerifiedAt,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if len(rawBytes) > 0 {
		_ = json.Unmarshal(rawBytes, &o.RawCallback)
	}

	return &o, nil
}

// UpdateOrderStatus updates order fulfillment status upon webhook notification.
func (r *TelebirrRepo) UpdateOrderStatus(ctx context.Context, outTradeNo, status string, txNo, payerPhone *string, rawCallback map[string]interface{}) error {
	var rawBytes []byte
	if rawCallback != nil {
		rawBytes, _ = json.Marshal(rawCallback)
	}

	query := `
		UPDATE telebirr_orders
		SET status = $1, transaction_no = COALESCE($2, transaction_no),
		    payer_phone = COALESCE($3, payer_phone), raw_callback = $4,
		    verified_at = CASE WHEN $1 = 'completed' THEN NOW() ELSE verified_at END,
		    updated_at = NOW()
		WHERE out_trade_no = $5
	`
	_, err := r.db.Exec(ctx, query, status, txNo, payerPhone, rawBytes, outTradeNo)
	return err
}

// GetOrdersByUser retrieves paginated orders for a user.
func (r *TelebirrRepo) GetOrdersByUser(ctx context.Context, userID int, limit, offset int) ([]models.TelebirrOrder, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM telebirr_orders WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, out_trade_no, user_id, amount, currency, title, order_type, status,
		       prepay_id, transaction_no, payer_phone, verified_at, created_at, updated_at
		FROM telebirr_orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.TelebirrOrder
	for rows.Next() {
		var o models.TelebirrOrder
		if err := rows.Scan(
			&o.ID,
			&o.OutTradeNo,
			&o.UserID,
			&o.Amount,
			&o.Currency,
			&o.Title,
			&o.OrderType,
			&o.Status,
			&o.PrepayID,
			&o.TransactionNo,
			&o.PayerPhone,
			&o.VerifiedAt,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		list = append(list, o)
	}

	return list, total, nil
}
