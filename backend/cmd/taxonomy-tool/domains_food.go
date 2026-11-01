package main

func getFoodDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 49. Root: Beverages, Coffee & Tea
		{
			ID:            "root-coffee-beverages",
			Slug:          "beverages-coffee-tea",
			Code:          "FOOD_COFFEE_BEV",
			Level:         1,
			CanonicalName: "Beverages, Specialty Coffee & Tea",
			Definition:    "Specialty Ethiopian washed/unwashed green coffee beans (Yirgacheffe, Sidama, Guji, Harar, Limu), artisan roasted beans, highland tea, and bottled mineral waters.",
			InclusionRules: []string{"Green export coffee beans", "Roasted artisan whole beans / ground coffee", "Highland black and green tea", "Natural mineral waters and soft drinks"},
			ExclusionRules: []string{"Raw whole grain crops (see Agriculture)", "Coffee shops and cafeteria dining services (see Restaurants & Catering)"},
			Examples:       []string{"Yirgacheffe Grade 1 washed green coffee", "Sidama natural roasted beans", "Wondo Genet herbal tea", "Ambo mineral water"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Coffee",
			Translations: map[string]string{
				"am": "መጠጦች፣ የኢትዮጵያ ቡና እና ሻይ",
				"or": "Dhugaatii, Buna fi Shaayii",
				"ti": "መስተታት፣ ናይ ኢትዮጵያ ቡንን ሻህን",
			},
			Aliases: map[string][]string{
				"en": {"coffee", "beverages", "tea", "yirgacheffe", "sidama coffee", "green coffee", "roasted coffee"},
				"am": {"ቡና", "የኢትዮጵያ ቡና", "ይርጋጨፌ ቡና", "ሲዳማ ቡና", "ሻይ", "መጠጥ"},
				"or": {"buna", "buna itoophiyaa", "shaayii", "dhugaatii"},
				"ti": {"ቡን", "ናይ ኢትዮጵያ ቡን", "ሻሂ", "መስተ"},
			},
			Children: []TaxonomyNode{
				// Super Family 49.1: Specialty Ethiopian Coffee Beans (ቡና)
				{
					ID:            "sf-ethiopian-specialty-coffee",
					Slug:          "ethiopian-specialty-coffee-beans",
					Code:          "COFFEE_BEANS",
					Level:         2,
					CanonicalName: "Specialty Ethiopian Single-Origin Coffee Beans (የኢትዮጵያ ቡና)",
					Definition:    "Arabica green and roasted beans sorted by terroir (Yirgacheffe, Sidama, Guji, Harar, Kaffa, Limu, Jimma).",
					Status:         "active",
					IsLeaf:         false,
					Identifiers: []TaxonomyIdentifier{
						{Type: "TAXONOMIC", Code: "Coffea arabica", Name: "Arabica Coffee"},
						{Type: "UNSPSC", Code: "50201706", Name: "Coffee beans"},
					},
					Translations: map[string]string{
						"am": "የተመረጡ የኢትዮጵያ የቡና ፍሬዎች",
						"or": "Buna Itoophiyaa Filatamaa",
						"ti": "ፍሉይ ናይ ኢትዮጵያ ቡን",
					},
					Attributes: []TaxonomyAttribute{
						{Slug: "coffee-origin-zone", NameEN: "Coffee Origin Terroir", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Yirgacheffe (Gedeo)", "Sidama", "Guji", "Harar (East/West)", "Limu / Jimma", "Kaffa", "Nekemte / Wollega"}}},
						{Slug: "coffee-processing-method", NameEN: "Processing Method", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Washed (Wet Process)", "Natural / Sundried (Dry Process)", "Honey Processed"}}},
						{Slug: "coffee-roast-state", NameEN: "State", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Green Beans (Unroasted)", "Roasted Whole Beans", "Ground Coffee Powder"}}},
					},
					Children: []TaxonomyNode{
						// Level 3: Yirgacheffe Washed
						{
							ID:            "cat-coffee-yirgacheffe",
							Slug:          "yirgacheffe-specialty-coffee",
							Code:          "COFFEE_YIRGA",
							Level:         3,
							CanonicalName: "Yirgacheffe Specialty Floral Arabica (ይርጋጨፌ ቡና)",
							Definition:    "Highland washed Arabica coffee recognized for intense jasmine floral aroma and citrus bergamot acidity.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የይርጋጨፌ ጥራት ያለው ቡና",
								"or": "Buna Yirgaacaffee",
								"ti": "ይርጋጨፈ ቡን",
							},
						},
						// Level 3: Sidama Single Origin
						{
							ID:            "cat-coffee-sidama",
							Slug:          "sidama-specialty-coffee",
							Code:          "COFFEE_SIDAMA",
							Level:         3,
							CanonicalName: "Sidama Single-Origin Arabica Coffee (ሲዳማ ቡና)",
							Definition:    "Rich berry, stone-fruit, and cane-sugar sweetness profiles from the Sidama highlands.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የሲዳማ ጥራት ያለው ቡና",
								"or": "Buna Sidaamaa",
								"ti": "ሲዳማ ቡን",
							},
						},
					},
				},
			},
		},

		// 50. Root: Food, Groceries & Traditional Spices
		{
			ID:            "root-food-groceries",
			Slug:          "food-groceries-traditional-spices",
			Code:          "FOOD_GROCERIES",
			Level:         1,
			CanonicalName: "Food, Groceries & Traditional Spices",
			Definition:    "Culinary spices (Berbere, Shiro, Korarima, Black Cumin, Fenugreek), natural honey, refined cooking oils, and staple grocery provisions.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Utensils",
			Translations: map[string]string{
				"am": "የምግብ ሸቀጦች፣ ቅመማ ቅመም፣ በርበሬ እና ማር",
				"or": "Nyaata, Mi'eessituu, Barbaree fi Damma",
				"ti": "መግብን ቀመማት፣ በርበረን መዓርን",
			},
			Children: []TaxonomyNode{
				// Super Family 50.1: Traditional Spices & Seasoning Mixes
				{
					ID:            "sf-traditional-spices-berbere",
					Slug:          "traditional-spices-berbere-shiro",
					Code:          "SPICE_TRADITIONAL",
					Level:         2,
					CanonicalName: "Traditional Spices, Berbere & Shiro Powders (በርበሬ እና ሽሮ)",
					Definition:    "Sun-dried red pepper blends (Berbere), spiced roasted chickpea legume flour (Shiro), Korarima, and Mekelesha.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የወጥ በርበሬ፣ የሽሮ ዱቄት እና የቅመም ዓይነቶች",
						"or": "Barbaree fi Daakuu Shiroo",
						"ti": "በርበረ ጸብሒ፣ ዱቄት ሽሮን ቀመማትን",
					},
					Aliases: map[string][]string{
						"en": {"berbere", "shiro", "spices", "korarima", "ethiopian spice", "chili powder"},
						"am": {"በርበሬ", "ሽሮ", "የወጥ በርበሬ", "ኮረሪማ", "ቅመም", "ሚጥሚጣ"},
						"or": {"barbaree", "daakuu shiroo", "mi'eessituu"},
						"ti": {"በርበረ", "ሽሮ", "ቀመማት", "ኮረሪማ"},
					},
				},
				// Super Family 50.2: Pure Natural Honey (ማር)
				{
					ID:            "sf-natural-honey-bee-products",
					Slug:          "pure-natural-honey-bee-products",
					Code:          "FOOD_HONEY",
					Level:         2,
					CanonicalName: "Pure Natural Honey & Apiculture Products (ማር)",
					Definition:    "Raw white honey (Tigray, Gonder highlands), yellow multi-floral honey, and beeswax.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ንፁህ የተፈጥሮ ማር (ነጭ ማር፣ ቀይ ማር) እና ሰም",
						"or": "Damma Uumamaa fi Gaagura",
						"ti": "ንጹር ተፈጥሮኣዊ መዓር (ጻዕዳ መዓር)",
					},
					Aliases: map[string][]string{
						"en": {"honey", "white honey", "raw honey", "beeswax"},
						"am": {"ማር", "ነጭ ማር", "የተፈጥሮ ማር", "ሰም"},
						"or": {"damma", "damma adii"},
						"ti": {"መዓር", "ጻዕዳ መዓር"},
					},
				},
			},
		},
	}
}
