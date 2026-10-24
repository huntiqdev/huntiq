package main

// GetFull60DomainData provides deep branches and leaf items for all 60 marketplace domains
func GetFull60DomainData() []TaxonomyDomain {
	return []TaxonomyDomain{
		// =====================================================================
		// DOMAIN 01: Electronics & Mobile Devices
		// =====================================================================
		{
			Slug: "electronics-mobile",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "smartphones-mobile", NameEN: "Smartphones and Mobile Handsets", NameAM: "ስማርት ስልኮች እና ሞባይል", NameOM: "Bilbila Ismaartii", NameTI: "ስማርት ስልኪ",
					Leaves: []TaxonomyLeaf{
						{Slug: "flagship-smartphones", NameEN: "Flagship High-Performance Smartphones", NameAM: "ከፍተኛ ደረጃ ስማርት ስልኮች", NameOM: "Bilbila Ismaartii Olaanaa", NameTI: "ላዕለዋይ ደረጃ ስማርት ስልኪ"},
						{Slug: "budget-smartphones", NameEN: "Budget & Entry-Level Smartphones", NameAM: "መካከለኛ እና የበጀት ስልኮች", NameOM: "Bilbila Baasii Xiqqaata", NameTI: "ናይ በጀት ስማርት ስልኪ"},
						{Slug: "foldable-smartphones", NameEN: "Foldable & Flip Screen Smartphones", NameAM: "ተጣጣፊ ስክሪን ስልኮች", NameOM: "Bilbila Dachaafamu", NameTI: "ዝዕጸፍ ስክሪን ስልኪ"},
						{Slug: "rugged-outdoor-phones", NameEN: "Waterproof & Shockproof Rugged Phones", NameAM: "ጠንካራ የውሃ መቋቋሚ ስልኮች", NameOM: "Bilbila Jabaa Bishaan Danda'u", NameTI: "ተሪር ማይ መከላኸሊ ስልኪ"},
						{Slug: "dual-sim-keypad-phones", NameEN: "Dual-SIM Basic Keypad Mobile Phones", NameAM: "ባለ ኪፓድ መደበኛ ስልኮች", NameOM: "Bilbila Qabduu Kiipaadii", NameTI: "ናይ ቁልፊ ሞባይል ስልኪ"},
					},
				},
				{
					Slug: "tablet-computers", NameEN: "Tablet Computers and Graphics Pads", NameAM: "ታብሌቶች እና ዲጂታል ፓዶች", NameOM: "Taableetota fi Paadiiwwan", NameTI: "ታብሌታትን ዲጂታል ፓድስን",
					Leaves: []TaxonomyLeaf{
						{Slug: "professional-tablets", NameEN: "11-inch to 13-inch Professional Tablets", NameAM: "ትላልቅ የቢሮ ታብሌቶች", NameOM: "Taableetii Hojii Gurguddaa", NameTI: "ዓበይቲ ናይ ስራሕ ታብሌታት"},
						{Slug: "compact-tablets-mini", NameEN: "7-inch to 9-inch Compact Handheld Tablets", NameAM: "አነስተኛ ተንቀሳቃሽ ታብሌቶች", NameOM: "Taableetii Xixinnoo", NameTI: "ንኣሽቱ ታብሌታት"},
						{Slug: "graphics-drawing-tablets", NameEN: "Digital Drawing Tablets & Stylus Pens", NameAM: "የስዕል እና የዲዛይን ታብሌቶች", NameOM: "Taableetii Fakkii Kaasu", NameTI: "ናይ ስእሊ ታብሌት"},
					},
				},
				{
					Slug: "power-charging-gear", NameEN: "Chargers, Power Banks and Cables", NameAM: "ቻርጀሮች፣ ፓወር ባንክ እና ገመዶች", NameOM: "Chaarcharoota fi Paawar Baankii", NameTI: "ቻርጀራትን ፓወር ባንክን",
					Leaves: []TaxonomyLeaf{
						{Slug: "gan-fast-chargers", NameEN: "65W-140W GaN Ultra-Fast USB-C Chargers", NameAM: "የጋን ፈጣን ቻርጀሮች (GaN)", NameOM: "Chaarchara Saffisaa GaN", NameTI: "GaN ፈጣን ቻርጀር"},
						{Slug: "high-capacity-power-banks", NameEN: "20,000mAh to 50,000mAh Power Banks", NameAM: "ከፍተኛ አቅም ያላቸው ፓወር ባንኮች", NameOM: "Paawar Baankii Dandeettii Guddaa", NameTI: "ዓቢይ ዓቕሚ ዘለዎ ፓወር ባንክ"},
						{Slug: "magsafe-wireless-chargers", NameEN: "Magnetic Wireless Qi & MagSafe Chargers", NameAM: "ገመድ አልባ ማግኔቲክ ቻርጀሮች", NameOM: "Chaarchara Wayaalessii", NameTI: "ገመድ ኣልቦ ቻርጀር"},
						{Slug: "braided-fast-data-cables", NameEN: "Braided 100W USB-C and Lightning Data Cables", NameAM: "የተሸረቡ ጠንካራ የዳታ ኬብሎች", NameOM: "Keebilii Daataa Jabaa", NameTI: "ተሪር ናይ ዳታ ኬብል"},
					},
				},
				{
					Slug: "audio-listening-wearables", NameEN: "Wireless Earbuds, Headphones and Headsets", NameAM: "የጆሮ ማዳመጫዎች እና ሄድፎኖች", NameOM: "Heedifoonii fi Qulqulleessaa", NameTI: "ናይ እዝኒ መዳመጺን ሄድፎንን",
					Leaves: []TaxonomyLeaf{
						{Slug: "anc-wireless-earbuds", NameEN: "Active Noise-Canceling (ANC) Earbuds", NameAM: "ድምጽ የሚቀንሱ ገመድ አልባ ኤርፖዶች", NameOM: "Earbuds Sagalee Qulqulleessu", NameTI: "ድምጺ መጽረዪ ኤርፖድስ"},
						{Slug: "over-ear-anc-headphones", NameEN: "Over-Ear Wireless Studio Headphones", NameAM: "ኦቨር-ኢየር ሄድፎኖች", NameOM: "Heedifoonii Gurraa Guutuu", NameTI: "ኦቨር-ኢየር ሄድፎን"},
						{Slug: "neckband-sports-earphones", NameEN: "Neckband Waterproof Sports Earphones", NameAM: "የስፖርት የጆሮ ማዳመጫዎች", NameOM: "Heedifoonii Ispoortii", NameTI: "ናይ ስፖርት ሄድፎን"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 02: Computers, IT & Office Technology
		// =====================================================================
		{
			Slug: "computers-it",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "portable-laptops", NameEN: "Laptops and Notebook Computers", NameAM: "ላፕቶፖች እና ኖትቡኮች", NameOM: "Laptooppiiwwan", NameTI: "ላፕቶፓት",
					Leaves: []TaxonomyLeaf{
						{Slug: "intel-core-ultra-laptops", NameEN: "Intel Core Ultra / Core i7 Laptops", NameAM: "ኢንቴል ኮር i7/Ultra ላፕቶፖች", NameOM: "Laptooppii Intel Core i7", NameTI: "ኢንቴል ኮር i7 ላፕቶፕ"},
						{Slug: "amd-ryzen-pro-laptops", NameEN: "AMD Ryzen 7 & 9 High-Performance Laptops", NameAM: "ኤኤምዲ ራይዘን ላፕቶፖች", NameOM: "Laptooppii AMD Ryzen", NameTI: "AMD ራይዘን ላፕቶፕ"},
						{Slug: "apple-silicon-macbooks", NameEN: "Apple MacBook Pro & Air (M-Series Silicon)", NameAM: "አፕል ማክቡክ (M1/M2/M3/M4)", NameOM: "MacBook Apple M-Series", NameTI: "ኣፕል ማክቡክ"},
						{Slug: "convertible-2in1-touch-laptops", NameEN: "2-in-1 Convertible 360-Touch Laptops", NameAM: "ተጣጣፊ የንክኪ ስክሪን ላፕቶፖች", NameOM: "Laptooppii Tuqaa Dachaafamu", NameTI: "ናይ ምትንካፍ ስክሪን ላፕቶፕ"},
					},
				},
				{
					Slug: "desktop-pc-hardware", NameEN: "Desktop Towers, All-in-Ones and Workstations", NameAM: "ዴስክቶፕ ኮምፒውተሮች እና ወርክስቴሽኖች", NameOM: "Kompiitara Minjaalaa", NameTI: "ናይ ዴስክቶፕ ኮምፒዩተር",
					Leaves: []TaxonomyLeaf{
						{Slug: "compact-all-in-one-computers", NameEN: "Core i5 / i7 All-in-One Desktop PCs", NameAM: "ኦል-ኢን-ዋን ዴስክቶፖች", NameOM: "Kompiitara All-in-One", NameTI: "ኦል-ኢን-ዋን ፒሲ"},
						{Slug: "cad-3d-workstation-towers", NameEN: "CAD & 3D Rendering Workstations (Xeon/Threadripper)", NameAM: "የ3ዲ እና አርክቴክቸር ወርክስቴሽኖች", NameOM: "Workstation CAD fi 3D", NameTI: "ናይ 3D ወርክስቴሽን"},
						{Slug: "tower-cpu-units", NameEN: "Core i3 to i7 Standalone Desktop Tower Units", NameAM: "የኮምፒውተር ሲፒዩ ሳጥኖች (CPU)", NameOM: "Saanduqa CPU Kompiitaraa", NameTI: "ናይ ኮምፒዩተር CPU ሳንዱቕ"},
					},
				},
				{
					Slug: "pc-hardware-components", NameEN: "Computer Storage, RAM, and Upgrades", NameAM: "የኮምፒውተር እቃዎች እና ማሻሻያዎች", NameOM: "Qaamolee Kompiitaraa", NameTI: "ክፍልታት ኮምፒዩተር",
					Leaves: []TaxonomyLeaf{
						{Slug: "gen4-nvme-pcie-ssd", NameEN: "PCIe Gen 4 & Gen 5 High Speed NVMe SSDs", NameAM: "ፈጣን ኤስኤስዲ (NVMe SSD)", NameOM: "Diraayivii NVMe SSD Saffisaa", NameTI: "ፈጣን NVMe SSD"},
						{Slug: "high-frequency-ddr5-ram", NameEN: "DDR5 High-Frequency Desktop & Laptop RAM", NameAM: "DDR5 የኮምፒውተር ራም", NameOM: "RAM Kompiitaraa DDR5", NameTI: "DDR5 ናይ ኮምፒዩተር ራም"},
						{Slug: "gaming-video-cards-rtx", NameEN: "NVIDIA GeForce RTX 30/40 Series Graphics Cards", NameAM: "የቪዲዮ ግራፊክስ ካርዶች (RTX)", NameOM: "Kaardii Giraafiiksii RTX", NameTI: "RTX ናይ ግራፊክስ ካርድ"},
						{Slug: "external-portable-backup-ssd", NameEN: "Rugged External Portable SSD Drives (1TB-4TB)", NameAM: "ተንቀሳቃሽ ኤክስተርናል ኤስኤስዲ", NameOM: "External SSD Bakka-buusaa", NameTI: "ተንቀሳቓሲ ኤክስተርናል SSD"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 11: Shoes, Bags and Accessories
		// =====================================================================
		{
			Slug: "shoes-bags",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "mens-footwear", NameEN: "Men's Footwear and Leather Shoes", NameAM: "የወንዶች ጫማዎች", NameOM: "Kophee Dhiiraa", NameTI: "ናይ ደቂ ተባዕትዮ ሳእኒ",
					Leaves: []TaxonomyLeaf{
						{Slug: "ethiopian-leather-oxford-shoes", NameEN: "Genuine Ethiopian Leather Oxford & Derby Shoes", NameAM: "የሀበሻ ቆዳ የወንዶች ጫማዎች", NameOM: "Kophee Gogaa Itoophiyaa", NameTI: "ናይ ኢትዮጵያ ቆርበት ሳእኒ"},
						{Slug: "mens-casual-loafers", NameEN: "Men's Casual Slip-On Leather Loafers", NameAM: "የወንዶች ሎፈር ጫማዎች", NameOM: "Kophee Lofaaraa", NameTI: "ናይ ወንዶች ሎፈር ሳእኒ"},
						{Slug: "athletic-running-sneakers", NameEN: "Athletic Running and Basketball Sneakers", NameAM: "የስፖርት እና የሩጫ ስኒከሮች", NameOM: "Kophee Fiigichaa fi Ispoortii", NameTI: "ናይ ሩጫ ስኒከር"},
						{Slug: "safety-steel-toe-boots", NameEN: "Industrial Steel-Toe Work Safety Boots", NameAM: "የግንባታ እና የፋብሪካ ሴፍቲ ጫማዎች", NameOM: "Kophee Eegumsaa Hojii", NameTI: "ናይ ስራሕ ሴፍቲ ሳእኒ"},
					},
				},
				{
					Slug: "womens-footwear", NameEN: "Women's Footwear and Heels", NameAM: "የሴቶች ጫማዎች", NameOM: "Kophee Dubartootaa", NameTI: "ናይ ደቀንስትዮ ሳእኒ",
					Leaves: []TaxonomyLeaf{
						{Slug: "pointed-high-heels", NameEN: "Pointed Stiletto & Block High Heels", NameAM: "ባለ ተረከዝ ጫማዎች (ሂልስ)", NameOM: "Kophee Koomee Dheeraa", NameTI: "ናይ ተረከዝ ሳእኒ"},
						{Slug: "comfortable-flat-pumps", NameEN: "Comfortable Leather Flat Pumps & Ballet Flats", NameAM: "ፍላት ጫማዎች", NameOM: "Kophee Dirree Salphaa", NameTI: "ፍላት ሳእኒ"},
						{Slug: "womens-ankle-boots", NameEN: "Fashion Ankle Boots and Platform Shoes", NameAM: "የሴቶች ቦቲ ጫማዎች", NameOM: "Kophee Bootii Dubartootaa", NameTI: "ናይ ደቀንስትዮ ቦቲ ሳእኒ"},
					},
				},
				{
					Slug: "luggage-travel-backpacks", NameEN: "Bags, Luggage and Travel Accessories", NameAM: "ቦርሳዎች፣ ሻንጣዎች እና የጉዞ እቃዎች", NameOM: "Boorsaa fi Saanduqa Imalaa", NameTI: "ቦርሳታትን ናይ መገሻ ሻንጣን",
					Leaves: []TaxonomyLeaf{
						{Slug: "hard-shell-spinner-suitcases", NameEN: "Hard-Shell 360-Spinner Suitcase Sets", NameAM: "ባለ ሮለር የጉዞ ሻንጣዎች", NameOM: "Saanduqa Imalaa Naanna'u", NameTI: "ናይ መገሻ ሻንጣ"},
						{Slug: "genuine-leather-handbags", NameEN: "Ethiopian Pure Leather Handbags & Totes", NameAM: "የሴቶች የእጅ ቆዳ ቦርሳዎች", NameOM: "Boorsaa Harka Gogaa", NameTI: "ናይ ቆርበት ናይ ኢድ ቦርሳ"},
						{Slug: "anti-theft-laptop-backpacks", NameEN: "Waterproof Anti-Theft Laptop Backpacks", NameAM: "የላፕቶፕ ቦርሳዎች", NameOM: "Boorsaa Dugdaa Laptooppii", NameTI: "ናይ ላፕቶፕ ቦርሳ"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 12: Jewelry, Watches and Luxury Goods
		// =====================================================================
		{
			Slug: "jewelry-luxury",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "ethiopian-gold-jewelry", NameEN: "Ethiopian Traditional and Modern Gold", NameAM: "የኢትዮጵያ ወርቅ ጌጣጌጦች", NameOM: "Faaya Warqee Itoophiyaa", NameTI: "ናይ ኢትዮጵያ ወርቂ ስልማት",
					Leaves: []TaxonomyLeaf{
						{Slug: "21k-traditional-gold-necklaces", NameEN: "21 Karat Handcrafted Traditional Gold Necklaces", NameAM: "ባህላዊ የ21 ካራት ወርቅ ሀብሎች", NameOM: "Mormaa Warqee Kaaraata 21", NameTI: "ናይ 21 ካራት ወርቂ ክሳድ ስልማት"},
						{Slug: "gold-wedding-bands-rings", NameEN: "Gold Wedding Bands & Engagement Rings", NameAM: "የጋብቻ እና የቁርባን የወርቅ ቀለበቶች", NameOM: "Qubee Cidhaa Warqee", NameTI: "ናይ መርዓ ወርቂ ቐለቤት"},
						{Slug: "traditional-gold-earrings", NameEN: "Ethiopian Filigree Gold Earrings and Chokers", NameAM: "የወርቅ የጆሮ ጌጦች", NameOM: "Gurraa Warqee Aadaa", NameTI: "ናይ ወርቂ ናይ እዝኒ ስልማት"},
					},
				},
				{
					Slug: "silver-traditional-crosses", NameEN: "Pure Silver Jewelry and Traditional Crosses", NameAM: "የብር ጌጣጌጦች እና መስቀሎች", NameOM: "Faaya Meetaa fi Fannoo", NameTI: "ናይ ብሩር ስልማትን መስቀልን",
					Leaves: []TaxonomyLeaf{
						{Slug: "ethiopian-silver-procession-crosses", NameEN: "Handmade Ethiopian Silver Neck Crosses (Meskel)", NameAM: "የብር መስቀሎች እና ሀብሎች", NameOM: "Fannoo Meetaa Harka Tolchame", NameTI: "ናይ ብሩር መስቀል"},
						{Slug: "925-sterling-silver-jewelry", NameEN: "925 Sterling Silver Chains and Bracelets", NameAM: "የ925 ብር ሰንሰለቶች እና አምባሮች", NameOM: "Alangaa Meetaa 925", NameTI: "ናይ 925 ብሩር ሰንሰለት"},
					},
				},
				{
					Slug: "luxury-wristwatches", NameEN: "Men's and Women's Luxury Watches", NameAM: "የቅንጦት የእጅ ሰዓቶች", NameOM: "Sa'aatii Harkaa Qananii", NameTI: "ናይ ቅንጦት ናይ ኢድ ሰዓት",
					Leaves: []TaxonomyLeaf{
						{Slug: "automatic-mechanical-watches", NameEN: "Automatic Self-Winding Mechanical Watches", NameAM: "አውቶማቲክ ሜካኒካል ሰዓቶች", NameOM: "Sa'aatii Meekaaniikaa Ofumaan Hojjetu", NameTI: "ኣውቶማቲክ ሜካኒካል ሰዓት"},
						{Slug: "quartz-chronograph-watches", NameEN: "Stainless Steel Chronograph Quartz Watches", NameAM: "ክሮኖግራፍ የብረት ሰዓቶች", NameOM: "Sa'aatii Kiroonogiraafii", NameTI: "ናይ ሓጺን ክሮኖግራፍ ሰዓት"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 13: Beauty, Cosmetics and Personal Care
		// =====================================================================
		{
			Slug: "beauty-personal-care",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "hair-care-styling", NameEN: "Hair Care, Braiding and Extensions", NameAM: "የጸጉር እንክብካቤ እና ቅባቶች", NameOM: "Kunuunsa fi Zayita Rifeensaa", NameTI: "ናይ ጸጉሪ ክንክንን ቅባታትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "natural-hair-growth-oils", NameEN: "Natural Rosemary & Castor Hair Growth Oils", NameAM: "የጸጉር ማሳደጊያ ተፈጥሯዊ ዘይቶች", NameOM: "Zayita Rifeensa Guddisu", NameTI: "ናይ ጸጉሪ መዕበዪ ዘይቲ"},
						{Slug: "braiding-human-hair-bundles", NameEN: "100% Human Hair Wigs, Weaves & Braids", NameAM: "ሰው ሰራሽ እና ተፈጥሯዊ ዊጎች", NameOM: "Rifeensa Nam-tolchee fi Uumamaa", NameTI: "ዊጋትን ናይ ጸጉሪ መመላኸዕን"},
						{Slug: "professional-hair-clippers-trimmers", NameEN: "Rechargeable Barbershop Hair Clippers & Trimmers", NameAM: "የወንዶች የጸጉር መላጫ ማሽኖች", NameOM: "Maashina Rifeensa Haadu", NameTI: "ናይ ጸጉሪ መላጸዪ ማሽን"},
					},
				},
				{
					Slug: "traditional-fragrances-incense", NameEN: "Ethiopian Incense, Itan and Perfumes", NameAM: "እጣን፣ ሱካር እና ባህላዊ ሽቶዎች", NameOM: "Ixaana fi Urgooftuu Aadaa", NameTI: "ዕጣንን ባህላዊ ሽቶን",
					Leaves: []TaxonomyLeaf{
						{Slug: "pure-ethiopian-frankincense-itan", NameEN: "Premium Tigray & Ogaden Frankincense (Itan / Lubanja)", NameAM: "የኦጋዴን እና የትግራይ ንጹህ እጣን", NameOM: "Ixaana Qulqulluu", NameTI: "ናይ ትግራይን ኦጋዴንን ጽሩይ ዕጣን"},
						{Slug: "traditional-scented-wood-chips", NameEN: "Scented Bukhur & Aromatic Wood Chips", NameAM: "ጥሩ መዓዛ ያላቸው የጢስ እንጨቶች", NameOM: "Muka Urgaa'u Aadaa", NameTI: "ምዑዝ ጨና ዘለዎ ዕንጨይቲ"},
						{Slug: "long-lasting-oil-perfumes", NameEN: "Alcohol-Free Long-Lasting Concentrated Attar Oils", NameAM: "ዘይት መሰል ረጅም ጊዜ የሚቆዩ ሽቶዎች", NameOM: "Urgooftuu Zayitaa Alkooolii Malee", NameTI: "ናይ ዘይቲ ሽቶታት"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 14: Health, Medical Equipment and Accessibility
		// =====================================================================
		{
			Slug: "health-medical",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "home-diagnostic-monitoring", NameEN: "Home Health Monitors and Diagnostics", NameAM: "የቤት ውስጥ የጤና መመርመሪያዎች", NameOM: "Meeshaa Qorannoo Fayyaa Manaa", NameTI: "ናይ ገዛ ጥዕና መመርመሪታት",
					Leaves: []TaxonomyLeaf{
						{Slug: "digital-blood-pressure-monitors", NameEN: "Digital Upper-Arm Blood Pressure Monitors", NameAM: "ዲጂታል የደም ግፊት መለኪያ", NameOM: "Safartuu Dhiibbaa Dhiigaa", NameTI: "ዲጂታል ናይ ደም ጸቕጢ መለክዒ"},
						{Slug: "glucometers-blood-sugar", NameEN: "Blood Glucose Meters and Test Strips", NameAM: "የስኳር መለኪያ መሳሪያዎች እና ስትሪፖች", NameOM: "Safartuu Sukkaara Dhiigaa", NameTI: "ናይ ሽኮር መለክዒ መሳርሒ"},
						{Slug: "pulse-oximeters-oxygen", NameEN: "Fingertip Pulse Oximeters & Digital Thermometers", NameAM: "የኦክስጅን እና የሙቀት መለኪያ", NameOM: "Safartuu Oksijiinii fi Ho'aa", NameTI: "ናይ ኦክስጅን መለክዒ"},
					},
				},
				{
					Slug: "mobility-accessibility-aids", NameEN: "Wheelchairs and Mobility Aids", NameAM: "ዊልቼሮች እና የመንቀሳቀሻ ድጋፎች", NameOM: "Barcuma Wiilcheerii fi Meeshaa Deeggarsaa", NameTI: "ዊልቸርን ናይ ምንቅስቓስ ሓገዝትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "folding-manual-wheelchairs", NameEN: "Foldable Lightweight Manual Wheelchairs", NameAM: "ተጣጣፊ ዊልቼሮች", NameOM: "Wiilcheerii Dachaafamu", NameTI: "ዝዕጸፍ ዊልቸር"},
						{Slug: "electric-powered-wheelchairs", NameEN: "Battery Powered Motorized Electric Wheelchairs", NameAM: "የኤሌክትሪክ ሞተር ዊልቼሮች", NameOM: "Wiilcheerii Elektiriikii", NameTI: "ናይ ኤሌክትሪክ ዊልቸር"},
						{Slug: "adjustable-walking-crutches", NameEN: "Adjustable Aluminum Crutches and Walking Canes", NameAM: "ክራንቾች እና ምርኩዞች", NameOM: "Ulee fi Qabduu Miilaa", NameTI: "ናይ ኢድ ምርኩዝን ክራንችን"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 15: Sports, Fitness and Outdoor Recreation
		// =====================================================================
		{
			Slug: "sports-fitness",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "gym-workout-equipment", NameEN: "Gym Machines and Free Weights", NameAM: "የጂም ማሽኖች እና ክብደቶች", NameOM: "Meeshaalee Jiimii fi Ulfaatinaa", NameTI: "ናይ ጂም መሳርሒታትን ክብደትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "motorized-running-treadmills", NameEN: "Folding Motorized Home Gym Treadmills", NameAM: "የኤሌክትሪክ የሩጫ ማሽኖች (ትሬድሚል)", NameOM: "Maashina Fiigichaa Jiimii", NameTI: "ናይ ሩጫ ማሽን (ትሬድሚል)"},
						{Slug: "adjustable-dumbbell-sets", NameEN: "Cast Iron & Rubber Encased Adjustable Dumbbells", NameAM: "የእጅ ክብደቶች (ደምበል)", NameOM: "Ulfaatina Harkaa (Dambelii)", NameTI: "ናይ ኢድ ክብደት (ደምበል)"},
						{Slug: "indoor-exercise-spin-bikes", NameEN: "Stationary Magnetic Exercise Spin Bikes", NameAM: "የቤት ውስጥ የአካል ብቃት ብስክሌቶች", NameOM: "Saayikilii Shaakala Jiimii", NameTI: "ናይ ውሽጢ ገዛ ስፖርት ብሽክለታ"},
					},
				},
				{
					Slug: "football-running-athletics", NameEN: "Football, Running and Team Sports", NameAM: "እግር ኳስ እና አትሌቲክስ", NameOM: "Kubbaa Miilaa fi Fiigicha", NameTI: "ኩዕሶ እግርን ኣትሌቲክስን",
					Leaves: []TaxonomyLeaf{
						{Slug: "ethiopian-national-team-jerseys", NameEN: "Ethiopian National Team & Club Football Jerseys", NameAM: "የዋልያዎቹ እና የክለብ ማሊያዎች", NameOM: "Uffata Garee Biyyooleessaa (Jeersii)", NameTI: "ናይ ዋልያታት ናይ ኩዕሶ ማልያ"},
						{Slug: "long-distance-running-shoes", NameEN: "Carbon Plate Long-Distance Marathon Running Shoes", NameAM: "የረጅም ርቀት ማራቶን የሩጫ ጫማዎች", NameOM: "Kophee Fiigicha Maaraatooniifi", NameTI: "ናይ ማራቶን ሩጫ ሳእኒ"},
						{Slug: "fifa-standard-footballs", NameEN: "Match Standard Size 5 Leather Footballs", NameAM: "ሳይዝ 5 የእግር ኳሶች", NameOM: "Kubbaa Miilaa Qulqulluu", NameTI: "ናይ እግሪ ኩዕሶ"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 21: Professional and Personal Services
		// =====================================================================
		{
			Slug: "professional-services",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "legal-audit-consulting", NameEN: "Legal, Tax and Accounting Services", NameAM: "የህግ፣ የግብር እና የሂሳብ አገልግሎቶች", NameOM: "Tajaajila Seeraa, Gibiraa fi Herregaa", NameTI: "ናይ ሕጊ፣ ግብሪን ሒሳብን ኣገልግሎት",
					Leaves: []TaxonomyLeaf{
						{Slug: "court-litigation-contract-drafting", NameEN: "Court Litigation & Commercial Contract Drafting", NameAM: "የውል ዝግጅት እና የፍርድ ቤት ጠበቃ", NameOM: "Waliigaltee fi Falmii Mana Murtii", NameTI: "ናይ ውዕሊ ምድላውን ናይ ቤት ፍርዲ ጠበቓን"},
						{Slug: "tax-filing-erm-accounting", NameEN: "Monthly Tax Declaration, Bookkeeping & Audit", NameAM: "የግብር ሂሳብ አሰራር እና ኦዲት", NameOM: "Herrega Gibiraa fi Oodiitii", NameTI: "ናይ ግብሪ ሒሳብን ኦዲትን"},
						{Slug: "trademark-business-registration", NameEN: "Ministry of Trade Licensing & Trademark Registration", NameAM: "የንግድ ፈቃድ ማውጣት እና ንግድ ምልክት ምዝገባ", NameOM: "Galmee Daldalaa fi Mallattoo", NameTI: "ናይ ንግዲ ፍቓድን ምልክትን ምዝገባ"},
					},
				},
				{
					Slug: "home-maintenance-repair", NameEN: "Plumbing, Electrical and Home Repairs", NameAM: "የቧንቧ፣ የኤሌክትሪክ እና የቤት ጥገና", NameOM: "Suphaa Ujummoo, Elektiriikii fi Manaa", NameTI: "ናይ ሻምብቆ፣ ኤሌክትሪክን ገዛን ጽገና",
					Leaves: []TaxonomyLeaf{
						{Slug: "emergency-plumbing-drainage-service", NameEN: "Emergency Water Pipe Leak & Drainage Unblocking", NameAM: "የፈጣን ቧንቧ እና ፍሳሽ ጥገና", NameOM: "Suphaa Ujummoo Bishaanii fi Dhangala'aa", NameTI: "ናይ ሻምብቆ ማይ ጽገና"},
						{Slug: "residential-electrical-installation", NameEN: "Residential Wiring & Circuit Breaker Upgrades", NameAM: "የቤት ውስጥ ኤሌክትሪክ መስመር ዝርጋታ", NameOM: "Hojii Sarara Elektirikii Manaa", NameTI: "ናይ ኤሌክትሪክ መስመር ዝርግሓ"},
						{Slug: "professional-house-cleaning-fumigation", NameEN: "Deep House Cleaning, Sofa Wash & Pest Control", NameAM: "የቤት ጽዳት እና የተባይ ማጥፋት አገልግሎት", NameOM: "Qulqullina Manaa fi Dhabamsiisa Ilbiisaa", NameTI: "ናይ ገዛ ጽሬትን ሓሳኹ ምክልኻልን"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 34: Property for Rent and Lease
		// =====================================================================
		{
			Slug: "property-rent",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "residential-apartment-rentals", NameEN: "Apartment Rentals (Furnished & Unfurnished)", NameAM: "የሚከራዩ አፓርታማዎች", NameOM: "Kiraa Apaartaamaa", NameTI: "ዝካረዩ ኣፓርታማታት",
					Leaves: []TaxonomyLeaf{
						{Slug: "fully-furnished-expat-apartments", NameEN: "Fully Furnished Luxury Expat Apartments (Bole/Kazanchis)", NameAM: "ሙሉ ፈርኒቸር ያላቸው የቅንጦት አፓርታማዎች", NameOM: "Apaartaamaa Meeshaa Guutuu Qabu", NameTI: "ሙሉእ ፈርኒቸር ዘለዎ ኣፓርታማ"},
						{Slug: "unfurnished-condominium-rentals", NameEN: "Unfurnished 1 to 3 Bedroom Condominiums", NameAM: "ያልተሟሉ የኮንዶሚኒየም ቤቶች ኪራይ", NameOM: "Kiraa Kandoominiyeemii Meeshaa Malee", NameTI: "ናይ ኮንዶሚኒየም ክራይ"},
						{Slug: "serviced-studio-guest-apartments", NameEN: "Daily & Weekly Serviced Studio Apartments", NameAM: "የቀን እና የሳምንት ሰርቪስ አፓርታማዎች", NameOM: "Apaartaamaa Istudiyoo Guyyaa", NameTI: "ናይ መዓልቲ ሰርቪስ ኣፓርታማ"},
					},
				},
				{
					Slug: "villa-compound-rentals", NameEN: "Standalone Residential Villas for Rent", NameAM: "የሚከራዩ ቪላዎች", NameOM: "Kiraa Viillaa Mana Jireenyaa", NameTI: "ዝካረዩ ቪላታት",
					Leaves: []TaxonomyLeaf{
						{Slug: "large-compound-embassy-villas", NameEN: "Large Compound Villas for Embassies & NGOs", NameAM: "ሰፊህ ግቢ ያላቸው ለድርጅቶች የሚሆኑ ቪላዎች", NameOM: "Viillaa Mooraa Bal'aa", NameTI: "ሰፊሕ ግቢ ዘለዎ ቪላ"},
						{Slug: "residential-family-villas", NameEN: "3 to 5 Bedroom Family Villas (CMC, Ayat, Lebu)", NameAM: "የመኖሪያ ቪላ ቤቶች", NameOM: "Viillaa Jireenya Maatii", NameTI: "ናይ ስድራቤት ቪላ"},
					},
				},
				{
					Slug: "commercial-retail-office-rentals", NameEN: "Commercial Office Spaces and Retail Shops", NameAM: "የንግድ ቢሮዎች እና ሱቆች ኪራይ", NameOM: "Kiraa Biiroo fi Dukaanaa", NameTI: "ናይ ንግዲ ቢሮታትን ድኳናትን ክራይ",
					Leaves: []TaxonomyLeaf{
						{Slug: "open-space-office-floors", NameEN: "Open-Plan Corporate Office Floors (Per m²)", NameAM: "ክፍት የቢሮ ወለሎች (በካሬ ሜትር)", NameOM: "Bakka Biiroo Bal'aa", NameTI: "ክፉት ናይ ቢሮ ቦታ"},
						{Slug: "mall-ground-floor-retail-shops", NameEN: "Mall Ground Floor Retail & Showroom Spaces", NameAM: "የሞል መሬት ወለል ሱቆች", NameOM: "Dukaana Mooolii Lafoo", NameTI: "ናይ ሞል መሬት ወለል ድኳን"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 38: Motorcycles, Scooters and Three-Wheelers
		// =====================================================================
		{
			Slug: "motorcycles-bajaj",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "three-wheeler-bajaj", NameEN: "Passenger and Cargo Bajaj Three-Wheelers", NameAM: "ተሳፋሪ እና የጭነት ባጃጆች", NameOM: "Baajaajii Fe'iinsaa fi Namaa", NameTI: "ናይ ተሳፋሪን ጽዕነትን ባጃጃት",
					Leaves: []TaxonomyLeaf{
						{Slug: "4-stroke-passenger-bajaj", NameEN: "4-Stroke Gasoline Passenger Bajaj (RE 205cc)", NameAM: "ባለ 4-ስትሮክ የመንገደኛ ባጃጅ", NameOM: "Baajaajii Imaltootaa 4-Stroke", NameTI: "ባዕል 4-ስትሮክ ናይ ተሳፈርቲ ባጃጅ"},
						{Slug: "cargo-delivery-three-wheelers", NameEN: "Heavy Duty Open Tray Cargo Three-Wheelers (Maxima)", NameAM: "የጭነት ባጃጆች", NameOM: "Baajaajii Fe'iinsa Meeshaa", NameTI: "ናይ ጽዕነት ባጃጅ"},
						{Slug: "electric-solar-tricycles", NameEN: "Electric Battery Powered Tricycles (Eco Bajaj)", NameAM: "የኤሌክትሪክ ባጃጆች", NameOM: "Baajaajii Elektiriikii", NameTI: "ናይ ኤሌክትሪክ ባጃጅ"},
					},
				},
				{
					Slug: "two-wheeler-motorcycles", NameEN: "Motorcycles, Scooters and Dirt Bikes", NameAM: "ሞተር ሳይክሎች እና ስኩተሮች", NameOM: "Moophedii fi Biskileetii Mootoraa", NameTI: "ሞተር ሳይክላትን ስኩተራትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "food-delivery-motorbikes", NameEN: "125cc-150cc Fuel Efficient Delivery Motorbikes", NameAM: "የዴሊቨሪ ሞተር ሳይክሎች", NameOM: "Moophedii Geejjiba Saffisaa", NameTI: "ናይ ዴሊቨሪ ሞተር ሳይክል"},
						{Slug: "offroad-dirt-trail-bikes", NameEN: "250cc Dual-Sport & Trail Dirt Motorcycles", NameAM: "የገጠር እና የሸካራ መንገድ ሞተሮች", NameOM: "Moophedii Daandii Baadiyyaa", NameTI: "ናይ ሸካር መገዲ ሞተር"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 41: Vehicle Parts, Tires and Accessories
		// =====================================================================
		{
			Slug: "vehicle-parts-accessories",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "tires-wheels-rims", NameEN: "Vehicle Tires, Tubes and Alloy Rims", NameAM: "የመኪና ጎማዎች እና ሪሞች", NameOM: "Gommaa fi Riimii Konkolaataa", NameTI: "ናይ መኪና ጎማታትን ሪማትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "passenger-car-radial-tires", NameEN: "Tubeless Radial Tires for Sedans (R14 to R18)", NameAM: "የአውቶሞቢል ራዲያል ጎማዎች", NameOM: "Gommaa Konkolaataa Xixinnoo", NameTI: "ናይ ንኣሽቱ መካይን ጎማ"},
						{Slug: "heavy-truck-trailer-tires", NameEN: "Heavy Truck Steel Radial Tires (12.00R20, 315/80R22.5)", NameAM: "የጭነት መኪና ጎማዎች", NameOM: "Gommaa Konkolaataa Fe'iinsaa", NameTI: "ናይ ጽዕነት መኪና ጎማ"},
						{Slug: "aluminum-alloy-car-rims", NameEN: "Custom Aluminum Alloy Sport Rims (Set of 4)", NameAM: "የአሉሚኒየም ስፖርት ሪሞች", NameOM: "Riimii Ispoortii Aluminiyeemii", NameTI: "ናይ ኣሉሚንየም ሪም"},
					},
				},
				{
					Slug: "engine-transmission-parts", NameEN: "Engine Components, Filters and Batteries", NameAM: "የሞተር መለዋወጫዎች፣ ባትሪዎች እና ፊልተሮች", NameOM: "Qaamolee Mootoraa fi Baatirii", NameTI: "ናይ ሞተር መለዋወጥን ባትሪን",
					Leaves: []TaxonomyLeaf{
						{Slug: "dry-charge-automotive-batteries", NameEN: "Maintenance-Free 12V 60Ah to 150Ah Automotive Batteries", NameAM: "የመኪና ባትሪዎች", NameOM: "Baatirii Konkolaataa", NameTI: "ናይ መኪና ባትሪ"},
						{Slug: "oil-fuel-air-filters", NameEN: "Original Oil, Fuel, and Air Filter Elements", NameAM: "የዘይት፣ የነዳጅ እና የንፋስ ፊልተሮች", NameOM: "Filtarii Zayitaa fi Boba'aa", NameTI: "ናይ ዘይትን ነዳድን ፊልተር"},
						{Slug: "brake-pads-shoes-rotors", NameEN: "Ceramic Brake Pads and Disc Brake Rotors", NameAM: "የፍሬን ፓዶች እና ዲስኮች", NameOM: "Paadii Firiinii", NameTI: "ናይ ፍሬን ፓድስ"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 47: Food, Groceries and Agricultural Produce
		// =====================================================================
		{
			Slug: "food-groceries",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "ethiopian-spices-condiments", NameEN: "Traditional Spices, Berbere and Shiro", NameAM: "ባህላዊ ቅመሞች፣ በርበሬ እና ሽሮ", NameOM: "Qooraa, Barbaree fi Shiroo", NameTI: "ባህላዊ ቀመማት፣ በርበረን ሽሮን",
					Leaves: []TaxonomyLeaf{
						{Slug: "authentic-red-berbere-powder", NameEN: "Authentic Sun-Dried Ethiopian Red Pepper (Berbere)", NameAM: "የተፈጨ ንጹህ የሀበሻ በርበሬ", NameOM: "Barbaree Diimaa Daakame", NameTI: "ዝተፈጨ ቀይሕ በርበረ"},
						{Slug: "roasted-spiced-shiro-powder", NameEN: "Roasted Chickpea Shiro Powder (Mitmita / Shiro)", NameAM: "የተፈተገ እና የተዘጋጀ የሽሮ ዱቄት", NameOM: "Daakuu Shiroo Qophaa'e", NameTI: "ዝተዳለወ ናይ ሽሮ ሓሩጭ"},
						{Slug: "hot-mitmita-spice-blend", NameEN: "Fiery Orange Mitmita Spice Blend", NameAM: "የሚጥሚጣ ቅመም", NameOM: "Qooraa Mitmixaa", NameTI: "ናይ ሚጥሚጣ ቀመም"},
						{Slug: "pure-ethiopian-kibe-butter", NameEN: "Clarified Traditional Spiced Butter (Niter Kibe)", NameAM: "የአገር ቤት ንጹህ የቅቤ ቅባት", NameOM: "Dhadhaa Qulqulluu Aadaa", NameTI: "ጽሩይ ናይ ዓዲ ጠስሚ"},
					},
				},
				{
					Slug: "bulk-groceries-edible-oils", NameEN: "Edible Oils, Sugar and Bulk Groceries", NameAM: "የምግብ ዘይት፣ ስኳር እና ሸቀጦች", NameOM: "Zayita Nyaataa fi Sukkaara", NameTI: "ናይ መግቢ ዘይትን ሽኮርን",
					Leaves: []TaxonomyLeaf{
						{Slug: "sunflower-edible-cooking-oil", NameEN: "Refined Sunflower & Palm Edible Cooking Oil (5L-20L)", NameAM: "የምግብ ማብሰያ ዘይት (5ሊ - 20ሊ)", NameOM: "Zayita Nyaataa Bilcheessaa", NameTI: "ናይ መግቢ ዘይቲ"},
						{Slug: "pure-natural-honey-gondar", NameEN: "100% Raw White & Amber Natural Honey (Gondar / Tigray)", NameAM: "የጎንደር እና የትግራይ ንጹህ ማር", NameOM: "Damma Uumamaa Qulqulluu", NameTI: "ጽሩይ ናይ ዓዲ መዓር"},
						{Slug: "fresh-red-onions-potatoes", NameEN: "Bulk Red Onions and Seed Potatoes (Per Quintal)", NameAM: "ቀይ ሽንኩርት እና ድንች (በኩንታል)", NameOM: "Qullubbii Diimaa fi Dinnicha", NameTI: "ቀይሕ ሽጉርትን ድንችን (ብኩንታል)"},
					},
				},
			},
		},
	}
}
