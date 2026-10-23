'use client';
import { useState, FormEvent } from 'react';
import { useLocale } from '@/hooks/useLocale';
import { api } from '@/lib/api';
import { Mail, Phone, MapPin, Clock, Send, CheckCircle, AlertCircle, Loader2 } from 'lucide-react';

export default function ContactPage() {
  const { t } = useLocale();
  const [sent, setSent] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [form, setForm] = useState({ name: '', email: '', message: '' });

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await api.submitContactMessage(form);
      setSent(true);
      setForm({ name: '', email: '', message: '' });
    } catch (err: any) {
      setError(err?.message || 'Failed to send message. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-5xl mx-auto px-4 py-12">
      <div className="text-center mb-12">
        <h1 className="text-3xl font-black mb-3">{(t as any).footer.contact}</h1>
        <p className="opacity-60 text-sm max-w-md mx-auto">{t.ui.message_placeholder}</p>
      </div>

      <div className="grid md:grid-cols-2 gap-8">
        {/* Info */}
        <div className="space-y-6">
          {[
            { icon: Phone, label: t.ui.company_phone,   value: (t as any).contact_page.phone },
            { icon: Mail,  label: t.ui.email,           value: (t as any).contact_page.email },
            { icon: MapPin,label: t.ui.address,         value: (t as any).contact_page.address },
            { icon: Clock, label: t.ui.working_hours,    value: (t as any).contact_page.working_hours },
          ].map(({ icon: Icon, label, value }) => (
            <div key={label} className="flex items-start gap-4 p-4 rounded-2xl bg-white/70 dark:bg-zinc-900 border border-violet-100 dark:border-zinc-800 backdrop-blur-sm">
              <div className="w-10 h-10 rounded-xl flex items-center justify-center shrink-0 text-white" style={{ background: 'linear-gradient(135deg,#7c3aed,#3b82f6)' }}>
                <Icon size={18} />
              </div>
              <div>
                <p className="text-xs font-black uppercase tracking-widest opacity-40 mb-0.5">{label}</p>
                <p className="font-bold text-sm">{value}</p>
              </div>
            </div>
          ))}
        </div>

        {/* Form */}
        <div className="bg-white/70 dark:bg-zinc-900 rounded-2xl border border-violet-100 dark:border-zinc-800 p-6 backdrop-blur-sm">
          {sent ? (
            <div className="flex flex-col items-center justify-center h-full gap-4 py-10 text-center">
              <CheckCircle size={48} className="text-green-500" />
              <h2 className="text-xl font-black">{t.ui.message_sent}</h2>
              <p className="opacity-60 text-sm">{t.ui.message_reply}</p>
              <button onClick={() => setSent(false)}
                className="mt-2 text-orange-600 font-bold text-sm hover:underline">
                {t.ui.send_another}
              </button>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              {error && (
                <div className="flex items-center gap-2 p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-600 dark:text-red-400 text-xs font-semibold">
                  <AlertCircle size={15} className="shrink-0" />
                  <span>{error}</span>
                </div>
              )}
              <div>
                <label className="text-xs font-black uppercase tracking-widest opacity-40 block mb-1">{t.ui.your_name}</label>
                <input required value={form.name} onChange={e => setForm(f => ({ ...f, name: e.target.value }))}
                  disabled={loading}
                  className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 disabled:opacity-50"
                  placeholder={t.ui.your_name} />
              </div>
              <div>
                <label className="text-xs font-black uppercase tracking-widest opacity-40 block mb-1">{t.ui.email_address}</label>
                <input required type="email" value={form.email} onChange={e => setForm(f => ({ ...f, email: e.target.value }))}
                  disabled={loading}
                  className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 disabled:opacity-50"
                  placeholder="you@example.com" />
              </div>
              <div>
                <label className="text-xs font-black uppercase tracking-widest opacity-40 block mb-1">{t.ui.message}</label>
                <textarea required rows={5} value={form.message} onChange={e => setForm(f => ({ ...f, message: e.target.value }))}
                  disabled={loading}
                  className="w-full border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2.5 text-sm outline-none focus:ring-2 focus:ring-orange-600 bg-zinc-50 dark:bg-zinc-800 resize-none disabled:opacity-50"
                  placeholder={t.ui.message_placeholder} />
              </div>
              <button type="submit"
                disabled={loading}
                className="w-full bg-orange-600 hover:bg-orange-700 disabled:opacity-60 text-white rounded-xl py-3 font-bold text-sm flex items-center justify-center gap-2 transition-colors">
                {loading ? (
                  <>
                    <Loader2 size={16} className="animate-spin" />
                    <span>{t.ui.please_wait || 'Please wait...'}</span>
                  </>
                ) : (
                  <>
                    <Send size={16} />
                    <span>{t.ui.send_message}</span>
                  </>
                )}
              </button>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}
