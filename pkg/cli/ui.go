package cli

import (
	_ "embed"
	"fmt"
)

//go:embed banners/banner-1.txt
var Banner string

// FormatValue converts a result value to a display string.
func FormatValue(value interface{}) string {
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

// ShowHelp displays the help message for the REPL.
func ShowHelp() {
	fmt.Print(`
Zylisp REPL - Command Reference

REPL Commands:
  (quit), (q)   - Exit the REPL
  (reset)       - Reset the environment
  (help)        - Show this help message

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
