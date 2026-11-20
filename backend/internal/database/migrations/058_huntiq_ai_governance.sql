-- Migration 058: HUNTIQ AI Governance, Conversations, and Audit Infrastructure
-- Creates conversation sessions, message history, and security/audit logging for HUNTIQ AI.

CREATE TABLE IF NOT EXISTS ai_conversations (
    id BIGSERIAL PRIMARY KEY,
    session_id VARCHAR(64) UNIQUE NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL DEFAULT 'New Conversation',
    language_code VARCHAR(10) NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '7 days'),
    is_archived BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_ai_conversations_user ON ai_conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_conversations_session ON ai_conversations(session_id);
CREATE INDEX IF NOT EXISTS idx_ai_conversations_expires ON ai_conversations(expires_at);

CREATE TABLE IF NOT EXISTS ai_messages (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL, -- 'user', 'assistant', 'system', 'tool'
    content TEXT NOT NULL,
    intent VARCHAR(50) NOT NULL DEFAULT 'MARKETPLACE_DISCOVERY',
    tool_calls JSONB,
    tool_results JSONB,
    metadata JSONB,
    token_count INT NOT NULL DEFAULT 0,
    latency_ms INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_messages_conv ON ai_messages(conversation_id, created_at ASC);

CREATE TABLE IF NOT EXISTS ai_audit_events (
    id BIGSERIAL PRIMARY KEY,
    session_id VARCHAR(64),
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    event_type VARCHAR(50) NOT NULL, -- 'QUERY', 'OUT_OF_DOMAIN', 'PROMPT_INJECTION_BLOCKED', 'TOOL_CALLED', 'PROVIDER_ERROR', 'RATE_LIMITED', 'UNSUPPORTED_LANGUAGE'
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_audit_events_type ON ai_audit_events(event_type, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ai_audit_events_created ON ai_audit_events(created_at DESC);

-- Seed default HUNTIQ AI platform settings if not already present
INSERT INTO platform_settings (key, value) VALUES
('ai_enabled', 'true'),
('ai_model', 'gemini-1.5-flash'),
('ai_timeout_seconds', '15'),
('ai_rate_limit_per_min', '30'),
('ai_maintenance_mode', 'false')
ON CONFLICT (key) DO NOTHING;
