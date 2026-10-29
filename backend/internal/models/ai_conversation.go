package models

import (
	"time"
)

// AIConversation represents an ongoing multi-turn conversational intelligence session.
type AIConversation struct {
	ID        string      `json:"id"`
	UserID    int         `json:"user_id"`
	Title     string      `json:"title"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Messages  []AIMessage `json:"messages,omitempty"`
}

// AIToolCall represents a tool invocation request emitted by the language model.
type AIToolCall struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Function AIFunctionInvocation   `json:"function"`
}

// AIFunctionInvocation contains the tool name and parsed JSON arguments.
type AIFunctionInvocation struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// AIMessage represents a single dialogue turn in an AI conversation.
type AIMessage struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversation_id"`
	Role           string                 `json:"role"` // 'user', 'assistant', 'system', 'tool'
	Content        string                 `json:"content"`
	ToolCalls      []AIToolCall           `json:"tool_calls,omitempty"`
	ToolCallID     *string                `json:"tool_call_id,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}
