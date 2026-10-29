'use client';

import React, { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import {
  Sparkles,
  Send,
  Plus,
  Trash2,
  Edit2,
  Check,
  X,
  Menu,
  ChevronLeft,
  ArrowLeft,
  ExternalLink,
  ShieldCheck,
  Compass,
  Code,
  Search,
  RotateCcw,
  Bot,
  User,
  Copy,
  CheckCheck,
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useLocale } from '@/context/LocaleContext';
import { aiApi, AIConversationItem, AIMessageItem, AIListingCard } from '@/lib/api/ai';
import { publicApi } from '@/lib/api/public';

const SUGGESTED_PROMPTS = [
  {
    icon: Search,
    title: 'Available Listings',
    prompt: 'What products and rentals are available right now?',
  },
  {
    icon: Sparkles,
    title: 'Find Phones in Addis',
    prompt: 'Find Samsung Galaxy or iPhone in Addis Ababa.',
  },
  {
    icon: Compass,
    title: 'Apartments for Rent',
    prompt: 'I need an apartment for rent in Addis Ababa.',
  },
  {
    icon: ShieldCheck,
    title: 'How to Post a Listing',
    prompt: 'How do I post an apartment or product on HUNTIQ?',
  },
];

function AIWorkspaceContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const initialQuery = searchParams.get('q') || '';

  const { user } = useAuth();
  const { locale } = useLocale();

  const [conversations, setConversations] = useState<AIConversationItem[]>([]);
  const [activeConversationId, setActiveConversationId] = useState<string | null>(null);
  const [messages, setMessages] = useState<AIMessageItem[]>([]);
  const [inputMessage, setInputMessage] = useState(initialQuery);
  const [loading, setLoading] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [searchFilter, setSearchFilter] = useState('');
  const [editingTitleId, setEditingTitleId] = useState<string | null>(null);
  const [editedTitle, setEditedTitle] = useState('');
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [copiedId, setCopiedId] = useState<string | null>(null);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Auto-scroll to bottom of messages
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, loading]);

  // Load user conversations if authenticated
  useEffect(() => {
    if (!user?.token) return;
    loadConversations();
  }, [user?.token]);

  const loadConversations = async () => {
    if (!user?.token) return;
    try {
      const list = await aiApi.listConversations(user.token);
      setConversations(list);
      if (list.length > 0 && !activeConversationId) {
        selectConversation(list[0].id);
      }
    } catch (err: any) {
      console.error('Failed to load conversations:', err);
    }
  };

  const selectConversation = async (id: string) => {
    setActiveConversationId(id);
    setErrorMsg(null);
    if (!user?.token) return;
    try {
      const conv = await aiApi.getConversation(user.token, id);
      setMessages(conv.messages || []);
    } catch (err: any) {
      setErrorMsg('Failed to load conversation history.');
    }
  };

  const handleNewChat = async () => {
    setErrorMsg(null);
    if (user?.token) {
      try {
        const newConv = await aiApi.createConversation(user.token, 'New Chat');
        setConversations((prev) => [newConv, ...prev]);
        setActiveConversationId(newConv.id);
        setMessages([]);
      } catch (err) {
        setActiveConversationId(null);
        setMessages([]);
      }
    } else {
      setActiveConversationId(null);
      setMessages([]);
    }
    setInputMessage('');
    if (textareaRef.current) textareaRef.current.focus();
  };

  const handleDeleteConversation = async (e: React.MouseEvent, id: string) => {
    e.stopPropagation();
    if (!user?.token) return;
    try {
      await aiApi.deleteConversation(user.token, id);
      const filtered = conversations.filter((c) => c.id !== id);
      setConversations(filtered);
      if (activeConversationId === id) {
        if (filtered.length > 0) {
          selectConversation(filtered[0].id);
        } else {
          handleNewChat();
        }
      }
    } catch (err) {
      setErrorMsg('Failed to delete conversation.');
    }
  };

  const handleStartRename = (e: React.MouseEvent, conv: AIConversationItem) => {
    e.stopPropagation();
    setEditingTitleId(conv.id);
    setEditedTitle(conv.title);
  };

  const handleSaveRename = async (e: React.MouseEvent, id: string) => {
    e.stopPropagation();
    if (!user?.token || !editedTitle.trim()) {
      setEditingTitleId(null);
      return;
    }
    try {
      await aiApi.updateConversation(user.token, id, editedTitle.trim());
      setConversations((prev) =>
        prev.map((c) => (c.id === id ? { ...c, title: editedTitle.trim() } : c))
      );
      setEditingTitleId(null);
    } catch (err) {
      setErrorMsg('Failed to update title.');
    }
  };

  const handleSendMessage = async (textToSend?: string) => {
    const text = (textToSend || inputMessage).trim();
    if (!text || loading) return;

    setInputMessage('');
    setErrorMsg(null);

    // Optimistically add user message to feed
    const tempUserMsg: AIMessageItem = {
      id: `temp_${Date.now()}`,
      conversation_id: activeConversationId || 'local',
      role: 'user',
      content: text,
      created_at: new Date().toISOString(),
    };
    setMessages((prev) => [...prev, tempUserMsg]);
    setLoading(true);

    try {
      if (!user?.token) {
        // Free guest mode: invoke stateless /ai/chat without requiring sign-in or DB writes
        try {
          const history = messages
            .filter((m) => m.role === 'user' || m.role === 'assistant')
            .slice(-4)
            .map((m) => ({ role: m.role, content: m.content }));

          const res = await aiApi.guestChat(text, history, locale);

          const guestResponse: AIMessageItem = {
            id: `guest_${Date.now()}`,
            conversation_id: 'guest',
            role: 'assistant',
            content: res.answer,
            metadata: {
              model: res.model || 'HUNTIQ Marketplace Assistant',
              latency_ms: res.latency_ms,
              listings: res.listings && res.listings.length > 0 ? res.listings : undefined,
            },
            created_at: new Date().toISOString(),
          };
          setMessages((prev) => [...prev, guestResponse]);
        } catch (err: any) {
          const guestResponse: AIMessageItem = {
            id: `guest_${Date.now()}`,
            conversation_id: 'guest',
            role: 'assistant',
            content: err?.message || 'I am temporarily unable to connect to the marketplace assistant. Please try asking again.',
            metadata: { model: 'HUNTIQ Assistant' },
            created_at: new Date().toISOString(),
          };
          setMessages((prev) => [...prev, guestResponse]);
        } finally {
          setLoading(false);
        }
        return;
      }


      let convId = activeConversationId;
      if (!convId) {
        const newConv = await aiApi.createConversation(user.token, text.slice(0, 35));
        convId = newConv.id;
        setActiveConversationId(newConv.id);
        setConversations((prev) => [newConv, ...prev]);
      }

      const assistantMsg = await aiApi.sendMessage(user.token, convId, text, locale);
      setMessages((prev) => [...prev, assistantMsg]);
      // Refresh conversation title if it was first message
      loadConversations();
    } catch (err: any) {
      setErrorMsg(err.message || 'AI service is temporarily busy. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleCopyCode = (id: string, text: string) => {
    navigator.clipboard.writeText(text);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
  };

  const filteredConversations = conversations.filter((c) =>
    c.title.toLowerCase().includes(searchFilter.toLowerCase())
  );

  return (
    <div className="flex h-screen w-screen overflow-hidden bg-zinc-50 dark:bg-zinc-950 text-zinc-900 dark:text-zinc-100 font-sans">
      {/* ── LEFT AREA: Conversation Sidebar ─────────────────────────────────── */}
      <aside
        className={`${
          sidebarOpen ? 'w-72 md:w-80' : 'w-0'
        } transition-all duration-300 ease-in-out shrink-0 flex flex-col border-r border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 relative z-30 overflow-hidden`}
      >
        {/* Sidebar Header */}
        <div className="p-4 border-b border-zinc-200 dark:border-zinc-800 flex items-center justify-between">
          <Link href="/products" className="flex items-center gap-2 text-zinc-900 dark:text-white font-black group">
            <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-orange-600 via-amber-600 to-orange-500 flex items-center justify-center text-white shadow-md shadow-orange-600/30">
              <Sparkles size={16} className="animate-pulse" />
            </div>
            <div className="flex flex-col">
              <span className="text-sm font-black tracking-tight leading-none">HUNTIQ AI</span>
              <span className="text-[10px] text-zinc-400 font-bold uppercase tracking-wider">Universal Workspace</span>
            </div>
          </Link>
          <button
            onClick={() => setSidebarOpen(false)}
            className="md:hidden p-1.5 rounded-lg hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-500"
          >
            <ChevronLeft size={18} />
          </button>
        </div>

        {/* Action: New Chat */}
        <div className="p-3 border-b border-zinc-100 dark:border-zinc-850">
          <button
            onClick={handleNewChat}
            className="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-700 hover:to-amber-700 text-white font-bold text-xs shadow-sm shadow-orange-600/20 active:scale-98 transition-all"
          >
            <Plus size={16} />
            <span>New Chat</span>
          </button>
        </div>

        {/* History Search Filter */}
        {conversations.length > 3 && (
          <div className="px-3 pt-2">
            <div className="relative">
              <Search size={13} className="absolute left-2.5 top-2.5 text-zinc-400" />
              <input
                type="text"
                value={searchFilter}
                onChange={(e) => setSearchFilter(e.target.value)}
                placeholder="Search chats..."
                className="w-full pl-8 pr-3 py-1.5 text-xs rounded-lg bg-zinc-100 dark:bg-zinc-800/60 border border-transparent focus:border-orange-500 dark:focus:border-orange-500 outline-none transition-colors"
              />
            </div>
          </div>
        )}

        {/* Conversation List */}
        <div className="flex-1 overflow-y-auto px-2 py-3 space-y-1">
          <div className="px-2 pb-1.5 text-[10px] font-black uppercase tracking-wider text-zinc-400 dark:text-zinc-500">
            Recent Conversations
          </div>

          {!user?.token && (
            <div className="p-3 m-1 rounded-xl bg-emerald-50 dark:bg-emerald-950/20 border border-emerald-200 dark:border-emerald-900/40 text-emerald-800 dark:text-emerald-200 text-xs">
              <div className="flex items-center gap-1.5 font-bold">
                <Sparkles size={13} className="text-emerald-600 dark:text-emerald-400" />
                <span>Free Guest Access</span>
              </div>
              <p className="text-[11px] opacity-80 mt-1">
                Ask questions and search verified live marketplace listings freely without an account.
              </p>
            </div>
          )}

          {filteredConversations.map((conv) => {
            const isActive = conv.id === activeConversationId;
            const isEditing = editingTitleId === conv.id;

            return (
              <div
                key={conv.id}
                onClick={() => selectConversation(conv.id)}
                className={`group flex items-center justify-between px-3 py-2.5 rounded-xl text-xs cursor-pointer transition-all ${
                  isActive
                    ? 'bg-orange-50 dark:bg-orange-950/40 text-orange-700 dark:text-orange-300 font-bold border border-orange-200 dark:border-orange-900/50'
                    : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800/50'
                }`}
              >
                {isEditing ? (
                  <div className="flex items-center gap-1 w-full" onClick={(e) => e.stopPropagation()}>
                    <input
                      type="text"
                      value={editedTitle}
                      onChange={(e) => setEditedTitle(e.target.value)}
                      className="flex-1 px-2 py-1 text-xs rounded bg-white dark:bg-zinc-800 border border-orange-500 outline-none"
                      autoFocus
                    />
                    <button
                      onClick={(e) => handleSaveRename(e, conv.id)}
                      className="p-1 text-emerald-600 hover:text-emerald-700"
                    >
                      <Check size={14} />
                    </button>
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        setEditingTitleId(null);
                      }}
                      className="p-1 text-zinc-400 hover:text-zinc-600"
                    >
                      <X size={14} />
                    </button>
                  </div>
                ) : (
                  <>
                    <span className="truncate flex-1 pr-2">{conv.title}</span>
                    <div className="hidden group-hover:flex items-center gap-1 shrink-0 text-zinc-400">
                      <button
                        onClick={(e) => handleStartRename(e, conv)}
                        title="Rename chat"
                        className="p-1 hover:text-zinc-700 dark:hover:text-zinc-200"
                      >
                        <Edit2 size={12} />
                      </button>
                      <button
                        onClick={(e) => handleDeleteConversation(e, conv.id)}
                        title="Delete chat"
                        className="p-1 hover:text-red-600"
                      >
                        <Trash2 size={12} />
                      </button>
                    </div>
                  </>
                )}
              </div>
            );
          })}
        </div>

        {/* Sidebar Footer */}
        <div className="p-3 border-t border-zinc-200 dark:border-zinc-800 flex items-center justify-between text-xs text-zinc-500">
          <Link
            href="/products"
            className="flex items-center gap-1.5 hover:text-orange-600 font-bold transition-colors"
          >
            <ArrowLeft size={14} />
            <span>Marketplace</span>
          </Link>
          <span className="text-[10px] px-2 py-0.5 rounded-full bg-zinc-100 dark:bg-zinc-800 font-mono">v2.0</span>
        </div>
      </aside>

      {/* ── CENTER AREA: Main Conversation ──────────────────────────────────── */}
      <main className="flex-1 flex flex-col h-full overflow-hidden relative">
        {/* Workspace Top Header */}
        <header className="h-14 border-b border-zinc-200 dark:border-zinc-800 bg-white/80 dark:bg-zinc-900/80 backdrop-blur-md px-4 flex items-center justify-between shrink-0 z-20">
          <div className="flex items-center gap-2">
            {!sidebarOpen && (
              <button
                onClick={() => setSidebarOpen(true)}
                title="Expand sidebar"
                className="p-2 rounded-xl border border-zinc-200 dark:border-zinc-800 hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-600 dark:text-zinc-300 transition-colors"
              >
                <Menu size={16} />
              </button>
            )}
            <div className="flex items-center gap-2">
              <span className="font-extrabold text-sm tracking-tight">HUNTIQ AI</span>
              <span className="text-[10px] px-2 py-0.5 rounded-full bg-emerald-100 dark:bg-emerald-950/60 text-emerald-700 dark:text-emerald-300 font-bold border border-emerald-300 dark:border-emerald-800">
                Connected
              </span>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <Link
              href="/products"
              className="text-xs font-bold text-zinc-500 hover:text-zinc-900 dark:hover:text-white flex items-center gap-1 transition-colors"
            >
              <span>Back to Marketplace</span>
              <ExternalLink size={13} />
            </Link>
          </div>
        </header>

        {/* Message Feed Scroll Area */}
        <div className="flex-1 overflow-y-auto px-4 md:px-8 py-6">
          <div className="max-w-3xl mx-auto space-y-6">
            {/* Initial Empty State */}
            {messages.length === 0 && (
              <div className="pt-8 pb-12 text-center animate-fade-in">
                <div className="w-16 h-16 mx-auto mb-4 rounded-3xl bg-gradient-to-tr from-orange-600 via-amber-600 to-orange-500 flex items-center justify-center text-white shadow-xl shadow-orange-600/25">
                  <Sparkles size={32} className="animate-pulse" />
                </div>
                <h1 className="text-2xl md:text-3xl font-black text-zinc-900 dark:text-white tracking-tight">
                  What can I help you with today?
                </h1>
                <p className="mt-2 text-sm text-zinc-500 dark:text-zinc-400 max-w-lg mx-auto leading-relaxed">
                  Ask general knowledge questions, generate code, explore Ethiopian market listings, or compare products across Addis Ababa and beyond.
                </p>

                {/* Suggested Prompt Cards */}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mt-8 text-left">
                  {SUGGESTED_PROMPTS.map((item, idx) => {
                    const Icon = item.icon;
                    return (
                      <button
                        key={idx}
                        onClick={() => handleSendMessage(item.prompt)}
                        className="p-3.5 rounded-2xl border border-zinc-200 dark:border-zinc-800/80 bg-white dark:bg-zinc-900/60 hover:border-orange-500 dark:hover:border-orange-500/80 hover:shadow-md hover:shadow-orange-500/5 transition-all group flex flex-col justify-between"
                      >
                        <div className="flex items-center gap-2 mb-1.5">
                          <div className="p-1.5 rounded-lg bg-orange-100 dark:bg-orange-950/60 text-orange-600 dark:text-orange-400">
                            <Icon size={14} />
                          </div>
                          <span className="text-xs font-bold text-zinc-900 dark:text-zinc-200 group-hover:text-orange-600 transition-colors">
                            {item.title}
                          </span>
                        </div>
                        <p className="text-[11px] text-zinc-500 dark:text-zinc-400 leading-snug line-clamp-2">
                          {item.prompt}
                        </p>
                      </button>
                    );
                  })}
                </div>
              </div>
            )}

            {/* Rendered Messages */}
            {messages.map((msg, index) => {
              const isUser = msg.role === 'user';
              if (msg.role === 'system' || msg.role === 'tool') return null;

              const listings = msg.metadata?.listings as AIListingCard[] | undefined;

              return (
                <div
                  key={msg.id || index}
                  className={`flex gap-3 md:gap-4 ${isUser ? 'justify-end' : 'justify-start'}`}
                >
                  {!isUser && (
                    <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-orange-600 to-amber-600 flex items-center justify-center text-white shrink-0 mt-1 shadow-sm shadow-orange-600/30">
                      <Bot size={16} />
                    </div>
                  )}

                  <div className={`max-w-[85%] md:max-w-[78%] flex flex-col ${isUser ? 'items-end' : 'items-start'}`}>
                    {/* Message Bubble */}
                    <div
                      className={`px-4 py-3 rounded-2xl text-xs sm:text-sm leading-relaxed ${
                        isUser
                          ? 'bg-orange-600 text-white rounded-br-sm shadow-sm shadow-orange-600/20'
                          : 'bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 text-zinc-800 dark:text-zinc-200 rounded-bl-sm shadow-sm'
                      }`}
                    >
                      {/* Code Block / Markdown Render */}
                      <div className="whitespace-pre-wrap font-sans space-y-2">
                        {msg.content}
                      </div>

                      {/* Copy Message Action */}
                      {!isUser && msg.content && (
                        <div className="mt-2 pt-2 border-t border-zinc-100 dark:border-zinc-800 flex items-center justify-between text-[10px] text-zinc-400">
                          <span>{msg.metadata?.model || 'HUNTIQ AI Core'}</span>
                          <button
                            onClick={() => handleCopyCode(msg.id, msg.content)}
                            className="flex items-center gap-1 hover:text-zinc-700 dark:hover:text-zinc-200 transition-colors"
                          >
                            {copiedId === msg.id ? <CheckCheck size={12} className="text-emerald-500" /> : <Copy size={12} />}
                            <span>{copiedId === msg.id ? 'Copied' : 'Copy'}</span>
                          </button>
                        </div>
                      )}
                    </div>

                    {/* Grounded Marketplace Listing Cards */}
                    {listings && listings.length > 0 && (
                      <div className="mt-3 w-full space-y-2">
                        <div className="flex items-center gap-1 text-[11px] font-bold text-zinc-500 dark:text-zinc-400">
                          <ShieldCheck size={13} className="text-emerald-500" />
                          <span>Authoritative Database Listings ({listings.length})</span>
                        </div>
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 w-full">
                          {listings.map((item) => (
                            <Link
                              key={item.id}
                              href={item.detail_url}
                              target="_blank"
                              className="p-2.5 rounded-xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900/80 hover:border-orange-500 transition-all flex gap-2.5 items-center group"
                            >
                              <div className="w-14 h-14 rounded-lg bg-zinc-100 dark:bg-zinc-800 overflow-hidden shrink-0 flex items-center justify-center text-zinc-400">
                                {item.image_url ? (
                                  <img
                                    src={item.image_url}
                                    alt={item.title}
                                    className="w-full h-full object-cover group-hover:scale-105 transition-transform"
                                  />
                                ) : (
                                  <Search size={18} />
                                )}
                              </div>
                              <div className="flex-1 min-w-0">
                                <p className="text-xs font-bold text-zinc-900 dark:text-white truncate group-hover:text-orange-600 transition-colors">
                                  {item.title}
                                </p>
                                <p className="text-xs font-black text-orange-600 mt-0.5">
                                  {item.price.toLocaleString()} {item.currency}
                                </p>
                                <div className="flex items-center gap-1.5 mt-1 text-[10px] text-zinc-400">
                                  <span>{item.city}</span>
                                  {item.seller_tier && item.seller_tier !== 'normal' && (
                                    <span className="px-1.5 py-0.2 rounded bg-amber-100 dark:bg-amber-950 text-amber-700 dark:text-amber-300 font-bold uppercase text-[9px]">
                                      {item.seller_tier}
                                    </span>
                                  )}
                                </div>
                              </div>
                            </Link>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>

                  {isUser && (
                    <div className="w-8 h-8 rounded-xl bg-zinc-200 dark:bg-zinc-800 flex items-center justify-center text-zinc-700 dark:text-zinc-300 shrink-0 mt-1">
                      <User size={16} />
                    </div>
                  )}
                </div>
              );
            })}

            {/* In-Flight Generation Loading State */}
            {loading && (
              <div className="flex gap-3 items-start animate-fade-in">
                <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-orange-600 to-amber-600 flex items-center justify-center text-white shrink-0 mt-1 shadow-sm shadow-orange-600/30">
                  <Bot size={16} />
                </div>
                <div className="px-4 py-3 rounded-2xl rounded-bl-sm bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 text-xs flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-orange-600 animate-ping" />
                  <span className="text-zinc-500 font-medium">HUNTIQ AI is thinking & querying marketplace services...</span>
                </div>
              </div>
            )}

            {/* Error & Retry State */}
            {errorMsg && (
              <div className="p-3 rounded-xl bg-red-50 dark:bg-red-950/20 border border-red-200 dark:border-red-900/40 text-red-700 dark:text-red-300 text-xs flex items-center justify-between">
                <span>{errorMsg}</span>
                <button
                  onClick={() => handleSendMessage()}
                  className="flex items-center gap-1 font-bold underline hover:no-underline ml-2"
                >
                  <RotateCcw size={12} />
                  <span>Retry</span>
                </button>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>
        </div>

        {/* ── BOTTOM AREA: Composer ─────────────────────────────────────────── */}
        <div className="p-4 border-t border-zinc-200 dark:border-zinc-800 bg-white/90 dark:bg-zinc-950/90 backdrop-blur-md">
          <div className="max-w-3xl mx-auto">
            <form
              onSubmit={(e) => {
                e.preventDefault();
                handleSendMessage();
              }}
              className="relative flex items-end gap-2 p-2 rounded-2xl border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-900 focus-within:border-orange-500 dark:focus-within:border-orange-500 shadow-sm transition-all"
            >
              <textarea
                ref={textareaRef}
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    handleSendMessage();
                  }
                }}
                placeholder="Ask HUNTIQ AI anything, or search products across Ethiopia..."
                rows={1}
                className="flex-1 max-h-32 p-2 text-xs sm:text-sm bg-transparent outline-none resize-none text-zinc-900 dark:text-zinc-100 placeholder:text-zinc-400"
              />
              <button
                type="submit"
                disabled={!inputMessage.trim() || loading}
                className="p-2.5 rounded-xl bg-gradient-to-r from-orange-600 to-amber-600 hover:from-orange-700 hover:to-amber-700 text-white disabled:opacity-40 disabled:pointer-events-none transition-all shadow-sm shadow-orange-600/30"
              >
                <Send size={15} />
              </button>
            </form>
            <div className="mt-2 text-center text-[10px] text-zinc-400 dark:text-zinc-500">
              HUNTIQ AI answers general questions and retrieves live marketplace listings. Always verify details before transactions.
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default function AIWorkspacePage() {
  return (
    <React.Suspense
      fallback={
        <div className="flex h-screen w-screen items-center justify-center bg-zinc-50 dark:bg-zinc-950 text-zinc-500 font-bold text-xs">
          Loading HUNTIQ AI...
        </div>
      }
    >
      <AIWorkspaceContent />
    </React.Suspense>
  );
}
