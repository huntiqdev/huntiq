package fpa

import (
	"math"
)

// RevenueBreakdown details all independent monetization streams and escrow separation
type RevenueBreakdown struct {
	TotalActiveSellers     float64
	CompletedTransactions  int
	GrossMerchandiseVolume float64 // Seller funds (GMV) — NOT HUNTIQ platform revenue
	CommissionRevenue      float64 // HUNTIQ commission (e.g. 2%)
	ListingFeeRevenue      float64
	SubscriptionRevenue    float64 // Gold & Premium tiers
	PromotedListingRevenue float64
	BusinessToolsRevenue   float64
	GrossBilledRevenue     float64
	DiscountsAndRefunds    float64
	NetRecognizedRevenue   float64
}

// CalculateRevenue models all independent revenue channels for month m
func CalculateRevenue(mau float64, cfg RevenueConfig) RevenueBreakdown {
	if mau <= 0 {
		return RevenueBreakdown{}
	}

	sellersPct := cfg.SellersPctOfUsers
	if sellersPct <= 0 {
		sellersPct = 0.15 // 15% of active users list items
	}
	totalSellers := mau * sellersPct

	// 1. Marketplace Commission & GMV (Seller Escrow Separation)
	txPerSeller := cfg.TransactionsPerSellerMo
	if txPerSeller <= 0 {
		txPerSeller = 3.2
	}
	completedTxFloat := totalSellers * txPerSeller
	completedTx := int(math.Round(completedTxFloat))

	atv := cfg.AverageTransactionValue
	if atv <= 0 {
		atv = 4500.0 // 4,500 ETB average transaction value
	}
	gmv := float64(completedTx) * atv

	commissionRate := cfg.CommissionRatePct
	if commissionRate <= 0 {
		commissionRate = 0.02 // 2.0% broker fee
	}
	grossCommission := gmv * commissionRate

	refundRate := cfg.RefundCancellationRate
	if refundRate < 0 {
		refundRate = 0.03
	}
	refunds := grossCommission * refundRate

	discountRate := cfg.CommissionDiscountRate
	if discountRate < 0 {
		discountRate = 0.05
	}
	discounts := grossCommission * discountRate

	netCommission := grossCommission - refunds - discounts
	if netCommission < 0 {
		netCommission = 0
	}

	// 2. Listing Fees
	listingsPerSeller := cfg.ListingsPerSellerMo
	if listingsPerSeller <= 0 {
		listingsPerSeller = 4.0
	}
	totalListings := totalSellers * listingsPerSeller

	paidListingAdoption := cfg.PaidListingAdoptionPct
	if paidListingAdoption <= 0 {
		paidListingAdoption = 0.20 // 20% beyond free quota
	}
	listingFeePerItem := cfg.PaidListingFeePerItem
	if listingFeePerItem <= 0 {
		listingFeePerItem = 50.0 // 50 ETB
	}
	listingFeeRevenue := totalListings * paidListingAdoption * listingFeePerItem

	// 3. Premium Membership Tiers (Gold & Premium)
	goldConversion := cfg.GoldTierConversionPct
	if goldConversion <= 0 {
		goldConversion = 0.02 // 2%
	}
	goldPrice := cfg.GoldTierMonthlyPrice
	if goldPrice <= 0 {
		goldPrice = 500.0 // 500 ETB/mo
	}
	goldRevenue := totalSellers * goldConversion * goldPrice

	premiumConversion := cfg.PremiumTierConversionPct
	if premiumConversion <= 0 {
		premiumConversion = 0.008 // 0.8%
	}
	premiumPrice := cfg.PremiumTierMonthlyPrice
	if premiumPrice <= 0 {
		premiumPrice = 1500.0 // 1,500 ETB/mo
	}
	premiumRevenue := totalSellers * premiumConversion * premiumPrice
	totalSubscriptionRevenue := goldRevenue + premiumRevenue

	// 4. Paid Promotion (Featured & Sponsored listings)
	featuredAdoption := cfg.FeaturedListingMonthlyAdoption
	if featuredAdoption <= 0 {
		featuredAdoption = 0.03
	}
	featuredPrice := cfg.FeaturedListingPrice
	if featuredPrice <= 0 {
		featuredPrice = 250.0 // 250 ETB per featured listing
	}
	promotedRevenue := totalListings * featuredAdoption * featuredPrice

	// 5. Business Tools & Verified Seller Subscriptions
	verifiedSellersPct := cfg.VerifiedBusinessSellersPct
	if verifiedSellersPct <= 0 {
		verifiedSellersPct = 0.015
	}
	verifiedToolPrice := cfg.VerifiedBusinessToolFeeMo
	if verifiedToolPrice <= 0 {
		verifiedToolPrice = 800.0 // 800 ETB/mo
	}
	businessToolsRevenue := totalSellers * verifiedSellersPct * verifiedToolPrice

	grossBilled := grossCommission + listingFeeRevenue + totalSubscriptionRevenue + promotedRevenue + businessToolsRevenue
	netRecognized := netCommission + listingFeeRevenue + totalSubscriptionRevenue + promotedRevenue + businessToolsRevenue

	return RevenueBreakdown{
		TotalActiveSellers:     round2(totalSellers),
		CompletedTransactions:  completedTx,
		GrossMerchandiseVolume: round2(gmv), // Strictly seller funds
		CommissionRevenue:      round2(netCommission),
		ListingFeeRevenue:      round2(listingFeeRevenue),
		SubscriptionRevenue:    round2(totalSubscriptionRevenue),
		PromotedListingRevenue: round2(promotedRevenue),
		BusinessToolsRevenue:   round2(businessToolsRevenue),
		GrossBilledRevenue:     round2(grossBilled),
		DiscountsAndRefunds:    round2(refunds + discounts),
		NetRecognizedRevenue:   round2(netRecognized),
	}
}
