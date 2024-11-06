package repositories

import (
	"backend-coding-challenge-enhanced/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type UserRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

// NewUserRepository initializes a new UserRepository with a DB connection.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// SetRedis sets the Redis client in the UserRepository.
func (r *UserRepository) SetRedis(rdb *redis.Client) {
	r.rdb = rdb
}

// FetchUserByID retrieves a user by ID from the database or Redis cache if available.
func (r *UserRepository) FetchUserByID(id int) (models.User, error) {
	var user models.User

	// Check cache
	cachedUser, err := r.rdb.Get(ctx, fmt.Sprintf("user:%d", id)).Result()
	if err == redis.Nil {
		log.Println("Cache miss, fetching from database")
	} else if err != nil {
		return user, fmt.Errorf("redis error: %v", err)
	} else {
		log.Println("Cache hit, returning user from Redis")
		json.Unmarshal([]byte(cachedUser), &user)
		return user, nil
	}

	// Fetch from PostgreSQL if not in cache
	row := r.db.QueryRow("SELECT id, name, created_at FROM users WHERE id=$1", id)
	err = row.Scan(&user.ID, &user.Name, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	// Cache the result for future requests
	userData, _ := json.Marshal(user)
	r.rdb.Set(ctx, fmt.Sprintf("user:%d", id), userData, 0)

	return user, nil
}

// FetchUserActionCount retrieves the total number of actions performed by a user.
func (r *UserRepository) FetchUserActionCount(userID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM actions WHERE user_id=$1", userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
