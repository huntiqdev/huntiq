package service

import (
	"regexp"
	"strings"

	"github.com/your-org/used-tools-marketplace/internal/service/aiclient"
)

var (
	rxCoding  = regexp.MustCompile(`(?i)\b(python|javascript|typescript|golang|c\+\+|java|html|css|sql|rust|php|bash script|powershell)\b|\b(reverse a string|binary search|fibonacci|bubble sort|leetcode)\b|\b(write|create|debug|generate)\b.*\b(code|function|program|script|regex|class|loop)\b`)
	rxScience = regexp.MustCompile(`(?i)\b(equation|derivative|integral|calculus|physics|quantum|chemistry|biology|homework|photosynthesis|mitochondria|thermodynamics|relativity|algebra)\b`)
	rxTrivia  = regexp.MustCompile(`(?i)\b(who wrote|capital of|tell me a.*story|write a.*poem|who is the president|who won the)\b`)
)

// IsOutOfDomain checks if the prompt is an unrelated general topic like coding, math, physics, or trivia.
func IsOutOfDomain(query string) bool {
	q := strings.TrimSpace(query)
	if q == "" {
		return false
	}

	// Always allow legitimate product questions even if they contain numbers, prices, or technical specs
	lower := strings.ToLower(q)
	if strings.Contains(lower, "phone") || strings.Contains(lower, "laptop") || strings.Contains(lower, "car") ||
		strings.Contains(lower, "apartment") || strings.Contains(lower, "house") || strings.Contains(lower, "rent") ||
		strings.Contains(lower, "price") || strings.Contains(lower, "cost") || strings.Contains(lower, "huntiq") ||
		strings.Contains(lower, "buy") || strings.Contains(lower, "sell") || strings.Contains(lower, "etb") ||
		strings.Contains(lower, "birr") || strings.Contains(lower, "specs") || strings.Contains(lower, "compare") {
		// If it's explicitly asking to write code or solve physics equations, it's still out of domain
		if rxCoding.MatchString(q) || rxScience.MatchString(q) {
			return true
		}
		return false
	}

	if rxCoding.MatchString(q) || rxScience.MatchString(q) || rxTrivia.MatchString(q) {
		return true
	}

	return false
}

// GetOutOfDomainRedirection returns a polite, concise marketplace-focused redirection in the requested locale.
func GetOutOfDomainRedirection(locale string) string {
	switch strings.ToLower(locale) {
	case "am":
		return "እኔ የHUNTIQ የገበያ ረዳት ነኝ። በኢትዮጵያ ውስጥ ያሉ ምርቶችን፣ መኪናዎችን፣ የቤትና አፓርታማ ኪራዮችን እና አገልግሎቶችን እንዲያገኙ ወይም ስለ HUNTIQ አጠቃቀም መረጃ ለመስጠት ዝግጁ ነኝ። ለትምህርት፣ ለፊዚክስ ወይም ለኮዲንግ ጥያቄዎች እባክዎን አጠቃላይ መድረኮችን ይጠቀሙ።"
	case "om":
		return "Ani gargaaraa gabaa HUNTIQ ti. Itoophiyaa keessatti oomishaalee, konkolaataa, kireeffannaa manaa fi tajaajiloota akka argattu ykn waa'ee HUNTIQ akka beektuuf nan gargaara. Gaaffilee barnootaa, koodii fi fiiziksiif maddoota biraa fayyadamaa."
	case "ti":
		return "ኣነ ናይ HUNTIQ ናይ ዕዳጋ ሓጋዚ እየ። ኣብ ኢትዮጵያ ፍርያት፣ መካይን፣ ናይ ገዛ ክራይን ካልኦት ኣገልግሎታትን ንምርካብ ወይ ብዛዕባ HUNTIQ ንምሕታት ክሕግዘኩም እኽእል እየ። ንትምህርታዊ፣ ፊዚክስ ወይ ኮዲንግ ሕቶታት ካልእ መድረኽ ተጠቐሙ።"
	default:
		return "I am HUNTIQ's marketplace assistant. I specialize in helping you find verified products, vehicles, property rentals, and services across Ethiopia, as well as answering questions about using HUNTIQ. For general coding, physics, or academic homework, please consult general-purpose educational tools."
	}
}

// CheckPlatformFAQ checks if the user is asking about HUNTIQ features, registration, listing rules, or fees.
func CheckPlatformFAQ(query, locale string) (bool, string) {
	lower := strings.ToLower(strings.TrimSpace(query))

	// How to post / sell / list
	if strings.Contains(lower, "how to post") || strings.Contains(lower, "how do i post") ||
		strings.Contains(lower, "how to sell") || strings.Contains(lower, "how to list") ||
		strings.Contains(lower, "post an apartment") || strings.Contains(lower, "post a product") ||
		strings.Contains(lower, "እንዴት መለጠፍ") || strings.Contains(lower, "እንዴት መሸጥ") ||
		strings.Contains(lower, "akkamitti maxxansu") || strings.Contains(lower, "akkamitti gurguru") ||
		strings.Contains(lower, "ከመይ ጌረ ይልጥፍ") {

		switch locale {
		case "am":
			return true, "በHUNTIQ ላይ ምርት ወይም የቤት ኪራይ ለመለጠፍ የሚከተሉትን ቀላል ደረጃዎች ይከተሉ:\n\n1. ከላይ በሚገኘው ሜኑ ላይ **'ማስታወቂያ ልጠፍ'** (Post Listing) የሚለውን ይጫኑ።\n2. የሻጭ/አከራይ መለያ ከሌለዎት በስልክ ቁጥርዎ በነጻ ይመዝገቡና በOTP ያረጋግጡ።\n3. ተገቢውን ምድብ (ለምሳሌ፡ Real Estate > Apartments) እና ዓይነት (ሽያጭ ወይም ኪራይ) ይምረጡ።\n4. ዋጋ፣ ክልል/ክፍለ ከተማ፣ ማብራሪያና ጥራት ያላቸው ፎቶዎችን ያስገቡ።\n5. ያስገቡትን መረጃ ሲያጸድቁ ማስታወቂያዎ በቀጥታ በመላው ኢትዮጵያ ለገዢዎች ይታያል!"
		case "om":
			return true, "HUNTIQ irratti oomisha ykn kireeffannaa manaa maxxansuuf tarkaanfiiwwan kanneen hordofaa:\n\n1. Fuula gubbaarraa **'Post Listing'** jedhu cuqaasaa.\n2. Lakkoofsa bilbilaa keessaniin galmaa'aa OTP tiin mirkaneessaa.\n3. Ramaddii fi gosa gurgurtaa (Gurgurtaa ykn Kireeffannaa) filadhaa.\n4. Gatii, bakka (Magaalaa/Kutaa Magaalaa), ibsa fi suuraalee qulqulluu galchaa.\n5. Erga mirkanaa'ee booda tarreeffamni keessan guutuu Itoophiyaatti ni mul'ata!"
		case "ti":
			return true, "ኣብ HUNTIQ ፍርያት ወይ ናይ ገዛ ክራይ ንምልጣፍ ዝስዕቡ ቅደም ሰዓባት ተኸተሉ:\n\n1. ካብ ላዕለዋይ መዐደሊ **'Post Listing'** ዝብል ጠውቑ።\n2. ብቑጽሪ ስልክኹም ተመዝጊብኩም ብOTP ኣረጋግጹ።\n3. ምድብን ዓይነትን (ሽያጥ ወይ ክራይ) ምረጹ።\n4. ዋጋ፣ ቦታ፣ መግለጺን ፎቶታትን ብምእታው ኣጽድቑ።\n5. ንጡፍ ኮይኑ ኣብ መላእ ኢትዮጵያ ንተጠቃምቲ ይርአ!"
		default:
			return true, "To post a product, vehicle, or rental on HUNTIQ:\n\n1. Click **'Post Listing'** on the top navigation bar.\n2. If you don't have a seller account, register with your phone number and verify via OTP (free for standard listings).\n3. Choose your category (e.g. Real Estate > Apartments) and listing type (**'For Rent'** or **'For Sale'**).\n4. Enter the price, city/subcity, detailed specifications, and upload clear photos.\n5. Submit your listing. Once verified, it goes live across Ethiopia instantly!"
		}
	}

	// How to register / create account
	if strings.Contains(lower, "how to register") || strings.Contains(lower, "how to create account") ||
		strings.Contains(lower, "sign up") || strings.Contains(lower, "እንዴት ልመዝገብ") ||
		strings.Contains(lower, "akkamitti galmaa'u") || strings.Contains(lower, "ከመይ ጌረ ይመዝገብ") {

		switch locale {
		case "am":
			return true, "በHUNTIQ ላይ ለመመዝገብ:\n\n1. ከላይ **'ይመዝገቡ'** (Register) የሚለውን ይጫኑ።\n2. የኢትዮጵያ ስልክ ቁጥርዎን (09...) እና የይለፍ ቃል ያስገቡ።\n3. በስልክዎ የሚደርሰውን ባለ 6-አሃዝ የይለፍ ቃል (OTP) በማስገባት መለያዎን ያረጋግጡ።\n4. ገዢ፣ ሻጭ ወይም ሁለቱንም መሆን ይችላሉ። መመዝገብ ሙሉ በሙሉ ነፃ ነው!"
		default:
			return true, "To create an account on HUNTIQ:\n\n1. Click **'Register'** on the top right navigation.\n2. Enter your Ethiopian phone number (e.g., `09...`) and choose a secure password.\n3. Enter the 6-digit verification code (OTP) sent to your phone.\n4. You can freely switch between buyer and seller capabilities anytime. Registration is 100% free!"
		}
	}

	// Membership tiers / pricing
	if strings.Contains(lower, "membership") || strings.Contains(lower, "subscription") ||
		strings.Contains(lower, "tier") || strings.Contains(lower, "premium plan") || strings.Contains(lower, "gold plan") {
		return true, "HUNTIQ offers three seller trust tiers:\n\n• **Standard / Normal (Free):** Post verified listings with basic search discovery.\n• **Gold Certified:** Enhanced search ranking, verified trust badge, and priority customer matching.\n• **Premium Verified:** Top search placement, full buyer-protection assurance, and featured hero display.\n\nBuyers and guest visitors can browse, search, and use the AI assistant completely free without any subscription."
	}

	return false, ""
}

// IsBroadProductSearch checks if the user is asking broadly for available products rather than a specific item.
func IsBroadProductSearch(query string) bool {
	lower := strings.ToLower(strings.TrimSpace(query))
	broadPhrases := []string{
		"what products are available now",
		"what products are available",
		"what are available now",
		"what is available now",
		"what is available",
		"what do you have",
		"show all products",
		"show me all products",
		"show all listings",
		"list all products",
		"all products",
		"available products",
		"ምን እቃዎች አሉ",
		"ምን ምርቶች አሉ",
		"ያሉ እቃዎች",
		"oomishaalee jiran",
		"maal qabdu",
		"እንታይ ፍርያት ኣለኩም",
		"ዘለዉ ፍርያት",
	}

	for _, phrase := range broadPhrases {
		if strings.Contains(lower, phrase) || lower == phrase {
			return true
		}
	}

	// Also check if cleaning strips all keywords leaving nothing
	clean := aiclient.CleanConversationalSearchTerm(query)
	return clean == "" || clean == "all" || clean == "products" || clean == "product" || clean == "items" || clean == "item"
}
