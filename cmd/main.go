package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/maxoov1/recorder/internal/handler"
	"github.com/maxoov1/recorder/internal/recorder/manager"
)

func env(key, value string) string {
	if e, ok := os.LookupEnv(key); ok {
		return e
	}
	return value
}

func main() {
	address := env("RECORDER_ADDRESS", "127.0.0.1:9000")
	base := env("RECORDER_BASE", "recording")

	manager := manager.New(base)
	handler := handler.New(base, manager)

	server := &http.Server{Addr: address, Handler: handler.RegisterRoutes()}

	log.Printf("starting server at %s", address)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("serve: %v", err)
	}
}
