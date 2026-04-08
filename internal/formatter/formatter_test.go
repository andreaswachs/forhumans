package formatter

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/andreaswachs/forhumans/internal/events"
)

func TestSessionEvent(t *testing.T) {
	f := New()
	event := events.Event{
		Type:      "session",
		ID:        "680697ea-1fbf-4f19-8fea-5570a158d944",
		CWD:       "/Users/awa/Source/forhumans",
		Timestamp: "2026-04-08T09:06:19.446Z",
	}

	output := f.Format(event)
	if !strings.Contains(output, "[SESSION]") {
		t.Errorf("expected [SESSION] in output, got: %s", output)
	}
	if !strings.Contains(output, "680697ea-1fbf-4f19-8fea-5570a158d944") {
		t.Errorf("expected session ID in output, got: %s", output)
	}
	if !strings.Contains(output, "/Users/awa/Source/forhumans") {
		t.Errorf("expected CWD in output, got: %s", output)
	}
}

func TestAgentStartEvent(t *testing.T) {
	f := NewWithOptions(Options{ShowLifecycle: true})
	event := events.Event{Type: "agent_start"}
	output := f.Format(event)

	if !strings.Contains(output, "[AGENT START]") {
		t.Errorf("expected [AGENT START] in output, got: %s", output)
	}
}

func TestAgentStartEventHidden(t *testing.T) {
	f := New() // ShowLifecycle disabled by default
	event := events.Event{Type: "agent_start"}
	output := f.Format(event)

	if output != "" {
		t.Errorf("expected empty output when ShowLifecycle disabled, got: %s", output)
	}
}

func TestAgentEndEvent(t *testing.T) {
	f := NewWithOptions(Options{ShowLifecycle: true})
	event := events.Event{Type: "agent_end"}
	output := f.Format(event)

	if !strings.Contains(output, "[AGENT END]") {
		t.Errorf("expected [AGENT END] in output, got: %s", output)
	}
}

func TestAgentEndEventHidden(t *testing.T) {
	f := New() // ShowLifecycle disabled by default
	event := events.Event{Type: "agent_end"}
	output := f.Format(event)

	if output != "" {
		t.Errorf("expected empty output when ShowLifecycle disabled, got: %s", output)
	}
}

func TestTurnStartEvent(t *testing.T) {
	f := NewWithOptions(Options{ShowLifecycle: true})
	event := events.Event{Type: "turn_start"}
	output := f.Format(event)

	if !strings.Contains(output, "[TURN START]") {
		t.Errorf("expected [TURN START] in output, got: %s", output)
	}
}

func TestTurnStartEventHidden(t *testing.T) {
	f := New() // ShowLifecycle disabled by default
	event := events.Event{Type: "turn_start"}
	output := f.Format(event)

	if output != "" {
		t.Errorf("expected empty output when ShowLifecycle disabled, got: %s", output)
	}
}

func TestTurnEndEvent(t *testing.T) {
	f := NewWithOptions(Options{ShowLifecycle: true})
	event := events.Event{Type: "turn_end"}
	output := f.Format(event)

	if !strings.Contains(output, "[TURN END]") {
		t.Errorf("expected [TURN END] in output, got: %s", output)
	}
}

func TestTurnEndEventHidden(t *testing.T) {
	f := New() // ShowLifecycle disabled by default
	event := events.Event{Type: "turn_end"}
	output := f.Format(event)

	if output != "" {
		t.Errorf("expected empty output when ShowLifecycle disabled, got: %s", output)
	}
}

func TestUserMessage(t *testing.T) {
	f := New()
	event := events.Event{
		Type: "message_start",
		Message: &events.Message{
			Role: "user",
			Content: []events.ContentBlock{
				{Type: "text", Text: "run"},
			},
			Timestamp: float64(1775639179472), // milliseconds epoch
		},
	}

	output := f.Format(event)
	if !strings.Contains(output, "[USER]") {
		t.Errorf("expected [USER] in output, got: %s", output)
	}
	if !strings.Contains(output, "run") {
		t.Errorf("expected 'run' in output, got: %s", output)
	}
}

func TestUserMessageIgnoresNonUserRole(t *testing.T) {
	f := New()
	event := events.Event{
		Type: "message_start",
		Message: &events.Message{
			Role: "assistant",
			Content: []events.ContentBlock{
				{Type: "text", Text: "response"},
			},
		},
	}

	output := f.Format(event)
	if output != "" {
		t.Errorf("expected empty output for non-user message, got: %s", output)
	}
}

func TestThinkingEndEvent(t *testing.T) {
	f := New()
	thinkingText := "User says \"run\". Probably want to list files? Let's run bash ls."
	event := events.Event{
		Type: "message_update",
		AssistantMessageEvent: &events.AssistantEvent{
			Type:    "thinking_end",
			Content: thinkingText,
		},
	}

	output := f.Format(event)
	if !strings.Contains(output, "[THINKING]") {
		t.Errorf("expected [THINKING] in output, got: %s", output)
	}
	if !strings.Contains(output, "User says") {
		t.Errorf("expected thinking content in output, got: %s", output)
	}
}

func TestThinkingEndTruncation(t *testing.T) {
	f := New()
	longThinking := strings.Repeat("x", 300)
	event := events.Event{
		Type: "message_update",
		AssistantMessageEvent: &events.AssistantEvent{
			Type:    "thinking_end",
			Content: longThinking,
		},
	}

	output := f.Format(event)
	if !strings.Contains(output, "...") {
		t.Errorf("expected truncation marker (...) in output for long thinking, got: %s", output)
	}
	if len(output) > 250 { // rough check that it's truncated
		t.Errorf("expected truncated thinking content, got length: %d", len(output))
	}
}

func TestResponseEvent(t *testing.T) {
	f := New()
	responseText := "The repository has been listed."
	event := events.Event{
		Type: "message_update",
		AssistantMessageEvent: &events.AssistantEvent{
			Type:    "text_end",
			Content: responseText,
		},
	}

	output := f.Format(event)
	if !strings.Contains(output, "[RESPONSE]") {
		t.Errorf("expected [RESPONSE] in output, got: %s", output)
	}
	if !strings.Contains(output, responseText) {
		t.Errorf("expected response text in output, got: %s", output)
	}
}

func TestToolCallEvent(t *testing.T) {
	f := New()
	argsJSON := json.RawMessage(`{"command":"ls -R ."}`)
	event := events.Event{
		Type: "message_update",
		AssistantMessageEvent: &events.AssistantEvent{
			Type: "toolcall_end",
			ToolCall: &events.ToolCall{
				Name:      "bash",
				Arguments: argsJSON,
			},
		},
	}

	output := f.Format(event)
	if !strings.Contains(output, "[TOOL CALL]") {
		t.Errorf("expected [TOOL CALL] in output, got: %s", output)
	}
	if !strings.Contains(output, "bash") {
		t.Errorf("expected tool name 'bash' in output, got: %s", output)
	}
	if !strings.Contains(output, "ls -R .") {
		t.Errorf("expected tool arguments in output, got: %s", output)
	}
}

func TestToolExecutionStart(t *testing.T) {
	f := New()
	argsJSON := json.RawMessage(`{"command":"ls -R ."}`)
	event := events.Event{
		Type:     "tool_execution_start",
		ToolName: "bash",
		Args:     argsJSON,
	}

	output := f.Format(event)
	if !strings.Contains(output, "[TOOL START]") {
		t.Errorf("expected [TOOL START] in output, got: %s", output)
	}
	if !strings.Contains(output, "bash") {
		t.Errorf("expected tool name in output, got: %s", output)
	}
	if !strings.Contains(output, "ls -R .") {
		t.Errorf("expected args in output, got: %s", output)
	}
}

func TestToolExecutionEnd(t *testing.T) {
	f := New()
	event := events.Event{
		Type:     "tool_execution_end",
		ToolName: "bash",
		IsError:  false,
		Result: &events.ToolResult{
			Content: []events.ContentBlock{
				{Type: "text", Text: "output.jsonl\nrandom_numbers.json"},
			},
		},
	}

	output := f.Format(event)
	if !strings.Contains(output, "[TOOL END]") {
		t.Errorf("expected [TOOL END] in output, got: %s", output)
	}
	if !strings.Contains(output, "bash") {
		t.Errorf("expected tool name in output, got: %s", output)
	}
	if !strings.Contains(output, "error=false") {
		t.Errorf("expected error=false in output, got: %s", output)
	}
	if !strings.Contains(output, "output.jsonl") {
		t.Errorf("expected tool result in output, got: %s", output)
	}
}

func TestToolExecutionEndWithError(t *testing.T) {
	f := New()
	event := events.Event{
		Type:     "tool_execution_end",
		ToolName: "bash",
		IsError:  true,
		Result: &events.ToolResult{
			Content: []events.ContentBlock{
				{Type: "text", Text: "command not found"},
			},
		},
	}

	output := f.Format(event)
	if !strings.Contains(output, "error=true") {
		t.Errorf("expected error=true in output, got: %s", output)
	}
}

func TestMessageUpdateWithNoEvent(t *testing.T) {
	f := New()
	event := events.Event{
		Type: "message_update",
		// AssistantMessageEvent is nil
	}

	output := f.Format(event)
	if output != "" {
		t.Errorf("expected empty output for message_update with no event, got: %s", output)
	}
}

func TestMessageUpdateWithUnknownSubType(t *testing.T) {
	f := New()
	event := events.Event{
		Type: "message_update",
		AssistantMessageEvent: &events.AssistantEvent{
			Type: "thinking_delta", // delta events are skipped
		},
	}

	output := f.Format(event)
	if output != "" {
		t.Errorf("expected empty output for thinking_delta, got: %s", output)
	}
}

func TestToolCallEventWithoutToolCall(t *testing.T) {
	f := New()
	event := events.Event{
		Type: "message_update",
		AssistantMessageEvent: &events.AssistantEvent{
			Type: "toolcall_end",
			// ToolCall is nil
		},
	}

	output := f.Format(event)
	if output != "" {
		t.Errorf("expected empty output for toolcall_end without ToolCall, got: %s", output)
	}
}

func TestUnknownEventType(t *testing.T) {
	f := New()
	event := events.Event{
		Type: "unknown_type",
		ID:   "test-id",
	}

	output := f.Format(event)
	if !strings.Contains(output, "[UNKNOWN]") {
		t.Errorf("expected [UNKNOWN] in output, got: %s", output)
	}
	if !strings.Contains(output, "unknown_type") {
		t.Errorf("expected unknown event type in raw JSON, got: %s", output)
	}
}

func TestComplexUserMessage(t *testing.T) {
	f := New()
	event := events.Event{
		Type: "message_start",
		Message: &events.Message{
			Role: "user",
			Content: []events.ContentBlock{
				{Type: "text", Text: "Please create a file"},
				{Type: "text", Text: "with some content"},
			},
			Timestamp: int64(1775639179472),
		},
	}

	output := f.Format(event)
	if !strings.Contains(output, "[USER]") {
		t.Errorf("expected [USER] in output, got: %s", output)
	}
	if !strings.Contains(output, "Please create a file") {
		t.Errorf("expected first text block in output, got: %s", output)
	}
	if !strings.Contains(output, "with some content") {
		t.Errorf("expected second text block in output, got: %s", output)
	}
}

func TestResponseNewlines(t *testing.T) {
	f := New()
	responseText := "Line 1\nLine 2\nLine 3"
	event := events.Event{
		Type: "message_update",
		AssistantMessageEvent: &events.AssistantEvent{
			Type:    "text_end",
			Content: responseText,
		},
	}

	output := f.Format(event)
	// Output should preserve newlines in response
	if !strings.Contains(output, "Line 1") {
		t.Errorf("expected multiline response content, got: %s", output)
	}
}

func TestTruncateWithNewlines(t *testing.T) {
	longThinking := "Line 1\nLine 2\n" + strings.Repeat("x", 300)
	f := New()
	event := events.Event{
		Type: "message_update",
		AssistantMessageEvent: &events.AssistantEvent{
			Type:    "thinking_end",
			Content: longThinking,
		},
	}

	output := f.Format(event)
	// Should replace newlines with spaces for single-line output
	if strings.Count(output, "\n") > 0 {
		t.Errorf("expected newlines to be replaced in truncated output, got: %s", output)
	}
	if !strings.Contains(output, "...") {
		t.Errorf("expected truncation marker, got: %s", output)
	}
}

func TestToolResultTruncation(t *testing.T) {
	f := New()
	longResult := strings.Repeat("output line\n", 100)
	event := events.Event{
		Type:     "tool_execution_end",
		ToolName: "bash",
		IsError:  false,
		Result: &events.ToolResult{
			Content: []events.ContentBlock{
				{Type: "text", Text: longResult},
			},
		},
	}

	output := f.Format(event)
	if !strings.Contains(output, "...") {
		t.Errorf("expected truncation marker in long result, got: %s", output)
	}
	// Verify it's actually truncated (500 char limit + some overhead)
	if len(output) > 700 {
		t.Errorf("expected truncated result, output length: %d", len(output))
	}
}

func TestTimestampFormatting(t *testing.T) {
	f := NewWithOptions(Options{ShowLifecycle: true})
	event := events.Event{
		Type: "agent_start",
	}

	output := f.Format(event)
	// Check for HH:MM:SS format
	parts := strings.Fields(output)
	if len(parts) < 1 {
		t.Errorf("expected timestamp in output, got: %s", output)
	}

	timeStr := parts[0]
	timeParts := strings.Split(timeStr, ":")
	if len(timeParts) != 3 {
		t.Errorf("expected HH:MM:SS format, got: %s", timeStr)
	}
}

func TestParseTimestampISO8601(t *testing.T) {
	// This tests the ParseTimestamp function indirectly through session event
	f := New()
	event := events.Event{
		Type:      "session",
		ID:        "test-id",
		CWD:       "/tmp",
		Timestamp: "2026-04-08T09:06:19.446Z",
	}

	output := f.Format(event)
	if !strings.Contains(output, "09:06:19") {
		t.Errorf("expected parsed timestamp 09:06:19 in output, got: %s", output)
	}
}

func TestParseTimestampMilliseconds(t *testing.T) {
	// Timestamp: 1775639179472 milliseconds = April 8, 2026, ~11:06:19
	f := New()
	event := events.Event{
		Type: "message_start",
		Message: &events.Message{
			Role:      "user",
			Timestamp: float64(1775639179472),
			Content: []events.ContentBlock{
				{Type: "text", Text: "test"},
			},
		},
	}

	output := f.Format(event)
	// Just verify we get a valid timestamp in output
	parts := strings.Fields(output)
	if len(parts) < 1 {
		t.Errorf("expected timestamp in output, got: %s", output)
	}
}

func TestMessageEndEvent(t *testing.T) {
	f := New()
	event := events.Event{
		Type: "message_end",
		Message: &events.Message{
			Role: "user",
			Content: []events.ContentBlock{
				{Type: "text", Text: "test"},
			},
		},
	}

	output := f.Format(event)
	if output != "" {
		t.Errorf("expected empty output for message_end, got: %s", output)
	}
}

func TestMessageEndUserEvent(t *testing.T) {
	f := New()
	event := events.Event{
		Type: "message_end",
		Message: &events.Message{
			Role: "user",
			Content: []events.ContentBlock{
				{Type: "text", Text: "run"},
			},
			Timestamp: float64(1775639179472),
		},
	}

	output := f.Format(event)
	if output != "" {
		t.Errorf("expected empty output for message_end (user), got: %s", output)
	}
}

func TestMessageEndAssistantEvent(t *testing.T) {
	f := New()
	event := events.Event{
		Type: "message_end",
		Message: &events.Message{
			Role: "assistant",
			Content: []events.ContentBlock{
				{Type: "text", Text: "response"},
			},
		},
	}

	output := f.Format(event)
	if output != "" {
		t.Errorf("expected empty output for message_end (assistant), got: %s", output)
	}
}

func TestToolExecutionUpdate(t *testing.T) {
	f := New()
	event := events.Event{
		Type:       "tool_execution_update",
		ToolName:   "bash",
		ToolCallID: "call-123",
		Args:       json.RawMessage(`{"command":"ls"}`),
		PartialResult: &events.ToolResult{
			Content: []events.ContentBlock{
				{Type: "text", Text: "partial output"},
			},
		},
	}

	output := f.Format(event)
	if output != "" {
		t.Errorf("expected empty output for tool_execution_update (intermediate), got: %s", output)
	}
}

func TestToolExecutionUpdateEmptyResult(t *testing.T) {
	f := New()
	event := events.Event{
		Type:     "tool_execution_update",
		ToolName: "bash",
		Args:     json.RawMessage(`{"command":"ls"}`),
		PartialResult: &events.ToolResult{
			Content: []events.ContentBlock{},
		},
	}

	output := f.Format(event)
	if output != "" {
		t.Errorf("expected empty output for tool_execution_update, got: %s", output)
	}
}

func TestResponseMultilineIndent(t *testing.T) {
	f := New()
	multilineResponse := "First line\nSecond line\nThird line"
	event := events.Event{
		Type: "message_update",
		AssistantMessageEvent: &events.AssistantEvent{
			Type:    "text_end",
			Content: multilineResponse,
		},
	}

	output := f.Format(event)
	lines := strings.Split(output, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines in output, got %d: %s", len(lines), output)
	}

	// First line should start with timestamp and [RESPONSE]
	if !strings.Contains(lines[0], "[RESPONSE]") {
		t.Errorf("expected [RESPONSE] in first line, got: %s", lines[0])
	}

	// Second and third lines should be indented to align with content start
	expectedIndent := strings.Index(lines[0], "First") // position where content starts
	if expectedIndent <= 0 {
		t.Errorf("couldn't find content start position in first line: %s", lines[0])
	}

	for i := 1; i < len(lines); i++ {
		actualIndent := len(lines[i]) - len(strings.TrimLeft(lines[i], " "))
		if actualIndent != expectedIndent {
			t.Errorf("line %d has wrong indentation: expected %d spaces, got %d. Line: %q", i+1, expectedIndent, actualIndent, lines[i])
		}
	}
}

func TestResponseSinglelineNoExtraIndent(t *testing.T) {
	f := New()
	singlelineResponse := "Single line response"
	event := events.Event{
		Type: "message_update",
		AssistantMessageEvent: &events.AssistantEvent{
			Type:    "text_end",
			Content: singlelineResponse,
		},
	}

	output := f.Format(event)
	if strings.Contains(output, "\n") {
		t.Errorf("expected single line output, got: %s", output)
	}
	if !strings.Contains(output, singlelineResponse) {
		t.Errorf("expected content in output, got: %s", output)
	}
}
