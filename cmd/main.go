package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jfilipedias/rinha-backend-2024-q1-go/internal/handler"
	"github.com/jfilipedias/rinha-backend-2024-q1-go/internal/repository"
)

func main() {
	db, err := pgxpool.New(context.Background(), "user=local_user password=local_password dbname=local_db host=localhost port=5432 sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	mux := http.NewServeMux()

	repository := repository.NewRepository(db)
	handler := handler.NewHandler(repository)
	mux.HandleFunc("GET /clientes/{id}/extrato", handler.FindStatementByCustomerID)
	mux.HandleFunc("POST /clientes/{id}/transacoes", handler.CreateTransaction)

	fmt.Println("Server listen into port 8080")
	http.ListenAndServe(":8080", mux)
}
