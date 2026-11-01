// I18nLocation: en is required, any other locale key is optional.
// This is open — adding a new language (e.g. zh, ar) just works automatically.
export interface I18nLocation {
  en: string;
  [locale: string]: string | undefined;
}

/**
 * Authoritative Ethiopian geographic name translations across all 4 official platform languages.
 * Backed by PostgreSQL location_hierarchy seed data.
 */
export const ETHIOPIAN_LOCATION_TRANSLATIONS: Record<string, Record<string, string>> = {
  // Regions & Chartered Cities
  'addis ababa': { am: 'አዲስ አበባ', or: 'Finfinnee', ti: 'ኣዲስ ኣበባ' },
  'dire dawa': { am: 'ድሬዳዋ', or: 'Dirree Dawaa', ti: 'ድሬዳዋ' },
  'oromia': { am: 'ኦሮሚያ', or: 'Oromiyaa', ti: 'ኦሮሚያ' },
  'amhara': { am: 'አማራ', or: 'Amaara', ti: 'ኣምሓራ' },
  'tigray': { am: 'ትግራይ', or: 'Tigraay', ti: 'ትግራይ' },
  'sidama': { am: 'ሲዳማ', or: 'Sidaama', ti: 'ሲዳማ' },
  'somali': { am: 'ሶማሌ', or: 'Somaalee', ti: 'ሶማሌ' },
  'afar': { am: 'አፋር', or: 'Afaar', ti: 'ዓፋር' },
  'benishangul-gumuz': { am: 'ቤኒሻንጉል ጉሙዝ', or: 'Benishaangul Gumuz', ti: 'ቤኒሻንጉል ጉሙዝ' },
  'gambela': { am: 'ጋምቤላ', or: 'Gambellaa', ti: 'ጋምቤላ' },
  'harari': { am: 'ሐረሪ', or: 'Hararii', ti: 'ሓረሪ' },
  'central ethiopia': { am: 'ማዕከላዊ ኢትዮጵያ', or: 'Itoophiyaa Giddugaleessaa', ti: 'ማእከላይ ኢትዮጵያ' },
  'south ethiopia': { am: 'ደቡብ ኢትዮጵያ', or: 'Itoophiyaa Kibbaa', ti: 'ደቡብ ኢትዮጵያ' },
  'southwest ethiopia': { am: 'ደቡብ ምዕራብ ኢትዮጵያ', or: 'Itoophiyaa Kibba Lixaa', ti: 'ደቡብ ምዕራብ ኢትዮጵያ' },

  // Addis Ababa Subcities
  'bole': { am: 'ቦሌ', or: 'Boolee', ti: 'ቦሌ' },
  'kirkos': { am: 'ቂርቆስ', or: 'Qirqoos', ti: 'ቂርቆስ' },
  'yeka': { am: 'የካ', or: 'Yakaa', ti: 'የካ' },
  'arada': { am: 'አራዳ', or: 'Aradaa', ti: 'ኣራዳ' },
  'lideta': { am: 'ልደታ', or: 'Lidetaa', ti: 'ልደታ' },
  'addis ketema': { am: 'አዲስ ከተማ', or: 'Magaalaa Haaraa', ti: 'ኣዲስ ከተማ' },
  'gullele': { am: 'ጉለሌ', or: 'Gullallee', ti: 'ጉለሌ' },
  'kolfe keranio': { am: 'ኮልፌ ቀራኒዮ', or: 'Kolfi Qaraaniyoo', ti: 'ኮልፌ ቀራንዮ' },
  'nifas silk-lafto': { am: 'ንፋስ ስልክ ላፍቶ', or: 'Nifas Silk Laaftoo', ti: 'ንፋስ ስልኪ ላፍቶ' },
  'nifas silk': { am: 'ንፋስ ስልክ', or: 'Nifas Silk', ti: 'ንፋስ ስልኪ' },
  'akaki kality': { am: 'አቃቂ ቃሊቲ', or: 'Aqaaqii Qaallittii', ti: 'ኣቓቒ ቓሊቲ' },
  'lemi kura': { am: 'ለሚ ኩራ', or: 'Lammii Kuraa', ti: 'ለሚ ኩራ' },

  // Major Economic Hubs & Cities
  'adama': { am: 'አዳማ', or: 'Adaamaa', ti: 'ኣዳማ' },
  'hawassa': { am: 'ሐዋሳ', or: 'Hawaasaa', ti: 'ሓዋሳ' },
  'bahir dar': { am: 'ባሕር ዳር', or: 'Baahir Daar', ti: 'ባሕር ዳር' },
  'mekelle': { am: 'መቐለ', or: 'Maqalee', ti: 'መቐለ' },
  'gondar': { am: 'ጎንደር', or: 'Goondar', ti: 'ጎንደር' },
  'bishoftu': { am: 'ቢሾፍቱ', or: 'Bishooftuu', ti: 'ቢሾፍቱ' },
  'jimma': { am: 'ጅማ', or: 'Jimmaa', ti: 'ጅማ' },
  'dessie': { am: 'ደሴ', or: 'Dasee', ti: 'ደሴ' },
  'shashemene': { am: 'ሻሸመኔ', or: 'Shaashamannee', ti: 'ሻሸመኔ' },
  'jigjiga': { am: 'ጅግጅጋ', or: 'Jigjiga', ti: 'ጅግጅጋ' },
  'harar': { am: 'ሐረር', or: 'Harar', ti: 'ሓረር' },
  'arba minch': { am: 'አርባ ምንጭ', or: 'Arba Minch', ti: 'ኣርባ ምንጭ' },
  'wolkite': { am: 'ወልቂጤ', or: 'Walqixxee', ti: 'ወልቂጤ' },
  'hosanna': { am: 'ሆሳዕና', or: 'Hosaa\'inaa', ti: 'ሆሳዕና' },
  'debre markos': { am: 'ደብረ ማርቆስ', or: 'Dabra Maarqos', ti: 'ደብረ ማርቆስ' },
  'debre berhan': { am: 'ደብረ ብርሃን', or: 'Dabra Birhaan', ti: 'ደብረ ብርሃን' },
  'kombolcha': { am: 'ኮምቦልቻ', or: 'Kombolcha', ti: 'ኮምቦልቻ' },
  'axum': { am: 'አክሱም', or: 'Aksum', ti: 'ኣክሱም' },
  'shire': { am: 'ሽሬ', or: 'Shiraa', ti: 'ሽሬ' },
  'nekemte': { am: 'ነቀምቴ', or: 'Naqamtee', ti: 'ነቀምቴ' },
  'sebeta': { am: 'ሰበታ', or: 'Sabataa', ti: 'ሰበታ' },
  'burayu': { am: 'ቡራዩ', or: 'Buraayyuu', ti: 'ቡራዩ' },
  'dukem': { am: 'ዱከም', or: 'Dukam', ti: 'ዱከም' },
  'bonga': { am: 'ቦንጋ', or: 'Boongaa', ti: 'ቦንጋ' },
  'mizan teferi': { am: 'ሚዛን ተፈሪ', or: 'Mizaan Tafarii', ti: 'ሚዛን ተፈሪ' },
  'dilla': { am: 'ዲላ', or: 'Dillaa', ti: 'ዲላ' },
  'sodo': { am: 'ሶዶ', or: 'Soddoo', ti: 'ሶዶ' },
  'east shewa': { am: 'ምስራቅ ሸዋ', or: 'Shawaa Bahaa', ti: 'ምብራቕ ሸዋ' },
  'north gondar': { am: 'ሰሜን ጎንደር', or: 'Goondar Kaabaa', ti: 'ሰሜን ጎንደር' },
};

/**
 * Returns the best display name for a location item given the current locale.
 * Falls back to English if the locale translation is missing.
 * Extensible: works for any locale key (am, or, ti, zh, ar, fr, etc.)
 */
export function locName(item: I18nLocation, locale: string): string {
  const norm = (locale === 'om' ? 'or' : locale) || 'en';
  if (norm && norm !== 'en' && item[norm]) return item[norm] as string;
  if (norm && norm !== 'en' && item.en) {
    const canonical = ETHIOPIAN_LOCATION_TRANSLATIONS[item.en.trim().toLowerCase()];
    if (canonical?.[norm]) return canonical[norm];
  }
  return item.en;
}

/**
 * Extract deduplicated and localized location parts from a row or location object.
 * Eliminates awkward duplicate hierarchies (e.g., Addis Ababa › Addis Ababa › Bole)
 * and formats properly in all 4 Ethiopian languages.
 */
export function formatLocationParts(row: Record<string, any>, locale: string): string[] {
  if (!row) return [];
  const normLocale = (locale === 'om' ? 'or' : locale) || 'en';

  const pick = (en: string | null | undefined, suffix: string): string => {
    if (!en) return '';
    const trimmed = typeof en === 'string' ? en.trim() : '';
    if (!trimmed) return '';

    const locKey = `${suffix}_${normLocale}`;
    if (normLocale !== 'en' && row[locKey]) {
      return typeof row[locKey] === 'string' ? row[locKey].trim() : '';
    }

    // Canonical dictionary lookup for English names
    if (normLocale !== 'en') {
      const canonical = ETHIOPIAN_LOCATION_TRANSLATIONS[trimmed.toLowerCase()];
      if (canonical?.[normLocale]) {
        return canonical[normLocale];
      }
      // Woreda localized formatting (e.g. Woreda 03 -> ወረዳ 03)
      if (suffix === 'woreda' || /^woreda\s+/i.test(trimmed)) {
        if (normLocale === 'am' || normLocale === 'ti') {
          return trimmed.replace(/^woreda\s+/i, 'ወረዳ ');
        }
      }
    }

    return trimmed;
  };

  const r = pick(row.region, 'region');
  const z = pick(row.zone, 'zone');
  const s = pick(row.subcity, 'subcity');
  const w = pick(row.woreda, 'woreda');
  const k = pick(row.kebele, 'kebele');

  const parts: string[] = [];
  const seen = new Set<string>();

  const add = (text: string) => {
    if (!text) return;
    const key = text.trim().toLowerCase();
    // Do not repeat if already added or identical
    if (!seen.has(key)) {
      seen.add(key);
      parts.push(text.trim());
    }
  };

  add(r);
  add(z);
  add(s);
  add(w);
  add(k);

  return parts;
}

/**
 * Build a human-readable location path string from a search result row.
 * Picks the right language for each level and eliminates duplicate levels.
 */
export function locPath(row: Record<string, any>, locale: string): string {
  const parts = formatLocationParts(row, locale);
  return parts.join(' › ');
}

export interface RecommendedHub {
  key: string;
  region: string;
  subcity?: string;
  zone?: string;
  name: {
    en: string;
    am: string;
    or: string;
    ti: string;
  };
}

/**
 * Primary commercial hubs across Ethiopia for reasonable quick recommendations
 * without confusing users with nested repetitions.
 */
export const RECOMMENDED_ETHIOPIA_HUBS: RecommendedHub[] = [
  {
    key: 'addis_all',
    region: 'Addis Ababa',
    name: {
      en: 'Addis Ababa (All)',
      am: 'አዲስ አበባ (ሁሉም)',
      or: 'Finfinnee (Hunda)',
      ti: 'ኣዲስ ኣበባ (ኩሉ)',
    },
  },
  {
    key: 'addis_bole',
    region: 'Addis Ababa',
    subcity: 'Bole',
    name: {
      en: 'Bole, Addis Ababa',
      am: 'ቦሌ፣ አዲስ አበባ',
      or: 'Boolee, Finfinnee',
      ti: 'ቦሌ፣ ኣዲስ ኣበባ',
    },
  },
  {
    key: 'addis_kirkos',
    region: 'Addis Ababa',
    subcity: 'Kirkos',
    name: {
      en: 'Kirkos / Mexico, Addis',
      am: 'ቂርቆስ / ሜክሲኮ፣ አዲስ',
      or: 'Qirqoos / Meeksiikoo',
      ti: 'ቂርቆስ / ሜክሲኮ',
    },
  },
  {
    key: 'addis_cmc',
    region: 'Addis Ababa',
    subcity: 'Lemi Kura',
    name: {
      en: 'CMC / Lemi Kura',
      am: 'ሲኤምሲ / ለሚ ኩራ',
      or: 'CMC / Lammii Kuraa',
      ti: 'ሲኤምሲ / ለሚ ኩራ',
    },
  },
  {
    key: 'addis_arada',
    region: 'Addis Ababa',
    subcity: 'Arada',
    name: {
      en: 'Piassa / Arada',
      am: 'ፒያሳ / አራዳ',
      or: 'Piyaasaa / Araadaa',
      ti: 'ፒያሳ / ኣራዳ',
    },
  },
  {
    key: 'adama',
    region: 'Oromia',
    zone: 'East Shewa',
    name: {
      en: 'Adama / Nazret',
      am: 'አዳማ / ናዝሬት',
      or: 'Adaamaa',
      ti: 'ኣዳማ',
    },
  },
  {
    key: 'hawassa',
    region: 'Sidama',
    zone: 'Hawassa',
    name: {
      en: 'Hawassa',
      am: 'ሐዋሳ',
      or: 'Hawaasaa',
      ti: 'ሓዋሳ',
    },
  },
  {
    key: 'bahir_dar',
    region: 'Amhara',
    zone: 'Bahir Dar',
    name: {
      en: 'Bahir Dar',
      am: 'ባሕር ዳር',
      or: 'Baahir Daar',
      ti: 'ባሕር ዳር',
    },
  },
  {
    key: 'mekelle',
    region: 'Tigray',
    zone: 'Mekelle',
    name: {
      en: 'Mekelle',
      am: 'መቐለ',
      or: 'Maqalee',
      ti: 'መቐለ',
    },
  },
  {
    key: 'dire_dawa',
    region: 'Dire Dawa',
    name: {
      en: 'Dire Dawa',
      am: 'ድሬዳዋ',
      or: 'Dirre Dawaa',
      ti: 'ድሬዳዋ',
    },
  },
  {
    key: 'bishoftu',
    region: 'Oromia',
    zone: 'East Shewa',
    name: {
      en: 'Bishoftu / Debre Zeyit',
      am: 'ቢሾፍቱ / ደብረ ዘይት',
      or: 'Bishooftuu',
      ti: 'ቢሾፍቱ',
    },
  },
  {
    key: 'gondar',
    region: 'Amhara',
    zone: 'North Gondar',
    name: {
      en: 'Gondar',
      am: 'ጎንደር',
      or: 'Goondar',
      ti: 'ጎንደር',
    },
  },
];
