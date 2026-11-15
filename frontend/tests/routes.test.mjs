import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const appDir = path.resolve(__dirname, '../src/app');

test('all redirect routes use Next.js server redirect without heavy client bundles', () => {
  const redirectRoutes = [
    { file: 'buy/page.tsx', target: '/products?type=sale' },
    { file: 'rent/page.tsx', target: '/products?type=rent' },
    { file: 'find/page.tsx', target: '/products' },
    { file: 'demand/page.tsx', target: '/products' },
    { file: 'sell/page.tsx', target: '/seller/product/new' },
    { file: 'list/page.tsx', target: '/seller/product/new' },
    { file: 'buyer/login/page.tsx', target: '/auth?as=buyer' },
    { file: 'buyer/register/page.tsx', target: '/auth?as=buyer&tab=register' },
    { file: 'seller/login/page.tsx', target: '/auth?as=seller' },
    { file: 'seller/register/page.tsx', target: '/auth?as=seller&tab=register' },
  ];

  for (const { file, target } of redirectRoutes) {
    const fullPath = path.join(appDir, file);
    assert.ok(fs.existsSync(fullPath), `Route file ${file} must exist`);
    const content = fs.readFileSync(fullPath, 'utf8');
    assert.ok(content.includes('redirect('), `${file} must call redirect()`);
    assert.ok(content.includes(target), `${file} must redirect to ${target}`);
    assert.ok(!content.includes("'use client'"), `${file} should be a server component (not use client)`);
  }
});

