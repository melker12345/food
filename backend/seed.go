package main

import (
	"context"
	"log"
	"time"
)

func seedDatabase() {
	// Check if meals already exist
	count, err := database.Collection("meals").CountDocuments(context.Background(), map[string]interface{}{})
	if err != nil {
		log.Printf("Error checking meal count: %v", err)
		return
	}

	if count > 0 {
		log.Printf("Seed skipped: %d meals already exist", count)
		return
	}

	// Sample meals
	meals := []Meal{
		{
			Name:     "Grilled Chicken Salad",
			ImageURL: "https://images.unsplash.com/photo-1568605114967-8130f3a36994",
			Ingredients: []Ingredient{
				{Name: "Chicken Breast", Quantity: "200 g"},
				{Name: "Mixed Greens", Quantity: "3 cups"},
				{Name: "Cherry Tomatoes", Quantity: "1 cup"},
				{Name: "Olive Oil", Quantity: "1 tbsp"},
			},
			Instructions: "Grill chicken, toss with greens and tomatoes, drizzle olive oil.",
			Nutrition:    Nutrition{Calories: 420, Protein: 40, Carbs: 12, Fat: 24},
			DietaryTags:  []string{"gluten-free"},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:     "Vegan Buddha Bowl",
			ImageURL: "https://images.unsplash.com/photo-1512621776951-a57141f2eefd",
			Ingredients: []Ingredient{
				{Name: "Quinoa", Quantity: "1 cup cooked"},
				{Name: "Roasted Chickpeas", Quantity: "1/2 cup"},
				{Name: "Avocado", Quantity: "1/2"},
				{Name: "Spinach", Quantity: "2 cups"},
			},
			Instructions: "Combine ingredients in a bowl, season to taste.",
			Nutrition:    Nutrition{Calories: 520, Protein: 18, Carbs: 60, Fat: 22},
			DietaryTags:  []string{"vegan", "gluten-free"},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Name:     "Pasta Primavera",
			ImageURL: "https://images.unsplash.com/photo-1473093295043-cdd812d0e601",
			Ingredients: []Ingredient{
				{Name: "Pasta", Quantity: "200 g"},
				{Name: "Mixed Vegetables", Quantity: "2 cups"},
				{Name: "Parmesan", Quantity: "2 tbsp"},
			},
			Instructions: "Cook pasta and sauté vegetables, toss together and top with cheese.",
			Nutrition:    Nutrition{Calories: 600, Protein: 20, Carbs: 90, Fat: 16},
			DietaryTags:  []string{"vegetarian"},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	// Convert to interface{} slice for insertion
	docs := make([]interface{}, len(meals))
	for i, meal := range meals {
		docs[i] = meal
	}

	_, err = database.Collection("meals").InsertMany(context.Background(), docs)
	if err != nil {
		log.Printf("Error seeding meals: %v", err)
		return
	}

	log.Println("Seed completed: inserted sample meals")
}
