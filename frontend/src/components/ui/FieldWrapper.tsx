import { ReactNode } from 'react';

interface FieldWrapperProps {
  label?: string;
  error?: string;
  children: ReactNode;
}

export function FieldWrapper({ label, error, children }: FieldWrapperProps) {
  return (
    <div className="flex flex-col gap-1">
      {label && <label className="text-xs font-black uppercase tracking-wider opacity-60">{label}</label>}
      {children}
      {error && <p className="text-xs text-red-500">{error}</p>}
    </div>
  );
}
