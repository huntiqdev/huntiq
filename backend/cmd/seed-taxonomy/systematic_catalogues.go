package main

// GetDomainSystematicCatalogues generates authentic, structured subsectors and leaves across all 60 domains
func GetDomainSystematicCatalogues() []DomainSystematicCatalog {
	var catalogues []DomainSystematicCatalog

	// Helper to add a domain catalogue
	addCat := func(domainSlug string, subs []SubSectorTemplate) {
		catalogues = append(catalogues, DomainSystematicCatalog{
			DomainSlug: domainSlug,
			SubSectors: subs,
		})
	}

	makeSub := func(slug, en, am, om, ti string, items [][4]string) SubSectorTemplate {
		return SubSectorTemplate{
			SubSlug: slug,
			NameEN:  en,
			NameAM:  am,
			NameOM:  om,
			NameTI:  ti,
			Leaves:  items,
		}
	}

	// -------------------------------------------------------------------------
	// 01. Electronics and Mobile Devices
	// -------------------------------------------------------------------------
	addCat("electronics-mobile", []SubSectorTemplate{
		makeSub("screen-protectors-finishes", "Screen Protectors and Film Coatings", "የስክሪን መከላከያ መስታወቶች", "Eegduu Iskiriinii", "ናይ ስክሪን መከላኸሊ", [][4]string{
			{"9h-tempered-glass-clear", "9H Clear HD Tempered Glass Screen Protectors", "ኦሪጅናል ጥርት ያለ የቴምፐርድ ስክሪን መከላከያ", "Iskiriinii Giilaasii 9H"},
			{"matte-anti-glare-glass", "Matte Gaming Anti-Glare Fingerprint-Resistant Glass", "የማይዝ ድምፅ አልባ የጌሚንግ ስክሪን መከላከያ", "Iskiriinii Maatii Geemiingii"},
			{"anti-peep-privacy-glass", "Anti-Spy 28-Degree Privacy Tempered Glass", "ፕራይቬሲ ከጎን የማያሳይ የስክሪን መከላከያ", "Iskiriinii Iccitii Eegu"},
			{"flexible-hydrogel-films", "Self-Healing Flexible Hydrogel Screen Protectors", "ተጣጣፊ የሀይድሮጀል ስክሪን ፊልሞች", "Fiilmii Haayidiroojeelii"},
			{"camera-lens-protector-rings", "Aluminum Alloy Camera Lens Protector Rings", "የካሜራ ሌንስ መከላከያ ቀለበቶች", "Qubee Eegduu Leensii"},
		}),
		makeSub("cases-protective-armor", "Phone Cases, Bumper Covers and Wallets", "የስልክ መከላከያ ከቨሮች እና ቦርሳዎች", "Koversii fi Eegduu Bilbilaa", "ናይ ሞባይል መከላኸሊ ከቨር", [][4]string{
			{"shockproof-military-armor-cases", "Dual-Layer Shockproof Heavy Duty Military Armor Cases", "የወታደራዊ ጠንካራ የስልክ ከቨሮች", "Koversii Jabaa Waraanaa"},
			{"magnetic-magsafe-clear-cases", "Crystal Clear Magnetic MagSafe Compatible Cases", "ማግሴፍ የሚሰሩ ጥርት ያሉ ከቨሮች", "Koversii Qulqulluu MagSafe"},
			{"genuine-leather-wallet-flip-cases", "Vintage Genuine Leather Wallet Flip Folio Cases", "የቆዳ የኪስ ቦርሳ መሰል የስልክ ከቨሮች", "Koversii Gogaa Kaardii Qabu"},
			{"liquid-silicone-soft-cases", "Slim Soft Touch Liquid Silicone Protective Cases", "ቀጫጭን ለስላሳ የሲሊኮን ከቨሮች", "Koversii Siilikonii Salphaa"},
			{"waterproof-underwater-pouches", "IPX8 Certified Floating Waterproof Phone Pouches", "የውሃ ውስጥ የስልክ መያዣ ኪሶች", "Boorsaa Bishaan Dhorku"},
		}),
		makeSub("cables-chargers-adapters", "Charging Cables, Adapters and Wireless Pads", "የቻርጅ ገመዶች፣ አዳፕተሮች እና ገመድ አልባ ቻርጀሮች", "Keebilii fi Chaarcharoota", "ናይ ቻርጅ ኬብላትን ኣዳፕተራትን", [][4]string{
			{"60w-100w-type-c-cables", "Braided 60W to 100W PD Fast Charging Type-C Cables", "ፈጣን የUSB-C ገመዶች (100 ዋት)", "Keebilii Type-C 100W"},
			{"lightning-to-usb-c-cables", "MFi-Certified USB-C to Lightning Fast Charge Cables", "የአይፎን የUSB-C ቻርጅ ገመዶች", "Keebilii Lightning Ayifoonii"},
			{"multi-device-charging-stations", "3-in-1 Foldable Wireless Charging Docks (Phone/Watch/Buds)", "የ3-በ-1 ገመድ አልባ ቻርጅ ማድረጊያ ዶኮች", "Chaarchara Wayaalessii 3-in-1"},
			{"car-fast-chargers-usb-pd", "Dual Port USB-C QC3.0 Fast Cigarette Lighter Car Chargers", "የመኪና ውስጥ ፈጣን ቻርጀሮች", "Chaarchara Konkolaataa Saffisaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 02. Computers, IT and Office Technology
	// -------------------------------------------------------------------------
	addCat("computers-it", []SubSectorTemplate{
		makeSub("pc-cooling-power-supplies", "Power Supplies (PSU), Case Fans and Liquid Coolers", "የኮምፒውተር ፓወር ሰፕላይ እና ማቀዝቀዣዎች", "Paawarii fi Qabbaneessituu Kompiitaraa", "ናይ ኮምፒዩተር ፓወርን መዝሓሊን", [][4]string{
			{"modular-gold-power-supplies", "80-Plus Gold Fully Modular Power Supplies (650W-1000W)", "ሞዱላር ፓወር ሰፕላይ (80 Plus Gold)", "Paawarii Kompiitaraa 80 Plus"},
			{"aio-liquid-cpu-coolers", "240mm and 360mm ARGB AIO Liquid CPU Water Coolers", "የውሃ ፈሳሽ የሲፒዩ ማቀዝቀዣዎች (AIO)", "Qabbaneessituu Bishaanii CPU"},
			{"tower-air-cpu-heatsinks", "Dual-Tower Heatpipe Air CPU Coolers with 120mm Fans", "የትልቅ ታወር የአየር ሲፒዩ ማቀዝቀዣዎች", "Qabbaneessituu Qilleensaa CPU"},
			{"argb-chassis-case-fans", "120mm PWM Addressable ARGB Hydraulic Bearing Case Fans", "የኮምፒውተር ሳጥን የኤልኢዲ ቬንቲሌተሮች", "Fiixee Saanduqa Kompiitaraa"},
		}),
		makeSub("networking-adapters-drives", "Wi-Fi Cards, Bluetooth Dongles and USB Enclosures", "የዋይፋይ ካርዶች፣ ብሉቱዝ እና የኤስኤስዲ ማስቀመጫዎች", "Kaardii Wi-Fi fi Saanduqa SSD", "ናይ ዋይፋይ ካርድን ናይ SSD መቐመጢን", [][4]string{
			{"pcie-wifi-6e-desktop-cards", "PCIe Wi-Fi 6E & Bluetooth 5.3 Desktop Expansion Cards", "የኮምፒውተር የዋይፋይ 6 ካርዶች", "Kaardii Wi-Fi 6E Kompiitaraa"},
			{"usb-nvme-external-enclosures", "10Gbps Aluminum M.2 NVMe SSD External Enclosures", "የኤም.2 ኤስኤስዲ ውጫዊ ማስቀመጫ ሳጥኖች", "Saanduqa NVMe SSD Dabalataa"},
			{"gigabit-usb-ethernet-adapters", "USB 3.0 to RJ45 Gigabit Ethernet Network Adapters", "የዩኤስቢ ወደ ኔትወርክ ኬብል መቀየሪያዎች", "Adaaptara USB gara Netwoorkii"},
		}),
	})

	// -------------------------------------------------------------------------
	// 08. Furniture and Interior Furnishings
	// -------------------------------------------------------------------------
	addCat("furniture-furnishings", []SubSectorTemplate{
		makeSub("study-kids-furniture", "Study Desks, Bookshelves and Bunk Beds", "የጥናት ጠረጴዛዎች፣ የመጽሐፍ መደርደሪያዎች እና ደርብ አልጋዎች", "Minjaala Qorannoo fi Siree Dabalataa", "ናይ መጽናዕቲ ጠረጴዛን ደርቢ ዓራትን", [][4]string{
			{"kids-wooden-bunk-beds", "Solid Hardwood Twin-over-Twin Bunk Beds with Ladders", "የልጆች ባለ ሁለት ደርብ የእንጨት አልጋዎች", "Siree Daa'immanii Lammaa Wal-irraa"},
			{"adjustable-study-desks-chairs", "Ergonomic Height-Adjustable Kids & Student Study Desks", "የተማሪዎች የጥናት ጠረጴዛዎች እና ወንበሮች", "Minjaala fi Barcuma Barattootaa"},
			{"multi-tier-wooden-bookshelves", "5-Tier Open Display Wooden Bookshelves and Library Units", "ባለ 5 ደረጃ የመጽሐፍ መደርደሪያዎች", "Madaallii Kitaabaa Mukaa"},
			{"compact-shoe-storage-cabinets", "Multi-Tier Dustproof Wooden & Metal Shoe Storage Racks", "የጫማ መደርደሪያ ካቢኔቶች", "Kaabineetii Kophee Kaa'an"},
		}),
	})

	// -------------------------------------------------------------------------
	// 09. Home, Kitchen and Household Goods
	// -------------------------------------------------------------------------
	addCat("kitchen-household", []SubSectorTemplate{
		makeSub("tableware-cutlery-glassware", "Porcelain Dinner Sets, Cutlery and Drinking Glasses", "የገበታ ሳህኖች፣ ማንኪያዎች እና የውሃ ብርጭቆዎች", "Saanii, Fal'aana fi Burcuqqoo", "ናይ መኣዲ ሳሕንታትን ማንካታትን ብርጭቆን", [][4]string{
			{"porcelain-dinnerware-sets", "32-Piece Fine Porcelain Ceramic Dinner Sets", "የሴራሚክ የገበታ ሳህን ስብስቦች", "Saanii Nyaataa Seeraamikii"},
			{"stainless-steel-cutlery-sets", "Stainless Steel Fork, Spoon and Steak Knife Cutlery Kits", "የማንኪያ፣ የሹካ እና የቢላ ስብስቦች", "Fal'aana fi Shookaa Sibiilaa"},
			{"crystal-drinking-water-glasses", "Lead-Free Crystal Drinking Glass Tumblers (Set of 6)", "የውሃ እና የጭማቂ ክሪስታል ብርጭቆዎች", "Burcuqqoo Bishaan fi Cuunfaa"},
			{"insulated-thermal-food-flasks", "Double-Wall Insulated Stainless Steel Food Flasks & Thermos", "ምግብ እና ሻይ የሚያሞቁ ቴርሞሶች", "Termoosii Nyaataa fi Shaayii"},
		}),
	})

	// -------------------------------------------------------------------------
	// 10. Clothing, Fashion and Apparel
	// -------------------------------------------------------------------------
	addCat("clothing-apparel", []SubSectorTemplate{
		makeSub("mens-casual-shirts-jeans", "Casual Shirts, Denim Jeans and Chinos", "የወንዶች ተራ ሸሚዞች፣ ጂንሶች እና ሱሪዎች", "Shaamizii fi Jiinsii Dhiiraa", "ናይ ደቂ ተባዕትዮ ሸሚዝን ጂንስን", [][4]string{
			{"linen-long-sleeve-casual-shirts", "100% Pure Linen Breathable Long Sleeve Casual Shirts", "የተልባ እግር (ሊነን) የወንዶች ሸሚዞች", "Shaamizii Liinenii Dhiiraa"},
			{"stretch-slim-fit-denim-jeans", "Comfort Stretch Cotton Slim Fit Denim Jeans", "የወንዶች ስሊም ፊት ጂንስ ሱሪዎች", "Jiinsii Dhiiraa Qalloo"},
			{"cotton-cargo-pants-multiple-pockets", "Heavy Duty Multiple-Pocket Outdoor Cargo Pants", "ባለ ብዙ ኪስ የካርጎ ሱሪዎች", "Surree Kaargoo Kiisii Hedduu"},
			{"cotton-pique-polo-shirts", "Classic Pique Knit Cotton Short-Sleeve Polo Shirts", "የፖሎ ቲሸርቶች", "Tiishartii Pooloo Dhiiraa"},
		}),
		makeSub("womens-dresses-skirts-modest", "Dresses, Midi Skirts, Cardigans and Abayas", "ቀሚሶች፣ ስከርት፣ ሹራቦች እና አባያዎች", "Qamisii, Gaabii fi Abaayaa", "ቀሚሳትን ኣባያታትን", [][4]string{
			{"chiffon-floral-maxi-dresses", "Elegant Floral Print Chiffon Long Maxi Dresses", "የቺፎን ረጅም ቀሚሶች", "Qamisii Dheeraa Shiifoonii"},
			{"high-waisted-pleated-midi-skirts", "High-Waisted A-Line Pleated Elegant Midi Skirts", "የተሸበሸቡ የሴቶች ስከርቶች", "Iskartii Dubartootaa"},
			{"luxury-dubai-nidha-abayas", "Imported Dubai Nidha Fabric Embroidered Open & Closed Abayas", "የዱባይ ጥልፍ አባያዎች", "Abaayaa Dubaay Faayame"},
			{"soft-knit-winter-cardigans", "Warm Button-Up Soft Knit Longline Winter Cardigans", "ረዣዥም የሴቶች ሹራቦች (ካርዲጋን)", "Shuraabii Dubartootaa Dheeraa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 11. Shoes, Bags and Accessories
	// -------------------------------------------------------------------------
	addCat("shoes-bags", []SubSectorTemplate{
		makeSub("leather-belts-wallets", "Genuine Leather Belts, Wallets and Keychains", "የቆዳ ቀበቶዎች፣ የኪስ ቦርሳዎች እና ቁልፍ መያዣዎች", "Qabattoo Gogaa fi Boorsaa Kiisii", "ናይ ቆርበት ቐበቶን ናይ ጁባ ቦርሳን", [][4]string{
			{"full-grain-mens-leather-belts", "100% Ethiopian Full-Grain Leather Formal Belts with Metal Buckle", "የወንዶች ንጹህ የቆዳ ቀበቶዎች", "Qabattoo Gogaa Qulqulluu Dhiiraa"},
			{"rfid-blocking-leather-bifold-wallets", "RFID-Blocking Slim Genuine Leather Bifold Wallets", "የወንዶች የቆዳ የኪስ ቦርሳዎች (RFID)", "Boorsaa Kiisii Gogaa (Wallet)"},
			{"womens-leather-zip-around-wallets", "Long Zip-Around Multi-Card Leather Clutch Wallets", "የሴቶች ረጅም የቆዳ ቦርሳዎች", "Boorsaa Kaardii Dubartootaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 13. Beauty, Cosmetics and Personal Care
	// -------------------------------------------------------------------------
	addCat("beauty-personal-care", []SubSectorTemplate{
		makeSub("skincare-facial-serums", "Facial Cleansers, Serums and Sunscreen", "የፊት መታጠቢያዎች፣ ሴረሞች እና የፀሐይ መከላከያ", "Qulqulleessituu Fuulaa fi Seeramii", "ናይ ገጽ መሐጸብን ጸሓይ መከላኸሊን", [][4]string{
			{"vitamin-c-hyaluronic-serums", "Vitamin C & Hyaluronic Acid Anti-Aging Brightening Serums", "ቫይታሚን ሲ እና ሀያሉሮኒክ የፊት ሴረም", "Seeramii Fuulaa Vaayitaaminii C"},
			{"broad-spectrum-spf50-sunscreens", "Invisible Broad Spectrum SPF 50+ Oil-Free Sunscreen Creams", "የፀሐይ መከላከያ ክሬሞች (SPF 50+)", "Kiriimii Aduu Dhorku SPF 50+"},
			{"gentle-foaming-facial-cleansers", "Gentle Hydrating Foaming Cleansers for Sensitive Skin", "የፊት መታጠቢያ ፈሳሽ ሳሙናዎች", "Saamunaa Fuulaa Dhangala'aa"},
			{"shea-butter-body-lotions", "Pure Shea Butter Deep Moisturizing Body Lotions", "የሺአ በተር የሰውነት ቅባቶች", "Looshinii Qullaa Shea Butter"},
		}),
	})

	// -------------------------------------------------------------------------
	// 15. Sports, Fitness and Outdoor Recreation
	// -------------------------------------------------------------------------
	addCat("sports-fitness", []SubSectorTemplate{
		makeSub("fitness-resistance-recovery", "Resistance Bands, Foam Rollers and Massage Guns", "የመወጠሪያ ሪባኖች፣ ፎም ሮለሮች እና የማሳጅ ሽጉጦች", "Meeshaa Shaakalaa fi Masajjii", "ናይ ስፖርት መወጠሪን ማሳጅን", [][4]string{
			{"heavy-duty-pullup-resistance-bands", "Natural Latex Heavy Loop Resistance Pull-Up Bands (5 Levels)", "የአካል ማጎልበቻ የመወጠሪያ ሪባኖች", "Rabaara Shaakalaa Humnaa"},
			{"percussion-deep-tissue-massage-guns", "Rechargeable Deep-Tissue Muscle Massage Therapy Guns", "የጡንቻ ማሳጅ ሽጉጦች (Massage Gun)", "Maashina Masajjii Maashaa"},
			{"high-density-eva-foam-rollers", "Deep Tissue Trigger Point EVA Foam Rollers for Recovery", "የአካል ማሸያ ፎም ሮለሮች", "Roolara Ispoortii Masajjii"},
		}),
	})

	// -------------------------------------------------------------------------
	// 16. Books, Literature and Publishing
	// -------------------------------------------------------------------------
	addCat("books-publishing", []SubSectorTemplate{
		makeSub("childrens-educational-books", "Children's Storybooks, Dictionaries and Reference", "የልጆች ተረቶች፣ መዝገበ ቃላት እና ማመሳከሪያዎች", "Kitaaba Daa'immanii fi Galmee Jechootaa", "ናይ ቆልዑ ዛንታታትን መዝገበ ቃላትን", [][4]string{
			{"illustrated-amharic-children-stories", "Illustrated Moral Storybooks in Amharic & English for Kids", "የህፃናት የስዕል ተረት መጻሕፍት", "Kitaaba Oduu Daa'immanii Fakkii Qabu"},
			{"trilingual-dictionaries-am-om-en", "Comprehensive Trilingual Dictionaries (Amharic - Oromo - English)", "የሶስት ቋንቋ መዝገበ ቃላት (አማርኛ-ኦሮምኛ-እንግሊዝኛ)", "Galmee Jechootaa Afaan Sadii"},
			{"geez-language-grammar-references", "Classical Ge'ez Grammar, Sawasew & Liturgical Text References", "የግዕዝ ሰዋሰው እና ቋንቋ መማሪያ መጻሕፍት", "Kitaaba Qorannoo Afaan Gi'iizii"},
		}),
	})

	// -------------------------------------------------------------------------
	// 17. Education, Courses and Training
	// -------------------------------------------------------------------------
	addCat("education-training", []SubSectorTemplate{
		makeSub("language-spoken-courses", "Spoken English, Arabic, Chinese and French", "የእንግሊዝኛ፣ የአረብኛ እና የውጭ ቋንቋ ስልጠናዎች", "Barnoota Afaan Alaa (English/Arabic)", "ናይ ወጻኢ ቋንቋታት ስልጠና", [][4]string{
			{"spoken-english-fluency-classes", "Intensive Spoken English Conversation & Accent Training", "የንግግር እንግሊዝኛ ተግባራዊ ስልጠና", "Leenjii Afaan Ingilizii Dubbachuu"},
			{"conversational-arabic-classes", "Modern Standard & Gulf Spoken Arabic Language Courses", "የአረብኛ ቋንቋ ተግባራዊ ስልጠና", "Barnoota Afaan Arabaa"},
			{"hsk-mandarin-chinese-classes", "HSK Mandarin Chinese for Business & Trade Courses", "የቻይንኛ ቋንቋ ስልጠና ለንግድ", "Barnoota Afaan Chaayinaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 21. Professional and Personal Services
	// -------------------------------------------------------------------------
	addCat("professional-services", []SubSectorTemplate{
		makeSub("appliance-electronics-repair", "Refrigerator, TV and Washing Machine Home Repair", "የፍሪጅ፣ የቴሌቪዥን እና የማጠቢያ ማሽን ጥገና", "Suphaa Firiijii, Televiziyinii fi Meeshaa", "ናይ ፍሪጅን ቴሌቪዥንን ጽገና", [][4]string{
			{"home-refrigerator-gas-refill-repair", "On-Site Refrigerator Compressor Repair and Gas Refill (R134a/R600a)", "የፍሪጅ ጋዝ መሙላት እና የኮምፕረሰር ጥገና በቤት", "Suphaa Firiijii fi Gaasii Guutuu"},
			{"led-smart-tv-screen-backlight-repair", "LED TV Motherboard & Backlight Strip Replacement Services", "የስማርት ቴሌቪዥን የቦርድ እና የስክሪን መብራት ጥገና", "Suphaa Televiziyinii Ismaartii"},
			{"washing-machine-motor-drum-service", "Automatic Washing Machine Motor, Pump & Bearing Replacement", "የልብስ ማጠቢያ ማሽን ጥገና በአድራሻዎ", "Suphaa Maashina Uffata Dhiquu"},
		}),
	})

	// -------------------------------------------------------------------------
	// 27. Industrial Machinery and Manufacturing
	// -------------------------------------------------------------------------
	addCat("industrial-machinery", []SubSectorTemplate{
		makeSub("metalworking-machine-tools", "Lathes, Milling Machines and Shapers", "የብረት መጥረቢያ ላቶች እና ማሽነሪዎች", "Maashina Sibila Qolchu (Laatii)", "ናይ ሓጺን መላጸዪ ላትን ማሽናትን", [][4]string{
			{"universal-horizontal-metal-lathes", "Heavy Duty Universal Gap-Bed Metal Turning Lathes", "የትላልቅ የብረት መጥረቢያ ላቶች (Lathe)", "Maashina Laatii Sibilaa Gurguddaa"},
			{"universal-vertical-milling-machines", "Multi-Axis Vertical & Horizontal Gear Milling Machines", "የማሽነሪ እቃዎች መስሪያ ሚሊንግ ማሽኖች", "Maashina Miiliingii Sibilaa"},
			{"heavy-hydraulic-guillotine-shears", "Hydraulic Sheet Metal Guillotine Shearing & Bending Machines", "የቆርቆሮ እና የብረት ሰሌዳ መቁረጫ ጊዮቲኖች", "Maashina Qorqoorroo Kutuu Jabaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 28. Tools, Workshop and Repair Equipment
	// -------------------------------------------------------------------------
	addCat("tools-workshop", []SubSectorTemplate{
		makeSub("air-compressors-pneumatics", "Piston Compressors, Spray Guns and Air Hoses", "የአየር ኮምፕሬሰሮች እና የቀለም መርጫዎች", "Kompreesara Qilleensaa fi Meeshaa Qalamaa", "ናይ ኣየር ኮምፕሬሰርን ቀለም መርጨውን", [][4]string{
			{"two-stage-belt-air-compressors", "100L to 300L Belt-Driven Industrial Piston Air Compressors", "ባለ ቀበቶ የአየር ኮምፕሬሰሮች (100L-300L)", "Kompreesara Qilleensaa 100L-300L"},
			{"high-volume-low-pressure-spray-guns", "HVLP Gravity Feed Automotive Paint Spray Guns", "የመኪና ቀለም መርጫ ሽጉጦች (Spray Gun)", "Meeshaa Qalama Konkolaataa Biifu"},
			{"reinforced-rubber-air-hoses", "Heavy Duty 20-Meter Braided Rubber Pneumatic Air Hoses", "የአየር ኮምፕሬሰር ጠንካራ ቱቦዎች", "Ujummoo Qilleensaa Jabaa 20m"},
		}),
	})

	// -------------------------------------------------------------------------
	// 32. Plumbing, Electrical and Building Systems
	// -------------------------------------------------------------------------
	addCat("plumbing-electrical", []SubSectorTemplate{
		makeSub("sanitary-fixtures-faucets", "Bathroom Sinks, Toilets, Showers and Chrome Mixers", "የመታጠቢያ ገንዳዎች፣ ቧንቧዎች እና ሽንት ቤቶች", "Meeshaalee Mana Finiinaa fi Bishaanii", "ናይ መሕጸብን ሽቃቕን ንዋያት", [][4]string{
			{"ceramic-dual-flush-toilet-sets", "Modern Water-Saving Ceramic Dual-Flush Water Closet Toilets", "ባለ ሁለት ፍላሽ የሴራሚክ ሽንት ቤቶች", "Waan Mana Finiinaa Seeraamikii"},
			{"brass-chrome-basin-mixer-taps", "Heavy Brass Chrome-Plated Hot & Cold Bathroom Basin Mixers", "የመታጠቢያ ቤት የሙቅ እና የቀዝቃዛ ቧንቧዎች", "Ujummoo Bishaanii Ho'aa/Qabbanaawaa"},
			{"overhead-rainfall-shower-columns", "Stainless Steel Overhead Rainfall Shower Head Column Sets", "የዝናብ መሰል የሻወር ራስ እና ቋሚዎች", "Shaaworii Bishaanii Roobaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 37. Vehicles and Automobiles
	// -------------------------------------------------------------------------
	addCat("vehicles-cars", []SubSectorTemplate{
		makeSub("luxury-suvs-crossovers", "High-End Luxury SUVs and Crossovers", "የቅንጦት 4x4 ጂፖች", "Konkolaattota Qananii 4x4", "ናይ ቅንጦት 4x4 ጂፓት", [][4]string{
			{"mercedes-gle-gls-suvs", "Mercedes-Benz GLE & GLS 4MATIC Luxury SUVs", "መርሴዲስ ቤንዝ GLE እና GLS ጂፖች", "Konkolaataa Marseedis GLE/GLS"},
			{"bmw-x5-x7-xdrive", "BMW X5 and X7 xDrive TwinPower Turbo SUVs", "ቢኤምደብሊው X5 እና X7 ጂፖች", "Konkolaataa BMW X5/X7"},
			{"range-rover-vogue-sport", "Land Rover Range Rover Vogue & Sport Autobiography", "ሬንጅ ሮቨር ቮግ እና ስፖርት", "Konkolaataa Range Rover Vogue"},
		}),
	})

	// -------------------------------------------------------------------------
	// 40. Trucks, Buses and Commercial Transport
	// -------------------------------------------------------------------------
	addCat("trucks-commercial-transport", []SubSectorTemplate{
		makeSub("bulk-liquid-chemical-tankers", "Fuel, Bitumen and Liquid Cargo Tankers", "የነዳጅ፣ የሬንጅ እና የፈሳሽ ቦቴዎች", "Konkolaataa Bootee Boba'aa fi Zayitaa", "ናይ ነዳድን ሬንጅን ቦቴታት", [][4]string{
			{"35000l-aluminum-fuel-tankers", "35,000 to 45,000 Liter Multi-Compartment Aluminum Fuel Tanker Trailers", "ትላልቅ የነዳጅ ማመላለሻ ቦቴዎች (45000 ሊትር)", "Bootee Boba'aa Gurguddaa 45000L"},
			{"insulated-hot-bitumen-tankers", "Insulated Heating Bitumen Asphalt Road Tankers", "የአስፋልት ሬንጅ ማመላለሻ ሙቅ ቦቴዎች", "Bootee Asfaaltii fi Bitumeenii"},
		}),
	})

	// -------------------------------------------------------------------------
	// 44. Crops, Grains, Seeds and Plants
	// -------------------------------------------------------------------------
	addCat("crops-grains-seeds", []SubSectorTemplate{
		makeSub("traditional-root-crops-enset", "Enset (False Banana), Kocho and Bulla", "የእንሰት፣ የቆጮ እና የቡላ ምርቶች", "Oomisha Qoccoo fi Buullaa", "ናይ ቆጮን ቡላን ፍርያት", [][4]string{
			{"fermented-white-kocho-bulk", "High-Quality Fermented White Gurage / Sidama Kocho (Per Quintal)", "የተዘጋጀ ንጹህ የጉራጌ እና የሲዳማ ቆጮ", "Qoccoo Adii Qulqulluu"},
			{"refined-pure-enset-bulla-powder", "Refined High-Grade Enset Starch Powder (Bulla)", "ንጹህ የቡላ ዱቄት በኩንታል", "Daakuu Buullaa Qulqulluu"},
			{"fresh-enset-corm-plants", "Mature Enset Plants and Corm Suckers for Planting", "የእንሰት ችግኞች እና ማሳደጊያዎች", "Biqiltuu Qoccoo fi Enset"},
		}),
	})

	// -------------------------------------------------------------------------
	// 45. Livestock, Poultry and Animal Husbandry
	// -------------------------------------------------------------------------
	addCat("livestock-poultry", []SubSectorTemplate{
		makeSub("dairy-cattle-genetics", "Heifers, Artificial Insemination and Milking Cows", "የወተት ላሞች፣ ጊደሮች እና ዝርያዎች", "Saawwaa Aannanii fi Gidara", "ናይ ጸባ ላሞትን ጣቓታትን", [][4]string{
			{"pure-holstein-friesian-heifers", "High-Index Pregnant Holstein Friesian Dairy Heifers", "እርጉዝ የሆልስታይን ፍሪዚያን የወተት ጊደሮች", "Gidara Saawwaa Holisteen Aannanii"},
			{"jersey-dairy-crossbred-cows", "High-Butterfat Jersey Crossbred Milking Cows", "የጀርሲ ዝርያ ከፍተኛ ቅቤ የወተት ላሞች", "Saawwaa Jeersii Aannaniifi Dhadhaa"},
			{"frozen-bovine-semen-straws", "Imported Pedigree Dairy Bull Frozen Semen Straws (AI)", "የላቀ የወተት ዝርያ የበሬ ዘር (AI)", "Sanyii Kormaa Saawwaa (AI)"},
		}),
	})

	// -------------------------------------------------------------------------
	// 49. Beverages, Coffee and Tea
	// -------------------------------------------------------------------------
	addCat("coffee-tea-beverages", []SubSectorTemplate{
		makeSub("ethiopian-highland-teas", "Black Tea, Green Tea and Herbal Infusions", "የሀበሻ ሻይ ቅጠሎች እና የዕፅዋት ሻይ", "Baala Shaayii fi Qoricha Aadaa", "ናይ ዓዲ ሻህን ናይ ዕጽዋት ሻህን", [][4]string{
			{"wushwush-gumaro-black-tea", "Wushwush & Gumaro Premium Highland Orthodox Black Tea", "የውሽውሽ እና የጉማሮ ጥቁር ሻይ ቅጠል", "Baala Shaayii Gurraacha Wushwush"},
			{"organic-ethiopian-green-tea", "Unfermented Pure Highland Organic Green Tea Leaves", "ንጹህ አረንጓዴ ሻይ (ግሪን ቲ)", "Shaayii Magariisa Qulqulluu"},
			{"traditional-herbal-tea-blends", "Wild Thyme (Tosign) & Spiced Herbal Infusion Blends", "የጦስኝ፣ የሎሚ ሳር እና የቅመም ሻይ", "Shaayii Xoosinnyii fi Lemoonii"},
		}),
	})

	// -------------------------------------------------------------------------
	// 51. Music, Musical Instruments and Entertainment
	// -------------------------------------------------------------------------
	addCat("music-instruments", []SubSectorTemplate{
		makeSub("modern-orchestral-instruments", "Guitars, Keyboards, Brass and Drum Kits", "ጊታሮች፣ ኪቦርዶች፣ ድራሞች እና መለከቶች", "Gitaara, Kiiboordii fi Diraama", "ጊታርን ኪቦርድን ድራምን", [][4]string{
			{"solid-top-acoustic-guitars", "Solid Spruce Top Steel-String Acoustic Guitars", "የአኮስቲክ ጊታሮች", "Gitaara Akustikii Sibiilaa"},
			{"arranger-synthesizer-keyboards", "61 & 88 Key Arranger Workstation Piano Keyboards (Yamaha/Korg)", "የሙዚቃ ማቀናበሪያ ፒያኖ ኪቦርዶች", "Kiiboordii Piyoanoo fi Meeshaa Muuziqaa"},
			{"5-piece-acoustic-drum-sets", "Complete 5-Piece Stage Acoustic Drum Kits with Cymbals", "የመድረክ አኮስቲክ ድራም ስብስቦች", "Diraamii Muuziqaa Guutuu"},
			{"brass-saxophones-trumpets", "Alto Saxophones and Brass Bb Concert Trumpets", "ሳክሶፎን እና ትራምፔት መለከቶች", "Saksoofoonii fi Malakata"},
		}),
	})

	// -------------------------------------------------------------------------
	// 53. Toys, Baby Products and Children's Goods
	// -------------------------------------------------------------------------
	addCat("toys-baby-products", []SubSectorTemplate{
		makeSub("educational-electronic-toys", "Educational Puzzles, Building Blocks and RC Cars", "ትምህርታዊ መጫወቻዎች እና የርቀት መቆጣጠሪያ መኪኖች", "Meeshaa Taphannaa fi Ijaarsaa", "ናይ ምስትውዓል መጻወትን መካይንን", [][4]string{
			{"creative-building-block-sets", "Multi-Piece Creative Interlocking Plastic Building Blocks", "የልጆች የመገጣጠሚያ ብሎኮች (ሌጎ መሰል)", "Meeshaa Ijaarsaa Ijoollee (Lego)"},
			{"rechargeable-rc-monster-trucks", "4WD High-Speed Rechargeable Remote Control Toy Cars", "የርቀት መቆጣጠሪያ ያላቸው መጫወቻ መኪኖች", "Konkolaataa Taphannaa Remootii Qabu"},
			{"early-learning-electronic-tablets", "Bilingual Early Word & Alphabet Educational Audio Tablets", "የልጆች የፊደል እና የቋንቋ መማሪያ መጫወቻ ፓዶች", "Taableetii Taphannaa Qubee fi Lakkoofsaa"},
		}),
	})

	// -------------------------------------------------------------------------
	// 54. Pet Supplies and Companion Animals
	// -------------------------------------------------------------------------
	addCat("pet-supplies", []SubSectorTemplate{
		makeSub("pet-cages-accessories", "Dog Houses, Kennels, Cages and Leashes", "የውሻ ቤቶች፣ አጥር፣ ካቴኖች እና ሰንሰለቶች", "Mana Saree, Saanduqa fi Sansalata", "ናይ ከልቢ ገዛን ሰንሰለትን", [][4]string{
			{"heavy-duty-folding-dog-crates", "Heavy Gauge Folding Double-Door Metal Dog Wire Crates", "የብረት የውሻ ማቆያ ካቴኖች (ሳጥኖች)", "Saanduqa Sibiilaa Saree Eegu"},
			{"retractable-heavy-duty-dog-leashes", "Tangle-Free Retractable Nylon Tape Walking Dog Leashes", "የውሻ መምሪያ ሰንሰለቶች እና ገመዶች", "Wadaroo fi Sansalata Saree"},
			{"stainless-steel-pet-food-bowls", "Non-Slip Weighted Stainless Steel Pet Food & Water Bowls", "የውሻ እና የድመት ምግብ እና ውሃ ሳህኖች", "Saanii Nyaataafi Bishaan Saree"},
		}),
	})

	// -------------------------------------------------------------------------
	// 57. Water, Irrigation and Environmental Equipment
	// -------------------------------------------------------------------------
	addCat("water-irrigation-equipment", []SubSectorTemplate{
		makeSub("water-filtration-purification", "Reverse Osmosis Units, Sand Filters and UV Purifiers", "የውሃ ማጣሪያ ማሽኖች እና አልትራቫዮሌት ዩኒቶች", "Maashina Bishaan Qulqulleessu", "ናይ ማይ መጻረዪ ማሽን", [][4]string{
			{"commercial-reverse-osmosis-plants", "500 to 2000 Liters/Hour Commercial Reverse Osmosis Water Plants", "የውሃ ማጣሪያ ትላልቅ ማሽኖች (RO)", "Maashina Bishaan Qulqulleessu RO"},
			{"multi-stage-under-sink-water-filters", "5-Stage Domestic Under-Sink Drinking Water Purifiers with Faucet", "የቤት ውስጥ የቧንቧ ውሃ ማጣሪያዎች", "Filtarii Bishaan Dhugatii Manaa"},
			{"ultraviolet-uv-water-sterilizers", "Stainless Steel Inline UV Germicidal Water Sterilizers", "የአልትራቫዮሌት ባክቴሪያ ማጥፊያ ማጣሪያዎች", "Qulqulleessituu Bakteeriyaa UV"},
		}),
	})

	// -------------------------------------------------------------------------
	// 58. Travel, Accommodation and Tourism Services
	// -------------------------------------------------------------------------
	addCat("travel-tourism", []SubSectorTemplate{
		makeSub("guest-houses-boutique-stays", "Furnished Guest Houses, Lodges and Boutique Stays", "የእንግዳ ማረፊያ ቤቶች፣ ሎጆች እና ቪላዎች", "Mana Boqonnaa fi Loojiiwwan", "ናይ ኣጋይሽ መዕረፍን ሎጃትን", [][4]string{
			{"addis-ababa-bole-guest-houses", "Luxury Furnished Daily Guest Houses in Bole and Kazanchis", "የቦሌ እና የካዛንቺስ የእንግዳ ማረፊያ ቤቶች", "Mana Boqonnaa Boolleefi Kaazaanchiis"},
			{"rift-valley-lakeside-resort-lodges", "Bishoftu & Langano Lakeside Eco-Resort Lodges & Cottages", "የቢሾፍቱ እና የላንጋኖ የሐይቅ ዳርቻ ሪዞርት ሎጆች", "Loojiiwwan Haroo Bishooftuu fi Laanganoo"},
			{"historic-route-boutique-hotels", "Lalibela & Gondar Historic Route Heritage Boutique Hotels", "የላሊበላ እና የጎንደር የቱሪስት ሆቴሎች", "Hoteelota Seenaa Laallibalaa fi Gondar"},
		}),
	})

	// -------------------------------------------------------------------------
	// 59. Packaging, Shipping and Warehousing
	// -------------------------------------------------------------------------
	addCat("packaging-shipping", []SubSectorTemplate{
		makeSub("shipping-seals-strapping-tools", "Security Seals, Steel Strapping and Pallet Trucks", "የኮንቴይነር ሲሎች፣ ማሰሪያዎች እና የእቃ መግፊያዎች", "Meeshaalee Saamsuu fi Harka Fe'uu", "ናይ ኮንቴይነር ሲልን መኣሰርን", [][4]string{
			{"container-bolt-security-seals", "ISO 17712 High-Security Container Door Bolt Seals", "የኮንቴይነር በር ብረት ሲሎች (Seals)", "Siilii Sibiilaa Konteeyinaraa"},
			{"hydraulic-manual-pallet-jacks", "2.5-Ton and 3-Ton Heavy Duty Manual Hydraulic Pallet Trucks", "የእቃ መግፊያ ሃይድሮሊክ ፎርኮች (Pallet Jack)", "Jaakii Paaleetii Fe'iinsaa 3-Ton"},
			{"polypropylene-strapping-tensioners", "Manual Steel & Plastic Box Strapping Band Tensioner Tools", "የካርቶን ማሰሪያ ፕላስቲክ እና መሳሪያዎች", "Meeshaa Wadaroo Kaartoonii Hidhu"},
		}),
	})

	// -------------------------------------------------------------------------
	// 60. General Classifieds, Donations and Exchanges
	// -------------------------------------------------------------------------
	addCat("general-classifieds", []SubSectorTemplate{
		makeSub("community-notices-announcements", "Public Notices, Tenders and Memorial Announcements", "የህዝብ ማስታወቂያዎች፣ ጨረታዎች እና መታሰቢያዎች", "Beeksisa Hawaasaa fi Caalbaasii", "ናይ ህዝቢ ምልክታን ጨረታን", [][4]string{
			{"commercial-public-tender-notices", "Government, NGO and Private Commercial Tender Procurement Notices", "የግል፣ የድርጅት እና የመንግስት ጨረታዎች", "Beeksisa Caalbaasii Mootummaa fi NGO"},
			{"obituary-memorial-notices", "Obituary, Funeral Tributes and Memorial Service Announcements", "የሀዘን እና የመታሰቢያ ማስታወቂያዎች", "Beeksisa Gadda fi Yaadannoo"},
			{"lost-pet-notices", "Lost & Found Family Dogs, Puppies and Companion Pets", "የጠፉ ውሾች እና የቤት እንስሳት ማስታወቂያዎች", "Beeksisa Saree fi Horii Manaa Bade"},
		}),
	})

	return catalogues
}
