import { fetchAPI, authHeaders } from './core';

export interface AIListingCard {
  id: number;
  title: string;
  price: number;
  currency: string;
  condition: string;
  seller_tier: string;
  city: string;
  subcity?: string;
  detail_url: string;
  image_url?: string;
}

export interface AIMessageItem {
  id: string;
  conversation_id: string;
  role: 'user' | 'assistant' | 'system' | 'tool';
  content: string;
  tool_calls?: any[];
  tool_call_id?: string;
  metadata?: {
    model?: string;
    latency_ms?: number;
    listings?: AIListingCard[];
    [key: string]: any;
  };
  created_at: string;
}

export interface AIConversationItem {
  id: string;
  user_id: number;
  title: string;
  created_at: string;
  updated_at: string;
  messages?: AIMessageItem[];
}

export const aiApi = {
  createConversation: (token: string, title?: string) =>
    fetchAPI<AIConversationItem>('/ai/conversations', {
      method: 'POST',
      body: JSON.stringify({ title: title || 'New Chat' }),
      headers: authHeaders(token),
    }),

  listConversations: (token: string, page: number = 1, limit: number = 30) =>
    fetchAPI<AIConversationItem[]>(`/ai/conversations?page=${page}&limit=${limit}`, {
      headers: authHeaders(token),
    }).then((r) => r ?? []),

  getConversation: (token: string, conversationId: string) =>
    fetchAPI<AIConversationItem>(`/ai/conversations/${conversationId}`, {
      headers: authHeaders(token),
    }),

  updateConversation: (token: string, conversationId: string, title: string) =>
    fetchAPI<{ updated: boolean; title: string }>(`/ai/conversations/${conversationId}`, {
      method: 'PATCH',
      body: JSON.stringify({ title }),
      headers: authHeaders(token),
    }),

  deleteConversation: (token: string, conversationId: string) =>
    fetchAPI<{ deleted: boolean }>(`/ai/conversations/${conversationId}`, {
      method: 'DELETE',
      headers: authHeaders(token),
    }),

  sendMessage: (token: string, conversationId: string, content: string, locale: string = 'en') =>
    fetchAPI<AIMessageItem>(`/ai/conversations/${conversationId}/messages`, {
      method: 'POST',
      body: JSON.stringify({ content, locale }),
      headers: authHeaders(token),
    }),

  guestChat: (message: string, history: { role: string; content: string }[] = [], locale: string = 'en') =>
    fetchAPI<{ answer: string; listings?: AIListingCard[]; latency_ms?: number; model?: string }>('/ai/chat', {
      method: 'POST',
      body: JSON.stringify({ message, history, locale }),
    }),
};

