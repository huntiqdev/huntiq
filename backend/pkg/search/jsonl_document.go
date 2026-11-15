package search

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// CategoryFacet represents the taxonomy lineage in the denormalized JSONL document.
type CategoryFacet struct {
	ID        int      `json:"id"`
	Slug      string   `json:"slug"`
	NameEn    string   `json:"name_en"`
	NameAm    string   `json:"name_am,omitempty"`
	NameOm    string   `json:"name_om,omitempty"`
	NameTi    string   `json:"name_ti,omitempty"`
	Level     int      `json:"level"`
	Path      string   `json:"path"`
	Ancestors []string `json:"ancestors,omitempty"`
}

// LocationFacet represents the geographic placement of the listing.
type LocationFacet struct {
	ID           int       `json:"id"`
	Region       string    `json:"region"`
	Zone         string    `json:"zone,omitempty"`
	Subcity      string    `json:"subcity,omitempty"`
	Woreda       string    `json:"woreda,omitempty"`
	SpecificArea string    `json:"specific_area,omitempty"`
	Coordinates  *GeoCoord `json:"coordinates,omitempty"`
}

// GeoCoord contains latitude and longitude for spatial proximity matching.
type GeoCoord struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// ProductSearchDoc is the modern denormalized JSONL representation of a marketplace listing.
// In JSONL format, each line is an independent, valid JSON object containing full product
// facets, precomputed search tokens, localized terms, and typed technical attributes.
type ProductSearchDoc struct {
	DocID          string                 `json:"doc_id"`
	ProductID      int                    `json:"product_id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description,omitempty"`
	Price          float64                `json:"price"`
	Currency       string                 `json:"currency"`
	Condition      string                 `json:"condition"`
	ListingType    string                 `json:"listing_type"`
	Status         string                 `json:"status"`
	SellerID       int                    `json:"seller_id"`
	SellerTier     string                 `json:"seller_tier"`
	Category       CategoryFacet          `json:"category"`
	Location       LocationFacet          `json:"location"`
	Attributes     map[string]interface{} `json:"attributes"`
	SearchKeywords []string               `json:"search_keywords,omitempty"`
	ImageURLs      []string               `json:"image_urls,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// NewProductSearchDoc constructs a validated, normalized ProductSearchDoc.
func NewProductSearchDoc(
	productID int,
	title string,
	description string,
	price float64,
	condition string,
	listingType string,
	status string,
	sellerID int,
	sellerTier string,
	cat CategoryFacet,
	loc LocationFacet,
	attrs map[string]interface{},
	imageURLs []string,
	createdAt, updatedAt time.Time,
) *ProductSearchDoc {
	if attrs == nil {
		attrs = make(map[string]interface{})
	}
	if sellerTier == "" {
		sellerTier = "normal"
	}
	if listingType == "" {
		listingType = "sale"
	}

	doc := &ProductSearchDoc{
		DocID:       fmt.Sprintf("prod_%d", productID),
		ProductID:   productID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Price:       price,
		Currency:    "ETB",
		Condition:   condition,
		ListingType: listingType,
		Status:      status,
		SellerID:    sellerID,
		SellerTier:  sellerTier,
		Category:    cat,
		Location:    loc,
		Attributes:  attrs,
		ImageURLs:   imageURLs,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	doc.SearchKeywords = doc.GenerateSearchKeywords()
	return doc
}

// GenerateSearchKeywords extracts tokenized search keywords from title, category names, and attributes.
func (d *ProductSearchDoc) GenerateSearchKeywords() []string {
	seen := make(map[string]bool)
	var keywords []string

	addWord := func(word string) {
		clean := strings.ToLower(strings.TrimSpace(word))
		clean = strings.Trim(clean, ".,!?;:'\"()[]{}")
		if len(clean) >= 2 && !seen[clean] {
			seen[clean] = true
			keywords = append(keywords, clean)
		}
	}

	// Tokenize Title
	for _, w := range strings.Fields(d.Title) {
		addWord(w)
	}

	// Add Category terms
	if d.Category.NameEn != "" {
		for _, w := range strings.Fields(d.Category.NameEn) {
			addWord(w)
		}
	}
	if d.Category.NameAm != "" {
		for _, w := range strings.Fields(d.Category.NameAm) {
			addWord(w)
		}
	}
	if d.Category.Slug != "" {
		for _, w := range strings.Split(d.Category.Slug, "-") {
			addWord(w)
		}
	}

	// Add Attribute values (strings, brands, models)
	for _, v := range d.Attributes {
		if s, ok := v.(string); ok && s != "" {
			for _, w := range strings.Fields(s) {
				addWord(w)
			}
		}
	}

	return keywords
}

// ToJSONL serializes the search document to a single newline-delimited JSON line.
func (d *ProductSearchDoc) ToJSONL() ([]byte, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ProductSearchDoc to JSONL: %w", err)
	}
	return append(b, '\n'), nil
}

// FromJSONL parses a single JSONL line into a ProductSearchDoc.
func FromJSONL(line []byte) (*ProductSearchDoc, error) {
	var doc ProductSearchDoc
	if err := json.Unmarshal(line, &doc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSONL line: %w", err)
	}
	return &doc, nil
}

// WriteJSONLStream streams multiple ProductSearchDocs into an io.Writer in JSONL format.
func WriteJSONLStream(w io.Writer, docs []*ProductSearchDoc) error {
	for _, doc := range docs {
		line, err := doc.ToJSONL()
		if err != nil {
			return err
		}
		if _, err := w.Write(line); err != nil {
			return fmt.Errorf("failed to write JSONL stream line: %w", err)
		}
	}
	return nil
}

