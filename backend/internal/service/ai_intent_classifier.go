package service

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/your-org/used-tools-marketplace/internal/service/aiclient"
)

// IntentType represents a distinct user intention within a message.
type IntentType string

const (
	IntentSocialGreeting      IntentType = "social_greeting"
	IntentSocialFarewell      IntentType = "social_farewell"
	IntentSocialGratitude     IntentType = "social_gratitude"
	IntentSocialSmalltalk     IntentType = "social_smalltalk"
	IntentIdentityInquiry     IntentType = "identity_inquiry"
	IntentCapabilityInquiry   IntentType = "capability_inquiry"
	IntentPlatformFAQ         IntentType = "platform_faq"
	IntentMarketplaceSearch   IntentType = "marketplace_search"
	IntentMarketplaceStats    IntentType = "marketplace_stats"
	IntentMarketplaceQuestion IntentType = "marketplace_question"
	IntentClarificationNeeded IntentType = "clarification_needed"
	IntentOutOfDomain         IntentType = "out_of_domain"
	IntentGeneralKnowledge    IntentType = "general_knowledge"
)

// ExtractedSearchRequest represents a single, independent category query.
type ExtractedSearchRequest struct {
	CategoryName string   `json:"category_name"` // "Mobile Phones", "Computers", "Houses", etc.
	CategorySlug string   `json:"category_slug,omitempty"`
	SearchTerm   string   `json:"search_term"`
	ListingType  *string  `json:"listing_type,omitempty"` // "sale" or "rent"
	MinPrice     *float64 `json:"min_price,omitempty"`
	MaxPrice     *float64 `json:"max_price,omitempty"`
	Location     string   `json:"location,omitempty"`
	Condition    string   `json:"condition,omitempty"`
}

// ConversationPlan is the decomposed, multi-dimensional reasoning plan.
type ConversationPlan struct {
	DetectedLanguage       string                   `json:"detected_language"` // en, am, om, ti
	Intents                []IntentType             `json:"intents"`
	SocialGreeting         string                   `json:"social_greeting,omitempty"`
	SocialFarewell         string                   `json:"social_farewell,omitempty"`
	IdentityResponse       string                   `json:"identity_response,omitempty"`
	PlatformFAQResponse    string                   `json:"platform_faq_response,omitempty"`
	GeneralResponse        string                   `json:"general_response,omitempty"`
	RequiresClarification  bool                     `json:"requires_clarification"`
	ClarificationQuestion  string                   `json:"clarification_question,omitempty"`
	Searches               []ExtractedSearchRequest `json:"searches,omitempty"`
	RequiresDatabaseSearch bool                     `json:"requires_database_search"`
}

// Regex patterns for price and budget extraction
var (
	rxPriceUnder = regexp.MustCompile(`(?i)\b(?:under|below|less than|max|budget of|up to|maximum of)\s*([\d,]+)\s*(?:etb|birr|ብር|qarshii)?\b`)
	rxPriceAbove = regexp.MustCompile(`(?i)\b(?:above|over|more than|min|minimum of|starting from|at least)\s*([\d,]+)\s*(?:etb|birr|ብር|qarshii)?\b`)
	rxPriceBirr  = regexp.MustCompile(`(?i)\b([\d,]+)\s*(?:etb|birr|ብር|qarshii)\b`)
	rxAmharicMax = regexp.MustCompile(`(?i)(?:ከ\s*)?([\d,]+)\s*(?:ብር|etb)?\s*(?:በታች|ያነሰ)`)
	rxOromoMax   = regexp.MustCompile(`(?i)(?:gatii\s*)?([\d,]+)\s*(?:qarshii|etb)?\s*(?:gadi)`)
	rxTigrinyaMax = regexp.MustCompile(`(?i)(?:ትሕቲ|ትሕቲ\s*)?([\d,]+)\s*(?:ቅርሺ|etb)?\s*(?:ትሕቲ)`)
)

// Known canonical category synonym registry for multi-category extraction.
type CategoryKeywordSpec struct {
	CanonicalName string
	CategorySlug  string
	Keywords      []string
	DefaultType   string // "sale" or "rent" (empty = any)
}

var categoryRegistry = []CategoryKeywordSpec{
	{
		CanonicalName: "Mobile Phones",
		CategorySlug:  "mobile-phones",
		Keywords: []string{
			"phone", "phones", "mobile", "mobiles", "smartphone", "smartphones", "iphone", "samsung", "redmi", "tecno", "infinix",
			"ስልክ", "ስልኮች", "ሞባይል", "ሞባይሎች", "አይፎን", "ሳምሰንግ",
			"bilbila", "bilbiloota", "moobaayilii", "ayifoonii",
			"ስልኪ", "ስልክታት", "ሞባይል", "ኣይፎን",
		},
	},
	{
		CanonicalName: "Computers & Laptops",
		CategorySlug:  "laptops-computers",
		Keywords: []string{
			"pc", "computer", "computers", "laptop", "laptops", "macbook", "desktop", "workstation", "dell", "hp", "lenovo",
			"ኮምፒውተር", "ኮምፒውተሮች", "ላፕቶፕ", "ላፕቶፖች", "ማክቡክ",
			"kompiitara", "kompiitaroota", "laaptooppii",
			"ኮምፒውተር", "ላፕቶፕ",
		},
	},
	{
		CanonicalName: "Residential Properties & Houses",
		CategorySlug:  "residential-properties-living-spaces",
		Keywords: []string{
			"house", "houses", "home", "homes", "villa", "villas", "apartment", "apartments", "flat", "flats", "condo", "condominium",
			"ቤት", "ቤቶች", "ቪላ", "ቪላዎች", "አፓርታማ", "አፓርታማዎች", "ኮንዶሚኒየም",
			"mana", "manneen", "viilaa", "appaartamaa", "kondominiyemii",
			"ገዛ", "ኣባይቲ", "ቪላ", "ኣፓርታማ", "ኮንዶሚኒየም",
		},
	},
	{
		CanonicalName: "Land & Plots",
		CategorySlug:  "land-plots-agricultural-development-acreage",
		Keywords: []string{
			"land", "plot", "plots", "acreage", "farmland",
			"መሬት", "ቦታ", "የእርሻ መሬት",
			"lafa", "lafoo",
			"መሬት", "ቦታ",
		},
	},
	{
		CanonicalName: "Vehicles & Automobiles",
		CategorySlug:  "passenger-vehicles-automobiles",
		Keywords: []string{
			"car", "cars", "automobile", "automobiles", "vehicle", "vehicles", "vitz", "corolla", "yaris", "pickup", "truck", "bajaj", "motorcycle",
			"መኪና", "መኪኖች", "ተሽከርካሪ", "ቪትዝ", "ኮሮላ", "ፒክአፕ", "ባጃጅ", "ሞተርሳይክል",
			"konkolaataa", "konkolaattota", "baajaajii",
			"መኪና", "መካይን", "ባጃጅ",
		},
	},
	{
		CanonicalName: "Agriculture & Tractors",
		CategorySlug:  "agriculture-crop-farming",
		Keywords: []string{
			"tractor", "tractors", "harvester", "plough", "farming equipment", "combine",
			"ትራክተር", "ትራክተሮች", "ማጨጃ", "የእርሻ መሳሪያ",
			"tiraaktara", "meeshaa qonnaa",
			"ትራክተር", "ናይ ሕርሻ መሳርሒ",
		},
	},
	{
		CanonicalName: "Jobs & Careers",
		CategorySlug:  "professional-occupations-careers-recruitment",
		Keywords: []string{
			"job", "jobs", "career", "employment", "hiring", "work", "vacancy", "vacancies",
			"ስራ", "ስራዎች", "ቅጥር", "ክፍት የስራ ቦታ",
			"hojii", "qacarraa",
			"ስራሕ", "ስራሕቲ", "ቅጥሪ",
		},
	},
	{
		CanonicalName: "Furniture & Home Appliances",
		CategorySlug:  "furniture-interior-furnishings",
		Keywords: []string{
			"furniture", "sofa", "bed", "chair", "table", "fridge", "refrigerator", "tv", "television", "washing machine",
			"ፈርኒቸር", "ሶፋ", "አልጋ", "ወንበር", "ጠረጴዛ", "ፍሪጅ", "ቲቪ", "ቴሌቪዥን", "ልብስ ማጠቢያ",
			"mi'a manaa", "firiijii", "televiziyoinii", "siree",
			"ፈርኒቸር", "ሶፋ", "ዓራት", "ፍሪጅ", "ተለቪዥን",
		},
	},
}

// Incomplete ambiguous word mappings for intelligent clarification
var clarificationMap = map[string]struct {
	ResolvedWord string
	Category     string
}{
	"pho":   {ResolvedWord: "phone", Category: "Mobile Phones"},
	"phon":  {ResolvedWord: "phone", Category: "Mobile Phones"},
	"lap":   {ResolvedWord: "laptop", Category: "Computers & Laptops"},
	"lapto": {ResolvedWord: "laptop", Category: "Computers & Laptops"},
	"comp":  {ResolvedWord: "computer", Category: "Computers & Laptops"},
	"hou":   {ResolvedWord: "house", Category: "Residential Properties & Houses"},
	"hous":  {ResolvedWord: "house", Category: "Residential Properties & Houses"},
	"apar":  {ResolvedWord: "apartment", Category: "Residential Properties & Houses"},
	"apart": {ResolvedWord: "apartment", Category: "Residential Properties & Houses"},
	"vil":   {ResolvedWord: "villa", Category: "Residential Properties & Houses"},
	"apa":   {ResolvedWord: "apartment", Category: "Residential Properties & Houses"},
	"sams":  {ResolvedWord: "Samsung", Category: "Mobile Phones"},
	"samsn": {ResolvedWord: "Samsung", Category: "Mobile Phones"},
	"iph":   {ResolvedWord: "iPhone", Category: "Mobile Phones"},
	"trac":  {ResolvedWord: "tractor", Category: "Agriculture & Tractors"},
	"mot":   {ResolvedWord: "motorcycle", Category: "Vehicles & Automobiles"},
}

// ClassifyMessage parses the complete utterance, resolves languages, multi-intent, entities, and context.
func ClassifyMessage(query string, locale string, context *ConversationContext) *ConversationPlan {
	cleanQuery := strings.TrimSpace(query)
	lower := strings.ToLower(cleanQuery)

	// 1. Language detection and normalization
	detectedLang := detectLanguage(cleanQuery, locale)

	plan := &ConversationPlan{
		DetectedLanguage: detectedLang,
		Intents:          make([]IntentType, 0),
		Searches:         make([]ExtractedSearchRequest, 0),
	}

	// 2. Check Out of Domain (Coding, science, physics, homework)
	if IsOutOfDomain(cleanQuery) {
		plan.Intents = append(plan.Intents, IntentOutOfDomain)
		plan.GeneralResponse = GetOutOfDomainRedirection(detectedLang)
		return plan
	}

	// 3. Social Conversational Dimension Check (Greeting, Farewell, Gratitude, Smalltalk)
	checkSocialIntents(lower, cleanQuery, detectedLang, plan)

	// 4. AI Identity & Capability Dimension Check
	checkIdentityAndCapability(lower, cleanQuery, detectedLang, plan)

	// 5. Platform FAQ Check
	if isFAQ, faqAnswer := CheckPlatformFAQ(cleanQuery, detectedLang); isFAQ {
		plan.Intents = append(plan.Intents, IntentPlatformFAQ)
		plan.PlatformFAQResponse = faqAnswer
	}

	// 6. Clarification for Incomplete / Ambiguous Short Tokens
	if checkClarification(cleanQuery, lower, detectedLang, plan) {
		return plan
	}

	// 7. Context Resolution for Follow-up Queries (e.g. "show me cheaper ones", "what about computers?")
	resolvedByContext := resolveContextFollowUp(cleanQuery, lower, detectedLang, context, plan)

	// 8. Marketplace Discovery & Multi-Category Extraction
	if !resolvedByContext {
		extractMarketplaceCategories(cleanQuery, lower, plan)
	}

	// 9. If searches were extracted, mark DB requirement
	if len(plan.Searches) > 0 {
		plan.RequiresDatabaseSearch = true
		plan.Intents = append(plan.Intents, IntentMarketplaceSearch)
	}

	// 10. Pure general conversation fallback (if no searches, no identity, no greeting, no FAQ)
	if len(plan.Intents) == 0 {
		handleGeneralConversation(lower, cleanQuery, detectedLang, plan)
	}

	return plan
}

func detectLanguage(text string, hintLocale string) string {
	hint := strings.ToLower(strings.TrimSpace(hintLocale))
	if hint == "or" {
		hint = "om"
	}
	if aiclient.ContainsGeez(text) {
		if hint == "ti" {
			return "ti"
		}
		// Check for specific Tigrinya Ge'ez characters or words
		if strings.Contains(text, "ገዛ") || strings.Contains(text, "ኣባይቲ") || strings.Contains(text, "ስራሕ") || strings.Contains(text, "ኣለኩም") || strings.Contains(text, "ክሕግዘኩም") {
			return "ti"
		}
		return "am"
	}

	lower := strings.ToLower(text)
	if strings.Contains(lower, "akkam") || strings.Contains(lower, "maal") || strings.Contains(lower, "mana") ||
		strings.Contains(lower, "bilbila") || strings.Contains(lower, "gatii") || strings.Contains(lower, "barbaada") {
		return "om"
	}

	if hint != "" && (hint == "am" || hint == "om" || hint == "ti" || hint == "en") {
		return hint
	}
	return "en"
}

func checkSocialIntents(lower string, raw string, lang string, plan *ConversationPlan) {
	// Greetings
	greetings := []string{
		"hello", "hi", "hey", "good morning", "good afternoon", "good evening", "howdy",
		"selam", "ሰላም", "እንደምን", "ጤና ይስጥልኝ", "እንደምን አደርክ", "እንደምን ዋልክ",
		"akkam", "akkam jirtu", "ashamaa", "nagayaa", "nagaa",
		"ሰላም", "ከመይ ኣለኻ", "ከመይ ቀኒኻ", "ደሓንዶ",
	}

	isGreeting := false
	for _, g := range greetings {
		if strings.HasPrefix(lower, g) || strings.Contains(lower, " "+g) || lower == g {
			isGreeting = true
			break
		}
	}

	if isGreeting {
		plan.Intents = append(plan.Intents, IntentSocialGreeting)
		switch lang {
		case "am":
			plan.SocialGreeting = "ሰላም! ወደ HUNTIQ እንኳን በደህና መጡ።"
		case "om":
			plan.SocialGreeting = "Akkam! Baga gara HUNTIQ nagaan dhuftan."
		case "ti":
			plan.SocialGreeting = "ሰላም! ናብ HUNTIQ ብደሓን መጻእኹም።"
		default:
			plan.SocialGreeting = "Hello! Welcome to HUNTIQ."
		}
	}

	// Gratitude
	if strings.Contains(lower, "thank") || strings.Contains(lower, "thanks") || strings.Contains(lower, "አመሰግናለሁ") ||
		strings.Contains(lower, "galatoomaa") || strings.Contains(lower, "የቐንየለይ") {
		plan.Intents = append(plan.Intents, IntentSocialGratitude)
	}

	// Farewell
	if strings.Contains(lower, "goodbye") || strings.Contains(lower, "bye") || strings.Contains(lower, "see you") ||
		strings.Contains(lower, "ደህና ሁን") || strings.Contains(lower, "ደህና ሁኑ") || strings.Contains(lower, "nagaatti") ||
		strings.Contains(lower, "ደሓን ኩን") {
		plan.Intents = append(plan.Intents, IntentSocialFarewell)
		switch lang {
		case "am":
			plan.SocialFarewell = "ደህና ይሁኑ! ሌላ ጊዜ በHUNTIQ ላይ እንደገና እንዲጎበኙን በደስታ እንጠብቃለን።"
		case "om":
			plan.SocialFarewell = "Nagaatti eegaa! Yeroo biraa deebitanii akka nu gaafattan abdanna."
		case "ti":
			plan.SocialFarewell = "ደሓን ኩኑ! ካልእ እዋን ኣብ HUNTIQ ተመሊስኩም ክትርእዩና ብሓጎስ ንጽበ።"
		default:
			plan.SocialFarewell = "Goodbye! Feel free to return anytime you need help discovering opportunities on HUNTIQ."
		}
	}

	// How are you / small talk
	if strings.Contains(lower, "how are you") || strings.Contains(lower, "how are you doing") ||
		strings.Contains(lower, "እንዴት ነህ") || strings.Contains(lower, "እንዴት ነሽ") || strings.Contains(lower, "እንደምን ነህ") ||
		strings.Contains(lower, "akkam jirta") || strings.Contains(lower, "ከመይ ኣለኻ") || strings.Contains(lower, "ከመይ ዲኻ") {
		plan.Intents = append(plan.Intents, IntentSocialSmalltalk)
	}
}

func checkIdentityAndCapability(lower string, raw string, lang string, plan *ConversationPlan) {
	// What is your name / Who are you / Creator
	isName := strings.Contains(lower, "what is your name") || strings.Contains(lower, "what's your name") ||
		strings.Contains(lower, "your name") || strings.Contains(lower, "ስምህ ማን ነው") || strings.Contains(lower, "ስምሽ ማን ነው") ||
		strings.Contains(lower, "maqaan kee eenyu") || strings.Contains(lower, "ሽምካ መን እዩ") || strings.Contains(lower, "መን እዩ ሽምካ")

	isWho := strings.Contains(lower, "who are you") || strings.Contains(lower, "who is this") ||
		strings.Contains(lower, "who made you") || strings.Contains(lower, "who created you") ||
		strings.Contains(lower, "who built you") || strings.Contains(lower, "who developed you") ||
		strings.Contains(lower, "who designed you") ||
		strings.Contains(lower, "አንተ ማን ነህ") || strings.Contains(lower, "አንቺ ማን ነሽ") ||
		strings.Contains(lower, "ማን ነህ") || strings.Contains(lower, "ማን ነሽ") ||
		strings.Contains(lower, "ማን ሰራህ") || strings.Contains(lower, "ማን ፈጠረህ") ||
		strings.Contains(lower, "ati eenyu") || strings.Contains(lower, "eenyutu si uume") ||
		strings.Contains(lower, "ንስኻ መን ኢኻ") || strings.Contains(lower, "መን ኢኻ")

	isWhat := strings.Contains(lower, "what is huntiq") || strings.Contains(lower, "what's huntiq") ||
		strings.Contains(lower, "tell me about huntiq") || strings.Contains(lower, "about huntiq") ||
		strings.Contains(lower, "huntiq ምንድን ነው") || strings.Contains(lower, "ስለ huntiq") ||
		strings.Contains(lower, "huntiq maali") || strings.Contains(lower, "huntiq እንታይ እዩ")

	isCapability := strings.Contains(lower, "what can you do") || strings.Contains(lower, "how can you help") ||
		strings.Contains(lower, "help me") || strings.Contains(lower, "ምን ማድረግ ትችላለህ") || strings.Contains(lower, "በምን ልትረዳኝ ትችላለህ") ||
		strings.Contains(lower, "maal gochuu dandeessa") || strings.Contains(lower, "እንታይ ክትገብር ትኽእል")

	if isName || isWho || isWhat {
		plan.Intents = append(plan.Intents, IntentIdentityInquiry)
		switch lang {
		case "am":
			plan.IdentityResponse = "እኔ HUNTIQ AI እባላለሁ — በኢትዮጵያ ውስጥ ያሉ እቃዎችን፣ መኪናዎችን፣ የቤት ኪራዮችን እና የተለያዩ አገልግሎቶችን በቀጥታ ከHUNTIQ ዳታቤዝ ፈልገው እንዲያገኙ የሚያግዝዎ ብልህ የገበያ ረዳት ነኝ።"
		case "om":
			plan.IdentityResponse = "Ani HUNTIQ AI jedhama — gargaaraa gabaa ammayyaa kan Itoophiyaa keessatti oomishaalee, konkolaattota, kireeffannaa manaa fi tajaajiloota kuusaa gabaa HUNTIQ irraa akka argattuuf si gargaaruudha."
		case "ti":
			plan.IdentityResponse = "ኣነ HUNTIQ AI ይበሃል — ኣብ ኢትዮጵያ ውሽጢ ዘለዉ ፍርያት፣ መካይን፣ ናይ ገዛ ክራይን ካልኦት ኣገልግሎታትን ካብ ዳታቤዝ HUNTIQ ብቐጥታ ክትረኽቡ ዝሕግዝ ናይ ዕዳጋ ሓጋዚኹም እየ።"
		default:
			plan.IdentityResponse = "I am HUNTIQ AI, your intelligent marketplace assistant. I help you discover verified products, vehicles, property rentals, and services across Ethiopia, compare market prices, and navigate the HUNTIQ platform."
		}
	} else if isCapability {
		plan.Intents = append(plan.Intents, IntentCapabilityInquiry)
		switch lang {
		case "am":
			plan.IdentityResponse = "በHUNTIQ ላይ የቀጥታ የገበያ ዝርዝሮችን መፈለግ፣ የዋጋና የቦታ ማጣሪያዎችን ማስተካከል፣ ስለ እቃዎች ዝርዝር መረጃዎችን ማወዳደር፣ እና በድረ-ገጹ ላይ እንዴት መለጠፍ እንደሚቻል ማብራሪያዎችን መስጠት እችላለሁ። ዛሬ ምን መፈለግ ይፈልጋሉ?"
		case "om":
			plan.IdentityResponse = "Kuusaa gabaa HUNTIQ irraa tarreeffamoota barbaaduu, gatii fi bakka filachuu, oomishaalee walbira qabuu, fi akkamitti akka maxxansan isiniif ibsuu nan danda'a. Har'a maal barbaaddu?"
		case "ti":
			plan.IdentityResponse = "ኣብ HUNTIQ ናይ ቀጥታ ዝርዝራት ክደልየልኩም፣ ብዋጋን ቦታን ከጻርየልኩም፣ ፍርያት ከነጻጽረልኩምን ከመይ ጌርኩም ከም እትልጥፉ ክሕብረኩምን እኽእል እየ። ሎሚ እንታይ ክትደልዩ ትደልዩ?"
		default:
			plan.IdentityResponse = "I can search live marketplace listings across Ethiopia, filter by price and location, compare different products, check seller trust tiers, and guide you through posting your own listings on HUNTIQ. What are you looking for today?"
		}
	}
}

func checkClarification(cleanQuery string, lower string, lang string, plan *ConversationPlan) bool {
	// If query is very short (< 4 chars in Latin) or 1-2 words that look like abbreviations
	trimmed := strings.TrimSpace(lower)
	if utf8.RuneCountInString(trimmed) < 2 {
		return false
	}

	// Check if matches our incomplete token ambiguity dictionary
	if entry, found := clarificationMap[trimmed]; found {
		plan.RequiresClarification = true
		plan.Intents = append(plan.Intents, IntentClarificationNeeded)
		switch lang {
		case "am":
			plan.ClarificationQuestion = "ለ'" + cleanQuery + "' ምርመራ ሳደርግ፤ '" + entry.ResolvedWord + "' (" + entry.Category + ") ማለትዎ ፈልገው ነው? አዎ ከሆነ የቀጥታ ዝርዝሮችን ወዲያውኑ ላሳይዎት እችላለሁ።"
		case "om":
			plan.ClarificationQuestion = "' " + cleanQuery + "' jechuun keessan '" + entry.ResolvedWord + "' (" + entry.Category + ") jechuudhaa? Eeyyee yoo ta'e tarreeffama jiru isiniif nan barbaada."
		case "ti":
			plan.ClarificationQuestion = "ብ'" + cleanQuery + "' ዝበልኩምዎ '" + entry.ResolvedWord + "' (" + entry.Category + ") ማለትኩም ድዩ? እወ እንተኾይኑ ዘለዉ ዝርዝራት ከርእየኩም እኽእል እየ።"
		default:
			plan.ClarificationQuestion = "Did you mean **" + entry.ResolvedWord + "** (" + entry.Category + ")? If so, I would be happy to search matching live listings for you."
		}
		return true
	}

	return false
}

func resolveContextFollowUp(cleanQuery string, lower string, lang string, ctx *ConversationContext, plan *ConversationPlan) bool {
	if ctx == nil || len(ctx.RecentSearches) == 0 {
		return false
	}

	lastSearch := ctx.RecentSearches[len(ctx.RecentSearches)-1]

	// Check for price follow-ups: "cheaper ones", "something cheaper", "lower price"
	isCheaper := strings.Contains(lower, "cheaper") || strings.Contains(lower, "cheapest") ||
		strings.Contains(lower, "lower price") || strings.Contains(lower, "less expensive") ||
		strings.Contains(lower, "ቅናሽ") || strings.Contains(lower, "ረከስ ያለ") || strings.Contains(lower, "ተመጣጣኝ") ||
		strings.Contains(lower, "rakasa") || strings.Contains(lower, "gatii gadi aanaa") ||
		strings.Contains(lower, "ሕሳር") || strings.Contains(lower, "ዝሓሰረ")

	// Check for brand follow-ups: "only samsung", "show me dell"
	var brandFound string
	brands := []string{"samsung", "apple", "iphone", "dell", "hp", "lenovo", "toyota", "hyundai", "nissan"}
	for _, b := range brands {
		if strings.Contains(lower, b) {
			brandFound = b
			break
		}
	}

	// Check for explicit category switch: "what about computers?", "find a house too"
	if strings.HasPrefix(lower, "what about ") || strings.HasPrefix(lower, "how about ") || strings.HasPrefix(lower, "and ") {
		// This is a switch, let normal extraction handle it
		return false
	}

	if isCheaper || brandFound != "" {
		newSearch := lastSearch // Inherit category and context!

		if brandFound != "" {
			newSearch.SearchTerm = brandFound + " " + lastSearch.SearchTerm
		}

		if isCheaper {
			// If last search had a max price, lower it or sort by price asc
			if lastSearch.MaxPrice != nil && *lastSearch.MaxPrice > 0 {
				newMax := *lastSearch.MaxPrice * 0.75
				newSearch.MaxPrice = &newMax
			}
		}

		plan.Searches = append(plan.Searches, newSearch)
		plan.Intents = append(plan.Intents, IntentMarketplaceSearch)
		plan.RequiresDatabaseSearch = true
		return true
	}

	return false
}

func extractMarketplaceCategories(cleanQuery string, lower string, plan *ConversationPlan) {
	// Extract potential price bounds across the sentence
	maxPrice, minPrice := extractPriceConstraints(lower)
	location := extractLocation(lower)
	listingType := extractListingType(lower)

	type catMatch struct {
		spec CategoryKeywordSpec
		pos  int
		kw   string
	}

	var matches []catMatch
	seenCat := make(map[string]bool)

	for _, spec := range categoryRegistry {
		earliestPos := -1
		matchedKw := ""
		for _, kw := range spec.Keywords {
			if containsWord(lower, kw) {
				idx := strings.Index(lower, kw)
				if earliestPos == -1 || (idx >= 0 && idx < earliestPos) {
					earliestPos = idx
					matchedKw = kw
				}
			}
		}
		if earliestPos >= 0 && !seenCat[spec.CanonicalName] {
			seenCat[spec.CanonicalName] = true
			matches = append(matches, catMatch{
				spec: spec,
				pos:  earliestPos,
				kw:   matchedKw,
			})
		}
	}

	// Sort matches strictly by their appearance order in the utterance
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].pos < matches[j].pos
	})

	// If categories matched (e.g. phone AND pc AND house):
	if len(matches) > 0 {
		for _, m := range matches {
			req := ExtractedSearchRequest{
				CategoryName: m.spec.CanonicalName,
				CategorySlug: m.spec.CategorySlug,
				SearchTerm:   m.spec.Keywords[0], // Clean canonical keyword (e.g. "phone", "pc", "house")
				MinPrice:     minPrice,
				MaxPrice:     maxPrice,
				Location:     location,
			}
			if listingType != nil {
				req.ListingType = listingType
			} else if m.spec.DefaultType != "" {
				req.ListingType = &m.spec.DefaultType
			}
			plan.Searches = append(plan.Searches, req)
		}
		return
	}

	// Check if the query has commercial buying/selling intent even without an exact category match
	commercialKeywords := []string{
		"buy", "sell", "rent", "looking for", "want", "need", "cost of", "price of",
		"መግዛት", "መሸጥ", "መከራየት", "እፈልጋለሁ", "ዋጋው ስንት ነው",
		"bituu", "gurguruu", "kireeffachuu", "barbaada", "gatiin meeqa",
		"ምግዛእ", "ምሻጥ", "ምክራይ", "እደሊ", "ዋጋ ክንደይ እዩ",
	}

	hasCommercial := false
	for _, cw := range commercialKeywords {
		if strings.Contains(lower, cw) {
			hasCommercial = true
			break
		}
	}

	if hasCommercial {
		cleaned := aiclient.CleanConversationalSearchTerm(cleanQuery)
		if cleaned != "" && len(cleaned) >= 2 {
			plan.Searches = append(plan.Searches, ExtractedSearchRequest{
				CategoryName: "Marketplace Search",
				SearchTerm:   cleaned,
				MinPrice:     minPrice,
				MaxPrice:     maxPrice,
				Location:     location,
				ListingType:  listingType,
			})
			return
		}
	}

	// If the user's intent is already identified as a pure social greeting or identity inquiry,
	// do NOT fall through into broad product search unless explicit commercial constraints exist.
	isPureSocialOrIdentity := containsIntent(plan.Intents, IntentSocialGreeting) ||
		containsIntent(plan.Intents, IntentIdentityInquiry) ||
		containsIntent(plan.Intents, IntentSocialFarewell) ||
		containsIntent(plan.Intents, IntentSocialSmalltalk)

	if isPureSocialOrIdentity {
		return
	}

	// If no predefined registry category matched, check if it's a broad product query
	if IsBroadProductSearch(cleanQuery) {
		plan.Searches = append(plan.Searches, ExtractedSearchRequest{
			CategoryName: "All Products",
			SearchTerm:   "",
			MinPrice:     minPrice,
			MaxPrice:     maxPrice,
			Location:     location,
		})
		return
	}
}

func extractPriceConstraints(lower string) (*float64, *float64) {
	var maxP *float64
	var minP *float64

	// Under / Max price
	if match := rxPriceUnder.FindStringSubmatch(lower); len(match) > 1 {
		val := parseCleanNumber(match[1])
		if val > 0 {
			maxP = &val
		}
	} else if match := rxAmharicMax.FindStringSubmatch(lower); len(match) > 1 {
		val := parseCleanNumber(match[1])
		if val > 0 {
			maxP = &val
		}
	} else if match := rxOromoMax.FindStringSubmatch(lower); len(match) > 1 {
		val := parseCleanNumber(match[1])
		if val > 0 {
			maxP = &val
		}
	} else if match := rxTigrinyaMax.FindStringSubmatch(lower); len(match) > 1 {
		val := parseCleanNumber(match[1])
		if val > 0 {
			maxP = &val
		}
	}

	// Above / Min price
	if match := rxPriceAbove.FindStringSubmatch(lower); len(match) > 1 {
		val := parseCleanNumber(match[1])
		if val > 0 {
			minP = &val
		}
	}

	return maxP, minP
}

func extractLocation(lower string) string {
	knownLocations := []string{
		"addis ababa", "bole", "cmc", "sarbet", "kazanchis", "piassa", "mexico", "gerji", "ayat", "lebu", "megenagna",
		"hawassa", "adama", "bahir dar", "mekelle", "dire dawa", "bishoftu",
		"አዲስ አበባ", "ቦሌ", "ሲኤምሲ", "ሳርቤት", "ካዛንቺስ", "ፒያሳ", "ገርጂ", "አያት", "መገናኛ",
		"finfinnee", "boolee",
		"ኣዲስ ኣበባ", "መቐለ",
	}

	for _, loc := range knownLocations {
		if strings.Contains(lower, loc) {
			// Normalize to canonical city / subcity string
			switch loc {
			case "addis ababa", "አዲስ አበባ", "finfinnee", "ኣዲስ ኣበባ":
				return "Addis Ababa"
			case "bole", "ቦሌ", "boolee":
				return "Bole"
			case "cmc", "ሲኤምሲ":
				return "CMC"
			case "hawassa":
				return "Hawassa"
			case "adama":
				return "Adama"
			case "bahir dar":
				return "Bahir Dar"
			case "mekelle", "መቐለ":
				return "Mekelle"
			case "dire dawa":
				return "Dire Dawa"
			default:
				return strings.Title(loc)
			}
		}
	}
	return ""
}

func extractListingType(lower string) *string {
	if strings.Contains(lower, "for rent") || strings.Contains(lower, "rent") || strings.Contains(lower, "rental") ||
		strings.Contains(lower, "ኪራይ") || strings.Contains(lower, "የሚከራይ") || strings.Contains(lower, "ክራይ") ||
		strings.Contains(lower, "kiraa") || strings.Contains(lower, "kireeffamu") {
		t := "rent"
		return &t
	}
	if strings.Contains(lower, "for sale") || strings.Contains(lower, "to buy") || strings.Contains(lower, "ሽያጭ") ||
		strings.Contains(lower, "የሚሸጥ") || strings.Contains(lower, "gurgurtaa") || strings.Contains(lower, "ሽያጥ") {
		t := "sale"
		return &t
	}
	return nil
}

func handleGeneralConversation(lower string, raw string, lang string, plan *ConversationPlan) {
	plan.Intents = append(plan.Intents, IntentGeneralKnowledge)
	switch lang {
	case "am":
		plan.GeneralResponse = "በHUNTIQ ላይ የቀጥታ የገበያ መረጃዎችን እና እቃዎችን መፈለግ ይችላሉ። ምን መፈለግ ይፈልጋሉ? ለምሳሌ፡ ስልኮች፣ መኪናዎች፣ ወይም የሚከራዩ ቤቶች ብለው ሊጠይቁኝ ይችላሉ።"
	case "om":
		plan.GeneralResponse = "HUNTIQ irratti oomishaalee fi tajaajiloota gabaa barbaaduu dandeessu. Maal barbaaddu? Fakkeenyaaf: bilbila, konkolaataa ykn mana kireeffamu jechuun na gaafachuu dandeessu."
	case "ti":
		plan.GeneralResponse = "ኣብ HUNTIQ ናይ ዕዳጋ ሓበሬታታትን ፍርያትን ክትረኽቡ ትኽእሉ ኢኹም። እንታይ ክትደልዩ ትደልዩ? ንኣብነት፡ ስልኪ፣ መኪና ወይ ናይ ገዛ ክራይ ብምባል ክትሓቱኒ ትኽእሉ ኢኹም።"
	default:
		plan.GeneralResponse = "I'm here to help you explore verified listings, check prices, or find items across Ethiopia on HUNTIQ. You can ask me to search for phones, laptops, cars, rentals, or explain how to post an ad."
	}
}

func parseCleanNumber(s string) float64 {
	clean := strings.ReplaceAll(s, ",", "")
	clean = strings.TrimSpace(clean)
	val, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0
	}
	return val
}

func containsWord(s, word string) bool {
	// For Ge'ez scripts, direct substring check with word boundary logic
	if aiclient.ContainsGeez(word) {
		return strings.Contains(s, word)
	}

	// For Latin tokens, enforce word boundary to avoid false positives (e.g. "car" in "carecare")
	pattern := `(?i)(?:^|[\s,.\-!?;/])` + regexp.QuoteMeta(word) + `(?:$|[\s,.\-!?;/])`
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return strings.Contains(s, word)
	}
	return matched
}
