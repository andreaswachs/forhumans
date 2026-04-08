# forhumans — AI Agent JSONL Log Formatter

A CLI tool that converts structured JSONL event streams from the pi AI agent harness into human-readable log output for CI workflows.

## Overview

When AI agents run in non-interactive CI sessions, they emit structured JSON events to stdout. This tool pipes those JSONL lines into readable formatted logs, showing the agent's thinking, responses, and tool calls.

## Building

```bash
go build -o forhumans-cli ./cmd/...
```

## Usage

Pipe JSONL event stream to stdin:

```bash
cat output.jsonl | ./forhumans-cli
```

Or in CI:

```bash
./agent-harness | ./forhumans-cli
```

### Flags

- `--lifecycle` — Show AGENT START/END and TURN START/END events (default: false)

```bash
# Show full lifecycle structure
cat output.jsonl | ./forhumans-cli --lifecycle

# Minimal output (default)
cat output.jsonl | ./forhumans-cli
```

## Output Format

Each log line follows the format: `HH:MM:SS [TAG] content`

### Event Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `[SESSION]` | Session metadata | `[SESSION] id=680697ea cwd=/Users/awa/Source/forhumans` |
| `[AGENT START/END]` | Agent lifecycle | `[AGENT START]` |
| `[TURN START/END]` | Conversation turn boundaries | `[TURN START]` |
| `[USER]` | User message | `[USER] run` |
| `[THINKING]` | LLM reasoning (truncated at 200 chars) | `[THINKING] User says "run". Probably want to list files...` |
| `[RESPONSE]` | LLM response | `[RESPONSE] Here's the output...` |
| `[TOOL CALL]` | Tool invocation with arguments | `[TOOL CALL] bash({"command":"ls -R ."})`  |
| `[TOOL START]` | Tool execution beginning | `[TOOL START] bash args={...}` |
| `[TOOL END]` | Tool execution complete with result | `[TOOL END] bash error=false result=output.jsonl` |
| `[UNKNOWN]` | Unknown event type (raw JSON) | `[UNKNOWN] {"type":"custom_event",...}` |

### Filtering Noise

The formatter skips delta events and intermediate updates to keep output clean:
- `*_delta` events (thinking_delta, text_delta)
- `*_start` / `*_update` intermediate events (except tool_execution_start/end)
- Non-user message_start events

Unknown event types are printed as `[UNKNOWN]` with the raw JSON, allowing visibility into new event types without code changes.

### Multi-line Content Indentation

Responses and other multi-line content are automatically indented to align continuation lines with the start of the content on the first line:

```
11:59:02 [RESPONSE] Created **`random_numbers.json`** with a JSON array:

                    ```json
                    [42, 7, 19]
                    ```
```

This improves readability when responses span multiple lines, especially in CI logs with many concurrent messages.

## Architecture

The project uses clean package separation:

- **`internal/events/`** — Event types and JSON unmarshaling, timestamp parsing, text extraction
- **`internal/formatter/`** — Formatting logic, truncation, output generation
- **`cmd/main.go`** — CLI entry point, JSONL reader

## Testing

Run all tests:

```bash
go test -v ./...
```

Coverage includes:
- **40 unit tests** across both packages
- All event types (session, agent, turn, user, thinking, response, tool call/execution)
- Edge cases (empty content, truncation, newlines, unknown types)
- Timestamp parsing (ISO 8601 and milliseconds epoch)
- Text extraction and block filtering

## Examples

### Default output (minimal, no lifecycle events)

```bash
$ cat output.jsonl | ./forhumans-cli
09:06:19 [SESSION] id=680697ea cwd=/Users/awa/Source/forhumans
09:06:19 [USER] list files
09:06:20 [THINKING] User wants to see files. I'll use bash ls.
09:06:20 [TOOL CALL] bash({"command":"ls -R ."})
09:06:20 [TOOL START] bash args={"command":"ls -R ."}
09:06:20 [TOOL END] bash error=false result=output.jsonl random_numbers.json
09:06:21 [RESPONSE] Here are the files in the directory...
```

### With `--lifecycle` flag (shows full structure)

```bash
$ cat output.jsonl | ./forhumans-cli --lifecycle
09:06:19 [SESSION] id=680697ea cwd=/Users/awa/Source/forhumans
09:06:19 [AGENT START]
09:06:19 [TURN START]
09:06:19 [USER] list files
09:06:20 [THINKING] User wants to see files. I'll use bash ls.
09:06:20 [TOOL CALL] bash({"command":"ls -R ."})
09:06:20 [TOOL START] bash args={"command":"ls -R ."}
09:06:20 [TOOL END] bash error=false result=output.jsonl random_numbers.json
09:06:21 [RESPONSE] Here are the files in the directory...
09:06:21 [TURN END]
09:06:21 [AGENT END]
```

### With unknown event types

```
09:06:19 [TURN START]
09:06:19 [USER] do something
09:06:19 [UNKNOWN] {"type":"custom_event","data":"..."}
09:06:20 [RESPONSE] Done
09:06:21 [TURN END]
```

## Development

### Adding a new event type

1. Add handling in `internal/formatter/formatter.go`'s `Format()` method
2. Add unit tests in `internal/formatter/formatter_test.go`
3. Run `go test -v ./...` to verify

Unknown types automatically fall through to the default case, which prints raw JSON.

## Dependencies

None — uses only Go standard library.
