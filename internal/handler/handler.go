package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jfilipedias/rinha-backend-2024-q1-go/internal/repository"
)

type Handler struct {
	repository *repository.Repository
}

func NewHandler(repository *repository.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) FindStatementByCustomerID(w http.ResponseWriter, r *http.Request) {
	customerId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
	}

	statement, err := h.repository.FindStatementByCustomerID(customerId)
	if err != nil {
		if err == repository.ErrCustomerNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(statement)
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	customerId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(customerId)
}
