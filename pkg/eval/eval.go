package eval

import (
	"fmt"

	"github.com/zylisp/lang/interpreter"
	"github.com/zylisp/lang/parser"
)

// globalEnv is the shared environment for the REPL
var globalEnv *interpreter.Env

func init() {
	globalEnv = interpreter.NewEnv(nil)
	interpreter.LoadPrimitives(globalEnv)
}

// Evaluator returns a function that evaluates Zylisp code in the global environment.
// This function is suitable for use with repl.ServerConfig.
func Evaluator() func(code string) (interface{}, string, error) {
	return func(code string) (interface{}, string, error) {
		return EvaluateWithEnv(globalEnv, code)
	}
}

// EvaluateWithEnv evaluates Zylisp code in a specific environment.
// This is primarily used for testing with isolated environments.
func EvaluateWithEnv(env *interpreter.Env, code string) (interface{}, string, error) {
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

// ResetGlobalEnv resets the global environment to its initial state.
// This is used by the (reset) command.
func ResetGlobalEnv() {
	globalEnv = interpreter.NewEnv(nil)
	interpreter.LoadPrimitives(globalEnv)
}

// GlobalEnv returns the global environment.
// This is primarily used for testing.
func GlobalEnv() *interpreter.Env {
	return globalEnv
}
