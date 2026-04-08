package events

import (
	"strings"
	"testing"
	"time"
)

func TestParseTimestampISO8601(t *testing.T) {
	ts := ParseTimestamp("2026-04-08T09:06:19.446Z")
	if ts.IsZero() {
		t.Errorf("expected valid timestamp, got zero time")
	}

	// Verify year is correct
	if ts.Year() != 2026 {
		t.Errorf("expected year 2026, got %d", ts.Year())
	}
	if ts.Month() != time.April {
		t.Errorf("expected month April, got %v", ts.Month())
	}
	if ts.Day() != 8 {
		t.Errorf("expected day 8, got %d", ts.Day())
	}
	if ts.Hour() != 9 {
		t.Errorf("expected hour 9, got %d", ts.Hour())
	}
}

func TestParseTimestampMilliseconds(t *testing.T) {
	// 1775639179472 milliseconds = April 8, 2026, ~11:06:19 UTC
	ts := ParseTimestamp(float64(1775639179472))
	if ts.IsZero() {
		t.Errorf("expected valid timestamp, got zero time")
	}

	// Verify it's within expected range
	if ts.Year() != 2026 {
		t.Errorf("expected year 2026, got %d", ts.Year())
	}
	if ts.Month() != time.April {
		t.Errorf("expected month April, got %v", ts.Month())
	}
	if ts.Day() != 8 {
		t.Errorf("expected day 8, got %d", ts.Day())
	}
}

func TestParseTimestampInt64Milliseconds(t *testing.T) {
	// Test with int64 input (some JSON parsers produce this)
	ts := ParseTimestamp(int64(1775639179472))
	if ts.IsZero() {
		t.Errorf("expected valid timestamp, got zero time")
	}
}

func TestParseTimestampNil(t *testing.T) {
	beforeCall := time.Now()
	ts := ParseTimestamp(nil)
	afterCall := time.Now()

	if ts.IsZero() {
		t.Errorf("expected non-zero timestamp for nil input")
	}

	// Should be close to now
	if ts.Before(beforeCall.Add(-time.Second)) || ts.After(afterCall.Add(time.Second)) {
		t.Errorf("expected timestamp close to now, got %v", ts)
	}
}

func TestParseTimestampInvalidString(t *testing.T) {
	beforeCall := time.Now()
	ts := ParseTimestamp("invalid-timestamp")
	afterCall := time.Now()

	if ts.IsZero() {
		t.Errorf("expected non-zero timestamp for invalid input")
	}

	// Should fall back to now
	if ts.Before(beforeCall.Add(-time.Second)) || ts.After(afterCall.Add(time.Second)) {
		t.Errorf("expected timestamp close to now for invalid input, got %v", ts)
	}
}

func TestExtractTextSingleBlock(t *testing.T) {
	content := []ContentBlock{
		{Type: "text", Text: "hello"},
	}

	text := ExtractText(content)
	if text != "hello" {
		t.Errorf("expected 'hello', got '%s'", text)
	}
}

func TestExtractTextMultipleBlocks(t *testing.T) {
	content := []ContentBlock{
		{Type: "text", Text: "line 1"},
		{Type: "text", Text: "line 2"},
		{Type: "text", Text: "line 3"},
	}

	text := ExtractText(content)
	expected := "line 1\nline 2\nline 3"
	if text != expected {
		t.Errorf("expected '%s', got '%s'", expected, text)
	}
}

func TestExtractTextIgnoresNonTextBlocks(t *testing.T) {
	content := []ContentBlock{
		{Type: "text", Text: "hello"},
		{Type: "thinking", Text: "internal thought"},
		{Type: "text", Text: "world"},
	}

	text := ExtractText(content)
	expected := "hello\nworld"
	if text != expected {
		t.Errorf("expected '%s', got '%s'", expected, text)
	}
}

func TestExtractTextEmptyBlocks(t *testing.T) {
	content := []ContentBlock{
		{Type: "text", Text: ""},
		{Type: "text", Text: "hello"},
		{Type: "text", Text: ""},
	}

	text := ExtractText(content)
	if text != "hello" {
		t.Errorf("expected 'hello', got '%s'", text)
	}
}

func TestExtractTextEmpty(t *testing.T) {
	content := []ContentBlock{}

	text := ExtractText(content)
	if text != "" {
		t.Errorf("expected empty string, got '%s'", text)
	}
}

func TestExtractTextNoTextBlocks(t *testing.T) {
	content := []ContentBlock{
		{Type: "thinking", Text: "thought 1"},
		{Type: "tool", Text: "tool result"},
	}

	text := ExtractText(content)
	if text != "" {
		t.Errorf("expected empty string, got '%s'", text)
	}
}

func TestExtractTextWithNewlines(t *testing.T) {
	content := []ContentBlock{
		{Type: "text", Text: "line 1\nline 2"},
		{Type: "text", Text: "line 3"},
	}

	text := ExtractText(content)
	if !strings.Contains(text, "line 1") {
		t.Errorf("expected content to contain 'line 1', got '%s'", text)
	}
	if !strings.Contains(text, "line 2") {
		t.Errorf("expected content to contain 'line 2', got '%s'", text)
	}
	if !strings.Contains(text, "line 3") {
		t.Errorf("expected content to contain 'line 3', got '%s'", text)
	}
}

func TestExtractTextPreservesFormatting(t *testing.T) {
	content := []ContentBlock{
		{Type: "text", Text: "**bold**"},
		{Type: "text", Text: "`code`"},
		{Type: "text", Text: "[link](url)"},
	}

	text := ExtractText(content)
	if !strings.Contains(text, "**bold**") {
		t.Errorf("expected markdown formatting preserved, got '%s'", text)
	}
}
