package repositories

import "backend-coding-challenge-enhanced/internal/models"

// UserRepositoryInterface defines methods for interacting with user data.
type UserRepositoryInterface interface {
	FetchUserByID(id int) (
		models.User,
		error,
	)
	FetchUserActionCount(userID int) (int, error)
}

// ActionRepositoryInterface defines methods for interacting with action data.
type ActionRepositoryInterface interface {
	FetchNextActionProbabilities(actionType string) (map[string]float64, error)
	FetchReferralIndex() (map[int]int, error)
}
