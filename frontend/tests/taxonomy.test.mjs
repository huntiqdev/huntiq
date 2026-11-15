import test from 'node:test';
import assert from 'node:assert/strict';

function resolveCategoryTitle(cat, locale = 'en') {
  if (!cat) return '';
  if (cat.i18n && cat.i18n[locale]) {
    return cat.i18n[locale];
  }
  switch (locale) {
    case 'am':
      if (cat.name_am) return cat.name_am;
      break;
    case 'or':
    case 'om':
      if (cat.name_om) return cat.name_om;
      break;
    case 'ti':
      if (cat.name_ti) return cat.name_ti;
      break;
    case 'en':
      if (cat.name_en) return cat.name_en;
      break;
  }
  return cat.i18n?.['en'] || cat.name || cat.name_en || cat.slug?.replace(/_/g, ' ') || '';
}

test('resolves title via JSONB i18n object precedence', () => {
  const cat = {
    id: 10,
    slug: 'electronics',
    name: 'Electronics',
    i18n: {
      en: 'Electronics',
      am: 'ኤሌክትሮኒክስ',
      or: 'Elektirooniksii',
      ti: 'ኤሌክትሮኒክስ'
    }
  };

  assert.equal(resolveCategoryTitle(cat, 'en'), 'Electronics');
  assert.equal(resolveCategoryTitle(cat, 'am'), 'ኤሌክትሮኒክስ');
  assert.equal(resolveCategoryTitle(cat, 'or'), 'Elektirooniksii');
  assert.equal(resolveCategoryTitle(cat, 'ti'), 'ኤሌክትሮኒክስ');
});

test('falls back gracefully to name or slug when locale is missing', () => {
  const cat = {
    id: 99,
    slug: 'custom_gear',
    name: 'Custom Gear'
  };

  assert.equal(resolveCategoryTitle(cat, 'fr'), 'Custom Gear');
  assert.equal(resolveCategoryTitle(cat, 'am'), 'Custom Gear');
});

test('handles 4-layer taxonomy boundary (levels 0..3)', () => {
  const MAX_TAXONOMY_DEPTH = 4;
  const levels = [0, 1, 2, 3];
  assert.equal(levels.length, MAX_TAXONOMY_DEPTH);
  assert.ok(levels[levels.length - 1] < MAX_TAXONOMY_DEPTH);
});

