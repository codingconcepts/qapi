package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/login", handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/api/get/{token}", handleGet).Methods(http.MethodGet)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": map[string]interface{}{
			"token": "MXU01u5KSMQ0SCNL4/6AFuP+DhZ7AoXWTIfmd7gl6Sp6vJQn0C2w6A/NsqZoBeGnZpw",
			"id":    "4977feb8-fac2-4c2a-b608-771ae8b0f081",
		},
	})
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)
	log.Println(r.Header)
}
