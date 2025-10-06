package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zylisp/repl/client"
	"github.com/zylisp/repl/server"
)

const banner = `
╔═══════════════════════════════════════╗
║                                       ║
║           Zylisp REPL v0.0.1          ║
║                                       ║
║      A Lisp that compiles to Go       ║
║                                       ║
╚═══════════════════════════════════════╝

Type expressions and press Enter.
Type 'exit' or 'quit' to leave.
Type ':reset' to clear the environment.
Type ':help' for more commands.

`

func main() {
	fmt.Print(banner)

	// Create server and client
	srv := server.NewServer()
	cli := client.NewClient(srv)

	// Create scanner for input
	scanner := bufio.NewScanner(os.Stdin)

	// REPL loop
	for {
		fmt.Print("> ")

		// Read input
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Handle special commands
		if handleCommand(line, cli) {
			continue
		}

		// Evaluate expression
		result, err := cli.Send(line)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println(result)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nGoodbye!")
}

// handleCommand handles special REPL commands
// Returns true if the command was handled, false otherwise
func handleCommand(line string, cli *client.Client) bool {
	switch line {
	case "exit", "quit":
		fmt.Println("\nGoodbye!")
		os.Exit(0)
		return true

	case ":reset":
		cli.Reset()
		fmt.Println("Environment reset")
		return true

	case ":help":
		showHelp()
		return true

	default:
		return false
	}
}

func showHelp() {
	fmt.Print(`
Available Commands:
  exit, quit    - Exit the REPL
  :reset        - Reset the environment
  :help         - Show this help message

Special Forms:
  define        - Define a variable: (define x 42)
  lambda        - Create a function: (lambda (x) (* x x))
  if            - Conditional: (if test then else)
  quote         - Quote an expression: (quote (1 2 3))

Primitives:
  Arithmetic:   +, -, *, /
  Comparison:   =, <, >, <=, >=
  Lists:        list, car, cdr, cons
  Predicates:   number?, symbol?, list?, null?

Examples:
  > (+ 1 2)
  3

  > (define square (lambda (x) (* x x)))
  <function>

  > (square 5)
  25

  > (if (> 5 3) "yes" "no")
  "yes"
`)
}
