package handlers_test

import (
	"backend-coding-challenge-enhanced/internal/handlers"
	"backend-coding-challenge-enhanced/internal/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService for testing UserHandler
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserActionCount(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockUserService) GetUserByID(userID int) (models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(models.User), args.Error(1)
}

func TestGetUserActionCount_Success(t *testing.T) {
	mockService := new(MockUserService)
	userHandler := handlers.NewUserHandler(mockService)

	mockUserID := 1
	expectedCount := 5

	mockService.On("GetUserActionCount", mockUserID).Return(expectedCount, nil)

	router := mux.NewRouter()
	router.HandleFunc("/user/{id}/actions/count", userHandler.GetUserActionCount)

	req := httptest.NewRequest("GET", "/user/1/actions/count", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	mockService.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), strconv.Itoa(expectedCount))
}

func TestGetUserByID_Success(t *testing.T) {
	mockService := new(MockUserService)
	userHandler := handlers.NewUserHandler(mockService)

	mockUserID := 1
	expectedUser := models.User{ID: mockUserID, Name: "John Doe", CreatedAt: "2023-10-06T11:12:22.758Z"}

	mockService.On("GetUserByID", mockUserID).Return(expectedUser, nil)

	router := mux.NewRouter()
	router.HandleFunc("/user/{id}", userHandler.GetUserByID)

	req := httptest.NewRequest("GET", "/user/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	mockService.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
}
