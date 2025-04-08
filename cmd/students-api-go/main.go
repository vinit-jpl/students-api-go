package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vinit-jpl/students-api-go/internal/config"
)

func main() {

	// load config
	cfg := config.MustLoad()

	// database setup

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome To Students API!"))
	})

	// setup server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)                                    // channel to receive os signal
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // notify the channel when os interrupt signal is received

	// creating go routine to run the server in the background
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server: %s", err.Error())
		}
	}()

	<-done
	// fmt.Printf("Server started %s", cfg.HTTPServer.Addr)
	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Fatalf("failed to start server: %s", err.Error())
	// }

	slog.Info("shutting down server...")

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancle()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown gracefully")

}
