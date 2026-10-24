package main

import "fmt"

// GetCatalogDomains41To50 returns comprehensive domain catalogues for roots 41 to 50
func GetCatalogDomains41To50() []DomainCatalogDef {
	var list []DomainCatalogDef

	makeLeaf := func(slug, en, am, om, ti string) LeafItemDef {
		if ti == "" {
			ti = fmt.Sprintf("ናይ %s", en)
		}
		return LeafItemDef{
			Slug:   slug,
			NameEN: en,
			NameAM: am,
			NameOM: om,
			NameTI: ti,
		}
	}

	// Domain: vehicle-parts-accessories
	list = append(list, DomainCatalogDef{
		RootSlug: "vehicle-parts-accessories",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "engine-brake-transmission-parts",
				NameEN: "Engine, Transmission, Brakes and Filters",
				NameAM: "የሞተር፣ የማርሽ፣ የፍሬን እና የነዳጅ ፊልተሮች",
				NameOM: "Mootora, Firiinii, Giyarii fi Filtarii",
				NameTI: "ናይ ሞተርን ፍሬንን ፊልተርን ኣቑሑት",
				Items: []LeafItemDef{
					makeLeaf("ceramic-semi-metallic-brake-pads", "OEM Heavy-Duty Ceramic & Semi-Metallic Front & Rear Brake Pads", "የፊት እና የኋላ ኦሪጅናል የፍሬን ፓዶች (Brake Pads)", "Baalii Firiinii Konkolaataa Qulqulluu", "ናይ OEM Heavy-Duty Ceramic & Semi-Metallic Front & Rear Brake Pads"),
					makeLeaf("complete-clutch-plate-pressure-kits", "Aisin/Valeo Heavy Duty Clutch Disc, Pressure Plate and Release Bearings", "የክላች ሰሃን፣ ፕሬሸር እና ቤሪንግ ሙሉ ስብስብ", "Kiteetii Kilaachii fi Piresaraa", "ናይ Aisin/Valeo Heavy Duty Clutch Disc, Pressure Plate and Release Bearings"),
					makeLeaf("oil-air-fuel-cabin-filters", "High-Flow Synthetic Engine Oil, Air and Diesel Fuel Filter Elements", "የዘይት፣ የነዳጅ እና የአየር ማጣሪያ ፊልተሮች", "Filtarii Zayitaa, Qilleensaa fi Boba'aa", "ናይ High-Flow Synthetic Engine Oil, Air and Diesel Fuel Filter Elements"),
					makeLeaf("high-performance-spark-glow-plugs", "Iridium Laser Spark Plugs & Fast-Heating Diesel Glow Plugs", "የቤንዚን ካንዴላዎች እና የናፍጣ ሻማዎች (Plugs)", "Kandelaa fi Shaamaa Mootoraa", "ናይ Iridium Laser Spark Plugs & Fast-Heating Diesel Glow Plugs"),
					makeLeaf("hydraulic-gas-shock-absorbers-struts", "Front & Rear Gas-Pressurized Heavy-Duty Shock Absorbers", "የመኪና የፊት እና የኋላ አሞርቲሳተሮች (Shock Absorbers)", "Amoortisaatara Konkolaataa Jabaa", "ናይ Front & Rear Gas-Pressurized Heavy-Duty Shock Absorbers"),
					makeLeaf("heavy-duty-timing-serpentine-belts", "EPDM Reinforced Rubber Engine Timing & Fan Alternator Belts", "የሞተር ታይሚንግ እና የዲናሞ ቀበቶዎች", "Qabattoo Mootoraa fi Tiimiingii", "ናይ EPDM Reinforced Rubber Engine Timing & Fan Alternator Belts"),
					makeLeaf("synthetic-engine-oil-10w40-5w30", "Full Synthetic High-Mileage API SN/CF Engine Motor Oils", "ሙሉ በሙሉ ሲንተቲክ የሞተር ዘይቶች (10W-40 / 5W-30)", "Zayita Mootoraa Qulqulluu", "ናይ Full Synthetic High-Mileage API SN/CF Engine Motor Oils"),
					makeLeaf("car-radiator-aluminum-cooling-cores", "Heavy-Duty Dual Core Aluminum Engine Water Radiators", "የመኪና የውሃ ራዲያተሮች እና ማቀዝቀዣዎች", "Raadiyeetara Mootora Qabbaneessu", "ናይ Heavy-Duty Dual Core Aluminum Engine Water Radiators"),
				},
			},
			{
				Slug:   "tires-rims-batteries-accessories",
				NameEN: "Tires, Alloy Rims, Batteries and Lighting",
				NameAM: "ጎማዎች፣ ሪሞች፣ ባትሪዎች እና መብራቶች",
				NameOM: "Gommaa, Riimii, Baatirii fi Ibsaa",
				NameTI: "ጎማታትን ባትሪታትን ናይ መኪና መብራህትን",
				Items: []LeafItemDef{
					makeLeaf("all-terrain-4x4-suv-tires", "265/65R17 & 285/60R18 Heavy-Ply All-Terrain 4x4 Off-Road Tires", "የ4x4 ጂፕ እና ፒክአፕ ጠንካራ ጎማዎች (All-Terrain)", "Gommaa Konkolaataa 4x4 Jabaa", "ናይ 265/65R17 & 285/60R18 Heavy-Ply All-Terrain 4x4 Off-Road Tires"),
					makeLeaf("commercial-truck-bus-tubeless-tires", "315/80R22.5 & 12.00R20 Heavy Truck Steel Radial Tubeless Tires", "የትላልቅ የጭነት መኪኖች የብረት ራዲያል ጎማዎች (315/80R22.5)", "Gommaa Konkolaataa Gurguddaa Sibilaa", "ናይ 315/80R22.5 & 12.00R20 Heavy Truck Steel Radial Tubeless Tires"),
					makeLeaf("passenger-car-radial-tires-r15-r16", "195/65R15 & 205/55R16 High-Mileage Fuel Efficient Passenger Car Tires", "የአውቶሞቢል መኪናዎች ጥራት ያላቸው ጎማዎች (R15/R16)", "Gommaa Konkolaataa Sasalphaa", "ናይ 195/65R15 & 205/55R16 High-Mileage Fuel Efficient Passenger Car Tires"),
					makeLeaf("maintenance-free-lead-acid-batteries", "12V 60Ah to 100Ah Sealed Maintenance-Free Calcium Car Batteries", "ውሃ የማይፈልጉ 12 ቮልት የመኪና ባትሪዎች (60Ah-100Ah)", "Baatirii Konkolaataa 12V Qulqulluu", "ናይ 12V 60Ah to 100Ah Sealed Maintenance-Free Calcium Car Batteries"),
					makeLeaf("heavy-commercial-truck-batteries-200ah", "12V 150Ah to 200Ah Heavy-Duty SinoTruk & Bus Starting Batteries", "የትላልቅ የጭነት መኪኖች እና አውቶቡሶች ባትሪዎች (150Ah-200Ah)", "Baatirii Konkolaataa Fe'iinsaa 200Ah", "ናይ 12V 150Ah to 200Ah Heavy-Duty SinoTruk & Bus Starting Batteries"),
					makeLeaf("high-lumen-canbus-led-headlight-bulbs", "Mini Projector Lens H4/H7 120W High-Lumen CSP LED Headlight Bulbs", "ጥራት ያላቸው እጅግ ደማቅ የፊት መብራት ኤልኢዲዎች", "Ibsaa Fuulduraa LED Humna Guddaa", "ናይ Mini Projector Lens H4/H7 120W High-Lumen CSP LED Headlight Bulbs"),
					makeLeaf("heavy-aluminum-alloy-car-rims", "15-Inch to 20-Inch Multi-Spoke Reinforced Aluminum Vehicle Wheel Rims", "የአውቶሞቢል እና የጂፕ ውብ የብረት ሪሞች (Alloy Rims)", "Riimii Sibilaa Konkolaataa", "ናይ 15-Inch to 20-Inch Multi-Spoke Reinforced Aluminum Vehicle Wheel Rims"),
					makeLeaf("car-multimedia-android-screen-radios", "9-10 Inch Android Touchscreen GPS Navigation Bluetooth Car Displays", "የመኪና ውስጥ አንድሮይድ ስክሪን፣ ሬዲዮ እና ጂፒኤስ", "Iskiriinii Andirooyidii Konkolaataa", "ናይ 9-10 Inch Android Touchscreen GPS Navigation Bluetooth Car Displays"),
				},
			},
		},
	})

	// Domain: vehicle-rental-logistics
	list = append(list, DomainCatalogDef{
		RootSlug: "vehicle-rental-logistics",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "car-rental-chauffeur-services",
				NameEN: "Car Rental, Field 4WD and Chauffeur Services",
				NameAM: "የመኪና ኪራይ፣ የገጠር 4x4 ጂፖች እና የአሽከርካሪ አገልግሎት",
				NameOM: "Kiree Konkolaataa, 4x4 fi Konkolaachisaa",
				NameTI: "ክራይ መካይንን ናይ 4x4 ጂፓትን",
				Items: []LeafItemDef{
					makeLeaf("field-work-4wd-suv-rentals-un-ngo", "Toyota V8 & Prado Daily / Monthly 4WD Rental for NGO & Field Work", "ለድርጅት እና የመስክ ስራ የሚከራዩ ቶዮታ V8 እና ፕራዶ ጂፖች", "Kiree Konkolaataa 4x4 Dirree fi NGO", "ናይ Toyota V8 & Prado Daily / Monthly 4WD Rental for NGO & Field Work"),
					makeLeaf("luxury-wedding-vip-convoy-car-rentals", "Mercedes S-Class, Limousine & V8 Decorated Wedding Car Hire", "የሰርግ እና የክብር እንግዶች የቅንጦት መኪኖች ኪራይ", "Kiree Konkolaataa Qananii Cidhaa", "ናይ Mercedes S-Class, Limousine & V8 Decorated Wedding Car Hire"),
					makeLeaf("self-drive-city-sedan-car-rentals", "Toyota Vitz & Corolla Daily Self-Drive Fuel-Efficient City Car Rentals", "ለከተማ ውስጥ አሽከርካሪ አልባ የቀን የመኪና ኪራይ", "Kiree Konkolaataa Oofaa Malee (Self-Drive)", "ናይ Toyota Vitz & Corolla Daily Self-Drive Fuel-Efficient City Car Rentals"),
					makeLeaf("chauffeur-driven-airport-pickup-transfers", "Addis Ababa Bole Airport to Hotel 24/7 Professional Chauffeur Pickups", "የቦሌ አውሮፕላን ማረፊያ አቀባበል እና የትራንስፖርት አገልግሎት", "Simannaa Buufata Xiyyaaraa Boolee", "ናይ Addis Ababa Bole Airport to Hotel 24/7 Professional Chauffeur Pickups"),
					makeLeaf("tourist-safari-customized-land-cruisers", "Pop-Up Roof 4x4 Safari Land Cruisers with Experienced Tour Guides", "ጣሪያቸው የሚከፈት የቱሪስት መመልከቻ 4x4 ላንድ ክሩዘሮች ኪራይ", "Konkolaataa Turiizimii fi Saafaarii", "ናይ Pop-Up Roof 4x4 Safari Land Cruisers with Experienced Tour Guides"),
					makeLeaf("monthly-corporate-staff-commute-buses", "Dedicated Driver 30-Seat Staff Morning & Evening Commuter Bus Leases", "የሰራተኞች ማመላለሻ ወርሃዊ የአውቶቡስ ኪራይ ኮንትራት", "Kiree Baasii Hojjattootaa Ji'aa", "ናይ Dedicated Driver 30-Seat Staff Morning & Evening Commuter Bus Leases"),
					makeLeaf("inter-city-one-way-drop-vehicle-rentals", "Addis Ababa to Hawassa, Bahir Dar & Dire Dawa Express Drop-offs", "ከአዲስ አበባ ወደ ክልል ከተሞች የአንድ መንገድ ማድረስ አገልግሎት", "Geejjiba Magaalota Gidduu", "ናይ Addis Ababa to Hawassa, Bahir Dar & Dire Dawa Express Drop-offs"),
					makeLeaf("armored-bulletproof-diplomatic-vehicle-hire", "B6/B7 Ballistic Level Armored Land Cruiser & Escort SUV Rentals", "ጥይት የማይበሳቸው የታጠቁ ዲፕሎማቲክ መኪኖች ኪራይ", "Konkolaataa Nageenyaa Sibilaa (Armored)", "ናይ B6/B7 Ballistic Level Armored Land Cruiser & Escort SUV Rentals"),
				},
			},
			{
				Slug:   "freight-trucking-heavy-haulage",
				NameEN: "Cargo Freight, SinoTruk and Moving Logistics",
				NameAM: "የጭነት ትራንስፖርት፣ የሲኖትራክ እና የከባድ እቃዎች ማጓጓዣ",
				NameOM: "Geejjiba Fe'iinsaa fi Lojistiiksii",
				NameTI: "ናይ ጽዕነት ትራንስፖርትን ምጉዕዓዝን",
				Items: []LeafItemDef{
					makeLeaf("djibouti-addis-ababa-cross-border-container-freight", "Fast Bonded Transit 20ft & 40ft Container Corridor Trucking", "የጅቡቲ ወደብ ኮንቴነሮች ወደ አዲስ አበባ የማጓጓዝ አገልግሎት", "Geejjiba Fe'iinsaa Buufata Jibuutii", "ናይ Fast Bonded Transit 20ft & 40ft Container Corridor Trucking"),
					makeLeaf("sinotruk-construction-sand-stone-haulage", "Crushed Basalt, Selected Material & River Sand Tipper Site Delivery", "የአሸዋ፣ ድንጋይ እና የኮንስትራክሽን ግብዓቶች ማጓጓዝ", "Fe'iinsa Cirrachaa fi Dhagaa SinoTruk", "ናይ Crushed Basalt, Selected Material & River Sand Tipper Site Delivery"),
					makeLeaf("heavy-machinery-lowbed-heavy-transport", "50-80 Ton Lowbed Trailer Mobilization for Excavators & Bulldozers", "ኤክስካቫተሮች እና ዶዘሮች መጫኛ ሎቤድ ተሳቢዎች (Lowbed)", "Geejjiba Maashinii Gurguddaa (Lowbed)", "ናይ 50-80 Ton Lowbed Trailer Mobilization for Excavators & Bulldozers"),
					makeLeaf("household-inter-city-furniture-moving-trucks", "Padded Enclosed Isuzu Truck Household Moving with Laborers Included", "የቤት እና የቢሮ እቃዎችን በጥንቃቄ የማጓጓዝ አገልግሎት", "Konkolaataa Qodaa Manaa Jijjiiru", "ናይ Padded Enclosed Isuzu Truck Household Moving with Laborers Included"),
					makeLeaf("refrigerated-cold-chain-logistics-trucking", "Daily Temperature-Monitored Dairy, Vaccine & Meat Cold Transit", "የስጋ፣ ወተት እና መድኃኒት ማቀዝቀዣ ትራንስፖርት አገልግሎት", "Geejjiba Qabbanaawaa Foonii fi Aanan", "ናይ Daily Temperature-Monitored Dairy, Vaccine & Meat Cold Transit"),
					makeLeaf("express-cargo-parcel-courier-services", "Overnight Addis Ababa to All Regional Capital Parcel Doorstep Delivery", "በ24 ሰዓት ውስጥ ለሁሉም የክልል ከተሞች እቃ የማድረስ አገልግሎት", "Dhaqqabsiisa Baasii fi Paakeejii Dafaa", "ናይ Overnight Addis Ababa to All Regional Capital Parcel Doorstep Delivery"),
					makeLeaf("flatbed-steel-cement-bulk-transport", "Long Haul 40-Ton Flatbed Trucking for Rebar, Coils & Cement Bags", "የብረት እና የሲሚንቶ የረጅም ርቀት የከባድ መኪና ትራንስፖርት", "Fe'iinsa Sibilaa fi Simintoo Tiraayilaraa", "ናይ Long Haul 40-Ton Flatbed Trucking for Rebar, Coils & Cement Bags"),
					makeLeaf("forklift-warehouse-crane-loading-hire", "3-Ton to 10-Ton Diesel Forklift & Mobile Crane Hourly Loading Rentals", "እቃ መጫኛ እና ማውረጃ ፎርክሊፍት እና ክሬን ኪራይ", "Kiree Forklift fi Kireenii Fe'iinsaa", "ናይ 3-Ton to 10-Ton Diesel Forklift & Mobile Crane Hourly Loading Rentals"),
				},
			},
		},
	})

	// Domain: agriculture-farming
	list = append(list, DomainCatalogDef{
		RootSlug: "agriculture-farming",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "tractors-combine-harvesters-equipment",
				NameEN: "Agricultural Tractors, Harvesters and Implements",
				NameAM: "የትራክተር፣ የኮምባይን እና የዘመናዊ እርሻ ማሽኖች",
				NameOM: "Tiraaktara, Maashinii Qonnaa fi Haamtuu",
				NameTI: "ትራክተርን ናይ ሕርሻ መሳርሒታትን",
				Items: []LeafItemDef{
					makeLeaf("4wd-agricultural-wheel-tractors-75hp", "75HP to 120HP Turbocharged 4WD Heavy Soil Cultivation Tractors", "ባለ 4x4 ኃይለኛ የእርሻ ትራክተሮች (75-120 የፈረስ ጉልበት)", "Tiraaktara Qonnaa Humna Guddaa 4WD", "ናይ 75HP to 120HP Turbocharged 4WD Heavy Soil Cultivation Tractors"),
					makeLeaf("self-propelled-grain-combine-harvesters", "Teff, Wheat & Barley Rotary Threshing Grain Combine Harvesters", "የጤፍ እና የስንዴ አጫጅ እና ወቃቂ ትላልቅ ኮምባይኖች", "Maashinii Midhaan Haamu fi Dha'u (Combine)", "ናይ Teff, Wheat & Barley Rotary Threshing Grain Combine Harvesters"),
					makeLeaf("heavy-hydraulic-disc-ploughs", "3-Bottom & 4-Bottom Reversible Heavy Tractor Disc Ploughs", "የትራክተር ማረሻ ባለ ዲስክ እርሻ መሳሪያዎች (Plough)", "Gomaa Qonnaa Tiraaktaraa (Disc Plough)", "ናይ 3-Bottom & 4-Bottom Reversible Heavy Tractor Disc Ploughs"),
					makeLeaf("heavy-duty-disc-harrows", "20-24 Disc Trailed Tandem Soil Pulverizing Tractor Harrows", "አፈር ማድቀቂያ የትራክተር ሀሮው መሳሪያዎች (Harrow)", "Meeshaa Biyyoo Daaku Tiraaktaraa", "ናይ 20-24 Disc Trailed Tandem Soil Pulverizing Tractor Harrows"),
					makeLeaf("walking-two-wheel-mini-tractors-power-tillers", "12HP-15HP Diesel Single-Axle Rotary Power Tillers for Smallholders", "ለአነስተኛ ገበሬዎች የሚሆኑ ባለ ሁለት እግር ሚኒ ትራክተሮች", "Tiraaktara Xiqqaataa Harkaa (Power Tiller)", "ናይ 12HP-15HP Diesel Single-Axle Rotary Power Tillers for Smallholders"),
					makeLeaf("tractor-mounted-seed-drill-planters", "Row Crop Maize, Soybean & Grain Precision Mechanical Planters", "ዘር መዝሪያ ዘመናዊ የትራክተር መሳሪያዎች (Planter)", "Maashinii Sanyii Facaasu Tiraaktaraa", "ናይ Row Crop Maize, Soybean & Grain Precision Mechanical Planters"),
					makeLeaf("boom-sprayers-tractor-mounted", "600L-1000L Chemical Pesticide & Fertilizer 12m Boom Sprayers", "መድኃኒት እና ማዳበሪያ መርጫ የትራክተር ቦም ስፕሬየሮች", "Meeshaa Qoricha Facaasu Tiraaktaraa", "ናይ 600L-1000L Chemical Pesticide & Fertilizer 12m Boom Sprayers"),
					makeLeaf("hydraulic-tipping-farm-tractor-trailers", "5-Ton to 10-Ton Single & Double Axle Hydraulic Dump Farm Trailers", "እህል እና ሰብል መጫኛ የትራክተር ተሳቢ ጋሪዎች", "Gaarii Tiraaktaraa Midhaan Fe'u", "ናይ 5-Ton to 10-Ton Single & Double Axle Hydraulic Dump Farm Trailers"),
				},
			},
			{
				Slug:   "irrigation-greenhouse-farm-inputs",
				NameEN: "Irrigation Systems, Greenhouses and Seedlings",
				NameAM: "የመስኖ እቃዎች፣ የግሪን ሀውስ ሼዶች እና ችግኞች",
				NameOM: "Bishaan Qonnaa, Giriinhoosii fi Biqiltuu",
				NameTI: "ናይ መስኖ ንዋያትን ግሪንሃውስን",
				Items: []LeafItemDef{
					makeLeaf("drip-irrigation-tubing-emitter-kits", "16mm UV-Resistant Pressure Compensating Drip Line Irrigation Rolls", "ቆጣቢ የጠብታ መስኖ ቧንቧዎች እና ኪቶች (Drip Irrigation)", "Ujummoo Bishaan Qorannoo Cophsuu (Drip)", "ናይ 16mm UV-Resistant Pressure Compensating Drip Line Irrigation Rolls"),
					makeLeaf("diesel-engine-centrifugal-irrigation-pumps", "3-Inch & 4-Inch High-Volume River Water Suction Diesel Pumps", "የወንዝ ውሃ መሳቢያ 3 እና 4 ኢንች የናፍጣ የመስኖ ፓምፖች", "Paampii Bishaanii Diiyezilaa Qonnaa", "ናይ 3-Inch & 4-Inch High-Volume River Water Suction Diesel Pumps"),
					makeLeaf("commercial-steel-frame-greenhouse-tunnels", "Galvanized Arch 200-Micron Anti-UV Horticultural Greenhouse Polyhouses", "የአትክልት እና የአበባ ማምረቻ የግሪን ሀውስ ሼዶች", "Giriinhoosii Abaaboo fi Kuduraa", "ናይ Galvanized Arch 200-Micron Anti-UV Horticultural Greenhouse Polyhouses"),
					makeLeaf("agricultural-high-pressure-rain-guns", "360-Degree Large Area Irrigation Rotary Impact Rain Gun Sprinklers", "ሰፊ እርሻን የሚያጠጡ የሚሽከረከሩ ሬይን ገን (Rain Gun) ስፕሪንክለሮች", "Facaastuu Bishaanii Rain Gun", "ናይ 360-Degree Large Area Irrigation Rotary Impact Rain Gun Sprinklers"),
					makeLeaf("grafted-hass-avocado-certified-seedlings", "Certified Disease-Free 1-Year Grafted Export Hass Avocado Saplings", "የተሻሻሉ የተተከሉ የሃስ አቮካዶ ችግኞች", "Biqiltuu Habaqadoo Hass Qulqulluu", "ናይ Certified Disease-Free 1-Year Grafted Export Hass Avocado Saplings"),
					makeLeaf("hybrid-apple-mango-orange-grafted-seedlings", "Tommy Atkins Mango & Washington Navel Citrus Grafted Fruit Seedlings", "የተዳቀሉ የማንጎ፣ የብርቱካን እና የፖም ችግኞች", "Biqiltuu Maangoo fi Birtukaanaa", "ናይ Tommy Atkins Mango & Washington Navel Citrus Grafted Fruit Seedlings"),
					makeLeaf("shade-netting-nursery-cloth-rolls", "50%-80% Sun Blocking Black & Green Agricultural Nursery Shade Nets", "የችግኝ ማፍያ የፀሐይ መከላከያ አረንጓዴ መረቦች", "Saaphana Gaaddisa Biqiltuu", "ናይ 50%-80% Sun Blocking Black & Green Agricultural Nursery Shade Nets"),
					makeLeaf("polyethylene-nursery-seedling-bags", "Perforated Black UV Plastic Seedling Grow Bags (Bulk Packs)", "ችግኝ ማፍያ ቀዳዳ ያላቸው ጥቁር ፌስታሎች (በኪሎ)", "Fastaala Biqiltuu Gurraacha", "ናይ Perforated Black UV Plastic Seedling Grow Bags (Bulk Packs)"),
				},
			},
		},
	})

	// Domain: crops-grains-seeds
	list = append(list, DomainCatalogDef{
		RootSlug: "crops-grains-seeds",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "ethiopian-grains-teff-wheat-barley",
				NameEN: "Teff, Wheat, Barley and Maize Grains",
				NameAM: "ጤፍ፣ ስንዴ፣ ገብስ እና በቆሎ እህሎች",
				NameOM: "Midhaan Xaafii, Qamadii, Garbuu fi Boqqoolloo",
				NameTI: "ናይ ጤፍን ስርናይን ስገምን ምሸላን",
				Items: []LeafItemDef{
					makeLeaf("magna-pure-white-teff-grain", "East Shewa Ada'a Verified Grade-1 Magna Pure White Teff (Per Quintal)", "የአደአ ምርጥ ማኛ ንጹህ ነጭ ጤፍ (በኩንታል)", "Xaafii Maanyaa Ada'aa Qulqulluu", "ናይ East Shewa Ada'a Verified Grade-1 Magna Pure White Teff (Per Quintal)"),
					makeLeaf("sergegna-mixed-teff-grain", "Cleaned Double-Winnowed Sergegna Brown/White Mixed Daily Teff", "የተመረቀ ሰርገኛ ጤፍ ለእንጀራ (በኩንታል)", "Xaafii Sargeenyaa Qulqulluu", "ናይ Cleaned Double-Winnowed Sergegna Brown/White Mixed Daily Teff"),
					makeLeaf("qey-red-teff-iron-rich-grain", "Nutrient-Dense Iron-Rich Red Teff for Health & Diabetics", "ቀይ ጤፍ ከፍተኛ የብረት ማዕድን ያለው (በኩንታል)", "Xaafii Diimaa Fayyaaf Gaarii", "ናይ Nutrient-Dense Iron-Rich Red Teff for Health & Diabetics"),
					makeLeaf("hard-amber-durum-wheat-grain", "Bale & Arsi Highlands High-Protein Amber Durum Pasta Wheat", "የባሌ እና አርሲ ጥራት ያለው የዱረም ስንዴ", "Qamadii Baalee Qulqulluu", "ናይ Bale & Arsi Highlands High-Protein Amber Durum Pasta Wheat"),
					makeLeaf("roasted-kollo-malting-barley", "Plump Kernel Premium Roasted Kollo & Brewery Malting Barley", "የቆሎ እና የቢራ ገብስ እህል በጥራት የተመረተ", "Garbuu Qolloofi Biiraa", "ናይ Plump Kernel Premium Roasted Kollo & Brewery Malting Barley"),
					makeLeaf("yellow-white-maize-grain-corn", "Mechanically Dried Clean Feed & Food Grade Yellow/White Maize", "ደረቅ ነጭ እና ቢጫ በቆሎ (በኩንታል)", "Boqqoolloo Adii fi Keelloo", "ናይ Mechanically Dried Clean Feed & Food Grade Yellow/White Maize"),
					makeLeaf("sorghum-mashila-grain-red-white", "Highland & Lowland Drought-Tolerant Red & White Mashila Grain", "የማሽላ እህል ለዳቦ እና ለመጠጥ", "Midhaan Mishingaa (Maashillaa)", "ናይ Highland & Lowland Drought-Tolerant Red & White Mashila Grain"),
					makeLeaf("finger-millet-dagussa-grain", "Traditional Fermented Beverage & Porridge Mineral-Rich Dagussa", "የዳጉሳ እህል ለገንፎ እና ለጠላ", "Midhaan Daagussaa Aadaa", "ናይ Traditional Fermented Beverage & Porridge Mineral-Rich Dagussa"),
				},
			},
			{
				Slug:   "pulses-oilseeds-spices",
				NameEN: "Oilseeds, Pulses, Lentils and Ethiopian Spices",
				NameAM: "የቅባት እህሎች፣ አተር፣ ምስር፣ ሽምብራ እና ቅመማ ቅመሞች",
				NameOM: "Midhaan Zayitaa, Misira, Shumburaa fi Qimiyyoo",
				NameTI: "ናይ ቅብኣት እኽልን ዓተርን ምስርን ቐመምትን",
				Items: []LeafItemDef{
					makeLeaf("humera-wollega-whitish-sesame-seeds", "Grade-1 Export Quality Pure Whitish Humera & Wollega Sesame Seeds", "የሁመራ እና የወለጋ አንደኛ ደረጃ ነጭ ሰሊጥ", "Sanyiin Saamunaa Saliii Xurree", "ናይ Grade-1 Export Quality Pure Whitish Humera & Wollega Sesame Seeds"),
					makeLeaf("ethiopian-niger-seed-neug-oilseed", "High-Oil Content Machine Cleaned Pure Neug (Niger Seed)", "የኑግ እህል ለንጹህ የምግብ ዘይት መጭመቂያ", "Nuugii Zayita Nyaataa", "ናይ High-Oil Content Machine Cleaned Pure Neug (Niger Seed)"),
					makeLeaf("brown-red-split-lentils-misir", "Debre Berhan Sifted Organic Whole & Split Red Lentils (Misir)", "የደብረ ብርሃን ንጹህ ቀይ ምስር እና ክክ ምስር", "Misira Diimaa fi Baaqelaa", "ናይ Debre Berhan Sifted Organic Whole & Split Red Lentils (Misir)"),
					makeLeaf("chickpeas-shimbra-desi-kabuli", "Large Kernel Kabuli & Desi High-Protein Chickpeas (Shimbra)", "የተመረጡ ነጭ እና ጥቁር ሽምብራ እህሎች", "Shumburaa Gurguddaa Qulqulluu", "ናይ Large Kernel Kabuli & Desi High-Protein Chickpeas (Shimbra)"),
					makeLeaf("faba-horse-beans-baqela", "Graded Ethiopian Highland Faba Beans (Baqela) for Ful & Stew", "የደጋ ንጹህ ባቄላ እና የባቄላ ክክ", "Baaqelaa Qulqulluu (Faba Bean)", "ናይ Graded Ethiopian Highland Faba Beans (Baqela) for Ful & Stew"),
					makeLeaf("ethiopian-red-chili-pepper-berbere", "Mareko Fana Sun-Dried Whole Pods & Ground Authentic Berbere Chili", "የማረቆ ፋና ንጹህ ቀይ ቃሪያ እና የተፈጨ በርበሬ", "Barbaree Diimaa Maaraqoo Faanaa", "ናይ Mareko Fana Sun-Dried Whole Pods & Ground Authentic Berbere Chili"),
					makeLeaf("black-cumin-tikur-azmud-seeds", "Aromatic Medicinal & Baking Pure Ethiopian Black Cumin (Tikur Azmud)", "ንጹህ የጥቁር አዝሙድ ፍሬዎች ለህክምና እና ለዳቦ", "Abasoodaa Gurraacha Qulqulluu", "ናይ Aromatic Medicinal & Baking Pure Ethiopian Black Cumin (Tikur Azmud)"),
					makeLeaf("cardamom-korerima-dried-pods", "Highland Natural Sun-Dried Giant Ethiopian False Cardamom (Korerima)", "ንጹህ የሀበሻ ኮረሪማ ለወጥ እና ለቡና ቅመም", "Qorariimaa Aadaa Qulqulluu", "ናይ Highland Natural Sun-Dried Giant Ethiopian False Cardamom (Korerima)"),
				},
			},
		},
	})

	// Domain: livestock-poultry
	list = append(list, DomainCatalogDef{
		RootSlug: "livestock-poultry",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "cattle-oxen-dairy-cows",
				NameEN: "Oxen, Bulls, Dairy Cows and Calves",
				NameAM: "በሬዎች፣ ሰንጋዎች፣ የወተት ላሞች እና ጥጆች",
				NameOM: "Sangaa, Sa'a Aannanii fi Jabbiilee",
				NameTI: "ናይ ስጋ ብዕራይን ጸባ ላሕምን",
				Items: []LeafItemDef{
					makeLeaf("fattened-holiday-slaughter-oxen-sanga", "Well-Fed 400kg-600kg Harar & Borana Holiday Fattened Bulls (Sanga)", "የበዓል ወፍራም የሀረር እና የቦረና ሰንጋ በሬዎች", "Sangaa Furdaa Ayyaanaa (Harar/Boorana)", "ናይ Well-Fed 400kg-600kg Harar & Borana Holiday Fattened Bulls (Sanga)"),
					makeLeaf("high-yield-holstein-friesian-dairy-cows", "80%-90% Crossbreed 20-30L Daily Milk Holstein Friesian Cows", "በቀን ከ20 እስከ 30 ሊትር ወተት የሚሰጡ የሆልስታይን የወተት ላሞች", "Sa'a Aannanii Faransayii (Holstein)", "ናይ 80%-90% Crossbreed 20-30L Daily Milk Holstein Friesian Cows"),
					makeLeaf("borana-beef-cattle-breeding-heifers", "Drought-Resistant Pure Borana Breed Breeding Bulls & Heifers", "ድርቅ የሚቋቋሙ ንጹህ የቦረና ዝርያ ወይፈኖች እና ጊደሮች", "Korma fi Raada Booranaa", "ናይ Drought-Resistant Pure Borana Breed Breeding Bulls & Heifers"),
					makeLeaf("traditional-plough-trained-oxen-pairs", "Yoke Trained High-Stamina Highland Farm Working Oxen Pairs", "የእርሻ ቀንበር የለመዱ ጠንካራ ጥማድ በሬዎች", "Cimdaa Sangaa Qonnaa Jabaa", "ናይ Yoke Trained High-Stamina Highland Farm Working Oxen Pairs"),
					makeLeaf("weaned-dairy-calves-heifers", "Vaccinated 3-6 Month Pedigree Female Dairy Replacement Calves", "የተከተቡ ከ3 እስከ 6 ወር ጥራት ያላቸው የላም ጥጆች", "Jabbilee Sa'a Aannanii Talaallii Qaban", "ናይ Vaccinated 3-6 Month Pedigree Female Dairy Replacement Calves"),
					makeLeaf("crossbred-jersey-dairy-cows", "High Butterfat Percentage Climate Resilient Jersey Crossbred Cows", "ወተታቸው ከፍተኛ ቅቤ ያለው የጀርሲ ዝርያ ላሞች", "Sa'a Aannanii Jeersii", "ናይ High Butterfat Percentage Climate Resilient Jersey Crossbred Cows"),
					makeLeaf("young-feedlot-fattening-bulls", "Lowland Pastoralist Sourced 18-24 Month Lean Feedlot Stock Bulls", "ለማድለብ የሚሆኑ ወጣት ጤናማ ወይፈኖች", "Goromsa Qophii Furdisaaf", "ናይ Lowland Pastoralist Sourced 18-24 Month Lean Feedlot Stock Bulls"),
					makeLeaf("veterinary-certified-breeding-bulls", "Proven Lineage Artificial Insemination / Natural Service Bulls", "የዘር ማሻሻያ የተረጋገጡ ምርጥ የበሬ ዝርያዎች", "Korma Sanyii Mirkanaa'e", "ናይ Proven Lineage Artificial Insemination / Natural Service Bulls"),
				},
			},
			{
				Slug:   "sheep-goats-poultry",
				NameEN: "Sheep, Goats, Chickens and Day-Old Chicks",
				NameAM: "በጎች፣ ፍየሎች፣ ዶሮዎች እና የጫጩት እርባታ",
				NameOM: "Hoolaa, Re'ee, Lukkuu fi Cuucii",
				NameTI: "ኣባጊዕን ኣጣልን ደራሁትን ጫጩታትን",
				Items: []LeafItemDef{
					makeLeaf("fat-tailed-highland-holiday-sheep", "Menz & Bonga Heavy Wool Fat-Tailed Holiday Slaughter Rams", "የመንዝ እና የቦንጋ ወፍራም የበዓል ሙክት በጎች", "Hoolaa Ayyaanaa Furdaa (Menz/Bonga)", "ናይ Menz & Bonga Heavy Wool Fat-Tailed Holiday Slaughter Rams"),
					makeLeaf("borana-somali-long-eared-meat-goats", "Export Grade Clean White & Brown Meat Buck Goats", "የቦረና እና የሶማሊ ምርጥ የስጋ ፍየሎች", "Re'ee Foonii Booranaa fi Somaalee", "ናይ Export Grade Clean White & Brown Meat Buck Goats"),
					makeLeaf("sasso-dual-purpose-chickens", "Hardy 2-Month Free-Range Dual Meat & Egg Sasso Chickens", "የሳሶ (Sasso) የስጋ እና የእንቁላል ዶሮዎች", "Lukkuu Saasoo Foonii fi Buuphaa", "ናይ Hardy 2-Month Free-Range Dual Meat & Egg Sasso Chickens"),
					makeLeaf("day-old-broiler-chick-boxes", "Vaccinated Fast-Growing Cobb 500 Day-Old Commercial Broiler Chicks", "በቀን የሚደርሱ የኮብ 500 የስጋ ዶሮ ጫጩቶች (በካርቶን)", "Cuucii Lukkuu Foonii Guyyaa Tokkoo (Cobb 500)", "ናይ Vaccinated Fast-Growing Cobb 500 Day-Old Commercial Broiler Chicks"),
					makeLeaf("day-old-layer-pullet-chicks", "High-Yield 300+ Egg Annual Bovans Brown Day-Old Layer Chicks", "በዓመት ከ300 በላይ እንቁላል የሚጥሉ የቦቫንስ ጫጩቶች", "Cuucii Lukkuu Buuphaa Guyyaa Tokkoo (Bovans)", "ናይ High-Yield 300+ Egg Annual Bovans Brown Day-Old Layer Chicks"),
					makeLeaf("fertile-fresh-hatching-eggs-crates", "Tested 85%+ Hatchability Purebreed Sasso & Bovans Hatching Eggs", "ለመፈልፈያ የሚሆኑ ምርጥ የዶሮ እንቁላሎች በክሬት", "Buuphaa Lukkuu Ilbiisaf Qophaa'e", "ናይ Tested 85%+ Hatchability Purebreed Sasso & Bovans Hatching Eggs"),
					makeLeaf("automatic-egg-incubator-hatchers", "500 to 5000 Egg Digital Temperature Humidity Turning Incubators", "ከ500 እስከ 5000 እንቁላል የሚፈልፍሉ አውቶማቲክ ኢንኩቤተሮች", "Maashinii Buuphaa Lukkuu Ilbiisu (Incubator)", "ናይ 500 to 5000 Egg Digital Temperature Humidity Turning Incubators"),
					makeLeaf("multi-tier-commercial-layer-chicken-cages", "Galvanized Steel Automatic Nipple Watered Battery Poultry Cages", "የዶሮ እርባታ ባለ ደረጃ ዘመናዊ የብረት ካኬቶች (Cages)", "Koorjii Lukkuu Sibilaa (Layer Cage)", "ናይ Galvanized Steel Automatic Nipple Watered Battery Poultry Cages"),
				},
			},
		},
	})

	// Domain: animal-feed-supplies
	list = append(list, DomainCatalogDef{
		RootSlug: "animal-feed-supplies",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "livestock-poultry-feed-nutrition",
				NameEN: "Animal Feed, Poultry Mash and Supplements",
				NameAM: "የከብት እና የዶሮ መኖ፣ ፋጉሎ፣ መኖ እና ማዕድናት",
				NameOM: "Nyaata Horii, Lukkuu fi Dabalata",
				NameTI: "ናይ ከብትን ደራሁትን መግቢን ፋጉሎን",
				Items: []LeafItemDef{
					makeLeaf("commercial-broiler-finisher-starter-feed", "Pelletized 20% Crude Protein Fast-Growth Commercial Broiler Feed", "የስጋ ዶሮ ማስጀመሪያ እና ማጠናቀቂያ መኖ በኩንታል", "Nyaata Lukkuu Foonii (Broiler Feed)", "ናይ Pelletized 20% Crude Protein Fast-Growth Commercial Broiler Feed"),
					makeLeaf("commercial-layer-mash-feed-sacks", "High-Calcium Strong Shell Production Commercial Layer Mash Feed", "የእንቁላል ዶሮ ጥራት ያለው መኖ በኩንታል", "Nyaata Lukkuu Buuphaa (Layer Mash)", "ናይ High-Calcium Strong Shell Production Commercial Layer Mash Feed"),
					makeLeaf("noug-cake-oilmeal-cattle-fattening", "High-Protein 30% Pure Niger Seed Noug Cake (Fagulo) Sacks", "የከብት ማድለቢያ ንጹህ የኑግ ፋጉሎ (በኩንታል)", "Nyaata Horii Faaguloo Nuugii", "ናይ High-Protein 30% Pure Niger Seed Noug Cake (Fagulo) Sacks"),
					makeLeaf("wheat-bran-frusheka-feed-bags", "Clean Flaked Coarse Wheat Bran (Frusheka) Dairy Cow Feed Sacks", "የስንዴ ፍሩሼካ ለወተት ላሞች እና ለከብቶች በኩንታል", "Furusheekaa Qamadii Horii", "ናይ Clean Flaked Coarse Wheat Bran (Frusheka) Dairy Cow Feed Sacks"),
					makeLeaf("mineral-multivitamin-salt-lick-blocks", "5kg Compressed Trace Element Mineral Livestock Salt Lick Blocks", "የከብቶች መላሻ የማዕድን ጨው ብሎኮች (Salt Lick)", "Dhagaa Ashaboo Horii Arrabsan (Salt Lick)", "ናይ 5kg Compressed Trace Element Mineral Livestock Salt Lick Blocks"),
					makeLeaf("compacted-hay-teff-straw-bales", "Dry Baled Alfalfa, Rhodes Grass & Teff Straw (Chid) Livestock Feed", "የታሰረ የሳር እና የጤፍ ጭድ መኖ በጥቅል (Bale)", "Cidii Xaafii fi Muka Horii", "ናይ Dry Baled Alfalfa, Rhodes Grass & Teff Straw (Chid) Livestock Feed"),
					makeLeaf("alfalfa-lucerne-high-protein-pellets", "Dehydrated 18% Protein Sun-Cured Alfalfa Cattle Pellets", "የአልፋልፋ ፕሮቲን የበለጸገ ዘመናዊ መኖ", "Nyaata Horii Alfaalfaa (Alfalfa Pellets)", "ናይ Dehydrated 18% Protein Sun-Cured Alfalfa Cattle Pellets"),
					makeLeaf("brewery-spent-grain-wet-atela", "Nutrient-Dense Wet Industrial Brewery Spent Barley Grains (Atela)", "የቢራ ፋብሪካ ተረፈ-ምርት አተላ ለከብት ማድለቢያ", "Atalaa Biiraa Horii Furdusu", "ናይ Nutrient-Dense Wet Industrial Brewery Spent Barley Grains (Atela)"),
				},
			},
			{
				Slug:   "veterinary-drugs-farm-supplies",
				NameEN: "Veterinary Drugs, Vaccines and Farm Supplies",
				NameAM: "የእንስሳት መድኃኒቶች፣ ክትባቶች እና እርሻ መገልገያዎች",
				NameOM: "Qoricha Horii, Talaallii fi Meeshaalee",
				NameTI: "ናይ እንስሳት ፈውስን ክትባትን",
				Items: []LeafItemDef{
					makeLeaf("broad-spectrum-ivermectin-dewormer-injections", "1% Sterile Veterinary Parasite & Tick Control Ivermectin Injections", "የእንስሳት ትል እና መዥገር ማጥፊያ አቨርሜክቲን መርፌ", "Qoricha Horii Raammoo fi Silmii (Ivermectin)", "ናይ 1% Sterile Veterinary Parasite & Tick Control Ivermectin Injections"),
					makeLeaf("oxytetracycline-long-acting-antibiotics", "20% Long-Acting Broad Spectrum Livestock Bacterial Antibiotics", "የእንስሳት አጠቃላይ የባክቴሪያ መከላከያ ኦክሲቴትራሳይክሊን መርፌ", "Qoricha Bakteeriyaa Horii (Oxytetracycline)", "ናይ 20% Long-Acting Broad Spectrum Livestock Bacterial Antibiotics"),
					makeLeaf("multivitamin-amino-acid-livestock-injections", "High-Potency Recovery & Appetite Boosting Livestock Vitamin Injections", "የእንስሳት የምግብ ፍላጎት መክፈቻ እና ማጠናከሪያ ቫይታሚን", "Vaayitaaminii fi Qoricha Horii", "ናይ High-Potency Recovery & Appetite Boosting Livestock Vitamin Injections"),
					makeLeaf("continuous-veterinary-vaccination-syringes", "Adjustable 2ml-5ml Automatic Self-Refilling Metal Livestock Syringes", "የከብቶች እና የዶሮ መከተቢያ አውቶማቲክ የብረት ሲሪንጆች", "Lilmoo Talaallii Horii Sibilaa", "ናይ Adjustable 2ml-5ml Automatic Self-Refilling Metal Livestock Syringes"),
					makeLeaf("electric-dairy-milking-machines-vacuum", "Pulsator Single & Double Bucket Portable Cow Milking Machines", "የላም ወተት መዛቢያ የኤሌክትሪክ ማሽኖች (Milking Machine)", "Maashinii Aanan Elmu Elektiriikii", "ናይ Pulsator Single & Double Bucket Portable Cow Milking Machines"),
					makeLeaf("heavy-duty-plastic-poultry-drinkers-feeders", "10L Inverted Bell Gravity Chicken Drinkers & Hanging Feeders", "የዶሮ ውሃ እና መኖ ማስቀመጫ ፕላስቲኮች", "Meeshaa Bishaan fi Nyaata Lukkuu", "ናይ 10L Inverted Bell Gravity Chicken Drinkers & Hanging Feeders"),
					makeLeaf("fly-tick-livestock-pour-on-solutions", "Deltamethrin External Ectoparasite Cattle Pour-On Liquid Dip", "የከብቶች መዥገር እና ዝንብ መከላከያ ኬሚካሎች", "Qoricha Silmii fi Titiisa Horii", "ናይ Deltamethrin External Ectoparasite Cattle Pour-On Liquid Dip"),
					makeLeaf("castration-burdizzo-pliers-veterinary", "Stainless Steel Bloodless Veterinary Livestock Castrator Clamps", "የከብት ማኮላሻ የብረት ቡርዲዞ መቆንጠጫዎች (Burdizzo)", "Buurdiizoo Korma Kolaasu (Burdizzo)", "ናይ Stainless Steel Bloodless Veterinary Livestock Castrator Clamps"),
				},
			},
		},
	})

	// Domain: food-groceries
	list = append(list, DomainCatalogDef{
		RootSlug: "food-groceries",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "fresh-vegetables-fruits-potatoes",
				NameEN: "Fresh Vegetables, Onions, Potatoes and Fruits",
				NameAM: "ትኩስ አትክልቶች፣ ሽንኩርት፣ ድንች እና ፍራፍሬዎች",
				NameOM: "Kuduraa, Qullubbii, Dinnicha fi Fuduraa",
				NameTI: "ሓደሽቲ ኣሕምልትን ሽጉርትን ድንሽትን ፍረታትን",
				Items: []LeafItemDef{
					makeLeaf("ethiopian-red-onions-shinkurt-sacks", "Highland Sun-Cured Firm Red Cooking Onions 100kg Mesh Sacks", "የሀበሻ ቀይ ሽንኩርት በጆንያ (100 ኪሎ)", "Qullubbii Diimaa Itoophiyaa (Kiiloo 100)", "ናይ Highland Sun-Cured Firm Red Cooking Onions 100kg Mesh Sacks"),
					makeLeaf("shashemene-clean-table-potatoes-sacks", "Smooth Skin Cream Flesh Fresh Table Potatoes (100kg Sacks)", "የሻሸመኔ ንጹህ ድንች በጅምላ (በጆንያ)", "Dinnicha Qulqulluu Shashemene (Kiiloo 100)", "ናይ Smooth Skin Cream Flesh Fresh Table Potatoes (100kg Sacks)"),
					makeLeaf("meki-ziway-fresh-salad-tomatoes", "Firm Red Greenhouse Salad Tomatoes in Standard Wooden Crates", "የመቂ እና ዝዋይ ትኩስ ቲማቲም በሳጥን", "Timaatima Qulqulluu Meqii fi Ziway", "ናይ Firm Red Greenhouse Salad Tomatoes in Standard Wooden Crates"),
					makeLeaf("ethiopian-garlic-nech-shinkurt", "Spicy High-Aroma Small Bulb Highland Natural Garlic (Nech Shinkurt)", "የሀበሻ ነጭ ሽንኩርት በጅምላ", "Qullubbii Adii Itoophiyaa Qulqulluu", "ናይ Spicy High-Aroma Small Bulb Highland Natural Garlic (Nech Shinkurt)"),
					makeLeaf("fresh-green-cabbage-tikil-gomen", "Tight Leaf Heavy Head Fresh Green Cabbage (Tikil Gomen)", "ትኩስ ጥቅል ጎመን በጅምላ", "Tibil Gomenii (Raafuu Geengoo)", "ናይ Tight Leaf Heavy Head Fresh Green Cabbage (Tikil Gomen)"),
					makeLeaf("habesha-collard-greens-yeabesha-gomen", "Fresh Tender Dark Green Organic Ethiopian Collard Greens", "የሀበሻ ቆስጣ እና የሀበሻ ጎመን በጥቅል", "Raafuu Aadaa Itoophiyaa", "ናይ Fresh Tender Dark Green Organic Ethiopian Collard Greens"),
					makeLeaf("fresh-ripe-hass-avocados-crates", "Butter Texture Rich Oil Creamy Green Hass Avocados in Wooden Boxes", "ቅቤ የሆኑ የሃስ አቮካዶ ፍራፍሬዎች በሳጥን", "Habaqadoo Hass Qulqulluu", "ናይ Butter Texture Rich Oil Creamy Green Hass Avocados in Wooden Boxes"),
					makeLeaf("sweet-apple-bananas-habesha-muz", "Naturally Ripened Sweet Yellow Apple Bananas (Muz)", "የሀበሻ ጣፋጭ ሙዝ በዘለላ", "Muzii Mi'aawaa Aadaa", "ናይ Naturally Ripened Sweet Yellow Apple Bananas (Muz)"),
				},
			},
			{
				Slug:   "meat-fish-dairy-traditional-honey",
				NameEN: "Meat, Fish, Dairy and Pure Organic Honey",
				NameAM: "ስጋ፣ ዓሳ፣ የወተት ተዋጽኦዎች እና ንጹህ የተፈጥሮ ማር",
				NameOM: "Foon, Qurxummii, Aanan fi Damma Aadaa",
				NameTI: "ስጋን ዓሳን ጸባን ማርን",
				Items: []LeafItemDef{
					makeLeaf("fresh-bone-in-beef-carcasses", "Freshly Slaughtered Municipal Abattoir Inspected Beef Quarters", "የቄራ የታረደ ንጹህ የበሬ ስጋ በኪሎ", "Foon Sangaa Qulqulluu Qalame", "ናይ Freshly Slaughtered Municipal Abattoir Inspected Beef Quarters"),
					makeLeaf("fresh-lake-tilapia-fish-fillets", "Lake Tana & Hawassa Cleaned Bone-Free Fresh Tilapia Fillets", "የጣና እና የሀዋሳ ትኩስ የቆረሶ ዓሳ (Tilapia)", "Qurxummii Xilaaphiyaa Qulqulluu", "ናይ Lake Tana & Hawassa Cleaned Bone-Free Fresh Tilapia Fillets"),
					makeLeaf("traditional-clarified-spiced-butter-nitir-qibe", "Aromatic Spiced Organic Ghee Clarified Butter (Nitir Qibe)", "የተነጠረ የሀበሻ ቅቤ ከነቅመሙ (ንጥር ቅቤ)", "Dhadhaa Nyaataa Qulqulluu (Niitir Qibee)", "ናይ Aromatic Spiced Organic Ghee Clarified Butter (Nitir Qibe)"),
					makeLeaf("fresh-ayib-ethiopian-cottage-cheese", "Soft White Whole Milk Traditional Ethiopian Cottage Cheese (Ayib)", "ንጹህ የሀበሻ አይብ በኪሎ", "Ayibii Aadaa Aanan Qulqulluu", "ናይ Soft White Whole Milk Traditional Ethiopian Cottage Cheese (Ayib)"),
					makeLeaf("pure-tigray-white-honey", "Pure Natural High-Altitude White Blossom Tigray Honey (Mar)", "የትግራይ ንጹህ ነጭ ማር", "Damma Adii Naannoo Tigraay", "ናይ Pure Natural High-Altitude White Blossom Tigray Honey (Mar)"),
					makeLeaf("yellow-amber-forest-honey", "Raw Wild Forest Amber Honey from Kaffa & Jimma Highlands", "የከፋ እና የጅማ የተፈጥሮ የጫካ ማር", "Damma Bosonaa Kaffa fi Jimmaa", "ናይ Raw Wild Forest Amber Honey from Kaffa & Jimma Highlands"),
					makeLeaf("pasteurized-fresh-cow-milk-pouches", "500ml Cold Chain Pasteurized Whole Cow Milk Pouches (Crate of 20)", "የታሸገ ንጹህ ላም ወተት በፓኬት", "Aanan Sa'aa Qulqulluu Saamunaan", "ናይ 500ml Cold Chain Pasteurized Whole Cow Milk Pouches (Crate of 20)"),
					makeLeaf("fresh-brown-organic-chicken-eggs", "Farm Collected Fresh Yolk Organic Brown Table Eggs (30-Egg Crates)", "ትኩስ የሀበሻ እና የፈረንጅ የዶሮ እንቁላሎች በክሬት", "Buuphaa Lukkuu Qulqulluu (Kireetii 30)", "ናይ Farm Collected Fresh Yolk Organic Brown Table Eggs (30-Egg Crates)"),
				},
			},
		},
	})

	// Domain: restaurants-catering
	list = append(list, DomainCatalogDef{
		RootSlug: "restaurants-catering",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "traditional-ethiopian-restaurant-catering",
				NameEN: "Traditional Restaurant Catering and Buffets",
				NameAM: "የባህል ምግብ ዝግጅት፣ የሰርግ ድግስ እና የቡፌ አገልግሎቶች",
				NameOM: "Qophii Nyaata Aadaa, Cidhaa fi Bufee",
				NameTI: "ባህላዊ ምግቢን ናይ መርዓ ድግስን ቡፌን",
				Items: []LeafItemDef{
					makeLeaf("wedding-holiday-traditional-buffet-catering", "Full Traditional Buffet (Doro Wat, Siga Wat, Kitfo & Injera) Catering", "የሰርግ እና የበዓል የተሟላ የባህል ቡፌ ምግብ ዝግጅት (የዶሮ ወጥ፣ ክትፎ)", "Qophii Nyaata Cidhaa Guutuu (Doro Wat, Kitfo)", "ናይ Full Traditional Buffet (Doro Wat, Siga Wat, Kitfo & Injera) Catering"),
					makeLeaf("gourmet-gurage-kitfo-catering-services", "Fresh Mفرድ Minced Beef Kitfo with Kocho, Gomen & Ayib Buffets", "የጉራጌ ልዩ ክትፎ ከነኮቾ፣ አይብ እና ጎመን ዝግጅት", "Qophii Kitfoo Aadaa fi Qoccoo", "ናይ Fresh Mفرድ Minced Beef Kitfo with Kocho, Gomen & Ayib Buffets"),
					makeLeaf("doro-wat-holiday-catering-pots", "Authentic 12-Cut Country Chicken & Boiled Egg Holiday Doro Wat Pots", "የበዓል እና የሰርግ የዶሮ ወጥ ዝግጅት በድስት", "Qophii Doro Wotii Ayyaanaa", "ናይ Authentic 12-Cut Country Chicken & Boiled Egg Holiday Doro Wat Pots"),
					makeLeaf("vegan-fasting-ye-tsom-buffet-catering", "10-Dish Ethiopian Orthodox Lent & Fasting Veggie Combo Buffets", "የጾም የባህል ምግብ ቡፌ ዝግጅት (የጾም በያይነቱ)", "Qophii Nyaata Soomaa (Ye Tsom Beyaynetu)", "ናይ 10-Dish Ethiopian Orthodox Lent & Fasting Veggie Combo Buffets"),
					makeLeaf("corporate-lunch-box-catering-delivery", "Hygienic Foil-Packed Balanced Executive Lunch Boxes for Conferences", "ለስብሰባ እና ለስልጠና የሚሆኑ የታሸጉ የምሳ ሳጥኖች", "Sanduuqa Nyaata Qophaa'ee Biiroof (Lunch Box)", "ናይ Hygienic Foil-Packed Balanced Executive Lunch Boxes for Conferences"),
					makeLeaf("roast-sheep-whole-carcass-catering", "Open-Pit Charcoal Roasted Whole Spiced Sheep (Yebeg Shekla)", "ሙሉ የተጠበሰ የበግ ስጋ ድግስ ዝግጅት (የበግ ጥብስ)", "Qophii Hoolaa Guutuu Waxalamee", "ናይ Open-Pit Charcoal Roasted Whole Spiced Sheep (Yebeg Shekla)"),
					makeLeaf("traditional-kocho-bulla-dish-preparations", "Fermented Enset Root Kocho Flatbreads & Bulla Porridge Services", "የተዘጋጀ የኮቾ እና የቡላ ባህላዊ ምግቦች", "Qophii Nyaata Qoccoo fi Buullaa", "ናይ Fermented Enset Root Kocho Flatbreads & Bulla Porridge Services"),
					makeLeaf("customized-wedding-cake-creations", "3-Tier to 5-Tier Fondant Decorated Gourmet Flavor Wedding Cakes", "ባለ 3 እና 5 ደርብ ውብ የሰርግ ኬኮች ዝግጅት", "Keekii Cidhaa Gurguddaa Qananii", "ናይ 3-Tier to 5-Tier Fondant Decorated Gourmet Flavor Wedding Cakes"),
				},
			},
			{
				Slug:   "pastry-baking-fast-food-supplies",
				NameEN: "Bakery Breads, Pastries and Fast Food Supplies",
				NameAM: "የዳቦ፣ የኬክ፣ የበርገር እና የፒዛ ዝግጅቶች",
				NameOM: "Daabboo, Keekii fi Nyaata Dafaa",
				NameTI: "ናይ ባኒን ኬክን በርገርን ፒዛን",
				Items: []LeafItemDef{
					makeLeaf("artisan-sourdough-french-baguette-breads", "Daily Baked Crispy Crust European Sourdough & Long French Baguettes", "ትኩስ የፈረንሳይ ባጌት እና የጣሊያን ዳቦዎች", "Daabboo Baageettii Faransaayii", "ናይ Daily Baked Crispy Crust European Sourdough & Long French Baguettes"),
					makeLeaf("custom-birthday-celebration-cakes", "Personalized Photo Print & Cream Chocolate Birthday Party Cakes", "የልደት እና የደስታ ቀን ኬኮች በስም እና ፎቶ የሚሰሩ", "Keekii Ayyaana Dhalootaa Bareedaa", "ናይ Personalized Photo Print & Cream Chocolate Birthday Party Cakes"),
					makeLeaf("commercial-burger-patties-buns-packs", "Frozen Seasoned Beef Patties & Sesame Seed Soft Burger Buns", "የበርገር ስጋ እና የበርገር ዳቦዎች በፓኬት", "Foon fi Daabboo Boorgarii (Burger)", "ናይ Frozen Seasoned Beef Patties & Sesame Seed Soft Burger Buns"),
					makeLeaf("authentic-italian-wood-fired-pizza-crusts", "Pre-Baked Thin Crust Mozzarella Ready Round Pizza Dough Shells", "የተዘጋጁ የፒዛ ሊጦች እና ፒዛዎች", "Daabboo Piizaa Xaaliyaanii", "ናይ Pre-Baked Thin Crust Mozzarella Ready Round Pizza Dough Shells"),
					makeLeaf("croissants-danish-puff-pastries", "Layered Butter Flaky French Croissants & Fruit Danish Pastries", "የቅቤ ክሩሳንቶች እና የፈረንሳይ ኬኮች (Croissant)", "Kuruwaasanii fi Paastirii", "ናይ Layered Butter Flaky French Croissants & Fruit Danish Pastries"),
					makeLeaf("sweet-donuts-cupcakes-catering-boxes", "Glazed Sprinkled Assorted Box of 12 Gourmet Party Doughnuts", "የዶናት እና የካፕኬክ ጣፋጮች በካርቶን", "Doonatii fi Mi'aawaa", "ናይ Glazed Sprinkled Assorted Box of 12 Gourmet Party Doughnuts"),
					makeLeaf("freshly-pressed-tropical-fruit-juices", "Avocado, Mango, Papaya & Guava Layered Mixed Special Spris Juices", "የአቮካዶ፣ ማንጎ እና ፓፓያ እስፕሪስ ጭማቂዎች", "Juusii Isphiriisii Fuduraa (Spris)", "ናይ Avocado, Mango, Papaya & Guava Layered Mixed Special Spris Juices"),
					makeLeaf("fried-samosas-spring-rolls-frozen", "Lentil, Meat & Vegetable Crispy Fried Triangular Samosa Packs", "የስጋ እና የምስር ሳምቡሳዎች (የተጠበሱ እና ፍሮዝን)", "Saambuusaa Foonii fi Misiraa", "ናይ Lentil, Meat & Vegetable Crispy Fried Triangular Samosa Packs"),
				},
			},
		},
	})

	// Domain: coffee-tea-beverages
	list = append(list, DomainCatalogDef{
		RootSlug: "coffee-tea-beverages",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "specialty-ethiopian-coffee-beans",
				NameEN: "Specialty Arabica Raw and Roasted Coffee Beans",
				NameAM: "የተመረጡ የኢትዮጵያ አረቢካ ቡናዎች (የተቆሉ እና ጥሬ)",
				NameOM: "Buna Arabica Itoophiyaa Filatamaa",
				NameTI: "ናይ ኢትዮጵያ ፍሉይ ቡን",
				Items: []LeafItemDef{
					makeLeaf("yirgacheffe-washed-grade-1-coffee", "Jasmine Floral & Citrus Notes Grade-1 Washed Yirgacheffe Arabica Beans", "የይርጋጨፌ አንደኛ ደረጃ የታጠበ ልዩ ቡና", "Buna Yirgaacaffee Sadarkaa 1ffaa", "ናይ Jasmine Floral & Citrus Notes Grade-1 Washed Yirgacheffe Arabica Beans"),
					makeLeaf("sidama-natural-specialty-coffee", "Berry Jam & Winey Sweetness Natural Sun-Dried Sidama Arabica Beans", "የሲዳማ የተፈጥሮ ፀሐይ የደረቀ ልዩ ቡና", "Buna Sidaamaa Aadaa Mi'aawaa", "ናይ Berry Jam & Winey Sweetness Natural Sun-Dried Sidama Arabica Beans"),
					makeLeaf("harar-longberry-wild-coffee-beans", "Wild Mocha Blueberry Notes Sun-Dried Natural Harar Longberry Beans", "የሀረር ምርጥ ሎንግበሪ የዱር ቡና", "Buna Harar Longberry Filatamaa", "ናይ Wild Mocha Blueberry Notes Sun-Dried Natural Harar Longberry Beans"),
					makeLeaf("limu-jimma-highland-arabica-coffee", "Balanced Low Acidity Sweet Brown Sugar Highland Limu/Jimma Beans", "የሊሙ እና የጅማ ጥራት ያለው የደጋ ቡና", "Buna Liimuu fi Jimmaa", "ናይ Balanced Low Acidity Sweet Brown Sugar Highland Limu/Jimma Beans"),
					makeLeaf("guji-natural-specialty-grade-coffee", "Peach, Apricot & Dark Chocolate Specialty Sun-Dried Guji Beans", "የጉጂ ልዩ የተፈጥሮ ጣዕም ያለው ቡና", "Buna Gujii Addunyaarratti Beekamaa", "ናይ Peach, Apricot & Dark Chocolate Specialty Sun-Dried Guji Beans"),
					makeLeaf("dark-roast-espresso-ground-coffee-bags", "Freshly Roasted & Vacuum Sealed 500g Italian Espresso Blend Coffee", "የተቆላ እና የተፈጨ ጥራት ያለው የኤስፕሬሶ ቡና (500 ግራም)", "Buna Akaawamee fi Daakame (500g)", "ናይ Freshly Roasted & Vacuum Sealed 500g Italian Espresso Blend Coffee"),
					makeLeaf("traditional-ethiopian-green-raw-coffee-quintals", "Export Standard 60kg Jute Bag Raw Green Arabica Coffee Quintals", "ያልተቆላ ጥሬ አረንጓዴ ቡና በኩንታል (በጆንያ)", "Buna Dheedhii Xurree (Kiiloo 60)", "ናይ Export Standard 60kg Jute Bag Raw Green Arabica Coffee Quintals"),
					makeLeaf("kefa-bonga-organic-forest-coffee", "Certified Biosphere Reserve Wild Organic Forest Arabica Beans", "የከፋ ቦንጋ የደን የተፈጥሮ ኦርጋኒክ ቡና", "Buna Bosonaa Kaaffa fi Boongaa", "ናይ Certified Biosphere Reserve Wild Organic Forest Arabica Beans"),
				},
			},
			{
				Slug:   "ethiopian-tea-traditional-drinks",
				NameEN: "Herbal Tea, Spices and Traditional Beverages (Tej/Tella)",
				NameAM: "ሻይ ቅጠሎች፣ ጠጅ፣ ጠላ እና ባህላዊ መጠጦች",
				NameOM: "Shaayii, Xajjii, Dhadhii fi Dhugaatii Aadaa",
				NameTI: "ናይ ሻሂ ቆጽልን ሜስን ስዋን",
				Items: []LeafItemDef{
					makeLeaf("gumero-wushwush-black-tea-leaves", "Highland Single-Origin Black Tea Leaves & Tea Bags (Wushwush)", "የዉሽዉሽ እና ጉመሮ ጥራት ያላቸው ጥቁር የሻይ ቅጠሎች", "Shaayii Gurraacha Wushwush", "ናይ Highland Single-Origin Black Tea Leaves & Tea Bags (Wushwush)"),
					makeLeaf("pure-ethiopian-honey-mead-tej", "Golden Fermented Organic Honey & Gesho Traditional Ethiopian Tej", "የተመረቀ ንጹህ የማር ጠጅ በመጠቅለያ (ጠርሙስ)", "Xajjii Damma Qulqulluu (Tej)", "ናይ Golden Fermented Organic Honey & Gesho Traditional Ethiopian Tej"),
					makeLeaf("traditional-filtered-ethiopian-tella", "Barley, Teff & Gesho Leaf Home-Brewed Traditional Filtered Tella", "የተጣራ የሀበሻ ባህላዊ ጠላ", "Farsoo fi Dhadhii Aadaa (Tella)", "ናይ Barley, Teff & Gesho Leaf Home-Brewed Traditional Filtered Tella"),
					makeLeaf("dried-gesho-leaves-stems-hops", "Sun-Dried Rhamnus Prinoides Gesho Leaves & Stems for Brewing", "የደረቀ የጌሾ ቅጠል እና እንጨት ለጠጅ እና ለጠላ", "Baala Geeshoo Aadaa", "ናይ Sun-Dried Rhamnus Prinoides Gesho Leaves & Stems for Brewing"),
					makeLeaf("ethiopian-spiced-tea-blend-chais", "Cinnamon, Cardamom, Clove & Ginger Scented Highland Black Tea", "ቀረፋ፣ ቅርንፉድ እና ጅንጅብል የተቀላቀለበት የሀበሻ ቅመም ሻይ", "Shaayii Qimiyyoo Aadaa", "ናይ Cinnamon, Cardamom, Clove & Ginger Scented Highland Black Tea"),
					makeLeaf("herbal-hibiscus-karkadeh-tea-flowers", "Deep Red Dried Hibiscus Roselle Flower Herbal Calming Tea", "የከርከዴህ ቀይ የሻይ አበባዎች", "Shaayii Abaaboo Karkadee", "ናይ Deep Red Dried Hibiscus Roselle Flower Herbal Calming Tea"),
					makeLeaf("fresh-green-moringa-powder-tea", "Nutrient-Dense 100% Pure Organic Dried Moringa Leaf Powders", "የሞሪንጋ የተፈጥሮ የጤና ሻይ ዱቄት", "Daakuu Moringaa Fayyaaf Gaarii", "ናይ Nutrient-Dense 100% Pure Organic Dried Moringa Leaf Powders"),
					makeLeaf("traditional-barley-besso-drink-mix", "Roasted Barley Flour (Besso) Powder Mix for Instant Shakes", "የገብስ በሶ ዱቄት ለመጠጥ እና ለቁርስ", "Daakuu Basso Garbuu", "ናይ Roasted Barley Flour (Besso) Powder Mix for Instant Shakes"),
				},
			},
		},
	})

	// Domain: art-crafts-antiques
	list = append(list, DomainCatalogDef{
		RootSlug: "art-crafts-antiques",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "traditional-paintings-wooden-crafts",
				NameEN: "Parchment Art, Paintings and Wood Carvings",
				NameAM: "የብራና ስዕሎች፣ የሸራ ቅቦች እና የእንጨት ቅርጻ ቅርጾች",
				NameOM: "Fakkii Birannaa, Mukaa fi Aadaa",
				NameTI: "ናይ ብራና ስእልን ናይ ዕንጨይቲ ቕርጽን",
				Items: []LeafItemDef{
					makeLeaf("goat-parchment-ethiopian-icon-paintings", "Traditional Hand-Painted Goat Skin Parchment (Biranna) Religious Icons", "የፍየል ቆዳ ብራና ላይ የተሳሉ ጥንታዊ የሃይማኖት ስዕሎች", "Fakkii Birannaa Gogaa Re'eerraa", "ናይ Traditional Hand-Painted Goat Skin Parchment (Biranna) Religious Icons"),
					makeLeaf("modern-ethiopian-canvas-oil-paintings", "Signed Original Addis Ababa Cityscape & Cultural Canvas Paintings", "በሸራ ላይ በዘይት ቀለም የተሳሉ ዘመናዊ የኢትዮጵያ ስዕሎች", "Fakkii Shaaraarraa Halluu Zayitaa", "ናይ Signed Original Addis Ababa Cityscape & Cultural Canvas Paintings"),
					makeLeaf("hand-carved-wooden-jimma-chairs", "3-Legged Solid Single-Log Hardwood Traditional Jimma Stools", "የጅማ ባለ 3 እግር ባህላዊ የእንጨት ወንበሮች (Jimma Stool)", "Barcuma Mukaa Jimmaa Aadaa", "ናይ 3-Legged Solid Single-Log Hardwood Traditional Jimma Stools"),
					makeLeaf("handcrafted-wooden-headrests-borkota", "Traditional Carved Wooden Pillows/Headrests (Borkota / Karra)", "የእንጨት ባህላዊ የራስ መደገፊያዎች (ቦርኮታ)", "Boraatii Mukaa Aadaa (Borkota)", "ናይ Traditional Carved Wooden Pillows/Headrests (Borkota / Karra)"),
					makeLeaf("ornate-wooden-processional-crosses", "Hand-Carved Intricate Lalibela Style Wooden Hand & Processional Crosses", "የላሊበላ ጥበብ የእንጨት እጅ እና የመስቀል ቅርጾች", "Fannoo Mukaa Aadaa Laalibalaa", "ናይ Hand-Carved Intricate Lalibela Style Wooden Hand & Processional Crosses"),
					makeLeaf("traditional-leather-painted-shields-gashe", "Embossed Buffalo & Ox Leather Traditional Warrior Battle Shields", "የቆዳ ባህላዊ የጦር ጋሻዎች (ጋሻ)", "Gaachana Gogaa Aadaa", "ናይ Embossed Buffalo & Ox Leather Traditional Warrior Battle Shields"),
					makeLeaf("handmade-clay-sculptures-figurines", "Terracotta Natural Clay Cultural Figures & Traditional Dancers", "የሸክላ ባህላዊ ቅርጻ ቅርጾች እና ምስሎች", "Fakkii Suphee Aadaa", "ናይ Terracotta Natural Clay Cultural Figures & Traditional Dancers"),
					makeLeaf("antique-ethiopian-silver-coptic-crosses", "Old Hand-Hammered Antique Gondar & Axumite Silver Neck Crosses", "የጥንት የጎንደር እና የአክሱም የብር መስቀሎች", "Fannoo Meetaa Durii Gondar", "ናይ Old Hand-Hammered Antique Gondar & Axumite Silver Neck Crosses"),
				},
			},
			{
				Slug:   "handwoven-baskets-pottery-antiques",
				NameEN: "Handwoven Baskets, Pottery and Cultural Antiques",
				NameAM: "በእጅ የተሰሩ ቅርጫቶች፣ የሸክላ እቃዎች እና ጥንታዊ ቅርሶች",
				NameOM: "Masobii, Suphee fi Qabeenya Durii",
				NameTI: "ናይ ኢድ ቅርጫታትን ናይ ሸኽላ ኣቑሑትን",
				Items: []LeafItemDef{
					makeLeaf("harari-colorful-woven-wall-plates", "Intricately Dyed Natural Grass & Straw Decorative Harari Wall Baskets", "የሀረሪ ውብ ባለቀለም ግድግዳ ላይ የሚሰቀሉ ቅርጫቶች", "Qodaa Aadaa Hararii Bareedaa", "ናይ Intricately Dyed Natural Grass & Straw Decorative Harari Wall Baskets"),
					makeLeaf("traditional-clay-dist-and-jebena-sets", "Handmade Scented Clay Coffee Pot (Jebena) & Charcoal Stove Sets", "የሸክላ ጀበናዎች እና የከሰል ማፍያ ምድጃዎች ስብስብ", "Jabanaa fi Majetii Suphee Aadaa", "ናይ Handmade Scented Clay Coffee Pot (Jebena) & Charcoal Stove Sets"),
					makeLeaf("vintage-horn-drinking-cups-wancha", "Polished Natural Ox Horn Traditional Wine & Tej Cups (Wancha)", "የቀንድ ዋንጫዎች ለጠጅ እና ለመጠጥ (ዋንጫ)", "Wancaa Qaanqee Aadaa", "ናይ Polished Natural Ox Horn Traditional Wine & Tej Cups (Wancha)"),
					makeLeaf("hand-beaded-cultural-leather-adornments", "Oromo & Afar Colorful Glass Bead Traditional Leather Belts & Neckpieces", "የኦሮሞ እና የአፋር ባህላዊ የዶቃ ጌጣጌጦች በቆዳ ላይ", "Faaya Callee fi Gogaa Aadaa", "ናይ Oromo & Afar Colorful Glass Bead Traditional Leather Belts & Neckpieces"),
					makeLeaf("antique-brass-liturgical-cymbals-bells", "Ancient Church Brass Bells, Incense Burners & Liturgical Artifacts", "የጥንት የነሐስ ቃጭሎች፣ ጸናጽሎች እና ቅርሶች", "Bilbila fi Meeshaalee Sibilaa Durii", "ናይ Ancient Church Brass Bells, Incense Burners & Liturgical Artifacts"),
					makeLeaf("woven-sisal-baskets-storage-lids", "Natural Plant Fiber Coiled Colorful Multi-Purpose Storage Baskets", "የቃጫ እና የሳር ውብ እቃ ማስቀመጫ ቅርጫቶች", "Kookii fi Saanduqa Aadaa", "ናይ Natural Plant Fiber Coiled Colorful Multi-Purpose Storage Baskets"),
					makeLeaf("traditional-spinning-wheels-injera-covers", "Handmade Bamboo Injera Lids (Akelo) & Cotton Spinning Spindles", "የእንጀራ አክባዮች (አኬሎ) እና የእንዝርት እቃዎች", "Akeeloo Injeeraa fi Inzirtii", "ናይ Handmade Bamboo Injera Lids (Akelo) & Cotton Spinning Spindles"),
					makeLeaf("framed-calligraphy-ancient-geez-manuscripts", "Replicated Ancient Geez Illuminated Vellum Script Framing Art", "የጥንታዊ ግዕዝ ጽሁፎች በክፈፍ የተዘጋጁ ውብ የጥበብ ስራዎች", "Barreeffama Ge'ez Durii Faayame", "ናይ Replicated Ancient Geez Illuminated Vellum Script Framing Art"),
				},
			},
		},
	})

	return list
}
