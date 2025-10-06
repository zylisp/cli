package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/zylisp/lang/interpreter"
	"github.com/zylisp/lang/parser"
	"github.com/zylisp/repl"
)

const banner = `
╔═══════════════════════════════════════╗
║                                       ║
║           Zylisp REPL v0.1.0          ║
║                                       ║
║      A Lisp that compiles to Go       ║
║                                       ║
╚═══════════════════════════════════════╝

Type expressions and press Enter.
Type 'exit' or 'quit' to leave.
Type ':reset' to clear the environment.
Type ':help' for more commands.

`

var (
	mode      = flag.String("mode", "local", "Mode: 'local', 'server', or 'client'")
	transport = flag.String("transport", "in-process", "Transport: 'in-process', 'unix', or 'tcp'")
	addr      = flag.String("addr", "", "Server address (for server/client modes)")
	codec     = flag.String("codec", "json", "Codec: 'json' or 'msgpack'")
)

// Global environment for the evaluator
var globalEnv *interpreter.Env

func init() {
	globalEnv = interpreter.NewEnv(nil)
	interpreter.LoadPrimitives(globalEnv)
}

// evaluateZylisp is the bridge to the actual Zylisp evaluator
func evaluateZylisp(code string) (interface{}, string, error) {
	return evaluateZylispWithEnv(globalEnv, code)
}

// evaluateZylispWithEnv evaluates code in a specific environment
func evaluateZylispWithEnv(env *interpreter.Env, code string) (interface{}, string, error) {
	// Tokenize
	tokens, err := parser.Tokenize(code)
	if err != nil {
		return nil, "", fmt.Errorf("tokenize error: %w", err)
	}

	// Parse
	expr, err := parser.Read(tokens)
	if err != nil {
		return nil, "", fmt.Errorf("parse error: %w", err)
	}

	// Evaluate
	result, err := interpreter.Eval(expr, env)
	if err != nil {
		return nil, "", fmt.Errorf("eval error: %w", err)
	}

	// Return result as interface{} and its string representation as output
	return result, "", nil
}

func main() {
	flag.Parse()

	// Set up context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Run appropriate mode
	var err error
	switch *mode {
	case "local":
		err = runLocal(ctx)
	case "server":
		err = runServer(ctx)
	case "client":
		err = runClient(ctx)
	default:
		fmt.Fprintf(os.Stderr, "Unknown mode: %s\n", *mode)
		flag.Usage()
		os.Exit(1)
	}

	if err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runLocal runs the REPL in local mode (in-process server/client)
func runLocal(ctx context.Context) error {
	// Create in-process server
	config := repl.ServerConfig{
		Transport: "in-process",
		Evaluator: evaluateZylisp,
	}

	server, err := repl.NewServer(config)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Start server in background
	serverCtx, serverCancel := context.WithCancel(ctx)
	defer serverCancel()

	go func() {
		if err := server.Start(serverCtx); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		}
	}()

	// For in-process mode, we need to use the in-process transport directly
	// since the universal client doesn't support in-process yet.
	// Instead, let's just run a simple REPL loop directly
	return runSimpleREPL(ctx)
}

// runSimpleREPL runs a simple REPL without client/server separation
func runSimpleREPL(ctx context.Context) error {
	fmt.Print(banner)

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
		if shouldExit := handleCommand(line); shouldExit {
			return nil
		}

		// Evaluate expression directly
		result, output, err := evaluateZylisp(line)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Display output if any
		if output != "" {
			fmt.Print(output)
		}

		// Display result
		fmt.Println(formatValue(result))
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// runServer starts a REPL server
func runServer(ctx context.Context) error {
	if *addr == "" {
		return fmt.Errorf("--addr required in server mode")
	}

	config := repl.ServerConfig{
		Transport: *transport,
		Addr:      *addr,
		Codec:     *codec,
		Evaluator: evaluateZylisp,
	}

	server, err := repl.NewServer(config)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	fmt.Printf("Starting Zylisp REPL server on %s (%s)\n", server.Addr(), *transport)

	return server.Start(ctx)
}

// runClient connects to a remote REPL server
func runClient(ctx context.Context) error {
	if *addr == "" {
		return fmt.Errorf("--addr required in client mode")
	}

	client := repl.NewClient()

	fmt.Printf("Connecting to Zylisp REPL at %s...\n", *addr)

	if err := client.Connect(ctx, *addr); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	fmt.Println("Connected!")

	// Run REPL loop
	return replLoop(ctx, client)
}

// replLoop runs the REPL loop with a connected client
func replLoop(ctx context.Context, client repl.Client) error {
	fmt.Print(banner)

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
		if shouldExit := handleCommand(line); shouldExit {
			return nil
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
		fmt.Println(formatValue(result.Value))
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// handleCommand processes special REPL commands
// Returns true if the REPL should exit
func handleCommand(line string) bool {
	switch line {
	case "exit", "quit":
		fmt.Println("\nGoodbye!")
		return true

	case ":reset":
		// Reset global environment
		globalEnv = interpreter.NewEnv(nil)
		interpreter.LoadPrimitives(globalEnv)
		fmt.Println("Environment reset")
		return false

	case ":help":
		showHelp()
		return false

	default:
		return false
	}
}

// formatValue converts the result value to a display string
func formatValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	// If it's a Zylisp value with String() method, use it
	if stringer, ok := value.(fmt.Stringer); ok {
		return stringer.String()
	}

	// Otherwise use default formatting
	return fmt.Sprintf("%v", value)
}

func showHelp() {
	fmt.Print(`
Zylisp REPL - Command Reference

REPL Commands:
  exit, quit    - Exit the REPL
  :reset        - Reset the environment
  :help         - Show this help message

Language Features:
  Special Forms: define, lambda, if, quote
  Arithmetic:    +, -, *, /
  Comparison:    =, <, >, <=, >=
  Lists:         list, car, cdr, cons
  Predicates:    number?, symbol?, list?, null?

Connection Modes:
  Local mode:  Built-in REPL (default)
    $ zylisp

  Server mode: Start a REPL server
    $ zylisp --mode=server --transport=tcp --addr=:5555

  Client mode: Connect to a remote REPL
    $ zylisp --mode=client --addr=localhost:5555

Examples:
  > (+ 1 2)
  3

  > (define square (lambda (x) (* x x)))
  <function>

  > (square 5)
  25

For more examples, see EXAMPLES.md
`)
}
