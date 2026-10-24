package main

import (
	"fmt"
)

// DomainSpec represents a domain with its subcategories and leaf items
type DomainSpec struct {
	DomainSlug string
	Branches   []BranchSpec
}

// BranchSpec represents a subcategory group containing leaf items
type BranchSpec struct {
	SubSlug string
	NameEN  string
	NameAM  string
	NameOM  string
	NameTI  string
	Items   []ItemSpec
}

// ItemSpec represents a terminal leaf category
type ItemSpec struct {
	Slug   string
	NameEN string
	NameAM string
	NameOM string
	NameTI string
}

// GetAllDomainSpecs returns the full commercial catalog branches for all 60 roots
func GetAllDomainSpecs() []DomainSpec {
	return []DomainSpec{
		// 04. Cameras, Photography and Optics
		{
			DomainSlug: "cameras-optics",
			Branches: []BranchSpec{
				{
					SubSlug: "dslr-mirrorless-cameras",
					NameEN:  "Digital Cameras (DSLR & Mirrorless)",
					NameAM:  "ዲጂታል ካሜራዎች",
					NameOM:  "Kaameeraa Dijitaalaa",
					NameTI:  "ዲጂታል ካሜራታት",
					Items: []ItemSpec{
						{Slug: "full-frame-mirrorless", NameEN: "Full-Frame Professional Mirrorless Cameras", NameAM: "ፕሮፌሽናል ሚረርለስ ካሜራዎች", NameOM: "Kaameeraa Miirarlesii", NameTI: "ፕሮፌሽናል ሚረርለስ ካሜራ"},
						{Slug: "aps-c-dslr-cameras", NameEN: "APS-C & Crop Sensor DSLR Cameras", NameAM: "ዲኤስኤልአር ካሜራዎች", NameOM: "Kaameeraa DSLR", NameTI: "DSLR ካሜራ"},
						{Slug: "compact-point-shoot", NameEN: "Compact Point and Shoot Cameras", NameAM: "ቀላል የኪስ ካሜራዎች", NameOM: "Kaameeraa Qaphxii fi Waraabii", NameTI: "ቀሊል ናይ ጁባ ካሜራ"},
						{Slug: "action-cameras-4k", NameEN: "4K Waterproof Action Cameras", NameAM: "የውሃ መቋቋሚ አክሽን ካሜራዎች", NameOM: "Kaameeraa Bishaan Danda'u", NameTI: "ናይ ማይ መኸላኸሊ ኣክሽን ካሜራ"},
					},
				},
				{
					SubSlug: "camera-lenses-optics",
					NameEN:  "Camera Lenses and Optical Filters",
					NameAM:  "የካሜራ ሌንሶች",
					NameOM:  "Leensii Kaameeraa",
					NameTI:  "ናይ ካሜራ ሌንሳት",
					Items: []ItemSpec{
						{Slug: "prime-portrait-lenses", NameEN: "Prime Portrait Lenses (50mm, 85mm)", NameAM: "የፖርትሬት ሌንሶች (50ሚሜ/85ሚሜ)", NameOM: "Leensii Suuraa Fuulaa", NameTI: "ናይ ፖርትሬት ሌንስ"},
						{Slug: "telephoto-zoom-lenses", NameEN: "Telephoto & Wildlife Zoom Lenses (70-200mm)", NameAM: "የቴሌፎቶ ማጉያ ሌንሶች", NameOM: "Leensii Fageenyaa", NameTI: "ናይ ርሑቕ መጉልሒ ሌንስ"},
						{Slug: "wide-angle-landscape-lenses", NameEN: "Ultra Wide-Angle Landscape Lenses", NameAM: "ሰፊህ እይታ ያላቸው ሌንሶች", NameOM: "Leensii Bal'aa", NameTI: "ሰፊሕ ምርኢት ዘለዎ ሌንስ"},
						{Slug: "macro-closeup-lenses", NameEN: "Macro & Close-Up Specialized Lenses", NameAM: "የማክሮ ሌንሶች", NameOM: "Leensii Maakroo", NameTI: "ማክሮ ሌንስ"},
					},
				},
				{
					SubSlug: "drones-aerial-imaging",
					NameEN:  "Camera Drones and Aerial Equipment",
					NameAM:  "ድሮኖች እና የአየር ላይ ካሜራዎች",
					NameOM:  "Diroonii fi Suuraa Qilleensaa",
					NameTI:  "ድሮናትን ናይ ኣየር ላይ ካሜራን",
					Items: []ItemSpec{
						{Slug: "professional-4k-drones", NameEN: "Professional 4K Video Drones (DJI Mavic)", NameAM: "ፕሮፌሽናል የቪዲዮ ድሮኖች", NameOM: "Diroonii Viidiyoo 4K", NameTI: "ናይ ቪዲዮ ድሮን 4K"},
						{Slug: "mini-travel-drones", NameEN: "Lightweight Mini Travel Drones (<249g)", NameAM: "ቀላል የጉዞ ድሮኖች", NameOM: "Diroonii Imalaa Salphaa", NameTI: "ቀሊል ናይ መገሻ ድሮን"},
						{Slug: "drone-batteries-propellers", NameEN: "Drone Intelligent Flight Batteries & Propellers", NameAM: "የድሮን ባትሪዎች እና ፕሮፔለሮች", NameOM: "Baatirii fi Meeshaa Diroonii", NameTI: "ናይ ድሮን ባትሪን ፕሮፔለርን"},
					},
				},
				{
					SubSlug: "studio-lighting-support",
					NameEN:  "Studio Lighting, Gimbals and Tripods",
					NameAM:  "የስቱዲዮ መብራቶች እና ትራይፖዶች",
					NameOM:  "Ibsaa Istudiyoo fi Tiraayipoodii",
					NameTI:  "ናይ ስቱዲዮ መብራህትን ትራይፖድን",
					Items: []ItemSpec{
						{Slug: "handheld-3axis-gimbals", NameEN: "3-Axis Handheld Gimbals & Stabilizers", NameAM: "የቪዲዮ ማረጋጊያ ጊምባሎች", NameOM: "Giimbaalii Tasgabbeessaa", NameTI: "ናይ ቪዲዮ መረጋግዒ ጊምባል"},
						{Slug: "heavy-duty-tripods-monopods", NameEN: "Professional Carbon Fiber Tripods & Fluid Heads", NameAM: "የካሜራ ትራይፖዶች", NameOM: "Tiraayipoodii Kaameeraa", NameTI: "ናይ ካሜራ ትራይፖድ"},
						{Slug: "studio-softbox-lights", NameEN: "Continuous LED Studio Softbox Kits", NameAM: "የስቱዲዮ ሶፍት ቦክስ መብራቶች", NameOM: "Ibsaa Istudiyoo Softbox", NameTI: "ናይ ስቱዲዮ ሶፍት ቦክስ መብራህቲ"},
						{Slug: "wireless-lapel-microphones", NameEN: "Wireless Lavalier & Shotgun Audio Microphones", NameAM: "የካሜራ ማይክሮፎኖች", NameOM: "Maayikiroofoonii Viidiyoo", NameTI: "ናይ ካሜራ ማይክሮፎን"},
					},
				},
			},
		},

		// 05. Gaming, Consoles and Virtual Reality
		{
			DomainSlug: "gaming-vr",
			Branches: []BranchSpec{
				{
					SubSlug: "gaming-consoles",
					NameEN:  "Video Game Consoles",
					NameAM:  "የቪዲዮ ጌም ኮንሶሎች",
					NameOM:  "Konsolii Geemiingii",
					NameTI:  "ናይ ጌም ኮንሶላት",
					Items: []ItemSpec{
						{Slug: "playstation-5-consoles", NameEN: "Sony PlayStation 5 (Disc / Digital)", NameAM: "ፕሌይስቴሽን 5 ኮንሶል", NameOM: "PlayStation 5", NameTI: "ፕሌይስቴሽን 5"},
						{Slug: "playstation-4-pro-slim", NameEN: "Sony PlayStation 4 (Pro / Slim)", NameAM: "ፕሌይስቴሽን 4", NameOM: "PlayStation 4", NameTI: "ፕሌይስቴሽን 4"},
						{Slug: "xbox-series-x-s", NameEN: "Microsoft Xbox Series X & Series S", NameAM: "ኤክስ ቦክስ ኮንሶሎች", NameOM: "Xbox Series X/S", NameTI: "ኤክስ ቦክስ"},
						{Slug: "nintendo-switch-oled", NameEN: "Nintendo Switch (OLED / Lite)", NameAM: "ኒንቴንዶ ስዊች", NameOM: "Nintendo Switch", NameTI: "ኒንቴንዶ ስዊች"},
					},
				},
				{
					SubSlug: "gaming-accessories",
					NameEN:  "Controllers, Headsets & Gaming Chairs",
					NameAM:  "የጌሚንግ መቆጣጠሪያዎች እና ወንበሮች",
					NameOM:  "Meeshaalee Geemiingii fi Barcuma",
					NameTI:  "ናይ ጌሚንግ መቆጻጸርን መንበርን",
					Items: []ItemSpec{
						{Slug: "wireless-game-controllers", NameEN: "DualSense & Xbox Wireless Controllers", NameAM: "ገመድ አልባ ጆይስቲክ / መቆጣጠሪያ", NameOM: "To'ataa Geemiingii (Joystiikii)", NameTI: "ናይ ጌም መቆጻጸሪ (ጆይስቲክ)"},
						{Slug: "surround-gaming-headsets", NameEN: "7.1 Surround Noise-Canceling Gaming Headsets", NameAM: "የጌሚንግ ሄድፎኖች", NameOM: "Heedifoonii Geemiingii", NameTI: "ናይ ጌሚንግ ሄድፎን"},
						{Slug: "ergonomic-gaming-chairs", NameEN: "Ergonomic High-Back Leather Gaming Chairs", NameAM: "የጌሚንግ ወንበሮች", NameOM: "Barcuma Geemiingii", NameTI: "ናይ ጌሚንግ መንበር"},
						{Slug: "racing-wheels-pedals", NameEN: "Force Feedback Steering Wheels & Pedals", NameAM: "የመኪና መዘወሪያ እና ፔዳል መቆጣጠሪያዎች", NameOM: "Meeshaa Konkolaachisaa Geemii", NameTI: "ናይ መኪና መዘወሪ ጌም"},
					},
				},
			},
		},

		// 08. Furniture and Interior Furnishings
		{
			DomainSlug: "furniture-furnishings",
			Branches: []BranchSpec{
				{
					SubSlug: "living-room-furniture",
					NameEN:  "Living Room Sofas, Couches and Coffee Tables",
					NameAM:  "የሳሎን ሶፋዎች እና ጠረጴዛዎች",
					NameOM:  "Soofaa fi Meeshaa Saaloonaa",
					NameTI:  "ናይ ሳሎን ሶፋታትን ጠረጴዛታትን",
					Items: []ItemSpec{
						{Slug: "l-shaped-corner-sofas", NameEN: "L-Shaped Sectional Corner Sofas", NameAM: "ኤል-ሼፕ የማእዘን ሶፋዎች", NameOM: "Soofaa Bocaa L", NameTI: "L ቅርጺ ዘለዎ ሶፋ"},
						{Slug: "genuine-leather-sofa-sets", NameEN: "Ethiopian Genuine Leather Sofa Sets (3+2+1)", NameAM: "የቆዳ ሶፋዎች", NameOM: "Soofaa Gogaa Qulqulluu", NameTI: "ናይ ቆርበት ሶፋታት"},
						{Slug: "modern-tv-stands-consoles", NameEN: "Modern Wall-Mounted TV Stands & Media Consoles", NameAM: "የቴሌቪዥን ማስቀመጫ ኮንሶሎች", NameOM: "Minjaala Televiziyinii", NameTI: "ናይ ቴሌቪዥን መቐመጢ"},
						{Slug: "wooden-glass-coffee-tables", NameEN: "Tempered Glass & Solid Wood Coffee Tables", NameAM: "የቡና እና የሳሎን ጠረጴዛዎች", NameOM: "Minjaala Bunaa", NameTI: "ናይ ቡን ጠረጴዛ"},
					},
				},
				{
					SubSlug: "bedroom-furniture",
					NameEN:  "Beds, Mattresses and Wardrobes",
					NameAM:  "አልጋዎች፣ ፍራሾች እና ቁም ሳጥኖች",
					NameOM:  "Siree, Firaashii fi Saanduqa Uffataa",
					NameTI:  "ዓራት፣ ፍራሽን ናይ ክዳን ሳንዱቕን",
					Items: []ItemSpec{
						{Slug: "king-queen-size-beds", NameEN: "King and Queen Size Hardwood Beds", NameAM: "ኪንግ እና ኩዊን ሳይዝ አልጋዎች", NameOM: "Siree Gurguddaa", NameTI: "ዓበይቲ ዓራታት"},
						{Slug: "orthopedic-foam-mattresses", NameEN: "High-Density Orthopedic Foam Mattresses", NameAM: "የኦርቶፔዲክ ስፖንጅ ፍራሾች", NameOM: "Firaashii Ispoonjii Jabaa", NameTI: "ኦርቶፔዲክ ናይ ሰፍነግ ፍራሽ"},
						{Slug: "multi-door-wooden-wardrobes", NameEN: "3-Door & 4-Door Wooden Clothes Wardrobes", NameAM: "ባለ 3 እና 4 በር የልብስ ቁምሳጥኖች", NameOM: "Saanduqa Uffataa Mukaa", NameTI: "ናይ ክዳን ቁምሳንዱቕ"},
						{Slug: "dressing-tables-mirrors", NameEN: "Bedroom Dressing Tables with Vanity Mirrors", NameAM: "የመዋቢያ ጠረጴዛዎች ከመስታወት ጋር", NameOM: "Minjaala Fuulaa fi Daawwitii", NameTI: "ናይ መመላኸዒ ጠረጴዛ ምስ መስትያት"},
					},
				},
				{
					SubSlug: "office-commercial-furniture",
					NameEN:  "Office Desks, Ergonomic Chairs and Filing",
					NameAM:  "የቢሮ ጠረጴዛዎች እና ወንበሮች",
					NameOM:  "Minjaala Biiroo fi Barcuma",
					NameTI:  "ናይ ቢሮ ጠረጴዛን መንበርን",
					Items: []ItemSpec{
						{Slug: "executive-wooden-desks", NameEN: "Managerial & Executive L-Shaped Desks", NameAM: "የስራ አስኪያጅ ጠረጴዛዎች", NameOM: "Minjaala Hojii Hooggana", NameTI: "ናይ ኣመሓዳሪ ጠረጴዛ"},
						{Slug: "mesh-ergonomic-office-chairs", NameEN: "High-Back Mesh Ergonomic Swivel Chairs", NameAM: "የሚሽከረከሩ ምቹ የቢሮ ወንበሮች", NameOM: "Barcuma Naanna'u Biiroo", NameTI: "ዝዘውር ናይ ቢሮ መንበር"},
						{Slug: "conference-meeting-tables", NameEN: "8 to 16 Person Conference Room Tables", NameAM: "የስብሰባ አዳራሽ ጠረጴዛዎች", NameOM: "Minjaala Walga'ii", NameTI: "ናይ ኣኼባ ጠረጴዛ"},
						{Slug: "metal-fireproof-filing-cabinets", NameEN: "Steel Filing Cabinets and Storage Lockers", NameAM: "የብረት ሰነድ ማስቀመጫ ካቢኔቶች", NameOM: "Kaabineetii Sanadaa Sibiilaa", NameTI: "ናይ ሓጺን ሰነድ መቐመጢ"},
					},
				},
			},
		},

		// 28. Tools, Workshop and Repair Equipment
		{
			DomainSlug: "tools-workshop",
			Branches: []BranchSpec{
				{
					SubSlug: "power-tools",
					NameEN:  "Power Drills, Grinders and Saws",
					NameAM:  "የኤሌክትሪክ ድራይሎች እና መቁረጫዎች",
					NameOM:  "Diriilii fi Meeshaalee Elektiriikii",
					NameTI:  "ናይ ኤሌክትሪክ መሳርሒታትን መቑረጽትን",
					Items: []ItemSpec{
						{Slug: "rotary-hammer-drills", NameEN: "Heavy Duty SDS-Plus Rotary Hammer Drills", NameAM: "የኮንክሪት መበሻ ሮተሪ ድራይሎች", NameOM: "Diriilii Konkiriitii Uru", NameTI: "ናይ ኮንክሪት መበስበሲ ድሪል"},
						{Slug: "angle-grinders-cutting", NameEN: "4.5-inch & 9-inch Electric Angle Grinders", NameAM: "የብረት እና የድንጋይ መቁረጫ ግራይንደሮች", NameOM: "Giraayinidara Sibilaa Muxu", NameTI: "ናይ ሓጺን መቑረጺ ግራይንደር"},
						{Slug: "cordless-impact-drivers", NameEN: "18V / 20V Lithium-Ion Cordless Impact Wrenches", NameAM: "የባትሪ ቦሎን አላቂ እና አስሪ ድራይሎች", NameOM: "Diriilii Baatirii Qabu", NameTI: "ናይ ባትሪ ቦሎን መፍቲሕ"},
						{Slug: "circular-miter-saws", NameEN: "Sliding Compound Miter and Circular Wood Saws", NameAM: "የእንጨት እና የአሉሚኒየም መጋዞች", NameOM: "Magaaza Mukaa fi Aluminiyemii", NameTI: "ናይ ዕንጨይቲ መጋዝ"},
					},
				},
				{
					SubSlug: "welding-metalworking",
					NameEN:  "Welding Machines and Plasma Cutters",
					NameAM:  "የብየዳ ማሽኖች እና መቁረጫዎች",
					NameOM:  "Maashina Sibila Wal-qabsiisuu",
					NameTI:  "ናይ ብየዳ ማሽናትን መቑረጽትን",
					Items: []ItemSpec{
						{Slug: "inverter-arc-welders", NameEN: "IGBT Inverter Arc (MMA) Welding Machines (200A-400A)", NameAM: "የኢንቨርተር ብየዳ ማሽኖች", NameOM: "Maashina Weeldara Inveertaraa", NameTI: "ናይ ኢንቨርተር ብየዳ ማሽን"},
						{Slug: "mig-mag-gas-welders", NameEN: "MIG/MAG Wire Feed Industrial Welders", NameAM: "የሽቦ ብየዳ ማሽኖች (MIG)", NameOM: "Maashina Weeldara MIG", NameTI: "ናይ ሽቦ ብየዳ ማሽን (MIG)"},
						{Slug: "air-plasma-cutting-machines", NameEN: "CNC & Handheld Air Plasma Metal Cutters", NameAM: "የፕላዝማ ብረት መቁረጫዎች", NameOM: "Kutaa Pilaazmaa Sibilaa", NameTI: "ናይ ፕላዝማ ሓጺን መቑረጺ"},
						{Slug: "welding-rods-safety-helmets", NameEN: "Welding Electrodes (E6013) & Auto-Darkening Helmets", NameAM: "የኤሌክትሮድ ብየዳ ሽቦዎች እና ጭምብሎች", NameOM: "Elekiroodii fi Koofiyyaa Eegumsaa", NameTI: "ኤሌክትሮድን መከላኸሊ ቆብን"},
					},
				},
				{
					SubSlug: "hand-tools-workshop",
					NameEN:  "Hand Tool Sets, Vises and Storage Cabinets",
					NameAM:  "የእጅ መሳሪያዎች እና ቱል ቦክሶች",
					NameOM:  "Meeshaalee Harkaa fi Saanduqa Meeshaa",
					NameTI:  "መሳርሒታት ኢድን ናይ ንዋይ ሳንዱቕን",
					Items: []ItemSpec{
						{Slug: "chrome-vanadium-socket-sets", NameEN: "1/4, 3/8, 1/2-inch Ratchet Socket Wrench Sets", NameAM: "የካቻቪቴ እና የቦሎን መፍቻ ኪቶች", NameOM: "Kachaaviitii fi Kilii Sibilaa", NameTI: "ናይ ካቻቪተን ቦሎንን መፍቲሕ"},
						{Slug: "heavy-bench-vises", NameEN: "Forged Steel Heavy Duty Swivel Bench Vises", NameAM: "የጠረጴዛ ብረት መያዣ መንጋጋዎች (ሞርሳ)", NameOM: "Qabduu Sibilaa Minjaalaa (Moorsaa)", NameTI: "ናይ ሓጺን መትሓዚ (ሞርሳ)"},
						{Slug: "rolling-tool-cabinets", NameEN: "7-Drawer Rolling Mechanic Tool Trolleys", NameAM: "ተንቀሳቃሽ የሜካኒክ ቱል ቦክስ ካቢኔቶች", NameOM: "Saanduqa Meeshaa Meekaaniikaa", NameTI: "ናይ መካኒክ ቱል ቦክስ ካቢኔት"},
						{Slug: "hydraulic-bottle-jacks", NameEN: "5-Ton to 50-Ton Hydraulic Floor & Bottle Jacks", NameAM: "የሃይድሮሊክ ክሪክ እና ጃኮች", NameOM: "Jaakii Hayidirooliikii", NameTI: "ናይ ሃይድሮሊክ ክሪክ"},
					},
				},
			},
		},

		// 29. Construction Materials and Building Supplies
		{
			DomainSlug: "construction-materials",
			Branches: []BranchSpec{
				{
					SubSlug: "cement-concrete-aggregates",
					NameEN:  "Cement, Sand and Crushed Stone Aggregates",
					NameAM:  "ሲሚንቶ፣ አሸዋ እና ጠጠር",
					NameOM:  "Simintoo, Cirracha fi Dhagaa Daakame",
					NameTI:  "ሲሚንቶ፣ ሑጻን ጠጠርን",
					Items: []ItemSpec{
						{Slug: "opc-grade-42-5-cement", NameEN: "Ordinary Portland Cement (OPC 42.5R - Dangote, Mugher)", NameAM: "ኦፒሲ 42.5 ጠንካራ ሲሚንቶ", NameOM: "Simintoo Jabaa OPC 42.5", NameTI: "OPC 42.5 ሓያል ሲሚንቶ"},
						{Slug: "ppc-grade-32-5-cement", NameEN: "Portland Pozzolana Cement (PPC 32.5N)", NameAM: "ፒፒሲ 32.5 የግንባታ ሲሚንቶ", NameOM: "Simintoo Ijaarsaa PPC 32.5", NameTI: "PPC 32.5 ናይ ህንጻ ሲሚንቶ"},
						{Slug: "coarse-river-sand", NameEN: "Washed River Concrete Sand (Per Truck Load)", NameAM: "የታጠበ የወንዝ አሸዋ (በጭነት)", NameOM: "Cirracha Lagaa Dhiqame", NameTI: "ዝተሓጸበ ናይ ሩባ ሑጻ"},
						{Slug: "crushed-basalt-gravel", NameEN: "Crushed Basalt Aggregate / Gravel (01, 02 Grade)", NameAM: "የባሳልት ጠጠር ለኮንክሪት", NameOM: "Dhagaa Cirracha Konkiriitii", NameTI: "ናይ ኮንክሪት ጠጠር"},
					},
				},
				{
					SubSlug: "reinforcement-steel-rebar",
					NameEN:  "Rebar Steel, Wire Mesh and Nails",
					NameAM:  "የግንባታ ብረት እና ቆርቆሮ",
					NameOM:  "Sibila Ijaarsaa fi Qorqoorroo",
					NameTI:  "ናይ ህንጻ ሓጺንን ቆርቆሮን",
					Items: []ItemSpec{
						{Slug: "deformed-rebar-ribbed-steel", NameEN: "Grade 60 Deformed Rebar Steel Bars (8mm to 32mm)", NameAM: "የኮንክሪት ብረት (8ሚሜ - 32ሚሜ)", NameOM: "Sibila Konkiriitii (8mm-32mm)", NameTI: "ናይ ኮንክሪት ሓጺን"},
						{Slug: "corrugated-iron-roofing-sheets", NameEN: "Galvanized Corrugated Iron Roofing Sheets (G-28 / G-32)", NameAM: "የጣሪያ ቆርቆሮዎች (G-28, G-32)", NameOM: "Qorqoorroo Baaxii Manaa", NameTI: "ናይ ናሕሲ ቆርቆሮ"},
						{Slug: "binding-wire-construction-nails", NameEN: "Black Annealed Binding Wire & Concrete Steel Nails", NameAM: "የማሰሪያ ሽቦ እና የግንባታ ሚስማር", NameOM: "Wadaroo fi Mismarii Ijaarsaa", NameTI: "ናይ ማእሰሪ ሽቦን መስማርን"},
					},
				},
				{
					SubSlug: "finishing-tiles-paints",
					NameEN:  "Ceramic Tiles, Granite and Wall Paints",
					NameAM:  "ሴራሚክ፣ ግራናይት እና ቀለሞች",
					NameOM:  "Seeraamikii, Giraanaayitii fi Qalama",
					NameTI:  "ሴራሚክ፣ ግራናይትን ቀለማትን",
					Items: []ItemSpec{
						{Slug: "porcelain-floor-tiles", NameEN: "60x60cm & 60x120cm Polished Porcelain Floor Tiles", NameAM: "የወለል ፖርሰላይን ሴራሚክ", NameOM: "Seeraamikii Wallee fi Lafaa", NameTI: "ናይ ምድሪ ሴራሚክ"},
						{Slug: "natural-granite-slabs", NameEN: "Natural Granite Countertops & Stair Slabs", NameAM: "የኪችን እና የደረጃ የተፈጥሮ ግራናይት ድንጋይ", NameOM: "Dhagaa Giraanaayitii", NameTI: "ናይ መደያይቦ ግራናይት"},
						{Slug: "interior-washable-emulsion", NameEN: "Interior Washable Acrylic Emulsion Wall Paints", NameAM: "የቤት ውስጥ የውሃ ቀለም", NameOM: "Qalama Mana Keessaa Dhiqamu", NameTI: "ናይ ውሽጢ ገዛ ቀለም"},
						{Slug: "exterior-weatherproof-quartz", NameEN: "Exterior Weatherproof Quartz & Textured Paint", NameAM: "የውጭ ግድግዳ ኳርትዝ ቀለም", NameOM: "Qalama Qilleensa Danda'u", NameTI: "ናይ ደገ ኳርትዝ ቀለም"},
					},
				},
			},
		},

		// 30. Construction Machinery and Heavy Equipment
		{
			DomainSlug: "construction-machinery",
			Branches: []BranchSpec{
				{
					SubSlug: "earthmoving-machinery",
					NameEN:  "Excavators, Wheel Loaders and Bulldozers",
					NameAM:  "ኤክስካቫተሮች፣ ሎደሮች እና ቡልዶዘሮች",
					NameOM:  "Maashinoota Qotuu fi Fe'uu",
					NameTI:  "ኤክስካቫተርን ሎደርን ቡልዶዘርን",
					Items: []ItemSpec{
						{Slug: "hydraulic-crawler-excavators", NameEN: "Hydraulic Crawler Excavators (20-30 Ton - CAT/Komatsu)", NameAM: "ትላልቅ ኤክስካቫተሮች (ቁፋሮ ማሽኖች)", NameOM: "Ekiskaavaatara Qotuu Gurguddaa", NameTI: "ናይ ቁፋሮ ማሽን (ኤክስካቫተር)"},
						{Slug: "four-wheel-drive-loaders", NameEN: "Articulated 4WD Wheel Loaders (3-5m³ Bucket)", NameAM: "የጭነት ሎደሮች (ዊል ሎደር)", NameOM: "Loodara Fe'iinsaa", NameTI: "ናይ ጽዕነት ሎደር"},
						{Slug: "crawler-track-bulldozers", NameEN: "Heavy Duty Track Bulldozers with Ripper", NameAM: "ቡልዶዘሮች ከመሬት ቆራጭ ጋር", NameOM: "Buldoozara Lafa Qixxeessu", NameTI: "ቡልዶዘር"},
						{Slug: "motor-road-graders", NameEN: "Motor Road Graders (140K / 140H)", NameAM: "የመንገድ ጠራጊ እና አስተካካይ ግሬደሮች", NameOM: "Gireedara Daandii Qixxeessu", NameTI: "ናይ መገዲ መስተኻኸሊ ግሬደር"},
					},
				},
				{
					SubSlug: "concrete-road-equipment",
					NameEN:  "Concrete Mixers, Pumps and Compaction Rollers",
					NameAM:  "የኮንክሪት ማቡኪያ እና መንገድ ማንከባለያዎች",
					NameOM:  "Maashina Konkiriitii fi Roolara",
					NameTI:  "ናይ ኮንክሪት ምሕዋስን መገዲ ማንከባለልን",
					Items: []ItemSpec{
						{Slug: "mobile-drum-concrete-mixers", NameEN: "Diesel & Electric Mobile Drum Concrete Mixers (350L-500L)", NameAM: "የናፍጣ ኮንክሪት ማቡኪያ ማሽኖች", NameOM: "Maashina Konkiriitii Makaa", NameTI: "ናይ ኮንክሪት መሕወሲ ማሽን"},
						{Slug: "hydraulic-concrete-pumps", NameEN: "Stationary & Truck-Mounted Hydraulic Concrete Pumps", NameAM: "የኮንክሪት ፓምፖች", NameOM: "Poompii Konkiriitii", NameTI: "ናይ ኮንክሪት ፓምፕ"},
						{Slug: "tandem-vibratory-rollers", NameEN: "Tandem & Single Drum Vibratory Soil Rollers", NameAM: "የአፈር እና የአስፋልት ማንከባለያዎች (ሮለር)", NameOM: "Roolara Biyyoo fi Asfaaltii", NameTI: "ናይ ሓመድን ኣስፋልትን ሮለር"},
					},
				},
			},
		},

		// 56. Energy, Solar and Electrical Power Equipment
		{
			DomainSlug: "energy-solar-power",
			Branches: []BranchSpec{
				{
					SubSlug: "solar-power-systems",
					NameEN:  "Solar Panels, Inverters and Solar Controllers",
					NameAM:  "የፀሐይ ብርሃን ኃይል እቃዎች (ሶላር)",
					NameOM:  "Meeshaalee Annisaa Soolaarii",
					NameTI:  "ናይ ጸሓይ ጸዓት መሳርሒታት (ሶላር)",
					Items: []ItemSpec{
						{Slug: "monocrystalline-solar-panels", NameEN: "Monocrystalline High-Efficiency Solar Panels (450W-650W)", NameAM: "ሞኖክሪስታላይን የሶላር ፕሌቶች", NameOM: "Paanaalii Soolaarii Humna Guddaa", NameTI: "ሞኖክሪስታላይን ናይ ሶላር ፓነል"},
						{Slug: "hybrid-pure-sine-inverters", NameEN: "Hybrid Solar Inverters with Built-in MPPT (3kVA-10kVA)", NameAM: "ሀይብሪድ የሶላር ኢንቨርተሮች", NameOM: "Inveertara Soolaarii Haayibriidii", NameTI: "ናይ ሶላር ኢንቨርተር"},
						{Slug: "lifepo4-lithium-batteries", NameEN: "LiFePO4 Lithium Solar Wall Batteries (48V 100Ah/200Ah)", NameAM: "የሊቲየም ሶላር ባትሪዎች", NameOM: "Baatirii Liitiyeemii Soolaarii", NameTI: "ናይ ሊቲየም ሶላር ባትሪ"},
						{Slug: "tubular-deep-cycle-batteries", NameEN: "Tubular Deep-Cycle Gel Batteries (12V 200Ah)", NameAM: "የዲፕ ሳይክል ጄል ባትሪዎች", NameOM: "Baatirii Diip Saayikilii", NameTI: "ናይ ጄል ባትሪ"},
					},
				},
				{
					SubSlug: "diesel-petrol-generators",
					NameEN:  "Diesel Generators and Power Backup Systems",
					NameAM:  "የናፍጣ ጀነሬተሮች እና ማመንጫዎች",
					NameOM:  "Jenereetara Diizilii fi Boba'aa",
					NameTI:  "ናይ ናፍጣ ጀነሬተራትን መመንጨዊን",
					Items: []ItemSpec{
						{Slug: "heavy-silent-diesel-generators", NameEN: "Soundproof Silent Diesel Generators (20kVA to 500kVA)", NameAM: "ድምፅ አልባ የናፍጣ ጀነሬተሮች", NameOM: "Jenereetara Diizilii Sagalee Hin Qabne", NameTI: "ድምጺ ኣልቦ ናይ ናፍጣ ጀነሬተር"},
						{Slug: "portable-petrol-generators", NameEN: "Portable Household Gasoline Generators (2.5kVA-7.5kVA)", NameAM: "ተንቀሳቃሽ የቤንዚን ጀነሬተሮች", NameOM: "Jenereetara Boba'aa Salphaa", NameTI: "ተንቀሳቓሲ ናይ ቤንዚን ጀነሬተር"},
						{Slug: "automatic-transfer-switches-ats", NameEN: "Automatic Transfer Switches (ATS Panels)", NameAM: "አውቶማቲክ የኤሌክትሪክ መቀየሪያ ሳጥኖች (ATS)", NameOM: "Iswiikii ATS Ofumaan Jijjiiru", NameTI: "ኣውቶማቲክ መቐየሪ ATS ሳንዱቕ"},
					},
				},
			},
		},

		// 57. Water, Irrigation and Environmental Equipment
		{
			DomainSlug: "water-irrigation-equipment",
			Branches: []BranchSpec{
				{
					SubSlug: "pumps-irrigation",
					NameEN:  "Submersible Pumps and Drip Irrigation Systems",
					NameAM:  "የውሃ ፓምፖች እና የመስኖ እቃዎች",
					NameOM:  "Poompii Bishaanii fi Jallisii",
					NameTI:  "ናይ ማይ ፓምፓትን ናይ መስኖ ንዋያትን",
					Items: []ItemSpec{
						{Slug: "deep-well-submersible-pumps", NameEN: "Solar & AC Deep-Well Submersible Borehole Pumps", NameAM: "የጥልቅ ጉድጓድ ሰብመርሲብል ፓምፖች", NameOM: "Poompii Boolla Bishaanii Gad-fagoo", NameTI: "ናይ ዓሚቕ ዒላ ሰብመርሲብል ፓምፕ"},
						{Slug: "drip-irrigation-kits", NameEN: "Gravity & Pressurized Drip Irrigation Line Kits", NameAM: "የጠብታ መስኖ መስመሮች እና ኪቶች", NameOM: "Ujummoo Jallisii Cophsuu", NameTI: "ናይ ነጥቢ መስኖ ሻምብቆታት"},
						{Slug: "hdpe-water-storage-tanks", NameEN: "Heavy Duty Plastic Water Storage Tanks (1,000L to 10,000L)", NameAM: "የሮቶ የፕላስቲክ ውሃ ማጠራቀሚያ ታንከሮች", NameOM: "Taanikara Bishaanii Pilaastikii", NameTI: "ናይ ማይ ፕላስቲክ ታንከር"},
						{Slug: "geomembrane-dam-liners", NameEN: "UV-Resistant Geomembrane Pond & Canal Liners", NameAM: "የኩሬ እና የቦይ መከላከያ ፕላስቲኮች (ጂኦሜምብሬን)", NameOM: "Pilaastikii Haroo Jallisiifi Bishaanii", NameTI: "ናይ ሓጽቢ መከላኸሊ ፕላስቲክ"},
					},
				},
			},
		},
	}
}

// GenerateAllCategories compiles the 60 roots and builds distinct, valid nodes
func GenerateAllCategories() ([]NestedCategoryDef, int) {
	categories := GetNestedCategories()
	specs := GetAllDomainSpecs()

	order := 100
	for _, spec := range specs {
		for _, branch := range spec.Branches {
			order++
			// Add branch as non-leaf subcategory under root
			categories = append(categories, NestedCategoryDef{
				ParentSlug:   spec.DomainSlug,
				Slug:         branch.SubSlug,
				NameEN:       branch.NameEN,
				NameAM:       branch.NameAM,
				NameOM:       branch.NameOM,
				NameTI:       branch.NameTI,
				IsLeaf:       false,
				DisplayOrder: order,
			})

			// Add items as terminal leaf categories under branch
			for itemIdx, item := range branch.Items {
				categories = append(categories, NestedCategoryDef{
					ParentSlug:   branch.SubSlug,
					Slug:         item.Slug,
					NameEN:       item.NameEN,
					NameAM:       item.NameAM,
					NameOM:       item.NameOM,
					NameTI:       item.NameTI,
					IsLeaf:       true,
					DisplayOrder: itemIdx + 1,
				})
			}
		}
	}

	// Count leaf nodes
	leafCount := 0
	for _, c := range categories {
		if c.IsLeaf {
			leafCount++
		}
	}

	return categories, leafCount
}

func main_dummy() {
	_, leaves := GenerateAllCategories()
	fmt.Printf("Generated leaves: %d\n", leaves)
}
