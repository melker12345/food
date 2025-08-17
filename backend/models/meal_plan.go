package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MealPlan struct {
	ID        uint             `json:"id" gorm:"primary_key"`
	UserID    uint             `json:"user_id"`
	Name      string           `json:"name"`
	WeekStart time.Time        `json:"week_start"`
	IsActive  bool             `json:"is_active" gorm:"default:true"`
	Meals     []MealPlanEntry  `json:"meals"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	User      User             `json:"user"`
}

// CurrentMealPlan represents the user's single active weekly meal plan
type CurrentMealPlan struct {
	ID                uint                  `json:"id" gorm:"primary_key"`
	UserID            uint                  `json:"user_id" gorm:"unique"`
	WeekStart         time.Time             `json:"week_start"`
	Meals             []MealPlanEntry       `json:"meals"`
	ShoppingList      *ShoppingList         `json:"shopping_list,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	User              User                  `json:"user"`
}

type MealPlanEntry struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	MealPlanID uint      `json:"meal_plan_id"`
	MealID     uint      `json:"meal_id"`
	Day        string    `json:"day"` // monday, tuesday, etc.
	MealType   string    `json:"meal_type"` // breakfast, lunch, dinner
	Servings   int       `json:"servings" gorm:"default:1"`
	CreatedAt  time.Time `json:"created_at"`
	Meal       Meal      `json:"meal"`
}

type ShoppingList struct {
	ID               uint                  `json:"id" gorm:"primary_key"`
	UserID           uint                  `json:"user_id"`
	MealPlanID       uint                  `json:"meal_plan_id"`
	Name             string                `json:"name"`
	Items            []ShoppingListItem    `json:"items"`
	IsCompleted      bool                  `json:"is_completed" gorm:"default:false"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	User             User                  `json:"user"`
	MealPlan         MealPlan              `json:"meal_plan"`
}

type ShoppingListItem struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	ShoppingListID uint      `json:"shopping_list_id"`
	IngredientID   uint      `json:"ingredient_id"`
	Quantity       float64   `json:"quantity"`
	Unit           string    `json:"unit"`
	IsPurchased    bool      `json:"is_purchased" gorm:"default:false"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Ingredient     Ingredient `json:"ingredient"`
}

type MealIngredient struct {
	MealID       uint      `json:"meal_id"`
	IngredientID uint      `json:"ingredient_id"`
	Quantity     float64   `json:"quantity"`
	Unit         string    `json:"unit"`
	Ingredient   Ingredient `json:"ingredient"`
}

func (mp *MealPlan) BeforeCreate(scope *gorm.Scope) error {
	return nil
}