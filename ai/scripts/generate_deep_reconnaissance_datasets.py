#!/usr/bin/env python3
"""
HUNTIQ Deep Reconnaissance Dataset Generator.
Synthesizes comprehensive, balanced, multilingual training corpora across
English, Amharic, Afaan Oromoo, and Tigrinya for:
1. Domain Boundary Classification (In-Domain Marketplace vs External Knowledge)
2. Safety, Scam, Phishing & Moderation Guard
3. Marketplace 8-Taxonomy Category Classification
4. Sentiment, Review & Dispute Resolution Classification
"""

import json
import os

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
AI_ROOT = os.path.dirname(SCRIPT_DIR)
DATA_DIR = os.path.join(AI_ROOT, "training", "training_data")
os.makedirs(DATA_DIR, exist_ok=True)


# -------------------------------------------------------------
# 1. DOMAIN BOUNDARY (In-Domain Marketplace vs Out-Of-Domain)
# -------------------------------------------------------------
def build_domain_dataset():
    in_domain_templates = [
        # Vehicles
        ("Looking to buy a used {car} in {loc}", "ethiopian_marketplace"),
        (
            "What is the price of {car} with automatic transmission in {loc}?",
            "ethiopian_marketplace",
        ),
        ("Second hand {car} for sale clean libre title", "ethiopian_marketplace"),
        ("Is the {car} price negotiable or fixed?", "ethiopian_marketplace"),
        ("Can I inspect the {car} tomorrow at your garage?", "ethiopian_marketplace"),
        ("በ{loc_am} {car_am} መኪና መግዛት እፈልጋለሁ", "ethiopian_marketplace"),
        ("የ{car_am} መኪና ዋጋ ስንት ነው በ{loc_am}?", "ethiopian_marketplace"),
        ("ያገለገለ ንጹህ {car_am} ኦሪጅናል ሊብሬ ያለው", "ethiopian_marketplace"),
        (
            "Konkolaataa {car} {loc_om} keessatti bituu barbaada",
            "ethiopian_marketplace",
        ),
        ("Gatiin konkolaataa {car} meeqa {loc_om}tti?", "ethiopian_marketplace"),
        ("ኣብ {loc_ti} {car_ti} መኪና ክገዝእ ይደሊ ኣለኹ", "ethiopian_marketplace"),
        ("ዋጋ መኪና {car_ti} ኣብ {loc_ti} ክንደይ እዩ?", "ethiopian_marketplace"),
        # Real Estate
        (
            "Find 2 bedroom apartment for rent in {loc} under {price} birr",
            "ethiopian_marketplace",
        ),
        (
            "Are there furnished condominiums for rent in {loc}?",
            "ethiopian_marketplace",
        ),
        (
            "Looking for villa house with garden for sale in {loc}",
            "ethiopian_marketplace",
        ),
        (
            "Commercial office and shop space for lease in {loc}",
            "ethiopian_marketplace",
        ),
        ("በ{loc_am} የሚከራይ ባለ ሁለት መኝታ ቤት ከ{price} ብር በታች", "ethiopian_marketplace"),
        ("በ{loc_am} የሚሸጥ ቪላ ቤት ወይም ኮንዶሚኒየም ፈልግልኝ", "ethiopian_marketplace"),
        ("የንግድ ሱቅ ወይም ቢሮ ኪራይ በ{loc_am}", "ethiopian_marketplace"),
        ("Mana kireeffamu {loc_om}tti gatii {price} gadi", "ethiopian_marketplace"),
        ("Mana jireenyaa fi viillaa {loc_om}tti gurguramu", "ethiopian_marketplace"),
        ("ኣብ {loc_ti} ዝካረ ገዛ ትሕቲ {price} ቅርሺ", "ethiopian_marketplace"),
        ("ቪላ ገዛ ዝሽየጥ ኣብ {loc_ti}", "ethiopian_marketplace"),
        # Electronics & Products
        (
            "How much is {phone} brand new sealed with warranty?",
            "ethiopian_marketplace",
        ),
        ("Looking for used {laptop} laptop 16GB RAM Core i7", "ethiopian_marketplace"),
        ("Smart TV 55 inch 4K UHD with wall mount", "ethiopian_marketplace"),
        ("Where can I buy original {phone} in {loc}?", "ethiopian_marketplace"),
        ("የ{phone_am} ስልክ ዋጋ ስንት ነው በ{loc_am}?", "ethiopian_marketplace"),
        ("ያገለገለ {laptop_am} ላፕቶፕ ኮምፒውተር መግዛት እፈልጋለሁ", "ethiopian_marketplace"),
        ("ስማርት ቴሌቪዥን 55 ኢንች በዋስትና የሚሸጥ", "ethiopian_marketplace"),
        ("Gatiin bilbila {phone} meeqa haaraa?", "ethiopian_marketplace"),
        ("Lappitooppii {laptop} fayyadamaa bituu barbaada", "ethiopian_marketplace"),
        ("ዋጋ ሞባይል {phone_ti} ሓድሽ ብዋሕስ", "ethiopian_marketplace"),
        ("ላፕቶፕ {laptop_ti} ኮምፒተር ክገዝእ ይደሊ", "ethiopian_marketplace"),
        # Agriculture
        ("Wholesale teff {price} quintals available in {loc}", "ethiopian_marketplace"),
        ("Looking to buy wheat and barley grains in bulk", "ethiopian_marketplace"),
        ("Drip irrigation equipment and water pumps for farm", "ethiopian_marketplace"),
        ("የማኛ ጤፍ በኩንታል ዋጋ በ{loc_am} ስንት ነው?", "ethiopian_marketplace"),
        ("የስንዴ እና ገብስ የእህል ምርት በጅምላ የሚሸጥ", "ethiopian_marketplace"),
        ("የእርሻ ውሃ መሳቢያ ፓምፕ እና ጄኔሬተር", "ethiopian_marketplace"),
        ("Xaafii fi qamadii baay'inaan bituu barbaada", "ethiopian_marketplace"),
        ("ጣፍ ብኩንታል ኣብ {loc_ti} ክንደይ ይሰርሕ?", "ethiopian_marketplace"),
        # Jobs & Services
        (
            "Looking for accountant and finance manager job opening",
            "ethiopian_marketplace",
        ),
        ("Need experienced heavy truck driver with license", "ethiopian_marketplace"),
        (
            "Plumbing and electrical maintenance technician in {loc}",
            "ethiopian_marketplace",
        ),
        (
            "Wedding catering and professional photography services",
            "ethiopian_marketplace",
        ),
        ("የሂሳብ ባለሙያ ወይም የሽያጭ ሰራተኛ ክፍት የስራ ቦታ", "ethiopian_marketplace"),
        ("የቤት ወይም የቢሮ ጥገና ባለሙያ በ{loc_am} እፈልጋለሁ", "ethiopian_marketplace"),
        ("Hojii konkolaachisummaa fi herregaa barbaada", "ethiopian_marketplace"),
        ("ናይ ሒሳብ ሰራሕተኛ ክፍቲ ስራሕ ቦታ ኣብ {loc_ti}", "ethiopian_marketplace"),
        # Platform FAQ & Trade Rules
        ("How do I post a listing on HUNTIQ?", "ethiopian_marketplace"),
        (
            "What payment methods does HUNTIQ support for transactions?",
            "ethiopian_marketplace",
        ),
        ("How to check vehicle documents before buying?", "ethiopian_marketplace"),
        ("በHUNTIQ ላይ እቃዎችን እንዴት መለጠፍ እና መሸጥ እችላለሁ?", "ethiopian_marketplace"),
        ("የመኪና ሊብሬ እና ህጋዊ ሰነዶችን እንዴት ማረጋገጥ እችላለሁ?", "ethiopian_marketplace"),
        (
            "HUNTIQ irratti oomisha akkamitti maxxansuu danda'a?",
            "ethiopian_marketplace",
        ),
        ("ኣብ HUNTIQ ፍርያት ከመይ ጌረ ይልጥፍ?", "ethiopian_marketplace"),
    ]

    out_of_domain_samples = [
        # Coding & Software Architecture
        (
            "Write a Python function to implement quicksort algorithm",
            "general_out_of_domain",
        ),
        (
            "How to configure Docker compose for PostgreSQL and Redis container?",
            "general_out_of_domain",
        ),
        (
            "Explain the difference between React useMemo and useCallback hooks",
            "general_out_of_domain",
        ),
        (
            "Write a recursive binary search tree traversal in C++",
            "general_out_of_domain",
        ),
        ("How to debug a memory leak in Golang goroutines?", "general_out_of_domain"),
        (
            "Create an SQL query to calculate 7-day moving average of sales",
            "general_out_of_domain",
        ),
        (
            "Explain how convolutional neural networks process images in PyTorch",
            "general_out_of_domain",
        ),
        (
            "Write a bash script to backup logs to remote S3 bucket",
            "general_out_of_domain",
        ),
        ("የፓይዘን ፕሮግራሚንግ ኮድ እንዴት እንደሚጻፍ አስተምረኝ", "general_out_of_domain"),
        # Additional Amharic, Oromo, Tigrinya Coding, Medicine & Science
        ("በጃቫስክሪፕት ውስጥ ሪአክት ዩዝኢፌክት እንዴት ይሰራል?", "general_out_of_domain"),
        ("የጃቫስክሪፕት ኮዲንግ እና የሶፍትዌር አሰራር", "general_out_of_domain"),
        ("በጃቫስክሪፕት እና ፓይዘን ፕሮግራሚንግ እንዴት ይጻፋል?", "general_out_of_domain"),
        ("የሪአክት ፍሬምወርክ እና የዌብሳይት ኮድ እንዴት ይሰራል?", "general_out_of_domain"),
        ("የኖድ ጄኤስ ሰርቨር አሰራር እና የኮምፒውተር ኮዲንግ", "general_out_of_domain"),
        ("የኤስኪውኤል ዳታቤዝ ቴብል ክሪኤት ማድረጊያ ኩዌሪ", "general_out_of_domain"),
        ("የኤችቲኤምኤል እና የሲኤስኤስ ድህረ-ገጽ ኮዲንግ አሰራር", "general_out_of_domain"),
        ("የአልጎሪዝም ስሌት እና የኮምፒውተር ሳይንስ ትምህርት", "general_out_of_domain"),
        ("የባዮሎጂ የሴል ክፍፍል ማይቶሲስ እና ማዮሲስ", "general_out_of_domain"),
        ("የኬሚስትሪ የኦክስጅን እና የሃይድሮጅን ውህድ", "general_out_of_domain"),
        ("የፊዚክስ የስበት ህግ ቀመር እና ስሌት", "general_out_of_domain"),
        (
            "Mallattoon dhukkuba dhiibbaa dhiigaa fi sukkaaraa maali?",
            "general_out_of_domain",
        ),
        ("Akkamitti fayyaa sammuu fi qaamaa eegganna?", "general_out_of_domain"),
        ("Qoricha dhukkuba mataa fi bowwoo naaf ibsaa?", "general_out_of_domain"),
        ("Barnoota herregaa fi aljebraa na barsiisi", "general_out_of_domain"),
        ("ሕማም ሽኮር ምልክታቱን ፈውሱን ብዝርዝር ንገረኒ", "general_out_of_domain"),
        ("ናይ ፊዚክስን ኬሚስትሪን ቀመራት ኣረድኣኒ", "general_out_of_domain"),
        # World History, Geography & Politics
        (
            "Who was the first emperor of the Roman Empire and when did he rule?",
            "general_out_of_domain",
        ),
        (
            "Explain the main causes and consequences of World War I",
            "general_out_of_domain",
        ),
        (
            "What is the capital city of Argentina and its primary export?",
            "general_out_of_domain",
        ),
        (
            "Who was the author of the epic novel War and Peace?",
            "general_out_of_domain",
        ),
        (
            "Discuss the political impact of the Industrial Revolution in Britain",
            "general_out_of_domain",
        ),
        (
            "What are the major tributaries of the Amazon River in South America?",
            "general_out_of_domain",
        ),
        ("የፈረንሳይ አብዮት ታሪክ እና ዋና ዋና ምክንያቶች", "general_out_of_domain"),
        ("የመጀመሪያው እና የሁለተኛው የዓለም ጦርነት ታሪክ", "general_out_of_domain"),
        ("የአውስትራሊያ ዋና ከተማ ማን ትባላለች?", "general_out_of_domain"),
        ("Seenaa fi qorannoo urjiilee fi lafaa na barsiisi", "general_out_of_domain"),
        ("ዛንታ ታሪኽ ካልኣይ ኲናት ዓለም ንገረኒ", "general_out_of_domain"),
        # Sports, Entertainment, Creative & Pop Culture
        (
            "Who won the UEFA Champions League football final in 2024?",
            "general_out_of_domain",
        ),
        (
            "Can you write a romantic poem about the desert stars?",
            "general_out_of_domain",
        ),
        (
            "What is the recipe for baking homemade French sourdough bread?",
            "general_out_of_domain",
        ),
        ("Who won the Academy Award for Best Actor in 2023?", "general_out_of_domain"),
        (
            "Give me the summary and main themes of Shakespeare's Hamlet",
            "general_out_of_domain",
        ),
        (
            "What are the best training exercises to run a marathon?",
            "general_out_of_domain",
        ),
        ("የእንግሊዝ ፕሪሚየር ሊግ የእግር ኳስ ሻምፒዮና ውጤት", "general_out_of_domain"),
        ("የዶሮ ወጥ እና የክትፎ አሰራር ቅደም ተከተል", "general_out_of_domain"),
        ("ስለ ሼክስፒር ሃምሌት ድራማ ታሪክ አጫውተኝ", "general_out_of_domain"),
        ("የፍቅር ግጥም ወይም አጭር ልብ ወለድ ጻፍልኝ", "general_out_of_domain"),
        ("የድሮ የኢትዮጵያ የባህል ዘፈን እና የሙዚቃ ታሪክ", "general_out_of_domain"),
        ("ስለ ቴዲ አፍሮ ወይም ጥላሁን ገሰሰ ዘፈኖች ንገረኝ", "general_out_of_domain"),
        ("Sirba Afaan Oromoo beekamaa naaf barreessi", "general_out_of_domain"),
        ("Sirba jaalalaa fi aadaa Oromoo naaf weellisi", "general_out_of_domain"),
        ("Walaloo jaalalaa fi aadaa Oromoo naaf barreessi", "general_out_of_domain"),
        ("Seenaa fi aadaa Oromoo bal'inaan na barsiisi", "general_out_of_domain"),
        (
            "Taba kubbaa miilaa fi taphattoota beekamoo addunyaa",
            "general_out_of_domain",
        ),
        ("Asoosama jaalalaa fi kitaaba naaf barreessi", "general_out_of_domain"),
        (
            "Sirboota Hacaaluu Hundeessaa fi Alii Birraa naaf tarreessi",
            "general_out_of_domain",
        ),
        ("ናይ ኤውሮጳ ሻምፕዮንስ ሊግ መን ተዓዊቱ?", "general_out_of_domain"),
        ("ግጥሚ ፍቕሪ ወይ ናይ ልቢ ወለድ ዛንታ ጸሓፈለይ", "general_out_of_domain"),
        ("ናይ ትግርኛ ደርፊ ግጥሚ ወይ ባህላዊ ሙዚቃ ንገረኒ", "general_out_of_domain"),
    ]

    # Additional Agricultural Machinery & Pumps In-Domain
    in_domain_machinery = [
        ("የእርሻ ውሃ መሳቢያ ፓምፕ እና ጄኔሬተር", "ethiopian_marketplace"),
        ("የውሃ ፓምፕ ማሽን በናፍጣ የሚሰራ የሚሸጥ", "ethiopian_marketplace"),
        ("ጄኔሬተር 15 ኬቪኤ ለፋብሪካ ወይም ለእርሻ ስራ", "ethiopian_marketplace"),
        ("የእርሻ ትራክተር እና የእርሻ ማረሻ እቃዎች በኢትዮጵያ", "ethiopian_marketplace"),
        ("Bishaan harkiftuu fi jeneraetara qonnaa gurguramu", "ethiopian_marketplace"),
        ("Tiraaktara qonnaa fi meeshaalee qonnaa haaraa", "ethiopian_marketplace"),
        ("ማሽን ማይ መሳቢን ጄነሬተርን ንሕርሻ ዝሽየጥ", "ethiopian_marketplace"),
        ("ትራክተርን ናይ ሕርሻ መሳርሕን ኣብ ዕዳጋ", "ethiopian_marketplace"),
    ]
    for text, label in in_domain_machinery:
        out_of_domain_samples.append((text, label))

    # Substitutions for generating rich combinatorial in-domain data
    cars = [
        ("Toyota Vitz", "ቶዮታ ቪትዝ", "ቶዮታ ቪትዝ"),
        ("Toyota Corolla", "ቶዮታ ኮሮላ", "ቶዮታ ኮሮላ"),
        ("Suzuki Dzire", "ሱዙኪ ዲዛይር", "ሱዙኪ ዲዛይር"),
        ("Hyundai Tucson", "ሂዩንዳይ ቱክሰን", "ሂዩንዳይ ቱክሰን"),
        ("Toyota Yaris", "ቶዮታ ያሪስ", "ቶዮታ ያሪስ"),
        ("Toyota Hilux Pickup", "ቶዮታ ሃይሉክስ", "ቶዮታ ሃይሉክስ"),
    ]
    locs = [
        ("Addis Ababa", "አዲስ አበባ", "Finfinnee", "ኣዲስ ኣበባ"),
        ("Bole", "ቦሌ", "Boolee", "ቦሌ"),
        ("CMC", "ሲኤምሲ", "CMC", "ሲኤምሲ"),
        ("Hawassa", "ሀዋሳ", "Hawaasaa", "ሓዋሳ"),
        ("Adama", "አዳማ", "Adaamaa", "ኣዳማ"),
        ("Bahir Dar", "ባህር ዳር", "Baahir Daar", "ባሕር ዳር"),
        ("Mekelle", "መቀሌ", "Maqalee", "መቐለ"),
    ]
    phones = [
        ("Samsung Galaxy S24", "ሳምሰንግ ጋላክሲ ኤስ24", "ሳምሰንግ ጋላክሲ"),
        ("iPhone 15 Pro Max", "አይፎን 15 ፕሮ ማክስ", "ኣይፎን 15 ፕሮ"),
        ("Redmi Note 13", "ሬድሚ ኖት 13", "ሬድሚ ኖት 13"),
    ]
    prices = ["20,000", "35,000", "50,000", "80,000", "15,000"]

    dataset = []
    # Generate combinatorial in-domain
    for c_en, c_am, c_ti in cars:
        for l_en, l_am, l_om, l_ti in locs:
            dataset.append(
                (f"Looking to buy a used {c_en} in {l_en}", "ethiopian_marketplace")
            )
            dataset.append((f"በ{l_am} {c_am} መኪና መግዛት እፈልጋለሁ", "ethiopian_marketplace"))
            dataset.append(
                (
                    f"Konkolaataa {c_en} {l_om}tti bituu barbaada",
                    "ethiopian_marketplace",
                )
            )
            dataset.append(
                (f"ኣብ {l_ti} {c_ti} መኪና ክገዝእ ይደሊ ኣለኹ", "ethiopian_marketplace")
            )

    for p_en, p_am, p_ti in phones:
        for l_en, l_am, l_om, l_ti in locs:
            dataset.append(
                (f"Where can I buy original {p_en} in {l_en}?", "ethiopian_marketplace")
            )
            dataset.append((f"የ{p_am} ስልክ ዋጋ ስንት ነው በ{l_am}?", "ethiopian_marketplace"))
            dataset.append(
                (
                    f"Bilbila {p_en} {l_om}tti meeqaan gurgurama?",
                    "ethiopian_marketplace",
                )
            )
            dataset.append(
                (f"ዋጋ ሞባይል {p_ti} ኣብ {l_ti} ክንደይ ይሰርሕ?", "ethiopian_marketplace")
            )

    for pr in prices:
        for l_en, l_am, l_om, l_ti in locs:
            dataset.append(
                (
                    f"Find 2 bedroom apartment for rent in {l_en} under {pr} birr",
                    "ethiopian_marketplace",
                )
            )
            dataset.append(
                (f"በ{l_am} የሚከራይ ባለ ሁለት መኝታ ቤት ከ{pr} ብር በታች", "ethiopian_marketplace")
            )
            dataset.append(
                (f"Mana kireeffamu {l_om}tti gatii {pr} gadi", "ethiopian_marketplace")
            )
            dataset.append((f"ኣብ {l_ti} ዝካረ ገዛ ትሕቲ {pr} ቅርሺ", "ethiopian_marketplace"))

    # Add core templates
    for text, label in in_domain_templates:
        dataset.append(
            (
                text.format(
                    car="Toyota Corolla",
                    loc="Bole",
                    car_am="ቶዮታ",
                    loc_am="ቦሌ",
                    car_ti="ቶዮታ",
                    loc_ti="ቦሌ",
                    loc_om="Finfinnee",
                    price="25,000",
                    phone="iPhone",
                    phone_am="አይፎን",
                    phone_ti="ኣይፎን",
                    laptop="Dell",
                    laptop_am="ዴል",
                    laptop_ti="ዴል",
                ),
                label,
            )
        )

    # Add out of domain
    for text, label in out_of_domain_samples:
        dataset.append((text, label))

    # Generate additional rich out of domain variations
    academic_topics = [
        "quantum entanglement and Einstein-Podolsky-Rosen paradox",
        "organic synthesis of aromatic hydrocarbons and benzene rings",
        "microeconomic utility maximization and indifference curves",
        "causes of the Roman empire collapse in Western Europe",
        "how mRNA vaccines stimulate immunological spike protein synthesis",
        "asymptotic computational complexity of Dijkstra shortest path algorithm",
        "astronomical measurement of cosmic microwave background radiation",
        "philosophical epistemology of Immanuel Kant categorical imperative",
        "mathematical derivation of Navier-Stokes fluid mechanics equations",
        "cognitive psychology of neuroplasticity and memory consolidation",
    ]
    for topic in academic_topics:
        dataset.append(
            (
                f"Explain in technical detail the concepts of {topic}",
                "general_out_of_domain",
            )
        )
        dataset.append(
            (
                f"Write a comprehensive research paper on {topic}",
                "general_out_of_domain",
            )
        )
        dataset.append(
            (
                f"Solve the homework assignment problem regarding {topic}",
                "general_out_of_domain",
            )
        )

    coding_languages = [
        "Python",
        "Rust",
        "Golang",
        "TypeScript",
        "C++",
        "Java",
        "Kotlin",
        "Swift",
    ]
    code_tasks = [
        "binary tree inversion",
        "LRU cache implementation",
        "matrix multiplication",
        "WebSocket connection pool",
        "OAuth2 authorization code flow",
        "regex parser",
    ]
    for lang in coding_languages:
        for task in code_tasks:
            dataset.append(
                (f"Write a {lang} code script for {task}", "general_out_of_domain")
            )
            dataset.append(
                (
                    f"How to optimize {task} performance in {lang}?",
                    "general_out_of_domain",
                )
            )

    # Additional legal and documentation in-domain samples
    legal_in_domain = [
        ("የመኪና ሊብሬ እና ህጋዊ ሰነዶችን እንዴት ማረጋገጥ እችላለሁ?", "ethiopian_marketplace"),
        ("የመኪና ሊብሬ ማረጋገጫ በውልና ማስረጃ እንዴት ይደረጋል?", "ethiopian_marketplace"),
        ("የባለቤትነት ማረጋገጫ ሊብሬ እና የውል ሰነድ ህጋዊነት", "ethiopian_marketplace"),
        ("መኪና ከመግዛት በፊት ሊብሬ እንዴት ይጣራል?", "ethiopian_marketplace"),
        ("የንግድ ውል እና የመኪና ሊብሬ ትክክለኛነት ማጣሪያ", "ethiopian_marketplace"),
        ("በHUNTIQ ላይ የመኪና ሊብሬ እና ሰነድ እንዴት ማረጋገጥ እችላለሁ?", "ethiopian_marketplace"),
    ]
    for text, label in legal_in_domain:
        dataset.append((text, label))

    # Additional poetry, songs, and literature out-of-domain samples
    general_creative = [
        ("የፍቅር ግጥም ወይም አጭር ልብ ወለድ ጻፍልኝ", "general_out_of_domain"),
        ("ግጥም ጻፍልኝ ስለ ፍቅር ወይም ስለ ሀገር ውበት", "general_out_of_domain"),
        ("ልብ ወለድ ድርሰት ወይም ልብ አንጠልጣይ ታሪክ ጻፍልኝ", "general_out_of_domain"),
        ("የባህል ሙዚቃ ግጥም ወይም ዜማ አዘጋጅልኝ", "general_out_of_domain"),
        ("የፍቅር ደብዳቤ ወይም አጭር ግጥም ጻፍልኝ", "general_out_of_domain"),
    ]
    for text, label in general_creative:
        dataset.append((text, label))

    return dataset


# -------------------------------------------------------------
# 2. SAFETY & SCAM MODERATION DATASET
# -------------------------------------------------------------
def build_safety_dataset():
    dataset = []

    # SAFE MARKETPLACE COMMUNICATIONS (EN, AM, OM, TI)
    safe_bases = [
        "I want to buy a used Toyota Vitz in Addis Ababa",
        "How much is Samsung Galaxy S23 in Bole?",
        "Are there 2 bedroom apartments for rent in Kazanchis?",
        "Looking for a professional wedding photographer",
        "Can I inspect the car tomorrow before buying?",
        "Is the laptop price negotiable?",
        "Do you offer delivery to Hawassa?",
        "Where is your shop located?",
        "Hello, is this still available?",
        "I am looking for teff for sale in Adama",
        "Need agricultural tractor rental for 3 days",
        "What is the warranty period for this generator?",
        "Do you have iPhone 15 Pro Max 256GB?",
        "Searching for office desk and ergonomic chair",
        "I want to buy second-hand Nike sneakers size 42",
        "Can we meet in public place like Edna Mall to inspect the phone?",
        "Please provide the IMEI number so I can check device status",
        "I will pay cash upon physical inspection and receipt verification",
        "Does the vehicle have all DARO authentication documents?",
        "Are maintenance records and service history available for this car?",
        "አዲስ አበባ ውስጥ ቪትዝ መኪና መግዛት እፈልጋለሁ",
        "ቦሌ የሚከራይ ባለ ሁለት መኝታ ቤት አለ?",
        "የሳምሰንግ ስልክ ዋጋ ስንት ነው?",
        "እቃውን አይቼ መክፈል እችላለሁ?",
        "ሱቃችሁ የት ነው የሚገኘው?",
        "ዋጋ ይቀነሳል ወይስ የመጨረሻ ነው?",
        "ዴል ላፕቶፕ ኮር አይ 7 ፈልጌ ነበር",
        "የሚሸጥ ንጹህ ጤፍ አለ?",
        "ለሰርግ ፎቶግራፍ አንሺ እፈልጋለሁ",
        "መኪናውን ነገ በአካል ማየት እችላለሁ?",
        "ስልኩ ኦሪጅናል መሆኑን በIMEI ማረጋገጥ እፈልጋለሁ",
        "እቃውን በሰው ፊት እና ደህንነቱ በተጠበቀ ቦታ እንፈትሽ",
        "ሰነዶች በሙሉ በውልና ማስረጃ የተረጋገጡ ናቸው?",
        "Konkolaataa Vitz Finfinneetti bituu barbaada",
        "Gatiin bilbila Samsung meeqa?",
        "Booleetti mana kireeffamu qabduu?",
        "Meeshicha ilaalee kaffaluu danda'aa?",
        "Suuqii keessan eessatti argama?",
        "Gatiin isaa xiqqoo ni hir'ataa?",
        "Baraartuu qonnaa kireeffachuu barbaada",
        "Lappitoppii Dell barbaada",
        "Bakka namoonni baay'atanitti wal agarree ilaaluu dandeenyaa?",
        "ኣብ ኣዲስ ኣበባ ቪትዝ መኪና ክገዝእ ይደሊ ኣለኹ",
        "ዋጋ ሳምሰንግ ሞባይል ክንደይ እዩ?",
        "ኣብ ቦሌ ዝካረ ገዛ ኣለኩምዶ?",
        "ነቲ ኣቕሓ ርእየ ክኸፍል ይኽእልዶ?",
        "ዱካንኩም ኣበይ እዩ ዝርከብ?",
        "ዋጋ ይንኪ ድዩ ዋላ ናይ መወዳእታ እዩ?",
        "ንጹህ ጣፍ ዝሽየጥ ኣለዶ?",
        "ነታ መኪና ኣብ ጋራጅ ኣእቲና ንመርምራ",
    ]
    for s in safe_bases:
        dataset.append((s, "safe"))

    # SCAM ADVANCE PAYMENT (Requests for upfront wire/Telebirr before viewing)
    scam_templates = [
        "Please send {pct}% advance deposit to my Telebirr account before seeing the {item}",
        "Transfer {amt} birr to my CBE account first, then I will ship the {item} by bus",
        "Deposit {amt} ETB reservation fee immediately via mobile money or the {item} will be sold",
        "The {item} is stored in a customs warehouse, wire {amt} Birr clearance fee first",
        "Send inspection fee of {amt} birr to my personal Telebirr before I show the house",
        "I am currently out of town, send {pct}% down payment and my driver will bring the {item}",
        "Deposit registration fee of {amt} Birr to secure this exclusive job position",
        "You must pay shipping fee of {amt} birr in advance before delivery driver leaves",
        "Wire money directly to my bank account, do not use HUNTIQ payment system",
        "Send token money via Awash Bank to hold the apartment until tomorrow",
        "መኪናውን ከማየትህ በፊት {pct}% በቴሌብር በቅድሚያ ላክ",
        "ስልኩን በፖስታ እንድልክልህ {amt} ብር በቅድሚያ በንግድ ባንክ አስገባ",
        "ቤቱን እንድይዝልዎት አሁኑኑ የቅድመ ክፍያ በቴሌብር ይላኩ",
        "ባለቤቱ ለጉዞ ስለቸኮለ ዋጋው ረክሷል፤ ገንዘቡን ወዲያውኑ በባንክ አስገቡ",
        "እቃውን ከመላካችን በፊት የትራንስፖርት ክፍያ በቴሌብር ቅድሚያ መላክ አለብዎት",
        "የምዝገባ ክፍያ {amt} ብር በቅድሚያ ካላስገቡ ስራውን አያገኙም",
        "ገንዘቡን በሲቢኢ ብር ላኩልኝ እና እቃውን በአውቶቡስ እልክላችኋለሁ",
        "ቤት ለማየት {amt} ብር የቦታ ማሳያ ቅድመ ክፍያ በቴሌብር ላክ",
        "እቃውን ሳይዙ በፊት በቴሌብር ግማሹን ገንዘብ አስገቡ",
        "Konkolaaticha osoo hin argin dura {pct}% kaffaltii telebirr dursa ergi",
        "Bilbila siif erguuf qarshii {amt} baankiin dursii ergi",
        "Mana siif qabuuf dursiiti telebirriin maallaqa dabarsi",
        "Dursitee kaffaltii geejjibaa ergi eegasii meeshaa siif ergina",
        "ነታ መኪና ቅድሚ ምርኣይካ {pct}% ቅድመ ክፍሊት ብቴሌብር ስደድ",
        "ስልኪ ቅድሚ ምልኣኸይ {amt} ቅርሺ ብባንኪ ስደድ",
        "ነቲ ገዛ ንኸይሸጦ ቅድመ ክፍሊት ብሲቢኢ ሕጂ ስደደለይ",
        "ቅድሚ እቲ ኣቕሓ ምብጻሑ ናይ ትራንስፖርት ገንዘብ ብቅድሚት ስደድ",
    ]
    for pct in ["20", "30", "50", "70"]:
        for amt in ["1,500", "2,000", "3,500", "5,000", "10,000"]:
            for item in ["car", "phone", "laptop", "generator", "apartment"]:
                for t in scam_templates[:10]:
                    dataset.append(
                        (t.format(pct=pct, amt=amt, item=item), "scam_advance")
                    )
    for pct in ["20", "30", "50"]:
        for amt in ["1000", "2000", "3000"]:
            for t in scam_templates[10:]:
                dataset.append((t.format(pct=pct, amt=amt, item="መኪና"), "scam_advance"))

    # PHISHING & CREDENTIAL HARVESTING
    phishing_bases = [
        "Click this link to verify your HUNTIQ bank account login: http://huntiq-fake.xyz",
        "Send your CBE Birr PIN and password to claim your cash prize",
        "Join our secret telegram channel t.me/free_crypto_ethiopia to get paid",
        "Enter your Telebirr secret code to confirm the payment receipt",
        "Download this APK file to activate your seller account",
        "Provide your bank account login credentials and OTP for direct deposit",
        "Click here to claim 50,000 Birr government gift on Telegram",
        "We need your ATM card number and CVV to process refund",
        "Verify your account by entering your 4-digit mobile banking PIN",
        "HUNTIQ security team: please send your SMS confirmation code immediately",
        "Click here to unlock your frozen bank account: http://cbe-verify.net",
        "የባንክ አካውንትዎን እና የቴሌብር ፒን ቁጥርዎን እዚህ ሊንክ ላይ ያስገቡ",
        "የቴሌብር ሚስጥር ቁጥርዎን (PIN) ይላኩልንና ሽልማትዎን ይውሰዱ",
        "ይህንን የቴሌግራም ሊንክ በመጫን 50 ሺህ ብር ተሸላሚ ይሁኑ",
        "አካውንትዎ እንዳይዘጋ ይህንን ፋይል ዳውንሎድ አድርገው ይክፈቱ",
        "የሲቢኢ ብር ኦቲፒ (OTP) ኮድዎን ላኩልኝ ክፍያውን ላጽድቅልዎት",
        "የባንክ ካርድ ቁጥር እና ከጀርባ ያለውን ባለ 3 ዲጂት ኮድ ላክ",
        "የይለፍ ቃልዎን እና ሚስጥር ቁጥርዎን እዚህ ድረ-ገጽ ላይ ያረጋግጡ",
        "የቴሌብር ኮድዎን ላኩልንና ገንዘብዎን በእጥፍ ያግኙ",
        "Koodii iccitii (PIN) telebirr keessanii ergaa badhaasa fudhadhaa",
        "Liinkii kana cuqaasaa herrega baankii keessan mirkaneessaa",
        "Koodii OTP baankii dhufe nuuf ergaa",
        "Faayila kana buufadhaa herrega keessan banaa",
        "ሚስጢራዊ ቁጽሪ (PIN) ቴሌብርካ ብምልኣኽ ሽልማትካ ውሰድ",
        "ናይ ባንኪ OTP ኮድ ስደደለይ ገንዘብካ ከእትወልካ",
        "ኣካውንትካ ንኸይዕጾ ነዚ ሊንክ ጠውቕ",
        "ናይ ባንኪ ካርድ ቁጽርን ፓስዎርድን ስደድ",
    ]
    for p in phishing_bases:
        for suffix in [
            "",
            " immediately",
            " urgently",
            " or account will be suspended",
            " to claim balance",
        ]:
            dataset.append((p + suffix, "phishing_harvesting"))

    # PROHIBITED CONTRABAND & WEAPONS
    contraband_bases = [
        "Looking to buy unregistered firearms and ammunition AK-47",
        "Where can I purchase fake counterfeit Birr banknotes?",
        "Selling cocaine and illegal narcotics discreet delivery",
        "Want to buy military weapons and pistols without license",
        "Fake Ethiopian passports and forged university degrees for sale",
        "Selling unregistered handguns and bullets in Addis",
        "Counterfeit 100 and 200 birr notes indistinguishable from real",
        "Smuggled illegal firearms available contact inbox",
        "Buy unlicensed sniper rifle and explosive gear",
        "Selling illegal marijuana weed and psychedelic drugs in bulk",
        "Forged government seals and fake driving licenses",
        "ያልተመዘገበ ክላሽ እና ሽጉጥ መግዛት እፈልጋለሁ",
        "የተሰረዘ ሀሰተኛ የብር ኖት 100 እና 200 ብር የሚሸጥ የት ይገኛል?",
        "ህገወጥ የጦር መሳሪያ ሽጉጥ እና ጥይት የሚሸጥ",
        "ሀሰተኛ የዩኒቨርሲቲ ዲግሪ እና ፓስፖርት የሚያዘጋጅ ሰው ፈልጌ ነበር",
        "ኮኬይን እና አደንዛዥ እፅ መግዛት እፈልጋለሁ",
        "የኮንትሮባንድ መሳሪያ ያለፈቃድ የሚሸጥ",
        "የተጭበረበረ የገንዘብ ኖት ሀሰተኛ ብር",
        "ሀሰተኛ መንጃ ፈቃድ እና የቀበሌ መታወቂያ እንሰራለን",
        "ህገወጥ አደንዛዥ እፅ ጫት እና ሄሮይን",
        "ያለ ህጋዊ ፈቃድ ሽጉጥ እና ጥይት መግዛት እፈልጋለሁ",
        "Qawwee fi rasaasa seeraan alaa bituu barbaada",
        "Qarshii sobaa Itoophiyaa bituu barbaada",
        "Qoricha sammuu nama hadoochu seeraan alaa",
        "Meesshaa waraanaa fi qawwee seeraan alaa",
        "Waraqaa eenyummaa fi heeyyama konkolaachisummaa sobaa",
        "ዘይሕጋዊ ኣጽዋርን ሽጉጥን ክገዝእ ይደሊ ኣለኹ",
        "ናይ ሓሶት ቅርሺ ኖት ዝሸይጥ ሰብ ኣለዶ?",
        "ዘይሕጋዊ ዕጸ-ፋርስ ክገዝእ ይደሊ",
        "ናይ ሓሶት ወረቐት መንነትን መንጃ ፍቓድን",
        "ናይ ኲናት ኣጽዋር ክላሽን ሽጉጥን ጥይትን",
    ]
    for c in contraband_bases:
        for prefix in [
            "",
            "Secret dealer: ",
            "Urgent: ",
            "Inbox price for: ",
            "Looking for: ",
        ]:
            dataset.append((prefix + c, "prohibited_contraband"))

    # HARASSMENT & ABUSE
    harassment_bases = [
        "You are an idiot scammer I will find you and harm you",
        "Shut up and stop messaging me you stupid fool",
        "I know where you live I will kill you if you do not pay",
        "Hate speech against ethnic groups in Ethiopia",
        "You worthless thief go die",
        "Get lost you stupid idiot",
        "I will destroy you and your family you bastard",
        "Stop calling me you filthy dog",
        "You ugly scammer I will track your IP address and beat you",
        "Stupid liar you are wasting my time go to hell",
        "You dirty piece of trash I will make you regret this",
        "Shut your mouth you shameless criminal",
        "ደደብ አጭበርባሪ አገኝህና አሳይሃለሁ",
        "የሌባ ልጅ አፍህን ዝጋ",
        "የምትኖርበትን አውቃለሁ ገድልሃለሁ",
        "ዘረኛ ስድብ እና ዛቻ",
        "የውሻ ልጅ ራስህን አጥፋ",
        "ከዚህ ጥፋ አንተ ደደብ ጅል",
        "አንተ ቆሻሻ ሌባ ዋጋህን ታገኛለህ",
        "አንገትህን እቀላዋለሁ እንዳላይህ",
        "የቆሻሻ ዘር አፍህን ዝጋ",
        "Ati hattuu fokkisaa si ajjeesa",
        "Afaan kee cufadhu doofaa",
        "Saree fokkisaa asii deemi",
        "Eessatti akka jirtu nan beeka si balleessa",
        "Si qabeen si tuma doofaa",
        "ንስኻ ሰራቒ ኢኻ ክቐትለካ እየ",
        "ኣፍካ ዕጾ ዓሻ ደደብ",
        "ኣነ ክቐትለካ እየ ሕነይ ክፈዲ እየ",
        "ናይ ከልቢ ወዲ ካብዚ ርሓቕ",
        "ክቐትለካ እየ ኣብዚ ከይትርከብ",
    ]
    for h in harassment_bases:
        for prefix in [
            "",
            "Hey listen: ",
            "Warning you: ",
            "You bastard: ",
            "Fuck off: ",
        ]:
            dataset.append((prefix + h, "harassment_abuse"))

    return dataset


# -------------------------------------------------------------
# 3. CATEGORY TAXONOMY DATASET (8 Classes)
# -------------------------------------------------------------
def build_category_dataset():
    dataset = []
    cat_items = {
        "products": [
            ("iPhone 14 Pro Max 256GB storage space gray", "products"),
            ("Samsung Galaxy S23 Ultra brand new sealed", "products"),
            ("Dell XPS 15 laptop with 32GB RAM Core i9", "products"),
            ("Sony WH-1000XM5 wireless noise cancelling headphones", "products"),
            ("Smart TV 55 inch 4K UHD with warranty", "products"),
            ("MacBook Pro M2 16 inch Apple Silicon", "products"),
            ("HP Pavilion gaming laptop RTX 3060", "products"),
            ("Canon EOS R6 mirrorless digital camera", "products"),
            ("iPad Air 5th generation 64GB WiFi", "products"),
            ("PlayStation 5 digital edition with 2 controllers", "products"),
            ("ስማርት ስልክ ሳምሰንግ ጋላክሲ በቦሌ", "products"),
            ("ላፕቶፕ ኮምፒውተር ዴል 16 ጂቢ ራም", "products"),
            ("ቴሌቪዥን 55 ኢንች ስማርት ቲቪ", "products"),
            ("የድምጽ ማጉያ ስፒከር ኦሪጅናል", "products"),
            ("አፕል አይፎን 15 ፕሮ ማክስ በካርቶን", "products"),
            ("ማክቡክ ፕሮ ኮምፒውተር የሚሸጥ", "products"),
            ("የጨዋታ ፕሌይስቴሽን 5 ኮንሶል", "products"),
            ("ዲጂታል ካሜራ ካኖን በሌንስ", "products"),
            ("Bilbila ammayyaa iPhone gurgurtaa", "products"),
            ("Kompitaroota laaptooppii Dell fi HP", "products"),
            ("Meeshaalee elektirooniksii Finfinnee keessatti", "products"),
            ("Tivii ammayyaa smart tv 4k", "products"),
            ("ስልኪ ሳምሰንግ ጋላክሲ ዝሽየጥ", "products"),
            ("ኮምፒተር ላፕቶፕ ዴል 16GB", "products"),
            ("ስማርት ቲቪ ተለቪዥን ሓድሽ", "products"),
            ("ኣይፎን 14 ፕሮ ማክስ ካርቶን ዘይተኸፍተ", "products"),
        ],
        "real_estate": [
            ("2 bedroom modern apartment for rent in Bole", "real_estate"),
            ("Luxury furnished villa with garden in CMC", "real_estate"),
            ("Studio condominium unit for rent in Kazanchis", "real_estate"),
            ("Commercial shop space for lease in Piassa", "real_estate"),
            ("3 bedroom townhouse for sale in Summit Addis", "real_estate"),
            ("Office floor space 200 sqm for rent in Bole", "real_estate"),
            ("Residential land plot 500 sqm with master plan", "real_estate"),
            ("በቦሌ የሚከራይ ባለ ሁለት መኝታ ዘመናዊ አፓርታማ", "real_estate"),
            ("በሲኤምሲ የሚሸጥ ቪላ ቤት ከጓሮ ጋር", "real_estate"),
            ("በካዛንቺስ የሚከራይ ስቱዲዮ ኮንዶሚኒየም", "real_estate"),
            ("በፒያሳ የሚከራይ የንግድ ሱቅ ወይም ቢሮ", "real_estate"),
            ("ባለ 3 መኝታ ቤት በሰሚት የሚሸጥ", "real_estate"),
            ("የመኖሪያ ወይም የንግድ ቦታ 500 ካሬ ሜትር", "real_estate"),
            ("በጀሞ የሚከራይ ባለ አንድ መኝታ ኮንዶሚኒየም", "real_estate"),
            ("Mana jireenyaa fi appaartamaa Booleetti kireeffamu", "real_estate"),
            ("Viillaa bareedaa fi lafa Finfinneetti gurguramu", "real_estate"),
            ("Mana kondoominiyeemii kiraa Kazanchis", "real_estate"),
            ("Suuqii fi kutaa daldalaa kireeffamu", "real_estate"),
            ("ኣብ ቦሌ ዝካረ ክልተ ክፍሊ ዘለዎ ኣፓርታማ", "real_estate"),
            ("ቪላ ገዛ ምስ ስፍሓት መሬት ዝሽየጥ ኣብ ሲኤምሲ", "real_estate"),
            ("ኮንዶሚንየም ገዛ ክራይ ኣብ ካዛንቺስ", "real_estate"),
            ("ናይ ንግዲ ዱካን ወይ ቢሮ ዝካረ", "real_estate"),
        ],
        "vehicles": [
            ("Toyota Corolla 2020 automatic transmission", "vehicles"),
            ("Toyota Vitz 2018 clean condition low mileage", "vehicles"),
            ("Suzuki Swift Dzire 2021 fuel efficient", "vehicles"),
            ("Hyundai Tucson 2019 AWD leather seats", "vehicles"),
            ("Commercial Isuzu NPR truck 3.5 ton for sale", "vehicles"),
            ("Toyota Hilux double cab pickup 4x4 diesel", "vehicles"),
            ("Bajaj RE 4-stroke passenger model 2022", "vehicles"),
            ("Honda CR-V compact SUV ready for transfer", "vehicles"),
            ("ቶዮታ ኮሮላ 2020 አውቶማቲክ መኪና", "vehicles"),
            ("ቶዮታ ቪትዝ ንጹህ ያገለገለ አነስተኛ ኪሎሜትር", "vehicles"),
            ("ሱዙኪ ዲዛይር 2021 የነዳጅ ቆጣቢ መኪና", "vehicles"),
            ("ሂዩንዳይ ቱክሰን 2019 ባለ ሙሉ አማራጭ", "vehicles"),
            ("የጭነት መኪና ኢሱዙ ኤንፒአር የሚሸጥ", "vehicles"),
            ("ቶዮታ ሃይሉክስ ድርብ ጋቢና ፒክአፕ", "vehicles"),
            ("የተሳፋሪ ባጃጅ አራት ስትሮክ የሚሸጥ", "vehicles"),
            ("Konkolaataa Toyota Corolla ootompaatikii", "vehicles"),
            ("Konkolaataa Toyota Vitz qulqulluu gurguramu", "vehicles"),
            ("Konkolaataa fe'umsaa Isuzu fi baajaajii", "vehicles"),
            ("Konkolaataa Suzuki Swift gatii salphaan", "vehicles"),
            ("መኪና ቶዮታ ኮሮላ 2020 ኦቶማቲክ", "vehicles"),
            ("ቶዮታ ቪትዝ ንጹህ መኪና ብሕሳብ", "vehicles"),
            ("ናይ ጽዕነት መኪና ኢሱዙ ዝሽየጥ", "vehicles"),
            ("ባጃጅ ተሳፋሪ ሓድሽ ሞዴል", "vehicles"),
        ],
        "jobs": [
            ("Experienced full stack web developer Python React", "jobs"),
            ("Senior accountant with IFRS certification required", "jobs"),
            ("Heavy duty truck driver with grade 4 license", "jobs"),
            ("Sales and marketing representative for retail FMCG", "jobs"),
            ("Graphic designer proficient in Photoshop Illustrator", "jobs"),
            ("Civil site engineer for construction building project", "jobs"),
            ("Nurse and health officer for private clinic vacancy", "jobs"),
            ("የድርጅት ከፍተኛ የሂሳብ ባለሙያ ክፍት የስራ ቦታ", "jobs"),
            ("የድር ድህረ ገጽ እና የሞባይል መተግበሪያ አልሚ", "jobs"),
            ("ደረጃ 4 መንጃ ፈቃድ ያለው የከባድ መኪና አሽከርካሪ", "jobs"),
            ("የሽያጭና ግብይት ሰራተኛ ለአስቸኳይ ስራ", "jobs"),
            ("የግራፊክስ ዲዛይነር እና የቪዲዮ ኤዲተር", "jobs"),
            ("የኮንስትራክሽን ሳይት መሃንዲስ የስራ ማስታወቂያ", "jobs"),
            ("ነርስ እና የጤና ባለሙያ ለክሊኒክ", "jobs"),
            ("Hojii ogeessa teeknoolojii fi sooftiweerii", "jobs"),
            ("Hojii ogeessa herregaa fi maallaqaa", "jobs"),
            ("Konkolaachisaa heeyyama qabu hojiif barbaadama", "jobs"),
            ("Hojii gurgurtaa fi beeksisa oomishaa", "jobs"),
            ("ናይ ሶፍትዌርን ኮምፒተርን ኢንጅነር ስራሕ", "jobs"),
            ("ናይ ሒሳብ ሰራሕተኛ ክፍቲ ስራሕ ቦታ", "jobs"),
            ("መንጃ ፍቓድ ዘለዎ መዋፈርቲ መኪና", "jobs"),
            ("ናይ ህንጻ መሃንድስ ክፍቲ ስራሕ መወዓውዒ", "jobs"),
        ],
        "agriculture": [
            ("Grade 1 white teff 100 quintals wholesale Gojjam", "agriculture"),
            ("Wheat grain and barley harvest bulk purchase", "agriculture"),
            ("Drip irrigation hose and water pump for horticulture", "agriculture"),
            ("Organic fresh honey from Bale mountains bulk container", "agriculture"),
            ("High yield coffee beans washed Sidama grade 2", "agriculture"),
            ("Livestock sheep and cattle fattened for holiday", "agriculture"),
            ("Agricultural tractor 75 horsepower four wheel drive", "agriculture"),
            ("የመጀመሪያ ደረጃ ነጭ ማኛ ጤፍ በኩንታል የጎጃም ምርት", "agriculture"),
            ("የስንዴ እና ገብስ የእህል ምርት በጅምላ አቅራቢ", "agriculture"),
            ("የእርሻ መስኖ ቱቦ እና የውሃ መሳቢያ ፓምፕ", "agriculture"),
            ("የባሌ ተፈጥሯዊ ንጹህ ማር በጀሪካን", "agriculture"),
            ("የሲዳማ ታጠበ ቡና በኮንቴነር ለገበያ", "agriculture"),
            ("የበሬ እና የሰንጋ ከብቶች ለበዓል የሚሸጡ", "agriculture"),
            ("የእርሻ ትራክተር ባለ 75 ፈረስ ጉልበት", "agriculture"),
            ("Oomisha xaafii adii qulqulluu baay'inaan", "agriculture"),
            ("Qamadii fi garbuu callaa gabaa qonnaa", "agriculture"),
            ("Buna Siidaamaa fi damma gaara Baalee", "agriculture"),
            ("Tiraaktara qonnaa fi bishaan harkiftuu", "agriculture"),
            ("ቀዳማይ ደረጃ ጻዕዳ ጣፍ ብኩንታል", "agriculture"),
            ("ስርናይን ስገምን ናይ ሕርሻ ፍርያት ብጅምላ", "agriculture"),
            ("ናይ ሲዳማ ጽሩይ ቡን ንዕዳጋ", "agriculture"),
            ("ትራክተር ሕርሻ ዝሽየጥ", "agriculture"),
        ],
        "services": [
            ("Professional plumbing and pipeline leak repair service", "services"),
            ("Electrical wiring installation for residential house", "services"),
            ("Wedding video and portrait photography studio service", "services"),
            ("Moving and relocation logistics truck with helpers", "services"),
            ("Furniture carpentry and custom cabinet woodwork", "services"),
            ("Legal consultation and business contract drafting", "services"),
            ("HVAC air conditioning repair and maintenance", "services"),
            ("የቧንቧ እና የፍሳሽ መስመር ጥገና ባለሙያ አገልግሎት", "services"),
            ("የቤት የኤሌክትሪክ መስመር ዝርጋታ እና ጥገና", "services"),
            ("የሰርግ ቪዲዮ እና ፎቶግራፍ ስቱዲዮ አገልግሎት", "services"),
            ("የቤት እቃ ማጓጓዝ እና የሎጅስቲክስ አገልግሎት", "services"),
            ("የእንጨት ስራ የፈርኒቸር እና ቁምሳጥን ማዘጋጀት", "services"),
            ("የህግ ማማከር እና የውል ሰነድ ዝግጅት አገልግሎት", "services"),
            ("የኤርኮንዲሽነር ጥገና በቤትዎ መጥተን እንሰራለን", "services"),
            ("Tajaajila suphaa ujummoo bishaanii fi lakkisaa", "services"),
            ("Tajaajila elektiriikaa fi sarara ibsaa manaa", "services"),
            ("Tajaajila suuraa fi viidiyoo sirna cidhaa", "services"),
            ("Tajaajila geejjibaa fi meeshaa manaa geessuu", "services"),
            ("ናይ ሻምቧን ማይን ጽገና ሞያዊ ኣገልግሎት", "services"),
            ("ናይ ኤሌክትሪክ መስመር ምዝርጋሕን ጽገናን", "services"),
            ("ናይ መርዓ ቪድዮን ፎቶን ኣገልግሎት", "services"),
            ("ናይ ኣቑሑት መጓዓዝያ ኣገልግሎት", "services"),
        ],
        "hotels": [
            ("Luxury 5 star hotel suite reservation in Addis Ababa", "hotels"),
            ("Budget guest house single bedroom daily rental Bole", "hotels"),
            ("Resort family bungalow stay near Lake Tana Bahir Dar", "hotels"),
            ("Hotel booking with free breakfast and airport shuttle", "hotels"),
            ("Furnished serviced apartment for weekly short stay", "hotels"),
            ("Lake view hotel room reservation in Hawassa", "hotels"),
            ("Executive conference hall booking with catering", "hotels"),
            ("ባለ 5 ኮከብ የሆቴል አልጋ ክፍል ማስያዝ በአዲስ አበባ", "hotels"),
            ("ለቀን ኪራይ ተስማሚ የቦሌ የእንግዳ ማረፊያ ክፍል", "hotels"),
            ("በጣና ሀይቅ ዳርቻ የሪዞርት ማረፊያ በባህር ዳር", "hotels"),
            ("ቁርስ እና የአየር ማረፊያ ትራንስፖርት የሚያካትት ሆቴል", "hotels"),
            ("ለሳምንት ወይም ለአጭር ጊዜ የተዘጋጀ ሙሉ አፓርታማ", "hotels"),
            ("በሀዋሳ ሀይቅ አጠገብ የሆቴል ክፍል ማስያዝ", "hotels"),
            ("ለስብሰባ እና ኮንፈረንስ አዳራሽ በምግብ አቅርቦት", "hotels"),
            ("Kutaa hoteelaa bareedaa Finfinneetti qabachuu", "hotels"),
            ("Mana keessummaa kiraa guyyaa Booleetti", "hotels"),
            ("Rizoortii haroo Xaanaatti boqonnaaf qabachuu", "hotels"),
            ("Hoteela ciree fi geejjiba buufata xiyyaaraa qabu", "hotels"),
            ("ናይ ሆቴል ክፍሊ ምሓዝ ኣብ ኣዲስ ኣበባ", "hotels"),
            ("ናይ መዓልቲ ክራይ መዕረፊ ኣጋይሽ ኣብ ቦሌ", "hotels"),
            ("ኣብ ባሕሪ ዳር ኣብ ቃላይ ጣና ዝርከብ ሪዞርት", "hotels"),
            ("ናይ ኣኼባ ኣዳራሽን ሆቴልን ምሓዝ", "hotels"),
        ],
        "tickets": [
            ("Live music concert VIP ticket Ethiopian New Year festival", "tickets"),
            ("National football team stadium match entrance ticket", "tickets"),
            ("Cinema movie ticket premiere screening Saturday night", "tickets"),
            ("Cultural theater drama show ticket at National Theatre", "tickets"),
            ("Domestic flight ticket Addis Ababa to Mekelle one way", "tickets"),
            ("International tech exhibition entry pass ticket", "tickets"),
            ("Charity gala dinner ticket reservation VIP table", "tickets"),
            ("የሙዚቃ ኮንሰርት ቪአይፒ መግቢያ ትኬት ለአዲስ አመት", "tickets"),
            ("የኢትዮጵያ ብሄራዊ ቡድን የእግር ኳስ ስታዲየም መግቢያ ትኬት", "tickets"),
            ("የሲኒማ ፊልም መግቢያ ትኬት ለቅዳሜ ምሽት", "tickets"),
            ("በብሄራዊ ቴአትር የባህል ቴአትር ትርኢት ትኬት", "tickets"),
            ("የአገር ውስጥ የበረራ ትኬት ከአዲስ አበባ ወደ መቀሌ", "tickets"),
            ("የቴክኖሎጂ ኤግዚቢሽን መግቢያ ፓስ ትኬት", "tickets"),
            ("የበጎ አድራጎት የእራት ግብዣ የክብር ጠረጴዛ ትኬት", "tickets"),
            ("Tikaattii konsartii muuziqaa ayyaana bara haaraa", "tickets"),
            ("Tikaattii tapha kubbaa miilaa istaadiyeemii", "tickets"),
            ("Tikaattii fiilmii sinimaa fi tiyaatiraa", "tickets"),
            ("Tikaattii xiyyaaraa Finfinnee gara Maqaleetti", "tickets"),
            ("ናይ ሙዚቃ ኮንሰርት መእተዊ ትኬት ሓድሽ ዓመት", "tickets"),
            ("ናይ ኩዕሶ እግሪ ስታድየም መእተዊ ትኬት", "tickets"),
            ("ናይ ፊልምን ሲኒማን መእተዊ ትኬት", "tickets"),
            ("ናይ ነፋሪት ትኬት ካብ ኣዲስ ኣበባ ናብ መቐለ", "tickets"),
        ],
    }

    for items in cat_items.values():
        for text, label in items:
            dataset.append((text, label))
            # Variations
            dataset.append((f"Where can I find {text}?", label))
            dataset.append((f"Show me available {text}", label))
            dataset.append((f"I am looking for {text}", label))

    return dataset


# -------------------------------------------------------------
# 4. SENTIMENT & DISPUTE DATASET (4 Classes)
# -------------------------------------------------------------
def build_sentiment_dispute_dataset():
    dataset = []
    bases = {
        "positive_review": [
            "Excellent seller! The phone arrived in mint condition exactly as described.",
            "Fast delivery, very honest dealer, highly recommended on HUNTIQ!",
            "Great experience renting this apartment, landlord was super helpful.",
            "The car runs perfectly, smooth transaction and genuine documents.",
            "Top quality teff grains, clean and well packaged. Five stars!",
            "Outstanding service! Will definitely buy from this trusted seller again.",
            "Authentic laptop with valid serial number, works flawlessly.",
            "Very polite and professional seller, gave me a great discount.",
            "Delivered right on time to Hawassa, item in pristine condition.",
            "Five stars review! Trusted platform transaction and verified seller.",
            "እጅግ በጣም ምርጥ ሻጭ! ስልኩ ልክ እንደተገለጸው ንጹህ ሆኖ ደርሶኛል።",
            "ፈጣን አቅርቦት እና ታማኝ ነጋዴ፣ በHUNTIQ ላይ በጣም እመክራቸዋለሁ!",
            "መኪናው ምንም እንከን የለውም፣ ሰነዶቹም ትክክለኛ ናቸው።",
            "ንጹህ ማኛ ጤፍ፣ ጥራቱ ከፍተኛ ነው። እናመሰግናለን!",
            "አስተማማኝ እና ታማኝ ሻጭ፣ እቃው በሰዓቱ ደርሶኛል እናመሰግናለን።",
            "ጥሩ መስተንግዶ እና ቅናሽ አድርገውልኛል፣ ባለ አምስት ኮከብ ደረጃ ሰጥቻለሁ።",
            "ላፕቶፑ ኦሪጅናል ነው፣ ምንም አይነት እንከን የለበትም።",
            "Gurguraa baay'ee gaarii! Bilbilli akkuma jedhame qulqulluudha.",
            "Tajaajila ariifataa fi amanamaa, galatoomaa!",
            "Konkolaataan kun baay'ee bareedaadha, nagaa guutuu qaba.",
            "Xaafiin dhufe baay'ee qulqulluu dha, ulfaadhaa!",
            "Oomisha qulqulluu fi gurguraa amanamaa argadhe.",
            "ብጣዕሚ ብሉጽ ሸያጢ! እቲ ስልኪ ጽሩይ ኮይኑ በጺሑኒ።",
            "እሙን ነጋዳይ፣ ብጣዕሚ የመስግነኩም!",
            "መኪና ጽብቕቲ እያ፣ ኩሉ ሰነዳታ ድማ ትክክል እዩ።",
            "ጽሩይ ጣፍ በጺሑኒ፣ ሓሙሽተ ኮከብ ሂበዮም ኣለኹ።",
            "እቲ ንብረት ጽቡቕ እዩ ቀልጢፉ በፂሑኒ የቐንየለይ",
            "ብጣዕሚ ጽቡቕ ኣቕሓ፣ ቀልጢፉ በጺሑኒ የቐንየለይ",
            "The laptop works perfectly, clean condition and quick response!",
            "Meeshaan baay'ee gaariidha, baay'ee na gammachiise",
            "እቃው በጣም ቆንጆ ነው እና ሻጩ ታማኝ ነው አመሰግናለሁ",
        ],
        "neutral_inquiry": [
            "Can you tell me if this laptop has international warranty?",
            "What are the office hours to view the apartment in Bole?",
            "Is the price inclusive of DARO document authentication fee?",
            "Does the generator come with extra fuel filter spare parts?",
            "What is the return policy if the item does not fit my size?",
            "How many kilometers has the vehicle driven according to the odometer?",
            "Does this smartphone support dual SIM and 5G network bands?",
            "Are utility bills for water and electricity included in the monthly rent?",
            "Can you provide additional close-up photos of the car interior?",
            "What payment channels are accepted: Telebirr, CBE Birr, or cash?",
            "Is this apartment still available for rent this weekend?",
            "የቤት እቃው አለ ወይስ ተሽጧል? ዋጋው ስንት ነው?",
            "Makiinaan kun amma yoonaa gurguraadhaaf jiraa?",
            "እዚ ገዛ እዚ ንክራይ ኣሎ ድዩ ዋጋኡ ክንደይ እዩ?",
            "የላፕቶፑ አለም አቀፍ ዋስትና ምን ያህል ጊዜ ነው?",
            "በቦሌ የሚገኘውን ቤት ለማየት የተመቹ ሰዓቶች መቼ ናቸው?",
            "ዋጋው የውልና ማስረጃ ማረጋገጫ ክፍያን ያካትታል?",
            "እቃው የማይመጥነኝ ከሆነ የመመለሻ ህጉ ምን ይመስላል?",
            "መኪናው እስካሁን ምን ያህል ኪሎሜትር ተጉዟል?",
            "ስልኩ ሁለት ሲም ካርድ እና የ5ጂ ኔትወርክ ይሰራል?",
            "የኤሌክትሪክ እና የውሃ ክፍያ በወርሃዊ ኪራዩ ውስጥ ተካቷል?",
            "Waa'ee kaffaltii fi heeyyama konkolaataa naaf ibsaa?",
            "Sa'aatii kamitti suuqii keessan dhufnee ilaaluu dandeenya?",
            "Bilbilli kun kaardii simii lama ni fudhataa?",
            "Gatiin kun kaffaltii qulqulleessa heeyyamaa ni dabalataa?",
            "Gatiin bishaanii fi ibsaa kiraa keessatti qabamaa?",
            "ናይዚ ላፕቶፕ ዋሕስ ክንደይ ዓመት እዩ?",
            "ነቲ ገዛ ንምርኣይ ዝጥዕም ሰዓት መዓስ እዩ?",
            "እቲ ስልኪ ክልተ ሲም ካርድ ይሰርሕ ድዩ?",
            "ናይ ማይን መብራህትን ክፍሊት ኣብቲ ክራይ ኣትዩ ድዩ?",
        ],
        "negative_dissatisfied": [
            "The product arrived with scratches and battery dies within an hour.",
            "Seller was 3 hours late to the meeting and refused to answer calls.",
            "The apartment was completely dirty and plumbing had constant leaks.",
            "Poor customer support, waited 2 days for a simple reply.",
            "The shoes are definitely fake replicas, not original Nike at all.",
            "Very bad experience: screen had cracked glass not shown in photos.",
            "Unprofessional behavior from seller, rude and dismissive tone.",
            "The vehicle makes strange engine knocking noises on test drive.",
            "Delivery was delayed by a full week without any explanation.",
            "Disappointed with the quality, does not match the advertised description.",
            "Product stopped working after 2 hours and seller refuses to answer",
            "እቃው ተበላሽቷል ሻጩም ስልክ አያነሳም ቅር ተሰኝቻለሁ",
            "Meeshaan cabeera garuu gurguraan bilbila hin kaasu",
            "እቲ ቴሌቪዥን ተበላሽዩ ሸያጢ ድማ ስልኪ ኣየልዕልን ኣሎ",
            "እቃው ተቧጭሮ ደርሶኛል፣ ባትሪውም በአንድ ሰዓት ውስጥ ያልቃል።",
            "ሻጩ ለስብሰባ 3 ሰዓት አርፍዶ ስልክም አያነሳም ነበር።",
            "ቤቱ በጣም ቆሻሻ ነው፣ የቧንቧ ፍሳሽም ያፈሳል።",
            "ጫማው ኮፒ እንጂ ኦሪጅናል ናይኪ አይደለም፣ አዝኛለሁ።",
            "የስልኩ ስክሪን ተሰብሯል፣ በፎቶው ላይ ግን ንጹህ ይመስል ነበር።",
            "ሻጩ በጣም ባለጌ እና አክብሮት የሌለው ሰው ነው፣ አላስደሰተኝም።",
            "መኪናው የሞተር ድምጽ ያሰማል፣ ጉድለቱ ተደብቆብኝ ነበር።",
            "እቃው በሳምንት ዘግይቶ ነው የደረሰው፣ አገልግሎታቸው ደካማ ነው።",
            "Meeshaan dhufe miidhamaa fi cabaa dha, hin jaalanne.",
            "Gurguraan yeroo qabame hin kabajne, bilbila hin kaasu.",
            "Suuraa irratti kan mul'ate fi kan dhufe wal hin simne.",
            "Tajaajilli keessan baay'ee laafaadha, guyyaa hedduu ture.",
            "እቲ ኣቕሓ ተሰይሩ በጺሑኒ፣ ባትሪ እውን ቀልጢፉ ይጠፍእ።",
            "እቲ ጫማ ናይ ሓሶት እዩ፣ ኦሪጅናል ኣይኮነን።",
            "እቲ ገዛ ጽሩይ ኣይኮነን፣ ናይ ማይ ሻምቧ ይፈስስ ኣሎ።",
            "እቲ ሸያጢ ብጣዕሚ ሕሱም እዩ፣ ኣይተረድኣንን።",
        ],
        "transaction_dispute": [
            "Seller stole my money via Telebirr without delivering, I want full refund or police action!",
            "ብር በቴሌብር ከፍዬ እቃውን አልላከም አጭበርብሮኛል በፍጥነት ብሬ ይመለስልኝ",
            "Maallaqa kaffalee meeshaa naaf hin ergine, kun hatumsa kaffaltiin koo haa deebi'u",
            "ገንዘብ ወሲዱ ንብረት ኣይሃበንን ሰሪቑኒ ገንዘበይ ይመለሰለይ ወይ ናብ ፖሊስ ክኸይድ",
            "I sent 10,000 birr payment but the seller blocked me and disappeared!",
            "Urgent fraud alert: seller provided fake invalid Libre title document!",
            "Dispute claim: received an empty box instead of iPhone 14 Pro Max.",
            "Seller took deposit and rented the same apartment to another tenant!",
            "I demand immediate refund escrow intervention, item was completely wrong.",
            "Official complaint: unauthorized transaction charge on my Telebirr account.",
            "Legal notice: seller refuses to return deposit after canceling sale.",
            "Escrow dispute opened: seller delivered defective damaged television set.",
            "Fraudulent listing: seller does not own the property listed for sale.",
            "I request formal mediation and police report documentation from HUNTIQ.",
            "10 ሺህ ብር ክፍያ ከላክሁ በኋላ ሻጩ ብሎክ አድርጎኝ ተሰወረ! እርዳታ እፈልጋለሁ።",
            "አስቸኳይ የማጭበርበር ክስ፡ ሻጩ የሰጠኝ የመኪና ሊብሬ ሰነድ ሀሰተኛ ነው!",
            "የHUNTIQ የደህንነት ቡድን ይግባኝ፡ የአይፎን ሳጥኑ ባዶ ሆኖ ደርሶኛል!",
            "ሻጩ የቅድመ ክፍያ ተቀብሎ ቤቱን ለሌላ ሰው አከራየው፣ ገንዘቤ ይመለስ!",
            "አስቸኳይ ይግባኝ፡ የተሰበረ ቴሌቪዥን ልከውልኛል፣ ክፍያዬ በአስቸኳይ ይመለስልኝ።",
            "ሻጩ ውሉን አፍርሶ ገንዘቤን አልመልስም ብሏል፣ የህግ እርዳታ እፈልጋለሁ።",
            "የማጭበርበር ወንጀል ተፈጽሞብኛል፣ ለፖሊስ ሪፖርት ማድረጊያ ሰነድ ስጡኝ።",
            "Maallaqa qarshii 10000 ergeen booda gurguraan na cufe, na gargaaraa!",
            "Heeyyami konkolaataa sobaa waan ta'eef maallaqni koo haa deebi'u!",
            "Gurguraan qabiyyee mana kireesse nama biraatti gurgure, himata qaba!",
            "Maallaqa koo telebirr irraa hatan, dafaa naaf deebisaa!",
            "HUNTIQ seeratti naaf haa dhiyeessu gurguraa kana!",
            "10000 ቅርሺ ምስ ሰደድኩ እቲ ሸያጢ ብሎክ ገይሩኒ፣ ምትላል እዩ ገንዘበይ ይመለስ!",
            "ሰነድ መኪና ናይ ሓሶት እዩ፣ ንHUNTIQ ክስ ከእቱ እየ!",
            "ባዶ ሳንዱቕ ስልኪ ልኢኾምለይ፣ ገንዘበይ ሕጂ ይመለሰለይ!",
            "ፍትሒ እደሊ ኣለኹ፣ እዚ ገበን እዩ ንፖሊስ ከረክቦ እየ!",
        ],
    }

    for label, items in bases.items():
        for text in items:
            dataset.append((text, label))
            dataset.append((f"Review statement: {text}", label))
            dataset.append((f"Customer feedback ticket: {text}", label))

    return dataset


if __name__ == "__main__":
    d_dataset = build_domain_dataset()
    s_dataset = build_safety_dataset()
    c_dataset = build_category_dataset()
    sd_dataset = build_sentiment_dispute_dataset()

    print(f"Domain dataset: {len(d_dataset)} samples")
    print(f"Safety dataset: {len(s_dataset)} samples")
    print(f"Category dataset: {len(c_dataset)} samples")
    print(f"Sentiment dispute dataset: {len(sd_dataset)} samples")

    for name, data in [
        ("domain_dataset.json", d_dataset),
        ("safety_dataset.json", s_dataset),
        ("category_dataset.json", c_dataset),
        ("sentiment_dispute_dataset.json", sd_dataset),
    ]:
        p = os.path.join(DATA_DIR, name)
        with open(p, "w", encoding="utf-8") as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
        print(f"Saved {name} ({len(data)} items)")
