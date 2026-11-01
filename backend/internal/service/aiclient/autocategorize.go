package aiclient

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// AutoCategorizeRequest contains parameters for listing auto-categorization.
type AutoCategorizeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Locale      string `json:"locale,omitempty"`
}

// AutoCategorizeResult represents AI intelligence for product listings.
type AutoCategorizeResult struct {
	DetectedTitle        string   `json:"detected_title"`
	PrimaryIntent        string   `json:"primary_intent"` // "sale" or "rent"
	CategoryHierarchyID  int      `json:"category_hierarchy_id"`
	CategorySlug         string   `json:"category_slug"`
	RootCategoryName     string   `json:"root_category_name"`
	SubcategoryName      string   `json:"subcategory_name"`
	LeafDestinationName  string   `json:"leaf_destination_name"`
	FullBranchPath       string   `json:"full_branch_path"`       // English path
	BranchPathLocalized  string   `json:"branch_path_localized"`  // Locale path
	LegacyCategoryID     int      `json:"legacy_category_id"`
	Condition            string   `json:"condition"`              // "new" or "used"
	SuggestedPriceETB    float64  `json:"suggested_price_etb"`
	PriceRangeMin        float64  `json:"price_range_min"`
	PriceRangeMax        float64  `json:"price_range_max"`
	SuggestedDescription string   `json:"suggested_description"`  // Localized ready-to-use explanation
	DetectedLocale       string   `json:"detected_locale"`
	Brand                string   `json:"brand,omitempty"`
	Model                string   `json:"model,omitempty"`
	Confidence           float64  `json:"confidence"`
	HumanErrorCorrected  bool     `json:"human_error_corrected"`
	OriginalTerm         string   `json:"original_term"`
	CorrectedTerm        string   `json:"corrected_term"`
}

// BranchDefinition defines a concrete terminal leaf destination with its 3-tier hierarchy.
type BranchDefinition struct {
	RootSlug            string
	RootNameEN          string
	RootNameAM          string
	RootNameOM          string
	RootNameTI          string

	SubcatSlug          string
	SubcatNameEN        string
	SubcatNameAM        string
	SubcatNameOM        string
	SubcatNameTI        string

	LeafSlug            string
	LeafNameEN          string
	LeafNameAM          string
	LeafNameOM          string
	LeafNameTI          string

	LegacyCategoryID    int
	FallbackHierarchyID int
	TypicalIntent       string // "sale" or "rent"
	DefaultCondition    string // "used" or "new"
	TypicalPriceETB     float64
	PriceMinETB         float64
	PriceMaxETB         float64

	Keywords            []string
	PhoneticTypos       []string
}

// canonicalBranches holds the master leaf destinations for Ethiopian marketplace categorization.
var canonicalBranches = []BranchDefinition{
	// 1. SMARTPHONES & MOBILE PHONES
	{
		RootSlug:            "electronics-mobile",
		RootNameEN:          "Electronics and Mobile Devices",
		RootNameAM:          "ኤሌክትሮኒክስ እና ሞባይል ስልኮች",
		RootNameOM:          "Elektirooniksii fi Bilbila",
		RootNameTI:          "ኤሌክትሮኒክስን ሞባይልን",

		SubcatSlug:          "smartphones-mobile",
		SubcatNameEN:        "Smartphones & Mobile Handsets",
		SubcatNameAM:        "ስማርትፎኖች እና ሞባይል ስልኮች",
		SubcatNameOM:        "Ismaartfoonii fi Bilbila",
		SubcatNameTI:        "ስማርትፎናትን ሞባይል ስልኪን",

		LeafSlug:            "ultra-flagship-smartphones",
		LeafNameEN:          "Smartphones & Touchscreen Handsets",
		LeafNameAM:          "ስማርትፎኖች እና ዘመናዊ ስልኮች",
		LeafNameOM:          "Ismaartfoonii fi Bilbila Ammayyaa",
		LeafNameTI:          "ስማርትፎናትን ዘመናዊ ስልክታትን",

		LegacyCategoryID:    1,
		FallbackHierarchyID: 101,
		TypicalIntent:       "sale",
		DefaultCondition:    "used",
		TypicalPriceETB:     45000,
		PriceMinETB:         12000,
		PriceMaxETB:         185000,

		Keywords: []string{
			"phone", "smartphone", "smartphones", "mobile", "cellphone", "android", "ios",
			"iphone", "apple", "samsung", "galaxy", "redmi", "xiaomi", "tecno", "infinix",
			"huawei", "itel", "pixel", "oneplus", "oppo", "vivo",
			"ስልክ", "ሞባይል", "ስማርትፎን", "ስማርት ፎን", "አይፎን", "ሳምሰንግ", "ሬድሚ", "ቴክኖ", "ኢንፊኒክስ", "ሁዋዌ",
			"bilbila", "moobaayilii", "ismaartfoonii",
			"ስልኪ", "ስማርትፎን",
		},
		PhoneticTypos: []string{
			"ayfon", "aifon", "aiphon", "i phone", "i-phone", "iphon", "ayfone", "ኣይፎን",
			"samsong", "sammsung", "samsun", "sumsung", "sam sung", "ሳምሶንግ",
			"redmy", "red mi", "shaomi", "ሻኦሚ",
			"tecnno", "tekno", "teckno", "ተክኖ",
			"infinx", "infinex", "እንፊኒክስ",
			"huawi", "huway", "ሁዋይ",
		},
	},

	// 2. TABLETS & IPADS
	{
		RootSlug:            "electronics-mobile",
		RootNameEN:          "Electronics and Mobile Devices",
		RootNameAM:          "ኤሌክትሮኒክስ እና ሞባይል ስልኮች",
		RootNameOM:          "Elektirooniksii fi Bilbila",
		RootNameTI:          "ኤሌክትሮኒክስን ሞባይልን",

		SubcatSlug:          "tablets-ereaders",
		SubcatNameEN:        "Tablets, iPads and E-Readers",
		SubcatNameAM:        "ታብሌቶች እና አይፓዶች",
		SubcatNameOM:        "Taableetii fi Aaypaadii",
		SubcatNameTI:        "ታብሌታትን ኣይፓዳትን",

		LeafSlug:            "touchscreen-tablets-ipads",
		LeafNameEN:          "Touchscreen Tablets & iPads",
		LeafNameAM:          "ታብሌቶች እና አይፓዶች",
		LeafNameOM:          "Taableetii fi Aaypaadii",
		LeafNameTI:          "ታብሌታትን ኣይፓዳትን",

		LegacyCategoryID:    1,
		FallbackHierarchyID: 102,
		TypicalIntent:       "sale",
		DefaultCondition:    "used",
		TypicalPriceETB:     35000,
		PriceMinETB:         14000,
		PriceMaxETB:         120000,

		Keywords: []string{
			"tablet", "tablets", "ipad", "ipad pro", "ipad air", "tab", "galaxy tab",
			"ታብሌት", "አይፓድ", "ታብ", "ኣይፓድ",
			"taableetii", "aaypaadii",
		},
		PhoneticTypos: []string{
			"tablit", "tablate", "aipad", "i pad", "i-pad", "ታብሌት",
		},
	},

	// 3. COMPUTERS & LAPTOPS
	{
		RootSlug:            "computers-it",
		RootNameEN:          "Computers, IT and Office Technology",
		RootNameAM:          "ኮምፒውተሮች እና የቢሮ ቴክኖሎጂ",
		RootNameOM:          "Kompiitara fi Teekinoolojii",
		RootNameTI:          "ኮምፒዩተራትን ናይ ቢሮ ቴክኖሎጂን",

		SubcatSlug:          "laptops-notebooks",
		SubcatNameEN:        "Laptops & Portable Computers",
		SubcatNameAM:        "ላፕቶፖች እና ፖርቴብል ኮምፒውተሮች",
		SubcatNameOM:        "Laaptooppii fi Kompiitara",
		SubcatNameTI:        "ላፕቶፓትን ኮምፒዩተራትን",

		LeafSlug:            "high-performance-laptops",
		LeafNameEN:          "Laptops & Notebook Computers",
		LeafNameAM:          "ላፕቶፖች እና ኖትቡኮች",
		LeafNameOM:          "Laaptooppii fi Nootbuukii",
		LeafNameTI:          "ላፕቶፓትን ኖትቡካትን",

		LegacyCategoryID:    1,
		FallbackHierarchyID: 201,
		TypicalIntent:       "sale",
		DefaultCondition:    "used",
		TypicalPriceETB:     38000,
		PriceMinETB:         18000,
		PriceMaxETB:         190000,

		Keywords: []string{
			"laptop", "laptops", "notebook", "macbook", "macbook pro", "macbook air",
			"dell", "hp", "lenovo", "thinkpad", "asus", "acer", "toshiba", "pc", "computer",
			"ላፕቶፕ", "ኮምፒውተር", "ማክቡክ", "ዴል", "ኤችፒ", "ሌኖቮ",
			"kompiitara", "laaptooppii", "ኖትቡክ",
		},
		PhoneticTypos: []string{
			"laptopp", "loptop", "labtop", "lab top", "lap top", "ላፕ ቶፕ", "ላብቶፕ",
			"macbok", "mac bok", "makbook", "ማክ ቡክ",
			"del", "del latitude", "dell inspiron",
		},
	},

	// 4. TELEVISIONS & HOME AUDIO
	{
		RootSlug:            "tv-audio",
		RootNameEN:          "TVs, Audio and Home Entertainment",
		RootNameAM:          "ቴሌቪዥን እና የድምጽ ሲስተሞች",
		RootNameOM:          "Televiziyinii fi Sagalee",
		RootNameTI:          "ቴሌቪዥንን ድምጺ ሲስተማትን",

		SubcatSlug:          "smart-tvs-displays",
		SubcatNameEN:        "Smart TVs & Flat Screens",
		SubcatNameAM:        "ስማርት ቴሌቪዥኖች እና ስክሪኖች",
		SubcatNameOM:        "Televiziyinii Ismaartii",
		SubcatNameTI:        "ስማርት ቴሌቪዥናትን ስክሪናትን",

		LeafSlug:            "led-smart-televisions",
		LeafNameEN:          "Smart TVs & Flat Panel Displays",
		LeafNameAM:          "ስማርት ቴሌቪዥኖች እና ጠፍጣፋ ስክሪኖች",
		LeafNameOM:          "Televiziyinii Ismaartii fi Iskiriinii",
		LeafNameTI:          "ስማርት ቴሌቪዥንን ስክሪናትን",

		LegacyCategoryID:    1,
		FallbackHierarchyID: 301,
		TypicalIntent:       "sale",
		DefaultCondition:    "used",
		TypicalPriceETB:     28000,
		PriceMinETB:         9000,
		PriceMaxETB:         110000,

		Keywords: []string{
			"tv", "television", "smart tv", "led tv", "4k tv", "oled", "screen",
			"ቴሌቪዥን", "ቲቪ", "ስማርት ቲቪ", "ሪሲቨር",
			"televizhinii", "tivii", "iskiriinii",
		},
		PhoneticTypos: []string{
			"televishin", "televisn", "televizhn", "ቲ ቪ", "ቴሌቪዥን",
		},
	},

	// 5. PASSENGER CARS (SEDANS & HATCHBACKS)
	{
		RootSlug:            "vehicles-cars",
		RootNameEN:          "Vehicles and Automobiles",
		RootNameAM:          "ተሽከርካሪዎች እና የመንገደኛ መኪኖች",
		RootNameOM:          "Konkolaattota fi Awutoomoobilii",
		RootNameTI:          "ተሽከርከርትን ናይ ህዝቢ መካይንን",

		SubcatSlug:          "passenger-cars",
		SubcatNameEN:        "Passenger Cars & Compact Vehicles",
		SubcatNameAM:        "የመንገደኛ መኪኖች እና አነስተኛ ተሽከርካሪዎች",
		SubcatNameOM:        "Konkolaattota Awutoomoobilii",
		SubcatNameTI:        "ናይ መንገደኛ መካይን",

		LeafSlug:            "compact-hatchback-cars",
		LeafNameEN:          "Compact & Hatchback Cars (Vitz / Yaris)",
		LeafNameAM:          "አነስተኛ የመንገደኛ መኪኖች (ቪትዝ / ያሪስ)",
		LeafNameOM:          "Konkolaattota Xixinnoo (Vitz / Yaris)",
		LeafNameTI:          "ቀለልቲ መካይን (ቪትዝ / ያሪስ)",

		LegacyCategoryID:    2,
		FallbackHierarchyID: 401,
		TypicalIntent:       "sale",
		DefaultCondition:    "used",
		TypicalPriceETB:     1350000,
		PriceMinETB:         850000,
		PriceMaxETB:         3200000,

		Keywords: []string{
			"car", "cars", "automobile", "toyota", "vitz", "yaris", "corolla", "suzuki",
			"dzire", "swift", "hyundai", "atos", "i10", "accent", "hatchback", "sedan",
			"መኪና", "ተሽከርካሪ", "ቶዮታ", "ቪትዝ", "ያሪስ", "ኮሮላ", "ሱዙኪ", "ዲዛየር", "ስዊፍት", "ሂዩንዳይ",
			"konkolaataa", "otoomobil", "awutoomoobilii",
			"መኪና", "ናይ ውልቂ መኪና",
		},
		PhoneticTypos: []string{
			"tovota", "toyotta", "toyyota", "toyta", "toyoto", "ቶዮታ",
			"vits", "vitzz", "vitss", "ቪትስ",
			"corola", "corrola", "ኮሮላ ኤስ",
			"suzki", "suzky", "ሱዝኪ",
			"hundai", "hyndai", "ህዩንዳይ",
		},
	},

	// 6. SUVS & 4WD VEHICLES
	{
		RootSlug:            "vehicles-cars",
		RootNameEN:          "Vehicles and Automobiles",
		RootNameAM:          "ተሽከርካሪዎች እና የመንገደኛ መኪኖች",
		RootNameOM:          "Konkolaattota fi Awutoomoobilii",
		RootNameTI:          "ተሽከርከርትን ናይ ህዝቢ መካይንን",

		SubcatSlug:          "suvs-crossovers",
		SubcatNameEN:        "SUVs, Crossovers & 4WD Off-Roaders",
		SubcatNameAM:        "ኤስዩቪ እና ባለ 4 ጎማ ተሽከርካሪዎች",
		SubcatNameOM:        "Konkolaattota SUV fi 4WD",
		SubcatNameTI:        "ኤስዩቪን ባለ 4 ጎማ መካይንን",

		LeafSlug:            "compact-midsize-suvs",
		LeafNameEN:          "SUVs & 4WD Off-Road Vehicles (RAV4 / Prado)",
		LeafNameAM:          "ኤስዩቪ እና ባለ 4 ጎማ ተሽከርካሪዎች (ራቭ4 / ፕራዶ)",
		LeafNameOM:          "Konkolaataa SUV fi 4WD (RAV4 / Prado)",
		LeafNameTI:          "ኤስዩቪን ባለ 4 ጎማ መካይንን (ራቭ4 / ፕራዶ)",

		LegacyCategoryID:    2,
		FallbackHierarchyID: 402,
		TypicalIntent:       "sale",
		DefaultCondition:    "used",
		TypicalPriceETB:     3800000,
		PriceMinETB:         1900000,
		PriceMaxETB:         14500000,

		Keywords: []string{
			"suv", "crossover", "4wd", "rav4", "prado", "land cruiser", "hilux", "d-max",
			"ራቭ4", "ፕራዶ", "ላንድ ክሩዘር", "ሃይሉክስ", "ፒክአፕ",
		},
		PhoneticTypos: []string{
			"rav 4", "rav-4", "ራቭ 4", "ራፍ4", "ፕራዶ", "ሃይ ሉክስ",
		},
	},

	// 7. MOTORCYCLES & THREE-WHEELERS (BAJAJ)
	{
		RootSlug:            "motorcycles-bajaj",
		RootNameEN:          "Motorcycles, Scooters and Three-Wheelers",
		RootNameAM:          "ሞተር ሳይክሎች እና ባጃጆች",
		RootNameOM:          "Moophedii fi Baajaajii",
		RootNameTI:          "ሞተር ሳይክላትን ባጃጃትን",

		SubcatSlug:          "three-wheelers-bajaj",
		SubcatNameEN:        "Three-Wheelers & Passenger Bajaj",
		SubcatNameAM:        "ባጃጆች እና ባለ ሦስት እግር ተሽከርካሪዎች",
		SubcatNameOM:        "Baajaajii fi Geejjiba Miila Sadii",
		SubcatNameTI:        "ባጃጃትን ባለ ሠለስተ እግሪ መጓዓዝያን",

		LeafSlug:            "passenger-cargo-bajaj",
		LeafNameEN:          "Passenger & Cargo Three-Wheelers (Bajaj / TVS)",
		LeafNameAM:          "ባጃጅ እና ባለ ሦስት እግር ተሽከርካሪ (ባጃጅ / ቲቪኤስ)",
		LeafNameOM:          "Baajaajii fi TVS",
		LeafNameTI:          "ባጃጅን ቲቪኤስን",

		LegacyCategoryID:    2,
		FallbackHierarchyID: 501,
		TypicalIntent:       "sale",
		DefaultCondition:    "used",
		TypicalPriceETB:     320000,
		PriceMinETB:         180000,
		PriceMaxETB:         520000,

		Keywords: []string{
			"bajaj", "three wheeler", "tvs", "motorcycle", "motorbike", "scooter",
			"ባጃጅ", "ቲቪኤስ", "ሞተር", "ሞተርሳይክል", "ሞተር ሳይክል",
			"baajaajii", "moophedii", "motora",
		},
		PhoneticTypos: []string{
			"bajaji", "bajajj", "baajaaji", "ባጃጅ",
			"moter", "motorsaikil", "motersaikil", "ሞተርሳይክል",
		},
	},

	// 8. PROPERTY FOR RENT (CONDOMINIUMS & APARTMENTS)
	{
		RootSlug:            "property-rent",
		RootNameEN:          "Property for Rent and Lease",
		RootNameAM:          "የሚከራዩ ቤቶች እና አፓርታማዎች",
		RootNameOM:          "Mana Kiraa",
		RootNameTI:          "ዝካረዩ ኣባይትን ኣፓርታማታትን",

		SubcatSlug:          "apartment-condo-rentals",
		SubcatNameEN:        "Apartment & Condominium Rentals",
		SubcatNameAM:        "የኮንዶሚኒየም እና የአፓርታማ ኪራይ",
		SubcatNameOM:        "Kiraa Mana Kondoominiyeemii",
		SubcatNameTI:        "ክራይ ኮንዶሚኒየምን ኣፓርታማን",

		LeafSlug:            "condominium-rentals",
		LeafNameEN:          "1-3 Bedroom Condominium Rentals",
		LeafNameAM:          "ባለ 1-3 መኝታ ኮንዶሚኒየም ኪራይ",
		LeafNameOM:          "Kiraa Mana Kondoominiyeemii Kutaa 1-3",
		LeafNameTI:          "ክራይ ኮንዶሚኒየም 1-3 ክፍሊ",

		LegacyCategoryID:    3,
		FallbackHierarchyID: 601,
		TypicalIntent:       "rent",
		DefaultCondition:    "used",
		TypicalPriceETB:     16000, // Monthly ETB
		PriceMinETB:         7000,
		PriceMaxETB:         45000,

		Keywords: []string{
			"rent", "rental", "condo", "condominium", "apartment", "studio", "flat",
			"ኪራይ", "የሚከራይ", "ኮንዶሚኒየም", "ኮንዶ", "አፓርታማ", "ስቱዲዮ", "ክፍል",
			"kiraa", "kondoominiyeemii", "apaartaamaa", "kutaa",
			"ክራይ", "ዝካረ", "ኮንዶሚኒየም",
		},
		PhoneticTypos: []string{
			"kondominiyem", "kondominyem", "kondominium", "ኮንዶምኒየም", "ኮንዶሚኒየም",
			"apartama", "apartma", "apartman", "ኣፓርታማ",
			"kiray", "kira", "የሚከራይ ቤት",
		},
	},

	// 9. PROPERTY FOR SALE (VILLAS & HOUSES)
	{
		RootSlug:            "property-sale",
		RootNameEN:          "Property for Sale",
		RootNameAM:          "የሚሸጡ ቤቶች እና ህንጻዎች",
		RootNameOM:          "Mana Gurguraa",
		RootNameTI:          "ዝሽየጡ ኣባይትን ህንጻታትን",

		SubcatSlug:          "residential-houses-villas",
		SubcatNameEN:        "Residential Houses & Modern Villas",
		SubcatNameAM:        "የመኖሪያ ቤቶች እና ቪላዎች",
		SubcatNameOM:        "Manneen Jireenyaa fi Viillaa",
		SubcatNameTI:        "ናይ መንበሪ ኣባይትን ቪላታትን",

		LeafSlug:            "residential-villas-houses",
		LeafNameEN:          "Residential Houses, G+1/G+2 Villas",
		LeafNameAM:          "የመኖሪያ ቤቶች እና ጂ+1/ጂ+2 ቪላዎች",
		LeafNameOM:          "Viillaa fi Mana Jireenyaa G+1/G+2",
		LeafNameTI:          "ናይ መንበሪ ኣባይትን ቪላታትን G+1",

		LegacyCategoryID:    3,
		FallbackHierarchyID: 701,
		TypicalIntent:       "sale",
		DefaultCondition:    "used",
		TypicalPriceETB:     18500000,
		PriceMinETB:         4500000,
		PriceMaxETB:         65000000,

		Keywords: []string{
			"house", "villa", "villas", "home", "real estate", "g+1", "g+2", "compound",
			"ቤት", "የሚሸጥ ቤት", "ቪላ", "ጂ+1", "ሪልስቴት", "ግቢ",
			"mana", "mana gurguraa", "viillaa",
			"ገዛ", "ዝሽየጥ ገዛ", "ቪላ",
		},
		PhoneticTypos: []string{
			"vila", "villah", "viilla", "ቪላ ቤት",
			"g1", "g 1", "g+1 villa",
		},
	},

	// 10. FURNITURE (SOFAS & LIVING ROOM)
	{
		RootSlug:            "furniture-furnishings",
		RootNameEN:          "Furniture and Interior Furnishings",
		RootNameAM:          "የቤት እቃዎች እና ፈርኒቸር",
		RootNameOM:          "Meeshaa Manaa fi Barcuma",
		RootNameTI:          "ፈርኒቸርን ናይ ውሽጢ ገዛ ኣቑሑትን",

		SubcatSlug:          "living-room-furniture",
		SubcatNameEN:        "Living Room Furniture & Sofas",
		SubcatNameAM:        "የሳሎን እቃዎች እና ሶፋዎች",
		SubcatNameOM:        "Meeshaa Saaloonaa fi Soofaa",
		SubcatNameTI:        "ናይ ሳሎን ኣቑሑትን ሶፋታትን",

		LeafSlug:            "living-room-sofas-couches",
		LeafNameEN:          "Living Room Sofas & L-Shape Couches",
		LeafNameAM:          "የሳሎን ሶፋዎች እና ኤል ቅርጽ ሶፋ",
		LeafNameOM:          "Soofaa Saaloonaa fi L-Shape",
		LeafNameTI:          "ናይ ሳሎን ሶፋታትን",

		LegacyCategoryID:    4,
		FallbackHierarchyID: 801,
		TypicalIntent:       "sale",
		DefaultCondition:    "used",
		TypicalPriceETB:     32000,
		PriceMinETB:         14000,
		PriceMaxETB:         95000,

		Keywords: []string{
			"sofa", "couch", "living room", "l-shape", "sectional", "chair", "table", "bed",
			"ሶፋ", "የሳሎን ሶፋ", "ኤል ቅርጽ ሶፋ", "አልጋ", "ቁምሳጥን", "የምግብ ጠረጴዛ", "ወንበር",
			"soofaa", "farnichara", "siree", "barcuma",
			"ሶፋ", "ዓራት", "ጠረጴዛ",
		},
		PhoneticTypos: []string{
			"sofaa", "sofa set", "l shape sofa", "ሶፋዎች",
			"alga", "algaa", "ፍራሽ",
			"qumsatin", "qumsaxin", "ቁም ሳጥን",
		},
	},

	// 11. CLOTHING & FASHION
	{
		RootSlug:            "clothing-apparel",
		RootNameEN:          "Clothing, Fashion and Apparel",
		RootNameAM:          "አልባሳት እና ፋሽን",
		RootNameOM:          "Uffata fi Faashinii",
		RootNameTI:          "ክዳውንትን ፋሽንን",

		SubcatSlug:          "traditional-modern-apparel",
		SubcatNameEN:        "Traditional & Modern Apparel",
		SubcatNameAM:        "ባህላዊ እና ዘመናዊ አልባሳት",
		SubcatNameOM:        "Uffata Aadaa fi Ammayyaa",
		SubcatNameTI:        "ባህላውን ዘመናውን ክዳውንቲ",

		LeafSlug:            "dresses-habesha-kemis",
		LeafNameEN:          "Traditional Habesha Kemis & Modern Dresses",
		LeafNameAM:          "የሐበሻ ቀሚስ እና ዘመናዊ ቀሚሶች",
		LeafNameOM:          "Uffata Aadaa fi Qamisii",
		LeafNameTI:          "ናይ ሓበሻ ቀሚስን ዘመናዊ ክዳውንትን",

		LegacyCategoryID:    5,
		FallbackHierarchyID: 901,
		TypicalIntent:       "sale",
		DefaultCondition:    "new",
		TypicalPriceETB:     6500,
		PriceMinETB:         2500,
		PriceMaxETB:         28000,

		Keywords: []string{
			"dress", "habesha kemis", "shirt", "t-shirt", "jeans", "pants", "suit", "clothes",
			"ቀሚስ", "የሐበሻ ቀሚስ", "ሸሚዝ", "ቲሸርት", "ሱሪ", "ልብስ", "ኮት",
			"uffata", "wayyaa", "qamisii",
			"ክዳን", "ቀሚስ", "ሸሚዝ",
		},
		PhoneticTypos: []string{
			"shert", "t-shert", "tshert", "ሸሚዝ",
			"surri", "surr", "jins", "ሱሪ",
			"habesha kemis", "habeshakemis", "ሐበሻ ቀሚስ",
		},
	},

	// 12. SHOES & FOOTWEAR
	{
		RootSlug:            "shoes-bags",
		RootNameEN:          "Shoes, Bags and Accessories",
		RootNameAM:          "ጫማዎች፣ ቦርሳዎች እና ጌጣጌጦች",
		RootNameOM:          "Kophee, Boorsaa fi Faaya",
		RootNameTI:          "ሳእንን ቦርሳታትን",

		SubcatSlug:          "footwear-shoes",
		SubcatNameEN:        "Footwear, Sneakers & Boots",
		SubcatNameAM:        "የስፖርት እና የቆዳ ጫማዎች",
		SubcatNameOM:        "Kophee Ispoortii fi Gogaa",
		SubcatNameTI:        "ሳእኒ ስፖርትን ቆርበትን",

		LeafSlug:            "sneakers-casual-shoes",
		LeafNameEN:          "Sneakers, Casual & Leather Shoes",
		LeafNameAM:          "ስኒከሮች እና የቆዳ ጫማዎች",
		LeafNameOM:          "Kophee Isniikarii fi Gogaa",
		LeafNameTI:          "ስኒከራትን ናይ ቆርበት ሳእንን",

		LegacyCategoryID:    5,
		FallbackHierarchyID: 902,
		TypicalIntent:       "sale",
		DefaultCondition:    "new",
		TypicalPriceETB:     3500,
		PriceMinETB:         1800,
		PriceMaxETB:         14000,

		Keywords: []string{
			"shoes", "sneakers", "boots", "sandals", "loafers", "leather shoes",
			"ጫማ", "ስኒከር", "ሳንደል", "የቆዳ ጫማ",
			"kophee", "sandala",
			"ሳእኒ", "ሳንደል",
		},
		PhoneticTypos: []string{
			"chama", "sniker", "sneker", "ጫማ",
		},
	},

	// 13. AGRICULTURE & LIVESTOCK
	{
		RootSlug:            "agriculture-farming",
		RootNameEN:          "Agriculture, Farming and Agricultural Land",
		RootNameAM:          "ግብርና፣ እርሻ እና የገጠር ልማት",
		RootNameOM:          "Qonnaa fi Misooma Baadiyyaa",
		RootNameTI:          "ሕርሻን ናይ ገጠር ልምዓትን",

		SubcatSlug:          "livestock-animals",
		SubcatNameEN:        "Livestock, Cattle and Poultry",
		SubcatNameAM:        "የቤት እንስሳት፣ በሬ እና ዶሮ",
		SubcatNameOM:        "Horii, Saawwan fi Lukkuu",
		SubcatNameTI:        "ናይ ገዛ እንስሳታት፣ ብዕራይን ደርሆን",

		LeafSlug:            "cattle-oxen-sheep",
		LeafNameEN:          "Cattle, Oxen, Sheep and Livestock",
		LeafNameAM:          "በሬ፣ ሰንጋ፣ በግ እና የቤት እንስሳት",
		LeafNameOM:          "Sangaa, Hoolaa fi Horii",
		LeafNameTI:          "ብዕራይ፣ ሰንጋ፣ በጊዕን",

		LegacyCategoryID:    6,
		FallbackHierarchyID: 1001,
		TypicalIntent:       "sale",
		DefaultCondition:    "new",
		TypicalPriceETB:     75000,
		PriceMinETB:         6000,
		PriceMaxETB:         180000,

		Keywords: []string{
			"ox", "oxen", "cattle", "sheep", "goat", "chicken", "poultry", "livestock",
			"በሬ", "ሰንጋ", "በግ", "ፍየል", "ዶሮ", "እንስሳት", "ላም",
			"sangaa", "hoolaa", "re'ee", "lukkuu", "horii",
			"ብዕራይ", "ሰንጋ", "በጊዕ", "ደርሆ",
		},
		PhoneticTypos: []string{
			"bere", "senga", "beg", "doro", "korma doro", "በሬ", "በግ", "ዶሮ",
		},
	},
}

// humanErrorPhoneticReplacements maps common Ethiopian phonetic misspellings to normalized search tokens.
var humanErrorPhoneticReplacements = map[string]string{
	"ayfon":        "iphone",
	"aifon":        "iphone",
	"aiphon":       "iphone",
	"i phone":      "iphone",
	"i-phone":      "iphone",
	"iphon":        "iphone",
	"ayfone":       "iphone",
	"samsong":      "samsung",
	"sammsung":     "samsung",
	"samsun":       "samsung",
	"sumsung":      "samsung",
	"redmy":        "redmi",
	"shaomi":       "xiaomi",
	"tecnno":       "tecno",
	"tekno":        "tecno",
	"teckno":       "tecno",
	"infinx":       "infinix",
	"infinex":      "infinix",
	"huawi":        "huawei",
	"huway":        "huawei",
	"laptopp":      "laptop",
	"loptop":       "laptop",
	"labtop":       "laptop",
	"macbok":       "macbook",
	"makbook":      "macbook",
	"tovota":       "toyota",
	"toyotta":      "toyota",
	"toyyota":      "toyota",
	"toyta":        "toyota",
	"toyoto":       "toyota",
	"vits":         "vitz",
	"vitzz":        "vitz",
	"corola":       "corolla",
	"corrola":      "corolla",
	"suzki":        "suzuki",
	"suzky":        "suzuki",
	"hundai":       "hyundai",
	"hyndai":       "hyundai",
	"bajaji":       "bajaj",
	"bajajj":       "bajaj",
	"moter":        "motorcycle",
	"motorsaikil":  "motorcycle",
	"motersaikil":  "motorcycle",
	"kondominiyem": "condominium",
	"kondominyem":  "condominium",
	"kondominium":  "condominium",
	"apartama":     "apartment",
	"apartma":      "apartment",
	"vila":         "villa",
	"alga":         "bed",
	"sofaa":        "sofa",
	"qumsatin":     "wardrobe",
	"qumsaxin":     "wardrobe",
	"shert":        "shirt",
	"t-shert":      "t-shirt",
	"tshert":       "t-shirt",
	"surri":        "trousers",
	"surr":         "trousers",
	"chama":        "shoes",
	"bere":         "ox",
	"senga":        "ox",
	"beg":          "sheep",
	"doro":         "chicken",
	"teff":         "teff grain",
}

// AutoCategorizeProduct analyzes seller input text and determines the concrete category hierarchy branch.
func (c *Client) AutoCategorizeProduct(ctx context.Context, req AutoCategorizeRequest) (*AutoCategorizeResult, error) {
	return FallbackAutoCategorize(req), nil
}

// FallbackAutoCategorize deterministically classifies any listing with human error tolerance.
func FallbackAutoCategorize(req AutoCategorizeRequest) *AutoCategorizeResult {
	rawInput := strings.TrimSpace(req.Title)
	if rawInput == "" {
		rawInput = strings.TrimSpace(req.Description)
	}
	if rawInput == "" {
		rawInput = "General Item"
	}

	normLocale := normalizeLocaleCode(req.Locale)
	if normLocale == "" {
		if ContainsGeez(rawInput) {
			normLocale = "am"
		} else if isOromoText(rawInput) {
			normLocale = "om"
		} else {
			normLocale = "en"
		}
	}

	lowerText := strings.ToLower(rawInput)

	// Step 1: Detect Human Errors and Phonetic Transliterations
	correctedTerm := ""
	originalTerm := ""
	humanErrorFound := false

	words := strings.Fields(lowerText)
	for _, w := range words {
		// Clean punctuation
		cleanedWord := strings.Trim(w, ",.!?\"'():-")
		if repl, ok := humanErrorPhoneticReplacements[cleanedWord]; ok {
			originalTerm = cleanedWord
			correctedTerm = repl
			humanErrorFound = true
			break
		}
	}

	// Step 2: Match Best Canonical Branch
	var matchedBranch BranchDefinition
	bestScore := -1

	for _, branch := range canonicalBranches {
		score := 0

		// Match phonetic typos
		for _, typo := range branch.PhoneticTypos {
			if strings.Contains(lowerText, strings.ToLower(typo)) {
				score += 8
			}
		}

		// Match keywords
		for _, kw := range branch.Keywords {
			kwLower := strings.ToLower(kw)
			if strings.Contains(lowerText, kwLower) {
				if len(kwLower) > 4 {
					score += 5
				} else {
					score += 3
				}
			}
		}

		// If human error corrected term matches
		if correctedTerm != "" {
			for _, kw := range branch.Keywords {
				if strings.Contains(strings.ToLower(kw), correctedTerm) {
					score += 6
				}
			}
		}

		if score > bestScore {
			bestScore = score
			matchedBranch = branch
		}
	}

	// If no strong match found, default to electronics/phones (most common listing domain)
	if bestScore <= 0 {
		matchedBranch = canonicalBranches[0]
	}

	// Step 3: Determine Condition
	condition := matchedBranch.DefaultCondition
	newKeywords := []string{"brand new", "brand-new", "unused", "አዲስ", "ሓዲሽ", "haaraa", "qulqulluu"}
	usedKeywords := []string{"second hand", "second-hand", "ያገለገለ", "ዝተጠቐመሉ", "tajaajila kan kenne"}

	matchedCondition := false
	for _, kw := range newKeywords {
		if strings.Contains(lowerText, kw) {
			condition = "new"
			matchedCondition = true
			break
		}
	}
	if !matchedCondition {
		for _, kw := range usedKeywords {
			if strings.Contains(lowerText, kw) {
				condition = "used"
				matchedCondition = true
				break
			}
		}
	}
	if !matchedCondition {
		// Single ASCII word check
		for _, w := range words {
			cw := strings.Trim(w, ",.!?\"'():-")
			if cw == "new" {
				condition = "new"
				break
			} else if cw == "used" {
				condition = "used"
				break
			}
		}
	}

	// Step 4: Determine Intent (Sale vs Rent)
	intent := matchedBranch.TypicalIntent
	rxRent := regexp.MustCompile(`(?i)\b(rent|rental|lease|ኪራይ|የሚከራይ|kiraa|ክራይ|ዝካረ)\b`)
	if rxRent.MatchString(lowerText) {
		intent = "rent"
	}

	// Step 5: Extract brand and model if recognized
	brand := ""
	model := ""
	if strings.Contains(lowerText, "iphone") || strings.Contains(lowerText, "አይፎን") || strings.Contains(lowerText, "ayfon") {
		brand = "Apple"
		model = extractModelSuffix(rawInput, "iphone", "አይፎን", "ayfon")
	} else if strings.Contains(lowerText, "samsung") || strings.Contains(lowerText, "ሳምሰንግ") || strings.Contains(lowerText, "samsong") {
		brand = "Samsung"
		model = extractModelSuffix(rawInput, "samsung", "ሳምሰንግ", "samsong")
	} else if strings.Contains(lowerText, "toyota") || strings.Contains(lowerText, "ቶዮታ") || strings.Contains(lowerText, "tovota") {
		brand = "Toyota"
		model = extractModelSuffix(rawInput, "toyota", "ቶዮታ", "tovota")
	} else if strings.Contains(lowerText, "dell") || strings.Contains(lowerText, "ዴል") {
		brand = "Dell"
	} else if strings.Contains(lowerText, "hp") || strings.Contains(lowerText, "ኤችፒ") {
		brand = "HP"
	}

	// Step 6: Build Localized Path and Description
	fullPathEN := fmt.Sprintf("%s › %s › %s", matchedBranch.RootNameEN, matchedBranch.SubcatNameEN, matchedBranch.LeafNameEN)
	pathLocalized := fullPathEN
	rootLocalized := matchedBranch.RootNameEN
	subcatLocalized := matchedBranch.SubcatNameEN
	leafLocalized := matchedBranch.LeafNameEN

	var suggestedDesc string

	switch normLocale {
	case "am":
		rootLocalized = matchedBranch.RootNameAM
		subcatLocalized = matchedBranch.SubcatNameAM
		leafLocalized = matchedBranch.LeafNameAM
		pathLocalized = fmt.Sprintf("%s › %s › %s", rootLocalized, subcatLocalized, leafLocalized)
		suggestedDesc = "በጣም በጥሩ ሁኔታ ላይ የሚገኝ እቃ። ተፈትሾ የተረጋገጠ እና ጥራቱ አስተማማኝ የሆነ። በኢትዮጵያ ውስጥ ደህንነቱ በተጠበቀ የግብይት ስርአት የሚሸጥ።"
	case "om":
		rootLocalized = matchedBranch.RootNameOM
		subcatLocalized = matchedBranch.SubcatNameOM
		leafLocalized = matchedBranch.LeafNameOM
		pathLocalized = fmt.Sprintf("%s › %s › %s", rootLocalized, subcatLocalized, leafLocalized)
		suggestedDesc = "Meeshaa haala gaarii fi qulqullina olaanaa qabu. Qoratamee kan mirkanaa'e. Gabaa HUNTIQ keessatti nagaan kan gurguramu."
	case "ti":
		rootLocalized = matchedBranch.RootNameTI
		subcatLocalized = matchedBranch.SubcatNameTI
		leafLocalized = matchedBranch.LeafNameTI
		pathLocalized = fmt.Sprintf("%s › %s › %s", rootLocalized, subcatLocalized, leafLocalized)
		suggestedDesc = "ኣብ ብሉፅ ኩነታት ዝርከብ ንብረት። ተፈቲሹ ዝተረጋገፀን ፅሬቱ ዝሓለወን። ኣብ ዕዳጋ HUNTIQ ብውሑስ ዝሽየጥ።"
	default:
		suggestedDesc = "Verified authentic item in excellent condition. Carefully inspected and available for safe, direct purchase and immediate handover."
	}

	confidence := 0.88
	if bestScore > 10 {
		confidence = 0.98
	} else if bestScore > 5 {
		confidence = 0.92
	}

	return &AutoCategorizeResult{
		DetectedTitle:        rawInput,
		PrimaryIntent:        intent,
		CategoryHierarchyID:  matchedBranch.FallbackHierarchyID,
		CategorySlug:         matchedBranch.LeafSlug,
		RootCategoryName:     rootLocalized,
		SubcategoryName:      subcatLocalized,
		LeafDestinationName:  leafLocalized,
		FullBranchPath:       fullPathEN,
		BranchPathLocalized:  pathLocalized,
		LegacyCategoryID:     matchedBranch.LegacyCategoryID,
		Condition:            condition,
		SuggestedPriceETB:    matchedBranch.TypicalPriceETB,
		PriceRangeMin:        matchedBranch.PriceMinETB,
		PriceRangeMax:        matchedBranch.PriceMaxETB,
		SuggestedDescription: suggestedDesc,
		DetectedLocale:       normLocale,
		Brand:                brand,
		Model:                model,
		Confidence:           confidence,
		HumanErrorCorrected:  humanErrorFound,
		OriginalTerm:         originalTerm,
		CorrectedTerm:        correctedTerm,
	}
}

func normalizeLocaleCode(loc string) string {
	switch strings.ToLower(strings.TrimSpace(loc)) {
	case "am", "amh":
		return "am"
	case "om", "or", "orm":
		return "om"
	case "ti", "tir":
		return "ti"
	case "en", "eng":
		return "en"
	default:
		return ""
	}
}

func extractModelSuffix(fullText string, triggerWords ...string) string {
	lower := strings.ToLower(fullText)
	for _, tw := range triggerWords {
		idx := strings.Index(lower, strings.ToLower(tw))
		if idx != -1 {
			remainder := strings.TrimSpace(fullText[idx+len(tw):])
			parts := strings.Fields(remainder)
			if len(parts) > 0 {
				if len(parts) >= 2 && (strings.EqualFold(parts[1], "pro") || strings.EqualFold(parts[1], "max") || strings.EqualFold(parts[1], "plus")) {
					return strings.Join(parts[:2], " ")
				}
				return parts[0]
			}
		}
	}
	return ""
}
