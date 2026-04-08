package events

import (
	"encoding/json"
	"strings"
	"time"
)

// Event is the top-level envelope for all JSONL events.
type Event struct {
	Type                  string          `json:"type"`
	ID                    string          `json:"id,omitempty"`
	CWD                   string          `json:"cwd,omitempty"`
	Timestamp             any             `json:"timestamp,omitempty"`
	Message               *Message        `json:"message,omitempty"`
	AssistantMessageEvent *AssistantEvent `json:"assistantMessageEvent,omitempty"`
	ToolCallID            string          `json:"toolCallId,omitempty"`
	ToolName              string          `json:"toolName,omitempty"`
	Args                  json.RawMessage `json:"args,omitempty"`
	Result                *ToolResult     `json:"result,omitempty"`
	IsError               bool            `json:"isError,omitempty"`
	PartialResult         *ToolResult     `json:"partialResult,omitempty"`
}

// Message represents a user or assistant message.
type Message struct {
	Role      string         `json:"role"`
	Content   []ContentBlock `json:"content"`
	Timestamp any            `json:"timestamp,omitempty"`
}

// ContentBlock represents a single content element (text, thinking, tool call, etc).
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// AssistantEvent represents incremental updates to assistant messages.
type AssistantEvent struct {
	Type         string    `json:"type"`
	Delta        string    `json:"delta,omitempty"`
	Content      string    `json:"content,omitempty"`
	ToolCall     *ToolCall `json:"toolCall,omitempty"`
	ContentIndex int       `json:"contentIndex,omitempty"`
}

// ToolCall represents a tool invocation.
type ToolCall struct {
	Type      string          `json:"type"`
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ToolResult represents the output of a tool execution.
type ToolResult struct {
	Content []ContentBlock `json:"content"`
}

// ParseTimestamp converts event timestamp (ISO 8601 or milliseconds epoch) to time.Time.
func ParseTimestamp(ts any) time.Time {
	if ts == nil {
		return time.Now()
	}

	switch v := ts.(type) {
	case string:
		// ISO 8601 format
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			return t
		}
	case float64:
		// milliseconds epoch
		return time.UnixMilli(int64(v))
	}

	return time.Now()
}

// ExtractText concatenates all text content blocks into a single string.
func ExtractText(content []ContentBlock) string {
	var result strings.Builder
	for _, block := range content {
		if block.Type == "text" && block.Text != "" {
			if result.Len() > 0 {
				result.WriteString("\n")
			}
			result.WriteString(block.Text)
		}
	}
	return result.String()
}
