package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/andreaswachs/forhumans/internal/events"
)

// Options controls formatter behavior.
type Options struct {
	// ShowLifecycle controls whether to output AGENT START/END and TURN START/END events.
	ShowLifecycle bool
}

// Formatter converts events into human-readable log output.
type Formatter struct {
	opts Options
}

// New creates a new event formatter with default options.
func New() *Formatter {
	return &Formatter{
		opts: Options{
			ShowLifecycle: false,
		},
	}
}

// NewWithOptions creates a new event formatter with custom options.
func NewWithOptions(opts Options) *Formatter {
	return &Formatter{opts: opts}
}

// Format processes an event and returns the formatted log line, or empty string if skipped.
func (f *Formatter) Format(event events.Event) string {
	switch event.Type {
	case "session":
		ts := formatTime(events.ParseTimestamp(event.Timestamp))
		return fmt.Sprintf("%s [SESSION] id=%s cwd=%s", ts, event.ID, event.CWD)

	case "agent_start":
		if !f.opts.ShowLifecycle {
			return ""
		}
		ts := formatTime(time.Now())
		return fmt.Sprintf("%s [AGENT START]", ts)

	case "agent_end":
		if !f.opts.ShowLifecycle {
			return ""
		}
		ts := formatTime(time.Now())
		return fmt.Sprintf("%s [AGENT END]", ts)

	case "turn_start":
		if !f.opts.ShowLifecycle {
			return ""
		}
		ts := formatTime(time.Now())
		return fmt.Sprintf("%s [TURN START]", ts)

	case "turn_end":
		if !f.opts.ShowLifecycle {
			return ""
		}
		ts := formatTime(time.Now())
		return fmt.Sprintf("%s [TURN END]", ts)

	case "message_start":
		if event.Message != nil && event.Message.Role == "user" {
			text := events.ExtractText(event.Message.Content)
			ts := formatTime(events.ParseTimestamp(event.Message.Timestamp))
			return fmt.Sprintf("%s [USER] %s", ts, text)
		}
		return ""

	case "message_update":
		if event.AssistantMessageEvent == nil {
			return ""
		}

		switch event.AssistantMessageEvent.Type {
		case "thinking_end":
			ts := formatTime(time.Now())
			content := truncate(event.AssistantMessageEvent.Content, 200)
			return fmt.Sprintf("%s [THINKING] %s", ts, content)

		case "text_end":
			ts := formatTime(time.Now())
			return formatWithIndent(ts, "[RESPONSE]", event.AssistantMessageEvent.Content)

		case "toolcall_end":
			if event.AssistantMessageEvent.ToolCall != nil {
				ts := formatTime(time.Now())
				tc := event.AssistantMessageEvent.ToolCall
				return fmt.Sprintf("%s [TOOL CALL] %s(%s)", ts, tc.Name, string(tc.Arguments))
			}
			return ""
		}
		return ""

	case "tool_execution_start":
		ts := formatTime(time.Now())
		argsStr := string(event.Args)
		return fmt.Sprintf("%s [TOOL START] %s args=%s", ts, event.ToolName, argsStr)

	case "tool_execution_end":
		ts := formatTime(time.Now())
		resultText := ""
		if event.Result != nil {
			resultText = events.ExtractText(event.Result.Content)
		}
		resultText = truncate(resultText, 500)
		return fmt.Sprintf("%s [TOOL END] %s error=%v result=%s", ts, event.ToolName, event.IsError, resultText)

	case "message_end":
		// message_end just marks completion of a message already output by message_start/message_update
		return ""

	case "tool_execution_update":
		// tool_execution_update is intermediate progress, skip to avoid noise
		// (final result will come from tool_execution_end)
		return ""

	default:
		// Handle unknown event types by printing raw JSON
		ts := formatTime(time.Now())
		rawJSON, err := json.Marshal(event)
		if err != nil {
			return fmt.Sprintf("%s [UNKNOWN] failed to marshal event", ts)
		}
		return fmt.Sprintf("%s [UNKNOWN] %s", ts, string(rawJSON))
	}
}

func formatTime(ts time.Time) string {
	return ts.Format("15:04:05")
}

// formatWithIndent formats a timestamp + tag + content, indenting multiline content
// to align with where the content starts on the first line.
func formatWithIndent(ts, tag, content string) string {
	prefix := fmt.Sprintf("%s %s ", ts, tag)
	lines := strings.Split(content, "\n")

	// First line doesn't need indent
	var result strings.Builder
	result.WriteString(prefix)
	result.WriteString(lines[0])

	// Subsequent lines get indented to align with first line content
	if len(lines) > 1 {
		indent := strings.Repeat(" ", len(prefix))
		for i := 1; i < len(lines); i++ {
			result.WriteString("\n")
			result.WriteString(indent)
			result.WriteString(lines[i])
		}
	}

	return result.String()
}

// truncate shortens a string to maxLen characters, replacing newlines with spaces.
func truncate(s string, maxLen int) string {
	// Replace newlines with spaces for single-line output
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
