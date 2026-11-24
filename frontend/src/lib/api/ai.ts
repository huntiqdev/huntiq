import { fetchAPI } from './core';
import type { AIChatRequest, AIChatResponse } from '@/types/ai';

export async function sendAIChat(req: AIChatRequest): Promise<AIChatResponse> {
  return fetchAPI<AIChatResponse>('/ai/chat', {
    method: 'POST',
    body: JSON.stringify(req),
  });
}

export async function getAIHealth(): Promise<{ status: string; service: string }> {
  return fetchAPI<{ status: string; service: string }>('/ai/health', {
    method: 'GET',
  });
}

export interface TranslationResult {
  source_language: string;
  translations: {
    en: string;
    am: string;
    or: string;
    ti: string;
  };
  fallback: boolean;
}

export async function translateListingText(text: string, fromLang?: string): Promise<TranslationResult> {
  return fetchAPI<TranslationResult>('/ai/translate', {
    method: 'POST',
    body: JSON.stringify({ text, from_lang: fromLang }),
  });
}
