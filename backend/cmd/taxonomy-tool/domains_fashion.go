package main

func getFashionDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 36. Root: Traditional, Cultural & Ceremonial Attire
		{
			ID:            "root-traditional-attire",
			Slug:          "traditional-cultural-ceremonial-attire",
			Code:          "FASH_TRADITIONAL",
			Level:         1,
			CanonicalName: "Traditional, Cultural & Ceremonial Attire",
			Definition:    "Handwoven cotton cultural garments, traditional Habesha Kemis with intricate silk embroidery (Tilet), Gabi, Netela, Kuta, and regional cultural attire (Oromo, Amhara, Tigray, Somali, Afar).",
			InclusionRules: []string{"Handwoven Habesha Kemis (የሀበሻ ቀሚስ)", "Pure cotton Gabi (ጋቢ)", "Shemma handloom Netela (ነጠላ)", "Oromo traditional attire (Wandaboo, Qomee)", "Tigray traditional Tilfi garments"},
			ExclusionRules: []string{"Modern factory jeans and t-shirts (see Men's/Women's Clothing)"},
			Examples:       []string{"Handwoven Habesha Kemis with pure silk Tilet", "Four-layer warm highland Gabi", "Fine cotton Netela with gold thread border"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Sparkles",
			Translations: map[string]string{
				"am": "ባህላዊ እና የሀገር ባህል አልባሳት (የሀበሻ ቀሚስ፣ ጋቢ)",
				"or": "Uffata Aadaa fi Sirnaa",
				"ti": "ባህላውን ናይ ዓዲ ክዳውንትን (ናይ ሓበሻ ቀሚስ፣ ጋቢ)",
			},
			Aliases: map[string][]string{
				"en": {"traditional clothing", "habesha kemis", "cultural dress", "gabi", "netela", "shemma", "tilet"},
				"am": {"የሀበሻ ቀሚስ", "ባህላዊ ልብስ", "ጋቢ", "ነጠላ", "ኩታ", "ጥልፍ", "የሀገር ባህል ልብስ"},
				"or": {"uffata aadaa", "gaabii", "xilfii"},
				"ti": {"ናይ ሓበሻ ቀሚስ", "ባህላዊ ክዳን", "ጋቢ", "ነጸላ"},
			},
			Children: []TaxonomyNode{
				// Super Family 36.1: Handwoven Habesha Dresses (Habesha Kemis)
				{
					ID:            "sf-habesha-kemis-dresses",
					Slug:          "handwoven-habesha-kemis-dresses",
					Code:          "TRAD_KEMIS",
					Level:         2,
					CanonicalName: "Handwoven Habesha Kemis & Embroidered Dresses (የሀበሻ ቀሚስ)",
					Definition:    "Finely spun raw cotton dresses adorned with hand-stitched silk and golden Tilet embroidery across neckline, sleeves, and hemline.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "53102000", Name: "Traditional and cultural dresses"},
					},
					Translations: map[string]string{
						"am": "የእጅ ጥልፍ የሀበሻ ቀሚሶች",
						"or": "Wandaboo Aadaa fi Xilfii",
						"ti": "ናይ ኢድ ጥልፊ ናይ ሓበሻ ቀሚሳት",
					},
					Attributes: []TaxonomyAttribute{
						{Slug: "apparel-material", NameEN: "Fabric & Cotton Type", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Pure Handspun Cotton (የእጅ ፈትል)", "Semi-Machine Shemma", "Imported Sifon with Hand Tilet"}}},
						{Slug: "embroidery-style", NameEN: "Embroidery Style (Tilet)", DataType: "select", IsRequired: false, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Pure Silk Tilet (የሀር ጥልፍ)", "Golden Zari Thread", "Traditional Multi-color Pattern"}}},
					},
				},
				// Super Family 36.2: Traditional Gabi, Netela & Kuta Wraps
				{
					ID:            "sf-traditional-wraps-gabi",
					Slug:          "traditional-cotton-wraps-gabi-netela",
					Code:          "TRAD_WRAPS",
					Level:         2,
					CanonicalName: "Traditional Cotton Wraps (Gabi, Netela, Kuta)",
					Definition:    "Multi-ply warm highland Gabi wraps, lightweight sheer Netela shoulder shawls, and ceremonial men's Kuta wraps.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ጋቢ፣ ነጠላ እና የክብር ኩታ",
						"or": "Gaabii, Netelaa fi Kuutaa",
						"ti": "ጋቢ፣ ነጸላን ኩታን",
					},
				},
			},
		},

		// 37. Root: Footwear & Shoes
		{
			ID:            "root-footwear-shoes",
			Slug:          "footwear-shoes-boots",
			Code:          "FASH_FOOTWEAR",
			Level:         1,
			CanonicalName: "Footwear, Shoes & Leather Boots",
			Definition:    "Highland Ethiopian leather formal dress shoes, casual sneakers, sports running shoes, work safety boots, and sandals.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Footprints",
			Translations: map[string]string{
				"am": "ጫማዎች እና የቆዳ ጫማዎች",
				"or": "Kophee fi Qophee Gogaa",
				"ti": "ሳእኒታትን ናይ ቆርበት ሳእኒን",
			},
			Aliases: map[string][]string{
				"en": {"shoes", "footwear", "sneakers", "leather shoes", "boots", "sandals"},
				"am": {"ጫማ", "የቆዳ ጫማ", "ስኒከር", "ቡትስ", "የስፖርት ጫማ"},
				"or": {"kophee", "qophee gogaa", "sniikara"},
				"ti": {"ሳእኒ", "ናይ ቆርበት ሳእኒ", "ስኒከር"},
			},
			Children: []TaxonomyNode{
				// Super Family 37.1: Ethiopian Genuine Leather Formal Shoes
				{
					ID:            "sf-leather-formal-shoes",
					Slug:          "ethiopian-leather-formal-shoes",
					Code:          "SHOE_LEATHER",
					Level:         2,
					CanonicalName: "Ethiopian Genuine Leather Formal & Dress Shoes",
					Definition:    "Full-grain calfskin and cowhide oxford, derby, and loafer shoes manufactured by Ethiopian footwear artisans.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "53111501", Name: "Men's leather shoes"},
					},
					Translations: map[string]string{
						"am": "የሀበሻ ቆዳ ጫማዎች (ኦክስፎርድ፣ ሎፈር)",
						"or": "Kophee Gogaa Daldalaa",
						"ti": "ናይ ዓዲ ቆርበት ሳእኒ",
					},
					Attributes: []TaxonomyAttribute{
						{Slug: "shoe-size-eu", NameEN: "Shoe Size (EU)", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"38", "39", "40", "41", "42", "43", "44", "45", "46"}}},
						{Slug: "leather-type", NameEN: "Leather Grade", DataType: "select", IsRequired: false, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Full Grain Leather", "Top Grain Leather", "Suede Leather", "Nubuck"}}},
					},
				},
				// Super Family 37.2: Athletic Sneakers & Casual Trainers
				{
					ID:            "sf-sneakers-athletic-shoes",
					Slug:          "athletic-sneakers-running-shoes",
					Code:          "SHOE_SNEAKERS",
					Level:         2,
					CanonicalName: "Athletic Sneakers, Running Shoes & Casual Trainers",
					Definition:    "Cushioned long-distance running shoes, lifestyle sneakers, and basketball high-tops.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የስፖርት እና የስኒከር ጫማዎች",
						"or": "Kophee Ispoortii fi Sniikara",
						"ti": "ናይ ስፖርትን ስኒከርን ሳእኒ",
					},
				},
			},
		},

		// 38. Root: Fine Jewelry, Precious Metals & Gemstones
		{
			ID:            "root-jewelry-gemstones",
			Slug:          "fine-jewelry-precious-metals-gemstones",
			Code:          "FASH_JEWELRY",
			Level:         1,
			CanonicalName: "Fine Jewelry, Precious Metals & Gemstones",
			Definition:    "Certified 21k/24k Ethiopian gold necklaces, wedding bands, Ethiopian Wollo precious crystal opals, and 925 sterling silver.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Gem",
			Translations: map[string]string{
				"am": "ጌጣጌጥ፣ የወርቅ ስራዎች እና ውድ ማዕድናት (ኦፓል)",
				"or": "Faaya Warqee, Meetaa fi Dhagaa Bareedduu",
				"ti": "ስልማት፣ ናይ ወርቂ ስራሕትን ክቡራት ማዕድናትን",
			},
			Aliases: map[string][]string{
				"en": {"gold", "jewelry", "wollo opal", "opal", "silver", "wedding ring", "necklace"},
				"am": {"ወርቅ", "ጌጣጌጥ", "የወርቅ ጌጥ", "የወሎ ኦፓል", "ብር", "ቀለበት", "ሀብል"},
				"or": {"warqee", "faaya", "dhagaa bareedduu", "meetii"},
				"ti": {"ወርቂ", "ስልማት", "ናይ ወርቂ ስልማት", "ቀለቤት"},
			},
			Children: []TaxonomyNode{
				// Super Family 38.1: Gold Jewelry & Bridal Sets
				{
					ID:            "sf-gold-jewelry-sets",
					Slug:          "certified-gold-jewelry-bridal-sets",
					Code:          "JEWEL_GOLD",
					Level:         2,
					CanonicalName: "Certified Gold Jewelry & Bridal Sets (ወርቅ)",
					Definition:    "Hallmarked 18k, 21k, and 24k gold chains, rings, bangles, and bridal wedding jewelry sets.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "54111501", Name: "Gold jewelry"},
					},
					Translations: map[string]string{
						"am": "የወርቅ ጌጣጌጦች እና የሰርግ ስብስቦች",
						"or": "Faaya Warqee fi Qophii Cidhaa",
						"ti": "ናይ ወርቂ ስልማታትን ናይ መርዓ ንዋያትን",
					},
					Attributes: []TaxonomyAttribute{
						{Slug: "gold-purity-karat", NameEN: "Gold Purity (Karat)", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"24 Karat (99.9%)", "21 Karat (87.5%)", "18 Karat (75.0%)", "14 Karat"}}},
						{Slug: "weight-grams", NameEN: "Gold Weight (grams)", DataType: "number", Unit: "g", IsRequired: true, IsFilterable: true},
					},
				},
				// Super Family 38.2: Ethiopian Wollo Opals & Raw Gemstones
				{
					ID:            "sf-wollo-opals-gemstones",
					Slug:          "ethiopian-wollo-opals-gemstones",
					Code:          "JEWEL_OPALS",
					Level:         2,
					CanonicalName: "Ethiopian Wollo Hydrophane Opals & Loose Gems (ኦፓል)",
					Definition:    "World-renowned multi-color play-of-color hydrophane opals from Delanta/Wollo, rough and faceted gems.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የወሎ የተፈጥሮ ኦፓል እና የከበሩ ድንጋዮች",
						"or": "Dhagaa Bareedduu Ooppaala Walloo",
						"ti": "ናይ ወሎ ተፈጥሮኣዊ ኦፓል",
					},
				},
			},
		},

		// 39. Root: Luggage, Bags, Backpacks & Leather Goods
		{
			ID:            "root-luggage-bags",
			Slug:          "luggage-bags-backpacks-leather",
			Code:          "FASH_BAGS",
			Level:         1,
			CanonicalName: "Luggage, Bags, Backpacks & Leather Goods",
			Definition:    "Handcrafted genuine leather handbags, laptop briefcases, school backpacks, duffel bags, and rolling travel suitcases.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Briefcase",
			Translations: map[string]string{
				"am": "ቦርሳዎች፣ የቆዳ ቦርሳ እና የጉዞ ሻንጣዎች",
				"or": "Boorsaa fi Shaaxaa Imalaa",
				"ti": "ቦርሳታትን ናይ መገሻ ሻንጣታትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-leather-handbags-totes",
					Slug:          "genuine-leather-handbags-briefcases",
					Code:          "BAG_LEATHER",
					Level:         2,
					CanonicalName: "Genuine Ethiopian Leather Handbags & Briefcases",
					Definition:    "Premium full-grain leather tote bags, cross-body purses, and professional executive laptop briefcases.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የእጅ የቆዳ ቦርሳዎች እና የላፕቶፕ ቦርሳ",
						"or": "Boorsaa Gogaa Harkaa",
						"ti": "ናይ ቆርበት ቦርሳታት",
					},
				},
			},
		},
	}
}
