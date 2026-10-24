package main

import "fmt"

// GetCatalogDomains21To30 returns comprehensive domain catalogues for roots 21 to 30
func GetCatalogDomains21To30() []DomainCatalogDef {
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

	// Domain: professional-services
	list = append(list, DomainCatalogDef{
		RootSlug: "professional-services",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "legal-accounting-translation-services",
				NameEN: "Legal, Accounting, Tax and Translation Services",
				NameAM: "የህግ፣ የሂሳብ፣ የታክስ እና የትርጉም አገልግሎቶች",
				NameOM: "Tajaajila Seeraa, Herregaa fi Turjumaanaa",
				NameTI: "ናይ ሕጊን ሕሳብን ትርጉምን ኣገልግሎታት",
				Items: []LeafItemDef{
					makeLeaf("business-incorporation-legal-services", "Company Formation, Trade License & Investment Legal Services", "የድርጅት ምስረታ፣ የንግድ ፈቃድ እና የኢንቨስትመንት ህግ አገልግሎቶች", "Tajaajila Seeraa Dhaabbata Hundeessuu", "ናይ Company Formation, Trade License & Investment Legal Services"),
					makeLeaf("tax-filing-vat-declaration-advisory", "Monthly VAT, Withholding Tax Declaration and Ministry Audit Prep", "የታክስ፣ የተጨማሪ እሴት ታክስ (VAT) እና የኦዲት ዝግጅት አገልግሎቶች", "Tajaajila Gibiraa fi VAT", "ናይ Monthly VAT, Withholding Tax Declaration and Ministry Audit Prep"),
					makeLeaf("annual-financial-audit-certification", "Certified Annual Financial Audit & External Reporting Services", "ዓመታዊ የሂሳብ ኦዲት እና ሪፖርት አገልግሎቶች", "Tajaajila Ooditii Waggaa", "ናይ Certified Annual Financial Audit & External Reporting Services"),
					makeLeaf("court-advocacy-litigation-services", "Civil & Commercial Court Litigation, Contract Drafting & Notary", "የፍርድ ቤት ክርክር፣ የውል ዝግጅት እና የውክልና አገልግሎቶች", "Tajaajila Mana Murtii fi Qophii Waliigaltee", "ናይ Civil & Commercial Court Litigation, Contract Drafting & Notary"),
					makeLeaf("certified-official-document-translation", "Embassy & Ministry Certified Document Translation & Notarization", "ለኤምባሲ እና ለመንግስት መስሪያ ቤት ሰነዶች የትርጉም አገልግሎት", "Tajaajila Turjumaanaa Waraqaa Ragaa", "ናይ Embassy & Ministry Certified Document Translation & Notarization"),
					makeLeaf("patent-trademark-ip-registration", "Trademark Registration, Patent Filing and IP Protection Advisory", "የንግድ ምልክት (ትሬድማርክ) እና የፓተንት ምዝገባ አገልግሎቶች", "Galmee Mallattoo Daldalaa (Trademark)", "ናይ Trademark Registration, Patent Filing and IP Protection Advisory"),
					makeLeaf("hr-recruitment-payroll-outsourcing", "Corporate Headhunting, Staff Placement & Monthly Payroll Outsourcing", "የሰራተኛ ቅጥር እና የደመወዝ ክፍያ ማኔጅመንት አገልግሎቶች", "Tajaajila Qacarrii fi Mindaa Hojjattootaa", "ናይ Corporate Headhunting, Staff Placement & Monthly Payroll Outsourcing"),
					makeLeaf("property-valuation-asset-appraisal", "Certified Fixed Asset, Real Estate and Collateral Valuation Services", "የቋሚ ንብረት፣ የቤት እና የዋስትና ግምት አገልግሎት", "Tajaajila Tilmaama Qabeenyaa fi Gamoo", "ናይ Certified Fixed Asset, Real Estate and Collateral Valuation Services"),
				},
			},
			{
				Slug:   "facility-cleaning-security-pest-control",
				NameEN: "Cleaning, Security, Pest Control and Moving Services",
				NameAM: "የጽዳት፣ የጥበቃ፣ የተባይ ማጥፋት እና የእቃ ማጓጓዝ አገልግሎቶች",
				NameOM: "Tajaajila Qulqullinaa, Eegumsaa fi Geejjibaa",
				NameTI: "ናይ ጽሬትን ሓለዋን ኣገልግሎታት",
				Items: []LeafItemDef{
					makeLeaf("commercial-office-deep-cleaning", "Post-Construction & Regular Office Deep Carpet & Tile Cleaning", "የቢሮ እና የህንፃ ጥልቅ የጽዳት አገልግሎቶች", "Tajaajila Qulqullina Biiroo fi Gamoo", "ናይ Post-Construction & Regular Office Deep Carpet & Tile Cleaning"),
					makeLeaf("corporate-security-guard-services", "Trained Uniformed Security Guards and Patrol Protection", "የሰለጠኑ የጥበቃ ሰራተኞች እና የደህንነት አገልግሎቶች", "Tajaajila Eegumsa fi Nageenyaa", "ናይ Trained Uniformed Security Guards and Patrol Protection"),
					makeLeaf("fumigation-pest-termite-control", "Industrial & Residential Insect, Cockroach and Termite Fumigation", "የበረሮ፣ የተባይ እና የነቀዝ ማጥፋት ኬሚካል ርጭት አገልግሎቶች", "Tajaajila Ilbiisota Balleessuu", "ናይ Industrial & Residential Insect, Cockroach and Termite Fumigation"),
					makeLeaf("professional-packers-movers-household", "Household Furniture Packing, Loading and Inter-City Moving Services", "የቤት እና የቢሮ እቃዎችን በጥንቃቄ የማሸግ እና የማጓጓዝ አገልግሎት", "Tajaajila Qodaa Manaa Geejjibuu", "ናይ Household Furniture Packing, Loading and Inter-City Moving Services"),
					makeLeaf("deep-groundwater-borehole-drilling", "Deep Water Well Drilling, Pump Lowering & Water Quality Testing", "የጥልቅ የከርሰ ምድር ውሃ ቁፋሮ እና የፓምፕ ተከላ አገልግሎት", "Tajaajila Boolla Bishaan Qotuu", "ናይ Deep Water Well Drilling, Pump Lowering & Water Quality Testing"),
					makeLeaf("customs-clearing-forwarding-transit", "Djibouti Corridor Import/Export Customs Clearance and Cargo Transit", "የጅቡቲ ኮሪደር የጉምሩክ ክሊራንስ እና የትራንዚት ወኪል አገልግሎቶች", "Tajaajila Gumurukii fi Dambaliqaa", "ናይ Djibouti Corridor Import/Export Customs Clearance and Cargo Transit"),
					makeLeaf("commercial-laundry-dry-cleaning-pickup", "Hotel Bedding, Curtain and Suit Dry Cleaning Pickup & Delivery", "የሆቴል አንሶላዎች፣ መጋረጃዎች እና የሱፍ እጥበት አገልግሎት", "Tajaajila Uffata Miicuu fi Diraykiliinii", "ናይ Hotel Bedding, Curtain and Suit Dry Cleaning Pickup & Delivery"),
					makeLeaf("electric-generator-maintenance-services", "Heavy Diesel Generator Overhaul, Filter Replacement & Rewinding", "የጄኔሬተር ጥገና እና ወቅታዊ የሰርቪስ አገልግሎቶች", "Tajaajila Suphaa Jenereetaraa", "ናይ Heavy Diesel Generator Overhaul, Filter Replacement & Rewinding"),
				},
			},
		},
	})

	// Domain: software-it-services
	list = append(list, DomainCatalogDef{
		RootSlug: "software-it-services",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "enterprise-software-erp-pos",
				NameEN: "ERP Systems, Accounting and POS Software",
				NameAM: "የኢአርፒ፣ የሂሳብ እና የሽያጭ ነጥብ (POS) ሶፍትዌሮች",
				NameOM: "Sooftiweerii ERP, POS fi Herregaa",
				NameTI: "ናይ ERPን POSን ሶፍትዌራት",
				Items: []LeafItemDef{
					makeLeaf("cloud-erp-business-management-software", "Inventory, Human Resources & Accounting Unified Cloud ERP Software", "የተቀናጀ የድርጅት ማኔጅመንት እና የሂሳብ ኢአርፒ ሶፍትዌር", "Sooftiweerii ERP Bulchiinsa Daldalaa", "ናይ Inventory, Human Resources & Accounting Unified Cloud ERP Software"),
					makeLeaf("supermarket-restaurant-pos-systems", "Touchscreen CBE-Birr & Telebirr Integrated Retail POS Software", "የሱፐርማርኬት እና ሬስቶራንት የክፍያ መቆጣጠሪያ POS ሶፍትዌር", "Sooftiweerii POS Suuqii fi Nyaataa", "ናይ Touchscreen CBE-Birr & Telebirr Integrated Retail POS Software"),
					makeLeaf("school-management-information-systems", "Student Attendance, Grading, Online Fee Payment & Portal Software", "የትምህርት ቤቶች የተማሪ ውጤት እና ክፍያ መቆጣጠሪያ ሶፍትዌር", "Sooftiweerii Bulchiinsa Mana Barumsaa", "ናይ Student Attendance, Grading, Online Fee Payment & Portal Software"),
					makeLeaf("hospital-clinic-patient-emr-software", "Electronic Medical Records, Pharmacy Stock & Inpatient Billing EMR", "የሆስፒታል እና የክሊኒክ የታካሚዎች መረጃ መቆጣጠሪያ ሶፍትዌር", "Sooftiweerii Hospitaalaa fi EMR", "ናይ Electronic Medical Records, Pharmacy Stock & Inpatient Billing EMR"),
					makeLeaf("hotel-property-management-systems-pms", "Room Booking, Front Desk Check-in & Housekeeping Hotel Software", "የሆቴል አልጋ እና የደንበኞች አገልግሎት መቆጣጠሪያ ሶፍትዌር", "Sooftiweerii Hoteelaa (PMS)", "ናይ Room Booking, Front Desk Check-in & Housekeeping Hotel Software"),
					makeLeaf("pharmacy-stock-expiry-management-software", "Batch Number Expiry Alert & Drug Sales Inventory Software", "የመድኃኒት መደብሮች የክምችት እና የመጠቀሚያ ጊዜ መቆጣጠሪያ ሶፍትዌር", "Sooftiweerii Qorichaa fi Faarmaasii", "ናይ Batch Number Expiry Alert & Drug Sales Inventory Software"),
					makeLeaf("microfinance-lending-credit-software", "Loan Disbursal, Interest Tracking & Collateral Registry Microfinance Systems", "የብድር እና ቁጠባ ተቋማት የሂሳብ መቆጣጠሪያ ሶፍትዌር", "Sooftiweerii Liqii fi Qusannaa", "ናይ Loan Disbursal, Interest Tracking & Collateral Registry Microfinance Systems"),
					makeLeaf("gps-fleet-tracking-management-software", "Live Fuel Sensor, Route Replay & Fleet Maintenance Dispatch Platforms", "የመኪናዎችን ጉዞ እና ነዳጅ በካርታ ላይ መከታተያ የጂፒኤስ ሶፍትዌር", "Sooftiweerii Hordoffii Konkolaataa (GPS)", "ናይ Live Fuel Sensor, Route Replay & Fleet Maintenance Dispatch Platforms"),
				},
			},
			{
				Slug:   "web-mobile-custom-development",
				NameEN: "Web Development, Mobile Apps and IT Infrastructure",
				NameAM: "የዌብሳይት፣ የሞባይል አፕሊኬሽን እና የኔትወርክ ዝርጋታ",
				NameOM: "Ijaarsa Weebsaayitii, Appii fi Netwoorkii",
				NameTI: "ናይ ዌብሳይትን ሞባይል ኣፕን ስራሕቲ",
				Items: []LeafItemDef{
					makeLeaf("custom-corporate-web-portal-development", "Next.js & Tailored Responsive Corporate Multilingual Web Portals", "የድርጅቶች ዘመናዊ ዌብሳይት እና የመረጃ ፖርታል ግንባታ", "Misooma Weebsaayitii Dhaabbilee", "ናይ Next.js & Tailored Responsive Corporate Multilingual Web Portals"),
					makeLeaf("fintech-mobile-wallet-payment-development", "Telebirr & Chapa API Integrated Secure iOS & Android Payment Apps", "የሞባይል ክፍያ እና የገንዘብ ማስተላለፊያ አፕሊኬሽኖች ስራ", "Hojii Appii Kaffaltii fi Baankii", "ናይ Telebirr & Chapa API Integrated Secure iOS & Android Payment Apps"),
					makeLeaf("enterprise-lan-fiber-cabling-installation", "Structured Copper & Fiber Network Installation, Rack Dressing & Testing", "የቢሮ እና የህንፃ ኔትወርክ እና ፋይበር ኬብል ዝርጋታ ስራዎች", "Hojii Diriirsa Netwoorkii fi Faayibaraa", "ናይ Structured Copper & Fiber Network Installation, Rack Dressing & Testing"),
					makeLeaf("cctv-ip-camera-installation-services", "High-Definition NVR IP Surveillance Camera Network Installation", "የሲሲቲቪ የደህንነት ካሜራዎች ገጠማ እና የኔትወርክ ስራ", "Hojii Dhaabbii Kaameeraa CCTV", "ናይ High-Definition NVR IP Surveillance Camera Network Installation"),
					makeLeaf("biometric-time-attendance-access-control", "Fingerprint & Face Recognition Door Access Control Systems", "የጣት አሻራ እና የፊት ገጽታ የሰራተኞች መግቢያ መቆጣጠሪያዎች", "Hojii To'annoo Seensaa Asheeraa", "ናይ Fingerprint & Face Recognition Door Access Control Systems"),
					makeLeaf("cybersecurity-vulnerability-penetration-testing", "Ethical Hacking Web, Mobile & Server Vulnerability Assessment Audits", "የሶፍትዌር እና የሰርቨር ደህንነት ፍተሻ (Penetration Testing)", "Qorannoo Nageenya Saayibarii", "ናይ Ethical Hacking Web, Mobile & Server Vulnerability Assessment Audits"),
					makeLeaf("linux-cloud-server-migration-maintenance", "Ubuntu/Debian Dedicated Server Setup, Docker & Automated Daily Backups", "የክላውድ ሰርቨር ዝግጅት እና የቴክኒክ ሰርቪስ ስራዎች", "Bulchiinsa Sarvarii fi Kilaawudii", "ናይ Ubuntu/Debian Dedicated Server Setup, Docker & Automated Daily Backups"),
					makeLeaf("email-hosting-domain-name-registration", "Corporate .et & .com Domain Names with Google Workspace/Office 365", "የድርጅት ኢሜይል እና የዌብሳይት ስም ምዝገባ (.et / .com)", "Galmee Doomeenii fi Imeelii Dhaabbataa", "ናይ Corporate .et & .com Domain Names with Google Workspace/Office 365"),
				},
			},
		},
	})

	// Domain: business-services
	list = append(list, DomainCatalogDef{
		RootSlug: "business-services",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "investment-feasibility-tax-consulting",
				NameEN: "Investment Feasibility, EIA and Management Consulting",
				NameAM: "የኢንቨስትመንት አዋጭነት ጥናት፣ የአካባቢ ተፅዕኖ እና የንግድ አማካሪነት",
				NameOM: "Qorannoo Daldalaa, Invastimantii fi Gorsaa",
				NameTI: "ናይ ኢንቨስትመንትን ኣዋጭነትን ምኽሪ",
				Items: []LeafItemDef{
					makeLeaf("bankable-investment-feasibility-studies", "Project Profiles & Bankable Feasibility Reports for Loan Approval", "ለባንክ ብድር የሚቀርቡ የኢንቨስትመንት አዋጭነት ጥናቶች (Feasibility)", "Qorannoo Facaatii Invastimantii (Feasibility)", "ናይ Project Profiles & Bankable Feasibility Reports for Loan Approval"),
					makeLeaf("environmental-social-impact-assessments", "Certified Environmental & Social Impact Assessment (ESIA) Reports", "የአካባቢ እና የማህበረሰብ ተፅዕኖ ግምገማ ጥናቶች (EIA)", "Qorannoo Dhiibbaa Naannoo (ESIA)", "ናይ Certified Environmental & Social Impact Assessment (ESIA) Reports"),
					makeLeaf("corporate-strategic-business-plans", "5-Year Organizational Growth Strategy & Financial Projection Models", "የአምስት ዓመት የድርጅት ስትራቴጂክ እቅድ ዝግጅት", "Karoora Tarsiimoo Daldalaa Waggaa 5", "ናይ 5-Year Organizational Growth Strategy & Financial Projection Models"),
					makeLeaf("iso-9001-22000-quality-certification", "Quality Management Systems (QMS) & HACCP Food Safety Audits", "የአይኤስኦ (ISO 9001/22000) የጥራት ማረጋገጫ ዝግጅት እና ኦዲት", "Waraqaa Ragaa Qulqullinaa ISO", "ናይ Quality Management Systems (QMS) & HACCP Food Safety Audits"),
					makeLeaf("customs-tariff-import-export-advisory", "Ethiopian Customs Commission Harmonized System Tariff Consulting", "የጉምሩክ ታሪፍ እና የወጪ ንግድ ህጋዊ የምክር አገልግሎት", "Gorsa Taarifa Gumurukii fi Daldala Alaa", "ናይ Ethiopian Customs Commission Harmonized System Tariff Consulting"),
					makeLeaf("company-valuation-mergers-acquisitions", "Due Diligence, Business Asset Valuation & Equity Investment Advisory", "የድርጅቶች ዋጋ ግምት እና የባለቤትነት ውህደት ምክር", "Tilmaama Gatii fi Walitti Dabalamuu Dhaabbilee", "ናይ Due Diligence, Business Asset Valuation & Equity Investment Advisory"),
					makeLeaf("ngo-project-monitoring-evaluation-m-and-e", "Baseline Surveys, Mid-Term Reviews & Final Grant Project Evaluations", "የልማት ድርጅቶች (NGO) የፕሮጀክት ክትትል እና ግምገማ (M&E)", "Hordoffii fi Gamaaggama Pirojektii (M&E)", "ናይ Baseline Surveys, Mid-Term Reviews & Final Grant Project Evaluations"),
					makeLeaf("market-research-consumer-sentiment-surveys", "Addis Ababa & Regional Field Focus Groups & Brand Awareness Surveys", "የገበያ ፍላጎት እና የሸማቾች አስተያየት የጥናት ዳሰሳዎች", "Qorannoo Gabaa fi Fedhii Maamiltootaa", "ናይ Addis Ababa & Regional Field Focus Groups & Brand Awareness Surveys"),
				},
			},
			{
				Slug:   "advertising-printing-media-production",
				NameEN: "Commercial Printing, Billboard and Media Advertising",
				NameAM: "የህትመት፣ የቢልቦርድ እና የሚዲያ ማስታወቂያ ስራዎች",
				NameOM: "Maxxansaa, Beeksisootaa fi Miidiyaa",
				NameTI: "ናይ ማሕተምን ምልክታን ሚድያን",
				Items: []LeafItemDef{
					makeLeaf("outdoor-large-format-billboard-printing", "Flex Banner, Vinyl Sticker & Backlit Billboard Graphic Printing", "የትላልቅ የውጪ ማስታወቂያ ባነሮች እና ስቲከሮች ህትመት", "Maxxansa Baanarii fi Beeksisootaa Bal'aa", "ናይ Flex Banner, Vinyl Sticker & Backlit Billboard Graphic Printing"),
					makeLeaf("lighted-3d-acrylic-storefront-signage", "Channel Letter 3D LED Illuminated Metal Commercial Storefront Signs", "የሱቅ እና የድርጅት 3D በኤልኢዲ የሚያበሩ ታፔላዎች (Signboard)", "Mallattoo Daldalaa 3D Ibsaa Qabu", "ናይ Channel Letter 3D LED Illuminated Metal Commercial Storefront Signs"),
					makeLeaf("corporate-branding-promotional-giveaways", "Branded T-Shirts, Umbrellas, Executive Pens, Notebooks & Mugs", "የድርጅት መለያ ስም የታተመባቸው ቲሸርቶች፣ ጃንጥላዎች እና ስጦታዎች", "Uffata fi Meeshaalee Mallattoo Dhaabbataa", "ናይ Branded T-Shirts, Umbrellas, Executive Pens, Notebooks & Mugs"),
					makeLeaf("professional-tv-radio-commercial-production", "Broadcast Quality Scriptwriting, 4K Filming & Radio Spot Mixing", "የቴሌቪዥን እና የሬዲዮ ማስታወቂያ ቪዲዮ ዝግጅት", "Hojii Viidiyoo fi Beeksisoota Raadiyoo", "ናይ Broadcast Quality Scriptwriting, 4K Filming & Radio Spot Mixing"),
					makeLeaf("digital-led-screen-truck-outdoor-advertising", "Mobile LED Truck Screen Rentals & Prime Highway Display Placements", "በተንቀሳቃሽ መኪና ኤልኢዲ ስክሪን የሚተላለፉ ማስታወቂያዎች", "Beeksisa Iskiriinii LED Konkolaatu", "ናይ Mobile LED Truck Screen Rentals & Prime Highway Display Placements"),
					makeLeaf("brochure-annual-report-book-printing", "Full-Color Gloss Laminated Annual Company Reports & Catalogs", "ዓመታዊ የድርጅት ሪፖርቶች፣ ብሮሸሮች እና መጽሔቶች ህትመት", "Maxxansa Gabaasa Waggaa fi Kitaabaa", "ናይ Full-Color Gloss Laminated Annual Company Reports & Catalogs"),
					makeLeaf("event-exhibition-booth-construction", "Trade Fair Custom Octanorm & Wooden Exhibition Pavilion Booths", "የኤግዚቢሽን እና የንግድ ትርኢት ድንኳኖች እና ቡቶች ግንባታ", "Ijaarsa Buufata Agarsiisa Daldalaa (Booth)", "ናይ Trade Fair Custom Octanorm & Wooden Exhibition Pavilion Booths"),
					makeLeaf("social-media-influencer-marketing-campaigns", "High-Reach TikTok & Telegram Creator Brand Endorsement Management", "የቲክቶክ እና ቴሌግራም ታዋቂ ሰዎች (Influencer) ማስታወቂያ ስራ", "Beeksisa Daldalaa Miidiyaa Hawaasaatiin", "ናይ High-Reach TikTok & Telegram Creator Brand Endorsement Management"),
				},
			},
		},
	})

	// Domain: office-supplies
	list = append(list, DomainCatalogDef{
		RootSlug: "office-supplies",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "paper-stationery-filing-supplies",
				NameEN: "Paper, Stationery and Document Filing Supplies",
				NameAM: "ወረቀቶች፣ እስክሪብቶዎች እና የፋይል ማስቀመጫዎች",
				NameOM: "Waraqaa, Qalama fi Qodaa Faayilaa",
				NameTI: "ወረቐትን ጽሕፈትን ናይ ሰነድ መዝገብን",
				Items: []LeafItemDef{
					makeLeaf("white-a4-copy-paper-reams", "80gsm High-Brightness Jam-Free 500-Sheet White A4 Copy Paper Reams", "80 ግራም ንጹህ ነጭ የA4 ፎቶኮፒ ወረቀቶች (በካርቶን)", "Waraqaa A4 Qulqulluu (Ream)", "ናይ 80gsm High-Brightness Jam-Free 500-Sheet White A4 Copy Paper Reams"),
					makeLeaf("a3-paper-rolls-architectural-tracing", "A3 White Laser Paper & 80g Translucent Architectural Tracing Rolls", "የA3 ፕሪንተር ወረቀቶች እና የፕላን ማተሚያ ወረቀቶች", "Waraqaa A3 fi Pilaanii", "ናይ A3 White Laser Paper & 80g Translucent Architectural Tracing Rolls"),
					makeLeaf("box-lever-arch-document-files", "Heavy Gauge 75mm Metal Bound Cardboard Lever Arch Document Folders", "ሰነድ ማስቀመጫ ጠንካራ የቦክስ ፋይሎች (Box File)", "Faayila Saanduqa Waraqaa (Box File)", "ናይ Heavy Gauge 75mm Metal Bound Cardboard Lever Arch Document Folders"),
					makeLeaf("smooth-ballpoint-pens-box-packs", "0.7mm Blue, Black & Red Smooth-Writing Office Ballpoint Pens (Box of 50)", "ሰማያዊ እና ጥቁር እስክሪብቶዎች በካርቶን (የ50 ፍሬ)", "Qalama Qulqulluu Biiroo (Kartonii)", "ናይ 0.7mm Blue, Black & Red Smooth-Writing Office Ballpoint Pens (Box of 50)"),
					makeLeaf("heavy-duty-desktop-staplers-perforators", "100-Sheet Capacity Steel Binding Staplers & 2-Hole Hole Punchers", "ትላልቅ ሰነዶችን የሚመቱ ስቴፕለሮች እና መበሻዎች (Puncher)", "Isteepilara fi Qaawwisa Waraqaa Jabaa", "ናይ 100-Sheet Capacity Steel Binding Staplers & 2-Hole Hole Punchers"),
					makeLeaf("transparent-brown-packing-tapes", "2-Inch Heavy Adhesion OPP Packaging & Box Sealing Tapes", "ካርቶን ማሸጊያ ትላልቅ ቴፖች (ቡናማ እና ግልጽ)", "Teepii Saanduqa Hidhu Jabaa", "ናይ 2-Inch Heavy Adhesion OPP Packaging & Box Sealing Tapes"),
					makeLeaf("self-adhesive-sticky-notes-memo-pads", "Neon Yellow 3x3 Inch Removable Adhesive Office Sticky Notes", "የማስታወሻ የሚለጠፉ ትናንሽ ወረቀቶች (Sticky Notes)", "Waraqaa Qabatu Yaadannoo (Sticky Notes)", "ናይ Neon Yellow 3x3 Inch Removable Adhesive Office Sticky Notes"),
					makeLeaf("plastic-id-card-laminating-pouches", "A4 & ID Size 125 Micron High-Gloss Thermal Laminating Pouches", "ሰነድ እና መታወቂያ የሚያትሙ የላሚኔሽን ፕላስቲኮች", "Pilaastikii Waraqaa Eenyummaa (Lamination)", "ናይ A4 & ID Size 125 Micron High-Gloss Thermal Laminating Pouches"),
				},
			},
			{
				Slug:   "commercial-money-counters-shredders",
				NameEN: "Cash Counters, Safe Boxes and Paper Shredders",
				NameAM: "የብር መቁጠሪያ ማሽኖች፣ ካዝናዎች እና ወረቀት መፍጫዎች",
				NameOM: "Maashinii Qarshii fi Saanduqa Sibilaa",
				NameTI: "ናይ ቅርሺ መቑጸርን ካዝናን",
				Items: []LeafItemDef{
					makeLeaf("multi-currency-fake-note-detector-counters", "UV, MG & IR Counterfeit Bill Detecting High-Speed Money Counters", "ሀሰተኛ ብር የሚለዩ ፈጣን የብር መቁጠሪያ ማሽኖች", "Maashinii Qarshii Lakkaa'u fi Sobaa Barbaadu", "ናይ UV, MG & IR Counterfeit Bill Detecting High-Speed Money Counters"),
					makeLeaf("heavy-duty-fireproof-digital-safes", "Biometric Fingerprint & Keypad Fire-Resistant Steel Valuables Safes", "በጣት አሻራ እና ኮድ የሚከፈቱ የእሳት አደጋ መከላከያ ጠንካራ ካዝናዎች", "Saanduqa Sibilaa Qabeenyaa (Safe Box)", "ናይ Biometric Fingerprint & Keypad Fire-Resistant Steel Valuables Safes"),
					makeLeaf("cross-cut-confidential-paper-shredders", "P-4 Security Level Continuous Duty Credit Card & Paper Shredders", "ሚስጥራዊ ሰነዶችን ቆራርጠው የሚያጠፉ የወረቀት መፍጫ ማሽኖች", "Maashinii Waraqaa Mummuru (Shredder)", "ናይ P-4 Security Level Continuous Duty Credit Card & Paper Shredders"),
					makeLeaf("thermal-roll-paper-supermarket-pos", "80x80mm BPA-Free High-Contrast POS Cash Register Receipt Rolls", "የደረሰኝ እና የካሽ ሬጅስተር ቴርማል ወረቀቶች በጥቅል", "Waraqaa Nagahee Kaffaltii (Thermal Roll)", "ናይ 80x80mm BPA-Free High-Contrast POS Cash Register Receipt Rolls"),
					makeLeaf("electronic-barcode-handheld-scanners", "2D QR Code & 1D Barcode USB / Bluetooth Supermarket Scanners", "የእቃ መሸጫ ባርኮድ ማንበቢያ መሳሪያዎች (Barcode Scanner)", "Iskaanara Baarkoodii Suuqii", "ናይ 2D QR Code & 1D Barcode USB / Bluetooth Supermarket Scanners"),
					makeLeaf("price-tag-labeling-guns-stickers", "Single Line Price Tag Marker Label Guns with Refill Roll Stickers", "ዋጋ መለጠፊያ ማሽኖች እና ስቲከሮች (Price Tag Gun)", "Meeshaa Istiikara Gatii Qabsiisu", "ናይ Single Line Price Tag Marker Label Guns with Refill Roll Stickers"),
					makeLeaf("heavy-duty-plastic-laminator-machines", "A3 Hot & Cold 4-Roller Photo Document Laminating Machines", "የA3 እና A4 ሰነድ ማተሚያ ላሚኔተር ማሽኖች", "Maashinii Waraqaa Lamiineetii Gochu", "ናይ A3 Hot & Cold 4-Roller Photo Document Laminating Machines"),
					makeLeaf("manual-coil-spiral-bookbinding-machines", "A4 Plastic Comb & Wire Spiral Document Punch Binding Machines", "የመጽሐፍ እና ሰነድ ማሰሪያ ስፓይራል ማሽኖች (Bookbinding)", "Maashinii Kitaaba Hidhu (Binding)", "ናይ A4 Plastic Comb & Wire Spiral Document Punch Binding Machines"),
				},
			},
		},
	})

	// Domain: retail-wholesale
	list = append(list, DomainCatalogDef{
		RootSlug: "retail-wholesale",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "fmcg-wholesale-packaged-goods",
				NameEN: "FMCG Wholesale, Edible Oils and Bulk Goods",
				NameAM: "የፈጣን ሸማች እቃዎች የጅምላ ንግድ፣ የምግብ ዘይት እና የታሸጉ እቃዎች",
				NameOM: "Daldala Jumlaa Nyaataa fi Zayitaa",
				NameTI: "ናይ ጅምላ ምግቢን ዘይትን",
				Items: []LeafItemDef{
					makeLeaf("refined-sunflower-palm-cooking-oils", "5-Liter & 20-Liter Factory Sealed Sunflower & Palm Cooking Oil Jerrycans", "ባለ 5 ሊትር እና 20 ሊትር የሱፍ እና የዘንባባ የምግብ ዘይቶች", "Zayita Nyaataa Litira 5 fi 20", "ናይ 5-Liter & 20-Liter Factory Sealed Sunflower & Palm Cooking Oil Jerrycans"),
					makeLeaf("bulk-white-pure-sugar-bags-50kg", "50kg Food-Grade Standard Refined White Cane Sugar Bags", "የ50 ኪሎ ነጭ የሸንኮራ አገዳ ስኳር በጅምላ", "Sukkaara Qulqulluu Kiiloo 50", "ናይ 50kg Food-Grade Standard Refined White Cane Sugar Bags"),
					makeLeaf("premium-wheat-flour-bakeries-50kg", "High-Gluten White Wheat Bakery Baking Flour 50kg Sacks", "የ50 ኪሎ የዳቦ እና ኬክ የፍርኖ ዱቄት በጅምላ", "Dawaa Qamadii Daabboo Kiiloo 50", "ናይ High-Gluten White Wheat Bakery Baking Flour 50kg Sacks"),
					makeLeaf("pasta-macaroni-factory-cartons", "Wholesale 500g Packs Spaghetti Pasta & Macaroni Full Cartons", "የፓስታ እና ማካሮኒ የፋብሪካ ካርቶኖች በጅምላ", "Paastaa fi Makkaroniin Kartooniidhaan", "ናይ Wholesale 500g Packs Spaghetti Pasta & Macaroni Full Cartons"),
					makeLeaf("packaged-mineral-water-bottle-packs", "Shrink-Wrapped Packs 0.5L, 1L & 2L Natural Bottled Spring Waters", "የታሸገ የተፈጥሮ ምንጭ የጠረጴዛ ውሃ በጅምላ (በፓክ)", "Bishaan Qulqulluu Saamunaan", "ናይ Shrink-Wrapped Packs 0.5L, 1L & 2L Natural Bottled Spring Waters"),
					makeLeaf("laundry-washing-powder-detergent-sacks", "10kg-25kg Enzyme Active Deep-Clean Industrial Washing Powders", "የልብስ ማጠቢያ የዱቄት ሳሙናዎች በጆንያ እና በካርቶን", "Saamunaa Daakuu Uffataa Gurguddaa", "ናይ 10kg-25kg Enzyme Active Deep-Clean Industrial Washing Powders"),
					makeLeaf("toilet-bath-soap-bar-cartons", "Antiseptic & Moisturizing Scented Bath Soap Bars (Carton of 72)", "የገላ ሳሙናዎች በካርቶን (የ72 ፍሬ)", "Saamunaa Qaamaa Kartoonii", "ናይ Antiseptic & Moisturizing Scented Bath Soap Bars (Carton of 72)"),
					makeLeaf("long-grain-basmati-rice-bags-25kg", "Aged Aromatic Super Kernel White Basmati Rice 25kg Sacks", "የ25 ኪሎ የባስማቲ ሩዝ በጅምላ", "Ruuzii Basmatii Kiiloo 25", "ናይ Aged Aromatic Super Kernel White Basmati Rice 25kg Sacks"),
				},
			},
			{
				Slug:   "retail-store-display-supermarket-fittings",
				NameEN: "Supermarket Shelves, Carts and Display Racks",
				NameAM: "የሱፐርማርኬት ሼልፎች፣ ጋሪዎች እና የማሳያ መደርደሪያዎች",
				NameOM: "Shalfii fi Gaarii Suuqii fi Supermarket",
				NameTI: "ናይ ሱፐርማርኬት ሼልፍን ጋሪን",
				Items: []LeafItemDef{
					makeLeaf("double-sided-supermarket-gondola-shelves", "Heavy-Duty Powder Coated Steel Double-Sided Island Shelving Racks", "የሱፐርማርኬት የመሃል ባለሁለት ገጽታ የብረት ሼልፎች", "Shalfii Sibilaa Supermarket", "ናይ Heavy-Duty Powder Coated Steel Double-Sided Island Shelving Racks"),
					makeLeaf("wall-mount-single-sided-store-shelves", "Perforated Pegboard Hook Compatible Commercial Wall Grocery Racks", "የግድግዳ ሱቅ መደርደሪያዎች ከነመንጠልጠያቸው", "Shalfii Girgidaa Suuqii", "ናይ Perforated Pegboard Hook Compatible Commercial Wall Grocery Racks"),
					makeLeaf("chrome-supermarket-shopping-trolleys", "100L-150L Swivel Bearing Wheels Metal Supermarket Grocery Carts", "የሱፐርማርኬት የሸማቾች የብረት ጋሪዎች (Shopping Cart)", "Gaarii Suuqii (Shopping Trolley)", "ናይ 100L-150L Swivel Bearing Wheels Metal Supermarket Grocery Carts"),
					makeLeaf("plastic-shopping-baskets-handles", "Ergonomic Handheld Double Metal Handle Heavy Plastic Store Baskets", "የእጅ መያዣ ያላቸው የፕላስቲክ ቅርጫቶች", "Kookii Harkaa Pilaastikii Suuqii", "ናይ Ergonomic Handheld Double Metal Handle Heavy Plastic Store Baskets"),
					makeLeaf("supermarket-checkout-counter-desks", "Stainless Steel Top Conveyor Belt Electric Retail Checkout Counters", "የሱፐርማርኬት እቃ መክፈያ እና ማሳለፊያ ዴስኮች", "Minjaala Kaffaltii Supermarket", "ናይ Stainless Steel Top Conveyor Belt Electric Retail Checkout Counters"),
					makeLeaf("slanted-wooden-vegetable-fruit-display-racks", "Ventilated Plastic Crate Fruit & Green Vegetable Display Stands", "የፍራፍሬ እና የአትክልት ማሳያ ዘመናዊ መደርደሪያዎች", "Dhaabbii Fuduraa fi Kuduraa", "ናይ Ventilated Plastic Crate Fruit & Green Vegetable Display Stands"),
					makeLeaf("chrome-wire-storage-mesh-bins", "Folding Promo Dump Bins for Supermarket Center Aisle Clearance Sales", "ለቅናሽ እቃዎች ማስቀመጫ የሚሆኑ የብረት ማሳያ ሳጥኖች", "Saanduqa Sibilaa Gurgurtaaf", "ናይ Folding Promo Dump Bins for Supermarket Center Aisle Clearance Sales"),
					makeLeaf("electronic-supermarket-pricing-scales", "Digital Thermal Barcode Printing Weight Price Computing Scales", "ክብደት መዝኖ ባርኮድ የሚያትም የሱፐርማርኬት ሚዛን", "Mizana Suuqii Baarkoodii Maxxansu", "ናይ Digital Thermal Barcode Printing Weight Price Computing Scales"),
				},
			},
		},
	})

	// Domain: business-opportunities
	list = append(list, DomainCatalogDef{
		RootSlug: "business-opportunities",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "established-businesses-for-sale",
				NameEN: "Operational Cafes, Boutiques and Clinics for Sale",
				NameAM: "የሚሸጡ ካፌዎች፣ ሬስቶራንቶች፣ ቡቲኮች እና ፋርማሲዎች",
				NameOM: "Daldala Dhaabbataa Gurguramu",
				NameTI: "ዝሽየጡ ካፌታትን ቡቲካትን",
				Items: []LeafItemDef{
					makeLeaf("turnkey-specialty-cafe-restaurant-sale", "Fully Furnished Prime Location Espresso Bar & Restaurant Operations", "የተሟላ እቃ ያላቸው የሚሸጡ ካፌዎች እና ሬስቶራንቶች", "Kaaffee fi Nyaata Guutuu Gurguramu", "ናይ Fully Furnished Prime Location Espresso Bar & Restaurant Operations"),
					makeLeaf("stocked-boutique-clothing-store-sale", "High-Footfall Mall Fashion Apparel Store with Stock & Decor Included", "እቃ እና ዲኮር የተሟላለት የሚሸጥ የልብስ ቡቲክ", "Buutikii Uffataa Guutuu Gurguramu", "ናይ High-Footfall Mall Fashion Apparel Store with Stock & Decor Included"),
					makeLeaf("licensed-active-pharmacy-drug-store-sale", "Turnkey Ministry Registered Community Retail Pharmacy Operations", "ፈቃድ እና መድኃኒት ያለው የሚሸጥ የፋርማሲ ስራ", "Faarmaasii Qorichaa Gurguramu", "ናይ Turnkey Ministry Registered Community Retail Pharmacy Operations"),
					makeLeaf("operational-beauty-salon-spa-sale", "Hydraulic Chairs, Hair Steamers and Nail Bar Fully Equipped Salons", "ሙሉ እቃ ያለው የሚሸጥ የሴቶች ጸጉር እና የውበት ሳሎን", "Saaloona Bareedinaa Guutuu Gurguramu", "ናይ Hydraulic Chairs, Hair Steamers and Nail Bar Fully Equipped Salons"),
					makeLeaf("supermarket-minimarket-operations-sale", "Active POS System, Display Chillers & Complete Grocery Store Asset Sale", "የተሟላ ሱፐርማርኬት ወይም ሚኒማርኬት ሽያጭ", "Supermarket Guutuu Gurguramu", "ናይ Active POS System, Display Chillers & Complete Grocery Store Asset Sale"),
					makeLeaf("water-purification-bottling-plant-sale", "Reverse Osmosis Mineral Water Packaging Factory Operational Sale", "የውሃ ማጣሪያ እና ማሸጊያ ፋብሪካ ሽያጭ", "Warshaa Bishaan Qulqulleessu Gurguramu", "ናይ Reverse Osmosis Mineral Water Packaging Factory Operational Sale"),
					makeLeaf("diagnostic-medical-laboratory-for-sale", "Hematology, Chemistry Analyzers & License Operational Lab Sale", "የተሟላ የህክምና ላቦራቶሪ ንግድ ሽያጭ", "Laaboraatorii Yaalaa Gurguramu", "ናይ Hematology, Chemistry Analyzers & License Operational Lab Sale"),
					makeLeaf("stationery-secretarial-service-shop-sale", "Heavy Copiers, Laminators, Stock & Prime University Gate Location", "የፎቶኮፒ እና የጽህፈት መሳሪያዎች ሱቅ ሽያጭ", "Suuqii Koompii fi Waraqaa Gurguramu", "ናይ Heavy Copiers, Laminators, Stock & Prime University Gate Location"),
				},
			},
			{
				Slug:   "franchises-joint-ventures-commercial-partnerships",
				NameEN: "Franchises, Joint Ventures and Investment Partnerships",
				NameAM: "ፍራንቻይዞች፣ የጋራ ኢንቨስትመንት እና የንግድ አጋርነት",
				NameOM: "Shariikummaa fi Hirmaannaa Daldalaa",
				NameTI: "ናይ ንግዲ ሽርክነትን ኢንቨስትመንትን",
				Items: []LeafItemDef{
					makeLeaf("fast-food-burger-pizza-franchise-license", "Brand Name, Menu Recipes & Standardized Fast Food Franchise Rights", "የፈጣን ምግብ (በርገር/ፒዛ) የፍራንቻይዝ መብት ሽያጭ", "Mirga Daldalaa Nyaata Dafaa (Franchise)", "ናይ Brand Name, Menu Recipes & Standardized Fast Food Franchise Rights"),
					makeLeaf("commercial-agricultural-farming-joint-venture", "Irrigated Land Owner Seeking Capital & Machinery Agricultural Partner", "በግብርና ስራ ላይ የመሬት እና የካፒታል የጋራ አጋርነት", "Shariikummaa Qonnaa fi Lafa", "ናይ Irrigated Land Owner Seeking Capital & Machinery Agricultural Partner"),
					makeLeaf("export-coffee-processing-equity-partnership", "ECX Registered Coffee Washing Station Equity Investment Offerings", "በቡና ኤክስፖርት እና ማቀነባበሪያ ላይ የጋራ ኢንቨስትመንት", "Hirmaannaa Daldala Bunaa Alaa", "ናይ ECX Registered Coffee Washing Station Equity Investment Offerings"),
					makeLeaf("pharmaceutical-manufacturing-share-offering", "Industrial Park Medical Device & IV Fluids Production Share Capital", "በመድኃኒት ፋብሪካ ግንባታ ላይ የአክሲዮን ሽያጭ", "Qooda Warshaa Qorichaa (Share)", "ናይ Industrial Park Medical Device & IV Fluids Production Share Capital"),
					makeLeaf("real-estate-commercial-development-syndicate", "Addis Ababa Prime Plot Co-Financing Mixed-Use Mall Syndicate Shares", "የንግድ ህንፃ ግንባታ የጋራ ባለቤትነት አክሲዮን", "Hirmaannaa Gamoo Daldalaa Ijaaruu", "ናይ Addis Ababa Prime Plot Co-Financing Mixed-Use Mall Syndicate Shares"),
					makeLeaf("commercial-transport-trucking-fleet-partner", "Freight Contract Secured Cross-Border SinoTruk Financing Partners", "የጭነት ትራንስፖርት መኪና የጋራ ንግድ አጋርነት", "Shariikummaa Konkolaataa Fe'iinsaa", "ናይ Freight Contract Secured Cross-Border SinoTruk Financing Partners"),
					makeLeaf("solar-mini-grid-rural-electrification-jv", "Off-Grid Solar Energy Micro-Utility Commercial Concession Rights", "የገጠር የሶላር መብራት ዝርጋታ የጋራ ኢንቨስትመንት", "Shariikummaa Elektiriikii Aduu", "ናይ Off-Grid Solar Energy Micro-Utility Commercial Concession Rights"),
					makeLeaf("distribution-agency-rights-international-brands", "Sole Ethiopian Regional Import & Distribution Representation Rights", "የአለምአቀፍ ብራንዶች ብቸኛ የኢትዮጵያ ወኪል መብት", "Mirga Dhiyeessaa Oomisha Addunyaa", "ናይ Sole Ethiopian Regional Import & Distribution Representation Rights"),
				},
			},
		},
	})

	// Domain: industrial-machinery
	list = append(list, DomainCatalogDef{
		RootSlug: "industrial-machinery",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "food-grain-processing-machinery",
				NameEN: "Flour Mills, Oil Pressers and Grain Processing Machinery",
				NameAM: "የእህል ወፍጮዎች፣ የዘይት መጭመቂያ እና የምግብ ማቀነባበሪያ ማሽኖች",
				NameOM: "Maashinii Midhaanii fi Nyaataa",
				NameTI: "ናይ እኽሊ መጥሓኒን ናይ ዘይቲ መጭመቕን",
				Items: []LeafItemDef{
					makeLeaf("commercial-diesel-electric-flour-mills", "High-Output Commercial Teff, Wheat & Maize Grinding Hammer Mills", "የጤፍ፣ ስንዴ እና የበቆሎ የኤሌክትሪክ እና የናፍጣ ወፍጮዎች", "Wafchoo Midhaanii Elektiriikii fi Diiyezilaa", "ናይ High-Output Commercial Teff, Wheat & Maize Grinding Hammer Mills"),
					makeLeaf("cold-press-spiral-oil-extraction-machines", "Sesame, Niger Seed & Peanut Automatic Screw Oil Press Machines", "የሰሊጥ እና የኑግ ዘይት መጭመቂያ አውቶማቲክ ማሽኖች", "Maashinii Zayita Nyaataa Cuunfu", "ናይ Sesame, Niger Seed & Peanut Automatic Screw Oil Press Machines"),
					makeLeaf("industrial-bakery-tunnel-rotary-ovens", "32-Tray Rotary Rack Diesel & Electric Commercial Bread Baking Ovens", "ትላልቅ የዳቦ መጋገሪያ ሮተሪ ኦቨን ማሽኖች (Rotary Oven)", "Oovonii Daabboo Gurguddaa (Rotary Oven)", "ናይ 32-Tray Rotary Rack Diesel & Electric Commercial Bread Baking Ovens"),
					makeLeaf("automatic-injera-baking-rotary-lines", "Conveyor Belt Gas/Electric Industrial Injera Production Systems", "በኮንቬየር የሚሰሩ አውቶማቲክ የእንጀራ መጋገሪያ ማሽኖች", "Maashinii Injeeraa Tolchu Ammayyaa", "ናይ Conveyor Belt Gas/Electric Industrial Injera Production Systems"),
					makeLeaf("automatic-mineral-water-rinse-fill-cap-monoblocks", "3-in-1 Automatic Washing, Filling & Capping Bottling Lines", "የታሸገ ውሃ ማጠቢያ፣ መሙያ እና መክደኛ አውቶማቲክ ማሽኖች", "Maashinii Bishaan Saamunaatti Naqu", "ናይ 3-in-1 Automatic Washing, Filling & Capping Bottling Lines"),
					makeLeaf("spaghetti-macaroni-pasta-extruders", "Stainless Steel Continuous Vacuum Dough Macaroni Extrusion Plants", "የፓስታ እና ማካሮኒ ማምረቻ የኢንዱስትሪ ማሽኖች", "Maashinii Paastaa fi Makkaroni Tolchu", "ናይ Stainless Steel Continuous Vacuum Dough Macaroni Extrusion Plants"),
					makeLeaf("coffee-roasting-drum-machines-commercial", "10kg-60kg Industrial Gas Heated Air-Cooled Coffee Bean Roasters", "የቡና መቁያ እና ማቀዝቀዣ የኢንዱስትሪ ማሽኖች", "Maashinii Buna Akaawu Gurguddaa", "ናይ 10kg-60kg Industrial Gas Heated Air-Cooled Coffee Bean Roasters"),
					makeLeaf("feed-pellet-mill-poultry-livestock", "Flat Die Animal Feed Mixing, Crushing & Pelletizing Extruders", "የዶሮ እና የከብት መኖ ማዘጋጃ እና ማድበያ ማሽኖች", "Maashinii Nyaata Horii fi Lukkuu Tolchu", "ናይ Flat Die Animal Feed Mixing, Crushing & Pelletizing Extruders"),
				},
			},
			{
				Slug:   "metal-plastic-packaging-manufacturing",
				NameEN: "Plastic Molding, Metal Lathes and Packaging Machinery",
				NameAM: "የፕላስቲክ መቅረጫ፣ የብረት ላዝ እና የማሸጊያ ማሽኖች",
				NameOM: "Maashinii Pilaastikii, Sibilaa fi Saanduqaa",
				NameTI: "ናይ ፕላስቲክን ማሸጊን ማሽናት",
				Items: []LeafItemDef{
					makeLeaf("plastic-injection-molding-machines", "Servo Hydraulic Precision Plastic Cup, Crate & Bucket Injection Machines", "የፕላስቲክ እቃዎች እና ባልዲዎች ማምረቻ ኢንጄክሽን ማሽኖች", "Maashinii Pilaastikii Tolchu (Injection)", "ናይ Servo Hydraulic Precision Plastic Cup, Crate & Bucket Injection Machines"),
					makeLeaf("metalworking-lathe-milling-machines", "Heavy-Duty Manual & CNC Metal Turning, Threading & Milling Lathes", "የብረት ቅርጽ ማውጫ እና ማስተካከያ የላዝ ማሽኖች (Lathe)", "Maashinii Sibila Qaru (Lathe)", "ናይ Heavy-Duty Manual & CNC Metal Turning, Threading & Milling Lathes"),
					makeLeaf("automatic-vffs-liquid-powder-packaging-machines", "Vertical Form Fill Seal (VFFS) Spice, Sugar & Milk Pouch Packers", "የዱቄት፣ ቅመም እና ፈሳሽ ማሸጊያ አውቶማቲክ ማሽኖች (VFFS)", "Maashinii Nyaata fi Daakuu Saamunaatti Hidhu", "ናይ Vertical Form Fill Seal (VFFS) Spice, Sugar & Milk Pouch Packers"),
					makeLeaf("plastic-film-blown-extrusion-machines", "High-Density Polyethylene (HDPE) Shopping Bag Film Blowing Lines", "የፌስታል እና የላስቲክ ሻንጣዎች ማምረቻ ማሽኖች", "Maashinii Fastaala fi Pilaastikii Tolchu", "ናይ High-Density Polyethylene (HDPE) Shopping Bag Film Blowing Lines"),
					makeLeaf("blow-molding-pet-bottle-blowing-machines", "Automatic 2-Cavity to 6-Cavity Infrared Heated PET Bottle Blowers", "የፕላስቲክ የውሃ እና የዘይት ኮዳዎች መንፊያ ማሽኖች (Blow Molding)", "Maashinii Saamunaa Bishaan Tolchu", "ናይ Automatic 2-Cavity to 6-Cavity Infrared Heated PET Bottle Blowers"),
					makeLeaf("corrugated-cardboard-carton-box-lines", "3-Ply & 5-Ply Corrugated Paper Box Slotting & Printing Machinery", "የካርቶን ሳጥኖች ማምረቻ እና ማተሚያ የፋብሪካ ማሽኖች", "Maashinii Saanduqa Kartoonii Tolchu", "ናይ 3-Ply & 5-Ply Corrugated Paper Box Slotting & Printing Machinery"),
					makeLeaf("hydraulic-scrap-metal-baling-presses", "High-Pressure Compacting Baling Presses for Steel & Aluminum Cans", "የብረታ ብረት ፍራሽ እና ቆርቆሮ መጨፍለቂያ ሀይድሮሊክ ማሽኖች", "Maashinii Sibila Caccabaa Dhiibu", "ናይ High-Pressure Compacting Baling Presses for Steel & Aluminum Cans"),
					makeLeaf("continuous-band-sealer-nitrogen-flush", "Horizontal Stainless Steel Foil Pouch Heat Sealing Conveyors", "የቺፕስ እና የምግብ እቃዎች አየር ሞልተው ማሸጊያ ማሽኖች", "Maashinii Saamunaa Qilleensa Guutu fi Cufu", "ናይ Horizontal Stainless Steel Foil Pouch Heat Sealing Conveyors"),
				},
			},
		},
	})

	// Domain: tools-workshop
	list = append(list, DomainCatalogDef{
		RootSlug: "tools-workshop",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "electric-cordless-power-tools",
				NameEN: "Electric and Cordless Drills, Grinders and Saws",
				NameAM: "የኤሌክትሪክ እና የባትሪ ድሪሎች፣ መቁረጫዎች እና መጋዞች",
				NameOM: "Meeshaalee Elektiriikii Qorannoo fi Muraa",
				NameTI: "ናይ ኤሌክትሪክ ድሪልን መቁረጽን",
				Items: []LeafItemDef{
					makeLeaf("brushless-cordless-drill-driver-kits", "20V Lithium-Ion Brushless 2-Speed Hammer Drill & Screwdriver Kits", "በባትሪ የሚሰሩ ኃይለኛ የኮድለስ ድሪሎች (Cordless Drill)", "Diriila Wayaalessii Baatirii 20V", "ናይ 20V Lithium-Ion Brushless 2-Speed Hammer Drill & Screwdriver Kits"),
					makeLeaf("heavy-angle-grinders-cutting-tools", "2400W 9-Inch Heavy-Duty Steel Cutting & Grinding Angle Grinders", "ትላልቅ የብረት መቁረጫ እና መሞረጃ ግራይንደሮች (ፍሌክስ)", "Filaaksii Sibila Muraa Humna Guddaa (Grinder)", "ናይ 2400W 9-Inch Heavy-Duty Steel Cutting & Grinding Angle Grinders"),
					makeLeaf("rotary-hammer-demolition-drills", "SDS-Plus 1500W Heavy Concrete Chipping Demolition Hammer Drills", "ኮንክሪት እና ድንጋይ የሚሰብሩ ትላልቅ የኤሌክትሪክ ሂልቲ ማሽኖች", "Maashinii Dhagaa fi Gamoo Caccabsu (Hilti)", "ናይ SDS-Plus 1500W Heavy Concrete Chipping Demolition Hammer Drills"),
					makeLeaf("electric-circular-wood-cutting-saws", "7-1/4 Inch Carbide Blade Precision Woodworking Circular Saws", "የእንጨት መቁረጫ የኤሌክትሪክ ሰርኩላር መጋዞች", "Magaaza Mukaa Elektiriikii", "ናይ 7-1/4 Inch Carbide Blade Precision Woodworking Circular Saws"),
					makeLeaf("heavy-duty-inverter-arc-welders", "250A-400A IGBT Digital Display Electric Arc Inverter Welding Machines", "አነስተኛ ኃይል የሚጠቀሙ ዘመናዊ የኢንቨርተር ብየዳ ማሽኖች", "Maashinii Biriqaa Sibilaa (Inverter Welder)", "ናይ 250A-400A IGBT Digital Display Electric Arc Inverter Welding Machines"),
					makeLeaf("laser-distance-meters-rangefinders", "100-Meter High-Precision Digital Laser Distance Measure Tools", "በሌዘር የመሬት እና የክፍል ርቀት መለኪያ መሳሪያዎች", "Meeshaa Safartuu Fageenyaa Leezerii", "ናይ 100-Meter High-Precision Digital Laser Distance Measure Tools"),
					makeLeaf("benchtop-wood-thickness-planers", "Dual-Blade 12.5-Inch Portable Timber Thicknesser Wood Planers", "የእንጨት ማለስለሻ እና ማስተካከያ ካስቴሪያ (Planer) ማሽኖች", "Kaasteeriyaa Mukaa Elektiriikii", "ናይ Dual-Blade 12.5-Inch Portable Timber Thicknesser Wood Planers"),
					makeLeaf("industrial-hot-air-heat-guns", "2000W Variable Temperature Bumper Plastic Welding Heat Guns", "የሙቀት አየር የሚነፉ ኃይለኛ የሂት ገን (Heat Gun) መሳሪያዎች", "Meeshaa Qilleensa Ho'aa (Heat Gun)", "ናይ 2000W Variable Temperature Bumper Plastic Welding Heat Guns"),
				},
			},
			{
				Slug:   "mechanic-hand-tools-workshop-gear",
				NameEN: "Mechanic Hand Tools, Sockets and Compressors",
				NameAM: "የመካኒክ እቃዎች፣ የካቻቪቴ እና የቁልፍ ስብስቦች",
				NameOM: "Meeshaalee Harkaa fi Kompireesaraa",
				NameTI: "ናይ መካኒክ ኣቑሑትን ኮምፕሬሰርን",
				Items: []LeafItemDef{
					makeLeaf("chrome-vanadium-socket-wrench-sets", "150-Piece 1/4, 3/8 & 1/2 Inch Drive Metric Ratchet Socket Toolboxes", "የቁልፍ እና የካቻቪቴ የተሟሉ ሳጥኖች (Socket Set)", "Kiteetii Qulfee Sibilaa (Socket Set)", "ናይ 150-Piece 1/4, 3/8 & 1/2 Inch Drive Metric Ratchet Socket Toolboxes"),
					makeLeaf("hydraulic-heavy-duty-floor-trolley-jacks", "3-Ton to 5-Ton Fast Lifting Vehicle Garage Hydraulic Floor Jacks", "የመኪና ማንሻ ባለ ጎማ ሀይድሮሊክ ክሪኮች (Floor Jack)", "Kirikii Konkolaataa Ol Kaasu (Jack)", "ናይ 3-Ton to 5-Ton Fast Lifting Vehicle Garage Hydraulic Floor Jacks"),
					makeLeaf("belt-drive-electric-air-compressors", "100L-300L 3HP Twin Cylinder Piston Workshop Air Compressors", "የአየር መሙያ እና ቀለም መርጫ ትላልቅ ኮምፕሬሰሮች", "Kompireesara Qilleensaa Gurguddaa", "ናይ 100L-300L 3HP Twin Cylinder Piston Workshop Air Compressors"),
					makeLeaf("pneumatic-impact-wrench-guns", "1/2-Inch & 3/4-Inch Drive Heavy Air Impact Lug Nut Wrench Guns", "ጎማ በፍጥነት የሚፈቱ የአየር መፍቻዎች (Impact Wrench)", "Meeshaa Qulfee Qilleensaa Konkolaataa", "ናይ 1/2-Inch & 3/4-Inch Drive Heavy Air Impact Lug Nut Wrench Guns"),
					makeLeaf("heavy-cast-iron-bench-vises", "6-Inch to 8-Inch 360-Degree Swivel Base Heavy Duty Workshop Vises", "የጠረጴዛ ላይ እቃ መቆንጠጫ ጠንካራ ሞርሶች (Bench Vise)", "Moorsoo Sibilaa Minjaalaa (Vise)", "ናይ 6-Inch to 8-Inch 360-Degree Swivel Base Heavy Duty Workshop Vises"),
					makeLeaf("combination-spanner-wrench-sets", "6mm to 32mm High-Torque Mirror Chrome Combination Ring Spanners", "የመኪና እና የማሽን መፍቻ ባለቀለበት ቁልፎች ስብስብ", "Qulfee Sibilaa Lakkoofsaa 6-32", "ናይ 6mm to 32mm High-Torque Mirror Chrome Combination Ring Spanners"),
					makeLeaf("digital-vernier-calipers-micrometers", "Stainless Steel 0.01mm Accuracy Electronic Vernier Calipers", "የብረት ውፍረት እና ስፋት መለኪያ ረቂቅ ቨርኒየር ካሊፐሮች", "Safartuu Qal'inaa (Vernier Caliper)", "ናይ Stainless Steel 0.01mm Accuracy Electronic Vernier Calipers"),
					makeLeaf("automotive-battery-charger-boosters", "12V/24V Heavy Commercial Fast Battery Chargers and Jump Starters", "የመኪና ባትሪ ቻርጀሮች እና ስታርተሮች", "Chaarchara Baatirii Konkolaataa", "ናይ 12V/24V Heavy Commercial Fast Battery Chargers and Jump Starters"),
				},
			},
		},
	})

	// Domain: construction-materials
	list = append(list, DomainCatalogDef{
		RootSlug: "construction-materials",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "cement-rebar-structural-steel",
				NameEN: "Cement, Rebar and Structural Steel",
				NameAM: "ሲሚንቶ፣ ብረት እና መዋቅራዊ ብረታ ብረቶች",
				NameOM: "Simintoo, Sibila fi Meeshaalee Gamoo",
				NameTI: "ሲሚንቶን ሓጺንን ናይ ህንጸት ኣቑሑትን",
				Items: []LeafItemDef{
					makeLeaf("portland-pozzolana-cement-42-5n", "Grade 42.5N/R High-Strength Bulk & Bagged Structural Cement", "ደረጃ 42.5 ከፍተኛ ጥንካሬ ያለው የህንፃ ሲሚንቶ (በኩንታል)", "Simintoo Sadarkaa 42.5 Jabaa", "ናይ Grade 42.5N/R High-Strength Bulk & Bagged Structural Cement"),
					makeLeaf("deformed-steel-rebar-rods-12mm-16mm", "Grade 60 High-Yield Thermo-Mechanically Treated Rebar Rods", "ደረጃ 60 የህንፃ ማሰሪያ ባለ ቅርጽ ብረቶች (12ሚሜ/16ሚሜ)", "Sibila Gamoo Muraa 12mm fi 16mm", "ናይ Grade 60 High-Yield Thermo-Mechanically Treated Rebar Rods"),
					makeLeaf("heavy-rebar-rods-20mm-24mm", "Grade 60 Structural Column & Beam Heavy Deformed Rebar Steel", "ለአምድ እና ለትላልቅ ህንፃዎች የሚሆኑ ወፍራም ብረቶች (20ሚሜ/24ሚሜ)", "Sibila Gamoo Utubaa 20mm fi 24mm", "ናይ Grade 60 Structural Column & Beam Heavy Deformed Rebar Steel"),
					makeLeaf("black-annealed-binding-wire-rolls", "1.6mm Flexible Soft Annealed Construction Rebar Binding Wire", "የብረት ማሰሪያ ለስላሳ ሽቦዎች በጥቅል", "Wayaara Sibila Hidhu", "ናይ 1.6mm Flexible Soft Annealed Construction Rebar Binding Wire"),
					makeLeaf("structural-steel-i-beams-h-profiles", "Hot-Rolled Structural Carbon Steel Universal Beams and Columns", "የህንፃ እና የመጋዘን ጣሪያ መስሪያ የብረት አይ-ቢሞች (I-Beam)", "Sibila Dheeraa Gamoo I-Beam", "ናይ Hot-Rolled Structural Carbon Steel Universal Beams and Columns"),
					makeLeaf("rectangular-square-hollow-tubes", "Mild Steel 2mm to 4mm Welded Tubular Structural Box Profiles", "የብረት ቱቦዎች እና ፕሮፋይሎች (ቦክስ ብረት)", "Sibila Booksii fi Ujummoo", "ናይ Mild Steel 2mm to 4mm Welded Tubular Structural Box Profiles"),
					makeLeaf("galvanized-corrugated-iron-roofing-sheets", "Gauge 28 & Gauge 30 Zinc Coated Corrugated Metal Roofing Sheets", "የጣሪያ ቆርቆሮዎች (ጌጅ 28 እና ጌጅ 30)", "Qorqorroo Manarraa G28 fi G30", "ናይ Gauge 28 & Gauge 30 Zinc Coated Corrugated Metal Roofing Sheets"),
					makeLeaf("expanded-metal-wire-mesh-rolls", "Welded BRC Concrete Slab Reinforcement Steel Wire Mesh Rolls", "ለወለል ኮንክሪት የሚነጠፉ የሽቦ መረቦች (BRC Mesh)", "Saaphana Sibilaa Dachee Gamoo (BRC)", "ናይ Welded BRC Concrete Slab Reinforcement Steel Wire Mesh Rolls"),
				},
			},
			{
				Slug:   "masonry-tiles-paints-finishing",
				NameEN: "Tiles, Paints, Glass and Finishing Materials",
				NameAM: "ሴራሚክ ፖርሲላይን፣ ቀለሞች እና የማስዋቢያ እቃዎች",
				NameOM: "Seeraamikii, Halluu fi Doowwituu Gamoo",
				NameTI: "ሴራሚክን ቀለማትን ናይ ፊኒሺንግ ኣቑሑትን",
				Items: []LeafItemDef{
					makeLeaf("polished-porcelain-floor-tiles-60x60", "Super-Gloss Nano-Treated Stain-Resistant Porcelain Floor Tiles", "የወለል ፖርሲላይን ንጹህ የሚያብረቀርቁ ሴራሚኮች (60x60)", "Seeraamikii Dachee Gamoo 60x60", "ናይ Super-Gloss Nano-Treated Stain-Resistant Porcelain Floor Tiles"),
					makeLeaf("anti-slip-ceramic-bathroom-wall-tiles", "Glazed Waterproof Anti-Bacterial Kitchen & Bathroom Wall Tiles", "የመታጠቢያ ቤት እና የኩሽና ግድግዳ ሴራሚኮች", "Seeraamikii Mana Finiinaa fi Kushiinaa", "ናይ Glazed Waterproof Anti-Bacterial Kitchen & Bathroom Wall Tiles"),
					makeLeaf("natural-granite-marble-countertop-slabs", "Polished 20mm Black Galaxy & White Quartz Kitchen Countertops", "የግራናይት እና እብነ-በረድ የወጥ ቤት እና ደረጃ ድንጋዮች", "Dhagaa Giraanaayitii fi Maarbiilii", "ናይ Polished 20mm Black Galaxy & White Quartz Kitchen Countertops"),
					makeLeaf("washable-acrylic-interior-wall-paints", "Anti-Fungal Matt Finish Premium Odorless Indoor Wall Emulsion Paints", "የሚታጠቡ የውስጥ ግድግዳ ጥራት ያላቸው ቀለሞች", "Halluu Girgidaa Keessoo Qulqulluu", "ናይ Anti-Fungal Matt Finish Premium Odorless Indoor Wall Emulsion Paints"),
					makeLeaf("weather-shield-exterior-wall-paints", "UV-Resistant Waterproof Exterior Acrylic Textured Facade Paints", "የፀሐይ እና የዝናብ መከላከያ የውጪ ግድግዳ ቀለሞች", "Halluu Girgidaa Alaa Qilleensa Danda'u", "ናይ UV-Resistant Waterproof Exterior Acrylic Textured Facade Paints"),
					makeLeaf("gypsum-plasterboard-ceiling-sheets", "9mm to 12mm Moisture-Resistant Gypsum False Ceiling Boards", "የኮርኒስ ጂፕሰም ቦርዶች እና የድጋፍ ፕሮፋይሎች", "Boordii Jiipsamii Gubbaa Manaa", "ናይ 9mm to 12mm Moisture-Resistant Gypsum False Ceiling Boards"),
					makeLeaf("tempered-architectural-float-glass-sheets", "6mm to 12mm Clear & Tinted Shatterproof Facade Curtain Glass", "የህንፃ መስኮት እና ገጽታ ጠንካራ የመስታወት ሉሆች", "Doowwituu Gamoo Jabaa (Tempered Glass)", "ናይ 6mm to 12mm Clear & Tinted Shatterproof Facade Curtain Glass"),
					makeLeaf("hollow-concrete-blocks-hcb-15-20", "Standard Crushed Basalt Machine-Vibrated Masonry Wall HCB Blocks", "የግድግዳ ማነጫ ጠንካራ ብሎኬቶች (HCB 15 እና HCB 20)", "Bilookeetii Gamoo Jabaa (HCB)", "ናይ Standard Crushed Basalt Machine-Vibrated Masonry Wall HCB Blocks"),
				},
			},
		},
	})

	// Domain: construction-machinery
	list = append(list, DomainCatalogDef{
		RootSlug: "construction-machinery",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "earthmoving-excavators-loaders",
				NameEN: "Excavators, Wheel Loaders and Bulldozers",
				NameAM: "ኤክስካቫተሮች፣ ሎደሮች እና ቡልዶዘሮች",
				NameOM: "Ekiskaaveetara, Loodara fi Buldoozara",
				NameTI: "ኤክስካቫተርን ሎደርን ቡልዶዘርን",
				Items: []LeafItemDef{
					makeLeaf("hydraulic-crawler-excavators-20-ton", "Heavy-Duty 20-30 Ton Hydraulic Earthmoving Crawler Excavators", "ከባድ የቁፋሮ ስራ የሚሰሩ ባለ ሰንሰለት ኤክስካቫተሮች", "Ekiskaaveetara Guddaa 20-Ton", "ናይ Heavy-Duty 20-30 Ton Hydraulic Earthmoving Crawler Excavators"),
					makeLeaf("wheel-loaders-quarry-buckets", "3-Cubic Meter Articulated Front-End Earthmoving Wheel Loaders", "የአሸዋ እና የድንጋይ መጫኛ ባለ ጎማ ትላልቅ ሎደሮች", "Loodara Qorii Bal'aa (Wheel Loader)", "ናይ 3-Cubic Meter Articulated Front-End Earthmoving Wheel Loaders"),
					makeLeaf("crawler-bulldozers-rippers", "Heavy Track-Type Bulldozers with Hydraulic Tilt Blades & Rippers", "የመሬት ማስተካከያ እና ገፊ ትላልቅ ቡልዶዘሮች (Bulldozer)", "Buldoozara Lafa Qopheessu", "ናይ Heavy Track-Type Bulldozers with Hydraulic Tilt Blades & Rippers"),
					makeLeaf("motor-graders-road-maintenance", "14-Foot Moldboard Articulated Tandem-Drive Road Motor Graders", "የመንገድ ጠጠር እና አፈር ማስተካከያ ግሬደሮች (Grader)", "Gireedara Daandii Hojjetu", "ናይ 14-Foot Moldboard Articulated Tandem-Drive Road Motor Graders"),
					makeLeaf("single-drum-vibratory-soil-compactors", "12-16 Ton Heavy Hydraulic Vibrating Drum Soil & Asphalt Rollers", "የአፈር እና የአስፋልት መደለዣ ትላልቅ ሮለሮች (Roller)", "Roolara Biyyoo fi Asfaaltii Dhiibu", "ናይ 12-16 Ton Heavy Hydraulic Vibrating Drum Soil & Asphalt Rollers"),
					makeLeaf("backhoe-loaders-4x4-tractor", "Multi-Purpose Front Loader & Rear Digging Bucket Backhoe Units", "የፊት ሎደር እና የኋላ መቆፈሪያ ያላቸው ባክሆ ሎደሮች", "Baakhoo Loodara (Backhoe Loader)", "ናይ Multi-Purpose Front Loader & Rear Digging Bucket Backhoe Units"),
					makeLeaf("mini-crawler-excavators-rubber-track", "Rubber Track Compact Tail-Swing Urban Utility Mini Excavators", "ለከተማ ውስጥ አነስተኛ ቁፋሮ የሚሆኑ ሚኒ ኤክስካቫተሮች", "Ekiskaaveetara Xiqqaataa (Mini Excavator)", "ናይ Rubber Track Compact Tail-Swing Urban Utility Mini Excavators"),
					makeLeaf("heavy-hydraulic-rock-breakers-hammers", "Excavator Attachment Heavy Hydraulic Stone & Concrete Breakers", "ከኤክስካቫተር ጋር የሚገጠሙ የድንጋይ መስበሪያ ጃክሀመሮች", "Hambara Dhagaa Caccabsu (Rock Breaker)", "ናይ Excavator Attachment Heavy Hydraulic Stone & Concrete Breakers"),
				},
			},
			{
				Slug:   "concrete-machinery-scaffolding-cranes",
				NameEN: "Concrete Mixers, Vibrators, Cranes and Scaffolding",
				NameAM: "የኮንክሪት ማቡኪያ፣ ቫይብሬተር፣ ክሬኖች እና ስካፎልዲንግ",
				NameOM: "Maashinii Konkiritii, Kireenii fi Iskaafoldiingii",
				NameTI: "ናይ ኮንክሪት መለወሽን ክሬንን",
				Items: []LeafItemDef{
					makeLeaf("diesel-drum-concrete-mixers-hopper", "Hydraulic Tipping Hopper 400L-500L Towable Site Concrete Mixers", "በናፍጣ የሚሰሩ የኮንክሪት ማቡኪያ ማሽኖች (ሚክሰር)", "Maashinii Konkiritii Maku (Mixer)", "ናይ Hydraulic Tipping Hopper 400L-500L Towable Site Concrete Mixers"),
					makeLeaf("gasoline-engine-concrete-vibrators", "Heavy Flexible Shaft Poker Needle Immersion Concrete Vibrators", "የኮንክሪት አየር ማውጫ እና ማደለዣ ቫይብሬተሮች", "Maashinii Konkiritii Sochoosu (Vibrator)", "ናይ Heavy Flexible Shaft Poker Needle Immersion Concrete Vibrators"),
					makeLeaf("telescopic-boom-rough-terrain-mobile-cranes", "25-50 Ton Heavy Lifting Hydraulic Mobile Truck Cranes", "ከ25 እስከ 50 ቶን እቃ የሚያነሱ ትላልቅ ተንቀሳቃሽ ክሬኖች", "Kireenii Meeshaa Ol Kaasu 25-50 Ton", "ናይ 25-50 Ton Heavy Lifting Hydraulic Mobile Truck Cranes"),
					makeLeaf("building-construction-tower-cranes", "60-Meter Jib Stationary High-Rise Construction Topkit Tower Cranes", "ለረጃጅም ህንፃዎች ግንባታ የሚተከሉ ታወር ክሬኖች", "Taawar Kireenii Gamoo Dheeraa", "ናይ 60-Meter Jib Stationary High-Rise Construction Topkit Tower Cranes"),
					makeLeaf("tubular-steel-scaffolding-frame-sets", "Walk-Through High-Strength Steel Exterior Building Scaffolding Sets", "የህንፃ ውጪ መስሪያ ጠንካራ የብረት ስካፎልዲንጎች", "Iskaafoldiingii Sibilaa Gamoo", "ናይ Walk-Through High-Strength Steel Exterior Building Scaffolding Sets"),
					makeLeaf("acrow-steel-telescopic-props-jacks", "Heavy Duty Concrete Slab Formwork Support Telescopic Steel Props", "የኮንክሪት ጣሪያ መደገፊያ የሚረዝሙ የብረት ጃኮች (Props)", "Jaakii Sibilaa Konkiritii Deeggaru", "ናይ Heavy Duty Concrete Slab Formwork Support Telescopic Steel Props"),
					makeLeaf("concrete-transit-mixer-trucks-10m3", "Heavy 6x4 & 8x4 Truck Chassis 10m3 Rotating Concrete Drum Trucks", "የተቦካ ኮንክሪት የሚያጓጉዙ ትላልቅ የኮንክሪት መኪኖች", "Konkolaataa Konkiritii Naannessu (Transit Mixer)", "ናይ Heavy 6x4 & 8x4 Truck Chassis 10m3 Rotating Concrete Drum Trucks"),
					makeLeaf("stationary-concrete-trailer-pumps", "High-Pressure Hydraulic Twin Cylinder Piston Concrete Pumps", "ኮንክሪት ወደ ከፍተኛ ህንፃ የሚረጩ የኮንክሪት ፓምፖች", "Paampii Konkiritii Gamoo Dheeraa", "ናይ High-Pressure Hydraulic Twin Cylinder Piston Concrete Pumps"),
				},
			},
		},
	})

	return list
}
