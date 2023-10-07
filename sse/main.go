package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/comp318/tutorials/sse/count"
)

func main() {
	var srv http.Server

	portFlag := flag.Int("p", 5318, "port number")
	flag.Parse()

	// Only handle "/count" endpoint.
	// Return 404 for all other paths.
	mux := http.NewServeMux()
	mux.Handle("/count", count.New())

	srv.Addr = fmt.Sprintf(":%d", *portFlag)
	srv.Handler = mux

	slog.Info("Count server listening", "port", *portFlag)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Server closed", "error", err)
	} else {
		slog.Info("Server closed", "error", err)
	}
}
