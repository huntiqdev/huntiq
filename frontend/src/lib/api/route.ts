import { NextRequest, NextResponse } from 'next/server';

const COOKIE = 'hm_token';
const SECURE = process.env.NODE_ENV === 'production';

// POST /api/auth/session — set token cookie after login
export async function POST(req: NextRequest) {
  const { token } = await req.json();
  if (!token || typeof token !== 'string' || token.split('.').length !== 3) {
    return NextResponse.json({ error: 'invalid token' }, { status: 400 });
  }
  let exp = Math.floor(Date.now() / 1000) + 60 * 60 * 24;
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    if (payload.exp) exp = payload.exp;
  } catch {}
  const maxAge = Math.max(0, exp - Math.floor(Date.now() / 1000));

  const res = NextResponse.json({ ok: true });
  res.cookies.set(COOKIE, token, {
    httpOnly: true,
    secure: SECURE,
    sameSite: 'lax',
    path: '/',
    maxAge: maxAge,
  });
  return res;
}

// DELETE /api/auth/session — clear cookie on logout
export async function DELETE() {
  const res = NextResponse.json({ ok: true });
  res.cookies.set(COOKIE, '', { httpOnly: true, maxAge: 0, path: '/' });
  return res;
}

// GET /api/auth/session — read token from cookie (server→client handoff on mount)
export async function GET(req: NextRequest) {
  const token = req.cookies.get(COOKIE)?.value ?? null;
  const res = NextResponse.json({ token });
  // Prevent browser/CDN from caching the session response
  res.headers.set('Cache-Control', 'no-store, no-cache, must-revalidate');
  return res;
}
