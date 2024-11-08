package handlers

import (
	"backend-coding-challenge-enhanced/internal/repositories"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type ActionHandler struct {
	repo repositories.ActionRepositoryInterface
}

func NewActionHandler(repo repositories.ActionRepositoryInterface) *ActionHandler {
	return &ActionHandler{repo: repo}
}

func (h *ActionHandler) GetNextActionProbabilities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actionType := vars["type"]

	probabilities, err := h.repo.FetchNextActionProbabilities(actionType)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidActionType) {
			http.Error(w, "Action type not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch probabilities", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(probabilities)
}

func (h *ActionHandler) GetReferralIndex(w http.ResponseWriter, r *http.Request) {
	referralIndex, err := h.repo.FetchReferralIndex()
	if err != nil {
		http.Error(w, "Failed to calculate referral index", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(referralIndex)
}
