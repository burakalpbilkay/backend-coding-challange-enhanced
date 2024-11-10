package services_test

import (
	"backend-coding-challenge-enhanced/internal/models"
	"backend-coding-challenge-enhanced/internal/repositories"
	"backend-coding-challenge-enhanced/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FetchUserByID(id int) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) FetchUserActionCount(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func TestGetUserActionCount_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	mockUserID := 1
	expectedCount := 10

	mockRepo.On("FetchUserActionCount", mockUserID).Return(expectedCount, nil)

	count, err := userService.GetUserActionCount(mockUserID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}

func TestGetUserActionCount_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	mockUserID := 999

	mockRepo.On("FetchUserActionCount", mockUserID).Return(0, repositories.ErrUserNotFound)

	count, err := userService.GetUserActionCount(mockUserID)

	assert.Equal(t, repositories.ErrUserNotFound, err)
	assert.Equal(t, 0, count)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	mockUserID := 1
	expectedUser := models.User{ID: mockUserID, Name: "John Doe", CreatedAt: "2023-10-06T11:12:22.758Z"}

	mockRepo.On("FetchUserByID", mockUserID).Return(expectedUser, nil)

	user, err := userService.GetUserByID(mockUserID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	mockUserID := 999
	mockRepo.On("FetchUserByID", mockUserID).Return(models.User{}, repositories.ErrUserNotFound)

	user, err := userService.GetUserByID(mockUserID)

	assert.Error(t, err)
	assert.Equal(t, repositories.ErrUserNotFound, err)
	assert.Equal(t, models.User{}, user)
	mockRepo.AssertExpectations(t)
}
