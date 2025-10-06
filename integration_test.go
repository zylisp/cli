package main

import (
	"fmt"
	"testing"

	"github.com/zylisp/repl/client"
	"github.com/zylisp/repl/server"
)

func TestIntegrationBasic(t *testing.T) {
	srv := server.NewServer()
	cli := client.NewClient(srv)

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
			result, err := cli.Send(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestIntegrationStateful(t *testing.T) {
	srv := server.NewServer()
	cli := client.NewClient(srv)

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
		result, err := cli.Send(step.input)
		if err != nil {
			t.Fatalf("step %d error: %v", i, err)
		}

		if result != step.expected {
			t.Errorf("step %d: got %q, want %q", i, result, step.expected)
		}
	}
}

func TestIntegrationFactorial(t *testing.T) {
	srv := server.NewServer()
	cli := client.NewClient(srv)

	// Define factorial function
	factorialDef := `
        (define factorial
          (lambda (n)
            (if (<= n 1)
                1
                (* n (factorial (- n 1))))))
    `

	_, err := cli.Send(factorialDef)
	if err != nil {
		t.Fatalf("define factorial error: %v", err)
	}

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
		result, err := cli.Send(input)
		if err != nil {
			t.Fatalf("factorial(%d) error: %v", tt.n, err)
		}

		if result != tt.expected {
			t.Errorf("factorial(%d): got %q, want %q", tt.n, result, tt.expected)
		}
	}
}

func TestIntegrationListProcessing(t *testing.T) {
	srv := server.NewServer()
	cli := client.NewClient(srv)

	// Define sum function
	sumDef := `
        (define sum
          (lambda (lst)
            (if (null? lst)
                0
                (+ (car lst) (sum (cdr lst))))))
    `

	_, err := cli.Send(sumDef)
	if err != nil {
		t.Fatalf("define sum error: %v", err)
	}

	result, err := cli.Send("(sum (list 1 2 3 4 5))")
	if err != nil {
		t.Fatalf("sum error: %v", err)
	}

	if result != "15" {
		t.Errorf("sum: got %q, want \"15\"", result)
	}
}
