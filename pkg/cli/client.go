package cli

import (
	"context"
	"fmt"

	"github.com/zylisp/repl"
)

// RunClient connects to a remote REPL server and starts a REPL loop.
func RunClient(ctx context.Context, addr string) error {
	if addr == "" {
		return fmt.Errorf("--addr required in client mode")
	}

	client := repl.NewClient()

	fmt.Printf("Connecting to Zylisp REPL at %s...\n", addr)

	if err := client.Connect(ctx, addr); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	fmt.Println("Connected!")

	// Run REPL loop
	return ReplLoop(ctx, client)
}
