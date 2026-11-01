package main

func getIndustrialDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 29. Root: Industrial Manufacturing Machinery & Production Lines
		{
			ID:            "root-industrial-machinery",
			Slug:          "industrial-manufacturing-machinery-production",
			Code:          "IND_MACHINERY",
			Level:         1,
			CanonicalName: "Industrial Manufacturing Machinery & Production Lines",
			Definition:    "Factory production machinery, plastic bottle blow molding systems, commercial grain flour mills, packaging equipment, and industrial boilers.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Cpu",
			Translations: map[string]string{
				"am": "የኢንዱስትሪ ማሽነሪዎች እና የማምረቻ መስመሮች",
				"or": "Maashinoota Warshaa fi Oomishaa",
				"ti": "ናይ ኢንዱስትሪ ማሽነሪታትን መስመር ምፍራይን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-grain-flour-mills",
					Slug:          "commercial-grain-flour-mills",
					Code:          "MACH_MILLS",
					Level:         2,
					CanonicalName: "Commercial Grain Flour Mills & Teff Cleaning Plants (የወፍጮ ማሽኖች)",
					Definition:    "High-capacity electric and diesel grain milling machines, teff destoners, seed cleaners, and hammer mills.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የእህል መፍጫ ወፍጮዎች እና የጤፍ ማበጠሪያዎች",
						"or": "Wafchoo Midhaanii fi Meeshaa Qulqulleessituu",
						"ti": "ናይ እኽሊ መጥሓኒታትን መጽረይትን",
					},
					Aliases: map[string][]string{
						"en": {"grain mill", "flour mill", "teff cleaner", "destoner", "hammer mill"},
						"am": {"ወፍጮ", "የእህል ወፍጮ", "የጤፍ ወፍጮ", "ማበጠሪያ"},
						"or": {"wafchoo", "maashina midhaanii"},
						"ti": {"መጥሓኒ", "ናይ እኽሊ ወፍጮ"},
					},
				},
			},
		},

		// 30. Root: Tools, Workshop & Mechanics Equipment
		{
			ID:            "root-tools-workshop",
			Slug:          "tools-workshop-mechanics-equipment",
			Code:          "TOOLS_WORKSHOP",
			Level:         1,
			CanonicalName: "Tools, Workshop & Mechanics Equipment",
			Definition:    "Professional hand toolsets, combination spanners, ratchets, socket sets, hydraulic floor jacks, and workbench vices.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Wrench",
			Translations: map[string]string{
				"am": "የእጅ መሳሪያዎች፣ የሜካኒክ እና የዎርክሾፕ እቃዎች",
				"or": "Meeshaalee Harka fi Suphaa",
				"ti": "መሳርሒታት ኢድን ናይ ዎርክሾፕ ንዋያትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-mechanic-hand-tools",
					Slug:          "mechanic-hand-toolsets-spanners",
					Code:          "TOOL_HANDSETS",
					Level:         2,
					CanonicalName: "Mechanics Toolsets, Spanners & Socket Kits",
					Definition:    "Chrome-vanadium combination wrenches, 1/2-inch drive metric socket sets, and screwdrivers.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የመፍቻ ቁልፎች፣ ካቻቪቴዎች እና የሜካኒክ እቃዎች",
						"or": "Qulfoo fi Meeshaalee Suphaa Konkolaataa",
						"ti": "መፍትሕታት፣ ካቻቪተታትን ናይ መካኒክ ኣቑሑትን",
					},
				},
			},
		},

		// 31. Root: Power Tools, Pneumatic Tools & Cutting Equipment
		{
			ID:            "root-power-tools",
			Slug:          "power-tools-pneumatic-cutting-equipment",
			Code:          "TOOLS_POWER",
			Level:         1,
			CanonicalName: "Power Tools, Pneumatic Tools & Cutting Equipment",
			Definition:    "Cordless brushless drills, rotary hammer drills, angle grinders, circular saws, and pneumatic impact wrenches (Bosch, DeWalt, Makita, Ingco).",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Zap",
			Translations: map[string]string{
				"am": "የኤሌክትሪክ እና የባትሪ እጅ ማሽኖች (ድሪል፣ ግራይንደር)",
				"or": "Maashinoota Harka Elektirikii (Diriilii fi Giraayindara)",
				"ti": "ናይ ኤሌክትሪክ መሳርሒታት (ድሪልን ግራይንደርን)",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-angle-grinders-drills",
					Slug:          "rotary-hammer-drills-angle-grinders",
					Code:          "POWER_DRILLS",
					Level:         2,
					CanonicalName: "Rotary Hammer Drills, Impact Drivers & Angle Grinders",
					Definition:    "SDS-plus rotary hammer drills, 4.5-inch and 9-inch metal cutting angle grinders, and cordless drivers.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ድሪሎች፣ ግራይንደሮች እና የኮንክሪት መበሻዎች",
						"or": "Diriilii fi Giraayindara",
						"ti": "ድሪላትን ግራይንደራትን",
					},
					Aliases: map[string][]string{
						"en": {"power tool", "drill", "grinder", "angle grinder", "rotary hammer", "bosch", "ingco"},
						"am": {"ድሪል", "ግራይንደር", "የብረት መቁረጫ", "የእንጨት መቁረጫ"},
						"or": {"diriilii", "giraayindara"},
						"ti": {"ድሪል", "ግራይንደር"},
					},
				},
			},
		},

		// 32. Root: Welding, Soldering & Thermal Fabrication
		{
			ID:            "root-welding-fabrication",
			Slug:          "welding-soldering-thermal-fabrication",
			Code:          "WELD_EQUIPMENT",
			Level:         1,
			CanonicalName: "Welding, Soldering & Thermal Fabrication Equipment",
			Definition:    "IGBT inverter arc welding machines, MIG/MAG gas welders, plasma cutters, welding electrodes, and autodarkening helmets.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Flame",
			Translations: map[string]string{
				"am": "የብየዳ ማሽኖች እና የብረት ማያያዣ እቃዎች (ሳልዳቱራ)",
				"or": "Maashina fi Meeshaalee Wal-qabsiisa Sibiilaa (Saldaaturaa)",
				"ti": "ናይ ሳልዳቱራ ማሽነሪታትን ብየዳን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-arc-inverter-welders",
					Slug:          "inverter-arc-welding-machines",
					Code:          "WELD_INVERTERS",
					Level:         2,
					CanonicalName: "IGBT Inverter Arc Welders & Electrodes (ሳልዳቱራ)",
					Definition:    "200A to 400A digital inverter MMA welding machines and mild steel welding rods (E6013).",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ኢንቨርተር የሳልዳቱራ ማሽኖች እና የኤሌክትሮድ ሽቦዎች",
						"or": "Maashina Saldaaturaa fi Eelektiroodii",
						"ti": "ናይ ሳልዳቱራ ማሽንን ኤሌክትሮድን",
					},
					Aliases: map[string][]string{
						"en": {"welding machine", "inverter welder", "arc welder", "saldatura", "welding rod"},
						"am": {"ሳልዳቱራ", "የሳልዳቱራ ማሽን", "ብየዳ", "የብየዳ ሽቦ"},
						"or": {"saldaaturaa", "maashina saldaaturaa"},
						"ti": {"ሳልዳቱራ", "ናይ ሳልዳቱራ ማሽን"},
					},
				},
			},
		},

		// 33. Root: Solar Energy Systems, Photovoltaics & Clean Tech
		{
			ID:            "root-solar-energy",
			Slug:          "solar-energy-photovoltaics-clean-tech",
			Code:          "ENERGY_SOLAR",
			Level:         1,
			CanonicalName: "Solar Energy Systems, Photovoltaics & Clean Tech",
			Definition:    "Monocrystalline tier-1 solar panels, pure sine wave hybrid inverters, MPPT charge controllers, and deep-cycle solar batteries.",
			InclusionRules: []string{"Monocrystalline solar panels", "Hybrid inverters (Growatt, Must, Felicity)", "Lithium LiFePO4 batteries", "Solar water heaters", "Solar submersible water pumps"},
			ExclusionRules: []string{"Diesel generators (see Power Generation)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Sun",
			Translations: map[string]string{
				"am": "የፀሐይ ብርሃን ኃይል (ሶላር) እና ንፁህ ኢነርጂ",
				"or": "Annisaa fi Soolaarii Aduu",
				"ti": "ናይ ጸሓይ ጸዓት (ሶላር) መሳርሒታትን",
			},
			Aliases: map[string][]string{
				"en": {"solar", "solar panel", "solar energy", "inverter", "hybrid inverter", "solar battery"},
				"am": {"ሶላር", "የፀሐይ ብርሃን", "ኢንቨርተር", "የሶላር ባትሪ", "ፓነል"},
				"or": {"soolaarii", "paaneelii aduu", "inveertara"},
				"ti": {"ሶላር", "ናይ ጸሓይ ፓነል", "ኢንቨርተር"},
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-photovoltaic-solar-panels",
					Slug:          "photovoltaic-solar-panels-modules",
					Code:          "SOLAR_PANELS",
					Level:         2,
					CanonicalName: "Monocrystalline & Bifacial Photovoltaic Solar Panels",
					Definition:    "High-efficiency solar modules (150W, 300W, 450W, 550W, 650W) for off-grid rural and urban backup systems.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "26111607", Name: "Solar panels"},
					},
					Translations: map[string]string{
						"am": "ሞኖክሪስታላይን የሶላር ፓነሎች",
						"or": "Paaneelii Soolaarii",
						"ti": "ናይ ሶላር ፓነላት",
					},
				},
				{
					ID:            "sf-solar-hybrid-inverters",
					Slug:          "hybrid-solar-inverters-controllers",
					Code:          "SOLAR_INVERTERS",
					Level:         2,
					CanonicalName: "Pure Sine Wave Hybrid Solar Inverters & MPPTs",
					Definition:    "Integrated solar inverter-chargers with grid-tie, battery-less operation, and automatic generator start capabilities.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ሀይብሪድ የሶላር ኢንቨርተሮች",
						"or": "Inveertara Soolaarii",
						"ti": "ናይ ሶላር ኢንቨርተራት",
					},
				},
			},
		},

		// 34. Root: Power Generation, Diesel Generators & Energy Storage
		{
			ID:            "root-power-generation",
			Slug:          "power-generation-diesel-generators-storage",
			Code:          "ENERGY_POWER_GEN",
			Level:         1,
			CanonicalName: "Power Generation, Diesel Generators & Energy Storage",
			Definition:    "Soundproof silent diesel generators (Perkins, Cummins, Baudouin, Ricardo engines from 10kVA to 1500kVA), automatic transfer switches (ATS), and industrial batteries.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Activity",
			Translations: map[string]string{
				"am": "የኤሌክትሪክ ጀነሬተሮች እና የኢነርጂ ማከማቻ",
				"or": "Jenereetara fi Meeshaalee Elektirikii",
				"ti": "ናይ ኤሌክትሪክ ጀነሬተራትን ጸዓትን",
			},
			Aliases: map[string][]string{
				"en": {"generator", "diesel generator", "perkins generator", "genset", "silent generator", "ats"},
				"am": {"ጀነሬተር", "የናፍጣ ጀነሬተር", "ፐርኪንስ ጀነሬተር", "ድምፅ አልባ ጀነሬተር"},
				"or": {"jenereetara", "jenereetara diizilii"},
				"ti": {"ጀነሬተር", "ናይ ዲዘል ጀነሬተር"},
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-silent-diesel-generators",
					Slug:          "soundproof-silent-diesel-generators",
					Code:          "GEN_DIESEL",
					Level:         2,
					CanonicalName: "Soundproof Silent Commercial Diesel Generators (ጀነሬተሮች)",
					Definition:    "Heavy-duty water-cooled diesel gensets equipped with sound-attenuated acoustic canopies and deep-sea digital controllers.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "26111601", Name: "Diesel generators"},
					},
					Translations: map[string]string{
						"am": "ድምጽ አልባ የናፍጣ ጀነሬተሮች (ፐርኪንስ / ካሚንስ)",
						"or": "Jenereetara Sagalee Hin Qabne (Diizilii)",
						"ti": "ድምጺ ዘይብሎም ናይ ዲዘል ጀነሬተራት",
					},
				},
			},
		},

		// 35. Root: Industrial Chemicals, Raw Substances & Laboratory Reagents
		{
			ID:            "root-industrial-chemicals",
			Slug:          "industrial-chemicals-raw-substances-reagents",
			Code:          "IND_CHEMICALS",
			Level:         1,
			CanonicalName: "Industrial Chemicals, Raw Substances & Reagents",
			Definition:    "Textile dyes, leather tanning salts, industrial surfactants, caustic soda, hydrochloric acid, laboratory reagents, and raw plastics.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "FlaskConical",
			Translations: map[string]string{
				"am": "የኢንዱስትሪ ኬሚካሎች፣ ጥሬ እቃዎች እና ሬጀንቶች",
				"or": "Keemikaala Warshaa fi Meeshaalee Dhedhii",
				"ti": "ናይ ኢንዱስትሪ ኬሚካላትን ጥረ ነገራትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-industrial-processing-chemicals",
					Slug:          "industrial-processing-chemicals-salts",
					Code:          "CHEM_PROCESSING",
					Level:         2,
					CanonicalName: "Industrial Processing Chemicals & Salts",
					Definition:    "Bulk caustic soda (sodium hydroxide, CAS 1310-73-2), sulfuric acid, sodium hypochlorite, and industrial detergents.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "CAS", Code: "1310-73-2", Name: "Sodium Hydroxide"},
						{Type: "UNSPSC", Code: "12352100", Name: "Industrial acids and bases"},
					},
					Translations: map[string]string{
						"am": "የማምረቻ ኬሚካሎች (ኮስቲክ ሶዳ፣ አሲድ)",
						"or": "Keemikaala Oomishaa fi Qoricha",
						"ti": "ናይ ምፍራይ ኬሚካላት",
					},
				},
			},
		},
	}
}
