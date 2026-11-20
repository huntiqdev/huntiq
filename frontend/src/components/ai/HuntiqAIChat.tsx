'use client';

import React, { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import {
  Sparkles,
  Send,
  RotateCcw,
  ArrowLeft,
  Search,
  ExternalLink,
  ShieldCheck,
  AlertCircle,
  HelpCircle,
  ChevronRight,
  Globe2,
} from 'lucide-react';
import { api, ChatResponse, ListingCardDTO } from '@/lib/api';

interface HuntiqAIChatProps {
  onCloseAI?: () => void;
  initialQuery?: string;
}

interface MessageItem {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  listings?: ListingCardDTO[];
  suggestedQueries?: string[];
  isOutOfDomain?: boolean;
  fallbackUsed?: boolean;
  timestamp: Date;
}

const SUPPORTED_LANGS = [
  { code: 'en', label: 'English', native: 'English' },
  { code: 'am', label: 'Amharic', native: 'አማርኛ' },
  { code: 'ti', label: 'Tigrinya', native: 'ትግርኛ' },
  { code: 'or', label: 'Afaan Oromo', native: 'Oromoo' },
];

const SUGGESTED_PROMPTS: Record<string, string[]> = {
  en: [
    'Find 2-bedroom furnished apartments for rent in Bole',
    'Show laptops under 45,000 ETB in Addis Ababa',
    'Find licensed electricians in Addis Ababa',
    'How do I post my property for rent on HUNTIQ?',
  ],
  am: [
    'በቦሌ የሚከራይ ባለ ሁለት መኝታ ቤት አሳየኝ',
    'ከ50,000 ብር በታች ያሉ ላፕቶፖች አሳየኝ',
    'በአዲስ አበባ ውስጥ ያሉ የኤሌክትሪክ ባለሙያዎችን ፈልግ',
    'በHUNTIQ ላይ ንብረቴን ለኪራይ እንዴት ማውጣት እችላለሁ?',
  ],
  ti: [
    'ኣብ መቐለ ወይ ኣዲስ ኣበባ ዝካረ ገዛ ኣርእየኒ',
    'ሕሳር ላፕቶፓት ክደሊ ደልየ',
    'ኣብ HUNTIQ ብኸመይ ንብረተይ ክሸይጥ ይኽእል?',
  ],
  or: [
    'Mana kiraa Finfinnee keessaa barbaada',
    'Konkolaattota gurguramuuf qophaa\'an naaf agarsiisi',
    'HUNTIQ irratti akkamitti gurguruu danda\'a?',
  ],
};

export function HuntiqAIChat({ onCloseAI, initialQuery }: HuntiqAIChatProps) {
  const [messages, setMessages] = useState<MessageItem[]>([]);
  const [inputText, setInputText] = useState(initialQuery || '');
  const [selectedLang, setSelectedLang] = useState('en');
  const [loading, setLoading] = useState(false);
  const [sessionId, setSessionId] = useState<string>('');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Initialize unique session identifier
  useEffect(() => {
    let sid = sessionStorage.getItem('huntiq_ai_session_id');
    if (!sid) {
      sid = 'huntiq_sess_' + Math.random().toString(36).substring(2, 15) + '_' + Date.now();
      sessionStorage.setItem('huntiq_ai_session_id', sid);
    }
    setSessionId(sid);

    // Initial greeting
    setMessages([
      {
        id: 'welcome-msg',
        role: 'assistant',
        content:
          'Welcome to HUNTIQ AI, Made by HUNTIQ. I am your verified marketplace assistant. Ask me to discover products, apartments, vehicles, professional services, gigs, or platform guidance across Ethiopia.',
        suggestedQueries: SUGGESTED_PROMPTS[selectedLang] || SUGGESTED_PROMPTS.en,
        timestamp: new Date(),
      },
    ]);
  }, []);

  // Auto-scroll to bottom of conversation
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, loading]);

  const handleSend = async (messageText?: string) => {
    const textToSend = (messageText || inputText).trim();
    if (!textToSend || loading) return;

    setInputText('');

    // Append user message
    const userMsg: MessageItem = {
      id: 'msg-' + Date.now(),
      role: 'user',
      content: textToSend,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMsg]);
    setLoading(true);

    try {
      const res: ChatResponse = await api.chatWithAI({
        session_id: sessionId,
        message: textToSend,
        language: selectedLang,
      });

      const assistantMsg: MessageItem = {
        id: 'reply-' + Date.now(),
        role: 'assistant',
        content: res.reply,
        listings: res.listings,
        suggestedQueries: res.suggested_queries,
        isOutOfDomain: res.is_out_of_domain,
        fallbackUsed: res.fallback_used,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, assistantMsg]);
    } catch (err: any) {
      setMessages((prev) => [
        ...prev,
        {
          id: 'err-' + Date.now(),
          role: 'assistant',
          content:
            'HUNTIQ AI is temporarily unable to reach the intelligence provider. You can still discover all verified listings directly using the manual search and filters.',
          fallbackUsed: true,
          timestamp: new Date(),
        },
      ]);
    } finally {
      setLoading(false);
      setTimeout(() => inputRef.current?.focus(), 100);
    }
  };

  const handleClearSession = async () => {
    if (confirm('Clear current HUNTIQ AI conversation?')) {
      try {
        if (sessionId) {
          await api.clearAISession(sessionId);
        }
      } catch {
        // Non-blocking
      }
      const newSid = 'huntiq_sess_' + Math.random().toString(36).substring(2, 15) + '_' + Date.now();
      sessionStorage.setItem('huntiq_ai_session_id', newSid);
      setSessionId(newSid);
      setMessages([
        {
          id: 'welcome-reset',
          role: 'assistant',
          content: 'Conversation cleared. How can HUNTIQ AI assist your marketplace search today?',
          suggestedQueries: SUGGESTED_PROMPTS[selectedLang] || SUGGESTED_PROMPTS.en,
          timestamp: new Date(),
        },
      ]);
    }
  };

  return (
    <div className="w-full max-w-5xl mx-auto rounded-3xl border border-zinc-200/80 dark:border-zinc-800 bg-white dark:bg-zinc-950 shadow-2xl overflow-hidden flex flex-col transition-all">
      {/* ── Top Header / Brand Bar ── */}
      <div className="px-4 py-3 sm:px-6 sm:py-4 bg-gradient-to-r from-blue-700 via-blue-600 to-indigo-700 text-white flex items-center justify-between gap-3 shadow-md">
        <div className="flex items-center gap-3">
          {onCloseAI && (
            <button
              onClick={onCloseAI}
              className="p-1.5 rounded-xl bg-white/10 hover:bg-white/20 text-white transition-all flex items-center gap-1 text-xs font-semibold"
              title="Return to manual browsing"
            >
              <ArrowLeft size={16} />
              <span className="hidden sm:inline">Manual Mode</span>
            </button>
          )}

          <div className="flex items-center gap-2">
            <div className="w-9 h-9 rounded-2xl bg-white/20 backdrop-blur-md flex items-center justify-center shadow-inner">
              <Sparkles size={20} className="text-yellow-300 animate-pulse" />
            </div>
            <div>
              <div className="flex items-center gap-1.5">
                <span className="font-black text-lg tracking-tight">HUNTIQ AI</span>
                <span className="text-[10px] font-bold px-2 py-0.5 rounded-full bg-white/20 text-yellow-200 border border-white/20 uppercase tracking-widest">
                  Made by HUNTIQ
                </span>
              </div>
              <p className="text-[11px] text-blue-100 hidden sm:block">
                Domain-Restricted Universal Marketplace Intelligence
              </p>
            </div>
          </div>
        </div>

        {/* ── Controls: Language Selector & Clear ── */}
        <div className="flex items-center gap-2">
          {/* Language Selector */}
          <div className="flex items-center gap-1 bg-white/15 backdrop-blur-md p-1 rounded-2xl border border-white/20 text-xs">
            <Globe2 size={13} className="text-blue-100 ml-1.5 hidden sm:inline" />
            {SUPPORTED_LANGS.map((lang) => (
              <button
                key={lang.code}
                onClick={() => setSelectedLang(lang.code)}
                className={`px-2 py-1 rounded-xl font-bold transition-all ${
                  selectedLang === lang.code
                    ? 'bg-white text-blue-800 shadow-sm'
                    : 'text-white/80 hover:text-white hover:bg-white/10'
                }`}
                title={lang.label}
              >
                {lang.native}
              </button>
            ))}
          </div>

          {/* Reset / Clear Button */}
          <button
            onClick={handleClearSession}
            className="p-2 rounded-xl bg-white/10 hover:bg-white/20 text-white transition-all"
            title="Clear Conversation"
          >
            <RotateCcw size={15} />
          </button>
        </div>
      </div>

      {/* ── Chat Messages Scroll Area ── */}
      <div className="p-4 sm:p-6 space-y-4 overflow-y-auto max-h-[600px] min-h-[380px] bg-zinc-50/60 dark:bg-zinc-900/40">
        {messages.map((msg) => {
          const isUser = msg.role === 'user';
          return (
            <div
              key={msg.id}
              className={`flex flex-col ${isUser ? 'items-end' : 'items-start'} space-y-2`}
            >
              <div className="flex items-start gap-2 max-w-[85%] sm:max-w-[80%]">
                {!isUser && (
                  <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-600 text-white flex items-center justify-center shrink-0 shadow-md mt-1">
                    <Sparkles size={16} />
                  </div>
                )}

                <div
                  className={`p-3.5 sm:p-4 rounded-2xl text-sm leading-relaxed shadow-sm ${
                    isUser
                      ? 'bg-blue-600 text-white rounded-tr-none'
                      : 'bg-white dark:bg-zinc-800 text-zinc-800 dark:text-zinc-100 border border-zinc-200/80 dark:border-zinc-700/80 rounded-tl-none'
                  }`}
                >
                  <p className="whitespace-pre-line">{msg.content}</p>

                  {/* Out of domain notice badge */}
                  {msg.isOutOfDomain && (
                    <div className="mt-2.5 pt-2 border-t border-zinc-200 dark:border-zinc-700 flex items-center gap-1.5 text-xs text-amber-600 dark:text-amber-400 font-semibold">
                      <AlertCircle size={13} />
                      <span>HUNTIQ AI only operates on marketplace listings and platform guidance.</span>
                    </div>
                  )}
                </div>
              </div>

              {/* ── Real Marketplace Listings Photographic Cards ── */}
              {msg.listings && msg.listings.length > 0 && (
                <div className="w-full mt-2 pl-10 pr-2">
                  <div className="flex items-center gap-2 mb-2">
                    <ShieldCheck size={14} className="text-emerald-500" />
                    <span className="text-xs font-black uppercase tracking-wider text-zinc-500 dark:text-zinc-400">
                      Verified Marketplace Listings ({msg.listings.length})
                    </span>
                  </div>

                  <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
                    {msg.listings.map((item) => (
                      <Link
                        key={item.id}
                        href={item.url || `/products/${item.id}`}
                        target="_blank"
                        className="group flex flex-col rounded-2xl border border-zinc-200/80 dark:border-zinc-800 bg-white dark:bg-zinc-900 overflow-hidden shadow-sm hover:shadow-md transition-all hover:-translate-y-0.5"
                      >
                        {/* Photograph Container */}
                        <div className="relative w-full h-36 bg-zinc-100 dark:bg-zinc-800 overflow-hidden">
                          {item.primary_image ? (
                            // eslint-disable-next-line @next/next/no-img-element
                            <img
                              src={item.primary_image}
                              alt={item.title}
                              className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                              onError={(e) => {
                                (e.target as HTMLImageElement).src = '/placeholder-image.png';
                              }}
                            />
                          ) : (
                            <div className="w-full h-full flex items-center justify-center text-zinc-400 text-xs">
                              No photograph
                            </div>
                          )}

                          {item.listing_type === 'rent' && (
                            <span className="absolute top-2 left-2 bg-purple-600 text-white text-[10px] font-black px-2 py-0.5 rounded-full uppercase tracking-wider shadow">
                              For Rent
                            </span>
                          )}
                        </div>

                        {/* Details */}
                        <div className="p-3 flex flex-col flex-1 justify-between">
                          <div>
                            <p className="text-xs font-black text-blue-600 dark:text-blue-400">
                              {item.price.toLocaleString()} {item.currency}
                            </p>
                            <h4 className="text-xs font-bold text-zinc-800 dark:text-zinc-100 line-clamp-1 group-hover:text-blue-600 transition-colors mt-0.5">
                              {item.title}
                            </h4>
                          </div>

                          <div className="mt-2 pt-2 border-t border-zinc-100 dark:border-zinc-800 flex items-center justify-between text-[10px] text-zinc-400">
                            <span>{item.location_name || 'Ethiopia'}</span>
                            <span className="flex items-center gap-0.5 text-blue-600 dark:text-blue-400 font-semibold group-hover:underline">
                              View <ChevronRight size={10} />
                            </span>
                          </div>
                        </div>
                      </Link>
                    ))}
                  </div>
                </div>
              )}

              {/* ── Suggested Follow-up Queries ── */}
              {msg.suggestedQueries && msg.suggestedQueries.length > 0 && (
                <div className="flex flex-wrap gap-1.5 pl-10 mt-1">
                  {msg.suggestedQueries.map((promptText, idx) => (
                    <button
                      key={idx}
                      onClick={() => handleSend(promptText)}
                      disabled={loading}
                      className="text-xs font-medium px-3 py-1.5 rounded-xl bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-200 hover:border-blue-500 hover:text-blue-600 dark:hover:text-blue-400 transition-all text-left"
                    >
                      💡 {promptText}
                    </button>
                  ))}
                </div>
              )}
            </div>
          );
        })}

        {/* Loading Indicator */}
        {loading && (
          <div className="flex items-center gap-2 text-zinc-500 dark:text-zinc-400 text-xs pl-2">
            <div className="w-7 h-7 rounded-xl bg-blue-600 text-white flex items-center justify-center animate-spin">
              <Sparkles size={14} />
            </div>
            <span className="font-semibold animate-pulse">
              HUNTIQ AI is retrieving verified listings & formulating response...
            </span>
          </div>
        )}

        <div ref={messagesEndRef} />
      </div>

      {/* ── Bottom Input Bar ── */}
      <div className="p-3 sm:p-4 border-t border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-950 flex items-center gap-2">
        <div className="relative flex-1">
          <input
            ref={inputRef}
            type="text"
            value={inputText}
            onChange={(e) => setInputText(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') handleSend();
            }}
            placeholder={
              selectedLang === 'am'
                ? 'በHUNTIQ ላይ እቃዎችን፣ ቤቶችን፣ መኪኖችን ወይንም አገልግሎቶችን ይጠይቁ...'
                : selectedLang === 'ti'
                ? 'ንብረታት፣ ኣባይቲ ወይ ኣገልግሎታት ኣብ HUNTIQ ሕተቱ...'
                : selectedLang === 'or'
                ? 'Meeshaalee, manneen ykn tajaajiloota HUNTIQ irratti gaafadhaa...'
                : 'Search products, rent a house, find services, or ask HUNTIQ policies...'
            }
            disabled={loading}
            className="w-full text-sm pl-4 pr-10 py-3 rounded-2xl border border-zinc-300 dark:border-zinc-700 bg-zinc-50/80 dark:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-blue-600 text-zinc-900 dark:text-zinc-100 disabled:opacity-50"
          />
          <Search size={16} className="absolute right-3.5 top-1/2 -translate-y-1/2 text-zinc-400" />
        </div>

        <button
          onClick={() => handleSend()}
          disabled={loading || !inputText.trim()}
          className="p-3 rounded-2xl bg-blue-600 hover:bg-blue-700 text-white font-bold transition-all disabled:opacity-40 disabled:hover:bg-blue-600 shadow-md shadow-blue-600/20"
        >
          <Send size={18} />
        </button>
      </div>

      {/* ── Privacy & Domain Disclaimer ── */}
      <div className="px-4 py-2 bg-zinc-100 dark:bg-zinc-900/80 border-t border-zinc-200/60 dark:border-zinc-800 text-[10px] text-zinc-400 flex flex-col sm:flex-row items-center justify-between gap-1">
        <span>HUNTIQ AI provides automated marketplace discovery from real platform records.</span>
        <span className="font-semibold text-zinc-500">Official Assistant · Made by HUNTIQ</span>
      </div>
    </div>
  );
}
