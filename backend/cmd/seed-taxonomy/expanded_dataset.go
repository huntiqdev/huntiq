package main

import (
	"fmt"
	"strings"
)

type leafItem struct {
	slug   string
	nameEN string
	nameAM string
	nameOM string
	nameTI string
}

type subTemplate struct {
	slug   string
	nameEN string
	nameAM string
	nameOM string
	nameTI string
	items  []leafItem
}

type domainTemplate struct {
	domainSlug string
	subs       []subTemplate
}

// GenerateExpandedDomainLeaves builds authentic subcategories and leaf items for marketplace domains
func GenerateExpandedDomainLeaves() []NestedCategoryDef {
	var list []NestedCategoryDef

	templates := []domainTemplate{
		// 01. Electronics and Mobile Devices
		{
			domainSlug: "electronics-mobile",
			subs: []subTemplate{
				{
					slug: "feature-phones-keypads", nameEN: "Keypad & Basic Mobile Phones", nameAM: "ባለ ኪፓድ ተንቀሳቃሽ ስልኮች", nameOM: "Bilbila Kiipaadii Salphaa", nameTI: "ናይ ቁልፊ ሞባይል ስልኪ",
					items: []leafItem{
						{"nokia-torch-phones", "Nokia Classic FM & Flashlight Keypad Phones", "የኖኪያ የእጅ ባትሪ ስልኮች", "Bilbila Nooriyaa Faanoo Qabu", "ናይ ኖኪያ ፍላሽ ዘለዎ ስልኪ"},
						{"itel-dual-sim-phones", "itel Long-Battery Dual-SIM Feature Phones", "የአይቴል ባትሪ ቆጣቢ ስልኮች", "Bilbila Itel Baatirii Dheeraa", "ናይ ኣይቴል ስልኪ"},
						{"elderly-large-keypad-phones", "Large Button & SOS Emergency Phones for Seniors", "ለትላልቅ ሰዎች የሚሆኑ ትላልቅ ኪፓድ ስልኮች", "Bilbila Kiipaadii Gurguddaa Manguddootaa", "ናይ ዓበይቲ ሰባት ስልኪ"},
					},
				},
				{
					slug: "mobile-repair-tools", nameEN: "Phone Repair Tools, Heat Guns and Microscopes", nameAM: "የሞባይል ስልክ ጥገና እቃዎች", nameOM: "Meeshaalee Suphaa Bilbilaa", nameTI: "ናይ ሞባይል ጽገና መሳርሒታት",
					items: []leafItem{
						{"smd-hot-air-rework-stations", "SMD Hot Air Rework Stations & Soldering Irons", "የኤስኤምዲ ሆት ኤር የብየዳ ማሽኖች", "Maashina Suphaa Ho'aa SMD", "ናይ SMD ሆት ኤር ብየዳ"},
						{"lcd-screen-separator-machines", "Vacuum LCD Screen Separator & Heating Plates", "የስክሪን ማለያያ ማሽኖች", "Maashina Iskiriinii Addaan Baasu", "ናይ ስክሪን መፈልለዪ ማሽን"},
						{"digital-repair-microscopes", "Trinocular Stereo Zoom Repair Microscopes", "የኤሌክትሮኒክስ ጥገና ማይክሮስኮፖች", "Maayikirooskooppii Suphaa", "ናይ ጽገና ማይክሮስኮፕ"},
						{"precision-screwdriver-sets", "Precision Screwdriver Sets for iPhone and Android", "የሞባይል መፍቻ ካቻቪቴ ኪቶች", "Kachaaviitii Qalloo Bilbilaa", "ናይ ሞባይል ካቻቪተ ኪት"},
					},
				},
			},
		},

		// 02. Computers, IT and Office Technology
		{
			domainSlug: "computers-it",
			subs: []subTemplate{
				{
					slug: "computer-networking-accessories", nameEN: "Ethernet Cables, USB Hubs and Adapters", nameAM: "የኮምፒውተር ገመዶች እና መገናኛዎች", nameOM: "Keebilii fi Meeshaalee Kompiitaraa", nameTI: "ናይ ኮምፒዩተር ኬብላትን ኣገናኘትን",
					items: []leafItem{
						{"usb-c-multiport-hubs", "USB-C to HDMI, Ethernet & Card Reader Hubs", "የዩኤስቢ-ሲ ማገናኛ ሃቦች (Hubs)", "USB-C Hubii Hedduu Qabu", "ናይ USB-C ኣገናኚ ሃብ"},
						{"high-speed-hdmi-displayport-cables", "4K/8K Braided HDMI and DisplayPort Video Cables", "የኤችዲኤምአይ የቪዲዮ ገመዶች", "Keebilii HDMI fi DisplayPort", "ናይ HDMI ቪዲዮ ኬብል"},
						{"cooling-pads-laptop-stands", "Adjustable Aluminum Laptop Stands & Cooling Pads", "የላፕቶፕ ማቀዝቀዣ እና ማስቀመጫዎች", "Qabbaneessituu fi Qabduu Laptooppii", "ናይ ላፕቶፕ መቐመጢን መዝሓሊን"},
					},
				},
				{
					slug: "uninterruptible-power-supplies", nameEN: "Uninterruptible Power Supplies (UPS) and Inverters", nameAM: "ዩፒኤስ (UPS) እና የሃይል ማከማቻዎች", nameOM: "Meeshaa Humna Elektiriikii UPS", nameTI: "ዩፒኤስ (UPS) መሳርሒታት",
					items: []leafItem{
						{"line-interactive-offline-ups", "650VA to 2000VA Line-Interactive Office UPS", "የቢሮ ዩፒኤስ (650VA - 2000VA)", "UPS Biiroo 650VA-2000VA", "ናይ ቢሮ ዩፒኤስ"},
						{"online-double-conversion-ups", "3kVA to 20kVA Online Double Conversion Server UPS", "የሰርቨር ኦንላይን ዩፒኤስ (3kVA+)", "UPS Sarvarii Humna Guddaa", "ናይ ሰርቨር ዩፒኤስ"},
						{"mini-router-backup-ups", "Mini 12V DC UPS for Wi-Fi Routers & Modems", "የዋይፋይ ራውተር ሚኒ ዩፒኤስ", "Mini UPS Raawutara Wi-Fi", "ናይ ዋይፋይ ራውተር ሚኒ ዩፒኤስ"},
					},
				},
			},
		},

		// 08. Furniture and Interior Furnishings
		{
			domainSlug: "furniture-furnishings",
			subs: []subTemplate{
				{
					slug: "dining-kitchen-furniture", nameEN: "Dining Tables, Chairs and Sideboards", nameAM: "የመመገቢያ ጠረጴዛዎች እና ወንበሮች", nameOM: "Minjaala Nyaataa fi Barcuma", nameTI: "ናይ መመገቢ ጠረጴዛን መንበርን",
					items: []leafItem{
						{"6-8-seater-dining-sets", "6-Seater and 8-Seater Hardwood Dining Table Sets", "የ6 እና የ8 ሰው መመገቢያ ጠረጴዛዎች", "Minjaala Nyaataa Nama 6-8", "ናይ 6-8 ሰብ መመገቢ ጠረጴዛ"},
						{"wooden-buffet-sideboards", "Modern Wooden Kitchen Buffets and Display Cabinets", "የቡፌ እና የዕቃ ማስቀመጫ ካቢኔቶች", "Kaabineetii Nyaataa fi Buufee", "ናይ ቡፌን ኣቑሑትን ካቢኔት"},
						{"bar-stools-counter-chairs", "High Swivel Bar Stools & Kitchen Island Chairs", "የባር እና የኪችን ደሴት ከፍታ ወንበሮች", "Barcuma Dheeraa Baarii", "ናይ ባር ነዋሕቲ መንበር"},
					},
				},
				{
					slug: "home-decor-carpets", nameEN: "Carpets, Rugs, Curtains and Wall Mirrors", nameAM: "ምንጣፎች፣ መጋረጃዎች እና የቤት ማስዋቢያዎች", nameOM: "Afata Lafaa, Golgaa fi Daawwitii", nameTI: "ምንጻፋት፣ መጋረጃታትን ናይ ገዛ መመላኸዕታትን",
					items: []leafItem{
						{"turkish-persian-style-carpets", "Large Turkish & Persian Pattern Living Room Carpets", "የሳሎን ምንጣፎች (ቱርክ)", "Afata Lafaa Saaloonaa", "ናይ ሳሎን ምንጻፍ"},
						{"blackout-luxury-curtains", "Double Layer Blackout Window Curtains and Sheers", "የመስኮት መጋረጃዎች እና ሳሳዎች", "Golgaa Foddaa Qananii", "ናይ መስኮት መጋረጃታት"},
						{"decorative-led-wall-mirrors", "LED Backlit Decorative Vanity Wall Mirrors", "የግድግዳ ጌጣጌጥ መስተዋቶች", "Daawwitii Wallee Miidhaginaa", "ናይ ግድግዳ መስትያት"},
					},
				},
			},
		},

		// 09. Home, Kitchen and Household Goods
		{
			domainSlug: "kitchen-household",
			subs: []subTemplate{
				{
					slug: "cookware-pots-pans", nameEN: "Stainless Steel Pots, Frying Pans and Pressure Cookers", nameAM: "ድስቶች፣ መጥበሻዎች እና የፕሬዠር ኩከሮች", nameOM: "Distii, Waadduu fi Meeshaa Bilcheessaa", nameTI: "ድስታት፣ መጥበስታትን ፕሬዠር ኩከርን",
					items: []leafItem{
						{"stainless-steel-pot-sets", "Multi-Piece Stainless Steel Cooking Pot Sets", "የስቲል ድስት ስብስቦች", "Distii Sibiila Qulqulluu", "ናይ ስቲል ድስቲታት"},
						{"granite-nonstick-frying-pans", "Granite & Ceramic Non-Stick Frying Pans", "የግራናይት መጥበሻዎች", "Waadduu Giraanaayitii Hin Qabanne", "ናይ ግራናይት መጥበሲ"},
						{"heavy-pressure-cookers", "Commercial & Home Aluminum Pressure Cookers (7L-15L)", "የአልሙኒየም ፕሬዠር ኩከሮች", "Kukarii Dhiibbaa Jabaa", "ናይ ፕሬዠር ኩከር ድስቲ"},
					},
				},
				{
					slug: "traditional-coffee-ceremony-items", nameEN: "Ethiopian Coffee Ceremony Sets, Jebena and Rekebot", nameAM: "የቡና ረከቦቶች፣ ጀበናዎች እና ሲኒዎች", nameOM: "Rakabootii, Jabaynaa fi Siinii Bunaa", nameTI: "ረከቦት፣ ጀበናን ሲኒን",
					items: []leafItem{
						{"wooden-handcrafted-rekebot", "Wooden Handcrafted Coffee Ceremony Rekebot Trays", "የእንጨት የቡና ረከቦቶች", "Rakabootii Mukaa Aadaa", "ናይ ዕንጨይቲ ረከቦት"},
						{"traditional-ceramic-sini-cups", "Classic Saba & Tilf Pattern Porcelain Coffee Sini Cups (Set of 12)", "የሳባ እና የጥልፍ ዲዛይን የቡና ሲኒዎች", "Siinii Bunaa Aadaa Saabaa", "ናይ ሳባ ዲዛይን ሲኒ ቡን"},
						{"clay-charcoal-stoves-gwadad", "Traditional Clay Charcoal Burners (Maderaja / Gwadad)", "የቡና ማፍያ የከሰል ማጨሻዎች", "Eelee Cilaattii Bunaa", "ናይ ከሰል መንደዲ ጓዳድ"},
					},
				},
			},
		},

		// 24. Office Supplies and Commercial Equipment
		{
			domainSlug: "office-supplies",
			subs: []subTemplate{
				{
					slug: "commercial-cash-handling", nameEN: "Banknote Counters, Safes and Cash Drawers", nameAM: "ገንዘብ ቆጣሪ ማሽኖች እና ካዝናዎች", nameOM: "Maashina Lakkoofsa Maallaqaa fi Kaaznaa", nameTI: "ገንዘብ መቑጸሪ ማሽንን ካዝናን",
					items: []leafItem{
						{"uv-fake-note-bill-counters", "High-Speed Cash Counting Machines with UV/MG Counterfeit Detection", "የሀሰተኛ ብር መለያ ገንዘብ ቆጣሪ ማሽኖች", "Maashina Maallaqa Lakkaawu fi Sooba Qoratu", "ናይ ሓሶት ቅርሺ መለሊ ገንዘብ መቑጸሪ"},
						{"fireproof-digital-safes", "Heavy Fireproof Electronic Digital Safe Boxes (30kg-150kg)", "የእሳት መቋቋሚ የዲጂታል ካዝናዎች", "Kaaznaa Balaa Ibiddarraa Eegu", "ሓዊ መኸላኸሊ ዲጂታል ካዝና"},
						{"heavy-metal-cash-drawers", "Heavy Gauge Steel Point-of-Sale (POS) Cash Drawers", "የሱቅ ገንዘብ መያዣ መሳቢያዎች", "Saanduqa Maallaqaa POS Sibiilaa", "ናይ POS ገንዘብ መቐመጢ"},
					},
				},
				{
					slug: "paper-stationery-filing", nameEN: "Printing Paper, Lever Arch Files and Binders", nameAM: "ኤ4 ወረቀቶች እና የሰነድ መያዣ ፎልደሮች", nameOM: "Waraqaa A4 fi Foolderii Sanadaa", nameTI: "A4 ወረቐትን ናይ ሰነድ ፎልደርን",
					items: []leafItem{
						{"a4-copier-paper-boxes", "80gsm A4 High-White Copier Paper (Box of 5 Reams)", "የኤ4 ፎቶኮፒ ወረቀት ካርቶኖች (80 ግራም)", "Kaartoonii Waraqaa A4 80gsm", "ካርቶን A4 ወረቐት"},
						{"heavy-lever-arch-box-files", "Cardboard & Plastic Lever Arch Box Files", "የሰነድ ማስቀመጫ ትላልቅ ፋይሎች (Box Files)", "Foolderii Sanadaa Gurguddaa", "ናይ ሰነድ ቦክስ ፋይል"},
						{"electric-thermal-binding-machines", "Comb & Thermal Document Binding Machines", "የወረቀት መበሻ እና ማሰሪያ ማሽኖች", "Maashina Sanada Wal-qabsiisu", "ናይ ሰነድ መበስበሲን መኣሰርን ማሽን"},
					},
				},
			},
		},

		// 25. Retail, Wholesale and Distribution
		{
			domainSlug: "retail-wholesale",
			subs: []subTemplate{
				{
					slug: "bulk-fmcg-commodities", nameEN: "Bulk Edible Oil, Sugar, Rice and FMCG Goods", nameAM: "የጅምላ የምግብ እና የፍጆታ እቃዎች", nameOM: "Daldala Jumlaa Nyaataa fi Zayitaa", nameTI: "ናይ ጅምላ መግብን ዘይትን",
					items: []leafItem{
						{"commercial-cooking-oil-drums", "Bulk Palm & Sunflower Cooking Oil Jerrycans (20L) & Drums (200L)", "የምግብ ዘይት በጀሪካን እና በበርሜል (ጅምላ)", "Zayita Nyaataa Jerikaana 20L fi Drum", "ናይ መግቢ ዘይቲ ብጅምላ"},
						{"refined-white-sugar-bags", "Refined White Granulated Sugar (50kg Factory Bags)", "የስኳር ከረጢቶች በጅምላ (50 ኪ.ግ)", "Sukkaara Adii Korojoo 50kg", "ናይ ሽኮር ከረጢት ብጅምላ (50 ኪ.ግ)"},
						{"imported-basmati-long-grain-rice", "Long Grain Basmati & Thai Parboiled Rice (25kg Bags)", "የባስማቲ እና የታይላንድ ሩዝ (25 ኪ.ግ)", "Ruuza Baasmaatii 25kg", "ናይ ባስማቲ ሩዝ (25 ኪ.ግ)"},
					},
				},
				{
					slug: "commercial-cleaning-detergents", nameEN: "Wholesale Soaps, Detergents and Sanitizers", nameAM: "የጅምላ ሳሙናዎች፣ ዱቄቶች እና ኬሚካሎች", nameOM: "Saamunaa fi Qulqulleessituu Jumlaa", nameTI: "ናይ ጅምላ ሳሙናታትን መጽረዪን",
					items: []leafItem{
						{"laundry-detergent-powder-bags", "Industrial Laundry Powder Detergent (20kg-25kg Sacks)", "የልብስ ማጠቢያ ዱቄት ሳሙና በጅምላ", "Daakuu Saamunaa Uffataa 20kg", "ናይ ክዳን ሕጻብ ሓሩጭ ሳሙና"},
						{"liquid-hand-dishwashing-soaps", "Bulk Liquid Handwash & Dishwashing Detergent (5L-20L)", "ፈሳሽ የእጅ እና የእቃ ማጠቢያ ሳሙና", "Saamunaa Dhangala'aa Harkaafi Qodaa", "ፈሳሲ ሳሙና"},
						{"multipurpose-bar-laundry-soaps", "Traditional 555 & Sunlight Multipurpose Laundry Bars", "የ555 እና የባር ሳሙና ካርቶኖች በጅምላ", "Saamunaa Dhagaa 555 Kaartoonii", "ናይ 555 ሳሙና ካርቶን"},
					},
				},
			},
		},

		// 27. Industrial Machinery and Manufacturing
		{
			domainSlug: "industrial-machinery",
			subs: []subTemplate{
				{
					slug: "food-processing-milling-plants", nameEN: "Grain Milling, Oil Press and Bakery Machines", nameAM: "የእህል ወፍጮ፣ የዘይት መጭመቂያ እና የዳቦ ማሽኖች", nameOM: "Maashina Midhaan Daakuu fi Zayitaa", nameTI: "ናይ እኽሊ መጥሓንን ዘይቲ መጭመቕን ማሽን",
					items: []leafItem{
						{"commercial-flour-hammer-mills", "Heavy Duty Electric Grain & Teff Hammer Mills", "የጤፍ እና የእህል ወፍጮ ማሽኖች", "Wafchoo Midhaanii fi Xaafii Elektiriikii", "ናይ ጣፍን እኽልን መጥሓን ማሽን"},
						{"screw-cold-press-oil-expellers", "Continuous Screw Oil Expeller Machines (Noug & Sesame)", "የኑግ እና የሰሊጥ ዘይት መጭመቂያ ማሽኖች", "Maashina Zayita Nuugii fi Saliixaa Cuunfu", "ናይ ኑግን ሰሊጥን ዘይቲ መጭመቒ"},
						{"rotary-deck-bakery-ovens", "Commercial Rotary Deck Bread & Pastry Baking Ovens", "የዳቦ እና የኬክ መጋገሪያ ትላልቅ እቶኖች", "Oovonii Daabboo fi Keekii Gurguddaa", "ናይ ባኒ መጋገዪ ዓቢይ እቶን"},
					},
				},
				{
					slug: "plastic-metal-packaging-machines", nameEN: "Plastic Molding, Blowing and Sealing Machinery", nameAM: "የፕላስቲክ ማምረቻ እና የማሸጊያ ማሽኖች", nameOM: "Maashina Pilaastikii fi Saamsuu", nameTI: "ናይ ፕላስቲክ መፍረዪን መዐሸግን ማሽን",
					items: []leafItem{
						{"plastic-injection-molding-machines", "Automatic Hydraulic Plastic Injection Molding Machines", "የፕላስቲክ እቃዎች ማምረቻ ማሽኖች (ኢንጄክሽን)", "Maashina Pilaastikii Tolchu Injeekshinii", "ናይ ፕላስቲክ መፍረዪ ማሽን"},
						{"pet-bottle-stretch-blow-molders", "Semi-Automatic PET Bottle Blowing Machines (Water/Juice)", "የውሃ ሃይላንድ ፕላስቲክ ጠርሙስ ማምረቻ ማሽኖች", "Maashina Qaruuraa Bishaanii Tolchu", "ናይ ሃይላንድ ጥርሙዝ መፍረዪ ማሽን"},
						{"continuous-band-bag-sealers", "Continuous Horizontal Nitrogen Flush Band Sealers", "የከረጢት እና የፕላስቲክ ማሸጊያ ማሽኖች", "Maashina Pilaastikii Saamsu", "ናይ ፕላስቲክ ከረጢት መዐሸጊ ማሽን"},
					},
				},
			},
		},

		// 31. Architecture, Engineering and Surveying
		{
			domainSlug: "architecture-engineering",
			subs: []subTemplate{
				{
					slug: "architectural-design-services", nameEN: "Building Blueprints, 3D Rendering and Interior Design", nameAM: "የህንጻ ፕላን፣ ዲዛይን እና የውስጥ ማስዋብ", nameOM: "Dizaayinii Pilaanii Manaa fi 3D", nameTI: "ናይ ህንጻ ፕላንን ዲዛይንን",
					items: []leafItem{
						{"full-architectural-construction-blueprints", "G+1 to G+10 Complete Architectural & Structural Blueprints", "የቪላ እና የህንጻ ሙሉ የግንባታ ፕላኖች", "Pilaanii Ijaarsaa Guutuu Viillaafi Gamoo", "ናይ ቪላን ህንጻን ሙሉእ ፕላን"},
						{"photorealistic-3d-exterior-rendering", "Photorealistic 3D Exterior & Interior Architectural Renderings", "የ3ዲ ህንጻ እይታ እና ዲዛይን ስራዎች", "Fakkii 3D Gamoo Miidhagaa", "ናይ 3D ህንጻ ስእላዊ መግለጺ"},
						{"bill-of-quantities-boq-takeoff", "Comprehensive Bill of Quantities (BOQ) & Cost Estimation", "የግንባታ ወጪ ግምት እና ቢኦኪው (BOQ)", "Shallaggii Baasii Ijaarsaa (BOQ)", "ናይ ህንጻ ወጻኢ ገማቲ (BOQ)"},
					},
				},
				{
					slug: "surveying-geotechnical-services", nameEN: "Topographic Land Surveying and Soil Testing", nameAM: "የመሬት ቅየሳ እና የአፈር ምርመራ", nameOM: "Qorannoo Lafaafi Biyyoo", nameTI: "ናይ መሬት ቅየሳን ሓመድ ምርመራን",
					items: []leafItem{
						{"total-station-topographic-surveys", "Total Station & DGPS High-Precision Boundary Surveys", "የቶታል ስቴሽን የመሬት ቅየሳ ስራዎች", "Safartuu Lafa DGPS fi Total Station", "ናይ ቶታል ስቴሽን መሬት ቅየሳ"},
						{"geotechnical-soil-investigation-boreholes", "Geotechnical Core Drilling & Foundation Soil Testing", "የመሰረት አፈር ጥናት እና ቁፋሮ", "Qorannoo Biyyoo Hundeessaa Gamoo", "ናይ መሠረት ሓመድ መጽናዕቲ"},
					},
				},
			},
		},
	}

	for _, dt := range templates {
		for subIdx, sub := range dt.subs {
			// Subcategory
			list = append(list, NestedCategoryDef{
				ParentSlug:   dt.domainSlug,
				Slug:         sub.slug,
				NameEN:       sub.nameEN,
				NameAM:       sub.nameAM,
				NameOM:       sub.nameOM,
				NameTI:       sub.nameTI,
				IsLeaf:       false,
				DisplayOrder: subIdx + 1,
			})

			// Items
			for itemIdx, item := range sub.items {
				list = append(list, NestedCategoryDef{
					ParentSlug:   sub.slug,
					Slug:         item.slug,
					NameEN:       item.nameEN,
					NameAM:       item.nameAM,
					NameOM:       item.nameOM,
					NameTI:       item.nameTI,
					IsLeaf:       true,
					DisplayOrder: itemIdx + 1,
				})
			}
		}
	}

	return list
}

// GenerateAllCategoriesExpanded merges all category definitions
func GenerateAllCategoriesExpanded() []NestedCategoryDef {
	var merged []NestedCategoryDef
	seen := make(map[string]bool)

	addUnique := func(item NestedCategoryDef) {
		slug := strings.ToLower(strings.TrimSpace(item.Slug))
		if seen[slug] {
			slug = fmt.Sprintf("%s-%s", item.ParentSlug, slug)
		}
		seen[slug] = true
		item.Slug = slug
		merged = append(merged, item)
	}

	for _, c := range Build60DomainCatalog() {
		addUnique(c)
	}
	for _, c := range GenerateExpandedDomainLeaves() {
		addUnique(c)
	}

	return merged
}
