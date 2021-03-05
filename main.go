package main

import (
	"context"
	"fmt"
	"github.com/peterkrauz/go-rest-api-playground/db"
	"github.com/peterkrauz/go-rest-api-playground/handler"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ADDRESS = ":8080"

func main() {
	listener := setupHttpListener()
	database := setupDatabase()
	server := setupServer(database)

	go func() { server.Serve(listener) }()
	defer Stop(server)
	log.Printf("Server started on %s", ADDRESS)

	setupServerChannel()
}

func setupHttpListener() net.Listener {
	listener, err := net.Listen("tcp", ADDRESS)

	if err != nil {
		log.Fatalf("Unexpected error: %s", err.Error())
	}
	return listener
}

func setupDatabase() db.Database {
	dbUser, dbPassword, dbName := os.Getenv("POSTGRES_USE"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	database, err := db.Initialize(dbUser, dbPassword, dbName)

	if err != nil {
		log.Fatalf("Could not setup database: %v", err)
	}

	defer database.Connection.Close()
	return database
}

func setupServer(database db.Database) *http.Server {
	httpHandler := handler.NewHandler(database)
	server := &http.Server{Handler: httpHandler}
	return server
}

func setupServerChannel() {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-channel))
	log.Println("Stopping server")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server unable to shutdown: %v", err)
		os.Exit(1)
	}
}
