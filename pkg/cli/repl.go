package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/zylisp/cli/pkg/eval"
	"github.com/zylisp/repl"
)

// getPrompt resolves the prompt string, handling special cases.
func getPrompt(promptFlag string) string {
	if promptFlag == "alt" {
		return "raskr> "
	}
	return promptFlag
}

// createReadline creates a readline instance with history support.
func createReadline(prompt string) (*readline.Instance, error) {
	// Get home directory for history file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	historyFile := filepath.Join(homeDir, ".zylisp_history")

	// Create readline config
	config := &readline.Config{
		Prompt:          prompt,
		HistoryFile:     historyFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	return readline.NewEx(config)
}

// SimpleReplLoop runs a simple REPL without client/server separation.
// This is used for local mode where evaluation happens directly.
func SimpleReplLoop(ctx context.Context, prompt string) error {
	fmt.Print(Banner)

	// Create readline instance
	rl, err := createReadline(prompt)
	if err != nil {
		return fmt.Errorf("failed to create readline: %w", err)
	}
	defer rl.Close()

	evaluator := eval.Evaluator()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Read line with readline
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Ctrl-C pressed
				continue
			} else if err == io.EOF {
				// Ctrl-D or EOF
				break
			}
			return fmt.Errorf("readline error: %w", err)
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// Handle special commands
		handled, shouldExit := handleCommand(line)
		if handled {
			if shouldExit {
				return nil
			}
			continue // Command was handled, skip evaluation
		}

		// Evaluate expression directly
		result, output, err := evaluator(line)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Display output if any
		if output != "" {
			fmt.Print(output)
		}

		// Display result
		fmt.Println(FormatValue(result))
	}

	return nil
}

// ReplLoop runs the REPL loop with a connected client.
// This is used for client mode where evaluation happens via protocol.
func ReplLoop(ctx context.Context, client repl.Client, prompt string) error {
	fmt.Print(Banner)

	// Create readline instance
	rl, err := createReadline(prompt)
	if err != nil {
		return fmt.Errorf("failed to create readline: %w", err)
	}
	defer rl.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Read line with readline
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Ctrl-C pressed
				continue
			} else if err == io.EOF {
				// Ctrl-D or EOF
				break
			}
			return fmt.Errorf("readline error: %w", err)
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// Handle special commands
		handled, shouldExit := handleCommand(line)
		if handled {
			if shouldExit {
				return nil
			}
			continue // Command was handled, skip evaluation
		}

		// Create context with timeout for evaluation
		evalCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		// Evaluate expression
		result, err := client.Eval(evalCtx, line)
		cancel()

		if err != nil {
			// Protocol error
			fmt.Printf("Protocol Error: %v\n", err)
			continue
		}

		// Check for output
		if result.Output != "" {
			fmt.Print(result.Output)
		}

		// Display result value
		fmt.Println(FormatValue(result.Value))
	}

	return nil
}

// handleCommand processes special REPL commands.
// Returns (handled, shouldExit) where:
//   - handled: true if this was a command (vs regular code to evaluate)
//   - shouldExit: true if the REPL should exit
func handleCommand(line string) (handled bool, shouldExit bool) {
	switch line {
	case "(quit)", "(q)":
		fmt.Println("\n\nSee you at RAGNARǪK!!\n")
		return true, true

	case "(reset)":
		eval.ResetGlobalEnv()
		fmt.Println("Environment reset")
		return true, false

	case "(help)":
		ShowHelp()
		return true, false

	default:
		return false, false
	}
}
