package main

// OccupationDef defines a standardized occupation entity
type OccupationDef struct {
	Slug        string
	TitleEN     string
	TitleAM     string
	TitleOM     string
	TitleTI     string
	CategorySlug string
	Code        string
}

// QualificationDef defines a qualification type / education level
type QualificationDef struct {
	NameEN     string
	NameAM     string
	NameOM     string
	NameTI     string
	LevelOrder int
	Category   string // 'degree', 'tvet', 'diploma', 'certificate', 'license'
}

// GetQualificationTypes returns verified Ethiopian and universal qualification levels
func GetQualificationTypes() []QualificationDef {
	return []QualificationDef{
		{NameEN: "Doctor of Philosophy (PhD / Doctorate)", NameAM: "የዶክትሬት ዲግሪ (PhD)", NameOM: "Digirii Doktoreetii (PhD)", NameTI: "ዶክትሬት ዲግሪ (PhD)", LevelOrder: 10, Category: "degree"},
		{NameEN: "Doctor of Medicine (MD)", NameAM: "የህክምና ዶክትሬት (MD)", NameOM: "Doktoora Yaalaa (MD)", NameTI: "ናይ ሕክምና ዶክተር (MD)", LevelOrder: 9, Category: "degree"},
		{NameEN: "Master of Science (MSc)", NameAM: "የሳይንስ ማስተርስ ዲግሪ (MSc)", NameOM: "Digirii Lammaffaa Saayinsii (MSc)", NameTI: "ማስተርስ ዲግሪ ሳይንስ (MSc)", LevelOrder: 8, Category: "degree"},
		{NameEN: "Master of Arts (MA)", NameAM: "የአርትስ ማስተርስ ዲግሪ (MA)", NameOM: "Digirii Lammaffaa Aartii (MA)", NameTI: "ማስተርስ ዲግሪ ኣርትስ (MA)", LevelOrder: 8, Category: "degree"},
		{NameEN: "Master of Business Administration (MBA)", NameAM: "የቢዝነስ አስተዳደር ማስተርስ (MBA)", NameOM: "Maastarsii Daldalaa (MBA)", NameTI: "ማስተርስ ዲግሪ ቢዝነስ (MBA)", LevelOrder: 8, Category: "degree"},
		{NameEN: "Master of Laws (LLM)", NameAM: "የህግ ማስተርስ ዲግሪ (LLM)", NameOM: "Maastarsii Seeraa (LLM)", NameTI: "ማስተርስ ዲግሪ ሕጊ (LLM)", LevelOrder: 8, Category: "degree"},
		{NameEN: "Postgraduate Diploma (PGD)", NameAM: "የድህረ ምረቃ ዲፕሎማ (PGD)", NameOM: "Dippiloomaa Eebba Boodaa (PGD)", NameTI: "ድሕረ ምረቓ ዲፕሎማ (PGD)", LevelOrder: 7, Category: "diploma"},
		{NameEN: "Bachelor of Science (BSc)", NameAM: "የመጀመሪያ ዲግሪ ሳይንስ (BSc)", NameOM: "Digirii Jalqabaa Saayinsii (BSc)", NameTI: "ቀዳማይ ዲግሪ ሳይንስ (BSc)", LevelOrder: 6, Category: "degree"},
		{NameEN: "Bachelor of Arts (BA)", NameAM: "የመጀመሪያ ዲግሪ አርትስ (BA)", NameOM: "Digirii Jalqabaa Aartii (BA)", NameTI: "ቀዳማይ ዲግሪ ኣርትስ (BA)", LevelOrder: 6, Category: "degree"},
		{NameEN: "Bachelor of Laws (LLB)", NameAM: "የህግ የመጀመሪያ ዲግሪ (LLB)", NameOM: "Digirii Jalqabaa Seeraa (LLB)", NameTI: "ቀዳማይ ዲግሪ ሕጊ (LLB)", LevelOrder: 6, Category: "degree"},
		{NameEN: "Bachelor of Pharmacy (BPharm)", NameAM: "የፋርማሲ የመጀመሪያ ዲግሪ (BPharm)", NameOM: "Digirii Jalqabaa Faarmasii (BPharm)", NameTI: "ቀዳማይ ዲግሪ ፋርማሲ (BPharm)", LevelOrder: 6, Category: "degree"},
		{NameEN: "Bachelor of Education (BEd)", NameAM: "የትምህርት የመጀመሪያ ዲግሪ (BEd)", NameOM: "Digirii Jalqabaa Barnootaa (BEd)", NameTI: "ቀዳማይ ዲግሪ ትምህርቲ (BEd)", LevelOrder: 6, Category: "degree"},
		{NameEN: "TVET Level V (National Advanced Diploma)", NameAM: "ቴክኒክና ሙያ ደረጃ 5 (ከፍተኛ ዲፕሎማ)", NameOM: "Sadarkaa TVET 5 (Dippiloomaa Olaanaa)", NameTI: "ቴክኒክን ሞያን ደረጃ 5 (ላዕለዋይ ዲፕሎማ)", LevelOrder: 5, Category: "tvet"},
		{NameEN: "TVET Level IV (National Diploma)", NameAM: "ቴክኒክና ሙያ ደረጃ 4 (ዲፕሎማ)", NameOM: "Sadarkaa TVET 4 (Dippiloomaa)", NameTI: "ቴክኒክን ሞያን ደረጃ 4 (ዲፕሎማ)", LevelOrder: 4, Category: "tvet"},
		{NameEN: "TVET Level III (Lead Craftsperson)", NameAM: "ቴክኒክና ሙያ ደረጃ 3", NameOM: "Sadarkaa TVET 3", NameTI: "ቴክኒክን ሞያን ደረጃ 3", LevelOrder: 3, Category: "tvet"},
		{NameEN: "TVET Level II (Skilled Operator)", NameAM: "ቴክኒክና ሙያ ደረጃ 2", NameOM: "Sadarkaa TVET 2", NameTI: "ቴክኒክን ሞያን ደረጃ 2", LevelOrder: 2, Category: "tvet"},
		{NameEN: "TVET Level I (Basic Certificate)", NameAM: "ቴክኒክና ሙያ ደረጃ 1", NameOM: "Sadarkaa TVET 1", NameTI: "ቴክኒክን ሞያን ደረጃ 1", LevelOrder: 1, Category: "tvet"},
		{NameEN: "Professional Certified Public Accountant (CPA / ACCA)", NameAM: "የተመሰከረለት የሂሳብ ባለሙያ (ACCA / CPA)", NameOM: "Ogeessa Herregaa Mirkanaa'e (ACCA)", NameTI: "ዝተመስከረሉ በዓል ሞያ ሒሳብ (ACCA)", LevelOrder: 7, Category: "certificate"},
		{NameEN: "Project Management Professional (PMP)", NameAM: "የፕሮጀክት ማኔጅመንት ሰርቲፊኬት (PMP)", NameOM: "Mirkaneessa Hooggansa Pirojektii (PMP)", NameTI: "ምሕደራ ፕሮጀክት ሰርቲፊኬት (PMP)", LevelOrder: 7, Category: "certificate"},
		{NameEN: "Professional Engineering Practicing License", NameAM: "የምህንድስና ስራ ፈቃድ", NameOM: "Eeyyama Hojii Injinariingii", NameTI: "ናይ ምህንድስና ስራሕ ፍቓድ", LevelOrder: 6, Category: "license"},
		{NameEN: "Medical Practicing License (Board Certified)", NameAM: "የህክምና ስራ ፈቃድ", NameOM: "Eeyyama Hojii Yaalaa", NameTI: "ናይ ሕክምና ስራሕ ፍቓድ", LevelOrder: 7, Category: "license"},
		{NameEN: "Commercial Heavy Vehicle Driving License (Grade 4 / Dry Cargo)", NameAM: "የደረቅ ጭነት 4ኛ መንጃ ፈቃድ", NameOM: "Eeyyama Konkolaachisummaa Sadarkaa 4", NameTI: "ናይ ጽዕነት 4ይ ደረጃ መንጃ ፍቓድ", LevelOrder: 3, Category: "license"},
		{NameEN: "Public Transport Bus Driving License (Grade 5)", NameAM: "የህዝብ ማመላለሻ 5ኛ መንጃ ፈቃድ", NameOM: "Eeyyama Konkolaachisummaa Sadarkaa 5", NameTI: "ናይ ህዝቢ መጓዓዝያ 5ይ ደረጃ መንጃ ፍቓድ", LevelOrder: 3, Category: "license"},
		{NameEN: "Automobile Driving License (Grade 2 / Private)", NameAM: "የህዝብ 1 / አውቶሞቢል መንጃ ፈቃድ", NameOM: "Eeyyama Konkolaachisummaa Awutoomoobilii", NameTI: "ናይ ውልቂ መኪና መንጃ ፍቓድ", LevelOrder: 2, Category: "license"},
	}
}

// GetStandardOccupations returns verified occupational roles aligned with Ethiopian and ISCO standards
func GetStandardOccupations() []OccupationDef {
	return []OccupationDef{
		// Technology & Software
		{Slug: "software-engineer", TitleEN: "Software Engineer", TitleAM: "ሶፍትዌር መሃንዲስ", TitleOM: "Injinara Sooftiweerii", TitleTI: "ሶፍትዌር መሃንድስ", CategorySlug: "jobs-recruitment", Code: "2512"},
		{Slug: "full-stack-developer", TitleEN: "Full Stack Developer", TitleAM: "ሙሉ የሶፍትዌር አበልጻጊ", TitleOM: "Oomishaa Sooftiweerii Guutuu", TitleTI: "ሙሉእ ሶፍትዌር ኣማዕባሊ", CategorySlug: "jobs-recruitment", Code: "2512"},
		{Slug: "frontend-developer", TitleEN: "Frontend Developer", TitleAM: "የፊት ለፊት ሶፍትዌር አበልጻጊ", TitleOM: "Oomishaa Fuula Duraa", TitleTI: "ናይ ቅድሚት ሶፍትዌር ኣማዕባሊ", CategorySlug: "jobs-recruitment", Code: "2513"},
		{Slug: "backend-developer", TitleEN: "Backend Developer", TitleAM: "የጀርባ ሲስተም አበልጻጊ", TitleOM: "Oomishaa Siistamii Duubaa", TitleTI: "ናይ ድሕሪት ሲስተም ኣማዕባሊ", CategorySlug: "jobs-recruitment", Code: "2512"},
		{Slug: "mobile-app-developer", TitleEN: "Mobile App Developer", TitleAM: "የሞባይል መተግበሪያ አበልጻጊ", TitleOM: "Oomishaa Appilikeeshinii Bilbilaa", TitleTI: "ናይ ሞባይል ኣፕ ኣማዕባሊ", CategorySlug: "jobs-recruitment", Code: "2514"},
		{Slug: "devops-engineer", TitleEN: "DevOps & Cloud Engineer", TitleAM: "የደመና ቴክኖሎጂ መሃንዲስ", TitleOM: "Injinara Kilaawudii", TitleTI: "ናይ ክላውድ ቴክኖሎጂ መሃንድስ", CategorySlug: "jobs-recruitment", Code: "2529"},
		{Slug: "database-administrator", TitleEN: "Database Administrator (DBA)", TitleAM: "የዳታቤዝ አስተዳዳሪ", TitleOM: "Bulchaa Daataabeesii", TitleTI: "ናይ ዳታቤዝ ኣመሓዳሪ", CategorySlug: "jobs-recruitment", Code: "2521"},
		{Slug: "network-systems-administrator", TitleEN: "Network & Systems Administrator", TitleAM: "የኔትወርክ እና ሲስተም አስተዳዳሪ", TitleOM: "Bulchaa Netwoorkii fi Siistamii", TitleTI: "ናይ መርበብ ሓበሬታ ኣመሓዳሪ", CategorySlug: "jobs-recruitment", Code: "2522"},
		{Slug: "cybersecurity-specialist", TitleEN: "Cybersecurity Specialist", TitleAM: "የሳይበር ደህንነት ባለሙያ", TitleOM: "Ogeessa Nageenya Saayibarii", TitleTI: "ናይ ሳይበር ድሕንነት በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "2529"},
		{Slug: "ui-ux-designer", TitleEN: "UI/UX Product Designer", TitleAM: "የUI/UX ዲዛይነር", TitleOM: "Dizaayinara UI/UX", TitleTI: "UI/UX ዲዛይነር", CategorySlug: "jobs-recruitment", Code: "2166"},
		{Slug: "qa-software-tester", TitleEN: "QA & Software Test Engineer", TitleAM: "የሶፍትዌር ጥራት ተቆጣጣሪ", TitleOM: "Injinara Qulqullina Sooftiweerii", TitleTI: "ናይ ሶፍትዌር ፅሬት መርማሪ", CategorySlug: "jobs-recruitment", Code: "2519"},
		{Slug: "data-scientist-analyst", TitleEN: "Data Analyst & Data Scientist", TitleAM: "የዳታ ተንታኝ እና ሳይንቲስት", TitleOM: "Qaaccessaa Daataa", TitleTI: "ናይ ዳታ ተንታኒ", CategorySlug: "jobs-recruitment", Code: "2523"},

		// Engineering & Construction
		{Slug: "civil-engineer", TitleEN: "Civil Engineer", TitleAM: "ሲቪል መሃንዲስ", TitleOM: "Injinara Siiviilii", TitleTI: "ሲቪል መሃንድስ", CategorySlug: "jobs-recruitment", Code: "2142"},
		{Slug: "structural-engineer", TitleEN: "Structural Engineer", TitleAM: "የመዋቅር መሃንዲስ", TitleOM: "Injinara Dhaabbataa", TitleTI: "ናይ ህንጻ መዋቕር መሃንድስ", CategorySlug: "jobs-recruitment", Code: "2142"},
		{Slug: "architect", TitleEN: "Architect", TitleAM: "አርክቴክት", TitleOM: "Arkiteektii", TitleTI: "ኣርክቴክት", CategorySlug: "jobs-recruitment", Code: "2161"},
		{Slug: "electrical-engineer", TitleEN: "Electrical Engineer", TitleAM: "የኤሌክትሪክ መሃንዲስ", TitleOM: "Injinara Elektiriikaa", TitleTI: "ኤሌክትሪካል መሃንድስ", CategorySlug: "jobs-recruitment", Code: "2151"},
		{Slug: "mechanical-engineer", TitleEN: "Mechanical Engineer", TitleAM: "መካኒካል መሃንዲስ", TitleOM: "Injinara Meekaanikaalaa", TitleTI: "መካኒካል መሃንድስ", CategorySlug: "jobs-recruitment", Code: "2144"},
		{Slug: "quantity-surveyor", TitleEN: "Quantity Surveyor (Estimator)", TitleAM: "የግንባታ ዋጋ እና መጠን ገማች", TitleOM: "Shallagaa Meeshaa Ijaarsaa", TitleTI: "ናይ ህንጻ ዋጋን መጠንን ገማቲ", CategorySlug: "jobs-recruitment", Code: "2149"},
		{Slug: "site-engineer", TitleEN: "Construction Site Engineer", TitleAM: "የሳይት መሃንዲስ", TitleOM: "Injinara Saayitii Ijaarsaa", TitleTI: "ናይ ሳይት መሃንድስ", CategorySlug: "jobs-recruitment", Code: "2142"},
		{Slug: "surveyor-geodesist", TitleEN: "Land Surveyor & Geodesist", TitleAM: "የመሬት ቅየሳ ባለሙያ", TitleOM: "Qorataa Lafa fi Kaartaa", TitleTI: "ናይ መሬት ቅየሳ በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "2165"},
		{Slug: "geotechnical-engineer", TitleEN: "Geotechnical Engineer", TitleAM: "የአፈር እና የጂኦቴክኒክ መሃንዲስ", TitleOM: "Injinara Biyyoo fi Dhagaa", TitleTI: "ናይ ሓመድ መሃንድስ", CategorySlug: "jobs-recruitment", Code: "2142"},
		{Slug: "sanitary-engineer", TitleEN: "Sanitary & Hydraulic Engineer", TitleAM: "የሳኒተሪ እና የውሃ ምህንድስና ባለሙያ", TitleOM: "Injinara Bishaanii fi Dhangala'aa", TitleTI: "ናይ ማይን ሳኒተሪን መሃንድስ", CategorySlug: "jobs-recruitment", Code: "2142"},

		// Healthcare & Medicine
		{Slug: "medical-doctor", TitleEN: "General Medical Practitioner (MD)", TitleAM: "አጠቃላይ የህክምና ሃኪም", TitleOM: "Hakiima Yaala Waliigalaa", TitleTI: "ሓፈሻዊ ሓኪም", CategorySlug: "jobs-recruitment", Code: "2211"},
		{Slug: "nurse-registered", TitleEN: "Registered Clinical Nurse", TitleAM: "ክሊኒካል ነርስ", TitleOM: "Narsii Kiliinikaa", TitleTI: "ክሊኒካል ነርስ", CategorySlug: "jobs-recruitment", Code: "2221"},
		{Slug: "pharmacist", TitleEN: "Pharmacist", TitleAM: "ፋርማሲስት (የመድሃኒት ባለሙያ)", TitleOM: "Ogeessa Qorichaa (Faarmaasistii)", TitleTI: "ፋርማሲስት", CategorySlug: "jobs-recruitment", Code: "2262"},
		{Slug: "medical-laboratory-technologist", TitleEN: "Medical Laboratory Technologist", TitleAM: "የላቦራቶሪ ቴክኖሎጂስት", TitleOM: "Teeknooloojistii Laabiraatorii", TitleTI: "ናይ ላቦራቶሪ በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "2240"},
		{Slug: "dentist", TitleEN: "Dental Surgeon & Dentist", TitleAM: "የጥርስ ሃኪም", TitleOM: "Hakiima Ilkaanii", TitleTI: "ሓኪም ስኒ", CategorySlug: "jobs-recruitment", Code: "2261"},
		{Slug: "public-health-officer", TitleEN: "Public Health Officer", TitleAM: "የህብረተሰብ ጤና መኮንን", TitleOM: "Ogeessa Fayyaa Hawaasaa", TitleTI: "ናይ ማሕበረሰብ ጥዕና በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "2263"},

		// Finance, Accounting & Business
		{Slug: "senior-accountant", TitleEN: "Senior Accountant", TitleAM: "ዋና የሂሳብ ባለሙያ", TitleOM: "Ogeessa Herregaa Olaanaa", TitleTI: "ላዕለዋይ በዓል ሞያ ሒሳብ", CategorySlug: "jobs-recruitment", Code: "2411"},
		{Slug: "financial-analyst", TitleEN: "Financial Analyst", TitleAM: "የፋይናንስ ተንታኝ", TitleOM: "Qaaccessaa Faayinaansii", TitleTI: "ናይ ፋይናንስ ተንታኒ", CategorySlug: "jobs-recruitment", Code: "2413"},
		{Slug: "internal-auditor", TitleEN: "Internal Auditor", TitleAM: "የውስጥ ኦዲተር", TitleOM: "Oodiitara Keessaa", TitleTI: "ናይ ውሽጢ ኦዲተር", CategorySlug: "jobs-recruitment", Code: "2411"},
		{Slug: "procurement-officer", TitleEN: "Procurement & Supply Chain Specialist", TitleAM: "የግዢ እና አቅርቦት ባለሙያ", TitleOM: "Ogeessa Bittaafi Dhiyeessii", TitleTI: "ናይ ዕደጋን ቀረብን በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "3323"},
		{Slug: "human-resources-manager", TitleEN: "Human Resources Manager", TitleAM: "የሰው ሀብት ስራ አመራር", TitleOM: "Hoogganaa Qabeenya Namaa", TitleTI: "ናይ ሰብ ሃፍቲ ኣመሓዳሪ", CategorySlug: "jobs-recruitment", Code: "1212"},
		{Slug: "marketing-manager", TitleEN: "Marketing & Growth Manager", TitleAM: "የገበያ ልማት እና ማርኬቲንግ ስራ አስኪያጅ", TitleOM: "Gaggeessaa Gabaa fi Beeksisoo", TitleTI: "ናይ ዕዳጋን ማርኬቲንግን ኣመሓዳሪ", CategorySlug: "jobs-recruitment", Code: "1221"},
		{Slug: "sales-executive", TitleEN: "B2B Sales Executive", TitleAM: "የሽያጭ ባለሙያ", TitleOM: "Ogeessa Gurgurtaa", TitleTI: "ናይ መሸጣ በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "3322"},

		// Agriculture & Environment
		{Slug: "agronomist", TitleEN: "Agronomist (Crop Specialist)", TitleAM: "የሰብል እና የአፈር ሳይንስ ባለሙያ (አግሮኖሚስት)", TitleOM: "Ogeessa Qonnaa fi Midhaanii", TitleTI: "ናይ ሰብልን ሕርሻን በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "2132"},
		{Slug: "veterinarian", TitleEN: "Veterinary Doctor (DVM)", TitleAM: "የእንስሳት ህክምና ሃኪም", TitleOM: "Doktoora Fayyaa Beelladaa", TitleTI: "ናይ እንስሳት ሓኪም", CategorySlug: "jobs-recruitment", Code: "2250"},
		{Slug: "animal-nutritionist", TitleEN: "Animal Nutritionist & Feed Specialist", TitleAM: "የእንስሳት መኖ እና ስነ-ምግብ ባለሙያ", TitleOM: "Ogeessa Nyaata Beelladaa", TitleTI: "ናይ እንስሳት መግቢ በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "2132"},
		{Slug: "irrigation-engineer", TitleEN: "Irrigation & Water Resources Engineer", TitleAM: "የመስኖ ምህንድስና ባለሙያ", TitleOM: "Injinara Jallisiifi Qabeenya Bishaanii", TitleTI: "ናይ መስኖ መሃንድስ", CategorySlug: "jobs-recruitment", Code: "2142"},

		// Technical Trades & Vocational (TVET)
		{Slug: "certified-electrician", TitleEN: "Licensed Industrial Electrician", TitleAM: "ፈቃድ ያለው የኤሌክትሪክ ባለሙያ", TitleOM: "Ogeessa Elektiriikii Mirkanaa'e", TitleTI: "ፍቓድ ዘለዎ ኤሌክትሪሽያን", CategorySlug: "jobs-recruitment", Code: "7411"},
		{Slug: "certified-plumber", TitleEN: "Master Plumber & Pipefitter", TitleAM: "የቧንቧ እና የሳኒተሪ ጥገና ባለሙያ", TitleOM: "Ogeessa Ujummoo Bishaanii", TitleTI: "ናይ ሻምብቆ ጽገና በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "7126"},
		{Slug: "automotive-mechanic", TitleEN: "Automotive Diagnostics & Engine Mechanic", TitleAM: "የመኪና ሞተር እና ኤሌክትሮኒክስ መካኒክ", TitleOM: "Meekaaniikii Konkolaataa", TitleTI: "ናይ መኪና መካኒክ", CategorySlug: "jobs-recruitment", Code: "7231"},
		{Slug: "welder-fabricator", TitleEN: "Certified Arc/MIG Welder & Fabricator", TitleAM: "የብረታ ብረት ብየዳ ባለሙያ", TitleOM: "Ogeessa Sibila Wal-qabsiisuu", TitleTI: "ናይ ሓጺን ብየዳ በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "7212"},
		{Slug: "heavy-equipment-operator", TitleEN: "Heavy Machinery Operator (Excavator / Loader)", TitleAM: "የከባድ ማሽነሪ ኦፕሬተር", TitleOM: "Ooppireetara Maashina Gurguddaa", TitleTI: "ናይ ከበድቲ ማሽነሪታት መዘወሪ", CategorySlug: "jobs-recruitment", Code: "8342"},
		{Slug: "professional-driver-cargo", TitleEN: "Heavy Dry Cargo Driver (Grade 4/5)", TitleAM: "የከባድ ጭነት መኪና ሹፌር", TitleOM: "Konkolaachisaa Fe'iinsa Gurguddaa", TitleTI: "ናይ ከቢድ ጽዕነት መኪና መራሒ", CategorySlug: "jobs-recruitment", Code: "8332"},
		{Slug: "carpenter-joiner", TitleEN: "Finish Carpenter & Furniture Maker", TitleAM: "የእንጨት ስራ እና የፈርኒቸር ባለሙያ", TitleOM: "Ogeessa Mukaa fi Barcumaa", TitleTI: "ናይ ዕንጨይቲ ስራሕ በዓል ሞያ", CategorySlug: "jobs-recruitment", Code: "7115"},
		{Slug: "solar-system-technician", TitleEN: "Solar PV Installation & Maintenance Technician", TitleAM: "የፀሐይ ብርሃን ኃይል ቴክኒሻን", TitleOM: "Teekiniishiyana Annisaa Soolaarii", TitleTI: "ናይ ጸሓይ ጸዓት ቴክኒሽያን", CategorySlug: "jobs-recruitment", Code: "7411"},

		// Legal, Media & Education
		{Slug: "licensed-attorney", TitleEN: "Licensed Attorney & Legal Counselor", TitleAM: "ጠበቃ እና የህግ አማካሪ", TitleOM: "Abukaatoo fi Gorsaa Seeraa", TitleTI: "ጠበቓን ናይ ሕጊ ኣማኻርን", CategorySlug: "jobs-recruitment", Code: "2611"},
		{Slug: "certified-translator", TitleEN: "Certified Multilingual Translator (EN/AM/OM/TI)", TitleAM: "የተመሰከረለት የቋንቋ ተርጓሚ", TitleOM: "Turjumaana Afaanii Mirkanaa'e", TitleTI: "ዝተመስከረሉ ተርጓሚ ቋንቋ", CategorySlug: "jobs-recruitment", Code: "2643"},
		{Slug: "high-school-teacher", TitleEN: "Secondary School Teacher", TitleAM: "የሁለተኛ ደረጃ መምህር", TitleOM: "Barsiisaa Sadarkaa Lammaffaa", TitleTI: "መምህር ካልኣይ ደረጃ", CategorySlug: "jobs-recruitment", Code: "2330"},
		{Slug: "university-lecturer", TitleEN: "University Lecturer & Academic Researcher", TitleAM: "የዩኒቨርሲቲ መምህር እና ተመራማሪ", TitleOM: "Barsiisaa fi Qorataa Yuunivarsiitii", TitleTI: "ናይ ዩኒቨርሲቲ መምህርን ተመራማርን", CategorySlug: "jobs-recruitment", Code: "2310"},
		{Slug: "commercial-photographer", TitleEN: "Professional Photographer & Cinematographer", TitleAM: "ፕሮፌሽናል ፎቶግራፈር እና ካሜራማን", TitleOM: "Suura-kaasaa fi Viidiyoo-kaasaa", TitleTI: "ፎቶግራፈርን ናይ ቪዲዮ ቀራጽን", CategorySlug: "jobs-recruitment", Code: "3431"},
		{Slug: "graphic-designer", TitleEN: "Graphic Designer & Brand Identity Specialist", TitleAM: "ግራፊክ ዲዛይነር", TitleOM: "Dizaayinara Giraafiiksii", TitleTI: "ግራፊክ ዲዛይነር", CategorySlug: "jobs-recruitment", Code: "2166"},
	}
}
