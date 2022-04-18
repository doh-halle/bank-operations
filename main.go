package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize Router
	r := mux.NewRouter()

	// Endpoints & Router Handlers
	r.HandleFunc("/api/create_customer", create_customer).Methods("POST")
	r.HandleFunc("/api/deposit_cash/{id}", deposit_cash).Methods("POST")
	r.HandleFunc("/api/withdraw_cash/{id}", withdraw_cash).Methods("POST")
	r.HandleFunc("/api/account_balance/{id}", account_balance).Methods("GET")

	// Web Server
	//port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":7000", r))
	//":"+port
}
