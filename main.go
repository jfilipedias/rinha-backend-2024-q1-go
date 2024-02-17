package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /clientes/{id}/extrato", FindStatementByCustomerID)

	mux.HandleFunc("POST /clientes/{id}/transacoes", CreateTransaction)

	fmt.Println("Server listen into port 8080")
	http.ListenAndServe(":8080", mux)
}

func FindStatementByCustomerID(w http.ResponseWriter, r *http.Request) {
	customerId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(customerId)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	customerId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(customerId)
}
