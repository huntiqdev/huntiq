package main

import "fmt"

// GetCatalogDomains31To40 returns comprehensive domain catalogues for roots 31 to 40
func GetCatalogDomains31To40() []DomainCatalogDef {
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

	// Domain: architecture-engineering
	list = append(list, DomainCatalogDef{
		RootSlug: "architecture-engineering",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "architectural-structural-design",
				NameEN: "Architectural Design, 3D Rendering and Structural Engineering",
				NameAM: "የአርክቴክቸር ዲዛይን፣ 3ዲ እይታ እና የመዋቅር ኢንጂነሪንግ",
				NameOM: "Dizaayinii Gamoo, 3D fi Injinariingii",
				NameTI: "ዲዛይን ኣርክቴክቸርን መዋቕራዊ ምህንድስናን",
				Items: []LeafItemDef{
					makeLeaf("residential-villa-g-plus-design", "G+1, G+2 & Luxury Residential Villa Architectural Plan Design", "የቪላ እና የመኖሪያ ቤቶች የአርክቴክቸር ፕላን ዲዛይን (G+1/G+2)", "Dizaayinii Viillaa fi Mana Jireenyaa", "ናይ G+1, G+2 & Luxury Residential Villa Architectural Plan Design"),
					makeLeaf("high-rise-commercial-building-architecture", "Multi-Storey Mixed-Use Commercial Mall & Office Architectural Design", "የንግድ ህንፃዎች እና የገበያ ማዕከላት የአርክቴክቸር ዲዛይን", "Dizaayinii Gamoo Daldalaa Dheeraa", "ናይ Multi-Storey Mixed-Use Commercial Mall & Office Architectural Design"),
					makeLeaf("reinforced-concrete-structural-engineering", "ETABS/SAP2000 Structural Calculation & Seismic Analysis Reports", "የኮንክሪት መዋቅር (ስትራክቸራል) ስሌት እና የዲዛይን ሪፖርቶች", "Herrega Injinariingii Gamoo (Structural)", "ናይ ETABS/SAP2000 Structural Calculation & Seismic Analysis Reports"),
					makeLeaf("photorealistic-3d-exterior-rendering", "Ultra-HD 3D Exterior Architectural Visualization & Walkthrough Animation", "የህንፃዎች የውጪ ገጽታ 3ዲ አኒሜሽን እና እይታ (3D Rendering)", "Agarsiisa Gamoo 3D fi Viidiyoo", "ናይ Ultra-HD 3D Exterior Architectural Visualization & Walkthrough Animation"),
					makeLeaf("luxury-interior-design-fitout-plans", "Modern Apartment & Office Interior Space Planning & Material Specs", "የአፓርታማ እና የቢሮ የውስጥ ማስዋብ ዲዛይን (Interior Design)", "Dizaayinii Keessa Manaa (Interior)", "ናይ Modern Apartment & Office Interior Space Planning & Material Specs"),
					makeLeaf("mep-engineering-services-design", "HVAC, Fire Fighting, Electrical & Sanitary MEP Comprehensive Engineering", "የእሳት መከላከያ፣ የአየር ማቀዝቀዣ እና የቧንቧ ዝርጋታ ኢንጂነሪንግ (MEP)", "Dizaayinii Elektiriikii fi Ujummoo (MEP)", "ናይ HVAC, Fire Fighting, Electrical & Sanitary MEP Comprehensive Engineering"),
					makeLeaf("landscape-architecture-urban-planning", "Public Park, Compound Hardscaping & Urban Greenery Master Planning", "የግቢ ማሳመር፣ አትክልት እና የመናፈሻ ቦታዎች ዲዛይን (Landscape)", "Dizaayinii Dirree fi Naannoo Magariisaa", "ናይ Public Park, Compound Hardscaping & Urban Greenery Master Planning"),
					makeLeaf("as-built-drawing-documentation-services", "Post-Construction CAD As-Built Verification & Municipality Submission", "የተገነቡ ህንፃዎች ትክክለኛ የአዝ-ቢልት (As-Built) ፕላን ሰነዶች", "Waraqaa Pilaanii Gamoo Dhaabbataa (As-Built)", "ናይ Post-Construction CAD As-Built Verification & Municipality Submission"),
				},
			},
			{
				Slug:   "surveying-geotechnical-soil-testing",
				NameEN: "Land Surveying, Soil Testing and GIS Services",
				NameAM: "የመሬት ቅየሳ፣ የአፈር ምርመራ እና ጂአይኤስ",
				NameOM: "Safarsa Lafaa, Qorannoo Biyyoo fi GIS",
				NameTI: "ምዕቃን መሬትን ምምርማር ሓመድን",
				Items: []LeafItemDef{
					makeLeaf("total-station-topographical-land-surveying", "Total Station & RTK GPS Boundary & Contour Topographical Surveying", "በትክክለኛ ቶታል ስቴሽን የመሬት ወሰን እና ኮንቱር ቅየሳ", "Safarsa Lafaa Tootaal Isteeshinii", "ናይ Total Station & RTK GPS Boundary & Contour Topographical Surveying"),
					makeLeaf("geotechnical-core-drilling-soil-investigation", "Borehole Soil Core Drilling, Standard Penetration Tests (SPT) & Reports", "ለህንፃ መሰረት የሚደረግ የከርሰ ምድር አፈር ምርመራ (Soil Test)", "Qorannoo Biyyoo Bu'uura Gamoo", "ናይ Borehole Soil Core Drilling, Standard Penetration Tests (SPT) & Reports"),
					makeLeaf("drone-photogrammetry-aerial-mapping", "High-Resolution RTK Drone Aerial Photogrammetry & Orthomosaic Maps", "በድሮን የሚሰራ የላይኛው የመሬት ካርታ እና የጂአይኤስ ቅየሳ", "Kaartaa Lafaa Dirooniin Hojjetamu", "ናይ High-Resolution RTK Drone Aerial Photogrammetry & Orthomosaic Maps"),
					makeLeaf("cadastral-boundary-demarcation-verification", "Urban Cadastral Property Boundary Verification and Dispute Demarcation", "የይዞታ ማረጋገጫ እና የመሬት ወሰን ፍተሻ አገልግሎቶች", "Mirkaneessa Daangaa Lafaa fi Qabeenyaa", "ናይ Urban Cadastral Property Boundary Verification and Dispute Demarcation"),
					makeLeaf("structural-integrity-audit-ultrasonic-testing", "Non-Destructive Rebound Hammer & Ultrasonic Concrete Integrity Audits", "የነባር ህንፃዎች ጥንካሬ እና ደህንነት ፍተሻ (Audit)", "Qorannoo Jabinna Gamoo Dhaabbataa", "ናይ Non-Destructive Rebound Hammer & Ultrasonic Concrete Integrity Audits"),
					makeLeaf("road-highway-geometric-alignment-design", "Civil 3D Highway Geometric Alignment, Earthwork Volume & Pavement Design", "የአስፋልት እና የጠጠር መንገድ አሰላለፍ እና የመሬት ስራ ዲዛይን", "Dizaayinii Daandii Asfaaltii fi Cirrachaa", "ናይ Civil 3D Highway Geometric Alignment, Earthwork Volume & Pavement Design"),
					makeLeaf("hydrological-drainage-flood-risk-modeling", "Watershed Delineation, Culvert Sizing & Urban Stormwater Flood Studies", "የጎርፍ አደጋ መከላከያ እና የፍሳሽ ቦይ ጥናትና ዲዛይን", "Qorannoo Loola fi Dhangala'aa Bishaanii", "ናይ Watershed Delineation, Culvert Sizing & Urban Stormwater Flood Studies"),
					makeLeaf("gis-spatial-database-cartography-services", "ArcGIS Enterprise Land Cover, Demographic & Spatial Thematic Mapping", "የጂአይኤስ (GIS) እና የካርቶግራፊ የመረጃ ካርታዎች ዝግጅት", "Tajaajila Kaartaa fi Daataa GIS", "ናይ ArcGIS Enterprise Land Cover, Demographic & Spatial Thematic Mapping"),
				},
			},
		},
	})

	// Domain: plumbing-electrical
	list = append(list, DomainCatalogDef{
		RootSlug: "plumbing-electrical",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "electrical-wiring-distribution-lighting",
				NameEN: "Wires, Breakers, Distribution Boards and LED Lighting",
				NameAM: "የኤሌክትሪክ ሽቦዎች፣ ሰርኪውት ብሬከሮች እና መብራቶች",
				NameOM: "Wayaara Elektiriikii, Bireekara fi Ibsaa LED",
				NameTI: "ናይ ኤሌክትሪክ ሽቦታትን መብራህቲ LEDን",
				Items: []LeafItemDef{
					makeLeaf("copper-single-core-building-wires-2-5mm", "Pure Annealed Copper PVC Insulated 2.5mm & 4mm House Wiring Coils", "2.5 እና 4 ሚሜ የቤት ውስጥ የኤሌክትሪክ ሽቦዎች በጥቅል", "Wayaara Elektiriikii Manaa 2.5mm", "ናይ Pure Annealed Copper PVC Insulated 2.5mm & 4mm House Wiring Coils"),
					makeLeaf("molded-case-circuit-breakers-mccb", "3-Phase 100A-400A Main Distribution Industrial MCCB Breakers", "የዋና መስመር የኤሌክትሪክ ሰርኪውት ብሬከሮች (MCCB)", "Bireekara Elektiriikii Gurguddaa (MCCB)", "ናይ 3-Phase 100A-400A Main Distribution Industrial MCCB Breakers"),
					makeLeaf("miniature-circuit-breakers-mcb-din-rail", "10A to 63A Single & Double Pole DIN-Rail Distribution Board MCBs", "የቤት ውስጥ ቆጣቢ አነስተኛ ብሬከሮች (MCB)", "Bireekara Elektiriikii Xiqqaataa (MCB)", "ናይ 10A to 63A Single & Double Pole DIN-Rail Distribution Board MCBs"),
					makeLeaf("led-ceiling-panel-recessed-downlights", "18W-36W Ultra-Slim Energy Saving Warm/Cool White LED Ceiling Panels", "የኮርኒስ ውስጥ ውብ የኤልኢዲ መብራቶች (Downlight)", "Ibsaa LED Gubbaa Manaa (Panel)", "ናይ 18W-36W Ultra-Slim Energy Saving Warm/Cool White LED Ceiling Panels"),
					makeLeaf("weatherproof-outdoor-security-floodlights", "100W-300W IP66 Waterproof Aluminum Compound LED Floodlights", "ለግቢ እና ለደህንነት የሚሆኑ ደማቅ የውጪ መብራቶች (Floodlight)", "Ibsaa Dirree LED Humna Guddaa", "ናይ 100W-300W IP66 Waterproof Aluminum Compound LED Floodlights"),
					makeLeaf("flush-mount-modular-light-switches-sockets", "Modern Luxury Tempered Glass Touch & Rocker Wall Socket Outlets", "ዘመናዊ ግድግዳ ላይ የሚገጠሙ የኤሌክትሪክ ስዊቾች እና ፕለጎች", "Iswiichii fi Sokkeettii Girgidaa", "ናይ Modern Luxury Tempered Glass Touch & Rocker Wall Socket Outlets"),
					makeLeaf("heavy-industrial-rubber-cable-extensions", "3-Core 2.5mm Oil-Resistant Drum Cable Extensions with Sockets", "ለግንባታ ስራ የሚሆኑ ጠንካራ የኤክስቴንሽን ገመዶች በጥቅል", "Keebilii Elektiriikii Dheeraa (Extension)", "ናይ 3-Core 2.5mm Oil-Resistant Drum Cable Extensions with Sockets"),
					makeLeaf("three-phase-manual-changeover-switches", "63A-250A Dual Power Generator-to-Grid Rotary Changeover Switches", "የጄኔሬተር እና የኤሌክትሪክ ኃይል መቀየሪያ ቼንጅ ኦቨሮች", "Iswiichii Jenereetara fi Humnaa (Changeover)", "ናይ 63A-250A Dual Power Generator-to-Grid Rotary Changeover Switches"),
				},
			},
			{
				Slug:   "plumbing-pipes-pumps-sanitary-fittings",
				NameEN: "PPR Pipes, Water Pumps, Tanks and Sanitary Ware",
				NameAM: "የውሃ ቧንቧዎች፣ የውሃ ታንከሮች፣ ፓምፖች እና መጸዳጃዎች",
				NameOM: "Ujummoo Bishaanii, Taankii fi Paampii",
				NameTI: "ናይ ማይ ሻምብቆታትን ታንክን ፓምፕን",
				Items: []LeafItemDef{
					makeLeaf("ppr-hot-cold-water-supply-pipes-25mm", "PN20 High-Pressure Heat-Fused Green PPR Plumbing Pipes & Sockets", "ለቀዝቃዛ እና ለሙቅ ውሃ የሚሆኑ ጠንካራ የፒፒአር (PPR) ቧንቧዎች", "Ujummoo Bishaanii PPR Magariisa", "ናይ PN20 High-Pressure Heat-Fused Green PPR Plumbing Pipes & Sockets"),
					makeLeaf("hdpe-corrugated-underground-drainage-pipes", "SN8 High-Density Polyethylene Heavy Sewer & Drainage Pipes", "ለመሬት ውስጥ ፍሳሽ የሚቀበሩ ጠንካራ የኤችዲፒኢ (HDPE) ቧንቧዎች", "Ujummoo Xurii Lafaa HDPE", "ናይ SN8 High-Density Polyethylene Heavy Sewer & Drainage Pipes"),
					makeLeaf("vertical-cylindrical-plastic-water-tanks-5000l", "Multi-Layer Food-Grade UV-Stabilized 5,000L-10,000L Water Rototanks", "ባለ 5,000 እና 10,000 ሊትር የውሃ ፕላስቲክ ታንከሮች (ሮቶ)", "Taankii Bishaanii Pilaastikii 5,000L", "ናይ Multi-Layer Food-Grade UV-Stabilized 5,000L-10,000L Water Rototanks"),
					makeLeaf("automatic-peripheral-booster-water-pumps", "0.5HP-1.5HP High-Head Cast Iron Clean Water Pressure Booster Pumps", "የቤት ውስጥ የውሃ ግፊት መጫኛ ኃይለኛ ፓምፖች", "Paampii Dhiibbaa Bishaanii (Water Pump)", "ናይ 0.5HP-1.5HP High-Head Cast Iron Clean Water Pressure Booster Pumps"),
					makeLeaf("submersible-deep-well-borehole-pumps", "Stainless Steel Multi-Stage Solar & Electric Deep Borehole Pumps", "የከርሰ ምድር ጥልቅ የውሃ ጉድጓድ ሰመጥ ፓምፖች (Submersible)", "Paampii Boolla Bishaan Keessaa", "ናይ Stainless Steel Multi-Stage Solar & Electric Deep Borehole Pumps"),
					makeLeaf("ceramic-dual-flush-close-coupled-toilets", "Water-Saving S-Trap Porcelain Bathroom Commode Toilet Sets", "የሸክላ ውብ የሽንት ቤት መጸዳጃዎች (ኮሞድ)", "Qodaa Mana Finiinaa Qulqulluu (Toilet)", "ናይ Water-Saving S-Trap Porcelain Bathroom Commode Toilet Sets"),
					makeLeaf("stainless-steel-double-bowl-kitchen-sinks", "Heavy 304 Stainless Steel Soundproof Undermount Kitchen Sinks", "ባለሁለት ገንዳ የማይዝግ የኩሽና የእቃ ማጠቢያ ሲንኮች", "Siinkii Kushiinaa Sibilaa Hin Naxne", "ናይ Heavy 304 Stainless Steel Soundproof Undermount Kitchen Sinks"),
					makeLeaf("luxury-brass-rainfall-shower-mixer-sets", "Thermostatic Mixer Valve Handheld & Overhead Rainfall Shower Sets", "የመታጠቢያ ቤት ዘመናዊ የሻወር መርጫዎች እና ቧንቧዎች", "Shaaworii fi Ujummoo Dhiqannaa Qananii", "ናይ Thermostatic Mixer Valve Handheld & Overhead Rainfall Shower Sets"),
				},
			},
		},
	})

	// Domain: property-sale
	list = append(list, DomainCatalogDef{
		RootSlug: "property-sale",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "apartments-condominiums-for-sale",
				NameEN: "Apartments and Condominiums for Sale",
				NameAM: "የሚሸጡ አፓርታማዎች እና ኮንዶሚኒየሞች",
				NameOM: "Apaartamaa fi Kondominiyemii Gurguramu",
				NameTI: "ዝሽየጡ ኣፓርታማታትን ኮንዶሚኒየምን",
				Items: []LeafItemDef{
					makeLeaf("luxury-three-bedroom-apartments-sale", "Master Bedroom En-Suite Modern 150m2 High-Rise Luxury Apartments", "ባለ 3 መኝታ ዘመናዊ የቅንጦት አፓርታማዎች ሽያጭ (150 ካሬ)", "Apaartamaa Qananii Kutaa 3 Gurguramu", "ናይ Master Bedroom En-Suite Modern 150m2 High-Rise Luxury Apartments"),
					makeLeaf("two-bedroom-semi-finished-apartments", "100m2 Flexible Layout Real Estate Developer Semi-Finished Units", "ባለ 2 መኝታ ከፊል ያለቀላቸው አፓርታማዎች ሽያጭ", "Apaartamaa Kutaa 2 Walakkaa Dhumate", "ናይ 100m2 Flexible Layout Real Estate Developer Semi-Finished Units"),
					makeLeaf("studio-one-bedroom-apartments-sale", "55m2 Compact Prime Urban Center Starter Apartments for Sale", "ስቱዲዮ እና ባለ 1 መኝታ አፓርታማዎች ሽያጭ", "Apaartamaa Isteediyoo fi Kutaa 1", "ናይ 55m2 Compact Prime Urban Center Starter Apartments for Sale"),
					makeLeaf("government-condominiums-40-60-sale", "Commercial Bank Mortgage Approved 40/60 Housing Units for Sale", "የ40/60 የመንግስት የኮንዶሚኒየም ቤቶች ሽያጭ", "Manneen Kondominiyemii 40/60 Gurguramu", "ናይ Commercial Bank Mortgage Approved 40/60 Housing Units for Sale"),
					makeLeaf("government-condominiums-20-80-sale", "Fully Titled 1-3 Bedroom 20/80 Residential Condominiums for Sale", "የ20/80 የኮንዶሚኒየም ቤቶች ሽያጭ", "Manneen Kondominiyemii 20/80 Gurguramu", "ናይ Fully Titled 1-3 Bedroom 20/80 Residential Condominiums for Sale"),
					makeLeaf("luxury-penthouse-panoramic-view-sale", "300m2 Top-Floor Duplex Private Terrace Luxury Penthouses", "ከተማውን የሚያሳዩ ባለ ቴራስ የቅንጦት ፔንት ሃውስ ቤቶች", "Apaartamaa Peentihaawusii Qananii", "ናይ 300m2 Top-Floor Duplex Private Terrace Luxury Penthouses"),
					makeLeaf("under-construction-off-plan-apartments", "Installment Payment Staged Milestone Off-Plan Real Estate Projects", "በግንባታ ላይ ያሉ በረጅም ጊዜ ክፍያ የሚሸጡ አፓርታማዎች", "Apaartamaa Ijaaramaa Jiru Kaffaltiin", "ናይ Installment Payment Staged Milestone Off-Plan Real Estate Projects"),
					makeLeaf("serviced-hotel-apartments-guaranteed-roi", "Fully Furnished Rental Yield Hotel Suite Apartments for Sale", "የተሟላ እቃ ያላቸው ገቢ የሚያስገኙ የሆቴል አፓርታማዎች", "Apaartamaa Hoteelaa Galii Qabu", "ናይ Fully Furnished Rental Yield Hotel Suite Apartments for Sale"),
				},
			},
			{
				Slug:   "residential-villas-houses-for-sale",
				NameEN: "Residential Villas, G+1, G+2 and Compounds for Sale",
				NameAM: "የሚሸጡ ቪላዎች፣ ጂ+1፣ ጂ+2 እና የመኖሪያ ቤቶች",
				NameOM: "Viillaa fi Mana Jireenyaa Gurguramu",
				NameTI: "ዝሽየጡ ቪላታትን ናይ መኖርያ ኣባይትን",
				Items: []LeafItemDef{
					makeLeaf("luxury-g-plus-two-residential-villas", "5-Bedroom Modern Architectural G+2 Standalone Luxury Villas", "ባለ 5 መኝታ ዘመናዊ ጂ+2 የመኖሪያ ቪላ ቤቶች ሽያጭ", "Viillaa Qananii G+2 Gurguramu", "ናይ 5-Bedroom Modern Architectural G+2 Standalone Luxury Villas"),
					makeLeaf("g-plus-one-family-villas-with-service-quarters", "Spacious Compound G+1 Villas with 3-Room Service Quarters", "የሰርቪስ ክፍል እና ሰፊ ግቢ ያላቸው ጂ+1 ቪላ ቤቶች", "Viillaa Maatii G+1 Kutaa Tajaajilaa Qabu", "ናይ Spacious Compound G+1 Villas with 3-Room Service Quarters"),
					makeLeaf("traditional-single-storey-residential-villas", "Solid Basalt Stone Built 250m2 Compound Traditional Villas", "በድንጋይ የተሰሩ 250 ካሬ ምቹ ቪላ ቤቶች ሽያጭ", "Viillaa Dachee Tokkoo Lafa Bal'aa", "ናይ Solid Basalt Stone Built 250m2 Compound Traditional Villas"),
					makeLeaf("gated-community-real-estate-villas", "24/7 Security CCTV Protected Private Neighborhood Villas", "ጥበቃ ባለበት ግቢ ውስጥ ያሉ ዘመናዊ ቪላ ቤቶች", "Viillaa Eegumsa Qabu (Gated Community)", "ናይ 24/7 Security CCTV Protected Private Neighborhood Villas"),
					makeLeaf("old-demolition-value-residential-houses", "Prime Location Residential Houses Sold for Plot Redevelopment", "ለመኖሪያ ወይም ለህንፃ መስሪያ የሚሆኑ የፈረሱ ቤቶች ሽያጭ", "Mana Jireenyaa Lafasaaf Gurguramu", "ናይ Prime Location Residential Houses Sold for Plot Redevelopment"),
					makeLeaf("townhouse-attached-cluster-villas", "Contemporary 3-Storey Attached Townhouse Urban Homes", "ዘመናዊ የታውን ሀውስ (Townhouse) የመኖሪያ ቤቶች", "Mana Jireenyaa Taawunhaawusii (Townhouse)", "ናይ Contemporary 3-Storey Attached Townhouse Urban Homes"),
					makeLeaf("unfinished-grey-structure-villas-sale", "Reinforced Concrete Frame Roofed Grey Structure Houses for Sale", "የኮንክሪት መዋቅሩ ያለቀለት የሚሸጥ ቪላ ቤት", "Viillaa Dhaabbanni Ijaarame Gurguramu", "ናይ Reinforced Concrete Frame Roofed Grey Structure Houses for Sale"),
					makeLeaf("luxury-countryside-resort-style-villas", "1000m2 Scenic Highland Fresh-Air Weekend Getaway Homes", "ከከተማ ውጪ ያሉ ሰፊ የተፈጥሮ አየር ማረፊያ ቪላዎች", "Viillaa Boqonnaa fi Turiizimii Bal'aa", "ናይ 1000m2 Scenic Highland Fresh-Air Weekend Getaway Homes"),
				},
			},
		},
	})

	// Domain: property-rent
	list = append(list, DomainCatalogDef{
		RootSlug: "property-rent",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "apartments-studios-for-rent",
				NameEN: "Furnished and Unfurnished Apartments for Rent",
				NameAM: "የሚከራዩ አፓርታማዎች፣ ስቱዲዮዎች እና ኮንዶሚኒየሞች",
				NameOM: "Apaartamaa fi Isteediyoo Kireeffamu",
				NameTI: "ዝካረዩ ኣፓርታማታትን ስቱዲዮን",
				Items: []LeafItemDef{
					makeLeaf("fully-furnished-expat-luxury-apartments", "Diplomatic Standard 2-3 Bedroom Fully Furnished Luxury Units", "የተሟላ እቃ ያላቸው የዲፕሎማት ደረጃ የቅንጦት አፓርታማዎች ኪራይ", "Apaartamaa Meeshaan Guutame Kireeffamu", "ናይ Diplomatic Standard 2-3 Bedroom Fully Furnished Luxury Units"),
					makeLeaf("unfurnished-spacious-family-apartments", "3-Bedroom 2-Bath Clean Modern Unfurnished Apartments for Lease", "እቃ የሌላቸው ሰፊ የቤተሰብ አፓርታማዎች ኪራይ", "Apaartamaa Meeshaa Malee Kireeffamu", "ናይ 3-Bedroom 2-Bath Clean Modern Unfurnished Apartments for Lease"),
					makeLeaf("condominium-apartments-for-rent", "1-2 Bedroom Tile Finished Secure Condominium Rentals", "የኮንዶሚኒየም የመኖሪያ ቤቶች ኪራይ", "Manneen Kondominiyemii Kireeffamu", "ናይ 1-2 Bedroom Tile Finished Secure Condominium Rentals"),
					makeLeaf("furnished-studio-compact-apartments", "All Utilities Included Monthly Studio Apartments for Single Expats", "መብራት እና ውሃ ያላቸው አነስተኛ ስቱዲዮ አፓርታማዎች", "Apaartamaa Isteediyoo Meeshaa Qabu", "ናይ All Utilities Included Monthly Studio Apartments for Single Expats"),
					makeLeaf("short-term-daily-weekly-guest-apartments", "Airbnb Style Wi-Fi, Kitchen & Housekeeping Serviced Suites", "በቀን እና በሳምንት የሚከራዩ የእንግዳ ማረፊያ አፓርታማዎች", "Apaartamaa Guyyaan Kireeffamu (Guest House)", "ናይ Airbnb Style Wi-Fi, Kitchen & Housekeeping Serviced Suites"),
					makeLeaf("penthouse-terrace-apartments-for-lease", "Exclusive Top-Floor Panoramic View Penthouses for Rent", "የቅንጦት ፔንት ሃውስ ቤቶች ኪራይ", "Apaartamaa Peentihaawusii Kireeffamu", "ናይ Exclusive Top-Floor Panoramic View Penthouses for Rent"),
					makeLeaf("shared-room-student-co-living-rentals", "Private Bedroom in Shared Kitchen Modern Apartment for Students", "የተማሪ እና የወጣቶች የክፍል ኪራይ (ከጋራ ኩሽና ጋር)", "Kutaa Manaa Qooddannaadhaan Kireeffamu", "ናይ Private Bedroom in Shared Kitchen Modern Apartment for Students"),
					makeLeaf("basement-service-room-rentals", "Budget Single Service Bedrooms in Compound with Shared Utilities", "በግቢ ውስጥ ያሉ አነስተኛ የሰርቪስ ክፍሎች ኪራይ", "Kutaa Tajaajilaa Gadi Aanaa Kireeffamu", "ናይ Budget Single Service Bedrooms in Compound with Shared Utilities"),
				},
			},
			{
				Slug:   "villas-commercial-compounds-for-rent",
				NameEN: "Villas, Offices and Commercial Compounds for Lease",
				NameAM: "የሚከራዩ ቪላዎች፣ የቢሮ ህንፃዎች እና ግቢዎች",
				NameOM: "Viillaa fi Gamoo Kireeffamu",
				NameTI: "ዝካረዩ ቪላታትን ናይ ቢሮ ህንጻታትን",
				Items: []LeafItemDef{
					makeLeaf("diplomatic-compound-luxury-villas", "Bulletproof Glass, Generator & Large Guarded Garden Embassy Villas", "ለኤምባሲ እና ዲፕሎማቶች የሚሆኑ ሰፊ የቪላ ግቢዎች ኪራይ", "Viillaa Embaasii fi Diippiloomaatotaa", "ናይ Bulletproof Glass, Generator & Large Guarded Garden Embassy Villas"),
					makeLeaf("residential-g-plus-one-family-villas", "Spacious 4-Bedroom G+1 Compound Residential Villa Rentals", "ባለ 4 መኝታ የቤተሰብ መኖሪያ ጂ+1 ቪላ ቤቶች ኪራይ", "Viillaa Jireenya Maatii G+1 Kireeffamu", "ናይ Spacious 4-Bedroom G+1 Compound Residential Villa Rentals"),
					makeLeaf("commercial-office-villa-conversions", "Paved Parking for 15+ Vehicles NGO Headquarters Office Villas", "ለድርጅት ወይም ለNGO ቢሮ የሚሆኑ ሰፊ ቪላዎች ኪራይ", "Viillaa Biiroo fi NGO Kireeffamu", "ናይ Paved Parking for 15+ Vehicles NGO Headquarters Office Villas"),
					makeLeaf("prime-streetfront-commercial-buildings", "Multi-Storey Open-Floor Bank & Insurance Commercial Showrooms", "ለባንክ እና ኢንሹራንስ ዋና መንገድ ላይ ያሉ የንግድ ህንፃዎች ኪራይ", "Gamoo Daldalaa Daandii Guddaarraa", "ናይ Multi-Storey Open-Floor Bank & Insurance Commercial Showrooms"),
					makeLeaf("open-floor-office-spaces-per-square-meter", "Fitted Elevator & Generator Backup Corporate Office Floor Leases", "በካሬ የሚከራዩ ዘመናዊ የቢሮ ክፍሎች", "Bakka Biiroo Kaaree Meetiiraan", "ናይ Fitted Elevator & Generator Backup Corporate Office Floor Leases"),
					makeLeaf("ground-floor-retail-shop-storefronts", "High-Traffic Pedestrian Shopping Center Retail Spaces for Rent", "የገበያ ማዕከል ውስጥ ያሉ የምድር ላይ ሱቆች ኪራይ", "Suuqii Daldalaa Seensarraa Kireeffamu", "ናይ High-Traffic Pedestrian Shopping Center Retail Spaces for Rent"),
					makeLeaf("walled-storage-warehouse-compounds", "Industrial Truck Accessible 1000m2 Dry Goods Logistics Compounds", "ትላልቅ የጭነት መኪኖች የሚገቡባቸው የመጋዘን ግቢዎች ኪራይ", "Giddugala Kuusaa fi Lojistiiksii (Warehouse)", "ናይ Industrial Truck Accessible 1000m2 Dry Goods Logistics Compounds"),
					makeLeaf("restaurant-lounge-commercial-properties", "Fully Permitted Kitchen, Bar & Outdoor Seating Food Venues for Lease", "ለሬስቶራንት እና ላውንጅ የሚሆኑ የተዘጋጁ ቦታዎች ኪራይ", "Bakka Nyaataa fi Laawunjiif Kireeffamu", "ናይ Fully Permitted Kitchen, Bar & Outdoor Seating Food Venues for Lease"),
				},
			},
		},
	})

	// Domain: land-farms
	list = append(list, DomainCatalogDef{
		RootSlug: "land-farms",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "residential-commercial-plots",
				NameEN: "Residential Plots, Commercial Land and Leases",
				NameAM: "የመኖሪያ፣ የንግድ እና የሊዝ መሬቶች",
				NameOM: "Lafa Jireenyaa fi Daldalaa",
				NameTI: "ናይ መኖርያን ንግድን መሬት",
				Items: []LeafItemDef{
					makeLeaf("residential-plots-with-title-deed-carta", "200m2 to 500m2 Fully Titled (Carta) Urban Residential Plots", "ካርታ እና ህጋዊ ሰነድ ያላቸው ከ200 እስከ 500 ካሬ የመኖሪያ ቦታዎች", "Lafa Jireenyaa Kaartaa Qabu (Carta)", "ናይ 200m2 to 500m2 Fully Titled (Carta) Urban Residential Plots"),
					makeLeaf("prime-commercial-investment-land-plots", "Roadside High-Density Zoning Commercial Building Plots for Sale", "ዋና መንገድ ዳር ለንግድ ህንፃ የሚሆኑ የኢንቨስትመንት መሬቶች", "Lafa Daldalaa fi Gamoo Daandiirraa", "ናይ Roadside High-Density Zoning Commercial Building Plots for Sale"),
					makeLeaf("real-estate-development-master-plots", "Multi-Hectare Subdivided Land Parcels for Housing Developments", "ለሪል ስቴት እና ሰፈር ግንባታ የሚሆኑ ትላልቅ ሄክታር መሬቶች", "Lafa Bal'aa Misooma Riil Isteetiif", "ናይ Multi-Hectare Subdivided Land Parcels for Housing Developments"),
					makeLeaf("industrial-park-manufacturing-land-leases", "Heavy Power & Water Connected Industrial Zone Manufacturing Plots", "በኢንዱስትሪ ፓርክ ውስጥ ለፋብሪካ የሚሆኑ የሊዝ መሬቶች", "Lafa Warshaa Paarkii Industirii", "ናይ Heavy Power & Water Connected Industrial Zone Manufacturing Plots"),
					makeLeaf("fuel-gas-station-approved-plots", "Highway Frontage Safety Distance Approved Petrol Station Land Sites", "ለነዳጅ ማደያ ግንባታ የተፈቀደላቸው ዋና መንገድ ዳር መሬቶች", "Lafa Buufata Boba'aa Eeyyamame", "ናይ Highway Frontage Safety Distance Approved Petrol Station Land Sites"),
					makeLeaf("suburban-weekend-villa-expansion-plots", "Scenic Hillside 1000m2 Suburban Fresh Air Plots with Water Lines", "ከከተማ ወጣ ብለው ለቪላ የሚሆኑ የተፈጥሮ አየር መሬቶች", "Lafa Viillaa Qilleensa Gaarii", "ናይ Scenic Hillside 1000m2 Suburban Fresh Air Plots with Water Lines"),
					makeLeaf("hospital-medical-center-zoned-land", "Accessible Urban Health Facility Zoned Institutional Plots", "ለሆስፒታል እና ጤና ተቋም ግንባታ የተፈቀደላቸው መሬቶች", "Lafa Hospitaalaa fi Fayyaa Eeyyamame", "ናይ Accessible Urban Health Facility Zoned Institutional Plots"),
					makeLeaf("warehouse-logistics-truck-terminal-land", "Wide Ring Road Accessible Heavy Fleet Logistics Depot Sites", "ለከባድ መኪኖች መጋዘን እና ተርሚናል የሚሆኑ ሰፊ መሬቶች", "Lafa Kuusaa fi Konkolaataa Fe'iinsaa", "ናይ Wide Ring Road Accessible Heavy Fleet Logistics Depot Sites"),
				},
			},
			{
				Slug:   "agricultural-farmland-horticulture-leases",
				NameEN: "Commercial Farmland, Horticulture and Irrigated Sites",
				NameAM: "የእርሻ መሬቶች፣ የመስኖ እና የሆርቲካልቸር እርሻዎች",
				NameOM: "Lafa Qonnaa fi Bishaan Qabaa",
				NameTI: "ናይ ሕርሻን መስኖን መሬት",
				Items: []LeafItemDef{
					makeLeaf("perennial-river-irrigated-commercial-farms", "50 to 500 Hectares Riverfront Year-Round Irrigated Farmlands", "በወንዝ ዳርቻ አመቱን ሙሉ በመስኖ የሚለሙ ትላልቅ እርሻዎች", "Lafa Qonnaa Bishaan Lagaa Qabu", "ናይ 50 to 500 Hectares Riverfront Year-Round Irrigated Farmlands"),
					makeLeaf("avocado-fruit-orchard-development-farms", "High-Yield Hass Avocado & Mango Plantation Farmlands for Lease", "የሃስ አቮካዶ እና የማንጎ ፍራፍሬ ማምረቻ እርሻዎች", "Lafa Oomisha Habaqadoo fi Maangoo", "ናይ High-Yield Hass Avocado & Mango Plantation Farmlands for Lease"),
					makeLeaf("export-greenhouse-flower-horticulture-farms", "Climate Controlled Polyhouse Cut-Flower Investment Farm Sites", "የአበባ እና የሆርቲካልቸር ኤክስፖርት የግሪን ሀውስ እርሻዎች", "Lafa Abaaboo fi Hooltikaalcharii", "ናይ Climate Controlled Polyhouse Cut-Flower Investment Farm Sites"),
					makeLeaf("dairy-cattle-fattening-livestock-ranches", "Fenced Grazing Pastureland & Deep Water Borehole Livestock Farms", "ለከብት እርባታ እና ማደለቢያ የሚሆኑ የታጠሩ ሰፊ እርሻዎች", "Lafa Horii fi Hoolaa Horsiisaa", "ናይ Fenced Grazing Pastureland & Deep Water Borehole Livestock Farms"),
					makeLeaf("poultry-egg-broiler-farm-facilities", "Biosecure Commercial Broiler & Layer Poultry Shed Facilities", "የዶሮ እርባታ እና የእንቁላል ማምረቻ የተሟላላቸው እርሻዎች", "Lafa Lukkuu fi Buuphaa Horsiisuu", "ናይ Biosecure Commercial Broiler & Layer Poultry Shed Facilities"),
					makeLeaf("teff-wheat-cereal-mechanized-farms", "Flat Fertile Black Soil Combine Harvester Compatible Grain Lands", "ለስንዴ እና ጤፍ እርሻ የሚሆኑ ሰፊ ለም የሜዳ እርሻዎች", "Lafa Midhaanii fi Qamadii Bal'aa", "ናይ Flat Fertile Black Soil Combine Harvester Compatible Grain Lands"),
					makeLeaf("commercial-coffee-forest-plantation-land", "Shade-Grown Arabica Organic Forest Specialty Coffee Farms", "የተፈጥሮ የደን የቡና እርሻ መሬቶች (አረቢካ)", "Lafa Qonna Buna Bosonaa Arabica", "ናይ Shade-Grown Arabica Organic Forest Specialty Coffee Farms"),
					makeLeaf("commercial-beekeeping-apiculture-forest-sites", "High-Flora Organic Forest Land Ideal for Modern Beehives", "ለንብ ማነብ እና ንጹህ ማር ማምረት ምቹ የሆኑ የደን ቦታዎች", "Lafa Damma fi Kannisa Horsiisuu", "ናይ High-Flora Organic Forest Land Ideal for Modern Beehives"),
				},
			},
		},
	})

	// Domain: commercial-real-estate
	list = append(list, DomainCatalogDef{
		RootSlug: "commercial-real-estate",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "warehouses-logistics-distribution-centers",
				NameEN: "Warehouses, Storage Sheds and Distribution Hubs",
				NameAM: "መጋዘኖች፣ የማከማቻ ሼዶች እና የሎጂስቲክስ ማዕከላት",
				NameOM: "Giddugala Kuusaa fi Lojistiiksii",
				NameTI: "መጋዘናትን ናይ ዕቑር ንዋያትን",
				Items: []LeafItemDef{
					makeLeaf("clear-span-industrial-steel-warehouses-1000m2", "10-Meter Eaves Height Fire Hydrant Clear Span Storage Warehouses", "ባለ 10 ሜትር ከፍታ ያላቸው 1000 ካሬ ዘመናዊ መጋዘኖች", "Kuusaa Sibilaa Bal'aa 1000m2", "ናይ 10-Meter Eaves Height Fire Hydrant Clear Span Storage Warehouses"),
					makeLeaf("bonded-customs-controlled-dry-ports", "Heavy Paved Container Yard Secured Bonded Customs Depots", "የጉምሩክ ክሊራንስ እና የኮንቴነር ማከማቻ ደረቅ ወደቦች", "Buufata Qorqorroo fi Gumurukii (Dry Port)", "ናይ Heavy Paved Container Yard Secured Bonded Customs Depots"),
					makeLeaf("temperature-controlled-cold-storage-facilities", "Sub-Zero Meat & Vegetable Controlled Atmosphere Cold Stores", "የስጋ እና ፍራፍሬ ማቀዝቀዣ ትላልቅ ኮልድ ስቶሬጅ መጋዘኖች", "Kuusaa Qabbanaawaa Foonii fi Fuduraa", "ናይ Sub-Zero Meat & Vegetable Controlled Atmosphere Cold Stores"),
					makeLeaf("pharmaceutical-gdp-compliant-warehouses", "Temperature Mapped Humidity Controlled Medical Storage Depots", "የመድኃኒት ማከማቻ ዓለምአቀፍ ደረጃቸውን የጠበቁ መጋዘኖች", "Kuusaa Qorichaa Qulqullina Eeggate", "ናይ Temperature Mapped Humidity Controlled Medical Storage Depots"),
					makeLeaf("commercial-truck-loading-dock-warehouses", "Hydraulic Leveler Loading Bays Multi-Door Distribution Centers", "የጭነት መኪና እቃ ማውረጃ ያላቸው ዘመናዊ መጋዘኖች", "Kuusaa Fe'iinsa Konkolaataa Gurguddaa", "ናይ Hydraulic Leveler Loading Bays Multi-Door Distribution Centers"),
					makeLeaf("small-business-storage-sheds-200m2", "Secure Urban Core 200m2 Light Inventory Mini Storage Units", "አነስተኛ እቃ ማከማቻ 200 ካሬ የከተማ መጋዘኖች", "Kuusaa Xiqqaataa Daldalaa 200m2", "ናይ Secure Urban Core 200m2 Light Inventory Mini Storage Units"),
					makeLeaf("bulk-grain-storage-metal-silos", "Aerated Temperature Sensor Industrial Grain Elevator Storage Silos", "የእህል ማከማቻ ዘመናዊ የብረት ሳይሎዎች (Silo)", "Saayiloo Kuusaa Midhaanii Sibilaa", "ናይ Aerated Temperature Sensor Industrial Grain Elevator Storage Silos"),
					makeLeaf("hazardous-chemical-approved-storage-sheds", "Explosion Proof Ventilated Spill Containment Industrial Depots", "የኬሚካል እና ተቀጣጣይ ነገሮች ማከማቻ ልዩ መጋዘኖች", "Kuusaa Keemikaalaa Nageenya Qabu", "ናይ Explosion Proof Ventilated Spill Containment Industrial Depots"),
				},
			},
			{
				Slug:   "commercial-complexes-office-parks-showrooms",
				NameEN: "Commercial Buildings, Showrooms and Office Towers",
				NameAM: "የንግድ ህንፃዎች፣ የመኪና ማሳያዎች እና የገበያ አዳራሾች",
				NameOM: "Gamoo Daldalaa, Mall fi Showroom",
				NameTI: "ናይ ንግዲ ህንጻታትን ሾውሩምን",
				Items: []LeafItemDef{
					makeLeaf("prime-roadside-vehicle-showrooms", "Floor-to-Ceiling Glass Frontage 40+ Vehicle Dealership Showrooms", "ዋና መንገድ ዳር ያሉ የመኪና መሸጫ ግዙፍ ሾውሩሞች (Showroom)", "Showroom Konkolaataa Daandii Guddaarraa", "ናይ Floor-to-Ceiling Glass Frontage 40+ Vehicle Dealership Showrooms"),
					makeLeaf("multi-storey-corporate-headquarter-towers", "Complete 10-Floor High-Visibility Prestigious Office Towers", "ሙሉ በሙሉ ለድርጅት የሚከራዩ ባለ 10 ፎቅ የቢሮ ህንፃዎች", "Gamoo Biiroo Dheeraa Dhaabbataa", "ናይ Complete 10-Floor High-Visibility Prestigious Office Towers"),
					makeLeaf("modern-retail-shopping-mall-buildings", "Multi-Level Escalator Fitted Food Court Shopping Mall Complexes", "ኤስካሌተር እና የመመገቢያ አዳራሽ ያላቸው ትላልቅ የገበያ ማዕከላት", "Gamoo Daldalaa Bal'aa (Shopping Mall)", "ናይ Multi-Level Escalator Fitted Food Court Shopping Mall Complexes"),
					makeLeaf("turnkey-hospital-medical-center-buildings", "Oxygen Piped Operating Theater Ready Private Hospital Complexes", "ለሆስፒታል ተብለው የተሰሩ የተሟሉ የህክምና ህንፃዎች", "Gamoo Hospitaalaa fi Yaalaaf Qophaa'e", "ናይ Oxygen Piped Operating Theater Ready Private Hospital Complexes"),
					makeLeaf("accredited-private-school-campus-properties", "Classrooms, Library, Science Labs & Sports Field School Campuses", "ለግል ትምህርት ቤት የተሟላላቸው የግቢ እና የህንፃ ንብረቶች", "Gamoo fi Dirree Mana Barumsaa", "ናይ Classrooms, Library, Science Labs & Sports Field School Campuses"),
					makeLeaf("banquet-hall-conference-event-centers", "1000+ Guest Capacity Soundproof Wedding & Corporate Event Venues", "የሰርግ እና የስብሰባ ትላልቅ አዳራሾች (Banquet Hall)", "Galma Cidhaa fi Walga'ii Gurguddaa", "ናይ 1000+ Guest Capacity Soundproof Wedding & Corporate Event Venues"),
					makeLeaf("boutique-hotel-resort-properties", "30-50 Room Operating Hotel Buildings with Bar & Restaurant Licenses", "ባለ 30 እና 50 ክፍል የሆቴል ህንፃዎች እና ሪዞርቶች", "Gamoo Hoteelaa fi Boqonnaa", "ናይ 30-50 Room Operating Hotel Buildings with Bar & Restaurant Licenses"),
					makeLeaf("commercial-bank-ready-branch-premises", "Armored Glass Teller Counter Strongroom Cash Vault Fitted Spaces", "ለባንክ ቅርንጫፍ ተስማሚ የሆኑ ካዝና እና መስተዋት ያላቸው ቦታዎች", "Bakka Baankii fi Kaaznaa Qophaa'e", "ናይ Armored Glass Teller Counter Strongroom Cash Vault Fitted Spaces"),
				},
			},
		},
	})

	// Domain: vehicles-cars
	list = append(list, DomainCatalogDef{
		RootSlug: "vehicles-cars",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "suvs-4wd-crossovers",
				NameEN: "SUVs, 4WD and Crossovers",
				NameAM: "የከተማ እና የገጠር 4x4 ጂፖች እና ኤስዩቪዎች",
				NameOM: "Konkolaataa SUV, 4x4 fi Crossover",
				NameTI: "ናይ 4x4 ጂፓትን SUVን",
				Items: []LeafItemDef{
					makeLeaf("toyota-land-cruiser-v8-series", "Toyota Land Cruiser 300, 200, Prado & Hardtop Rugged 4WD SUVs", "ቶዮታ ላንድ ክሩዘር V8፣ ፕራዶ እና ሀርድቶፕ ጂፖች", "Toyota Land Cruiser V8 fi Prado", "ናይ Toyota Land Cruiser 300, 200, Prado & Hardtop Rugged 4WD SUVs"),
					makeLeaf("compact-crossover-suvs", "Toyota RAV4, Hyundai Tucson & Kia Sportage Fuel Efficient Crossovers", "ቶዮታ ራቭ4 እና ሃዩንዳይ ቱክሰን የቤተሰብ መኪኖች", "Toyota RAV4 fi Hyundai Tucson", "ናይ Toyota RAV4, Hyundai Tucson & Kia Sportage Fuel Efficient Crossovers"),
					makeLeaf("mitsubishi-pajero-nissan-patrol", "Mitsubishi Pajero & Nissan Patrol Heavy 4x4 Expedition Vehicles", "ሚትሱቢሺ ፓጄሮ እና ኒሳን ፓትሮል ጠንካራ ጂፖች", "Mitsubishi Pajero fi Nissan Patrol", "ናይ Mitsubishi Pajero & Nissan Patrol Heavy 4x4 Expedition Vehicles"),
					makeLeaf("luxury-german-suvs-mercedes-bmw", "Mercedes-Benz G-Class, GLE, BMW X5 & Audi Q7 Premium Luxury SUVs", "መርሴዲስ ቤንዝ እና ቢኤምደብሊው የቅንጦት ጂፖች", "Mercedes-Benz fi BMW X5 Qananii", "ናይ Mercedes-Benz G-Class, GLE, BMW X5 & Audi Q7 Premium Luxury SUVs"),
					makeLeaf("suzuki-jimny-compact-4x4", "Suzuki Jimny 3-Door & 5-Door AllGrip Lightweight Off-Roaders", "ሱዙኪ ጂምኒ አነስተኛ እና ጠንካራ የ4x4 መኪኖች", "Suzuki Jimny 4x4 Xiqqaataa", "ናይ Suzuki Jimny 3-Door & 5-Door AllGrip Lightweight Off-Roaders"),
					makeLeaf("land-rover-range-rover-defenders", "Land Rover Defender 110 & Range Rover Autobiography Air Suspension SUVs", "ላንድ ሮቨር ዲፌንደር እና ሬንጅ ሮቨር የቅንጦት መኪኖች", "Land Rover Defender fi Range Rover", "ናይ Land Rover Defender 110 & Range Rover Autobiography Air Suspension SUVs"),
					makeLeaf("changan-geely-chinese-crossovers", "Changan CS55, Geely Coolray & Jetour Dashing Modern Tech Crossovers", "ቻንጋን እና ጂሊ ዘመናዊ የቻይና ቴክኖሎጂ መኪኖች", "Changan fi Geely Crossover Ammayyaa", "ናይ Changan CS55, Geely Coolray & Jetour Dashing Modern Tech Crossovers"),
					makeLeaf("ford-explorer-everest-suvs", "Ford Everest & Explorer 7-Seater Turbo Diesel Family Expedition SUVs", "ፎርድ ኤቨረስት ባለ 7 መቀመጫ የቤተሰብ ጂፖች", "Ford Everest SUV Maatii", "ናይ Ford Everest & Explorer 7-Seater Turbo Diesel Family Expedition SUVs"),
				},
			},
			{
				Slug:   "sedans-hatchbacks-compact-cars",
				NameEN: "Sedans, Hatchbacks and City Cars",
				NameAM: "አውቶሞቢል ሴዳኖች፣ ሀችባኮች እና የከተማ መኪኖች",
				NameOM: "Konkolaataa Sasalphaa fi Seedaanii",
				NameTI: "ናይ ከተማ መካይንን ሴዳንን",
				Items: []LeafItemDef{
					makeLeaf("toyota-corolla-compact-sedans", "Toyota Corolla Executive, Executive Plus & Dual VVT-i Reliable Sedans", "ቶዮታ ኮሮላ ጥራት ያላቸው እና አስተማማኝ አውቶሞቢሎች", "Toyota Corolla Seedaanii Filatamaa", "ናይ Toyota Corolla Executive, Executive Plus & Dual VVT-i Reliable Sedans"),
					makeLeaf("toyota-vitz-yaris-hatchbacks", "Toyota Vitz, Yaris & Starlet Fuel Efficient City Daily Hatchbacks", "ቶዮታ ቪትዝ እና ያሪስ ነዳጅ ቆጣቢ አነስተኛ መኪኖች", "Toyota Vitz fi Yaris Qusataa", "ናይ Toyota Vitz, Yaris & Starlet Fuel Efficient City Daily Hatchbacks"),
					makeLeaf("suzuki-dzire-swift-city-cars", "Suzuki Dzire & Swift 1.2L DualJet Low Fuel Urban Runabout Cars", "ሱዙኪ ዲዛየር እና ስዊፍት እጅግ ቆጣቢ መኪኖች", "Suzuki Dzire fi Swift Saffisaa", "ናይ Suzuki Dzire & Swift 1.2L DualJet Low Fuel Urban Runabout Cars"),
					makeLeaf("hyundai-elantra-accent-atos", "Hyundai Elantra, Accent & Atos Affordable Urban Daily Commuters", "ሃዩንዳይ ኤላንትራ፣ አክሰንት እና አቶስ የከተማ መኪኖች", "Hyundai Elantra, Accent fi Atos", "ናይ Hyundai Elantra, Accent & Atos Affordable Urban Daily Commuters"),
					makeLeaf("mercedes-benz-c-class-e-class", "Mercedes-Benz C200, E300 & S-Class Executive Luxury Smooth Sedans", "መርሴዲስ ቤንዝ ሲ-ክላስ እና ኢ-ክላስ የቅንጦት አውቶሞቢሎች", "Mercedes-Benz C-Class fi E-Class", "ናይ Mercedes-Benz C200, E300 & S-Class Executive Luxury Smooth Sedans"),
					makeLeaf("volkswagen-id4-id6-electric-vehicles", "Volkswagen ID.4, ID.6 & BYD Atto 3 100% Electric Modern Vehicles", "ፎልክስዋገን ID.4 እና ቢዋይዲ ሙሉ በሙሉ የኤሌክትሪክ መኪኖች (EV)", "Konkolaataa Elektiriikii Volkswagen ID.4", "ናይ Volkswagen ID.4, ID.6 & BYD Atto 3 100% Electric Modern Vehicles"),
					makeLeaf("toyota-camry-executive-sedans", "Toyota Camry Hybrid 2.5L Dynamic Force Smooth Business Sedans", "ቶዮታ ካምሪ የቢሮ እና የቤተሰብ የቅንጦት መኪኖች", "Toyota Camry Seedaanii Qananii", "ናይ Toyota Camry Hybrid 2.5L Dynamic Force Smooth Business Sedans"),
					makeLeaf("lifan-geely-byd-local-assembled-cars", "Lifan 520, Geely CK & BYD F3 Economical Local Assembled Sedans", "ሊፋን እና ጂሊ በሀገር ውስጥ የተገጣጠሙ ተመጣጣኝ መኪኖች", "Lifan fi Geely Konkolaataa Biyya Keessaa", "ናይ Lifan 520, Geely CK & BYD F3 Economical Local Assembled Sedans"),
				},
			},
		},
	})

	// Domain: motorcycles-bajaj
	list = append(list, DomainCatalogDef{
		RootSlug: "motorcycles-bajaj",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "three-wheelers-bajaj-tuktuk",
				NameEN: "Bajaj, TVS and Cargo Three-Wheelers (Tuk-Tuk)",
				NameAM: "ባጃጆች፣ ቲቪኤስ እና ባለ ሶስት እግር ጭነቶች",
				NameOM: "Bajaajii, TVS fi Konkolaataa Miila Sadii",
				NameTI: "ባጃጅን ቲቪኤስን ባለ ሰለስተ እግሪ መካይንን",
				Items: []LeafItemDef{
					makeLeaf("bajaj-re-compact-4-stroke-passenger", "Bajaj RE Compact 205cc DTS-i 4-Stroke Urban Passenger Three-Wheelers", "ባጃጅ RE ባለ 4-ስትሮክ የከተማ የህዝብ ማመላለሻዎች", "Bajaajii RE 4-Stroke Imaltootaa", "ናይ Bajaj RE Compact 205cc DTS-i 4-Stroke Urban Passenger Three-Wheelers"),
					makeLeaf("tvs-king-deluxe-passenger-tuktuk", "TVS King Deluxe 200cc Electric Start DuraLife Passenger Bajaj Units", "ቲቪኤስ ኪንግ (TVS King) ዘመናዊ የህዝብ ባጃጆች", "TVS King Bajaajii Ammayyaa", "ናይ TVS King Deluxe 200cc Electric Start DuraLife Passenger Bajaj Units"),
					makeLeaf("piaggio-ape-commercial-cargo-three-wheelers", "Piaggio Ape Diesel Heavy Cargo Open Bed Commercial Three-Wheelers", "ፒያጆ የጭነት ባጃጆች ከነሳጥናቸው (Cargo Bajaj)", "Bajaajii Fe'iinsaa Piaggio Ape", "ናይ Piaggio Ape Diesel Heavy Cargo Open Bed Commercial Three-Wheelers"),
					makeLeaf("electric-battery-powered-passenger-bajaj", "Lithium Battery 100km Range Zero-Emission Electric Passenger Tuk-Tuks", "በባትሪ የሚሰሩ ዘመናዊ የኤሌክትሪክ ባጃጆች", "Bajaajii Elektiriikii Baatirii", "ናይ Lithium Battery 100km Range Zero-Emission Electric Passenger Tuk-Tuks"),
					makeLeaf("re-maxima-heavy-cargo-delivery-three-wheelers", "Bajaj Maxima C Heavy Frame 236cc Diesel Farm & Urban Delivery Units", "ባጃጅ ማክሲማ የጭነት ማመላለሻ ጠንካራ ባጃጆች", "Bajaajii Fe'iinsaa Bajaj Maxima C", "ናይ Bajaj Maxima C Heavy Frame 236cc Diesel Farm & Urban Delivery Units"),
					makeLeaf("modified-refrigerated-milk-meat-bajaj", "Insulated Cold Box Body Three-Wheelers for Dairy & Meat Transit", "የወተት እና የስጋ ማመላለሻ ማቀዝቀዣ ያላቸው ባጃጆች", "Bajaajii Foon fi Aanan Qabbaneessu", "ናይ Insulated Cold Box Body Three-Wheelers for Dairy & Meat Transit"),
					makeLeaf("three-wheeler-steel-cabin-roof-canopies", "Heavy Gauge Rain & Sun Weatherproof Bajaj Canvas and Steel Canopies", "የባጃጅ ጣሪያ፣ ሸራዎች እና የዝናብ መከላከያ ከቨሮች", "Shaaraa fi Gubbaa Bajaajii", "ናይ Heavy Gauge Rain & Sun Weatherproof Bajaj Canvas and Steel Canopies"),
					makeLeaf("bajaj-genuine-spare-parts-clutch-engine", "OEM DTS-i Pistons, Piston Rings, Clutch Cables and Brake Drums", "ኦሪጅናል የባጃጅ መለዋወጫ እቃዎች እና ኬብሎች", "Meeshaalee Jijjiirraa Bajaajii Qulqulluu", "ናይ OEM DTS-i Pistons, Piston Rings, Clutch Cables and Brake Drums"),
				},
			},
			{
				Slug:   "motorcycles-scooters-dirt-bikes",
				NameEN: "Motorcycles, Delivery Bikes and Electric Scooters",
				NameAM: "ሞተር ሳይክሎች፣ የዴሊቨሪ ሞተሮች እና ስኩተሮች",
				NameOM: "Mootorsayikilii, Iskuutara fi Meeshaa Geejjibaa",
				NameTI: "ሞተርሳይክልን ናይ ዴሊቨሪ ሞተርን",
				Items: []LeafItemDef{
					makeLeaf("boxer-bm150-heavy-duty-motorcycles", "Bajaj Boxer BM150 Rugged Rural Mud & Hill Transportation Motorcycles", "ባጃጅ ቦክሰር 150 ለገጠር እና ዳገት የሚሆን ጠንካራ ሞተር", "Bajaj Boxer BM150 Mootorsayikilii Jabaa", "ናይ Bajaj Boxer BM150 Rugged Rural Mud & Hill Transportation Motorcycles"),
					makeLeaf("tvs-hlx-125-150-reliable-bikes", "TVS HLX Plus Heavy Load Carrying High-Mileage Commercial Bikes", "ቲቪኤስ ኤችኤልኤክስ (TVS HLX) ጥራት ያላቸው ሞተር ሳይክሎች", "TVS HLX Mootorsayikilii Hojii", "ናይ TVS HLX Plus Heavy Load Carrying High-Mileage Commercial Bikes"),
					makeLeaf("yamaha-dt-125-trail-dirt-bikes", "Yamaha DT125 / AG100 2-Stroke Countryside Agricultural Field Bikes", "ያማሃ ዲቲ (DT125) የገጠር እና የዱር ጉዞ ሞተር ሳይክሎች", "Yamaha DT125 Mootorsayikilii Dirree", "ናይ Yamaha DT125 / AG100 2-Stroke Countryside Agricultural Field Bikes"),
					makeLeaf("express-food-delivery-cargo-motorcycles", "Rear Insulated Heated Food Box Courier Delivery 125cc Motorbikes", "የፈጣን ምግብ ማድረሻ ሳጥን ያላቸው የዴሊቨሪ ሞተሮች", "Mootorsayikilii Nyaata Dhaqqabsiisu (Delivery)", "ናይ Rear Insulated Heated Food Box Courier Delivery 125cc Motorbikes"),
					makeLeaf("foldable-adult-electric-commuter-scooters", "Dual Motor 45km/h Long-Range Lithium Foldable Electric Kick Scooters", "የሚታጠፉ ዘመናዊ የኤሌክትሪክ ስኩተሮች (Electric Scooter)", "Iskuutara Elektiriikii Dachaafamu", "ናይ Dual Motor 45km/h Long-Range Lithium Foldable Electric Kick Scooters"),
					makeLeaf("125cc-automatic-twist-and-go-moped-scooters", "Vespa Style Retro Automatic Transmission City Commuter Scooters", "ቀላል የሚነዱ አውቶማቲክ የስኩተር ሞተሮች (Moped)", "Iskuutara Salphaa Awutoomaatikii", "ናይ Vespa Style Retro Automatic Transmission City Commuter Scooters"),
					makeLeaf("high-power-electric-motorcycles-ebikes", "72V 3000W Fast Acceleration Zero-Fuel Electric Motorbikes", "ነዳጅ የማይፈልጉ ኃይለኛ የኤሌክትሪክ ሞተር ሳይክሎች", "Mootorsayikilii Elektiriikii Boba'aa Malee", "ናይ 72V 3000W Fast Acceleration Zero-Fuel Electric Motorbikes"),
					makeLeaf("dot-ece-certified-motorcycle-helmets", "Full-Face Dual Visor Anti-Fog Certified Motorcycle Safety Helmets", "የሞተር ሳይክል አሽከርካሪ የራስ ቁር ቆቦች (Helmet)", "Qubba Nageenyaa Mootorsayikilii (Helmet)", "ናይ Full-Face Dual Visor Anti-Fog Certified Motorcycle Safety Helmets"),
				},
			},
		},
	})

	// Domain: bicycles-mobility
	list = append(list, DomainCatalogDef{
		RootSlug: "bicycles-mobility",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "mountain-road-commuter-bicycles",
				NameEN: "Mountain Bikes, Road Racers and Commuter Cycles",
				NameAM: "የተራራ፣ የከተማ እና የስፖርት ብስክሌቶች",
				NameOM: "Biskileetii Gaaraa, Daandii fi Ispoortii",
				NameTI: "ናይ ጎቦን ከተማን ስፖርትን ብሽክለታታት",
				Items: []LeafItemDef{
					makeLeaf("aluminum-frame-mountain-bikes-29er", "21-27 Speed Shimano Gear Hydraulic Disc Brake 29-Inch Mountain Bikes", "የሺማኖ ማርሽ ያላቸው 29 ኢንች የተራራ ብስክሌቶች (MTB)", "Biskileetii Gaaraa Shimano 29-Inch", "ናይ 21-27 Speed Shimano Gear Hydraulic Disc Brake 29-Inch Mountain Bikes"),
					makeLeaf("lightweight-carbon-fiber-road-racing-bikes", "700C Drop Bar Lightweight Aerodynamic Road Racing Bicycles", "እጅግ ፈጣን የውድድር እና የአስፋልት ብስክሌቶች (Road Bike)", "Biskileetii Fiigichaa Daandii (Road Bike)", "ናይ 700C Drop Bar Lightweight Aerodynamic Road Racing Bicycles"),
					makeLeaf("hybrid-city-commuter-bicycles", "Upright Comfort Saddle Front Suspension Urban Commuter Bikes", "ለዕለት ተዕለት የከተማ ጉዞ የሚሆኑ ምቹ ብስክሌቶች", "Biskileetii Magaalaa fi Hojii Guyyaa", "ናይ Upright Comfort Saddle Front Suspension Urban Commuter Bikes"),
					makeLeaf("folding-compact-travel-bicycles", "20-Inch Quick-Fold Alloy Trunk-Fit Commuter Pocket Bikes", "በቀላሉ የሚታጠፉ እና መኪና ውስጥ የሚገቡ ብስክሌቶች", "Biskileetii Dachaafamu Konkolaataa Keessa Seenu", "ናይ 20-Inch Quick-Fold Alloy Trunk-Fit Commuter Pocket Bikes"),
					makeLeaf("pedal-assist-electric-mountain-ebikes", "36V/48V Torque Sensor 250W-500W Mid-Drive Electric Bicycles", "በባትሪ እና በፔዳል የሚሰሩ ዘመናዊ የኤሌክትሪክ ብስክሌቶች", "Biskileetii Elektiriikii (E-Bike)", "ናይ 36V/48V Torque Sensor 250W-500W Mid-Drive Electric Bicycles"),
					makeLeaf("heavy-duty-cargo-utility-bicycles", "Reinforced Rear Rack 150kg Payload Farmer & Delivery Bicycles", "እቃ እና ከባድ ጭነት የሚያጓጉዙ ጠንካራ ብስክሌቶች", "Biskileetii Fe'iinsaa fi Ba'aa Jabaa", "ናይ Reinforced Rear Rack 150kg Payload Farmer & Delivery Bicycles"),
					makeLeaf("fixed-gear-single-speed-urban-fixies", "Flip-Flop Hub Minimalist Lightweight City Track Bicycles", "ቀላል እና ጥገና የማይበዙ ነጠላ ማርሽ ብስክሌቶች", "Biskileetii Salphaa Fiiksii (Fixie)", "ናይ Flip-Flop Hub Minimalist Lightweight City Track Bicycles"),
					makeLeaf("cruiser-classic-wide-handlebar-bicycles", "Cushioned Spring Seat Wide Balloon Tire Leisure Bicycles", "የመዝናኛ ሰፊ መሪ እና ምቹ ኮርቻ ያላቸው ክላሲክ ብስክሌቶች", "Biskileetii Qananii fi Bashannanaa", "ናይ Cushioned Spring Seat Wide Balloon Tire Leisure Bicycles"),
				},
			},
			{
				Slug:   "children-bikes-mobility-gear",
				NameEN: "Children's Bicycles, Tricycles and Gear",
				NameAM: "የልጆች ብስክሌቶች፣ ባለ ሶስት እግሮች እና መለዋወጫዎች",
				NameOM: "Biskileetii Daa'immanii fi Meeshaalee",
				NameTI: "ናይ ቆልዑ ብሽክለታትን መለዋወጢ ንዋያትን",
				Items: []LeafItemDef{
					makeLeaf("kids-training-wheel-bicycles-16-inch", "Safety Enclosed Chain Guard Removable Training Wheels Kids Bikes", "ለጀማሪ ልጆች የሚሆኑ ረዳት ጎማ ያላቸው ብስክሌቶች", "Biskileetii Daa'immanii Miila Afurii", "ናይ Safety Enclosed Chain Guard Removable Training Wheels Kids Bikes"),
					makeLeaf("toddler-balance-bikes-no-pedal", "Lightweight Magnesium Alloy Foot-to-Floor Early Balance Bicycles", "ለትንንሽ ልጆች ሚዛን መጠበቂያ ፔዳል አልባ ብስክሌቶች", "Biskileetii Madaallii Daa'immanii", "ናይ Lightweight Magnesium Alloy Foot-to-Floor Early Balance Bicycles"),
					makeLeaf("multi-stage-push-handle-tricycles", "Canopy Sunshade Parental Steering Adjustable Toddler Tricycles", "ወላጆች የሚገፏቸው የፀሐይ መከላከያ ያላቸው የህፃናት ባለ 3 እግር ብስክሌቶች", "Biskileetii Miila Sadii Daa'immanii", "ናይ Canopy Sunshade Parental Steering Adjustable Toddler Tricycles"),
					makeLeaf("kids-safety-helmets-knee-elbow-pads", "Impact Absorbing Multi-Sport CPSC Certified Kids Protective Sets", "የልጆች ብስክሌት መከላከያ ቆቦች፣ የጉልበት እና የክንድ ፓዶች", "Qubba fi Eegduu Jilbaa Daa'immanii", "ናይ Impact Absorbing Multi-Sport CPSC Certified Kids Protective Sets"),
					makeLeaf("bicycle-shimano-derailleur-gear-systems", "7-Speed to 12-Speed Indexed Rear Derailleurs, Shifters & Chains", "የብስክሌት የሺማኖ ማርሾች፣ ቼኖች እና መለዋወጫዎች", "Maarshii fi Saamunaa Biskileetii Shimano", "ናይ 7-Speed to 12-Speed Indexed Rear Derailleurs, Shifters & Chains"),
					makeLeaf("high-pressure-floor-bicycle-tire-pumps", "Universal Presta & Schrader Dual Valve Gauge Floor Foot Pumps", "የብስክሌት ጎማ መንፊያ ትላልቅ የእግር ፓምፖች (ካቴራ)", "Paampii Gommaa Biskileetii", "ናይ Universal Presta & Schrader Dual Valve Gauge Floor Foot Pumps"),
					makeLeaf("puncture-resistant-bicycle-tires-tubes", "26, 27.5 & 29-Inch Kevlar Anti-Flat Mountain Bicycle Tires", "የማይበሱ ጠንካራ የብስክሌት ጎማዎች እና ካናቶች", "Gommaa fi Kaanaa Biskileetii Jabaa", "ናይ 26, 27.5 & 29-Inch Kevlar Anti-Flat Mountain Bicycle Tires"),
					makeLeaf("waterproof-usb-rechargeable-bike-lights", "High-Lumen Front White Headlight & Red Flashing Warning Tail Lights", "በዩኤስቢ ቻርጅ የሚደረጉ የፊት እና የኋላ የብስክሌት መብራቶች", "Ibsaa Biskileetii USB Chaarjamu", "ናይ High-Lumen Front White Headlight & Red Flashing Warning Tail Lights"),
				},
			},
		},
	})

	// Domain: trucks-commercial-transport
	list = append(list, DomainCatalogDef{
		RootSlug: "trucks-commercial-transport",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "heavy-commercial-dump-cargo-trucks",
				NameEN: "Dump Trucks, Cargo Trucks and Prime Movers",
				NameAM: "የጭነት መኪኖች፣ የዳምፕ ቆሻሻ ማመላለሻዎች እና ከባድ ተሳቢዎች",
				NameOM: "Konkolaataa Fe'iinsaa fi Daampii",
				NameTI: "ናይ ጽዕነት መካይንን ዳምፐርን ተሳቢን",
				Items: []LeafItemDef{
					makeLeaf("sinotruk-howo-dump-tipper-trucks", "Sinotruk Howo 371HP 6x4 & 8x4 Heavy Construction Tipper Trucks", "ሲኖትራክ ሆዎ (SinoTruk Howo) የድንጋይ እና አሸዋ መጫኛ ዳምፐሮች", "SinoTruk Howo Daampii Fe'iinsaa", "ናይ Sinotruk Howo 371HP 6x4 & 8x4 Heavy Construction Tipper Trucks"),
					makeLeaf("isuzu-fvr-fsr-heavy-cargo-trucks", "Isuzu FVR, FSR & Forward 10-15 Ton Dry Goods Freight Cargo Trucks", "ኢሱዙ FVR እና FSR ጠንካራ የ10-15 ቶን የጭነት መኪኖች", "Isuzu FVR fi FSR Konkolaataa Fe'iinsaa", "ናይ Isuzu FVR, FSR & Forward 10-15 Ton Dry Goods Freight Cargo Trucks"),
					makeLeaf("isuzu-npr-nqr-medium-cargo-trucks", "Isuzu NPR 4-6 Ton City Cargo Delivery & Distribution Trucks", "ኢሱዙ NPR ከ4 እስከ 6 ቶን የከተማ ውስጥ የጭነት መኪኖች", "Isuzu NPR Fe'iinsa Giddugaleessaa", "ናይ Isuzu NPR 4-6 Ton City Cargo Delivery & Distribution Trucks"),
					makeLeaf("heavy-prime-mover-tractor-trucks-6x4", "Mercedes Actros, Volvo FH & Renault Long-Haul 6x4 Prime Movers", "ቮልቮ፣ መርሴዲስ እና ረኖ የጅቡቲ ኮሪደር ረጅም ተሳቢ ቦቴዎች", "Volvo fi Mercedes Konkolaataa Dheeraa (Trailer)", "ናይ Mercedes Actros, Volvo FH & Renault Long-Haul 6x4 Prime Movers"),
					makeLeaf("heavy-fuel-tanker-trucks-40000l", "Multi-Compartment Steel & Aluminum Petroleum Fuel Tanker Trucks", "ባለ 40,000 ሊትር የነዳጅ እና የናፍጣ ማመላለሻ ቦቴዎች", "Konkolaataa Bootee Boba'aa 40,000L", "ናይ Multi-Compartment Steel & Aluminum Petroleum Fuel Tanker Trucks"),
					makeLeaf("refrigerated-insulated-cold-cargo-trucks", "Thermo King Unit Insulated Meat & Perishable Cargo Box Trucks", "ስጋ እና የወተት ተዋጽኦዎችን የሚያጓጉዙ ማቀዝቀዣ ያላቸው መኪኖች", "Konkolaataa Fe'iinsaa Qabbaneessaa Qabu", "ናይ Thermo King Unit Insulated Meat & Perishable Cargo Box Trucks"),
					makeLeaf("flatbed-container-skeletal-semi-trailers", "40-Foot 3-Axle Container Locking Heavy Transport Skeletal Trailers", "ባለ 40 ጫማ የኮንቴነር መጫኛ ባለ 3 አክሰል ተሳቢዎች (Trailer)", "Tiraayilara Fe'iinsa Qorqorroo (40-Foot)", "ናይ 40-Foot 3-Axle Container Locking Heavy Transport Skeletal Trailers"),
					makeLeaf("faw-shacman-heavy-duty-cargo-trucks", "Shacman F3000 & FAW J6P 380HP Heavy Duty Diesel Haulage Trucks", "ሻክማን እና ፋው ጠንካራ የኮንስትራክሽን እና የጭነት መኪኖች", "Shacman fi FAW Konkolaataa Fe'iinsaa", "ናይ Shacman F3000 & FAW J6P 380HP Heavy Duty Diesel Haulage Trucks"),
				},
			},
			{
				Slug:   "passenger-buses-commercial-vans",
				NameEN: "Minibuses, City Buses and Passenger Coaches",
				NameAM: "የህዝብ ማመላለሻ ሚኒባሶች፣ አውቶቡሶች እና ቫኖች",
				NameOM: "Baasii fi Minibaasii Imaltootaa",
				NameTI: "ናይ ህዝቢ መመላለሲ ሚኒባስን ኣውቶቡስን",
				Items: []LeafItemDef{
					makeLeaf("toyota-hiace-high-roof-minibuses", "Toyota HiAce 15-Seat High Roof Commuter Van Passenger Minibuses", "ቶዮታ ሀይኤስ (Toyota HiAce) ባለ 15 መቀመጫ የህዝብ ሚኒባሶች", "Toyota HiAce Minibaasii Imaltootaa", "ናይ Toyota HiAce 15-Seat High Roof Commuter Van Passenger Minibuses"),
					makeLeaf("isuzu-npr-medium-public-buses", "Isuzu NPR Chassis 30-35 Seat Inter-City Public Transport Buses", "ኢሱዙ NPR ባለ 30-35 ሰው የሀገር አቋራጭ አውቶቡሶች (ታሪካ)", "Isuzu NPR Baasii Imaltootaa 30-Seat", "ናይ Isuzu NPR Chassis 30-35 Seat Inter-City Public Transport Buses"),
					makeLeaf("golden-dragon-higer-luxury-coaches", "50-60 Reclining Seat AC & Video Inter-City Luxury Travel Coaches", "ጎልደን ድራጎን እና ሃይገር ባለ 50-60 ሰው የቅንጦት አውቶቡሶች", "Baasii Qananii Golden Dragon Imalaa", "ናይ 50-60 Reclining Seat AC & Video Inter-City Luxury Travel Coaches"),
					makeLeaf("city-low-floor-mass-transit-buses", "Ankai & Yutong High-Capacity Urban Low-Floor City Passenger Buses", "የከተማ የህዝብ ማመላለሻ ትላልቅ ሰፋፊ አውቶቡሶች (አንበሳ/ሸገር)", "Baasii Magaalaa Bal'aa (Sheger/Anbessa)", "ናይ Ankai & Yutong High-Capacity Urban Low-Floor City Passenger Buses"),
					makeLeaf("hyundai-h1-starex-executive-vans", "Hyundai H-1 Grand Starex 9-12 Seat Tourist & Family Executive Vans", "ሃዩንዳይ ኤች-1 (H-1) የቱሪስት እና የቤተሰብ የቅንጦት ቫኖች", "Hyundai H-1 Vaanii Turiizimii", "ናይ Hyundai H-1 Grand Starex 9-12 Seat Tourist & Family Executive Vans"),
					makeLeaf("ford-transit-commercial-cargo-panel-vans", "High-Roof Long-Wheelbase Enclosed Commercial Logistics Panel Vans", "ፎርድ ትራንዚት ዕቃ መጫኛ የተዘጉ የከተማ ቫኖች (Panel Van)", "Ford Transit Vaanii Fe'iinsaa", "ናይ High-Roof Long-Wheelbase Enclosed Commercial Logistics Panel Vans"),
					makeLeaf("electric-mass-transit-city-buses", "Zero-Emission Fast-Charge City Commuter Modern Electric Buses", "በኤሌክትሪክ የሚሰሩ ዘመናዊ የከተማ አውቶቡሶች", "Baasii Magaalaa Elektiriikii", "ናይ Zero-Emission Fast-Charge City Commuter Modern Electric Buses"),
					makeLeaf("school-bus-safety-equipped-coaches", "Yellow High-Visibility Seat-Belted 30-Seat Student Transport Buses", "የተማሪዎች ማመላለሻ ቢጫ አውቶቡሶች ከነደህንነት ቀበቷቸው", "Baasii Barattootaa Nageenya Qabu", "ናይ Yellow High-Visibility Seat-Belted 30-Seat Student Transport Buses"),
				},
			},
		},
	})

	return list
}
