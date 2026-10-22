'use client';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import { Spinner } from '@/components/ui/Spinner';
import { ShieldCheck, Mail, Lock, KeyRound } from 'lucide-react';

export default function AdminLoginPage() {
  const router = useRouter();
  const { user, login } = useAuth();
  
  const [step, setStep] = useState<1 | 2>(1);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [code, setCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (user?.role === 'admin') {
      router.replace('/admin/dashboard');
    }
  }, [user, router]);

  const handleStep1 = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await api.adminLogin({ email, password });
      if (res.is_super) {
        await login(res.token!);
        router.push('/admin/dashboard');
      } else {
        setStep(2);
      }
    } catch (err: any) {
      setError(err?.message || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const handleStep2 = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await api.adminLoginVerifyOTP({ email, code });
      await login(res.token);
      router.push('/admin/dashboard');
    } catch (err: any) {
      setError(err?.message || 'OTP verification failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center px-4 bg-zinc-950">
      <div className="w-full max-w-md bg-zinc-900 border border-zinc-800 rounded-3xl p-8 shadow-2xl relative overflow-hidden">
        {/* Glow effect */}
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-1/2 bg-orange-500/10 blur-[100px] pointer-events-none" />
        
        <div className="relative z-10">
          <div className="flex flex-col items-center mb-8">
            <div className="w-16 h-16 bg-gradient-to-br from-orange-500 to-red-600 rounded-2xl flex items-center justify-center shadow-lg shadow-orange-500/20 mb-4">
              <ShieldCheck className="text-white" size={32} />
            </div>
            <h1 className="text-2xl font-black text-white text-center">Admin Access</h1>
            <p className="text-zinc-400 text-sm mt-1 text-center font-semibold">
              {step === 1 ? 'Authorized personnel only' : 'Enter the OTP sent to your email'}
            </p>
          </div>

          {error && (
            <div className="bg-red-500/10 border border-red-500/20 text-red-400 text-sm px-4 py-3 rounded-xl mb-6 font-bold text-center">
              {error}
            </div>
          )}

          {step === 1 ? (
            <form onSubmit={handleStep1} className="space-y-4">
              <div>
                <label className="text-xs font-black text-zinc-500 uppercase tracking-widest mb-1 block">Email</label>
                <div className="relative">
                  <Mail className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
                  <input type="email" required
                    className="w-full bg-zinc-950 border border-zinc-800 text-white rounded-xl pl-12 pr-4 py-3 font-semibold focus:outline-none focus:border-orange-500 focus:ring-1 focus:ring-orange-500 transition-all"
                    value={email} onChange={e => setEmail(e.target.value)}
                  />
                </div>
              </div>
              <div>
                <label className="text-xs font-black text-zinc-500 uppercase tracking-widest mb-1 block">Password</label>
                <div className="relative">
                  <Lock className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
                  <input type="password" required
                    className="w-full bg-zinc-950 border border-zinc-800 text-white rounded-xl pl-12 pr-4 py-3 font-semibold focus:outline-none focus:border-orange-500 focus:ring-1 focus:ring-orange-500 transition-all"
                    value={password} onChange={e => setPassword(e.target.value)}
                  />
                </div>
              </div>
              <button disabled={loading}
                className="w-full bg-orange-600 hover:bg-orange-500 text-white font-black py-4 rounded-xl mt-4 shadow-lg shadow-orange-600/20 transition-all active:scale-[0.98] flex justify-center disabled:opacity-50">
                {loading ? <Spinner size={20} className="text-white" /> : 'Sign In'}
              </button>
            </form>
          ) : (
            <form onSubmit={handleStep2} className="space-y-4">
              <div>
                <label className="text-xs font-black text-zinc-500 uppercase tracking-widest mb-1 block">Authentication Code</label>
                <div className="relative">
                  <KeyRound className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
                  <input type="text" required
                    className="w-full bg-zinc-950 border border-zinc-800 text-white rounded-xl pl-12 pr-4 py-3 font-mono font-bold tracking-widest focus:outline-none focus:border-orange-500 focus:ring-1 focus:ring-orange-500 transition-all text-center text-lg"
                    value={code} onChange={e => setCode(e.target.value)}
                  />
                </div>
              </div>
              <button disabled={loading}
                className="w-full bg-orange-600 hover:bg-orange-500 text-white font-black py-4 rounded-xl mt-4 shadow-lg shadow-orange-600/20 transition-all active:scale-[0.98] flex justify-center disabled:opacity-50">
                {loading ? <Spinner size={20} className="text-white" /> : 'Verify Code'}
              </button>
              <button type="button" onClick={() => setStep(1)}
                className="w-full text-zinc-500 font-bold text-sm mt-4 hover:text-zinc-300">
                Back to Login
              </button>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}
