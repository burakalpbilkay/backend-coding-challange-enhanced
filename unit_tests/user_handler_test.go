package handlers

import (
	"backend-coding-challenge-enhanced/internal/handlers"
	"backend-coding-challenge-enhanced/internal/repositories"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetUserByID tests the GetUserByID function in UserHandler
func TestGetUserByID(t *testing.T) {
	mockUserRepo := &repositories.MockUserRepo{}
	userHandler := handlers.NewUserHandler(mockUserRepo)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler.GetUserByID)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}
}

// TestGetUserActionCount tests the GetUserActionCount function in UserHandler
func TestGetUserActionCount(t *testing.T) {
	mockUserRepo := &repositories.MockUserRepo{}
	userHandler := handlers.NewUserHandler(mockUserRepo)

	req, _ := http.NewRequest("GET", "/user/1/actions/count", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler.GetUserActionCount)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}
}
