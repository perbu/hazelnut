package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/perbu/hazelnut/config"
	"github.com/perbu/hazelnut/server"
)

func main() {
	// Create a configuration programmatically
	cfg := getConfig()
	// Set up a logger with the log level from config
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.GetLogLevel(),
	})
	logger := slog.New(handler)

	// Create context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create a new Hazelnut server
	hazelnut, err := server.New(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to create hazelnut server", "error", err)
		os.Exit(1)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	// Start the server in a goroutine
	go func() {
		defer wg.Done()
		if err := hazelnut.Run(ctx); err != nil {
			logger.Error("hazelnut server error", "error", err)
		}
	}()

	logger.Info("hazelnut server started at :8080")

	// You can do other things in your application here
	// ...

	// Wait for termination signal
	<-ctx.Done()
	logger.Info("shutting down")
	wg.Wait()
	// We could perform additional cleanup here if needed

	fmt.Println("Goodbye!")
}

func getConfig() *config.Config {
	return &config.Config{
		DefaultBackend: config.BackendConfig{
			Target:  "https://example.com:443",
			Timeout: 30 * time.Second,
		},
		Frontend: config.FrontendConfig{
			BaseURL: "http://localhost:8888",
		},
		Cache: config.CacheConfig{
			MaxObj:  "1M",
			MaxCost: "1G",
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "text",
		},
	}
}
