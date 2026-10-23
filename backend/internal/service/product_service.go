package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/your-org/used-tools-marketplace/pkg/constants"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/internal/repository"
)

// isAllowedImageURL accepts /uploads/ paths from our own backend only.
// When backendBase is configured, only absolute URLs from that base are accepted —
// relative paths are rejected to prevent mixed URL storage in the DB.
func (s *ProductService) isAllowedImageURL(u string) bool {
	// If backendBase is set, require absolute URLs from our own origin only.
	if s.backendBase != "" {
		if strings.HasPrefix(u, s.backendBase+"/uploads/") {
			return true
		}
		// Allow picsum.photos for seeded demo products only
		return strings.HasPrefix(u, "https://picsum.photos/")
	}
	// backendBase not configured (local dev) — only accept relative /uploads/ paths.
	// Never allow absolute URLs from unknown domains even in dev.
	return strings.HasPrefix(u, "/uploads/")
}

type ProductService struct {
	repo        *repository.ProductRepo
	catRepo     *repository.CategoryRepo
	locRepo     *repository.LocationRepo
	userRepo    *repository.UserRepo
	attrRepo    *repository.AttributeRepo
	backendBase string // base URL of this server, used to validate /uploads/ image URLs
}

func NewProductService(repo repository.ProductRepository, catRepo repository.CategoryRepository, locRepo repository.LocationRepository, userRepo repository.UserRepository) *ProductService {
	return &ProductService{repo: repo, catRepo: catRepo, locRepo: locRepo, userRepo: userRepo}
}

// WithAttributeRepo sets the attribute repository for dynamic product attributes.
func (s *ProductService) WithAttributeRepo(attrRepo *repository.AttributeRepo) *ProductService {
	s.attrRepo = attrRepo
	return s
}

// WithBackendBase sets the backend base URL used to validate absolute /uploads/ image URLs.
// Call this at startup after config is loaded.
func (s *ProductService) WithBackendBase(base string) *ProductService {
	s.backendBase = base
	return s
}

// loadFeeSettings reads fee parameters from platform_settings.
// Falls back to defaults if any key is missing, unparseable, or out of bounds.
// Bounds: fee percentages must be in (0, 100], free_weeks must be >= 0.
func (s *ProductService) loadFeeSettings(ctx context.Context) FeeSettings {
	fs := DefaultFeeSettings()
	get := func(key string) (float64, bool) {
		v, err := s.userRepo.GetSetting(ctx, key)
		if err != nil || v == "" {
			return 0, false
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil || f < 0 || f > 100 {
			return 0, false // reject out-of-bounds fee percentages; 0 is valid (fee-free tier)
		}
		return f, true
	}
	getInt := func(key string) (int, bool) {
		v, err := s.userRepo.GetSetting(ctx, key)
		if err != nil || v == "" {
			return 0, false
		}
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return 0, false // reject negative values
		}
		return n, true
	}
	if v, ok := get("normal_startup_fee"); ok {
		fs.NormalStartupFee = v
	}
	if v, ok := get("gold_startup_fee"); ok {
		fs.GoldStartupFee = v
	}
	if v, ok := get("premium_startup_fee"); ok {
		fs.PremiumStartupFee = v
	}
	if v, ok := get("normal_weekly_increment"); ok {
		fs.NormalWeeklyIncrement = v
	}
	if v, ok := get("gold_weekly_increment"); ok {
		fs.GoldWeeklyIncrement = v
	}
	if v, ok := get("premium_weekly_increment"); ok {
		fs.PremiumWeeklyIncrement = v
	}
	if v, ok := getInt("listing_free_weeks"); ok {
		fs.FreeWeeks = v
	}
	return fs
}

func (s *ProductService) ListApproved(ctx context.Context, f dto.ProductFilter) ([]models.Product, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}
	return s.repo.ListApproved(ctx, f)
}

func (s *ProductService) GetByID(ctx context.Context, id int) (*models.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product id")
	}
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Deep Scan: Block access if the seller is banned/inactive.
	// This ensures consistency even if someone has a direct link to a product.
	if !p.SellerActive {
		return nil, errors.New("product no longer available (seller account inactive)")
	}
	// Also ensure public only sees approved/sold products.
	// Pending/Rejected are for seller/admin eyes only.
	if p.Status != "approved" && p.Status != "sold" {
		return nil, errors.New("product is not yet available for public view")
	}
	if s.attrRepo != nil {
		attrs, err := s.attrRepo.GetProductAttributes(ctx, p.ID)
		if err == nil && len(attrs) > 0 {
			p.Attributes = attrs
		}
	}
	return p, nil
}

func (s *ProductService) ListCategories(ctx context.Context) ([]models.Category, error) {
	return s.catRepo.List(ctx)
}

// PreviewListingFee returns the fee breakdown for a given seller, price, and duration.
func (s *ProductService) PreviewListingFee(ctx context.Context, sellerID int, price float64, durationDays int) (FeeBreakdown, error) {
	tierStr, _, err := s.userRepo.GetTier(ctx, sellerID)
	if err != nil {
		tierStr = string(TierNormal)
	}
	feeSettings := s.loadFeeSettings(ctx)
	return CalculateListingFeeWithSettings(SellerTier(tierStr), price, durationDays, feeSettings), nil
}

// prohibitedKeywords defines illicit and dangerous goods blocked across EN, AM, OR, and TI.
var prohibitedKeywords = []struct {
	term     string
	category string
}{
	// Weapons & Explosives (EN)
	{"firearm", "weapons"}, {"handgun", "weapons"}, {"pistol", "weapons"},
	{"rifle", "weapons"}, {"ammunition", "weapons"}, {"explosive", "explosives"},
	{"grenade", "explosives"}, {"bomb", "explosives"},
	// Weapons & Explosives (AM, OR, TI)
	{"ሽጉጥ", "weapons"}, {"ጠመንጃ", "weapons"}, {"ጥይት", "weapons"}, {"ቦምብ", "explosives"},
	{"qawwee", "weapons"}, {"rasaasa", "weapons"}, {"bombaa", "explosives"},
	{"ቦምባ", "explosives"},
	// Narcotics & Controlled Substances
	{"cocaine", "controlled_substances"}, {"heroin", "controlled_substances"},
	{"methamphetamine", "controlled_substances"}, {"opium", "controlled_substances"},
	{"ኮኬይን", "controlled_substances"}, {"ሄሮይን", "controlled_substances"},
	{"kookeeyinii", "controlled_substances"}, {"heeroyinii", "controlled_substances"},
	// Counterfeits & Financial Fraud
	{"counterfeit money", "fraud"}, {"fake currency", "fraud"}, {"fake notes", "fraud"},
	{"የሐሰት ገንዘብ", "fraud"}, {"qarshii sobaa", "fraud"}, {"ናይ ሓሶት ቅርሺ", "fraud"},
}

func checkProhibitedContent(title, description string) error {
	combined := strings.ToLower(title + " " + description)
	for _, pk := range prohibitedKeywords {
		if strings.Contains(combined, strings.ToLower(pk.term)) {
			return fmt.Errorf("listing contains prohibited item (%s: '%s') — policy violation", pk.category, pk.term)
		}
	}
	return nil
}

func (s *ProductService) CreateForSeller(ctx context.Context, sellerID int, req dto.CreateProductRequest) (*models.Product, error) {
	// Load all content limits from DB — admin can change these at any time
	minImages := s.getSettingInt(ctx, constants.SettingMinImagesPerListing, constants.DefaultMinImages)
	maxTitle := s.getSettingInt(ctx, constants.SettingMaxTitleLength, constants.DefaultMaxTitleLength)
	maxDesc := s.getSettingInt(ctx, constants.SettingMaxDescriptionLength, constants.DefaultMaxDescriptionLength)

	if len(req.ImageURLs) < minImages {
		return nil, fmt.Errorf("minimum %d images required", minImages)
	}
	if req.Title == "" || req.Price <= 0 {
		return nil, errors.New("title and a positive price are required")
	}
	if len(req.Title) > maxTitle {
		return nil, fmt.Errorf("title must be %d characters or fewer", maxTitle)
	}

	// Trust & Safety: Block prohibited illegal goods and dangerous items
	if err := checkProhibitedContent(req.Title, req.Description); err != nil {
		return nil, err
	}

	// Anti-Spam: Check for duplicates within the last 24h
	isDup, err := s.repo.CheckDuplicate(ctx, sellerID, req.Title, req.Price)
	if err == nil && isDup {
		return nil, errors.New("you already posted a similar product recently — please wait or edit the existing listing")
	}

	if len(req.Description) > maxDesc {
		return nil, fmt.Errorf("description must be %d characters or fewer", maxDesc)
	}
	
	// ── Listing type validation ──
	listingType := req.ListingType
	if listingType == "" {
		listingType = "sale" // default
	}
	if listingType != "sale" && listingType != "rent" {
		return nil, errors.New("listing_type must be 'sale' or 'rent'")
	}
	
	// ── Rental-specific validation ──
	if listingType == "rent" {
		if req.RentalPrice == nil || *req.RentalPrice <= 0 {
			return nil, errors.New("rental_price is required for rent listings")
		}
		if req.RentalDurationValue == nil || *req.RentalDurationValue <= 0 {
			return nil, errors.New("rental_duration_value is required for rent listings")
		}
		if req.RentalDurationUnit == nil || *req.RentalDurationUnit == "" {
			return nil, errors.New("rental_duration_unit is required for rent listings")
		}
		validUnits := map[string]bool{"hour": true, "day": true, "week": true, "month": true, "year": true}
		if !validUnits[*req.RentalDurationUnit] {
			return nil, errors.New("rental_duration_unit must be: hour, day, week, month, or year")
		}
		if req.AvailableFrom == nil || *req.AvailableFrom == "" {
			return nil, errors.New("available_from date is required for rent listings")
		}
	}
	
	// Validate image URLs — must be /uploads/ paths or trusted CDN URLs
	for _, u := range req.ImageURLs {
		if !s.isAllowedImageURL(u) {
			return nil, errors.New("invalid image URL: only uploaded or CDN images are allowed")
		}
	}
	if req.Region == "" && req.LocationID == nil {
		return nil, errors.New("location (region) is required")
	}
	// Resolve location_id from hierarchy if not provided directly
	if req.LocationID == nil && req.Region != "" {
		id, err := s.locRepo.ResolveID(ctx, req.Region, req.Zone, req.Subcity, req.Woreda)
		if err != nil {
			return nil, errors.New("invalid location: " + err.Error())
		}
		req.LocationID = &id
	}
	// Resolve seller tier and calculate listing fee from DB settings
	tierStr, _, err := s.userRepo.GetTier(ctx, sellerID)
	if err != nil {
		tierStr = string(TierNormal)
	}

	// Enforce: seller cannot gate their product to a buyer tier higher than their own.
	// A normal seller creating a premium-only product creates unreachable inventory.
	if req.RequiredBuyerTier != "" && req.RequiredBuyerTier != "normal" {
		tierRank := map[string]int{"normal": 0, "gold": 1, "premium": 2}
		if tierRank[req.RequiredBuyerTier] > tierRank[tierStr] {
			return nil, fmt.Errorf("your %s seller tier cannot restrict products to %s buyers — you must first upgrade your own tier", tierStr, req.RequiredBuyerTier)
		}
	}

	// Enforce per-tier image limits from DB settings
	maxImages := s.maxImagesForTier(ctx, SellerTier(tierStr))
	if len(req.ImageURLs) > maxImages {
		return nil, errors.New("too many images for your tier")
	}
	duration := req.ListingDurationDays
	if duration < 1 {
		duration = 30 // default 1 month
	}
	// Validate duration is one of the allowed values.
	// If not a match, snap to the nearest valid option (permissive — don't reject).
	validDurations := constants.ListingDurations
	isValidDuration := false
	for _, d := range validDurations {
		if duration == d {
			isValidDuration = true
			break
		}
	}
	if !isValidDuration {
		// Snap to nearest valid duration rather than rejecting
		nearest := validDurations[0]
		for _, d := range validDurations {
			if abs(d-duration) < abs(nearest-duration) {
				nearest = d
			}
		}
		duration = nearest
	}
	feeSettings := s.loadFeeSettings(ctx)
	fee := CalculateListingFeeWithSettings(SellerTier(tierStr), req.Price, duration, feeSettings)
	cond := req.Condition
	// Map frontend condition values to DB enum (new | like_new | good | fair | poor)
	switch cond {
	case "new":
		cond = "new"
	case "like_new":
		cond = "like_new"
	case "fair":
		cond = "fair"
	case "poor":
		cond = "poor"
	default:
		cond = "good" // "used", "" or anything else → good
	}
	catHierarchyID := req.CategoryHierarchyID
	if catHierarchyID == nil {
		catHierarchyID = req.CategoryID
	}
	p := &models.Product{
		SellerID:               sellerID,
		Title:                  req.Title,
		Price:                  req.Price,
		Description:            req.Description,
		CategoryID:             req.CategoryID,
		CategoryHierarchyID:    catHierarchyID,
		Subcategory:            req.Subcategory,
		Item:                   req.Item,
		Condition:              cond,
		LocationID:             req.LocationID,
		LocationHierarchyID:    req.LocationHierarchyID,
		Kebele:                 req.Kebele,
		SpecificArea:           req.SpecificArea,
		Latitude:               req.Latitude,
		Longitude:              req.Longitude,
		ListingDurationDays:    &duration,
		StartupFeePercent:      &fee.StartupPct,
		WeeklyIncrementPercent: &fee.WeeklyIncrementPct,
		TotalFeePercent:        &fee.TotalFeePercent,
		TotalFeeAmount:         &fee.TotalFeeAmount,
		ListingType:            listingType,
		RentalPrice:            req.RentalPrice,
		RentalDurationValue:    req.RentalDurationValue,
		RentalDurationUnit:     req.RentalDurationUnit,
		DepositAmount:          req.DepositAmount,
		AvailableFrom:          req.AvailableFrom,
		AvailableUntil:         req.AvailableUntil,
	}
	id, err := s.repo.Create(ctx, p, req.ImageURLs)
	if err != nil {
		return nil, err
	}
	p.ID = id
	p.Status = "pending"

	if len(req.Attributes) > 0 && s.attrRepo != nil {
		for key, val := range req.Attributes {
			def, err := s.attrRepo.GetAttributeBySlug(ctx, key)
			if err != nil || def == nil {
				continue
			}
			pa := &models.ProductAttribute{
				ProductID:   id,
				AttributeID: def.ID,
			}
			switch v := val.(type) {
			case string:
				pa.ValueText = &v
			case float64:
				pa.ValueNumber = &v
			case bool:
				pa.ValueBool = &v
			case map[string]interface{}:
				pa.ValueJSON = v
			case []interface{}:
				pa.ValueJSON = map[string]interface{}{"list": v}
			default:
				strVal := fmt.Sprintf("%v", v)
				pa.ValueText = &strVal
			}
			_ = s.attrRepo.SetProductAttribute(ctx, pa)
		}
	}
	return p, nil
}

func (s *ProductService) UpdateForSeller(ctx context.Context, productID, sellerID int, req dto.UpdateProductRequest) error {
	existing, err := s.repo.GetByIDAndSeller(ctx, productID, sellerID)
	if err != nil {
		return errors.New("product not found")
	}
	if existing.Status != "pending" && existing.Status != "rejected" {
		return errors.New("only pending or rejected products can be edited")
	}
	if len(req.ImageURLs) > 0 {
		minImages := s.getSettingInt(ctx, constants.SettingMinImagesPerListing, constants.DefaultMinImages)
		if len(req.ImageURLs) < minImages {
			return fmt.Errorf("minimum %d images required", minImages)
		}
		for _, u := range req.ImageURLs {
			if !s.isAllowedImageURL(u) {
				return errors.New("invalid image URL: only uploaded or CDN images are allowed")
			}
		}
		tierStr, _, terr := s.userRepo.GetTier(ctx, sellerID)
		if terr != nil {
			tierStr = string(TierNormal)
		}
		if len(req.ImageURLs) > s.maxImagesForTier(ctx, SellerTier(tierStr)) {
			return errors.New("too many images for your tier")
		}
	}

	maxTitle := s.getSettingInt(ctx, constants.SettingMaxTitleLength, constants.DefaultMaxTitleLength)
	maxDesc := s.getSettingInt(ctx, constants.SettingMaxDescriptionLength, constants.DefaultMaxDescriptionLength)

	if req.Title != nil {
		if *req.Title == "" {
			return errors.New("title cannot be empty")
		}
		if len(*req.Title) > maxTitle {
			return fmt.Errorf("title must be %d characters or fewer", maxTitle)
		}
		existing.Title = *req.Title
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return errors.New("price must be greater than zero")
		}
		existing.Price = *req.Price
	}
	if req.Description != nil {
		if len(*req.Description) > maxDesc {
			return fmt.Errorf("description must be %d characters or fewer", maxDesc)
		}
		existing.Description = *req.Description
	}
	if req.CategoryID != nil {
		existing.CategoryID = req.CategoryID
	}
	if req.Condition != nil {
		cond := *req.Condition
		switch cond {
		case "new", "like_new", "fair", "poor":
			// valid DB enum values
		default:
			cond = "good" // "used" or anything else → good
		}
		existing.Condition = cond
	}
	if req.Subcategory != nil {
		existing.Subcategory = *req.Subcategory
	}
	if req.Item != nil {
		existing.Item = *req.Item
	}
	if req.LocationID != nil {
		existing.LocationID = req.LocationID
	}
	if req.LocationHierarchyID != nil {
		existing.LocationHierarchyID = req.LocationHierarchyID
	}
	if req.Kebele != nil {
		existing.Kebele = *req.Kebele
	}
	if req.SpecificArea != nil {
		existing.SpecificArea = *req.SpecificArea
	}
	if req.Latitude != nil {
		existing.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		existing.Longitude = req.Longitude
	}
	// Rental fields — only update if listing_type is 'rent'
	if existing.ListingType == "rent" {
		if req.RentalPrice != nil {
			existing.RentalPrice = req.RentalPrice
		}
		if req.RentalDurationValue != nil {
			existing.RentalDurationValue = req.RentalDurationValue
		}
		if req.RentalDurationUnit != nil {
			existing.RentalDurationUnit = req.RentalDurationUnit
		}
		if req.DepositAmount != nil {
			existing.DepositAmount = req.DepositAmount
		}
		if req.AvailableFrom != nil {
			existing.AvailableFrom = req.AvailableFrom
		}
		if req.AvailableUntil != nil {
			existing.AvailableUntil = req.AvailableUntil
		}
	}

	if req.CategoryHierarchyID != nil {
		existing.CategoryHierarchyID = req.CategoryHierarchyID
	}

	if err := s.repo.Update(ctx, existing, req.ImageURLs); err != nil {
		return err
	}

	if len(req.Attributes) > 0 && s.attrRepo != nil {
		for key, val := range req.Attributes {
			def, err := s.attrRepo.GetAttributeBySlug(ctx, key)
			if err != nil || def == nil {
				continue
			}
			pa := &models.ProductAttribute{
				ProductID:   productID,
				AttributeID: def.ID,
			}
			switch v := val.(type) {
			case string:
				pa.ValueText = &v
			case float64:
				pa.ValueNumber = &v
			case bool:
				pa.ValueBool = &v
			case map[string]interface{}:
				pa.ValueJSON = v
			case []interface{}:
				pa.ValueJSON = map[string]interface{}{"list": v}
			default:
				strVal := fmt.Sprintf("%v", v)
				pa.ValueText = &strVal
			}
			_ = s.attrRepo.SetProductAttribute(ctx, pa)
		}
	}
	return nil
}

func (s *ProductService) ListForSeller(ctx context.Context, sellerID int) ([]models.Product, error) {
	return s.repo.ListBySeller(ctx, sellerID)
}


// getSettingInt reads an integer setting from DB, returning defaultVal if missing/invalid.
func (s *ProductService) getSettingInt(ctx context.Context, key string, defaultVal int) int {
	v, _ := s.userRepo.GetSetting(ctx, key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}

// maxImagesForTier returns the max images allowed for a tier from DB settings.
func (s *ProductService) maxImagesForTier(ctx context.Context, tier SellerTier) int {
	key := constants.SettingNormalMaxImages
	defaultVal := constants.DefaultNormalMax
	switch tier {
	case TierGold:
		key = constants.SettingGoldMaxImages
		defaultVal = constants.DefaultGoldMax
	case TierPremium:
		key = constants.SettingPremiumMaxImages
		defaultVal = constants.DefaultPremiumMax
	}
	return s.getSettingInt(ctx, key, defaultVal)
}

// FeeOptions returns all duration options with fee breakdown for a given price+tier, using DB settings.
func (s *ProductService) FeeOptions(ctx context.Context, tier SellerTier, price float64) []FeeBreakdown {
	feeSettings := s.loadFeeSettings(ctx)
	var out []FeeBreakdown
	for _, d := range ListingDurations {
		out = append(out, CalculateListingFeeWithSettings(tier, price, d, feeSettings))
	}
	return out
}

// PopularCategories returns categories ranked by activity in the last N days.
func (s *ProductService) PopularCategories(ctx context.Context, days, limit int) ([]models.PopularCategory, error) {
	return s.repo.PopularCategories(ctx, days, limit)
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
