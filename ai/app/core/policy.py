"""Platform FAQ and domain boundary policies for HUNTIQ AI."""

import re

# Domain boundary regexes
RX_CODING = re.compile(
    r"(?i)\b(python|javascript|typescript|golang|c\+\+|java|html|css|sql|rust|php|bash script|powershell)\b|"
    r"\b(reverse a string|binary search|fibonacci|bubble sort|leetcode)\b|"
    r"\b(write|create|debug|generate)\b.*\b(code|function|program|script|regex|class|loop)\b"
)
RX_SCIENCE = re.compile(
    r"(?i)\b(quadratic|formula|equation|derivative|integral|calculus|physics|quantum|chemistry|biology|homework|"
    r"photosynthesis|mitochondria|thermodynamics|relativity|algebra|fiiziksii|herrega)\b|"
    r"(ፊዚክስ|ኳንተም|ቀመር|ኬሚስትሪ|ባዮሎጂ|ሂሳብ|ስሌት|ተዋጽኦ|ኢንቴግራል)"
)
RX_TRIVIA = re.compile(
    r"(?i)\b(who wrote|capital of|tell me a.*story|write a.*poem|who is the president|who won the)\b"
)

# Legitimate marketplace product and domain keywords
MARKETPLACE_TOKENS = {
    "phone",
    "laptop",
    "car",
    "apartment",
    "house",
    "rent",
    "price",
    "cost",
    "huntiq",
    "buy",
    "sell",
    "etb",
    "birr",
    "specs",
    "compare",
    "drill",
    "tools",
    "ስልክ",
    "ሞባይል",
    "መኪና",
    "ቤት",
    "አፓርታማ",
    "ኮንዶሚኒየም",
    "ኮምፒውተር",
    "ላፕቶፕ",
    "ጫማ",
    "ልብስ",
    "ቲቪ",
    "ቴሌቪዥን",
    "ፈርኒቸር",
    "ዋጋ",
    "ብር",
    "ኪራይ",
    "ሽያጭ",
    "እቃ",
    "ምርት",
    "bilbila",
    "konkolaataa",
    "mana",
    "kompiitara",
    "kiraa",
    "gatii",
    "meeshaa",
    "bituu",
    "ገዛ",
    "ክራይ",
    "ንብረት",
    "ክገዝእ",
    "ቅርሺ",
    "ቦሌ",
}


def is_out_of_domain(query: str) -> bool:
    """Check if the user request is outside the HUNTIQ marketplace scope."""
    q = query.strip()
    if not q:
        return False

    lower = q.lower()
    has_market_token = any(token in lower for token in MARKETPLACE_TOKENS)

    if has_market_token:
        # If user explicitly asks to write code or solve physics equations, it's still out of domain
        return bool(RX_CODING.search(q) or RX_SCIENCE.search(q))

    return bool(RX_CODING.search(q) or RX_SCIENCE.search(q) or RX_TRIVIA.search(q))


def get_out_of_domain_redirection(locale: str) -> str:
    """Return a polite, domain-specialized redirection in the target locale."""
    loc = (locale or "en").lower()
    if loc in ("om", "or"):
        return (
            "Ani gargaaraa gabaa HUNTIQ ti. Itoophiyaa keessatti oomishaalee, konkolaataa, "
            "kireeffannaa manaa fi tajaajiloota akka argattu ykn waa'ee HUNTIQ akka beektuuf "
            "nan gargaara. Gaaffilee barnootaa, koodii fi fiiziksiif maddoota biraa fayyadamaa."
        )
    if loc == "am":
        return (
            "እኔ የHUNTIQ የገበያ ረዳት ነኝ። በኢትዮጵያ ውስጥ ያሉ ምርቶችን፣ መኪናዎችን፣ የቤትና አፓርታማ ኪራዮችን "
            "እና አገልግሎቶችን እንዲያገኙ ወይም ስለ HUNTIQ አጠቃቀም መረጃ ለመስጠት ዝግጁ ነኝ። ለትምህርት፣ "
            "ለፊዚክስ ወይም ለኮዲንግ ጥያቄዎች እባክዎን አጠቃላይ መድረኮችን ይጠቀሙ።"
        )
    if loc == "ti":
        return (
            "ኣነ ናይ HUNTIQ ናይ ዕዳጋ ሓጋዚ እየ። ኣብ ኢትዮጵያ ፍርያት፣ መካይን፣ ናይ ገዛ ክራይን ካልኦት "
            "ኣገልግሎታትን ንምርካብ ወይ ብዛዕባ HUNTIQ ንምሕታት ክሕግዘኩም እኽእል እየ። ንትምህርታዊ፣ "
            "ፊዚክስ ወይ ኮዲንግ ሕቶታት ካልእ መድረኽ ተጠቐሙ።"
        )
    return (
        "I am HUNTIQ's marketplace assistant. I specialize in helping you find verified "
        "products, vehicles, property rentals, and services across Ethiopia, as well as "
        "answering questions about using HUNTIQ. For general coding, physics, or academic "
        "homework, please consult general-purpose educational tools."
    )


def check_platform_faq(query: str, locale: str) -> tuple[bool, str]:
    """Check if the user is asking about platform usage, listing rules, verification, or fees."""
    lower = query.strip().lower()
    loc = (locale or "en").lower()
    if loc == "or":
        loc = "om"

    # 1. How to post / publish a listing
    if any(
        k in lower
        for k in [
            "how to post",
            "how do i post",
            "how to sell",
            "how to list",
            "how do i publish",
            "how to publish",
            "publish a listing",
            "publish an ad",
            "post an apartment",
            "post a product",
            "እንዴት መለጠፍ",
            "እንዴት መሸጥ",
            "akkamitti maxxansu",
            "akkamitti gurguru",
            "ከመይ ጌረ ይልጥፍ",
        ]
    ):
        if loc == "am":
            return True, (
                "በHUNTIQ ላይ ምርት ወይም የቤት ኪራይ ለመለጠፍ የሚከተሉትን ቀላል ደረጃዎች ይከተሉ:\n\n"
                "1. ከላይ በሚገኘው ሜኑ ላይ **'ማስታወቂያ ልጠፍ'** (Post Listing) የሚለውን ይጫኑ።\n"
                "2. የሻጭ/አከራይ መለያ ከሌለዎት በስልክ ቁጥርዎ በነጻ ይመዝገቡና በOTP ያረጋግጡ።\n"
                "3. ተገቢውን ምድብ (ለምሳሌ፡ Real Estate > Apartments) እና ዓይነት (ሽያጭ ወይም ኪራይ) ይምረጡ።\n"
                "4. ዋጋ፣ ክልል/ክፍለ ከተማ፣ ማብራሪያና ጥራት ያላቸው ፎቶዎችን ያስገቡ።\n"
                "5. ያስገቡትን መረጃ ሲያጸድቁ ማስታወቂያዎ በቀጥታ በመላው ኢትዮጵያ ለገዢዎች ይታያል!"
            )
        if loc == "om":
            return True, (
                "HUNTIQ irratti oomisha ykn kireeffannaa manaa maxxansuuf tarkaanfiiwwan kanneen hordofaa:\n\n"
                "1. Fuula gubbaarraa **'Post Listing'** jedhu cuqaasaa.\n"
                "2. Lakkoofsa bilbilaa keessaniin galmaa'aa OTP tiin mirkaneessaa.\n"
                "3. Ramaddii fi gosa gurgurtaa (Gurgurtaa ykn Kireeffannaa) filadhaa.\n"
                "4. Gatii, bakka (Magaalaa/Kutaa Magaalaa), ibsa fi suuraalee qulqulluu galchaa.\n"
                "5. Erga mirkanaa'ee booda tarreeffamni keessan guutuu Itoophiyaatti ni mul'ata!"
            )
        if loc == "ti":
            return True, (
                "ኣብ HUNTIQ ፍርያት ወይ ናይ ገዛ ክራይ ንምልጣፍ ዝስዕቡ ቅደም ሰዓባት ተኸተሉ:\n\n"
                "1. ካብ ላዕለዋይ መዐደሊ **'Post Listing'** ዝብል ጠውቑ።\n"
                "2. ብቑጽሪ ስልክኹም ተመዝጊብኩም ብOTP ኣረጋግጹ።\n"
                "3. ምድብን ዓይነትን (ሽያጥ ወይ ክራይ) ምረጹ።\n"
                "4. ዋጋ፣ ቦታ፣ መግለጺን ፎቶታትን ብምእታው ኣጽድቑ።\n"
                "5. ንጡፍ ኮይኑ ኣብ መላእ ኢትዮጵያ ንተጠቃምቲ ይርአ!"
            )
        return True, (
            "To post a product, vehicle, or rental on HUNTIQ:\n\n"
            "1. Click **'Post Listing'** on the top navigation bar.\n"
            "2. If you don't have a seller account, register with your phone number and verify via OTP (free for standard listings).\n"
            "3. Choose your category (e.g. Real Estate > Apartments) and listing type (**'For Rent'** or **'For Sale'**).\n"
            "4. Enter the price, city/subcity, detailed specifications, and upload clear photos.\n"
            "5. Submit your listing. Once verified, it goes live across Ethiopia instantly!"
        )

    # 2. How to register / create account
    if any(
        k in lower
        for k in [
            "how to register",
            "how to create account",
            "sign up",
            "create an account",
            "እንዴት ልመዝገብ",
            "መመዝገብ እፈልጋለሁ",
            "akkamitti galmaa'u",
            "ከመይ ጌረ ይመዝገብ",
        ]
    ):
        if loc == "am":
            return True, (
                "በHUNTIQ ላይ መለያ ለመክፈት:\n\n"
                "1. ከላይ በቀኝ በኩል **'ግባ / ተመዝገብ'** (Sign In) የሚለውን ይጫኑ።\n"
                "2. የኢትዮጵያ ስልክ ቁጥርዎን ያስገቡ።\n"
                "3. በSMS የሚደርስዎትን ባለ 6 አሃዝ የማረጋገጫ ኮድ (OTP) ያስገቡ።\n"
                "4. መለያዎ ወዲያውኑ ይከፈታል! ገዢም ሆኑ ሻጭ ተጨማሪ ክፍያ ሳያስፈልግ መገበያየት ይችላሉ።"
            )
        if loc == "om":
            return True, (
                "HUNTIQ irratti akkaawuntii banachuuf:\n\n"
                "1. Mirga gubbaarraa **'Sign In / Register'** kan jedhu cuqaasaa.\n"
                "2. Lakkoofsa bilbila Itoophiyaa keessanii galchaa.\n"
                "3. Koodii OTP sarara SMS tiin isiniif dhufu galchuun mirkaneessaa.\n"
                "4. Akkaawuntiin keessan yeruma sana qophii ta'a!"
            )
        if loc == "ti":
            return True, (
                "ኣብ HUNTIQ ኣካውንት ንምኽፋት:\n\n"
                "1. ካብ የማን ላዕሊ **'ተመዝገብ'** ዝብል ጠውቑ።\n"
                "2. ናይ ኢትዮጵያ ቁጽሪ ስልክኹም ኣእትዉ።\n"
                "3. ብSMS ዝመጸኩም ናይ 6 ኣሃዝ መረጋገጺ ኮድ (OTP) ብምእታው ኣጽድቑ።\n"
                "4. ብቕጽበት ይኽፈተልኩም!"
            )
        return True, (
            "To create an account on HUNTIQ:\n\n"
            "1. Click **'Sign In / Register'** in the top navigation bar.\n"
            "2. Enter your Ethiopian mobile phone number (+251...).\n"
            "3. Enter the 6-digit SMS verification code (OTP).\n"
            "4. Your account is ready immediately for buying, renting, or selling!"
        )

    # 3. Vehicle Libre / Ownership verification
    if any(
        k in lower
        for k in [
            "libre",
            "legal papers",
            "car papers",
            "ownership document",
            "ሊብሬ",
            "የመኪና ሰነድ",
            "liibiree",
            "waraqaa abbaa qabeenyummaa",
            "ናይ መኪና ሰነድ",
        ]
    ):
        if loc == "am":
            return True, (
                "በኢትዮጵያ ውስጥ መኪና ከመግዛትዎ በፊት የሚከተሉትን ህጋዊ ሰነዶች ያረጋግጡ:\n\n"
                "1. **ኦሪጅናል ሊብሬ (Libre):** የስም ማዛወሪያ ታሪክና የመኪናው ሻሲ (Chassis) ቁጥር ከሰነዱ ጋር መመሳከሩን ያረጋግጡ።\n"
                "2. **የትራንስፖርት ባለስልጣን ምርመራ:** መኪናው ላይ የብድር እገዳ (Bank Loan/Pledge) ወይም የፍርድ ቤት እገዳ አለመኖሩን ያጣሩ።\n"
                "3. **የሻጭ ህጋዊ ውክልና:** ሻጩ ባለቤቱ ካልሆነ ህጋዊ እና የታደሰ የውክልና ስልጣን (Power of Attorney) ሊኖረው ይገባል።\n"
                "4. ክፍያ የሚፈጸመው ሰነዶች በሙሉ በውልና ማስረጃ ተረጋግጠው ስም ሲዛወር ብቻ መሆን አለበት።"
            )
        return True, (
            "Before purchasing any vehicle in Ethiopia, verify these critical legal safeguards:\n\n"
            "1. **Original Libre (Title Document):** Cross-check engine & chassis VIN numbers against the physical car.\n"
            "2. **Transport Authority Clearance:** Confirm at the Subcity Transport Bureau that there are no bank liens, collateral pledges, or court injunctions.\n"
            "3. **Power of Attorney:** If the seller is not the sole titleholder, ensure they hold a valid, notarized agency document.\n"
            "4. **Secure Settlement:** Never transfer payment before signing verified transfer contracts at the Documents Authentication & Registration Agency (DARO)."
        )

    # 4. Electronics Anti-Fraud / IMEI check
    if any(
        k in lower
        for k in [
            "imei",
            "anti-fraud",
            "scam",
            "safe test",
            "test used electronics",
            "የስልክ ማጭበርበር",
            "የimei ቁጥር",
            "ደህንነት",
        ]
    ):
        if loc == "am":
            return True, (
                "ያገለገሉ ስልኮችን ወይም ኤሌክትሮኒክስ በደህና ለመግዛት:\n\n"
                "1. **IMEI ቁጥርን ያረጋግጡ:** በስልኩ ላይ `*#06#` በመጫን የሚመጣውን IMEI ከስልኩ ጀርባና ከኦሪጅናል ካርቶኑ ጋር ያመሳክሩ።\n"
                "2. **iCloud / Google FRP መለያ:** የሻጩ መለያ ሙሉ በሙሉ መውጣቱን (Sign Out) እና ስልኩ Reset መደረጉን ያረጋግጡ።\n"
                "3. **የባትሪ ጤንነትና ስክሪን:** ኦሪጅናል ስክሪን መሆኑን እና ባትሪው ቶሎ እንደማይወርድ ይፈትሹ።\n"
                "4. ሁልጊዜ በአደባባይ እና ደህንነቱ በተጠበቀ ቦታ ይገናኙ፤ እቃውን በሚገባ ሳይፈትሹ ገንዘብ አያስተላልፉ።"
            )
        return True, (
            "Checklist for safely buying used electronics in Ethiopia:\n\n"
            "1. **Verify IMEI:** Dial `*#06#` on the phone and verify the IMEI matches the box and internal settings.\n"
            "2. **Account Sign-Out:** Ensure Apple iCloud or Google FRP accounts are completely removed before purchase.\n"
            "3. **Hardware Diagnostic:** Test touch screen, cameras, cellular signal, and battery health.\n"
            "4. **Public Meetup:** Always inspect devices in a well-lit public venue before settling payment."
        )

    return False, ""


def get_safety_warning(safety_class: str, locale: str) -> str:
    """Return localized safety warning message for security/policy violations."""
    loc = (locale or "en").lower()
    if loc == "or":
        loc = "om"

    if safety_class == "scam_advance":
        if loc == "am":
            return "⚠️ የደህንነት ማስጠንቀቂያ፡ እቃውን በአካል አይተው ሳይፈትሹ በቴሌብርም ሆነ በባንክ ምንም አይነት የቅድመ ክፍያ (advance payment) አያስተላልፉ! HUNTIQ አስቀድሞ ገንዘብ እንዲልኩ በጭራሽ አይጠይቅዎትም።"
        if loc == "om":
            return "⚠️ Akeekkachiisa Nageenyaa: Meeshicha ijaan argitanii osoo hin qoratin dura maallaqa dursaatiin kaffaltii telebirr ykn baankii hin erginaa! HUNTIQ dursitanii akka kaffaltan hin gaafatu."
        if loc == "ti":
            return "⚠️ ናይ ደሕንነት መጠንቀቕታ፡ ነቲ ኣቕሓ ብኣካል ርኢኹም ከይተመርመረ ዝኾነ ናይ ቅድመ ክፍሊት (advance payment) ብቴሌብር ይኹን ብባንኪ ኣይትስደዱ! HUNTIQ ቅድሚት ክትከፍሉ ኣይሓተኩምን እዩ።"
        return "⚠️ Safety Warning: NEVER send advance payments or deposits via Telebirr/CBE before physically inspecting the item! HUNTIQ never asks for advance payments."

    if safety_class == "phishing_harvesting":
        if loc == "am":
            return "🛡️ የጥበቃ ማስጠንቀቂያ፡ የቴሌብር፣ የባንክ ወይም የHUNTIQ መለያ ሚስጥር ቁጥርዎን (PIN / OTP) ለማንም አያጋሩ። ያልታወቁ የቴሌግራም ወይም ድረ-ገጽ ሊንኮችን አይክፈቱ።"
        if loc == "om":
            return "🛡️ Akeekkachiisa: Koodii iccitii (PIN / OTP) telebirr, baankii ykn HUNTIQ keessanii eenyuufiyyuu hin qoodinaa. Liinkii hin beekamne hin baninaa."
        if loc == "ti":
            return "🛡️ መጠንቀቕታ፡ ሚስጢራዊ ቁጽሪ (PIN / OTP) ቴሌብር፣ ባንኪ ወይ HUNTIQ ንዝኾነ ሰብ ኣይትሃቡ። ዘይተፈልጡ ናይ ቴሌግራም ወይ መርበብ ሊንክታት ኣይትኽፈቱ።"
        return "🛡️ Security Warning: NEVER share your bank PIN, Telebirr OTP, or login passwords with anyone. Beware of external links."

    if safety_class == "prohibited_contraband":
        if loc == "am":
            return "🚫 የተከለከለ እቃ፡ የጦር መሳሪያዎች፣ ጥይት፣ ሀሰተኛ ገንዘብ እና ህገወጥ እቃዎች በHUNTIQ መድረክ ላይ በጥብቅ የተከለከሉ ናቸው።"
        if loc == "om":
            return "🚫 Meeshaa dhowwame: Meeshaan waraanaa, qarshiin sobaa fi meeshaaleen seeraan alaa HUNTIQ irratti gonkumaa dhowwamaadha."
        if loc == "ti":
            return "🚫 ዝተኸልከለ ኣቕሓ፡ ናይ ኲናት ኣጽዋር፣ ናይ ሓሶት ቅርሽን ዘይሕጋዊ ኣቑሑትን ኣብ መድረኽ HUNTIQ ብጥብቂ ዝተኸልከሉ እዮም።"
        return "🚫 Prohibited Items: Unlicensed weapons, counterfeit currency, illicit drugs, and contraband are strictly prohibited on HUNTIQ."

    if safety_class == "harassment_abuse":
        if loc == "am":
            return "⚠️ የማህበረሰብ ደንብ፡ ስድብ፣ ዛቻ እና የስነምግባር ጥሰት በHUNTIQ መድረክ ላይ የተከለከለ ነው። እባክዎን በትህትና ይነጋገሩ።"
        if loc == "om":
            return "⚠️ Seera Hawaasaa: Arraba, doorsisaa fi kabaja dhabuun HUNTIQ irratti hin hayyamamu. Kabajaan haasa'aa."
        if loc == "ti":
            return "⚠️ ናይ ማሕበረሰብ ደንቢ፡ ጸርፊ፣ ምፍርራህን ሕድገት ዘይብሉ ተግባራትን ኣብ HUNTIQ ክልኩል እዩ። ብትሕትና ተዘራረቡ።"
        return "⚠️ Community Guidelines: Harassment, abusive language, and threats violate HUNTIQ community standards."

    return "⚠️ Please follow HUNTIQ safe trading practices."
