package main

// BrandDef defines a manufacturer or commercial trademark
type BrandDef struct {
	Name            string
	Slug            string
	CountryOfOrigin string
	CategorySlugs   []string
}

// ProductFamilyDef defines a family or series under a brand
type ProductFamilyDef struct {
	BrandSlug    string
	Name         string
	Slug         string
	SeriesName   string
	CategorySlug string
}

// ProductModelDef defines a real, verified product model with default specs
type ProductModelDef struct {
	BrandSlug    string
	FamilySlug   string
	CategorySlug string
	Name         string
	Slug         string
	ModelNumber  string
	ReleaseYear  int
	DefaultSpecs string // JSONB string
}

// GetVerifiedBrands returns verified manufacturers across marketplace sectors
func GetVerifiedBrands() []BrandDef {
	return []BrandDef{
		// Mobile & Electronics
		{Name: "Samsung", Slug: "samsung", CountryOfOrigin: "South Korea", CategorySlugs: []string{"electronics-mobile", "home-appliances", "tv-audio"}},
		{Name: "Apple", Slug: "apple", CountryOfOrigin: "United States", CategorySlugs: []string{"electronics-mobile", "computers-it"}},
		{Name: "Tecno", Slug: "tecno", CountryOfOrigin: "China", CategorySlugs: []string{"electronics-mobile"}},
		{Name: "Infinix", Slug: "infinix", CountryOfOrigin: "China", CategorySlugs: []string{"electronics-mobile"}},
		{Name: "Xiaomi", Slug: "xiaomi", CountryOfOrigin: "China", CategorySlugs: []string{"electronics-mobile", "home-appliances", "tv-audio"}},
		{Name: "Huawei", Slug: "huawei", CountryOfOrigin: "China", CategorySlugs: []string{"electronics-mobile", "telecom-networking"}},
		{Name: "Oppo", Slug: "oppo", CountryOfOrigin: "China", CategorySlugs: []string{"electronics-mobile"}},
		{Name: "Vivo", Slug: "vivo", CountryOfOrigin: "China", CategorySlugs: []string{"electronics-mobile"}},
		{Name: "Google", Slug: "google", CountryOfOrigin: "United States", CategorySlugs: []string{"electronics-mobile"}},
		{Name: "Realme", Slug: "realme", CountryOfOrigin: "China", CategorySlugs: []string{"electronics-mobile"}},
		{Name: "OnePlus", Slug: "oneplus", CountryOfOrigin: "China", CategorySlugs: []string{"electronics-mobile"}},
		{Name: "Nokia", Slug: "nokia", CountryOfOrigin: "Finland", CategorySlugs: []string{"electronics-mobile", "telecom-networking"}},
		{Name: "Sony", Slug: "sony", CountryOfOrigin: "Japan", CategorySlugs: []string{"tv-audio", "cameras-optics", "gaming-vr", "electronics-mobile"}},
		{Name: "LG", Slug: "lg", CountryOfOrigin: "South Korea", CategorySlugs: []string{"tv-audio", "home-appliances"}},
		{Name: "Hisense", Slug: "hisense", CountryOfOrigin: "China", CategorySlugs: []string{"tv-audio", "home-appliances"}},
		{Name: "TCL", Slug: "tcl", CountryOfOrigin: "China", CategorySlugs: []string{"tv-audio", "home-appliances"}},

		// Computers & IT
		{Name: "Dell", Slug: "dell", CountryOfOrigin: "United States", CategorySlugs: []string{"computers-it"}},
		{Name: "HP", Slug: "hp", CountryOfOrigin: "United States", CategorySlugs: []string{"computers-it"}},
		{Name: "Lenovo", Slug: "lenovo", CountryOfOrigin: "China", CategorySlugs: []string{"computers-it"}},
		{Name: "Asus", Slug: "asus", CountryOfOrigin: "Taiwan", CategorySlugs: []string{"computers-it", "gaming-vr"}},
		{Name: "Acer", Slug: "acer", CountryOfOrigin: "Taiwan", CategorySlugs: []string{"computers-it"}},
		{Name: "MSI", Slug: "msi", CountryOfOrigin: "Taiwan", CategorySlugs: []string{"computers-it", "gaming-vr"}},
		{Name: "Toshiba", Slug: "toshiba", CountryOfOrigin: "Japan", CategorySlugs: []string{"computers-it", "home-appliances"}},
		{Name: "Canon", Slug: "canon", CountryOfOrigin: "Japan", CategorySlugs: []string{"cameras-optics", "computers-it"}},
		{Name: "Epson", Slug: "epson", CountryOfOrigin: "Japan", CategorySlugs: []string{"computers-it"}},
		{Name: "Logitech", Slug: "logitech", CountryOfOrigin: "Switzerland", CategorySlugs: []string{"computers-it", "gaming-vr"}},

		// Networking & Telecom
		{Name: "TP-Link", Slug: "tp-link", CountryOfOrigin: "China", CategorySlugs: []string{"telecom-networking"}},
		{Name: "Cisco", Slug: "cisco", CountryOfOrigin: "United States", CategorySlugs: []string{"telecom-networking"}},
		{Name: "MikroTik", Slug: "mikrotik", CountryOfOrigin: "Latvia", CategorySlugs: []string{"telecom-networking"}},
		{Name: "Ubiquiti", Slug: "ubiquiti", CountryOfOrigin: "United States", CategorySlugs: []string{"telecom-networking"}},
		{Name: "D-Link", Slug: "d-link", CountryOfOrigin: "Taiwan", CategorySlugs: []string{"telecom-networking"}},
		{Name: "ZTE", Slug: "zte", CountryOfOrigin: "China", CategorySlugs: []string{"telecom-networking", "electronics-mobile"}},

		// Vehicles & Automobiles
		{Name: "Toyota", Slug: "toyota", CountryOfOrigin: "Japan", CategorySlugs: []string{"vehicles-cars", "trucks-commercial-transport", "vehicle-parts-accessories"}},
		{Name: "Isuzu", Slug: "isuzu", CountryOfOrigin: "Japan", CategorySlugs: []string{"trucks-commercial-transport", "vehicles-cars", "vehicle-parts-accessories"}},
		{Name: "Suzuki", Slug: "suzuki", CountryOfOrigin: "Japan", CategorySlugs: []string{"vehicles-cars", "motorcycles-bajaj"}},
		{Name: "Hyundai", Slug: "hyundai", CountryOfOrigin: "South Korea", CategorySlugs: []string{"vehicles-cars", "trucks-commercial-transport", "construction-machinery"}},
		{Name: "Kia", Slug: "kia", CountryOfOrigin: "South Korea", CategorySlugs: []string{"vehicles-cars"}},
		{Name: "Nissan", Slug: "nissan", CountryOfOrigin: "Japan", CategorySlugs: []string{"vehicles-cars", "trucks-commercial-transport"}},
		{Name: "Mitsubishi", Slug: "mitsubishi", CountryOfOrigin: "Japan", CategorySlugs: []string{"vehicles-cars", "trucks-commercial-transport"}},
		{Name: "Honda", Slug: "honda", CountryOfOrigin: "Japan", CategorySlugs: []string{"vehicles-cars", "motorcycles-bajaj"}},
		{Name: "Mercedes-Benz", Slug: "mercedes-benz", CountryOfOrigin: "Germany", CategorySlugs: []string{"vehicles-cars", "trucks-commercial-transport"}},
		{Name: "BMW", Slug: "bmw", CountryOfOrigin: "Germany", CategorySlugs: []string{"vehicles-cars", "motorcycles-bajaj"}},
		{Name: "Volkswagen", Slug: "volkswagen", CountryOfOrigin: "Germany", CategorySlugs: []string{"vehicles-cars"}},
		{Name: "Ford", Slug: "ford", CountryOfOrigin: "United States", CategorySlugs: []string{"vehicles-cars", "trucks-commercial-transport"}},
		{Name: "Sinotruk (HOWO)", Slug: "sinotruk", CountryOfOrigin: "China", CategorySlugs: []string{"trucks-commercial-transport"}},
		{Name: "FAW", Slug: "faw", CountryOfOrigin: "China", CategorySlugs: []string{"trucks-commercial-transport"}},
		{Name: "Bajaj Auto", Slug: "bajaj", CountryOfOrigin: "India", CategorySlugs: []string{"motorcycles-bajaj"}},
		{Name: "TVS Motor", Slug: "tvs", CountryOfOrigin: "India", CategorySlugs: []string{"motorcycles-bajaj"}},
		{Name: "Piaggio", Slug: "piaggio", CountryOfOrigin: "Italy", CategorySlugs: []string{"motorcycles-bajaj"}},
		{Name: "Yamaha", Slug: "yamaha", CountryOfOrigin: "Japan", CategorySlugs: []string{"motorcycles-bajaj", "music-instruments"}},

		// Agriculture & Heavy Machinery
		{Name: "Caterpillar (CAT)", Slug: "caterpillar", CountryOfOrigin: "United States", CategorySlugs: []string{"construction-machinery", "energy-solar-power"}},
		{Name: "Komatsu", Slug: "komatsu", CountryOfOrigin: "Japan", CategorySlugs: []string{"construction-machinery"}},
		{Name: "John Deere", Slug: "john-deere", CountryOfOrigin: "United States", CategorySlugs: []string{"agriculture-farming", "construction-machinery"}},
		{Name: "Massey Ferguson", Slug: "massey-ferguson", CountryOfOrigin: "United States", CategorySlugs: []string{"agriculture-farming"}},
		{Name: "New Holland", Slug: "new-holland", CountryOfOrigin: "Italy", CategorySlugs: []string{"agriculture-farming", "construction-machinery"}},
		{Name: "Kubota", Slug: "kubota", CountryOfOrigin: "Japan", CategorySlugs: []string{"agriculture-farming"}},
		{Name: "SANY", Slug: "sany", CountryOfOrigin: "China", CategorySlugs: []string{"construction-machinery"}},
		{Name: "XCMG", Slug: "xcmg", CountryOfOrigin: "China", CategorySlugs: []string{"construction-machinery"}},

		// Tools & Construction Materials
		{Name: "Bosch", Slug: "bosch", CountryOfOrigin: "Germany", CategorySlugs: []string{"tools-workshop", "home-appliances", "vehicle-parts-accessories"}},
		{Name: "Makita", Slug: "makita", CountryOfOrigin: "Japan", CategorySlugs: []string{"tools-workshop"}},
		{Name: "DeWalt", Slug: "dewalt", CountryOfOrigin: "United States", CategorySlugs: []string{"tools-workshop"}},
		{Name: "Total Tools", Slug: "total-tools", CountryOfOrigin: "China", CategorySlugs: []string{"tools-workshop"}},
		{Name: "Ingco", Slug: "ingco", CountryOfOrigin: "China", CategorySlugs: []string{"tools-workshop"}},
		{Name: "Dangote Cement", Slug: "dangote-cement", CountryOfOrigin: "Nigeria", CategorySlugs: []string{"construction-materials"}},
		{Name: "Mugher Cement", Slug: "mugher-cement", CountryOfOrigin: "Ethiopia", CategorySlugs: []string{"construction-materials"}},
		{Name: "Derba Cement", Slug: "derba-cement", CountryOfOrigin: "Ethiopia", CategorySlugs: []string{"construction-materials"}},
		{Name: "Habesha Cement", Slug: "habesha-cement", CountryOfOrigin: "Ethiopia", CategorySlugs: []string{"construction-materials"}},
		{Name: "Sika", Slug: "sika", CountryOfOrigin: "Switzerland", CategorySlugs: []string{"construction-materials"}},

		// Domestic Ethiopian Footwear & Apparel
		{Name: "Anbessa Shoe", Slug: "anbessa-shoe", CountryOfOrigin: "Ethiopia", CategorySlugs: []string{"shoes-bags"}},
		{Name: "Tikur Abay Shoe", Slug: "tikur-abay-shoe", CountryOfOrigin: "Ethiopia", CategorySlugs: []string{"shoes-bags"}},
		{Name: "Kangaroo Shoe", Slug: "kangaroo-shoe", CountryOfOrigin: "Ethiopia", CategorySlugs: []string{"shoes-bags"}},
		{Name: "Sheba Leather", Slug: "sheba-leather", CountryOfOrigin: "Ethiopia", CategorySlugs: []string{"shoes-bags", "clothing-apparel"}},
	}
}

// GetProductFamilies returns product lines under major brands
func GetProductFamilies() []ProductFamilyDef {
	return []ProductFamilyDef{
		// Samsung Mobile Families
		{BrandSlug: "samsung", Name: "Galaxy S Series (Flagship)", Slug: "samsung-galaxy-s", SeriesName: "Galaxy S Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "samsung", Name: "Galaxy A Series (Midrange & Essential)", Slug: "samsung-galaxy-a", SeriesName: "Galaxy A Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "samsung", Name: "Galaxy Z Series (Foldables)", Slug: "samsung-galaxy-z", SeriesName: "Galaxy Z Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "samsung", Name: "Galaxy M Series (Performance Battery)", Slug: "samsung-galaxy-m", SeriesName: "Galaxy M Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "samsung", Name: "Galaxy Note Series (Stylus)", Slug: "samsung-galaxy-note", SeriesName: "Galaxy Note Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "samsung", Name: "Galaxy Tab Series (Tablets)", Slug: "samsung-galaxy-tab", SeriesName: "Galaxy Tab Series", CategorySlug: "electronics-mobile"},

		// Apple Families
		{BrandSlug: "apple", Name: "iPhone Pro / Pro Max Series", Slug: "apple-iphone-pro", SeriesName: "iPhone Pro Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "apple", Name: "iPhone Standard / Plus Series", Slug: "apple-iphone-standard", SeriesName: "iPhone Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "apple", Name: "iPad Pro / Air / Mini", Slug: "apple-ipad-series", SeriesName: "iPad Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "apple", Name: "MacBook Pro / Air", Slug: "apple-macbook-series", SeriesName: "MacBook Series", CategorySlug: "computers-it"},

		// Tecno Families
		{BrandSlug: "tecno", Name: "Camon Series (Photography)", Slug: "tecno-camon-series", SeriesName: "Camon Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "tecno", Name: "Spark Series (Youth & Value)", Slug: "tecno-spark-series", SeriesName: "Spark Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "tecno", Name: "Phantom Series (Flagship Fold/Flip)", Slug: "tecno-phantom-series", SeriesName: "Phantom Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "tecno", Name: "Pova Series (Battery & Gaming)", Slug: "tecno-pova-series", SeriesName: "Pova Series", CategorySlug: "electronics-mobile"},

		// Infinix Families
		{BrandSlug: "infinix", Name: "Note Series (All-Round Fast Charge)", Slug: "infinix-note-series", SeriesName: "Note Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "infinix", Name: "Hot Series (Everyday Smartphone)", Slug: "infinix-hot-series", SeriesName: "Hot Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "infinix", Name: "Zero Series (Flagship Ultra)", Slug: "infinix-zero-series", SeriesName: "Zero Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "infinix", Name: "Smart Series (Entry Level)", Slug: "infinix-smart-series", SeriesName: "Smart Series", CategorySlug: "electronics-mobile"},

		// Xiaomi Families
		{BrandSlug: "xiaomi", Name: "Redmi Note Series", Slug: "xiaomi-redmi-note", SeriesName: "Redmi Note Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "xiaomi", Name: "Redmi Numbered / C Series", Slug: "xiaomi-redmi-standard", SeriesName: "Redmi Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "xiaomi", Name: "Poco X & F Performance Series", Slug: "xiaomi-poco-series", SeriesName: "Poco Series", CategorySlug: "electronics-mobile"},
		{BrandSlug: "xiaomi", Name: "Xiaomi Flagship Number Series", Slug: "xiaomi-flagship-series", SeriesName: "Xiaomi Series", CategorySlug: "electronics-mobile"},

		// Toyota Vehicle Families
		{BrandSlug: "toyota", Name: "Land Cruiser Family (300 / 200 / 70)", Slug: "toyota-land-cruiser", SeriesName: "Land Cruiser", CategorySlug: "vehicles-cars"},
		{BrandSlug: "toyota", Name: "Land Cruiser Prado (TX / VX)", Slug: "toyota-prado", SeriesName: "Prado", CategorySlug: "vehicles-cars"},
		{BrandSlug: "toyota", Name: "Hilux Pickup Family", Slug: "toyota-hilux", SeriesName: "Hilux", CategorySlug: "trucks-commercial-transport"},
		{BrandSlug: "toyota", Name: "Corolla / Corolla Cross", Slug: "toyota-corolla", SeriesName: "Corolla", CategorySlug: "vehicles-cars"},
		{BrandSlug: "toyota", Name: "Vitz / Yaris Compact Hatchback", Slug: "toyota-vitz-yaris", SeriesName: "Vitz / Yaris", CategorySlug: "vehicles-cars"},
		{BrandSlug: "toyota", Name: "RAV4 Compact SUV", Slug: "toyota-rav4", SeriesName: "RAV4", CategorySlug: "vehicles-cars"},
		{BrandSlug: "toyota", Name: "HiAce Van / Minibus", Slug: "toyota-hiace", SeriesName: "HiAce", CategorySlug: "trucks-commercial-transport"},
		{BrandSlug: "toyota", Name: "Coaster Public Bus", Slug: "toyota-coaster", SeriesName: "Coaster", CategorySlug: "trucks-commercial-transport"},

		// Isuzu Commercial Families
		{BrandSlug: "isuzu", Name: "N-Series Light Truck (NPR / NQR / NKR)", Slug: "isuzu-n-series", SeriesName: "N-Series (NPR)", CategorySlug: "trucks-commercial-transport"},
		{BrandSlug: "isuzu", Name: "F-Series Medium Truck (FSR / FTR / Forward)", Slug: "isuzu-f-series", SeriesName: "F-Series (Forward)", CategorySlug: "trucks-commercial-transport"},
		{BrandSlug: "isuzu", Name: "D-Max Pickup Series", Slug: "isuzu-d-max", SeriesName: "D-Max", CategorySlug: "trucks-commercial-transport"},

		// Suzuki Families
		{BrandSlug: "suzuki", Name: "Dzire Compact Sedan", Slug: "suzuki-dzire", SeriesName: "Dzire", CategorySlug: "vehicles-cars"},
		{BrandSlug: "suzuki", Name: "Swift Hatchback", Slug: "suzuki-swift", SeriesName: "Swift", CategorySlug: "vehicles-cars"},
		{BrandSlug: "suzuki", Name: "Jimny 4x4", Slug: "suzuki-jimny", SeriesName: "Jimny", CategorySlug: "vehicles-cars"},

		// Bajaj & TVS Three-Wheelers
		{BrandSlug: "bajaj", Name: "RE Compact Passenger Bajaj", Slug: "bajaj-re-passenger", SeriesName: "Bajaj RE", CategorySlug: "motorcycles-bajaj"},
		{BrandSlug: "bajaj", Name: "Maxima Cargo Three-Wheeler", Slug: "bajaj-maxima-cargo", SeriesName: "Maxima Cargo", CategorySlug: "motorcycles-bajaj"},
		{BrandSlug: "tvs", Name: "King Deluxe Passenger", Slug: "tvs-king-deluxe", SeriesName: "TVS King", CategorySlug: "motorcycles-bajaj"},

		// Heavy Machinery & Tractors
		{BrandSlug: "caterpillar", Name: "CAT 320 Hydraulic Excavator Series", Slug: "cat-320-excavator", SeriesName: "320 Series", CategorySlug: "construction-machinery"},
		{BrandSlug: "caterpillar", Name: "CAT 950 Wheel Loader Series", Slug: "cat-950-loader", SeriesName: "950 Series", CategorySlug: "construction-machinery"},
		{BrandSlug: "massey-ferguson", Name: "MF High Utility Farm Tractors (MF 385 / 375)", Slug: "mf-utility-tractors", SeriesName: "MF Utility", CategorySlug: "agriculture-farming"},
		{BrandSlug: "john-deere", Name: "5E Series Agricultural Tractors (5075E / 5050D)", Slug: "john-deere-5e", SeriesName: "5E Series", CategorySlug: "agriculture-farming"},
	}
}

// GetVerifiedModels returns authentic, verified product models with technical specs
func GetVerifiedModels() []ProductModelDef {
	return []ProductModelDef{
		// Samsung Galaxy S Series
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-s", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy S24 Ultra", Slug: "samsung-galaxy-s24-ultra", ModelNumber: "SM-S928B", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.8, "ram": "12GB", "storage": "512GB", "battery": 5000, "processor": "Snapdragon 8 Gen 3 for Galaxy", "os": "Android 14, One UI 6.1", "network": "5G", "sim": "Dual Physical SIM + eSIM"}`,
		},
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-s", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy S24+", Slug: "samsung-galaxy-s24-plus", ModelNumber: "SM-S926B", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.7, "ram": "12GB", "storage": "256GB", "battery": 4900, "processor": "Exynos 2400 / Snapdragon 8 Gen 3", "os": "Android 14", "network": "5G"}`,
		},
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-s", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy S24", Slug: "samsung-galaxy-s24", ModelNumber: "SM-S921B", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.2, "ram": "8GB", "storage": "128GB", "battery": 4000, "processor": "Exynos 2400", "os": "Android 14", "network": "5G"}`,
		},
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-s", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy S23 Ultra", Slug: "samsung-galaxy-s23-ultra", ModelNumber: "SM-S918B", ReleaseYear: 2023,
			DefaultSpecs: `{"screen_size": 6.8, "ram": "12GB", "storage": "256GB", "battery": 5000, "processor": "Snapdragon 8 Gen 2", "os": "Android 13 (upgradable)", "network": "5G"}`,
		},

		// Samsung Galaxy A Series
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-a", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy A55 5G", Slug: "samsung-galaxy-a55-5g", ModelNumber: "SM-A556B", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.6, "ram": "8GB", "storage": "128GB", "battery": 5000, "processor": "Exynos 1480", "network": "5G", "os": "Android 14"}`,
		},
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-a", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy A35 5G", Slug: "samsung-galaxy-a35-5g", ModelNumber: "SM-A356B", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.6, "ram": "8GB", "storage": "128GB", "battery": 5000, "processor": "Exynos 1380", "network": "5G", "os": "Android 14"}`,
		},
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-a", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy A25 5G", Slug: "samsung-galaxy-a25-5g", ModelNumber: "SM-A256B", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.5, "ram": "6GB", "storage": "128GB", "battery": 5000, "processor": "Exynos 1280", "network": "5G", "os": "Android 14"}`,
		},
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-a", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy A16 5G", Slug: "samsung-galaxy-a16-5g", ModelNumber: "SM-A166B", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.7, "ram": "6GB", "storage": "128GB", "battery": 5000, "processor": "Dimensity 6300", "network": "5G", "os": "Android 14"}`,
		},
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-a", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy A15", Slug: "samsung-galaxy-a15", ModelNumber: "SM-A155F", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.5, "ram": "4GB", "storage": "128GB", "battery": 5000, "processor": "Helio G99", "network": "4G LTE", "os": "Android 14"}`,
		},
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-a", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy A05s", Slug: "samsung-galaxy-a05s", ModelNumber: "SM-A057F", ReleaseYear: 2023,
			DefaultSpecs: `{"screen_size": 6.7, "ram": "4GB", "storage": "64GB", "battery": 5000, "processor": "Snapdragon 680 4G", "network": "4G LTE"}`,
		},

		// Samsung Foldables
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-z", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy Z Fold6", Slug: "samsung-galaxy-z-fold6", ModelNumber: "SM-F956B", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 7.6, "ram": "12GB", "storage": "512GB", "battery": 4400, "processor": "Snapdragon 8 Gen 3 for Galaxy", "network": "5G"}`,
		},
		{
			BrandSlug: "samsung", FamilySlug: "samsung-galaxy-z", CategorySlug: "electronics-mobile",
			Name: "Samsung Galaxy Z Flip6", Slug: "samsung-galaxy-z-flip6", ModelNumber: "SM-F741B", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.7, "ram": "12GB", "storage": "256GB", "battery": 4000, "processor": "Snapdragon 8 Gen 3 for Galaxy", "network": "5G"}`,
		},

		// Apple iPhone Models
		{
			BrandSlug: "apple", FamilySlug: "apple-iphone-pro", CategorySlug: "electronics-mobile",
			Name: "Apple iPhone 16 Pro Max", Slug: "apple-iphone-16-pro-max", ModelNumber: "A3296", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.9, "ram": "8GB", "storage": "256GB", "battery": 4685, "processor": "Apple A18 Pro", "os": "iOS 18", "network": "5G"}`,
		},
		{
			BrandSlug: "apple", FamilySlug: "apple-iphone-pro", CategorySlug: "electronics-mobile",
			Name: "Apple iPhone 16 Pro", Slug: "apple-iphone-16-pro", ModelNumber: "A3293", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.3, "ram": "8GB", "storage": "128GB", "battery": 3582, "processor": "Apple A18 Pro", "os": "iOS 18", "network": "5G"}`,
		},
		{
			BrandSlug: "apple", FamilySlug: "apple-iphone-standard", CategorySlug: "electronics-mobile",
			Name: "Apple iPhone 16", Slug: "apple-iphone-16", ModelNumber: "A3287", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.1, "ram": "8GB", "storage": "128GB", "battery": 3561, "processor": "Apple A18", "os": "iOS 18", "network": "5G"}`,
		},
		{
			BrandSlug: "apple", FamilySlug: "apple-iphone-pro", CategorySlug: "electronics-mobile",
			Name: "Apple iPhone 15 Pro Max", Slug: "apple-iphone-15-pro-max", ModelNumber: "A3106", ReleaseYear: 2023,
			DefaultSpecs: `{"screen_size": 6.7, "ram": "8GB", "storage": "256GB", "battery": 4441, "processor": "Apple A17 Pro", "os": "iOS 17", "network": "5G"}`,
		},
		{
			BrandSlug: "apple", FamilySlug: "apple-iphone-standard", CategorySlug: "electronics-mobile",
			Name: "Apple iPhone 15", Slug: "apple-iphone-15", ModelNumber: "A3090", ReleaseYear: 2023,
			DefaultSpecs: `{"screen_size": 6.1, "ram": "6GB", "storage": "128GB", "battery": 3349, "processor": "Apple A16 Bionic", "os": "iOS 17", "network": "5G"}`,
		},
		{
			BrandSlug: "apple", FamilySlug: "apple-iphone-standard", CategorySlug: "electronics-mobile",
			Name: "Apple iPhone 13", Slug: "apple-iphone-13", ModelNumber: "A2633", ReleaseYear: 2021,
			DefaultSpecs: `{"screen_size": 6.1, "ram": "4GB", "storage": "128GB", "battery": 3240, "processor": "Apple A15 Bionic", "os": "iOS 15", "network": "5G"}`,
		},
		{
			BrandSlug: "apple", FamilySlug: "apple-iphone-standard", CategorySlug: "electronics-mobile",
			Name: "Apple iPhone 11", Slug: "apple-iphone-11", ModelNumber: "A2221", ReleaseYear: 2019,
			DefaultSpecs: `{"screen_size": 6.1, "ram": "4GB", "storage": "64GB", "battery": 3110, "processor": "Apple A13 Bionic", "network": "4G LTE"}`,
		},

		// Tecno Models
		{
			BrandSlug: "tecno", FamilySlug: "tecno-camon-series", CategorySlug: "electronics-mobile",
			Name: "Tecno Camon 30 Pro 5G", Slug: "tecno-camon-30-pro-5g", ModelNumber: "CL8", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.78, "ram": "12GB", "storage": "512GB", "battery": 5000, "processor": "Dimensity 8200 Ultimate", "network": "5G"}`,
		},
		{
			BrandSlug: "tecno", FamilySlug: "tecno-spark-series", CategorySlug: "electronics-mobile",
			Name: "Tecno Spark 20 Pro+", Slug: "tecno-spark-20-pro-plus", ModelNumber: "KJ7", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.78, "ram": "8GB", "storage": "256GB", "battery": 5000, "processor": "Helio G99 Ultimate", "network": "4G LTE"}`,
		},
		{
			BrandSlug: "tecno", FamilySlug: "tecno-spark-series", CategorySlug: "electronics-mobile",
			Name: "Tecno Spark 20", Slug: "tecno-spark-20", ModelNumber: "KJ5", ReleaseYear: 2023,
			DefaultSpecs: `{"screen_size": 6.6, "ram": "8GB", "storage": "128GB", "battery": 5000, "processor": "Helio G85", "network": "4G LTE"}`,
		},

		// Infinix Models
		{
			BrandSlug: "infinix", FamilySlug: "infinix-note-series", CategorySlug: "electronics-mobile",
			Name: "Infinix Note 40 Pro+ 5G", Slug: "infinix-note-40-pro-plus-5g", ModelNumber: "X6851B", ReleaseYear: 2024,
			DefaultSpecs: `{"screen_size": 6.78, "ram": "12GB", "storage": "256GB", "battery": 4600, "processor": "Dimensity 7020", "network": "5G"}`,
		},
		{
			BrandSlug: "infinix", FamilySlug: "infinix-hot-series", CategorySlug: "electronics-mobile",
			Name: "Infinix Hot 40 Pro", Slug: "infinix-hot-40-pro", ModelNumber: "X6837", ReleaseYear: 2023,
			DefaultSpecs: `{"screen_size": 6.78, "ram": "8GB", "storage": "128GB", "battery": 5000, "processor": "Helio G99", "network": "4G LTE"}`,
		},

		// Toyota Vehicles
		{
			BrandSlug: "toyota", FamilySlug: "toyota-land-cruiser", CategorySlug: "vehicles-cars",
			Name: "Toyota Land Cruiser 300 VXR 3.5L Twin-Turbo", Slug: "toyota-land-cruiser-300-vxr", ModelNumber: "VJA300", ReleaseYear: 2024,
			DefaultSpecs: `{"fuel_type": "Benzine / Petrol", "engine_capacity": 3500, "transmission": "Automatic", "drivetrain": "4WD", "seats": 7}`,
		},
		{
			BrandSlug: "toyota", FamilySlug: "toyota-prado", CategorySlug: "vehicles-cars",
			Name: "Toyota Land Cruiser Prado TX 2.8L Diesel", Slug: "toyota-prado-tx-28-diesel", ModelNumber: "GDJ150", ReleaseYear: 2023,
			DefaultSpecs: `{"fuel_type": "Diesel", "engine_capacity": 2800, "transmission": "Automatic", "drivetrain": "4WD", "seats": 7}`,
		},
		{
			BrandSlug: "toyota", FamilySlug: "toyota-hilux", CategorySlug: "trucks-commercial-transport",
			Name: "Toyota Hilux Revo 2.8 D-4D 4x4 Double Cab", Slug: "toyota-hilux-revo-double-cab", ModelNumber: "GUN126", ReleaseYear: 2024,
			DefaultSpecs: `{"fuel_type": "Diesel", "engine_capacity": 2800, "transmission": "Automatic", "drivetrain": "4WD", "body_type": "Double Cabin Pickup"}`,
		},
		{
			BrandSlug: "toyota", FamilySlug: "toyota-corolla", CategorySlug: "vehicles-cars",
			Name: "Toyota Corolla Cross 1.8L Hybrid", Slug: "toyota-corolla-cross-hybrid", ModelNumber: "ZVG10", ReleaseYear: 2023,
			DefaultSpecs: `{"fuel_type": "Hybrid (Benzine/Electric)", "engine_capacity": 1800, "transmission": "Automatic (CVT)", "drivetrain": "FWD"}`,
		},
		{
			BrandSlug: "toyota", FamilySlug: "toyota-corolla", CategorySlug: "vehicles-cars",
			Name: "Toyota Corolla 1.6L Dual VVT-i (Executive)", Slug: "toyota-corolla-16-executive", ModelNumber: "ZRE171", ReleaseYear: 2020,
			DefaultSpecs: `{"fuel_type": "Benzine / Petrol", "engine_capacity": 1600, "transmission": "Automatic", "drivetrain": "FWD"}`,
		},
		{
			BrandSlug: "toyota", FamilySlug: "toyota-vitz-yaris", CategorySlug: "vehicles-cars",
			Name: "Toyota Vitz / Yaris 1.0L / 1.3L Hatchback", Slug: "toyota-vitz-hatchback", ModelNumber: "KSP130", ReleaseYear: 2018,
			DefaultSpecs: `{"fuel_type": "Benzine / Petrol", "engine_capacity": 1000, "transmission": "Automatic (CVT)", "drivetrain": "FWD"}`,
		},

		// Isuzu Commercial Trucks
		{
			BrandSlug: "isuzu", FamilySlug: "isuzu-n-series", CategorySlug: "trucks-commercial-transport",
			Name: "Isuzu NPR 75 4.5-Ton Dry Cargo Truck", Slug: "isuzu-npr-75-cargo", ModelNumber: "NPR75", ReleaseYear: 2023,
			DefaultSpecs: `{"fuel_type": "Diesel", "engine_capacity": 5193, "transmission": "Manual 6-Speed", "payload_capacity_kg": 4500, "axles": 2}`,
		},
		{
			BrandSlug: "isuzu", FamilySlug: "isuzu-f-series", CategorySlug: "trucks-commercial-transport",
			Name: "Isuzu FSR 33 8-Ton Medium Duty Truck", Slug: "isuzu-fsr-33-medium", ModelNumber: "FSR33", ReleaseYear: 2022,
			DefaultSpecs: `{"fuel_type": "Diesel", "engine_capacity": 8226, "transmission": "Manual 6-Speed", "payload_capacity_kg": 8000}`,
		},

		// Suzuki Models
		{
			BrandSlug: "suzuki", FamilySlug: "suzuki-dzire", CategorySlug: "vehicles-cars",
			Name: "Suzuki Dzire 1.2L GLX Sedan", Slug: "suzuki-dzire-12-glx", ModelNumber: "Dzire-GLX", ReleaseYear: 2023,
			DefaultSpecs: `{"fuel_type": "Benzine / Petrol", "engine_capacity": 1197, "transmission": "Automatic (AGS)", "drivetrain": "FWD"}`,
		},
		{
			BrandSlug: "suzuki", FamilySlug: "suzuki-swift", CategorySlug: "vehicles-cars",
			Name: "Suzuki Swift 1.2L DualJet Hatchback", Slug: "suzuki-swift-12-dualjet", ModelNumber: "Swift-GL", ReleaseYear: 2023,
			DefaultSpecs: `{"fuel_type": "Benzine / Petrol", "engine_capacity": 1197, "transmission": "Automatic (CVT)", "drivetrain": "FWD"}`,
		},

		// Bajaj RE Three-Wheeler
		{
			BrandSlug: "bajaj", FamilySlug: "bajaj-re-passenger", CategorySlug: "motorcycles-bajaj",
			Name: "Bajaj RE 4S 205cc Passenger Three-Wheeler", Slug: "bajaj-re-4s-passenger", ModelNumber: "RE-205", ReleaseYear: 2024,
			DefaultSpecs: `{"fuel_type": "Benzine / Petrol", "engine_capacity": 205, "transmission": "Manual 4-Forward + 1-Reverse", "seats": 4, "fuel_tank_liters": 8}`,
		},
		{
			BrandSlug: "tvs", FamilySlug: "tvs-king-deluxe", CategorySlug: "motorcycles-bajaj",
			Name: "TVS King Deluxe 200cc Passenger Bajaj", Slug: "tvs-king-deluxe-200", ModelNumber: "King-Deluxe", ReleaseYear: 2024,
			DefaultSpecs: `{"fuel_type": "Benzine / Petrol", "engine_capacity": 199.2, "transmission": "Manual 4-Speed", "seats": 4}`,
		},

		// Agricultural Tractors
		{
			BrandSlug: "massey-ferguson", FamilySlug: "mf-utility-tractors", CategorySlug: "agriculture-farming",
			Name: "Massey Ferguson MF 385 4WD 85HP Farm Tractor", Slug: "mf-385-4wd-tractor", ModelNumber: "MF-385", ReleaseYear: 2023,
			DefaultSpecs: `{"fuel_type": "Diesel", "horsepower": 85, "engine_capacity": 4100, "transmission": "Manual 8-Speed", "drivetrain": "4WD"}`,
		},
		{
			BrandSlug: "john-deere", FamilySlug: "john-deere-5e", CategorySlug: "agriculture-farming",
			Name: "John Deere 5075E 75HP Utility Tractor", Slug: "john-deere-5075e-tractor", ModelNumber: "5075E", ReleaseYear: 2023,
			DefaultSpecs: `{"fuel_type": "Diesel", "horsepower": 75, "transmission": "PowerReverser 12F/12R", "drivetrain": "4WD"}`,
		},

		// Heavy Construction Machinery
		{
			BrandSlug: "caterpillar", FamilySlug: "cat-320-excavator", CategorySlug: "construction-machinery",
			Name: "Caterpillar 320 GC Hydraulic Crawler Excavator", Slug: "cat-320-gc-excavator", ModelNumber: "320-GC", ReleaseYear: 2023,
			DefaultSpecs: `{"operating_weight_kg": 20500, "bucket_capacity_m3": 1.0, "net_power_hp": 146, "fuel_type": "Diesel"}`,
		},
	}
}
