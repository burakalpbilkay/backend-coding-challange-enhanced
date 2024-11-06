package handlers

import (
	"backend-coding-challenge-enhanced/internal/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	repo repositories.UserRepositoryInterface
}

func NewUserHandler(repo repositories.UserRepositoryInterface) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	user, err := h.repo.FetchUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUserActionCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	count, err := h.repo.FetchUserActionCount(id)
	if err != nil {
		http.Error(w, "Failed to fetch action count", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}
