package main

func getTransportDomains() []TaxonomyNode {
	return []TaxonomyNode{
		// 14. Root: Passenger Vehicles & Automobiles
		{
			ID:            "root-passenger-vehicles",
			Slug:          "passenger-vehicles-automobiles",
			Code:          "TRANS_PASSENGER",
			Level:         1,
			CanonicalName: "Passenger Vehicles & Automobiles",
			Definition:    "Complete registered passenger cars, sedans, hatchbacks, compact crossovers, SUVs, and luxury 4WD automobiles.",
			InclusionRules: []string{"Sedans", "Hatchbacks", "SUVs / 4x4", "Station wagons", "Electric passenger cars"},
			ExclusionRules: []string{"Heavy commercial trucks (see Commercial Trucks)", "Vehicle replacement spare parts (see Vehicle Parts)", "Automobile rental services (see Transportation Services)"},
			Examples:       []string{"Toyota Vitz", "Toyota Corolla", "Toyota Land Cruiser", "Hyundai Tucson", "Suzuki Dzire", "BYD Atto 3 EV"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Car",
			Translations: map[string]string{
				"am": "ተሽከርካሪዎች እና የመንገደኛ መኪኖች",
				"or": "Konkolaattota fi Awutoomoobilii",
				"ti": "ተሽከርከርትን ናይ ህዝቢ መካይንን",
			},
			Aliases: map[string][]string{
				"en": {"car", "cars", "automobile", "passenger car", "sedan", "suv", "vehicle"},
				"am": {"መኪና", "ተሽከርካሪ", "አውቶሞቢል", "ቪትዝ", "ኮሮላ", "ላንድክሩዘር"},
				"or": {"konkolaataa", "awutoomoobilii", "makiinaa"},
				"ti": {"መኪና", "ተሽከርካሪ", "ኣውቶሞቢል"},
			},
			Children: []TaxonomyNode{
				// Super Family 14.1: Sedans & Saloon Cars
				{
					ID:            "sf-sedans-saloons",
					Slug:          "sedans-saloon-automobiles",
					Code:          "AUTO_SEDANS",
					Level:         2,
					CanonicalName: "Sedans & Saloon Automobiles",
					Definition:    "Three-box passenger automobiles with separate compartments for engine, passenger cabin, and boot/trunk.",
					Status:         "active",
					IsLeaf:         true,
					Identifiers: []TaxonomyIdentifier{
						{Type: "UNSPSC", Code: "25101503", Name: "Passenger automobiles"},
					},
					Translations: map[string]string{
						"am": "ሴዳን እና ሳሎን መኪኖች",
						"or": "Konkolaataa Seedaanii",
						"ti": "ሴዳን መካይን",
					},
					Attributes: []TaxonomyAttribute{
						{Slug: "vehicle-year", NameEN: "Manufacture Year", DataType: "number", IsRequired: true, IsFilterable: true},
						{Slug: "vehicle-mileage-km", NameEN: "Mileage (km)", DataType: "number", Unit: "km", IsRequired: true, IsFilterable: true},
						{Slug: "transmission-type", NameEN: "Transmission", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Automatic", "Manual", "CVT"}}},
						{Slug: "fuel-type", NameEN: "Fuel / Powertrain", DataType: "select", IsRequired: true, IsFilterable: true, ValidationRules: map[string]interface{}{"options": []string{"Benzine (Petrol)", "Diesel", "Electric (EV)", "Hybrid"}}},
					},
				},
				// Super Family 14.2: Hatchbacks & Compact Cars
				{
					ID:            "sf-hatchbacks-compacts",
					Slug:          "hatchbacks-compact-cars",
					Code:          "AUTO_HATCHBACKS",
					Level:         2,
					CanonicalName: "Hatchbacks & Compact City Cars (ቪትዝ እና ትናንሽ መኪኖች)",
					Definition:    "Two-box city runarounds (Toyota Vitz, Yaris, Suzuki Swift, Hyundai Grand i10) popular for fuel economy.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ሀችባክ እና የታመቁ ትናንሽ መኪኖች (ቪትዝ)",
						"or": "Konkolaattota Qancoo fi Haachbaakii",
						"ti": "ሃችባክን ናይ ከተማ መካይንን",
					},
					Aliases: map[string][]string{
						"en": {"hatchback", "compact car", "toyota vitz", "suzuki swift", "small car"},
						"am": {"ቪትዝ", "ያሪስ", "ሀችባክ", "ትንሽ መኪና"},
						"or": {"viitsii", "konkolaataa xiqqoo"},
						"ti": {"ቪትዝ", "ንእሽቶ መኪና"},
					},
				},
				// Super Family 14.3: SUVs & 4x4 Off-Road Vehicles
				{
					ID:            "sf-suvs-offroad-4x4",
					Slug:          "suvs-4x4-offroad-vehicles",
					Code:          "AUTO_SUVS",
					Level:         2,
					CanonicalName: "SUVs, Crossovers & 4x4 Off-Road Vehicles",
					Definition:    "High-clearance all-wheel drive and four-wheel drive wagons (Land Cruiser, RAV4, Prado, Fortuner).",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ኤስዩቪ፣ ፎርዊል ድራይቭ እና 4x4 መኪኖች",
						"or": "Konkolaataa Foor-wiilii fi Ees-Yuu-Vii",
						"ti": "ኤስዩቪን 4x4 መካይንን",
					},
				},
			},
		},

		// 15. Root: Commercial Trucks, Buses & Heavy Transport
		{
			ID:            "root-commercial-transport",
			Slug:          "commercial-trucks-buses-heavy-transport",
			Code:          "TRANS_COMMERCIAL",
			Level:         1,
			CanonicalName: "Commercial Trucks, Buses & Heavy Transport",
			Definition:    "Medium/heavy cargo trucks, Sino dump trucks, Isuzu FSR/NPR freight haulers, minibuses, and passenger coaches.",
			InclusionRules: []string{"Isuzu Forward / NPR / FSR", "Sino Truck dumpers", "Toyota HiAce minibuses", "Coaster midibuses", "Long-haul prime movers"},
			ExclusionRules: []string{"Earthmoving excavators without road certification (see Heavy Construction Machinery)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Truck",
			Translations: map[string]string{
				"am": "የጭነት፣ የአውቶቡስ እና የህዝብ ትራንስፖርት ተሽከርካሪዎች",
				"or": "Konkolaataa Fe'iinsaa fi Baasii Gurguddoo",
				"ti": "ናይ ጽዕነትን ኣውቶቡስን ከበድቲ ተሽከርከርቲ",
			},
			Children: []TaxonomyNode{
				// Super Family 15.1: Medium Freight Haulers & Pickups
				{
					ID:            "sf-trucks-pickups-freight",
					Slug:          "light-medium-cargo-trucks-pickups",
					Code:          "TRUCK_MEDIUM",
					Level:         2,
					CanonicalName: "Light & Medium Cargo Trucks & Double-Cab Pickups",
					Definition:    "Isuzu NPR/NQR/FSR box trucks, dropsides, and Toyota Hilux double/single cab work utility vehicles.",
					Status:         "active",
					IsLeaf:         false,
					Translations: map[string]string{
						"am": "ቀላል እና መካከለኛ የጭነት መኪኖች እና ፒክአፕ",
						"or": "Konkolaataa Fe'iinsa Giddu-galeessaa fi Piik-appii",
						"ti": "ቀለልትን ማእከላይን ናይ ጽዕነት መካይንን ፒካፕን",
					},
					Children: []TaxonomyNode{
						// Level 3: Utility Pickups (Hilux, etc.)
						{
							ID:            "cat-utility-pickups",
							Slug:          "double-single-cab-pickups",
							Code:          "PICKUP_TRUCKS",
							Level:         3,
							CanonicalName: "Double & Single Cabin Utility Pickups (ፒክአፕ)",
							Definition:    "4x4 and 4x2 light commercial double-cabin and single-cabin pickup trucks.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "ባለ ሁለት እና ባለ አንድ ጋቢና ፒክአፕ",
								"or": "Konkolaataa Piik-appii Gaabinaa Lamaa fi Tokkoo",
								"ti": "ናይ ጽዕነት ፒካፕ",
							},
							Aliases: map[string][]string{
								"en": {"pickup", "hilux", "double cabin", "single cabin", "4x4 pickup"},
								"am": {"ፒክአፕ", "ሃይሉክስ", "ባለ ሁለት ጋቢና", "የጭነት ፒክአፕ"},
								"or": {"piik-appii", "haayilaksii"},
								"ti": {"ፒካፕ", "ሃይሉክስ"},
							},
						},
						// Level 3: Isuzu Cargo Haulers
						{
							ID:            "cat-isuzu-haulers",
							Slug:          "isuzu-cargo-trucks-npr-fsr",
							Code:          "ISUZU_HAULERS",
							Level:         3,
							CanonicalName: "Isuzu Commercial Cargo Trucks (NPR / FSR / Forward)",
							Definition:    "High-volume commercial dry cargo transport vehicles (FSR, NPR, Forward, Reward).",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "አይሱዙ የጭነት መኪኖች (ኤፍኤስአር / ኤንፒአር)",
								"or": "Konkolaataa Fe'iinsaa Isuzuu (FSR / NPR)",
								"ti": "ኣይሱዙ ናይ ጽዕነት መካይን",
							},
							Aliases: map[string][]string{
								"en": {"isuzu", "isuzu fsr", "isuzu npr", "isuzu forward", "cargo truck"},
								"am": {"አይሱዙ", "ኤፍኤስአር", "ኤንፒአር", "የጭነት አይሱዙ"},
								"or": {"isuzuu", "isuzuu fe'iinsaa"},
								"ti": {"ኣይሱዙ", "ናይ ጽዕነት ኣይሱዙ"},
							},
						},
					},
				},
			},
		},

		// 16. Root: Motorcycles, Scooters & Three-Wheelers (Bajaj)
		{
			ID:            "root-motorcycles-bajaj",
			Slug:          "motorcycles-scooters-three-wheelers-bajaj",
			Code:          "TRANS_TWO_THREE_WHEEL",
			Level:         1,
			CanonicalName: "Motorcycles, Scooters & Three-Wheelers (Bajaj)",
			Definition:    "Three-wheeled passenger auto-rickshaws (Bajaj RE, TVS King, Piaggio), cargo tricycles, and delivery motorbikes.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Bike",
			Translations: map[string]string{
				"am": "ሞተር ሳይክሎች፣ ስኩተሮች እና ባጃጆች",
				"or": "Moophedii fi Baajaajii",
				"ti": "ሞተር ሳይክላትን ባጃጃትን",
			},
			Children: []TaxonomyNode{
				// Super Family 16.1: Three-Wheeled Auto Rickshaws (Bajaj)
				{
					ID:            "sf-bajaj-three-wheelers",
					Slug:          "three-wheeled-auto-rickshaws-bajaj",
					Code:          "AUTO_BAJAJ",
					Level:         2,
					CanonicalName: "Three-Wheeled Auto-Rickshaws & Bajaj",
					Definition:    "Passenger urban transport three-wheelers (Bajaj RE 4S, TVS King Deluxe) and rear cargo tricycles.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "ባጃጆች እና ባለ ሶስት እግር ተሽከርካሪዎች",
						"or": "Baajaajii fi Konkolaataa Miila Sadii",
						"ti": "ባጃጃትን ባለ ሰለስተ እግሪ ተሽከርከርትን",
					},
					Aliases: map[string][]string{
						"en": {"bajaj", "auto rickshaw", "three wheeler", "tvs king", "keke"},
						"am": {"ባጃጅ", "ባለ 3 እግር", "ቲቪኤስ", "የከተማ ባጃጅ"},
						"or": {"baajaajii", "baajaajii miila sadii"},
						"ti": {"ባጃጅ", "ባለ 3 እግሪ"},
					},
				},
			},
		},

		// 17. Root: Bicycles, Micro-Mobility & Non-Motorized Transport
		{
			ID:            "root-bicycles-mobility",
			Slug:          "bicycles-micromobility-transport",
			Code:          "TRANS_BICYCLES",
			Level:         1,
			CanonicalName: "Bicycles, Micro-Mobility & Non-Motorized Transport",
			Definition:    "Mountain bicycles, urban commuter road bikes, children's pedal cycles, electric standup scooters, and cargo pushcarts.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Compass",
			Translations: map[string]string{
				"am": "ብስክሌቶች እና ቀላል ተንቀሳቃሽ ተሽከርካሪዎች",
				"or": "Saayikilii fi Geejjiba Salphaa",
				"ti": "ብሽክለታታትን ቀለልቲ መጓዓዝያን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-pedal-bicycles",
					Slug:          "pedal-bicycles-mountain-bikes",
					Code:          "BIKE_PEDAL",
					Level:         2,
					CanonicalName: "Mountain Bikes, City Road Cycles & Hybrid Bicycles",
					Definition:    "Multi-gear aluminum mountain bicycles, delivery cycles, and commuter road bikes.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የተራራ እና የመንገድ ብስክሌቶች",
						"or": "Saayikilii Gaaraa fi Daandii",
						"ti": "ናይ ጎቦን መንገዲን ብሽክለታታት",
					},
				},
			},
		},

		// 18. Root: Vehicle Parts, Replacement Components & Powertrain
		{
			ID:            "root-vehicle-parts",
			Slug:          "vehicle-parts-replacement-powertrain",
			Code:          "AUTO_PARTS",
			Level:         1,
			CanonicalName: "Vehicle Parts, Replacement Components & Powertrain",
			Definition:    "Physical mechanical vehicle replacement components, complete crate engines, gearboxes, brake pads, suspension shocks, and filters.",
			InclusionRules: []string{"Brake pads, rotors and calipers", "Alternators and starter motors", "Shock absorbers and coil springs", "Air/oil/fuel filters", "Engine pistons, timing belts and gaskets"},
			ExclusionRules: []string{"Mechanic garage repair labor (see Vehicle Repair Services)", "Vehicle tires (see Vehicle Tires)"},
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Wrench",
			Translations: map[string]string{
				"am": "የተሽከርካሪ መለዋወጫዎች እና የሞተር እቃዎች",
				"or": "Qaama Jijjiirraa Konkolaataa fi Mootoraa",
				"ti": "ናይ መካይን መለዋወጢ ኣቑሑትን ሞተርን",
			},
			Aliases: map[string][]string{
				"en": {"spare parts", "car parts", "auto parts", "brake pads", "engine parts"},
				"am": {"መለዋወጫ", "የመኪና መለዋወጫ", "ፍሬን", "ሾክ", "ፊልተር"},
				"or": {"qaama jijjiirraa", "sparii", "firiinii"},
				"ti": {"መለዋወጢ", "ናይ መኪና መለዋወጢ", "ፍሬን"},
			},
			Children: []TaxonomyNode{
				// Super Family 18.1: Brake Systems & Friction Components
				{
					ID:            "sf-vehicle-brake-systems",
					Slug:          "vehicle-brake-systems-friction",
					Code:          "PART_BRAKES",
					Level:         2,
					CanonicalName: "Brake Systems & Friction Components",
					Definition:    "Ceramic/semi-metallic brake pads, vented disc rotors, brake master cylinders, and ABS hydraulic modules.",
					Status:         "active",
					IsLeaf:         false,
					Translations: map[string]string{
						"am": "የፍሬን ሲስተም እና የፍሬን ጫማዎች",
						"or": "Sirna Firiinii Konkolaataa",
						"ti": "ናይ ፍሬን ሲስተምን ኣቑሑትን",
					},
					Children: []TaxonomyNode{
						// Level 3: Brake Pads & Shoes
						{
							ID:            "cat-brake-pads-shoes",
							Slug:          "automotive-brake-pads-shoes",
							Code:          "BRAKE_PADS",
							Level:         3,
							CanonicalName: "Automotive Brake Pads & Shoes (የፍሬን ጫማ / ባሌታ)",
							Definition:    "Front and rear axle friction brake pads and drum brake shoes for passenger and commercial vehicles.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የፊት እና የኋላ የፍሬን ጫማዎች (ባሌታ)",
								"or": "Kophee Firiinii Konkolaataa",
								"ti": "ናይ ፍሬን ባሌታታትን ጫማታትን",
							},
							Aliases: map[string][]string{
								"en": {"brake pad", "brake shoe", "baleta", "front brake pads"},
								"am": {"ባሌታ", "የፍሬን ጫማ", "የፍሬን ባሌታ"},
								"or": {"firiinii", "kophee firiinii"},
								"ti": {"ባሌታ", "ናይ ፍሬን ጫማ"},
							},
						},
						// Level 3: Brake Rotors & Drums
						{
							ID:            "cat-brake-rotors-discs",
							Slug:          "brake-rotors-disc-drums",
							Code:          "BRAKE_ROTORS",
							Level:         3,
							CanonicalName: "Brake Rotors, Discs & Wheel Drums",
							Definition:    "Ventilated, solid, and slotted iron brake discs and cast rear wheel brake drums.",
							Status:         "active",
							IsLeaf:         true,
							Translations: map[string]string{
								"am": "የፍሬን ዲስኮች እና ድራሞች",
								"or": "Diskiilee Firiinii",
								"ti": "ናይ ፍሬን ዲስክታትን ድራምን",
							},
						},
					},
				},
				// Super Family 18.2: Filtration & Engine Lubrication
				{
					ID:            "sf-vehicle-filters-lube",
					Slug:          "vehicle-filters-filtration",
					Code:          "PART_FILTERS",
					Level:         2,
					CanonicalName: "Automotive Filters (Air, Oil, Fuel, Cabin)",
					Definition:    "Spin-on oil filters, pleated engine air intake filters, diesel fuel water-separators, and AC cabin filters.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የዘይት፣ የነዳጅ እና የአየር ፊልተሮች",
						"or": "Filterii Zayitaa, Boba'aa fi Qilleensaa",
						"ti": "ናይ ዘይትን ነዳድን ኣየርን ፊልተራት",
					},
					Aliases: map[string][]string{
						"en": {"filter", "oil filter", "fuel filter", "air filter"},
						"am": {"ፊልተር", "የዘይት ፊልተር", "የነዳጅ ፊልተር", "የአየር ፊልተር"},
						"or": {"filterii", "filterii zayitaa"},
						"ti": {"ፊልተር", "ናይ ዘይቲ ፊልተር"},
					},
				},
			},
		},

		// 19. Root: Vehicle Tires, Wheels, Rims & Traction Systems
		{
			ID:            "root-vehicle-tires-wheels",
			Slug:          "vehicle-tires-wheels-rims",
			Code:          "AUTO_TIRES",
			Level:         1,
			CanonicalName: "Vehicle Tires, Wheels, Rims & Traction Systems",
			Definition:    "Passenger car radial tires, commercial truck tubeless tires, alloy rims, steel wheels, inner tubes, and tire repair patches.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Disc",
			Translations: map[string]string{
				"am": "የተሽከርካሪ ጎማዎች፣ ሪሞች እና ውስጠ-ጎማዎች",
				"or": "Gommaa Konkolaataa fi Riimii",
				"ti": "ናይ መካይን ጎማታትን ሪማትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-radial-passenger-tires",
					Slug:          "passenger-car-radial-tires",
					Code:          "TIRE_PASSENGER",
					Level:         2,
					CanonicalName: "Passenger Car & SUV Radial Tires (የመኪና ጎማዎች)",
					Definition:    "All-season, highway terrain, and all-terrain radial tires for cars and crossovers (13-inch to 22-inch).",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የመንገደኛ መኪና እና የኤስዩቪ ጎማዎች",
						"or": "Gommaa Awutoomoobilii",
						"ti": "ናይ መኪና ጎማታት",
					},
					Aliases: map[string][]string{
						"en": {"tire", "tires", "tyres", "car tire", "bridgestone", "michelin"},
						"am": {"ጎማ", "የመኪና ጎማ", "ጎማዎች", "አዲስ ጎማ"},
						"or": {"gommaa", "gommaa makiinaa"},
						"ti": {"ጎማ", "ናይ መኪና ጎማ"},
					},
				},
			},
		},

		// 20. Root: Aircraft, Aviation Parts & Drone Systems
		{
			ID:            "root-aviation-drones",
			Slug:          "aircraft-aviation-drones-parts",
			Code:          "TRANS_AVIATION",
			Level:         1,
			CanonicalName: "Aircraft, Aviation Parts & Drone Systems",
			Definition:    "Commercial UAVs, agricultural spray drones, aerial inspection quadcopters, avionics, and general aviation components.",
			Status:         "active",
			IsLeaf:         false,
			Icon:           "Plane",
			Translations: map[string]string{
				"am": "አውሮፕላኖች፣ ድሮኖች እና የአቪዬሽን እቃዎች",
				"or": "Xiyyaara, Diroonii fi Meeshaalee Hawaa",
				"ti": "ነፈርቲ፣ ድሮናትን ናይ ኣቪዬሽን ኣቑሑትን",
			},
			Children: []TaxonomyNode{
				{
					ID:            "sf-commercial-drones",
					Slug:          "commercial-uavs-drones",
					Code:          "DRONE_UAV",
					Level:         2,
					CanonicalName: "Commercial Drones & Agricultural Spraying UAVs",
					Definition:    "DJI enterprise drones, crop-spraying multicopters, and topographical photogrammetry UAVs.",
					Status:         "active",
					IsLeaf:         true,
					Translations: map[string]string{
						"am": "የግብርና እና የቪዲዮ ድሮኖች",
						"or": "Diroonii Qonnaa fi Suuraa",
						"ti": "ናይ ሕርሻን ቪዲዮን ድሮናት",
					},
				},
			},
		},
	}
}
