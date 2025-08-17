package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                uint      `json:"id" gorm:"primary_key"`
	Email             string    `json:"email" gorm:"unique;not null"`
	Username          string    `json:"username" gorm:"unique;not null"`
	Password          string    `json:"-" gorm:"not null"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	DietaryRestrictions StringArray `json:"dietary_restrictions" gorm:"type:text[]"`
	PreferredMealTypes  StringArray `json:"preferred_meal_types" gorm:"type:text[]"`
	Allergies           StringArray `json:"allergies" gorm:"type:text[]"`
	CalorieGoal         int      `json:"calorie_goal"`
	IsActive            bool     `json:"is_active" gorm:"default:true"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at" sql:"index"`
}

type UserPreferences struct {
	UserID              uint        `json:"user_id"`
	DietaryRestrictions StringArray `json:"dietary_restrictions"`
	PreferredMealTypes  StringArray `json:"preferred_meal_types"`
	Allergies           StringArray `json:"allergies"`
	CalorieGoal         int         `json:"calorie_goal"`
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("CreatedAt", time.Now())
}
