package core

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ckalagara/group-a-accounts/model"
)

type Handler interface {
	UpdateAccount(w http.ResponseWriter, r *http.Request)
	GetAccount(w http.ResponseWriter, r *http.Request)
	DeleteAccount(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service Service
}

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

// PUT /accounts
func (h *handler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var acc model.Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := h.service.UpdateAccount(acc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, updated)
}

// GET /accounts/{id}
func (h *handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "missing account ID", http.StatusBadRequest)
		return
	}

	account, err := h.service.GetAccount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, account)
}

// DELETE /accounts/{id}
func (h *handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		http.Error(w, "missing account ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAccount(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper: write JSON response
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Extract ID from /accounts/{id}
func getIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 2 && parts[0] == "accounts" {
		return parts[1]
	}
	return ""
}
