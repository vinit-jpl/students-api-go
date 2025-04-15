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
	"github.com/vinit-jpl/students-api-go/internal/http/handlers/student"
	"github.com/vinit-jpl/students-api-go/internal/storage/sqlite"
)

func main() {

	// load config
	cfg := config.MustLoad()

	// database setup

	storage, err := sqlite.New(cfg) // create a new sqlite connection
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))         // register the student handler
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage)) // register the student handler

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown gracefully")

}

// 1:30:21
