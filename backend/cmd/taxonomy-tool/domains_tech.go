package main

func getTechDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 7. Root: Consumer Electronics & Mobile Devices
		{
			ID:            "root-electronics-mobile",
			Slug:          "electronics-mobile",
			Code:          "ELEC_MOBILE",
			Level:         1,
			CanonicalName: "Consumer Electronics & Mobile Devices",
			Definition:    "Handheld mobile telephony equipment, smartphones, feature phones, cellular tablets, and wearable electronics.",
			InclusionRules: []string{"Smartphones", "Feature phones / keypad phones", "Tablet computers", "Smartwatches", "Mobile phone accessories"},
			ExclusionRules: []string{"Laptop and desktop computers (see Computers)", "Phone repair and screen replacement services (see Services)"},
			Examples:       []string{"iPhone 15 Pro", "Samsung Galaxy S24", "Tecno Spark", "Itel keypad phone", "iPad Air", "Tempered glass screen protector"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Smartphone",
			Translations: map[string]string{
				"am": "ኤሌክትሮኒክስ እና የሞባይል ስልኮች",
				"or": "Elektirooniksii fi Meeshaalee Moobaayilii",
				"ti": "ኤሌክትሮኒክስን ናይ ሞባይል ስልኪታትን",
			},
			Aliases: map[string][]string{
				"en": {"electronics", "mobile", "phones", "smartphones", "cellphones", "cell phone"},
				"am": {"ኤሌክትሮኒክስ", "ስልክ", "ሞባይል", "ስማርትፎን", "ተንቀሳቃሽ ስልክ", "ስልኮች"},
				"or": {"elektirooniksii", "bilbila", "moobaayilii", "ismaartfooni"},
				"ti": {"ኤሌክትሮኒክስ", "ስልኪ", "ሞባይል", "ስማርትፎን"},
			},
			Children: []TaxonomyNode{
				// Super Family 7.1: Mobile Handsets & Smartphones
				{
					ID:            "sf-smartphones-handsets",
					Slug:          "mobile-phones",
					Code:          "ELEC_HANDSETS",
					Level:         2,
					CanonicalName: "Mobile Handsets & Smartphones",
					Definition:    "Touchscreen smartphones running iOS, Android, or HarmonyOS, as well as ruggedized keypad feature phones.",
					Status:         "active",
					IsLeaf:         false,
					Translations: map[string]string{
						"am": "ስማርትፎኖች እና የሞባይል ቀፎዎች",
						"or": "Ismaartfoonii fi Bilbila Harkaa",
						"ti": "ስማርትፎናትን ናይ ኢድ ስልኪታትን",
					},
					Aliases: map[string][]string{
						"en": {"smartphones", "mobile handsets", "touchscreen phones", "cellphones"},
						"am": {"ስማርትፎን", "ቀፎ", "የሞባይል ቀፎ", "አንድሮይድ ስልክ", "አይፎን"},
						"or": {"ismaartfooni", "bilbila harkaa"},
						"ti": {"ስማርትፎን", "ቀፎ ስልኪ"},
					},
					Children: []TaxonomyNode{
						// Level 3: Touchscreen Smartphones
						{
							ID:            "cat-touchscreen-smartphones",
							Slug:          "touchscreen-smartphones",
							Code:          "PHONE_SMART",
							Level:         3,
							CanonicalName: "Touchscreen Smartphones (ዘመናዊ ስልኮች)",
							Definition:    "Capacitive touchscreen cellular devices supporting 4G/5G broadband, mobile applications, and high-resolution cameras.",
							Status:         "active",
							IsLeaf:         false,
							Identifiers: []TaxonomyIdentifier{
								{Type: "UNSPSC", Code: "43191501", Name: "Mobile phones"},
								{Type: "GS1_GPC", Code: "10001454", Name: "Mobile Telephones/Smartphones"},
							},
							Translations: map[string]string{
								"am": "ዘመናዊ የንክኪ ስክሪን ስልኮች",
								"or": "Bilbila Ismaartii Iskiriinii",
								"ti": "ናይ ንክኪ ስክሪን ዘመናዊ ስልኪታት",
							},
							Aliases: map[string][]string{
								"en": {"smartphone", "android phone", "iphone", "galaxy phone", "5g phone"},
								"am": {"ስማርት ስልክ", "የንክኪ ስልክ", "አንድሮይድ", "አይፎን"},
								"or": {"bilbila ammayyaa", "bilbila iskiriinii"},
								"ti": {"ዘመናዊ ስልኪ", "ስልኪ ንክኪ"},
							},
							Attributes: []TaxonomyAttribute{
								{Slug: "storage-capacity", NameEN: "Internal Storage", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"32GB", "64GB", "128GB", "256GB", "512GB", "1TB"}}},
								{Slug: "ram-capacity", NameEN: "RAM Memory", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"2GB", "3GB", "4GB", "6GB", "8GB", "12GB", "16GB", "24GB"}}},
								{Slug: "cellular-network-gen", NameEN: "Network Generation", DataType: "select", IsRequired: false, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"5G", "4G LTE", "3G"}}},
								{Slug: "sim-slot-type", NameEN: "SIM Configuration", DataType: "select", IsRequired: false, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Dual SIM (Physical)", "Single SIM + eSIM", "Dual eSIM", "Single SIM"}}},
								{Slug: "item-condition", NameEN: "Condition", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"brand_new", "like_new", "good", "fair"}}},
							},
							Children: []TaxonomyNode{
								// Level 4: Flagship Tier Smartphones
								{
									ID:            "leaf-flagship-smartphones",
									Slug:          "flagship-tier-smartphones",
									Code:          "PHONE_FLAGSHIP",
									Level:         4,
									CanonicalName: "Flagship & Ultra-Premium Smartphones",
									Definition:    "Top-tier performance devices featuring advanced optical zoom cameras, titanium/aluminum chassis, and high-refresh OLED displays.",
									Status:         "active",
									IsLeaf:         true,
									Translations: map[string]string{
										"am": "ፍላግሺፕ እና ከፍተኛ ደረጃ ስማርትፎኖች",
										"or": "Ismaartfoonii Sadarkaa Olaanaa",
										"ti": "ላዕለዋይ ደረጃ ዘመናዊ ስልኪታት",
									},
								},
								// Level 4: Mid-Range & Budget Smartphones
								{
									ID:            "leaf-midrange-smartphones",
									Slug:          "midrange-budget-smartphones",
									Code:          "PHONE_MIDRANGE",
									Level:         4,
									CanonicalName: "Mid-Range & Everyday Budget Smartphones",
									Definition:    "Value-optimized smartphones balancing battery endurance, dual/triple camera configurations, and cost accessibility.",
									Status:         "active",
									IsLeaf:         true,
									Translations: map[string]string{
										"am": "መካከለኛ እና ተመጣጣኝ ዋጋ ያላቸው ስልኮች",
										"or": "Ismaartfoonii Giddu-galeessaa fi Gatii Madaalawaa",
										"ti": "ማእከላይን ህዝባዊ ዋጋን ዘለዎም ስልኪታት",
									},
								},
							},
						},
						// Level 3: Basic & Feature Phones (Keypad)
						{
							ID:            "cat-feature-keypad-phones",
							Slug:          "feature-keypad-phones",
							Code:          "PHONE_KEYPAD",
							Level:         3,
							CanonicalName: "Feature Phones & Keypad Handsets (የቁልፍ ስልኮች)",
							Definition:    "Push-button physical keypad mobile handsets with ultra-long battery standby, FM radios, and flashlights.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የቁልፍ / ባትሪ ስልኮች",
								"or": "Bilbila Qabduu (Batirii)",
								"ti": "ናይ ቁልፊ / ባትሪ ስልኪታት",
							},
							Aliases: map[string][]string{
								"en": {"keypad phone", "button phone", "feature phone", "itel phone", "nokia torch"},
								"am": {"የቁልፍ ስልክ", "የባትሪ ስልክ", "አይቴል", "ኖኪያ ባትሪ"},
								"or": {"bilbila qabduu", "bilbila batirii"},
								"ti": {"ናይ ቁልፊ ስልኪ", "ናይ ባትሪ ስልኪ"},
							},
						},
					},
				},
				// Super Family 7.2: Mobile Physical Parts & Replacement Components
				{
					ID:            "sf-mobile-replacement-parts",
					Slug:          "mobile-phone-replacement-parts",
					Code:          "ELEC_PHONE_PARTS",
					Level:         2,
					CanonicalName: "Mobile Phone Physical Replacement Parts",
					Definition:    "Physical hardware replacement items and spare modules for mobile handset repair.",
					InclusionRules: []string{"Replacement LCD/OLED display panels", "Internal replacement batteries", "Charging port flex cables", "Back glass housings", "Camera modules"},
					ExclusionRules: []string{"Mobile repair labor and screen fitting services (see Professional Repair Services)"},
					Status:         "active",
					IsLeaf:         false,
					Translations: map[string]string{
						"am": "የሞባይል ስልክ መለዋወጫ እቃዎች እና ስክሪኖች",
						"or": "Qaama Jijjiirraa Bilbila Harkaa",
						"ti": "ናይ ሞባይል ስልኪ መለዋወጢ ኣቑሑት",
					},
					Children: []TaxonomyNode{
						// Level 3: Replacement Display Assemblies & Glass
						{
							ID:            "cat-phone-replacement-screens",
							Slug:          "replacement-screens-lcd-oled",
							Code:          "PART_PHONE_SCREEN",
							Level:         3,
							CanonicalName: "Replacement Phone Screens & OLED Displays",
							Definition:    "Physical digitizer glass and complete LCD/AMOLED replacement display screens.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የሞባይል መለዋወጫ ስክሪኖች",
								"or": "Iskiriinii Jijjiirraa Bilbilaa",
								"ti": "ናይ ሞባይል መለዋወጢ ስክሪናት",
							},
							Aliases: map[string][]string{
								"en": {"replacement screen", "phone display", "lcd screen", "oled panel", "touch glass"},
								"am": {"የስልክ ስክሪን", "ስክሪን መለዋወጫ", "ኤልሲዲ", "የሞባይል መስታወት"},
								"or": {"iskiriinii bilbilaa", "fuula bilbilaa"},
								"ti": {"ናይ ስልኪ ስክሪን", "መስታወት ስልኪ"},
							},
						},
						// Level 3: Internal Replacement Batteries
						{
							ID:            "cat-phone-replacement-batteries",
							Slug:          "phone-replacement-batteries",
							Code:          "PART_PHONE_BATTERY",
							Level:         3,
							CanonicalName: "Internal Replacement Phone Batteries",
							Definition:    "Lithium-ion and lithium-polymer internal replacement cells for mobile handsets.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የሞባይል ባትሪዎች",
								"or": "Baatirii Jijjiirraa Bilbilaa",
								"ti": "ናይ ሞባይል መለዋወጢ ባትሪታት",
							},
						},
					},
				},
				// Super Family 7.3: Mobile Accessories
				{
					ID:            "sf-mobile-accessories",
					Slug:          "mobile-phone-accessories",
					Code:          "ELEC_PHONE_ACC",
					Level:         2,
					CanonicalName: "Mobile Phone Accessories",
					Definition:    "Protective cases, tempered glass, wall chargers, power banks, and cables for mobile devices.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የሞባይል ስልክ ተጨማሪ እቃዎች",
						"or": "Meeshaalee Dabalataa Bilbilaa",
						"ti": "ናይ ሞባይል ስልኪ ተወሳኺ ኣቑሑት",
					},
					Aliases: map[string][]string{
						"en": {"phone accessories", "phone case", "charger", "power bank", "usb cable"},
						"am": {"የስልክ ቻርጀር", "ፓወር ባንክ", "የስልክ ከቨር", "ኬብል"},
						"or": {"chaarjara", "paawor baankii", "hagooggii bilbilaa"},
						"ti": {"ቻርጀር", "ፓወር ባንክ", "ናይ ስልኪ ከቨር"},
					},
				},
			},
		},

		// 8. Root: Computers, Computing Hardware & Peripherals
		{
			ID:            "root-computers-it",
			Slug:          "computers-it",
			Code:          "COMP_HARDWARE",
			Level:         1,
			CanonicalName: "Computers, Computing Hardware & Peripherals",
			Definition:    "Laptops, desktops, enterprise servers, computer internal hardware, monitors, printers, and storage media.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Laptop",
			Translations: map[string]string{
				"am": "ኮምፒውተሮች፣ የኮምፒውተር ሃርድዌር እና ፔሪፌራልስ",
				"or": "Kompiitara, Haardiweerii fi Meeshaalee Walqabatan",
				"ti": "ኮምፒዩተራት፣ ናይ ኮምፒዩተር ሃርድዌርን ተዛመድቲ ኣቑሑትን",
			},
			Children: []TaxonomyNode{
				// Super Family 8.1: Laptops & Notebooks
				{
					ID:            "sf-laptops-notebooks",
					Slug:          "computers-laptops",
					Code:          "COMP_LAPTOPS",
					Level:         2,
					CanonicalName: "Laptops, Notebooks & Mobile Workstations",
					Definition:    "Portable personal computers including business ultrabooks, mobile workstations, and gaming laptops.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "43211503", Name: "Notebook computers"},
						{Type: "GS1_GPC", Code: "10001452", Name: "Laptops/Notebooks"},
					},
					Translations: map[string]string{
						"am": "ላፕቶፖች እና ተንቀሳቃሽ ኮምፒውተሮች",
						"or": "Laaptooppii fi Kompiitara Harkaa",
						"ti": "ላፕቶፓትን ተንቀሳቐስቲ ኮምፒዩተራትን",
					},
					Aliases: map[string][]string{
						"en": {"laptop", "notebook", "macbook", "thinkpad", "workstation laptop"},
						"am": {"ላፕቶፕ", "ኮምፒውተር", "ማክቡክ", "ቲንክፓድ"},
						"or": {"laaptooppii", "kompiitara"},
						"ti": {"ላፕቶፕ", "ኮምፒዩተር"},
					},
					Attributes: []TaxonomyAttribute{
						{Slug: "processor-cpu", NameEN: "Processor (CPU)", DataType: "text", IsRequired: true, IsFilterable: true},
						{Slug: "ram-capacity", NameEN: "RAM Memory", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"8GB", "16GB", "32GB", "64GB", "128GB"}}},
						{Slug: "storage-capacity", NameEN: "Solid State Drive (SSD)", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"256GB", "512GB", "1TB", "2TB", "4TB"}}},
					},
				},
				// Super Family 8.2: Desktop Computers & All-in-Ones
				{
					ID:            "sf-desktop-computers",
					Slug:          "desktop-computers-all-in-ones",
					Code:          "COMP_DESKTOPS",
					Level:         2,
					CanonicalName: "Desktop Computers, Workstations & All-in-Ones",
					Definition:    "Tower desktop PCs, mini PCs, all-in-one workstations, and gaming battle-stations.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የጠረጴዛ ላይ ኮምፒውተሮች እና ወርክስቴሽኖች",
						"or": "Kompiitara Minjaalaa",
						"ti": "ናይ ጣውላ ኮምፒዩተራትን ዎርክስቴሽናትን",
					},
				},
				// Super Family 8.3: Printers, Scanners & Office Copiers
				{
					ID:            "sf-printers-scanners",
					Slug:          "printers-scanners-copiers",
					Code:          "COMP_PRINTERS",
					Level:         2,
					CanonicalName: "Printers, High-Speed Scanners & Copiers",
					Definition:    "Laser printers, continuous ink-tank color printers, POS receipt printers, and flatbed scanners.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ፕሪንተሮች፣ ስካነሮች እና የፎቶ ኮፒ ማሽኖች",
						"or": "Pirintara, Iskaanara fi Koppii",
						"ti": "ፕሪንተራት፣ ስካነራትን ፎቶ ኮፒን",
					},
				},
			},
		},

		// 9. Root: Software, Operating Systems & Digital Applications
		{
			ID:            "root-software-digital",
			Slug:          "software-operating-systems-applications",
			Code:          "SOFT_DIGITAL",
			Level:         1,
			CanonicalName: "Software, Operating Systems & Digital Licenses",
			Definition:    "Packaged software licenses, enterprise accounting software (Peachtree), CAD packages, OS keys, and office suites.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Code",
			Translations: map[string]string{
				"am": "ሶፍትዌር፣ ኦፕሬቲንግ ሲስተሞች እና የዲጂታል ፍቃዶች",
				"or": "Sooftiweerii fi Heyyama Dijitaalaa",
				"ti": "ሶፍትዌር፣ ኦፕሬቲንግ ሲስተማትን ዲጂታል ፍቓዳትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-accounting-erp-software",
					Slug:          "accounting-business-erp-software",
					Code:          "SOFT_ERP",
					Level:         2,
					CanonicalName: "Accounting, Payroll & Enterprise ERP Software",
					Definition:    "Commercial accounting software (Peachtree/Sage 50, QuickBooks, Odoo modules) and ERCA-compliant tax suites.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የሂሳብ እና የድርጅት አስተዳደር ሶፍትዌሮች (ፒችትሪ)",
						"or": "Sooftiweerii Herregaa fi Daldalaa (Peachtree)",
						"ti": "ናይ ሒሳብን ንግድን ሶፍትዌራት (ፒችትሪ)",
					},
				},
			},
		},

		// 10. Root: Telecommunications & Networking Infrastructure
		{
			ID:            "root-telecom-networking",
			Slug:          "telecommunications-networking-infrastructure",
			Code:          "TELE_NETWORKING",
			Level:         1,
			CanonicalName: "Telecommunications & Networking Infrastructure",
			Definition:    "4G/5G Wi-Fi modems, enterprise routers, managed switches, fiber optic patch cables, and telecom antennas.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Wifi",
			Translations: map[string]string{
				"am": "ቴሌኮሙኒኬሽን እና የኔትወርክ መሰረተ ልማት",
				"or": "Telekominikeeshinii fi Netwoorkii",
				"ti": "ቴሌኮሙኒኬሽንን ናይ መርበብ ሓበሬታ ትሕተ-ቅርጽን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-wifi-routers-modems",
					Slug:          "wireless-routers-4g-5g-modems",
					Code:          "NET_ROUTERS",
					Level:         2,
					CanonicalName: "Wireless Routers, 4G/5G Pocket Modems & Access Points",
					Definition:    "Safaricom & Ethio Telecom compatible 4G/5G SIM routers, Wi-Fi 6 mesh nodes, and outdoor wireless bridges.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ዋይፋይ ራውተሮች እና የሞባይል ኢንተርኔት ሞደሞች",
						"or": "Rawtara fi Moodeemii Intarneetii",
						"ti": "ዋይፋይ ራውተራትን ናይ ኢንተርኔት ሞደማትን",
					},
				},
			},
		},

		// 11. Root: Audio, Video, TV & Home Entertainment
		{
			ID:            "root-tv-audio",
			Slug:          "audio-video-tv-home-entertainment",
			Code:          "AV_ENTERTAINMENT",
			Level:         1,
			CanonicalName: "Audio, Video, TV & Home Entertainment",
			Definition:    "Smart LED/QLED televisions, satellite receivers/decoders, home theatre speaker arrays, and studio headphones.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Tv",
			Translations: map[string]string{
				"am": "ቴሌቪዥን፣ የድምጽ ሲስተሞች እና የቤት መዝናኛ",
				"or": "Televiziyinii fi Meeshaalee Sagalee",
				"ti": "ቴሌቪዥን፣ ናይ ድምጺ ሲስተማትን ናይ ገዛ መዘናግዕን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-smart-televisions",
					Slug:          "smart-tvs-flat-panels",
					Code:          "TV_SMART",
					Level:         2,
					CanonicalName: "Smart 4K UHD Televisions & Flat Panels",
					Definition:    "32-inch to 85-inch 4K UHD smart screens with Android TV, Tizen, and WebOS streaming capabilities.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ስማርት ቴሌቪዥኖች እና ፍላት ስክሪን",
						"or": "Televiziyinii Ismaartii fi Iskiriinii Balleessaa",
						"ti": "ስማርት ቴሌቪዥናትን ፍላት ስክሪንን",
					},
					Aliases: map[string][]string{
						"en": {"tv", "smart tv", "television", "4k tv", "oled tv"},
						"am": {"ቲቪ", "ቴሌቪዥን", "ስማርት ቲቪ", "ፍላት ቲቪ"},
						"or": {"tiivii", "televiziyinii"},
						"ti": {"ቲቪ", "ቴሌቪዥን"},
					},
				},
			},
		},

		// 12. Root: Photography, Videography & Optical Equipment
		{
			ID:            "root-cameras-optics",
			Slug:          "cameras-optics",
			Code:          "PHOTO_OPTICS",
			Level:         1,
			CanonicalName: "Photography, Videography & Optical Equipment",
			Definition:    "Full-frame mirrorless digital cameras, DSLR systems, professional video cinema gear, and camera lenses.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Camera",
			Translations: map[string]string{
				"am": "ካሜራዎች፣ ቪዲዮ እና የፎቶግራፊ እቃዎች",
				"or": "Kaameraa fi Meeshaalee Suuraa Kaasuu",
				"ti": "ካሜራታት፣ ቪዲዮን ናይ ፎቶግራፍ ኣቑሑትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-mirrorless-dslr-cameras",
					Slug:          "mirrorless-dslr-cameras",
					Code:          "PHOTO_CAMERAS",
					Level:         2,
					CanonicalName: "Mirrorless & DSLR Digital Cameras",
					Definition:    "Professional full-frame and APS-C digital camera bodies for wedding, studio, and commercial media.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የፎቶ እና የቪዲዮ ካሜራዎች",
						"or": "Kaameraalee Suuraa Dijitaalaa",
						"ti": "ናይ ፎቶን ቪዲዮን ካሜራታት",
					},
					Aliases: map[string][]string{
						"en": {"camera", "dslr", "mirrorless", "sony alpha", "canon camera"},
						"am": {"ካሜራ", "የፎቶ ካሜራ", "የቪዲዮ ካሜራ"},
						"or": {"kaameraa"},
						"ti": {"ካሜራ"},
					},
				},
			},
		},

		// 13. Root: Video Gaming, Consoles & Esports Equipment
		{
			ID:            "root-gaming-esports",
			Slug:          "video-gaming-consoles-esports",
			Code:          "GAME_CONSOLES",
			Level:         1,
			CanonicalName: "Video Gaming, Consoles & Esports Equipment",
			Definition:    "Sony PlayStation, Microsoft Xbox, Nintendo consoles, game discs, gaming controllers, and esports accessories.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Gamepad2",
			Translations: map[string]string{
				"am": "ጌሚንግ፣ የቪዲዮ ጌም ኮንሶሎች እና መለዋወጫዎች",
				"or": "Konsolii Tapha Viidiyoo fi Geemiingii",
				"ti": "ጌሚንግ፣ ናይ ቪዲዮ ጸወታ ኮንሶላትን ኣቑሑትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-gaming-consoles",
					Slug:          "home-gaming-consoles",
					Code:          "GAME_HARDWARE",
					Level:         2,
					CanonicalName: "Video Game Consoles (PlayStation & Xbox)",
					Definition:    "PlayStation 4, PlayStation 5, Xbox Series X/S hardware consoles and wireless gamepads.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ፕሌይስቴሽን እና የጌም ኮንሶሎች",
						"or": "Pileyisteeshinii fi Konsolii Taphaa",
						"ti": "ፕሌይስቴሽንን ናይ ጸወታ ኮንሶላትን",
					},
					Aliases: map[string][]string{
						"en": {"playstation", "ps5", "ps4", "xbox", "gaming console", "controller"},
						"am": {"ፕሌይስቴሽን", "ፒኤስ4", "ፒኤስ5", "የጌም እቃ", "ጆይስቲክ"},
						"or": {"pileyisteeshinii", "geemii"},
						"ti": {"ፕሌይስቴሽን", "ናይ ጸወታ እቃ"},
					},
				},
			},
		},
	}
}
