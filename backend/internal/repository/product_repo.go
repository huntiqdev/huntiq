package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/pkg/dto"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/pkg/search"
)

type ProductRepo struct{ db *pgxpool.Pool }

func NewProductRepo(db *pgxpool.Pool) *ProductRepo { return &ProductRepo{db} }

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullStrPtr(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}

func addWhere(b *strings.Builder, args []interface{}, col string, val interface{}) ([]interface{}, int) {
	idx := len(args) + 1
	b.WriteString(col)
	b.WriteString(fmt.Sprintf("$%d", idx)) //nolint:gosec // G201: Building static SQL fragment with computed $N parameter placeholder
	return append(args, val), idx
}

const sqlProductSelectCols = `
	SELECT p.id, p.seller_id, p.location_id, p.title, p.price,
	       COALESCE(p.description,''), p.category_id,
	       COALESCE(p.subcategory,''), COALESCE(p.item,''), p.condition,
	       p.status, p.created_at, p.updated_at,
	       p.expires_at,
	       c.id, c.name, c.slug,
	       l.id, l.region, l.zone, l.subcity, l.woreda,
	       COALESCE(u.tier::text,'normal') AS seller_tier,
	       COALESCE(p.access_tier,'normal') AS required_buyer_tier,
	       COALESCE(p.listing_type::text,'sale') AS listing_type,
	       p.rental_price, p.rental_duration_value, p.rental_duration_unit,
	       p.deposit_amount, p.available_from, p.available_until,
	       COALESCE(p.is_currently_rented,FALSE) AS is_currently_rented,
	       p.category_hierarchy_id,
	       p.location_hierarchy_id
	FROM products p
	LEFT JOIN categories c ON c.id = p.category_id
	LEFT JOIN locations  l ON l.id = p.location_id
	LEFT JOIN users      u ON u.id = p.seller_id`

const sqlGetProductByID = `
	SELECT p.id, p.seller_id, p.location_id, p.title, p.price,
	        COALESCE(p.description,''), p.category_id,
	        COALESCE(p.subcategory,''), COALESCE(p.item,''), p.condition,
	        p.status, COALESCE(p.admin_notes,''), p.created_at, p.updated_at,
	        COALESCE(p.kebele,''), COALESCE(p.specific_area,''), p.latitude, p.longitude,
	        p.expires_at,
	        c.id, c.name, c.slug,
	        l.id, l.region, l.zone, l.subcity, l.woreda,
	        COALESCE(u.tier::text,'normal') AS seller_tier,
	        COALESCE(p.access_tier,'normal') AS required_buyer_tier,
	        COALESCE(u.is_active, true) AS seller_active,
	        COALESCE(p.listing_type::text,'sale') AS listing_type,
	        p.rental_price, p.rental_duration_value, p.rental_duration_unit,
	        p.deposit_amount, p.available_from, p.available_until,
	        COALESCE(p.is_currently_rented,FALSE) AS is_currently_rented,
	        p.current_tenant_id, p.rented_until,
	        p.category_hierarchy_id,
	        p.location_hierarchy_id
	 FROM products p
	 LEFT JOIN categories c ON c.id = p.category_id
	 LEFT JOIN locations  l ON l.id = p.location_id
	 LEFT JOIN users      u ON u.id = p.seller_id
	 WHERE p.id = $1`

const sqlCheckDuplicate = `
	SELECT EXISTS(
		SELECT 1 FROM products
		WHERE seller_id = $1
		  AND title = $2
		  AND price = $3
		  AND created_at > NOW() - INTERVAL '24 hours'
		  AND status != 'rejected'
	)`

const sqlCreateProduct = `
	INSERT INTO products (
		seller_id, location_id, title, price, description, category_id, subcategory, item, condition, status,
		listing_duration_days, startup_fee_percent, weekly_increment_percent, total_fee_percent, total_fee_amount,
		kebele, specific_area, latitude, longitude,
		listing_type, rental_price, rental_duration_value, rental_duration_unit, deposit_amount, available_from, available_until,
		category_hierarchy_id, location_hierarchy_id
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, 'pending',
		$10, $11, $12, $13, $14,
		$15, $16, $17, $18,
		$19, $20, $21, $22, $23, $24, $25,
		$26, $27
	) RETURNING id`

const sqlCreateProductImage = `
	INSERT INTO product_images (product_id, image_url, order_number) VALUES ($1, $2, $3)`

const sqlUpdateProduct = `
	UPDATE products
	 SET title=$1, price=$2, description=$3, category_id=COALESCE($4::int, category_id), subcategory=$5, item=$6,
	     condition=$7, location_id=COALESCE($8::int, location_id), kebele=$9, specific_area=$10,
	     latitude=$11, longitude=$12,
	     rental_price=$13, rental_duration_value=$14, rental_duration_unit=$15,
	     deposit_amount=$16, available_from=$17, available_until=$18,
	     category_hierarchy_id=COALESCE($19::int, category_hierarchy_id),
	     location_hierarchy_id=COALESCE($20::int, location_hierarchy_id),
	     startup_fee_percent=COALESCE($21, startup_fee_percent),
	     weekly_increment_percent=COALESCE($22, weekly_increment_percent),
	     total_fee_percent=COALESCE($23, total_fee_percent),
	     total_fee_amount=COALESCE($24, total_fee_amount),
	     fee_paid=$25,
	     status='pending', admin_notes='', updated_at=NOW()
	 WHERE id=$26 AND seller_id=$27 AND status IN ('pending','rejected')`

const sqlDeleteProductImages = `DELETE FROM product_images WHERE product_id=$1`

const sqlListBySeller = `
	SELECT p.id, p.seller_id, p.location_id, p.title, p.price,
	        COALESCE(p.description,''), p.category_id,
	        COALESCE(p.subcategory,''), COALESCE(p.item,''), p.condition,
	        p.status, COALESCE(p.admin_notes,''), p.created_at, p.updated_at,
	        COALESCE(p.listing_type::text,'sale') AS listing_type,
	        p.rental_price, p.rental_duration_value, p.rental_duration_unit,
	        p.deposit_amount, p.available_from, p.available_until,
	        COALESCE(p.is_currently_rented,FALSE) AS is_currently_rented,
	        p.listing_duration_days, p.total_fee_percent, p.total_fee_amount,
	        COALESCE(p.fee_paid,FALSE) AS fee_paid,
	        COALESCE(p.fee_payment_pending,FALSE) AS fee_payment_pending,
	        p.expires_at
	 FROM products p WHERE p.seller_id = $1 ORDER BY p.created_at DESC`

const sqlListImagesBatch = `
	SELECT DISTINCT ON (product_id) id, product_id, image_url, order_number, created_at
	 FROM product_images WHERE product_id = ANY($1::int[])
	 ORDER BY product_id, order_number ASC`

const sqlGetProductByIDAndSeller = `
	SELECT id, seller_id, location_id, title, price,
	        COALESCE(description,''), category_id,
	        COALESCE(subcategory,''), COALESCE(item,''), condition,
	        status, COALESCE(admin_notes,''), created_at, updated_at,
	        COALESCE(listing_type::text,'sale') AS listing_type,
	        rental_price, rental_duration_value, rental_duration_unit,
	        deposit_amount, available_from, available_until,
	        COALESCE(is_currently_rented,FALSE) AS is_currently_rented,
	        listing_duration_days, total_fee_percent, total_fee_amount,
	        COALESCE(fee_paid,FALSE), COALESCE(fee_payment_pending,FALSE)
	 FROM products WHERE id=$1 AND seller_id=$2`

const sqlGetProductImages = `
	SELECT id, product_id, image_url, order_number, created_at
	 FROM product_images WHERE product_id = $1 ORDER BY order_number`

const sqlPopularCategories = `
	WITH listing_counts AS (
		SELECT category_id, COUNT(*) AS new_listings
		FROM products
		WHERE created_at >= NOW() - make_interval(days => $1)
		  AND category_id IS NOT NULL
		GROUP BY category_id
	),
	order_counts AS (
		SELECT p.category_id, COUNT(o.id) AS order_count
		FROM order_requests o
		JOIN products p ON p.id = o.product_id
		WHERE o.created_at >= NOW() - make_interval(days => $1)
		  AND p.category_id IS NOT NULL
		GROUP BY p.category_id
	)
	SELECT
		c.id, c.name, c.slug,
		COALESCE(l.new_listings, 0)  AS new_listings,
		COALESCE(o.order_count, 0)   AS order_count,
		COALESCE(l.new_listings, 0) + COALESCE(o.order_count, 0) * 2 AS activity_score
	FROM categories c
	LEFT JOIN listing_counts l ON l.category_id = c.id
	LEFT JOIN order_counts   o ON o.category_id = c.id
	WHERE COALESCE(c.is_active, true) = true
	  AND (COALESCE(l.new_listings, 0) + COALESCE(o.order_count, 0)) > 0
	ORDER BY activity_score DESC
	LIMIT $2`

func (r *ProductRepo) buildApprovedWhere(ctx context.Context, f dto.ProductFilter) (string, []interface{}) {
	args := []interface{}{"approved"}
	var wb strings.Builder
	wb.WriteString("p.status = $1 AND (p.expires_at IS NULL OR p.expires_at > NOW()) AND u.is_active = true")

	// Support both legacy category_id and new category_hierarchy_id
	if f.CategoryHierarchyID != nil {
		if f.IncludeDescendants {
			// Filter by category and all descendants using database function, with fallback to legacy category_id
			idx := len(args) + 1
			wb.WriteString(fmt.Sprintf(" AND (p.category_hierarchy_id IN (SELECT id FROM get_category_descendants($%d)) OR (p.category_hierarchy_id IS NULL AND p.category_id = $%d))", idx, idx))
			args = append(args, *f.CategoryHierarchyID)
		} else {
			idx := len(args) + 1
			wb.WriteString(fmt.Sprintf(" AND (p.category_hierarchy_id = $%d OR (p.category_hierarchy_id IS NULL AND p.category_id = $%d))", idx, idx))
			args = append(args, *f.CategoryHierarchyID)
		}
	} else if f.CategoryID != nil {
		// Legacy category_id filter (flat categories), with compatibility for new category_hierarchy_id
		idx := len(args) + 1
		wb.WriteString(fmt.Sprintf(" AND (p.category_id = $%d OR p.category_hierarchy_id = $%d)", idx, idx))
		args = append(args, *f.CategoryID)
	}
	
	if f.MinPrice != nil {
		args, _ = addWhere(&wb, args, " AND p.price >= ", *f.MinPrice)
	}
	if f.MaxPrice != nil {
		args, _ = addWhere(&wb, args, " AND p.price <= ", *f.MaxPrice)
	}
	if f.LocationHierarchyID != nil {
		idx := len(args) + 1
		wb.WriteString(fmt.Sprintf(` AND p.location_hierarchy_id IN (SELECT id FROM get_location_descendants($%d))`, idx))
		args = append(args, *f.LocationHierarchyID)
	} else if f.LocationID != nil {
		args, _ = addWhere(&wb, args, " AND p.location_id = ", *f.LocationID)
	}
	if f.Woreda != nil && *f.Woreda != "" {
		args, _ = addWhere(&wb, args, " AND l.woreda ILIKE ", likeWrap(*f.Woreda))
	} else if f.Subcity != nil && *f.Subcity != "" {
		idx := len(args) + 1
		wb.WriteString(fmt.Sprintf(` AND (l.subcity ILIKE $%d OR p.location_hierarchy_id IN (SELECT id FROM location_hierarchy WHERE name_en ILIKE $%d))`, idx, idx))
		args = append(args, likeWrap(*f.Subcity))
	} else if f.Zone != nil && *f.Zone != "" {
		args, _ = addWhere(&wb, args, " AND l.zone ILIKE ", likeWrap(*f.Zone))
	} else if f.Region != nil && *f.Region != "" {
		idx := len(args) + 1
		wb.WriteString(fmt.Sprintf(` AND (
			l.region ILIKE $%d
			OR p.location_hierarchy_id IN (
				SELECT d.id FROM location_hierarchy r
				CROSS JOIN LATERAL get_location_descendants(r.id) d
				WHERE r.name_en ILIKE $%d
			)
		)`, idx, idx))
		args = append(args, likeWrap(*f.Region))
	}
	
	// Enhanced search with brand detection, conversational cleaning, and multilingual synonym expansion
	if f.Search != nil && *f.Search != "" {
		searchTerm := *f.Search
		cleanTerm := search.CleanConversationalSearchTerm(searchTerm)
		if cleanTerm == "" {
			cleanTerm = strings.TrimSpace(searchTerm)
		}

		synTsQuery := search.BuildSynonymTsQuery(cleanTerm)

		// Check for brand mentions using clean term
		brandQuery := `SELECT brand_name FROM extract_brands($` + fmt.Sprintf("%d", len(args)+1) + `) LIMIT 1`
		var detectedBrand *string
		_ = r.db.QueryRow(ctx, brandQuery, cleanTerm).Scan(&detectedBrand)

		idxTerm := len(args) + 1
		args = append(args, cleanTerm)

		if synTsQuery != "" {
			idxSyn := len(args) + 1
			args = append(args, synTsQuery)

			wb.WriteString(fmt.Sprintf(` AND (
				search_vector @@ websearch_to_tsquery('simple', $%d)
				OR search_vector @@ to_tsquery('simple', $%d)
				OR p.title ILIKE '%%' || $%d || '%%'
				OR EXISTS (
					SELECT 1 FROM product_attributes pa
					WHERE pa.product_id = p.id
					  AND pa.value_text ILIKE '%%' || $%d || '%%'
				)
				OR p.category_hierarchy_id IN (SELECT id FROM get_category_descendants_by_term($%d))
				OR (p.category_hierarchy_id IS NULL AND p.category_id IN (SELECT id FROM get_category_descendants_by_term($%d)))
			)`, idxTerm, idxSyn, idxTerm, idxTerm, idxTerm, idxTerm))
		} else {
			wb.WriteString(fmt.Sprintf(` AND (
				search_vector @@ websearch_to_tsquery('simple', $%d)
				OR p.title ILIKE '%%' || $%d || '%%'
				OR EXISTS (
					SELECT 1 FROM product_attributes pa
					WHERE pa.product_id = p.id
					  AND pa.value_text ILIKE '%%' || $%d || '%%'
				)
				OR p.category_hierarchy_id IN (SELECT id FROM get_category_descendants_by_term($%d))
				OR (p.category_hierarchy_id IS NULL AND p.category_id IN (SELECT id FROM get_category_descendants_by_term($%d)))
			)`, idxTerm, idxTerm, idxTerm, idxTerm, idxTerm))
		}
	}
	
	if f.Condition != nil && *f.Condition != "" {
		if *f.Condition == "used" {
			args, _ = addWhere(&wb, args, " AND p.condition != ", "new")
		} else {
			args, _ = addWhere(&wb, args, " AND p.condition = ", *f.Condition)
		}
	}
	if f.Subcategory != nil && *f.Subcategory != "" {
		args, _ = addWhere(&wb, args, " AND p.subcategory = ", *f.Subcategory)
	}
	if f.Item != nil && *f.Item != "" {
		args, _ = addWhere(&wb, args, " AND p.item = ", *f.Item)
	}
	if f.ListingType != nil && *f.ListingType != "" {
		args, _ = addWhere(&wb, args, " AND p.listing_type = ", *f.ListingType)
	}
	if f.SellerTier != nil && *f.SellerTier != "" {
		st := strings.ToLower(strings.TrimSpace(*f.SellerTier))
		if st == "gold,premium" || st == "premium,gold" {
			wb.WriteString(" AND u.tier IN ('gold', 'premium')")
		} else if st == "gold" || st == "premium" || st == "normal" {
			args, _ = addWhere(&wb, args, " AND u.tier = ", st)
		}
	}

	return wb.String(), args
}

// CountApproved returns the exact total count of matching approved products in the database.
func (r *ProductRepo) CountApproved(ctx context.Context, f dto.ProductFilter) (int, error) {
	whereClause, args := r.buildApprovedWhere(ctx, f)
	countQuery := `
		SELECT COUNT(*)
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
		LEFT JOIN locations  l ON l.id = p.location_id
		LEFT JOIN users      u ON u.id = p.seller_id
		WHERE ` + whereClause
	var count int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&count)
	return count, err
}

func (r *ProductRepo) ListApproved(ctx context.Context, f dto.ProductFilter) ([]models.Product, error) {
	whereClause, args := r.buildApprovedWhere(ctx, f)

	idx := len(args) + 1
	offset := (f.Page - 1) * f.Limit
	args = append(args, f.Limit, offset)

	orderClause := "ORDER BY p.created_at DESC"
	if f.SortBy != nil {
		switch *f.SortBy {
		case "price_asc":
			orderClause = "ORDER BY p.price ASC, p.created_at DESC"
		case "price_desc":
			orderClause = "ORDER BY p.price DESC, p.created_at DESC"
		case "newest":
			orderClause = "ORDER BY p.created_at DESC"
		}
	}

	query := sqlProductSelectCols + "\n\tWHERE " + whereClause +
		"\n\t" + orderClause +
		"\n\tLIMIT $" + fmt.Sprintf("%d", idx) + " OFFSET $" + fmt.Sprintf("%d", idx+1) //nolint:gosec // G201: Building static SQL with computed $N placeholders
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var cat models.Category
		var catID *int
		var catName, catSlug *string
		var loc models.Location
		var locID *int
		var locRegion *string
		var sellerTier, requiredBuyerTier, listingType string
		var rentalDurationUnit *string
		var availableFrom, availableUntil *string

		if err := rows.Scan(
			&p.ID, &p.SellerID, &p.LocationID, &p.Title, &p.Price,
			&p.Description, &p.CategoryID, &p.Subcategory, &p.Item, &p.Condition,
			&p.Status, &p.CreatedAt, &p.UpdatedAt,
			&p.ExpiresAt,
			&catID, &catName, &catSlug,
			&locID, &locRegion, &loc.Zone, &loc.Subcity, &loc.Woreda,
			&sellerTier, &requiredBuyerTier,
			&listingType, &p.RentalPrice, &p.RentalDurationValue, &rentalDurationUnit,
			&p.DepositAmount, &availableFrom, &availableUntil, &p.IsCurrentlyRented,
			&p.CategoryHierarchyID,
			&p.LocationHierarchyID,
		); err != nil {
			return nil, err
		}
		p.SellerTier = sellerTier
		p.RequiredBuyerTier = requiredBuyerTier
		p.ListingType = listingType
		p.RentalDurationUnit = rentalDurationUnit
		p.AvailableFrom = availableFrom
		p.AvailableUntil = availableUntil
		if catID != nil && catName != nil && catSlug != nil {
			cat.ID = *catID
			cat.Name = *catName
			cat.Slug = *catSlug
			p.Category = &cat
		}
		if locID != nil {
			loc.ID = *locID
			if locRegion != nil {
				loc.Region = *locRegion
			}
			p.Location = &loc
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return products, nil
	}

	idIndex := make(map[int]int, len(products))
	productIDs := make([]int, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
		idIndex[p.ID] = i
	}
	imgRows, err := r.db.Query(ctx, sqlListImagesBatch, productIDs)
	if err != nil {
		return products, nil
	}
	defer imgRows.Close()
	for imgRows.Next() {
		var img models.ProductImage
		if err := imgRows.Scan(&img.ID, &img.ProductID, &img.ImageURL, &img.OrderNumber, &img.CreatedAt); err == nil {
			if i, ok := idIndex[img.ProductID]; ok {
				products[i].Images = []models.ProductImage{img}
			}
		}
	}
	return products, nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id int) (*models.Product, error) {
	var p models.Product
	var cat models.Category
	var catID *int
	var catName, catSlug *string
	var loc models.Location
	var locID *int
	var locRegion *string
	var sellerTier, requiredBuyerTier, listingType string
	var rentalDurationUnit *string
	var availableFrom, availableUntil *string

	err := r.db.QueryRow(ctx, sqlGetProductByID, id).Scan(
		&p.ID, &p.SellerID, &p.LocationID, &p.Title, &p.Price,
		&p.Description, &p.CategoryID, &p.Subcategory, &p.Item, &p.Condition,
		&p.Status, &p.AdminNotes, &p.CreatedAt, &p.UpdatedAt,
		&p.Kebele, &p.SpecificArea, &p.Latitude, &p.Longitude,
		&p.ExpiresAt,
		&catID, &catName, &catSlug,
		&locID, &locRegion, &loc.Zone, &loc.Subcity, &loc.Woreda,
		&sellerTier, &requiredBuyerTier, &p.SellerActive,
		&listingType, &p.RentalPrice, &p.RentalDurationValue, &rentalDurationUnit,
		&p.DepositAmount, &availableFrom, &availableUntil, &p.IsCurrentlyRented,
		&p.CurrentTenantID, &p.RentedUntil,
		&p.CategoryHierarchyID,
		&p.LocationHierarchyID,
	)
	if err != nil {
		return nil, err
	}
	p.SellerTier = sellerTier
	p.RequiredBuyerTier = requiredBuyerTier
	p.ListingType = listingType
	p.RentalDurationUnit = rentalDurationUnit
	p.AvailableFrom = availableFrom
	p.AvailableUntil = availableUntil
	if catID != nil {
		cat.ID = *catID
		cat.Name = *catName
		cat.Slug = *catSlug
		p.Category = &cat
	}
	if locID != nil {
		loc.ID = *locID
		if locRegion != nil {
			loc.Region = *locRegion
		}
		p.Location = &loc
	}

	imgRows, err := r.db.Query(ctx, sqlGetProductImages, id)
	if err != nil {
		return nil, err
	}
	defer imgRows.Close()
	for imgRows.Next() {
		var img models.ProductImage
		if err := imgRows.Scan(&img.ID, &img.ProductID, &img.ImageURL, &img.OrderNumber, &img.CreatedAt); err != nil {
			return nil, err
		}
		p.Images = append(p.Images, img)
	}
	return &p, imgRows.Err()
}

func (r *ProductRepo) CheckDuplicate(ctx context.Context, sellerID int, title string, price float64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, sqlCheckDuplicate, sellerID, title, price).Scan(&exists)
	return exists, err
}

func (r *ProductRepo) Create(ctx context.Context, p *models.Product, imageURLs []string) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	catHierarchyID := p.CategoryHierarchyID
	if catHierarchyID == nil {
		catHierarchyID = p.CategoryID
	}
	var id int
	err = tx.QueryRow(ctx, sqlCreateProduct,
		p.SellerID, p.LocationID, p.Title, p.Price, p.Description,
		p.CategoryID, p.Subcategory, p.Item, p.Condition,
		p.ListingDurationDays, p.StartupFeePercent, p.WeeklyIncrementPercent, p.TotalFeePercent, p.TotalFeeAmount,
		nullStr(p.Kebele), nullStr(p.SpecificArea), p.Latitude, p.Longitude,
		p.ListingType, p.RentalPrice, p.RentalDurationValue, nullStrPtr(p.RentalDurationUnit), p.DepositAmount, nullStrPtr(p.AvailableFrom), nullStrPtr(p.AvailableUntil),
		catHierarchyID, p.LocationHierarchyID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	for i, url := range imageURLs {
		if _, err = tx.Exec(ctx, sqlCreateProductImage, id, url, i+1); err != nil {
			return 0, err
		}
	}
	return id, tx.Commit(ctx)
}

func (r *ProductRepo) Update(ctx context.Context, p *models.Product, imageURLs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, sqlUpdateProduct,
		p.Title, p.Price, p.Description, p.CategoryID, p.Subcategory, p.Item,
		p.Condition, p.LocationID, nullStr(p.Kebele), nullStr(p.SpecificArea),
		p.Latitude, p.Longitude,
		p.RentalPrice, p.RentalDurationValue, nullStrPtr(p.RentalDurationUnit),
		p.DepositAmount, nullStrPtr(p.AvailableFrom), nullStrPtr(p.AvailableUntil),
		p.CategoryHierarchyID,
		p.LocationHierarchyID,
		p.StartupFeePercent, p.WeeklyIncrementPercent, p.TotalFeePercent, p.TotalFeeAmount, p.FeePaid,
		p.ID, p.SellerID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product not found or cannot be edited in its current state")
	}
	if len(imageURLs) > 0 {
		if _, err = tx.Exec(ctx, sqlDeleteProductImages, p.ID); err != nil {
			return err
		}
		for i, url := range imageURLs {
			if _, err = tx.Exec(ctx, sqlCreateProductImage, p.ID, url, i+1); err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}

func (r *ProductRepo) ListBySeller(ctx context.Context, sellerID int) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, sqlListBySeller, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var p models.Product
		var listingType string
		var rentalDurationUnit *string
		var availableFrom, availableUntil *string
		if err := rows.Scan(
			&p.ID, &p.SellerID, &p.LocationID, &p.Title, &p.Price,
			&p.Description, &p.CategoryID, &p.Subcategory, &p.Item, &p.Condition,
			&p.Status, &p.AdminNotes, &p.CreatedAt, &p.UpdatedAt,
			&listingType, &p.RentalPrice, &p.RentalDurationValue, &rentalDurationUnit,
			&p.DepositAmount, &availableFrom, &availableUntil, &p.IsCurrentlyRented,
			&p.ListingDurationDays, &p.TotalFeePercent, &p.TotalFeeAmount,
			&p.FeePaid, &p.FeePaymentPending, &p.ExpiresAt,
		); err != nil {
			return nil, err
		}
		p.ListingType = listingType
		p.RentalDurationUnit = rentalDurationUnit
		p.AvailableFrom = availableFrom
		p.AvailableUntil = availableUntil
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return products, nil
	}
	idIndex := make(map[int]int, len(products))
	productIDs := make([]int, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
		idIndex[p.ID] = i
	}
	imgRows, err := r.db.Query(ctx, sqlListImagesBatch, productIDs)
	if err == nil {
		defer imgRows.Close()
		for imgRows.Next() {
			var img models.ProductImage
			if err := imgRows.Scan(&img.ID, &img.ProductID, &img.ImageURL, &img.OrderNumber, &img.CreatedAt); err == nil {
				if i, ok := idIndex[img.ProductID]; ok {
					products[i].Images = []models.ProductImage{img}
				}
			}
		}
	}
	return products, nil
}

func (r *ProductRepo) GetByIDAndSeller(ctx context.Context, id, sellerID int) (*models.Product, error) {
	var p models.Product
	var listingType string
	var rentalDurationUnit *string
	var availableFrom, availableUntil *string
	err := r.db.QueryRow(ctx, sqlGetProductByIDAndSeller, id, sellerID).Scan(
		&p.ID, &p.SellerID, &p.LocationID, &p.Title, &p.Price,
		&p.Description, &p.CategoryID, &p.Subcategory, &p.Item, &p.Condition,
		&p.Status, &p.AdminNotes, &p.CreatedAt, &p.UpdatedAt,
		&listingType, &p.RentalPrice, &p.RentalDurationValue, &rentalDurationUnit,
		&p.DepositAmount, &availableFrom, &availableUntil, &p.IsCurrentlyRented,
		&p.ListingDurationDays, &p.TotalFeePercent, &p.TotalFeeAmount,
		&p.FeePaid, &p.FeePaymentPending,
	)
	if err != nil {
		return nil, err
	}
	p.ListingType = listingType
	p.RentalDurationUnit = rentalDurationUnit
	p.AvailableFrom = availableFrom
	p.AvailableUntil = availableUntil
	imgRows, err := r.db.Query(ctx, sqlGetProductImages, id)
	if err == nil {
		defer imgRows.Close()
		for imgRows.Next() {
			var img models.ProductImage
			if err := imgRows.Scan(&img.ID, &img.ProductID, &img.ImageURL, &img.OrderNumber, &img.CreatedAt); err == nil {
				p.Images = append(p.Images, img)
			}
		}
	}
	return &p, nil
}

func (r *ProductRepo) PopularCategories(ctx context.Context, days, limit int) ([]models.PopularCategory, error) {
	rows, err := r.db.Query(ctx, sqlPopularCategories, days, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.PopularCategory
	for rows.Next() {
		var p models.PopularCategory
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.NewListings, &p.OrderCount, &p.ActivityScore); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
