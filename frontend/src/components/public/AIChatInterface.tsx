'use client';

import React, { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { api } from '@/lib/api';
import { SearchAIMeta } from '@/lib/api/public';
import { Product } from '@/types';
import {
  Sparkles,
  Send,
  ArrowLeft,
  Brain,
  ChevronDown,
  ChevronRight,
  ShieldCheck,
  AlertCircle,
  HelpCircle,
  X,
  RotateCcw,
  Tag,
  Star,
  CheckCircle2,
  ExternalLink,
} from 'lucide-react';

interface ChatMessage {
  id: string;
  sender: 'user' | 'ai';
  text: string;
  timestamp: string;
  meta?: SearchAIMeta;
  products?: Product[];
}

interface AIChatInterfaceProps {
  onExit: () => void;
  initialQuery?: string;
  locale: string;
}

const FAQ_SUGGESTIONS_BY_LOCALE: Record<string, Array<{ icon: string; label: string; prompt: string }>> = {
  en: [
    { icon: '📱', label: 'iPhone 15 Pro Price Check', prompt: 'What is the fair market price for an Apple iPhone 15 Pro in Addis Ababa?' },
    { icon: '🚗', label: 'Toyota Vitz Libre Verification', prompt: 'Toyota Vitz: What legal papers and Libre ownership documents should I check before paying?' },
    { icon: '🏠', label: 'Apartments in Bole / CMC', prompt: 'Find furnished 2-bedroom rental apartments in Bole or CMC under 40,000 ETB.' },
    { icon: '🛡️', label: 'Electronics Anti-Fraud Checklist', prompt: 'How to safely test used electronics and verify phone IMEI to avoid scams in Ethiopia?' },
    { icon: '👑', label: 'Verified Merchant Deals', prompt: 'Show me active listings from Gold and Premium verified sellers in Addis Ababa.' },
    { icon: '💻', label: 'MacBook Pro for Developers', prompt: 'Looking for a MacBook Pro M3 or powerful workstation laptop with warranty.' },
  ],
  am: [
    { icon: '📱', label: 'የአይፎን 15 ፕሮ ዋጋ', prompt: 'በአዲስ አበባ ውስጥ የአፕል አይፎን 15 ፕሮ ትክክለኛ የገበያ ዋጋ ስንት ነው?' },
    { icon: '🚗', label: 'የቶዮታ ቪትዝ ሊብሬ ምርመራ', prompt: 'ቶዮታ ቪትዝ፡ ከመክፈሌ በፊት ምን አይነት ህጋዊ ሰነዶችን እና ሊብሬ ማረጋገጥ አለብኝ?' },
    { icon: '🏠', label: 'አፓርታማ በቦሌ / ሲኤምሲ', prompt: 'በቦሌ ወይም ሲኤምሲ ከ40,000 ብር በታች ባለ 2 መኝታ የተሟላ ኪራይ ቤት ፈልግልኝ።' },
    { icon: '🛡️', label: 'የኤሌክትሮኒክስ ማጭበርበር መከላከያ', prompt: 'ስልኮችን ስገዛ ማጭበርበርን ለመከላከል የIMEI ቁጥር እና የደህንነት ሁኔታ እንዴት ማረጋገጥ እችላለሁ?' },
    { icon: '👑', label: 'የተረጋገጡ ነጋዴዎች', prompt: 'በአዲስ አበባ ውስጥ ከወርቅ እና ፕሪሚየም የተረጋገጡ ሻጮች የቀረቡ እቃዎችን አሳየኝ።' },
    { icon: '💻', label: 'ላፕቶፕ ለስራ', prompt: 'ዋስትና ያለው ማክቡክ ፕሮ ወይም ጠንካራ ላፕቶፕ እፈልጋለሁ።' },
  ],
  or: [
    { icon: '📱', label: 'Gatii iPhone 15 Pro', prompt: 'Finfinnee keessatti gatiin gabaa sirrii Apple iPhone 15 Pro meeqa?' },
    { icon: '🚗', label: 'Qorannoo Liibiree Toyota Vitz', prompt: 'Toyota Vitz: Kaffaluu koo dura waraqaalee seeraa fi liibiree akkamii qorachuu qaba?' },
    { icon: '🏠', label: 'Mana Kireeffamu Boolee / CMC', prompt: 'Boolee ykn CMC keessatti mana kireeffamu kutaa 2 qabu gatiin 40,000 ETB gadi ta\'e naaf barbaadi.' },
    { icon: '🛡️', label: 'Hanna Meeshaalee Elektirooniksii', prompt: 'Hanna irraa of eeguuf bilbila yoon bitu IMEI fi nageenya isaa akkamittan mirkaneessa?' },
    { icon: '👑', label: 'Gurgurtoota Mirkanaa\'an', prompt: 'Gurgurtoota sadarkaa Warqii fi Piriimiyem Finfinnee keessa jiran irraa oomisha jiru na agarsiisi.' },
    { icon: '💻', label: 'Laaptooppii Hojii', prompt: 'MacBook Pro ykn laaptooppii hojii qulqullina fi wabii qabu barbaada.' },
  ],
  ti: [
    { icon: '📱', label: 'ዋጋ ኣይፎን 15 ፕሮ', prompt: 'ኣብ ኣዲስ ኣበባ ትኽክለኛ ናይ ዕዳጋ ዋጋ ኣፕል ኣይፎን 15 ፕሮ ክንደይ እዩ?' },
    { icon: '🚗', label: 'ናይ ቶዮታ ቪትዝ ሊብሬ መርመራ', prompt: 'ቶዮታ ቪትዝ፡ ቅድሚ ምኽፋለይ እንታይ ሕጋዊ ሰነዳትን ሊብሬን ክምርምር ኣለኒ?' },
    { icon: '🏠', label: 'ኣፓርታማ ኣብ ቦሌ / ሲኤምሲ', prompt: 'ኣብ ቦሌ ወይ ሲኤምሲ ትሕቲ 40,000 ቅርሺ ዝካረ ባለ 2 መደቀሲ ኣፓርታማ ድለየለይ።' },
    { icon: '🛡️', label: 'ምክልኻል ምትላል ኤሌክትሮኒክስ', prompt: 'ምትላል ንምክልኻል ተንቀሳቓሲ ስልኪ ክገዝእ ከለኹ IMEI ቁጽርን ደሕንነትን ብኸመይ ይረጋገጽ?' },
    { icon: '👑', label: 'ዝተረጋገፁ ነጋዶ', prompt: 'ኣብ ኣዲስ ኣበባ ካብ ወርቂን ፕሪሚየምን ዝተረጋገፁ ነጋዶ ዘለዉ ፍርያት ኣርእየኒ።' },
    { icon: '💻', label: 'ላፕቶፕ ንስራሕ', prompt: 'ውሕስና ዘለዎ ማክቡክ ፕሮ ወይ ጽኑዕ ላፕቶፕ እደሊ ኣለኹ።' },
  ],
};

const WELCOME_MESSAGES: Record<string, string> = {
  en: 'Hello! I am your HUNTIQ Super AI Marketplace Advisor. I can analyze market prices, inspect seller trust tiers, guide vehicle Libre checks, and verify electronic specs in English, አማርኛ, Afaan Oromoo, or ትግርኛ. What are you looking to buy, rent, or investigate today?',
  am: 'ሰላም! እኔ የHUNTIQ ሱፐር AI የገበያ አማካሪ ነኝ። የገበያ ዋጋን ለመተንተን፣ የሻጮችን አስተማማኝነት ደረጃ ለመፈተሽ፣ የመኪና ሊብሬ ምርመራዎችን ለመምራት እና የኤሌክትሮኒክስ እቃዎችን ለማረጋገጥ ዝግጁ ነኝ። ዛሬ ምን መግዛት፣ መከራየት ወይም መመርመር ይፈልጋሉ?',
  or: 'Akkam! Ani gorsa gabaa HUNTIQ Super AI ti. Gatii gabaa xiinxaluu, sadarkaa amanamummaa gurgurtootaa qorachuu, qorannoo liibiree konkolaataa geggeessuu fi qulqullina meeshaalee mirkaneessuuf qophiidha. Har\'a maal bituu, kireeffachuu ykn qorachuu barbaaddu?',
  ti: 'ሰላም! ኣነ ናይ HUNTIQ ሱፐር AI ናይ ዕዳጋ ኣማኻሪ እየ። ናይ ዕዳጋ ዋጋታት ንምትንታን፣ ደረጃ ተኣማንነት ነጋዶ ንምፍታሽ፣ ናይ መኪና ሊብሬ መርመራ ንምምራሕን ናይ ኤሌክትሮኒክስ ፍርያት ንምርግጋጽን ድሉው እየ። ሎሚ እንታይ ክትገዝኡ፣ ክትካረዩ ወይ ክትምርምሩ ትደልዩ?',
};

const RESET_MESSAGES: Record<string, string> = {
  en: 'Chat history cleared. How else can I help your marketplace search today?',
  am: 'የውይይት ታሪክ ጸድቷል። ዛሬ በገበያ ፍለጋዎ በምን ልርዳዎት?',
  or: 'Seenaan haasaa haqameera. Har\'a gabaa keessatti akkamittan isin gargaaruu danda\'a?',
  ti: 'ናይ ዕላል ታሪኽ ተደምሲሱ። ሎሚ ኣብ ናይ ዕዳጋ ድላይኩም ብምንታይ ክሕግዘኩም?',
};

const INTERFACE_UI: Record<string, Record<string, string>> = {
  en: {
    backCatalog: 'Back to Catalog',
    title: 'HUNTIQ AI Advisor',
    badge: 'Live Marketplace',
    subtitle: 'Ethiopian Marketplace Intelligence · Live Database Search',
    safetyBtn: 'Buyer Safety & Guide',
    quickFaq: 'Quick Inquiries & Market FAQ Starters',
    clickPrompt: 'Click any prompt to ask instantly',
    placeholder: 'Ask HUNTIQ AI anything in English, አማርኛ, Afaan Oromoo, or ትግርኛ...',
    ask: 'Ask',
    safetyTitle: 'HUNTIQ Buyer Safety & Market Guide',
    understood: 'Understood',
  },
  am: {
    backCatalog: 'ወደ ካታሎግ ተመለስ',
    title: 'HUNTIQ AI አማካሪ',
    badge: 'የቀጥታ ገበያ',
    subtitle: 'የኢትዮጵያ የገበያ መረጃ · የቀጥታ ዳታቤዝ ፍለጋ',
    safetyBtn: 'የገዢዎች ደህንነት መመሪያ',
    quickFaq: 'ፈጣን ጥያቄዎች እና የገበያ መመሪያዎች',
    clickPrompt: 'ወዲያውኑ ለመጠየቅ ማንኛውንም ይጫኑ',
    placeholder: 'የሚፈልጉትን ማንኛውንም ጥያቄ ይጠይቁ (ለምሳሌ፦ የአይፎን 15 ዋጋ፣ የቪትዝ ሊብሬ ምርመራ...)',
    ask: 'ጠይቅ',
    safetyTitle: 'የHUNTIQ የገዢዎች ደህንነት እና የገበያ መመሪያ',
    understood: 'ተረድቻለሁ',
  },
  or: {
    backCatalog: 'Gara Kaatalaagiitti Deebi\'aa',
    title: 'Gorsa HUNTIQ AI',
    badge: 'Gabaa Kallattii',
    subtitle: 'Odeeffannoo Gabaa Itoophiyaa · Sakatta\'a Kuusaa Kallattii',
    safetyBtn: 'Nageenya Bittootaa',
    quickFaq: 'Gaaffilee Ariifachiisaa fi Gorsa Gabaa',
    clickPrompt: 'Battalumatti gaafachuuf cuqaasaa',
    placeholder: 'Waan barbaaddan kamiyyuu gaafadhaa (fkn gatii iPhone 15, qorannoo konkolaataa)...',
    ask: 'Gaafadhu',
    safetyTitle: 'Qajeelfama Nageenya Bittootaa HUNTIQ',
    understood: 'Hubadheera',
  },
  ti: {
    backCatalog: 'ናብ ካታሎግ ተመለስ',
    title: 'HUNTIQ AI ኣማኻሪ',
    badge: 'ቀጥታዊ ዕዳጋ',
    subtitle: 'ናይ ኢትዮጵያ ዕዳጋ መረዳእታ · ቀጥታዊ ናይ ዳታቤዝ ድላይ',
    safetyBtn: 'ደሕንነት ገዛእቲ',
    quickFaq: 'ቅልጡፍ ሕቶታትን መምርሒታትን',
    clickPrompt: 'ብቕጽበት ንምሕታት ጠውቑ',
    placeholder: 'ዝደለይዎ ሕቶ ሕተቱ (ንኣብነት፦ ዋጋ ኣይፎን 15፣ መርመራ ሊብሬ ቪትዝ...)',
    ask: 'ሕተት',
    safetyTitle: 'ናይ HUNTIQ ናይ ገዛእቲ ደሕንነት መምርሒ',
    understood: 'ተረዲኡኒ',
  },
};

interface ConversationContext {
  lastItem?: string;
  lastBrand?: string;
  lastMaxPrice?: number;
  lastMinPrice?: number;
  lastLocation?: string;
}

export function AIChatInterface({ onExit, initialQuery = '', locale }: AIChatInterfaceProps) {
  const activeLocale = (locale === 'om' ? 'or' : locale) || 'en';
  const ui = INTERFACE_UI[activeLocale] || INTERFACE_UI.en;
  const faqs = FAQ_SUGGESTIONS_BY_LOCALE[activeLocale] || FAQ_SUGGESTIONS_BY_LOCALE.en;
  const welcomeText = WELCOME_MESSAGES[activeLocale] || WELCOME_MESSAGES.en;

  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputVal, setInputVal] = useState('');
  const [loading, setLoading] = useState(false);
  const [openReasoningMap, setOpenReasoningMap] = useState<Record<string, boolean>>({});
  const [showFaqModal, setShowFaqModal] = useState(false);
  const chatBottomRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const contextRef = useRef<ConversationContext>({});

  // Auto-scroll chat to bottom
  const scrollToBottom = () => {
    setTimeout(() => {
      chatBottomRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, 100);
  };

  // Initial welcome message or run initial query
  useEffect(() => {
    if (initialQuery.trim()) {
      handleSend(initialQuery.trim());
    } else {
      setMessages([
        {
          id: 'welcome',
          sender: 'ai',
          text: welcomeText,
          timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
        },
      ]);
    }
  }, []);

  const resolveEffectiveQuery = (rawQuery: string): string => {
    const q = rawQuery.trim();
    const lower = q.toLowerCase();
    const ctx = contextRef.current;

    // Follow-up refinement detection (e.g., "only lenovo", "cheaper one", "how many", "in bole")
    const isRefinement =
      lower.startsWith('only ') ||
      lower.startsWith('ብቻ ') ||
      lower.startsWith('qofa ') ||
      lower.startsWith('ጥራይ ') ||
      lower.startsWith('show me ') ||
      lower.startsWith('cheaper') ||
      lower.startsWith('what about') ||
      lower.startsWith('how many') ||
      lower.startsWith('ስንት') ||
      lower.startsWith('meeqa') ||
      lower.startsWith('ክንደይ') ||
      lower.startsWith('compare') ||
      lower.startsWith('አወዳድር') ||
      lower.split(/\s+/).length <= 2;

    if (isRefinement && ctx.lastItem) {
      // If query does not contain the previous item, blend it
      if (!lower.includes(ctx.lastItem.toLowerCase())) {
        return `${q} ${ctx.lastItem}`;
      }
    }
    return q;
  };

  const handleSend = async (questionText: string) => {
    const query = questionText.trim();
    if (!query || loading) return;

    const userMsgId = `user-${Date.now()}`;
    const userMsg: ChatMessage = {
      id: userMsgId,
      sender: 'user',
      text: query,
      timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
    };

    setMessages(prev => [...prev, userMsg]);
    setInputVal('');
    setLoading(true);
    scrollToBottom();

    const effectiveQuery = resolveEffectiveQuery(query);

    try {
      const res = await api.getProductsWithMeta({
        search: effectiveQuery,
        mode: 'ai',
        locale,
        preferred_language: locale,
      });

      // Update multi-turn conversational context
      if (res.meta?.cleaned_search_term) {
        contextRef.current.lastItem = res.meta.cleaned_search_term;
      }
      if (res.meta?.deep_parameters?.extracted_entities?.brand) {
        contextRef.current.lastBrand = res.meta.deep_parameters.extracted_entities.brand;
      }
      if (res.meta?.deep_parameters?.extracted_entities?.price_range?.max_price) {
        contextRef.current.lastMaxPrice =
          res.meta.deep_parameters.extracted_entities.price_range.max_price;
      }
      if (res.meta?.deep_parameters?.extracted_entities?.location) {
        contextRef.current.lastLocation =
          res.meta.deep_parameters.extracted_entities.location;
      }

      const aiText =
        res.meta?.localized_answer ||
        res.meta?.market_verdict ||
        (activeLocale === 'am'
          ? `የቀጥታ ዳታቤዝ ዝርዝሮችን ለ"${query}" ተንትኛለሁ።`
          : activeLocale === 'om'
          ? `Tarreeffamoota gabaa "${query}" tiif xiinxaleera.`
          : activeLocale === 'ti'
          ? `ናይ ቀጥታ ዝርዝራት ን"${query}" ተተንቲኑ ኣሎ።`
          : `I analyzed our live database listings for "${query}". Here are the market insights and verified options.`);

      const aiMsgId = `ai-${Date.now()}`;
      const aiMsg: ChatMessage = {
        id: aiMsgId,
        sender: 'ai',
        text: aiText,
        meta: res.meta,
        products: res.products,
        timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      };

      setMessages(prev => [...prev, aiMsg]);
      // Open reasoning accordion by default for fresh answers
      if (res.meta?.reasoning_chain) {
        setOpenReasoningMap(prev => ({ ...prev, [aiMsgId]: true }));
      }
    } catch {
      const errorMsgText =
        activeLocale === 'am'
          ? 'ከ-AI ረዳት ጋር መገናኘት አልተቻለም። እባክዎ ጥቂት ቆይተው እንደገና ይሞክሩ።'
          : activeLocale === 'om'
          ? "Gargaaraa AI wajjin walqunnamuu hin dandeenye. Maaloo xiqqoo turanii irra deebi'aa yaalaa."
          : activeLocale === 'ti'
          ? 'ምስ AI ሓጋዚ ምርኻብ ኣይተኻእለን። በጃኹም ቁሩብ ጸኒሕኩም ደጊምኩም ፈትኑ።'
          : 'I encountered an issue connecting to the AI inference worker. Please make sure the service is online or try asking your question again.';
      const errorMsg: ChatMessage = {
        id: `err-${Date.now()}`,
        sender: 'ai',
        text: errorMsgText,
        timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      };
      setMessages(prev => [...prev, errorMsg]);
    } finally {
      setLoading(false);
      scrollToBottom();
      inputRef.current?.focus();
    }
  };

  const toggleReasoning = (msgId: string) => {
    setOpenReasoningMap(prev => ({ ...prev, [msgId]: !prev[msgId] }));
  };

  const clearChat = () => {
    contextRef.current = {};
    setMessages([
      {
        id: 'welcome-reset',
        sender: 'ai',
        text: RESET_MESSAGES[activeLocale] || RESET_MESSAGES.en,
        timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      },
    ]);
  };

  return (
    <div className="flex flex-col h-[calc(100vh-68px)] max-h-[calc(100vh-68px)] bg-zinc-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-100 transition-colors">
      {/* ── Top Bar ── */}
      <div className="shrink-0 flex items-center justify-between px-4 sm:px-6 py-3 bg-white/80 dark:bg-zinc-900/80 backdrop-blur-md border-b border-zinc-200 dark:border-zinc-800 z-20">
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={onExit}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-zinc-200 dark:border-zinc-700 hover:bg-zinc-100 dark:hover:bg-zinc-800 text-xs font-bold transition-all text-zinc-700 dark:text-zinc-300"
          >
            <ArrowLeft size={14} />
            <span>{ui.backCatalog}</span>
          </button>

          <div className="h-4 w-[1px] bg-zinc-200 dark:bg-zinc-800 hidden sm:block" />

          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-orange-600 to-amber-500 text-white flex items-center justify-center shadow-md shadow-orange-600/30">
              <Sparkles size={16} className="animate-pulse" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <span className="font-black text-sm text-zinc-900 dark:text-white">
                  {ui.title}
                </span>
                <span className="px-1.5 py-0.5 rounded text-[9px] font-black uppercase tracking-wider bg-orange-100 dark:bg-orange-950/60 text-orange-600 dark:text-orange-400 border border-orange-200 dark:border-orange-800/60">
                  {ui.badge}
                </span>
              </div>
              <p className="text-[10px] text-zinc-400 hidden sm:block">
                {ui.subtitle}
              </p>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => setShowFaqModal(true)}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-bold text-orange-600 dark:text-orange-400 hover:bg-orange-50 dark:hover:bg-orange-950/30 transition-colors"
          >
            <HelpCircle size={14} />
            <span className="hidden sm:inline">{ui.safetyBtn}</span>
          </button>

          <button
            type="button"
            onClick={clearChat}
            title="Reset Conversation"
            className="p-1.5 rounded-xl text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
          >
            <RotateCcw size={15} />
          </button>
        </div>
      </div>

      {/* ── Main Chat Scrollable Stream ── */}
      <div className="flex-1 overflow-y-auto px-4 sm:px-6 py-6 space-y-6">
        <div className="max-w-4xl mx-auto space-y-6">
          {/* Quick FAQ Suggestion Chips at the Top */}
          <div className="p-4 rounded-2xl bg-gradient-to-r from-orange-500/10 via-amber-500/10 to-teal-500/10 dark:from-orange-950/20 dark:via-zinc-900 dark:to-teal-950/20 border border-orange-200/60 dark:border-orange-900/40">
            <div className="flex items-center justify-between mb-2">
              <span className="text-xs font-black uppercase tracking-wider text-orange-700 dark:text-orange-400 flex items-center gap-1.5">
                <Sparkles size={13} />
                {ui.quickFaq}
              </span>
              <span className="text-[10px] font-bold text-zinc-500">
                {ui.clickPrompt}
              </span>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-2">
              {faqs.map((faq, idx) => (
                <button
                  key={idx}
                  type="button"
                  onClick={() => handleSend(faq.prompt)}
                  className="flex items-center gap-2 p-2.5 rounded-xl bg-white dark:bg-zinc-900 border border-zinc-200/80 dark:border-zinc-800 hover:border-orange-400 dark:hover:border-orange-500 hover:shadow-md transition-all text-left group"
                >
                  <span className="text-lg shrink-0 group-hover:scale-110 transition-transform">
                    {faq.icon}
                  </span>
                  <div className="min-w-0 flex-1">
                    <p className="text-xs font-bold text-zinc-800 dark:text-zinc-200 truncate">
                      {faq.label}
                    </p>
                    <p className="text-[10px] text-zinc-400 truncate">{faq.prompt}</p>
                  </div>
                </button>
              ))}
            </div>
          </div>

          {/* Conversation Messages */}
          {messages.map(msg => {
            const isUser = msg.sender === 'user';
            const msgReasoningOpen = openReasoningMap[msg.id] ?? false;

            return (
              <div
                key={msg.id}
                className={`flex gap-3 ${isUser ? 'justify-end' : 'justify-start'}`}
              >
                {!isUser && (
                  <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-orange-600 to-amber-500 text-white flex items-center justify-center shrink-0 shadow-sm mt-1">
                    <Sparkles size={15} />
                  </div>
                )}

                <div
                  className={`max-w-[88%] sm:max-w-[80%] rounded-2xl p-4 shadow-sm space-y-3 ${
                    isUser
                      ? 'bg-orange-600 text-white rounded-tr-none'
                      : 'bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-tl-none'
                  }`}
                >
                  {/* Message Text */}
                  <div className="text-sm font-medium leading-relaxed whitespace-pre-line">
                    {msg.text}
                  </div>

                  {/* If AI Message: Real-Time DB Analysis Card */}
                  {!isUser &&
                    msg.meta?.real_time_analysis &&
                    msg.meta.real_time_analysis.total_analyzed > 0 && (
                      <div className="pt-3 border-t border-zinc-100 dark:border-zinc-800 space-y-2.5">
                        <div className="flex flex-wrap items-center justify-between gap-2">
                          <span className="text-[11px] font-black uppercase tracking-wider text-orange-600 dark:text-orange-400 flex items-center gap-1.5">
                            📊 Live Market Intelligence (
                            {msg.meta?.total_matching !== undefined &&
                            msg.meta.total_matching > msg.meta.real_time_analysis.total_analyzed
                              ? `${msg.meta.total_matching} Matches · ${msg.meta.real_time_analysis.total_analyzed} Analyzed`
                              : `${msg.meta.real_time_analysis.total_analyzed} Database Listings`}
                            )
                          </span>
                          <div className="flex items-center gap-1.5 text-[10px] font-black">
                            <span className="px-2 py-0.5 rounded-full bg-amber-100 dark:bg-amber-900/40 text-amber-800 dark:text-amber-300 border border-amber-300/50">
                              👑 {msg.meta.real_time_analysis.premium_count} Premium Verified
                            </span>
                            <span className="px-2 py-0.5 rounded-full bg-yellow-100 dark:bg-yellow-900/40 text-yellow-800 dark:text-yellow-300 border border-yellow-300/50">
                              ⭐ {msg.meta.real_time_analysis.gold_count} Gold Certified
                            </span>
                          </div>
                        </div>

                        {/* Price Statistics */}
                        {msg.meta.real_time_analysis.price_stats.avg_price > 0 && (
                          <div className="grid grid-cols-3 gap-2 text-xs">
                            <div className="bg-zinc-50 dark:bg-zinc-800/60 rounded-xl p-2 border border-zinc-200/60 dark:border-zinc-700/60 text-center">
                              <span className="text-[10px] text-zinc-400 block">Lowest</span>
                              <span className="font-bold text-emerald-600 dark:text-emerald-400">
                                {Number(
                                  msg.meta.real_time_analysis.price_stats.min_price
                                ).toLocaleString()}{' '}
                                ETB
                              </span>
                            </div>
                            <div className="bg-orange-50/60 dark:bg-orange-950/40 rounded-xl p-2 border border-orange-200 dark:border-orange-900/40 text-center">
                              <span className="text-[10px] text-orange-600 dark:text-orange-400 block font-bold">
                                Average
                              </span>
                              <span className="font-black text-orange-600 dark:text-orange-400">
                                {Number(
                                  msg.meta.real_time_analysis.price_stats.avg_price
                                ).toLocaleString()}{' '}
                                ETB
                              </span>
                            </div>
                            <div className="bg-zinc-50 dark:bg-zinc-800/60 rounded-xl p-2 border border-zinc-200/60 dark:border-zinc-700/60 text-center">
                              <span className="text-[10px] text-zinc-400 block">Highest</span>
                              <span className="font-bold text-zinc-800 dark:text-zinc-200">
                                {Number(
                                  msg.meta.real_time_analysis.price_stats.max_price
                                ).toLocaleString()}{' '}
                                ETB
                              </span>
                            </div>
                          </div>
                        )}

                        {/* Top Verified Highlights */}
                        {(msg.meta.real_time_analysis.top_premium_pick ||
                          msg.meta.real_time_analysis.top_gold_pick) && (
                          <div className="flex flex-wrap gap-2 pt-1 text-xs">
                            {msg.meta.real_time_analysis.top_premium_pick && (
                              <div className="flex items-center gap-1.5 px-3 py-1 rounded-xl bg-amber-500/10 border border-amber-500/30 text-amber-900 dark:text-amber-200 text-[11px]">
                                <span className="font-black">👑 Top Premium Pick:</span>
                                <span className="font-semibold truncate max-w-[180px]">
                                  {msg.meta.real_time_analysis.top_premium_pick.title}
                                </span>
                                <span className="font-black text-orange-600 dark:text-orange-400">
                                  {Number(
                                    msg.meta.real_time_analysis.top_premium_pick.price
                                  ).toLocaleString()}{' '}
                                  ETB
                                </span>
                              </div>
                            )}
                            {msg.meta.real_time_analysis.top_gold_pick && (
                              <div className="flex items-center gap-1.5 px-3 py-1 rounded-xl bg-yellow-500/10 border border-yellow-500/30 text-yellow-900 dark:text-yellow-200 text-[11px]">
                                <span className="font-black">⭐ Top Gold Pick:</span>
                                <span className="font-semibold truncate max-w-[180px]">
                                  {msg.meta.real_time_analysis.top_gold_pick.title}
                                </span>
                                <span className="font-black text-amber-600 dark:text-amber-400">
                                  {Number(
                                    msg.meta.real_time_analysis.top_gold_pick.price
                                  ).toLocaleString()}{' '}
                                  ETB
                                </span>
                              </div>
                            )}
                          </div>
                        )}
                      </div>
                    )}

                  {/* If AI Message: Expandable 4-Step Chain-of-Thought Reasoning */}
                  {!isUser &&
                    msg.meta?.reasoning_chain &&
                    msg.meta.reasoning_chain.steps &&
                    msg.meta.reasoning_chain.steps.length > 0 && (
                      <div className="pt-2 border-t border-zinc-100 dark:border-zinc-800">
                        <button
                          type="button"
                          onClick={() => toggleReasoning(msg.id)}
                          className="flex items-center gap-1.5 text-xs font-bold text-orange-600 dark:text-orange-400 hover:text-orange-700 transition-colors"
                        >
                          <Brain size={14} />
                          <span>
                            {msgReasoningOpen
                              ? 'Hide AI Logical Reasoning'
                              : 'View AI Multi-Step Logical Reasoning (Chain-of-Thought)'}
                          </span>
                          {msgReasoningOpen ? <ChevronDown size={12} /> : <ChevronRight size={12} />}
                        </button>

                        {msgReasoningOpen && (
                          <div className="mt-2.5 p-3 rounded-xl bg-orange-50/50 dark:bg-zinc-950/60 border border-orange-200/50 dark:border-zinc-800 space-y-2.5 text-xs">
                            {msg.meta.reasoning_chain.steps.map(st => (
                              <div key={st.step} className="flex flex-col gap-0.5">
                                <span className="font-black text-orange-700 dark:text-orange-400 flex items-center gap-1">
                                  {st.step === 1 && '🎯'}
                                  {st.step === 2 && '📈'}
                                  {st.step === 3 && '⚖️'}
                                  {st.step === 4 && '🛡️'}
                                  {st.title}
                                </span>
                                <p className="text-zinc-600 dark:text-zinc-300 pl-4 border-l-2 border-orange-300 dark:border-orange-800/80 leading-relaxed">
                                  {st.detail}
                                </p>
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    )}

                  {/* If AI Message: Matching Live Database Products */}
                  {!isUser && msg.products && msg.products.length > 0 && (
                    <div className="pt-3 border-t border-zinc-100 dark:border-zinc-800 space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="text-xs font-black uppercase tracking-wider text-zinc-500 flex items-center gap-1">
                          <Tag size={12} />
                          Recommended Database Listings ({msg.products.length})
                        </span>
                      </div>

                      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        {msg.products.map(p => {
                          const thumb =
                            p.images?.[0]?.image_url ||
                            'https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=400';
                          const tier = p.seller_tier || 'normal';
                          const deal = msg.meta?.real_time_analysis?.deal_ratings?.find(
                            d => d.product_id === p.id
                          );

                          return (
                            <div
                              key={p.id}
                              className="flex gap-2.5 p-2 rounded-xl bg-zinc-50 dark:bg-zinc-800/60 border border-zinc-200 dark:border-zinc-700/60 hover:border-orange-400 transition-all group"
                            >
                              <div className="relative w-16 h-16 rounded-lg overflow-hidden shrink-0 bg-zinc-200 dark:bg-zinc-700">
                                <Image
                                  src={thumb}
                                  alt={p.title}
                                  fill
                                  sizes="64px"
                                  className="object-cover group-hover:scale-105 transition-transform"
                                />
                              </div>
                              <div className="flex-1 min-w-0 flex flex-col justify-between">
                                <div>
                                  <div className="flex items-center gap-1 mb-0.5 flex-wrap">
                                    {tier === 'premium' ? (
                                      <span className="px-1.5 py-0.2 rounded text-[8px] font-black bg-amber-500 text-white uppercase">
                                        👑 Premium
                                      </span>
                                    ) : tier === 'gold' ? (
                                      <span className="px-1.5 py-0.2 rounded text-[8px] font-black bg-yellow-500 text-black uppercase">
                                        ⭐ Gold
                                      </span>
                                    ) : null}
                                    {deal && (
                                      <span
                                        className={`px-1.5 py-0.2 rounded text-[8px] font-black uppercase ${
                                          deal.deal_score >= 80
                                            ? 'bg-emerald-600 text-white'
                                            : deal.deal_score >= 60
                                            ? 'bg-blue-600 text-white'
                                            : 'bg-zinc-600 text-white'
                                        }`}
                                      >
                                        {deal.badge} ({deal.deal_score})
                                      </span>
                                    )}
                                    <span className="text-[9px] text-zinc-400 capitalize">
                                      {p.condition}
                                    </span>
                                  </div>
                                  <p className="text-xs font-bold text-zinc-900 dark:text-zinc-100 truncate">
                                    {p.title}
                                  </p>
                                </div>
                                <div className="flex items-center justify-between mt-1">
                                  <span className="font-black text-xs text-orange-600 dark:text-orange-400">
                                    {Number(p.price).toLocaleString()} ETB
                                  </span>
                                  <Link
                                    href={`/products/${p.id}`}
                                    className="flex items-center gap-0.5 text-[10px] font-black text-orange-600 hover:underline"
                                  >
                                    <span>View</span>
                                    <ExternalLink size={10} />
                                  </Link>
                                </div>
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  )}

                  <div className="flex items-center justify-end gap-1 text-[9px] opacity-40">
                    <span>{msg.timestamp}</span>
                  </div>
                </div>
              </div>
            );
          })}

          {/* Loading Indicator */}
          {loading && (
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-orange-600 to-amber-500 text-white flex items-center justify-center shrink-0 animate-pulse">
                <Brain size={16} />
              </div>
              <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl rounded-tl-none p-3 shadow-sm flex items-center gap-3 text-xs text-zinc-500 font-medium">
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-orange-500 rounded-full animate-bounce [animation-delay:-0.3s]" />
                  <div className="w-2 h-2 bg-orange-500 rounded-full animate-bounce [animation-delay:-0.15s]" />
                  <div className="w-2 h-2 bg-orange-500 rounded-full animate-bounce" />
                </div>
                <span>
                  HUNTIQ AI is deconstructing query, checking market feasibility & querying live database...
                </span>
              </div>
            </div>
          )}

          <div ref={chatBottomRef} />
        </div>
      </div>

      {/* ── Fixed Chat Input Bar ── */}
      <div className="shrink-0 p-4 bg-white/90 dark:bg-zinc-900/90 backdrop-blur-md border-t border-zinc-200 dark:border-zinc-800">
        <form
          onSubmit={e => {
            e.preventDefault();
            handleSend(inputVal);
          }}
          className="max-w-4xl mx-auto flex items-center gap-2"
        >
          <input
            ref={inputRef}
            type="text"
            value={inputVal}
            onChange={e => setInputVal(e.target.value)}
            disabled={loading}
            placeholder={ui.placeholder}
            className="flex-1 px-4 py-3 rounded-2xl bg-zinc-100 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 text-sm outline-none focus:ring-2 focus:ring-orange-600 transition-all"
          />
          <button
            type="submit"
            disabled={!inputVal.trim() || loading}
            className="px-5 py-3 rounded-2xl bg-gradient-to-r from-orange-600 via-amber-600 to-orange-500 hover:from-orange-700 hover:to-amber-600 disabled:opacity-40 text-white font-black text-sm shadow-md shadow-orange-600/30 flex items-center gap-1.5 transition-all active:scale-95"
          >
            <span>{ui.ask}</span>
            <Send size={15} />
          </button>
        </form>
      </div>

      {/* ── Safety & FAQ Modal ── */}
      {showFaqModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="w-full max-w-lg bg-white dark:bg-zinc-900 rounded-3xl p-6 shadow-2xl border border-zinc-200 dark:border-zinc-800 space-y-4 max-h-[85vh] overflow-y-auto">
            <div className="flex items-center justify-between pb-3 border-b border-zinc-100 dark:border-zinc-800">
              <div className="flex items-center gap-2">
                <ShieldCheck size={20} className="text-orange-600" />
                <h3 className="font-black text-base">{ui.safetyTitle}</h3>
              </div>
              <button
                type="button"
                onClick={() => setShowFaqModal(false)}
                className="p-1 rounded-lg text-zinc-400 hover:text-zinc-600"
              >
                <X size={18} />
              </button>
            </div>

            <div className="space-y-3 text-xs leading-relaxed text-zinc-600 dark:text-zinc-300">
              <div className="p-3 rounded-xl bg-orange-50 dark:bg-orange-950/40 border border-orange-200 dark:border-orange-800 space-y-1">
                <p className="font-black text-orange-700 dark:text-orange-400">
                  1. Advance Wire Scams Warning
                </p>
                <p>
                  Never transfer money in advance via Telebirr or CBE Birr to individual unverified
                  sellers before physically inspecting goods in person.
                </p>
              </div>

              <div className="p-3 rounded-xl bg-zinc-50 dark:bg-zinc-800/60 border border-zinc-200 dark:border-zinc-700 space-y-1">
                <p className="font-black text-zinc-800 dark:text-zinc-200">
                  2. Electronics & Smartphones Protocol
                </p>
                <p>
                  Always dial *#06# to verify the IMEI against the device chassis and box. Verify
                  iCloud or Google FRP removal before paying.
                </p>
              </div>

              <div className="p-3 rounded-xl bg-zinc-50 dark:bg-zinc-800/60 border border-zinc-200 dark:border-zinc-700 space-y-1">
                <p className="font-black text-zinc-800 dark:text-zinc-200">
                  3. Vehicle Libre Ownership Verification
                </p>
                <p>
                  For cars (Toyota, Suzuki, Isuzu), always verify original Libre documents with the
                  Addis Ababa Transport Bureau to verify genuine ownership.
                </p>
              </div>

              <div className="p-3 rounded-xl bg-amber-50 dark:bg-amber-950/30 border border-amber-200 dark:border-amber-800 space-y-1">
                <p className="font-black text-amber-800 dark:text-amber-300">
                  4. HUNTIQ Verified Tiers (Gold & Premium)
                </p>
                <p>
                  Look for 👑 Premium and ⭐ Gold badges. These merchants have submitted verified
                  business licenses and national IDs to HUNTIQ platform administrators.
                </p>
              </div>
            </div>

            <button
              type="button"
              onClick={() => setShowFaqModal(false)}
              className="w-full py-2.5 rounded-xl bg-orange-600 text-white font-bold text-xs"
            >
              {ui.understood}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
