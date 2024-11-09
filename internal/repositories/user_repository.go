package repositories

import (
	"backend-coding-challenge-enhanced/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type UserRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

var ErrUserNotFound = errors.New("user not found")

// Initialize a new UserRepository with a DB connection.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Set the Redis client in the UserRepository.
func (r *UserRepository) SetRedis(rdb *redis.Client) {
	r.rdb = rdb
}

// Retrieve a user by ID from the database or Redis cache if available.
func (r *UserRepository) FetchUserByID(id int) (models.User, error) {
	var user models.User

	// Get the cached result from Redis
	cacheKey := "user_id:" + strconv.Itoa(id)
	cachedUser, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		log.Println("Cache miss, fetching from database")
	} else if err != nil {
		return user, fmt.Errorf("redis error: %v", err)
	} else {
		log.Println("Cache hit, returning user from Redis")
		json.Unmarshal([]byte(cachedUser), &user)
		return user, nil
	}

	// Go to DB if Redis could not find the result
	row := r.db.QueryRow("SELECT id, name, created_at FROM users WHERE id=$1", id)
	err = row.Scan(&user.ID, &user.Name, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	// Cache the result for a min
	userData, _ := json.Marshal(user)
	r.rdb.Set(ctx, fmt.Sprintf("user:%d", id), userData, 60*time.Second)

	return user, nil
}

// Retrieve the total number of actions performed by a user from the database or Redis cache if available.
func (r *UserRepository) FetchUserActionCount(userID int) (int, error) {
	cacheKey := "user_action_count:" + strconv.Itoa(userID)

	// Get the cached result from Redis
	cachedCount, err := r.rdb.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil {
		log.Println("Cache miss, fetching from database")
	} else if err == nil {
		log.Println("Cache hit, returning user from Redis")
		if count, parseErr := strconv.Atoi(cachedCount); parseErr == nil {
			return count, nil
		}
	}

	// Go to DB if Redis could not find the result
	// Validate that user exists
	var exists bool
	err = r.db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE id=$1)", userID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("error checking user existence: %v", err)
	}
	if !exists {
		return 0, ErrUserNotFound
	}

	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM actions WHERE user_id=$1", userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	// Cache the result for a min
	r.rdb.Set(ctx, cacheKey, strconv.Itoa(count), 60*time.Second)
	return count, nil
}
