'use client';

import React, { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import {
  Sparkles,
  X,
  Send,
  RotateCcw,
  ExternalLink,
  ChevronDown,
  Tag,
  MapPin,
  Search,
  MessageSquare,
  ShieldCheck,
  AlertCircle,
  HelpCircle,
} from 'lucide-react';
import { useLocale } from '@/hooks/useLocale';
import { sendAIChat } from '@/lib/api/ai';
import type { ChatMessage, ListingCard } from '@/types/ai';

const INITIAL_SUGGESTIONS: Record<string, string[]> = {
  am: [
    'ላፕቶፕ ከ 50,000 ብር በታች',
    'የሚከራይ ቤት ቦሌ',
    'የደላላ ኮሚሽን ሞዴል እንዴት ይሰራል?',
    'መኪና ከመግዛት በፊት ምን ልመርምር?',
  ],
  ti: [
    'ላፕቶፕ ትሕቲ 50,000 ብር',
    'ክካረ ገዛ ኣብ ቦሌ',
    'ናይ ደላላ ኮሚሽን ከመይ ይሰርሕ?',
    'መኪና ቅድሚ ምግዛእ እንታይ ክምርመር ኣለዎ?',
  ],
  or: [
    'Kompiitara gatii 50,000 gadi',
    'Mana kireeffamu Boolee',
    'Dilaalli HUNTIQ akkamitti hojjeta?',
    'Konkolaataa bituu dura maal qoradhu?',
  ],
  en: [
    'Laptops under 50,000 ETB',
    'Apartments for rent in Bole',
    'How does the broker model work?',
    'What should I inspect before buying a car?',
  ],
};

const WELCOME_MESSAGES: Record<string, string> = {
  am: 'ሰላም! እኔ HUNTIQ AI ነኝ። በHUNTIQ ላይ እቃዎችን፣ ቤቶችን፣ መኪኖችን፣ አገልግሎቶችን ለማግኘት ወይም የገበያና ንግድ ምክር ለመጠየቅ ዝግጁ ነኝ። ምን ልርዳዎት?',
  ti: 'ከመይ! ኣነ HUNTIQ AI እየ። ኣብ HUNTIQ ንብረታት፣ ኣባይቲ፣ መካይን፣ ኣገልግሎታት ንምርካብ ወይ ናይ ዕዳጋን ንግድን ምኽሪ ንምሕታት ድሉው እየ። እንታይ ክሕግዘኩም?',
  or: 'Akkam! Ani HUNTIQ AI dha. Gabaa HUNTIQ irratti meeshaalee, manneen, konkolaattota, tajaajiloota barbaaduuf ykn gorsa daldalaa argachuuf qophiidha. Maal isin gargaaru?',
  en: 'Hello! I am HUNTIQ AI — your dedicated marketplace discovery and commercial advisory assistant. How can I help you find listings, evaluate pricing, or explore opportunities today?',
};

const ERROR_MESSAGES: Record<string, Record<string, string>> = {
  NO_INTERNET: {
    en: 'Your internet connection appears to be unavailable. Please reconnect and try again.',
    am: 'የኢንተርኔት ግንኙነትዎ የለም ወይም ተቋርጧል። እባክዎ ግንኙነትዎን መልሰው ከዚያ እንደገና ይሞክሩ።',
    or: 'Walitti hidhamiinsi interneetii keessanii hin jiru ykn citeera. Maaloo walitti hidhamiinsa keessan deebisaa, achiis irra deebi\'aa yaalaa.',
    ti: 'ምትእስሳር ኢንተርነትካ የለን ወይ ተቛሪጹ ኣሎ። በጃኻ ምትእስሳርካ መሊስካ ኣራኽብ፣ ድሕሪኡ እንደገና ፈትን።',
  },
  BACKEND_UNREACHABLE: {
    en: 'Unable to reach HUNTIQ servers. Please verify that the platform is accessible and retry.',
    am: 'ከ HUNTIQ ሰርቨሮች ጋር መገናኘት አልተቻለም። እባክዎ ፕላትፎርሙ እንደሚሰራ አረጋግጠው እንደገና ይሞክሩ።',
    or: 'Sarvaroota HUNTIQ waliin walqunnamuu hin dandeenye. Maaloo waltajjichi hojjetaa jiraachuu isaa mirkaneessaa irra deebi\'aa yaalaa.',
    ti: 'ምስ HUNTIQ ሰርቨራት ክራኸብ ኣይከኣለን። በጃኹም እቲ ፕላትፎርም ይሰርሕ ምህላዉ ኣረጋጊጽኩም ደጊምኩም ፈትኑ።',
  },
  AI_UNAVAILABLE: {
    en: 'HUNTIQ AI is temporarily unavailable. Your internet connection may be working, but our AI service cannot currently complete your request. Please try again shortly.',
    am: 'HUNTIQ AI በጊዜያዊነት አይሰራም። የኢንተርኔት ግንኙነትዎ ሊሰራ ይችላል ነገር ግን የ AI አገልግሎታችን በአሁኑ ጊዜ ጥያቄዎን ማጠናቀቅ አልቻለም። እባክዎ ጥቂት ቆይተው እንደገና ይሞክሩ።',
    or: 'HUNTIQ AI yeroof hojii ala dha. Walitti hidhamiinsi interneetii keessanii hojjechuu danda\'a, garuu tajaajilli AI keenya yeroo ammaa gaaffii keessan xumuruu hin dandeenye. Maaloo xiqqoo turee irra deebi\'aa yaalaa.',
    ti: 'HUNTIQ AI ንግዚኡ ኣየገልግልን ኣሎ። ምትእስሳር ኢንተርነትኩም ይሰርሕ ይኸውን፣ እንተኾነ ግና ናይ AI ኣገልግሎትና ኣብዚ እዋን\'ዚ ሕቶኹም ክውድእ ኣይከኣለን። በጃኹም ጸኒሕኩም ደጊምኩም ፈትኑ።',
  },
  TIMEOUT: {
    en: 'The AI request timed out. Please check your connection and retry.',
    am: 'የ AI ጥያቄው ጊዜው አልቋል። እባክዎ ግንኙነትዎን ፈትሸው እንደገና ይሞክሩ።',
    or: 'Gaaffiin AI yeroo fudhateera. Maaloo walitti hidhamiinsa keessan qoradhaa irra deebi\'aa yaalaa.',
    ti: 'ናይ AI ሕቶ ግዜ ወሲዱ ተቛሪጹ። በጃኹም ምትእስሳርኩም መርሚርኩም ደጊምኩም ፈትኑ።',
  },
  RATE_LIMITED: {
    en: 'Too many AI requests. Please wait a moment before asking again.',
    am: 'በጣም ብዙ የ AI ጥያቄዎች ቀርበዋል። እባክዎ ሌላ ጥያቄ ከመጠየቅዎ በፊት ትንሽ ይጠብቁ።',
    or: 'Gaaffileen AI baay\'ee dhiyaataniiru. Maaloo irra deebitanii gaafachuu dura xiqqoo eegaa.',
    ti: 'ብጣዕሚ ብዙሕ ናይ AI ሕቶታት ቀሪቡ። በጃኹም ካልእ ሕቶ ቅድሚ ምሕታትኩም ሒደት ተጸበዩ።',
  },
  GENERIC_ERROR: {
    en: 'An unexpected issue occurred while processing your request. Please try again.',
    am: 'ጥያቄዎን በማስኬድ ላይ ያልተጠበቀ ችግር አጋጥሟል። እባክዎ እንደገና ይሞክሩ።',
    or: 'Gaaffii keessan raawwachuu irratti rakkoon hin eegamne uumameera. Maaloo irra deebi\'aa yaalaa.',
    ti: 'ሕቶኹም ኣብ ምፍጻም ዘይተጸበዮ ጸገም ኣጋጢሙ። በጃኹም ደጊምኩም ፈትኑ።',
  },
};

const OFFLINE_BANNER: Record<string, string> = {
  en: 'You are currently offline. Changes and messages will be sent when reconnected.',
  am: 'በአሁኑ ጊዜ ከመስመር ውጭ ነዎት። ሲገናኙ መልዕክቶች ይላካሉ።',
  or: 'Yeroo ammaa toora irraa ala jirtu. Yeroo walqunnamtan ergaawwan ni ergamu.',
  ti: 'ኣብዚ እዋን ካብ መስመር ወጻኢ ኢኹም ዘለኹም። ምስ ተራኸብኩም መልእኽትታት ክስደድ እዩ።',
};

const RETRY_LABELS: Record<string, string> = {
  en: 'Retry',
  am: 'እንደገና ሞክር',
  or: 'Irra deebi\'i',
  ti: 'ደጊምካ ፈትን',
};

export function AIAssistant() {
  const { locale } = useLocale();
  const currentLang = (locale in INITIAL_SUGGESTIONS ? locale : 'en') as 'en' | 'am' | 'ti' | 'or';

  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [hasUnread, setHasUnread] = useState(false);
  const [isOnline, setIsOnline] = useState(true);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Monitor real network status
  useEffect(() => {
    if (typeof window === 'undefined') return;
    setIsOnline(navigator.onLine);

    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);
    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  // Initialize with localized welcome message
  useEffect(() => {
    if (messages.length === 0) {
      setMessages([
        {
          id: 'welcome',
          role: 'assistant',
          content: WELCOME_MESSAGES[currentLang] || WELCOME_MESSAGES.en,
          timestamp: Date.now(),
          suggestedQueries: INITIAL_SUGGESTIONS[currentLang] || INITIAL_SUGGESTIONS.en,
        },
      ]);
    }
  }, [currentLang, messages.length]);

  // Scroll to bottom on new messages
  useEffect(() => {
    if (isOpen) {
      messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages, isOpen, isLoading]);

  // Focus input when opened
  useEffect(() => {
    if (isOpen) {
      setHasUnread(false);
      setTimeout(() => inputRef.current?.focus(), 150);
    }
  }, [isOpen]);

  const handleClearChat = () => {
    setMessages([
      {
        id: 'welcome-' + Date.now(),
        role: 'assistant',
        content: WELCOME_MESSAGES[currentLang] || WELCOME_MESSAGES.en,
        timestamp: Date.now(),
        suggestedQueries: INITIAL_SUGGESTIONS[currentLang] || INITIAL_SUGGESTIONS.en,
      },
    ]);
  };

  const handleRetry = (failedMsgId: string, textToRetry: string) => {
    // Remove the failed assistant message to preserve clean context
    setMessages((prev) => prev.filter((m) => m.id !== failedMsgId));
    handleSend(textToRetry);
  };

  const handleSend = async (textToSend?: string) => {
    const text = (textToSend || input).trim();
    if (!text || isLoading) return;

    setInput('');

    // Append user message
    const userMsg: ChatMessage = {
      id: 'u-' + Date.now(),
      role: 'user',
      content: text,
      timestamp: Date.now(),
    };

    const newHistory = [...messages, userMsg];
    setMessages(newHistory);
    setIsLoading(true);

    try {
      // Check immediate network signal
      if (typeof navigator !== 'undefined' && !navigator.onLine) {
        throw new Error('NO_INTERNET');
      }

      // Build lightweight conversation history for Gemini (last 4 turns)
      const apiHistory = newHistory
        .filter((m) => m.id !== 'welcome' && !m.isError)
        .slice(-4)
        .map((m) => ({
          role: m.role === 'user' ? 'user' : 'model',
          parts: [{ text: m.content }],
        }));

      const res = await sendAIChat({
        message: text,
        language: currentLang,
        history: apiHistory,
      });

      const assistantMsg: ChatMessage = {
        id: 'a-' + Date.now(),
        role: 'assistant',
        content: res.reply || 'No response generated.',
        timestamp: Date.now(),
        path: res.path,
        intent: res.intent,
        isOutOfDomain: res.is_out_of_domain,
        listings: res.listings,
        suggestedQueries: res.suggested_queries,
        latencyMs: res.latency_ms,
      };

      setMessages((prev) => [...prev, assistantMsg]);

      if (!isOpen) {
        setHasUnread(true);
      }
    } catch (err: any) {
      let code: 'NO_INTERNET' | 'BACKEND_UNREACHABLE' | 'AI_UNAVAILABLE' | 'TIMEOUT' | 'RATE_LIMITED' | 'GENERIC_ERROR' = 'GENERIC_ERROR';

      const isNavOffline = typeof navigator !== 'undefined' && !navigator.onLine;
      const errMsg = (err?.message || err?.error || '').toLowerCase();
      const errStatus = err?.status || err?.statusCode || 0;

      if (isNavOffline || errMsg === 'no_internet') {
        code = 'NO_INTERNET';
      } else if (errStatus === 429 || errMsg.includes('too many') || errMsg.includes('rate limit')) {
        code = 'RATE_LIMITED';
      } else if (errStatus === 504 || errMsg.includes('timeout') || errMsg.includes('deadline exceeded')) {
        code = 'TIMEOUT';
      } else if (errStatus === 502 || errStatus === 503 || errMsg.includes('engine') || errMsg.includes('gemini') || errMsg.includes('ai service') || errMsg.includes('provider')) {
        code = 'AI_UNAVAILABLE';
      } else if (errStatus === 0 || errMsg.includes('failed to fetch') || errMsg.includes('networkerror') || errMsg.includes('load failed')) {
        code = 'BACKEND_UNREACHABLE';
      }

      const dict = ERROR_MESSAGES[code] || ERROR_MESSAGES.GENERIC_ERROR;
      const localizedContent = dict[currentLang] || dict.en;

      const errorMsg: ChatMessage = {
        id: 'err-' + Date.now(),
        role: 'assistant',
        content: localizedContent,
        timestamp: Date.now(),
        isError: true,
        errorCode: code,
        retryText: text,
      };
      setMessages((prev) => [...prev, errorMsg]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <>
      {/* ── FLOATING AI BUTTON (Bottom-Right) ───────────────────────────── */}
      {/* Positioned at bottom-[76px] on mobile (above 56px bottom nav) and bottom-6 on desktop */}
      <div className="fixed z-40 bottom-[76px] right-4 md:bottom-6 md:right-6 select-none print:hidden">
        <button
          type="button"
          onClick={() => setIsOpen((prev) => !prev)}
          aria-label={isOpen ? 'Close HUNTIQ AI' : 'Open HUNTIQ AI Assistant'}
          className={`group relative flex items-center justify-center transition-all duration-300 ease-out active:scale-90 ${
            isOpen
              ? 'w-12 h-12 md:w-14 md:h-14 bg-slate-800 dark:bg-slate-700 text-white shadow-lg ring-2 ring-slate-400/20'
              : 'w-13 h-13 md:w-14 md:h-14 bg-gradient-to-tr from-[#1769E0] via-[#207bf8] to-[#38A5FF] text-white shadow-xl shadow-[#1769E0]/35 hover:shadow-2xl hover:shadow-[#1769E0]/50 hover:scale-105 ring-4 ring-white/90 dark:ring-[#091426]/90'
          } rounded-full`}
        >
          {isOpen ? (
            <X size={24} className="stroke-[2.5px] transition-transform duration-200 group-hover:rotate-90" />
          ) : (
            <>
              <Sparkles size={24} className="stroke-[2.2px] animate-pulse" />
              {/* Pulsing online indicator dot */}
              <span className="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3.5 w-3.5 bg-emerald-500 border-2 border-white dark:border-[#091426]"></span>
              </span>

              {/* Unread message badge */}
              {hasUnread && (
                <span className="absolute -top-1 -left-1 flex h-4 w-4 bg-amber-500 rounded-full items-center justify-center text-[9px] font-black text-white ring-2 ring-white dark:ring-[#091426]">
                  !
                </span>
              )}
            </>
          )}

          {/* Desktop Hover Label Pill */}
          {!isOpen && (
            <span className="hidden md:flex absolute right-16 px-3 py-1.5 rounded-full bg-slate-900/90 dark:bg-white/95 text-white dark:text-slate-900 text-xs font-bold shadow-md whitespace-nowrap opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none items-center gap-1.5">
              <Sparkles size={13} className="text-[#38A5FF] dark:text-[#1769E0]" />
              HUNTIQ AI
            </span>
          )}
        </button>
      </div>

      {/* ── CHATBOT POPUP DIALOG / DRAWER ─────────────────────────────────── */}
      {/* Small screens: nicely contained card floating above bottom nav & below header */}
      {/* Large screens: docked card at bottom-24 right-6 (420px x 620px) */}
      {isOpen && (
        <div
          role="dialog"
          aria-labelledby="huntiq-ai-title"
          className="fixed z-50 transition-all duration-300 ease-out flex flex-col overflow-hidden bg-white/98 dark:bg-[#0c182c]/98 backdrop-blur-2xl border border-slate-200/90 dark:border-slate-800/90 shadow-2xl
            inset-x-2 bottom-[76px] top-16 rounded-3xl
            sm:inset-auto sm:bottom-24 sm:right-6 sm:w-[420px] sm:max-w-[calc(100vw-3rem)] sm:h-[620px] sm:max-h-[calc(100vh-8rem)] sm:rounded-3xl animate-in fade-in zoom-in-95"
        >
          {/* ── Top Header ── */}
          <div className="flex items-center justify-between px-4 py-3 bg-gradient-to-r from-slate-50 via-white to-slate-50 dark:from-[#091426] dark:via-[#0c182c] dark:to-[#091426] border-b border-slate-200/80 dark:border-slate-800/80 select-none">
            <div className="flex items-center gap-2.5">
              <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-[#1769E0] to-[#38A5FF] flex items-center justify-center text-white shadow-md shadow-[#1769E0]/20">
                <Sparkles size={16} className="stroke-[2.5px]" />
              </div>
              <div>
                <div className="flex items-center gap-1.5">
                  <h3 id="huntiq-ai-title" className="text-xs sm:text-sm font-black tracking-tight text-slate-900 dark:text-white">
                    HUNTIQ AI
                  </h3>
                  <span className={`px-1.5 py-0.5 rounded text-[9px] font-black uppercase tracking-wider border ${
                    isOnline
                      ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20'
                      : 'bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20'
                  }`}>
                    {isOnline ? 'Live' : 'Offline'}
                  </span>
                </div>
                <p className="text-[10px] text-slate-500 dark:text-slate-400 font-medium">
                  Ethiopia Broker & Marketplace Advisor
                </p>
              </div>
            </div>

            <div className="flex items-center gap-1 text-slate-400 dark:text-slate-500">
              {/* Clear conversation */}
              <button
                type="button"
                onClick={handleClearChat}
                title="Clear conversation"
                className="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors"
                aria-label="Restart chat"
              >
                <RotateCcw size={15} />
              </button>

              {/* Close/Minimize */}
              <button
                type="button"
                onClick={() => setIsOpen(false)}
                title="Minimize chat"
                className="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors"
                aria-label="Close chat"
              >
                <ChevronDown size={18} />
              </button>
            </div>
          </div>

          {/* ── Offline Banner ── */}
          {!isOnline && (
            <div className="bg-amber-500/10 border-b border-amber-500/20 px-3.5 py-1.5 flex items-center gap-2 text-[11px] font-semibold text-amber-600 dark:text-amber-400 select-none">
              <AlertCircle size={13} className="flex-shrink-0" />
              <span>{OFFLINE_BANNER[currentLang] || OFFLINE_BANNER.en}</span>
            </div>
          )}

          {/* ── Messages Scroll Area ── */}
          <div className="flex-1 overflow-y-auto p-3 sm:p-4 space-y-3.5 overscroll-contain text-slate-800 dark:text-slate-200">
            {messages.map((m) => {
              const isAssistant = m.role === 'assistant';

              return (
                <div
                  key={m.id}
                  className={`flex flex-col ${isAssistant ? 'items-start' : 'items-end'}`}
                >
                  <div className="flex items-end gap-1.5 max-w-[88%] sm:max-w-[85%]">
                    {/* Assistant Avatar */}
                    {isAssistant && (
                      <div className="w-6 h-6 rounded-full bg-[#1769E0]/10 dark:bg-[#38A5FF]/10 text-[#1769E0] dark:text-[#38A5FF] flex-shrink-0 flex items-center justify-center mb-1">
                        <Sparkles size={13} />
                      </div>
                    )}

                    {/* Chat Bubble */}
                    <div
                      className={`rounded-2xl px-3.5 py-2.5 text-xs sm:text-sm leading-relaxed ${
                        isAssistant
                          ? m.isError
                            ? 'bg-rose-500/10 border border-rose-500/25 text-slate-800 dark:text-slate-100 rounded-bl-sm'
                            : m.isOutOfDomain
                            ? 'bg-amber-500/10 border border-amber-500/20 text-slate-800 dark:text-slate-100 rounded-bl-sm'
                            : 'bg-slate-100 dark:bg-slate-800/90 text-slate-900 dark:text-slate-100 rounded-bl-sm border border-slate-200/50 dark:border-slate-700/50'
                          : 'bg-gradient-to-r from-[#1769E0] to-[#207bf8] text-white rounded-br-sm shadow-sm'
                      }`}
                    >
                      {/* Out of domain badge */}
                      {isAssistant && m.isOutOfDomain && (
                        <div className="flex items-center gap-1 text-[11px] font-bold text-amber-600 dark:text-amber-400 mb-1">
                          <AlertCircle size={13} />
                          <span>HUNTIQ Marketplace Only</span>
                        </div>
                      )}

                      <p className="whitespace-pre-wrap font-normal">{m.content}</p>

                      {/* Error Card with Retry Button */}
                      {isAssistant && m.isError && m.retryText && (
                        <div className="mt-2.5 pt-2 border-t border-rose-500/20 flex items-center justify-between gap-2">
                          <span className="text-[10px] text-rose-500 dark:text-rose-400 font-semibold flex items-center gap-1">
                            <AlertCircle size={12} />
                            {m.errorCode}
                          </span>
                          <button
                            type="button"
                            onClick={() => handleRetry(m.id, m.retryText!)}
                            className="flex items-center gap-1 px-2.5 py-1 rounded-lg bg-rose-500 hover:bg-rose-600 text-white text-[11px] font-bold shadow-sm transition-all active:scale-95"
                          >
                            <RotateCcw size={11} className="stroke-[2.5px]" />
                            <span>{RETRY_LABELS[currentLang] || RETRY_LABELS.en}</span>
                          </button>
                        </div>
                      )}

                      {/* ── PRODUCT LISTINGS (Non-intrusive clean cards that open on normal screen) ── */}
                      {isAssistant && m.listings && m.listings.length > 0 && (
                        <div className="mt-3 pt-2.5 border-t border-slate-200/80 dark:border-slate-700/80 space-y-2">
                          <div className="flex items-center justify-between text-[11px] font-bold text-slate-600 dark:text-slate-400">
                            <span>Verified Marketplace Listings ({m.listings.length}):</span>
                          </div>

                          <div className="space-y-2 max-h-60 overflow-y-auto pr-1 scrollbar-thin">
                            {m.listings.map((prod) => (
                              <div
                                key={prod.id}
                                className="group/item flex items-center gap-2.5 p-2 rounded-xl bg-white dark:bg-[#091426] border border-slate-200 dark:border-slate-700/70 hover:border-[#1769E0] dark:hover:border-[#38A5FF] transition-all shadow-sm"
                              >
                                {/* Thumbnail */}
                                <div className="w-12 h-12 rounded-lg bg-slate-100 dark:bg-slate-800 overflow-hidden flex-shrink-0 flex items-center justify-center border border-slate-200/60 dark:border-slate-700/60">
                                  {prod.image_url ? (
                                    <img
                                      src={prod.image_url}
                                      alt={prod.title}
                                      className="w-full h-full object-cover group-hover/item:scale-105 transition-transform"
                                      loading="lazy"
                                      onError={(e) => {
                                        (e.target as HTMLElement).style.display = 'none';
                                      }}
                                    />
                                  ) : (
                                    <Tag size={16} className="text-slate-400" />
                                  )}
                                </div>

                                {/* Details */}
                                <div className="flex-1 min-w-0">
                                  <h4 className="text-xs font-bold text-slate-900 dark:text-white truncate">
                                    {prod.title}
                                  </h4>
                                  <div className="flex items-center gap-1.5 mt-0.5">
                                    <span className="text-xs font-black text-[#1769E0] dark:text-[#38A5FF]">
                                      {prod.price.toLocaleString()} {prod.currency || 'ETB'}
                                    </span>
                                    {prod.listing_type === 'rent' && (
                                      <span className="text-[10px] text-slate-400">/rent</span>
                                    )}
                                  </div>
                                  <div className="flex items-center gap-2 text-[10px] text-slate-500 dark:text-slate-400 truncate mt-0.5">
                                    {prod.location_name && (
                                      <span className="flex items-center gap-0.5 truncate">
                                        <MapPin size={10} />
                                        {prod.location_name}
                                      </span>
                                    )}
                                  </div>
                                </div>

                                {/* Action: Open on normal screen */}
                                <Link
                                  href={`/products/${prod.id}`}
                                  className="flex items-center gap-1 px-2.5 py-1.5 rounded-lg bg-[#1769E0]/10 dark:bg-[#38A5FF]/10 text-[#1769E0] dark:text-[#38A5FF] hover:bg-[#1769E0] hover:text-white dark:hover:bg-[#38A5FF] dark:hover:text-slate-900 text-[11px] font-bold transition-all flex-shrink-0"
                                >
                                  <span>View</span>
                                  <ExternalLink size={12} />
                                </Link>
                              </div>
                            ))}
                          </div>

                          {/* Primary CTA: View all results on main screen grid */}
                          <Link
                            href={`/products`}
                            className="w-full mt-2 py-2 px-3 rounded-xl bg-slate-200/70 hover:bg-slate-200 dark:bg-slate-700/60 dark:hover:bg-slate-700 text-slate-800 dark:text-slate-200 text-xs font-bold flex items-center justify-center gap-1.5 transition-colors"
                          >
                            <Search size={13} className="text-[#1769E0] dark:text-[#38A5FF]" />
                            <span>Browse full catalog on screen</span>
                          </Link>
                        </div>
                      )}
                    </div>
                  </div>

                  {/* Latency and path info */}
                  {isAssistant && m.latencyMs !== undefined && (
                    <span className="text-[9px] text-slate-400 dark:text-slate-500 mt-1 ml-8">
                      {m.latencyMs}ms {m.path ? `• ${m.path}` : ''}
                    </span>
                  )}

                  {/* ── Suggested query chips ── */}
                  {isAssistant && m.suggestedQueries && m.suggestedQueries.length > 0 && (
                    <div className="flex flex-wrap gap-1.5 mt-2 ml-7 max-w-[95%]">
                      {m.suggestedQueries.map((chip, idx) => (
                        <button
                          key={idx}
                          type="button"
                          onClick={() => handleSend(chip)}
                          className="text-[11px] sm:text-xs font-medium px-2.5 py-1 rounded-full bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-[#1769E0] dark:text-[#38A5FF] border border-slate-200/80 dark:border-slate-700/80 transition-all hover:scale-[1.02] active:scale-95 text-left"
                        >
                          {chip}
                        </button>
                      ))}
                    </div>
                  )}
                </div>
              );
            })}

            {/* Typing animation indicator */}
            {isLoading && (
              <div className="flex items-end gap-1.5">
                <div className="w-6 h-6 rounded-full bg-[#1769E0]/10 dark:bg-[#38A5FF]/10 text-[#1769E0] dark:text-[#38A5FF] flex-shrink-0 flex items-center justify-center mb-1">
                  <Sparkles size={13} />
                </div>
                <div className="bg-slate-100 dark:bg-slate-800/90 rounded-2xl rounded-bl-sm px-4 py-2.5 border border-slate-200/50 dark:border-slate-700/50 flex items-center gap-1.5">
                  <span className="w-2 h-2 rounded-full bg-[#1769E0] animate-bounce" style={{ animationDelay: '0ms' }} />
                  <span className="w-2 h-2 rounded-full bg-[#1769E0] animate-bounce" style={{ animationDelay: '150ms' }} />
                  <span className="w-2 h-2 rounded-full bg-[#1769E0] animate-bounce" style={{ animationDelay: '300ms' }} />
                  <span className="text-[11px] text-slate-500 dark:text-slate-400 ml-1 font-medium">
                    HUNTIQ AI is thinking...
                  </span>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* ── Input Bar ── */}
          <div className="p-2.5 sm:p-3 bg-white dark:bg-[#0c182c] border-t border-slate-200/80 dark:border-slate-800/80">
            <form
              onSubmit={(e) => {
                e.preventDefault();
                handleSend();
              }}
              className="flex items-center gap-1.5 bg-slate-100 dark:bg-slate-800/80 rounded-2xl px-3 py-1.5 border border-slate-200/60 dark:border-slate-700/60 focus-within:border-[#1769E0] dark:focus-within:border-[#38A5FF] focus-within:ring-2 focus-within:ring-[#1769E0]/15 transition-all"
            >
              <input
                ref={inputRef}
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder={
                  currentLang === 'am'
                    ? 'በHUNTIQ ላይ እቃ፣ ቤት ወይም የገበያ መረጃ ይጠይቁ...'
                    : currentLang === 'ti'
                    ? 'ኣብ HUNTIQ ንብረት፣ ገዛ ወይ ሓበሬታ ሕተቱ...'
                    : currentLang === 'or'
                    ? 'HUNTIQ irratti meeshaa ykn odeeffannoo gaafadhaa...'
                    : 'Search listings, ask broker advice, prices...'
                }
                disabled={isLoading}
                className="flex-1 bg-transparent text-xs sm:text-sm text-slate-900 dark:text-slate-100 placeholder-slate-400 dark:placeholder-slate-500 focus:outline-none min-w-0"
              />

              <button
                type="submit"
                disabled={!input.trim() || isLoading}
                aria-label="Send inquiry"
                className="w-8 h-8 rounded-xl bg-[#1769E0] hover:bg-[#1359c4] disabled:opacity-40 disabled:hover:bg-[#1769E0] text-white flex items-center justify-center transition-all flex-shrink-0 active:scale-95 shadow-sm"
              >
                <Send size={14} className="stroke-[2.5px]" />
              </button>
            </form>

            <div className="flex items-center justify-between px-1 mt-1.5 text-[10px] text-slate-400 dark:text-slate-500">
              <span className="flex items-center gap-1">
                <ShieldCheck size={11} className="text-emerald-500" />
                Verified Broker AI
              </span>
              <span>HUNTIQ Marketplace v1.0</span>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
