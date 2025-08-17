package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var database *mongo.Database

func main() {
	// Load environment variables
	port := getEnv("PORT", "4000")
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017/food")
	corsOrigin := getEnv("CORS_ORIGIN", "*")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	mongoClient = client
	database = client.Database("food")
	log.Println("Connected to MongoDB")

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	if corsOrigin == "*" {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = []string{corsOrigin}
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Routes
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// User routes
		api.POST("/users/upsert", upsertUser)
		api.GET("/users/me", getMe)

		// Meal routes
		api.POST("/meals", createMeal)
		api.GET("/meals/random", getRandomMeal)
		api.POST("/meals/like", likeMeal)

		// Plan routes
		api.POST("/plans/generate", generateWeeklyPlan)
		api.GET("/plans/weekly", getWeeklyPlan)

		// Shopping routes
		api.GET("/shopping/list", getShoppingList)
	}

	// Seed database on startup
	seedDatabase()

	log.Printf("API listening on http://localhost:%s", port)
	r.Run(":" + port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
