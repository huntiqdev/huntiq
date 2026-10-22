package models

import "time"

type Product struct {
	ID          int        `json:"id"`
	SellerID    int        `json:"seller_id"`
	LocationID  *int       `json:"location_id"`
	Title       string     `json:"title"`
	Price               float64            `json:"price"`
	Description         string             `json:"description"`
	CategoryID          *int               `json:"category_id"`
	CategoryHierarchyID *int               `json:"category_hierarchy_id,omitempty"`
	Subcategory         string             `json:"subcategory,omitempty"`
	Item                string             `json:"item,omitempty"`
	Condition           string             `json:"condition"`
	Status              string             `json:"status"`
	AdminNotes  string     `json:"admin_notes,omitempty"`
	ReviewedBy  *int       `json:"reviewed_by,omitempty"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`

	// Precise location (optional)
	Kebele       string   `json:"kebele,omitempty"`
	SpecificArea string   `json:"specific_area,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`

	// Listing fee fields
	ListingDurationDays    *int     `json:"listing_duration_days,omitempty"`
	StartupFeePercent      *float64 `json:"startup_fee_percent,omitempty"`
	WeeklyIncrementPercent *float64 `json:"weekly_increment_percent,omitempty"`
	TotalFeePercent        *float64 `json:"total_fee_percent,omitempty"`
	TotalFeeAmount         *float64 `json:"total_fee_amount,omitempty"`
	FeePaid                bool     `json:"fee_paid"`
	FeePaymentPending      bool     `json:"fee_payment_pending"` // seller submitted claim, awaiting admin verification

	// Joined fields
	Images     []ProductImage     `json:"images,omitempty"`
	Category   *Category          `json:"category,omitempty"`
	Location   *Location          `json:"location,omitempty"`
	SellerTier string             `json:"seller_tier,omitempty"` // normal | gold | premium
	Attributes []ProductAttribute `json:"attributes,omitempty"`

	// Tier access control
	RequiredBuyerTier string `json:"required_buyer_tier,omitempty"`
	SellerActive      bool   `json:"seller_active"`

	// Sold tracking
	SoldAt *time.Time `json:"sold_at,omitempty"`

	// ── Rental fields (only for listing_type='rent') ──
	ListingType         string   `json:"listing_type"`                    // 'sale' or 'rent'
	RentalPrice         *float64 `json:"rental_price,omitempty"`          // price per duration
	RentalDurationValue *int     `json:"rental_duration_value,omitempty"` // e.g., 1, 7, 30
	RentalDurationUnit  *string  `json:"rental_duration_unit,omitempty"`  // hour, day, week, month, year
	DepositAmount       *float64 `json:"deposit_amount,omitempty"`        // security deposit
	AvailableFrom       *string  `json:"available_from,omitempty"`        // DATE
	AvailableUntil      *string  `json:"available_until,omitempty"`       // DATE
	IsCurrentlyRented   bool     `json:"is_currently_rented,omitempty"`
	CurrentTenantID     *int     `json:"current_tenant_id,omitempty"`
	RentedUntil         *string  `json:"rented_until,omitempty"` // DATE
}

type ProductImage struct {
	ID          int        `json:"id"`
	ProductID   int        `json:"product_id"`
	ImageURL    string     `json:"image_url"`
	OrderNumber int        `json:"order_number"`
	IsThumbnail bool       `json:"is_thumbnail"`
	IsBlurred   bool       `json:"is_blurred"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
