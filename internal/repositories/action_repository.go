package repositories

import (
	"backend-coding-challenge-enhanced/internal/constants"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type ActionRepository struct {
	db *sql.DB
}

var ErrInvalidActionType = errors.New("invalid action type")

// Initialize a new ActionRepository with a DB connection.
func NewActionRepository(db *sql.DB) *ActionRepository {
	return &ActionRepository{db: db}
}

// Calculate the probability of each action following a given action type
func (r *ActionRepository) FetchNextActionProbabilities(actionType string) (map[string]float64, error) {

	// Validate that actionType is valid
	if !constants.ValidActionTypes[constants.ActionType(actionType)] {
		return nil, ErrInvalidActionType
	}

	probabilities := make(map[string]float64)
	totalCount := 0

	query := `
		SELECT next_action.type AS next_action, COUNT(*) 
		FROM actions AS current_action
		JOIN actions AS next_action
		ON current_action.user_id = next_action.user_id
		AND next_action.created_at = (
			SELECT MIN(created_at)
			FROM actions
			WHERE user_id = current_action.user_id
			AND created_at > current_action.created_at
		)
		WHERE current_action.type = $1
		GROUP BY next_action.type;
    `

	log.Println("Executing query to get action count")
	rows, err := r.db.Query(query, actionType)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	// Calculate counts
	actionCounts := make(map[string]int)
	for rows.Next() {
		var nextAction string
		var count int
		if err := rows.Scan(&nextAction, &count); err != nil {
			log.Printf("Scan error: %v", err)
			return nil, fmt.Errorf("scan error: %v", err)
		}

		actionCounts[nextAction] = count
		totalCount += count
	}

	// Return empty map if no results
	if totalCount == 0 {
		return probabilities, nil
	}

	// Calculate probabilities for each next action
	for action, count := range actionCounts {
		probabilities[action] = float64(count) / float64(totalCount)
	}

	log.Println("FetchNextActionProbabilities completed successfully")
	return probabilities, nil
}

// Calculate the referral index for each user
func (r *ActionRepository) FetchReferralIndex() (map[int]int, error) {
	referralMap := make(map[int][]int)

	// Get the data from the "actions" table
	rows, err := r.db.Query("SELECT user_id, target_user FROM actions WHERE type = 'REFER_USER'")
	if err != nil {
		return nil, fmt.Errorf("error fetching referrals: %v", err)
	}
	defer rows.Close()

	// Populate referralMap with direct referrals
	for rows.Next() {
		var userID, targetUser int
		if err := rows.Scan(&userID, &targetUser); err != nil {
			return nil, err
		}
		referralMap[userID] = append(referralMap[userID], targetUser)
	}

	referralCountCache := make(map[int]int)

	// Calculate the referral count for each user
	for userID := range referralMap {
		//Skip if already calculated for the user
		if _, exists := referralCountCache[userID]; exists {
			continue
		}

		// Using BFS approach
		// Initialize a queue
		queue := []int{userID}
		visited := make(map[int]bool)
		totalReferrals := 0

		for len(queue) > 0 {
			currentUser := queue[0]
			queue = queue[1:]

			// Skip if this user has already processed
			if visited[currentUser] {
				continue
			}
			visited[currentUser] = true

			// Add direct referrals to the count and queue
			totalReferrals += len(referralMap[currentUser])
			for _, referredUser := range referralMap[currentUser] {
				if _, calculated := referralCountCache[referredUser]; !calculated {
					queue = append(queue, referredUser)
				}
			}
		}

		referralCountCache[userID] = totalReferrals
	}

	return referralCountCache, nil
}
