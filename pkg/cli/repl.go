package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zylisp/cli/pkg/eval"
	"github.com/zylisp/repl"
)

// SimpleReplLoop runs a simple REPL without client/server separation.
// This is used for local mode where evaluation happens directly.
func SimpleReplLoop(ctx context.Context) error {
	fmt.Print(Banner)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

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
		evaluator := eval.Evaluator()
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

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// ReplLoop runs the REPL loop with a connected client.
// This is used for client mode where evaluation happens via protocol.
func ReplLoop(ctx context.Context, client repl.Client) error {
	fmt.Print(Banner)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

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

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
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
		fmt.Println("\nGoodbye!")
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
