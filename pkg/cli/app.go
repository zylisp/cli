package cli

import (
	"context"
	"fmt"
)

// Config holds the CLI configuration.
type Config struct {
	Mode      string
	Transport string
	Addr      string
	Codec     string
	Prompt    string
}

// Run executes the CLI application with the given configuration.
func Run(ctx context.Context, cfg Config) error {
	switch cfg.Mode {
	case "local":
		return RunLocal(ctx, cfg.Prompt)
	case "server":
		return RunServer(ctx, cfg.Transport, cfg.Addr, cfg.Codec)
	case "client":
		return RunClient(ctx, cfg.Addr, cfg.Prompt)
	default:
		return fmt.Errorf("unknown mode: %s", cfg.Mode)
	}
}
