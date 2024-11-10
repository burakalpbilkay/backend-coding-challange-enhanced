package handlers_test

import (
	"backend-coding-challenge-enhanced/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockActionService for testing ActionHandler
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

	router := mux.NewRouter()
	router.HandleFunc("/action/{type}/next", actionHandler.GetNextActionProbabilities)

	req := httptest.NewRequest("GET", "/action/REFER_USER/next", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	mockService.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"ADD_TO_CRM":0.7`)
	assert.Contains(t, w.Body.String(), `"VIEW_CONVERSATION":0.3`)
}

func TestGetReferralIndex_Success(t *testing.T) {
	mockService := new(MockActionService)
	actionHandler := handlers.NewActionHandler(mockService)

	expectedReferralIndex := map[int]int{1: 3, 2: 0, 3: 7}
	mockService.On("GetReferralIndex").Return(expectedReferralIndex, nil)

	router := mux.NewRouter()
	router.HandleFunc("/users/referral-index", actionHandler.GetReferralIndex)

	req := httptest.NewRequest("GET", "/users/referral-index", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	mockService.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"1":3`)
	assert.Contains(t, w.Body.String(), `"2":0`)
	assert.Contains(t, w.Body.String(), `"3":7`)
}
