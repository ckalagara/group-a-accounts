package main

import (
	"encoding/json"
	"log"
	"net/http"
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

	log.Println("Setting up routes")

	srvMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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

	srv := &http.Server{
		Addr:    Addr,
		Handler: srvMux,
	}
	log.Println("Listening on", Addr)

	srv.ListenAndServe()

}
