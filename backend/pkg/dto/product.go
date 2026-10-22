package dto

// ProductFilter holds query params for the product list endpoint.
type ProductFilter struct {
	CategoryID          *int
	CategoryHierarchyID *int  // New category hierarchy ID
	IncludeDescendants  bool  // Include products from child categories
	Subcategory         *string
	Item                *string
	MinPrice            *float64
	MaxPrice            *float64
	LocationID          *int
	LocationHierarchyID *int  // Data-driven location hierarchy node ID
	Region              *string
	Zone                *string
	Subcity             *string
	Woreda              *string
	Condition           *string
	Search              *string
	Page                int
	Limit               int
	ListingType         *string // 'sale' or 'rent' filter
}

// CreateProductRequest is used by sellers to submit a new product.
type CreateProductRequest struct {
	Title               string                 `json:"title"`
	Price               float64                `json:"price"`
	Description         string                 `json:"description"`
	CategoryID          *int                   `json:"category_id"`
	CategoryHierarchyID *int                   `json:"category_hierarchy_id,omitempty"`
	Subcategory         string                 `json:"subcategory"`
	Item                string                 `json:"item"`
	Condition           string                 `json:"condition"`
	LocationID          *int                   `json:"location_id"`
	LocationHierarchyID *int                   `json:"location_hierarchy_id,omitempty"`
	Region              string                 `json:"region"`
	Zone                string                 `json:"zone,omitempty"`
	Subcity             string                 `json:"subcity,omitempty"`
	Woreda              string                 `json:"woreda,omitempty"`
	Kebele              string                 `json:"kebele,omitempty"`
	SpecificArea        string                 `json:"specific_area,omitempty"`
	Latitude            *float64               `json:"latitude,omitempty"`
	Longitude           *float64               `json:"longitude,omitempty"`
	ImageURLs           []string               `json:"image_urls"`
	ListingDurationDays int                    `json:"listing_duration_days"`
	RequiredBuyerTier   string                 `json:"required_buyer_tier,omitempty"`
	Attributes          map[string]interface{} `json:"attributes,omitempty"`
	
	// Listing type: 'sale' (default) or 'rent'
	ListingType string `json:"listing_type"`
	
	// Rental-specific fields (required if listing_type='rent')
	RentalPrice         *float64 `json:"rental_price,omitempty"`
	RentalDurationValue *int     `json:"rental_duration_value,omitempty"`
	RentalDurationUnit  *string  `json:"rental_duration_unit,omitempty"` // hour, day, week, month, year
	DepositAmount       *float64 `json:"deposit_amount,omitempty"`
	AvailableFrom       *string  `json:"available_from,omitempty"` // DATE
	AvailableUntil      *string  `json:"available_until,omitempty"` // DATE
}

// UpdateProductRequest allows partial updates by the seller.
type UpdateProductRequest struct {
	Title               *string                `json:"title"`
	Price               *float64               `json:"price"`
	Description         *string                `json:"description"`
	CategoryID          *int                   `json:"category_id"`
	CategoryHierarchyID *int                   `json:"category_hierarchy_id,omitempty"`
	Subcategory         *string                `json:"subcategory"`
	Item                *string                `json:"item"`
	Attributes          map[string]interface{} `json:"attributes,omitempty"`
	Condition           *string                `json:"condition"`
	LocationID          *int                   `json:"location_id"`
	LocationHierarchyID *int                   `json:"location_hierarchy_id,omitempty"`
	Kebele              *string                `json:"kebele,omitempty"`
	SpecificArea        *string                `json:"specific_area,omitempty"`
	Latitude            *float64               `json:"latitude,omitempty"`
	Longitude           *float64               `json:"longitude,omitempty"`
	ImageURLs           []string               `json:"image_urls"`
	
	// Rental fields (only for listing_type='rent')
	RentalPrice         *float64 `json:"rental_price,omitempty"`
	RentalDurationValue *int     `json:"rental_duration_value,omitempty"`
	RentalDurationUnit  *string  `json:"rental_duration_unit,omitempty"`
	DepositAmount       *float64 `json:"deposit_amount,omitempty"`
	AvailableFrom       *string  `json:"available_from,omitempty"`
	AvailableUntil      *string  `json:"available_until,omitempty"`
}



