package main

import (
	"backend-coding-challenge-enhanced/internal/handlers"
	"backend-coding-challenge-enhanced/internal/middleware"
	"backend-coding-challenge-enhanced/internal/repositories"
	"backend-coding-challenge-enhanced/internal/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize Redis
	rdb := repositories.InitRedis()

	// Initialize PostgreSQL database
	db := repositories.InitDB()
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	userRepo.SetRedis(rdb)
	actionRepo := repositories.NewActionRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	actionService := services.NewActionService(actionRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	actionHandler := handlers.NewActionHandler(actionService)

	// Initialize rate limit middleware
	rateLimiter := middleware.RateLimit(rdb, 100, time.Minute)

	// Set up router and routes
	router := mux.NewRouter()

	// User related endpoints
	router.Handle("/user/{id}", rateLimiter(http.HandlerFunc(userHandler.GetUserByID))).Methods("GET")
	router.Handle("/user/{id}/actions/count", rateLimiter(http.HandlerFunc(userHandler.GetUserActionCount))).Methods("GET")

	// Action related endpoints
	router.Handle("/action/{type}/next", rateLimiter(http.HandlerFunc(actionHandler.GetNextActionProbabilities))).Methods("GET")
	router.Handle("/users/referral-index", rateLimiter(http.HandlerFunc(actionHandler.GetReferralIndex))).Methods("GET")

	// Server setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
