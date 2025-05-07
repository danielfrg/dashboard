package main

import (
	"os"
	"log/slog"

	"dashboard/internal/server"
)

func main() {
	// Run your server.
	if err := server.RunServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}
}
