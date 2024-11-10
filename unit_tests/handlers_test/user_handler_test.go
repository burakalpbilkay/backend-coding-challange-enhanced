package handlers_test

import (
	"backend-coding-challenge-enhanced/internal/handlers"
	"backend-coding-challenge-enhanced/internal/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserServiceInterface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserActionCount(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

// New GetUserByID method for MockUserService
func (m *MockUserService) GetUserByID(userID int) (models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(models.User), args.Error(1)
}

func TestGetUserActionCount_Success(t *testing.T) {
	mockService := new(MockUserService)
	userHandler := handlers.NewUserHandler(mockService) // Now accepts UserServiceInterface

	mockUserID := 1
	expectedCount := 5

	mockService.On("GetUserActionCount", mockUserID).Return(expectedCount, nil)

	req := httptest.NewRequest("GET", "/user/"+strconv.Itoa(mockUserID)+"/actions/count", nil)
	w := httptest.NewRecorder()

	userHandler.GetUserActionCount(w, req)

	mockService.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), strconv.Itoa(expectedCount))
}

func TestGetUserByID_Success(t *testing.T) {
	mockService := new(MockUserService)
	userHandler := handlers.NewUserHandler(mockService)

	mockUserID := 1
	expectedUser := models.User{ID: mockUserID, Name: "TestName TestSurname", CreatedAt: "2024-10-10T11:12:99.999Z"}

	mockService.On("GetUserByID", mockUserID).Return(expectedUser, nil)

	req := httptest.NewRequest("GET", "/user/"+strconv.Itoa(mockUserID), nil)
	w := httptest.NewRecorder()

	userHandler.GetUserByID(w, req)

	mockService.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TestName TestSurname")
}
