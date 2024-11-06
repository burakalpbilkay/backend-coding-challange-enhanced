package repositories

type MockActionRepo struct{}

func (m *MockActionRepo) FetchNextActionProbabilities(actionType string) (map[string]float64, error) {
	// Return mock action probability data for testing purposes
	return map[string]float64{
		"ADD_TO_CRM":        0.70,
		"REFER_USER":        0.20,
		"VIEW_CONVERSATION": 0.10,
	}, nil
}

func (m *MockActionRepo) FetchReferralIndex() (map[int]int, error) {
	// Return mock referral index data for testing purposes
	return map[int]int{
		1: 3,
		2: 0,
		3: 7,
	}, nil
}
