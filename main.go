package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/andreaswachs/forhumans/internal/events"
	"github.com/andreaswachs/forhumans/internal/formatter"
)

func main() {
	showLifecycle := flag.Bool("lifecycle", false, "show AGENT and TURN start/end events")
	flag.Parse()

	opts := formatter.Options{
		ShowLifecycle: *showLifecycle,
	}
	f := formatter.NewWithOptions(opts)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0), 10485760) // 10 MB buffer size

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var event events.Event
		if err := json.Unmarshal(line, &event); err != nil {
			// Skip malformed JSON lines
			continue
		}

		output := f.Format(event)
		if output != "" {
			fmt.Println(output)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}
