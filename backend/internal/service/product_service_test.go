package service

import (
	"testing"
)

func TestCheckProhibitedContent(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		wantBlocked bool
	}{
		// Clean items
		{"Safe electronics", "Toyota Corolla 2020", "Clean used car in good condition", false},
		{"Safe tools", "Bosch Cordless Drill 18V", "Professional hammer drill with 2 batteries", false},
		{"Safe property", "2 Bedroom Apartment in Bole", "Modern apartment for rent with parking", false},

		// Weapons (English)
		{"Weapons - Firearm", "Hunting Rifle with scope", "High precision firearm", true},
		{"Weapons - Handgun", "Pistol 9mm", "Barely used with holster", true},
		{"Weapons - Explosive", "Military grenade", "Collector item explosive", true},

		// Weapons (Amharic)
		{"Weapons - AM 1", "ጥሩ ሽጉጥ ለሽያጭ", "የግል መከላከያ", true},
		{"Weapons - AM 2", "አሮጌ ጠመንጃ", "ጥይት ያለው", true},

		// Weapons (Afaan Oromo)
		{"Weapons - OR 1", "Qawwee gurguraaf", "Haala gaarii irra jira", true},
		{"Weapons - OR 2", "Rasaasa gurgurtaa", "Meeshaa waraanaa", true},

		// Weapons (Tigrinya)
		{"Weapons - TI 1", "ቦምባ ንመሸጣ", "ሓደገኛ ኣቕሓ", true},

		// Controlled Substances
		{"Narcotics - Cocaine", "Pure cocaine pack", "Fast delivery", true},
		{"Narcotics - Heroin", "Brown heroin", "High grade", true},
		{"Narcotics - AM", "ኮኬይን ለሽያጭ", "ፈጣን ማድረስ", true},

		// Fraud / Counterfeit
		{"Counterfeit Money EN", "Counterfeit money for sale", "Undetectable fake notes", true},
		{"Counterfeit Money AM", "የሐሰት ገንዘብ", "እውነተኛ የሚመስል ኖት", true},
		{"Counterfeit Money OR", "Qarshii sobaa", "Baankii kan darbu", true},
		{"Counterfeit Money TI", "ናይ ሓሶት ቅርሺ", "ሓቀኛ ዝመስል", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkProhibitedContent(tc.title, tc.description)
			if tc.wantBlocked && err == nil {
				t.Errorf("Expected prohibited item to be blocked: title='%s', desc='%s'", tc.title, tc.description)
			}
			if !tc.wantBlocked && err != nil {
				t.Errorf("Expected clean item to pass, but got error: %v", err)
			}
		})
	}
}

func TestMapRootSlugToLegacyCategoryID(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		rootID   int
		expected int
	}{
		// Legacy ID preservation (1..11)
		{"Legacy 1 Electronics", "electronics", 1, 1},
		{"Legacy 2 Vehicles", "vehicles", 2, 2},
		{"Legacy 3 Property", "property", 3, 3},
		{"Legacy 4 Home", "home", 4, 4},
		{"Legacy 5 Fashion", "fashion", 5, 5},
		{"Legacy 6 Agriculture", "agriculture", 6, 6},
		{"Legacy 7 Business", "business", 7, 7},
		{"Legacy 8 Services", "services", 8, 8},
		{"Legacy 9 Lifestyle", "lifestyle", 9, 9},
		{"Legacy 10 Education", "education", 10, 10},
		{"Legacy 11 Other", "other", 11, 11},

		// Root domain slug mappings when rootID > 11
		{"Electronics mobile domain", "electronics-mobile", 50, 1},
		{"Computers IT domain", "computers-it", 51, 1},
		{"Vehicles cars domain", "vehicles-cars", 60, 2},
		{"Motorcycles domain", "motorcycles-bajaj", 61, 2},
		{"Property rent domain", "property-rent", 70, 3},
		{"Home appliances domain", "home-appliances", 80, 4},
		{"Clothing apparel domain", "clothing-apparel", 90, 5},
		{"Agriculture farming domain", "crops-grains-seeds", 100, 6},
		{"Industrial machinery domain", "industrial-machinery", 110, 7},
		{"Jobs recruitment domain", "jobs-recruitment", 120, 8},
		{"Professional services domain", "professional-services", 121, 8},
		{"Health medical domain", "health-medical", 130, 9},
		{"Books publishing domain", "books-publishing", 140, 10},
		{"General classifieds domain", "general-classifieds", 150, 11},

		// Unmapped / unknown slugs must safely fall back to 11 (Other)
		{"Unknown slug", "unknown-domain-xyz", 999, 11},
		{"Empty slug with large root ID", "", 163, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapRootSlugToLegacyCategoryID(tt.slug, tt.rootID)
			if got != tt.expected {
				t.Errorf("MapRootSlugToLegacyCategoryID(%q, %d) = %d; want %d", tt.slug, tt.rootID, got, tt.expected)
			}
			if got < 1 || got > 11 {
				t.Fatalf("CRITICAL: Category mapping returned out-of-range ID %d (must be between 1 and 11 to satisfy categories(id) FK)", got)
			}
		})
	}
}

func TestCategoryHierarchyLegacyRangeGuarantee(t *testing.T) {
	// Child category IDs (e.g. 163 for Refrigerators, 250 for Laptops) must never be assigned
	// directly to the legacy categories table foreign key.
	childCategoryIDs := []int{12, 13, 50, 100, 163, 250, 500, 1200}
	for _, childID := range childCategoryIDs {
		// Calling MapRootSlugToLegacyCategoryID with an unmapped slug or child ID
		legacyID := MapRootSlugToLegacyCategoryID("", childID)
		if legacyID > 11 {
			t.Fatalf("Child category ID %d mapped to invalid legacy category ID %d; must not exceed 11", childID, legacyID)
		}
	}
}
