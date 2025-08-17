package handlers

import (
	"net/http"
	"time"
	"math/rand"

	"food-app/database"
	"food-app/models"

	"github.com/gin-gonic/gin"
)

// GetCurrentMealPlan gets the user's single active meal plan
func GetCurrentMealPlan(c *gin.Context) {
	userID := c.GetUint("userID")

	var mealPlan models.CurrentMealPlan
	if database.DB.Preload("Meals").Preload("Meals.Meal").Preload("Meals.Meal.Ingredients").
		Preload("ShoppingList").Preload("ShoppingList.Items").Preload("ShoppingList.Items.Ingredient").
		Where("user_id = ?", userID).First(&mealPlan).RecordNotFound() {
		
		// If no current meal plan exists, create one for this week
		weekStart := getCurrentWeekStart()
		mealPlan = models.CurrentMealPlan{
			UserID:    userID,
			WeekStart: weekStart,
		}
		
		if err := database.DB.Create(&mealPlan).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal plan"})
			return
		}
	}

	c.JSON(http.StatusOK, mealPlan)
}

// PopulateFromLikedMeals auto-populates the current meal plan with liked meals
func PopulateFromLikedMeals(c *gin.Context) {
	userID := c.GetUint("userID")

	// Get or create current meal plan
	var mealPlan models.CurrentMealPlan
	if database.DB.Where("user_id = ?", userID).First(&mealPlan).RecordNotFound() {
		weekStart := getCurrentWeekStart()
		mealPlan = models.CurrentMealPlan{
			UserID:    userID,
			WeekStart: weekStart,
		}
		
		if err := database.DB.Create(&mealPlan).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal plan"})
			return
		}
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

	// Clear existing meals
	database.DB.Where("meal_plan_id = ?", mealPlan.ID).Delete(&models.MealPlanEntry{})

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

	// Update shopping list automatically
	updateShoppingListForCurrentPlan(userID, mealPlan.ID)

	// Load updated meal plan
	database.DB.Preload("Meals").Preload("Meals.Meal").Preload("Meals.Meal.Ingredients").
		Preload("ShoppingList").Preload("ShoppingList.Items").Preload("ShoppingList.Items.Ingredient").
		First(&mealPlan, mealPlan.ID)

	c.JSON(http.StatusOK, mealPlan)
}

// UpdateMealInPlan updates a specific meal in the current plan
func UpdateMealInPlan(c *gin.Context) {
	userID := c.GetUint("userID")
	
	type UpdateMealRequest struct {
		Day      string `json:"day" binding:"required"`
		MealType string `json:"meal_type" binding:"required"`
		MealID   *uint  `json:"meal_id"` // nil to remove meal
		Servings int    `json:"servings"`
	}

	var req UpdateMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current meal plan
	var mealPlan models.CurrentMealPlan
	if database.DB.Where("user_id = ?", userID).First(&mealPlan).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active meal plan found"})
		return
	}

	// Remove existing meal for this day/type
	database.DB.Where("meal_plan_id = ? AND day = ? AND meal_type = ?", 
		mealPlan.ID, req.Day, req.MealType).Delete(&models.MealPlanEntry{})

	// Add new meal if provided
	if req.MealID != nil {
		servings := req.Servings
		if servings == 0 {
			servings = 1
		}

		entry := models.MealPlanEntry{
			MealPlanID: mealPlan.ID,
			MealID:     *req.MealID,
			Day:        req.Day,
			MealType:   req.MealType,
			Servings:   servings,
		}

		if err := database.DB.Create(&entry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add meal to plan"})
			return
		}
	}

	// Update shopping list automatically
	updateShoppingListForCurrentPlan(userID, mealPlan.ID)

	// Return updated meal plan
	database.DB.Preload("Meals").Preload("Meals.Meal").Preload("Meals.Meal.Ingredients").
		Preload("ShoppingList").Preload("ShoppingList.Items").Preload("ShoppingList.Items.Ingredient").
		First(&mealPlan, mealPlan.ID)

	c.JSON(http.StatusOK, mealPlan)
}

// ToggleShoppingItem toggles the purchased status of a shopping list item
func ToggleShoppingItem(c *gin.Context) {
	userID := c.GetUint("userID")
	itemID := c.Param("item_id")

	type ToggleRequest struct {
		IsPurchased bool   `json:"is_purchased"`
		Notes       string `json:"notes"`
	}

	var req ToggleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the item belongs to the user's shopping list
	var item models.ShoppingListItem
	if database.DB.Joins("JOIN shopping_lists ON shopping_list_items.shopping_list_id = shopping_lists.id").
		Where("shopping_list_items.id = ? AND shopping_lists.user_id = ?", itemID, userID).
		First(&item).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shopping list item not found"})
		return
	}

	// Update the item
	item.IsPurchased = req.IsPurchased
	item.Notes = req.Notes
	
	if err := database.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shopping list item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Helper functions
func getCurrentWeekStart() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	daysFromMonday := weekday - 1
	monday := now.AddDate(0, 0, -daysFromMonday)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
}

func updateShoppingListForCurrentPlan(userID, mealPlanID uint) {
	// Delete existing shopping list for this meal plan
	var existingList models.ShoppingList
	if !database.DB.Where("meal_plan_id = ?", mealPlanID).First(&existingList).RecordNotFound() {
		database.DB.Where("shopping_list_id = ?", existingList.ID).Delete(&models.ShoppingListItem{})
		database.DB.Delete(&existingList)
	}

	// Get meal plan with all meals
	var mealPlan models.CurrentMealPlan
	if database.DB.Preload("Meals").Preload("Meals.Meal").Where("id = ?", mealPlanID).First(&mealPlan).RecordNotFound() {
		return
	}

	// Create new shopping list
	shoppingList := models.ShoppingList{
		UserID:     userID,
		MealPlanID: mealPlanID,
		Name:       "Week of " + mealPlan.WeekStart.Format("Jan 2, 2006"),
	}

	if err := database.DB.Create(&shoppingList).Error; err != nil {
		return
	}

	// Aggregate ingredients from all meals
	ingredientQuantities := make(map[uint]map[string]float64) // ingredient_id -> unit -> total_quantity

	for _, entry := range mealPlan.Meals {
		// Get meal ingredients with quantities
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
				IsPurchased:    false,
			}

			database.DB.Create(&item)
		}
	}
}
