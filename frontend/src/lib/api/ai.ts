import { fetchAPI } from './core';

export interface ListingCardDTO {
  id: number;
  title: string;
  price: number;
  currency: string;
  category_name?: string;
  location_name?: string;
  primary_image?: string;
  image_urls?: string[];
  url: string;
  condition?: string;
  listing_type?: string;
  created_at: string;
}

export interface ChatRequest {
  session_id: string;
  message: string;
  language?: string;
  category?: string;
  location?: string;
}

export interface ChatResponse {
  session_id: string;
  language: string;
  reply: string;
  intent: string;
  listings?: ListingCardDTO[];
  suggested_queries?: string[];
  is_out_of_domain: boolean;
  fallback_used: boolean;
  tokens_used?: number;
  latency_ms: number;
}

export interface AIHealthStatus {
  status: string;
  provider: string;
  model: string;
  enabled: boolean;
  latency_ms: number;
  last_checked: string;
  supported_languages: string[];
}

export interface AIMetricsReport {
  total_queries: number;
  out_of_domain_queries: number;
  adversarial_blocks: number;
  provider_errors: number;
  total_tokens: number;
  estimated_cost_usd: number;
  average_latency_ms: number;
  active_sessions_24h: number;
}

export const aiApi = {
  chatWithAI: async (req: ChatRequest): Promise<ChatResponse> => {
    return fetchAPI<ChatResponse>('/ai/chat', {
      method: 'POST',
      body: JSON.stringify(req),
    });
  },

  getAIHealth: async (): Promise<AIHealthStatus> => {
    return fetchAPI<AIHealthStatus>('/ai/health', { method: 'GET' });
  },

  getAISessionHistory: async (sessionId: string): Promise<{ conversation: any; messages: any[] }> => {
    return fetchAPI<{ conversation: any; messages: any[] }>(`/ai/session/${sessionId}`, { method: 'GET' });
  },

  clearAISession: async (sessionId: string): Promise<{ success: boolean; message: string }> => {
    return fetchAPI<{ success: boolean; message: string }>(`/ai/session/${sessionId}`, { method: 'DELETE' });
  },

  getAIAdminMetrics: async (): Promise<AIMetricsReport> => {
    return fetchAPI<AIMetricsReport>('/admin/ai/metrics', { method: 'GET' });
  },
};
