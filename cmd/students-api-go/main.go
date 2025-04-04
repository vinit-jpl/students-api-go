package main

import (
	"fmt"
	"log"
	"net/http"

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
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Println("Server started")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}

}
