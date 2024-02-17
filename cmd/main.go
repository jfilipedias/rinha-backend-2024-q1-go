package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jfilipedias/rinha-backend-2024-q1-go/internal/handler"
	"github.com/jfilipedias/rinha-backend-2024-q1-go/internal/repository"
)

func main() {
	db, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_DSN"))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	mux := http.NewServeMux()

	repository := repository.NewRepository(db)
	handler := handler.NewHandler(repository)
	mux.HandleFunc("GET /clientes/{id}/extrato", handler.FindStatementByCustomerID)
	mux.HandleFunc("POST /clientes/{id}/transacoes", handler.CreateTransaction)

	port := os.Getenv("API_PORT")
	fmt.Println("Server listen into port " + port)
	http.ListenAndServe(":"+port, mux)
}
