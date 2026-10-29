-- Migration: 036_ai_conversations.sql
-- Description: Additive schema for full-page multi-turn conversational AI sessions and messages.
-- Safe, additive migration with foreign keys, ownership indexes, and JSONB metadata for tool tracking.

CREATE TABLE IF NOT EXISTS ai_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL DEFAULT 'New Chat',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_conversations_user_updated
    ON ai_conversations (user_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS ai_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
    role VARCHAR(32) NOT NULL, -- 'user', 'assistant', 'system', 'tool'
    content TEXT NOT NULL DEFAULT '',
    tool_calls JSONB DEFAULT NULL,
    tool_call_id VARCHAR(128) DEFAULT NULL,
    metadata JSONB DEFAULT NULL, -- model, token counts, latency_ms, inline product cards
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_messages_conversation_created
    ON ai_messages (conversation_id, created_at ASC);

-- Rollback Guidance:
-- DROP TABLE IF EXISTS ai_messages;
-- DROP TABLE IF EXISTS ai_conversations;
