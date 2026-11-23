export interface ListingCard {
  id: number;
  title: string;
  price: number;
  currency: string;
  category_name: string;
  location_name: string;
  condition: string;
  listing_type: string;
  image_url?: string;
  view_count?: number;
}

export interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: number;
  path?: string;
  intent?: string;
  isOutOfDomain?: boolean;
  listings?: ListingCard[];
  suggestedQueries?: string[];
  latencyMs?: number;
  isError?: boolean;
  errorCode?: 'NO_INTERNET' | 'BACKEND_UNREACHABLE' | 'AI_UNAVAILABLE' | 'TIMEOUT' | 'RATE_LIMITED' | 'GENERIC_ERROR';
  retryText?: string;
}

export interface AIChatRequest {
  message: string;
  language?: string;
  history?: Array<{
    role: string;
    parts?: Array<{ text: string }>;
  }>;
}

export interface AIChatResponse {
  reply: string;
  language: string;
  intent: string;
  path: string;
  is_out_of_domain: boolean;
  fallback_used?: boolean;
  tokens_used: number;
  latency_ms: number;
  suggested_queries?: string[];
  listings?: ListingCard[];
  error?: string;
  error_code?: string;
}
