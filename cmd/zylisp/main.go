package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zylisp/cli/pkg/cli"
)

var (
	mode      = flag.String("mode", "local", "Mode: 'local', 'server', or 'client'")
	transport = flag.String("transport", "in-process", "Transport: 'in-process', 'unix', or 'tcp'")
	addr      = flag.String("addr", "", "Server address (for server/client modes)")
	codec     = flag.String("codec", "json", "Codec: 'json' or 'msgpack'")
)

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

	// Create configuration
	cfg := cli.Config{
		Mode:      *mode,
		Transport: *transport,
		Addr:      *addr,
		Codec:     *codec,
	}

	// Run the CLI application
	if err := cli.Run(ctx, cfg); err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
