package models

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Meal struct {
	ID               uint           `json:"id" gorm:"primary_key"`
	Name             string         `json:"name" gorm:"not null"`
	Description      string         `json:"description"`
	ImageURL         string         `json:"image_url"`
	PrepTime         int            `json:"prep_time"` // in minutes
	CookTime         int            `json:"cook_time"` // in minutes
	Servings         int            `json:"servings" gorm:"default:4"`
	Difficulty       string         `json:"difficulty"` // easy, medium, hard
	Cuisine          string         `json:"cuisine"`
	MealType         string         `json:"meal_type"` // breakfast, lunch, dinner, snack
	Instructions     string `json:"instructions" gorm:"type:text"`
	Ingredients      []Ingredient   `json:"ingredients" gorm:"many2many:meal_ingredients;"`
	NutritionInfo    NutritionInfo  `json:"nutrition_info" gorm:"embedded"`
	DietaryTags      StringArray `json:"dietary_tags" gorm:"type:text[]"`
	Allergens        StringArray `json:"allergens" gorm:"type:text[]"`
	LikesCount       int            `json:"likes_count" gorm:"default:0"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        *time.Time     `json:"deleted_at" sql:"index"`
}

type Ingredient struct {
	ID          uint    `json:"id" gorm:"primary_key"`
	Name        string  `json:"name" gorm:"unique;not null"`
	Category    string  `json:"category"` // protein, vegetable, grain, etc.
	Unit        string  `json:"unit"`     // cup, tbsp, piece, etc.
	CaloriesPer100g float64 `json:"calories_per_100g"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MealIngredient struct {
	MealID       uint        `json:"meal_id"`
	IngredientID uint        `json:"ingredient_id"`
	Quantity     float64     `json:"quantity"`
	Unit         string      `json:"unit"`
	Ingredient   Ingredient  `json:"ingredient"`
}

type NutritionInfo struct {
	Calories      float64 `json:"calories"`
	Protein       float64 `json:"protein"`      // grams
	Carbohydrates float64 `json:"carbohydrates"` // grams
	Fat           float64 `json:"fat"`          // grams
	Fiber         float64 `json:"fiber"`        // grams
	Sugar         float64 `json:"sugar"`        // grams
	Sodium        float64 `json:"sodium"`       // milligrams
}

type UserMealInteraction struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UserID    uint      `json:"user_id"`
	MealID    uint      `json:"meal_id"`
	Liked     bool      `json:"liked"`
	Disliked  bool      `json:"disliked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user"`
	Meal      Meal      `json:"meal"`
}

type MealReview struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UserID    uint      `json:"user_id"`
	MealID    uint      `json:"meal_id"`
	Rating    int       `json:"rating"` // 1-5 stars
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user"`
	Meal      Meal      `json:"meal"`
}

// StringArray is a custom type for PostgreSQL string arrays
type StringArray []string

// Value implements the driver.Valuer interface for database storage
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return "{}", nil
	}
	return pq.Array(sa).Value()
}

// Scan implements the sql.Scanner interface for database retrieval
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = StringArray{}
		return nil
	}

	switch v := value.(type) {
	case string:
		// Handle JSON string format for backward compatibility
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			var arr []string
			if err := json.Unmarshal([]byte(v), &arr); err == nil {
				*sa = StringArray(arr)
				return nil
			}
		}
		// Handle PostgreSQL array format
		return pq.Array(sa).Scan(value)
	case []byte:
		// Handle JSON byte format
		if len(v) > 0 && v[0] == '[' {
			var arr []string
			if err := json.Unmarshal(v, &arr); err == nil {
				*sa = StringArray(arr)
				return nil
			}
		}
		return pq.Array(sa).Scan(value)
	default:
		return pq.Array(sa).Scan(value)
	}
}

// MarshalJSON implements json.Marshaler
func (sa StringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(sa))
}

// UnmarshalJSON implements json.Unmarshaler
func (sa *StringArray) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*sa = StringArray(arr)
	return nil
}

func (m *Meal) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("CreatedAt", time.Now())
}
