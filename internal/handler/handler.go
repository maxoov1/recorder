package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/maxoov1/recorder/internal/recorder/manager"
)

type Handler struct {
	base    string
	manager *manager.RecorderManager
}

func New(base string, manager *manager.RecorderManager) *Handler {
	return &Handler{base: base, manager: manager}
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /records/", http.StripPrefix(
		"/records/", http.FileServer(http.Dir(h.base))),
	)

	mux.HandleFunc("POST /{identifier}", h.runInstanceHandler)
	mux.HandleFunc("DELETE /{identifier}", h.stopInstanceHandler)

	return mux
}

type Request struct {
	Endpoint string `json:"endpoint"`
}

func (r Request) Validate() error {
	if r.Endpoint == "" {
		return fmt.Errorf("endpoint is empty")
	}

	return nil
}

func (h *Handler) runInstanceHandler(w http.ResponseWriter, r *http.Request) {
	identifier := r.PathValue("identifier")

	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := h.manager.Run(context.Background(), identifier, request.Endpoint); err != nil {
		if errors.Is(err, manager.ErrExist) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		log.Printf("failed to run instance %q (%s): %v", identifier, request.Endpoint, err)

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) stopInstanceHandler(w http.ResponseWriter, r *http.Request) {
	identifier := r.PathValue("identifier")

	if err := h.manager.Stop(identifier); err != nil {
		if errors.Is(err, manager.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Printf("failed to stop instance %q: %v", identifier, err)

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
