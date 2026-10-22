'use client';
import { useState, useRef, useEffect } from 'react';
import { ShieldCheck, RefreshCw } from 'lucide-react';
import { useLocale } from '@/hooks/useLocale';

interface Props {
  recipient: string;  // phone number or email address — shown in the "sent to" label
  onVerified: (code: string) => void;
  onResend: () => void;
  error?: string;
  loading?: boolean;
}

export function OTPInput({ recipient, onVerified, onResend, error, loading }: Props) {
  const { t } = useLocale();
  const [digits, setDigits] = useState(['', '', '', '', '', '']);
  const [countdown, setCountdown] = useState(300);
  const refs = useRef<(HTMLInputElement | null)[]>([]);
  const submittedRef = useRef(false); // prevents double-submit on fast input

  // Reset digits + timer whenever recipient changes (component re-use)
  useEffect(() => {
    setDigits(['', '', '', '', '', '']);
    setCountdown(300);
    submittedRef.current = false;
    refs.current[0]?.focus();
    const t = setInterval(() => setCountdown(c => Math.max(0, c - 1)), 1000);
    return () => clearInterval(t);
  }, [recipient]);

  // Reset submitted guard when loading clears (e.g. wrong code, user can retry)
  useEffect(() => {
    if (!loading) submittedRef.current = false;
  }, [loading]);

  const trySubmit = (code: string) => {
    if (submittedRef.current || loading) return;
    submittedRef.current = true;
    onVerified(code);
  };

  const handleChange = (i: number, val: string) => {
    const d = val.replace(/\D/g, '').slice(-1);
    const next = [...digits];
    next[i] = d;
    setDigits(next);
    if (d && i < 5) refs.current[i + 1]?.focus();
    if (next.every(x => x) && next.join('').length === 6) {
      trySubmit(next.join(''));
    }
  };

  const handleKeyDown = (i: number, e: React.KeyboardEvent) => {
    if (e.key === 'Backspace' && !digits[i] && i > 0) {
      refs.current[i - 1]?.focus();
    }
  };

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault();
    const text = e.clipboardData.getData('text').replace(/\D/g, '').slice(0, 6);
    if (text.length !== 6) return;
    const next = text.split('');
    setDigits(next);
    refs.current[5]?.focus();
    trySubmit(text);
  };

  const mm = String(Math.floor(countdown / 60)).padStart(2, '0');
  const ss = String(countdown % 60).padStart(2, '0');

  return (
    <div className="flex flex-col gap-4">
      <div className="text-center">
        <div className="w-12 h-12 bg-orange-100 dark:bg-orange-900/30 rounded-2xl flex items-center justify-center mx-auto mb-3">
          <ShieldCheck size={22} className="text-orange-600" />
        </div>
        <p className="font-black text-lg">{t.ui.enter_otp}</p>
        <p className="text-xs text-zinc-500 mt-1">
          {t.ui.sent_to} <span className="font-bold text-zinc-700 dark:text-zinc-300">{recipient}</span>
        </p>
      </div>

      <div className="flex gap-2 justify-center" onPaste={handlePaste}>
        {digits.map((d, i) => (
          <input
            key={i}
            ref={el => { refs.current[i] = el; }}
            type="text"
            inputMode="numeric"
            maxLength={1}
            value={d}
            aria-label={`OTP digit ${i + 1}`}
            onChange={e => handleChange(i, e.target.value)}
            onKeyDown={e => handleKeyDown(i, e)}
            disabled={loading}
            className={`w-11 h-12 text-center text-xl font-black border-2 rounded-xl outline-none transition-all
              ${d ? 'border-orange-600 bg-orange-50 dark:bg-orange-900/20' : 'border-zinc-200 dark:border-zinc-700 bg-zinc-50 dark:bg-zinc-800'}
              focus:border-orange-600 focus:ring-2 focus:ring-orange-600/20 disabled:opacity-50`}
          />
        ))}
      </div>

      {error && (
        <p className="text-xs text-red-500 bg-red-50 dark:bg-red-900/20 px-3 py-2 rounded-xl text-center">{error}</p>
      )}

      <div className="flex items-center justify-between text-xs">
        <span className={`font-bold ${countdown <= 30 ? 'text-red-500' : 'text-zinc-400'}`}>
          {t.ui.expires_in} {mm}:{ss}
        </span>
        {countdown === 0 && (
          <button
            onClick={() => {
              setCountdown(300);
              setDigits(['', '', '', '', '', '']);
              submittedRef.current = false;
              refs.current[0]?.focus();
              onResend();
            }}
            className="flex items-center gap-1 text-orange-600 font-bold hover:underline">
            <RefreshCw size={11} /> {t.ui.resend_otp}
          </button>
        )}
      </div>
    </div>
  );
}
