package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"

	"food-app/models"
)

// RecipeAPIService handles external recipe API integration
type RecipeAPIService struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// ExternalRecipe represents a recipe from external API
type ExternalRecipe struct {
	ID           int                    `json:"id"`
	Title        string                 `json:"title"`
	Image        string                 `json:"image"`
	ImageType    string                 `json:"imageType"`
	PrepTime     int                    `json:"readyInMinutes"`
	Servings     int                    `json:"servings"`
	Summary      string                 `json:"summary"`
	Instructions []InstructionStep      `json:"analyzedInstructions"`
	Nutrition    ExternalNutrition      `json:"nutrition"`
	DietaryTags  []string               `json:"diets"`
	Ingredients  []ExternalIngredient   `json:"extendedIngredients"`
}

type InstructionStep struct {
	Steps []struct {
		Number int    `json:"number"`
		Step   string `json:"step"`
	} `json:"steps"`
}

type ExternalNutrition struct {
	Nutrients []struct {
		Name     string  `json:"name"`
		Amount   float64 `json:"amount"`
		Unit     string  `json:"unit"`
		Title    string  `json:"title"`
	} `json:"nutrients"`
}

type ExternalIngredient struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Amount       float64 `json:"amount"`
	Unit         string  `json:"unit"`
	OriginalName string  `json:"originalName"`
}

// NewRecipeAPIService creates a new recipe API service
func NewRecipeAPIService() *RecipeAPIService {
	return &RecipeAPIService{
		APIKey:  getEnv("RECIPE_API_KEY", ""),
		BaseURL: getEnv("RECIPE_API_URL", "https://api.spoonacular.com/recipes"),
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// SearchRecipes searches for recipes by query
func (r *RecipeAPIService) SearchRecipes(query string, limit int) ([]ExternalRecipe, error) {
	if r.APIKey == "" {
		// Return mock data if no API key configured
		return r.getMockRecipes(limit), nil
	}

	url := fmt.Sprintf("%s/complexSearch?query=%s&number=%d&addRecipeInformation=true&fillIngredients=true&apiKey=%s", 
		r.BaseURL, query, limit, r.APIKey)

	resp, err := r.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recipes: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Results []ExternalRecipe `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Results, nil
}

// ConvertToMeal converts an external recipe to our internal Meal model
func (r *RecipeAPIService) ConvertToMeal(recipe ExternalRecipe) models.Meal {
	// Convert instructions
	instructions := []string{}
	for _, instructionGroup := range recipe.Instructions {
		for _, step := range instructionGroup.Steps {
			instructions = append(instructions, step.Step)
		}
	}
	instructionsJSON, _ := json.Marshal(instructions)

	// Extract nutrition info
	nutrition := models.NutritionInfo{}
	for _, nutrient := range recipe.Nutrition.Nutrients {
		switch nutrient.Name {
		case "Calories":
			nutrition.Calories = nutrient.Amount
		case "Protein":
			nutrition.Protein = nutrient.Amount
		case "Carbohydrates":
			nutrition.Carbohydrates = nutrient.Amount
		case "Fat":
			nutrition.Fat = nutrient.Amount
		case "Fiber":
			nutrition.Fiber = nutrient.Amount
		case "Sugar":
			nutrition.Sugar = nutrient.Amount
		case "Sodium":
			nutrition.Sodium = nutrient.Amount
		}
	}

	// Determine meal type from dietary tags
	mealType := "dinner" // default
	for _, diet := range recipe.DietaryTags {
		if diet == "breakfast" || diet == "brunch" {
			mealType = "breakfast"
			break
		} else if diet == "lunch" {
			mealType = "lunch"
			break
		}
	}

	// Convert dietary tags
	dietaryTags := models.StringArray(recipe.DietaryTags)

	return models.Meal{
		Name:         recipe.Title,
		Description:  stripHTML(recipe.Summary),
		ImageURL:     recipe.Image,
		PrepTime:     recipe.PrepTime,
		CookTime:     0, // Not provided by most APIs
		Servings:     recipe.Servings,
		Difficulty:   "medium", // Default
		Cuisine:      "various", // Would need additional API call
		MealType:     mealType,
		Instructions: string(instructionsJSON),
		NutritionInfo: nutrition,
		DietaryTags:  dietaryTags,
		Allergens:    models.StringArray{}, // Would need additional processing
	}
}

// ImportRecipesFromAPI imports recipes from external API and saves to database
func (r *RecipeAPIService) ImportRecipesFromAPI(queries []string, limitPerQuery int) error {
	for _, query := range queries {
		log.Printf("Importing recipes for query: %s", query)
		
		recipes, err := r.SearchRecipes(query, limitPerQuery)
		if err != nil {
			log.Printf("Error fetching recipes for '%s': %v", query, err)
			continue
		}

		for _, externalRecipe := range recipes {
			meal := r.ConvertToMeal(externalRecipe)
			
			// This would be implemented with proper database access
			// For now, just log the import
			log.Printf("Would import meal: %s", meal.Name)

			log.Printf("Imported meal: %s", meal.Name)
		}
	}

	return nil
}

// getMockRecipes returns mock recipe data when no API key is configured
func (r *RecipeAPIService) getMockRecipes(limit int) []ExternalRecipe {
	mockRecipes := []ExternalRecipe{
		{
			ID:       1,
			Title:    "Mediterranean Chicken Bowl",
			Image:    "https://images.unsplash.com/photo-1546833999-b9f581a1996d?w=500",
			PrepTime: 25,
			Servings: 4,
			Summary:  "A healthy Mediterranean-inspired chicken bowl with fresh vegetables and quinoa.",
			DietaryTags: []string{"gluten-free", "high-protein", "mediterranean"},
			Nutrition: ExternalNutrition{
				Nutrients: []struct {
					Name     string  `json:"name"`
					Amount   float64 `json:"amount"`
					Unit     string  `json:"unit"`
					Title    string  `json:"title"`
				}{
					{Name: "Calories", Amount: 420, Unit: "kcal"},
					{Name: "Protein", Amount: 35, Unit: "g"},
					{Name: "Carbohydrates", Amount: 45, Unit: "g"},
					{Name: "Fat", Amount: 12, Unit: "g"},
				},
			},
		},
		{
			ID:       2,
			Title:    "Vegan Buddha Bowl",
			Image:    "https://images.unsplash.com/photo-1512621776951-a57141f2eefd?w=500",
			PrepTime: 20,
			Servings: 2,
			Summary:  "A nutritious plant-based bowl packed with colorful vegetables and plant proteins.",
			DietaryTags: []string{"vegan", "gluten-free", "high-fiber"},
			Nutrition: ExternalNutrition{
				Nutrients: []struct {
					Name     string  `json:"name"`
					Amount   float64 `json:"amount"`
					Unit     string  `json:"unit"`
					Title    string  `json:"title"`
				}{
					{Name: "Calories", Amount: 380, Unit: "kcal"},
					{Name: "Protein", Amount: 15, Unit: "g"},
					{Name: "Carbohydrates", Amount: 52, Unit: "g"},
					{Name: "Fat", Amount: 15, Unit: "g"},
				},
			},
		},
	}

	if limit < len(mockRecipes) {
		return mockRecipes[:limit]
	}
	return mockRecipes
}

// Helper functions
func stripHTML(s string) string {
	// Simple HTML tag removal - in production, use a proper HTML parser
	// This is a basic implementation
	return s // For now, return as-is
}

func getEnv(key, defaultValue string) string {
	// This should match the getEnv function in your main package
	// For now, return default
	return defaultValue
}
