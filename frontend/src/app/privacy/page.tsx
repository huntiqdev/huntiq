'use client';
import { useLocale } from '@/hooks/useLocale';
import { Shield } from 'lucide-react';

export default function PrivacyPage() {
  const { t } = useLocale();
  const pp = (t as any).privacy_page;

  const sections = [
    { title: pp.s1_title, body: pp.s1_content },
    { title: pp.s2_title, body: pp.s2_content },
    { title: pp.s3_title, body: pp.s3_content },
    { title: pp.s4_title, body: pp.s4_content },
    { title: pp.s5_title, body: pp.s5_content },
    { title: pp.s6_title, body: pp.s6_content },
  ];

  return (
    <div className="max-w-3xl mx-auto px-4 py-12">
      <div className="flex items-center gap-3 mb-8">
        <div className="w-10 h-10 rounded-xl bg-orange-600/10 text-orange-600 flex items-center justify-center">
          <Shield size={20} />
        </div>
        <div>
          <h1 className="text-2xl font-black">{pp.title}</h1>
          <p className="text-xs opacity-40 font-bold">{pp.last_updated}</p>
        </div>
      </div>

      <p className="text-sm opacity-60 mb-8 leading-relaxed">{pp.content}</p>

      <div className="space-y-6">
        {sections.map((s, i) => (
          <div key={i} className="bg-white/70 dark:bg-zinc-900 rounded-2xl border border-violet-100 dark:border-zinc-800 p-5 backdrop-blur-sm">
            <h2 className="font-black mb-2 text-sm">{s.title}</h2>
            <p className="text-sm opacity-70 leading-relaxed">{s.body}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
