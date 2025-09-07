package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ckalagara/group-a-accounts/core"
	"github.com/ckalagara/group-a-accounts/middleware"
)

const (
	Addr = ":8085"
)

type GenResponse struct {
	Status      string `json:"status,omitempty"`
	RefType     string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

func main() {
	log.Println("Starting service group-a-accounts")

	srvMux := http.NewServeMux()
	h := core.NewHandler(context.Background(), "mongodb://mongodb:27017")

	log.Println("Setting up routes")
	srvMux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		res := GenResponse{
			Status:      "OK",
			RefType:     "HealthCheck",
			Description: "Service is up and running",
		}
		data, err := json.Marshal(&res)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	srvMux.HandleFunc("GET /v1/accounts/{id}", h.GetAccount)
	srvMux.HandleFunc("PUT /v1/accounts", h.UpdateAccount)

	srvMuxDes := http.NewServeMux()
	srvMuxDes.HandleFunc("DELETE /v1/accounts/{id}", h.DeleteAccount)
	srvMux.Handle("/", middleware.ValidateClearence(srvMuxDes))

	srv := &http.Server{
		Addr: Addr,
		Handler: middleware.NewMiddlewareChain(
			middleware.LogRequest,
			middleware.ValidateJWT)(srvMux),
	}
	log.Println("Listening on", Addr)

	srv.ListenAndServe()

}
