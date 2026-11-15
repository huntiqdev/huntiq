import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const langDir = path.resolve(__dirname, '../src/lib/language');

test('i18n translation files exist and are valid JSON', () => {
  const locales = ['en', 'am', 'or', 'ti'];
  for (const loc of locales) {
    const filePath = path.join(langDir, `${loc}.json`);
    assert.ok(fs.existsSync(filePath), `Translation file ${loc}.json must exist`);
    const content = fs.readFileSync(filePath, 'utf8');
    const parsed = JSON.parse(content);
    assert.ok(typeof parsed === 'object' && parsed !== null, `${loc}.json must parse into an object`);
  }
});

test('i18n contains essential navigation and header keys', () => {
  const enPath = path.join(langDir, 'en.json');
  const en = JSON.parse(fs.readFileSync(enPath, 'utf8'));

  assert.ok(en.nav, 'en.json must have nav section');
  assert.ok(en.nav.buy, 'en.json must have nav.buy');
  assert.ok(en.nav.rent, 'en.json must have nav.rent');
  assert.ok(en.nav.sell, 'en.json must have nav.sell');
});

test('non-English locales have valid top-level keys', () => {
  const locales = ['am', 'or', 'ti'];
  for (const loc of locales) {
    const filePath = path.join(langDir, `${loc}.json`);
    const data = JSON.parse(fs.readFileSync(filePath, 'utf8'));
    assert.ok(Object.keys(data).length > 0, `${loc}.json must have keys`);
  }
});

