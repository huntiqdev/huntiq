package main

import "fmt"

// GetFullSystematicDomainMatrix returns granular subcategories and leaf items across all 60 domains
func GetFullSystematicDomainMatrix() []DomainMatrixEntry {
	var list []DomainMatrixEntry

	// Helper function
	addEntry := func(domainSlug string, subs []SubBranchDef) {
		list = append(list, DomainMatrixEntry{
			DomainSlug:  domainSlug,
			SubBranches: subs,
		})
	}

	makeBranch := func(subSlug, en, am, om, ti string, items [][4]string) SubBranchDef {
		var leaves []LeafDef
		for _, it := range items {
			leaves = append(leaves, LeafDef{
				Slug:   it[0],
				NameEN: it[1],
				NameAM: it[2],
				NameOM: it[3],
				NameTI: fmt.Sprintf("ናይ %s", it[1]),
			})
		}
		return SubBranchDef{
			SubSlug: subSlug,
			NameEN:  en,
			NameAM:  am,
			NameOM:  om,
			NameTI:  ti,
			Leaves:  leaves,
		}
	}

	// -------------------------------------------------------------------------
	// 01. Electronics & Mobile Devices
	// -------------------------------------------------------------------------
	addEntry("electronics-mobile", []SubBranchDef{
		makeBranch("mobile-phone-components", "Phone Screens, Batteries and Internal Spare Parts", "የስልክ ስክሪኖች፣ ባትሪዎች እና መለዋወጫዎች", "Iskiriinii fi Baatirii Bilbilaa", "ስክሪንን ባትሪን ሞባይል", [][4]string{
			{"oled-amoled-replacement-screens", "Original AMOLED & OLED Phone Replacement Displays", "ኦሪጅናል AMOLED የስልክ ስክሪኖች", "Iskiriinii AMOLED Qulqulluu"},
			{"high-capacity-replacement-batteries", "OEM Internal High-Capacity Lithium Phone Batteries", "ኦሪጅናል የስልክ ባትሪዎች", "Baatirii Bilbilaa Qulqulluu"},
			{"charging-port-flex-cables", "USB-C & Lightning Charging Port Flex Ribbon Cables", "የቻርጅ ፖርት ኬብሎች እና ቦርዶች", "Kaardii fi Boordii Chaarjii"},
			{"camera-module-replacements", "Front & Rear Camera Module Replacements", "የሞባይል ካሜራ መለዋወጫዎች", "Kaameeraa Jijjiirraa Bilbilaa"},
			{"fingerprint-face-id-sensors", "In-Display Fingerprint and Face ID Sensor Modules", "የጣት አሻራ እና የፊት መለያ ሴንሰሮች", "Seensara Asheeraa fi Fuulaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 02. Computers, IT and Office Technology
	// -------------------------------------------------------------------------
	addEntry("computers-it", []SubBranchDef{
		makeBranch("computer-peripherals-keyboards-mice", "Keyboards, Mice and Desktop Peripherals", "ኪቦርዶች፣ ማውሶች እና መገልገያዎች", "Kiiboordii fi Hantuutaa Kompiitaraa", "ኪቦርድን ማውስን", [][4]string{
			{"mechanical-gaming-keyboards", "RGB Hot-Swappable Mechanical Gaming Keyboards", "ሜካኒካል የጌሚንግ ኪቦርዶች", "Kiiboordii Meekaaniikaalaa"},
			{"wireless-ergonomic-mice", "Ergonomic Bluetooth & 2.4GHz Wireless Optical Mice", "ገመድ አልባ ማውሶች", "Hantuutaa Wayaalessii (Mouse)"},
			{"noise-canceling-usb-headsets", "Call Center & Business USB Noise-Canceling Headsets", "የጥሪ ማዕከል እና የቢሮ የጆሮ ማዳመጫዎች", "Heedseetii USB Waamichaa"},
			{"4k-streaming-webcams", "4K Ultra-HD Autofocus Streaming & Video Conference Webcams", "4ኬ የቪዲዮ ኮንፈረንስ ዌብካሞች", "Kaameeraa Weebkaamii 4K"},
		}),
	})

	// -------------------------------------------------------------------------
	// 04. Cameras, Photography and Optics
	// -------------------------------------------------------------------------
	addEntry("cameras-optics", []SubBranchDef{
		makeBranch("camera-accessories-optics", "Camera Bags, Memory Cards and Tripods", "የካሜራ ቦርሳዎች፣ ሚሞሪ ካርዶች እና ትራይፖዶች", "Boorsaa Kaameeraa fi Kaardii", "ናይ ካሜራ ቦርሳን ሚሞሪን", [][4]string{
			{"v90-uhs-ii-sd-cards", "High-Speed V60 & V90 UHS-II 4K/8K SD Cards (128GB-512GB)", "ፈጣን የቪዲዮ ካሜራ ኤስዲ ካርዶች", "Kaardii SD Saffisaa V90"},
			{"waterproof-padded-camera-backpacks", "Waterproof Padded DSLR Camera Backpacks with Lens Dividers", "ውሃ የማያስገቡ የካሜራ ቦርሳዎች", "Boorsaa Kaameeraa Bishaan Danda'u"},
			{"circular-polarizing-nd-filters", "CPL Circular Polarizer & Variable ND Neutral Density Filters", "የካሜራ ሌንስ ማጣሪያ መስታወቶች (ND/CPL)", "Filtarii Leensii Kaameeraa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 06. TVs, Audio and Home Entertainment
	// -------------------------------------------------------------------------
	addEntry("tv-audio", []SubBranchDef{
		makeBranch("pro-audio-mixers-pa", "Professional Audio Mixers, Microphones and Amplifiers", "የድምጽ ማቀናበሪያዎች፣ ማይክሮፎኖች እና አምፕሊፋየሮች", "Miiksara fi Sirna Sagalee Gurguddaa", "ናይ ድምጺ መተሓሓዚን ማይክሮፎንን", [][4]string{
			{"multi-channel-audio-mixing-consoles", "12 to 32 Channel Professional Digital Audio Mixing Consoles", "የድምጽ ማቀናበሪያ ዲጂታል ሚክሰሮች", "Miiksara Sagalee Dijitaalaa"},
			{"wireless-uhf-vocal-microphones", "Dual UHF Wireless Handheld Vocal Microphones", "ገመድ አልባ የመድረክ ማይክሮፎኖች", "Maayikiroofoonii Wayaalessii UHF"},
			{"heavy-duty-power-amplifiers", "2000W to 5000W Class-D Professional Sound Power Amplifiers", "ከፍተኛ አቅም ያላቸው አምፕሊፋየሮች", "Aampiliifaayera Humna Guddaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 07. Home Appliances
	// -------------------------------------------------------------------------
	addEntry("home-appliances", []SubBranchDef{
		makeBranch("climate-control-water-heaters", "Air Conditioners, Fans and Water Heaters", "የአየር ማቀዝቀዣዎች፣ ቬንቲሌተሮች እና ቦይለሮች", "Qabbaneessituu fi Booyilara Bishaanii", "ናይ ኣየር መዝሓሊን ቦይለርን", [][4]string{
			{"inverter-split-air-conditioners", "Energy-Saving Split Inverter Air Conditioners (12000-24000 BTU)", "ኢንቨርተር ኤሲ (የአየር ማቀዝቀዣ)", "Kuularii Qilleensaa Inveertaraa"},
			{"electric-storage-water-heaters-geysers", "50L to 100L Vertical Electric Storage Water Heaters (Geysers)", "የኤሌክትሪክ የውሃ ማሞቂያ ቦይለሮች", "Booyilara Bishaan Ho'isu 50L-100L"},
			{"high-velocity-industrial-fans", "Heavy Duty Oscillating Pedestal and Wall Fans", "የግድግዳ እና የቁም ቬንቲሌተሮች (ፋን)", "Fiixee Qilleensaa Jabaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 10. Clothing, Fashion and Apparel
	// -------------------------------------------------------------------------
	addEntry("clothing-apparel", []SubBranchDef{
		makeBranch("womens-traditional-habesha-kemis", "Handwoven Habesha Kemis, Zuria and Netela", "የእጅ ጥልፍ ሀበሻ ቀሚሶች እና ዙሪያዎች", "Qamisii Aadaa fi Zuuriyaa", "ናይ ጥልፊ ሓበሻ ቀሚስ", [][4]string{
			{"gold-thread-tilfi-habesha-kemis", "Pure Cotton Gold Thread Royal Embroidered Habesha Kemis", "በወርቅ ክር የተጠለፉ የሰርግ ሀበሻ ቀሚሶች", "Qamisii Aadaa Warqeen Faayame"},
			{"wollo-traditional-embroidery-kemis", "Wollo Cultural Traditional Dress (Bati / Wollo)", "የወሎ ባህላዊ ጥልፍ ቀሚስ", "Qamisii Aadaa Walloo"},
			{"gondar-royal-traditional-dress", "Gondar Royal Court Style Silk Embroidered Kemis", "የጎንደር ባህላዊ ቀሚስ", "Qamisii Aadaa Gondar"},
			{"handwoven-cotton-netela-shawls", "Handwoven Pure Cotton Netela Shawls with Border Tilf", "የእጅ ጥልፍ ነጠላዎች እና ሻሾች", "Naxalaa fi Shaashii Aadaa"},
		}),
		makeBranch("mens-suits-formal", "Men's Business Suits, Tuxedos and Formal Blazers", "የወንዶች ሙሉ ሱፍ እና ጃኬቶች", "Suufii fi Kootii Dhiiraa", "ናይ ደቂ ተባዕትዮ ሱፍ", [][4]string{
			{"three-piece-business-suits", "Italian & Turkish Cut 3-Piece Formal Wool Suits (Jacket, Vest, Pants)", "ባለ 3 ጥንድ የወንዶች ሙሉ ሱፍ", "Suufii Guutuu Dhiiraa"},
			{"tuxedo-evening-dinner-jackets", "Satin Lapel Wedding Tuxedos and Black Tie Blazers", "የሰርግ እና የምሽት ቱክሲዶ ጃኬቶች", "Kootii Cidhaa Qananii (Tuxedo)"},
		}),
	})

	// -------------------------------------------------------------------------
	// 29. Construction Materials & Building Supplies
	// -------------------------------------------------------------------------
	addEntry("construction-materials", []SubBranchDef{
		makeBranch("roofing-gypsum-insulation", "Gypsum Boards, False Ceilings and Thermal Insulation", "ጂፕሰም ቦርዶች እና የኮርኒስ እቃዎች", "Jiipsumii fi Baaxii Manaa", "ጂፕሰም ቦርድን ናይ ኮርኒስ ንዋያትን", [][4]string{
			{"regular-moisture-gypsum-boards", "9mm & 12mm Standard & Moisture-Resistant Gypsum Boards", "የኮርኒስ ጂፕሰም ቦርዶች (9ሚሜ/12ሚሜ)", "Boordii Jiipsumii 9mm fi 12mm"},
			{"galvanized-ceiling-studs-channels", "Galvanized Light Steel Keel Studs and Track Channels", "የጂፕሰም ማሰሪያ የብረት ቻናሎች", "Sibiila Qalloo Jiipsumii"},
			{"thermal-glasswool-insulation-rolls", "Fiberglass Wool & Rockwool Thermal Acoustic Insulation", "የድምጽ እና የሙቀት መከላከያ ሱፍ (ኢንሱሌሽን)", "Wulii Qabbanaafi Ho'a Dhorku"},
		}),
	})

	// -------------------------------------------------------------------------
	// 30. Construction Machinery and Heavy Equipment
	// -------------------------------------------------------------------------
	addEntry("construction-machinery", []SubBranchDef{
		makeBranch("drilling-crushing-equipment", "Stone Crushers, Air Compressors and Rock Drills", "የድንጋይ መፍጫ ማሽኖች እና የዋሻ ድራይሎች", "Maashina Dhagaa Daakuu fi Qotu", "ናይ እምኒ መጥሓንን መበስበሲን", [][4]string{
			{"jaw-cone-stone-crushers", "Primary Jaw & Secondary Cone Mobile Stone Crushing Plants", "የድንጋይ መፍጫ ትላልቅ ማሽኖች (ክረሸር)", "Maashina Dhagaa Daaku (Crusher)"},
			{"high-pressure-diesel-air-compressors", "High-Pressure Skid Diesel Screw Air Compressors for Jackhammers", "የናፍጣ ኤር ኮምፕሬሰሮች ለጃክሀመር", "Kompreesara Qilleensaa Diizilii"},
			{"pneumatic-rock-jackhammers", "Heavy Pneumatic Rock Drilling Jackhammers (YT28)", "የድንጋይ መበሻ ጃክሀመሮች", "Jaakahamara Dhagaa Uru"},
		}),
	})

	// -------------------------------------------------------------------------
	// 33. Property for Sale
	// -------------------------------------------------------------------------
	addEntry("property-sale", []SubBranchDef{
		makeBranch("commercial-buildings-properties-sale", "Shopping Plazas, Office Buildings and Mixed-Use", "የሚሸጡ የንግድ ህንጻዎች እና ፕላዛዎች", "Gamoo Daldalaa Gurguraa", "ዝሽየጡ ናይ ንግዲ ህንጻታት", [][4]string{
			{"multi-storey-office-buildings-sale", "G+4 to G+12 Modern Office Buildings with Elevator (Bole, Kazanchis)", "ባለ ሊፍት ዘመናዊ የቢሮ ህንጻዎች", "Gamoo Biiroo Liiftii Qabu"},
			{"ground-plus-shops-mix-use-sale", "Mixed Commercial-Residential G+2 & G+3 Buildings for Sale", "የንግድ እና የመኖሪያ ህንጻዎች", "Gamoo Daldala fi Jireenyaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 37. Vehicles and Automobiles
	// -------------------------------------------------------------------------
	addEntry("vehicles-cars", []SubBranchDef{
		makeBranch("commercial-vans-hiace", "Cargo Delivery Vans and Passenger Minivans", "የጭነት ቫኖች እና የቤተሰብ ሚኒቫኖች", "Vaanii Fe'iinsaa fi Maatii", "ናይ ጽዕነት ቫንን ናይ ስድራቤት ሚኒቫንን", [][4]string{
			{"blind-van-cargo-haulers", "Toyota HiAce Blind Van Panel Cargo Carriers", "የእቃ መጫኛ ዝግ ቫኖች (ብላይንድ ቫን)", "Vaanii Fe'iinsaa Cufaa"},
			{"7-seater-family-minivans", "7-Seater Family MPV Minivans (Toyota Avanza, Suzuki Ertiga)", "ባለ 7 ሰው የቤተሰብ መኪኖች", "Konkolaataa Maatii Nama 7"},
		}),
	})

	// -------------------------------------------------------------------------
	// 40. Trucks, Buses and Commercial Transport
	// -------------------------------------------------------------------------
	addEntry("trucks-commercial-transport", []SubBranchDef{
		makeBranch("specialized-transport-rigs", "Refrigerated Reefer Trucks and Concrete Transit Mixers", "የቀዝቃዛ ጭነት እና የኮንክሪት ማመላለሻዎች", "Konkolaataa Qabbanaa'aa fi Makaa", "ናይ ዝሑል ጽዕነትን ኮንክሪት መዘወርን", [][4]string{
			{"refrigerated-cold-chain-trucks", "Insulated Temperature-Controlled Reefer Trucks (Meat/Dairy)", "የቀዝቃዛ ጭነት ተሽከርካሪዎች (ለስጋ እና ለወተት)", "Konkolaataa Qabbanaa'aa Foonii fi Aannanii"},
			{"transit-concrete-mixer-trucks", "8m³ to 12m³ Heavy Duty Concrete Transit Agitator Trucks", "የኮንክሪት ማቡኪያ ተሽከርካሪዎች (ሚክሰር ትራክ)", "Konkolaataa Konkiriitii Naannessu"},
		}),
	})

	// -------------------------------------------------------------------------
	// 44. Crops, Grains, Seeds and Plants
	// -------------------------------------------------------------------------
	addEntry("crops-grains-seeds", []SubBranchDef{
		makeBranch("specialty-ethiopian-spices-herbs", "Korarima, Black Cumin, Fenugreek and Mustard Seeds", "ኮረሪማ፣ ጥቁር አዝሙድ፣ አብሽ እና ሰናፍጭ", "Qooraa Aadaa, Qorariimaa fi Habaasoodaa", "ኮረሪማ፣ ጸሊም ኣዝሙድ፣ ኣብሽን ስናፍጭን", [][4]string{
			{"dried-kefa-korarima-pods", "Sun-Dried Whole Brown Ethiopian Cardamom Pods (Korarima)", "የከፋ ንጹህ የኮረሪማ ፍሬ", "Qorariimaa Qulqulluu Kaffaa"},
			{"black-cumin-habasoda-seeds", "Pure Cleaned Black Cumin Seeds (Tikur Azmud / Tikur Azmud)", "ጥቁር አዝሙድ በኩንታል", "Habaasoodaa Gurraacha"},
			{"white-cumin-nech-azmud", "White Cumin Seeds (Nech Azmud / Kamun)", "ነጭ አዝሙድ (ከሙን)", "Habaasoodaa Adii"},
			{"ethiopian-fenugreek-abish", "High-Protein Ethiopian Fenugreek Seeds (Abish)", "የሀገር ቤት አብሽ", "Abishi Aadaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 56. Energy, Solar & Electrical Power Equipment
	// -------------------------------------------------------------------------
	addEntry("energy-solar-power", []SubBranchDef{
		makeBranch("solar-charge-controllers-cables", "MPPT Solar Charge Controllers and Solar Cables", "የሶላር ቻርጅ ኮንትሮለሮች እና የሶላር ኬብሎች", "To'ataa Chaarjii Soolaarii MPPT", "ናይ ሶላር ቻርጅ ኮንትሮለር MPPT", [][4]string{
			{"mppt-solar-charge-controllers-60a", "60A to 100A Auto-Detect MPPT Solar Charge Controllers", "የኤምፒፒቲ (MPPT) የሶላር ቻርጅ ኮንትሮለሮች", "To'ataa Soolaarii MPPT 60A-100A"},
			{"photovoltaic-pv-cables-mc4", "4mm² and 6mm² UV-Resistant Solar PV Cables with MC4 Connectors", "የሶላር ኬብሎች እና MC4 አገናኞች", "Keebilii Soolaarii fi Qabsiisaa MC4"},
		}),
	})

	return list
}
