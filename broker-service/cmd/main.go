package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vocal-tracker/broker-service/config"
	"github.com/vocal-tracker/broker-service/routes"
)

const webPort = "8080"

func main() {
	log.Printf("Starting broker service on port %s\n", webPort)

	// Initialize config with gRPC connections
	cfg := config.NewConfig()
	if cfg == nil {
		log.Fatal("Failed to initialize config")
	}
	defer cfg.Close()

	// Create router instance
	router := routes.NewRouter(cfg)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: router,
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
