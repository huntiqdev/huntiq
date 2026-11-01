package main

func getHomeDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 40. Root: Major Home & Kitchen Appliances
		{
			ID:            "root-major-appliances",
			Slug:          "major-home-kitchen-appliances",
			Code:          "HOME_APPLIANCES",
			Level:         1,
			CanonicalName: "Major Home & Kitchen Appliances",
			Definition:    "Large white goods, double-door frost-free refrigerators, chest deep-freezers, front-load washing machines, and gas/electric cooking ranges.",
			InclusionRules: []string{"Refrigerators", "Deep freezers", "Washing machines", "Electric stoves / cookers", "Microwave ovens"},
			ExclusionRules: []string{"Small blenders and electric mitad (see Small Domestic Appliances)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Tv",
			Translations: map[string]string{
				"am": "የቤት ውስጥ ትላልቅ ኤሌክትሮኒክስ (ፍሪጅ፣ ማጠቢያ ማሽን)",
				"or": "Meeshaalee Manaa Gurguddoo (Firiijii fi Meeshaa Dhiqqaa)",
				"ti": "ናይ ገዛ ዓበይቲ ኤሌክትሮኒክስ (ፍሪጅ፣ ማሽን ሕጽበ)",
			},
			Aliases: map[string][]string{
				"en": {"appliances", "refrigerator", "fridge", "freezer", "washing machine", "stove", "oven"},
				"am": {"ፍሪጅ", "ማቀዝቀዣ", "የልብስ ማጠቢያ", "ምድጃ", "ኦቨን"},
				"or": {"firiijii", "maashina uffataa", "eeksaasii"},
				"ti": {"ፍሪጅ", "ናይ ክዳን ማሽን", "እቶን"},
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-refrigerators-freezers",
					Slug:          "refrigerators-deep-freezers",
					Code:          "APP_REFRIGERATION",
					Level:         2,
					CanonicalName: "Refrigerators, Chest Freezers & Coolers (ፍሪጆች)",
					Definition:    "Energy-efficient inverter compressor refrigerators, upright frost-free freezers, and horizontal commercial chest freezers.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "52141501", Name: "Domestic refrigerators"},
					},
					Translations: map[string]string{
						"am": "ፍሪጆች እና የበረዶ ማቀዝቀዣዎች",
						"or": "Firiijii fi Kuusaa Qabbanaawaa",
						"ti": "ፍሪጃትን ናይ በረድ መዝቀጥታትን",
					},
				},
				{
					ID:            "sf-washing-machines",
					Slug:          "automatic-washing-machines-dryers",
					Code:          "APP_LAUNDRY",
					Level:         2,
					CanonicalName: "Automatic Washing Machines & Tumble Dryers",
					Definition:    "Front-loading and top-loading automatic laundry washing machines and combo washer-dryers (6kg to 16kg).",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "አውቶማቲክ የልብስ ማጠቢያ ማሽኖች",
						"or": "Maashina Uffata Dhiqu",
						"ti": "ኣውቶማቲክ ናይ ክዳውንቲ መሕጸቢ ማሽናት",
					},
				},
			},
		},

		// 41. Root: Small Domestic Appliances & Traditional Cookware
		{
			ID:            "root-small-appliances",
			Slug:          "small-domestic-appliances-cookware",
			Code:          "HOME_SMALL_APP",
			Level:         1,
			CanonicalName: "Small Domestic Appliances & Traditional Cookware",
			Definition:    "Electric Injera baking pans (Mitad), electric coffee roasters, food blenders, meat mincers, water dispensers, and traditional clay coffee pots (Jebena).",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Flame",
			Translations: map[string]string{
				"am": "ትናንሽ የቤት እቃዎች፣ የኤሌክትሪክ ምጣድ እና ጀበና",
				"or": "Meeshaalee Kushiinaa fi Eelee Elektirikii",
				"ti": "ናይ ክሽነ ንዋያት፣ ናይ ኤሌክትሪክ ምጣድን ጀበናን",
			},
			Children: []TaxonomyNode{
				// Super Family 41.1: Electric Injera Baking Pans (Mitad / ምጣድ)
				{
					ID:            "sf-electric-injera-mitad",
					Slug:          "electric-injera-baking-pans-mitad",
					Code:          "APP_MITAD",
					Level:         2,
					CanonicalName: "Electric Injera Baking Pans & Clay Plates (የኤሌክትሪክ ምጣድ)",
					Definition:    "Traditional heavy-clay and digital thermostat electric baking griddles essential for Ethiopian injera preparation.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የኤሌክትሪክ የእንጀራ ምጣዶች",
						"or": "Eelee Buddeenaa Elektirikii",
						"ti": "ናይ ኤሌክትሪክ ናይ እንጀራ ምጣድ",
					},
					Aliases: map[string][]string{
						"en": {"injera mitad", "electric mitad", "mitad", "baking pan", "clay pan"},
						"am": {"ምጣድ", "የኤሌክትሪክ ምጣድ", "የእንጀራ ምጣድ", "የሸክላ ምጣድ"},
						"or": {"eelee", "eelee elektirikii", "eelee buddeenaa"},
						"ti": {"ምጣድ", "ናይ ኤሌክትሪክ ምጣድ"},
					},
				},
				// Super Family 41.2: Traditional Clay Coffee Pots & Ceremony Sets
				{
					ID:            "sf-traditional-jebena-sets",
					Slug:          "clay-coffee-pots-jebena-ceremony",
					Code:          "APP_JEBENA",
					Level:         2,
					CanonicalName: "Traditional Clay Coffee Pots (Jebena) & Ceremony Sets",
					Definition:    "Handcrafted earthenware clay Jebenas, charcoal stoves (Keshkesh), ceramic Sini coffee cups, and woven Rekebot trays.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የቡና ጀበና፣ ሲኒዎች እና ረከቦቶች",
						"or": "Jabanaa Bunaa fi Siinii",
						"ti": "ጀበና ቡን፣ ሲኒታትን ረከቦትን",
					},
					Aliases: map[string][]string{
						"en": {"jebena", "coffee pot", "sini", "rekebot", "coffee ceremony"},
						"am": {"ጀበና", "የቡና ጀበና", "ሲኒ", "ረከቦት"},
						"or": {"jabanaa", "siinii"},
						"ti": {"ጀበና", "ሲኒ", "ረከቦት"},
					},
				},
			},
		},

		// 42. Root: Furniture & Interior Architectural Furnishings
		{
			ID:            "root-furniture-furnishings",
			Slug:          "furniture-interior-furnishings",
			Code:          "HOME_FURNITURE",
			Level:         1,
			CanonicalName: "Furniture & Interior Architectural Furnishings",
			Definition:    "Living room sofa sets, highland leather couches, master beds, orthopedic spring mattresses, dining tables, and executive office desks.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Armchair",
			Translations: map[string]string{
				"am": "የቤት እቃዎች፣ ፈርኒቸር እና ሶፋዎች",
				"or": "Meeshaalee Manaa, Barcuma fi Soofaa",
				"ti": "ፈርኒቸር፣ ሶፋታትን ናይ ገዛ ኣቑሑትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-living-room-sofas",
					Slug:          "living-room-sofas-sectional-couches",
					Code:          "FURN_SOFAS",
					Level:         2,
					CanonicalName: "Living Room Sofa Sets & L-Shape Sectionals (ሶፋዎች)",
					Definition:    "Fabric upholstered and genuine leather living room suites, recliners, and wooden carved salon seating.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የሳሎን ሶፋዎች እና ኤል-ቅርጽ ሶፋ",
						"or": "Soofaa Saaloonaa fi Barcuma",
						"ti": "ናይ ሳሎን ሶፋታትን መቐመጢታትን",
					},
					Aliases: map[string][]string{
						"en": {"sofa", "couch", "living room set", "sectional", "leather sofa"},
						"am": {"ሶፋ", "የሳሎን ሶፋ", "ኤል ሶፋ", "የቆዳ ሶፋ"},
						"or": {"soofaa", "barcuma saaloonaa"},
						"ti": {"ሶፋ", "ናይ ሳሎን ሶፋ"},
					},
				},
				{
					ID:            "sf-bedroom-beds-mattresses",
					Slug:          "bedroom-beds-orthopedic-mattresses",
					Code:          "FURN_BEDS",
					Level:         2,
					CanonicalName: "Master Beds, Headboards & Orthopedic Mattresses",
					Definition:    "Solid wood and upholstered master beds (King, Queen, Single), memory foam and pocket spring mattresses.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "አልጋዎች እና የስፖንጅ ፍራሾች",
						"or": "Siree fi Firaashii",
						"ti": "ዓራትታትን ፍራሻትን",
					},
				},
			},
		},
	}
}
