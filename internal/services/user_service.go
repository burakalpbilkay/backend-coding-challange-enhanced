package services

import (
	"backend-coding-challenge-enhanced/internal/models"
	"backend-coding-challenge-enhanced/internal/repositories"
)

type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(userID int) (models.User, error) {
	return s.userRepo.FetchUserByID(userID)
}

func (s *UserService) GetUserActionCount(userID int) (int, error) {
	return s.userRepo.FetchUserActionCount(userID)
}
