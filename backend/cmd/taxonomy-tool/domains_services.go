package main

func getServicesDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 55. Root: Consumer Electronics & Appliance Repair Services
		{
			ID:            "root-electronics-repair-services",
			Slug:          "electronics-appliance-repair-services",
			Code:          "SERV_REPAIR_ELEC",
			Level:         1,
			CanonicalName: "Consumer Electronics & Appliance Repair Services",
			Definition:    "Professional diagnostic, maintenance, screen replacement, and hardware repair services for smartphones, laptops, TVs, and home appliances.",
			InclusionRules: []string{"Smartphone screen fitting and glass replacement labor", "Laptop motherboard chip-level repair", "Refrigerator gas refilling and compressor replacement", "TV panel troubleshooting"},
			ExclusionRules: []string{"Physical replacement screens and display parts without labor (see Consumer Electronics)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Wrench",
			Translations: map[string]string{
				"am": "የኤሌክትሮኒክስ፣ የስልክ እና የቤት እቃዎች ጥገና አገልግሎት",
				"or": "Tajaajila Suphaa Elektirooniksii fi Bilbilaa",
				"ti": "ናይ ኤሌክትሮኒክስን ስልኪን ጽገና ኣገልግሎት",
			},
			Aliases: map[string][]string{
				"en": {"phone repair", "screen replacement", "laptop repair", "appliance repair", "electronics repair service"},
				"am": {"ስልክ ጥገና", "የስልክ ጥገና", "ስክሪን ቅየራ", "ላፕቶፕ ጥገና", "ፍሪጅ ጥገና"},
				"or": {"suphaa bilbilaa", "suphaa kompiitaraa"},
				"ti": {"ጽገና ስልኪ", "ምቕያር ስክሪን", "ጽገና ላፕቶፕ"},
			},
			Children: []TaxonomyNode{
				// Super Family 55.1: Smartphone Screen & Battery Replacement Services
				{
					ID:            "sf-phone-repair-labor",
					Slug:          "smartphone-screen-hardware-repair-services",
					Code:          "SERV_PHONE_FIX",
					Level:         2,
					CanonicalName: "Smartphone Screen & Hardware Repair Services (የስልክ ጥገና)",
					Definition:    "Certified bench technician labor for cracked touchscreen digitizer replacement, IC chip micro-soldering, and charging port repair.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "81112403", Name: "Computer and telephone hardware repair service"},
					},
					Translations: map[string]string{
						"am": "የስማርትፎን ስክሪን ቅየራ እና የሃርድዌር ጥገና አገልግሎት",
						"or": "Tajaajila Suphaa Iskiriinii fi Haardiweerii Bilbilaa",
						"ti": "ናይ ስማርትፎን ስክሪን ምቕያርን ጽገናን ኣገልግሎት",
					},
					Aliases: map[string][]string{
						"en": {"iphone repair", "screen replacement service", "display fix", "broken glass repair"},
						"am": {"የአይፎን ጥገና", "ስክሪን መቀየር", "የተሰበረ መስታወት ጥገና"},
						"or": {"suphaa iskiriinii"},
						"ti": {"ጽገና ስክሪን"},
					},
				},
				// Super Family 55.2: Refrigerator & Major Appliance Maintenance Services
				{
					ID:            "sf-refrigerator-appliance-service",
					Slug:          "refrigerator-ac-appliance-maintenance",
					Code:          "SERV_APPLIANCE_FIX",
					Level:         2,
					CanonicalName: "Refrigerator, AC & Major Appliance Repair Services (ፍሪጅ ጥገና)",
					Definition:    "On-site diagnostic technicians for refrigeration gas charging, thermostat calibration, and motor rewinding.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የፍሪጅ ጋዝ ሙሌት እና የቤት እቃዎች ጥገና አገልግሎት",
						"or": "Tajaajila Suphaa Firiijii fi Meeshaa Manaa",
						"ti": "ናይ ፍሪጅ ጋዝ ምምላእን ጽገናን ኣገልግሎት",
					},
				},
			},
		},

		// 56. Root: Professional Automotive Repair & Mechanical Services
		{
			ID:            "root-automotive-repair-services",
			Slug:          "automotive-repair-mechanical-services",
			Code:          "SERV_AUTO_REPAIR",
			Level:         1,
			CanonicalName: "Automotive Repair, Mechanics & Bodywork Services",
			Definition:    "Professional garage mechanic services, engine overhaul, computerized OBD2 diagnostics, transmission repair, and accident panel beating/spray painting.",
			InclusionRules: []string{"Engine overhauling labor", "Computerized automotive diagnostics", "Panel beating and spray oven painting", "Suspension repair labor"},
			ExclusionRules: []string{"Physical spare parts without service (see Vehicle Parts)", "Complete automobiles (see Passenger Vehicles)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Car",
			Translations: map[string]string{
				"am": "የመኪና ጋራዥ ጥገና እና የሜካኒክ አገልግሎቶች",
				"or": "Tajaajila Suphaa Konkolaataa fi Gaaraajii",
				"ti": "ናይ መኪና ጋራጅ ጽገናን ሜካኒክን ኣገልግሎት",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-garage-engine-mechanics",
					Slug:          "garage-engine-mechanical-overhaul",
					Code:          "AUTO_MECHANICS",
					Level:         2,
					CanonicalName: "Automotive Garage Mechanics & Engine Overhauls (ጋራዥ)",
					Definition:    "Complete engine rebuilding, cylinder head skimming, timing belt replacement, and computerized vehicle diagnostics.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የመኪና ሞተር እና የሜካኒክ ጥገና (ጋራዥ)",
						"or": "Suphaa Mootora Konkolaataa (Gaaraajii)",
						"ti": "ናይ መኪና ሞተር ጽገና (ጋራጅ)",
					},
				},
			},
		},

		// 57. Root: Freight Logistics, Cargo Moving & Transportation Services
		{
			ID:            "root-freight-logistics-services",
			Slug:          "freight-logistics-cargo-moving-services",
			Code:          "SERV_LOGISTICS",
			Level:         1,
			CanonicalName: "Freight Logistics, Cargo Moving & Transportation Services",
			Definition:    "Cross-country dry cargo trucking services, residential furniture house-moving, refrigerated food transport, and customs freight forwarding.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Truck",
			Translations: map[string]string{
				"am": "የጭነት ትራንስፖርት፣ የቤት እቃ ማጓጓዝ እና ሎጂስቲክስ",
				"or": "Tajaajila Geejjiba Meeshaa fi Mobaayilaayizeeshinii",
				"ti": "ናይ ጽዕነት መጓዓዝያ፣ ናይ ገዛ ኣቑሑት ምግዓዝን ሎጂስቲክስን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-house-moving-relocation",
					Slug:          "residential-house-moving-relocation",
					Code:          "SERV_HOUSE_MOVING",
					Level:         2,
					CanonicalName: "Residential House-Moving & Office Relocation Services",
					Definition:    "Professional loading, padded packing, transport, and unloading for household furniture and corporate offices.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የቤት እና የቢሮ እቃ ማጓጓዝ አገልግሎት",
						"or": "Tajaajila Meeshaa Manaa fi Biiroo Geessuu",
						"ti": "ናይ ገዛን ቢሮን ኣቑሑት ምግዓዝ ኣገልግሎት",
					},
				},
			},
		},

		// 58. Root: Travel, Tourism & Hospitality Services
		{
			ID:            "root-tourism-travel-services",
			Slug:          "travel-tourism-hospitality-services",
			Code:          "SERV_TOURISM",
			Level:         1,
			CanonicalName: "Travel, Tourism & Tour Operating Services",
			Definition:    "Certified tourist itineraries (Lalibela, Simien Mountains, Danakil Depression, Omo Valley), tour guide services, and historic route travel packages.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Compass",
			Translations: map[string]string{
				"am": "ጉዞ፣ ቱሪዝም እና የሀገር ውስጥ ጉብኝት አገልግሎቶች",
				"or": "Tajaajila Imala fi Tuurizimii",
				"ti": "መገሻ፣ ቱሪዝምን ናይ ውሽጢ ዓዲ ምብጻሕን ኣገልግሎት",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-historic-route-tours",
					Slug:          "historic-route-nature-tour-packages",
					Code:          "TOUR_PACKAGES",
					Level:         2,
					CanonicalName: "Historic Route & National Park Guided Tour Packages",
					Definition:    "Organized all-inclusive trekking, historical exploration, and cultural photography expeditions.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የታሪካዊ ስፍራዎች እና የብሔራዊ ፓርክ የጉዞ ጥቅሎች",
						"or": "Qophii Imala Bakkeewwan Seenaa fi Paarkii",
						"ti": "ናይ ታሪኻዊ ቦታታትን ፓርክታትን ናይ መገሻ ፓኬጃት",
					},
				},
			},
		},
	}
}
