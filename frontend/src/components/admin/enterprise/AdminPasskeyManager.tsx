'use client';
import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { api } from '@/lib/api';
import { isWebAuthnSupported, isPlatformAuthenticatorAvailable, performPasskeyRegistration } from '@/lib/webauthn';
import { Spinner } from '@/components/ui/Spinner';
import {
  Fingerprint,
  KeyRound,
  ShieldCheck,
  Plus,
  Trash2,
  AlertCircle,
  CheckCircle2,
  Cpu,
  Laptop,
  Smartphone,
  Shield,
  Clock,
} from 'lucide-react';

export function AdminPasskeyManager() {
  const { user } = useAuth();
  const [authenticators, setAuthenticators] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [registering, setRegistering] = useState(false);
  const [friendlyName, setFriendlyName] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [platformAvail, setPlatformAvail] = useState<boolean | null>(null);

  const loadAuthenticators = useCallback(async () => {
    if (!user?.token) return;
    setLoading(true);
    setError('');
    try {
      const res = await api.adminPasskeyList(user.token);
      setAuthenticators(res?.authenticators || []);
    } catch (err: any) {
      console.error('Failed to load authenticators:', err);
      // Non-fatal if no passkeys yet
    } finally {
      setLoading(false);
    }
  }, [user?.token]);

  useEffect(() => {
    loadAuthenticators();
    if (typeof window !== 'undefined') {
      isPlatformAuthenticatorAvailable().then((avail) => setPlatformAvail(avail));
    }
  }, [loadAuthenticators]);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user?.token) return;
    setError('');
    setSuccess('');

    if (!isWebAuthnSupported()) {
      setError('WebAuthn is not supported on this browser or platform.');
      return;
    }

    const name = friendlyName.trim() || 'My Passkey Authenticator';
    setRegistering(true);

    try {
      // 1. Begin registration with backend
      const beginRes = await api.adminPasskeyRegisterBegin(user.token, name);
      if (!beginRes || !beginRes.challenge_id || !beginRes.options) {
        throw new Error('Invalid registration challenge received from server.');
      }

      // 2. Perform WebAuthn ceremony with browser hardware authenticator
      const attestationResponse = await performPasskeyRegistration(beginRes.options);

      // 3. Verify attestation on Go backend
      const verifyRes = await api.adminPasskeyRegisterVerify(
        user.token,
        beginRes.challenge_id,
        name,
        attestationResponse
      );

      if (!verifyRes || !verifyRes.success) {
        throw new Error(verifyRes?.message || 'Failed to verify passkey registration.');
      }

      setSuccess('Passkey registered successfully! You can now use your biometric sensor or security key to log in.');
      setFriendlyName('');
      setIsModalOpen(false);
      await loadAuthenticators();
    } catch (err: any) {
      console.error('Passkey registration error:', err);
      setError(err?.message || 'Failed to register passkey. Please ensure your device authenticator is ready.');
    } finally {
      setRegistering(false);
    }
  };

  const handleDelete = async (id: number, name: string) => {
    if (!user?.token) return;
    if (!confirm(`Are you sure you want to revoke "${name}"? You will no longer be able to use this authenticator.`)) {
      return;
    }

    try {
      await api.adminPasskeyDelete(user.token, id);
      setSuccess(`Passkey "${name}" revoked.`);
      await loadAuthenticators();
    } catch (err: any) {
      setError(err?.message || 'Failed to revoke authenticator.');
    }
  };

  return (
    <div className="space-y-6">
      {/* Overview Banner */}
      <div className="bg-zinc-900 border border-zinc-800 rounded-3xl p-6 lg:p-8 relative overflow-hidden">
        <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          <div className="space-y-1">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-black uppercase tracking-wider">
              <ShieldCheck size={14} />
              <span>Phishing-Resistant FIDO2 / WebAuthn</span>
            </div>
            <h2 className="text-xl lg:text-2xl font-black text-white">Registered Passkeys & Authenticators</h2>
            <p className="text-zinc-400 text-xs max-w-xl">
              Biometric sensors (Touch ID, Windows Hello, Face ID) and external security keys (YubiKey) provide cryptographic user verification without storing passwords or biometric data.
            </p>
          </div>

          <button
            onClick={() => {
              setError('');
              setSuccess('');
              setIsModalOpen(true);
            }}
            className="flex items-center gap-2 px-5 py-3 bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 text-white font-black text-xs rounded-xl shadow-lg shadow-orange-600/20 transition-all shrink-0"
          >
            <Plus size={16} />
            <span>Enroll New Passkey</span>
          </button>
        </div>

        {/* Biometric Privacy Callout */}
        <div className="mt-6 pt-5 border-t border-zinc-800/80 grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
          <div className="flex items-center gap-2.5 text-zinc-400">
            <Cpu size={16} className="text-orange-400 shrink-0" />
            <span>On-device biometric verification</span>
          </div>
          <div className="flex items-center gap-2.5 text-zinc-400">
            <Shield size={16} className="text-emerald-400 shrink-0" />
            <span>Zero raw biometrics stored</span>
          </div>
          <div className="flex items-center gap-2.5 text-zinc-400">
            <KeyRound size={16} className="text-indigo-400 shrink-0" />
            <span>Cryptographic assertion signatures</span>
          </div>
        </div>
      </div>

      {/* Status Messages */}
      {error && (
        <div className="bg-red-500/10 border border-red-500/20 text-red-400 text-xs px-4 py-3 rounded-xl font-semibold flex items-center gap-2.5">
          <AlertCircle size={16} className="shrink-0" />
          <span>{error}</span>
        </div>
      )}

      {success && (
        <div className="bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs px-4 py-3 rounded-xl font-semibold flex items-center gap-2.5">
          <CheckCircle2 size={16} className="shrink-0" />
          <span>{success}</span>
        </div>
      )}

      {/* Authenticators List */}
      <div className="bg-zinc-900 border border-zinc-800 rounded-3xl overflow-hidden shadow-xl">
        <div className="p-5 border-b border-zinc-800 flex items-center justify-between">
          <h3 className="font-extrabold text-sm text-white">Active Authenticators ({authenticators.length})</h3>
          <span className="text-[11px] font-mono text-zinc-500">Tier-0 Hardware Trust</span>
        </div>

        {loading ? (
          <div className="h-40 flex items-center justify-center">
            <Spinner />
          </div>
        ) : authenticators.length === 0 ? (
          <div className="p-10 text-center space-y-3">
            <div className="w-12 h-12 bg-zinc-800 rounded-2xl flex items-center justify-center text-zinc-500 mx-auto">
              <Fingerprint size={24} />
            </div>
            <p className="text-sm font-bold text-zinc-400">No passkeys enrolled yet.</p>
            <p className="text-xs text-zinc-500 max-w-md mx-auto">
              Enroll your laptop fingerprint reader, mobile phone, or hardware security key to enable instant, phishing-resistant privileged logins.
            </p>
          </div>
        ) : (
          <div className="divide-y divide-zinc-800">
            {authenticators.map((auth) => (
              <div key={auth.id} className="p-5 flex items-center justify-between hover:bg-zinc-800/40 transition-colors">
                <div className="flex items-center gap-4">
                  <div className="w-10 h-10 bg-orange-600/10 border border-orange-500/20 rounded-xl flex items-center justify-center text-orange-400 shrink-0">
                    <Fingerprint size={20} />
                  </div>
                  <div>
                    <h4 className="font-bold text-sm text-white">{auth.friendly_name}</h4>
                    <div className="flex items-center gap-3 text-[11px] text-zinc-400 mt-1 font-mono">
                      <span>ID: {auth.credential_id_preview}</span>
                      <span>•</span>
                      <span className="flex items-center gap-1">
                        <Clock size={11} />
                        Enrolled: {new Date(auth.created_at).toLocaleDateString()}
                      </span>
                      {auth.last_used_at && (
                        <>
                          <span>•</span>
                          <span className="text-emerald-400">
                            Last used: {new Date(auth.last_used_at).toLocaleDateString()}
                          </span>
                        </>
                      )}
                    </div>
                  </div>
                </div>

                <div className="flex items-center gap-3">
                  <span className="text-[10px] font-black uppercase tracking-wider px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                    UV Enforced
                  </span>
                  <button
                    onClick={() => handleDelete(auth.id, auth.friendly_name)}
                    className="p-2 text-zinc-500 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-colors"
                    title="Revoke passkey"
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Enrollment Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-md">
          <div className="w-full max-w-md bg-zinc-950 border border-zinc-800 rounded-3xl p-6 shadow-2xl space-y-5 animate-in zoom-in-95">
            <div className="flex items-center justify-between border-b border-zinc-800/80 pb-4">
              <div className="flex items-center gap-2.5">
                <div className="p-2 bg-orange-600/10 text-orange-500 rounded-xl border border-orange-500/20">
                  <Fingerprint size={20} />
                </div>
                <div>
                  <h3 className="text-base font-black text-white">Enroll Hardware Passkey</h3>
                  <p className="text-xs text-zinc-400">Touch ID, Windows Hello, or Security Key</p>
                </div>
              </div>
              <button
                onClick={() => setIsModalOpen(false)}
                className="text-zinc-500 hover:text-white text-sm font-bold"
              >
                ✕
              </button>
            </div>

            <form onSubmit={handleRegister} className="space-y-4">
              <div>
                <label className="text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5 block">
                  Authenticator Friendly Name
                </label>
                <input
                  type="text"
                  required
                  placeholder="e.g. Solomon MacBook Touch ID / YubiKey 5C"
                  value={friendlyName}
                  onChange={(e) => setFriendlyName(e.target.value)}
                  className="w-full bg-zinc-900 border border-zinc-800 rounded-xl px-4 py-3 text-sm text-white focus:outline-none focus:border-orange-500 transition-all font-medium"
                  disabled={registering}
                />
                <p className="text-[11px] text-zinc-500 mt-1.5">
                  Give this authenticator a name so you can identify it in access logs.
                </p>
              </div>

              <div className="bg-zinc-900/60 border border-zinc-800 rounded-xl p-3 text-zinc-400 text-xs space-y-1">
                <p className="font-bold text-white text-[11px]">When you click Enroll:</p>
                <p className="text-[11px] leading-relaxed">
                  Your browser will prompt you to verify with your fingerprint, face recognition, PIN, or tap your hardware key.
                </p>
              </div>

              <div className="flex items-center justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  disabled={registering}
                  className="px-4 py-2.5 rounded-xl text-xs font-bold text-zinc-400 hover:text-white hover:bg-zinc-900 transition-all"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={registering}
                  className="px-5 py-2.5 bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-500 hover:to-amber-500 text-white font-bold text-xs rounded-xl shadow-lg shadow-orange-600/20 transition-all flex items-center gap-2 disabled:opacity-50"
                >
                  {registering ? (
                    <>
                      <Spinner />
                      <span>Waiting for Device Sensor...</span>
                    </>
                  ) : (
                    <>
                      <Fingerprint size={16} />
                      <span>Start Device Ceremony</span>
                    </>
                  )}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
