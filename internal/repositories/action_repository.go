package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ActionProbabilitiesCacheKey = "action_probabilities:"
	ReferralIndexCacheKey       = "referral_index"
)

type ActionRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

// Initialize a new ActionRepository with a DB connection.
func NewActionRepository(db *sql.DB) *ActionRepository {
	return &ActionRepository{db: db}
}

// Set the Redis client in the ActionRepository.
func (r *ActionRepository) SetRedis(rdb *redis.Client) {
	r.rdb = rdb
}

// Calculate the probability of each action following a given action type
func (r *ActionRepository) FetchNextActionProbabilities(actionType string) (map[string]float64, error) {

	cacheKey := ActionProbabilitiesCacheKey + actionType

	// Get the cached result from Redis
	cachedData, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		log.Println("Cache miss, fetching from database")
	} else if err == nil {
		log.Println("Cache hit, returning user from Redis")
		var probabilities map[string]float64
		if json.Unmarshal([]byte(cachedData), &probabilities) == nil {
			return probabilities, nil
		}
	}

	// Go to DB if Redis could not find the result
	probabilities := make(map[string]float64)
	totalCount := 0

	query := `
	SELECT next_action, COUNT(*)
	FROM (
		SELECT 
			type AS current_action,
			LEAD(type) OVER (PARTITION BY user_id ORDER BY created_at) AS next_action
		FROM actions
	) AS sequential_actions
	WHERE current_action = $1
	AND next_action IS NOT NULL
	GROUP BY next_action;
	
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

	// Cache the result for a min
	data, _ := json.Marshal(probabilities)
	r.rdb.Set(ctx, cacheKey, data, 60*time.Second)
	return probabilities, nil
}

// Calculate the referral index for each user
func (r *ActionRepository) FetchReferralIndex() (map[int]int, error) {

	cacheKey := ReferralIndexCacheKey

	// Get the cached result from Redis
	cachedData, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		log.Println("Cache miss, fetching from database")
	} else if err == nil {
		log.Println("Cache hit, returning user from Redis")
		var referralIndex map[int]int
		if json.Unmarshal([]byte(cachedData), &referralIndex) == nil {
			return referralIndex, nil
		}
	}

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

	// Cache the result for a min
	data, _ := json.Marshal(referralCountCache)
	r.rdb.Set(ctx, cacheKey, data, 60*time.Second)

	return referralCountCache, nil
}
