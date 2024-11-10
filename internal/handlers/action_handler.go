package handlers

import (
	"backend-coding-challenge-enhanced/internal/helpers"
	"backend-coding-challenge-enhanced/internal/services"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type ActionHandler struct {
	actionService services.ActionServiceInterface
}

func NewActionHandler(actionService services.ActionServiceInterface) *ActionHandler {
	return &ActionHandler{actionService: actionService}
}

func (h *ActionHandler) GetNextActionProbabilities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actionType := vars["type"]

	probabilities, err := h.actionService.GetNextActionProbabilities(actionType)
	if err != nil {
		if errors.Is(err, services.ErrInvalidActionType) {
			helpers.JSONError(w, "Action type not found", http.StatusBadRequest)
			return
		}
		helpers.JSONError(w, "Failed to fetch probabilities", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(probabilities)
}

func (h *ActionHandler) GetReferralIndex(w http.ResponseWriter, r *http.Request) {
	referralIndex, err := h.actionService.GetReferralIndex()
	if err != nil {
		helpers.JSONError(w, "Failed to calculate referral index", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(referralIndex)
}
