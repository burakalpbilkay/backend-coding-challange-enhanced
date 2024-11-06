package repositories

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	connStr := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") +
		"@" + os.Getenv("POSTGRES_HOST") + "/" + os.Getenv("POSTGRES_DB") + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func InitRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	return redis.NewClient(&redis.Options{
		Addr: redisHost,
	})
}
