package main

func getExpandedDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 57. Root: Security Systems, CCTV & Access Control
		{
			ID:            "root-security-surveillance",
			Slug:          "security-systems-surveillance-cctv",
			Code:          "SEC_SURVEILLANCE",
			Level:         1,
			CanonicalName: "Security Systems, Surveillance, CCTV & Access Control",
			Definition:    "IP surveillance cameras, NVR/DVR recording units, biometric time-attendance terminals, magnetic door locks, and perimeter intruder alarms.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "ShieldCheck",
			Translations: map[string]string{
				"am": "የደህንነት ሲስተሞች፣ ሲሲቲቪ ካሜራዎች እና የበር መቆጣጠሪያዎች",
				"or": "Sirna Nageenyaa, Kaameraa CCTV fi Eegumsaa",
				"ti": "ናይ ድሕንነት ሲስተማት፣ ሲሲቲቪ ካሜራታትን ናይ ኣፍደገ ምቁጽጻርን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-cctv-cameras-nvrs",
					Slug:          "cctv-ip-cameras-recording-nvrs",
					Code:          "SEC_CCTV",
					Level:         2,
					CanonicalName: "CCTV IP Cameras & Network Video Recorders (ሲሲቲቪ)",
					Definition:    "Hikvision, Dahua 4MP/8MP night-vision dome and bullet cameras, and multi-channel PoE NVRs.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የሲሲቲቪ ካሜራዎች እና ሪኮርደሮች (NVR/DVR)",
						"or": "Kaameraalee CCTV fi Rikooridara",
						"ti": "ናይ ሲሲቲቪ ካሜራታትን ሪኮርደራትን",
					},
					Aliases: map[string][]string{
						"en": {"cctv", "security camera", "ip camera", "nvr", "dvr", "hikvision"},
						"am": {"ሲሲቲቪ", "የደህንነት ካሜራ", "ካሜራ", "የበር ካሜራ"},
						"or": {"kaameraa cctv", "kaameraa nageenyaa"},
						"ti": {"ሲሲቲቪ", "ናይ ድሕንነት ካሜራ"},
					},
				},
			},
		},

		// 58. Root: Fire Safety, Firefighting & Disaster Prevention
		{
			ID:            "root-fire-safety",
			Slug:          "fire-safety-firefighting-disaster-prevention",
			Code:          "FIRE_SAFETY",
			Level:         1,
			CanonicalName: "Fire Safety, Firefighting & Disaster Prevention",
			Definition:    "Dry powder and CO2 fire extinguishers, optical smoke detectors, fire alarm sounders, and commercial fire hose reels.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Flame",
			Translations: map[string]string{
				"am": "የእሳት አደጋ መከላከያ እና የማስጠንቀቂያ እቃዎች",
				"or": "Meeshaalee Balaa Ibidaa Ittisuu",
				"ti": "ናይ ሓዊ ሓደጋ መከላኸልን መጥፍእን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-fire-extinguishers",
					Slug:          "portable-fire-extinguishers",
					Code:          "FIRE_EXTINGUISH",
					Level:         2,
					CanonicalName: "Portable Fire Extinguishers (CO2 & Dry Powder)",
					Definition:    "ABC dry chemical powder and carbon dioxide portable fire canisters (2kg, 6kg, 9kg, 50kg trolley).",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ተንቀሳቃሽ የእሳት ማጥፊያ ሲሊንደሮች",
						"or": "Meeshaalee Ibida Dhaamsan",
						"ti": "ተንቀሳቐስቲ ናይ ሓዊ መጥፍኢ ሲሊንደራት",
					},
				},
			},
		},

		// 59. Root: Occupational Safety, PPE & Workwear
		{
			ID:            "root-occupational-safety-ppe",
			Slug:          "occupational-safety-ppe-workwear",
			Code:          "SAFETY_PPE",
			Level:         1,
			CanonicalName: "Occupational Safety, Personal Protective Equipment & Workwear",
			Definition:    "Steel-toe industrial safety shoes, high-visibility reflective vests, hard hats, welding safety goggles, and chemical respirator masks.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "HardHat",
			Translations: map[string]string{
				"am": "የስራ ላይ ደህንነት እና የጥበቃ አልባሳት (ሴፍቲ ጫማ፣ ቆብ)",
				"or": "Meeshaalee Nageenya Hojii fi Eegumsaa (PPE)",
				"ti": "ናይ ስራሕ ድሕንነትን መከላኸሊ ክዳውንትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-safety-boots-helmets",
					Slug:          "steel-toe-safety-boots-hard-hats",
					Code:          "PPE_FOOT_HEAD",
					Level:         2,
					CanonicalName: "Steel-Toe Safety Boots & Industrial Hard Hats (ሴፍቲ ጫማ)",
					Definition:    "Puncture-resistant steel-toe work boots (Caterpillar, Delta Plus) and ABS dielectric construction safety helmets.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የብረት አፍ ሴፍቲ ጫማዎች እና የግንባታ ቆቦች",
						"or": "Kophee Nageenyaa fi Qobfoo Ijaarsaa",
						"ti": "ናይ ብረት ኣፍ ሴፍቲ ሳእኒታትን ቆብዕን",
					},
					Aliases: map[string][]string{
						"en": {"safety shoes", "safety boots", "steel toe", "hard hat", "ppe", "helmet"},
						"am": {"ሴፍቲ ጫማ", "የስራ ጫማ", "የግንባታ ቆብ", "ሄልሜት"},
						"or": {"kophee nageenyaa", "qobfoo"},
						"ti": {"ሴፍቲ ሳእኒ", "ቆብዕ"},
					},
				},
			},
		},

		// 60. Root: Water Treatment, Borehole Pumps & Irrigation Systems
		{
			ID:            "root-water-irrigation",
			Slug:          "water-treatment-borehole-pumps-irrigation",
			Code:          "WATER_TECH",
			Level:         1,
			CanonicalName: "Water Treatment, Borehole Pumps & Irrigation Systems",
			Definition:    "Deep-well submersible borehole pumps (Franklin, Pedrollo), rotational water drip irrigation kits, and commercial reverse osmosis water filtration.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Droplet",
			Translations: map[string]string{
				"am": "የውሃ ማጣሪያ፣ ጥልቅ የውሃ ጉድጓድ ፓምፖች እና የመስኖ እቃዎች",
				"or": "Paampii Bishaanii, Qulqulleessituu fi Jallisii",
				"ti": "ናይ ማይ መጽረዪ፣ ናይ ዓሚቕ ዒላ ፓምፓታትን መስኖን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-submersible-borehole-pumps",
					Slug:          "submersible-borehole-water-pumps",
					Code:          "PUMP_SUBMERSIBLE",
					Level:         2,
					CanonicalName: "Submersible Borehole & Well Water Pumps (ሰመሪሴብል ፓምፕ)",
					Definition:    "Multi-stage high-head stainless steel submersible pumps for agricultural deep-wells and urban multi-story building water delivery.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ሰመሪሴብል የውሃ መሳቢያ ፓምፖች (ፔድሮሎ)",
						"or": "Paampii Bishaanii Lafa Jalaa (Samarsiibilii)",
						"ti": "ሰመርሲብል ናይ ማይ ፓምፓታት",
					},
					Aliases: map[string][]string{
						"en": {"submersible pump", "water pump", "borehole pump", "pedrollo", "well pump"},
						"am": {"ሰመሪሴብል ፓምፕ", "የውሃ ፓምፕ", "የጉድጓድ ፓምፕ", "ፔድሮሎ"},
						"or": {"paampii bishaanii", "samarsiibilii"},
						"ti": {"ሰመርሲብል ፓምፕ", "ናይ ማይ ፓምፕ"},
					},
				},
			},
		},

		// 61. Root: Office Supplies, Stationery & Consumables
		{
			ID:            "root-office-supplies",
			Slug:          "office-supplies-stationery-paper-consumables",
			Code:          "OFFICE_SUPPLIES",
			Level:         1,
			CanonicalName: "Office Supplies, Stationery & Printing Consumables",
			Definition:    "A4 copy paper reams, laser printer toner cartridges, arch lever box files, binding machines, and heavy-duty paper shredders.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "FileText",
			Translations: map[string]string{
				"am": "የቢሮ እቃዎች፣ ስቴሽነሪ እና የህትመት ቶነሮች",
				"or": "Meeshaalee Biiroo fi Waraqaa",
				"ti": "ናይ ቢሮ ንዋያት፣ ስቴሽነሪን ቶነርን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-copy-paper-reams",
					Slug:          "a4-copy-paper-cartons-reams",
					Code:          "PAPER_REAMS",
					Level:         2,
					CanonicalName: "A4 White Copy Paper Cartons & Reams (A4 ወረቀት)",
					Definition:    "80gsm and 75gsm multi-purpose high-speed laser photocopying paper reams (Double A, PaperOne, Chamex).",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ኤ 4 የፎቶ ኮፒ ወረቀቶች በካርቶን",
						"or": "Waraqaa Fotokoppii A4",
						"ti": "A4 ናይ ፎቶ ኮፒ ወረቐት",
					},
				},
			},
		},

		// 62. Root: Retail POS Equipment, Cash Registers & Store Fixtures
		{
			ID:            "root-retail-pos",
			Slug:          "retail-pos-cash-registers-store-fixtures",
			Code:          "RETAIL_POS",
			Level:         1,
			CanonicalName: "Retail Point of Sale (POS), Cash Registers & Store Fixtures",
			Definition:    "Fiscal electronic cash registers (ERCA certified F-POS), barcode laser scanners, thermal receipt printers, and supermarket shelving gondolas.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "CreditCard",
			Translations: map[string]string{
				"am": "የችርቻሮ ሽያጭ መመዝገቢያ (ካሽ ሬጂስተር) እና ባርኮድ ስካነር",
				"or": "Meeshaalee Galii Daldalaa (Kaash Reejistara)",
				"ti": "ናይ ንግዲ መመዝገቢ (ካሽ ሬጂስተር)ን ባርኮድ ስካነርን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-fiscal-cash-registers",
					Slug:          "fiscal-electronic-cash-registers",
					Code:          "POS_REGISTERS",
					Level:         2,
					CanonicalName: "ERCA-Certified Fiscal Electronic Cash Registers (ካሽ ሬጂስተር)",
					Definition:    "Government-accredited electronic sales register machines (Datecs, Daisy, Tremol) with GPRS telecommunication modules.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የገቢዎች እውቅና ያላቸው የሽያጭ መመዝገቢያ ማሽኖች (ካሽ ሬጂስተር)",
						"or": "Maashina Galii Galmeessu (Kaash Reejistara)",
						"ti": "ናይ ግብሪ ፍቓድ ዘለዎም ናይ መሸጢ ማሽናት",
					},
					Aliases: map[string][]string{
						"en": {"cash register", "pos machine", "fiscal machine", "datecs", "receipt printer"},
						"am": {"ካሽ ሬጂስተር", "የሽያጭ ማሽን", "የገቢዎች ማሽን", "ቴርማል ፕሪንተር"},
						"or": {"kaash reejistara"},
						"ti": {"ካሽ ሬጂስተር"},
					},
				},
			},
		},

		// 63. Root: Commercial Kitchen, Bakery & Food Processing Equipment
		{
			ID:            "root-commercial-kitchen",
			Slug:          "commercial-kitchen-bakery-food-equipment",
			Code:          "REST_EQUIPMENT",
			Level:         1,
			CanonicalName: "Commercial Kitchen, Bakery & Catering Equipment",
			Definition:    "Commercial deck bakery ovens, spiral dough kneaders, stainless steel food prep counters, shawarma broilers, and industrial deep fryers.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "UtensilsCrossed",
			Translations: map[string]string{
				"am": "የዳቦ መጋገሪያ፣ የሆቴል እና የሬስቶራንት ማብሰያ እቃዎች",
				"or": "Meeshaalee Kushiinaa Hoteelaa fi Daabboo",
				"ti": "ናይ ዳቦ መጋገርያ፣ ናይ ሆቴልን ሬስቶራንትን መብሰሊ ንዋያት",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-bakery-deck-ovens",
					Slug:          "bakery-deck-ovens-dough-mixers",
					Code:          "BAKE_EQUIPMENT",
					Level:         2,
					CanonicalName: "Commercial Bakery Deck Ovens & Spiral Dough Mixers",
					Definition:    "Electric and gas stone-hearth bread and pastry deck ovens, planetary cake mixers, and bread slicers.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የዳቦ መጋገሪያ ኦቨኖች እና የሊጥ ማቡኪያ ማሽኖች",
						"or": "Ovenii Daabboo fi Maashina Boocaa",
						"ti": "ናይ ዳቦ መጋገርያ እቶናትን መሐወሲታትን",
					},
				},
			},
		},

		// 64. Root: Packaging, Shipping & Industrial Containers
		{
			ID:            "root-packaging-containers",
			Slug:          "packaging-shipping-industrial-containers",
			Code:          "PACK_CONTAINERS",
			Level:         1,
			CanonicalName: "Packaging, Shipping & Industrial Containers",
			Definition:    "Corrugated cardboard cartons, polypropylene (PP) woven agricultural grain bags, pallet stretch wrap, and IBC liquid storage containers.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Box",
			Translations: map[string]string{
				"am": "የማሸጊያ ካርቶኖች፣ የሰብል ከረጢቶች እና ኮንቴነሮች",
				"or": "Meeshaalee Qophii Saamsuu fi Saandooqa",
				"ti": "መዐሸጊ ካርቶናት፣ ናይ እኽሊ ማዳበርያታትን መትሓዚታትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-pp-grain-sacks",
					Slug:          "polypropylene-woven-grain-sacks",
					Code:          "PACK_SACKS",
					Level:         2,
					CanonicalName: "Polypropylene (PP) Woven Grain Sacks (የእህል ማዳበሪያ ከረጢት)",
					Definition:    "50kg and 100kg reinforced woven plastic sacks for teff, wheat, coffee, and fertilizer transport.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የ100 ኪሎ እና 50 ኪሎ የፕላስቲክ የእህል ከረጢቶች (ማዳበሪያ)",
						"or": "Qaraaxii Midhaanii fi Xaa'oo (Kiiloo 100 fi 50)",
						"ti": "ናይ 100 ኪሎን 50 ኪሎን ናይ እኽሊ ማዳበርያታት",
					},
					Aliases: map[string][]string{
						"en": {"grain sacks", "pp bags", "madaberia", "jute bags", "packaging bags"},
						"am": {"ማዳበሪያ", "የእህል ከረጢት", "ጆንያ", "የቡና ከረጢት"},
						"or": {"qaraaxii", "keeshaa"},
						"ti": {"ማዳበርያ", "ናይ እኽሊ ከረጢት"},
					},
				},
			},
		},

		// 65. Root: Professional Legal, Financial & Corporate Advisory Services
		{
			ID:            "root-legal-financial-services",
			Slug:          "professional-legal-financial-advisory-services",
			Code:          "SERV_LEGAL_FIN",
			Level:         1,
			CanonicalName: "Professional Legal, Financial, Audit & Advisory Services",
			Definition:    "Chartered certified accounting and auditing, Ethiopian corporate tax filing, commercial law contract preparation, and business feasibility studies.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Scale",
			Translations: map[string]string{
				"am": "የህግ፣ የኦዲት፣ የግብር እና የሂሳብ አያያዝ አገልግሎቶች",
				"or": "Tajaajila Seeraa, Oodiitii fi Herregaa",
				"ti": "ናይ ሕጊ፣ ኦዲት፣ ግብርን ሒሳብን ኣገልግሎት",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-accounting-tax-audit-services",
					Slug:          "chartered-accounting-tax-audit-services",
					Code:          "SERV_AUDIT",
					Level:         2,
					CanonicalName: "Chartered Accounting, Tax Filing & External Audit Services",
					Definition:    "Certified audit reports, monthly VAT/TOT declarations to Ministry of Revenues, and statutory business financial statements.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የውጭ ኦዲት፣ የታክስ እና የቫት ማሳወቂያ አገልግሎት",
						"or": "Tajaajila Oodiitii Alaa fi Taaksii",
						"ti": "ናይ ደገ ኦዲትን ግብርን ኣገልግሎት",
					},
				},
			},
		},

		// 66. Root: Information Technology, Cloud & Software Services
		{
			ID:            "root-it-cloud-services",
			Slug:          "information-technology-cloud-software-services",
			Code:          "SERV_IT_CLOUD",
			Level:         1,
			CanonicalName: "Information Technology, Cloud & Software Services",
			Definition:    "Custom web and mobile application programming, cloud hosting infrastructure, cybersecurity penetration testing, and IT systems integration.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Terminal",
			Translations: map[string]string{
				"am": "የሶፍትዌር ልማት፣ የድረ-ገጽ እና የክላውድ ቴክኖሎጂ አገልግሎቶች",
				"or": "Tajaajila Misooma Sooftiweerii fi Teekinoolojii",
				"ti": "ናይ ሶፍትዌር ምዕባለ፣ ዌብሳይትን ቴክኖሎጂን ኣገልግሎት",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-web-mobile-app-dev",
					Slug:          "custom-web-mobile-app-development",
					Code:          "SERV_DEV",
					Level:         2,
					CanonicalName: "Custom Web Application & Mobile App Development",
					Definition:    "Full-stack React/Next.js and Flutter/Android development with payment gateway integration (Telebirr, CBE Birr, Chapa).",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የሞባይል አፕሊኬሽን እና የድረ-ገጽ ማበልጸግ አገልግሎት",
						"or": "Misooma Weebsaayitii fi Moobaayilii Appii",
						"ti": "ናይ ሞባይል ኣፕሊኬሽንን ዌብሳይትን ምዕባለ",
					},
				},
			},
		},

		// 67. Root: Architecture, Engineering & Civil Drafting Services
		{
			ID:            "root-architecture-engineering-services",
			Slug:          "architecture-engineering-civil-drafting-services",
			Code:          "SERV_ENG_ARCH",
			Level:         1,
			CanonicalName: "Architecture, Engineering & Civil Drafting Services",
			Definition:    "Architectural building blueprints, structural engineering calculations, 3D visualization, soil geotechnical testing, and topographical surveying.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Compass",
			Translations: map[string]string{
				"am": "አርክቴክቸር፣ የህንፃ ምህንድስና ፕላን እና የቅየሳ አገልግሎቶች",
				"or": "Dizaayinii Gamoo, Injinariingii fi Kaartaa",
				"ti": "ኣርክቴክቸር፣ ናይ ህንጻ ፕላንን መሃንድስነትን ኣገልግሎት",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-architectural-plans-3d",
					Slug:          "architectural-building-plans-3d-renders",
					Code:          "SERV_ARCH_PLANS",
					Level:         2,
					CanonicalName: "Architectural Building Plans, 3D Renders & Permits",
					Definition:    "Municipality-approved architectural drawings, structural engineering detail drawings, and electrical/sanitary building schematics.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የቤት ፕላን ዝግጅት እና የህንፃ ዲዛይን ስራዎች",
						"or": "Pilaanii Ijaarsa Gamoo fi Dizaayinii",
						"ti": "ናይ ቤት ፕላንን ዲዛይን ህንጻን",
					},
				},
			},
		},

		// 68. Root: Building Construction & Craft Trade Services
		{
			ID:            "root-construction-craft-services",
			Slug:          "building-construction-craft-trade-services",
			Code:          "SERV_CONSTRUCTION",
			Level:         1,
			CanonicalName: "Building Construction, Contracting & Craft Trades",
			Definition:    "General construction contracting, certified plumbing installation, licensed electrical wiring, masonry, and gypsum interior finishing.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Hammer",
			Translations: map[string]string{
				"am": "የግንባታ ኮንትራት፣ የኤሌክትሪክ እና የቧንቧ ሙያ አገልግሎቶች",
				"or": "Tajaajila Kontiraata Ijaarsaa fi Ogeeyyii",
				"ti": "ናይ ህንጻ ኮንትራት፣ ኤሌክትሪክን ሻምብቆን ሞያዊ ኣገልግሎት",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-electrical-installation-trade",
					Slug:          "certified-electrical-wiring-installation",
					Code:          "TRADE_ELECTRICAL",
					Level:         2,
					CanonicalName: "Certified Electrical Wiring & Generator Installation",
					Definition:    "Conduit pulling, circuit breaker panel wiring, generator automatic transfer switch wiring, and earthing installation.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የኤሌክትሪክ መስመር ዝርጋታ እና የጀነሬተር ማገናኘት ስራ",
						"or": "Sarara Elektirikii fi Jenereetara Hidhuu",
						"ti": "ናይ ኤሌክትሪክ መስመር ምዝርጋሕን ጀነሬተር ምትእስሳርን",
					},
				},
			},
		},

		// 69. Root: Events, Staging & Audio-Visual Production Services
		{
			ID:            "root-event-production-services",
			Slug:          "events-staging-audiovisual-production-services",
			Code:          "SERV_EVENTS",
			Level:         1,
			CanonicalName: "Events, Wedding Staging & Audio-Visual Production",
			Definition:    "Wedding hall floral decorations, wedding tent rentals, professional audio concert line-arrays, event photography, and master of ceremonies.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "PartyPopper",
			Translations: map[string]string{
				"am": "የሰርግ፣ የዝግጅት ማስዋብ እና የድምጽ ኪራይ አገልግሎቶች",
				"or": "Qophii Sirna Cidhaa fi Miidhagina Sagantaa",
				"ti": "ናይ መርዓ፣ ናይ መደባት ምድላውን ናይ ድምጺ ክራይን ኣገልግሎት",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-wedding-decor-tents",
					Slug:          "wedding-ceremony-decor-tent-rentals",
					Code:          "SERV_WEDDING_DECOR",
					Level:         2,
					CanonicalName: "Wedding Ceremony Decor, Stage Lighting & Tent Rentals",
					Definition:    "Luxury bridal stage floral backdrops, Chiavari chairs, red carpets, and waterproof banquet canopy tents.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የሰርግ ማድመቂያዎች፣ ድንኳኖች እና የመድረክ ጌጦች",
						"or": "Dunkana Cidhaa fi Miidhagina Waltajjii",
						"ti": "ናይ መርዓ ድንኳናት፣ መንበርን ናይ መድረኽ ስልማትን",
					},
				},
			},
		},

		// 70. Root: General Classifieds, Community & Donations
		{
			ID:            "root-general-classifieds",
			Slug:          "general-classifieds-community-donations",
			Code:          "COMMUNITY_CLASSIFIEDS",
			Level:         1,
			CanonicalName: "General Classifieds, Community & Item Donations",
			Definition:    "Community announcements, lost & found document notices (National ID, Kebele ID, passport, driver license), and charitable non-profit item donations.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Megaphone",
			Translations: map[string]string{
				"am": "አጠቃላይ ማስታወቂያዎች፣ የጠፉ እቃዎች እና ልገሳ",
				"or": "Beeksisa Waliigalaa, Meeshaa Bade fi Kennaa",
				"ti": "ሓፈሻዊ ምልክታታት፣ ዝጠፍኡ ኣቑሑትን ልገሳን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-lost-and-found-documents",
					Slug:          "lost-found-official-documents",
					Code:          "COMM_LOST_FOUND",
					Level:         2,
					CanonicalName: "Lost & Found Official Credentials & Documents (የጠፉ ሰነዶች)",
					Definition:    "Community lost and found notices for Fayda digital IDs, passports, vehicle libre titles, and academic degree certificates.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የጠፉ መታወቂያዎች፣ ፓስፖርት እና ሰነዶች",
						"or": "Waraqaa Eenyummaa fi Paaspoortii Bade",
						"ti": "ዝጠፍኡ መታወቒታት፣ ፓስፖርትን ሰነዳትን",
					},
				},
			},
		},
	}
}
