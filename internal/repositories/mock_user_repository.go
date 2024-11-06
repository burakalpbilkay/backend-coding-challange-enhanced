package repositories

import "backend-coding-challenge-enhanced/internal/models"

type MockUserRepo struct{}

func (m *MockUserRepo) FetchUserByID(id int) (models.User, error) {
	return models.User{ID: id, Name: "Burak Bilkay", CreatedAt: "2024-10-06T11:12:22.758Z"}, nil
}

func (m *MockUserRepo) FetchUserActionCount(userID int) (int, error) {
	return 100, nil
}
