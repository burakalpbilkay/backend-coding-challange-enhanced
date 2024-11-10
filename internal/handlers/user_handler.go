package handlers

import (
	"backend-coding-challenge-enhanced/internal/helpers"
	"backend-coding-challenge-enhanced/internal/repositories"
	"backend-coding-challenge-enhanced/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		helpers.JSONError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		helpers.JSONError(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUserActionCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	// Convert userID to an integer, validating that it's numeric
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		helpers.JSONError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	count, err := h.userService.GetUserActionCount(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			helpers.JSONError(w, "User not found", http.StatusNotFound)
			return
		}
		helpers.JSONError(w, "Failed to retrieve action count", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}
