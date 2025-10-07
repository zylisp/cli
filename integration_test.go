package main

import (
	"fmt"
	"testing"

	"github.com/zylisp/cli/pkg/cli"
	"github.com/zylisp/cli/pkg/eval"
	"github.com/zylisp/lang/interpreter"
)

// setupTestEnv creates a fresh test environment
func setupTestEnv(t *testing.T) *interpreter.Env {
	t.Helper()
	env := interpreter.NewEnv(nil)
	interpreter.LoadPrimitives(env)
	return env
}

// evalTestCode evaluates code in the given environment
func evalTestCode(t *testing.T, env *interpreter.Env, code string) string {
	t.Helper()

	result, _, err := eval.EvaluateWithEnv(env, code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return cli.FormatValue(result)
}

func TestIntegrationBasic(t *testing.T) {
	env := setupTestEnv(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"number", "42", "42"},
		{"add", "(+ 1 2)", "3"},
		{"nested", "(+ (* 2 3) 4)", "10"},
		{"string", `"hello"`, `"hello"`},
		{"boolean", "true", "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalTestCode(t, env, tt.input)

			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestIntegrationStateful(t *testing.T) {
	env := setupTestEnv(t)

	steps := []struct {
		input    string
		expected string
	}{
		{"(define x 10)", "10"},
		{"x", "10"},
		{"(define y 20)", "20"},
		{"(+ x y)", "30"},
		{"(define add (lambda (a b) (+ a b)))", "<function>"},
		{"(add x y)", "30"},
		{"(add 5 7)", "12"},
	}

	for i, step := range steps {
		result := evalTestCode(t, env, step.input)

		if result != step.expected {
			t.Errorf("step %d: got %q, want %q", i, result, step.expected)
		}
	}
}

func TestIntegrationFactorial(t *testing.T) {
	env := setupTestEnv(t)

	// Define factorial function
	factorialDef := `
        (define factorial
          (lambda (n)
            (if (<= n 1)
                1
                (* n (factorial (- n 1))))))
    `

	_ = evalTestCode(t, env, factorialDef)

	tests := []struct {
		n        int
		expected string
	}{
		{0, "1"},
		{1, "1"},
		{5, "120"},
		{6, "720"},
	}

	for _, tt := range tests {
		input := fmt.Sprintf("(factorial %d)", tt.n)
		result := evalTestCode(t, env, input)

		if result != tt.expected {
			t.Errorf("factorial(%d): got %q, want %q", tt.n, result, tt.expected)
		}
	}
}

func TestIntegrationListProcessing(t *testing.T) {
	env := setupTestEnv(t)

	// Define sum function
	sumDef := `
        (define sum
          (lambda (lst)
            (if (null? lst)
                0
                (+ (car lst) (sum (cdr lst))))))
    `

	_ = evalTestCode(t, env, sumDef)

	result := evalTestCode(t, env, "(sum (list 1 2 3 4 5))")

	if result != "15" {
		t.Errorf("sum: got %q, want \"15\"", result)
	}
}
