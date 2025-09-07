package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"redis-learning/internal/server"
)

func main() {
	// Parse command line flags
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "6379", "Server port")
	flag.Parse()

	// Create server
	srv := server.NewServer(*host, *port)

	// Handle graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down server...")
		srv.Stop()
		os.Exit(0)
	}()

	// Start server
	log.Printf("Starting Redis server on %s:%s", *host, *port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
