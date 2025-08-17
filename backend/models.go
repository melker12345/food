package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user document
type User struct {
	ID                  primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	AuthProviderID      string               `bson:"authProviderId" json:"authProviderId"`
	Name                string               `bson:"name" json:"name"`
	Email               string               `bson:"email" json:"email"`
	DietaryPreferences  []string             `bson:"dietaryPreferences" json:"dietaryPreferences"`
	HealthGoals         []string             `bson:"healthGoals" json:"healthGoals"`
	Goal                string               `bson:"goal" json:"goal"` // maintenance, cutting, bulking
	LikedMeals          []primitive.ObjectID `bson:"likedMeals" json:"likedMeals"`
	CreatedAt           time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time            `bson:"updatedAt" json:"updatedAt"`
}

// Ingredient represents an ingredient item
type Ingredient struct {
	Name     string `bson:"name" json:"name"`
	Quantity string `bson:"quantity" json:"quantity"`
}

// Nutrition represents nutritional information
type Nutrition struct {
	Calories int `bson:"calories" json:"calories"`
	Protein  int `bson:"protein" json:"protein"`
	Carbs    int `bson:"carbs" json:"carbs"`
	Fat      int `bson:"fat" json:"fat"`
}

// Meal represents a meal document
type Meal struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	ImageURL     string             `bson:"imageUrl" json:"imageUrl"`
	Ingredients  []Ingredient       `bson:"ingredients" json:"ingredients"`
	Instructions string             `bson:"instructions" json:"instructions"`
	Nutrition    Nutrition          `bson:"nutrition" json:"nutrition"`
	DietaryTags  []string           `bson:"dietaryTags" json:"dietaryTags"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// DailyMealEntry represents a meal entry for a specific day
type DailyMealEntry struct {
	Meal     primitive.ObjectID `bson:"meal" json:"meal"`
	MealType string             `bson:"mealType" json:"mealType"` // breakfast, lunch, dinner, snack
}

// DayPlan represents a single day's meal plan
type DayPlan struct {
	Date  time.Time        `bson:"date" json:"date"`
	Meals []DailyMealEntry `bson:"meals" json:"meals"`
}

// WeeklyPlan represents a weekly meal plan
type WeeklyPlan struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	User          primitive.ObjectID `bson:"user" json:"user"`
	WeekStartDate time.Time          `bson:"weekStartDate" json:"weekStartDate"`
	Days          []DayPlan          `bson:"days" json:"days"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Request/Response types
type UpsertUserRequest struct {
	AuthProviderID     string   `json:"authProviderId" binding:"required"`
	Name               string   `json:"name" binding:"required"`
	Email              string   `json:"email" binding:"required,email"`
	DietaryPreferences []string `json:"dietaryPreferences"`
	HealthGoals        []string `json:"healthGoals"`
	Goal               string   `json:"goal"`
}

type CreateMealRequest struct {
	Name         string       `json:"name" binding:"required"`
	ImageURL     string       `json:"imageUrl" binding:"required,url"`
	Ingredients  []Ingredient `json:"ingredients" binding:"required"`
	Instructions string       `json:"instructions" binding:"required"`
	Nutrition    Nutrition    `json:"nutrition" binding:"required"`
	DietaryTags  []string     `json:"dietaryTags"`
}

type LikeMealRequest struct {
	AuthProviderID string `json:"authProviderId" binding:"required"`
	MealID         string `json:"mealId" binding:"required"`
}

type GeneratePlanRequest struct {
	AuthProviderID string `json:"authProviderId" binding:"required"`
}

type ShoppingListItem struct {
	Name       string   `json:"name"`
	Quantities []string `json:"quantities"`
}

type ShoppingListResponse struct {
	WeekStartDate time.Time          `json:"weekStartDate"`
	Items         []ShoppingListItem `json:"items"`
}
