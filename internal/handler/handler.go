package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jfilipedias/rinha-backend-2024-q1-go/internal/entity"
	"github.com/jfilipedias/rinha-backend-2024-q1-go/internal/repository"
)

type CreateTransactionRequestBody struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type Handler struct {
	repository *repository.Repository
}

func NewHandler(repository *repository.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) FindStatementByCustomerID(w http.ResponseWriter, r *http.Request) {
	customerID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
		return
	}

	statement, err := h.repository.FindStatementByCustomerID(customerID)
	if err != nil {
		if err == repository.ErrCustomerNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(statement)
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	customerID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
	}

	var body CreateTransactionRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if (body.Type != "c" && body.Type != "d") || (body.Description == "" || len(body.Description) > 10) {
		http.Error(w, "invalid body", http.StatusUnprocessableEntity)
		return
	}

	transaction := &entity.Transaction{
		Value:       body.Value,
		Type:        body.Type,
		Description: body.Description,
		CustomerID:  customerID,
	}

	result, err := h.repository.CreateTransaction(transaction)
	if err != nil {
		if err == repository.ErrCustomerNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
