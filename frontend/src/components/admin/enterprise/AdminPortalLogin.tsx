'use client';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { isWebAuthnSupported, isPlatformAuthenticatorAvailable, performPasskeyAuthentication } from '@/lib/webauthn';
import { Spinner } from '@/components/ui/Spinner';
import {
  Fingerprint,
  KeyRound,
  Mail,
  Lock,
  ShieldCheck,
  AlertCircle,
  CheckCircle2,
  Cpu,
  ArrowRight,
  ShieldAlert,
  Server,
  Users,
} from 'lucide-react';
import Link from 'next/link';

export interface AdminPortalLoginProps {
  portal: 'executive' | 'platform' | 'security' | 'staff';
  portalTitle: string;
  portalBadge: string;
  portalTheme: 'crimson' | 'indigo' | 'emerald' | 'orange';
  portalSubtitle: string;
  portalDescription: string;
  allowPasswordMfaFallback?: boolean;
}

const THEME_STYLES = {
  crimson: {
    badge: 'bg-red-500/10 text-red-400 border-red-500/30',
    iconBg: 'from-red-600 to-rose-700 shadow-red-600/30',
    primaryBtn: 'bg-gradient-to-r from-red-600 to-rose-600 hover:from-red-500 hover:to-rose-500 shadow-red-600/25',
    focusRing: 'focus:border-red-500 focus:ring-red-500',
    glow: 'bg-red-500/10',
    accentText: 'text-red-400',
    border: 'border-red-900/40',
  },
  indigo: {
    badge: 'bg-indigo-500/10 text-indigo-400 border-indigo-500/30',
    iconBg: 'from-indigo-600 to-blue-700 shadow-indigo-600/30',
    primaryBtn: 'bg-gradient-to-r from-indigo-600 to-blue-600 hover:from-indigo-500 hover:to-blue-500 shadow-indigo-600/25',
    focusRing: 'focus:border-indigo-500 focus:ring-indigo-500',
    glow: 'bg-indigo-500/10',
    accentText: 'text-indigo-400',
    border: 'border-indigo-900/40',
  },
  emerald: {
    badge: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30',
    iconBg: 'from-emerald-600 to-teal-700 shadow-emerald-600/30',
    primaryBtn: 'bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 shadow-emerald-600/25',
    focusRing: 'focus:border-emerald-500 focus:ring-emerald-500',
    glow: 'bg-emerald-500/10',
    accentText: 'text-emerald-400',
    border: 'border-emerald-900/40',
  },
  orange: {
    badge: 'bg-orange-500/10 text-orange-400 border-orange-500/30',
    iconBg: 'from-orange-600 to-amber-700 shadow-orange-600/30',
    primaryBtn: 'bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 shadow-orange-600/25',
    focusRing: 'focus:border-orange-500 focus:ring-orange-500',
    glow: 'bg-orange-500/10',
    accentText: 'text-orange-400',
    border: 'border-orange-900/40',
  },
};

export function AdminPortalLogin({
  portal,
  portalTitle,
  portalBadge,
  portalTheme,
  portalSubtitle,
  portalDescription,
  allowPasswordMfaFallback = false,
}: AdminPortalLoginProps) {
  const router = useRouter();
  const { user, login } = useAuth();

  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [biometricAvailable, setBiometricAvailable] = useState<boolean | null>(null);

  // Fallback password + OTP state (only active when enabled)
  const [authMode, setAuthMode] = useState<'passkey' | 'password'>('passkey');
  const [password, setPassword] = useState('');
  const [otpCode, setOtpCode] = useState('');
  const [otpStep, setOtpStep] = useState(false);

  const theme = THEME_STYLES[portalTheme];

  useEffect(() => {
    // If already logged in as admin, redirect to dashboard
    if (user?.role === 'admin') {
      router.replace('/admin/dashboard');
    }
  }, [user, router]);

  useEffect(() => {
    if (typeof window !== 'undefined') {
      isPlatformAuthenticatorAvailable().then((avail) => setBiometricAvailable(avail));
    }
  }, []);

  const handlePasskeyLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!isWebAuthnSupported()) {
      setError('WebAuthn / Passkey authentication is not supported by this browser. Please use Chrome, Safari, Edge, or Firefox.');
      return;
    }

    if (!email.trim()) {
      setError('Please enter your authorized administrative email address.');
      return;
    }

    setLoading(true);
    try {
      // 1. Begin ceremony with Go backend
      const beginRes = await api.adminPasskeyLoginBegin({
        email: email.trim().toLowerCase(),
        portal,
      });

      if (!beginRes || !beginRes.challenge_id || !beginRes.options) {
        throw new Error('Invalid challenge received from authentication server.');
      }

      // 2. Browser WebAuthn ceremony (Biometric / Security Key)
      const credentialAssertion = await performPasskeyAuthentication(beginRes.options);

      // 3. Verify assertion on Go backend
      const verifyRes = await api.adminPasskeyLoginVerify({
        challenge_id: beginRes.challenge_id,
        response: credentialAssertion,
      });

      if (!verifyRes || !verifyRes.token) {
        throw new Error('Authentication assertion verification failed.');
      }

      // 4. Set session token in AuthContext and cookie
      await login(verifyRes.token, verifyRes.session_token);

      // 5. Navigate to authorized portal
      router.push(verifyRes.redirect_url || '/admin/dashboard');
    } catch (err: any) {
      console.error('Passkey authentication error:', err);
      setError(err?.message || 'Authentication failed or access restricted. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handlePasswordStep1 = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await api.adminLogin({ email: email.trim().toLowerCase(), password });
      if (res.token && res.is_super) {
        await login(res.token);
        router.push('/admin/dashboard');
      } else {
        setOtpStep(true);
      }
    } catch (err: any) {
      setError(err?.message || 'Invalid credentials or access denied.');
    } finally {
      setLoading(false);
    }
  };

  const handlePasswordStep2 = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await api.adminLoginVerifyOTP({ email: email.trim().toLowerCase(), code: otpCode.trim() });
      await login(res.token);
      router.push('/admin/dashboard');
    } catch (err: any) {
      setError(err?.message || 'Multi-factor verification code is invalid or expired.');
    } finally {
      setLoading(false);
    }
  };

  const getPortalIcon = () => {
    switch (portal) {
      case 'executive':
        return <ShieldAlert className="text-white" size={32} />;
      case 'platform':
        return <Server className="text-white" size={32} />;
      case 'security':
        return <ShieldCheck className="text-white" size={32} />;
      case 'staff':
      default:
        return <Users className="text-white" size={32} />;
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center px-4 bg-zinc-950 py-12 relative overflow-hidden font-sans">
      {/* Background ambient lighting */}
      <div className={`absolute top-1/4 left-1/2 -translate-x-1/2 w-[600px] h-[350px] ${theme.glow} blur-[120px] pointer-events-none rounded-full`} />

      <div className="w-full max-w-lg bg-zinc-900/90 border border-zinc-800/80 backdrop-blur-xl rounded-3xl p-8 lg:p-10 shadow-2xl relative z-10">
        {/* Top Header Badge */}
        <div className="flex flex-col items-center mb-8">
          <div className="flex items-center gap-2 mb-4">
            <span className={`text-[11px] font-black uppercase tracking-widest px-3 py-1 rounded-full border ${theme.badge}`}>
              {portalBadge}
            </span>
          </div>

          <div className={`w-16 h-16 bg-gradient-to-br ${theme.iconBg} rounded-2xl flex items-center justify-center shadow-lg mb-4`}>
            {getPortalIcon()}
          </div>

          <h1 className="text-2xl lg:text-3xl font-black text-white text-center tracking-tight">
            {portalTitle}
          </h1>
          <p className="text-zinc-400 text-xs mt-1.5 text-center font-medium max-w-sm">
            {portalSubtitle}
          </p>
        </div>

        {/* Security Notification Banner */}
        <div className="bg-zinc-950/60 border border-zinc-800/80 rounded-2xl p-3.5 mb-6 text-zinc-300 text-xs flex items-start gap-3">
          <Cpu className={`shrink-0 ${theme.accentText} mt-0.5`} size={16} />
          <div className="space-y-0.5">
            <p className="font-bold text-white text-[11px] uppercase tracking-wider">
              Hardware Biometrics & FIDO2
            </p>
            <p className="text-zinc-400 text-[11px] leading-relaxed">
              Touch ID, Face ID, Windows Hello or Security Keys. Zero biometric data is stored on HUNTIQ servers.
            </p>
          </div>
        </div>

        {/* Error Alert */}
        {error && (
          <div className="bg-red-500/10 border border-red-500/20 text-red-400 text-xs px-4 py-3 rounded-xl mb-6 font-semibold flex items-center gap-2.5">
            <AlertCircle size={16} className="shrink-0" />
            <span>{error}</span>
          </div>
        )}

        {/* Primary Form: Passkey Authentication */}
        {authMode === 'passkey' && (
          <form onSubmit={handlePasskeyLogin} className="space-y-5">
            <div>
              <label htmlFor="admin-email" className="text-[11px] font-black text-zinc-400 uppercase tracking-wider mb-1.5 block">
                Privileged Identity Email
              </label>
              <div className="relative">
                <Mail className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
                <input
                  id="admin-email"
                  type="email"
                  required
                  autoComplete="username webauthn"
                  placeholder="admin@huntiq.internal"
                  className={`w-full bg-zinc-950/80 border border-zinc-800 text-white rounded-xl pl-12 pr-4 py-3.5 text-sm font-medium focus:outline-none focus:ring-1 ${theme.focusRing} transition-all placeholder:text-zinc-600`}
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  disabled={loading}
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className={`w-full py-4 rounded-xl text-white font-bold text-sm tracking-wide shadow-lg flex items-center justify-center gap-3 transition-all ${theme.primaryBtn} disabled:opacity-50 disabled:cursor-not-allowed`}
            >
              {loading ? (
                <>
                  <Spinner />
                  <span>Verifying Hardware Assertion...</span>
                </>
              ) : (
                <>
                  <Fingerprint size={20} />
                  <span>Sign in with Passkey / Biometrics</span>
                  <ArrowRight size={16} className="opacity-70" />
                </>
              )}
            </button>
          </form>
        )}

        {/* Optional Password + MFA Fallback for Staff */}
        {allowPasswordMfaFallback && authMode === 'password' && !otpStep && (
          <form onSubmit={handlePasswordStep1} className="space-y-4">
            <div>
              <label className="text-[11px] font-black text-zinc-400 uppercase tracking-wider mb-1 block">Staff Email</label>
              <div className="relative">
                <Mail className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
                <input
                  type="email"
                  required
                  placeholder="staff@huntiq.internal"
                  className={`w-full bg-zinc-950 border border-zinc-800 text-white rounded-xl pl-12 pr-4 py-3 text-sm font-medium focus:outline-none focus:ring-1 ${theme.focusRing}`}
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  disabled={loading}
                />
              </div>
            </div>
            <div>
              <label className="text-[11px] font-black text-zinc-400 uppercase tracking-wider mb-1 block">Password</label>
              <div className="relative">
                <Lock className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
                <input
                  type="password"
                  required
                  placeholder="••••••••••••"
                  className={`w-full bg-zinc-950 border border-zinc-800 text-white rounded-xl pl-12 pr-4 py-3 text-sm font-medium focus:outline-none focus:ring-1 ${theme.focusRing}`}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  disabled={loading}
                />
              </div>
            </div>
            <button
              type="submit"
              disabled={loading}
              className={`w-full py-3.5 rounded-xl text-white font-bold text-sm tracking-wide transition-all ${theme.primaryBtn} disabled:opacity-50 flex items-center justify-center gap-2`}
            >
              {loading ? <Spinner /> : <span>Continue with Password</span>}
            </button>
          </form>
        )}

        {allowPasswordMfaFallback && authMode === 'password' && otpStep && (
          <form onSubmit={handlePasswordStep2} className="space-y-4">
            <div>
              <label className="text-[11px] font-black text-zinc-400 uppercase tracking-wider mb-1 block">
                MFA Verification Code
              </label>
              <div className="relative">
                <KeyRound className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
                <input
                  type="text"
                  required
                  maxLength={6}
                  placeholder="123456"
                  className={`w-full bg-zinc-950 border border-zinc-800 text-white text-center tracking-[0.5em] text-lg font-bold rounded-xl py-3 focus:outline-none focus:ring-1 ${theme.focusRing}`}
                  value={otpCode}
                  onChange={(e) => setOtpCode(e.target.value)}
                  disabled={loading}
                />
              </div>
            </div>
            <button
              type="submit"
              disabled={loading}
              className={`w-full py-3.5 rounded-xl text-white font-bold text-sm tracking-wide transition-all ${theme.primaryBtn} disabled:opacity-50 flex items-center justify-center gap-2`}
            >
              {loading ? <Spinner /> : <span>Verify MFA & Enter</span>}
            </button>
          </form>
        )}

        {/* Fallback Toggle Switcher */}
        {allowPasswordMfaFallback && (
          <div className="mt-6 pt-5 border-t border-zinc-800 text-center">
            {authMode === 'passkey' ? (
              <button
                type="button"
                onClick={() => setAuthMode('password')}
                className="text-xs text-zinc-400 hover:text-white font-semibold transition-colors inline-flex items-center gap-1.5"
              >
                <Lock size={13} />
                <span>Sign in with Staff Password & MFA instead</span>
              </button>
            ) : (
              <button
                type="button"
                onClick={() => setAuthMode('passkey')}
                className={`text-xs ${theme.accentText} hover:underline font-bold transition-colors inline-flex items-center gap-1.5`}
              >
                <Fingerprint size={14} />
                <span>Return to Passkey / Biometrics Authentication</span>
              </button>
            )}
          </div>
        )}

        {/* Directory Navigator Link */}
        <div className="mt-8 pt-4 border-t border-zinc-900 flex items-center justify-between text-[11px] text-zinc-500">
          <span>HUNTIQ Privileged Security Matrix</span>
          <Link href="/admin/login" className="text-zinc-400 hover:text-white underline transition-colors">
            Portal Directory
          </Link>
        </div>
      </div>
    </div>
  );
}
