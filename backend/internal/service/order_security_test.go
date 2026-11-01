package service

import (
	"context"
	"testing"

	"github.com/your-org/used-tools-marketplace/pkg/dto"
)

func TestOrderRequest_ValidationAndSecurity(t *testing.T) {
	svc := &SellerService{}

	ctx := context.Background()

	t.Run("missing buyer name", func(t *testing.T) {
		req := dto.CreateOrderRequest{
			ProductID:  1,
			BuyerName:  "",
			BuyerPhone: "+251911223344",
		}
		_, err := svc.CreateOrderRequest(ctx, req)
		if err == nil || err.Error() != "buyer_name and buyer_phone are required" {
			t.Fatalf("expected buyer_name required error, got %v", err)
		}
	})

	t.Run("invalid Ethiopian phone number", func(t *testing.T) {
		req := dto.CreateOrderRequest{
			ProductID:  1,
			BuyerName:  "Abebe",
			BuyerPhone: "invalid-phone",
		}
		_, err := svc.CreateOrderRequest(ctx, req)
		if err == nil || err.Error() != "invalid Ethiopian phone number for buyer_phone" {
			t.Fatalf("expected invalid Ethiopian phone error, got %v", err)
		}
	})

	t.Run("negative offer price", func(t *testing.T) {
		negPrice := -500.0
		req := dto.CreateOrderRequest{
			ProductID:  1,
			BuyerName:  "Abebe",
			BuyerPhone: "0911223344",
			OfferPrice: &negPrice,
		}
		_, err := svc.CreateOrderRequest(ctx, req)
		if err == nil || err.Error() != "offer_price must be greater than zero" {
			t.Fatalf("expected offer_price error, got %v", err)
		}
	})

	t.Run("zero offer price", func(t *testing.T) {
		zeroPrice := 0.0
		req := dto.CreateOrderRequest{
			ProductID:  1,
			BuyerName:  "Abebe",
			BuyerPhone: "+251911223344",
			OfferPrice: &zeroPrice,
		}
		_, err := svc.CreateOrderRequest(ctx, req)
		if err == nil || err.Error() != "offer_price must be greater than zero" {
			t.Fatalf("expected offer_price error, got %v", err)
		}
	})
}

func TestProductService_GetByIDAndSeller_Guards(t *testing.T) {
	svc := &ProductService{}
	ctx := context.Background()

	t.Run("invalid product ID", func(t *testing.T) {
		_, err := svc.GetByIDAndSeller(ctx, 0, 1)
		if err == nil {
			t.Fatal("expected error for product_id <= 0")
		}
	})

	t.Run("invalid seller ID", func(t *testing.T) {
		_, err := svc.GetByIDAndSeller(ctx, 10, 0)
		if err == nil {
			t.Fatal("expected error for seller_id <= 0")
		}
	})
}

func TestTransactionService_Validation(t *testing.T) {
	svc := &TransactionService{}
	ctx := context.Background()

	t.Run("invalid buyer phone", func(t *testing.T) {
		req := dto.CreateTransactionRequest{
			ProductID:  1,
			BuyerName:  "Dawit",
			BuyerPhone: "not-a-phone",
			Amount:     1000,
		}
		_, err := svc.Create(ctx, 1, req)
		if err == nil || err.Error() != "invalid Ethiopian phone number for buyer_phone" {
			t.Fatalf("expected invalid Ethiopian phone error, got %v", err)
		}
	})

	t.Run("negative amount", func(t *testing.T) {
		req := dto.CreateTransactionRequest{
			ProductID:  1,
			BuyerName:  "Dawit",
			BuyerPhone: "+251911223344",
			Amount:     -100,
		}
		_, err := svc.Create(ctx, 1, req)
		if err == nil || err.Error() != "amount must be greater than zero" {
			t.Fatalf("expected amount error, got %v", err)
		}
	})

	t.Run("zero amount", func(t *testing.T) {
		req := dto.CreateTransactionRequest{
			ProductID:  1,
			BuyerName:  "Dawit",
			BuyerPhone: "+251911223344",
			Amount:     0,
		}
		_, err := svc.Create(ctx, 1, req)
		if err == nil || err.Error() != "amount must be greater than zero" {
			t.Fatalf("expected amount error, got %v", err)
		}
	})

	t.Run("missing product id", func(t *testing.T) {
		req := dto.CreateTransactionRequest{
			ProductID:  0,
			BuyerName:  "Dawit",
			BuyerPhone: "+251911223344",
			Amount:     500,
		}
		_, err := svc.Create(ctx, 1, req)
		if err == nil || err.Error() != "product_id is required" {
			t.Fatalf("expected product_id required error, got %v", err)
		}
	})
}

func TestSellerService_UpdateOrderStatus_StatusValidation(t *testing.T) {
	svc := &SellerService{}
	ctx := context.Background()

	t.Run("invalid order status", func(t *testing.T) {
		req := dto.UpdateOrderRequest{
			Status: "arbitrary_status",
		}
		err := svc.UpdateOrderStatus(ctx, 1, 1, req)
		if err == nil || err.Error() != "invalid status" {
			t.Fatalf("expected 'invalid status' error, got %v", err)
		}
	})
}
