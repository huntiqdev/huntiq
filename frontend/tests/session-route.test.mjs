import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';

const require = createRequire(import.meta.url);
const ts = require('typescript');
const { NextRequest } = require('next/server');

// Dynamically load the TypeScript route handler
const routePath = path.resolve('src/app/api/auth/session/route.ts');
const routeSource = fs.readFileSync(routePath, 'utf8');
const transpiled = ts.transpileModule(routeSource, {
  compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2022 },
});
const moduleScope = { exports: {} };
const runner = new Function('require', 'module', 'exports', transpiled.outputText);
runner(require, moduleScope, moduleScope.exports);

const { GET, POST, DELETE } = moduleScope.exports;

test('GET /api/auth/session: ignores forged query parameters token and session_token', async () => {
  const forgedUrl = 'http://localhost:3000/api/auth/session?token=forged.jwt.token&session_token=forged_session&redirect=/admin';
  const req = new NextRequest(forgedUrl, {
    method: 'GET',
    headers: {
      host: 'localhost:3000',
    },
  });

  const response = await GET(req);

  // Status must be 200 JSON, NOT a redirect
  assert.equal(response.status, 200);

  const data = await response.json();

  // Forged query tokens must NEVER be returned or accepted
  assert.equal(data.token, null, 'Forged ?token query parameter must not establish session token');
  assert.equal(data.session_token, null, 'Forged ?session_token query parameter must not establish admin session token');

  // Must not set authentication cookies via query parameters
  const setCookie = response.headers.get('set-cookie');
  assert.equal(setCookie, null, 'GET request with query params must not set session cookies');

  // Cache-Control must prevent caching
  const cacheControl = response.headers.get('cache-control');
  assert.ok(cacheControl?.includes('no-store'), 'Cache-Control header must include no-store');
});

test('GET /api/auth/session: reads legitimate cookies when present', async () => {
  const req = new NextRequest('http://localhost:3000/api/auth/session', {
    method: 'GET',
    headers: {
      cookie: 'hm_token=legitimate.jwt.token; hm_admin_session=legitimate_admin_session_token',
    },
  });

  const response = await GET(req);
  assert.equal(response.status, 200);

  const data = await response.json();
  assert.equal(data.token, 'legitimate.jwt.token');
  assert.equal(data.session_token, 'legitimate_admin_session_token');
});

test('POST /api/auth/session: legitimate login flow sets secure cookies', async () => {
  // Construct a dummy 3-part JWT
  const header = Buffer.from(JSON.stringify({ alg: 'HS256', typ: 'JWT' })).toString('base64url');
  const payload = Buffer.from(JSON.stringify({ user_id: 1, exp: Math.floor(Date.now() / 1000) + 3600 })).toString('base64url');
  const signature = 'test_signature';
  const dummyJwt = `${header}.${payload}.${signature}`;

  const req = new NextRequest('http://localhost:3000/api/auth/session', {
    method: 'POST',
    headers: {
      'content-type': 'application/json',
    },
    body: JSON.stringify({
      token: dummyJwt,
      session_token: 'test_admin_session_token',
    }),
  });

  const response = await POST(req);
  assert.equal(response.status, 200);

  const data = await response.json();
  assert.equal(data.ok, true);

  const setCookie = response.headers.get('set-cookie');
  assert.ok(setCookie, 'POST must set session cookies');
  assert.ok(setCookie.includes('hm_token='), 'Cookie hm_token must be set');
  assert.ok(setCookie.includes('hm_admin_session='), 'Cookie hm_admin_session must be set');
});

test('DELETE /api/auth/session: logout flow clears both cookies', async () => {
  const response = await DELETE();
  assert.equal(response.status, 200);

  const setCookie = response.headers.get('set-cookie');
  assert.ok(setCookie, 'DELETE must set clearing cookies');
  assert.ok(setCookie.includes('Max-Age=0') || setCookie.includes('max-age=0') || setCookie.includes('expires='), 'Cookies must expire immediately');
});
