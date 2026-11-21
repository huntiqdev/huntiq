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
