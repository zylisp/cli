package cli

import (
	"context"
	"fmt"

	"github.com/zylisp/cli/pkg/eval"
	"github.com/zylisp/repl"
)

// RunServer starts a REPL server.
func RunServer(ctx context.Context, transport, addr, codec string) error {
	if addr == "" {
		return fmt.Errorf("--addr required in server mode")
	}

	config := repl.ServerConfig{
		Transport: transport,
		Addr:      addr,
		Codec:     codec,
		Evaluator: eval.Evaluator(),
	}

	server, err := repl.NewServer(config)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	fmt.Printf("Starting Zylisp REPL server on %s (%s)\n", server.Addr(), transport)

	return server.Start(ctx)
}
