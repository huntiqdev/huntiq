package main

// GetExpandedVerifiedModels returns hundreds of authentic verified models across phones, laptops, vehicles, and machinery
func GetExpandedVerifiedModels() []ProductModelDef {
	var models []ProductModelDef

	addM := func(brand, family, category, name, slug, modelNum string, year int, specs string) {
		models = append(models, ProductModelDef{
			BrandSlug:    brand,
			FamilySlug:   family,
			CategorySlug: category,
			Name:         name,
			Slug:         slug,
			ModelNumber:  modelNum,
			ReleaseYear:  year,
			DefaultSpecs: specs,
		})
	}

	// -------------------------------------------------------------------------
	// SAMSUNG GALAXY S SERIES (Flagship)
	// -------------------------------------------------------------------------
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S23+", "samsung-galaxy-s23-plus", "SM-S916B", 2023,
		`{"screen_size": 6.6, "ram": "8GB", "storage": "256GB", "battery": 4700, "processor": "Snapdragon 8 Gen 2", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S23", "samsung-galaxy-s23", "SM-S911B", 2023,
		`{"screen_size": 6.1, "ram": "8GB", "storage": "128GB", "battery": 3900, "processor": "Snapdragon 8 Gen 2", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S22 Ultra", "samsung-galaxy-s22-ultra", "SM-S908B", 2022,
		`{"screen_size": 6.8, "ram": "12GB", "storage": "256GB", "battery": 5000, "processor": "Snapdragon 8 Gen 1 / Exynos 2200", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S22+", "samsung-galaxy-s22-plus", "SM-S906B", 2022,
		`{"screen_size": 6.6, "ram": "8GB", "storage": "128GB", "battery": 4500, "processor": "Snapdragon 8 Gen 1", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S22", "samsung-galaxy-s22", "SM-S901B", 2022,
		`{"screen_size": 6.1, "ram": "8GB", "storage": "128GB", "battery": 3700, "processor": "Snapdragon 8 Gen 1", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S21 Ultra 5G", "samsung-galaxy-s21-ultra-5g", "SM-G998B", 2021,
		`{"screen_size": 6.8, "ram": "12GB", "storage": "256GB", "battery": 5000, "processor": "Exynos 2100", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S21+ 5G", "samsung-galaxy-s21-plus-5g", "SM-G996B", 2021,
		`{"screen_size": 6.7, "ram": "8GB", "storage": "128GB", "battery": 4800, "processor": "Exynos 2100", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S21 5G", "samsung-galaxy-s21-5g", "SM-G991B", 2021,
		`{"screen_size": 6.2, "ram": "8GB", "storage": "128GB", "battery": 4000, "processor": "Exynos 2100", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S20 Ultra 5G", "samsung-galaxy-s20-ultra-5g", "SM-G988B", 2020,
		`{"screen_size": 6.9, "ram": "12GB", "storage": "128GB", "battery": 5000, "processor": "Exynos 990", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S20 FE 5G", "samsung-galaxy-s20-fe-5g", "SM-G781B", 2020,
		`{"screen_size": 6.5, "ram": "6GB", "storage": "128GB", "battery": 4500, "processor": "Snapdragon 865", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S10+", "samsung-galaxy-s10-plus", "SM-G975F", 2019,
		`{"screen_size": 6.4, "ram": "8GB", "storage": "128GB", "battery": 4100, "processor": "Exynos 9820", "network": "4G LTE"}`)
	addM("samsung", "samsung-galaxy-s", "electronics-mobile", "Samsung Galaxy S10", "samsung-galaxy-s10", "SM-G973F", 2019,
		`{"screen_size": 6.1, "ram": "8GB", "storage": "128GB", "battery": 3400, "processor": "Exynos 9820", "network": "4G LTE"}`)

	// -------------------------------------------------------------------------
	// SAMSUNG GALAXY A SERIES (Popular in Ethiopia)
	// -------------------------------------------------------------------------
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A54 5G", "samsung-galaxy-a54-5g", "SM-A546B", 2023,
		`{"screen_size": 6.4, "ram": "8GB", "storage": "128GB", "battery": 5000, "processor": "Exynos 1380", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A53 5G", "samsung-galaxy-a53-5g", "SM-A536B", 2022,
		`{"screen_size": 6.5, "ram": "6GB", "storage": "128GB", "battery": 5000, "processor": "Exynos 1280", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A52s 5G", "samsung-galaxy-a52s-5g", "SM-A528B", 2021,
		`{"screen_size": 6.5, "ram": "6GB", "storage": "128GB", "battery": 4500, "processor": "Snapdragon 778G", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A51", "samsung-galaxy-a51", "SM-A515F", 2020,
		`{"screen_size": 6.5, "ram": "4GB", "storage": "128GB", "battery": 4000, "processor": "Exynos 9611", "network": "4G LTE"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A34 5G", "samsung-galaxy-a34-5g", "SM-A346B", 2023,
		`{"screen_size": 6.6, "ram": "6GB", "storage": "128GB", "battery": 5000, "processor": "Dimensity 1080", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A33 5G", "samsung-galaxy-a33-5g", "SM-A336B", 2022,
		`{"screen_size": 6.4, "ram": "6GB", "storage": "128GB", "battery": 5000, "processor": "Exynos 1280", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A24", "samsung-galaxy-a24", "SM-A245F", 2023,
		`{"screen_size": 6.5, "ram": "4GB", "storage": "128GB", "battery": 5000, "processor": "Helio G99", "network": "4G LTE"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A23", "samsung-galaxy-a23", "SM-A235F", 2022,
		`{"screen_size": 6.6, "ram": "4GB", "storage": "64GB", "battery": 5000, "processor": "Snapdragon 680", "network": "4G LTE"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A14 5G", "samsung-galaxy-a14-5g", "SM-A146B", 2023,
		`{"screen_size": 6.6, "ram": "4GB", "storage": "64GB", "battery": 5000, "processor": "Exynos 1330", "network": "5G"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A14 (4G)", "samsung-galaxy-a14-4g", "SM-A145F", 2023,
		`{"screen_size": 6.6, "ram": "4GB", "storage": "64GB", "battery": 5000, "processor": "Helio G80", "network": "4G LTE"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A13", "samsung-galaxy-a13", "SM-A135F", 2022,
		`{"screen_size": 6.6, "ram": "4GB", "storage": "64GB", "battery": 5000, "processor": "Exynos 850", "network": "4G LTE"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A12", "samsung-galaxy-a12", "SM-A125F", 2021,
		`{"screen_size": 6.5, "ram": "4GB", "storage": "64GB", "battery": 5000, "processor": "Helio P35", "network": "4G LTE"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A04s", "samsung-galaxy-a04s", "SM-A047F", 2022,
		`{"screen_size": 6.5, "ram": "3GB", "storage": "32GB", "battery": 5000, "processor": "Exynos 850", "network": "4G LTE"}`)
	addM("samsung", "samsung-galaxy-a", "electronics-mobile", "Samsung Galaxy A03 Core", "samsung-galaxy-a03-core", "SM-A032F", 2021,
		`{"screen_size": 6.5, "ram": "2GB", "storage": "32GB", "battery": 5000, "processor": "Unisoc SC9863A", "network": "4G LTE"}`)

	// -------------------------------------------------------------------------
	// SAMSUNG GALAXY NOTE SERIES
	// -------------------------------------------------------------------------
	addM("samsung", "samsung-galaxy-note", "electronics-mobile", "Samsung Galaxy Note 20 Ultra 5G", "samsung-galaxy-note-20-ultra-5g", "SM-N986B", 2020,
		`{"screen_size": 6.9, "ram": "12GB", "storage": "256GB", "battery": 4500, "processor": "Exynos 990", "network": "5G", "stylus": true}`)
	addM("samsung", "samsung-galaxy-note", "electronics-mobile", "Samsung Galaxy Note 20", "samsung-galaxy-note-20", "SM-N980F", 2020,
		`{"screen_size": 6.7, "ram": "8GB", "storage": "256GB", "battery": 4300, "processor": "Exynos 990", "network": "4G LTE", "stylus": true}`)
	addM("samsung", "samsung-galaxy-note", "electronics-mobile", "Samsung Galaxy Note 10+", "samsung-galaxy-note-10-plus", "SM-N975F", 2019,
		`{"screen_size": 6.8, "ram": "12GB", "storage": "256GB", "battery": 4300, "processor": "Exynos 9825", "network": "4G LTE", "stylus": true}`)

	// -------------------------------------------------------------------------
	// APPLE IPHONES
	// -------------------------------------------------------------------------
	addM("apple", "apple-iphone-standard", "electronics-mobile", "Apple iPhone 15 Plus", "apple-iphone-15-plus", "A3094", 2023,
		`{"screen_size": 6.7, "ram": "6GB", "storage": "128GB", "battery": 4383, "processor": "Apple A16 Bionic", "os": "iOS 17", "network": "5G"}`)
	addM("apple", "apple-iphone-pro", "electronics-mobile", "Apple iPhone 15 Pro", "apple-iphone-15-pro", "A3102", 2023,
		`{"screen_size": 6.1, "ram": "8GB", "storage": "128GB", "battery": 3274, "processor": "Apple A17 Pro", "os": "iOS 17", "network": "5G"}`)
	addM("apple", "apple-iphone-pro", "electronics-mobile", "Apple iPhone 14 Pro Max", "apple-iphone-14-pro-max", "A2894", 2022,
		`{"screen_size": 6.7, "ram": "6GB", "storage": "128GB", "battery": 4323, "processor": "Apple A16 Bionic", "os": "iOS 16", "network": "5G"}`)
	addM("apple", "apple-iphone-pro", "electronics-mobile", "Apple iPhone 14 Pro", "apple-iphone-14-pro", "A2890", 2022,
		`{"screen_size": 6.1, "ram": "6GB", "storage": "128GB", "battery": 3200, "processor": "Apple A16 Bionic", "os": "iOS 16", "network": "5G"}`)
	addM("apple", "apple-iphone-standard", "electronics-mobile", "Apple iPhone 14", "apple-iphone-14", "A2882", 2022,
		`{"screen_size": 6.1, "ram": "6GB", "storage": "128GB", "battery": 3279, "processor": "Apple A15 Bionic", "os": "iOS 16", "network": "5G"}`)
	addM("apple", "apple-iphone-pro", "electronics-mobile", "Apple iPhone 13 Pro Max", "apple-iphone-13-pro-max", "A2643", 2021,
		`{"screen_size": 6.7, "ram": "6GB", "storage": "128GB", "battery": 4352, "processor": "Apple A15 Bionic", "os": "iOS 15", "network": "5G"}`)
	addM("apple", "apple-iphone-pro", "electronics-mobile", "Apple iPhone 13 Pro", "apple-iphone-13-pro", "A2638", 2021,
		`{"screen_size": 6.1, "ram": "6GB", "storage": "128GB", "battery": 3095, "processor": "Apple A15 Bionic", "os": "iOS 15", "network": "5G"}`)
	addM("apple", "apple-iphone-standard", "electronics-mobile", "Apple iPhone 12 Pro Max", "apple-iphone-12-pro-max", "A2411", 2020,
		`{"screen_size": 6.7, "ram": "6GB", "storage": "128GB", "battery": 3687, "processor": "Apple A14 Bionic", "os": "iOS 14", "network": "5G"}`)
	addM("apple", "apple-iphone-standard", "electronics-mobile", "Apple iPhone 12", "apple-iphone-12", "A2403", 2020,
		`{"screen_size": 6.1, "ram": "4GB", "storage": "64GB", "battery": 2815, "processor": "Apple A14 Bionic", "os": "iOS 14", "network": "5G"}`)
	addM("apple", "apple-iphone-standard", "electronics-mobile", "Apple iPhone 11 Pro Max", "apple-iphone-11-pro-max", "A2218", 2019,
		`{"screen_size": 6.5, "ram": "4GB", "storage": "64GB", "battery": 3969, "processor": "Apple A13 Bionic", "os": "iOS 13", "network": "4G LTE"}`)
	addM("apple", "apple-iphone-standard", "electronics-mobile", "Apple iPhone XR", "apple-iphone-xr", "A2105", 2018,
		`{"screen_size": 6.1, "ram": "3GB", "storage": "64GB", "battery": 2942, "processor": "Apple A12 Bionic", "os": "iOS 12", "network": "4G LTE"}`)

	// -------------------------------------------------------------------------
	// TOYOTA VEHICLES (Ethiopia Market Workhorses)
	// -------------------------------------------------------------------------
	addM("toyota", "toyota-land-cruiser", "vehicles-cars", "Toyota Land Cruiser 200 VXR 4.5L V8 Diesel", "toyota-lc200-vxr-v8", "VDJ200", 2020,
		`{"fuel_type": "Diesel", "engine_capacity": 4500, "transmission": "Automatic 6-Speed", "drivetrain": "4WD", "seats": 8}`)
	addM("toyota", "toyota-land-cruiser", "vehicles-cars", "Toyota Land Cruiser 76 Hardtop 4.2L 1HZ Diesel", "toyota-lc76-hardtop-1hz", "HZJ76", 2023,
		`{"fuel_type": "Diesel", "engine_capacity": 4164, "transmission": "Manual 5-Speed", "drivetrain": "4WD", "body_type": "5-Door Hardtop"}`)
	addM("toyota", "toyota-land-cruiser", "vehicles-cars", "Toyota Land Cruiser 79 Single Cabin Pickup 4.2L Diesel", "toyota-lc79-single-cabin", "HZJ79", 2023,
		`{"fuel_type": "Diesel", "engine_capacity": 4164, "transmission": "Manual 5-Speed", "drivetrain": "4WD", "body_type": "Single Cabin Pickup"}`)
	addM("toyota", "toyota-prado", "vehicles-cars", "Toyota Land Cruiser Prado TX-L 2.8L Diesel", "toyota-prado-tx-l-28", "GDJ150L", 2023,
		`{"fuel_type": "Diesel", "engine_capacity": 2800, "transmission": "Automatic 6-Speed", "drivetrain": "4WD", "seats": 7}`)
	addM("toyota", "toyota-prado", "vehicles-cars", "Toyota Land Cruiser Prado 120 2.7L Petrol", "toyota-prado-120-27", "TRJ120", 2008,
		`{"fuel_type": "Benzine / Petrol", "engine_capacity": 2700, "transmission": "Automatic 4-Speed", "drivetrain": "4WD", "seats": 7}`)
	addM("toyota", "toyota-hilux", "trucks-commercial-transport", "Toyota Hilux Single Cabin 2.4 D-4D 4x4", "toyota-hilux-single-cab-24", "GUN125", 2023,
		`{"fuel_type": "Diesel", "engine_capacity": 2400, "transmission": "Manual 6-Speed", "drivetrain": "4WD", "payload_capacity_kg": 1000}`)
	addM("toyota", "toyota-vitz-yaris", "vehicles-cars", "Toyota Vitz 1.0L F-Package (XP130)", "toyota-vitz-10-xp130", "KSP130", 2017,
		`{"fuel_type": "Benzine / Petrol", "engine_capacity": 996, "transmission": "Automatic (CVT)", "drivetrain": "FWD", "body_type": "5-Door Hatchback"}`)
	addM("toyota", "toyota-vitz-yaris", "vehicles-cars", "Toyota Vitz 1.3L Jewela (XP90)", "toyota-vitz-13-xp90", "SCP90", 2010,
		`{"fuel_type": "Benzine / Petrol", "engine_capacity": 1296, "transmission": "Automatic (CVT)", "drivetrain": "FWD"}`)
	addM("toyota", "toyota-rav4", "vehicles-cars", "Toyota RAV4 2.5L Adventure AWD", "toyota-rav4-25-adventure", "XA50", 2022,
		`{"fuel_type": "Benzine / Petrol", "engine_capacity": 2500, "transmission": "Automatic 8-Speed", "drivetrain": "AWD"}`)
	addM("toyota", "toyota-hiace", "trucks-commercial-transport", "Toyota HiAce Commuter 15-Seater Minibus (2.5 D-4D)", "toyota-hiace-commuter-15", "KDH200", 2020,
		`{"fuel_type": "Diesel", "engine_capacity": 2500, "transmission": "Manual 5-Speed", "seats": 15}`)
	addM("toyota", "toyota-coaster", "trucks-commercial-transport", "Toyota Coaster 30-Seater Passenger Bus (4.2L 1HZ)", "toyota-coaster-30-seat", "HZB50", 2022,
		`{"fuel_type": "Diesel", "engine_capacity": 4164, "transmission": "Manual 5-Speed", "seats": 30}`)

	// -------------------------------------------------------------------------
	// ISUZU TRUCKS
	// -------------------------------------------------------------------------
	addM("isuzu", "isuzu-n-series", "trucks-commercial-transport", "Isuzu NPR 71 4-Ton Cargo Truck", "isuzu-npr-71-cargo", "NPR71", 2019,
		`{"fuel_type": "Diesel", "engine_capacity": 4570, "transmission": "Manual 5-Speed", "payload_capacity_kg": 4000}`)
	addM("isuzu", "isuzu-n-series", "trucks-commercial-transport", "Isuzu NQR 75 5.5-Ton Medium Cargo Truck", "isuzu-nqr-75-cargo", "NQR75", 2022,
		`{"fuel_type": "Diesel", "engine_capacity": 5193, "transmission": "Manual 6-Speed", "payload_capacity_kg": 5500}`)
	addM("isuzu", "isuzu-f-series", "trucks-commercial-transport", "Isuzu FTR 33 12-Ton Heavy Duty Truck", "isuzu-ftr-33-heavy", "FTR33", 2021,
		`{"fuel_type": "Diesel", "engine_capacity": 8226, "transmission": "Manual 6-Speed", "payload_capacity_kg": 12000}`)
	addM("isuzu", "isuzu-d-max", "trucks-commercial-transport", "Isuzu D-Max 3.0L 4x4 Double Cab V-Cross", "isuzu-dmax-30-vcross", "RG01", 2023,
		`{"fuel_type": "Diesel", "engine_capacity": 2999, "transmission": "Automatic 6-Speed", "drivetrain": "4WD"}`)

	// -------------------------------------------------------------------------
	// SUZUKI VEHICLES
	// -------------------------------------------------------------------------
	addM("suzuki", "suzuki-dzire", "vehicles-cars", "Suzuki Dzire 1.2L GL Manual", "suzuki-dzire-12-gl-manual", "Dzire-GL-M", 2022,
		`{"fuel_type": "Benzine / Petrol", "engine_capacity": 1197, "transmission": "Manual 5-Speed", "drivetrain": "FWD"}`)
	addM("suzuki", "suzuki-swift", "vehicles-cars", "Suzuki Swift 1.2L GLX DualJet", "suzuki-swift-12-glx", "Swift-GLX", 2023,
		`{"fuel_type": "Benzine / Petrol", "engine_capacity": 1197, "transmission": "Automatic (CVT)", "drivetrain": "FWD"}`)
	addM("suzuki", "suzuki-jimny", "vehicles-cars", "Suzuki Jimny 1.5L Sierra 4x4 3-Door", "suzuki-jimny-15-sierra", "JB74W", 2023,
		`{"fuel_type": "Benzine / Petrol", "engine_capacity": 1462, "transmission": "Manual 5-Speed", "drivetrain": "4WD"}`)

	// -------------------------------------------------------------------------
	// HEAVY DUMP TRUCKS (SinoTruck)
	// -------------------------------------------------------------------------
	addM("sinotruk", "cat-950-loader", "trucks-commercial-transport", "Sinotruk HOWO 371 6x4 Dump Truck Tipper", "sinotruk-howo-371-tipper", "ZZ3257", 2022,
		`{"fuel_type": "Diesel", "engine_capacity": 9726, "transmission": "Manual 10-Speed", "payload_capacity_kg": 25000, "horsepower": 371}`)
	addM("sinotruk", "cat-950-loader", "trucks-commercial-transport", "Sinotruk HOWO 420 6x4 Heavy Prime Mover Tractor Head", "sinotruk-howo-420-tractor", "ZZ4257", 2023,
		`{"fuel_type": "Diesel", "engine_capacity": 11596, "transmission": "Manual 12-Speed", "payload_capacity_kg": 40000, "horsepower": 420}`)

	return models
}

// GenerateExpandedOccupations adds more certified occupations across health, logistics, trades, and education
func GenerateExpandedOccupations() []OccupationDef {
	var list []OccupationDef

	add := func(slug, en, am, om, ti, cat, code string) {
		list = append(list, OccupationDef{
			Slug:         slug,
			TitleEN:      en,
			TitleAM:      am,
			TitleOM:      om,
			TitleTI:      ti,
			CategorySlug: cat,
			Code:         code,
		})
	}

	add("mobile-repair-technician", "Mobile Phone Hardware & Software Repair Technician", "የሞባይል ስልክ ጥገና ቴክኒሻን", "Teekiniishiyana Suphaa Bilbilaa", "ናይ ሞባይል ጽገና ቴክኒሽያን", "jobs-recruitment", "7422")
	add("computer-hardware-technician", "Computer Maintenance & Hardware Technician", "የኮምፒውተር ጥገና ቴክኒሻን", "Teekiniishiyana Suphaa Kompiitaraa", "ናይ ኮምፒዩተር ጽገና ቴክኒሽያን", "jobs-recruitment", "7422")
	add("warehouse-inventory-supervisor", "Warehouse Logistics & Inventory Supervisor", "የመጋዘን እና የዕቃ ዝርዝር ተቆጣጣሪ", "To'ataa Mankuusa fi Qabeenyaa", "ናይ መኽዘን ተቖጻጻሪ", "jobs-recruitment", "4321")
	add("forklift-operator", "Certified Industrial Forklift Operator", "የፎርክሊፍት ማሽነሪ ኦፕሬተር", "Ooppireetara Foorikiliiftii", "ናይ ፎርክሊፍት መዘወሪ", "jobs-recruitment", "8344")
	add("commercial-bank-manager", "Commercial Bank Branch Manager", "የባንክ ቅርንጫፍ ስራ አስኪያጅ", "Meneejara Damee Baankii", "ናይ ባንክ ጨንፈር ሓላፊ", "jobs-recruitment", "1211")
	add("hotel-executive-chef", "Executive Head Chef (National & International Cuisine)", "ዋና ሼፍ (የምግብ አዘጋጅ)", "Sheefii Ol-aanaa Nyaataa", "ዋና ሼፍ (መግቢ ኣዳላዊ)", "jobs-recruitment", "3434")
	add("food-beverage-waiter", "Professional Hotel & Restaurant Waiter / Waitress", "አስተናጋጅ (ወንድ/ሴት)", "Keeshummeessaa Nyaataa", "ኣሰናዳኢ መግቢ", "jobs-recruitment", "5131")
	add("barista-coffee-master", "Professional Coffee Barista & Roaster", "የቡና ባሪስታ ባለሙያ", "Ogeessa Buna Qopheessuu (Baariistaa)", "ናይ ቡን ባሪስታ", "jobs-recruitment", "5132")
	add("pharmacy-druggist-technician", "Licensed Pharmacy Druggist / Dispenser", "ድራጊስት (የመድሃኒት አደላጊ)", "Diraagiistii Faarmasii", "ናይ ፋርማሲ ድራጊስት", "jobs-recruitment", "3213")
	add("radiographer-xray-technician", "Medical Radiographer & X-Ray Technician", "የራጅ እና የራዲዮሎጂ ቴክኒሻን", "Teekiniishiyana Raajii fi X-Ray", "ናይ ራጅ ቴክኒሽያን", "jobs-recruitment", "3211")
	add("optometrist-eye-specialist", "Optometrist & Ophthalmic Lens Specialist", "የአይን ምርመራ እና የብርጽቆሽ ባለሙያ", "Ogeessa Yaala Ijaa", "ናይ ዓይኒ ሓኪም", "jobs-recruitment", "2267")
	add("tvet-vocational-instructor", "Technical TVET Vocational Skills Instructor", "የቴክኒክና ሙያ መምህር", "Barsiisaa Ogummaa TVET", "ናይ ቴክኒክን ሞያን መምህር", "jobs-recruitment", "2320")
	add("primary-school-teacher", "Primary School Teacher (Grades 1-8)", "የመጀመሪያ ደረጃ መምህር", "Barsiisaa Sadarkaa Tokkoffaa", "መምህር ቀዳማይ ደረጃ", "jobs-recruitment", "2341")
	add("kindergarten-early-childhood-educator", "Early Childhood & KG Educator", "የህፃናት መዋያ እና መዋለ ህፃናት መምህርት", "Barsiiftuu Daa'immanii", "ናይ ህጻናት መውዓሊ መምህር", "jobs-recruitment", "2342")
	add("security-guard-supervisor", "Corporate Security Operations Supervisor", "የጥበቃ ስራ አመራር እና ተቆጣጣሪ", "To'ataa Nageenya Dhaabbataa", "ናይ ሓለዋ ተቖጻጻሪ", "jobs-recruitment", "5414")
	add("cctv-surveillance-operator", "CCTV Surveillance & Control Room Operator", "የሲሲቲቪ ካሜራ ኦፕሬተር", "Ooppireetara Kaameeraa CCTV", "ናይ CCTV ካሜራ መከታተሊ", "jobs-recruitment", "5414")

	return list
}
