package database

import (
	"fmt"
	"log"
	"os"

	"food-app/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func Connect() {
	var err error
	
	// Check if we should use SQLite for development
	if getEnv("DB_TYPE", "sqlite") == "sqlite" {
		dbPath := getEnv("DB_PATH", "./food_app.db")
		DB, err = gorm.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal("Failed to connect to SQLite database:", err)
		}
		log.Println("SQLite database connection established")
	} else {
		// PostgreSQL configuration
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "password")
		dbname := getEnv("DB_NAME", "food_app")
		sslmode := getEnv("DB_SSLMODE", "disable")

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)

		DB, err = gorm.Open("postgres", dsn)
		if err != nil {
			log.Fatal("Failed to connect to PostgreSQL database:", err)
		}
		log.Println("PostgreSQL database connection established")
	}
}

func Migrate() {
	// Auto-migrate all models
	DB.AutoMigrate(
		&models.User{},
		&models.Meal{},
		&models.Ingredient{},
		&models.MealIngredient{},
		&models.UserMealInteraction{},
		&models.MealReview{},
		&models.MealPlan{},
		&models.MealPlanEntry{},
		&models.CurrentMealPlan{},
		&models.ShoppingList{},
		&models.ShoppingListItem{},
	)

	log.Println("Database migration completed")
}

func SeedData() {
	// Check if data already exists
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Database already contains data, skipping seed")
		return
	}

	// Seed ingredients
	ingredients := []models.Ingredient{
		{Name: "Chicken Breast", Category: "protein", Unit: "piece", CaloriesPer100g: 165},
		{Name: "Rice", Category: "grain", Unit: "cup", CaloriesPer100g: 130},
		{Name: "Broccoli", Category: "vegetable", Unit: "cup", CaloriesPer100g: 34},
		{Name: "Salmon", Category: "protein", Unit: "fillet", CaloriesPer100g: 208},
		{Name: "Sweet Potato", Category: "vegetable", Unit: "piece", CaloriesPer100g: 86},
		{Name: "Spinach", Category: "vegetable", Unit: "cup", CaloriesPer100g: 23},
		{Name: "Quinoa", Category: "grain", Unit: "cup", CaloriesPer100g: 222},
		{Name: "Olive Oil", Category: "fat", Unit: "tbsp", CaloriesPer100g: 884},
		{Name: "Garlic", Category: "seasoning", Unit: "clove", CaloriesPer100g: 149},
		{Name: "Onion", Category: "vegetable", Unit: "piece", CaloriesPer100g: 40},
		{Name: "Tomato", Category: "vegetable", Unit: "piece", CaloriesPer100g: 18},
		{Name: "Bell Pepper", Category: "vegetable", Unit: "piece", CaloriesPer100g: 31},
		{Name: "Black Beans", Category: "protein", Unit: "cup", CaloriesPer100g: 132},
		{Name: "Avocado", Category: "fat", Unit: "piece", CaloriesPer100g: 160},
		{Name: "Lemon", Category: "fruit", Unit: "piece", CaloriesPer100g: 29},
	}

	for _, ingredient := range ingredients {
		var existing models.Ingredient
		if DB.Where("name = ?", ingredient.Name).First(&existing).RecordNotFound() {
			DB.Create(&ingredient)
		}
	}

	// Seed meals
	meals := []models.Meal{
		{
			Name:        "Grilled Chicken with Rice and Broccoli",
			Description: "A healthy and balanced meal with lean protein, complex carbs, and vegetables",
			ImageURL:    "https://images.unsplash.com/photo-1546833999-b9f581a1996d?w=500",
			PrepTime:    15,
			CookTime:    25,
			Servings:    4,
			Difficulty:  "easy",
			Cuisine:     "American",
			MealType:    "dinner",
			Instructions: `["Season chicken breast with salt and pepper","Heat grill pan over medium-high heat","Cook chicken for 6-7 minutes per side","Steam broccoli for 5 minutes","Cook rice according to package instructions","Serve together with a drizzle of olive oil"]`,
			NutritionInfo: models.NutritionInfo{
				Calories:      450,
				Protein:       35,
				Carbohydrates: 45,
				Fat:           12,
				Fiber:         4,
				Sugar:         3,
				Sodium:        320,
			},
			DietaryTags: models.StringArray{"gluten-free", "high-protein"},
			Allergens:   models.StringArray{},
		},
		{
			Name:        "Quinoa Buddha Bowl",
			Description: "A nutritious vegetarian bowl packed with quinoa, vegetables, and healthy fats",
			ImageURL:    "https://images.unsplash.com/photo-1512621776951-a57141f2eefd?w=500",
			PrepTime:    20,
			CookTime:    15,
			Servings:    2,
			Difficulty:  "easy",
			Cuisine:     "Mediterranean",
			MealType:    "lunch",
			Instructions: `["Cook quinoa according to package instructions","Roast sweet potato cubes at 400°F for 20 minutes","Sauté spinach with garlic","Slice avocado","Arrange all ingredients in a bowl","Drizzle with olive oil and lemon juice"]`,
			NutritionInfo: models.NutritionInfo{
				Calories:      380,
				Protein:       12,
				Carbohydrates: 52,
				Fat:           15,
				Fiber:         8,
				Sugar:         8,
				Sodium:        180,
			},
			DietaryTags: models.StringArray{"vegan", "gluten-free", "high-fiber"},
			Allergens:   models.StringArray{},
		},
		{
			Name:        "Baked Salmon with Sweet Potato",
			Description: "Heart-healthy salmon with roasted sweet potato and vegetables",
			ImageURL:    "https://images.unsplash.com/photo-1467003909585-2f8a72700288?w=500",
			PrepTime:    10,
			CookTime:    20,
			Servings:    4,
			Difficulty:  "medium",
			Cuisine:     "Mediterranean",
			MealType:    "dinner",
			Instructions: `["Preheat oven to 425°F","Season salmon with herbs and lemon","Cut sweet potato into wedges","Toss vegetables with olive oil","Bake salmon for 12-15 minutes","Roast vegetables for 25 minutes"]`,
			NutritionInfo: models.NutritionInfo{
				Calories:      420,
				Protein:       28,
				Carbohydrates: 35,
				Fat:           18,
				Fiber:         6,
				Sugar:         12,
				Sodium:        280,
			},
			DietaryTags: models.StringArray{"gluten-free", "omega-3", "heart-healthy"},
			Allergens:   models.StringArray{"fish"},
		},
	}

			for _, meal := range meals {
			var existing models.Meal
			if DB.Where("name = ?", meal.Name).First(&existing).RecordNotFound() {
				if err := DB.Create(&meal).Error; err != nil {
					log.Printf("Error creating meal %s: %v", meal.Name, err)
					continue
				}

				// Add meal ingredients with quantities
				seedMealIngredients(meal.ID, meal.Name)
			}
		}

	log.Println("Database seeded with initial data")
}

func seedMealIngredients(mealID uint, mealName string) {
	// Define ingredient mappings for each meal
	mealIngredientMappings := map[string][]struct {
		name     string
		quantity float64
		unit     string
	}{
		"Grilled Chicken with Rice and Broccoli": {
			{"Chicken Breast", 1, "lb"},
			{"White Rice", 1, "cup"},
			{"Broccoli", 1, "head"},
			{"Olive Oil", 2, "tbsp"},
		},
		"Quinoa Buddha Bowl": {
			{"Quinoa", 1, "cup"},
			{"Sweet Potato", 1, "medium"},
			{"Spinach", 2, "cups"},
			{"Avocado", 1, "medium"},
			{"Olive Oil", 1, "tbsp"},
		},
		"Salmon with Vegetables": {
			{"Salmon Fillet", 6, "oz"},
			{"Asparagus", 1, "bunch"},
			{"Lemon", 1, "medium"},
			{"Olive Oil", 1, "tbsp"},
		},
	}

	ingredients, exists := mealIngredientMappings[mealName]
	if !exists {
		return // Skip if no mapping defined
	}

	for _, ing := range ingredients {
		var ingredient models.Ingredient
		if !DB.Where("name = ?", ing.name).First(&ingredient).RecordNotFound() {
			// Create meal ingredient relationship
			mealIngredient := models.MealIngredient{
				MealID:       mealID,
				IngredientID: ingredient.ID,
				Quantity:     ing.quantity,
				Unit:         ing.unit,
			}
			
			if err := DB.Create(&mealIngredient).Error; err != nil {
				log.Printf("Error creating meal ingredient for %s: %v", ing.name, err)
			}
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
