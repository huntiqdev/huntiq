package ai

import (
	"regexp"
	"strings"
	"unicode"
)

// ClassificationResult holds the evaluation outcome of user prompt inspection.
type ClassificationResult struct {
	Intent          DomainIntent
	Language        string
	IsSupportedLang bool
	IsAllowed       bool
	Confidence      float64
	Reason          string
}

// Pre-compiled regex patterns for deterministic security & domain classification
var (
	// Adversarial bypass and jailbreak patterns
	reJailbreak = regexp.MustCompile(`(?i)(ignore\s+(all\s+)?previous\s+instructions|disregard\s+(all\s+)?prior|you\s+are\s+now\s+(an\s+)?unrestricted|pretend\s+to\s+be|dan\s+mode|jailbreak|reveal\s+(the\s+)?system\s+prompt|reveal\s+.*secret|developer\s+mode|print\s+(the\s+)?hidden|show\s+me\s+your\s+instructions|base64|rot13)`)

	// Math and equations (e.g., 2+2, 100/5, e=mc^2, derivative of x^2)
	reMathEquations = regexp.MustCompile(`(?i)(e\s*=\s*mc\s*[\^²2]|derivative\s+of|integral\s+of|solve\s+for\s+[a-z]|\b\d+\s*[\+\-\*\/]\s*\d+\b|log\s*\(\s*\d+\s*\)|sin\s*\(\s*\d+\s*\)|cos\s*\(\s*\d+\s*\)|quadratic\s+equation)`)

	// Academic physics, chemistry, biology
	reAcademicSciences = regexp.MustCompile(`(?i)(quantum\s+mechanics|quantum\s+physics|thermodynamics|general\s+relativity|string\s+theory|photosynthesis|mitochondria|chemical\s+reaction|avogadro|periodic\s+table)`)

	// General programming tutorials / code generation unrelated to marketplace
	reGenericCoding = regexp.MustCompile(`(?i)(write\s+a\s+(python|javascript|golang|c\+\+|java)\s+script|create\s+a\s+binary\s+tree|implement\s+(quick\s*sort|merge\s*sort|bubble\s*sort)|write\s+code\s+to|how\s+to\s+install\s+linux|leetcode|hackerrank)`)

	// General trivia, world politics, entertainment
	reGeneralTrivia = regexp.MustCompile(`(?i)(who\s+won\s+.*world\s+cup|capital\s+of\s+[a-z]+|who\s+is\s+the\s+president\s+of\s+[a-z]+|tell\s+me\s+a\s+joke|write\s+an\s+essay\s+about|horoscope|celebrity\s+gossip)`)

	// Unsupported foreign scripts: Arabic, Cyrillic, CJK, Devanagari, Hebrew
	reUnsupportedScripts = regexp.MustCompile(`[\p{Arabic}\p{Cyrillic}\p{Han}\p{Hiragana}\p{Katakana}\p{Devanagari}\p{Hebrew}]`)
)

// LanguageClassifier determines whether text is in one of the 4 supported languages.
type LanguageClassifier struct{}

func NewLanguageClassifier() *LanguageClassifier {
	return &LanguageClassifier{}
}

// ClassifyLanguage inspects the text and returns language code and whether it is supported.
func (c *LanguageClassifier) ClassifyLanguage(text, requestedLang string) (string, bool) {
	req := strings.ToLower(strings.TrimSpace(requestedLang))
	if req == "om" {
		req = LangOR
	}
	if IsSupportedLanguage(req) {
		return req, true
	}

	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return LangEN, true
	}

	// Check for unsupported foreign scripts
	if reUnsupportedScripts.MatchString(trimmed) {
		return "unsupported", false
	}

	// Count Ge'ez / Ethiopic runes
	ethiopicCount := 0
	latinCount := 0
	totalRunes := 0

	for _, r := range trimmed {
		if unicode.IsLetter(r) {
			totalRunes++
			if r >= 0x1200 && r <= 0x137F { // Ethiopic Unicode block
				ethiopicCount++
			} else if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				latinCount++
			}
		}
	}

	// Ge'ez script analysis (Amharic vs Tigrinya)
	if ethiopicCount > 0 && float64(ethiopicCount)/float64(totalRunes+1) > 0.3 {
		lower := strings.ToLower(trimmed)
		// Specific Tigrinya markers
		tigrinyaKeywords := []string{
			"ከመይ", "ኣለኹም", "ደሊየ", "ኣሎ", "እንታይ", "ምግዛእ", "ምሻጥ", "ክራይ", "ገዛ", "ዋጋ",
			"ሓበሬታ", "ብኸመይ", "ጽቡቕ", "ስራሕ", "መኪና", "ስልኪ", "ኣብ", "ናይ",
		}
		for _, kw := range tigrinyaKeywords {
			if strings.Contains(lower, kw) {
				return LangTI, true
			}
		}
		return LangAM, true
	}

	// Latin script analysis (Afaan Oromo vs English vs other Latin)
	lower := strings.ToLower(trimmed)
	afaanOromoKeywords := []string{
		"akkam", "jirtu", "barbaada", "biyya", "maaloo", "bituu", "gurguruu", "kireeffachuu",
		"mana", "konkolaataa", "bilbila", "hojii", "caalbaasii", "gabaa", "meeshaa", "itti",
		"keessan", "argachuuf", "akkamitti", "waan", "gaarii", "nagaa",
	}
	oromoHits := 0
	for _, kw := range afaanOromoKeywords {
		if strings.Contains(lower, kw) {
			oromoHits++
		}
	}
	if oromoHits >= 1 {
		return LangOR, true
	}

	return LangEN, true
}

// DomainClassifier strictly controls what requests HUNTIQ AI is authorized to process.
type DomainClassifier struct {
	langClassifier *LanguageClassifier
}

func NewDomainClassifier() *DomainClassifier {
	return &DomainClassifier{
		langClassifier: NewLanguageClassifier(),
	}
}

// Classify performs comprehensive pre-screening on incoming user messages.
func (d *DomainClassifier) Classify(prompt, requestedLang string) ClassificationResult {
	res := ClassificationResult{
		Intent:          IntentMarketplaceDiscovery,
		IsAllowed:       true,
		IsSupportedLang: true,
		Confidence:      1.0,
	}

	// 1. Language validation
	lang, isSupported := d.langClassifier.ClassifyLanguage(prompt, requestedLang)
	res.Language = lang
	res.IsSupportedLang = isSupported

	if !isSupported {
		res.IsAllowed = false
		res.Intent = IntentOutOfDomain
		res.Reason = "Unsupported language detected. HUNTIQ AI strictly operates in English, Amharic, Tigrinya, and Afaan Oromo."
		return res
	}

	clean := strings.TrimSpace(prompt)
	if clean == "" {
		res.Intent = IntentAmbiguous
		res.Reason = "Empty prompt"
		return res
	}

	// 2. Adversarial Injection Detection
	if reJailbreak.MatchString(clean) {
		res.IsAllowed = false
		res.Intent = IntentAdversarialBypass
		res.Confidence = 0.99
		res.Reason = "Adversarial prompt injection attempt detected and blocked."
		return res
	}

	// 3. Out-of-Domain Filter: Mathematics & Equations
	if reMathEquations.MatchString(clean) {
		res.IsAllowed = false
		res.Intent = IntentOutOfDomain
		res.Confidence = 0.95
		res.Reason = "Academic mathematics or equation question out of marketplace domain."
		return res
	}

	// 4. Out-of-Domain Filter: Academic Sciences (Physics, Chemistry, Biology)
	if reAcademicSciences.MatchString(clean) {
		res.IsAllowed = false
		res.Intent = IntentOutOfDomain
		res.Confidence = 0.95
		res.Reason = "Academic physical science question out of marketplace domain."
		return res
	}

	// 5. Out-of-Domain Filter: Generic Programming Tutorials
	if reGenericCoding.MatchString(clean) {
		res.IsAllowed = false
		res.Intent = IntentOutOfDomain
		res.Confidence = 0.95
		res.Reason = "General software development tutorial question out of marketplace domain."
		return res
	}

	// 6. Out-of-Domain Filter: General Trivia / Entertainment
	if reGeneralTrivia.MatchString(clean) {
		res.IsAllowed = false
		res.Intent = IntentOutOfDomain
		res.Confidence = 0.90
		res.Reason = "General non-marketplace trivia question out of domain."
		return res
	}

	// 7. Policy / Account Assistance In-Domain Classification
	lower := strings.ToLower(clean)
	if strings.Contains(lower, "how to register") || strings.Contains(lower, "how to post") ||
		strings.Contains(lower, "how to sell") || strings.Contains(lower, "listing fee") ||
		strings.Contains(lower, "payment policy") || strings.Contains(lower, "refund policy") ||
		strings.Contains(lower, "እንዴት መመዝገብ") || strings.Contains(lower, "የክፍያ ፖሊሲ") ||
		strings.Contains(lower, "ብኸመይ ክምዝገብ") || strings.Contains(lower, "akkamitti galmaa'uu") {
		res.Intent = IntentPolicyHelp
		res.IsAllowed = true
		return res
	}

	// 8. In-Domain Marketplace Intent (Products, Properties, Vehicles, Services, Gigs, Tenders)
	res.Intent = IntentMarketplaceDiscovery
	res.IsAllowed = true
	return res
}

// GenerateUnsupportedLanguageReply crafts a message directing the user to the 4 supported languages.
func GenerateUnsupportedLanguageReply() string {
	return "HUNTIQ AI supports English, Amharic (አማርኛ), Tigrinya (ትግርኛ), and Afaan Oromo. Please select or type in one of these four languages to discover products, properties, vehicles, and services on HUNTIQ.\n\n" +
		"HUNTIQ AI በእንግሊዝኛ፣ በአማርኛ፣ በትግርኛ እና በአፋን ኦሮሞ ቋንቋዎች ብቻ ያገለግላል። እባክዎ ከእነዚህ ቋንቋዎች አንዱን ይምረጡ።\n\n" +
		"HUNTIQ AI ብእንግሊዝኛ፣ ብትግርኛ፣ ብኣምሓርኛን ብኦሮምኛን ጥራይ የገልግል። በጃኹም ካብዞም ቋንቋታት ሓደ ምረጹ።\n\n" +
		"HUNTIQ AI afaan Ingiliffaa, Oromoo, Amaaraa fi Tigriffaa qofaan tajaajila kenna. Maaloo afaanota kanneen keessaa tokko filadhaa."
}
