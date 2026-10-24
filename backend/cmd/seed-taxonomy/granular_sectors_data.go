package main

import "fmt"

// GetGranularSectorsData returns rich, real-world marketplace sector hierarchies
func GetGranularSectorsData() []SubcategorySpec {
	var sectors []SubcategorySpec

	// Helper to add sector with items
	addSector := func(domain, slug, en, am, om, ti string, items [][4]string) {
		var leaves []LeafCategorySpec
		for _, it := range items {
			leaves = append(leaves, LeafCategorySpec{
				ParentSubSlug: slug,
				Slug:          it[0],
				NameEN:        it[1],
				NameAM:        it[2],
				NameOM:        it[3],
				NameTI:        fmt.Sprintf("ናይ %s", it[1]),
			})
		}
		sectors = append(sectors, SubcategorySpec{
			DomainSlug: domain,
			Slug:       slug,
			NameEN:     en,
			NameAM:     am,
			NameOM:     om,
			NameTI:     ti,
			Leaves:     leaves,
		})
	}

	// =========================================================================
	// 41. Vehicle Parts, Tires and Accessories (vehicle-parts-accessories)
	// =========================================================================
	addSector("vehicle-parts-accessories", "auto-engine-components", "Internal Engine Parts and Components", "የሞተር መለዋወጫዎች", "Qaamolee Mootora Konkolaataa", "ናይ ሞተር መለዋወጢ", [][4]string{
		{"pistons-piston-rings", "Engine Pistons, Rings and Gudgeon Pins", "የሞተር ፒስተኖች እና ሪንጎች", "Piistoonii fi Riingii Mootoraa"},
		{"cylinder-head-gaskets", "Cylinder Heads and Multi-Layer Steel (MLS) Gaskets", "የሲሊንደር ራስ ጋስኬቶች", "Gaaskeetii fi Kilaandaaraa"},
		{"crankshafts-camshafts", "Forged Steel Crankshafts and Camshafts", "ክራንክ ሻፍት እና ካም ሻፍት", "Kiraankishaawufii fi Kaamshaawufii"},
		{"timing-belt-chain-kits", "Engine Timing Belts, Chains and Tensioner Pulley Kits", "የታይሚንግ ሻንሰሮች እና ቀበቶዎች", "Qabattoo fi Sansalata Taayimiingii"},
		{"engine-valves-lifters", "Intake and Exhaust Engine Valves and Hydraulic Lifters", "የሞተር ቫልቮች", "Vaalvii Mootoraa"},
		{"engine-mounts-bushings", "Hydraulic & Rubber Anti-Vibration Engine Mounts", "የሞተር ጎማዎች (ሞተር ማውንት)", "Gommaa Mootora Qabu (Mount)"},
	})

	addSector("vehicle-parts-accessories", "cooling-lubrication-systems", "Radiators, Water Pumps and Oil Coolers", "የራዲያተር እና የማቀዝቀዣ ሲስተም", "Raadiyeetara fi Qabbaneessituu", "ናይ ራዲያተርን መዝሓሊን", [][4]string{
		{"aluminum-car-radiators", "Heavy Duty Aluminum Engine Cooling Radiators", "የአሉሚኒየም ራዲያተሮች", "Raadiyeetara Aluminiyeemii"},
		{"automotive-water-pumps", "Engine Mechanical & Electric Water Cooling Pumps", "የውሃ ፓምፖች (ዋተር ፓምፕ)", "Poompii Bishaan Mootoraa"},
		{"cooling-radiator-fans", "Electric Radiator Cooling Fans & Viscous Fan Clutches", "የራዲያተር ራዲያል ማቀዝቀዣ ቬንቲሌተሮች", "Fiixee Raadiyeetaraa"},
		{"engine-coolant-antifreeze", "Long-Life Glycol Radiator Coolant & Antifreeze", "የራዲያተር ማቀዝቀዣ ውሃ (አንቲፍሪዝ)", "Bishaan Raadiyeetaraa Qabbaneessu"},
	})

	addSector("vehicle-parts-accessories", "suspension-steering-parts", "Shock Absorbers, Springs and Steering Racks", "ሾክ አምፖርሳሮች እና እገዳዎች", "Shookii fi Isteeriingii", "ናይ ሾክ ኣምፖርሳርን እገዳን", [][4]string{
		{"gas-hydraulic-shock-absorbers", "Front & Rear Gas-Pressurized Shock Absorbers", "ሾክ አምፖርሳሮች (ሾክ)", "Shookiifi Isteeriingii Konkolaataa"},
		{"heavy-coil-leaf-springs", "Suspension Coil Springs and Multi-Leaf Springs", "የተሽከርካሪ ባሌስትራዎች (ሊፍ ስፕሪንግ)", "Baaleestiraa fi Isphiiriingii"},
		{"power-steering-racks-pumps", "Hydraulic Power Steering Racks and Pressure Pumps", "የሃይድሮሊክ መዘወሪያ (ስቲሪንግ ራክ)", "Isteeriingii fi Poompii Harkaa"},
		{"ball-joints-tie-rod-ends", "Front Suspension Ball Joints, Tie Rod Ends & Bushings", "የፊት እገዳ ቦል ጆይንቶች እና ቴርሚናሎች", "Boojjeyintii fi Terminaala"},
	})

	addSector("vehicle-parts-accessories", "auto-electrical-ignition", "Alternators, Starter Motors and Sensors", "የመኪና ኤሌክትሪክ እና ስታርተሮች", "Elektiriikii fi Istaatara Konkolaataa", "ናይ መኪና ኤሌክትሪክን ስታርተርን", [][4]string{
		{"12v-24v-alternators", "High-Output 12V and 24V Alternators and Dynamos", "የመኪና ዳይናሞዎች እና አልተርኔተሮች", "Daayinaamoo fi Altarneetara"},
		{"automotive-starter-motors", "Heavy Duty Engine Starter Motors and Bendix Gears", "የመኪና ማስነሻ ስታርተሮች", "Istaatara Konkolaataa"},
		{"iridium-spark-glow-plugs", "Laser Iridium Spark Plugs & Diesel Glow Plugs", "የሻማ እና የማሞቂያ ፕላጎች (ካንዴላ)", "Kaandelaa fi Kaandiilee Diizilii"},
		{"ecu-electronic-sensors", "Oxygen Sensors, Mass Airflow (MAF) & Crankshaft Sensors", "የመኪና ሴንሰሮች (ኦክስጅን፣ MAF)", "Seensara fi Kompiitara ECU"},
	})

	addSector("vehicle-parts-accessories", "vehicle-body-lighting-accessories", "Headlights, Mirrors, Bumpers and Body Parts", "የፊት መብራቶች፣ መስተዋቶች እና ባምፐሮች", "Ibsaa Fuulduraa fi Daawwitii", "ናይ ቅድሚት መብራህትን መስትያትን", [][4]string{
		{"led-halogen-headlight-assemblies", "OEM LED and Halogen Front Headlight Assemblies", "የፊት ዋና መብራቶች (ሄድላይት)", "Ibsaa Fuulduraa Gurguddaa"},
		{"tail-lights-brake-lamps", "Rear Tail Lights, Reverse Lamps & Indicators", "የኋላ ፍሬን እና ሪቨርስ መብራቶች", "Ibsaa Duubaa fi Firiinii"},
		{"exterior-side-view-mirrors", "Power-Folding Heated Side View Mirrors", "የጎን መመልከቻ መስተዋቶች", "Daawwitii Cinaa Konkolaataa"},
		{"front-rear-bumper-bars", "ABS Plastic and Steel Bull Bars & Bumpers", "የፊት እና የኋላ ባምፐሮች (መከላከያዎች)", "Baampara fi Eegduu Fuulduraa"},
		{"custom-floor-mats-seat-covers", "Leatherette Custom Seat Covers & All-Weather Floor Mats", "የመቀመጫ ልብሶች እና የወለል ምንጣፎች", "Uffata Teessumaa fi Afata Lafaa"},
	})

	// =========================================================================
	// 44. Crops, Grains, Seeds and Plants (crops-grains-seeds)
	// =========================================================================
	addSector("crops-grains-seeds", "pulses-legumes-beans", "Lentils, Field Peas, Faba Beans and Pulses", "ምስር፣ አተር፣ ባቄላ እና ጥራጥሬዎች", "Misira, Atara, Baqelaa fi Booloqqee", "ምስሪ፣ ዓተርን ባልደንጓን", [][4]string{
		{"ethiopian-red-lentils-misir", "Grade 1 Cleaned Ethiopian Red Lentils (Misir)", "የተፈተገ እና ድፍን ቀይ ምስር", "Misira Diimaa Qulqulluu"},
		{"green-split-field-peas-ater", "Green & Yellow Split Field Peas (Ater)", "አረንጓዴ እና ቢጫ ክክ አተር", "Atara Magariisa fi Keelloo"},
		{"large-faba-broad-beans-baqela", "Large Grade Faba Broad Beans (Baqela)", "ትላልቅ የሀገር ቤት ባቄላ", "Baqelaa Gurguddaa"},
		{"white-pea-navy-haricot-beans", "White Pea Navy Beans & Canning Haricot Beans", "ነጭ ቦሎቄ (የኤክስፖርት ጥራት)", "Booloqqee Adii"},
		{"soya-beans-feed-grade", "Non-GMO Feed Grade & Processing Soybeans", "አኩሪ አተር በኩንታል", "Akkurii Ataraa"},
	})

	addSector("crops-grains-seeds", "oilseeds-cash-crops", "Sesame, Linseed, Groundnuts and Oilseeds", "ሰሊጥ፣ ተልባ፣ ለውዝ እና የቅባት እህሎች", "Saliixa, Talbaa fi Loozii", "ሰሊጥ፣ ተልባን ለውዝን", [][4]string{
		{"whitish-humera-sesame-export", "Whitish Humera Grade 1 Export Sesame", "የሁመራ አንደኛ ደረጃ ነጭ ሰሊጥ", "Saliixa Adii Humaratti"},
		{"wollega-mixed-sesame", "Wollega Type High-Oil Mixed Sesame Seeds", "የወለጋ ሰሊጥ", "Saliixa Wallaggaa"},
		{"brown-linseed-telba", "High-Purity Brown Linseed / Flaxseed (Telba)", "ንጹህ የሀገር ቤት ተልባ", "Talbaa Qulqulluu"},
		{"shelled-groundnuts-peanuts", "Dry Shelled Groundnuts / Peanuts (Babile)", "የባቢሌ የተላጠ ለውዝ", "Loozii Baabilee"},
		{"sunflower-seeds-commercial", "Black & Striped Oil-Rich Sunflower Seeds", "የሱፍ ፍሬ በኩንታል", "Sanyii Suufii Zayitaa"},
	})

	addSector("crops-grains-seeds", "fresh-vegetables-horticulture", "Fresh Onions, Tomatoes, Potatoes and Vegetables", "ሽንኩርት፣ ቲማቲም፣ ድንች እና አትክልቶች", "Qullubbii, Timaatima fi Dinnicha", "ሽጉርቲ፣ ቲማቲምን ድንችን", [][4]string{
		{"red-bombay-onions", "Red Bombay Dry Bulb Onions (Ziway / Meki)", "የዝዋይ እና የመቂ ቀይ ሽንኩርት", "Qullubbii Diimaa Zeway"},
		{"processing-tomatoes-roma", "Fresh Farm Harvested Roma & Salad Tomatoes", "ሮማ እና የሰላጣ ቲማቲም", "Timaatima Oomisha Haaraa"},
		{"table-seed-potatoes", "White & Red Table Potatoes and Seed Tubers", "የምግብ እና የዘር ድንች", "Dinnicha Nyaataa fi Sanyii"},
		{"green-chili-peppers-karia", "Fresh Green Hot Chili Peppers (Karia)", "ትኩስ የቃሪያ ምርት", "Qaraa fi Barbaree Magariisa"},
		{"white-cabbage-green-kale", "Fresh White Cabbage & Gomen Greens (Tikle Gomen)", "ጥቅል ጎመን እና የሀበሻ ጎመን", "Raafuu fi Tikle Goomana"},
		{"table-carrots-beetroots", "Fresh Washed Orange Carrots and Red Beetroots", "ካሮት እና ቀይ ስር", "Kaarotii fi Qayisaara"},
		{"fresh-garlic-bulbs", "Debre Berhan & Gondar Fresh Dry White Garlic", "የደብረ ብርሃን እና የጎንደር ነጭ ሽንኩርት", "Qullubbii Adii"},
	})

	addSector("crops-grains-seeds", "fresh-fruits-orchards", "Avocados, Mangoes, Bananas and Citrus", "አቮካዶ፣ ማንጎ፣ ሙዝ እና ብርቱካን", "Avokaadoo, Maangoo fi Muuzii", "ኣቮካዶ፣ ማንጎን ባናናን", [][4]string{
		{"hass-avocados-export", "Hass & Fuerte Export Quality Green Avocados", "ሀስ እና ፉዌርቴ አቮካዶ", "Avokaadoo Haasi fi Fuyeertee"},
		{"apple-mangoes-dilla-arba-minch", "Sweet Apple Mangoes & Kent Mangoes (Arba Minch)", "የአርባ ምንጭ እና የዲላ ማንጎ", "Maangoo Mi'aawaa Arba Minch"},
		{"cavendish-bananas-crates", "Sweet Cavendish Bananas (By Crate & Truck Load)", "ጣፋጭ የሀበሻ ሙዝ በሳጥን", "Muuzii Kaawendishi"},
		{"sweet-valencia-oranges", "Juicy Sweet Valencia Oranges & Mandarins", "ብርቱካን እና መንደሪን", "Burtukaana fi Mandariinii"},
		{"fresh-papayas-watermelons", "Red Lady Papayas and Large Sweet Watermelons", "ፓፓያ እና ሀብሀብ", "Paappaayyaa fi Qamaha"},
	})

	// =========================================================================
	// 29. Construction Materials & Building Supplies (construction-materials)
	// =========================================================================
	addSector("construction-materials", "concrete-blocks-masonry", "Hollow Concrete Blocks, Bricks and Stone", "የሲሚንቶ ብሎኬቶች፣ ጡቦች እና ድንጋዮች", "Bilookeetii Konkiriitii fi Xoophii", "ብሎኬት፣ ጡብን እምንን", [][4]string{
		{"hcb-machine-cut-blocks-15-20", "Class A Machine-Vibrated Hollow Concrete Blocks (15cm & 20cm)", "የ15 እና የ20 ማሽን ቆራጭ ብሎኬት", "Bilookeetii Konkiriitii 15cm fi 20cm"},
		{"kiln-fired-clay-red-bricks", "Kiln-Fired Solid & Perforated Red Clay Bricks", "የተጋገሩ ቀይ የሸክላ ጡቦች", "Xoophii Diimaa Suphee"},
		{"chiseled-basalt-foundation-stone", "Hand-Chiseled Basalt Building Stones (Per Truck Load)", "የተፈለፈለ የመሰረት ድንጋይ በጭነት", "Dhagaa Hundeessaa Qotame"},
		{"paving-interlocking-blocks", "Heavy Duty Interlocking Concrete Paver Blocks (6cm & 8cm)", "የአስፋልት እና የግቢ ንጣፍ ኢንተርሎኪንግ ድንጋይ", "Kobbilii fi Dhagaa Lafaa"},
	})

	addSector("construction-materials", "doors-windows-glass", "Steel Security Doors, Aluminum Profiles and Glass", "የደህንነት በሮች፣ የአሉሚኒየም መስኮቶች እና መስታወት", "Balbala Sibiilaa fi Foddaa Aluminiyeemii", "ናይ ሓጺን ማዕጾን መስኮትን", [][4]string{
		{"turkish-steel-security-doors", "Imported Turkish Multi-Lock Armored Steel Security Doors", "የቱርክ ባለ ብዙ ቁልፍ የደህንነት የብረት በሮች", "Balbala Sibiilaa Tirkii Nageenyaa"},
		{"mdf-flush-wooden-interior-doors", "Pre-Hung MDF and Teak Wood Interior Bedroom Doors", "የክፍል የእንጨት በሮች (MDF)", "Balbala Mukaa Kutaa Jireenyaa"},
		{"aluminum-sliding-window-frames", "Powder-Coated Aluminum Sliding Window Profiles with Glass", "የአሉሚኒየም ተንሸራታች የመስኮት ፍሬሞች", "Foddaa Aluminiyeemii Dachaafamu"},
		{"tinted-tempered-building-glass", "6mm, 8mm & 10mm Clear, Tinted & Tempered Architectural Glass", "የህንጻ እና የመስኮት የጠነከረ መስታወት", "Daawwitii Foddaa Jabaa (Tempered)"},
	})

	addSector("construction-materials", "timber-roofing-adhesives", "Wood Timber, Tile Adhesives and Waterproofing", "እንጨት፣ የሴራሚክ ማጣበቂያ እና ዋተርፕሩፊንግ", "Muka Ijaarsaa fi Zayita Bishaan Dhorku", "ዕንጨይቲ፣ ናይ ሴራሚክ መጣበቕን ዋተርፕሩፍን", [][4]string{
		{"eucalyptus-zigba-construction-timber", "Treated Eucalyptus Poles, Zigba Timber & Formwork Planks", "የባህር ዛፍ ቋሚዎች እና የዝግባ ሳንቃዎች", "Muka Ijaarsaa Baahrzaafii fi Ziqbaa"},
		{"polymer-modified-tile-adhesives", "C2TE Polymer-Modified Flexible Ceramic Tile Adhesive (25kg)", "የሴራሚክ ማጣበቂያ ዱቄት (25 ኪ.ግ)", "Qalama Seeraamikii Qabsiisu 25kg"},
		{"waterproof-cementitious-slurry", "2-Component Cementitious Waterproofing Slurry (SikaTop)", "የመሰረት እና የጣሪያ የውሃ መከላከያ (ሲካ)", "Qoricha Bishaan Dhorku (Sika)"},
		{"bitumen-torch-on-roofing-membranes", "4mm Bituminous Torch-on Waterproofing Roofing Membranes", "የጣሪያ ሬንጅ ፕላስቲክ (ሮል)", "Membraanii Bitumeenii Baaxii"},
	})

	// =========================================================================
	// 10. Clothing, Fashion & Apparel (clothing-apparel)
	// =========================================================================
	addSector("clothing-apparel", "childrens-infants-clothing", "Kids, Toddlers and Children's Fashion", "የልጆች እና የህፃናት ልብሶች", "Uffata Daa'immanii fi Ijoollee", "ናይ ቆልዑን ህጻናትን ክዳን", [][4]string{
		{"baby-cotton-rompers-bodysuits", "100% Cotton Newborn Rompers and Onesies", "የጨቅላ ህፃናት የጥጥ ልብሶች", "Uffata Daa'imman Dhalatan Qulqulluu"},
		{"boys-casual-suits-sets", "Boys Casual T-Shirts, Shorts and Formal Suit Sets", "የወንዶች ልጆች ልብሶች እና ሱፎች", "Uffata Ijoollee Dhiiraa"},
		{"girls-party-dresses-skirts", "Girls Party Princess Dresses and Casual Skirts", "የሴቶች ልጆች ቀሚሶች", "Qamisii Ijoollee Dubaraa"},
		{"school-uniform-sets", "Standard Primary & Secondary School Uniform Sets", "የትምህርት ቤት ዩኒፎርሞች", "Yuunifoormii Mana Barumsaa"},
	})

	addSector("clothing-apparel", "workwear-safety-garments", "Overalls, Uniforms and Safety Jackets", "የስራ ልብሶች፣ ኦቨሮሎች እና አንጸባራቂ ጃኬቶች", "Uffata Hojii fi Eegumsaa", "ናይ ስራሕ ክዳን (ኦቨሮል)", [][4]string{
		{"industrial-cotton-overalls-boilersuits", "100% Cotton Heavy Duty Industrial Overalls (Boilersuits)", "የፋብሪካ እና የጋራዥ ኦቨሮል ልብሶች", "Oovaroolii Hojii Sibiilaafi Warshaa"},
		{"high-visibility-reflective-vests", "High-Visibility Fluorescent Reflective Safety Vests", "አንጸባራቂ የሴፍቲ ጃኬቶች (ቬስት)", "Jaakeetii Eegumsaa Ifaa (Reflective)"},
		{"medical-scrubs-doctor-coats", "Medical Doctors Lab Coats & Clinical Scrub Sets", "የሃኪም ነጭ ጋውኖች እና የነርሶች ልብሶች", "Gawonii Adii Hakiimotaa"},
		{"security-guard-uniform-sets", "Security Guard Uniform Sets with Badges & Caps", "የጥበቃ ሰራተኞች ዩኒፎርሞች", "Yuunifoormii Eegdota Nageenyaa"},
	})

	// =========================================================================
	// 15. Sports, Fitness and Outdoor Recreation (sports-fitness)
	// =========================================================================
	addSector("sports-fitness", "outdoor-camping-hiking", "Camping Tents, Sleeping Bags and Hiking Gear", "የካምፒንግ ድንኳኖች እና የሽርሽር እቃዎች", "Dunkaana Kaampiingii fi Imalaa", "ናይ ካምፒንግ ድንኳንን መገሻን", [][4]string{
		{"waterproof-camping-dome-tents", "2 to 6 Person Waterproof Windproof Camping Dome Tents", "የሽርሽር እና የካምፒንግ ድንኳኖች", "Dunkaana Kaampiingii Bishaan Danda'u"},
		{"thermal-hiking-sleeping-bags", "Sub-Zero Thermal Hiking Sleeping Bags", "የካምፒንግ መኝታ ቦርሳዎች (Sleeping Bags)", "Saanduqa Ciisichaa Kaampiingii"},
		{"trekking-poles-hiking-backpacks", "Aluminum Hiking Trekking Poles & 60L Internal Frame Backpacks", "የእግር ጉዞ ምርኩዞች እና የጉዞ ቦርሳዎች", "Ulee Imalaa fi Boorsaa Dugdaa 60L"},
		{"portable-outdoor-gas-stoves", "Mini Portable Camping Gas Stoves & Cookware Sets", "ተንቀሳቃሽ የሽርሽር ጋዞች", "Istoovii Kaampiingii Salphaa"},
	})

	// =========================================================================
	// 18. Jobs, Careers and Recruitment (jobs-recruitment)
	// =========================================================================
	addSector("jobs-recruitment", "hospitality-catering-jobs", "Hotel, Restaurant, Chef and Hospitality Roles", "የሆቴል፣ የቱሪዝም እና የሬስቶራንት ስራዎች", "Hojiiwwan Hoteelaa fi Tajaajilaa", "ናይ ሆቴልን ቱሪዝምን ስራሕቲ", [][4]string{
		{"executive-head-chefs-cooks", "Head Chefs, Pastry Chefs and Line Cooks", "ዋና የምግብ አብሳዮች (ሼፎች)", "Banaa Hojii Sheefii Nyaataa"},
		{"hotel-front-desk-receptionists", "Hotel Front Desk Receptionists and Guest Services", "የሆቴል የፊት ለፊት አስተናጋጆች (ሪሴፕሽን)", "Banaa Hojii Risaapshinii Hoteelaa"},
		{"hotel-housekeeping-supervisors", "Housekeeping Room Attendants and Supervisors", "የጽዳት እና የቤት ሰራተኞች ተቆጣጣሪዎች", "To'attoota Qulqullina Hoteelaa"},
		{"food-beverage-service-waiters", "Restaurant Waiters, Waitresses and Bartenders", "የምግብ እና የመጠጥ አስተናጋጆች", "Keeshummeessitoota Nyaataa fi Dhugaatii"},
	})

	addSector("jobs-recruitment", "transport-logistics-driver-jobs", "Commercial Drivers, Dispatchers and Logistics", "የአሽከርካሪዎች እና የትራንስፖርት ስራዎች", "Hojii Konkolaachisummaa fi Loojistiiksii", "ናይ መራሕቲ መካይን ስራሕቲ", [][4]string{
		{"heavy-truck-trailer-driver-jobs", "Dry Cargo & Fuel Tanker Grade 4/5 Truck Drivers", "የደረቅ ጭነት እና የቦቴ አሽከርካሪዎች", "Konkolaachiftoota Fe'iinsa Gurguddaa"},
		{"private-chauffeur-driver-jobs", "Private Family & Corporate Automobile Chauffeurs", "የግል እና የድርጅት ሹፌሮች", "Konkolaachiftoota Dhuunfaa fi Dhaabbataa"},
		{"logistics-import-export-officers", "Customs Clearance, Transit & Import/Export Logistics Specialists", "የጉምሩክ ክሊራንስ እና የትራንዚት ባለሙያዎች", "Ogeessota Loojistiiksii fi Gumurukii"},
	})

	// =========================================================================
	// 56. Energy, Solar & Electrical Power Equipment (energy-solar-power)
	// =========================================================================
	addSector("energy-solar-power", "solar-lighting-pico", "Solar Home Lighting Kits and Streetlights", "የሶላር የመንገድ መብራቶች እና የቤት ኪቶች", "Ibsaa Soolaarii Daandiifi Manaa", "ናይ ሶላር መብራህቲታት", [][4]string{
		{"all-in-one-solar-street-lights", "Integrated 100W-300W LED All-in-One Solar Streetlights", "የመንገድ ላይ የሶላር መብራቶች (All-in-One)", "Ibsaa Daandii Soolaarii LED"},
		{"rural-solar-home-lighting-kits", "Pico Solar Home Lighting Kits with Bulbs & Phone Charger", "የገጠር የቤት ውስጥ የሶላር መብራት ኪቶች", "Meeshaa Ibsaa Soolaarii Baadiyyaa"},
		{"solar-floodlights-remote", "Commercial Remote-Controlled Solar LED Floodlights", "የግቢ ማብሪያ የሶላር ፕሮጀክተሮች", "Piroojeektara Soolaarii Mooraa"},
	})

	return sectors
}
