package main

import "fmt"

// DomainCatalogDef defines a domain with structured subcategories and leaf items
type DomainCatalogDef struct {
	RootSlug      string
	Subcategories []SubcatCatalogDef
}

// SubcatCatalogDef defines a subcategory with terminal leaves
type SubcatCatalogDef struct {
	Slug   string
	NameEN string
	NameAM string
	NameOM string
	NameTI string
	Items  []LeafItemDef
}

// LeafItemDef defines a leaf entity
type LeafItemDef struct {
	Slug   string
	NameEN string
	NameAM string
	NameOM string
	NameTI string
}

// GetUniversalDomainLeafDatabase produces comprehensive, authentic leaf categories across all 60 domains
func GetUniversalDomainLeafDatabase() []DomainCatalogDef {
	var domains []DomainCatalogDef

	// Helper function
	addD := func(rootSlug string, subs []SubcatCatalogDef) {
		domains = append(domains, DomainCatalogDef{
			RootSlug:      rootSlug,
			Subcategories: subs,
		})
	}

	// Helper to create a subcategory with items
	makeSub := func(slug, en, am, om, ti string, items [][4]string) SubcatCatalogDef {
		var leafList []LeafItemDef
		for _, it := range items {
			leafList = append(leafList, LeafItemDef{
				Slug:   it[0],
				NameEN: it[1],
				NameAM: it[2],
				NameOM: it[3],
				NameTI: fmt.Sprintf("ናይ %s", it[1]),
			})
		}
		return SubcatCatalogDef{
			Slug:   slug,
			NameEN: en,
			NameAM: am,
			NameOM: om,
			NameTI: ti,
			Items:  leafList,
		}
	}

	// -------------------------------------------------------------------------
	// 03. Telecommunications and Networking
	// -------------------------------------------------------------------------
	addD("telecom-networking", []SubcatCatalogDef{
		makeSub("fiber-broadband-cabling", "Fiber Optic Equipment and Fusion Splicers", "የፋይበር ኦፕቲክ እቃዎች እና ስፕላይሰሮች", "Meeshaalee Faayibaraa", "ናይ ፋይበር ኦፕቲክ መሳርሒታት", [][4]string{
			{"optical-fiber-fusion-splicers", "Automatic Core-Alignment Fiber Fusion Splicing Machines", "የፋይበር ማገናኛ ስፕላይሰሮች", "Maashina Splaayisara Faayibaraa"},
			{"optical-time-domain-reflectometers", "Handheld OTDR Fiber Optic Fault Locators", "የፋይበር ብልሽት መለያ መሳሪያዎች (OTDR)", "Meeshaa Qorannoo OTDR Faayibaraa"},
			{"armored-outdoor-fiber-cables", "2 to 96 Core Outdoor Armored Fiber Optic Trunk Cables", "የውጭ የፋይበር ኦፕቲክ ኬብሎች", "Keebilii Faayibaraa Ala Dirree"},
			{"fiber-optic-patch-panels-odf", "Rack Mount Optical Distribution Frames (ODF 24/48 Port)", "የፋይበር ማከፋፈያ ሳጥኖች (ODF)", "Saanduqa ODF Faayibaraa"},
		}),
		makeSub("enterprise-wireless-ptp", "Point-to-Point Wireless Links and Antennas", "የረጅም ርቀት ገመድ አልባ ኔትወርክ (PTP)", "Netwoorkii PTP Fageenya Dheeraa", "ናይ ነዊሕ ርሕቐት PTP ኔትወርክ", [][4]string{
			{"ptp-wireless-bridges-airfiber", "5GHz to 60GHz Point-to-Point Wireless Gigabit Bridges", "የረጅም ርቀት ጊጋቢት ገመድ አልባ ድልድዮች", "Bridge Wayaalessii PTP"},
			{"high-gain-parabolic-antennas", "High-Gain Parabolic Grid & Dish Antennas (30dBi+)", "ፓራቦሊክ የኔትወርክ አንቴናዎች", "Anteenaa Paaraaboolikii"},
			{"outdoor-enterprise-access-points", "IP67 Weatherproof Outdoor Wi-Fi 6 Access Points", "የውጭ የዋይፋይ አክሰስ ፖይንቶች", "Access Point Wi-Fi 6 Ala"},
		}),
	})

	// -------------------------------------------------------------------------
	// 07. Home Appliances (Deep categories)
	// -------------------------------------------------------------------------
	addD("home-appliances", []SubcatCatalogDef{
		makeSub("small-kitchen-electric-appliances", "Blenders, Juicers, Kettles and Toasters", "ትናንሽ የወጥ ቤት ኤሌክትሮኒክስ እቃዎች", "Meeshaalee Kushiinaa Xixinnoo", "ንኣሽቱ ናይ ክሽነ ኤሌክትሮኒክስ", [][4]string{
			{"commercial-heavy-duty-smoothie-blenders", "Commercial High-Power 2200W Ice-Crushing Smoothie Blenders", "ጠንካራ የንግድ የፈሳሽ መፍጫዎች (ብሌንደር)", "Bileendara Humna Guddaa 2200W"},
			{"slow-masticating-cold-press-juicers", "Cold-Press Slow Masticating Fruit and Vegetable Juicers", "የፍራፍሬ ጭማቂ መጭመቂያ ማሽኖች (ጁሰር)", "Maashina Cuunfaa Kuduraa fi Fuduraa"},
			{"rapid-boil-cordless-electric-kettles", "2-Liter Stainless Steel Fast-Boil Electric Kettles", "የኤሌክትሪክ የውሃ ማንቆርቆሪያዎች (ኬትል)", "Mankoorkooraa Elektiriikii 2L"},
			{"non-stick-electric-sandwich-makers", "Non-Stick 4-Slice Sandwich Makers & Panini Grills", "የሳንድዊች መጋገሪያዎች (ቶስተር)", "Toostara Saandwiichii"},
			{"rotisserie-electric-countertop-ovens", "45L to 60L Convection Rotisserie Countertop Electric Ovens", "የኤሌክትሪክ መጋገሪያ ምድጃዎች", "Oovonii Minjaalaa Elektiriikii"},
		}),
		makeSub("vacuum-cleaners-floor-care", "Vacuum Cleaners, Steam Mops and Floor Scrubbers", "ቫኪዩም ክሊነሮች እና የወለል ማጽጃዎች", "Qulqulleessituu Lafaa fi Faakiyuumii", "ቫኪዩም ክሊነርን መጽረዪን", [][4]string{
			{"heavy-wet-dry-vacuum-cleaners", "Commercial Wet and Dry Industrial Drum Vacuum Cleaners (30L-80L)", "የደረቅ እና የፈሳሽ ቫኪዩም ክሊነሮች", "Faakiyuumii Qulqullinaa Guddaa"},
			{"bagless-cyclonic-home-vacuums", "Bagless Cyclonic Multi-Surface Home Vacuum Cleaners", "የቤት ውስጥ አቧራ መጣያ ቫኪዩም ክሊነሮች", "Faakiyuumii Mana Keessaa"},
			{"cordless-stick-vacuum-cleaners", "Lightweight Rechargeable Cordless Stick Vacuum Cleaners", "የባትሪ ተንቀሳቃሽ ቫኪዩም ክሊነሮች", "Faakiyuumii Baatirii Qabu Salphaa"},
			{"high-temperature-steam-mops", "Sanitizing High-Temperature Floor Steam Cleaners", "የእንፋሎት የወለል እና የሰድር ማጽጃዎች", "Qulqulleessituu Hurkaa Ho'aa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 12. Jewelry, Watches and Luxury Goods
	// -------------------------------------------------------------------------
	addD("jewelry-luxury", []SubcatCatalogDef{
		makeSub("ethiopian-gemstones-minerals", "Ethiopian Welo Opals, Emeralds and Sapphires", "የኢትዮጵያ የከበሩ ድንጋዮች (ኦፓል)", "Dhagaa Gatii Guddaa Itoophiyaa (Ooppaalii)", "ናይ ኢትዮጵያ ክቡራት ኣእማን (ኦፓል)", [][4]string{
			{"rough-uncut-welo-fire-opals", "Rough Uncut Natural Ethiopian Welo Hydrophane Fire Opals", "ያልተፈተገ የደላንታ እና የወሎ ኦፓል", "Ooppaalii Uumamaa Wolloo"},
			{"polished-cabochon-faceted-opals", "Faceted & Oval Cut Polished Precious Welo Opals", "የተወለወለ እና የተጠረበ ኦፓል", "Ooppaalii Qallame fi Qulqullaa'e"},
			{"ethiopian-shakiso-emeralds", "Natural Green Ethiopian Shakiso Emerald Gemstones", "የሻኪሶ ተፈጥሯዊ ኤመራልድ (ዘመሩድ)", "Dhagaa Eemeraaldii Shakkisoo"},
			{"natural-blue-sapphires", "Unheated Blue Sapphire Crystals (Tigray / Oromia)", "ተፈጥሯዊ ብሉ ሳፋየር ድንጋይ", "Safayira Cuquliisa Uumamaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 14. Health, Medical Equipment and Accessibility
	// -------------------------------------------------------------------------
	addD("health-medical", []SubcatCatalogDef{
		makeSub("clinical-hospital-furniture", "Hospital Beds, Examination Couches and IV Stands", "የሆስፒታል እና የክሊኒክ አልጋዎች", "Siree Hospitaalaa fi Meeshaa Yaalaa", "ናይ ሆስፒታል ዓራታትን መሳርሒታትን", [][4]string{
			{"multi-function-electric-hospital-beds", "3-Crank & 5-Function Electric ICU Medical Hospital Beds", "ባለ ብዙ እንቅስቃሴ የሆስፒታል አልጋዎች", "Siree Hospitaalaa Elektiriikii"},
			{"hydraulic-patient-examination-tables", "Adjustable Hydraulic Doctor Examination Couches", "የሃኪም መመርመሪያ ጠረጴዛዎች", "Minjaala Qorannoo Hakiimaa"},
			{"stainless-steel-medical-iv-stands", "5-Caster Stainless Steel Adjustable IV Infusion Stands", "የሴረም መስቀያ የብረት ቋሚዎች (IV Stand)", "Fannoo Seeramii Sibiila Qulqulluu"},
			{"emergency-first-responder-stretchers", "Foldable Ambulance Emergency Patient Transport Stretchers", "የአምቡላንስ እና የአደጋ ጊዜ በረኖች (Stretcher)", "Siree Baattuu Dhukkubsataa Saffisaa"},
		}),
		makeSub("oxygen-respiratory-therapy", "Oxygen Concentrators, Cylinders and Nebulizers", "የኦክስጅን ማሽኖች እና ሲሊንደሮች", "Maashina Oksijiinii fi Siliandara", "ናይ ኦክስጅን ማሽንን ሲሊንደርን", [][4]string{
			{"5l-10l-medical-oxygen-concentrators", "5-Liter and 10-Liter Home Medical Oxygen Concentrators", "የቤት ውስጥ የኦክስጅን ኮንሰንትሬተሮች", "Maashina Oksijiinii Manaa 5L-10L"},
			{"high-pressure-medical-oxygen-cylinders", "40-Liter High-Pressure Refillable Medical Oxygen Gas Tanks", "ትላልቅ የኦክስጅን ሲሊንደሮች (ታንኮች)", "Siliandara Oksijiinii Yaalaa 40L"},
			{"compressor-mesh-nebulizers", "Quiet Compressor & Ultrasonic Mesh Asthma Nebulizer Inhalers", "የአስም ማከሚያ ኔቡላይዘር ማሽኖች", "Maashina Niibuulaayizara Asmii"},
		}),
	})

	// -------------------------------------------------------------------------
	// 20. Freelance, Remote Work and Contract Gigs
	// -------------------------------------------------------------------------
	addD("freelance-gigs", []SubcatCatalogDef{
		makeSub("translation-transcription-gigs", "Amharic, Oromo, Tigrinya and English Language Services", "የቋንቋ ትርጉም እና ጽሁፍ አገልግሎት", "Tajaajila Turjumaana Afaanii", "ናይ ቋንቋ ትርጉምን ጽሑፍን ኣገልግሎት", [][4]string{
			{"amharic-english-certified-document-translation", "Certified Amharic to English Legal Document Translation", "የአማርኛ-እንግሊዝኛ የሰነድ ትርጉም", "Turjumaana Afaan Amaaraa fi Ingilizii"},
			{"afan-oromo-english-translation-transcription", "Professional Afaan Oromo to English Translation & Audio Transcription", "የኦሮምኛ ቋንቋ የድምጽ እና የሰነድ ትርጉም", "Turjumaana Afaan Oromoo fi Sagalee"},
			{"tigrinya-english-translation-localization", "Certified Tigrinya to English Document & Website Localization", "የትግርኛ ቋንቋ ሰነድ እና ድረ ገጽ ትርጉም", "Turjumaana Afaan Tigriffaa fi Weebsaayitii"},
			{"geez-fidel-optical-character-recognition", "Ge'ez Fidel OCR Data Entry, Scanning and Archiving", "የግዕዝ እና የፊደል ጽሁፍ መተየብ እና ዲጂታይዜሽን", "Galmee fi Barreeffama Qubee Gi'iizii"},
		}),
		makeSub("remote-software-digital-gigs", "Remote Software, Web, Design & Data Entry Projects", "የርቀት የሶፍትዌር እና የዲጂታል ስራዎች", "Hojii Fageenyaa Sooftiweerii fi Daataa", "ናይ ርሑቕ ሶፍትዌርን ዳታን ስራሕቲ", [][4]string{
			{"react-nextjs-frontend-contract-development", "Freelance React / Next.js / Tailwind Frontend Web Development", "የፍሪላንስ ዌብሳይት ግንባታ (React/Next.js)", "Oomisha Weebsaayitii Fageenyaa (React)"},
			{"golang-python-backend-api-contracts", "Golang, Node.js & Python REST/gRPC API Microservices Development", "የጀርባ ሲስተም እና የኤፒአይ ግንባታ (Go/Python)", "Hojii API fi Siistamii Duubaa (Go/Python)"},
			{"social-media-content-creation-flyers", "Graphic Flyers, Social Media Banners & Motion Graphics Gigs", "የማህበራዊ ሚዲያ ፍላየር እና ፖስተር ዲዛይን", "Dizaayinii Biriishara fi Weebsaayitii"},
			{"data-entry-spreadsheet-cleanup", "Bulk Excel Data Entry, Virtual Assistance & PDF Formatting", "የኤክሴል ዳታ መሙላት እና የቨርቹዋል ረዳት ስራዎች", "Galmee Daataa fi Gargaarsa Faayilaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 23. Business Services and Consulting
	// -------------------------------------------------------------------------
	addD("business-services", []SubcatCatalogDef{
		makeSub("investment-feasibility-advisory", "Bank Feasibility Studies and Project Proposals", "የኢንቨስትመንት ጥናት እና የብድር ፕሮፖዛል", "Qorannoo Faayidaa Daldalaa fi Baankii", "ናይ ባንክ ብድርን ኢንቨስትመንትን ፕሮፖዛል", [][4]string{
			{"bank-loan-feasibility-study-documents", "Comprehensive Bank Loan Feasibility Studies (CBE, Awash, Dashen)", "ለባንክ ብድር የሚቀርቡ የቢዝነስ ፕላኖች", "Pirooppoozaala Liqii Baankii"},
			{"eic-investment-permit-consulting", "Ethiopian Investment Commission (EIC) Licensing & Duty-Free Clearance", "የኢንቨስትመንት ፈቃድ እና ከቀረጥ ነፃ ማማከር", "Gorsa Eeyyama Invastimantii fi Gumurukii"},
			{"environmental-social-impact-assessments", "Certified Environmental & Social Impact Assessment (ESIA) Reports", "የአካባቢ እና ማህበራዊ ተፅዕኖ ግምገማ (ESIA)", "Qorannoo Dhiibbaa Naannoo (ESIA)"},
		}),
	})

	// -------------------------------------------------------------------------
	// 38. Motorcycles, Scooters and Three-Wheelers
	// -------------------------------------------------------------------------
	addD("motorcycles-bajaj", []SubcatCatalogDef{
		makeSub("bajaj-motorcycle-spare-parts", "Bajaj RE and Motorcycle Engine Spare Parts", "የባጃጅ እና የሞተር መለዋወጫዎች", "Qaama Jijjiirraa Baajaajii", "ናይ ባጃጅ መለዋወጢ ክፍልታት", [][4]string{
			{"bajaj-re-cylinder-piston-block-kits", "Bajaj RE 205cc Cylinder Block, Piston & Ring Assemblies", "የባጃጅ ሲሊንደር እና ፒስተን ኪት", "Kilaandara fi Piistoonii Baajaajii RE"},
			{"bajaj-clutch-plates-cables", "Bajaj Friction Clutch Plates, Clutch Hubs & Control Cables", "የባጃጅ ክላች ፕሌት እና ኬብሎች", "Kilaachii fi Wadaroo Baajaajii"},
			{"three-wheeler-tires-tubes-400-8", "4.00-8 & 4.50-10 Bajaj Nylon Heavy Ply Tires & Tubes", "የባጃጅ ጎማዎች እና ካሜሪዎች (4.00-8)", "Gommaa fi Kaameeraa Baajaajii"},
			{"bajaj-carburetors-fuel-pumps", "Bajaj RE 4S Float Carburetors and Electronic Fuel Pumps", "የባጃጅ ካርቡሬተሮች", "Kaarbireetara Baajaajii"},
		}),
	})

	// -------------------------------------------------------------------------
	// 45. Livestock, Poultry and Animal Husbandry
	// -------------------------------------------------------------------------
	addD("livestock-poultry", []SubcatCatalogDef{
		makeSub("indigenous-livestock-breeds", "Ethiopian Native Cattle, Sheep and Goat Breeds", "የሀገር በቀል ከብቶች፣ በጎች እና ፍየሎች", "Gosa Beellada Aadaa Itoophiyaa", "ናይ ዓዲ ከብቲ፣ ኣባጊዕን ኣጣልን", [][4]string{
			{"fogera-horro-native-cows", "Fogera and Horro Indigenous High-Adaptability Cattle", "የፎገራ እና የሆሮ ዝርያ ላሞች", "Saawwaa Foogaraa fi Horroo"},
			{"bonga-menz-indigenous-sheep", "Bonga and Menz Highland Meat and Wool Sheep", "የቦንጋ እና የመንዝ ዝርያ በጎች", "Hoolaa Boongaa fi Manzii"},
			{"afar-somali-dryland-goats", "Afar and Somali Highly Drought-Resistant Dairy Goats", "የአፋር እና የሶማሌ ፍየሎች", "Re'ee Afaarii fi Somaalee"},
		}),
		makeSub("modern-beekeeping-apiculture", "Modern Beehives, Honey Extractors and Bee Colonies", "ዘመናዊ የንብ ቀፎዎች እና የንብ ማነብ እቃዎች", "Gaagura Dammaa Ammayyaa fi Kanniisa", "ዘመናዊ ናይ ንህቢ ቆፎን መዓርን", [][4]string{
			{"langstroth-frame-wooden-beehives", "3-Box Frame Langstroth Modern Timber Beehives", "ላንግስትሮዝ ዘመናዊ የእንጨት የንብ ቀፎዎች", "Gaagura Dammaa Laangistirooz"},
			{"stainless-steel-manual-honey-extractors", "4-Frame & 6-Frame Stainless Steel Radial Honey Extractors", "የማር ማጣሪያ እና ማወጫ ማሽኖች", "Maashina Damma Cuunfu (Extractor)"},
			{"beekeeping-smokers-protective-suits", "Protective Beekeeping Veil Suits, Gloves and Stainless Smokers", "የንብ አልባሳት፣ ጓንቶች እና የጢስ ማጨሻዎች", "Uffata fi Qodaa Eegumsa Kanniisaa"},
			{"pure-queen-bee-colonies-swarms", "Fertile Queen Bees and Established Honeybee Colonies (Apis mellifera)", "የሰሩ ንቦች እና የንግስት ንብ መንጋዎች", "Kanniisa Dammaa fi Giiftii Kanniisaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 49. Beverages, Coffee and Tea
	// -------------------------------------------------------------------------
	addD("coffee-tea-beverages", []SubcatCatalogDef{
		makeSub("traditional-ethiopian-beverages", "Honey Wine (Tej), Tella and Local Drinks", "ጠጅ፣ ጠላ እና ባህላዊ መጠጦች", "Daadhii fi Farsoo Aadaa", "ሜስን ስዋን ባህላዊ መስተታትን", [][4]string{
			{"pure-fermented-honey-wine-tej", "Fermented Pure Raw Honey Wine (Tej) with Gesho", "ንጹህ የማር ጠጅ (ከጌሾ ጋር የተዘጋጀ)", "Daadhii Dammaa Qulqulluu (Tej)"},
			{"traditional-barley-tella-supply", "Traditional Fermented Barley & Gesho Dark Tella", "የገብስ እና የጌሾ ድፍድፍ ጠላ", "Farsoo Garbuu fi Geeshoo"},
			{"dry-gesho-leaves-twigs", "Sun-Dried Highland Gesho Leaves and Crushed Twigs (Rhamnus prinoides)", "የደረቀ የጌሾ ቅጠል እና እንጨት", "Baala fi Muka Geeshoo Goge"},
		}),
	})

	return domains
}
