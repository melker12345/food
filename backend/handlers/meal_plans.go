package handlers

import (
	"net/http"
	"time"

	"food-app/database"
	"food-app/models"

	"github.com/gin-gonic/gin"
	"math/rand"
)

type CreateMealPlanRequest struct {
	Name      string             `json:"name" binding:"required"`
	WeekStart string             `json:"week_start" binding:"required"`
	Meals     []MealPlanEntryReq `json:"meals"`
}

type MealPlanEntryReq struct {
	MealID   uint   `json:"meal_id" binding:"required"`
	Day      string `json:"day" binding:"required"`
	MealType string `json:"meal_type" binding:"required"`
	Servings int    `json:"servings"`
}

type AutoGenerateMealPlanRequest struct {
	Name      string `json:"name" binding:"required"`
	WeekStart string `json:"week_start" binding:"required"`
}

func CreateMealPlan(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateMealPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse week start date
	weekStart, err := time.Parse("2006-01-02", req.WeekStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Create meal plan
	mealPlan := models.MealPlan{
		UserID:    userID,
		Name:      req.Name,
		WeekStart: weekStart,
		IsActive:  true,
	}

	if err := database.DB.Create(&mealPlan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal plan"})
		return
	}

	// Add meal plan entries
	for _, mealReq := range req.Meals {
		servings := mealReq.Servings
		if servings == 0 {
			servings = 1
		}

		entry := models.MealPlanEntry{
			MealPlanID: mealPlan.ID,
			MealID:     mealReq.MealID,
			Day:        mealReq.Day,
			MealType:   mealReq.MealType,
			Servings:   servings,
		}

		database.DB.Create(&entry)
	}

	// Load the complete meal plan with relationships
	database.DB.Preload("Meals").Preload("Meals.Meal").Preload("Meals.Meal.Ingredients").First(&mealPlan, mealPlan.ID)

	c.JSON(http.StatusCreated, mealPlan)
}

func AutoGenerateMealPlan(c *gin.Context) {
	userID := c.GetUint("userID")

	var req AutoGenerateMealPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse week start date
	weekStart, err := time.Parse("2006-01-02", req.WeekStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Get user's liked meals
	var interactions []models.UserMealInteraction
	if err := database.DB.Preload("Meal").Preload("Meal.Ingredients").
		Where("user_id = ? AND liked = true", userID).Find(&interactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch liked meals"})
		return
	}

	if len(interactions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No liked meals found. Please like some meals first."})
		return
	}

	// Separate meals by type
	mealsByType := map[string][]models.Meal{
		"breakfast": {},
		"lunch":     {},
		"dinner":    {},
	}

	for _, interaction := range interactions {
		meal := interaction.Meal
		if meal.MealType != "" {
			mealsByType[meal.MealType] = append(mealsByType[meal.MealType], meal)
		}
	}

	// If no specific meal types, use all meals for all types
	allMeals := []models.Meal{}
	for _, interaction := range interactions {
		allMeals = append(allMeals, interaction.Meal)
	}

	// Create meal plan
	mealPlan := models.MealPlan{
		UserID:    userID,
		Name:      req.Name,
		WeekStart: weekStart,
		IsActive:  true,
	}

	if err := database.DB.Create(&mealPlan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal plan"})
		return
	}

	// Generate meal entries for the week
	days := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	mealTypes := []string{"breakfast", "lunch", "dinner"}

	for _, day := range days {
		for _, mealType := range mealTypes {
			// Choose meals for this meal type
			availableMeals := mealsByType[mealType]
			if len(availableMeals) == 0 {
				// Fallback to all meals if no specific type available
				availableMeals = allMeals
			}

			if len(availableMeals) > 0 {
				// Randomly select a meal
				selectedMeal := availableMeals[rand.Intn(len(availableMeals))]

				entry := models.MealPlanEntry{
					MealPlanID: mealPlan.ID,
					MealID:     selectedMeal.ID,
					Day:        day,
					MealType:   mealType,
					Servings:   1,
				}

				database.DB.Create(&entry)
			}
		}
	}

	// Load the complete meal plan with relationships
	database.DB.Preload("Meals").Preload("Meals.Meal").Preload("Meals.Meal.Ingredients").First(&mealPlan, mealPlan.ID)

	c.JSON(http.StatusCreated, mealPlan)
}

func GetMealPlans(c *gin.Context) {
	userID := c.GetUint("userID")

	var mealPlans []models.MealPlan
	if err := database.DB.Preload("Meals").Preload("Meals.Meal").
		Where("user_id = ?", userID).Order("created_at DESC").Find(&mealPlans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meal plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meal_plans": mealPlans})
}

func GetMealPlan(c *gin.Context) {
	userID := c.GetUint("userID")
	planID := c.Param("id")

	var mealPlan models.MealPlan
	if database.DB.Preload("Meals").Preload("Meals.Meal").Preload("Meals.Meal.Ingredients").
		Where("id = ? AND user_id = ?", planID, userID).First(&mealPlan).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
		return
	}

	c.JSON(http.StatusOK, mealPlan)
}

func UpdateMealPlan(c *gin.Context) {
	userID := c.GetUint("userID")
	planID := c.Param("id")

	var mealPlan models.MealPlan
	if database.DB.Where("id = ? AND user_id = ?", planID, userID).First(&mealPlan).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
		return
	}

	var req CreateMealPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update meal plan
	if req.Name != "" {
		mealPlan.Name = req.Name
	}

	if req.WeekStart != "" {
		if weekStart, err := time.Parse("2006-01-02", req.WeekStart); err == nil {
			mealPlan.WeekStart = weekStart
		}
	}

	database.DB.Save(&mealPlan)

	// Update meal entries if provided
	if len(req.Meals) > 0 {
		// Delete existing entries
		database.DB.Where("meal_plan_id = ?", mealPlan.ID).Delete(&models.MealPlanEntry{})

		// Add new entries
		for _, mealReq := range req.Meals {
			servings := mealReq.Servings
			if servings == 0 {
				servings = 1
			}

			entry := models.MealPlanEntry{
				MealPlanID: mealPlan.ID,
				MealID:     mealReq.MealID,
				Day:        mealReq.Day,
				MealType:   mealReq.MealType,
				Servings:   servings,
			}

			database.DB.Create(&entry)
		}
	}

	// Load updated meal plan
	database.DB.Preload("Meals").Preload("Meals.Meal").Preload("Meals.Meal.Ingredients").First(&mealPlan, mealPlan.ID)

	c.JSON(http.StatusOK, mealPlan)
}

func DeleteMealPlan(c *gin.Context) {
	userID := c.GetUint("userID")
	planID := c.Param("id")

	var mealPlan models.MealPlan
	if database.DB.Where("id = ? AND user_id = ?", planID, userID).First(&mealPlan).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
		return
	}

	// Delete meal plan entries first
	database.DB.Where("meal_plan_id = ?", mealPlan.ID).Delete(&models.MealPlanEntry{})

	// Delete meal plan
	database.DB.Delete(&mealPlan)

	c.JSON(http.StatusOK, gin.H{"message": "Meal plan deleted successfully"})
}

func GenerateShoppingList(c *gin.Context) {
	userID := c.GetUint("userID")
	planID := c.Param("id")

	var mealPlan models.MealPlan
	if database.DB.Preload("Meals").Preload("Meals.Meal").Preload("Meals.Meal.Ingredients").
		Where("id = ? AND user_id = ?", planID, userID).First(&mealPlan).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal plan not found"})
		return
	}

	// Create shopping list
	shoppingList := models.ShoppingList{
		UserID:     userID,
		MealPlanID: mealPlan.ID,
		Name:       "Shopping List for " + mealPlan.Name,
	}

	if err := database.DB.Create(&shoppingList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shopping list"})
		return
	}

	// Aggregate ingredients from all meals
	ingredientQuantities := make(map[uint]map[string]float64) // ingredient_id -> unit -> total_quantity

	for _, entry := range mealPlan.Meals {
		// Get meal ingredients with quantities (this would need to be implemented in meal seeding)
		var mealIngredients []models.MealIngredient
		database.DB.Preload("Ingredient").Where("meal_id = ?", entry.MealID).Find(&mealIngredients)

		for _, mealIngredient := range mealIngredients {
			ingredientID := mealIngredient.IngredientID
			unit := mealIngredient.Unit
			quantity := mealIngredient.Quantity * float64(entry.Servings)

			if ingredientQuantities[ingredientID] == nil {
				ingredientQuantities[ingredientID] = make(map[string]float64)
			}
			ingredientQuantities[ingredientID][unit] += quantity
		}
	}

	// Create shopping list items
	for ingredientID, units := range ingredientQuantities {
		for unit, quantity := range units {
			item := models.ShoppingListItem{
				ShoppingListID: shoppingList.ID,
				IngredientID:   ingredientID,
				Quantity:       quantity,
				Unit:           unit,
			}
			database.DB.Create(&item)
		}
	}

	// Load complete shopping list
	database.DB.Preload("Items").Preload("Items.Ingredient").First(&shoppingList, shoppingList.ID)

	c.JSON(http.StatusCreated, shoppingList)
}

func GetShoppingLists(c *gin.Context) {
	userID := c.GetUint("userID")

	var shoppingLists []models.ShoppingList
	if err := database.DB.Preload("Items").Preload("Items.Ingredient").
		Where("user_id = ?", userID).Order("created_at DESC").Find(&shoppingLists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping lists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shopping_lists": shoppingLists})
}

func UpdateShoppingListItem(c *gin.Context) {
	userID := c.GetUint("userID")
	itemID := c.Param("item_id")

	var item models.ShoppingListItem
	if database.DB.Joins("JOIN shopping_lists ON shopping_list_items.shopping_list_id = shopping_lists.id").
		Where("shopping_list_items.id = ? AND shopping_lists.user_id = ?", itemID, userID).
		First(&item).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shopping list item not found"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&item).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, item)
}
