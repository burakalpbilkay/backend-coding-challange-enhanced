package handlers_test

import (
	"backend-coding-challenge-enhanced/internal/handlers"
	"backend-coding-challenge-enhanced/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockActionService is a mock implementation of ActionServiceInterface
type MockActionService struct {
	mock.Mock
}

func (m *MockActionService) GetNextActionProbabilities(actionType string) (map[string]float64, error) {
	args := m.Called(actionType)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *MockActionService) GetReferralIndex() (map[int]int, error) {
	args := m.Called()
	return args.Get(0).(map[int]int), args.Error(1)
}

func TestGetNextActionProbabilities_Success(t *testing.T) {
	mockService := new(MockActionService)
	actionHandler := handlers.NewActionHandler(mockService)

	mockActionType := "REFER_USER"
	expectedProbabilities := map[string]float64{"ADD_TO_CRM": 0.7, "VIEW_CONVERSATION": 0.3}

	mockService.On("GetNextActionProbabilities", mockActionType).Return(expectedProbabilities, nil)

	req := httptest.NewRequest("GET", "/action/"+mockActionType+"/next", nil)
	w := httptest.NewRecorder()

	// Simulate calling GetNextActionProbabilities handler
	actionHandler.GetNextActionProbabilities(w, req)

	mockService.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"ADD_TO_CRM":0.7`)
	assert.Contains(t, w.Body.String(), `"VIEW_CONVERSATION":0.3`)
}

func TestGetNextActionProbabilities_NotFound(t *testing.T) {
	mockService := new(MockActionService)
	actionHandler := handlers.NewActionHandler(mockService)

	mockActionType := "NON_EXISTENT_TYPE"
	mockService.On("GetNextActionProbabilities", mockActionType).Return(nil, services.ErrInvalidActionType)

	req := httptest.NewRequest("GET", "/action/"+mockActionType+"/next", nil)
	w := httptest.NewRecorder()

	// Simulate calling GetNextActionProbabilities handler
	actionHandler.GetNextActionProbabilities(w, req)

	mockService.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid action type")
}

func TestGetReferralIndex_Success(t *testing.T) {
	mockService := new(MockActionService)
	actionHandler := handlers.NewActionHandler(mockService)

	expectedReferralIndex := map[int]int{1: 3, 2: 0, 3: 7}
	mockService.On("GetReferralIndex").Return(expectedReferralIndex, nil)

	req := httptest.NewRequest("GET", "/users/referral-index", nil)
	w := httptest.NewRecorder()

	// Simulate calling GetReferralIndex handler
	actionHandler.GetReferralIndex(w, req)

	mockService.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"1":3`)
	assert.Contains(t, w.Body.String(), `"2":0`)
	assert.Contains(t, w.Body.String(), `"3":7`)
}
