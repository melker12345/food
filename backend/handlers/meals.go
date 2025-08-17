package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"food-app/database"
	"food-app/models"

	"github.com/gin-gonic/gin"
)

func GetMeals(c *gin.Context) {
	var meals []models.Meal
	query := database.DB.Preload("Ingredients")

	// Apply filters
	if cuisine := c.Query("cuisine"); cuisine != "" {
		query = query.Where("cuisine = ?", cuisine)
	}

	if mealType := c.Query("meal_type"); mealType != "" {
		query = query.Where("meal_type = ?", mealType)
	}

	if difficulty := c.Query("difficulty"); difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	if maxPrepTime := c.Query("max_prep_time"); maxPrepTime != "" {
		if time, err := strconv.Atoi(maxPrepTime); err == nil {
			query = query.Where("prep_time <= ?", time)
		}
	}

	if dietaryTags := c.Query("dietary_tags"); dietaryTags != "" {
		tags := strings.Split(dietaryTags, ",")
		for _, tag := range tags {
			query = query.Where("? = ANY(dietary_tags)", strings.TrimSpace(tag))
		}
	}

	// Exclude allergens
	if allergens := c.Query("exclude_allergens"); allergens != "" {
		excludeList := strings.Split(allergens, ",")
		for _, allergen := range excludeList {
			query = query.Where("NOT (? = ANY(allergens))", strings.TrimSpace(allergen))
		}
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&meals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"meals": meals,
		"page":  page,
		"limit": limit,
	})
}

func GetMeal(c *gin.Context) {
	id := c.Param("id")
	
	var meal models.Meal
	if database.DB.Preload("Ingredients").First(&meal, id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal not found"})
		return
	}

	c.JSON(http.StatusOK, meal)
}

func GetPersonalizedMeals(c *gin.Context) {
	userID := c.GetUint("userID")
	
	var user models.User
	if database.DB.First(&user, userID).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var meals []models.Meal
	query := database.DB.Preload("Ingredients")

	// Filter by user's dietary restrictions
	for _, tag := range user.DietaryRestrictions {
		query = query.Where("? = ANY(dietary_tags)", tag)
	}

	// Exclude user's allergens
	for _, allergen := range user.Allergies {
		query = query.Where("NOT (? = ANY(allergens))", allergen)
	}

	// Filter by preferred meal types
	if len(user.PreferredMealTypes) > 0 {
		query = query.Where("meal_type IN (?)", user.PreferredMealTypes)
	}

	// Exclude meals the user has disliked
	var dislikedMealIDs []uint
	database.DB.Model(&models.UserMealInteraction{}).
		Where("user_id = ? AND disliked = true", userID).
		Pluck("meal_id", &dislikedMealIDs)
	
	if len(dislikedMealIDs) > 0 {
		query = query.Where("id NOT IN (?)", dislikedMealIDs)
	}

	// Order by likes count and randomize a bit
	query = query.Order("likes_count DESC, RANDOM()")

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&meals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch personalized meals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"meals": meals,
		"page":  page,
		"limit": limit,
	})
}

func LikeMeal(c *gin.Context) {
	userID := c.GetUint("userID")
	mealID := c.Param("id")

	// Check if interaction already exists
	var interaction models.UserMealInteraction
	if database.DB.Where("user_id = ? AND meal_id = ?", userID, mealID).First(&interaction).RecordNotFound() {
		// Create new interaction
		interaction = models.UserMealInteraction{
			UserID: userID,
			MealID: parseUint(mealID),
			Liked:  true,
		}
		database.DB.Create(&interaction)
	} else {
		// Update existing interaction
		interaction.Liked = true
		interaction.Disliked = false
		database.DB.Save(&interaction)
	}

	// Update meal likes count
	database.DB.Model(&models.Meal{}).Where("id = ?", mealID).
		UpdateColumn("likes_count", database.DB.Model(&models.UserMealInteraction{}).
			Where("meal_id = ? AND liked = true", mealID).Select("count(*)"))

	c.JSON(http.StatusOK, gin.H{"message": "Meal liked successfully"})
}

func DislikeMeal(c *gin.Context) {
	userID := c.GetUint("userID")
	mealID := c.Param("id")

	// Check if interaction already exists
	var interaction models.UserMealInteraction
	if database.DB.Where("user_id = ? AND meal_id = ?", userID, mealID).First(&interaction).RecordNotFound() {
		// Create new interaction
		interaction = models.UserMealInteraction{
			UserID:   userID,
			MealID:   parseUint(mealID),
			Disliked: true,
		}
		database.DB.Create(&interaction)
	} else {
		// Update existing interaction
		interaction.Liked = false
		interaction.Disliked = true
		database.DB.Save(&interaction)
	}

	// Update meal likes count
	database.DB.Model(&models.Meal{}).Where("id = ?", mealID).
		UpdateColumn("likes_count", database.DB.Model(&models.UserMealInteraction{}).
			Where("meal_id = ? AND liked = true", mealID).Select("count(*)"))

	c.JSON(http.StatusOK, gin.H{"message": "Meal disliked successfully"})
}

func GetLikedMeals(c *gin.Context) {
	userID := c.GetUint("userID")

	var interactions []models.UserMealInteraction
	if err := database.DB.Preload("Meal").Preload("Meal.Ingredients").
		Where("user_id = ? AND liked = true", userID).Find(&interactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch liked meals"})
		return
	}

	var meals []models.Meal
	for _, interaction := range interactions {
		meals = append(meals, interaction.Meal)
	}

	c.JSON(http.StatusOK, gin.H{"meals": meals})
}

func GetTrendingMeals(c *gin.Context) {
	var meals []models.Meal
	query := database.DB.Preload("Ingredients").Order("likes_count DESC")

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&meals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trending meals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"meals": meals,
		"page":  page,
		"limit": limit,
	})
}

func AddMealReview(c *gin.Context) {
	userID := c.GetUint("userID")
	mealID := c.Param("id")

	var review models.MealReview
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review.UserID = userID
	review.MealID = parseUint(mealID)

	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func GetMealReviews(c *gin.Context) {
	mealID := c.Param("id")

	var reviews []models.MealReview
	if err := database.DB.Preload("User").Where("meal_id = ?", mealID).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func parseUint(s string) uint {
	if val, err := strconv.ParseUint(s, 10, 32); err == nil {
		return uint(val)
	}
	return 0
}
