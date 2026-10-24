package main

// GetExpandedDomainsData supplies deep branches for remaining domains
func GetExpandedDomainsData() []TaxonomyDomain {
	return []TaxonomyDomain{
		// =====================================================================
		// DOMAIN 16: Books, Literature and Publishing
		// =====================================================================
		{
			Slug: "books-publishing",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "ethiopian-literature-novels", NameEN: "Ethiopian Literature and Amharic Novels", NameAM: "የኢትዮጵያ ስነ-ጽሁፍ እና ልብ ወለዶች", NameOM: "Asoosama fi Ogbarruu", NameTI: "ናይ ኢትዮጵያ ስነ ጽሑፍን ልብ ወለድን",
					Leaves: []TaxonomyLeaf{
						{Slug: "classic-amharic-novels", NameEN: "Classic & Contemporary Amharic Fiction", NameAM: "የአማርኛ ልብወለድ መጻሕፍት", NameOM: "Kitaaba Asoosama Afaan Amaaraa", NameTI: "ናይ ኣምሓርኛ ልብ ወለድ መጻሕፍቲ"},
						{Slug: "oromo-literature-books", NameEN: "Afaan Oromo Literature & Poetry", NameAM: "የኦሮምኛ ስነ-ጽሁፍ መጻሕፍት", NameOM: "Ogbarruu fi Walaloo Afaan Oromoo", NameTI: "ናይ ኦሮምኛ ስነ ጽሑፍ"},
						{Slug: "tigrinya-literature-history", NameEN: "Tigrinya Novels and Historical Texts", NameAM: "የትግርኛ ልብ ወለዶች እና ታሪኮች", NameOM: "Kitaabota Seenaa fi Asoosama Tigriffaa", NameTI: "ናይ ትግርኛ ልብ ወለድን ታሪክን"},
					},
				},
				{
					Slug: "academic-religious-texts", NameEN: "Academic Textbooks and Spiritual Literature", NameAM: "የትምህርት እና የሃይማኖት መጻሕፍት", NameOM: "Kitaabota Amantii fi Barnootaa", NameTI: "ናይ ትምህርትን ሃይማኖትን መጻሕፍቲ",
					Leaves: []TaxonomyLeaf{
						{Slug: "orthodox-tewahedo-books", NameEN: "Ethiopian Orthodox Tewahedo Spiritual Books & Bibles", NameAM: "የኦርቶዶክስ ተዋሕዶ መንፈሳዊ መጻሕፍት", NameOM: "Kitaaba Qulqulluu fi Amantii Ortodooksii", NameTI: "ናይ ኦርቶዶክስ ተዋህዶ መንፈሳዊ መጻሕፍቲ"},
						{Slug: "islamic-quran-literature", NameEN: "Holy Quran, Hadith & Islamic Jurisprudence", NameAM: "ቅዱስ ቁርኣን እና ኢስላማዊ መጻሕፍት", NameOM: "Qur'aana Qulqulluu fi Kitaabota Islaamaa", NameTI: "ቅዱስ ቁርኣንን እስላማዊ መጻሕፍትን"},
						{Slug: "evangelical-spiritual-literature", NameEN: "Protestant & Evangelical Christian Books", NameAM: "የወንጌላውያን ክርስቲያን መጻሕፍት", NameOM: "Kitaabota Kirstaanaa Wangeelaa", NameTI: "ናይ ወንጌላውያን ክርስትያን መጻሕፍቲ"},
						{Slug: "university-academic-textbooks", NameEN: "Higher Education Medical, Engineering & Business Textbooks", NameAM: "የዩኒቨርሲቲ እና የኮሌጅ መማሪያ መጻሕፍት", NameOM: "Kitaabota Yuunivarsiitii", NameTI: "ናይ ዩኒቨርሲቲ መማህሪ መጻሕፍቲ"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 17: Education, Courses and Training
		// =====================================================================
		{
			Slug: "education-training",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "professional-skills-training", NameEN: "Professional Certification & Software Bootcamps", NameAM: "የሙያ ማረጋገጫ እና የስልጠና ኮርሶች", NameOM: "Leenjii Ogummaa fi Teeknooloojii", NameTI: "ናይ ሞያ ምስክር ወረቐትን ስልጠናን",
					Leaves: []TaxonomyLeaf{
						{Slug: "software-coding-bootcamps", NameEN: "Web Development, Python & Full Stack Bootcamps", NameAM: "የኮዲንግ እና የሶፍትዌር ስልጠናዎች", NameOM: "Leenjii Barreeffama Koodii (Coding)", NameTI: "ናይ ሶፍትዌር ስልጠና (Coding)"},
						{Slug: "accounting-peachtree-quickbooks", NameEN: "Practical Accounting, Peachtree & IFRS Training", NameAM: "የፒችትሪ እና የሂሳብ አያያዝ ስልጠና", NameOM: "Leenjii Herregaa fi Piiyichtirii", NameTI: "ናይ ፒችትሪን ሒሳብን ስልጠና"},
						{Slug: "graphic-design-video-editing-classes", NameEN: "Adobe Photoshop, Premiere & Graphic Design Classes", NameAM: "የግራፊክ ዲዛይን እና የቪዲዮ ኤዲቲንግ ስልጠና", NameOM: "Barnoota Giraafiiksii fi Viidiyoo", NameTI: "ናይ ግራፊክ ዲዛይን ስልጠና"},
					},
				},
				{
					Slug: "vocational-tvet-driving", NameEN: "Vocational TVET and Driving Lessons", NameAM: "የቴክኒክና ሙያ እንዲሁም የመንጃ ፈቃድ ስልጠና", NameOM: "Leenjii Konkolaachisummaa fi TVET", NameTI: "ናይ ቴክኒክን ሞያን ናይ መኪና መዝወርን ስልጠና",
					Leaves: []TaxonomyLeaf{
						{Slug: "commercial-driving-license-training", NameEN: "Auto & Heavy Truck Driving License Schools (Grade 2-5)", NameAM: "የመንጃ ፈቃድ ትምህርት ቤቶች", NameOM: "Mana Barumsa Konkolaachisummaa", NameTI: "ቤት ትምህርቲ መዝወሪ መኪና"},
						{Slug: "tvet-electrical-plumbing-trades", NameEN: "Hands-on Electrical & Plumbing Trade Workshops", NameAM: "የተግባር የኤሌክትሪክ እና የቧንቧ ስልጠና", NameOM: "Shaakala TVET Elektiriikii fi Ujummoo", NameTI: "ናይ ኤሌክትሪክን ሻምብቆን ተግባራዊ ስልጠና"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 22: Software, Digital Products and IT Services
		// =====================================================================
		{
			Slug: "software-it-services",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "custom-software-engineering", NameEN: "Custom Web, Mobile & Enterprise ERP Development", NameAM: "የድር ጣቢያ እና የሞባይል አፕሊኬሽን ልማት", NameOM: "Oomisha Weebsaayitii fi Appilikeeshinii", NameTI: "ናይ ዌብሳይትን ሞባይል ኣፕሊኬሽንን ምዕባለ",
					Leaves: []TaxonomyLeaf{
						{Slug: "enterprise-erp-crm-development", NameEN: "Custom Enterprise Resource Planning (ERP) Systems", NameAM: "የድርጅት ማኔጅመንት ሲስተሞች (ERP)", NameOM: "Siistamii Hooggansa Dhaabbataa (ERP)", NameTI: "ናይ ትካል መመሓደሪ ሲስተም (ERP)"},
						{Slug: "mobile-ios-android-apps", NameEN: "Native iOS & Android Mobile App Development", NameAM: "የአይኦኤስ እና አንድሮይድ አፕሊኬሽን ስራ", NameOM: "Oomisha Appilikeeshinii Bilbilaa", NameTI: "ናይ ሞባይል ኣፕሊኬሽን ስራሕ"},
						{Slug: "e-commerce-portal-solutions", NameEN: "E-Commerce Portals & Chapa / Telebirr Integrations", NameAM: "የኦንላይን መገበያያ ድረ ገጾች ከቴሌብር ክፍያ ጋር", NameOM: "Weebsaayitii Daldala Intarneetii", NameTI: "ናይ ኦንላይን ዕዳጋ ዌብሳይት"},
					},
				},
				{
					Slug: "cloud-cybersecurity-services", NameEN: "Cloud Infrastructure, Domains and Hosting", NameAM: "ክላውድ ሆስቲንግ እና የደህንነት ኦዲት", NameOM: "Hoostiingii fi Nageenya Kilaawudii", NameTI: "ናይ ክላውድ ሆስቲንግን ድሕንነትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "vps-cloud-hosting-servers", NameEN: "VPS Cloud Hosting & Linux Server Administration", NameAM: "የቪፒኤስ ክላውድ ሰርቨር አገልግሎት", NameOM: "Tajaajila Sarvarii Kilaawudii VPS", NameTI: "ናይ VPS ክላውድ ሰርቨር ኣገልግሎት"},
						{Slug: "domain-registration-ssl", NameEN: ".et and International Domain Name Registration", NameAM: "የዶሜይን ስም ምዝገባ (.et) እና SSL", NameOM: "Galmee Maqaa Doomeeyinii (.et)", NameTI: "ናይ ዶሜይን ምዝገባ (.et)"},
						{Slug: "cybersecurity-vulnerability-auditing", NameEN: "Penetration Testing & Cybersecurity Auditing", NameAM: "የሲስተም ደህንነት ፍተሻ እና ኦዲት", NameOM: "Qorannoo Nageenya Saayibarii", NameTI: "ናይ ሳይበር ድሕንነት ፍተሻ"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 26: Business Opportunities and Franchises
		// =====================================================================
		{
			Slug: "business-opportunities",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "operational-businesses-for-sale", NameEN: "Operational Businesses & Retail Stores for Sale", NameAM: "ተሰርተው ያለቁ የሚሸጡ ንግዶች እና ሱቆች", NameOM: "Daldalawwan Hojjechaa Jiran Gurguraa", NameTI: "ዝሰርሑ ዘለዉ ዝሽየጡ ንግድታትን ድኳናትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "supermarkets-grocery-stores-sale", NameEN: "Turnkey Supermarkets & Minimarts for Transfer", NameAM: "የሚሸጡ ሱፐርማርኬቶች እና ሚኒማርቶች", NameOM: "Suupparmaarkeetii Gurguraa", NameTI: "ዝሽየጡ ሱፐርማርኬታት"},
						{Slug: "cafes-restaurants-bars-sale", NameEN: "Established Cafes, Lounges and Restaurants for Sale", NameAM: "የሚሸጡ ካፌዎች እና ሬስቶራንቶች", NameOM: "Kaafeefi Resteeraantii Gurguraa", NameTI: "ዝሽየጡ ካፌታትን ሬስቶራንታትን"},
						{Slug: "licensed-pharmacies-clinics-sale", NameEN: "Licensed Retail Pharmacies and Specialty Clinics for Transfer", NameAM: "የሚሸጡ ፈቃድ ያላቸው ፋርማሲዎች", NameOM: "Faarmasii Eeyyama Qabu Gurguraa", NameTI: "ፍቓድ ዘለዎም ዝሽየጡ ፋርማሲታት"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 50: Art, Crafts, Collectibles and Antiques
		// =====================================================================
		{
			Slug: "art-crafts-antiques",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "traditional-handicrafts", NameEN: "Handwoven Mesob, Baskets and Pottery", NameAM: "መሶቦች፣ ቅርጫቶች እና የሸክላ ስራዎች", NameOM: "Masobii fi Meeshaa Suphee", NameTI: "መሶብን ናይ ሸኽላ ኣቑሑትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "handwoven-colorful-mesob", NameEN: "Large Handwoven Traditional Dining Mesob", NameAM: "ባህላዊ የገበታ መሶቦች", NameOM: "Masobii Nyaataa Aadaa", NameTI: "ባህላዊ ናይ መኣዲ መሶብ"},
						{Slug: "clay-coffee-pots-jebena", NameEN: "Handcrafted Clay Coffee Brewing Pots (Jebena)", NameAM: "የሸክላ የቡና ጀበናዎች", NameOM: "Jabaynaa Suphee Bunaa", NameTI: "ናይ ሸኽላ ናይ ቡን ጀበና"},
						{Slug: "hand-carved-wooden-stools-crosses", NameEN: "Traditional Wooden Jimma Chairs & Carved Crosses", NameAM: "የጅማ ባህላዊ በርጩማዎች እና ቅርጾች", NameOM: "Barcuma Aadaa Jimmaa", NameTI: "ናይ ጅማ ባህላዊ በርጩማ"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 51: Music, Musical Instruments and Entertainment
		// =====================================================================
		{
			Slug: "music-instruments",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "ethiopian-traditional-instruments", NameEN: "Kirar, Masinko, Kebero and Washint", NameAM: "ክራር፣ ማሲንቆ፣ ከበሮ እና ዋሽንት", NameOM: "Kiraara, Maasiinqoo fi Kabaroo", NameTI: "ክራር፣ ማሲንቆ፣ ከበሮን ዋሽንትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "acoustic-electric-kirar", NameEN: "6-String Acoustic & Electric Concert Kirar", NameAM: "አኮስቲክ እና ኤሌክትሪክ ክራር", NameOM: "Kiraara Elektiriikii fi Aadaa", NameTI: "ኣኮስቲክን ኤሌክትሪክን ክራር"},
						{Slug: "single-string-horsehair-masinko", NameEN: "Traditional Handcrafted Horsehair Masinko", NameAM: "ባህላዊ ማሲንቆ", NameOM: "Maasiinqoo Aadaa", NameTI: "ባህላዊ ማሲንቆ"},
						{Slug: "large-church-concert-kebero", NameEN: "Rawhide Covered Church & Cultural Kebero Drums", NameAM: "የቆዳ የከበሮ ከበሮዎች", NameOM: "Kabaroo Gogaa Aadaa", NameTI: "ናይ ቆርበት ከበሮ"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 53: Toys, Baby Products and Children's Goods
		// =====================================================================
		{
			Slug: "toys-baby-products",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "baby-gear-strollers", NameEN: "Baby Strollers, Car Seats and Cribs", NameAM: "የህፃናት ጋሪዎች፣ አልጋዎች እና መኪና ወንበሮች", NameOM: "Gaarii Daa'immanii fi Siree", NameTI: "ናይ ህጻናት ጋሪታትን ዓራትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "foldable-lightweight-baby-strollers", NameEN: "Foldable Lightweight Airplane Baby Strollers", NameAM: "ተጣጣፊ የህፃናት ጋሪዎች", NameOM: "Gaarii Daa'immanii Dachaafamu", NameTI: "ዝዕጸፍ ናይ ህጻናት ጋሪ"},
						{Slug: "infant-convertible-car-seats", NameEN: "ISOFIX Infant Convertible Safety Car Seats", NameAM: "የህፃናት የመኪና ውስጥ ደህንነት ወንበሮች", NameOM: "Barcuma Eegumsaa Konkolaataa Daa'immanii", NameTI: "ናይ ህጻናት መኪና መንበር"},
						{Slug: "wooden-baby-cribs-bassinet", NameEN: "Solid Wood Baby Cribs with Rocking Bassinets", NameAM: "የእንጨት የህፃናት አልጋዎች (ክሪብ)", NameOM: "Siree Mukaa Daa'immanii", NameTI: "ናይ ዕንጨይቲ ናይ ህጻናት ዓራት"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 54: Pet Supplies and Companion Animals
		// =====================================================================
		{
			Slug: "pet-supplies",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "companion-dogs-cats", NameEN: "Purebred Dogs, Cats and Puppies", NameAM: "የቤት ውሾች፣ ድመቶች እና ቡችሎች", NameOM: "Saree fi Adurree Manaa", NameTI: "ናይ ገዛ ኣኽላባትን ድሙታትን",
					Leaves: []TaxonomyLeaf{
						{Slug: "german-shepherd-guard-puppies", NameEN: "German Shepherd Working Line Guard Puppies", NameAM: "ጀርመን ሼፐርድ የጥበቃ ቡችሎች", NameOM: "Cuucii Saree Jarman Sheppardi", NameTI: "ጀርመን ሸፐርድ ናይ ሓለዋ ቡችላታት"},
						{Slug: "rottweiler-boerboel-puppies", NameEN: "Rottweiler & Boerboel Heavy Guard Dogs", NameAM: "ሮትዋይለር እና ቦርቦል የጥበቃ ውሾች", NameOM: "Saree Eegumsaa Jabaa", NameTI: "ሮትዋይለር ናይ ሓለዋ ኣኽላባት"},
						{Slug: "premium-dry-dog-food-kibble", NameEN: "High-Protein Beef & Chicken Dry Dog Food (15kg-20kg)", NameAM: "የታሸገ የውሻ ምግብ (ኪብል)", NameOM: "Nyaata Saree Qophaa'e", NameTI: "ዝተዳለወ ናይ ከልቢ መግቢ"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 58: Travel, Accommodation and Tourism Services
		// =====================================================================
		{
			Slug: "travel-tourism",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "ethiopian-tour-packages", NameEN: "Historic Route & Trekking Tour Packages", NameAM: "የሀገር ውስጥ የቱሪዝም እና የጉዞ ፓኬጆች", NameOM: "Tajaajila Imala fi Tuurizimii", NameTI: "ናይ ውሽጢ ዓዲ ቱሪዝምን መገሻን",
					Leaves: []TaxonomyLeaf{
						{Slug: "lalibela-gondar-historic-tours", NameEN: "Lalibela Rock-Hewn Churches & Gondar Castles Tours", NameAM: "የላሊበላ እና የጎንደር ታሪካዊ ጉብኝት", NameOM: "Daawwanna Seenaa Laallibalaa fi Gondar", NameTI: "ናይ ላሊበላን ጎንደርን ታሪካዊ ዑደት"},
						{Slug: "simien-bale-mountains-trekking", NameEN: "Simien Mountains & Bale National Park Trekking Expeditions", NameAM: "የሰሜን እና የባሌ ተራሮች የእግር ጉዞ", NameOM: "Imala Gaara Semeen fi Baalee", NameTI: "ናይ ስሜን ጎቦታት መገሻ"},
						{Slug: "danakil-depression-erta-ale-tours", NameEN: "Danakil Depression, Dallol & Erta Ale Volcano Expeditions", NameAM: "የዳናኪል በረሃ እና የኤርታሌ እሳተ ጎመራ ጉዞ", NameOM: "Daawwanna Gammoojjii Danakil", NameTI: "ናይ ዳናኪልን ኤርታለን ዑደት"},
					},
				},
			},
		},

		// =====================================================================
		// DOMAIN 59: Packaging, Shipping and Warehousing
		// =====================================================================
		{
			Slug: "packaging-shipping",
			Subdomains: []TaxonomySubdomain{
				{
					Slug: "industrial-bulk-packaging", NameEN: "Grain Bags, Cartons and Industrial Packaging", NameAM: "የእህል ማዳበሪያ ከረጢቶች እና ካርቶኖች", NameOM: "Qodaa Saamsuu fi Korojoo", NameTI: "ናይ እኽሊ መዐሸጊ ከረጢትን ካርቶንን",
					Leaves: []TaxonomyLeaf{
						{Slug: "woven-polypropylene-grain-bags", NameEN: "50kg & 100kg Woven Polypropylene (PP) Grain Sacks", NameAM: "የ50 ኪ.ግ እና የ100 ኪ.ግ የእህል ማዳበሪያ ከረጢቶች", NameOM: "Korojoo Midhaanii PP (50kg/100kg)", NameTI: "ናይ 50/100 ኪ.ግ ናይ እኽሊ ከረጢት"},
						{Slug: "corrugated-master-carton-boxes", NameEN: "3-Ply & 5-Ply Corrugated Shipping Carton Boxes", NameAM: "ባለ 3 እና ባለ 5 ፕላይ የካርቶን ሳጥኖች", NameOM: "Saanduqa Kaartoonii Jabaa", NameTI: "ናይ ካርቶን ሳንዱቕ"},
						{Slug: "industrial-pallet-stretch-film", NameEN: "Heavy Duty Machine & Hand Pallet Stretch Wrap Film", NameAM: "የእቃ ማሸጊያ ፕላስቲክ (ስትሬች ፊልም)", NameOM: "Pilaastikii Meeshaa Saamsu", NameTI: "መዐሸጊ ፕላስቲክ ስትሬች ፊልም"},
					},
				},
			},
		},
	}
}
