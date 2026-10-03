package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/maxoov1/recorder/internal/handler"
	"github.com/maxoov1/recorder/internal/recorder/manager"
)

var (
	_defaultAddress = "127.0.0.1:9000"
	_defaultBase    = "recording"
)

func main() {
	manager := manager.New(_defaultBase)
	handler := handler.New(manager)

	server := &http.Server{Addr: _defaultAddress, Handler: handler.RegisterRoutes()}

	log.Printf("starting server at %s", _defaultAddress)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("serve: %v", err)
	}
}
