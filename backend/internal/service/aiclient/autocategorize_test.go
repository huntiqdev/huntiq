package aiclient

import (
	"context"
	"testing"
)

func TestAutoCategorizeProduct_HumanErrorTolerance(t *testing.T) {
	c := NewClient("", "")

	// Test 1: Phonetic typo "ayfon 13 pro max"
	res, err := c.AutoCategorizeProduct(context.Background(), AutoCategorizeRequest{
		Title:  "ayfon 13 pro max",
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.HumanErrorCorrected {
		t.Errorf("expected HumanErrorCorrected to be true for 'ayfon'")
	}
	if res.CorrectedTerm != "iphone" {
		t.Errorf("expected CorrectedTerm 'iphone', got '%s'", res.CorrectedTerm)
	}
	if res.LegacyCategoryID != 1 {
		t.Errorf("expected LegacyCategoryID 1 (Electronics), got %d", res.LegacyCategoryID)
	}
	if res.Brand != "Apple" {
		t.Errorf("expected Brand 'Apple', got '%s'", res.Brand)
	}

	// Test 2: Phonetic typo "tovota vitz"
	resCar, err := c.AutoCategorizeProduct(context.Background(), AutoCategorizeRequest{
		Title:  "tovota vitz 2008 used",
		Locale: "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resCar.HumanErrorCorrected {
		t.Errorf("expected HumanErrorCorrected to be true for 'tovota'")
	}
	if resCar.LegacyCategoryID != 2 {
		t.Errorf("expected LegacyCategoryID 2 (Vehicles), got %d", resCar.LegacyCategoryID)
	}
	if resCar.Condition != "used" {
		t.Errorf("expected Condition 'used', got '%s'", resCar.Condition)
	}
	if resCar.SuggestedPriceETB <= 0 {
		t.Errorf("expected positive suggested ETB price, got %f", resCar.SuggestedPriceETB)
	}
}

func TestAutoCategorizeProduct_MultiLanguageParity(t *testing.T) {
	c := NewClient("", "")

	// Amharic rental test
	resAm, err := c.AutoCategorizeProduct(context.Background(), AutoCategorizeRequest{
		Title:  "ኮንዶሚኒየም 2 ክፍል ኪራይ",
		Locale: "am",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resAm.PrimaryIntent != "rent" {
		t.Errorf("expected PrimaryIntent 'rent', got '%s'", resAm.PrimaryIntent)
	}
	if resAm.LegacyCategoryID != 3 {
		t.Errorf("expected LegacyCategoryID 3 (Property), got %d", resAm.LegacyCategoryID)
	}
	if resAm.DetectedLocale != "am" {
		t.Errorf("expected DetectedLocale 'am', got '%s'", resAm.DetectedLocale)
	}
	if resAm.SuggestedDescription == "" {
		t.Errorf("expected non-empty suggested description in Amharic")
	}

	// Afaan Oromoo test
	resOm, err := c.AutoCategorizeProduct(context.Background(), AutoCategorizeRequest{
		Title:  "bilbila samsong haaraa",
		Locale: "om",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resOm.LegacyCategoryID != 1 {
		t.Errorf("expected LegacyCategoryID 1, got %d", resOm.LegacyCategoryID)
	}

	// Tigrinya test
	resTi, err := c.AutoCategorizeProduct(context.Background(), AutoCategorizeRequest{
		Title:  "ዝሽየጥ ሓዲሽ ገዛ ቪላ",
		Locale: "ti",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resTi.LegacyCategoryID != 3 {
		t.Errorf("expected LegacyCategoryID 3 (Property), got %d", resTi.LegacyCategoryID)
	}
	if resTi.Condition != "new" {
		t.Errorf("expected Condition 'new' for ሓዲሽ, got '%s'", resTi.Condition)
	}
}
