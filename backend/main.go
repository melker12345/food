package main

import (
	"food-app/database"
	"food-app/handlers"
	"food-app/middleware"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// getEnv gets environment variable with fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Initialize database
	database.Connect()
	database.Migrate()
	database.SeedData()

	// Create Gin router
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	corsOrigins := getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001,http://localhost:5173,http://localhost:5174,http://127.0.0.1:3000,http://127.0.0.1:3001,http://127.0.0.1:5173,http://127.0.0.1:5174")
	if corsOrigins == "*" {
		config.AllowAllOrigins = true
	} else {
		origins := strings.Split(corsOrigins, ",")
		for i, origin := range origins {
			origins[i] = strings.TrimSpace(origin)
		}
		config.AllowOrigins = origins
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Food Planning API is running"})
	})

	// API routes
	api := r.Group("/api/v1")

	// Public routes
	public := api.Group("/")
	{
		// Authentication
		public.POST("/register", handlers.Register)
		public.POST("/login", handlers.Login)

		// Public meal browsing
		public.GET("/meals", handlers.GetMeals)
		public.GET("/meals/:id", handlers.GetMeal)
		public.GET("/meals/trending", handlers.GetTrendingMeals)
		public.GET("/meals/:id/reviews", handlers.GetMealReviews)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// User profile
		protected.GET("/profile", handlers.GetProfile)
		protected.PUT("/profile", handlers.UpdateProfile)
		protected.PUT("/profile/preferences", handlers.UpdatePreferences)

		// Personalized meals
		protected.GET("/meals/personalized", handlers.GetPersonalizedMeals)
		protected.GET("/meals/liked", handlers.GetLikedMeals)
		protected.POST("/meals/:id/like", handlers.LikeMeal)
		protected.POST("/meals/:id/dislike", handlers.DislikeMeal)
		protected.POST("/meals/:id/reviews", handlers.AddMealReview)

		// Current Meal Plan (Single Plan Approach)
		protected.GET("/current-meal-plan", handlers.GetCurrentMealPlan)
		protected.POST("/current-meal-plan/populate-from-liked", handlers.PopulateFromLikedMeals)
		protected.PUT("/current-meal-plan/meals", handlers.UpdateMealInPlan)
		protected.PUT("/shopping-items/:item_id", handlers.ToggleShoppingItem)

		// Legacy Meal planning (Multiple Plans)
		protected.POST("/meal-plans", handlers.CreateMealPlan)
		protected.POST("/meal-plans/auto-generate", handlers.AutoGenerateMealPlan)
		protected.GET("/meal-plans", handlers.GetMealPlans)
		protected.GET("/meal-plans/:id", handlers.GetMealPlan)
		protected.PUT("/meal-plans/:id", handlers.UpdateMealPlan)
		protected.DELETE("/meal-plans/:id", handlers.DeleteMealPlan)

		// Shopping lists
		protected.POST("/meal-plans/:id/shopping-list", handlers.GenerateShoppingList)
		protected.GET("/shopping-lists", handlers.GetShoppingLists)
		protected.PUT("/shopping-list-items/:item_id", handlers.UpdateShoppingListItem)
	}

	// Start server
	port := "8080"
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
