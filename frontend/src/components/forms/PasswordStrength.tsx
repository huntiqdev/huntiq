'use client';

interface PasswordStrengthProps {
  password: string;
  className?: string;
}

interface Rule {
  label: string;
  test: (pw: string) => boolean;
}

const RULES: Rule[] = [
  { label: 'At least 8 characters',    test: pw => pw.length >= 8 },
  { label: 'One uppercase letter (A-Z)', test: pw => /[A-Z]/.test(pw) },
  { label: 'One lowercase letter (a-z)', test: pw => /[a-z]/.test(pw) },
  { label: 'One number (0-9)',          test: pw => /\d/.test(pw) },
];

export function validatePasswordStrength(password: string): { valid: boolean; message: string } {
  if (password.length < 8)       return { valid: false, message: 'Password must be at least 8 characters' };
  if (!/[A-Z]/.test(password))   return { valid: false, message: 'Password must contain at least one uppercase letter' };
  if (!/[a-z]/.test(password))   return { valid: false, message: 'Password must contain at least one lowercase letter' };
  if (!/\d/.test(password))      return { valid: false, message: 'Password must contain at least one number' };
  return { valid: true, message: '' };
}

export default function PasswordStrength({ password, className = '' }: PasswordStrengthProps) {
  if (!password) return null;

  const passed = RULES.filter(r => r.test(password)).length;
  const total  = RULES.length;
  const pct    = (passed / total) * 100;

  const strengthColor =
    pct <= 25  ? 'bg-red-500' :
    pct <= 50  ? 'bg-orange-500' :
    pct <= 75  ? 'bg-yellow-500' :
                 'bg-green-500';

  const strengthLabel =
    pct <= 25  ? 'Weak' :
    pct <= 50  ? 'Fair' :
    pct <= 75  ? 'Good' :
                 'Strong';

  return (
    <div className={`space-y-2 ${className}`}>
      {/* Bar */}
      <div className="flex items-center gap-2">
        <div className="flex-1 h-1.5 bg-zinc-200 dark:bg-white/10 rounded-full overflow-hidden">
          <div
            className={`h-full rounded-full transition-all duration-300 ${strengthColor}`}
            style={{ width: `${pct}%` }}
          />
        </div>
        <span className={`text-[10px] font-black uppercase tracking-wider ${
          pct <= 25  ? 'text-red-500' :
          pct <= 50  ? 'text-orange-500' :
          pct <= 75  ? 'text-yellow-500' :
                       'text-green-500'
        }`}>{strengthLabel}</span>
      </div>

      {/* Rules checklist */}
      <div className="grid grid-cols-2 gap-x-3 gap-y-1">
        {RULES.map(rule => {
          const ok = rule.test(password);
          return (
            <div key={rule.label} className="flex items-center gap-1.5">
              <div className={`w-3 h-3 rounded-full flex items-center justify-center shrink-0 ${
                ok ? 'bg-green-500' : 'bg-zinc-300 dark:bg-white/10'
              }`}>
                {ok && (
                  <svg width="8" height="8" viewBox="0 0 8 8" fill="none">
                    <path d="M1.5 4L3 5.5L6.5 2" stroke="white" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
                  </svg>
                )}
              </div>
              <span className={`text-[10px] leading-tight ${
                ok ? 'text-green-600 dark:text-green-400' : 'text-zinc-400 dark:text-white/30'
              }`}>{rule.label}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
