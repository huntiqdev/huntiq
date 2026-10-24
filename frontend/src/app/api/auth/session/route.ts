import { NextRequest, NextResponse } from 'next/server';

const JWT_COOKIE = 'hm_token';
const ADMIN_SESSION_COOKIE = 'hm_admin_session';
const SECURE = process.env.NODE_ENV === 'production';

// POST /api/auth/session — set authentication cookies after login
export async function POST(req: NextRequest) {
  try {
    const body = await req.json();
    const token = body?.token;
    const sessionToken = body?.session_token;

    if (!token && !sessionToken) {
      return NextResponse.json({ error: 'token or session_token required' }, { status: 400 });
    }

    const res = NextResponse.json({ ok: true });

    // Set JWT token cookie if present
    if (token && typeof token === 'string' && token.split('.').length === 3) {
      let exp = Math.floor(Date.now() / 1000) + 60 * 60 * 24;
      try {
        const payload = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')));
        if (payload.exp) exp = payload.exp;
      } catch {}
      const maxAge = Math.max(0, exp - Math.floor(Date.now() / 1000));

      res.cookies.set(JWT_COOKIE, token, {
        httpOnly: true,
        secure: SECURE,
        sameSite: 'lax',
        path: '/',
        maxAge: maxAge,
      });
    }

    // Set Admin Session cookie if present (WebAuthn / Passkey enterprise session)
    if (sessionToken && typeof sessionToken === 'string') {
      res.cookies.set(ADMIN_SESSION_COOKIE, sessionToken, {
        httpOnly: true,
        secure: SECURE,
        sameSite: 'lax',
        path: '/',
        maxAge: 60 * 60 * 8, // 8 hour absolute maximum
      });
    }

    return res;
  } catch (err: any) {
    return NextResponse.json({ error: 'invalid session request' }, { status: 400 });
  }
}

// DELETE /api/auth/session — clear both cookies on logout
export async function DELETE() {
  const res = NextResponse.json({ ok: true });
  res.cookies.set(JWT_COOKIE, '', { httpOnly: true, maxAge: 0, path: '/' });
  res.cookies.set(ADMIN_SESSION_COOKIE, '', { httpOnly: true, maxAge: 0, path: '/' });
  return res;
}

// GET /api/auth/session — read tokens from cookies (server→client handoff on mount)
export async function GET(req: NextRequest) {
  let token = req.cookies.get(JWT_COOKIE)?.value ?? null;
  let sessionToken = req.cookies.get(ADMIN_SESSION_COOKIE)?.value ?? null;

  const qToken = req.nextUrl.searchParams.get('token');
  const qSession = req.nextUrl.searchParams.get('session_token');
  const redirectUrl = req.nextUrl.searchParams.get('redirect');

  if (process.env.NODE_ENV !== 'production' && (qToken || qSession)) {
    if (qToken) token = qToken;
    if (qSession) sessionToken = qSession;
  }

  const res = redirectUrl
    ? NextResponse.redirect(new URL(redirectUrl, req.url))
    : NextResponse.json({ token, session_token: sessionToken });

  if (process.env.NODE_ENV !== 'production' && (qToken || qSession)) {
    if (token) {
      res.cookies.set(JWT_COOKIE, token, {
        httpOnly: true,
        secure: SECURE,
        sameSite: 'lax',
        path: '/',
        maxAge: 60 * 60 * 24,
      });
    }
    if (sessionToken) {
      res.cookies.set(ADMIN_SESSION_COOKIE, sessionToken, {
        httpOnly: true,
        secure: SECURE,
        sameSite: 'lax',
        path: '/',
        maxAge: 60 * 60 * 8,
      });
    }
  }

  res.headers.set('Cache-Control', 'no-store, no-cache, must-revalidate');
  return res;
}

