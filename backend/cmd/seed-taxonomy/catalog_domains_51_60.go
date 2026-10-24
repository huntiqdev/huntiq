package main

import "fmt"

// GetCatalogDomains51To60 returns comprehensive domain catalogues for roots 51 to 60
func GetCatalogDomains51To60() []DomainCatalogDef {
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

	// Domain: music-instruments
	list = append(list, DomainCatalogDef{
		RootSlug: "music-instruments",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "traditional-ethiopian-instruments",
				NameEN: "Krar, Masinko, Kebero Drums and Washint",
				NameAM: "ክራር፣ ማሲንቆ፣ ከበሮ እና ዋሽንት",
				NameOM: "Kiraara, Maasinqoo, Kabbara fi Waashintii",
				NameTI: "ክራርን ማሲንቆን ከበሮን ዋሽንትን",
				Items: []LeafItemDef{
					makeLeaf("electric-acoustic-6-string-krar", "Electric-Acoustic 6-String Wooden Ethiopian Krar with Built-in Pickup", "የኤሌክትሪክ እና የአኮስቲክ ባለ 6 ገመድ የሀበሻ ክራሮች", "Kiraara Aadaa fi Elektiriikii 6-String", "ናይ Electric-Acoustic 6-String Wooden Ethiopian Krar with Built-in Pickup"),
					makeLeaf("handcrafted-horsehair-masinko", "Handcrafted Solid Wood & Horsehair Traditional Masinko with Bow", "የእንጨት እና የፈረስ ጭራ ባህላዊ ማሲንቆ ከመምቻው ጋር", "Maasinqoo Aadaa Harkaan Hojjetame", "ናይ Handcrafted Solid Wood & Horsehair Traditional Masinko with Bow"),
					makeLeaf("traditional-church-cultural-kebero-drums", "Double-Headed Calfskin Wooden Church & Wedding Kebero Drums", "የቆዳ ባህላዊ እና የቤተክርስቲያን ከበሮዎች", "Kabbara Aadaa fi Waldaa Gogaa", "ናይ Double-Headed Calfskin Wooden Church & Wedding Kebero Drums"),
					makeLeaf("bamboo-shepherd-washint-flutes", "4-Hole and 6-Hole Natural Bamboo Highland Washint Flutes", "የቀርከሃ ባህላዊ የዋሽንት ዜማ መሳሪያዎች", "Waashintii Aadaa Leemmanaa", "ናይ 4-Hole and 6-Hole Natural Bamboo Highland Washint Flutes"),
					makeLeaf("traditional-wooden-begena-harp", "10-String Sacred Wooden King David Harp (Begena / Begenna)", "ባለ 10 አውታር ጥንታዊ በገና", "Baganndaa Aadaa (Begena)", "ናይ 10-String Sacred Wooden King David Harp (Begena / Begenna)"),
					makeLeaf("brass-tsebtsebo-cymbals-sistrum-tsenatsil", "Hand-Hammered Brass Liturgical Tsenatsil (Sistrum) & Cymbals", "የነሐስ ጸናጽል እና ቃጭሎች", "Sinaasila Sibiilaa Aadaa (Tsenatsil)", "ናይ Hand-Hammered Brass Liturgical Tsenatsil (Sistrum) & Cymbals"),
					makeLeaf("embilta-traditional-horn-pipes", "Traditional Bamboo & Horn Royal Embilta Wind Instruments", "ባህላዊ የእምቢልታ የዜማ መሳሪያዎች", "Imbiltaa Aadaa", "ናይ Traditional Bamboo & Horn Royal Embilta Wind Instruments"),
					makeLeaf("traditional-animal-skin-dubuq-hand-drums", "Terracotta & Rawhide Small Cultural Folk Dabuka Drums", "የሸክላ እና የቆዳ ባህላዊ ዳቡቃ ከበሮዎች", "Dabuuqaa Aadaa Suphee", "ናይ Terracotta & Rawhide Small Cultural Folk Dabuka Drums"),
				},
			},
			{
				Slug:   "modern-keyboards-guitars-drums",
				NameEN: "Guitars, Arranger Keyboards and Audio Production",
				NameAM: "ጊታሮች፣ ኪቦርዶች፣ ድራሞች እና የስቱዲዮ እቃዎች",
				NameOM: "Gitaara, Kiiboordii fi Diraama Ammayyaa",
				NameTI: "ጊታርን ኪቦርድን ድራምን ናይ ስቱዲዮ ንዋያትን",
				Items: []LeafItemDef{
					makeLeaf("solid-top-steel-acoustic-guitars", "Spruce Solid Top Cutaway Steel-String Acoustic Guitars with Preamp", "ጥራት ያላቸው የአኮስቲክ ጊታሮች ከነ አምፕሊፋየሩ", "Gitaara Akustikii Qulqulluu", "ናይ Spruce Solid Top Cutaway Steel-String Acoustic Guitars with Preamp"),
					makeLeaf("pro-arranger-workstation-keyboards", "61-Key & 88-Key Oriental / Arranger Workstation Synthesizer Keyboards (Korg/Yamaha)", "የሙዚቃ ማቀናበሪያ ፒያኖ ኪቦርዶች (Korg/Yamaha)", "Kiiboordii Muuziqaa Ammayyaa", "ናይ 61-Key & 88-Key Oriental / Arranger Workstation Synthesizer Keyboards (Korg/Yamaha)"),
					makeLeaf("solid-body-electric-lead-guitars", "Dual Humbucker Solid Alder Body Electric Rock & Jazz Guitars", "የኤሌክትሪክ ሊድ ጊታሮች", "Gitaara Elektiriikii Qulqulluu", "ናይ Dual Humbucker Solid Alder Body Electric Rock & Jazz Guitars"),
					makeLeaf("5-piece-acoustic-stage-drum-kits", "Complete 5-Piece Birch Stage Acoustic Drum Kits with Brass Cymbals", "ባለ 5 ድራም የተሟላ የመድረክ አኮስቲክ ድራም ስብስብ", "Diraamii Muuziqaa Guutuu 5-Piece", "ናይ Complete 5-Piece Birch Stage Acoustic Drum Kits with Brass Cymbals"),
					makeLeaf("electronic-mesh-head-digital-drums", "Quiet Mesh-Head USB/MIDI Electronic Practice Digital Drum Sets", "ድምፅ የማይረብሹ የኤሌክትሮኒክስ ዲጂታል ድራሞች", "Diraamii Dijitaalaa Elektirooniiksii", "ናይ Quiet Mesh-Head USB/MIDI Electronic Practice Digital Drum Sets"),
					makeLeaf("studio-reference-monitor-speakers", "Active Bi-Amplified 5-Inch to 8-Inch Studio Audio Reference Monitors", "የስቱዲዮ የሙዚቃ ማቀናበሪያ ድምጽ ማስተካከያ ሞኒተር ስፒከሮች", "Isppiikara Istuudiyoo Qulqullinaa", "ናይ Active Bi-Amplified 5-Inch to 8-Inch Studio Audio Reference Monitors"),
					makeLeaf("large-diaphragm-condenser-microphones", "Gold-Sputtered Vocal Studio Recording XLR Condenser Mics", "ለድምጽ መቅረጫ የሚሆኑ ጥራት ያላቸው ኮንደንሰር ማይክራፎኖች", "Maayikiroofoonii Waraabbii Istuudiyoo", "ናይ Gold-Sputtered Vocal Studio Recording XLR Condenser Mics"),
					makeLeaf("usb-audio-interface-preamps", "24-Bit 192kHz 2-In/2-Out Ultra-Low Noise USB Audio Interfaces", "ከኮምፒውተር ጋር የሚገናኙ የድምጽ ማቀናበሪያ ካርዶች (Audio Interface)", "Kaardii Sagalee Istuudiyoo USB", "ናይ 24-Bit 192kHz 2-In/2-Out Ultra-Low Noise USB Audio Interfaces"),
				},
			},
		},
	})

	// Domain: events-production
	list = append(list, DomainCatalogDef{
		RootSlug: "events-production",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "wedding-tent-stage-decor",
				NameEN: "Tents, Wedding Decoration and Seating Rental",
				NameAM: "የሰርግ ድንኳኖች፣ የዲኮር እቃዎች እና የወንበር ኪራይ",
				NameOM: "Dunkaana Cidhaa, Faaya fi Kiree Barcumaa",
				NameTI: "ድንኳን መርዓን ናይ ዲኮር ንዋያትን",
				Items: []LeafItemDef{
					makeLeaf("heavy-aluminum-marquee-wedding-tents", "100 to 1000 Guest Clear-Span PVC Waterproof Luxury Marquee Tents", "ከ100 እስከ 1000 ሰው የሚይዙ የሰርግ የቅንጦት ድንኳኖች ኪራይ", "Dunkaana Cidhaa Qananii Bal'aa", "ናይ 100 to 1000 Guest Clear-Span PVC Waterproof Luxury Marquee Tents"),
					makeLeaf("chiavari-wedding-banquet-chairs", "Crystal Resin & Gold Wood Padded Chiavari Banquet Event Chairs", "የሰርግ እና የድግስ ውብ የሺቫሪ ወንበሮች ኪራይ (Chiavari Chairs)", "Barcuma Cidhaa Qananii (Chiavari)", "ናይ Crystal Resin & Gold Wood Padded Chiavari Banquet Event Chairs"),
					makeLeaf("led-wedding-stage-floral-backdrops", "Custom Romantic Arch Floral Stages & Fairy Light Head Table Decor", "የሙሽሮች መድረክ ማድመቂያ አበቦች እና መብራቶች (Stage Decor)", "Faaya Masaraa Cidhaa fi Abaaboo", "ናይ Custom Romantic Arch Floral Stages & Fairy Light Head Table Decor"),
					makeLeaf("commercial-event-trussing-stage-platforms", "Modular Heavy Aluminum Outdoor Concert & Wedding Stage Trussing", "የመድረክ መስሪያ ጠንካራ የብረት ትረሶች እና ፕላትፎርሞች", "Masaraa fi Sibila Sagantaa (Trussing)", "ናይ Modular Heavy Aluminum Outdoor Concert & Wedding Stage Trussing"),
					makeLeaf("round-family-banquet-tables-tablecloths", "10-Person Heavy Wooden Round Tables with Satin Tablecloths", "ባለ 10 ሰው ክብ የሰርግ ጠረጴዛዎች እና ጨርቆች", "Minjaala Geengoo Cidhaa", "ናይ 10-Person Heavy Wooden Round Tables with Satin Tablecloths"),
					makeLeaf("moving-head-beam-stage-event-lights", "DMX512 Controlled Sharpy Moving Head Concert Stage Lights", "የመድረክ ማድመቂያ የሚሽከረከሩ ቢም መብራቶች (Moving Head)", "Ibsaa Masaraa fi Sagantaa (Moving Head)", "ናይ DMX512 Controlled Sharpy Moving Head Concert Stage Lights"),
					makeLeaf("outdoor-smoke-bubble-fog-machines", "Low-Lying Dry Ice Cloud & Bubble Wedding Entrance FX Machines", "የሙሽሮች መግቢያ የጭስ እና የአረፋ ማሽኖች (Fog Machine)", "Maashinii Aaraa fi Hurrii Cidhaa", "ናይ Low-Lying Dry Ice Cloud & Bubble Wedding Entrance FX Machines"),
					makeLeaf("red-carpet-vip-stanchion-velvet-ropes", "Polished Chrome Stanchions with Red Velvet Crowd Control Ropes", "የክብር እንግዶች ቀይ ምንጣፍ እና የወርቅ መለያ ብረቶች", "Afata Diimaa fi Sibila Faayaa (Red Carpet)", "ናይ Polished Chrome Stanchions with Red Velvet Crowd Control Ropes"),
				},
			},
			{
				Slug:   "photography-video-production-services",
				NameEN: "Wedding Photography, Cinematography and Livestreaming",
				NameAM: "የሰርግ ፎቶግራፍ፣ ቪዲዮ ቀረጻ እና ቀጥታ ስርጭት",
				NameOM: "Suuraa Cidhaa, Viidiyoo fi Qophii",
				NameTI: "ናይ መርዓ ፎቶን ቪድዮን ቀጥታ ፈነወን",
				Items: []LeafItemDef{
					makeLeaf("cinematic-4k-wedding-video-filming", "Multi-Camera 4K Drone Aerial, Gimbal & Crane Wedding Video Coverage", "በድሮን እና በዘመናዊ ካሜራዎች የሰርግ ቪዲዮ ቀረጻ (4K)", "Waraabbii Viidiyoo Cidhaa 4K", "ናይ Multi-Camera 4K Drone Aerial, Gimbal & Crane Wedding Video Coverage"),
					makeLeaf("luxury-flush-mount-wedding-photo-albums", "Handmade Crystal Acrylic Cover Panoramic High-Gloss Storybooks", "የሰርግ ፎቶ ማስቀመጫ የክሪስታል እና የቆዳ ትላልቅ አልበሞች", "Albamii Suuraa Cidhaa Qananii", "ናይ Handmade Crystal Acrylic Cover Panoramic High-Gloss Storybooks"),
					makeLeaf("pre-wedding-studio-outdoor-photoshoots", "Scenic Destination Engagement & Studio High-End Portrait Sessions", "ከሰርግ በፊት የሚደረጉ የውጪ እና የስቱዲዮ ውብ ፎቶዎች", "Waraabbii Suuraa Kaadhimannaa", "ናይ Scenic Destination Engagement & Studio High-End Portrait Sessions"),
					makeLeaf("multi-camera-live-youtube-streaming", "Blackmagic ATEM Multi-Cam Facebook/YouTube Live Broadcast Rigs", "የሰርግ እና የስብሰባ ቀጥታ የዩቲዩብ ስርጭት አገልግሎት (Livestream)", "Tamsaasa Qophii Kallattii (Livestream)", "ናይ Blackmagic ATEM Multi-Cam Facebook/YouTube Live Broadcast Rigs"),
					makeLeaf("360-degree-photo-booth-spin-platforms", "Automatic Rotating Video Selfie 360 Photo Booth Party Platforms", "360 ዲግሪ ዙሪያውን የሚቀርጹ የፎቶ ቡዝ መድረኮች", "Masaraa Suuraa Digirii 360", "ናይ Automatic Rotating Video Selfie 360 Photo Booth Party Platforms"),
					makeLeaf("event-dj-sound-system-packages", "High-Wattage Line Array Speakers, Subwoofers & Pro Wedding DJs", "የድግስ ድምጽ ማጉያዎች እና ፕሮፌሽናል ዲጄዎች (DJ & Sound)", "Sagalee fi DJ Cidhaa Humna Guddaa", "ናይ High-Wattage Line Array Speakers, Subwoofers & Pro Wedding DJs"),
					makeLeaf("led-video-wall-display-screens-indoor", "P2.5 & P3.9 High-Resolution Seamless Indoor Wedding LED Walls", "በሰርግ እና በስብሰባ ላይ ቪዲዮ የሚያሳዩ ትላልቅ የኤልኢዲ ግድግዳዎች", "Iskiriinii LED Guddaa Cidhaa", "ናይ P2.5 & P3.9 High-Resolution Seamless Indoor Wedding LED Walls"),
					makeLeaf("event-instant-photo-print-kiosks", "Touchscreen Dye-Sublimation 10-Second Instant Magnet Photo Prints", "ፎቶ እንዳነሱ ወዲያውኑ አትሞ የሚያወጣ የፎቶ ማሽን", "Maxxansaa Suuraa Battalaa Sagantaa", "ናይ Touchscreen Dye-Sublimation 10-Second Instant Magnet Photo Prints"),
				},
			},
		},
	})

	// Domain: toys-baby-products
	list = append(list, DomainCatalogDef{
		RootSlug: "toys-baby-products",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "baby-gear-strollers-cribs",
				NameEN: "Baby Strollers, Cribs, Car Seats and Walkers",
				NameAM: "የህፃናት ጋሪዎች፣ አልጋዎች፣ የመኪና ወንበሮች እና መራመጃዎች",
				NameOM: "Gaarii Daa'immanii, Siree fi Barcuma",
				NameTI: "ናይ ህጻናት ጋሪን ዓራትን መንበርን",
				Items: []LeafItemDef{
					makeLeaf("multi-position-foldable-baby-strollers", "One-Hand Gravity Fold Shock-Absorbing Reclining Infant Strollers", "በቀላሉ የሚታጠፉ ምቹ የህፃናት መግፊያ ጋሪዎች", "Gaarii Daa'immanii Dachaafamu", "ናይ One-Hand Gravity Fold Shock-Absorbing Reclining Infant Strollers"),
					makeLeaf("solid-wood-convertible-baby-cribs", "3-in-1 Teething Rail Guardrail Height Adjustable Timber Baby Cribs", "የእንጨት አስተማማኝ የህፃናት አልጋዎች ከነፍራሻቸው", "Siree Daa'immanii Mukaa Jabaa", "ናይ 3-in-1 Teething Rail Guardrail Height Adjustable Timber Baby Cribs"),
					makeLeaf("isofix-certified-baby-car-seats", "360-Degree Swivel Rear & Forward Facing Newborn to 12-Year Car Seats", "ለመኪና ውስጥ የህፃናት የደህንነት ወንበሮች (Car Seat)", "Barcuma Daa'immanii Konkolaataa (Car Seat)", "ናይ 360-Degree Swivel Rear & Forward Facing Newborn to 12-Year Car Seats"),
					makeLeaf("musical-interactive-baby-walkers", "Anti-Rollover Multi-Height Play Tray Early Development Walkers", "ሙዚቃ ያላቸው የማይንሸራተቱ የህፃናት መራመጃዎች (Walker)", "Meeshaa Deemsaa Daa'immanii (Walker)", "ናይ Anti-Rollover Multi-Height Play Tray Early Development Walkers"),
					makeLeaf("ergonomic-baby-carriers-hipseats", "Breathable Mesh Kangaroo Pouch Front & Back Ergonomic Hip Seat Carriers", "ህፃናትን በጀርባ እና በደረቱ የሚያዝሉ ምቹ ካንጋሮዎች", "Baattuu Daa'immanii Qananii (Carrier)", "ናይ Breathable Mesh Kangaroo Pouch Front & Back Ergonomic Hip Seat Carriers"),
					makeLeaf("foldable-portable-baby-playpens", "Breathable Mesh Zipper Gate Indoor & Outdoor Baby Safety Playards", "ህፃናት እንዳይወድቁ የሚጫወቱባቸው የታጠሩ አልጋዎች (Playpen)", "Mana Tapha Daa'immanii Nageenya Qabu", "ናይ Breathable Mesh Zipper Gate Indoor & Outdoor Baby Safety Playards"),
					makeLeaf("electric-baby-swings-rockers", "Bluetooth Music Remote Controlled Automatic Bouncer Rocking Chairs", "በኤሌክትሪክ የሚወዛወዙ እና የሚያጫውቱ የህፃናት ወንበሮች", "Barcuma Daa'immanii Socho'u (Swing)", "ናይ Bluetooth Music Remote Controlled Automatic Bouncer Rocking Chairs"),
					makeLeaf("folding-baby-bathtubs-temperature-sensor", "Collapsible Non-Slip Newborn Bath Basin with Real-Time Thermometer", "የሙቀት መጠን የሚያሳዩ የሚታጠፉ የህፃናት መታጠቢያ ገንዳዎች", "Qodaa Dhiqannaa Daa'immanii Dachaafamu", "ናይ Collapsible Non-Slip Newborn Bath Basin with Real-Time Thermometer"),
				},
			},
			{
				Slug:   "kids-educational-toys-remote-cars",
				NameEN: "Educational Toys, STEM Kits and Ride-On Electric Cars",
				NameAM: "ትምህርታዊ አሻንጉሊቶች፣ የሮቦቲክስ እቃዎች እና የህፃናት መኪኖች",
				NameOM: "Meeshaa Tapha Daa'immanii fi Konkolaataa",
				NameTI: "ናይ ቆልዑ መጻወትን ናይ ኤሌክትሪክ መካይንን",
				Items: []LeafItemDef{
					makeLeaf("battery-powered-kids-ride-on-cars", "12V Dual Motor Remote Controlled Licensed SUV Ride-On Electric Cars", "በባትሪ የሚሰሩ እና በሪሞት የሚነዱ የህፃናት መኪኖች", "Konkolaataa Daa'immanii Elektiriikii", "ናይ 12V Dual Motor Remote Controlled Licensed SUV Ride-On Electric Cars"),
					makeLeaf("classic-wooden-educational-building-blocks", "Montessori Color & Shape Sorting Solid Hardwood Stacking Blocks", "የእንጨት ቅርጽ እና ቀለም መለያ መማሪያ ብሎኮች", "Muka Ijaarsa Tapha Daa'immanii", "ናይ Montessori Color & Shape Sorting Solid Hardwood Stacking Blocks"),
					makeLeaf("stem-robotics-electronic-science-kits", "Solar Powered Programmable DIY Mechanical Science Lab Toys for Kids", "የሳይንስ እና ሮቦቲክስ መስሪያ የልጆች መማሪያ ኪቶች", "Meeshaa Barnoota Saayinsii fi Roobootiksii", "ናይ Solar Powered Programmable DIY Mechanical Science Lab Toys for Kids"),
					makeLeaf("high-speed-rc-remote-control-monster-trucks", "4WD Waterproof Off-Road 40km/h Hobby Grade RC Buggy Trucks", "በሪሞት የሚሽከረከሩ ፈጣን የስፖርት መኪኖች (RC Cars)", "Konkolaataa Taphaa Rimoottii (RC Car)", "ናይ 4WD Waterproof Off-Road 40km/h Hobby Grade RC Buggy Trucks"),
					makeLeaf("plush-soft-cuddle-stuffed-animal-dolls", "Hypoallergenic Giant Teddy Bear & Cartoon Character Plush Toys", "ለስላሳ ትላልቅ የቴዲ ቤር አሻንጉሊቶች (Teddy Bear)", "Taphaa Daa'immanii Lallaafaa (Teddy Bear)", "ናይ Hypoallergenic Giant Teddy Bear & Cartoon Character Plush Toys"),
					makeLeaf("kids-play-kitchen-cooking-food-sets", "Light & Sound Simulation Pretend Chef Cooking Kitchen Playsets", "የምግብ ማብሰያ እና የኩሽና እቃዎች መጫወቻ ስብስቦች", "Kushiinaa Tapha Daa'immanii", "ናይ Light & Sound Simulation Pretend Chef Cooking Kitchen Playsets"),
					makeLeaf("magnetic-building-tiles-3d-puzzle-sets", "100-Piece Diamond Magnet Architectural Geometric Castle Tiles", "በማግኔት የሚገጣጠሙ ዘመናዊ የልጆች መገንቢያ ኪቶች", "Taphaa Magneetii Ijaarsa Gamoo", "ናይ 100-Piece Diamond Magnet Architectural Geometric Castle Tiles"),
					makeLeaf("children-electric-keyboards-drums-karaoke", "Microphone Integrated Early Music Learning Mini Piano Keyboards", "የልጆች የፒያኖ እና የከበሮ መጫወቻ መሳሪያዎች", "Kiiboordii fi Diraamii Tapha Daa'immanii", "ናይ Microphone Integrated Early Music Learning Mini Piano Keyboards"),
				},
			},
		},
	})

	// Domain: pet-supplies
	list = append(list, DomainCatalogDef{
		RootSlug: "pet-supplies",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "dog-cat-food-nutrition",
				NameEN: "Dog Food, Cat Food and Pet Care Nutrition",
				NameAM: "የውሻ እና የድመት ምግቦች፣ ብስኩቶች እና አልሚ ምግቦች",
				NameOM: "Nyaata Sarootaa fi Adurree",
				NameTI: "ናይ ከልቢን ድሙን ምግቢ",
				Items: []LeafItemDef{
					makeLeaf("dry-kibble-adult-dog-food-bags-20kg", "High-Protein 26% Chicken & Beef Complete Balanced Dog Kibble (20kg)", "የውሻ የደረቅ ምግብ (ኪብል) በ20 ኪሎ ከረጢት", "Nyaata Sarootaa Gogaa (Kiiloo 20)", "ናይ High-Protein 26% Chicken & Beef Complete Balanced Dog Kibble (20kg)"),
					makeLeaf("canned-wet-dog-food-chunks-in-gravy", "Real Meat Chunk Nutritious Canned Wet Dog Food (400g Tins)", "የታሸገ የስጋ እና የዶሮ የውሻ እርጥብ ምግብ በጣሳ", "Nyaata Sarootaa Saamunaa Foonii", "ናይ Real Meat Chunk Nutritious Canned Wet Dog Food (400g Tins)"),
					makeLeaf("dry-cat-food-hairball-control-10kg", "Salmon & Tuna Omega-3 Urinary & Hairball Care Cat Food (10kg)", "የድመት የዓሳ እና የዶሮ ደረቅ ምግብ በከረጢት", "Nyaata Adurree Qurxummii (Kiiloo 10)", "ናይ Salmon & Tuna Omega-3 Urinary & Hairball Care Cat Food (10kg)"),
					makeLeaf("puppy-growth-formula-milk-replacers", "Calcium Enriched Weaning Puppy Milk Replacer Formula Powder", "ለትንንሽ ቡችሎች የሚሆን አጥቢ የዱቄት ወተት", "Aanan Cuucii Sarootaa Daakuu", "ናይ Calcium Enriched Weaning Puppy Milk Replacer Formula Powder"),
					makeLeaf("dental-care-chew-sticks-dog-treats", "Tartar Removing Mint Flavored Beef Rawhide Dog Chew Bones", "ጥርስ የሚያጸዱ እና የሚታመኩ የውሻ አጥንቶች እና ብስኩቶች", "Lafee fi Nyaata Sarootaa Qulqulleessu", "ናይ Tartar Removing Mint Flavored Beef Rawhide Dog Chew Bones"),
					makeLeaf("clumping-odor-control-bentonite-cat-litter", "Dust-Free Lavender Scented Instant Clumping Clay Cat Litter (10kg)", "ሽታ የማያመጣ የድመት ቆሻሻ መጣያ አሸዋ (Cat Litter)", "Cirracha Adurree Xurii Qulqulleessu", "ናይ Dust-Free Lavender Scented Instant Clumping Clay Cat Litter (10kg)"),
					makeLeaf("glucosamine-joint-supplements-dogs", "Chondroitin Hip & Joint Mobility Chewable Tablets for Dogs", "የውሻ መገጣጠሚያ እና አጥንት ማጠናከሪያ ቫይታሚኖች", "Qoricha fi Vaayitaaminii Lafee Sarootaa", "ናይ Chondroitin Hip & Joint Mobility Chewable Tablets for Dogs"),
					makeLeaf("balanced-caged-bird-seed-mixes", "Millet, Sunflower & Canary Seed Cleaned Bird Feed Mixes", "የወፍ እና የፓሮት የተመጣጠነ እህል መኖ በጥቅል", "Nyaata Simbirroo fi Baarotii", "ናይ Millet, Sunflower & Canary Seed Cleaned Bird Feed Mixes"),
				},
			},
			{
				Slug:   "pet-accessories-housing-grooming",
				NameEN: "Pet Cages, Collars, Leashes and Grooming",
				NameAM: "የውሻ ቤቶች፣ አንገትጌዎች፣ ማሰሪያዎች እና ማበጠሪያዎች",
				NameOM: "Mana Sarootaa, Qabduu fi Meeshaalee",
				NameTI: "ናይ ከልቢ ገዛን ማእሰርን መመሸጥን",
				Items: []LeafItemDef{
					makeLeaf("heavy-duty-metal-wire-folding-dog-crates", "Double Door Leak-Proof Plastic Pan Large Dog Travel Training Crates", "የብረት የሚታጠፉ የውሻ ማቆያ ካኬቶች (Crate)", "Koorjii Sibilaa Sarootaa Dachaafamu", "ናይ Double Door Leak-Proof Plastic Pan Large Dog Travel Training Crates"),
					makeLeaf("insulated-weatherproof-wooden-outdoor-dog-houses", "Asphalt Roof Treated Timber Warm Outdoor Kennel Dog Houses", "የእንጨት ጥራት ያላቸው የውሻ ማደሪያ ቤቶች (Kennel)", "Mana Sarootaa Mukaa Alaa (Kennel)", "ናይ Asphalt Roof Treated Timber Warm Outdoor Kennel Dog Houses"),
					makeLeaf("retractable-heavy-duty-dog-leashes", "5-Meter One-Button Brake Heavy Nylon Cord Walking Dog Leashes", "የሚረዝም እና የሚያጥር የውሻ መራመጃ ጠንካራ ገመድ", "Haada Deemsaa Sarootaa (Leash)", "ናይ 5-Meter One-Button Brake Heavy Nylon Cord Walking Dog Leashes"),
					makeLeaf("tactical-padded-k9-dog-harnesses", "Military Grade No-Pull Handle Molle Training Vest Dog Harnesses", "የደረት ማሰሪያ ጠንካራ የውሻ ልጓሞች (Harness)", "Qabattoo Qomaa Sarootaa Jabaa", "ናይ Military Grade No-Pull Handle Molle Training Vest Dog Harnesses"),
					makeLeaf("professional-electric-pet-hair-clippers", "Low-Noise Rechargeable Ceramic Blade Dog & Cat Grooming Shavers", "የውሻ እና የድመት ጸጉር መቁረጫ የኤሌክትሪክ ማሽኖች", "Maashinii Rifeensa Sarootaa Muru", "ናይ Low-Noise Rechargeable Ceramic Blade Dog & Cat Grooming Shavers"),
					makeLeaf("self-cleaning-pet-slicker-brushes", "Stainless Steel Pin Hair Shedding & Tangling Removal Brushes", "የጸጉር መርገፍን የሚከላከሉ የውሻ ማበጠሪያዎች", "Filtooma Rifeensa Sarootaa fi Adurree", "ናይ Stainless Steel Pin Hair Shedding & Tangling Removal Brushes"),
					makeLeaf("stainless-steel-double-pet-feeding-bowls", "Non-Skid Silicone Mat Food & Water Elevated Double Pet Bowls", "የማይንሸራተቱ የውሻ እና ድመት ምግብ እና ውሃ ሳህኖች", "Qodaa Nyaataa fi Bishaanii Sarootaa", "ናይ Non-Skid Silicone Mat Food & Water Elevated Double Pet Bowls"),
					makeLeaf("flea-tick-medicated-collars-dogs-cats", "8-Month Continuous Protection Waterproof Anti-Flea Neck Collars", "ቁንጫ እና መዥገር የሚከላከሉ የውሻ አንገትጌዎች", "Qabdoo Mormaa Silmii fi Qunxaa Ittisu", "ናይ 8-Month Continuous Protection Waterproof Anti-Flea Neck Collars"),
				},
			},
		},
	})

	// Domain: raw-materials-chemicals
	list = append(list, DomainCatalogDef{
		RootSlug: "raw-materials-chemicals",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "industrial-chemicals-solvents-resins",
				NameEN: "Industrial Chemicals, Solvents and Water Treatment",
				NameAM: "የኢንዱስትሪ ኬሚካሎች፣ ሳሙና ማምረቻ እና የውሃ ማከሚያዎች",
				NameOM: "Keemikaala Warshaa fi Qulqulleessa Bishaanii",
				NameTI: "ናይ ፋብሪካ ኬሚካላትን ናይ ማይ መጽረይን",
				Items: []LeafItemDef{
					makeLeaf("caustic-soda-pearls-sodium-hydroxide", "99% Pure Industrial Grade Solid Caustic Soda Flakes & Pearls (25kg)", "ኮስቲክ ሶዳ ለሳሙና ማምረቻ (99% ንጹህ በ25 ኪሎ ጆንያ)", "Koostik Soodaa Warshaa Saamunaa (25kg)", "ናይ 99% Pure Industrial Grade Solid Caustic Soda Flakes & Pearls (25kg)"),
					makeLeaf("sodium-lauryl-ether-sulfate-sles-70", "High-Foaming Liquid Detergent & Shampoo Raw Material SLES 70% Drums", "ፈሳሽ ሳሙና እና ሻምፑ ማምረቻ ጥሬ እቃ (SLES 70% በበርሜል)", "Keemikaala Saamunaa Dhangala'aa (SLES 70%)", "ናይ High-Foaming Liquid Detergent & Shampoo Raw Material SLES 70% Drums"),
					makeLeaf("calcium-hypochlorite-chlorine-65-70", "Granular Swimming Pool & Water Treatment Disinfectant Chlorine", "የውሃ ማከሚያ እና ጀርም ማጥፊያ ንጹህ ክሎሪን (በበርሜል)", "Kilooriinii Bishaan Qulqulleessu (Chlorine)", "ናይ Granular Swimming Pool & Water Treatment Disinfectant Chlorine"),
					makeLeaf("sulfuric-hydrochloric-industrial-acids", "Industrial Grade Technical Concentrated Acid Drums for Mining & Tannery", "የቆዳ እና የማዕድን ፋብሪካ ሰልፈሪክ እና ሃይድሮክሎሪክ አሲዶች", "Asiidii Warshaa fi Gogaa (Acid)", "ናይ Industrial Grade Technical Concentrated Acid Drums for Mining & Tannery"),
					makeLeaf("linear-alkylbenzene-sulfonic-acid-labsa", "96% Active Detergent Powder & Liquid Washing Raw Material LABSA", "የዱቄት ሳሙና እና የፈሳሽ ሳሙና ማምረቻ ላብሳ (LABSA)", "Keemikaala Saamunaa Daakuu (LABSA)", "ናይ 96% Active Detergent Powder & Liquid Washing Raw Material LABSA"),
					makeLeaf("virgin-polypropylene-pp-plastic-granules", "Injection & Raffia Grade Virgin PP Plastic Raw Material Pellets", "የላስቲክ እና ፕላስቲክ እቃዎች ማምረቻ ፒፒ ፍሬዎች (PP Granules)", "Daakuu fi Firii Pilaastikii (PP Granules)", "ናይ Injection & Raffia Grade Virgin PP Plastic Raw Material Pellets"),
					makeLeaf("high-density-polyethylene-hdpe-granules", "Blow Molding Pipe & Bottle Grade Virgin HDPE Resin Pellets", "ቧንቧ እና ጀሪካን ማምረቻ ኤችዲፒኢ ጥሬ እቃዎች (HDPE)", "Firii Pilaastikii Ujummoo (HDPE Granules)", "ናይ Blow Molding Pipe & Bottle Grade Virgin HDPE Resin Pellets"),
					makeLeaf("industrial-grade-ethanol-denatured-alcohol", "96% Pure Ethyl Alcohol Solvent for Disinfectant & Paint Industry", "የህክምና እና የፋብሪካ 96% አልኮል በበርሜል", "Alkoolii Warshaa fi Qorichaa 96%", "ናይ 96% Pure Ethyl Alcohol Solvent for Disinfectant & Paint Industry"),
				},
			},
			{
				Slug:   "metals-minerals-raw-tannery",
				NameEN: "Raw Metals, Minerals, Timber and Leather Hides",
				NameAM: "ጥሬ ብረቶች፣ ማዕድናት፣ ጣውላ እና ጥሬ ቆዳዎች",
				NameOM: "Sibila, Dhagaa, Mukaa fi Gogaa Dheedhii",
				NameTI: "ጥረ ሓጺንን ማዕድናትን ዕንጨይትን ቆርበትን",
				Items: []LeafItemDef{
					makeLeaf("pure-copper-ingots-cathode-plates", "Grade-A 99.99% Electrolytic Copper Cathode Sheets & Raw Ingots", "ንጹህ የመዳብ ሰሌዳዎች እና ጥሬ መዳብ (Copper)", "Sibila Qulqulluu Koompparaa (Copper)", "ናይ Grade-A 99.99% Electrolytic Copper Cathode Sheets & Raw Ingots"),
					makeLeaf("virgin-aluminum-ingots-extrusion-billets", "99.7% Purity Aluminum Alloy Smelting Ingots for Foundries", "የአሉሚኒየም ጥሬ ማቅለጫ ብረቶች (Aluminum Ingots)", "Sibila Alumiiniyemii Dheedhii", "ናይ 99.7% Purity Aluminum Alloy Smelting Ingots for Foundries"),
					makeLeaf("pure-limestone-calcium-carbonate-powder", "Micronized 400-Mesh Coated Calcium Carbonate Industrial Minerals", "የካልሲየም ካርቦኔት እና የኖራ ድንጋይ ዱቄት በጅምላ", "Daakuu Dhagaa Kaalsiyeemii (Limestone)", "ናይ Micronized 400-Mesh Coated Calcium Carbonate Industrial Minerals"),
					makeLeaf("industrial-grade-gypsum-raw-rocks", "Crushed High-Purity Raw Gypsum for Plaster & Cement Factories", "ለሲሚንቶ እና ጂፕሰም ፋብሪካ የሚሆን ጥሬ የጂፕሰም ድንጋይ", "Dhagaa Jiipsamii Dheedhii", "ናይ Crushed High-Purity Raw Gypsum for Plaster & Cement Factories"),
					makeLeaf("wet-blue-ethiopian-sheep-leather-hides", "Export Tanned High-Tensile Soft Ethiopian Highland Sheep Hides", "ጥራት ያላቸው የኢትዮጵያ በግ እና ፍየል ቆዳዎች (Wet Blue)", "Gogaa Hoolaa fi Re'ee Warshaa (Wet Blue)", "ናይ Export Tanned High-Tensile Soft Ethiopian Highland Sheep Hides"),
					makeLeaf("eucalyptus-treated-construction-timber-poles", "Termite Resistant 10-12cm Diameter Eucalyptus Wood Scaffold Poles", "የባህር ዛፍ ጠንካራ የህንፃ መስሪያ ማገር እና ምሶሶዎች", "Muka Baargamoo Gamoo Jabaa", "ናይ Termite Resistant 10-12cm Diameter Eucalyptus Wood Scaffold Poles"),
					makeLeaf("kiln-dried-hardwood-timber-lumber-planks", "Planed High-Grade Tid, Wanza & Zigba Solid Furniture Planks", "የደረቁ የጥድ፣ ዋንዛ እና ዝግባ ጥራት ያላቸው የእንጨት ሳንቃዎች", "Muka Saanqaa Wanzaa fi Qirxii", "ናይ Planed High-Grade Tid, Wanza & Zigba Solid Furniture Planks"),
					makeLeaf("pure-silica-sand-glass-manufacturing", "Washed 99% SiO2 High Purity White Silica Quartz Sand for Glass", "ለመስታወት ፋብሪካ የሚሆን ንጹህ ነጭ የሲሊካ አሸዋ", "Cirracha Doowwituu Qulqulluu (Silica)", "ናይ Washed 99% SiO2 High Purity White Silica Quartz Sand for Glass"),
				},
			},
		},
	})

	// Domain: energy-solar-power
	list = append(list, DomainCatalogDef{
		RootSlug: "energy-solar-power",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "solar-panels-inverters-batteries",
				NameEN: "Solar Panels, Lithium Batteries and Inverters",
				NameAM: "የፀሐይ ብርሃን ፓነሎች፣ ሊቲየም ባትሪዎች እና ኢንቨርተሮች",
				NameOM: "Paaneelii Aduu, Baatirii fi Invartara",
				NameTI: "ናይ ፀሓይ ፓነላትን ባትሪን ኢንቨርተርን",
				Items: []LeafItemDef{
					makeLeaf("monocrystalline-tier-1-solar-panels-550w", "Half-Cut Cell High-Efficiency 550W Tier-1 Monocrystalline Solar Panels", "ከፍተኛ ኃይል የሚያመነጩ 550 ዋት የሞኖክሪስታላይን የሶላር ፓነሎች", "Paaneelii Aduu Humna Guddaa 550W", "ናይ Half-Cut Cell High-Efficiency 550W Tier-1 Monocrystalline Solar Panels"),
					makeLeaf("bifacial-dual-glass-solar-panels-650w", "Front & Rear Generation Utility Grade 650W Bifacial Solar Modules", "ከፊትና ከኋላ ኃይል የሚስቡ 650 ዋት ግዙፍ የሶላር ፓነሎች", "Paaneelii Aduu Fuula Lamee 650W", "ናይ Front & Rear Generation Utility Grade 650W Bifacial Solar Modules"),
					makeLeaf("lifepo4-lithium-wall-mount-batteries-48v", "100Ah/200Ah 5kWh-10kWh 6000-Cycle Smart BMS Lithium Solar Batteries", "ለቤት እና ለድርጅት የሚሆኑ ዘመናዊ የሊቲየም ባትሪዎች (LiFePO4)", "Baatirii Liitiyeemii Aduu (LiFePO4)", "ናይ 100Ah/200Ah 5kWh-10kWh 6000-Cycle Smart BMS Lithium Solar Batteries"),
					makeLeaf("pure-sine-wave-hybrid-solar-inverters", "5kW-10kW MPPT 48V Off-Grid & Grid-Tie Pure Sine Wave Inverters", "መብራት እና ሶላር በአንድ ላይ የሚያቀናጁ ሃይብሪድ ኢንቨርተሮች (5kW/10kW)", "Invartara Aduu fi Humnaa (Hybrid)", "ናይ 5kW-10kW MPPT 48V Off-Grid & Grid-Tie Pure Sine Wave Inverters"),
					makeLeaf("deep-cycle-tubular-gel-batteries-200ah", "12V 200Ah Heavy Duty Sealed Solar Gel Lead-Carbon Energy Batteries", "12 ቮልት 200Ah ጥራት ያላቸው የሶላር ጄል ባትሪዎች", "Baatirii Jelii Aduu 200Ah", "ናይ 12V 200Ah Heavy Duty Sealed Solar Gel Lead-Carbon Energy Batteries"),
					makeLeaf("solar-water-pumping-inverters-vfd", "3-Phase 380V Solar Powered Agricultural Borehole VFD Inverters", "ያለ ባትሪ በቀጥታ በሶላር የውሃ ፓምፕ የሚያሰሩ ኢንቨርተሮች", "Invartara Paampii Bishaanii Aduu", "ናይ 3-Phase 380V Solar Powered Agricultural Borehole VFD Inverters"),
					makeLeaf("solar-street-light-all-in-one-led", "100W-300W Radar Motion Sensor Integrated Solar Street Lights", "በራሳቸው የሚበሩ እና የሚጠፉ የፀሐይ ብርሃን የመንገድ መብራቶች", "Ibsaa Daandii Aduu Ofiin Hojjetu", "ናይ 100W-300W Radar Motion Sensor Integrated Solar Street Lights"),
					makeLeaf("solar-home-dc-lighting-kits-payg", "3-Bulb USB Phone Charger FM Radio Portable Rural Solar Kits", "ለገጠር ቤት የሚሆኑ 3 አምፑል፣ ሬዲዮ እና ስልክ ቻርጀር ያላቸው ሶላሮች", "Kiteetii Ibsaa Aduu Manaa fi Raadiyoo", "ናይ 3-Bulb USB Phone Charger FM Radio Portable Rural Solar Kits"),
				},
			},
			{
				Slug:   "diesel-generators-voltage-stabilizers",
				NameEN: "Diesel Generators, AVR Stabilizers and Wind Power",
				NameAM: "የናፍጣ ጄኔሬተሮች፣ የቮልቴጅ ማረጋጊያዎች (AVR)",
				NameOM: "Jenereetara Diiyezilaa fi Istabilaayizara",
				NameTI: "ናይ ናፍጣ ጄኔሬተርን ስታቢላይዘርን",
				Items: []LeafItemDef{
					makeLeaf("silent-soundproof-diesel-generators-20kva", "Perkins/Cummins 20kVA-50kVA Ultra-Quiet Canopy Diesel Generators", "ድምጽ የማይረብሹ 20kVA-50kVA የቢሮ እና የቤት ጄኔሬተሮች", "Jenereetara Sagalee Qabbanaawaa 20kVA", "ናይ Perkins/Cummins 20kVA-50kVA Ultra-Quiet Canopy Diesel Generators"),
					makeLeaf("heavy-industrial-diesel-generators-250kva", "Heavy Duty Factory Power 250kVA-500kVA Water Cooled Generators", "ለፋብሪካ እና ትላልቅ ህንፃዎች የሚሆኑ 250kVA-500kVA ጄኔሬተሮች", "Jenereetara Warshaa Humna Guddaa 250kVA", "ናይ Heavy Duty Factory Power 250kVA-500kVA Water Cooled Generators"),
					makeLeaf("automatic-mains-failure-ats-panels", "Automatic Transfer Switch (ATS) 100A-400A Generator Starting Panels", "መብራት ሲጠፋ ጄኔሬተር በራሱ የሚያስነሱ አውቶማቲክ የATS ሳጥኖች", "Panalii Jenereetara Ofiin Kaasu (ATS)", "ናይ Automatic Transfer Switch (ATS) 100A-400A Generator Starting Panels"),
					makeLeaf("servo-motor-automatic-voltage-regulators", "5kVA-30kVA 3-Phase Precision Servo Voltage Stabilizers (AVR)", "የመብራት መውረድ እና መወላገድ የሚያስተካክሉ ስታቢላይዘሮች (AVR)", "Istabilaayizara Elektiriikii (AVR)", "ናይ 5kVA-30kVA 3-Phase Precision Servo Voltage Stabilizers (AVR)"),
					makeLeaf("portable-silent-petrol-inverter-generators", "2kVA-3kVA Lightweight Low Fuel Consumption Petrol Inverters", "ቀላል እና ነዳጅ ቆጣቢ ተንቀሳቃሽ የቤንዚን ጄኔሬተሮች", "Jenereetara Benziinaa Salphaa 2kVA", "ናይ 2kVA-3kVA Lightweight Low Fuel Consumption Petrol Inverters"),
					makeLeaf("high-capacity-portable-power-stations", "1000W-2400W Fast AC Charging Lithium Battery Power Stations", "በቀላሉ የሚንቀሳቀሱ በኤሌክትሪክ የሚሞሉ ፓወር ስቴሽኖች", "Buufata Humnaa Liitiyeemii (Power Station)", "ናይ 1000W-2400W Fast AC Charging Lithium Battery Power Stations"),
					makeLeaf("industrial-battery-storage-systems-bess", "Containerized Megawatt Scale Commercial Peak Shaving Solar BESS", "ለትላልቅ ፋብሪካዎች የሚሆኑ ግዙፍ የኮንቴነር ባትሪ ሲስተሞች", "Kuusaa Humnaa Gurguddaa (BESS)", "ናይ Containerized Megawatt Scale Commercial Peak Shaving Solar BESS"),
					makeLeaf("micro-wind-turbine-generators-1kw", "Off-Grid Low Wind Speed Starting Wind Turbine Generators", "በንፋስ ኃይል መብራት የሚያመነጩ አነስተኛ ተርባይኖች", "Maashinii Qilleensaan Humna Maddisiisu", "ናይ Off-Grid Low Wind Speed Starting Wind Turbine Generators"),
				},
			},
		},
	})

	// Domain: water-irrigation-equipment
	list = append(list, DomainCatalogDef{
		RootSlug: "water-irrigation-equipment",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "water-purification-filtration-plants",
				NameEN: "Water Filtration, RO Plants and Dispensers",
				NameAM: "የውሃ ማጣሪያዎች፣ የሮተር ፕላንቶች እና የውሃ ማቀዝቀዣዎች",
				NameOM: "Bishaan Qulqulleessaa fi Meeshaalee",
				NameTI: "ናይ ማይ መጽረይን መዝሓልን",
				Items: []LeafItemDef{
					makeLeaf("commercial-reverse-osmosis-ro-water-plants", "1,000 to 5,000 Liters/Hour Stainless Steel RO Water Purifiers", "በሰዓት ከ1,000 እስከ 5,000 ሊትር የሚያጣሩ የፋብሪካ RO ማጣሪያዎች", "Warshaa Bishaan Qulqulleessu (RO Plant)", "ናይ 1,000 to 5,000 Liters/Hour Stainless Steel RO Water Purifiers"),
					makeLeaf("home-under-sink-5-stage-ro-purifiers", "Multi-Stage Sediment, Carbon & Membrane Direct Drink Purifiers", "ከወጥ ቤት ሲንክ ስር የሚገጠሙ የቤት ውስጥ የውሃ ማጣሪያዎች", "Qulqulleessaa Bishaanii Kushiinaa (5-Stage RO)", "ናይ Multi-Stage Sediment, Carbon & Membrane Direct Drink Purifiers"),
					makeLeaf("hot-cold-compressor-water-dispensers", "Top-Loading & Bottom-Loading Floor Standing Clean Water Dispensers", "ሙቅ እና ቀዝቃዛ ውሃ የሚያወጡ የቢሮ የውሃ ዲስፔንሰሮች", "Dispeensara Bishaanii Ho'aa fi Qabbanaawaa", "ናይ Top-Loading & Bottom-Loading Floor Standing Clean Water Dispensers"),
					makeLeaf("heavy-fiberglass-sand-carbon-filter-tanks", "Commercial Multimedia Sand, Carbon & Iron Removal Filter Vessels", "የአሸዋ እና የካርቦን ውሃ ማጣሪያ ትላልቅ የፋይበር ታንከሮች", "Taankii Qulqulleessaa Cirrachaa fi Kaarboonii", "ናይ Commercial Multimedia Sand, Carbon & Iron Removal Filter Vessels"),
					makeLeaf("ultra-filtration-uf-membrane-modules", "High Flux Hollow Fiber Low Pressure Suspended Solid Filters", "የውሃ ቆሻሻ እና ዝቃጭ የሚያጣሩ አልትራ ፊልትሬሽን ሜምብሬኖች", "Meembireenii Qulqulleessaa Bishaanii (UF)", "ናይ High Flux Hollow Fiber Low Pressure Suspended Solid Filters"),
					makeLeaf("ultraviolet-uv-water-sterilizers", "Stainless Steel Chamber Fast Bacteria Killing Disinfection UV Lamps", "ጀርም እና ባክቴሪያ በብርሃን የሚያጠፉ የUV ማምከኛ መብራቶች", "Ibsaa Bakteeriyaa Ajjeesu (UV Sterilizer)", "ናይ Stainless Steel Chamber Fast Bacteria Killing Disinfection UV Lamps"),
					makeLeaf("water-softener-ion-exchange-systems", "Hard Well Water Calcium Lime Scaling Removal Automatic Softeners", "ጨዋማ እና ጠጣር የጉድጓድ ውሃ የሚያለሰልሱ አውቶማቲክ ማጣሪያዎች", "Meeshaa Bishaan Cabaa Lallafisu", "ናይ Hard Well Water Calcium Lime Scaling Removal Automatic Softeners"),
					makeLeaf("portable-gravity-camping-water-filter-buckets", "Ceramic Candle Emergency Bacteria Free Drinking Water Purifiers", "መብራት የማይፈልጉ የሴራሚክ ሻማ የጠረጴዛ ላይ ውሃ ማጣሪያዎች", "Qulqulleessaa Bishaanii Seeraamikii", "ናይ Ceramic Candle Emergency Bacteria Free Drinking Water Purifiers"),
				},
			},
			{
				Slug:   "waste-management-recycling-equipment",
				NameEN: "Waste Management, Recycling and Sewage Treatment",
				NameAM: "የቆሻሻ አወጋገድ፣ የፕላስቲክ መፍጫ እና የፍሳሽ ማጣሪያ እቃዎች",
				NameOM: "Qulqullina Balfaa, Deebisanii Hojjechuu fi Ujummoo",
				NameTI: "ናይ ጐሓፍ ምእላይን ሪሳይክሊንግን",
				Items: []LeafItemDef{
					makeLeaf("heavy-waste-plastic-crusher-granulators", "Heavy-Duty Rotary Knife Plastic Bottle, Pipe & Scrap Shredders", "የፕላስቲክ ኮዳዎች እና ፌስታሎች መፍጫ ኃይለኛ ማሽኖች", "Maashinii Pilaastikii Caccabsu (Crusher)", "ናይ Heavy-Duty Rotary Knife Plastic Bottle, Pipe & Scrap Shredders"),
					makeLeaf("hospital-medical-waste-incinerators", "Smokeless Dual-Chamber High Temperature Healthcare Incinerators", "የሆስፒታል እና የክሊኒክ ቆሻሻ ማቃጠያ ከጭስ ነጻ ማሽኖች (Incinerator)", "Gubituu Balfa Hospitaalaa (Incinerator)", "ናይ Smokeless Dual-Chamber High Temperature Healthcare Incinerators"),
					makeLeaf("outdoor-wheelie-waste-bins-240l", "Heavy-Duty UV Resistant HDPE 2-Wheeled Municipal Garbage Carts", "የግቢ እና የከተማ ቆሻሻ መጣያ ባለ ጎማ ትላልቅ ባልዲዎች (240L)", "Saanduqa Balfaa Konkolaatu 240L", "ናይ Heavy-Duty UV Resistant HDPE 2-Wheeled Municipal Garbage Carts"),
					makeLeaf("commercial-bio-gas-digester-plants", "Prefabricated Membrane Farm Livestock Manure Clean Energy Digesters", "ከከብት ፍግ ጋዝ እና ማዳበሪያ የሚያመርቱ የባዮጋዝ ፕላንቶች", "Warshaa Baayogaasii Fandii Horiirraa", "ናይ Prefabricated Membrane Farm Livestock Manure Clean Energy Digesters"),
					makeLeaf("sewage-sludge-heavy-submersible-pumps", "Cast Iron Non-Clogging Cutter Impeller Waste Drainage Pumps", "የፍሳሽ እና የቆሻሻ ውሃ መሳቢያ ቆራጭ ሰመጥ ፓምፖች", "Paampii Bishaan Xuraawaa fi Dhangala'aa", "ናይ Cast Iron Non-Clogging Cutter Impeller Waste Drainage Pumps"),
					makeLeaf("grease-trap-interceptor-units-restaurant", "Stainless Steel Commercial Kitchen Oil and Grease Waste Separators", "የሬስቶራንት እና የሆቴል ኩሽና ስብ እና ዘይት መያዣ ሳጥኖች", "Qabduu Zayitaa fi Coomaa Kushiinaa", "ናይ Stainless Steel Commercial Kitchen Oil and Grease Waste Separators"),
					makeLeaf("compact-package-sewage-treatment-plants", "Aerobic MBBR High-Rise Building Treated Graywater Recycling Plants", "የህንፃ ፍሳሽ አጣርተው ለአትክልት ውሃነት የሚቀይሩ ማጣሪያዎች", "Warshaa Bishaan Xuraawaa Qulqulleessu", "ናይ Aerobic MBBR High-Rise Building Treated Graywater Recycling Plants"),
					makeLeaf("hydraulic-cardboard-paper-balers", "Vertical Waste Paper, Carton & Pet Plastic Compactor Balers", "ካርቶን እና ወረቀት መጨፍለቂያ ሀይድሮሊክ ማሽኖች", "Maashinii Kartoonii fi Waraqaa Dhiibu", "ናይ Vertical Waste Paper, Carton & Pet Plastic Compactor Balers"),
				},
			},
		},
	})

	// Domain: travel-tourism
	list = append(list, DomainCatalogDef{
		RootSlug: "travel-tourism",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "tours-safari-heritage-packages",
				NameEN: "Historic Route, Safari and Trekking Packages",
				NameAM: "የታሪክ መስመር ጉብኝቶች፣ የዱር እንስሳት እና የተራራ የእግር ጉዞ",
				NameOM: "Daawwanna Seenaa, Bosonaa fi Gaaraa",
				NameTI: "ናይ ታሪክን ሳፋሪን መገሻታት",
				Items: []LeafItemDef{
					makeLeaf("lalibela-rock-hewn-churches-guided-tours", "Comprehensive 3-Day Flight & Expert Guided Historic Lalibela Packages", "የላሊበላ ውቅር አብያተ ክርስቲያናት የተሟላ የጉብኝት ጥቅል", "Daawwanna Gamoo Dhagaa Laalibalaa", "ናይ Comprehensive 3-Day Flight & Expert Guided Historic Lalibela Packages"),
					makeLeaf("simien-mountains-trekking-expeditions", "Ras Dashen & Walia Ibex Scenery Mountain Trekking Expeditions", "የሰሜን ተራሮች እና የራስ ዳሸን የእግር ጉዞ እና የካምፒንግ ጉብኝት", "Imala Gaara Semeen fi Waaliyaa", "ናይ Ras Dashen & Walia Ibex Scenery Mountain Trekking Expeditions"),
					makeLeaf("danakil-depression-erta-ale-volcano-tours", "4x4 Desert Expedition to Erta Ale Active Lava Lake & Dallol Salt Flats", "የዳሎል እና ኤርታሌ እሳተ ጎመራ የ4x4 የበረሃ ጉዞ ጥቅል", "Imala Gammoojjii Daalol fi Irtale", "ናይ 4x4 Desert Expedition to Erta Ale Active Lava Lake & Dallol Salt Flats"),
					makeLeaf("omo-valley-cultural-tribal-photography-tours", "Mursi, Hamer & Karo Traditional Cultural Photography Safaris", "የኦሞ ሸለቆ ባህላዊ ብሔረሰቦች የፎቶግራፍ እና የባህል ጉብኝት", "Daawwanna Aadaa Sulula Omoo", "ናይ Mursi, Hamer & Karo Traditional Cultural Photography Safaris"),
					makeLeaf("bale-mountains-red-fox-wildlife-safaris", "Sanetti Plateau Endemic Ethiopian Red Wolf & Harenna Forest Treks", "የባሌ ተራሮች እና የቀይ ቀበሮ መመልከቻ የዱር እንስሳት ጉብኝት", "Daawwanna Jeedala Diimaa fi Gaara Baalee", "ናይ Sanetti Plateau Endemic Ethiopian Red Wolf & Harenna Forest Treks"),
					makeLeaf("gondar-castles-lake-tana-monastery-cruises", "Fasil Ghebbi Imperial Enclosure & Island Monastery Boat Cruises", "የፋሲል ግቢ ግንቦች እና የጣና ደሴት ገዳማት የጀልባ ጉዞ", "Daawwanna Masaraa Faasil fi Xaanaa", "ናይ Fasil Ghebbi Imperial Enclosure & Island Monastery Boat Cruises"),
					makeLeaf("harar-walled-city-hyena-feeding-tours", "Jugol UNESCO Walled City, Rimbaud House & Night Hyena Feeding", "የሀረር ጁጎል ግንብ እና የጅብ መመገብ ባህላዊ ጉብኝት", "Daawwanna Dabaqaa Harar Jugol", "ናይ Jugol UNESCO Walled City, Rimbaud House & Night Hyena Feeding"),
					makeLeaf("addis-ababa-city-highlights-stopover-tours", "National Museum Lucy, Mount Entoto & Mercato Guided Day Tours", "የአዲስ አበባ ከተማ እይታዎች (ሉሲ፣ እንጦጦ፣ መርካቶ) የ1 ቀን ጉብኝት", "Daawwanna Magaalaa Finfinnee fi Luusii", "ናይ National Museum Lucy, Mount Entoto & Mercato Guided Day Tours"),
				},
			},
			{
				Slug:   "hotel-lodging-resort-reservations",
				NameEN: "Resorts, Eco-Lodges and Hotel Bookings",
				NameAM: "ሪዞርቶች፣ የሀይቅ ዳር ማረፊያዎች እና የሆቴል ክፍሎች",
				NameOM: "Hoteela, Boqonnaa fi Bakka Ciisichaa",
				NameTI: "ናይ ሆቴልን ሪዞርትን ምዝገባ",
				Items: []LeafItemDef{
					makeLeaf("luxury-rift-valley-lake-resort-suites", "Lake Langano & Bishoftu Crater Waterfront Private Jacuzzi Suites", "የቢሾፍቱ እና ላንጋኖ ሀይቅ ዳርቻ የቅንጦት ሪዞርቶች", "Boqonnaa Haroo Bishoofituu fi Langaanoo", "ናይ Lake Langano & Bishoftu Crater Waterfront Private Jacuzzi Suites"),
					makeLeaf("eco-friendly-national-park-mountain-lodges", "Solar Powered Sustainable Stone & Thatched Mountain Eco-Lodges", "ከተፈጥሮ ጋር የተስማሙ የፓርክ ውስጥ ሎጆች (Eco-Lodge)", "Loojii Bosonaa fi Gaaraa Qananii", "ናይ Solar Powered Sustainable Stone & Thatched Mountain Eco-Lodges"),
					makeLeaf("addis-ababa-5-star-international-hotels", "Luxury Bole Diplomatic District 5-Star Hotel King Suites", "የአዲስ አበባ የ5 ኮከብ አለምአቀፍ ሆቴሎች የክብር ክፍሎች", "Hoteela Qananii Sadarkaa 5ffaa Finfinnee", "ናይ Luxury Bole Diplomatic District 5-Star Hotel King Suites"),
					makeLeaf("budget-friendly-clean-city-guest-houses", "En-Suite Bath Hot Water Wi-Fi Budget City Transit Guest Rooms", "ንጹህ እና ተመጣጣኝ ዋጋ ያላቸው የእንግዳ ማረፊያዎች", "Mana Boqonnaa Qulqulluu fi Baasii Xiqqaataa", "ናይ En-Suite Bath Hot Water Wi-Fi Budget City Transit Guest Rooms"),
					makeLeaf("hot-spring-wellness-spa-resort-stays", "Natural Thermal Mineral Water Healing Bath Pool Wellness Resorts", "የሶዶሬ እና ወንዶ ገነት የተፈጥሮ ፍልውሃ የጤና ሪዞርቶች", "Boqonnaa Bishaan Ho'aa Sodoree", "ናይ Natural Thermal Mineral Water Healing Bath Pool Wellness Resorts"),
					makeLeaf("conference-delegate-hotel-accommodations", "Multi-Day UN ECA & African Union Conference Hotel Room Blocks", "ለአፍሪካ ህብረት እና ስብሰባዎች የሚሆኑ የሆቴል ክፍሎች", "Iddoo Ciisicha Walga'ii Gamtaa Afrikaa", "ናይ Multi-Day UN ECA & African Union Conference Hotel Room Blocks"),
					makeLeaf("traditional-tukul-cultural-village-stays", "Authentic Grass-Roofed Conical Tukul Experience Stays with Meals", "በሳር ክዳን ባህላዊ ቱኩል ጎጆዎች ውስጥ ማረፍ እና መመገብ", "Mana Tukula Aadaa fi Boqonnaa", "ናይ Authentic Grass-Roofed Conical Tukul Experience Stays with Meals"),
					makeLeaf("camping-gear-rental-mountain-campsite-permits", "Complete 2-Person Trekking Tent, Foam & Stove Campsite Packages", "ለተራራ ጉዞ የሚሆኑ የካምፒንግ ድንኳኖች እና እቃዎች ኪራይ", "Kiree Dunkaanaa fi Meeshaa Kaampii", "ናይ Complete 2-Person Trekking Tent, Foam & Stove Campsite Packages"),
				},
			},
		},
	})

	// Domain: packaging-shipping
	list = append(list, DomainCatalogDef{
		RootSlug: "packaging-shipping",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "boxes-sacks-protective-packaging",
				NameEN: "Boxes, Sacks, Strapping and Protective Packaging",
				NameAM: "ካርቶኖች፣ ጆንያዎች፣ ፕላስቲኮች እና እቃ መጠቅለያዎች",
				NameOM: "Saanduqa Kartoonii, Qodaa fi Meeshaalee",
				NameTI: "ካርቶናትን ጆንያታትን ናይ መሸፈኒ ንዋያትን",
				Items: []LeafItemDef{
					makeLeaf("corrugated-shipping-carton-boxes-standard", "3-Ply & 5-Ply Heavy-Duty Export Shipping Brown Cardboard Boxes", "ጠንካራ ባለ 3 እና 5 ደርብ እቃ መጫኛ የካርቶን ሳጥኖች", "Saanduqa Kartoonii Jabaa (Carton Box)", "ናይ 3-Ply & 5-Ply Heavy-Duty Export Shipping Brown Cardboard Boxes"),
					makeLeaf("woven-polypropylene-pp-sacks-100kg", "High-Strength Breathable Grain, Fertilizer & Sugar PP Woven Sacks", "የእህል እና የስኳር ጠንካራ የፕላስቲክ ጆንያዎች (100 ኪሎ)", "Qodaa Fastaalaa Midhaanii fi Sukkaaraa", "ናይ High-Strength Breathable Grain, Fertilizer & Sugar PP Woven Sacks"),
					makeLeaf("air-bubble-cushion-wrap-protective-rolls", "Shock Absorbing Fragile Glass & Electronics Bubble Wrap Rolls", "የሚሰበሩ እቃዎችን መጠቅለያ አየር የተሞላባቸው ፕላስቲኮች (Bubble Wrap)", "Pilaastikii Qilleensaa Meeshaa Eegu (Bubble Wrap)", "ናይ Shock Absorbing Fragile Glass & Electronics Bubble Wrap Rolls"),
					makeLeaf("heavy-pallet-stretch-film-rolls", "23-Micron Hand & Machine Industrial Clear Pallet Stretch Wrap", "እቃዎችን በፓሌት ላይ አጥብቆ ማሰሪያ ስትሬች ፊልሞች (Stretch Wrap)", "Fiilmii Saanduqa Hidhu (Stretch Film)", "ናይ 23-Micron Hand & Machine Industrial Clear Pallet Stretch Wrap"),
					makeLeaf("pet-plastic-pallet-strapping-bands", "Heavy Duty High-Tension Green Polyester Cargo Packaging Straps", "የጭነት እቃዎችን ማሰሪያ ጠንካራ የፕላስቲክ ቴፖች (Strapping Band)", "Qabattoo Pilaastikii Fe'iinsaa", "ናይ Heavy Duty High-Tension Green Polyester Cargo Packaging Straps"),
					makeLeaf("sealed-air-dunnage-cargo-airbags", "Heavy Kraft Paper Truck & Container Inflatable Void Fill Airbags", "ኮንቴነር ውስጥ እቃዎች እንዳይወዛወዙ የሚነፉ የደህንነት ኤርባጎች", "Eerbaagii Qorqorroo Fe'iinsaa", "ናይ Heavy Kraft Paper Truck & Container Inflatable Void Fill Airbags"),
					makeLeaf("heavy-duty-wooden-shipping-pallets", "Standard Euro & Four-Way Heat-Treated Forklift Wooden Pallets", "የፎርክሊፍት መጫኛ ጠንካራ የእንጨት ፓሌቶች (Pallet)", "Paaleetii Mukaa Forklift", "ናይ Standard Euro & Four-Way Heat-Treated Forklift Wooden Pallets"),
					makeLeaf("sealed-plastic-jerrycans-20l-chemical", "UN Certified High-Density HDPE Industrial Chemical & Oil Jerrycans", "ባለ 20 ሊትር ጠንካራ የኬሚካል እና ዘይት ጀሪካኖች", "Jerikaanii Pilaastikii Jabaa 20L", "ናይ UN Certified High-Density HDPE Industrial Chemical & Oil Jerrycans"),
				},
			},
			{
				Slug:   "warehouse-material-handling-storage",
				NameEN: "Forklifts, Hand Pallet Trucks and Warehouse Racking",
				NameAM: "ፎርክሊፍቶች፣ የእጅ ፓሌት ጃኮች እና የመጋዘን ሼልፎች",
				NameOM: "Forklift, Jaakii Paaleetii fi Shalfii",
				NameTI: "ፎርክሊፍትን ናይ ኢድ ፓሌት ጃክን ሼልፍን",
				Items: []LeafItemDef{
					makeLeaf("hydraulic-manual-hand-pallet-trucks", "2.5-Ton to 3-Ton Nylon/Polyurethane Wheel Hydraulic Pallet Jacks", "ከ2.5 እስከ 3 ቶን እቃ የሚያነሱ የእጅ ፓሌት ጃኮች", "Jaakii Paaleetii Harkaa 3-Ton", "ናይ 2.5-Ton to 3-Ton Nylon/Polyurethane Wheel Hydraulic Pallet Jacks"),
					makeLeaf("electric-reach-trucks-warehouse-stackers", "1.5-Ton 6-Meter High Narrow Aisle Electric Pallet Reach Trucks", "በመጋዘን ውስጥ ፎቅ ላይ እቃ የሚያወጡ የኤሌክትሪክ ስታከሮች", "Istaakara Elektiriikii Kuusaa", "ናይ 1.5-Ton 6-Meter High Narrow Aisle Electric Pallet Reach Trucks"),
					makeLeaf("heavy-duty-selective-pallet-racking-systems", "Adjustable Multi-Level Steel Heavy Pallet Beam Storage Racks", "ትላልቅ ፓሌቶችን በረድፍ ማስቀመጫ የብረት ሼልፎች", "Shalfii Sibilaa Paaleetii Kuusaa", "ናይ Adjustable Multi-Level Steel Heavy Pallet Beam Storage Racks"),
					makeLeaf("hand-push-heavy-folding-platform-trolleys", "300kg Capacity Mute Wheel Steel Folding Handle Flatbed Carts", "በእጅ የሚገፉ ባለ አራት ጎማ ጠፍጣፋ የጭነት ጋሪዎች", "Gaarii Harkaa Dachaafamu 300kg", "ናይ 300kg Capacity Mute Wheel Steel Folding Handle Flatbed Carts"),
					makeLeaf("dock-leveler-hydraulic-loading-ramps", "6-Ton to 10-Ton Stationary & Mobile Warehouse Forklift Loading Ramps", "ኮንቴነር ውስጥ ፎርክሊፍት እንዲገባ የሚያደርጉ ሀይድሮሊክ ራምፖች", "Raampii Fe'iinsa Forklift", "ናይ 6-Ton to 10-Ton Stationary & Mobile Warehouse Forklift Loading Ramps"),
					makeLeaf("heavy-duty-roller-gravity-conveyor-lines", "Extendable Flexible Steel Unloading Warehouse Gravity Roller Tracks", "እቃዎችን በቀላሉ ከጭነት መኪና የሚያንሸራትቱ የሮለር መስመሮች", "Moorsoo Fe'iinsa Saanduqaa (Conveyor)", "ናይ Extendable Flexible Steel Unloading Warehouse Gravity Roller Tracks"),
					makeLeaf("industrial-drum-lifter-tilter-trolleys", "Hydraulic 200L Steel Oil Drum Transport & Pouring Dispensers", "የነዳጅ እና የኬሚካል በርሜል ማንሻ እና ማዘንበያ ጋሪዎች", "Meeshaa Barmeelii Ol Kaasu", "ናይ Hydraulic 200L Steel Oil Drum Transport & Pouring Dispensers"),
					makeLeaf("warehouse-digital-crane-hanging-scales", "1-Ton to 5-Ton Wireless Remote Industrial Overhead Crane Scales", "በክሬን የተሰቀሉ ከባድ እቃዎችን የሚመዝኑ ዲጂታል ሚዛኖች", "Mizana Kireenii Rimoottiidhaan", "ናይ 1-Ton to 5-Ton Wireless Remote Industrial Overhead Crane Scales"),
				},
			},
		},
	})

	// Domain: general-classifieds
	list = append(list, DomainCatalogDef{
		RootSlug: "general-classifieds",
		Subcategories: []SubcatCatalogDef{
			{
				Slug:   "charity-donations-humanitarian-aid",
				NameEN: "Charity Donations, Community Aid and Mutual Assistance",
				NameAM: "የበጎ አድራጎት ልገሳዎች፣ የእርዳታ እቃዎች እና የማህበረሰብ ድጋፎች",
				NameOM: "Kennaa Tola-Ooltummaa fi Gargaarsa",
				NameTI: "ናይ በጎ ኣድራጎት ልገሳን ናይ ማሕበረሰብ ሓገዝን",
				Items: []LeafItemDef{
					makeLeaf("bulk-used-clothing-shoe-charity-donations", "Sorted Clean Wearable Clothes & Shoes Bundles for Displaced Families", "ለተፈናቃዮች እና ለተቸገሩ የሚለገሱ ንጹህ ልብሶች እና ጫማዎች", "Uffata fi Kophee Kennaa Tola-Ooltummaa", "ናይ Sorted Clean Wearable Clothes & Shoes Bundles for Displaced Families"),
					makeLeaf("school-supplies-backpack-stationery-donations", "Notebooks, Pens, Rulers & School Bags Bundles for Needy Students", "ለአቅመ ደካማ ተማሪዎች የሚለገሱ የደብተር እና የቦርሳ ስብስቦች", "Daftara fi Meeshaalee Barnootaa Kennaaf", "ናይ Notebooks, Pens, Rulers & School Bags Bundles for Needy Students"),
					makeLeaf("wheelchair-crutches-medical-aid-donations", "Functional Mobility Wheelchairs & Orthopedic Braces for Disability Aid", "ለአካል ጉዳተኞች የሚለገሱ ዊልቸሮች እና ምርኩዞች", "Barcuma Naafaa fi Ulee Kennaaf", "ናይ Functional Mobility Wheelchairs & Orthopedic Braces for Disability Aid"),
					makeLeaf("emergency-food-ration-flour-oil-relief", "Wheat Flour, Edible Oil & Split Pea Emergency Disaster Relief Packs", "የአስቸኳይ ጊዜ የምግብ እርዳታ እሽጎች (ዱቄት፣ ዘይት፣ እህል)", "Nyaata Gargaarsa Balaa fi Yeroo Cimaa", "ናይ Wheat Flour, Edible Oil & Split Pea Emergency Disaster Relief Packs"),
					makeLeaf("children-books-educational-library-donations", "Children's Illustrated Storybooks & Encyclopedias for Rural Libraries", "ለገጠር ትምህርት ቤቶች የሚለገሱ የህፃናት መጻሕፍት", "Kitaabolee Daa'immanii Kennaaf", "ናይ Children's Illustrated Storybooks & Encyclopedias for Rural Libraries"),
					makeLeaf("infant-baby-formula-diaper-care-packs", "Clean Baby Diapers, Nutrition Milk Powder & Soap for Orphanages", "ለህፃናት ማሳደጊያዎች የሚሰጡ ዳይፐሮች እና የህፃናት ወተቶች", "Daayparii fi Aanan Daa'immanii Kennaaf", "ናይ Clean Baby Diapers, Nutrition Milk Powder & Soap for Orphanages"),
					makeLeaf("community-tree-seedling-greening-donations", "Indigenous Tree Seedlings for Green Legacy Community Planting", "ለአረንጓዴ አሻራ የሚለገሱ የሀገር በቀል ዛፎች ችግኞች", "Biqiltuu Muka Aadaa Ashaaraa Magariisaaf", "ናይ Indigenous Tree Seedlings for Green Legacy Community Planting"),
					makeLeaf("clean-drinking-water-truck-aid-services", "Water Tanker Truck Deployments for Drought Affected Communities", "በድርቅ ለተጎዱ ወገኖች የውሃ ቦቴዎችን የማሰማራት ድጋፍ", "Bootee Bishaan Dhugaatii Gargaarsaaf", "ናይ Water Tanker Truck Deployments for Drought Affected Communities"),
				},
			},
			{
				Slug:   "barter-exchanges-lost-and-found",
				NameEN: "Barter Exchanges, Swaps and Lost & Found Notices",
				NameAM: "የእቃ በዕቃ ልውውጦች፣ ቅያሪ እና የጠፉ/የተገኙ እቃዎች",
				NameOM: "Jijjiirraa Meeshaa fi Meeshaa Bade/Argame",
				NameTI: "ናይ ንዋይ ብንዋይ ምልውዋጥን ጥፉእ/ዝተረኽበ ሓበሬታን",
				Items: []LeafItemDef{
					makeLeaf("smartphone-laptop-direct-barter-exchanges", "Equal Value Electronics Swap: Smartphone for Laptop Exchanges", "ስልክ በላፕቶፕ ወይም ላፕቶፕ በስልክ የእቃ በዕቃ ልውውጥ", "Jijjiirraa Bilbilaa fi Laptooppii Meeshaan", "ናይ Equal Value Electronics Swap: Smartphone for Laptop Exchanges"),
					makeLeaf("vehicle-upgrade-trade-in-plus-cash", "Trade-In Older Model Toyota Vehicle Plus Cash for Newer Car", "ነባር መኪና ሰጥቶ ተጨማሪ ከፍሎ አዲስ መኪና መለወጥ (Trade-in)", "Konkolaataa Duraa Kennanii Haaraa Fuchuu", "ናይ Trade-In Older Model Toyota Vehicle Plus Cash for Newer Car"),
					makeLeaf("residential-land-house-property-swap", "Exchange Regional House or Plot for Addis Ababa Real Estate Property", "የክልል ቤት ሰጥቶ በአዲስ አበባ ቤት ወይም መሬት መለዋወጥ", "Jijjiirraa Manaa fi Lafa Magaalota Gidduu", "ናይ Exchange Regional House or Plot for Addis Ababa Real Estate Property"),
					makeLeaf("lost-official-documents-identification-notices", "Lost Passports, Title Deeds, Kebele IDs & Driving License Alerts", "የጠፉ ፓስፖርቶች፣ መንጃ ፈቃዶች፣ ካርታዎች እና መታወቂያዎች ማስታወቂያ", "Beeksisa Waraqaa Eenyummaa fi Paaspoortii Bade", "ናይ Lost Passports, Title Deeds, Kebele IDs & Driving License Alerts"),
					makeLeaf("found-keys-wallets-phones-owner-search", "Safe Custody Found House Keys, Pocket Wallets & Phone Owner Search", "የተገኙ የቤት ቁልፎች፣ የኪስ ቦርሳዎች እና ስልኮች ፍለጋ", "Qabduu fi Boorsaa Kiisii Argame", "ናይ Safe Custody Found House Keys, Pocket Wallets & Phone Owner Search"),
					makeLeaf("livestock-for-grain-barter-exchanges", "Equated Value Exchange: Fattened Cattle for Cereal Grain Quintals", "ከብት ሰጥቶ እህል መውሰድ (ባህላዊ የእቃ በዕቃ ልውውጥ)", "Jijjiirraa Horii fi Midhaanii Aadaa", "ናይ Equated Value Exchange: Fattened Cattle for Cereal Grain Quintals"),
					makeLeaf("musical-instrument-guitar-keyboard-trades", "Trade Lead Electric Guitar for Pro Arranger Keyboard Synthesizers", "የሙዚቃ መሳሪያዎችን እርስ በእርስ መለዋወጥ", "Jijjiirraa Meeshaa Muuziqaa (Gitaara fi Kiiboordii)", "ናይ Trade Lead Electric Guitar for Pro Arranger Keyboard Synthesizers"),
					makeLeaf("heavy-machinery-equipment-idle-swap", "Exchange Concrete Mixer for Agricultural Power Tiller Farm Machinery", "የማይሰሩ ወይም የተፈቱ ማሽኖችን በሌላ ማሽን መለዋወጥ", "Jijjiirraa Maashinii Gamoo fi Qonnaa", "ናይ Exchange Concrete Mixer for Agricultural Power Tiller Farm Machinery"),
				},
			},
		},
	})

	return list
}
