package main

// NestedCategoryDef defines a child or leaf category in the hierarchy
type NestedCategoryDef struct {
	ParentSlug   string
	Slug         string
	NameEN       string
	NameAM       string
	NameOM       string
	NameTI       string
	Description  string
	IsLeaf       bool
	DisplayOrder int
}

// GetNestedCategories returns the deep hierarchy across all 60 domains
func GetNestedCategories() []NestedCategoryDef {
	var list []NestedCategoryDef

	// Helper function to append nodes
	add := func(parent, slug, en, am, om, ti string, isLeaf bool, order int) {
		list = append(list, NestedCategoryDef{
			ParentSlug:   parent,
			Slug:         slug,
			NameEN:       en,
			NameAM:       am,
			NameOM:       om,
			NameTI:       ti,
			IsLeaf:       isLeaf,
			DisplayOrder: order,
		})
	}

	// =========================================================================
	// DOMAIN 01: Electronics and Mobile Devices (electronics-mobile)
	// =========================================================================
	add("electronics-mobile", "mobile-phones", "Mobile Phones", "ሞባይል ስልኮች", "Bilbila Moobaayilii", "ሞባይል ስልክታት", false, 1)
	add("mobile-phones", "smartphones-android", "Android Smartphones", "አንድሮይድ ስማርት ስልኮች", "Bilbila Ismaartii Andirooyidii", "ኣንድሮይድ ስማርት ስልኪ", true, 1)
	add("mobile-phones", "smartphones-ios", "iPhones & iOS Devices", "አይፎን ስልኮች", "Bilbila Ismaartii Ayifoonii", "ኣይፎን ስልኪ", true, 2)
	add("mobile-phones", "feature-phones", "Feature Phones & Keypad Devices", "ባለ ኪፓድ ስልኮች", "Bilbila Kiipaadii", "ናይ ቁልፊ ሞባይል ስልኪ", true, 3)
	add("mobile-phones", "refurbished-phones", "Certified Refurbished Phones", "የታደሱ ሞባይል ስልኮች", "Bilbila Haaromfame", "ዝተሓደሰ ሞባይል ስልኪ", true, 4)

	add("electronics-mobile", "tablets-ereaders", "Tablets and E-Readers", "ታብሌቶች እና ኢ-ሪደር", "Taableetii fi E-Reader", "ታብሌታትን ኢ-ሪደርን", false, 2)
	add("tablets-ereaders", "android-tablets", "Android Tablets", "አንድሮይድ ታብሌቶች", "Taableetii Andirooyidii", "ኣንድሮይድ ታብሌት", true, 1)
	add("tablets-ereaders", "apple-ipads", "Apple iPad Series", "አፕል አይፓድ", "Ayipaadii Apple", "ኣይፓድ", true, 2)
	add("tablets-ereaders", "ereaders-kindle", "E-Book Readers (Kindle/Kobo)", "የመጽሐፍ አንባቢ ታብሌቶች", "Dubaasaa Kitaabaa Dijitaalaa", "ናይ መጽሓፍ ምንባብ ታብሌት", true, 3)
	add("tablets-ereaders", "kids-tablets", "Children & Educational Tablets", "የልጆች ትምህርታዊ ታብሌቶች", "Taableetii Daa'immanii", "ናይ ቆልዑ ታብሌት", true, 4)

	add("electronics-mobile", "smartwatches-wearables", "Smartwatches and Wearables", "ስማርት ሰዓቶች እና ተለባሾች", "Sa'aatii Ismaartii", "ስማርት ሰዓታትን", false, 3)
	add("smartwatches-wearables", "apple-watches", "Apple Watch", "አፕል ስማርት ሰዓት", "Sa'aatii Apple", "ኣፕል ስማርት ሰዓት", true, 1)
	add("smartwatches-wearables", "android-smartwatches", "Android / WearOS Smartwatches", "አንድሮይድ ስማርት ሰዓቶች", "Sa'aatii Andirooyidii", "ኣንድሮይድ ስማርት ሰዓት", true, 2)
	add("smartwatches-wearables", "fitness-trackers", "Fitness Trackers & Bands", "የአካል ብቃት መከታተያ ባንዶች", "Hordoftuu Fayyaa", "ናይ ጥዕና መከታተሊ ባንድ", true, 3)

	add("electronics-mobile", "mobile-accessories", "Mobile Accessories", "የስልክ ተጨማሪ እቃዎች", "Meeshaalee Dabalataa Bilbilaa", "ናይ ሞባይል መተሓሓዚ ኣቑሑት", false, 4)
	add("mobile-accessories", "chargers-adapters", "Fast Chargers and Adapters", "ፈጣን ቻርጀሮች", "Chaarcharoota Dafan", "ፈጣናት ቻርጀራት", true, 1)
	add("mobile-accessories", "power-banks", "Power Banks & Portable Chargers", "ፓወር ባንኮች", "Paawar Baankii", "ፓወር ባንክታት", true, 2)
	add("mobile-accessories", "cables-converters", "USB Cables and Type-C Converters", "የዩኤስቢ ኬብሎች", "Keebilii USB", "ኬብላት ዩኤስቢ", true, 3)
	add("mobile-accessories", "cases-covers", "Protective Cases and Back Covers", "የስልክ ከቨሮች", "Koversii Bilbilaa", "ናይ ሞባይል ከቨራት", true, 4)
	add("mobile-accessories", "screen-protectors", "Tempered Glass Screen Protectors", "የስክሪን መከላከያ መስታወቶች", "Eegduu Iskiriinii", "መከላኸሊ ስክሪን መስተዋት", true, 5)
	add("mobile-accessories", "tws-earbuds", "True Wireless Earbuds (TWS)", "ገመድ አልባ ኤርፖዶች", "Airpods Sagalee", "ገመድ ኣልቦ ኤርፖድስ", true, 6)
	add("mobile-accessories", "car-phone-holders", "Car Phone Mounts and Wireless Chargers", "የመኪና ውስጥ የስልክ መያዣዎች", "Qabduu Bilbila Konkolaataa", "ናይ መኪና ሞባይል መትሓዚ", true, 7)
	add("mobile-accessories", "phone-repair-parts", "Replacement Screens and Internal Parts", "የስልክ ስክሪኖች እና መለዋወጫዎች", "Qaama Jijjiirraa Bilbilaa", "ናይ ሞባይል መለዋወጢ ክፍልታት", true, 8)

	// =========================================================================
	// DOMAIN 02: Computers, IT and Office Technology (computers-it)
	// =========================================================================
	add("computers-it", "laptops-notebooks", "Laptops and Notebooks", "ላፕቶፖች", "Laptooppiiwwan", "ላፕቶፓት", false, 1)
	add("laptops-notebooks", "business-laptops", "Business & Professional Laptops", "የስራ እና የቢሮ ላፕቶፖች", "Laptooppii Hojii", "ናይ ስራሕ ላፕቶፕ", true, 1)
	add("laptops-notebooks", "gaming-laptops", "High-End Gaming Laptops", "የጌሚንግ ላፕቶፖች", "Laptooppii Geemiingii", "ናይ ጌሚንግ ላፕቶፕ", true, 2)
	add("laptops-notebooks", "ultrabooks-macbooks", "Ultrabooks & Apple MacBooks", "አልትራቡክ እና ማክቡክ", "MacBook fi Ultrabook", "ማክቡክን ኣልትራቡክን", true, 3)
	add("laptops-notebooks", "budget-student-laptops", "Student & Budget Laptops", "የተማሪ እና የበጀት ላፕቶፖች", "Laptooppii Barattootaa", "ናይ ተምሃሮ ላፕቶፕ", true, 4)

	add("computers-it", "desktop-computers", "Desktop Computers & Workstations", "የዴስክቶፕ ኮምፒውተሮች", "Kompiitara Minjaalaa", "ናይ ዴስክቶፕ ኮምፒዩተር", false, 2)
	add("desktop-computers", "all-in-one-pcs", "All-in-One Desktop PCs", "ኦል-ኢን-ዋን ኮምፒውተሮች", "Kompiitara All-in-One", "ኦል-ኢን-ዋን ፒሲ", true, 1)
	add("desktop-computers", "tower-workstations", "Engineering & Design Workstations", "የምህንድስና ወርክስቴሽኖች", "Workstation Injinariingii", "ናይ ምህንድስና ወርክስቴሽን", true, 2)
	add("desktop-computers", "custom-gaming-rigs", "Custom Built Gaming Rigs", "የተገጣጠሙ ጌሚንግ ዴስክቶፖች", "Desktop Geemiingii", "ናይ ጌሚንግ ዴስክቶፕ", true, 3)
	add("desktop-computers", "mini-pcs", "Mini PCs and Thin Clients", "ሚኒ ፒሲዎች", "Mini PC", "ሚኒ ፒሲታት", true, 4)

	add("computers-it", "computer-monitors", "Computer Monitors & Displays", "ሞኒተሮች", "Moniitaroota", "ሞኒተራት", false, 3)
	add("computer-monitors", "ultrawide-4k-monitors", "4K & Ultrawide Productivity Displays", "4ኬ ሞኒተሮች", "Moniitara 4K", "4ኬ ሞኒተራት", true, 1)
	add("computer-monitors", "gaming-high-refresh-monitors", "High Refresh Rate Gaming Monitors", "የጌሚንግ ሞኒተሮች (144Hz+)", "Moniitara Geemiingii", "ናይ ጌሚንግ ሞኒተራት", true, 2)
	add("computer-monitors", "standard-office-monitors", "Standard Office LED Monitors", "መደበኛ የቢሮ ሞኒተሮች", "Moniitara Hojii Biiroo", "ናይ ቢሮ ሞኒተራት", true, 3)

	add("computers-it", "storage-components", "Internal Storage and PC Components", "የኮምፒውተር እቃዎች እና ማከማቻ", "Qaamolee Kompiitaraa", "ክፍልታት ኮምፒዩተር", false, 4)
	add("storage-components", "nvme-ssd-drives", "M.2 NVMe and SATA SSDs", "ኤስኤስዲ ድራይቮች (SSD)", "Diraayivii SSD", "ኤስኤስዲ ድራይቭታት", true, 1)
	add("storage-components", "ram-modules-pc", "DDR4 & DDR5 RAM Modules", "የኮምፒውተር ራም (RAM)", "RAM Kompiitaraa", "ናይ ኮምፒዩተር ራም", true, 2)
	add("storage-components", "graphics-cards-gpu", "Dedicated Graphics Cards (NVIDIA/AMD)", "የግራፊክስ ካርዶች (GPU)", "Kaardii Giraafiiksii (GPU)", "ናይ ግራፊክስ ካርድታት", true, 3)
	add("storage-components", "external-hard-drives", "External Portable Hard Drives & Flash Disks", "ኤክስተርናል ሃርድ ዲስኮች እና ፍላሽ", "Hard Drive fi Flash Ala", "ናይ ወጻኢ ሃርድ ዲስክን ፍላሽን", true, 4)

	add("computers-it", "printers-scanners-pos", "Printers, Scanners and POS Systems", "ፕሪንተሮች እና ስካነሮች", "Piriintara fi Iskaanara", "ፕሪንተራትን ስካነራትን", false, 5)
	add("printers-scanners-pos", "laser-inkjet-printers", "Office Laser & Color Inkjet Printers", "የቢሮ ሌዘር ፕሪንተሮች", "Piriintara Leezarii", "ናይ ቢሮ ሌዘር ፕሪንተር", true, 1)
	add("printers-scanners-pos", "receipt-pos-printers", "Thermal Receipt & POS Printers", "የደረሰኝ እና የቲኬት ፕሪንተሮች", "Piriintara Rasiidaa POS", "ናይ ቅብሊት ፕሪንተር", true, 2)
	add("printers-scanners-pos", "high-speed-scanners", "Flatbed and Document Scanners", "ሰነድ ስካነሮች", "Iskaanara Sanadaa", "ናይ ሰነድ ስካነር", true, 3)
	add("printers-scanners-pos", "toner-cartridges", "Toner and Ink Cartridges", "የፕሪንተር ቀለም እና ቶነር", "Waraqaa fi Qalama Piriintaraa", "ቀለም ፕሪንተርን ቶነርን", true, 4)

	// =========================================================================
	// DOMAIN 03: Telecommunications and Networking (telecom-networking)
	// =========================================================================
	add("telecom-networking", "wireless-routers-modems", "Wireless Routers and Modems", "ዋይፋይ ራውተሮች እና ሞደሞች", "Raawutara fi Moodamii", "ራውተራትን ሞደማትን", false, 1)
	add("wireless-routers-modems", "4g-5g-sim-routers", "4G / 5G LTE SIM Wireless Routers", "የሲም ካርድ 4ጂ/5ጂ ራውተሮች", "Raawutara SIM 4G/5G", "ናይ ሲም ካርድ ራውተር 4G/5G", true, 1)
	add("wireless-routers-modems", "mesh-wifi-extenders", "Mesh Wi-Fi Systems and Signal Extenders", "የዋይፋይ ሲግናል ማጉያ መሳሪያዎች", "Bal'isaa Mallattoo Wi-Fi", "ናይ ዋይፋይ ምልክት መጉልሒ", true, 2)
	add("wireless-routers-modems", "adsl-vdsl-modems", "Fiber ONU & Broadband Modems", "የፋይበር እና ብሮድባንድ ሞደሞች", "Moodamii Faayibaraa", "ናይ ፋይበር ሞደም", true, 3)

	add("telecom-networking", "enterprise-networking", "Switches, Racks and Patch Panels", "ኔትወርክ ስዊቾች እና ራኮች", "Iswiikii fi Raakii Netwoorkii", "ስዊቻትን ራካትን", false, 2)
	add("enterprise-networking", "managed-poe-switches", "Managed Gigabit & PoE Switches", "ጊጋቢት ኔትወርክ ስዊቾች", "Iswiikii Giigaabiitii", "ጊጋቢት ኔትወርክ ስዊች", true, 1)
	add("enterprise-networking", "server-cabinets-racks", "Server Cabinets and Data Racks", "የሰርቨር ራኮች እና ሳጥኖች", "Saanduqa Sarvarii", "ናይ ሰርቨር ራክን ካቢኔትን", true, 2)
	add("enterprise-networking", "patch-panels-cords", "Cat6/Cat6A Cables, Patch Panels & Connectors", "የኔትወርክ ኬብሎች እና አገናኞች", "Keebilii fi Qabsiisaa Cat6", "ኬብላት ኔትወርክን ኣገናኘትን", true, 3)
	add("enterprise-networking", "fiber-optic-cables-tools", "Fiber Optic Patch Cords, SFP Modules & Splicers", "የፋይበር ኦፕቲክ ኬብሎች እና መሳሪያዎች", "Meeshaa Faayibar Ooptikii", "ናይ ፋይበር ኦፕቲክ መሳርሒታት", true, 4)

	// =========================================================================
	// DOMAIN 06: TVs, Audio and Home Entertainment (tv-audio)
	// =========================================================================
	add("tv-audio", "televisions-screens", "Televisions and Smart Screens", "ቴሌቪዥኖች", "Televiziyinii", "ቴሌቪዥናት", false, 1)
	add("televisions-screens", "4k-smart-tvs", "4K Ultra HD Smart LED TVs", "4ኬ ስማርት ቴሌቪዥኖች", "Televiziyinii 4K Ismaartii", "4ኬ ስማርት ቴሌቪዥን", true, 1)
	add("televisions-screens", "oled-qled-tvs", "OLED & QLED Premium TVs", "ኦሌድ እና ኪውሌድ ቴሌቪዥኖች", "Televiziyinii OLED/QLED", "OLED/QLED ቴሌቪዥን", true, 2)
	add("televisions-screens", "compact-tvs-monitors", "32-inch to 43-inch Compact Bedroom TVs", "መካከለኛ ቴሌቪዥኖች (32-43 ኢንች)", "Televiziyinii Giddu-galeessaa", "ማእከላይ ቴሌቪዥን", true, 3)

	add("tv-audio", "home-audio-theaters", "Soundbars and Home Audio Systems", "የድምጽ ማጉያዎች እና ሆም ቲያትር", "Sagalee fi Saawundibaarii", "ናይ ድምጺ ሲስተማት", false, 2)
	add("home-audio-theaters", "soundbars-subwoofers", "Soundbars with Wireless Subwoofers", "ሳውንድባር ከሰብዉፈር ጋር", "Saawundibaarii", "ሳውንድባር", true, 1)
	add("home-audio-theaters", "home-theaters-5-1", "5.1 & 7.1 Surround Home Theaters", "ሆም ቲያትር ሲስተሞች", "Sirna Sagalee Manaa", "ሆም ቲያትር ሲስተም", true, 2)
	add("home-audio-theaters", "bluetooth-party-speakers", "Rechargeable Bluetooth Party Speakers", "ተንቀሳቃሽ የብሉቱዝ ስፒከሮች", "Ispiikara Biliituutii", "ተንቀሳቓሲ ብሉቱዝ ስፒከር", true, 3)
	add("home-audio-theaters", "satellite-decoders-dishes", "Satellite Receivers, Decoders and Parabolic Dishes", "ዲኮደሮች እና የሳተላይት ዲሾች", "Diishiifi Dekoodarii", "ዲኮደርን ሳተላይት ዲሽን", true, 4)

	// =========================================================================
	// DOMAIN 07: Home Appliances (home-appliances)
	// =========================================================================
	add("home-appliances", "refrigeration", "Refrigerators and Freezers", "ፍሪጆች እና ዲፕ ፍሪዘሮች", "Firiijii fi Qabbaneessituu", "ፍሪጃትን ዲፕ ፍሪዘራትን", false, 1)
	add("refrigeration", "double-door-refrigerators", "Double Door & French Door Refrigerators", "ባለ ሁለት በር ፍሪጆች", "Firiijii Balbala Lammaa", "ባዕል ክልተ ማዕጾ ፍሪጅ", true, 1)
	add("refrigeration", "chest-deep-freezers", "Chest Deep Freezers (Meat & Commercial)", "ዲፕ ፍሪዘሮች", "Firiizara Gurguddaa", "ዲፕ ፍሪዘር", true, 2)
	add("refrigeration", "single-door-mini-fridges", "Compact & Mini Bar Refrigerators", "ሚኒ ፍሪጆች", "Firiijii Xiqqaa", "ንእሽቶ ፍሪጅ", true, 3)

	add("home-appliances", "washers-dryers", "Washing Machines and Dryers", "የልብስ ማጠቢያ ማሽኖች", "Maashina Uffata Dhiqu", "ናይ ክዳን ምሕጻብ ማሽን", false, 2)
	add("washers-dryers", "front-load-washers", "Automatic Front-Load Washing Machines", "ከፊት የሚከፈቱ አውቶማቲክ ማጠቢያዎች", "Maashina Dhiqaa Fuulduraa", "ኣውቶማቲክ ናይ ክዳን ምሕጻብ ማሽን", true, 1)
	add("washers-dryers", "top-load-washers", "Top-Load Washing Machines", "ከላይ የሚከፈቱ የልብስ ማጠቢያዎች", "Maashina Dhiqaa Gubbaa", "ናይ ላዕሊ ማጠቢ ማሽን", true, 2)
	add("washers-dryers", "twin-tub-semi-automatic", "Semi-Automatic Twin-Tub Washers", "ባለ ሁለት ገንዳ ከፊል አውቶማቲክ ማጠቢያዎች", "Maashina Qulqulleessaa Uffataa", "መንታ ጋን ማጠቢ ማሽን", true, 3)

	add("home-appliances", "cooking-appliances", "Stoves, Ovens and Ethiopian Mitad", "የማብሰያ እቃዎች እና ምጣዶች", "Eelee fi Meeshaalee Bilcheessaa", "ናይ መብሰሊ ንዋያትን ምጣድን", false, 3)
	add("cooking-appliances", "electric-injera-mitad", "Electric Teff Injera Baking Mitad (Energy Saving)", "የኤሌክትሪክ የእንጀራ ምጣዶች", "Eelee Injeeraa Elektiriikii", "ናይ ኤሌክትሪክ ናይ እንጀራ ምጣድ", true, 1)
	add("cooking-appliances", "gas-electric-stoves", "Standing Gas & Electric Cookers / Ovens", "የጋዝ እና የኤሌክትሪክ ምድጃዎች", "Istoovii Gaasii fi Elektiriikii", "ናይ ጋዝን ኤሌክትሪክን እቶናት", true, 2)
	add("cooking-appliances", "microwave-airfryers", "Microwave Ovens and Air Fryers", "ማይክሮዌቭ እና ኤር ፍራየር", "Maayikirooweevii", "ማይክሮዌቭን ኤር ፍራየርን", true, 3)
	add("cooking-appliances", "blenders-food-processors", "Heavy Duty Blenders & Food Processors", "ብሌንደሮች እና ፈጫጮች", "Bulandoota fi Miiksara", "ብሌንደራትን መፍሓኒታትን", true, 4)
	add("cooking-appliances", "water-dispensers-purifiers", "Hot & Cold Water Dispensers and Filters", "የውሃ ማጣሪያ እና ማሞቂያ/ማቀዝቀዣዎች", "Dispeensara Bishaanii", "ናይ ማይ መዘርግሒን መጻረዪን", true, 5)

	// =========================================================================
	// DOMAIN 10: Clothing, Fashion and Apparel (clothing-apparel)
	// =========================================================================
	add("clothing-apparel", "traditional-ethiopian-wear", "Traditional Ethiopian Attire", "ባህላዊ የኢትዮጵያ አልባሳት", "Uffata Aadaa Itoophiyaa", "ባህላዊ ናይ ኢትዮጵያ ክዳውንቲ", false, 1)
	add("traditional-ethiopian-wear", "habesha-kemis-dresses", "Women's Habesha Kemis (Tilf / Zuria)", "የሴቶች ሀበሻ ቀሚስ", "Qamisii Aadaa Dubartootaa", "ናይ ደቀንስትዮ ሓበሻ ቀሚስ", true, 1)
	add("traditional-ethiopian-wear", "mens-traditional-gabi-netela", "Men's Traditional Suits, Gabi and Netela", "የወንዶች ጋቢ፣ ነጠላ እና ባህላዊ ልብሶች", "Gaabii, Naxalaa fi Uffata Dhiiraa", "ናይ ደቂ ተባዕትዮ ጋቢን ነጠላን", true, 2)
	add("traditional-ethiopian-wear", "oromo-cultural-attire", "Oromo Cultural Outfits (Wandaboo / Qomee)", "የኦሮሞ ባህላዊ ልብሶች", "Uffata Aadaa Oromoo", "ናይ ኦሮሞ ባህላዊ ክዳውንቲ", true, 3)
	add("traditional-ethiopian-wear", "tigray-traditional-attire", "Tigray Cultural Dresses (Tilfi / Kidan Habesha)", "የትግራይ ባህላዊ ልብሶች", "Uffata Aadaa Tigraay", "ናይ ትግራይ ባህላዊ ቀሚስ (ጥልፊ)", true, 4)

	add("clothing-apparel", "mens-modern-clothing", "Men's Modern Fashion", "የወንዶች ዘመናዊ ፋሽን", "Uffata Dhiiraa Ammayyaa", "ናይ ደቂ ተባዕትዮ ዘመናዊ ክዳን", false, 2)
	add("mens-modern-clothing", "mens-suits-blazers", "Men's Formal Suits, Tuxedos & Blazers", "የወንዶች ሱፎች እና ጃኬቶች", "Kootii fi Suufii Dhiiraa", "ሱፋትን ጃኬታትን", true, 1)
	add("mens-modern-clothing", "mens-shirts-polos", "Button-Down Shirts, Polos and T-Shirts", "ሸሚዞች እና ቲሸርቶች", "Shaamiziifi Tiishartii", "ሸሚዛትን ቲሸርታትን", true, 2)
	add("mens-modern-clothing", "mens-jeans-trousers", "Denim Jeans, Chinos and Casual Trousers", "ሱሪዎች እና ጂንሶች", "Surree fi Jiinsii Dhiiraa", "ሱሪታትን ጂንስን", true, 3)
	add("mens-modern-clothing", "mens-jackets-hoodies", "Winter Jackets, Leather Coats and Hoodies", "ኮቶች እና ሹራቦች", "Jaakeetii fi Shuraabii", "ጃኬታትን ሹራባትን", true, 4)

	add("clothing-apparel", "womens-modern-clothing", "Women's Modern Fashion", "የሴቶች ዘመናዊ ፋሽን", "Uffata Dubartootaa", "ናይ ደቀንስትዮ ዘመናዊ ክዳን", false, 3)
	add("womens-modern-clothing", "womens-dresses-gowns", "Casual Dresses, Evening Gowns & Skirts", "ቀሚሶች እና ቀሚስ መሰል ልብሶች", "Qamisii Ammayyaa", "ዘመናዊ ቀሚሳት", true, 1)
	add("womens-modern-clothing", "womens-blouses-tops", "Blouses, Shirts and Tops", "ብላውዞች እና ሸሚዞች", "Bilaawuzii Dubartootaa", "ብላውዛትን ሸሚዛትን", true, 2)
	add("womens-modern-clothing", "womens-abayas-hijabs", "Modest Fashion, Abayas and Hijabs", "አባያዎች እና ሂጃቦች", "Hijaabii fi Abaayaa", "ኣባያታትን ሒጃባትን", true, 3)

	// =========================================================================
	// DOMAIN 18: Jobs, Careers and Recruitment (jobs-recruitment)
	// =========================================================================
	add("jobs-recruitment", "tech-software-jobs", "Technology and Software Engineering Jobs", "የቴክኖሎጂ እና የሶፍትዌር ስራዎች", "Hojiiwwan Teeknooloojii", "ናይ ቴክኖሎጂን ሶፍትዌርን ስራሕቲ", false, 1)
	add("tech-software-jobs", "software-developer-vacancies", "Software Developer & Engineer Positions", "የሶፍትዌር ኢንጂነሪንግ ክፍት ቦታዎች", "Banaa Hojii Oomisha Sooftiweerii", "ክፍቲ ናይ ሶፍትዌር መሃንድስ ቦታታት", true, 1)
	add("tech-software-jobs", "it-support-network-vacancies", "IT Support and Network Administration", "የኔትወርክ እና የአይቲ ድጋፍ ሰጪዎች", "Gargaaraa IT fi Netwoorkii", "ናይ IT ደገፍን መርበብ ሓበሬታን", true, 2)
	add("tech-software-jobs", "product-design-uiux-vacancies", "UI/UX & Product Design Roles", "የUI/UX ዲዛይን ክፍት ቦታዎች", "Banaa Hojii Dizaayinii UI/UX", "ክፍቲ ናይ UI/UX ዲዛይን ስራሕቲ", true, 3)

	add("jobs-recruitment", "engineering-construction-jobs", "Engineering and Construction Vacancies", "የምህንድስና እና የግንባታ ስራዎች", "Hojiiwwan Injinariingii fi Ijaarsaa", "ናይ ምህንድስናን ህንጻን ስራሕቲ", false, 2)
	add("engineering-construction-jobs", "civil-structural-jobs", "Civil & Structural Engineer Positions", "የሲቪል መሃንዲስ የስራ ማስታወቂያዎች", "Banaa Hojii Injinara Siiviilii", "ናይ ሲቪል መሃንድስ ክፍቲ ቦታታት", true, 1)
	add("engineering-construction-jobs", "site-supervisors-foremen", "Site Foremen and Construction Supervisors", "የሳይት ተቆጣጣሪዎች እና ፎርማኖች", "To'attoota Saayitii Ijaarsaa", "ናይ ሳይት ተቖጻጸርቲ", true, 2)
	add("engineering-construction-jobs", "surveyors-cad-draftsmen", "Land Surveyors and CAD Draftsmen", "የቅየሳ እና የCAD ረቂቅ ባለሙያዎች", "Ogeessa Kaartaa fi CAD", "ናይ ቅየሳን CAD ረቂቕን ሰብ ሞያ", true, 3)

	add("jobs-recruitment", "accounting-finance-jobs", "Accounting, Banking and Finance Jobs", "የሂሳብ፣ የባንክ እና የፋይናንስ ስራዎች", "Hojiiwwan Herregaa fi Baankii", "ናይ ሒሳብን ባንክን ፋይናንስን ስራሕቲ", false, 3)
	add("accounting-finance-jobs", "senior-accountant-vacancies", "Senior Accountant & Auditor Roles", "የአካውንታንት እና ኦዲተር ክፍት ቦታዎች", "Banaa Hojii Ogeessa Herregaa", "ናይ ሒሳብ ሓላፊ ክፍቲ ቦታታት", true, 1)
	add("accounting-finance-jobs", "bank-teller-cashier-jobs", "Bank Tellers and Cashiers", "የባንክ ቴለር እና ገንዘብ ያዥ ስራዎች", "Hojiiwwan Kaffaltii fi Baankii", "ናይ ባንክ ቴለርን ተቐባል ገንዘብን", true, 2)

	add("jobs-recruitment", "healthcare-medical-jobs", "Healthcare, Nursing and Medical Jobs", "የጤና እና የህክምና ስራዎች", "Hojiiwwan Fayyaa fi Yaalaa", "ናይ ጥዕናን ሕክምናን ስራሕቲ", false, 4)
	add("healthcare-medical-jobs", "clinical-nurses-midwives", "Clinical Nurses and Midwives", "የነርሲንግ እና አዋላጅ ክፍት ቦታዎች", "Banaa Hojii Narsii fi Deessiftuu", "ናይ ነርሲንግ ክፍቲ ቦታታት", true, 1)
	add("healthcare-medical-jobs", "medical-doctors-specialists", "Medical Doctors and Clinical Officers", "የሀኪሞች ክፍት ቦታዎች", "Banaa Hojii Hakiimotaa", "ናይ ሓካይም ክፍቲ ቦታታት", true, 2)
	add("healthcare-medical-jobs", "pharmacists-druggists", "Pharmacists and Druggists", "የፋርማሲ ባለሙያዎች", "Banaa Hojii Faarmaasistootaa", "ናይ ፋርማሲስት ክፍቲ ቦታታት", true, 3)

	// =========================================================================
	// DOMAIN 33: Property for Sale (property-sale)
	// =========================================================================
	add("property-sale", "residential-villas-houses-sale", "Residential Villas and Houses for Sale", "የሚሸጡ ቪላዎች እና የመኖሪያ ቤቶች", "Gurgurtaa Mana Jireenyaa fi Viillaa", "ዝሽየጡ ቪላታትን ኣባይትን", false, 1)
	add("residential-villas-houses-sale", "luxury-villas-sale", "Luxury Standalone Villas (Bole, CMC, Old Airport)", "የቅንጦት ቪላዎች", "Viillaa Qananii Gurguraa", "ናይ ቅንጦት ቪላታት", true, 1)
	add("residential-villas-houses-sale", "g-plus-residential-sale", "G+1, G+2 and G+3 Residential Houses for Sale", "ጂ+1፣ ጂ+2 እና ጂ+3 ቤቶች", "Mana G+1, G+2 fi G+3", "G+1, G+2, G+3 ኣባይቲ", true, 2)
	add("residential-villas-houses-sale", "traditional-service-houses-sale", "Compound Houses with Service Quarters", "ከነሰርቪሱ ያሉ የግቢ ቤቶች", "Mana Qubannaa fi Tajaajilaa", "ምስ ሰርቪሶም ዘለዉ ገዛውቲ", true, 3)

	add("property-sale", "apartments-condos-sale", "Apartments and Condominium Units for Sale", "የሚሸጡ አፓርታማዎች እና ኮንዶሚኒየሞች", "Gurgurtaa Kandoominiyeemii fi Apaartaamaa", "ዝሽየጡ ኣፓርታማታትን ኮንዶሚኒየማትን", false, 2)
	add("apartments-condos-sale", "real-estate-apartments-sale", "Private Real Estate Apartments (Finished & Semi-Finished)", "የሪል እስቴት አፓርታማዎች", "Apaartaamaa Dhaabbata Daldalaa", "ናይ ሪል እስቴት ኣፓርታማታት", true, 1)
	add("apartments-condos-sale", "condominium-studios-sale", "Addis Ababa Condominium Studio Units", "ስቱዲዮ ኮንዶሚኒየም ቤቶች", "Kandoominiyeemii Istudiyoo", "ስቱዲዮ ኮንዶሚኒየም", true, 2)
	add("apartments-condos-sale", "condominium-1-bed-sale", "1-Bedroom Condominium Units", "ባለ 1 መኝታ ኮንዶሚኒየም", "Kandoominiyeemii Kutaa 1", "ባዕል 1 መደቀሲ ኮንዶሚኒየም", true, 3)
	add("apartments-condos-sale", "condominium-2-bed-sale", "2-Bedroom Condominium Units", "ባለ 2 መኝታ ኮንዶሚኒየም", "Kandoominiyeemii Kutaa 2", "ባዕል 2 መደቀሲ ኮንዶሚኒየም", true, 4)
	add("apartments-condos-sale", "condominium-3-bed-sale", "3-Bedroom Condominium Units", "ባለ 3 መኝታ ኮንዶሚኒየም", "Kandoominiyeemii Kutaa 3", "ባዕል 3 መደቀሲ ኮንዶሚኒየም", true, 5)

	// =========================================================================
	// DOMAIN 37: Vehicles and Automobiles (vehicles-cars)
	// =========================================================================
	add("vehicles-cars", "passenger-sedans-compacts", "Sedans and Compact Hatchbacks", "የመንገደኛ ሴዳኖች እና ኮምፓክት መኪኖች", "Konkolaattota Seedaanii fi Haakpaakii", "ናይ መንገደኛ መካይን", false, 1)
	add("passenger-sedans-compacts", "compact-hatchbacks", "Compact City Hatchbacks (Vitz, Yaris, Swift, Alto)", "ትናንሽ ከተማ መኪኖች (ቪትዝ፣ ያሪስ)", "Konkolaattota Magaalaa Xixinnoo", "ንኣሽቱ ናይ ከተማ መካይን", true, 1)
	add("passenger-sedans-compacts", "economy-sedans", "Economy Sedans (Corolla, Dzire, Accent, Elantra)", "ሴዳን መኪኖች (ኮሮላ፣ ዲዛየር)", "Seedaanii Qusannaa", "ሴዳን መካይን (ኮሮላ፣ ዲዛየር)", true, 2)
	add("passenger-sedans-compacts", "executive-luxury-sedans", "Executive & Luxury Sedans (Mercedes C/E-Class, BMW 3/5-Series)", "የቅንጦት ሴዳኖች", "Seedaanii Qananii", "ናይ ቅንጦት ሴዳናት", true, 3)

	add("vehicles-cars", "suvs-four-wheel-drives", "SUVs, Crossovers and 4WD Vehicles", "አይነተ ብዙ 4WD መኪኖች እና ጂፖች", "Konkolaattota SUV fi 4WD", "SUVን 4WDን መካይን", false, 2)
	add("suvs-four-wheel-drives", "full-size-4wd-suvs", "Full-Size 4WD SUVs (Land Cruiser 300/200, Patrol)", "ሙሉ 4x4 ጂፖች (ላንድ ክሩዘር)", "Konkolaattota 4x4 Gurguddaa", "ሙሉእ 4x4 ጂፓት", true, 1)
	add("suvs-four-wheel-drives", "compact-crossovers", "Compact Crossovers (RAV4, Tucson, Sportage, Creta)", "መካከለኛ ጂፖች (ራቭ4፣ ቱክሰን)", "Crossover Giddu-galeessaa", "ማእከላይ ጂፓት (ራቭ4)", true, 2)
	add("suvs-four-wheel-drives", "prado-suvs", "Toyota Land Cruiser Prado Family", "የቶዮታ ፕራዶ ቤተሰብ", "Maatii Tooyootaa Piraadoo", "ቶዮታ ፕራዶ", true, 3)

	add("vehicles-cars", "electric-hybrid-vehicles", "Electric and Hybrid Vehicles", "የኤሌክትሪክ እና የሀይብሪድ መኪኖች", "Konkolaattota Elektiriikii fi Haayibriidii", "ናይ ኤሌክትሪክን ሃይብሪድን መካይን", false, 3)
	add("electric-hybrid-vehicles", "pure-electric-cars", "100% Battery Electric Vehicles (EVs)", "ሙሉ የኤሌክትሪክ መኪኖች (EV)", "Konkolaataa Elektiriikii Qulqulluu", "ናይ ባትሪ ኤሌክትሪክ መካይን", true, 1)
	add("electric-hybrid-vehicles", "petrol-hybrid-cars", "Gasoline-Electric Hybrid Cars (Corolla Cross, Prius)", "ሀይብሪድ መኪኖች", "Konkolaataa Haayibriidii", "ሃይብሪድ መካይን", true, 2)

	// =========================================================================
	// DOMAIN 40: Trucks, Buses and Commercial Transport (trucks-commercial-transport)
	// =========================================================================
	add("trucks-commercial-transport", "light-medium-commercial-trucks", "Light and Medium Cargo Trucks", "ቀላል እና መካከለኛ የጭነት መኪኖች", "Konkolaattota Fe'iinsaa Giddu-galeessaa", "ቀለልትን ማእከሎትን ናይ ጽዕነት መካይን", false, 1)
	add("light-medium-commercial-trucks", "isuzu-npr-nkr-light", "Isuzu NPR / NKR Light Cargo Trucks (3.5 - 5 Tons)", "አይሱዙ ኤንፒአር የጭነት መኪኖች", "Konkolaataa Isuzu NPR", "ኣይሱዙ NPR ናይ ጽዕነት መኪና", true, 1)
	add("light-medium-commercial-trucks", "isuzu-fsr-ftr-medium", "Isuzu FSR / Forward Medium Trucks (7 - 10 Tons)", "አይሱዙ ኤፍኤስአር (ፎርዋርድ)", "Konkolaataa Isuzu FSR", "ኣይሱዙ FSR (ፎርዋርድ)", true, 2)
	add("light-medium-commercial-trucks", "pickups-double-cab", "Double Cabin Pickups (Hilux, D-Max)", "ባለ ሁለት ረድፍ ፒክአፖች", "Piik-aappii Balbala Lammaa", "ድርብ ካቢን ፒካፕታት", true, 3)
	add("light-medium-commercial-trucks", "pickups-single-cab", "Single Cabin Work Pickups", "ባለ አንድ ረድፍ ፒክአፖች", "Piik-aappii Balbala Tokkoo", "ንጽል ካቢን ፒካፕታት", true, 4)

	add("trucks-commercial-transport", "heavy-duty-tippers-trailers", "Heavy Duty Dump Trucks and Trailers", "ከባድ የደረቅ ጭነት እና የቆሻሻ ማመላለሻዎች", "Konkolaataa Fe'iinsaa Gurguddaa", "ከበድቲ ናይ ጽዕነትን ደምፐራትን መካይን", false, 2)
	add("heavy-duty-tippers-trailers", "sino-howo-dump-trucks", "Sino Truck (HOWO) Tipper Dump Trucks", "ሲኖ ትራክ የቆሻሻ / አሸዋ ማመላለሻዎች", "Sinotruk Dampara", "ሲኖ ትራክ ደምፐር", true, 1)
	add("heavy-duty-tippers-trailers", "long-haul-trailers", "Long-Haul Prime Movers & Freight Trailers", "ረጅም ተሳቢ የጭነት መኪኖች", "Konkolaataa Harkisaa Dheeraa", "ናይ ርሑቕ መገዲ ተሳሓብቲ", true, 2)
	add("heavy-duty-tippers-trailers", "fuel-water-tankers", "Bulk Fuel and Water Tanker Trucks", "የነዳጅ እና የውሃ ቦቴዎች", "Konkolaataa Bootee Bishaanii fi Boba'aa", "ናይ ነዳድን ማይን ቦቴታት", true, 3)

	add("trucks-commercial-transport", "passenger-buses-vans", "Passenger Buses, Minibuses and Vans", "የህዝብ ማመላለሻ አውቶቡሶች እና ሚኒባሶች", "Baasota fi Miiniibaasii Fe'iinsa Namaa", "ናይ ህዝቢ ኣውቶቡሳትን ሚኒባሳትን", false, 3)
	add("passenger-buses-vans", "toyota-hiace-minibuses", "Toyota HiAce 12-15 Seater Minibuses", "ቶዮታ ሀይኤስ ሚኒባሶች", "Miiniibaasii Tooyootaa HiAce", "ቶዮታ ሃይኤስ ሚኒባስ", true, 1)
	add("passenger-buses-vans", "coaster-midibuses", "Toyota Coaster Midibuses (26 - 32 Seats)", "ኮስተር አውቶቡሶች", "Baasii Koostaraa", "ኮስተር ኣውቶቡስ", true, 2)
	add("passenger-buses-vans", "intercity-coach-buses", "Large Intercity High-Capacity Coach Buses", "የረጅም ርቀት የህዝብ አውቶቡሶች", "Baasii Gurguddaa Magaalaa Gidduu", "ናይ ነዊሕ መገዲ ኣውቶቡሳት", true, 3)

	// =========================================================================
	// DOMAIN 44: Crops, Grains, Seeds and Plants (crops-grains-seeds)
	// =========================================================================
	add("crops-grains-seeds", "teff-grain-varieties", "Ethiopian Teff Varieties", "የኢትዮጵያ የጤፍ ዝርያዎች", "Gosa Qamadii fi Xaafii", "ናይ ጣፍ ዓሌታት", false, 1)
	add("teff-grain-varieties", "magna-white-teff", "Magna (Super White) Premium Teff", "ማግና (ነጭ) ጤፍ", "Xaafii Adii Magnaa", "ማግና (ጻዕዳ) ጣፍ", true, 1)
	add("teff-grain-varieties", "nech-standard-teff", "Nech (Standard White) Teff", "ነጭ ጤፍ", "Xaafii Adii", "ጻዕዳ ጣፍ", true, 2)
	add("teff-grain-varieties", "sergegna-mixed-teff", "Sergegna (Mixed White & Brown) Teff", "ሰርገኛ ጤፍ", "Xaafii Sergeenyaa", "ሰርገኛ ጣፍ", true, 3)
	add("teff-grain-varieties", "key-brown-teff", "Key (Red/Brown) High-Iron Teff", "ቀይ ጤፍ", "Xaafii Diimaa", "ቀይሕ ጣፍ", true, 4)

	add("crops-grains-seeds", "cereals-wheat-maize", "Cereals, Wheat, Barley and Maize", "ስንዴ፣ ገብስ እና በቆሎ", "Qamadii, Garbuu fi Boqqolloo", "ስርናይ፣ ስገምን ዕፉንን", false, 2)
	add("cereals-wheat-maize", "durum-bread-wheat", "Durum & Bread Wheat (Bulk Grain)", "የዳቦ እና የፓስታ ስንዴ", "Qamadii Daabboo", "ናይ ባኒ ስርናይ", true, 1)
	add("cereals-wheat-maize", "malted-food-barley", "Food Barley and Malting Barley", "የምግብ እና የቢራ ገብስ", "Garbuu Nyaataa fi Biiraa", "ናይ መግቢ ስገም", true, 2)
	add("crops-grains-seeds", "white-yellow-maize", "White and Yellow Feed Maize", "ነጭ እና ቢጫ በቆሎ", "Boqqolloo Adii fi Keelloo", "ጻዕዳን ቢጫን ዕፉን", true, 3)
	add("crops-grains-seeds", "sorghum-millet", "Sorghum and Finger Millet (Dagussa)", "ማሽላ እና ዳጉሳ", "Mishingaa fi Daagussaa", "ማሽላን ዳጉሻን", true, 4)

	add("crops-grains-seeds", "oilseeds-pulses", "Oilseeds and Pulses (Export Quality)", "የቅባት እህሎች እና ጥራጥሬዎች", "Midhaan Zayitaa fi Nyaataa", "ናይ ዘይቲ እኽልን ጥራጥሬን", false, 3)
	add("oilseeds-pulses", "whitish-humera-sesame", "Whitish Humera / Wollega Sesame Seeds", "የሁመራ እና የወለጋ ሰሊጥ", "Saliixa Humaratti", "ናይ ሁመራ ሰሊጥ", true, 1)
	add("oilseeds-pulses", "kabuli-desi-chickpeas", "Kabuli & Desi Chickpeas (Shimbra)", "የካቡሊ እና የአገር ቤት ሽምብራ", "Shumburaa", "ሽምብራ", true, 2)
	add("oilseeds-pulses", "red-kidney-faba-beans", "Red Kidney Beans & Faba Broad Beans (Baqela)", "ቦሎቄ እና ባቄላ", "Baqelaa fi Booloqqee", "ዓተርን ባልደንጓን", true, 3)
	add("oilseeds-pulses", "niger-seed-noug", "Niger Seed (Noug) for Edible Oil", "የኑግ ፍሬ", "Nuugii", "ኑግ", true, 4)
	add("oilseeds-pulses", "soybeans-bulk", "Soybeans (Bulk Commercial Grain)", "አኩሪ አተር", "Akkurii Atara", "ኣኩሪ ዓተር", true, 5)

	// =========================================================================
	// DOMAIN 45: Livestock, Poultry and Animal Husbandry (livestock-poultry)
	// =========================================================================
	add("livestock-poultry", "cattle-dairy-beef", "Cattle, Dairy Cows and Beef Oxen", "የቁም ከብቶች፣ ላሞች እና ሰንጋዎች", "Loon, Saawwaa fi Sangaa", "ጥሪት፣ ላሞትን ሰንጋታትን", false, 1)
	add("cattle-dairy-beef", "fattened-beef-oxen", "Fattened Beef Oxen (Harar, Boran, Jimma Senga)", "የሰቡ ሰንጋዎች ለበዓል እና ለስጋ", "Sangaa Furdate", "ዝሰብሑ ሰንጋታት", true, 1)
	add("cattle-dairy-beef", "high-yield-dairy-cows", "Holstein Friesian & Crossbred Dairy Cows", "የወተት ላሞች (ሆልስታይን)", "Saawwaa Aannanii", "ናይ ጸባ ላሞት", true, 2)
	add("cattle-dairy-beef", "boran-heifers-calves", "Boran Breed Heifers and Young Calves", "የቦረና ዝርያ ጊደሮች እና ጥጆች", "Gidarii fi Jabbii Booranaa", "ናይ ቦረና ዓሌት ጣቓታት", true, 3)

	add("livestock-poultry", "sheep-goats", "Sheep and Goats (Shoats)", "በጎች እና ፍየሎች", "Hoolaa fi Re'ee", "ኣባጊዕን ኣጣልን", false, 2)
	add("sheep-goats", "harar-highland-sheep", "Harar & Arsi Fat-Tailed Highland Sheep", "የሀረር እና የአርሲ ወፍራም ጭራ በጎች", "Hoolaa Arsiifi Harar", "ናይ ሓረርን ኣርሲን ኣባጊዕ", true, 1)
	add("sheep-goats", "boer-somali-goats", "Boer Cross and Somali Meat Goats", "የሶማሌ እና የቦር ፍየሎች", "Re'ee Somaalee", "ናይ ሶማል ኣጣል", true, 2)

	add("livestock-poultry", "poultry-farming", "Poultry, Chickens and Hatchery Eggs", "የዶሮ እርባታ እና እንቁላል", "Lukkuu fi Hanqaaquu", "ናይ ደርሆ መራብሕን እንቋቑሖን", false, 3)
	add("poultry-farming", "day-old-chicks-layer-broiler", "Day-Old Chicks (Bovans Brown, Sasso, Cobb 500)", "የአንድ ቀን ጫጩቶች (ሳሶ፣ ቦቫንስ)", "Cuucii Guyyaa Tokkoo", "ናይ ሓደ መዓልቲ ጨጨውቲ", true, 1)
	add("poultry-farming", "ready-layer-hens", "Point-of-Lay Pullets & Egg Layers", "እንቁላል ጣይ ዶሮዎች", "Lukkuu Hanqaaquu Buuftu", "እንቋቑሖ ዝወድቓ ደራሁ", true, 2)
	add("poultry-farming", "local-habesha-chickens", "Traditional Habesha Free-Range Chickens (Kero)", "የሀበሻ የቤት ዶሮዎች", "Lukkuu Aadaa Habashaa", "ናይ ሓበሻ ደራሁ", true, 3)

	// =========================================================================
	// DOMAIN 49: Beverages, Coffee and Tea (coffee-tea-beverages)
	// =========================================================================
	add("coffee-tea-beverages", "specialty-ethiopian-coffee", "Specialty Green and Roasted Ethiopian Coffee", "የኢትዮጵያ የቡና ፍሬ እና የተቆላ ቡና", "Buna Itoophiyaa Addaa", "ፍሉይ ናይ ኢትዮጵያ ቡን", false, 1)
	add("specialty-ethiopian-coffee", "yirgacheffe-washed-coffee", "Yirgacheffe Grade 1 & 2 Washed Coffee Beans", "ይርጋጨፌ ታጥቦ የተዘጋጀ ቡና", "Buna Yirgaacaffee Dhiqame", "ይርጋጨፈ ዝተሓጸበ ቡን", true, 1)
	add("specialty-ethiopian-coffee", "sidama-specialty-coffee", "Sidama Grade 1 Specialty Coffee", "የሲዳማ ስፔሻሊቲ ቡና", "Buna Sidaamaa", "ናይ ሲዳማ ቡን", true, 2)
	add("specialty-ethiopian-coffee", "harar-longberry-natural", "Harar Longberry Natural Sundried Coffee", "የሀረር ሎንግበሪ የተፈጥሮ ቡና", "Buna Harar Aduun Goge", "ናይ ሓረር ቡን", true, 3)
	add("specialty-ethiopian-coffee", "jimma-limu-coffee", "Limu & Jimma Specialty Arabica", "የሊሙ እና የጅማ አረቢካ ቡና", "Buna Jimmaa fi Liimuu", "ናይ ጅማን ሊሙን ቡን", true, 4)
	add("specialty-ethiopian-coffee", "roasted-ground-coffee-retail", "Freshly Roasted & Ground Packaged Coffee", "የተቆላ እና የተፈጨ የታሸገ ቡና", "Buna Akaawwamee Daakame", "ዝተጠብሰን ዝተፈጨን ቡን", true, 5)

	return list
}
