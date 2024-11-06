package handlers

import (
	"backend-coding-challenge-enhanced/internal/handlers"
	"backend-coding-challenge-enhanced/internal/repositories"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// TestGetNextActionProbabilities tests the GetNextActionProbabilities function in ActionHandler
func TestGetNextActionProbabilities(t *testing.T) {
	mockActionRepo := &repositories.MockActionRepo{}
	actionHandler := handlers.NewActionHandler(mockActionRepo)

	req, err := http.NewRequest("GET", "/action/ADD_TO_CRM/next", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/action/{type}/next", actionHandler.GetNextActionProbabilities)
	router.ServeHTTP(rr, req)

	// Check the HTTP status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code is wrong: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := map[string]float64{
		"ADD_TO_CRM":        0.70,
		"REFER_USER":        0.20,
		"VIEW_CONVERSATION": 0.10,
	}
	var actual map[string]float64
	json.NewDecoder(rr.Body).Decode(&actual)

	for key, value := range expected {
		if actual[key] != value {
			t.Errorf("Wrong result for %s: got %v want %v", key, actual[key], value)
		}
	}
}

// TestGetReferralIndex tests the GetReferralIndex function in ActionHandler
func TestGetReferralIndex(t *testing.T) {
	mockActionRepo := &repositories.MockActionRepo{}
	actionHandler := handlers.NewActionHandler(mockActionRepo)

	req, err := http.NewRequest("GET", "/users/referral-index", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(actionHandler.GetReferralIndex)
	handler.ServeHTTP(rr, req)

	// Check the HTTP status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code is wrong: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := map[int]int{
		1: 3,
		2: 0,
		3: 7,
	}
	var actual map[int]int
	json.NewDecoder(rr.Body).Decode(&actual)

	for key, value := range expected {
		if actual[key] != value {
			t.Errorf("Wrong result for user %d: got %v want %v", key, actual[key], value)
		}
	}
}
