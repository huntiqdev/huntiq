'use client';
import { useLocale } from '@/hooks/useLocale';
import { FileText } from 'lucide-react';

export default function TermsPage() {
  const { t } = useLocale();
  const tp = (t as any).terms_page;

  const sections = [
    { title: tp.s1_title, body: tp.s1_content },
    { title: tp.s2_title, body: tp.s2_content },
    { title: tp.s3_title, body: tp.s3_content },
    { title: tp.s4_title, body: tp.s4_content },
    { title: tp.s5_title, body: tp.s5_content },
    { title: tp.s6_title, body: tp.s6_content },
    { title: tp.s7_title, body: tp.s7_content },
    { title: tp.s8_title, body: tp.s8_content },
    { title: tp.s9_title, body: tp.s9_content },
    { title: tp.s10_title, body: tp.s10_content },
  ];

  return (
    <div className="max-w-3xl mx-auto px-4 py-8 sm:py-12 pb-safe">
      <div className="flex items-center gap-3 mb-6 sm:mb-8">
        <div className="w-10 h-10 rounded-xl bg-orange-600/10 text-orange-600 flex items-center justify-center">
          <FileText size={20} />
        </div>
        <div>
          <h1 className="text-xl sm:text-2xl font-black text-zinc-900 dark:text-white">{tp.title}</h1>
          <p className="text-xs opacity-40 font-bold">{tp.last_updated}</p>
        </div>
      </div>

      <p className="text-sm opacity-60 mb-6 sm:mb-8 leading-relaxed">{tp.content}</p>

      <div className="space-y-4 sm:space-y-6">
        {sections.map((s, i) => (
          <div key={i} className="bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 p-4 sm:p-5 shadow-sm">
            <h2 className="font-black mb-2 text-sm">{s.title}</h2>
            <p className="text-sm opacity-70 leading-relaxed">{s.body}</p>
          </div>
        ))}
      </div>

      <p className="mt-8 text-xs opacity-40 text-center">
        {tp.s10_content}
      </p>
    </div>
  );
}
