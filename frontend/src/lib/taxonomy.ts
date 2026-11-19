/**
 * Dynamic category taxonomy utilities
 * All categories are fetched from database - NO hardcoded data
 */

interface CategoryLike {
  name?: string;
  name_en?: string;
  name_am?: string;
  name_om?: string;
  name_ti?: string;
  slug?: string;
  i18n?: Record<string, string>;
}

/**
 * Resolves the category title in the requested language with proper fallback
 * @param category - Category object from database
 * @param locale - Language code (en, am, om/or, ti)
 * @returns Localized category name or slug as fallback
 */
export function resolveCategoryTitle(
  category: CategoryLike | null | undefined,
  locale: string = 'en'
): string {
  if (!category) return '';

  // Check i18n field first if available
  if (category.i18n && category.i18n[locale]) {
    return category.i18n[locale];
  }

  // Fallback to specific language fields based on locale
  switch (locale) {
    case 'am':
      return category.name_am || category.name || category.name_en || category.slug || '';
    case 'om':
    case 'or': // Oromo can be 'om' or 'or'
      return category.name_om || category.name || category.name_en || category.slug || '';
    case 'ti':
      return category.name_ti || category.name || category.name_en || category.slug || '';
    case 'en':
    default:
      return category.name_en || category.name || category.slug || '';
  }
}

/**
 * Gets the category name with proper formatting
 * @param category - Category object from database
 * @param locale - Language code
 * @returns Formatted category name
 */
export function getCategoryName(
  category: CategoryLike | null | undefined,
  locale: string = 'en'
): string {
  const title = resolveCategoryTitle(category, locale);
  // Clean up slug formatting if that's what we got
  return title.replace(/_/g, ' ').replace(/-/g, ' ');
}

/**
 * Traverses a category tree to find the full path of ancestor nodes leading to the target
 * @param nodes - Tree of category nodes with optional children
 * @param target - Target ID or slug to locate
 * @returns Array of category nodes in lineage order [root, ..., target]
 */
export function findCategoryPath<T extends { id?: number; slug?: string; children?: T[] }>(
  nodes: T[],
  target?: number | string | null
): T[] {
  if (!target && target !== 0) return [];
  const targetStr = String(target).trim().toLowerCase();
  const targetNum = typeof target === 'number' ? target : (!isNaN(Number(target)) ? Number(target) : undefined);

  for (const node of nodes) {
    const matchId = targetNum !== undefined && node.id === targetNum;
    const matchSlug = node.slug && node.slug.toLowerCase() === targetStr;
    const matchIdStr = node.id !== undefined && String(node.id) === targetStr;

    if (matchId || matchSlug || matchIdStr) {
      return [node];
    }
    if (node.children && node.children.length > 0) {
      const childPath = findCategoryPath(node.children, target);
      if (childPath.length > 0) {
        return [node, ...childPath];
      }
    }
  }
  return [];
}

/**
 * Locates a single category node anywhere in a tree
 */
export function findCategoryNode<T extends { id?: number; slug?: string; children?: T[] }>(
  nodes: T[],
  target?: number | string | null
): T | undefined {
  const path = findCategoryPath(nodes, target);
  return path.length > 0 ? path[path.length - 1] : undefined;
}

/**
 * Returns a high-affinity contextual emoji for visual cues
 */
export function getCategoryEmoji(slug: string = '', icon?: string): string {
  if (icon && icon.trim().length > 0 && icon.trim().length <= 4) return icon.trim();
  const s = (slug || '').toLowerCase();
  if (s.includes('phone') || s.includes('smartphone') || s.includes('mobile')) return '📱';
  if (s.includes('laptop') || s.includes('comput') || s.includes('pc')) return '💻';
  if (s.includes('tablet') || s.includes('ipad') || s.includes('ereader')) return '📟';
  if (s.includes('audio') || s.includes('sound') || s.includes('headphone') || s.includes('speaker')) return '🎧';
  if (s.includes('camera') || s.includes('photo') || s.includes('lens')) return '📷';
  if (s.includes('tv') || s.includes('television') || s.includes('display') || s.includes('video')) return '📺';
  if (s.includes('gaming') || s.includes('console') || s.includes('esport')) return '🎮';
  if (s.includes('wearable') || s.includes('watch')) return '⌚';
  if (s.includes('electronic')) return '⚡';
  if (s.includes('vehicle') || s.includes('automotiv') || s.includes('car')) return '🚗';
  if (s.includes('motorcycle') || s.includes('bike') || s.includes('scooter')) return '🏍️';
  if (s.includes('truck') || s.includes('lorry') || s.includes('commercial-vehicle')) return '🚚';
  if (s.includes('part') || s.includes('spare') || s.includes('tire') || s.includes('tyre')) return '⚙️';
  if (s.includes('property') || s.includes('real-estate') || s.includes('hous') || s.includes('home')) return '🏠';
  if (s.includes('apartment') || s.includes('flat') || s.includes('condo')) return '🏢';
  if (s.includes('land') || s.includes('plot')) return '🏞️';
  if (s.includes('cloth') || s.includes('apparel') || s.includes('fashion') || s.includes('dress') || s.includes('suit')) return '👗';
  if (s.includes('footwear') || s.includes('shoe') || s.includes('boot')) return '👟';
  if (s.includes('bag') || s.includes('luggage') || s.includes('backpack')) return '🎒';
  if (s.includes('jewel') || s.includes('gem') || s.includes('gold')) return '💎';
  if (s.includes('food') || s.includes('beverage') || s.includes('grocery') || s.includes('snack')) return '🍎';
  if (s.includes('coffee') || s.includes('tea') || s.includes('cafe')) return '☕';
  if (s.includes('grain') || s.includes('cereal') || s.includes('teff')) return '🌾';
  if (s.includes('meat') || s.includes('livestock') || s.includes('cattle')) return '🥩';
  if (s.includes('agri') || s.includes('farm') || s.includes('crop') || s.includes('seed')) return '🌱';
  if (s.includes('construction') || s.includes('build') || s.includes('cement')) return '🧱';
  if (s.includes('tool') || s.includes('drill') || s.includes('saw') || s.includes('hardware')) return '🔨';
  if (s.includes('machin') || s.includes('industrial') || s.includes('manufactur')) return '🏭';
  if (s.includes('furniture') || s.includes('chair') || s.includes('table') || s.includes('bed') || s.includes('sofa')) return '🛋️';
  if (s.includes('kitchen') || s.includes('appliance') || s.includes('cookware')) return '🍳';
  if (s.includes('beauty') || s.includes('personal-care') || s.includes('cosmetic') || s.includes('skincare')) return '💄';
  if (s.includes('health') || s.includes('medic') || s.includes('pharma') || s.includes('clinic')) return '💊';
  if (s.includes('sport') || s.includes('fitness') || s.includes('gym')) return '⚽';
  if (s.includes('book') || s.includes('stationery') || s.includes('educat') || s.includes('course')) return '📚';
  if (s.includes('baby') || s.includes('kid') || s.includes('toy')) return '🧸';
  if (s.includes('service') || s.includes('gig') || s.includes('work') || s.includes('job') || s.includes('consult')) return '💼';
  if (s.includes('energy') || s.includes('solar') || s.includes('utilit') || s.includes('battery')) return '🔋';
  if (s.includes('chemical') || s.includes('petro')) return '🧪';
  if (s.includes('event') || s.includes('exhibit') || s.includes('party')) return '🎉';
  if (s.includes('art') || s.includes('craft') || s.includes('antique')) return '🎨';
  if (s.includes('security') || s.includes('surveillance') || s.includes('cctv')) return '🛡️';
  return '📦';
}

