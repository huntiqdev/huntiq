package main

func getAgriDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 1. Root: Agriculture & Crop Farming
		{
			ID:            "root-agri-crops",
			Slug:          "agriculture-crop-farming",
			Code:          "AGRI_CROPS",
			Level:         1,
			CanonicalName: "Agriculture & Crop Farming",
			Definition:    "Commercial and smallholder cultivation of cereal grains, oilseeds, pulses, stimulants, and staple field crops.",
			InclusionRules: []string{"Raw and dried whole grains", "Legumes and pulse crops", "Oilseed harvests", "Cultivated field staple crops"},
			ExclusionRules: []string{"Prepared grocery foods or packaged flour (see Food & Groceries)", "Agricultural machinery and tractors (see Heavy Machinery)"},
			Examples:       []string{"Teff grain", "Wheat grain", "Sesame seed", "Faba beans", "Chickpeas", "Soybeans"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Wheat",
			Translations: map[string]string{
				"am": "ግብርና እና የእህል ሰብሎች ልማት",
				"or": "Qonnaa fi Oomisha Midhaanii",
				"ti": "ሕርሻን ናይ እኽሊ ሰብልታትን",
			},
			Aliases: map[string][]string{
				"en": {"agriculture", "crops", "farming", "grains", "harvest", "field crops"},
				"am": {"ግብርና", "እህል", "ሰብል", "የእህል ምርት", "የእርሻ ሰብል"},
				"or": {"qonnaa", "midhaan", "oomisha qonnaa", "sanyii midhaanii"},
				"ti": {"ሕርሻ", "እኽሊ", "ሰብሊ", "ፍርያት ሕርሻ"},
			},
			Children: []TaxonomyNode{
				// Super Family 1.1: Cereal Grains & Staple Crops
				{
					ID:            "sf-agri-cereals",
					Slug:          "cereal-grains-staples",
					Code:          "AGRI_CEREALS",
					Level:         2,
					CanonicalName: "Cereal Grains & Staple Crops",
					Definition:    "Major cereal grasses grown for edible starchy grain components.",
					InclusionRules: []string{"Teff", "Wheat", "Barley", "Maize", "Sorghum", "Millet", "Rice"},
					ExclusionRules: []string{"Processed breakfast cereals", "Wheat flour"},
					Status:         "active",
					IsLeaf:         false,
					Translations: map[string]string{
						"am": "የእህል እህሎች እና ዋና ሰብሎች",
						"or": "Midhaan Dhedhii fi Nyaata Bu'uuraa",
						"ti": "እኽልታትን መሰረታዊ ሰብልታትን",
					},
					Aliases: map[string][]string{
						"en": {"cereals", "staple grains", "cereal crops"},
						"am": {"የእህል ዓይነቶች", "ጥራጥሬዎች"},
						"or": {"gosoota midhaanii"},
						"ti": {"ዓይነታት እኽሊ"},
					},
					Children: []TaxonomyNode{
						// Level 3: Teff (Eragrostis tef)
						{
							ID:            "cat-cereal-teff",
							Slug:          "teff-grain",
							Code:          "CEREAL_TEFF",
							Level:         3,
							CanonicalName: "Teff Grain (ጤፍ)",
							Definition:    "Endemic Ethiopian cereal grain essential for injera preparation.",
							Status:         "active",
							IsLeaf:         false,
							Identifiers: []TaxonomyIdentifier{
								{Type: "TAXONOMIC", Code: "Eragrostis tef", Name: "Teff Botanical Species"},
								{Type: "UNSPSC", Code: "10151500", Name: "Raw cereals and grains"},
							},
							Translations: map[string]string{
								"am": "የጤፍ እህል",
								"or": "Midhaan Xaafii",
								"ti": "እኽሊ ጣፍ",
							},
							Aliases: map[string][]string{
								"en": {"teff", "tef", "eragrostis tef", "white teff", "brown teff"},
								"am": {"ጤፍ", "የጤፍ እህል", "ነጭ ጤፍ", "ሰርገኛ ጤፍ", "ቀይ ጤፍ", "ማግና ጤፍ"},
								"or": {"xaafii", "midhaan xaafii", "xaafii adii", "xaafii diimaa"},
								"ti": {"ጣፍ", "እኽሊ ጣፍ", "ጻዕዳ ጣፍ", "ቀይሕ ጣፍ"},
							},
							Attributes: []TaxonomyAttribute{
								{Slug: "crop-grade", NameEN: "Teff Grade", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Magna (Super White)", "Nech (White)", "Sergegna (Mixed)", "Key (Brown/Red)"}}},
								{Slug: "crop-origin-region", NameEN: "Harvest Region", DataType: "select", IsRequired: false, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Gojjam (Amhara)", "Shewa (Oromia)", "Wollo (Amhara)", "Tigray", "Other"}}},
								{Slug: "moisture-content-percent", NameEN: "Moisture Content (%)", DataType: "number", Unit: "%", IsRequired: false, IsFilterable: true},
							},
							Children: []TaxonomyNode{
								// Level 4: Magna Teff
								{
									ID:            "leaf-teff-magna",
									Slug:          "magna-super-white-teff",
									Code:          "TEFF_MAGNA",
									Level:         4,
									CanonicalName: "Magna Super-White Teff Grain",
									Definition:    "Premium grade bright-white teff grain cultivated in fertile highlands (East Gojjam, Ada'a/Bishoftu).",
									Status:         "active",
									IsLeaf:         true,
									Translations: map[string]string{
										"am": "ማግና ሱፐር ነጭ ጤፍ",
										"or": "Xaafii Adii Magnaa",
										"ti": "ማግና ጻዕዳ ጣፍ",
									},
								},
								// Level 4: Sergegna Teff
								{
									ID:            "leaf-teff-sergegna",
									Slug:          "sergegna-mixed-teff",
									Code:          "TEFF_SERGEGNA",
									Level:         4,
									CanonicalName: "Sergegna Mixed Red & White Teff Grain",
									Definition:    "Nutrient-dense blended white and red grain teff offering optimal traditional sour-dough fermentation.",
									Status:         "active",
									IsLeaf:         true,
									Translations: map[string]string{
										"am": "ሰርገኛ ቅይጥ ጤፍ",
										"or": "Xaafii Makamaa Sargeenyaa",
										"ti": "ሰርገኛ ሕውስዋስ ጣፍ",
									},
								},
								// Level 4: Red/Brown Teff
								{
									ID:            "leaf-teff-key",
									Slug:          "key-brown-red-teff",
									Code:          "TEFF_KEY",
									Level:         4,
									CanonicalName: "Key Brown & Red Iron-Rich Teff Grain",
									Definition:    "High-iron content red and brown varieties cultivated across diverse microclimates.",
									Status:         "active",
									IsLeaf:         true,
									Translations: map[string]string{
										"am": "ቀይ / ቡናማ ጤፍ",
										"or": "Xaafii Diimaa",
										"ti": "ቀይሕ ጣፍ",
									},
								},
							},
						},
						// Level 3: Wheat Grain (Triticum)
						{
							ID:            "cat-cereal-wheat",
							Slug:          "wheat-grain",
							Code:          "CEREAL_WHEAT",
							Level:         3,
							CanonicalName: "Wheat Grain (ስንዴ)",
							Definition:    "Bread wheat and durum wheat grains for milling, pasta production, and food manufacture.",
							Status:         "active",
							IsLeaf:         true,
							Identifiers: []TaxonomyIdentifier{
								{Type: "TAXONOMIC", Code: "Triticum aestivum", Name: "Bread Wheat"},
								{Type: "UNSPSC", Code: "10151501", Name: "Raw wheat"},
							},
							Translations: map[string]string{
								"am": "የስንዴ እህል",
								"or": "Midhaan Qamadii",
								"ti": "እኽሊ ስርናይ",
							},
							Aliases: map[string][]string{
								"en": {"wheat", "raw wheat", "durum wheat", "bread wheat"},
								"am": {"ስንዴ", "የስንዴ እህል", "የዳቦ ስንዴ"},
								"or": {"qamadii", "midhaan qamadii"},
								"ti": {"ስርናይ", "እኽሊ ስርናይ"},
							},
						},
						// Level 3: Barley Grain (Hordeum vulgare)
						{
							ID:            "cat-cereal-barley",
							Slug:          "barley-grain",
							Code:          "CEREAL_BARLEY",
							Level:         3,
							CanonicalName: "Barley Grain (ገብስ)",
							Definition:    "Food barley, malting barley for brewing, and feed barley.",
							Status:         "active",
							IsLeaf:         true,
							Identifiers: []TaxonomyIdentifier{
								{Type: "TAXONOMIC", Code: "Hordeum vulgare", Name: "Barley"},
								{Type: "UNSPSC", Code: "10151502", Name: "Raw barley"},
							},
							Translations: map[string]string{
								"am": "የገብስ እህል",
								"or": "Midhaan Garbuu",
								"ti": "እኽሊ ስገም",
							},
							Aliases: map[string][]string{
								"en": {"barley", "malt barley", "food barley"},
								"am": {"ገብስ", "የገብስ እህል", "የቢራ ገብስ"},
								"or": {"garbuu", "midhaan garbuu"},
								"ti": {"ስገም", "እኽሊ ስገም"},
							},
						},
						// Level 3: Maize / Corn (Zea mays)
						{
							ID:            "cat-cereal-maize",
							Slug:          "maize-corn-grain",
							Code:          "CEREAL_MAIZE",
							Level:         3,
							CanonicalName: "Maize & Corn Grain (በቆሎ)",
							Definition:    "Field corn, white maize, and yellow maize for human consumption, starch extraction, and animal feeds.",
							Status:         "active",
							IsLeaf:         true,
							Identifiers: []TaxonomyIdentifier{
								{Type: "TAXONOMIC", Code: "Zea mays", Name: "Maize"},
								{Type: "UNSPSC", Code: "10151503", Name: "Raw corn"},
							},
							Translations: map[string]string{
								"am": "የበቆሎ እህል",
								"or": "Boqqolloo",
								"ti": "እኽሊ ኣፉን",
							},
							Aliases: map[string][]string{
								"en": {"maize", "corn", "grain corn", "white maize"},
								"am": {"በቆሎ", "የበቆሎ እህል", "ነጭ በቆሎ"},
								"or": {"boqqolloo", "midhaan boqqolloo"},
								"ti": {"ኣፉን", "ዕፉን", "እኽሊ ኣፉን"},
							},
						},
						// Level 3: Sorghum Grain (Sorghum bicolor)
						{
							ID:            "cat-cereal-sorghum",
							Slug:          "sorghum-grain",
							Code:          "CEREAL_SORGHUM",
							Level:         3,
							CanonicalName: "Sorghum Grain (ማሽላ)",
							Definition:    "Drought-tolerant sorghum grain for flatbreads, porridge, and traditional beverages.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የማሽላ እህል",
								"or": "Mishingaa",
								"ti": "እኽሊ መሸላ",
							},
							Aliases: map[string][]string{
								"en": {"sorghum", "milo", "guinea corn"},
								"am": {"ማሽላ", "የማሽላ እህል"},
								"or": {"mishingaa"},
								"ti": {"መሸላ", "እኽሊ መሸላ"},
							},
						},
					},
				},
				// Super Family 1.2: Oilseeds & Industrial Crops
				{
					ID:            "sf-agri-oilseeds",
					Slug:          "oilseeds-industrial-crops",
					Code:          "AGRI_OILSEEDS",
					Level:         2,
					CanonicalName: "Oilseeds & Industrial Export Crops",
					Definition:    "Seeds cultivated primarily for vegetable edible oil extraction and export trade.",
					InclusionRules: []string{"Sesame seed", "Niger seed (Noug)", "Linseed/Flaxseed", "Soybeans", "Sunflower seeds", "Castor beans"},
					ExclusionRules: []string{"Extracted and refined cooking oils (see Food & Groceries)"},
					Status:         "active",
					IsLeaf:         false,
					Translations: map[string]string{
						"am": "የቅባት እህሎች እና የኢንዱስትሪ ሰብሎች",
						"or": "Sanyii Zayitaa fi Oomisha Daldalaa",
						"ti": "ናይ ቅብኣት እኽልታትን ንግዳዊ ሰብልታትን",
					},
					Aliases: map[string][]string{
						"en": {"oilseeds", "sesame", "noug", "flaxseed", "linseed"},
						"am": {"የቅባት እህል", "ሰሊጥ", "ኑግ", "ተልባ", "ሱፍ"},
						"or": {"sanyii zayitaa", "saliixa", "nuugii"},
						"ti": {"ናይ ቅብኣት እኽሊ", "ሰሊጥ", "ኑግ", "ጣሊያን"},
					},
					Children: []TaxonomyNode{
						// Level 3: Sesame Seed (Sesamum indicum)
						{
							ID:            "cat-oilseed-sesame",
							Slug:          "sesame-seed",
							Code:          "OILSEED_SESAME",
							Level:         3,
							CanonicalName: "Sesame Seed (ሰሊጥ)",
							Definition:    "White Humera, Wollega, and Metema export sesame seeds.",
							Status:         "active",
							IsLeaf:         true,
							Identifiers: []TaxonomyIdentifier{
								{Type: "TAXONOMIC", Code: "Sesamum indicum", Name: "Sesame"},
								{Type: "UNSPSC", Code: "10151515", Name: "Sesame seed"},
							},
							Translations: map[string]string{
								"am": "የሰሊጥ ፍሬ",
								"or": "Sanyii Saliixaa",
								"ti": "ፍረ ሰሊጥ",
							},
							Aliases: map[string][]string{
								"en": {"sesame", "sesame seed", "humera sesame", "wollega sesame"},
								"am": {"ሰሊጥ", "የሰሊጥ ፍሬ", "ሁመራ ሰሊጥ", "ወለጋ ሰሊጥ"},
								"or": {"saliixa", "sanyii saliixaa"},
								"ti": {"ሰሊጥ", "ፍረ ሰሊጥ", "ሁመራ ሰሊጥ"},
							},
						},
						// Level 3: Niger Seed / Noug (Guizotia abyssinica)
						{
							ID:            "cat-oilseed-noug",
							Slug:          "niger-seed-noug",
							Code:          "OILSEED_NOUG",
							Level:         3,
							CanonicalName: "Niger Seed / Noug (ኑግ)",
							Definition:    "Indigenous Ethiopian high-oil content seed used for domestic culinary oil processing.",
							Status:         "active",
							IsLeaf:         true,
							Identifiers: []TaxonomyIdentifier{
								{Type: "TAXONOMIC", Code: "Guizotia abyssinica", Name: "Niger seed"},
							},
							Translations: map[string]string{
								"am": "የኑግ ፍሬ",
								"or": "Sanyii Nuugii",
								"ti": "ፍረ ኑግ",
							},
							Aliases: map[string][]string{
								"en": {"noug", "niger seed", "neug"},
								"am": {"ኑግ", "የኑግ ፍሬ", "ንፁህ ኑግ"},
								"or": {"nuugii", "sanyii nuugii"},
								"ti": {"ኑግ", "ፍረ ኑግ"},
							},
						},
						// Level 3: Linseed / Flaxseed (Linum usitatissimum)
						{
							ID:            "cat-oilseed-flaxseed",
							Slug:          "linseed-flaxseed",
							Code:          "OILSEED_FLAX",
							Level:         3,
							CanonicalName: "Linseed & Flaxseed (ተልባ)",
							Definition:    "Brown and golden flaxseed harvested for food beverage preparation and omega-3 oil pressing.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የተልባ ፍሬ",
								"or": "Talbaa",
								"ti": "ፍረ እንጣጢዕ",
							},
							Aliases: map[string][]string{
								"en": {"flaxseed", "linseed", "flax"},
								"am": {"ተልባ", "የተልባ ፍሬ"},
								"or": {"talbaa", "sanyii talbaa"},
								"ti": {"እንጣጢዕ", "ፍረ እንጣጢዕ"},
							},
						},
						// Level 3: Soybeans (Glycine max)
						{
							ID:            "cat-oilseed-soybean",
							Slug:          "soybeans-grain",
							Code:          "OILSEED_SOY",
							Level:         3,
							CanonicalName: "Soybeans (አኩሪ አተር)",
							Definition:    "Commercial edible legumes cultivated for industrial soy processing, lecithin, and protein feed.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "አኩሪ አተር",
								"or": "Soyaabinii",
								"ti": "ኣኩሪ ኣተር",
							},
						},
					},
				},
				// Super Family 1.3: Pulses & Legumes
				{
					ID:            "sf-agri-pulses",
					Slug:          "pulses-legumes-beans",
					Code:          "AGRI_PULSES",
					Level:         2,
					CanonicalName: "Pulses, Legumes & Dried Beans",
					Definition:    "High-protein dried legume seeds harvested from pods for human food staples.",
					InclusionRules: []string{"Faba beans (Bakela)", "Chickpeas (Shimbra)", "Lentils (Misir)", "Field peas (Ater)", "Haricot beans (Boleqe)"},
					ExclusionRules: []string{"Cooked stews or packaged shiro powder (see Food & Groceries)"},
					Status:         "active",
					IsLeaf:         false,
					Translations: map[string]string{
						"am": "የጥራጥሬ ሰብሎች እና አተር/ምስር",
						"or": "Oomisha Mi'eessituu fi Boqqolloo",
						"ti": "ጥራጥሬታትን ዓተርን ምስርን",
					},
					Children: []TaxonomyNode{
						// Level 3: Lentils (Lens culinaris)
						{
							ID:            "cat-pulse-lentils",
							Slug:          "dried-lentils",
							Code:          "PULSE_LENTIL",
							Level:         3,
							CanonicalName: "Dried Lentils (ምስር)",
							Definition:    "Red, brown, and green whole and split dried lentils.",
							Status:         "active",
							IsLeaf:         true,
							Identifiers: []TaxonomyIdentifier{
								{Type: "TAXONOMIC", Code: "Lens culinaris", Name: "Lentil"},
							},
							Translations: map[string]string{
								"am": "የምስር ክክ / ፍሬ",
								"or": "Misira",
								"ti": "ምስሪ",
							},
							Aliases: map[string][]string{
								"en": {"lentils", "red lentils", "brown lentils", "misir"},
								"am": {"ምስር", "የምስር ክክ", "ቀይ ምስር"},
								"or": {"misira", "misira diimaa"},
								"ti": {"ምስሪ", "ቀይሕ ምስሪ"},
							},
						},
						// Level 3: Chickpeas (Cicer arietinum)
						{
							ID:            "cat-pulse-chickpeas",
							Slug:          "chickpeas-shimbra",
							Code:          "PULSE_CHICKPEA",
							Level:         3,
							CanonicalName: "Chickpeas (ሽምብራ)",
							Definition:    "Desi and Kabuli whole dried chickpeas used for roasting, boiling, and shiro base.",
							Status:         "active",
							IsLeaf:         true,
							Identifiers: []TaxonomyIdentifier{
								{Type: "TAXONOMIC", Code: "Cicer arietinum", Name: "Chickpea"},
							},
							Translations: map[string]string{
								"am": "የሽምብራ ፍሬ",
								"or": "Shumburaa",
								"ti": "ዓተር ጓጓ",
							},
							Aliases: map[string][]string{
								"en": {"chickpeas", "garbanzo beans", "shimbra"},
								"am": {"ሽምብራ", "የሽምብራ ፍሬ"},
								"or": {"shumburaa"},
								"ti": {"ዓተር ጓጓ", "ሽምብራ"},
							},
						},
						// Level 3: Faba Beans / Broad Beans (Vicia faba)
						{
							ID:            "cat-pulse-faba-beans",
							Slug:          "faba-beans-bakela",
							Code:          "PULSE_FABA",
							Level:         3,
							CanonicalName: "Faba Beans / Broad Beans (ባቄላ)",
							Definition:    "Dried large faba beans for foul preparation and food staples.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የባቄላ ፍሬ",
								"or": "Baqqelaa",
								"ti": "ባልዶንጓ",
							},
							Aliases: map[string][]string{
								"en": {"faba beans", "broad beans", "bakela"},
								"am": {"ባቄላ", "የባቄላ ፍሬ"},
								"or": {"baqqelaa"},
								"ti": {"ባልዶንጓ", "ባቄላ"},
							},
						},
					},
				},
			},
		},

		// 2. Root: Horticulture, Vegetables & Floriculture
		{
			ID:            "root-horticulture",
			Slug:          "horticulture-vegetables-floriculture",
			Code:          "HORT_VEG_FLOR",
			Level:         1,
			CanonicalName: "Horticulture, Vegetables & Floriculture",
			Definition:    "Cultivation and harvesting of fresh vegetables, root tubers, orchard fruits, cut roses, and greenhouse flora.",
			InclusionRules: []string{"Fresh onions", "Potatoes", "Tomatoes", "Garlic", "Avocados", "Cut export roses", "Nursery seedlings"},
			ExclusionRules: []string{"Cooked salads or canned vegetables (see Food & Groceries)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Flower2",
			Translations: map[string]string{
				"am": "ሆርቲካልቸር፣ አትክልት እና የአበባ ልማት",
				"or": "Oomisha Kuduraa fi Biqiltootaa",
				"ti": "ሆርቲካልቸር፣ ኣሕምልትን ዕምባባን",
			},
			Aliases: map[string][]string{
				"en": {"horticulture", "vegetables", "fruits", "flowers", "roses", "greenhouse"},
				"am": {"አትክልት", "ፍራፍሬ", "አበባ", "የአበባ ልማት", "የጓሮ አትክልት"},
				"or": {"kuduraa", "muduraa", "abaaboo", "biqiltuu"},
				"ti": {"ኣሕምልቲ", "ፍረታት", "ዕምባባ", "ናይ ገዛ ኣሕምልቲ"},
			},
			Children: []TaxonomyNode{
				// Super Family 2.1: Fresh Vegetables & Alliums
				{
					ID:            "sf-hort-vegetables",
					Slug:          "fresh-vegetables-alliums",
					Code:          "HORT_VEGETABLES",
					Level:         2,
					CanonicalName: "Fresh Vegetables & Allium Crops",
					Definition:    "Bulb vegetables, leafy greens, brassicas, and solanaceous garden vegetables.",
					InclusionRules: []string{"Onions", "Garlic", "Shallots", "Cabbage", "Kale/Gomen", "Tomatoes", "Peppers"},
					Status:         "active",
					IsLeaf:         false,
					Translations: map[string]string{
						"am": "ትኩስ አትክልቶች እና የሽንኩርት ዘሮች",
						"or": "Kuduraalee Dhedhii",
						"ti": "ሓደሽቲ ኣሕምልትን ሽጉርትን",
					},
					Children: []TaxonomyNode{
						// Level 3: Onions & Alliums (Allium cepa)
						{
							ID:            "cat-veg-onions",
							Slug:          "onions-alliums",
							Code:          "VEG_ONION",
							Level:         3,
							CanonicalName: "Onions & Bulb Alliums (ሽንኩርት)",
							Definition:    "Red Bombay onions, yellow onions, shallots, and spring onions harvested for market trade.",
							Status:         "active",
							IsLeaf:         false,
							Identifiers: []TaxonomyIdentifier{
								{Type: "TAXONOMIC", Code: "Allium cepa", Name: "Onion"},
								{Type: "UNSPSC", Code: "50101504", Name: "Fresh onions"},
								{Type: "GS1_GPC", Code: "10005845", Name: "Vegetables - Fresh Unprocessed"},
							},
							Translations: map[string]string{
								"am": "ቀይ እና ነጭ ሽንኩርት",
								"or": "Qullubbii Dimaa fi Adii",
								"ti": "ቀይሕን ጻዕዳን ሽጉርቲ",
							},
							Aliases: map[string][]string{
								"en": {"onion", "onions", "red onion", "shallot", "bulb onion"},
								"am": {"ሽንኩርት", "ቀይ ሽንኩርት", "የሀበሻ ሽንኩርት", "ነጭ ሽንኩርት"},
								"or": {"qullubbii", "qullubbii dimaa", "qullubbii adii"},
								"ti": {"ሽጉርቲ", "ቀይሕ ሽጉርቲ", "ጻዕዳ ሽጉርቲ"},
							},
							Attributes: []TaxonomyAttribute{
								{Slug: "crop-origin-region", NameEN: "Production Region", DataType: "select", IsRequired: false, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Rift Valley (Meki/Ziway)", "Amhara (Kobo/Fogera)", "Tigray", "Oromia", "Other"}}},
								{Slug: "packaging-unit", NameEN: "Trade Unit", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Quintal (100kg bag)", "Net Bag (25kg)", "Net Bag (10kg)", "Kilogram (kg)"}}},
							},
							Children: []TaxonomyNode{
								// Level 4: Red Bombay Dry Bulb Onion
								{
									ID:            "leaf-onion-red-bombay",
									Slug:          "red-bombay-dry-bulb-onion",
									Code:          "ONION_RED_BOMBAY",
									Level:         4,
									CanonicalName: "Red Bombay Dry Bulb Onion",
									Definition:    "Standard commercial pungent red bulb onion cultivated extensively in Ziway, Meki, and Awash basin.",
									Status:         "active",
									IsLeaf:         true,
									Translations: map[string]string{
										"am": "ቀይ ቦምቤይ ደረቅ ሽንኩርት",
										"or": "Qullubbii Dimaa Boombeeyii",
										"ti": "ቀይሕ ቦምበይ ደረቕ ሽጉርቲ",
									},
								},
								// Level 4: Garlic Bulbs (Allium sativum)
								{
									ID:            "leaf-garlic-bulbs",
									Slug:          "fresh-garlic-bulbs",
									Code:          "ALLIUM_GARLIC",
									Level:         4,
									CanonicalName: "Fresh Garlic Bulbs (ነጭ ሽንኩርት)",
									Definition:    "Whole cured garlic bulbs cultivated in highland regions.",
									Status:         "active",
									IsLeaf:         true,
									Identifiers: []TaxonomyIdentifier{
										{Type: "TAXONOMIC", Code: "Allium sativum", Name: "Garlic"},
									},
									Translations: map[string]string{
										"am": "ነጭ ሽንኩርት (ነጭ እጓል)",
										"or": "Qullubbii Adii",
										"ti": "ጻዕዳ ሽጉርቲ",
									},
								},
							},
						},
						// Level 3: Solanaceous Vegetables (Tomato & Peppers)
						{
							ID:            "cat-veg-solanaceous",
							Slug:          "tomatoes-peppers-potatoes",
							Code:          "VEG_SOLANUM",
							Level:         3,
							CanonicalName: "Tomatoes, Peppers & Irish Potatoes",
							Definition:    "Fresh vine-ripened tomatoes, green chili peppers, red chili, and Irish potatoes.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "ቲማቲም፣ ቃሪያ እና ድንች",
								"or": "Timaatima, Qaraa fi Dinnicha",
								"ti": "ኮሚደረ፣ በርበረን ድንሽ",
							},
							Aliases: map[string][]string{
								"en": {"tomatoes", "peppers", "potatoes", "chili"},
								"am": {"ቲማቲም", "ድንች", "ቃሪያ", "ሚጥሚጣ"},
								"or": {"timaatima", "dinnicha", "qaraa"},
								"ti": {"ኮሚደረ", "ድንሽ", "በርበረ"},
							},
						},
					},
				},
				// Super Family 2.2: Cut Flowers & Floriculture
				{
					ID:            "sf-hort-floriculture",
					Slug:          "floriculture-cut-flowers",
					Code:          "HORT_FLOWERS",
					Level:         2,
					CanonicalName: "Floriculture, Cut Flowers & Greenhouse Flora",
					Definition:    "Export and domestic greenhouse cut roses, carnations, gypsophila, and decorative ornamental plants.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የአበባ ልማት እና የተቆረጡ አበቦች",
						"or": "Misooma Abaaboo fi Biqiltoota Miidhaginaa",
						"ti": "ምልማዕ ዕምባባታትን ጌጻዊ ተኽልታትን",
					},
				},
			},
		},

		// 3. Root: Livestock, Pastoral Farming & Animal Husbandry
		{
			ID:            "root-livestock",
			Slug:          "livestock-animal-husbandry",
			Code:          "LIVE_STOCK",
			Level:         1,
			CanonicalName: "Livestock, Pastoral Farming & Animal Husbandry",
			Definition:    "Live domestic ruminants, cattle, sheep, goats, camels, and equines for dairy, beef, breeding, and draft work.",
			InclusionRules: []string{"Live cattle and bulls", "Dairy cows", "Sheep and rams", "Goats", "Camels", "Equines (horses, donkeys)"},
			ExclusionRules: []string{"Prepared cuts of meat and butcher beef (see Food & Groceries)", "Animal feed and fodder (see Animal Feed)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Footprints",
			Translations: map[string]string{
				"am": "የቁም እንስሳት እና የእንስሳት እርባታ",
				"or": "Beelladaa fi Horsiisa Horii",
				"ti": "ጥሪትን ናይ እንስሳት መራብሕን",
			},
			Aliases: map[string][]string{
				"en": {"livestock", "cattle", "cows", "sheep", "goats", "camels", "bulls"},
				"am": {"የቁም እንስሳት", "ከብቶች", "ላሞች", "በሬዎች", "በጎች", "ፍየሎች", "ግመሎች"},
				"or": {"beellada", "horii", "sa'a", "sangaa", "hoolaa", "ree'ee", "gaala"},
				"ti": {"ጥሪት", "ከብቲ", "ላም", "ብዕራይ", "በጊዕ", "ጤል", "ግመል"},
			},
			Children: []TaxonomyNode{
				// Super Family 3.1: Cattle & Bovines (Bos taurus / Bos indicus)
				{
					ID:            "sf-live-cattle",
					Slug:          "cattle-bovines",
					Code:          "LIVE_CATTLE",
					Level:         2,
					CanonicalName: "Cattle & Bovines (የከብት እርባታ)",
					Definition:    "Boran beef cattle, Holstein-Friesian dairy crosses, and indigenous highland zebu cattle.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "TAXONOMIC", Code: "Bos taurus / Bos indicus", Name: "Cattle"},
					},
					Translations: map[string]string{
						"am": "የቀንድ ከብቶች (ላሞች እና በሬዎች)",
						"or": "Loon (Sa'a fi Sangaa)",
						"ti": "ከብቲ (ላምን ብዕራይን)",
					},
					Attributes: []TaxonomyAttribute{
						{Slug: "livestock-breed", NameEN: "Cattle Breed", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Boran", "Holstein Friesian Cross", "Jersey Cross", "Highland Zebu", "Fogera", "Barka"}}},
						{Slug: "livestock-purpose", NameEN: "Purpose", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Dairy Production", "Fattening & Beef", "Breeding Stock", "Draft / Ploughing"}}},
					},
				},
				// Super Family 3.2: Small Ruminants (Sheep & Goats)
				{
					ID:            "sf-live-sheep-goats",
					Slug:          "sheep-goats-small-ruminants",
					Code:          "LIVE_SHEEP_GOAT",
					Level:         2,
					CanonicalName: "Sheep & Goats (በጎች እና ፍየሎች)",
					Definition:    "Hararghe highland sheep, Menz black sheep, Bonga sheep, and Boer cross goats.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "በጎች እና ፍየሎች",
						"or": "Hoolota fi Re'oota",
						"ti": "በጊዕን ጤልን",
					},
				},
			},
		},

		// 4. Root: Poultry & Aviculture
		{
			ID:            "root-poultry",
			Slug:          "poultry-aviculture",
			Code:          "POULTRY_AVIC",
			Level:         1,
			CanonicalName: "Poultry & Aviculture",
			Definition:    "Domestic fowls, layer chickens, broiler meat birds, day-old chicks, turkeys, and hatching eggs.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Egg",
			Translations: map[string]string{
				"am": "የዶሮ እርባታ እና አእዋፋት",
				"or": "Horsiisa Lukkuu",
				"ti": "ናይ ደርሆ መራብሕን ኣዕዋፍን",
			},
			Children: []TaxonomyNode{
				// Super Family 4.1: Chickens & Poultry Birds
				{
					ID:            "sf-poultry-chickens",
					Slug:          "chickens-layers-broilers",
					Code:          "POULTRY_CHICKEN",
					Level:         2,
					CanonicalName: "Commercial & Backyard Chickens (ዶሮዎች)",
					Definition:    "Day-old chicks, Sasso, Bovans Brown, and Cobb 500 layers and broilers.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የእንቁላል እና የስጋ ዶሮዎች",
						"or": "Lukkuu Buuphaa fi Foonii",
						"ti": "ናይ እንቋቑሖን ስጋን ደራሁ",
					},
				},
			},
		},

		// 5. Root: Fisheries & Aquaculture
		{
			ID:            "root-fisheries",
			Slug:          "fisheries-aquaculture",
			Code:          "FISH_AQUA",
			Level:         1,
			CanonicalName: "Fisheries & Aquaculture",
			Definition:    "Freshwater commercial fish catches, lake tilapia, African catfish, Nile perch, fingerling seed stock, and hatchery gear.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Fish",
			Translations: map[string]string{
				"am": "ዓሳ ሀብት እና የውሃ ህይወት እርባታ",
				"or": "Misooma Qurxummii fi Bishaanii",
				"ti": "ናይ ዓሳ ሃብትን ሕርሻ ማይን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-freshwater-fish",
					Slug:          "freshwater-commercial-fish",
					Code:          "FISH_FRESHWATER",
					Level:         2,
					CanonicalName: "Freshwater Whole Fish (ትኩስ ዓሳ)",
					Definition:    "Lake Tana, Ziway, and Hawassa Tilapia and Nile Perch catches.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ትኩስ የሐይቅ ዓሳዎች",
						"or": "Qurxummii Haroo",
						"ti": "ሓድሽ ናይ ቀላይ ዓሳ",
					},
				},
			},
		},

		// 6. Root: Forestry, Logging & Timber
		{
			ID:            "root-forestry",
			Slug:          "forestry-logging-timber",
			Code:          "FORESTRY_TIMBER",
			Level:         1,
			CanonicalName: "Forestry, Logging & Raw Timber",
			Definition:    "Commercial timber plantations, eucalyptus poles, raw sawn timber logs, bamboo culms, and certified wood harvests.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Trees",
			Translations: map[string]string{
				"am": "ደን ልማት፣ የዛፍ ቁራጮች እና ግንድ",
				"or": "Misooma Bosonaa fi Muka",
				"ti": "ልምዓት ደንን ዕንጨይትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-raw-timber-poles",
					Slug:          "raw-timber-poles-logs",
					Code:          "FORESTRY_LOGS",
					Level:         2,
					CanonicalName: "Eucalyptus Poles, Bamboo & Sawmill Logs",
					Definition:    "Construction poles, highland bamboo culms, and unprocessed raw hardwood logs.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የባህር ዛፍ ቋሚዎች፣ የቀርከሃ ዛፍ እና ግንድ",
						"or": "Utubaa Baargamoo fi Mukkeen Ijaarsaa",
						"ti": "ናይ ባሕሪ ዛፍ ኣዕኑድን ቃንጫን",
					},
				},
			},
		},
	}
}
