package main

func getPropertyDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 21. Root: Residential Properties & Living Spaces
		{
			ID:            "root-residential-properties",
			Slug:          "residential-properties-living-spaces",
			Code:          "REAL_RESIDENTIAL",
			Level:         1,
			CanonicalName: "Residential Properties & Living Spaces",
			Definition:    "Physical residential housing units, condominium apartments, detached luxury villas, townhouses, and studio suites.",
			InclusionRules: []string{"Condominium apartments", "Real estate developer apartments", "Detached villas", "Townhouses", "Studio apartments"},
			ExclusionRules: []string{"Commercial office spaces (see Commercial Properties)", "Bare residential plots without building structures (see Land & Acreage)"},
			Examples:       []string{"2-bedroom CMC apartment", "Luxury villa in Bole Atlas", "Studio apartment in Yeka", "40/60 condominium unit"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Home",
			Translations: map[string]string{
				"am": "የመኖሪያ ቤቶች እና አፓርታማዎች",
				"or": "Manni Jireenyaa fi Appaartamaa",
				"ti": "ናይ መንበሪ ኣባይትን ኣፓርታማታትን",
			},
			Aliases: map[string][]string{
				"en": {"residential", "house", "apartment", "villa", "condo", "condominium", "real estate", "living space"},
				"am": {"ቤት", "የመኖሪያ ቤት", "አፓርታማ", "ቪላ", "ኮንዶሚኒየም", "ኮንዶ", "ስቱዲዮ"},
				"or": {"mana", "mana jireenyaa", "appaartamaa", "viilaa", "kondominiyemii"},
				"ti": {"ቤት", "ናይ መንበሪ ቤት", "ኣፓርታማ", "ቪላ", "ኮንዶሚኒየም"},
			},
			Children: []TaxonomyNode{
				// Super Family 21.1: Multi-Unit Apartments & Condominiums
				{
					ID:            "sf-apartments-condos",
					Slug:          "apartments-condominiums-flats",
					Code:          "PROP_APARTMENTS",
					Level:         2,
					CanonicalName: "Multi-Unit Apartments & Condominium Flats",
					Definition:    "Multi-story residential apartment units, luxury penthouses, developer apartment flats, and government 20/80 and 40/60 condominium housing.",
					Status:         "active",
					IsLeaf:         false,
					Translations: map[string]string{
						"am": "አፓርታማዎች እና የኮንዶሚኒየም ቤቶች",
						"or": "Appaartamaa fi Kondominiyemii",
						"ti": "ኣፓርታማታትን ኮንዶሚኒየምን",
					},
					Children: []TaxonomyNode{
						// Level 3: Standard & Developer Apartments
						{
							ID:            "cat-standard-apartments",
							Slug:          "modern-developer-apartments",
							Code:          "APT_MODERN",
							Level:         3,
							CanonicalName: "Modern Developer Apartments & Flats",
							Definition:    "Private real estate developments (Ayat, Tsehay, Noah, Gift, Flintstone, Emaar) with elevators, backup generators, and designated parking.",
							Status:         "active",
							IsLeaf:         true,
							Identifiers: []TaxonomyIdentifier{
								{Type: "UNSPSC", Code: "80131502", Name: "Apartment or condominium rental or sale"},
							},
							Translations: map[string]string{
								"am": "ዘመናዊ የሪል እስቴት አፓርታማዎች",
								"or": "Appaartamaa Ammayyaa Daldalaa",
								"ti": "ዘመናዊ ናይ ሪል እስቴት ኣፓርታማታት",
							},
							Attributes: []TaxonomyAttribute{
								{Slug: "property-area-sqm", NameEN: "Floor Area (m²)", DataType: "number", Unit: "m²", IsRequired: true, IsFilterable: true},
								{Slug: "bedrooms-count", NameEN: "Number of Bedrooms", DataType: "number", IsRequired: true, IsFilterable: true},
								{Slug: "bathrooms-count", NameEN: "Number of Bathrooms", DataType: "number", IsRequired: true, IsFilterable: true},
								{Slug: "furnishing-status", NameEN: "Furnishing Status", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Furnished", "Semi-Furnished", "Unfurnished"}}},
								{Slug: "property-deed-type", NameEN: "Title Deed / Documentation", DataType: "select", IsRequired: false, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Digital Land Certificate / Carte", "Holding Certificate", "Sub-lease Agreement", "Under Construction Contract"}}},
							},
						},
						// Level 3: Condominiums (20/80 & 40/60)
						{
							ID:            "cat-condominium-units",
							Slug:          "public-condominium-units",
							Code:          "APT_CONDO",
							Level:         3,
							CanonicalName: "Public Condominium Units (ኮንዶሚኒየም)",
							Definition:    "Studio, 1-bedroom, 2-bedroom, and 3-bedroom government scheme condominium units.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የኮንዶሚኒየም ቤቶች (20/80 እና 40/60)",
								"or": "Kondominiyemii (20/80 fi 40/60)",
								"ti": "ናይ ኮንዶሚኒየም ኣባይቲ",
							},
						},
					},
				},
				// Super Family 21.2: Detached Villas & Compound Houses
				{
					ID:            "sf-villas-compound-houses",
					Slug:          "detached-villas-compound-houses",
					Code:          "PROP_VILLAS",
					Level:         2,
					CanonicalName: "Detached Villas & Private Compound Residences",
					Definition:    "Single-family standalone homes, compound villas with private gardens, security gates, and servant quarters.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ቪላዎች እና የጊቢ ውስጥ ቤቶች",
						"or": "Viilaa fi Manni Compound Qabu",
						"ti": "ቪላታትን ናይ ውሽጢ ቀጽሪ ኣባይትን",
					},
					Aliases: map[string][]string{
						"en": {"villa", "detached house", "compound house", "mansion"},
						"am": {"ቪላ", "የጊቢ ቤት", "ቪላ ቤት", "ቪላዎች"},
						"or": {"viilaa", "mana qe'ee"},
						"ti": {"ቪላ", "ናይ ቀጽሪ ቤት"},
					},
				},
			},
		},

		// 22. Root: Commercial Properties, Offices & Retail Spaces
		{
			ID:            "root-commercial-properties",
			Slug:          "commercial-properties-offices-retail",
			Code:          "REAL_COMMERCIAL",
			Level:         1,
			CanonicalName: "Commercial Properties, Offices & Retail Spaces",
			Definition:    "Office spaces, commercial retail mall shops, banks, showrooms, and entire mixed-use commercial buildings.",
			InclusionRules: []string{"Office spaces in corporate towers", "Retail shops in shopping centers (Merkato, Bole, Megenagna)", "Commercial showrooms", "Whole commercial buildings"},
			ExclusionRules: []string{"Residential apartments (see Residential Properties)", "Industrial manufacturing factories (see Industrial Real Estate)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Building2",
			Translations: map[string]string{
				"am": "የንግድ ቤቶች፣ ቢሮዎች እና ሱቆች",
				"or": "Manni Daldalaa, Biiroo fi Suuqii",
				"ti": "ናይ ንግዲ ኣባይቲ፣ ቢሮታትን ድኳናትን",
			},
			Aliases: map[string][]string{
				"en": {"commercial", "office", "shop", "retail space", "commercial building", "mall shop"},
				"am": {"የንግድ ቤት", "ሱቅ", "ቢሮ", "የገበያ አዳራሽ", "ሞል", "ሾውሩም"},
				"or": {"mana daldalaa", "biiroo", "suuqii", "moolii"},
				"ti": {"ናይ ንግዲ ቤት", "ድኳን", "ቢሮ", "ሞል"},
			},
			Children: []TaxonomyNode{
				// Super Family 22.1: Retail Mall Shops & Boutiques
				{
					ID:            "sf-retail-shops-malls",
					Slug:          "retail-shops-shopping-mall-spaces",
					Code:          "COMM_SHOPS",
					Level:         2,
					CanonicalName: "Retail Shops, Boutiques & Shopping Mall Spaces (ሱቆች)",
					Definition:    "Enclosed retail shops inside commercial malls, street-level retail boutiques, and Merkato wholesale stalls.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "80131503", Name: "Commercial retail space rental or sale"},
					},
					Translations: map[string]string{
						"am": "የሞል ሱቆች እና የችርቻሮ ንግድ ቦታዎች",
						"or": "Suuqii Daldalaa fi Moolii",
						"ti": "ናይ ሞል ድኳናትን ናይ ንግዲ ቦታታትን",
					},
					Aliases: map[string][]string{
						"en": {"shop", "retail shop", "boutique", "mall space", "merkato shop"},
						"am": {"ሱቅ", "የሞል ሱቅ", "ቡቲክ", "የመደብር ቦታ", "የመርካቶ ሱቅ"},
						"or": {"suuqii", "suuqii moolii"},
						"ti": {"ድኳን", "ናይ ሞል ድኳን"},
					},
				},
				// Super Family 22.2: Corporate Office Suites
				{
					ID:            "sf-corporate-offices",
					Slug:          "corporate-office-suites-floors",
					Code:          "COMM_OFFICES",
					Level:         2,
					CanonicalName: "Corporate Office Suites & Entire Commercial Floors",
					Definition:    "Partitioned professional office units, executive suites, open-plan workspaces, and whole commercial floors.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የቢሮ ክፍሎች እና ሙሉ ህንጻዎች",
						"or": "Kutaa Biiroo fi Gamoo Daldalaa",
						"ti": "ናይ ቢሮ ክፍልታትን ሙሉእ ህንጻታትን",
					},
				},
			},
		},

		// 23. Root: Industrial Real Estate, Warehouses & Logistics Hubs
		{
			ID:            "root-industrial-real-estate",
			Slug:          "industrial-real-estate-warehouses",
			Code:          "REAL_INDUSTRIAL",
			Level:         1,
			CanonicalName: "Industrial Real Estate, Warehouses & Logistics Hubs",
			Definition:    "Storage warehouses (Kality, Dukem, Mojo, Gelan), industrial park manufacturing sheds, and refrigerated cold rooms.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Warehouse",
			Translations: map[string]string{
				"am": "የኢንዱስትሪ ይዞታዎች፣ መጋዘኖች እና ፋብሪካዎች",
				"or": "Mankuusa, Warshaa fi Iddoo Daldalaa Gurguddoo",
				"ti": "ናይ ኢንዱስትሪ ትሕተ-ቅርጽ፣ መኽዘናትን ፋብሪካታትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-storage-warehouses",
					Slug:          "dry-cargo-storage-warehouses",
					Code:          "IND_WAREHOUSES",
					Level:         2,
					CanonicalName: "Dry Cargo Warehouses & Logistics Hubs (መጋዘኖች)",
					Definition:    "High-ceiling dry storage depots, bonded customs warehouses, and logistics transit yards.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የእቃ መጋዘኖች እና የሎጂስቲክስ ማዕከላት",
						"or": "Mankuusa Kuusaa Meeshaa",
						"ti": "ናይ ኣቑሑት መኽዘናት",
					},
				},
			},
		},

		// 24. Root: Land Plots, Agricultural Estates & Development Acreage
		{
			ID:            "root-land-acreage",
			Slug:          "land-plots-agricultural-development-acreage",
			Code:          "REAL_LAND",
			Level:         1,
			CanonicalName: "Land Plots, Agricultural Estates & Development Acreage",
			Definition:    "Zoned residential plots, commercial development land, leased agricultural acreage, and industrial plots.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "MapPin",
			Translations: map[string]string{
				"am": "መሬት፣ እርሻ ቦታዎች እና የልማት ይዞታዎች",
				"or": "Lafa Jireenyaa, Qonnaa fi Daldalaa",
				"ti": "መሬት፣ ናይ ሕርሻ ቦታታትን ናይ ልምዓት ይዞታታትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-residential-land-plots",
					Slug:          "residential-building-land-plots",
					Code:          "LAND_RESIDENTIAL",
					Level:         2,
					CanonicalName: "Residential Building Land Plots (የመኖሪያ መሬት)",
					Definition:    "Surveyed residential plots with certified municipal coordinates and registered lease rights.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የመኖሪያ ቤት መሬት ይዞታዎች",
						"or": "Lafa Ijaarsa Jireenyaa",
						"ti": "ናይ መንበሪ ቤት መሬት",
					},
				},
				{
					ID:            "sf-agricultural-farmland-lease",
					Slug:          "commercial-agricultural-farmland",
					Code:          "LAND_AGRICULTURAL",
					Level:         2,
					CanonicalName: "Commercial Agricultural Farmland & Irrigation Acreage",
					Definition:    "Hectares of arable agricultural farmland with river/groundwater access for commercial crops.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የንግድ እርሻ መሬት እና የመስኖ ይዞታዎች",
						"or": "Lafa Qonnaa Daldalaa fi Jallisii",
						"ti": "ናይ ንግዲ ሕርሻ መሬትን መስኖን",
					},
				},
			},
		},
	}
}
