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
