package main

func getHealthDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 44. Root: Healthcare, Medical Devices & Diagnostic Equipment
		{
			ID:            "root-medical-devices",
			Slug:          "healthcare-medical-devices-diagnostics",
			Code:          "MED_DEVICES",
			Level:         1,
			CanonicalName: "Healthcare, Medical Devices & Diagnostic Equipment",
			Definition:    "Hospital clinical equipment, diagnostic ultrasound scanners, digital blood pressure monitors, pulse oximeters, and surgical instruments.",
			InclusionRules: []string{"Diagnostic ultrasound scanners", "Digital BP monitors", "Blood glucose meters", "Patient monitors and ECG", "Hospital beds and IV stands"},
			ExclusionRules: []string{"Wheelchairs and mobility canes (see Mobility Aids)", "Pharmaceutical drugs (see Pharmaceuticals)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Activity",
			Translations: map[string]string{
				"am": "የህክምና መሳሪያዎች እና የጤና ምርመራ እቃዎች",
				"or": "Meeshaalee Yaalaa fi Qorannoo Fayyaa",
				"ti": "ናይ ሕክምና መሳርሒታትን መርመራ ጥዕናን",
			},
			Aliases: map[string][]string{
				"en": {"medical equipment", "healthcare", "bp monitor", "ultrasound", "hospital equipment", "stethoscope"},
				"am": {"የህክምና እቃ", "የህክምና መሳሪያ", "የግፊት መለኪያ", "የስኳር መለኪያ", "አልትራሳውንድ"},
				"or": {"meeshaa yaalaa", "meeshaa dhiibbaa dhiigaa"},
				"ti": {"ናይ ሕክምና ኣቑሑት", "ናይ ደም ጸቕጢ መዐቀኒ"},
			},
			Children: []TaxonomyNode{
				// Super Family 44.1: Home & Clinic Diagnostic Monitors
				{
					ID:            "sf-clinical-monitors-diagnostic",
					Slug:          "home-clinic-diagnostic-monitors",
					Code:          "DIAG_MONITORS",
					Level:         2,
					CanonicalName: "Diagnostic Monitors & Vitals Checkers (BP, Glucose, Oximeter)",
					Definition:    "Automated upper-arm digital blood pressure monitors (Omron, Beurer), fingertip pulse oximeters, and blood glucometers.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "42181601", Name: "Blood pressure monitors"},
					},
					Translations: map[string]string{
						"am": "የደም ግፊት እና የስኳር መለኪያ መሳሪያዎች",
						"or": "Meeshaalee Dhiibbaa Dhiigaa fi Sukkaaraa",
						"ti": "ናይ ደም ጸቕጥን ሽኮርን መዐቀኒታት",
					},
					Aliases: map[string][]string{
						"en": {"blood pressure monitor", "bp machine", "glucometer", "pulse oximeter"},
						"am": {"የግፊት መለኪያ", "የስኳር መለኪያ", "የልብ ምት መለኪያ"},
						"or": {"safartuu dhiibbaa dhiigaa"},
						"ti": {"ናይ ጸቕጢ መዐቀኒ"},
					},
				},
				// Super Family 44.2: Ultrasound & Imaging Equipment
				{
					ID:            "sf-ultrasound-imaging",
					Slug:          "ultrasound-scanners-imaging-systems",
					Code:          "MED_IMAGING",
					Level:         2,
					CanonicalName: "Diagnostic Ultrasound Scanners & Radiology Imaging",
					Definition:    "Portable and trolley-mounted color Doppler ultrasound diagnostic scanners for obstetrics, cardiology, and abdominal imaging.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የአልትራሳውንድ ማሽኖች እና የራዲዮሎጂ እቃዎች",
						"or": "Maashina Altraasawundii fi Suuraa Yaalaa",
						"ti": "ናይ ኣልትራሳውንድ ማሽናትን ራዲዮሎጂን",
					},
				},
			},
		},

		// 45. Root: Mobility Aids, Accessibility & Orthopedic Equipment
		{
			ID:            "root-mobility-aids",
			Slug:          "mobility-aids-accessibility-orthopedics",
			Code:          "MED_MOBILITY",
			Level:         1,
			CanonicalName: "Mobility Aids, Accessibility & Orthopedic Equipment",
			Definition:    "Manual folding wheelchairs, motorized power wheelchairs, aluminum underarm crutches, walker rollators, and orthopedic braces.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Accessibility",
			Translations: map[string]string{
				"am": "የመንቀሳቀሻ ድጋፍ ሰጪ እቃዎች (ዊልቸር፣ ክራንች)",
				"or": "Meeshaalee Deeggarsa Sochii (Wiilcheerii)",
				"ti": "ናይ ምንቅስቓስ ደገፍቲ ኣቑሑት (ዊልቸር፣ ክራንች)",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-wheelchairs-walkers",
					Slug:          "manual-motorized-wheelchairs",
					Code:          "AID_WHEELCHAIRS",
					Level:         2,
					CanonicalName: "Manual & Electric Wheelchairs & Mobility Walkers (ዊልቸር)",
					Definition:    "Chrome steel folding wheelchairs, commode wheelchairs, and electric powerchairs.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "42192210", Name: "Wheelchairs"},
					},
					Translations: map[string]string{
						"am": "የእጅ እና የሞተር ዊልቸሮች",
						"or": "Wiilcheerii Harka fi Elektirikii",
						"ti": "ናይ ኢድን ሞተርን ዊልቸራት",
					},
					Aliases: map[string][]string{
						"en": {"wheelchair", "electric wheelchair", "crutches", "walker"},
						"am": {"ዊልቸር", "ክራንች", "የአካል ጉዳተኞች ወንበር"},
						"or": {"wiilcheerii"},
						"ti": {"ዊልቸር", "ክራንች"},
					},
				},
			},
		},
	}
}
