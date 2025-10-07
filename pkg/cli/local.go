package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/zylisp/cli/pkg/eval"
	"github.com/zylisp/repl"
)

// RunLocal runs the REPL in local mode (in-process server/client).
func RunLocal(ctx context.Context, promptFlag string) error {
	// Create in-process server
	config := repl.ServerConfig{
		Transport: "in-process",
		Evaluator: eval.Evaluator(),
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
	// Instead, we run a simple REPL loop directly
	prompt := getPrompt(promptFlag)
	return SimpleReplLoop(ctx, prompt)
}
