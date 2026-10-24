package main

import "fmt"

// GetCatalogDomains11To20 returns comprehensive domain catalogues for roots 11 to 20
func GetCatalogDomains11To20() []DomainCatalogDef {
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

	// Domain: shoes-bags
	list = append(list, DomainCatalogDef{
		RootSlug: "shoes-bags",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "mens-leather-formal-casual-shoes",
				NameEN: "Men's Formal Dress & Casual Leather Shoes",
				NameAM: "የወንዶች የቆዳ የቢሮ እና ተራ ጫማዎች",
				NameOM: "Kophee Gogaa Dhiiraa Hojii",
				NameTI: "ናይ ደቂ ተባዕትዮ ናይ ቆርበት ጫማ",
				Items: []LeafItemDef{
					makeLeaf("mens-pure-leather-oxford-shoes", "100% Genuine Ethiopian Leather Formal Oxford Dress Shoes", "የወንዶች ንጹህ የቆዳ የቢሮ ጫማዎች", "Kophee Gogaa Qulqulluu Hojii Dhiiraa", "ናይ 100% Genuine Ethiopian Leather Formal Oxford Dress Shoes"),
					makeLeaf("slip-on-leather-loafer-shoes", "Handcrafted Leather Slip-On Moccasin Driving Loafers", "ቀላል የሚደረጉ የቆዳ ጫማዎች (Loafers)", "Kophee Gogaa Salphaa (Loafer)", "ናይ Handcrafted Leather Slip-On Moccasin Driving Loafers"),
					makeLeaf("heavy-duty-leather-work-boots", "Steel-Toe Waterproof Heavy-Duty Construction Work Boots", "የብረት አፍ የደህንነት የስራ ቦቶች (Safety Boots)", "Kophee Nageenyaa Hojii (Safety)", "ናይ Steel-Toe Waterproof Heavy-Duty Construction Work Boots"),
					makeLeaf("cushioned-running-walking-sneakers", "Breathable Mesh Air-Cushioned Athletic Running Sneakers", "ስፖርታዊ የሩጫ እና የእግር ጉዞ እስኒከሮች", "Kophee Fiigichaa fi Ispoortii", "ናይ Breathable Mesh Air-Cushioned Athletic Running Sneakers"),
					makeLeaf("casual-streetwear-skate-sneakers", "Durable Canvas & Suede Streetwear Casual Sneakers", "ዘመናዊ የወጣቶች ተራ እስኒከሮች", "Kophee Isniikara Dargaggootaa", "ናይ Durable Canvas & Suede Streetwear Casual Sneakers"),
					makeLeaf("genuine-leather-summer-sandals", "Adjustable Strap Genuine Leather Comfort Summer Sandals", "የወንዶች የቆዳ ክፍት ጫማዎች (ሰንደል)", "Kophee Banaa Gogaa Dhiiraa", "ናይ Adjustable Strap Genuine Leather Comfort Summer Sandals"),
					makeLeaf("chelsea-ankle-boots-mens-leather", "Handcrafted Pull-On Elastic Side Leather Chelsea Boots", "የወንዶች ቼልሲ የቆዳ ቦቶች (Chelsea Boots)", "Bootii Gogaa Chelsea Dhiiraa", "ናይ Handcrafted Pull-On Elastic Side Leather Chelsea Boots"),
					makeLeaf("classic-brogue-wingtip-leather-shoes", "Wingtip Perforated Classic Brogue Formal Dress Shoes", "የብሮግ ጥልፍ ዲዛይን ያላቸው የወንዶች ጫማዎች", "Kophee Gogaa Brogue Dhiiraa", "ናይ Wingtip Perforated Classic Brogue Formal Dress Shoes"),
				},
			},
			{
				Slug:   "womens-footwear-heels-sandals",
				NameEN: "Women's Heels, Flats, Boots and Bridal Shoes",
				NameAM: "የሴቶች ተረከዝ ጫማዎች፣ ፍላት፣ ቦቶች እና ሰንደሎች",
				NameOM: "Kophee Qeensaa, Banaa fi Kaalsii Dubartootaa",
				NameTI: "ናይ ደቂ ኣንስትዮ ተረከዝ ጫማን ሰንደልን",
				Items: []LeafItemDef{
					makeLeaf("pointed-toe-high-stiletto-heels", "Classic Pointed-Toe Elegant High Stiletto Party Heels", "የሴቶች ረጃጅም የድግስ ተረከዝ ጫማዎች", "Kophee Qeenxee Dheeraa Cidhaa", "ናይ Classic Pointed-Toe Elegant High Stiletto Party Heels"),
					makeLeaf("block-heel-comfortable-office-pumps", "Mid-Block Heel Comfortable Everyday Office Pumps", "መካከለኛ ተረከዝ ያላቸው የቢሮ ጫማዎች", "Kophee Hojii Dubartootaa Qananii", "ናይ Mid-Block Heel Comfortable Everyday Office Pumps"),
					makeLeaf("cushioned-ballerina-flat-shoes", "Soft Memory Foam Ballerina Flat Slip-On Shoes", "ፍላት ለስላሳ ጫማዎች", "Kophee Dirraamaa Salphaa (Flats)", "ናይ Soft Memory Foam Ballerina Flat Slip-On Shoes"),
					makeLeaf("ankle-length-leather-chelsea-boots", "Genuine Leather Ankle-High Elastic Chelsea Boots", "የሴቶች አጫጭር የቆዳ ቦቶች (Chelsea Boots)", "Bootii Gogaa Gabaabaa Dubartootaa", "ናይ Genuine Leather Ankle-High Elastic Chelsea Boots"),
					makeLeaf("embellished-wedding-bridal-sandals", "Crystal Embellished Open-Toe Bridal Wedding Sandals", "የሰርግ እና የድግስ ክፍት ጫማዎች", "Kophee Banaa Cidhaa Faayame", "ናይ Crystal Embellished Open-Toe Bridal Wedding Sandals"),
					makeLeaf("orthopedic-walking-slip-on-shoes", "Arch Support Lightweight Orthopedic Comfort Shoes", "ለእግር ምቾት ተስማሚ ኦርቶፔዲክ ጫማዎች", "Kophee Fayyaa Miilaa", "ናይ Arch Support Lightweight Orthopedic Comfort Shoes"),
					makeLeaf("knee-high-fashion-riding-boots", "Side-Zip Genuine Leather Knee-High Fashion Riding Boots", "እስከ ጉልበት የሚደርሱ ረዣዥም የሴቶች ቦቶች", "Bootii Dheeraa Jilbaa Dubartootaa", "ናይ Side-Zip Genuine Leather Knee-High Fashion Riding Boots"),
					makeLeaf("casual-wedge-platform-summer-sandals", "Cork Sole Espadrille Wedge Platform Casual Summer Sandals", "ከፍ ያሉ የፕላትፎርም ሰንደሎች (Wedges)", "Kophee Banaa Weejii (Wedge)", "ናይ Cork Sole Espadrille Wedge Platform Casual Summer Sandals"),
				},
			},
			{
				Slug:   "handbags-totes-backpacks-luggage",
				NameEN: "Handbags, Travel Backpacks and Luggage Sets",
				NameAM: "የእጅ ቦርሳዎች፣ የጉዞ ቦርሳዎች እና ሻንጣዎች",
				NameOM: "Boorsaa Harkaa, Imalaa fi Saanduqa",
				NameTI: "ናይ ኢድ ቦርሳታትን ናይ መገሻ ሻንጣታትን",
				Items: []LeafItemDef{
					makeLeaf("genuine-ethiopian-leather-tote-bags", "Handcrafted Full-Grain Ethiopian Leather Shoulder Tote Bags", "የሴቶች ትላልቅ ንጹህ የቆዳ የእጅ ቦርሳዎች", "Boorsaa Gogaa Qulqulluu Dubartootaa", "ናይ Handcrafted Full-Grain Ethiopian Leather Shoulder Tote Bags"),
					makeLeaf("luxury-designer-crossbody-bags", "Gold Hardware Chain Strap Luxury Crossbody Evening Bags", "የቅንጦት አነስተኛ የትከሻ ቦርሳዎች", "Boorsaa Qananii Sarara Qabu", "ናይ Gold Hardware Chain Strap Luxury Crossbody Evening Bags"),
					makeLeaf("anti-theft-travel-business-backpacks", "Waterproof USB Charging Anti-Theft Business Travel Backpacks", "ስርቆት መከላከያ ያላቸው የጉዞ እና የቢሮ ቦርሳዎች", "Boorsaa Dugdaa Hattuuraa Eegu", "ናይ Waterproof USB Charging Anti-Theft Business Travel Backpacks"),
					makeLeaf("hardshell-polycarbonate-spinner-suitcases", "3-Piece Polycarbonate 360-Degree Spinner Luggage Suitcase Sets", "የጉዞ የፕላስቲክ ሻንጣዎች ስብስብ (3-Piece)", "Saanduqa Imalaa Naanna'u (Suitcase)", "ናይ 3-Piece Polycarbonate 360-Degree Spinner Luggage Suitcase Sets"),
					makeLeaf("waterproof-sports-gym-duffel-bags", "Multi-Compartment Wet/Dry Shoe Pocket Sports Gym Duffel Bags", "የስፖርት እና የጂም ቦርሳዎች ከጫማ ማስቀመጫ ጋር", "Boorsaa Ispoortii fi Jiimii", "ናይ Multi-Compartment Wet/Dry Shoe Pocket Sports Gym Duffel Bags"),
					makeLeaf("rfid-blocking-leather-bifold-wallets-unisex", "Genuine Leather Multi-Card RFID Protected Pocket Wallets", "የቆዳ የኪስ ቦርሳዎች (RFID መከላከያ)", "Boorsaa Kiisii Gogaa (Wallet)", "ናይ Genuine Leather Multi-Card RFID Protected Pocket Wallets"),
					makeLeaf("convertible-laptop-messenger-briefcases", "Water-Resistant Shoulder Laptop Messenger Briefcase Bags", "የላፕቶፕ እና የሰነድ መያዣ የትከሻ ቦርሳዎች", "Boorsaa Hojii Laaptooppii fi Waraqaa", "ናይ Water-Resistant Shoulder Laptop Messenger Briefcase Bags"),
					makeLeaf("canvas-duffel-weekender-travel-bags", "Heavy Washed Canvas & Leather Trim Overnight Weekender Bags", "የሸራ እና የቆዳ የሳምንት መጨረሻ የጉዞ ቦርሳዎች", "Boorsaa Imalaa Shaaraa fi Gogaa", "ናይ Heavy Washed Canvas & Leather Trim Overnight Weekender Bags"),
				},
			},
		},
	})

	// Domain: jewelry-luxury
	list = append(list, DomainCatalogDef{
		RootSlug: "jewelry-luxury",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "gold-silver-traditional-jewelry",
				NameEN: "Gold, Silver and Traditional Ethiopian Jewelry",
				NameAM: "የወርቅ፣ የብር እና የባህል ጌጣጌጦች",
				NameOM: "Faaya Warqee, Meetaa fi Aadaa",
				NameTI: "ናይ ወርቅን ብሩርን ባህላዊ ስልማትን",
				Items: []LeafItemDef{
					makeLeaf("traditional-ethiopian-gold-cross-pendants", "21K Hand-Filigreed Ethiopian Traditional Gold Cross Pendants", "ባህላዊ የሀበሻ የወርቅ መስቀሎች (21 ካራት)", "Faaya Fannoo Warqee Aadaa (21K)", "ናይ 21K Hand-Filigreed Ethiopian Traditional Gold Cross Pendants"),
					makeLeaf("ethiopian-amber-lomi-silver-necklaces", "Authentic Silver Bead Traditional Lomi & Amber Bridal Necklaces", "የባህል የብር ሎሚ የአንገት ጌጦች", "Faaya Mormaa Meetaa fi Loomii Aadaa", "ናይ Authentic Silver Bead Traditional Lomi & Amber Bridal Necklaces"),
					makeLeaf("24k-pure-gold-wedding-bands", "Hallmarked Solid 24K Classic Curved Wedding Rings & Bands", "የ24 ካራት ንጹህ የወርቅ የሰርግ ቀለበቶች", "Qubee Cidhaa Warqee Qulqulluu (24K)", "ናይ Hallmarked Solid 24K Classic Curved Wedding Rings & Bands"),
					makeLeaf("925-sterling-silver-bracelets-chains", "Rhodium Plated 925 Italian Sterling Silver Curb Link Chains", "የ925 ንጹህ የብር ሰንሰለቶች እና አምባሮች", "Saamunaa Meetaa 925 Qulqulluu", "ናይ Rhodium Plated 925 Italian Sterling Silver Curb Link Chains"),
					makeLeaf("natural-opal-gemstone-jewelry", "Certified Natural Welo Ethiopian Fire Opal Rings & Pendants", "የወሎ የተፈጥሮ ኦፓል የከበረ ድንጋይ ጌጣጌጦች", "Faaya Dhagaa Gati-Jabeessa Ooppaalaa", "ናይ Certified Natural Welo Ethiopian Fire Opal Rings & Pendants"),
					makeLeaf("diamond-solitaire-engagement-rings", "VVS1 Clarity Certified Brilliant-Cut Solitaire Diamond Rings", "የዳይመንድ የጋብቻ ጥያቄ ቀለበቶች", "Qubee Kaadhimannaa Daayimondii", "ናይ VVS1 Clarity Certified Brilliant-Cut Solitaire Diamond Rings"),
					makeLeaf("solid-gold-bangles-bracelets", "Embossed Traditional Gold Cuff Bangles and Wristlets", "የወርቅ አምባሮች እና የእጅ ጌጦች", "Albee Warqee Harkaa", "ናይ Embossed Traditional Gold Cuff Bangles and Wristlets"),
					makeLeaf("pearl-necklace-earring-sets", "Cultured Freshwater White Pearl Necklace & Stud Earring Sets", "የእንቁ (ፐርል) የአንገት እና የጆሮ ጌጦች", "Faaya Luulii (Pearl) Qulqulluu", "ናይ Cultured Freshwater White Pearl Necklace & Stud Earring Sets"),
				},
			},
			{
				Slug:   "luxury-watches-timepieces",
				NameEN: "Luxury Swiss, Automatic and Quartz Watches",
				NameAM: "የቅንጦት የስዊስ እና አውቶማቲክ እጅ ሰዓቶች",
				NameOM: "Sa'aatii Qananii Harkaa",
				NameTI: "ናይ ቅንጦት ናይ ኢድ ሰዓታት",
				Items: []LeafItemDef{
					makeLeaf("swiss-automatic-chronograph-watches", "COSC Chronometer Certified Swiss Movement Stainless Steel Watches", "የስዊስ አውቶማቲክ ክሮኖግራፍ የቅንጦት እጅ ሰዓቶች", "Sa'aatii Harkaa Siwiisii Awutoomaatikii", "ናይ COSC Chronometer Certified Swiss Movement Stainless Steel Watches"),
					makeLeaf("classic-leather-strap-dress-watches", "Ultra-Slim Rose Gold Case Sapphire Crystal Men's Dress Watches", "የቆዳ ማሰሪያ ያላቸው የቢሮ እና የድግስ ሰዓቶች", "Sa'aatii Qalla'aa Gogaa Dhiiraa", "ናይ Ultra-Slim Rose Gold Case Sapphire Crystal Men's Dress Watches"),
					makeLeaf("heavy-stainless-steel-diver-watches", "300m Water-Resistant Ceramic Bezel Automatic Diver Watches", "ውሃ የማይገባቸው ጠንካራ የዳይቨር እጅ ሰዓቶች", "Sa'aatii Bishaan Keessaa Jabaa", "ናይ 300m Water-Resistant Ceramic Bezel Automatic Diver Watches"),
					makeLeaf("women-diamond-bezel-luxury-watches", "Mother-of-Pearl Dial Diamond Bezel Gold Women's Watches", "የሴቶች የወርቅ እና የዳይመንድ የቅንጦት ሰዓቶች", "Sa'aatii Qananii Dubartootaa", "ናይ Mother-of-Pearl Dial Diamond Bezel Gold Women's Watches"),
					makeLeaf("vintage-mechanical-pocket-watches", "Skeleton Roman Numeral Open-Face Mechanical Pocket Watches", "የኪስ ሰዓቶች ከነሰንሰለታቸው", "Sa'aatii Kiisii Aadaa Saamunaa", "ናይ Skeleton Roman Numeral Open-Face Mechanical Pocket Watches"),
					makeLeaf("solar-powered-perpetual-calendar-watches", "Eco-Drive Titanium Tough Radio-Controlled Solar Wristwatches", "በፀሐይ ብርሃን የሚሰሩ ዘመናዊ ሰዓቶች", "Sa'aatii Aduudhaan Hojjetu", "ናይ Eco-Drive Titanium Tough Radio-Controlled Solar Wristwatches"),
					makeLeaf("watch-repair-tools-leather-straps", "Watchmakers Spring Bar Remover Kits & Replacement Leather Bands", "የሰዓት መለዋወጫ ማሰሪያዎች እና የጥገና እቃዎች", "Meeshaa Suphaa fi Qabattoo Sa'aatii", "ናይ Watchmakers Spring Bar Remover Kits & Replacement Leather Bands"),
					makeLeaf("wood-glass-watch-collector-boxes", "12-Slot Velvet Lined High-Gloss Wooden Watch Display Cases", "ሰዓት ማስቀመጫ የቅንጦት የእንጨት ሳጥኖች", "Saanduqa Sa'aatii Kuufamu", "ናይ 12-Slot Velvet Lined High-Gloss Wooden Watch Display Cases"),
				},
			},
		},
	})

	// Domain: beauty-personal-care
	list = append(list, DomainCatalogDef{
		RootSlug: "beauty-personal-care",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "skincare-facial-treatments",
				NameEN: "Skincare, Serums and Sunscreens",
				NameAM: "የፊት እና የቆዳ እንክብካቤ፣ ክሬሞች እና ሰነዶች",
				NameOM: "Kunsuunsa Gogaa fi Fuulaa",
				NameTI: "ናይ ገጽን ቆርበትን ምክንኻን",
				Items: []LeafItemDef{
					makeLeaf("hyaluronic-acid-hydrating-serums", "Intensive Moisture Plumping Pure Hyaluronic Acid Serums", "ቆዳን የሚያረሰርሱ የሃያሉሮኒክ አሲድ ሴረሞች", "Seeramii Fuulaa Jiidhaa (Hyaluronic)", "ናይ Intensive Moisture Plumping Pure Hyaluronic Acid Serums"),
					makeLeaf("broad-spectrum-spf-sunscreen-lotions", "SPF 50+ PA++++ Non-Greasy Invisible Broad Spectrum Sunscreens", "የፀሐይ መከላከያ ክሬሞች (SPF 50+ Sunscreen)", "Kiriimii Aduuraa Eegu (SPF 50+)", "ናይ SPF 50+ PA++++ Non-Greasy Invisible Broad Spectrum Sunscreens"),
					makeLeaf("vitamin-c-brightening-facial-serums", "Dark Spot Correcting Stabilized 20% Vitamin C Facial Serums", "ጠቆርቆር ያለ ቆዳን የሚያቀሉ የቫይታሚን ሲ ሴረሞች", "Seeramii Fuulaa Ifsu (Vitamin C)", "ናይ Dark Spot Correcting Stabilized 20% Vitamin C Facial Serums"),
					makeLeaf("salicylic-acid-anti-acne-cleansers", "Pore-Clearing Gentle BHA Salicylic Acid Foam Face Washes", "ብጉር እና ቆሻሻን የሚያጸዱ የፊት ሳሙናዎች", "Saamunaa Finnee Fuulaa Qulqulleessu", "ናይ Pore-Clearing Gentle BHA Salicylic Acid Foam Face Washes"),
					makeLeaf("retinol-anti-aging-night-creams", "Collagen Boosting Fine Line Smoothing Night Face Creams", "የእድሜ መግፋት ምልክቶችን የሚከላከሉ የሌሊት ክሬሞች", "Kiriimii Halkanii Fuulaa (Retinol)", "ናይ Collagen Boosting Fine Line Smoothing Night Face Creams"),
					makeLeaf("soothing-centella-barrier-repair-moisturizers", "Cica Ceramide Deep Moisture Barrier Repairing Creams", "የተጎዳ ቆዳን የሚያክሙ የሴራማይድ ሞይስቸራይዘሮች", "Kiriimii Gogaa Miidhame Suphu", "ናይ Cica Ceramide Deep Moisture Barrier Repairing Creams"),
					makeLeaf("exfoliating-glycolic-aha-facial-toners", "7% Glycolic Acid Clarifying Brightening Skin Toners", "የቆዳን የላይኛውን ክፍል የሚያድሱ ቶነሮች", "Toonarii Fuulaa Qulqulleessu", "ናይ 7% Glycolic Acid Clarifying Brightening Skin Toners"),
					makeLeaf("collagen-hydrogel-under-eye-patches", "Puffiness & Dark Circle Reducing Under-Eye Moisture Masks", "ከዓይን ስር ያለ ጥቁረት የሚያጠፉ አይ ፓቾች", "Maaskii Ija Jala Kiriimii", "ናይ Puffiness & Dark Circle Reducing Under-Eye Moisture Masks"),
				},
			},
			{
				Slug:   "hair-care-traditional-oils",
				NameEN: "Hair Care, Braiding Extensions and Traditional Oils",
				NameAM: "የጸጉር እንክብካቤ፣ የጸጉር ቅባቶች እና ዊጎች",
				NameOM: "Rifeensa, Zayita Aadaa fi Faaya",
				NameTI: "ናይ ጸጉሪ ምክንኻንን ቅብኣትን ዊግን",
				Items: []LeafItemDef{
					makeLeaf("traditional-ethiopian-qibe-butter-korerima", "Traditional Scented Butter & Herbal Hair Nourishing Ointments", "የጸጉር ባህላዊ የቅቤ እና የኮረሪማ ቅባቶች", "Dhadhaa Rifeensaa Aadaa", "ናይ Traditional Scented Butter & Herbal Hair Nourishing Ointments"),
					makeLeaf("pure-rosemary-mint-scalp-growth-oils", "Essential Hair Density & Scalp Strengthening Rosemary Biotin Drops", "ጸጉር የሚያበቅሉ የሮዝመሪ እና የባዮቲን ዘይቶች", "Zayita Rifeensa Biqilchu (Rosemary)", "ናይ Essential Hair Density & Scalp Strengthening Rosemary Biotin Drops"),
					makeLeaf("virgin-castor-oil-black-seed-elixirs", "Cold-Pressed Ethiopian Pure Castor & Black Seed Hair Oils", "ንጹህ የጉሎ እና የጥቁር አዝሙድ የጸጉር ዘይቶች", "Zayita Qulqulluu Guloofi Abasoodaa", "ናይ Cold-Pressed Ethiopian Pure Castor & Black Seed Hair Oils"),
					makeLeaf("synthetic-braiding-hair-kanekalon-packs", "Pre-Stretched Jumbo Kanekalon Braid Hair Extensions", "የሹሩባ መስሪያ ጥራት ያላቸው የካኔካሎን ጸጉሮች", "Rifeensa Shurrubbaa (Kanekalon)", "ናይ Pre-Stretched Jumbo Kanekalon Braid Hair Extensions"),
					makeLeaf("100-percent-human-hair-lace-front-wigs", "HD Transparent Lace Front Brazilian Virgin Human Hair Wigs", "ኦሪጅናል የሰው ጸጉር ዊጎች (Human Hair Wigs)", "Wiigii Rifeensa Namaa Qulqulluu", "ናይ HD Transparent Lace Front Brazilian Virgin Human Hair Wigs"),
					makeLeaf("deep-conditioning-keratin-hair-masks", "Salon Grade Damage Repair Argan & Keratin Treatment Masks", "የተጎዳ ጸጉርን የሚያክሙ የኬራቲን ማከሚያዎች", "Maaskii Rifeensaa Keeraatiinii", "ናይ Salon Grade Damage Repair Argan & Keratin Treatment Masks"),
					makeLeaf("sulfate-free-curl-defining-shampoos", "Moisturizing Sulfate-Free Shampoos & Leave-In Curl Conditioners", "ሰልፌት የሌለባቸው ለስላሳ ሻምፑዎች እና ኮንዲሽነሮች", "Shaampuu fi Kondiishinara Rifeensaa", "ናይ Moisturizing Sulfate-Free Shampoos & Leave-In Curl Conditioners"),
					makeLeaf("pro-ionic-hair-dryers-straighteners", "Ceramic Tourmaline Fast-Heating Flat Iron Hair Straighteners", "የጸጉር ማድረቂያ እና መተኮሻ ኃይለኛ ካውያዎች", "Kaawuyaa fi Gorsituu Rifeensaa", "ናይ Ceramic Tourmaline Fast-Heating Flat Iron Hair Straighteners"),
				},
			},
			{
				Slug:   "fragrances-oud-traditional-etotan",
				NameEN: "Perfumes, Ouds and Traditional Incense (Itan)",
				NameAM: "ሽቶዎች፣ ዑዶች እና ባህላዊ እጣኖች (እጣን)",
				NameOM: "Ushaa, Ixaana Aadaa fi Shittoo",
				NameTI: "ሽቶታትን ዑድን ባህላዊ ዕጣንን",
				Items: []LeafItemDef{
					makeLeaf("traditional-fragrant-ethiopian-itan-frankincense", "Ogaden & Tigray Premium Grade Church & Home White Frankincense", "የኦጋዴን እና የትግራይ ምርጥ ነጭ እጣኖች", "Ixaana Aadaa Qulqulluu", "ናይ Ogaden & Tigray Premium Grade Church & Home White Frankincense"),
					makeLeaf("traditional-clay-etotan-charcoal-burners", "Handcrafted Ornate Clay Incense Burners (Mebkera / Girgira)", "ባህላዊ የሸክላ የእጣን ማጨሻ ጊርጊራዎች", "Girgiraa Ixaanaa Suphee Aadaa", "ናይ Handcrafted Ornate Clay Incense Burners (Mebkera / Girgira)"),
					makeLeaf("concentrated-long-lasting-arabic-oud-perfumes", "Pure Agarwood Cambodian & Amber Concentrated Perfume Oils", "ለረጅም ጊዜ የሚቆዩ የአረብ ዑድ ሽቶዎች", "Shittoo Qananii Uudii (Oud)", "ናይ Pure Agarwood Cambodian & Amber Concentrated Perfume Oils"),
					makeLeaf("designer-french-eau-de-parfum-men", "Original 100ml Woody Aromatic Fresh Men's Eau de Parfum", "ኦሪጅናል የወንዶች የፈረንሳይ ሽቶዎች", "Shittoo Dhiiraa Qulqulluu 100ml", "ናይ Original 100ml Woody Aromatic Fresh Men's Eau de Parfum"),
					makeLeaf("floral-fruity-women-eau-de-parfum", "Sweet Vanilla & Jasmine Long-Lasting Women's Luxury Perfumes", "ጣፋጭ መአዛ ያላቸው የሴቶች ሽቶዎች", "Shittoo Dubartootaa Bareedaa", "ናይ Sweet Vanilla & Jasmine Long-Lasting Women's Luxury Perfumes"),
					makeLeaf("pocket-perfume-oils-roll-ons", "6ml Non-Alcoholic Concentrated Fragrance Attar Roll-Ons", "አልኮል የሌለባቸው አነስተኛ የኪስ ሽቶዎች (Attar)", "Shittoo Kiisii Alkoolii Malee (Attar)", "ናይ 6ml Non-Alcoholic Concentrated Fragrance Attar Roll-Ons"),
					makeLeaf("scented-bakhoor-wood-chips", "Sugar & Essential Oil Scented Arab Home Fragrance Bakhoor", "የቤት ውስጥ መአዛ ማጣፈጫ የባኩር እንጨቶች", "Qoraan Shittoo Manaa (Bakhoor)", "ናይ Sugar & Essential Oil Scented Arab Home Fragrance Bakhoor"),
					makeLeaf("electric-smokeless-incense-burners", "Thermostat Controlled Safe Electric Bakhoor Ceramic Heaters", "ጭስ የማያበዙ ዘመናዊ የኤሌክትሪክ እጣን ማጨሻዎች", "Gorsituu Ixaanaa Elektiriikii", "ናይ Thermostat Controlled Safe Electric Bakhoor Ceramic Heaters"),
				},
			},
		},
	})

	// Domain: health-medical
	list = append(list, DomainCatalogDef{
		RootSlug: "health-medical",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "diagnostic-home-medical-devices",
				NameEN: "Home Health Monitors and Diagnostic Devices",
				NameAM: "የቤት ውስጥ የጤና መመርመሪያ መሳሪያዎች",
				NameOM: "Meeshaalee Qorannoo Fayyaa",
				NameTI: "ናይ ገዛ ጥዕና መመርመሪ መሳርሒታት",
				Items: []LeafItemDef{
					makeLeaf("digital-arm-blood-pressure-monitors", "Clinical Accuracy Automatic Upper Arm Voice Blood Pressure Monitors", "የክንድ የደም ግፊት መለኪያ ዲጂታል መሳሪያዎች", "Meeshaa Safartuu Dhiibbaa Dhiigaa", "ናይ Clinical Accuracy Automatic Upper Arm Voice Blood Pressure Monitors"),
					makeLeaf("blood-glucose-monitoring-meter-kits", "Fasting & Post-Meal Blood Sugar Diabetes Test Meters with Strips", "የስኳር በሽታ መመርመሪያ መሳሪያዎች ከነስትሪፑ", "Meeshaa Safartuu Sukkaara Dhiigaa", "ናይ Fasting & Post-Meal Blood Sugar Diabetes Test Meters with Strips"),
					makeLeaf("fingertip-pulse-oximeters-spo2", "High Accuracy OLED SpO2 Blood Oxygen Saturation & Pulse Meters", "የኦክስጅን እና የልብ ምት መለኪያ የጣት መሳሪያዎች", "Meeshaa Safartuu Oksijiinii (SpO2)", "ናይ High Accuracy OLED SpO2 Blood Oxygen Saturation & Pulse Meters"),
					makeLeaf("contactless-infrared-forehead-thermometers", "Instant 1-Second Body & Object Temperature Infrared Scanners", "ያለ ንክኪ በ1 ሰከንድ የሰውነት ሙቀት መለኪያ ቴርሞሜትሮች", "Termoomeetira Hoo'a Qaamaa", "ናይ Instant 1-Second Body & Object Temperature Infrared Scanners"),
					makeLeaf("portable-mesh-ultrasonic-nebulizers", "Quiet Handheld USB Rechargeable Asthma Mist Inhaler Nebulizers", "የአስም እና የመተንፈሻ አካላት መድኃኒት መተንፈሻ ኔቡላይዘሮች", "Meeshaa Haphisa Qorichaa (Nebulizer)", "ናይ Quiet Handheld USB Rechargeable Asthma Mist Inhaler Nebulizers"),
					makeLeaf("medical-grade-fetal-doppler-monitors", "Pocket Prenatal Baby Heartbeat Sound Listening Doppler Ultrasound", "የጽንስ የልብ ምት ማዳመጫ ዶፕለር መሳሪያዎች", "Meeshaa Sagalee Onnee Da'imaa", "ናይ Pocket Prenatal Baby Heartbeat Sound Listening Doppler Ultrasound"),
					makeLeaf("multi-frequency-body-composition-scales", "Bluetooth Smartphone Sync Fat Muscle BMI Smart Bathroom Scales", "የሰውነት ስብ፣ ጡንቻ እና ክብደት መመዘኛ ስማርት ሚዛኖች", "Mizana Ulfaatinna Qaamaa Ismaartii", "ናይ Bluetooth Smartphone Sync Fat Muscle BMI Smart Bathroom Scales"),
					makeLeaf("ecg-handheld-heart-rhythm-monitors", "Direct Palm EKG Lead-I Arrhythmia Detection Recorders", "የልብ ምት ምልከታ ማድረጊያ በእጅ የሚያዙ የEKG መሳሪያዎች", "Meeshaa EKG Onnee Harkaa", "ናይ Direct Palm EKG Lead-I Arrhythmia Detection Recorders"),
				},
			},
			{
				Slug:   "mobility-orthopedic-rehabilitation",
				NameEN: "Wheelchairs, Crutches and Orthopedic Braces",
				NameAM: "የዊልቸር፣ ምርኩዝ እና ኦርቶፔዲክ ድጋፎች",
				NameOM: "Barcuma Naafaa, Ulee fi Gargaarsa Qaamaa",
				NameTI: "ናይ መንበር ስንኩላንን ምርኩዝን ኦርቶፔዲክን",
				Items: []LeafItemDef{
					makeLeaf("folding-lightweight-steel-wheelchairs", "Heavy-Duty Chrome Plated Manual Self-Propelled Folding Wheelchairs", "የሚታጠፉ ጠንካራ የብረት ዊልቸሮች", "Barcuma Naafaa Dachaafamu (Wheelchair)", "ናይ Heavy-Duty Chrome Plated Manual Self-Propelled Folding Wheelchairs"),
					makeLeaf("motorized-electric-power-wheelchairs", "Joystick Controlled Dual Motor Airline-Approved Electric Wheelchairs", "በጆይስቲክ የሚሰሩ የኤሌክትሪክ ዊልቸሮች", "Barcuma Naafaa Elektiriikii", "ናይ Joystick Controlled Dual Motor Airline-Approved Electric Wheelchairs"),
					makeLeaf("adjustable-aluminum-underarm-crutches", "Push-Button Height Adjustable Lightweight Underarm Crutches", "ከፍ እና ዝቅ የሚሉ የአሉሚኒየም ምርኩዞች (ክራንች)", "Ulee Hirkannoo Qaamaa (Crutches)", "ናይ Push-Button Height Adjustable Lightweight Underarm Crutches"),
					makeLeaf("quad-cane-four-legged-walking-sticks", "Non-Slip Rubber Footed Ergonomic Grip Quad Walking Canes", "አራት እግር ያላቸው የማይንሸራተቱ ምርኩዞች", "Ulee Deemsaa Miila Afurii", "ናይ Non-Slip Rubber Footed Ergonomic Grip Quad Walking Canes"),
					makeLeaf("folding-mobility-walkers-rollators", "Padded Seat Hand Brake 4-Wheeled Elderly Mobility Walkers", "መቀመጫ እና ብሬክ ያላቸው ባለ ጎማ የድጋፍ መራመጃዎች", "Meeshaa Deemsaa Gargaaru (Rollator)", "ናይ Padded Seat Hand Brake 4-Wheeled Elderly Mobility Walkers"),
					makeLeaf("hinged-knee-braces-ligament-support", "Dual Side Stabilizer ACL/MCL Post-Surgery Adjustable Knee Braces", "የጉልበት መገጣጠሚያ ማሰሪያ ጠንካራ ብሬሶች", "Qabduu Jilbaa Sibilaa (Knee Brace)", "ናይ Dual Side Stabilizer ACL/MCL Post-Surgery Adjustable Knee Braces"),
					makeLeaf("lumbar-back-support-posture-belts", "Dual Pull Breathable Lumbar Spine Decompression Belts", "ለወገብ እና ለጀርባ ህመም ማስታገሻ ጠንካራ ቀበቶዎች", "Qabattoo Mudhii Dugdaa", "ናይ Dual Pull Breathable Lumbar Spine Decompression Belts"),
					makeLeaf("cervical-neck-traction-collars", "Inflatable Neck Spinal Cervical Immobilization Pain Relief Collars", "የአንገት ህመም ማስታገሻ የድጋፍ ኮላሮች", "Qabduu Mormaa (Cervical Collar)", "ናይ Inflatable Neck Spinal Cervical Immobilization Pain Relief Collars"),
				},
			},
			{
				Slug:   "hospital-supplies-clinical-ppe",
				NameEN: "Clinical Supplies, First Aid and PPE",
				NameAM: "የሆስፒታል እቃዎች፣ የመጀመሪያ እርዳታ እና ጓንቶች",
				NameOM: "Meeshaalee Hospitaalaa fi Gargaarsa Jalqabaa",
				NameTI: "ናይ ሆስፒታል ንዋያትን ናይ ቀዳማይ ረድኤትን",
				Items: []LeafItemDef{
					makeLeaf("disposable-nitrile-examination-gloves", "Powder-Free Medical Grade Textured Fingertip Blue Nitrile Gloves", "የህክምና የናይትራይል ጓንቶች (100 ፍሬ በአንድ ፓኬት)", "Gilaavii Hospitaalaa Naayitiiraayilii", "ናይ Powder-Free Medical Grade Textured Fingertip Blue Nitrile Gloves"),
					makeLeaf("comprehensive-emergency-first-aid-kits", "Workplace & Vehicle Certified 120-Piece Emergency First Aid Bags", "የመጀመሪያ እርዳታ መስጫ የተሟሉ ኪቶች (First Aid Kit)", "Boorsaa Gargaarsa Jalqabaa Guutuu", "ናይ Workplace & Vehicle Certified 120-Piece Emergency First Aid Bags"),
					makeLeaf("3-ply-surgical-face-masks", "BFE 99% Meltblown Filter Certified Earloop Medical Face Masks", "የአፍ እና የአፍንጫ መሸፈኛ የህክምና ማስኮች", "Maaskii Afaanii fi Funyaanii", "ናይ BFE 99% Meltblown Filter Certified Earloop Medical Face Masks"),
					makeLeaf("medical-oxygen-concentrators-10l", "Continuous High-Purity 93% Molecular Sieve 10L Oxygen Machines", "ለታካሚ ኦክስጅን የሚያመርቱ 10 ሊትር ማሽኖች (Concentrator)", "Maashinii Oksijiinii Hospitaalaa 10L", "ናይ Continuous High-Purity 93% Molecular Sieve 10L Oxygen Machines"),
					makeLeaf("anti-bedsore-alternating-air-mattresses", "Quiet Medical Pump Ripple Pressure Sore Hospital Air Overlays", "የአልጋ ቁስለትን የሚከላከሉ አየር የሚነፉ የህክምና ፍራሾች", "Firaashii Qilleensaa Madaa Ittisu", "ናይ Quiet Medical Pump Ripple Pressure Sore Hospital Air Overlays"),
					makeLeaf("manual-multi-function-hospital-beds", "ABS Headboard Two-Crank Adjustable Nursing Clinic Patient Beds", "ለክሊኒክ እና ለሆስፒታል ታካሚ የሚሆኑ ባለ ሁለት ክራንክ አልጋዎች", "Siree Dhukkubsataa Hospitaalaa", "ናይ ABS Headboard Two-Crank Adjustable Nursing Clinic Patient Beds"),
					makeLeaf("portable-medical-suction-machines", "Oil-Free Lubrication High Vacuum Sputum Phlegm Suction Units", "የአክታ እና ፈሳሽ መምጠጫ የህክምና ማሽኖች (Suction)", "Maashinii Hancufaa Xuuxu (Suction)", "ናይ Oil-Free Lubrication High Vacuum Sputum Phlegm Suction Units"),
					makeLeaf("sterile-surgical-gauze-cotton-rolls", "100% Bleached Absorbent Hospital Surgical Cotton & Gauze Swabs", "ንጹህ የህክምና ጥጥ እና ጋዞች በጥቅል", "Jirbii fi Faashetaa Hospitaalaa", "ናይ 100% Bleached Absorbent Hospital Surgical Cotton & Gauze Swabs"),
				},
			},
		},
	})

	// Domain: sports-fitness
	list = append(list, DomainCatalogDef{
		RootSlug: "sports-fitness",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "gym-weight-training-equipment",
				NameEN: "Gym Weights, Dumbbells and Barbells",
				NameAM: "የጂም እቃዎች፣ ደምብሎች እና የክብደት ማንሻዎች",
				NameOM: "Meeshaalee Jiimii fi Ulfaatinaa",
				NameTI: "ናይ ጂም ንዋያትን ደምበላትን",
				Items: []LeafItemDef{
					makeLeaf("adjustable-quick-lock-dumbbell-pairs", "2.5kg to 24kg Single Handle Quick Turn Dial Adjustable Dumbbells", "ክብደታቸው የሚቀያየር ዘመናዊ ደምብሎች", "Dambalii Ulfaatina Jijjiiramu", "ናይ 2.5kg to 24kg Single Handle Quick Turn Dial Adjustable Dumbbells"),
					makeLeaf("solid-cast-iron-kettlebells", "Color Coded Vinyl Dipped 4kg to 32kg Strength Kettlebells", "የብረት ኬትልቤል የክብደት ማንሻዎች", "Kettlebell Sibilaa Ulfaatinaa", "ናይ Color Coded Vinyl Dipped 4kg to 32kg Strength Kettlebells"),
					makeLeaf("olympic-barbell-weight-plates-sets", "2-Inch Hole Solid Cast Iron & Rubber Bumper Weight Plates", "የኦሊምፒክ ባርቤል እና የክብደት ሰሃኖች (Plates)", "Balaqqaa Ulfaatinaa Oolimpikii", "ናይ 2-Inch Hole Solid Cast Iron & Rubber Bumper Weight Plates"),
					makeLeaf("multi-position-flat-incline-decline-benches", "Commercial 800lb Capacity Heavy Steel Adjustable Workout Benches", "የክብደት ማንሻ የሚታጠፉ ጠንካራ ወንበሮች (Bench)", "Barcuma Ispoortii Ulfaatinaa", "ናይ Commercial 800lb Capacity Heavy Steel Adjustable Workout Benches"),
					makeLeaf("heavy-duty-doorway-pull-up-bars", "No-Screw Multi-Grip Chin-Up Upper Body Exercise Door Bars", "በር ላይ የሚገጠሙ የፑል አፕ መስሪያ ጠንካራ ብረቶች", "Sibila Dhaabbii Pull-Up", "ናይ No-Screw Multi-Grip Chin-Up Upper Body Exercise Door Bars"),
					makeLeaf("natural-latex-resistance-loop-bands", "5-Level Heavy Tension Booty & Mobility Resistance Exercise Bands", "የስፖርት መወጠሪያ የላስቲክ ገመዶች (Resistance Bands)", "Gummata Ispoortii Qorannoo", "ናይ 5-Level Heavy Tension Booty & Mobility Resistance Exercise Bands"),
					makeLeaf("commercial-cable-cross-functional-trainers", "Dual Weight Stack Multi-Gym Cable Crossover Pulley Machines", "ሁሉንም የሰውነት ክፍል የሚያሰሩ ሁለገብ የጂም ማሽኖች", "Maashinii Jiimii Bal'aa (Cable Cross)", "ናይ Dual Weight Stack Multi-Gym Cable Crossover Pulley Machines"),
					makeLeaf("heavy-rubber-gym-floor-tiles", "50x50cm 20mm High-Density Impact Sound Absorbing Rubber Tiles", "የጂም ወለል ላይ የሚነጠፉ ወፍራም የላስቲክ ምንጣፎች", "Afata Dachee Jiimii Jabaa", "ናይ 50x50cm 20mm High-Density Impact Sound Absorbing Rubber Tiles"),
				},
			},
			{
				Slug:   "cardio-machines-endurance",
				NameEN: "Treadmills, Exercise Bikes and Cardio Gear",
				NameAM: "ትሬድሚሎች፣ የስፖርት ብስክሌቶች እና የካርዲዮ እቃዎች",
				NameOM: "Maashinii Fiigichaa fi Kaardiyoo",
				NameTI: "ትሬድሚልን ናይ ስፖርት ብሽክለታን",
				Items: []LeafItemDef{
					makeLeaf("heavy-motorized-folding-treadmills", "3.0HP Peak Auto-Incline Shock-Absorbing Running Treadmills", "በኤሌክትሪክ የሚሰሩ የሚታጠፉ የሩጫ ማሽኖች (Treadmill)", "Maashinii Fiigichaa Elektiriikii", "ናይ 3.0HP Peak Auto-Incline Shock-Absorbing Running Treadmills"),
					makeLeaf("magnetic-resistance-spin-exercise-bikes", "Flywheel Belt-Drive Silent Studio Cycling Indoor Bikes", "የቤት ውስጥ የክብደት መቀነሻ የስፖርት ብስክሌቶች (Spin Bike)", "Biskileetii Ispoortii Mana Keessaa", "ናይ Flywheel Belt-Drive Silent Studio Cycling Indoor Bikes"),
					makeLeaf("dual-action-elliptical-cross-trainers", "Smooth Magnetic Low-Impact Joint-Safe Elliptical Gliders", "የእግር እና የእጅ እንቅስቃሴ መስሪያ ኤሊፕቲካል ማሽኖች", "Maashinii Ispoortii Miilaa fi Harkaa", "ናይ Smooth Magnetic Low-Impact Joint-Safe Elliptical Gliders"),
					makeLeaf("water-resistance-wooden-rowing-machines", "Natural Oak Wood Real-Water Sound Aerobic Rowing Trainers", "የቀዘፋ ስፖርት ማሽኖች (Rowing Machine)", "Maashinii Rowiingii Bishaanii", "ናይ Natural Oak Wood Real-Water Sound Aerobic Rowing Trainers"),
					makeLeaf("steel-bearing-fast-jump-ropes", "Tangle-Free Adjustable Cable Speed Skipping Ropes", "የብረት ገመድ ያላቸው ፈጣን የስፖርት መዝለያ ገመዶች", "Haada Utaalchaa Saffisaa", "ናይ Tangle-Free Adjustable Cable Speed Skipping Ropes"),
					makeLeaf("high-density-eva-exercise-foam-rollers", "Deep Tissue Muscle Trigger Point Massage Recovery Rollers", "ጡንቻ ለማፍታት የሚያገለግሉ ፎም ሮለሮች", "Meeshaa Sukkuummaa Murree (Roller)", "ናይ Deep Tissue Muscle Trigger Point Massage Recovery Rollers"),
					makeLeaf("thick-anti-slip-yoga-pilates-mats", "10mm Eco-Friendly NBR Alignment Line Non-Slip Floor Mats", "የዮጋ እና የአካል ብቃት እንቅስቃሴ ምንጣፎች (Yoga Mat)", "Afata Yoogaa fi Ispoortii", "ናይ 10mm Eco-Friendly NBR Alignment Line Non-Slip Floor Mats"),
					makeLeaf("stepper-aerobic-step-platforms", "Height Adjustable Non-Skid Home Fitness Step Platforms", "የኤሮቢክስ ስፖርት መወጣጫ ሰሌዳዎች", "Boordii Ispoortii Eerobiiksii", "ናይ Height Adjustable Non-Skid Home Fitness Step Platforms"),
				},
			},
			{
				Slug:   "football-athletics-team-sports",
				NameEN: "Football, Athletics, Camping and Outdoor Recreation",
				NameAM: "የእግር ኳስ፣ የአትሌቲክስ እና የካምፒንግ እቃዎች",
				NameOM: "Kubbaa Miilaa, Atileetiksii fi Imalaa",
				NameTI: "ኩዕሶ እግርን ኣትሌቲክስን ናይ ካምፒንግ ንዋያትን",
				Items: []LeafItemDef{
					makeLeaf("fifa-quality-thermal-bonded-footballs", "Seamless Waterproof Official Size 5 Match Footballs", "ኦሪጅናል የፊፋ ጥራት ያላቸው የእግር ኳሶች (Size 5)", "Kubbaa Miilaa Qulqulluu (Size 5)", "ናይ Seamless Waterproof Official Size 5 Match Footballs"),
					makeLeaf("synthetic-turf-football-boots-cleats", "Agility TPU Stud Firm Ground Outdoor Soccer Cleats", "የሜዳ እና የሰው ሰራሽ ሳር የእግር ኳስ ጫማዎች", "Kophee Kubbaa Miilaa Qara", "ናይ Agility TPU Stud Firm Ground Outdoor Soccer Cleats"),
					makeLeaf("pro-goalkeeper-gloves-finger-spines", "4mm Contact Latex High-Grip Protection Goalkeeper Gloves", "የበረኛ ጓንቶች ከነጣት መከላከያው ጋር", "Gilaavii Eegduu Goolii (Goalkeeper)", "ናይ 4mm Contact Latex High-Grip Protection Goalkeeper Gloves"),
					makeLeaf("breathable-club-football-jerseys", "Dri-FIT Embroidered Club & National Team Football Kits", "የእግር ኳስ ክለቦች እና የብሔራዊ ቡድን ማሊያዎች", "Uffata Maaliyaa Kubbaa Miilaa", "ናይ Dri-FIT Embroidered Club & National Team Football Kits"),
					makeLeaf("shin-guards-compression-sleeves", "Impact Resistant Hard Shell Soccer Protective Shin Pads", "የእግር ኳስ የባት መከላከያ እቃዎች (Shin Guard)", "Eegduu Qal'oo Miilaa (Shin Guard)", "ናይ Impact Resistant Hard Shell Soccer Protective Shin Pads"),
					makeLeaf("automatic-pop-up-camping-tents", "4-6 Person Waterproof UV-Protected Instant Family Camping Tents", "የሚተከሉ የካምፒንግ እና የጉዞ ድንኳኖች", "Dunkaana Imalaa fi Bosonaa (Tent)", "ናይ 4-6 Person Waterproof UV-Protected Instant Family Camping Tents"),
					makeLeaf("waterproof-down-sleeping-bags", "Sub-Zero Mummy Lightweight Outdoor Trekking Sleeping Bags", "ብርድ የሚከላከሉ የጉዞ መኝታ ቦርሳዎች (Sleeping Bag)", "Boorsaa Rafuu Bosonaa (Sleeping Bag)", "ናይ Sub-Zero Mummy Lightweight Outdoor Trekking Sleeping Bags"),
					makeLeaf("solar-rechargeable-led-camping-lanterns", "Power Bank Output Collapsible Bright Tent Hanging Lanterns", "በፀሐይ የሚሞሉ የጉዞ እና የካምፕ ፋኖሶች", "Faanoosii Imalaa Aduudhaan Chaarjamu", "ናይ Power Bank Output Collapsible Bright Tent Hanging Lanterns"),
				},
			},
		},
	})

	// Domain: books-publishing
	list = append(list, DomainCatalogDef{
		RootSlug: "books-publishing",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "ethiopian-history-literature",
				NameEN: "Ethiopian History, Philosophy and Literature",
				NameAM: "የኢትዮጵያ ታሪክ፣ ፍልስፍና እና ስነ-ጽሁፍ መጻሕፍት",
				NameOM: "Kitaabolee Seenaa fi Aadaa Itoophiyaa",
				NameTI: "መጻሕፍቲ ታሪክን ስነ-ጽሑፍን ኢትዮጵያ",
				Items: []LeafItemDef{
					makeLeaf("ancient-ethiopian-history-chronicles", "Royal Chronicles, Axumite & Medieval Ethiopian Historical Studies", "የጥንታዊ ኢትዮጵያ ታሪክ እና የነገስታት ዜና መዋዕል መጻሕፍት", "Kitaaba Seenaa Itoophiyaa Durii", "ናይ Royal Chronicles, Axumite & Medieval Ethiopian Historical Studies"),
					makeLeaf("classic-amharic-fictional-novels", "Timeless Amharic Literary Masterpieces & Romantic Fiction", "ታዋቂ የአማርኛ ልብወለድ እና ስነ-ጽሁፍ መጻሕፍት", "Kitaaba Asoosama Afaan Amaaraa", "ናይ Timeless Amharic Literary Masterpieces & Romantic Fiction"),
					makeLeaf("oromo-literature-history-gadaa-studies", "Gadaa Democratic Philosophy & Afaan Oromo Literary Texts", "የገዳ ስርዓት ፍልስፍና እና የአፋን ኦሮሞ ስነ-ጽሁፍ መጻሕፍት", "Kitaaba Sirna Gadaa fi Seenaa", "ናይ Gadaa Democratic Philosophy & Afaan Oromo Literary Texts"),
					makeLeaf("tigrinya-literature-northern-history", "Tigrinya Novels, Poetry, Cultural Proverbs and Northern History", "የትግርኛ ስነ-ጽሁፍ፣ ግጥሞች እና የታሪክ መጻሕፍት", "መጻሕፍቲ ስነ-ጽሑፍ ትግርኛ", "ናይ Tigrinya Novels, Poetry, Cultural Proverbs and Northern History"),
					makeLeaf("ethiopian-philosophy-zara-yaqob-texts", "Hatata Zera Yacob & Medieval Rational Philosophical Treatises", "የዘርአ ያዕቆብ ፍልስፍና (ሐተታ) እና ጥንታዊ የፍልስፍና መጻሕፍት", "Kitaaba Falaasama Itoophiyaa", "ናይ Hatata Zera Yacob & Medieval Rational Philosophical Treatises"),
					makeLeaf("ethiopian-biographies-memoirs", "Autobiographies & Memoirs of Notable Ethiopian Leaders and Scholars", "የታዋቂ ኢትዮጵያውያን መሪዎች እና ምሁራን የህይወት ታሪኮች", "Kitaaba Seenaa Jireenyaa", "ናይ Autobiographies & Memoirs of Notable Ethiopian Leaders and Scholars"),
					makeLeaf("traditional-ethiopian-poetry-qene", "Ge'ez & Amharic Sem-enna-Worq Wax and Gold Poetry Collections", "የሰምና ወርቅ ቅኔያት እና የግጥም መድብሎች", "Kitaaba Walaloo fi Qene", "ናይ Ge'ez & Amharic Sem-enna-Worq Wax and Gold Poetry Collections"),
					makeLeaf("contemporary-political-economic-commentary", "Modern Ethiopian Geopolitics, Federalism and Development Studies", "የወቅታዊ የኢትዮጵያ ፖለቲካ እና ኢኮኖሚ ጥናት መጻሕፍት", "Kitaaba Siyaasaa fi Diinagdee", "ናይ Modern Ethiopian Geopolitics, Federalism and Development Studies"),
				},
			},
			{
				Slug:   "religious-theological-texts",
				NameEN: "Religious, Spiritual and Theological Texts",
				NameAM: "የሃይማኖት፣ የቅዱሳት መጻሕፍት እና የመንፈሳዊ ትምህርት መጻሕፍት",
				NameOM: "Kitaabolee Amantaa fi Qulqullaa'aa",
				NameTI: "ሃይማኖታዊን መንፈሳውን መጻሕፍቲ",
				Items: []LeafItemDef{
					makeLeaf("ethiopian-orthodox-bibles-metsehaf-kedus", "81-Canon Complete Geez-Amharic Holy Bibles (Metsehaf Kedus)", "የ81 አሃዱ የኦርቶዶክስ ተዋህዶ ቅዱስ መጽሐፍ", "Macaafa Qulqulluu Ortodoksii (81)", "ናይ 81-Canon Complete Geez-Amharic Holy Bibles (Metsehaf Kedus)"),
					makeLeaf("protestant-christian-bibles-amharic-oromoo", "Old and New Testament Protestant Holy Bibles in Ethiopian Languages", "የወንጌላውያን ቅዱስ መጽሐፍ (በአማርኛ እና አፋን ኦሮሞ)", "Macaafa Qulqulluu Wangeelaa", "ናይ Old and New Testament Protestant Holy Bibles in Ethiopian Languages"),
					makeLeaf("holy-quran-arabic-amharic-translations", "Tajweed Arabic Text with Clear Side-by-Side Amharic/Oromo Meanings", "ቅዱስ ቁርኣን (ከነአማርኛ እና ኦሮምኛ ትርጉሙ)", "Qur'aana Qulqulluu Afaan Oromootiin", "ናይ Tajweed Arabic Text with Clear Side-by-Side Amharic/Oromo Meanings"),
					makeLeaf("hadith-collections-islamic-jurisprudence", "Sahih Bukhari, Sahih Muslim & Fiqh Islamic Reference Books", "የሀዲስ ስብስቦች እና የእስልምና ህግጋት መጻሕፍት", "Kitaaba Hadiisaa fi Fiwhii", "ናይ Sahih Bukhari, Sahih Muslim & Fiqh Islamic Reference Books"),
					makeLeaf("traditional-geez-prayer-books-widdase-maryam", "Leather Bound Geez-Amharic Liturgical & Daily Prayer Books", "የውዳሴ ማርያም እና የጸሎት መጻሕፍት በቆዳ የተዘጋጁ", "Kitaaba Kadhannaa Ge'ez fi Amaaraa", "ናይ Leather Bound Geez-Amharic Liturgical & Daily Prayer Books"),
					makeLeaf("spiritual-devotionals-fasting-guides", "Christian Fasting, Meditation and Daily Devotional Booklets", "የዕለት ተዕለት መንፈሳዊ ማሰላሰያ እና የጾም መመሪያዎች", "Kitaaba Guyyaa Guyyaa Hafuuraa", "ናይ Christian Fasting, Meditation and Daily Devotional Booklets"),
					makeLeaf("islamic-duaa-hisnul-muslim-pocket-books", "Fortress of the Muslim Daily Invocations & Supplications", "የዕለት ተዕለት የዱዓ እና የዚክር አነስተኛ መጻሕፍት", "Kitaaba Du'aa'ii fi Zikrii", "ናይ Fortress of the Muslim Daily Invocations & Supplications"),
					makeLeaf("theological-scholarly-commentaries", "Biblical Exegesis, Church History and Comparative Religion Studies", "የሃይማኖት ትምህርት እና የንጽጽር ሃይማኖት ጥናት መጻሕፍት", "Kitaaba Qorannoo Amantaa", "ናይ Biblical Exegesis, Church History and Comparative Religion Studies"),
				},
			},
			{
				Slug:   "academic-textbooks-exam-prep",
				NameEN: "Academic Textbooks, Children's and Self-Help Books",
				NameAM: "የትምህርት፣ የፈተና መዘጋጃ፣ የልጆች እና የልማት መጻሕፍት",
				NameOM: "Kitaabolee Barnootaa fi Daa'immanii",
				NameTI: "ናይ ትምህርቲ መጻሕፍትን ናይ ቆልዑን",
				Items: []LeafItemDef{
					makeLeaf("ethiopian-ministry-curriculum-textbooks", "Grades 9-12 New Ethiopian Curriculum Subject Reference Textbooks", "የ9ኛ እስከ 12ኛ ክፍል አዲሱ ስርዓተ ትምህርት መማሪያ መጻሕፍት", "Kitaaba Barnootaa Kutaa 9-12", "ናይ Grades 9-12 New Ethiopian Curriculum Subject Reference Textbooks"),
					makeLeaf("grade-12-national-exam-preparation-guides", "Past Matriculation Questions & Worked Solutions Exam Booklets", "የ12ኛ ክፍል ሀገር አቀፍ የዩኒቨርሲቲ መግቢያ ፈተና መዘጋጃ መጻሕፍት", "Kitaaba Qophii Qormaata Kutaa 12", "ናይ Past Matriculation Questions & Worked Solutions Exam Booklets"),
					makeLeaf("university-engineering-medical-textbooks", "International Standard Calculus, Anatomy and Coding Textbooks", "የዩኒቨርሲቲ የህክምና፣ ኢንጂነሪንግ እና የቴክኖሎጂ መጻሕፍት", "Kitaaba Yuunivarsitii Injinariingii", "ናይ International Standard Calculus, Anatomy and Coding Textbooks"),
					makeLeaf("bilingual-children-illustrated-storybooks", "Amharic-English Colorful Illustrated Moral Fairy Tales and Fables", "የልጆች ባለቀለም የስዕል እና የተረት መጻሕፍት", "Kitaaba Sheekkoo Daa'immanii Fakkii", "ናይ Amharic-English Colorful Illustrated Moral Fairy Tales and Fables"),
					makeLeaf("personal-finance-business-self-help", "Wealth Creation, Entrepreneurship and Mindset Growth Best Sellers", "የግል ህይወት ለውጥ፣ የንግድ እና የገንዘብ አስተዳደር መጻሕፍት", "Kitaaba Karoora Jireenyaa fi Daldalaa", "ናይ Wealth Creation, Entrepreneurship and Mindset Growth Best Sellers"),
					makeLeaf("english-language-grammar-ielts-prep", "Advanced Grammar in Use, Vocabulary & IELTS Academic Guides", "የእንግሊዝኛ ቋንቋ ሰዋሰው እና የIELTS ፈተና መዘጋጃ መጻሕፍት", "Kitaaba Afaan Ingilizii fi IELTS", "ናይ Advanced Grammar in Use, Vocabulary & IELTS Academic Guides"),
					makeLeaf("amharic-english-oromoo-dictionaries", "Comprehensive Trilingual Dictionaries and Phrasebooks", "የአማርኛ፣ ኦሮምኛ እና እንግሊዝኛ መዝገበ ቃላት (Dictionary)", "Galmee Jechootaa Afaan Sadii", "ናይ Comprehensive Trilingual Dictionaries and Phrasebooks"),
					makeLeaf("high-capacity-executive-planners-notebooks", "Hardcover Undated Daily Goal Planner Leather Organizers", "የስራ እቅድ ማስታወሻ የቆዳ አጀንዳዎች (Notebook/Planner)", "Daftara Qananii Karoora Hojii", "ናይ Hardcover Undated Daily Goal Planner Leather Organizers"),
				},
			},
		},
	})

	// Domain: education-training
	list = append(list, DomainCatalogDef{
		RootSlug: "education-training",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "technology-software-certifications",
				NameEN: "Coding Bootcamps, IT and Cyber Security Training",
				NameAM: "የኮዲንግ፣ የሶፍትዌር እና የሳይበር ደህንነት ስልጠናዎች",
				NameOM: "Leenjii Sooftiweerii fi Teekinoolojii",
				NameTI: "ናይ ኮምፒዩተር ስልጠናታትን ሰርቲፊኬትን",
				Items: []LeafItemDef{
					makeLeaf("full-stack-web-development-bootcamps", "Hands-On React, Node.js and PostgreSQL Full-Stack Web Bootcamps", "የሙሉ የዌብሳይት ስራ እና ሶፍትዌር ማበልጸጊያ ስልጠናዎች", "Leenjii Full-Stack Web Development", "ናይ Hands-On React, Node.js and PostgreSQL Full-Stack Web Bootcamps"),
					makeLeaf("python-data-science-machine-learning-courses", "Data Analysis, Pandas, SQL and Machine Learning Certification Training", "የዳታ ሳይንስ እና ፓይተን ፕሮግራሚንግ ስልጠናዎች", "Leenjii Data Science fi Python", "ናይ Data Analysis, Pandas, SQL and Machine Learning Certification Training"),
					makeLeaf("comptia-cisco-ccna-networking-classes", "Live Lab Cisco CCNA & Network Troubleshooting Certification Preparation", "የሲስኮ CCNA እና የኔትወርክ ሰርቲፊኬሽን ስልጠናዎች", "Leenjii CCNA fi Netwoorkii", "ናይ Live Lab Cisco CCNA & Network Troubleshooting Certification Preparation"),
					makeLeaf("ethical-hacking-cyber-security-defense", "Offensive Security Penetration Testing & Defense Certification Courses", "የሳይበር ደህንነት እና የኢቲካል ሄኪንግ ስልጠናዎች", "Leenjii Nageenya Saayibarii (Cyber Security)", "ናይ Offensive Security Penetration Testing & Defense Certification Courses"),
					makeLeaf("mobile-app-development-flutter-react-native", "Cross-Platform iOS & Android Mobile Application Development Classes", "የሞባይል አፕሊኬሽን ስራ ስልጠናዎች (Flutter/React Native)", "Leenjii Hojii Appii Moobaayilii", "ናይ Cross-Platform iOS & Android Mobile Application Development Classes"),
					makeLeaf("cloud-computing-aws-azure-certification", "Amazon AWS & Microsoft Azure Solutions Architect Exam Prep Classes", "የክላውድ ኮምፒውቲንግ ስልጠናዎች (AWS/Azure)", "Leenjii Kilaawudii (AWS/Azure)", "ናይ Amazon AWS & Microsoft Azure Solutions Architect Exam Prep Classes"),
					makeLeaf("ui-ux-product-design-figma-mastery", "User Research, Wireframing & Interactive Figma Prototyping Courses", "የUI/UX ዲዛይን እና ፊግማ ስልጠናዎች", "Leenjii Dizaayinii UI/UX fi Figma", "ናይ User Research, Wireframing & Interactive Figma Prototyping Courses"),
					makeLeaf("ai-prompt-engineering-generative-ai", "Enterprise LLM Integration & Generative AI Workflow Training", "የሰው ሰራሽ አስተውሎት (AI) እና የፕሮምፕት ስልጠናዎች", "Leenjii Ogummaa AI Ammayyaa", "ናይ Enterprise LLM Integration & Generative AI Workflow Training"),
				},
			},
			{
				Slug:   "vocational-tvet-technical-trades",
				NameEN: "TVET Vocational Skills, Electric, Auto and Plumbing",
				NameAM: "የቴክኒክና ሙያ፣ የመኪና ጥገና እና የኤሌክትሪክ ስልጠናዎች",
				NameOM: "Leenjii Ogummaa fi Teeknikaa (TVET)",
				NameTI: "ናይ ቴክኒክን ሞያን ስልጠናታት",
				Items: []LeafItemDef{
					makeLeaf("auto-mechanics-electronic-diagnostics-course", "Vehicle Engine Rebuilding & OBD-II Scanner Diagnostic Masterclasses", "የመኪና ሞተር እና የኤሌክትሮኒክስ ጥገና የሙያ ስልጠና", "Leenjii Makaanikii Konkolaataa", "ናይ Vehicle Engine Rebuilding & OBD-II Scanner Diagnostic Masterclasses"),
					makeLeaf("certified-building-electrical-installation", "Domestic & Industrial 3-Phase Electrical Wiring Certification Courses", "የህንፃ እና የኢንዱስትሪ ኤሌክትሪክ ዝርጋታ የሙያ ስልጠና", "Leenjii Elektiriikii Gamoo fi Warshaa", "ናይ Domestic & Industrial 3-Phase Electrical Wiring Certification Courses"),
					makeLeaf("sanitary-plumbing-pipefitting-training", "PPR & HDPE Water Supply, Drainage & Pump Maintenance Trade Classes", "የቧንቧ እና የፍሳሽ መስመር ዝርጋታ የቴክኒክ ስልጠና", "Leenjii Ujummoo Bishaanii (Plumbing)", "ናይ PPR & HDPE Water Supply, Drainage & Pump Maintenance Trade Classes"),
					makeLeaf("mig-tig-electric-arc-welding-courses", "Structural Steel Arc, MIG & TIG Metal Fabrication Welding Training", "የብረታ ብረት እና የአርክ ብየዳ ስልጠናዎች", "Leenjii Biriqaa Sibilaa (Welding)", "ናይ Structural Steel Arc, MIG & TIG Metal Fabrication Welding Training"),
					makeLeaf("refrigeration-air-conditioning-technician", "HVAC Inverter AC & Commercial Refrigerator Gas Charging Courses", "የፍሪጅ እና ኤሲ ጥገና የሙያ ስልጠና", "Leenjii Suphaa Firiijii fi AC", "ናይ HVAC Inverter AC & Commercial Refrigerator Gas Charging Courses"),
					makeLeaf("heavy-machinery-excavator-loader-operating", "Safe Hydraulic Excavator, Forklift & Wheel Loader Operator Training", "የትላልቅ ኤክስካቫተር እና ሎደር ማሽኖች የማሽከርከር ስልጠና", "Leenjii Konkolaachisaa Maashinii Gurguddaa", "ናይ Safe Hydraulic Excavator, Forklift & Wheel Loader Operator Training"),
					makeLeaf("professional-tailoring-fashion-design-classes", "Pattern Making, Garment Cutting and Industrial Sewing Machine Classes", "የስፌት፣ የፋሽን ዲዛይን እና የልብስ ጥልፍ ስልጠናዎች", "Leenjii Hodhaa fi Dizaayinii Uffataa", "ናይ Pattern Making, Garment Cutting and Industrial Sewing Machine Classes"),
					makeLeaf("woodworking-cabinetry-furniture-fabrication", "Precision Joinery, CNC Wood Carving and Modern Cabinetry Courses", "የአናጺነት እና ዘመናዊ የእንጨት እቃዎች ስራ ስልጠና", "Leenjii Hojii Mukaa fi Soofaa", "ናይ Precision Joinery, CNC Wood Carving and Modern Cabinetry Courses"),
				},
			},
			{
				Slug:   "languages-business-music-tutoring",
				NameEN: "Languages, Business Management and Tutoring",
				NameAM: "የቋንቋ፣ የቢዝነስ ማኔጅመንት እና የማጠናከሪያ ትምህርቶች",
				NameOM: "Leenjii Afaanii, Daldalaa fi Barnootaa",
				NameTI: "ናይ ቋንቋን ቢዝነስን ስልጠናታት",
				Items: []LeafItemDef{
					makeLeaf("spoken-english-fluency-ielts-training", "Accent Reduction, Public Speaking & Academic IELTS Prep Classes", "የእንግሊዝኛ አነጋገር እና የIELTS ፈተና ማሰልጠኛ", "Leenjii Dubbii Afaan Ingilizii fi IELTS", "ናይ Accent Reduction, Public Speaking & Academic IELTS Prep Classes"),
					makeLeaf("conversational-amharic-for-expats", "Practical Daily Vocabulary & Grammar for Foreign Residents in Ethiopia", "ለዲፕሎማቶች እና የውጭ ዜጎች የሚሆን የአማርኛ ቋንቋ ስልጠና", "Leenjii Afaan Amaaraa Alagootaaf", "ናይ Practical Daily Vocabulary & Grammar for Foreign Residents in Ethiopia"),
					makeLeaf("practical-afaan-oromoo-language-classes", "Interactive Conversational Oromo Grammar & Business Terminology", "የአፋን ኦሮሞ ቋንቋ አነጋገር እና የጽህፈት ስልጠና", "Leenjii Afaan Oromoo Hundaaf", "ናይ Interactive Conversational Oromo Grammar & Business Terminology"),
					makeLeaf("chinese-mandarin-hsk-classes", "Business Mandarin Speaking & HSK International Exam Preparation", "የቻይንኛ ቋንቋ የቢዝነስ አነጋገር ስልጠና", "Leenjii Afaan Chaayinaa (Mandarin)", "ናይ Business Mandarin Speaking & HSK International Exam Preparation"),
					makeLeaf("accounting-quickbooks-peachtree-mastery", "Practical Ethiopian Tax, Payroll & QuickBooks Bookkeeping Training", "የፒችትሪ እና የኩዊክቡክስ የሂሳብ አያያዝ ስልጠና", "Leenjii Herregaa fi QuickBooks", "ናይ Practical Ethiopian Tax, Payroll & QuickBooks Bookkeeping Training"),
					makeLeaf("digital-marketing-social-media-management", "Meta Ads, Google SEO, Content Creation and Funnel Scaling Classes", "የዲጂታል ማርኬቲንግ እና የሶሻል ሚዲያ ማስታወቂያ ስልጠና", "Leenjii Daldala Dijitaalaa (Digital Marketing)", "ናይ Meta Ads, Google SEO, Content Creation and Funnel Scaling Classes"),
					makeLeaf("private-home-tutors-math-physics", "One-on-One Elementary & High School Exam Preparation Home Tutoring", "የቤት ለቤት የግል አስተማሪዎች (ሒሳብ፣ ፊዚክስ፣ እንግሊዝኛ)", "Barsiisaa Manaa (Home Tutor)", "ናይ One-on-One Elementary & High School Exam Preparation Home Tutoring"),
					makeLeaf("acoustic-guitar-piano-traditional-krar-lessons", "Hands-On Instrumental Music Theory, Keyboard & Krar Lessons", "የፒያኖ፣ የጊታር እና የባህል ክራር የሙዚቃ ትምህርቶች", "Leenjii Muuziqaa Gitaaraa fi Kiraaraa", "ናይ Hands-On Instrumental Music Theory, Keyboard & Krar Lessons"),
				},
			},
		},
	})

	// Domain: jobs-recruitment
	list = append(list, DomainCatalogDef{
		RootSlug: "jobs-recruitment",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "it-software-engineering-jobs",
				NameEN: "Software Engineering, Cloud and IT Vacancies",
				NameAM: "የሶፍትዌር ኢንጂነሪንግ እና የቴክኖሎጂ የስራ ማስታወቂያዎች",
				NameOM: "Banaa Hojii Sooftiweerii fi IT",
				NameTI: "ክፉት ስራሕቲ ሶፍትዌርን ITን",
				Items: []LeafItemDef{
					makeLeaf("senior-backend-go-java-engineers", "Senior Go, Java or Node.js Distributed Microservices Backend Developers", "ከፍተኛ የጀርባ ሶፍትዌር አበልጻጊዎች (Backend Developers)", "Ogeessa Sooftiweerii Backend", "ናይ Senior Go, Java or Node.js Distributed Microservices Backend Developers"),
					makeLeaf("frontend-react-nextjs-developers", "TypeScript, React, Next.js and Responsive Web Frontend Engineers", "የፊት ለፊት ዌብ አበልጻጊዎች (Frontend React/Next.js)", "Ogeessa Frontend React/Next.js", "ናይ TypeScript, React, Next.js and Responsive Web Frontend Engineers"),
					makeLeaf("mobile-application-engineers", "Native Android, iOS or Flutter Cross-Platform Mobile Engineers", "የሞባይል አፕሊኬሽን ኢንጂነሮች (Flutter/Android)", "Injinara Appii Moobaayilii", "ናይ Native Android, iOS or Flutter Cross-Platform Mobile Engineers"),
					makeLeaf("cloud-devops-site-reliability-engineers", "Kubernetes, Docker, CI/CD Pipelines & Cloud Infrastructure Engineers", "የክላውድ እና የዴቭኦፕስ መሃንዲሶች (DevOps Engineers)", "Injinara Kilaawudii fi DevOps", "ናይ Kubernetes, Docker, CI/CD Pipelines & Cloud Infrastructure Engineers"),
					makeLeaf("cybersecurity-threat-analysts", "SOC Security Analysts, Vulnerability Testers & ISO Compliance Leads", "የሳይበር ደህንነት እና መረጃ ጥበቃ ባለሙያዎች", "Ogeessa Nageenya Saayibarii", "ናይ SOC Security Analysts, Vulnerability Testers & ISO Compliance Leads"),
					makeLeaf("database-administrators-dba", "PostgreSQL, Oracle and MySQL High Availability Database Architects", "የዳታቤዝ አስተዳዳሪዎች (Database Administrators)", "Bulchaa Kuusaa Daataa (DBA)", "ናይ PostgreSQL, Oracle and MySQL High Availability Database Architects"),
					makeLeaf("it-support-system-administrators", "Active Directory, LAN/WAN Hardware & Desktop Support Specialists", "የአይቲ እና የኔትወርክ ሰፖርት ሰራተኞች", "Ogeessa Gargaarsa IT fi Netwoorkii", "ናይ Active Directory, LAN/WAN Hardware & Desktop Support Specialists"),
					makeLeaf("product-managers-agile-scrum-leads", "Tech Product Managers, Technical Writers and Scrum Master Leads", "የሶፍትዌር ፕሮዳክት ማኔጀሮች እና ስክረም መሪዎች", "Gaggeessaa Pirodaaktii IT", "ናይ Tech Product Managers, Technical Writers and Scrum Master Leads"),
				},
			},
			{
				Slug:   "engineering-construction-architecture-jobs",
				NameEN: "Civil Engineering, Architecture and Construction Vacancies",
				NameAM: "የሲቪል ኢንጂነሪንግ፣ አርክቴክቸር እና ኮንስትራክሽን ስራዎች",
				NameOM: "Banaa Hojii Injinariingii fi Gamoo",
				NameTI: "ክፉት ስራሕቲ ምህንድስናን ህንጸትን",
				Items: []LeafItemDef{
					makeLeaf("senior-civil-resident-site-engineers", "BSc Civil Engineering 5+ Years High-Rise Concrete Site Supervision", "ከፍተኛ የሳይት ሲቪል መሃንዲሶች (Site Engineers)", "Injinara Siiviilii Saayitii", "ናይ BSc Civil Engineering 5+ Years High-Rise Concrete Site Supervision"),
					makeLeaf("licensed-architectural-designers", "B.Arch Revit, ArchiCAD & Lumion Residential Architectural Designers", "ባለፈቃድ አርክቴክቶች እና የህንፃ ፕላን ዲዛይነሮች", "Dizaayinara Gamoo (Architect)", "ናይ B.Arch Revit, ArchiCAD & Lumion Residential Architectural Designers"),
					makeLeaf("structural-design-engineers", "MSc Structural Engineering ETABS Seismic Analysis Design Specialists", "የመዋቅር (ስትራክቸራል) ኢንጂነሮች", "Injinara Caasaa Gamoo (Structural)", "ናይ MSc Structural Engineering ETABS Seismic Analysis Design Specialists"),
					makeLeaf("quantity-surveyors-cost-estimators", "Take-Off Sheet, Payment Certificate & BoQ Preparation Surveyors", "የስራ ዝርዝር እና ዋጋ ስሌት ባለሙያዎች (Quantity Surveyors)", "Ogeessa Safara fi Baasii Gamoo", "ናይ Take-Off Sheet, Payment Certificate & BoQ Preparation Surveyors"),
					makeLeaf("sanitary-mep-building-engineers", "Plumbing, Firefighting & Building Mechanical Systems Specialists", "የህንፃ ሳኒተሪ እና ሜካኒካል ኢንጂነሮች", "Injinara Saanitarii fi Ujummoo", "ናይ Plumbing, Firefighting & Building Mechanical Systems Specialists"),
					makeLeaf("electrical-site-supervisors", "Building Power Distribution, Substation & Generator Installation Leads", "የህንፃ ኤሌክትሪክ ሳይት ተቆጣጣሪዎች", "To'ataa Elektiriikii Saayitii", "ናይ Building Power Distribution, Substation & Generator Installation Leads"),
					makeLeaf("geotechnical-materials-engineers", "Soil Boring, Concrete Crushing & Compaction Quality Control Techs", "የአፈር እና የኮንስትራክሽን ማቴሪያል ተመራማሪዎች", "Ogeessa Qorannoo Biyyoo fi Gamoo", "ናይ Soil Boring, Concrete Crushing & Compaction Quality Control Techs"),
					makeLeaf("construction-project-managers", "PMP Certified Real Estate & High-Rise General Project Directors", "የኮንስትራክሽን እና የሪል ስቴት ፕሮጀክት ማኔጀሮች", "Gaggeessaa Pirojektii Gamoo", "ናይ PMP Certified Real Estate & High-Rise General Project Directors"),
				},
			},
			{
				Slug:   "finance-health-education-sales-jobs",
				NameEN: "Finance, Healthcare, Sales and Education Vacancies",
				NameAM: "የፋይናንስ፣ የጤና፣ የሽያጭ እና የመምህርነት ስራዎች",
				NameOM: "Banaa Hojii Faayinaansii, Fayyaa fi Barnootaa",
				NameTI: "ክፉት ስራሕቲ ፋይናንስን ጥዕናን ትምህርትን",
				Items: []LeafItemDef{
					makeLeaf("senior-corporate-accountants", "BA Accounting & Finance IFRS Reporting & Ministry Tax Compliance", "ከፍተኛ የሂሳብ ሹሞች እና አካውንታንቶች", "Ogeessa Herregaa Olaanaa", "ናይ BA Accounting & Finance IFRS Reporting & Ministry Tax Compliance"),
					makeLeaf("certified-internal-auditors", "ACCA / CIA Internal Control & Financial Risk Assurance Auditors", "የውስጥ ኦዲተሮች (Internal Auditors)", "Ooditara Keessaa", "ናይ ACCA / CIA Internal Control & Financial Risk Assurance Auditors"),
					makeLeaf("general-medical-practitioners", "MD Certified General Physicians for Private Clinics and Hospitals", "አጠቃላይ የህክምና ዶክተሮች (Medical Doctors)", "Ogeessa Yaalaa (Doktara)", "ናይ MD Certified General Physicians for Private Clinics and Hospitals"),
					makeLeaf("registered-clinical-nurses", "BSc Nursing Licensed Emergency & Inpatient Care Clinical Nurses", "የተመዘገቡ ባለሙያ ነርሶች (BSc Nurses)", "Ogeettii Narsii Hospitaalaa", "ናይ BSc Nursing Licensed Emergency & Inpatient Care Clinical Nurses"),
					makeLeaf("licensed-retail-pharmacists", "B.Pharm Licensed Community Drug Dispensing Pharmacists", "ባለፈቃድ የፋርማሲ ባለሙያዎች (Pharmacists)", "Ogeessa Qorichaa (Faarmaasistii)", "ናይ B.Pharm Licensed Community Drug Dispensing Pharmacists"),
					makeLeaf("corporate-b2b-sales-executives", "Commission-Driven High-Value Software & Equipment Sales Agents", "የድርጅት ሽያጭ እና የገበያ ማስፋፊያ ባለሙያዎች", "Gurguraa Daldala Dhaabbilee (B2B)", "ናይ Commission-Driven High-Value Software & Equipment Sales Agents"),
					makeLeaf("secondary-school-stem-teachers", "Bed Math, Physics, Chemistry & Biology Secondary School Teachers", "የሁለተኛ ደረጃ የሳይንስ እና ሒሳብ መምህራን", "Barsiisaa Saayinsii Kutaa 9-12", "ናይ Bed Math, Physics, Chemistry & Biology Secondary School Teachers"),
					makeLeaf("customer-service-call-center-agents", "Fluent Amharic, Oromo, Tigrinya & English Front Desk Representatives", "የደንበኞች አገልግሎት እና የጥሪ ማዕከል ባለሙያዎች", "Tajaajila Maamiltootaa fi Waamichaa", "ናይ Fluent Amharic, Oromo, Tigrinya & English Front Desk Representatives"),
				},
			},
		},
	})

	// Domain: cvs-talent
	list = append(list, DomainCatalogDef{
		RootSlug: "cvs-talent",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "tech-engineering-talent-profiles",
				NameEN: "Software, Tech and Engineering Talent Profiles",
				NameAM: "የሶፍትዌር፣ የቴክኖሎጂ እና የኢንጂነሪንግ ባለሙያ መገለጫዎች",
				NameOM: "Piroofayilii Ogeessota IT fi Injinariingii",
				NameTI: "ናይ ITን ምህንድስናን ክኢላታት መገለጺ",
				Items: []LeafItemDef{
					makeLeaf("verified-full-stack-engineer-profiles", "Verified Senior Full-Stack Engineers Open to Full-Time Hiring", "የተረጋገጡ የሙሉ ሶፍትዌር ኢንጂነሮች ፕሮፋይል", "Piroofayilii Ogeessa Full-Stack", "ናይ Verified Senior Full-Stack Engineers Open to Full-Time Hiring"),
					makeLeaf("licensed-civil-structural-engineer-profiles", "Registered Professional Civil Engineers with Verified Portfolios", "የተረጋገጡ የሲቪል እና ስትራክቸራል ኢንጂነሮች ፕሮፋይል", "Piroofayilii Injinara Siiviilii", "ናይ Registered Professional Civil Engineers with Verified Portfolios"),
					makeLeaf("ui-ux-design-lead-portfolios", "Senior Figma UI/UX Product Designers with Live App Case Studies", "የዩአይ/ዩኤክስ ፕሮዳክት ዲዛይነሮች ፖርትፎሊዮ", "Piroofayilii Dizaayinara UI/UX", "ናይ Senior Figma UI/UX Product Designers with Live App Case Studies"),
					makeLeaf("devops-cloud-architect-cvs", "AWS Certified Solutions Architects with Infrastructure as Code Skills", "የክላውድ እና ዴቭኦፕስ ባለሙያዎች ሲቪ", "Piroofayilii Ogeessa Kilaawudii", "ናይ AWS Certified Solutions Architects with Infrastructure as Code Skills"),
					makeLeaf("mobile-flutter-react-native-specialists", "Published App Store Mobile Engineers with 4+ Years Experience", "የሞባይል አፕሊኬሽን ስራ ባለሙያዎች ፕሮፋይል", "Piroofayilii Ogeessa Appii Moobaayilii", "ናይ Published App Store Mobile Engineers with 4+ Years Experience"),
					makeLeaf("data-analyst-bi-specialist-profiles", "PowerBI, Tableau & Advanced SQL Business Intelligence Specialists", "የዳታ አናሊስት እና ቢአይ ባለሙያዎች መገለጫ", "Piroofayilii Ogeessa Xiinxala Daataa", "ናይ PowerBI, Tableau & Advanced SQL Business Intelligence Specialists"),
					makeLeaf("senior-qa-test-automation-engineers", "Selenium, Cypress & Postman Automated Testing Quality Engineers", "የሶፍትዌር ጥራት እና አውቶሜሽን ፈታሾች ፕሮፋይል", "Piroofayilii Ogeessa Qorannoo Sooftiweerii", "ናይ Selenium, Cypress & Postman Automated Testing Quality Engineers"),
					makeLeaf("licensed-architect-cad-draftsman-cvs", "Certified 3D Architectural Visualizers & BIM Detailers", "የአርክቴክቶች እና የ3ዲ ዲዛይን ባለሙያዎች ሲቪ", "Piroofayilii Dizaayinara Gamoo 3D", "ናይ Certified 3D Architectural Visualizers & BIM Detailers"),
				},
			},
			{
				Slug:   "business-finance-legal-talent-profiles",
				NameEN: "Finance, Accounting, Legal and HR Talent Profiles",
				NameAM: "የፋይናንስ፣ የሂሳብ፣ የህግ እና የሰው ሃብት ባለሙያ መገለጫዎች",
				NameOM: "Piroofayilii Ogeessota Herregaa, Seeraa fi HR",
				NameTI: "ናይ ፋይናንስን ሕግን HRን ክኢላታት መገለጺ",
				Items: []LeafItemDef{
					makeLeaf("acca-certified-financial-managers", "Chartered Certified Accountants with Multinational ERP Experience", "የኤሲሲኤ (ACCA) የተረጋገጡ የፋይናንስ ስራ አስኪያጆች", "Piroofayilii Ogeessa ACCA Herregaa", "ናይ Chartered Certified Accountants with Multinational ERP Experience"),
					makeLeaf("experienced-tax-compliance-accountants", "Ministry of Revenue Certified Tax Filing & Audit Senior Specialists", "የታክስ እና የተጨማሪ እሴት ታክስ ከፍተኛ አካውንታንቶች", "Piroofayilii Ogeessa Gibiraa fi VAT", "ናይ Ministry of Revenue Certified Tax Filing & Audit Senior Specialists"),
					makeLeaf("licensed-corporate-commercial-lawyers", "Commercial Arbitration, Trade Licensing & Investment Attorneys", "የንግድ እና የኮርፖሬት ህግ ጠበቆች እና የህግ አማካሪዎች", "Piroofayilii Abukaatoo fi Ogeessa Seeraa", "ናይ Commercial Arbitration, Trade Licensing & Investment Attorneys"),
					makeLeaf("senior-human-resources-people-directors", "Labor Proclamation Compliance, Talent Acquisition & Payroll Heads", "የሰው ሃብት (HR) ዳይሬክተሮች እና አስተዳዳሪዎች", "Piroofayilii Gaggeessaa Qabeenya Namaa (HR)", "ናይ Labor Proclamation Compliance, Talent Acquisition & Payroll Heads"),
					makeLeaf("certified-procurement-supply-chain-heads", "International Import Letter of Credit (LC) & Logistics Specialists", "የግዥ፣ የውጭ ንግድ እና የኤልሲ (LC) ስፔሻሊስቶች", "Piroofayilii Ogeessa Bittaa fi Lojistiiksii", "ናይ International Import Letter of Credit (LC) & Logistics Specialists"),
					makeLeaf("strategic-marketing-brand-managers", "FMCG Brand Growth, Corporate PR and Campaign Directors", "የብራንድ እና የገበያ ጥናት ስራ አስኪያጆች", "Piroofayilii Gaggeessaa Daldalaa fi Maarkeetiingii", "ናይ FMCG Brand Growth, Corporate PR and Campaign Directors"),
					makeLeaf("bilingual-executive-administrative-assistants", "Fluent English & Amharic C-Suite Office Executive Assistants", "የከፍተኛ ስራ አስፈፃሚዎች ረዳት አስተዳዳሪዎች", "Piroofayilii Gargaaraa Hojii Gaggeessaa", "ናይ Fluent English & Amharic C-Suite Office Executive Assistants"),
					makeLeaf("certified-pmp-project-managers", "PMP Credentialed Infrastructure & Tech Transformation Leads", "የተረጋገጡ የፕሮጀክት ማኔጅመንት ባለሙያዎች (PMP)", "Piroofayilii Gaggeessaa Pirojektii PMP", "ናይ PMP Credentialed Infrastructure & Tech Transformation Leads"),
				},
			},
			{
				Slug:   "health-education-hospitality-talent-profiles",
				NameEN: "Healthcare, Education, Hospitality and Operations Talent",
				NameAM: "የጤና፣ የትምህርት፣ የሆቴል እና የአሰራር ባለሙያ መገለጫዎች",
				NameOM: "Piroofayilii Ogeessota Fayyaa, Barnootaa fi Hoteelaa",
				NameTI: "ናይ ጥዕናን ትምህርትን ሆስፒታሊቲን ክኢላታት መገለጺ",
				Items: []LeafItemDef{
					makeLeaf("licensed-medical-officer-physician-cvs", "5+ Years Post-Internship Hospital General Practitioners", "የህክምና ዶክተሮች የተሟላ ሲቪ እና የስራ ታሪክ", "Piroofayilii Doktara Yaalaa", "ናይ 5+ Years Post-Internship Hospital General Practitioners"),
					makeLeaf("experienced-icu-pediatric-nurses", "Licensed Clinical Hospital Inpatient & Neonatal ICU Nurses", "የጽኑ ህሙማን እና የህፃናት ህክምና ነርሶች", "Piroofayilii Narsii Yaala Cimaa", "ናይ Licensed Clinical Hospital Inpatient & Neonatal ICU Nurses"),
					makeLeaf("secondary-high-school-stem-educators", "Master of Education Subject Matter Specialist High School Teachers", "የትምህርት ማስተርስ ያላቸው የሁለተኛ ደረጃ መምህራን", "Piroofayilii Barsiisaa Saayinsii fi Herregaa", "ናይ Master of Education Subject Matter Specialist High School Teachers"),
					makeLeaf("certified-executive-chefs-culinary-leads", "5-Star Hotel Continental & Traditional Ethiopian Master Chefs", "የአምስት ኮከብ ሆቴል ዋና ሼፎች እና የምግብ አዘጋጆች", "Piroofayilii Sheefii Nyaata Qopheessu", "ናይ 5-Star Hotel Continental & Traditional Ethiopian Master Chefs"),
					makeLeaf("hotel-general-managers-hospitality-leads", "Front Office, F&B and Housekeeping Operations General Directors", "የሆቴል እና ቱሪዝም ስራ አስኪያጆች", "Piroofayilii Gaggeessaa Hoteelaa fi Turiizimii", "ናይ Front Office, F&B and Housekeeping Operations General Directors"),
					makeLeaf("certified-commercial-heavy-truck-drivers", "Grade-5 Licensed Inter-City Trailer & Fuel Tanker Master Drivers", "የደረጃ 5 መንጃ ፈቃድ ያላቸው የከባድ መኪና አሽከርካሪዎች", "Piroofayilii Konkolaachisaa Fe'iinsaa Sadarkaa 5", "ናይ Grade-5 Licensed Inter-City Trailer & Fuel Tanker Master Drivers"),
					makeLeaf("expert-baristas-roasters-coffee-cuppers", "SCA Certified Specialty Coffee Cuppers, Roasters and Baristas", "የተረጋገጡ የቡና አዋቂዎች፣ ቆዪዎች እና ባሪስታዎች", "Piroofayilii Ogeessa Bunaa fi Baaristichaa", "ናይ SCA Certified Specialty Coffee Cuppers, Roasters and Baristas"),
					makeLeaf("tvet-certified-master-electricians", "Level-IV TVET Certified Substation & Industrial Electricians", "የደረጃ 4 ቴክኒክና ሙያ የተረጋገጡ ማስተር ኤሌክትሪሻኖች", "Piroofayilii Ogeessa Elektiriikii Sadarkaa 4", "ናይ Level-IV TVET Certified Substation & Industrial Electricians"),
				},
			},
		},
	})

	// Domain: freelance-gigs
	list = append(list, DomainCatalogDef{
		RootSlug: "freelance-gigs",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "digital-design-creative-gigs",
				NameEN: "Graphic Design, Video Editing and Creative Services",
				NameAM: "የግራፊክስ ዲዛይን፣ የቪዲዮ ኤዲቲንግ እና የፈጠራ ስራዎች",
				NameOM: "Hojiiwwan Dizaayinii, Viidiyoo fi Ogummaa",
				NameTI: "ናይ ግራፊክስ ዲዛይንን ቪድዮ ኤዲቲንግን",
				Items: []LeafItemDef{
					makeLeaf("custom-business-logo-brand-identity-gigs", "Vector Minimalist Logo, Brand Book & Color Palette Gigs", "የድርጅት ሎጎ እና መለያ ምልክት ዲዛይን ስራዎች", "Dizaayinii Loogoo fi Mallattoo Daldalaa", "ናይ Vector Minimalist Logo, Brand Book & Color Palette Gigs"),
					makeLeaf("youtube-tiktok-reels-video-editing-gigs", "CapCut, Premiere Pro 4K Fast Viral Social Video Editing", "የቲክቶክ፣ ዩቲዩብ እና ሪልስ ቪዲዮ ኤዲቲንግ ስራዎች", "Gulaallii Viidiyoo YouTube fi TikTok", "ናይ CapCut, Premiere Pro 4K Fast Viral Social Video Editing"),
					makeLeaf("flyer-poster-billboard-print-design-gigs", "High-Resolution Event Banner, Brochure & Social Media Post Graphics", "የማስታወቂያ ፍላየር፣ ፖስተር እና ባነር ዲዛይን ስራዎች", "Dizaayinii Poostaraa fi Biroosharaa", "ናይ High-Resolution Event Banner, Brochure & Social Media Post Graphics"),
					makeLeaf("3d-product-modeling-packaging-renders", "Blender & 3ds Max Realistic Product Bottle & Box 3D Renders", "የምርቶች 3D ሞዴሊንግ እና ፓኬጅንግ ዲዛይን ስራዎች", "Dizaayinii 3D Saanduqa Oomishaa", "ናይ Blender & 3ds Max Realistic Product Bottle & Box 3D Renders"),
					makeLeaf("custom-hand-drawn-digital-illustrations", "Procreate Children Book Characters & Portrait Digital Art Gigs", "የመጽሐፍ እና የቁም ስዕል ዲጂታል ኢላስትሬሽን ስራዎች", "Fakkii Dijitaalaa fi Kitaaba Daa'immanii", "ናይ Procreate Children Book Characters & Portrait Digital Art Gigs"),
					makeLeaf("ui-mockup-mobile-app-wireframe-gigs", "Figma Mobile Screen Redesigns & Clean Modern Landing Pages", "የስልክ አፕሊኬሽን እና የዌብሳይት ስክሪን ዲዛይን ስራዎች", "Dizaayinii Iskiriinii Appii fi Weebsaayitii", "ናይ Figma Mobile Screen Redesigns & Clean Modern Landing Pages"),
					makeLeaf("voiceover-audio-narration-amharic-oromo", "Pro Studio Clear Male/Female Commercial & Documentary Voiceovers", "የማስታወቂያ እና የዶክመንተሪ የድምጽ ንባብ ስራዎች", "Waraabbii Sagalee Beeksisootaa", "ናይ Pro Studio Clear Male/Female Commercial & Documentary Voiceovers"),
					makeLeaf("podcast-audio-cleaning-noise-removal", "Echo Reduction, Level Balancing & Mastering for Clean Podcast Audio", "የፖድካስት ድምጽ ማስተካከል እና ማጥራት ስራዎች", "Qulqulleessa Sagalee fi Poodkaastii", "ናይ Echo Reduction, Level Balancing & Mastering for Clean Podcast Audio"),
				},
			},
			{
				Slug:   "content-writing-translation-gigs",
				NameEN: "Content Writing, Translation and Virtual Assistance",
				NameAM: "የጽሁፍ ዝግጅት፣ የትርጉም እና የቨርቹዋል ረዳት ስራዎች",
				NameOM: "Barreeffama, Turjumaana fi Gargaarsa Online",
				NameTI: "ናይ ጽሑፍን ትርጉምን ኦንላይን ረድኤትን",
				Items: []LeafItemDef{
					makeLeaf("certified-amharic-english-translation-gigs", "Fast Notarized Embassy Document, Agreement & Certificate Translation", "የእንግሊዝኛ ወደ አማርኛ እና በተቃራኒው የሰነዶች ትርጉም ስራዎች", "Turjumaana Afaan Ingilizii gara Amaaraatti", "ናይ Fast Notarized Embassy Document, Agreement & Certificate Translation"),
					makeLeaf("certified-afaan-oromoo-translation-gigs", "Accurate Legal, Medical & Marketing English to Afaan Oromoo Tasks", "ከእንግሊዝኛ ወደ አፋን ኦሮሞ የሰነዶች ትርጉም ስራዎች", "Turjumaana Afaan Ingilizii gara Oromootti", "ናይ Accurate Legal, Medical & Marketing English to Afaan Oromoo Tasks"),
					makeLeaf("tigrinya-english-translation-services", "High-Accuracy Official Document and Text Tigrinya Translations", "የትግርኛ ወደ እንግሊዝኛ እና በተቃራኒው የትርጉም ስራዎች", "ትርጉም ትግርኛ ናብ እንግሊዝኛን", "ናይ High-Accuracy Official Document and Text Tigrinya Translations"),
					makeLeaf("seo-optimized-blog-article-writing", "Niche Researched 1500-Word Google Ranking Blog Articles", "ለዌብሳይት የሚሆኑ ማራኪ ጽሁፎች እና ብሎጎች ዝግጅት", "Barreeffama Weebsaayitii fi Seenaa", "ናይ Niche Researched 1500-Word Google Ranking Blog Articles"),
					makeLeaf("professional-cv-resume-linkedin-rewrite", "ATS-Friendly Executive Resume Polish, Cover Letter & LinkedIn Revamp", "ዘመናዊ የሲቪ (CV) አጻጻፍ እና የሊንክድኢን ፕሮፋይል ማስተካከል", "Barreeffama CV fi Xalayaa Hojii", "ናይ ATS-Friendly Executive Resume Polish, Cover Letter & LinkedIn Revamp"),
					makeLeaf("data-entry-web-research-lead-generation", "Accurate Excel Spreadsheet Data Entry, Cleaning & Web Scraping", "የዳታ ማስገባት፣ ማደራጀት እና መረጃ የማሰባሰብ ስራዎች", "Galmee Daataa fi Qorannoo Weebsaayitii", "ናይ Accurate Excel Spreadsheet Data Entry, Cleaning & Web Scraping"),
					makeLeaf("ecommerce-product-listing-upload-gigs", "Bulk Shopify, Amazon & Huntiq Product Descriptions & Image Edits", "ምርቶችን ወደ ኦንላይን መደብር የመጫን እና ዝርዝር መጻፍ ስራዎች", "Oomisha Weebsaayitiitti Fe'uu", "ናይ Bulk Shopify, Amazon & Huntiq Product Descriptions & Image Edits"),
					makeLeaf("remote-virtual-executive-assistant-services", "Calendar Scheduling, Email Inbox Management & Travel Booking Support", "የርቀት የቢሮ ረዳትነት እና የቀጠሮ ማስተናገድ ስራዎች", "Tajaajila Gargaaraa Biiroo Online", "ናይ Calendar Scheduling, Email Inbox Management & Travel Booking Support"),
				},
			},
			{
				Slug:   "software-tech-freelance-tasks",
				NameEN: "Web Development, Bug Fixing and Tech Gigs",
				NameAM: "የዌብሳይት ስራ፣ የቴክኒክ ጥገና እና ፈጣን የቴክኖሎጂ ስራዎች",
				NameOM: "Hojiiwwan Weebsaayitii fi Teekinoolojii",
				NameTI: "ናይ ዌብሳይት ስራሕትን ቴክኒካዊ ጽገናን",
				Items: []LeafItemDef{
					makeLeaf("wordpress-elementor-business-website-setup", "Responsive 5-Page Modern Corporate WordPress Website Launch", "የድርጅት ዌብሳይት በዎርድፕረስ በአጭር ቀን ማዘጋጀት", "Ijaarsa Weebsaayitii WordPress", "ናይ Responsive 5-Page Modern Corporate WordPress Website Launch"),
					makeLeaf("shopify-woocommerce-ecommerce-store-design", "Payment Gateway Integrated Online Store Setup & Product Import", "የኦንላይን የገበያ መደብር (ኢኮሜርስ) ዌብሳይት ዝግጅት", "Ijaarsa Suuqii Online (E-commerce)", "ናይ Payment Gateway Integrated Online Store Setup & Product Import"),
					makeLeaf("bug-fixing-javascript-python-scripts", "Fast 24-Hour HTML/CSS/JS Layout & Backend Error Troubleshooting", "የኮድ ስህተቶችን እና የዌብሳይት ችግሮችን ማስተካከል", "Sirreeffama Badii Koodii (Bug Fixing)", "ናይ Fast 24-Hour HTML/CSS/JS Layout & Backend Error Troubleshooting"),
					makeLeaf("telegram-automation-sales-bots", "Interactive Inline Keyboard Ethiopian Payment Integrated Telegram Bots", "የቴሌግራም ቦቶች እና አውቶማቲክ የሽያጭ ሲስተም ስራ", "Hojii Boottii Telegiraamii (Telegram Bot)", "ናይ Interactive Inline Keyboard Ethiopian Payment Integrated Telegram Bots"),
					makeLeaf("api-integration-payment-chapa-telebirr", "Telebirr, Chapa & CBE Birr Checkout Webhook API Integrations", "የቴሌብር፣ ቻፓ እና ሲቢኢ ብር የክፍያ ሲስተሞች ማገናኘት", "Waliin Qabsiisa Kaffaltii Telebirr fi Chapa", "ናይ Telebirr, Chapa & CBE Birr Checkout Webhook API Integrations"),
					makeLeaf("excel-macro-vba-dashboard-automation", "Advanced Financial Modeling, Dynamic Pivot Tables & VBA Formulas", "የኤክሴል ማክሮ እና አውቶማቲክ ዳሽቦርድ ማዘጋጀት", "Hojii Excel fi Foormulaa Humna Guddaa", "ናይ Advanced Financial Modeling, Dynamic Pivot Tables & VBA Formulas"),
					makeLeaf("sql-database-query-optimization-backup", "Slow Query Indexing, Database Migration & Automated Offsite Backups", "የዳታቤዝ ፍጥነት ማስተካከል እና የባክአፕ ስራዎች", "Sirreeffama Saffisa Kuusaa Daataa (SQL)", "ናይ Slow Query Indexing, Database Migration & Automated Offsite Backups"),
					makeLeaf("speed-optimization-google-pagespeed-fix", "Core Web Vitals 90+ Score Caching, Compression & Image Optimization", "የዌብሳይት ፍጥነትን እጅግ ፈጣን የማድረግ ስራዎች", "Saffisiisa Weebsaayitii Google PageSpeed", "ናይ Core Web Vitals 90+ Score Caching, Compression & Image Optimization"),
				},
			},
		},
	})

	return list
}
