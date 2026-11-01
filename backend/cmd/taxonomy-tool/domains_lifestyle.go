package main

func getLifestyleDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 52. Root: Musical Instruments & Sound Gear
		{
			ID:            "root-musical-instruments",
			Slug:          "musical-instruments-sound-gear",
			Code:          "MUSIC_INSTRUMENTS",
			Level:         1,
			CanonicalName: "Musical Instruments & Audio Staging Gear",
			Definition:    "Traditional Ethiopian acoustic folk instruments (Kirar, Masinko, Kebero, Washint), modern electric synthesizers, keyboards, studio microphones, and stage audio.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Music",
			Translations: map[string]string{
				"am": "ሙዚቃ፣ የሙዚቃ መሳሪያዎች እና የድምጽ እቃዎች",
				"or": "Muuziqaa fi Meeshaalee Muuziqaa",
				"ti": "ሙዚቃ፣ ናይ ሙዚቃ መሳርሒታትን ናይ ድምጺ ንዋያትን",
			},
			Children: []TaxonomyNode{
				// Super Family 52.1: Traditional Ethiopian Acoustic Instruments
				{
					ID:            "sf-traditional-instruments",
					Slug:          "traditional-ethiopian-folk-instruments",
					Code:          "MUSIC_TRADITIONAL",
					Level:         2,
					CanonicalName: "Traditional Ethiopian Folk Instruments (ክራር፣ ማሲንቆ፣ ከበሮ)",
					Definition:    "Handcrafted five/six-string acoustic and electric Kirars, horsehair single-string Masinko fiddles, and leather-stretched church Kebero drums.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ባህላዊ የሙዚቃ መሳሪያዎች (ክራር፣ ማሲንቆ፣ ከበሮ፣ ዋሽንት)",
						"or": "Meeshaalee Muuziqaa Aadaa (Kiraara, Masinqoo)",
						"ti": "ባህላዊ መሳርሒታት ሙዚቃ (ክራር፣ ማሲንቆ፣ ከበሮ)",
					},
					Aliases: map[string][]string{
						"en": {"kirar", "krar", "masinko", "kebero", "washint", "traditional music"},
						"am": {"ክራር", "ማሲንቆ", "ከበሮ", "ዋሽንት", "የኤሌክትሪክ ክራር"},
						"or": {"kiraara", "masinqoo", "kabaaroo"},
						"ti": {"ክራር", "ማሲንቆ", "ከበሮ"},
					},
				},
				// Super Family 52.2: Modern Electronic Keyboards & Synths
				{
					ID:            "sf-electronic-keyboards",
					Slug:          "electronic-keyboards-synthesizers",
					Code:          "MUSIC_KEYBOARDS",
					Level:         2,
					CanonicalName: "Electronic Keyboards, Organs & Synthesizers (ኪቦርዶች)",
					Definition:    "Yamaha and Korg professional arranger workstations preloaded with Ethiopian oriental scales and rhythms.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የኤሌክትሮኒክስ ኪቦርዶች እና ኦርጋኖች (ያማሃ / ኮርግ)",
						"or": "Kiiboordii fi Oorgaanii",
						"ti": "ኪቦርዳትን ኦርጋናትን",
					},
				},
			},
		},

		// 53. Root: Arts, Crafts, Collectibles & Antiques
		{
			ID:            "root-arts-collectibles",
			Slug:          "arts-crafts-collectibles-antiques",
			Code:          "ARTS_COLLECTIBLES",
			Level:         1,
			CanonicalName: "Arts, Crafts, Collectibles & Antiques",
			Definition:    "Handmade processional brass and wooden crosses (Lalibela, Gondar, Axum), traditional parchment paintings, and historical numismatics.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Palette",
			Translations: map[string]string{
				"am": "የጥበብ ስራዎች፣ የእጅ ጥበብ እና ቅርሶች (መስቀሎች)",
				"or": "Hojii Aartii, Aadaa fi Qabeenya Seenaa",
				"ti": "ስነ ጥበብ፣ ናይ ኢድ ጥበባትን ቅርስታትን (መስቀላት)",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-traditional-crosses-artifacts",
					Slug:          "traditional-crosses-historical-artifacts",
					Code:          "ART_CROSSES",
					Level:         2,
					CanonicalName: "Handmade Processional Crosses & Cultural Artifacts (መስቀሎች)",
					Definition:    "Lost-wax cast brass, silver, and carved olive-wood hand, neck, and processional crosses.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የእጅ እና የብር መስቀሎች (የላሊበላ፣ የአክሱም መስቀል)",
						"or": "Fannoo Aadaa fi Qabeenya Seenaa",
						"ti": "ናይ ኢድን ብሩርን መስቀላት",
					},
				},
			},
		},

		// 54. Root: Pet Supplies & Companion Animals
		{
			ID:            "root-pet-supplies",
			Slug:          "pet-supplies-companion-animals",
			Code:          "PET_SUPPLIES",
			Level:         1,
			CanonicalName: "Pet Supplies, Veterinary Care & Companion Animals",
			Definition:    "Registered domestic dogs (German Shepherd, Rottweiler, Boerboel puppies), cat food, dog biscuits, and metal cages.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Dog",
			Translations: map[string]string{
				"am": "የቤት እንስሳት (ውሾች) እና የእንክብካቤ እቃዎች",
				"or": "Horii Manaa (Saree) fi Kunuunsa Isaanii",
				"ti": "ናይ ገዛ እንስሳታት (ኣኽላባት)ን መግብታትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-companion-dogs-puppies",
					Slug:          "pedigree-dogs-puppies",
					Code:          "PET_DOGS",
					Level:         2,
					CanonicalName: "Purebred Dogs, Guard Dogs & Puppies (ውሾች)",
					Definition:    "Vaccinated German Shepherd, Rottweiler, Caucasian Shepherd, and Boerboel guard puppies and trained adult canines.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የዘር ውሾች እና ቡችሎች (ጀርመን ሼፐርድ፣ ሮትዋይለር)",
						"or": "Ilmoo Saree fi Saree Eegumsaa",
						"ti": "ናይ ዘርኢ ኣኽላባትን ሒሉታትን",
					},
					Aliases: map[string][]string{
						"en": {"dog", "dogs", "puppy", "puppies", "german shepherd", "rottweiler"},
						"am": {"ውሻ", "ቡችላ", "የቤት ውሻ", "የጥበቃ ውሻ", "ጀርመን ሼፐርድ"},
						"or": {"saree", "ilmoo saree"},
						"ti": {"ከልቢ", "ሒሉ", "ናይ ዘርኢ ከልቢ"},
					},
				},
			},
		},
	}
}
