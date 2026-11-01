package main

func getConstructionDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 25. Root: Construction Materials, Masonry & Structural Aggregates
		{
			ID:            "root-construction-materials",
			Slug:          "construction-materials-masonry-aggregates",
			Code:          "CONST_MATERIALS",
			Level:         1,
			CanonicalName: "Construction Materials, Masonry & Aggregates",
			Definition:    "Primary building construction materials including Portland cement, red clay bricks, hollow concrete blocks (HCB), aggregate stone, and washed river sand.",
			InclusionRules: []string{"Portland cement (Dangote, Mugher, Derba, National)", "Hollow concrete blocks (HCB)", "River sand and crushed basalt aggregate", "Hydrated lime and gypsum plaster"},
			ExclusionRules: []string{"Structural steel rebar (see Structural Steel)", "Earthmoving excavators (see Heavy Construction Machinery)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Layers",
			Translations: map[string]string{
				"am": "የግንባታ እቃዎች፣ ሲሚንቶ፣ አሸዋ እና ብሎኬት",
				"or": "Meeshaalee Ijaarsaa, Simentoo fi Cirracha",
				"ti": "ናይ ህንጻ መሳርሒታት፣ ሲሚንቶ፣ ሑጻን ብሎኬትን",
			},
			Aliases: map[string][]string{
				"en": {"construction materials", "building materials", "cement", "sand", "hcb", "concrete blocks"},
				"am": {"የግንባታ እቃ", "ሲሚንቶ", "ብሎኬት", "አሸዋ", "ጠጠር", "ዳንጎቴ", "ሙገር"},
				"or": {"meeshaa ijaarsaa", "simentoo", "cirracha", "bilookeettii"},
				"ti": {"ናይ ህንጻ ኣቑሑት", "ሲሚንቶ", "ብሎኬት", "ሑጻ"},
			},
			Children: []TaxonomyNode{
				// Super Family 25.1: Cement & Binding Agents
				{
					ID:            "sf-cement-binders",
					Slug:          "portland-cement-binding-agents",
					Code:          "MAT_CEMENT",
					Level:         2,
					CanonicalName: "Portland Cement & Hydraulic Binding Agents (ሲሚንቶ)",
					Definition:    "OPC (Ordinary Portland Cement 42.5R) and PPC (Pozzolana Portland Cement 32.5N) for structural concrete and plastering.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "30111601", Name: "Portland cement"},
					},
					Translations: map[string]string{
						"am": "ፖርትላንድ ሲሚንቶ (ዳንጎቴ፣ ሙገር፣ ደርባ)",
						"or": "Simentoo Poortilaandii",
						"ti": "ፖርትላንድ ሲሚንቶ",
					},
					Aliases: map[string][]string{
						"en": {"cement", "dangote cement", "mugher cement", "derba cement", "ppc", "opc"},
						"am": {"ሲሚንቶ", "የዳንጎቴ ሲሚንቶ", "የሙገር ሲሚንቶ", "የደርባ ሲሚንቶ"},
						"or": {"simentoo", "simentoo daangotee"},
						"ti": {"ሲሚንቶ", "ናይ ዳንጎተ ሲሚንቶ"},
					},
				},
				// Super Family 25.2: Masonry Blocks & Precast Concrete
				{
					ID:            "sf-masonry-blocks",
					Slug:          "hollow-concrete-blocks-masonry",
					Code:          "MAT_BLOCKS",
					Level:         2,
					CanonicalName: "Hollow Concrete Blocks (HCB) & Masonry Bricks",
					Definition:    "Standard machine-vibrated concrete blocks (10cm, 15cm, 20cm thickness), paving stones, and interlocking curb stones.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የኮንክሪት ብሎኬቶች (HCB) እና ጡቦች",
						"or": "Bilookeettii Konkiraatii fi Xuubii",
						"ti": "ናይ ኮንክሪት ብሎኬታትን ጡብን",
					},
				},
			},
		},

		// 26. Root: Structural Steel, Metal Framing & Rebar
		{
			ID:            "root-structural-steel",
			Slug:          "structural-steel-metal-framing-rebar",
			Code:          "CONST_STEEL",
			Level:         1,
			CanonicalName: "Structural Steel, Metal Framing & Rebar",
			Definition:    "Deformed high-yield steel reinforcement rebar, RHS/SHS hollow tubing, H-beams, angle iron, and corrugated iron sheet framing.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Grid",
			Translations: map[string]string{
				"am": "የብረት ዘንጎች (ቆርቆሮ እና ባስቶሎ)",
				"or": "Sibiila Ijaarsaa fi Qorqorroo",
				"ti": "ናይ ህንጻ ብረታትን ቆርቆሮን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-reinforcement-rebar",
					Slug:          "deformed-reinforcement-steel-rebar",
					Code:          "STEEL_REBAR",
					Level:         2,
					CanonicalName: "Deformed Concrete Reinforcement Steel Rebar (ባስቶሎ)",
					Definition:    "Grade 60 and Grade 75 deformed rebars in diameters: 8mm, 10mm, 12mm, 14mm, 16mm, 20mm, 24mm, and 32mm.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "30103623", Name: "Concrete reinforcement rebar"},
					},
					Translations: map[string]string{
						"am": "የኮንክሪት ብረት ዘንጎች (ባስቶሎ)",
						"or": "Sibiila Ijaarsa Konkiraatii (Baastoloo)",
						"ti": "ናይ ህንጻ ብረት ዘንጊታት (ባስቶሎ)",
					},
					Aliases: map[string][]string{
						"en": {"rebar", "reinforcement steel", "steel bar", "bastolo", "deformed bar"},
						"am": {"ባስቶሎ", "የግንባታ ብረት", "የኮንክሪት ብረት", "ባለ 12 ብረት", "ባለ 16 ብረት"},
						"or": {"baastoloo", "sibiila ijaarsaa"},
						"ti": {"ባስቶሎ", "ናይ ህንጻ ብረት"},
					},
				},
			},
		},

		// 27. Root: Heavy Construction Machinery & Earthmoving Equipment
		{
			ID:            "root-construction-machinery",
			Slug:          "heavy-construction-machinery-earthmoving",
			Code:          "CONST_MACHINERY",
			Level:         1,
			CanonicalName: "Heavy Construction Machinery & Earthmoving Equipment",
			Definition:    "Hydraulic crawler excavators, track-type bulldozers, motor graders, wheel loaders, mobile cranes, and transit concrete mixers.",
			InclusionRules: []string{"Hydraulic excavators (CAT, Komatsu, Sany)", "Track bulldozers (D6, D8)", "Motor graders", "Transit concrete mixers", "Tower cranes and mobile truck cranes"},
			ExclusionRules: []string{"Hand drills and personal power tools (see Power Tools)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Hammer",
			Translations: map[string]string{
				"am": "የኮንስትራክሽን እና ከባድ ማሽነሪዎች",
				"or": "Maashinoota Ijaarsaa Gurguddoo",
				"ti": "ናይ ህንጻን ከበድቲ ማሽነሪታትን",
			},
			Aliases: map[string][]string{
				"en": {"construction machinery", "heavy equipment", "excavator", "bulldozer", "loader", "crane", "grader"},
				"am": {"ከባድ ማሽነሪ", "ኤክስካቫተር", "ቡልዶዘር", "ሎደር", "ግሬደር", "ክሬን"},
				"or": {"maashina gurguddaa", "ekiskaaveetara", "buldoozara"},
				"ti": {"ከበድቲ ማሽነሪታት", "ኤክስካቫተር", "ቡልዶዘር"},
			},
			Children: []TaxonomyNode{
				// Super Family 27.1: Earthmoving Excavators & Bulldozers
				{
					ID:            "sf-earthmoving-excavators",
					Slug:          "hydraulic-excavators-bulldozers",
					Code:          "MACH_EARTHMOVING",
					Level:         2,
					CanonicalName: "Hydraulic Excavators, Bulldozers & Wheel Loaders",
					Definition:    "Heavy earthmoving machines (CAT 320, Komatsu PC200, SANY SY215, Liugong loaders).",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "22101505", Name: "Hydraulic excavators"},
					},
					Translations: map[string]string{
						"am": "ኤክስካቫተሮች፣ ቡልዶዘሮች እና ሎደሮች",
						"or": "Ekiskaaveetara fi Buldoozara",
						"ti": "ኤክስካቫተራትን ቡልዶዘራትን",
					},
				},
				// Super Family 27.2: Concrete Transit Mixers & Batching Plants
				{
					ID:            "sf-concrete-mixers-pumps",
					Slug:          "concrete-transit-mixers-batching",
					Code:          "MACH_CONCRETE",
					Level:         2,
					CanonicalName: "Concrete Transit Mixers, Pumps & Batching Plants",
					Definition:    "Truck-mounted rotating drum transit concrete mixers (6m³ to 12m³), boom concrete pumps, and stationary pan mixers.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የኮንክሪት ማደባለቂያዎች እና ፓምፖች",
						"or": "Maashina Konkiraatii fi Paampii",
						"ti": "ናይ ኮንክሪት መሐወሲታትን ፓምፓትን",
					},
				},
			},
		},

		// 28. Root: Plumbing, Piping & Sanitary Ware
		{
			ID:            "root-plumbing-sanitary",
			Slug:          "plumbing-piping-sanitary-ware",
			Code:          "CONST_PLUMBING",
			Level:         1,
			CanonicalName: "Plumbing, Piping, Valves & Sanitary Ware",
			Definition:    "PPR hot/cold water pipes, PVC drainage tubing, submersible water well pumps, brass gate valves, and ceramic toilets/sinks.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Pipette",
			Translations: map[string]string{
				"am": "የውሃ ቧንቧ እቃዎች፣ ቫልቮች እና የንፅህና መጠበቂያዎች",
				"or": "Ujummoo Bishaanii fi Meeshaalee Qulqullinaa",
				"ti": "ናይ ማይ ሻምብቆታትን ጽሬት መሳርሒታትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-water-pipes-fittings",
					Slug:          "ppr-pvc-water-pipes-fittings",
					Code:          "PLUMB_PIPES",
					Level:         2,
					CanonicalName: "PPR, HDPE & PVC Water Pipes & Compression Fittings",
					Definition:    "High-density polyethylene irrigation pipes, green PPR potable water pipes, and PVC sewer drainage tubes.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የውሃ ቧንቧዎች እና ማያያዣ እቃዎች (PPR / PVC)",
						"or": "Ujummoolee Bishaanii fi Wal-qabsiistota",
						"ti": "ናይ ማይ ሻምብቆታትን መላገቢታትን",
					},
				},
			},
		},
	}
}
