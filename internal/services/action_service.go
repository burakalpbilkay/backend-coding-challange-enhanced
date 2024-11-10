package services

import (
	"backend-coding-challenge-enhanced/internal/constants"
	"backend-coding-challenge-enhanced/internal/repositories"
	"errors"
)

type ActionService struct {
	actionRepo repositories.ActionRepositoryInterface
}

var ErrInvalidActionType = errors.New("invalid action type")

func NewActionService(actionRepo repositories.ActionRepositoryInterface) *ActionService {
	return &ActionService{actionRepo: actionRepo}
}

func (s *ActionService) GetNextActionProbabilities(actionType string) (map[string]float64, error) {

	// Validate that actionType is valid
	if !constants.ValidActionTypes[constants.ActionType(actionType)] {
		return nil, ErrInvalidActionType
	}

	return s.actionRepo.FetchNextActionProbabilities(actionType)
}

func (s *ActionService) GetReferralIndex() (map[int]int, error) {
	return s.actionRepo.FetchReferralIndex()
}
