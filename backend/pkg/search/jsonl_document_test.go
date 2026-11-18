package search

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestProductSearchDoc_ToJSONL(t *testing.T) {
	now := time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)
	cat := CategoryFacet{
		ID:        1012947,
		Slug:      "sony-fe-24-70mm-f28-gm-ii-lens",
		NameEn:    "Sony FE 24-70mm f/2.8 GM II Lens",
		NameAm:    "ሶኒ FE 24-70 ሚሜ f/2.8 GM II ሌንስ",
		Level:     6,
		Path:      "/1/1378/1380/1012922/1012931/1012942/1012947/",
		Ancestors: []string{"electronics", "cameras-photography", "camera-lenses"},
	}
	loc := LocationFacet{
		ID:           1,
		Region:       "Addis Ababa",
		Subcity:      "Bole",
		SpecificArea: "Atlas",
		Coordinates:  &GeoCoord{Lat: 9.0123, Lng: 38.7890},
	}
	attrs := map[string]interface{}{
		"brand":        "Sony",
		"model":        "FE 24-70mm f/2.8 GM II",
		"lens_mount":   "Sony E-Mount",
		"max_aperture": "f/2.8",
		"filter_size":  "82mm",
	}

	doc := NewProductSearchDoc(
		40,
		"Sony FE 24-70mm f/2.8 GM II Lens",
		"Mint condition Sony G Master II standard zoom lens.",
		185000,
		"like_new",
		"sale",
		"approved",
		17,
		"gold",
		cat,
		loc,
		attrs,
		[]string{"https://images.unsplash.com/photo-1617005082133-548c4dd27f35"},
		now,
		now,
	)

	line, err := doc.ToJSONL()
	if err != nil {
		t.Fatalf("expected ToJSONL to succeed, got %v", err)
	}

	// Must end with a single newline character
	if !strings.HasSuffix(string(line), "\n") {
		t.Errorf("JSONL line must terminate with a newline character")
	}

	// Round-trip parse
	parsed, err := FromJSONL(line)
	if err != nil {
		t.Fatalf("expected FromJSONL to parse line, got %v", err)
	}

	if parsed.ProductID != 40 {
		t.Errorf("expected product_id 40, got %d", parsed.ProductID)
	}
	if parsed.Category.Slug != "sony-fe-24-70mm-f28-gm-ii-lens" {
		t.Errorf("expected category slug sony-fe-24-70mm-f28-gm-ii-lens, got %s", parsed.Category.Slug)
	}
	if parsed.Attributes["brand"] != "Sony" {
		t.Errorf("expected brand Sony, got %v", parsed.Attributes["brand"])
	}
	if len(parsed.SearchKeywords) == 0 {
		t.Errorf("expected generated search keywords, got 0")
	}
}

func TestWriteJSONLStream(t *testing.T) {
	now := time.Now()
	docs := []*ProductSearchDoc{
		NewProductSearchDoc(1, "Item 1", "Desc 1", 100, "new", "sale", "approved", 1, "normal", CategoryFacet{Slug: "cat1"}, LocationFacet{Region: "AA"}, nil, nil, now, now),
		NewProductSearchDoc(2, "Item 2", "Desc 2", 200, "used", "sale", "approved", 2, "normal", CategoryFacet{Slug: "cat2"}, LocationFacet{Region: "AA"}, nil, nil, now, now),
	}

	var buf bytes.Buffer
	err := WriteJSONLStream(&buf, docs)
	if err != nil {
		t.Fatalf("WriteJSONLStream failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 JSONL lines, got %d", len(lines))
	}
}

