// Package search provides language-agnostic text utilities for Ethiopian marketplace
// search query processing. It is intentionally AI-free and has no dependency on the
// aiclient package so that repositories, services, and handlers can import it without
// pulling in any AI-specific types or infrastructure.
package search

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ContainsGeez returns true if s contains any Ethiopic / Ge'ez Unicode code points (U+1200–U+137F).
func ContainsGeez(s string) bool {
	for _, r := range s {
		if r >= 0x1200 && r <= 0x137F {
			return true
		}
	}
	return false
}

// knownEthiopianLocations is used to strip location tokens from search queries.
var knownEthiopianLocations = []string{
	"bole", "piassa", "merkato", "kazanchis", "megenagna", "cmc", "sarbet", "gerji",
	"jomo", "summit", "yeka", "kirkos", "arada", "gullele", "kolfe", "nifas silk",
	"addis ababa", "hawassa", "adama", "bahir dar", "gondar", "mekelle", "dire dawa",
	"jimma", "dessie", "assosa", "shashamane", "arba minch", "harar",
}

var (
	rxPriceUnder = regexp.MustCompile(`(?i)(?:under|below|less than|max(?:imum)?|up to|ከ|በታች|hanga)\s*(\d+(?:,\d+)*(?:\.\d+)?)\s*(k|thousand|ሺህ|etb|birr)?`)
	rxPriceAbove = regexp.MustCompile(`(?i)(?:above|over|more than|min(?:imum)?|starting from|ከ|በላይ)\s*(\d+(?:,\d+)*(?:\.\d+)?)\s*(k|thousand|ሺህ|etb|birr)?`)
	rxPriceK     = regexp.MustCompile(`(?i)\b(\d+(?:,\d+)*(?:\.\d+)?)\s*k\b`)
	rxPriceETB   = regexp.MustCompile(`(?i)\b(\d+(?:,\d+)*(?:\.\d+)?)\s*(?:etb|birr)\b`)
)

// ethiopianKeywordMap maps Ethiopian-language product terms to canonical English equivalents.
var ethiopianKeywordMap = map[string]string{
	"ስልክ": "phone", "ሞባይል": "phone", "ስማርትፎን": "phone", "አይፎን": "phone", "ሳምሰንግ": "phone",
	"bilbila": "phone", "moobaayilii": "phone", "ስልኪ": "phone",
	"መኪና": "car", "ተሽከርካሪ": "car", "ቶዮታ": "car", "ቪትዝ": "car", "ኮሮላ": "car",
	"konkolaataa": "car", "otoomobil": "car",
	"ቤት": "house", "ቪላ": "villa", "ሪልስቴት": "house", "mana": "house", "viillaa": "villa",
	"ገዛ": "house", "አፓርታማ": "apartment", "ኮንዶሚኒየም": "apartment", "ኮንዶ": "apartment",
	"ስቱዲዮ": "apartment", "apaartaamaa": "apartment", "kondoominiyeemii": "apartment",
	"ኣፓርታማ": "apartment", "condo": "apartment",
	"ኮምፒውተር": "laptop", "ላፕቶፕ": "laptop", "ማክቡክ": "laptop",
	"kompiitara": "laptop", "laaptooppii": "laptop",
	"ጫማ": "shoes", "ስኒከር": "shoes", "kophee": "shoes",
	"ልብስ": "clothes", "ቀሚስ": "clothes", "ሱሪ": "clothes", "wayyaa": "clothes",
	"uffata": "clothes", "ክዳን": "clothes",
	"ቲቪ": "tv", "ቴሌቪዥን": "tv", "televizhinii": "tv", "tivii": "tv",
	"ሰዓት": "watch", "sa'aatii": "watch",
	"ፈርኒቸር": "furniture", "ሶፋ": "furniture", "አልጋ": "furniture",
	"ጠረጴዛ": "furniture", "ወንበር": "furniture", "farnichara": "furniture",
	"siree": "furniture", "ዓራት": "furniture",
	"ጀነሬተር": "generator", "መሰርሰሪያ": "drill",
}

// marketplaceSynonymGroups defines bidirectional cross-lingual synonyms for EN, AM, OM, TI.
var marketplaceSynonymGroups = [][]string{
	{"phone", "phones", "mobile", "smartphone", "smartphones", "cellphone", "iphone", "samsung", "tecno", "infinix", "redmi", "xiaomi", "ስልክ", "ሞባይል", "ስማርትፎን", "አይፎን", "ሳምሰንግ", "bilbila", "moobaayilii", "ስልኪ"},
	{"car", "cars", "vehicle", "vehicles", "automobile", "toyota", "vitz", "corolla", "yaris", "suzuki", "hyundai", "መኪና", "ተሽከርካሪ", "ቶዮታ", "ቪትዝ", "ኮሮላ", "konkolaataa", "otoomobil"},
	{"villa", "villas", "mansion", "ቪላ", "ቪላዎች", "viillaa"},
	{"house", "houses", "home", "real estate", "property", "residence", "ቤት", "ሪልስቴት", "mana", "ገዛ"},
	{"apartment", "apartments", "condo", "condos", "condominium", "studio", "flat", "አፓርታማ", "ኮንዶሚኒየም", "ኮንዶ", "ስቱዲዮ", "apaartaamaa", "kondoominiyeemii", "ኣፓርታማ"},
	{"laptop", "laptops", "computer", "computers", "pc", "macbook", "dell", "hp", "lenovo", "ኮምፒውተር", "ላፕቶፕ", "ማክቡክ", "kompiitara", "laaptooppii"},
	{"shoe", "shoes", "sneakers", "boots", "loafers", "ጫማ", "ስኒከር", "ቡትስ", "kophee"},
	{"clothes", "clothing", "dress", "shirt", "pants", "suit", "jacket", "t-shirt", "ልብስ", "ቀሚስ", "ሱሪ", "ጃኬት", "ቲሸርት", "wayyaa", "uffata", "ክዳን"},
	{"tv", "television", "smart tv", "screen", "ቲቪ", "ቴሌቪዥን", "televizhinii", "tivii"},
	{"watch", "watches", "smartwatch", "ሰዓት", "smart watch", "sa'aatii"},
	{"furniture", "sofa", "bed", "table", "chair", "desk", "ፈርኒቸር", "ሶፋ", "አልጋ", "ጠረጴዛ", "ወንበር", "farnichara", "siree", "ዓራት"},
	{"tool", "tools", "drill", "generator", "hammer", "saw", "መሣሪያ", "እቃ", "ጀነሬተር", "መሰርሰሪያ", "meeshaa"},
}

// GetSearchSynonyms returns cross-lingual synonyms for the given term.
func GetSearchSynonyms(term string) []string {
	clean := strings.ToLower(strings.TrimSpace(term))
	if clean == "" {
		return nil
	}
	if mapped, ok := ethiopianKeywordMap[clean]; ok {
		clean = mapped
	}
	words := strings.Fields(clean)
	for _, group := range marketplaceSynonymGroups {
		matched := false
		cleanRunes := len([]rune(clean))
		for _, member := range group {
			memberRunes := len([]rune(member))
			if clean == member || (memberRunes >= 2 && strings.Contains(clean, member)) || (cleanRunes >= 3 && strings.Contains(member, clean)) {
				matched = true
				break
			}
			for _, w := range words {
				wRunes := len([]rune(w))
				if w == member || (memberRunes >= 2 && strings.Contains(w, member)) || (wRunes >= 3 && strings.Contains(member, w)) {
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}
		if matched {
			return group
		}
	}
	return []string{clean}
}

// BuildSynonymTsQuery constructs a sanitized PostgreSQL tsquery disjunction string
// (e.g. "'phone' | 'ስልክ' | 'bilbila'") for fast GIN index search across languages.
func BuildSynonymTsQuery(term string) string {
	syns := GetSearchSynonyms(term)
	if len(syns) <= 1 {
		return ""
	}
	seen := make(map[string]bool)
	var tokens []string
	for _, s := range syns {
		s = strings.TrimSpace(s)
		if s == "" || strings.Contains(s, " ") {
			continue
		}
		sanitized := strings.Map(func(r rune) rune {
			if r == '\'' || r == '"' || r == '\\' || r == '|' || r == '&' || r == '!' || r == '(' || r == ')' {
				return -1
			}
			return r
		}, s)
		sanitized = strings.TrimSpace(sanitized)
		if sanitized == "" || seen[sanitized] {
			continue
		}
		seen[sanitized] = true
		tokens = append(tokens, fmt.Sprintf("'%s'", sanitized))
	}
	if len(tokens) <= 1 {
		return ""
	}
	return strings.Join(tokens, " | ")
}

// GetCanonicalEnglishTerm resolves any localized product term to its primary English name.
func GetCanonicalEnglishTerm(term string) string {
	clean := strings.ToLower(strings.TrimSpace(term))
	if clean == "" {
		return ""
	}
	if mapped, ok := ethiopianKeywordMap[clean]; ok {
		return mapped
	}
	cleanRunes := len([]rune(clean))
	for _, group := range marketplaceSynonymGroups {
		for _, m := range group {
			mRunes := len([]rune(m))
			if clean == m || (mRunes >= 2 && strings.Contains(clean, m)) || (cleanRunes >= 3 && strings.Contains(m, clean)) {
				return group[0]
			}
		}
	}
	return clean
}

// CleanConversationalSearchTerm strips conversational filler, pronouns, prepositions,
// and price/location tokens from a user query, isolating the item name for PostgreSQL full-text search.
func CleanConversationalSearchTerm(q string) string {
	cleaned := strings.TrimSpace(q)
	if cleaned == "" {
		return ""
	}

	for _, ch := range []string{
		"?", "!", ".", ",", "\"", "'", ":", ";", "(", ")", "[", "]", "{", "}",
		"-", "_", "/", "\\", "@", "#", "$", "%", "^", "&", "*", "+", "=", "~", "`",
		"፧", "፡", "።", "፣", "፤", "፦", "፠", "፨",
	} {
		cleaned = strings.ReplaceAll(cleaned, ch, " ")
	}

	lower := strings.ToLower(cleaned)

	lower = rxPriceUnder.ReplaceAllString(lower, " ")
	lower = rxPriceAbove.ReplaceAllString(lower, " ")
	lower = rxPriceK.ReplaceAllString(lower, " ")
	lower = rxPriceETB.ReplaceAllString(lower, " ")

	fillerPhrases := []string{
		"what products are available now", "what products are available",
		"what are available now", "what are available", "what products do you have",
		"what is available now", "what is available", "what do you have",
		"can you find me a", "can you find me an", "can you find me",
		"can you find a", "can you find an", "can you find",
		"can you show me all", "can you show me a", "can you show me an", "can you show me",
		"i want to buy a", "i want to buy an", "i want to buy",
		"i want to get a", "i want to get an", "i want to get",
		"i want to rent a", "i want to rent an", "i want to rent",
		"i want to find", "i want a", "i want an", "i want",
		"i need to buy a", "i need to buy", "i need to get",
		"i need to find a", "i need to find", "i need to rent a", "i need to rent",
		"i need a", "i need an", "i need",
		"looking for a", "looking for an", "looking for",
		"look for a", "look for an", "look for",
		"find me a", "find me an", "find me", "find a", "find an",
		"search for a", "search for an", "search for",
		"show me all", "show me a", "show me an", "show me",
		"please find me a", "please find me", "please find a", "please find",
		"where can i buy a", "where can i buy an", "where can i buy",
		"where can i find a", "where can i find an", "where can i find",
		"where can i get a", "where can i get an", "where can i get",
		"is there any", "is there a", "is there an", "is there",
		"are there any", "are there",
		"give me a", "give me an", "give me",
		"near me", "around me", "for rent", "for sale",
		"brand new", "second hand", "secondhand", "used only",
		// Amharic
		"መግዛት እፈልጋለሁ", "እፈልጋለሁ", "ማግኘት እፈልጋለሁ", "ፈልጉልኝ", "ፈልግልኝ",
		"ዋጋው ስንት ነው", "አሳዩኝ", "አሳየኝ", "ምን አለ", "ምን እቃ አለ",
		// Afaan Oromo
		"bituu barbaada", "bituu barbaadu", "barbaada", "barbaadu",
		"naaf barbaadi", "naaf barbaadaa", "na agarsiisi",
		"gatiin meeqa", "maal qabdu",
		// Tigrinya
		"ክገዝእ ደልየ", "ከቢ ደልየ", "ኣርእዩኒ", "ድለዩለይ",
		"ክንዲይ እዩ", "እንታይ ኣለኩም",
	}
	sort.Slice(fillerPhrases, func(i, j int) bool { return len(fillerPhrases[i]) > len(fillerPhrases[j]) })
	for _, phrase := range fillerPhrases {
		lower = strings.ReplaceAll(lower, phrase, " ")
	}

	for _, loc := range knownEthiopianLocations {
		lower = strings.ReplaceAll(lower, "in "+loc, " ")
		lower = strings.ReplaceAll(lower, "at "+loc, " ")
		lower = strings.ReplaceAll(lower, loc, " ")
	}

	stopWords := map[string]bool{
		"a": true, "an": true, "the": true, "in": true, "at": true, "on": true,
		"for": true, "to": true, "of": true, "with": true, "and": true, "or": true,
		"is": true, "are": true, "was": true, "be": true, "been": true,
		"now": true, "me": true, "my": true, "you": true, "i": true, "we": true,
		"please": true, "some": true, "any": true, "all": true, "only": true, "ones": true, "one": true,
		"used": true, "new": true, "product": true, "products": true, "item": true, "items": true,
		"listing": true, "listings": true, "price": true, "cost": true,
		"have": true, "has": true, "do": true, "does": true, "can": true, "will": true, "get": true, "got": true,
		"there": true,
		// Amharic
		"አለ": true, "አሉ": true, "ነው": true, "ነበር": true, "ናት": true,
		"ወይ": true, "ስለ": true, "ጋር": true, "ከ": true, "በ": true, "የ": true, "ወደ": true,
		"ምን": true, "የት": true, "ዋጋ": true, "ስንት": true,
		// Afaan Oromo
		"fi": true, "ykn": true, "jira": true, "jiru": true, "naaf": true, "na": true,
		"maal": true, "akkam": true, "gatii": true, "meeqa": true,
		// Tigrinya
		"እዩ": true, "ኣሎ": true, "ዶ": true, "ድዩ": true, "እንታይ": true,
	}

	words := strings.Fields(lower)
	var finalWords []string
	for _, w := range words {
		if strings.HasPrefix(w, "የ") && len([]rune(w)) > 2 {
			candidate := strings.TrimPrefix(w, "የ")
			if _, ok := ethiopianKeywordMap[candidate]; ok || ContainsGeez(candidate) {
				w = candidate
			}
		}
		if !stopWords[w] && len(w) > 0 {
			finalWords = append(finalWords, w)
		}
	}
	return strings.TrimSpace(strings.Join(finalWords, " "))
}
