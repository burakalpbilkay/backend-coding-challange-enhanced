package services_test

import (
	"backend-coding-challenge-enhanced/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockActionRepository is a mock implementation of ActionRepositoryInterface
type MockActionRepository struct {
	mock.Mock
}

func (m *MockActionRepository) FetchNextActionProbabilities(actionType string) (map[string]float64, error) {
	args := m.Called(actionType)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *MockActionRepository) FetchReferralIndex() (map[int]int, error) {
	args := m.Called()
	return args.Get(0).(map[int]int), args.Error(1)
}

func TestGetNextActionProbabilities_Success(t *testing.T) {
	mockRepo := new(MockActionRepository)
	actionService := services.NewActionService(mockRepo)

	mockActionType := "REFER_USER"
	expectedProbabilities := map[string]float64{"ADD_TO_CRM": 0.7, "VIEW_CONVERSATION": 0.3}

	mockRepo.On("FetchNextActionProbabilities", mockActionType).Return(expectedProbabilities, nil)

	probabilities, err := actionService.GetNextActionProbabilities(mockActionType)

	assert.NoError(t, err)
	assert.Equal(t, expectedProbabilities, probabilities)
	mockRepo.AssertExpectations(t)
}

func TestGetNextActionProbabilities_NotFound(t *testing.T) {
	mockRepo := new(MockActionRepository)
	actionService := services.NewActionService(mockRepo)

	mockActionType := "NON_EXISTENT_TYPE"
	mockRepo.On("FetchNextActionProbabilities", mockActionType).Return(map[string]float64{}, services.ErrInvalidActionType)

	probabilities, err := actionService.GetNextActionProbabilities(mockActionType)

	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidActionType, err)
	assert.Empty(t, probabilities)
}

func TestGetReferralIndex_Success(t *testing.T) {
	mockRepo := new(MockActionRepository)
	actionService := services.NewActionService(mockRepo)

	expectedReferralIndex := map[int]int{1: 3, 2: 0, 3: 7}
	mockRepo.On("FetchReferralIndex").Return(expectedReferralIndex, nil)

	referralIndex, err := actionService.GetReferralIndex()

	assert.NoError(t, err)
	assert.Equal(t, expectedReferralIndex, referralIndex)
	mockRepo.AssertExpectations(t)
}
