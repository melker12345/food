package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User handlers
func upsertUser(c *gin.Context) {
	var req UpsertUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Set default goal if not provided
	if req.Goal == "" {
		req.Goal = "maintenance"
	}

	// Validate goal enum
	validGoals := []string{"maintenance", "cutting", "bulking"}
	isValidGoal := false
	for _, goal := range validGoals {
		if req.Goal == goal {
			isValidGoal = true
			break
		}
	}
	if !isValidGoal {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal. Must be one of: maintenance, cutting, bulking"})
		return
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":               req.Name,
			"email":              req.Email,
			"dietaryPreferences": req.DietaryPreferences,
			"healthGoals":        req.HealthGoals,
			"goal":               req.Goal,
			"updatedAt":          now,
		},
		"$setOnInsert": bson.M{
			"createdAt":  now,
			"likedMeals": []primitive.ObjectID{},
		},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.M{"authProviderId": req.AuthProviderID}

	var user User
	err := database.Collection("users").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func getMe(c *gin.Context) {
	authProviderID := c.Query("authProviderId")
	if authProviderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authProviderId is required"})
		return
	}

	var user User
	err := database.Collection("users").FindOne(context.Background(), bson.M{"authProviderId": authProviderID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// Meal handlers
func createMeal(c *gin.Context) {
	var req CreateMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal payload", "details": err.Error()})
		return
	}

	now := time.Now()
	meal := Meal{
		Name:         req.Name,
		ImageURL:     req.ImageURL,
		Ingredients:  req.Ingredients,
		Instructions: req.Instructions,
		Nutrition:    req.Nutrition,
		DietaryTags:  req.DietaryTags,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result, err := database.Collection("meals").InsertOne(context.Background(), meal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meal"})
		return
	}

	meal.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, meal)
}

func getRandomMeal(c *gin.Context) {
	dietary := c.Query("dietary")
	goal := c.Query("goal")

	// Build filter
	filter := bson.M{}

	// Add dietary filter
	if dietary != "" {
		tags := strings.Split(dietary, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		filter["dietaryTags"] = bson.M{"$all": tags}
	}

	// Add goal-based nutrition filter
	switch goal {
	case "cutting":
		// Prefer lower calorie meals (<= 550)
		filter["nutrition.calories"] = bson.M{"$lte": 550}
	case "bulking":
		// Prefer higher protein (>= 25g)
		filter["nutrition.protein"] = bson.M{"$gte": 25}
	}

	// Use aggregation to get random meal with projection
	pipeline := []bson.M{
		{"$match": filter},
		{"$sample": bson.M{"size": 1}},
		{"$project": bson.M{
			"name":         1,
			"imageUrl":     1,
			"instructions": 1,
			"ingredients":  1,
		}},
	}

	cursor, err := database.Collection("meals").Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch random meal"})
		return
	}
	defer cursor.Close(context.Background())

	var meals []bson.M
	if err = cursor.All(context.Background(), &meals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode meal"})
		return
	}

	if len(meals) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No meals found"})
		return
	}

	c.JSON(http.StatusOK, meals[0])
}

func likeMeal(c *gin.Context) {
	var req LikeMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "details": err.Error()})
		return
	}

	mealID, err := primitive.ObjectIDFromHex(req.MealID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal ID"})
		return
	}

	filter := bson.M{"authProviderId": req.AuthProviderID}
	update := bson.M{"$addToSet": bson.M{"likedMeals": mealID}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var user User
	err = database.Collection("users").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like meal"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Plan handlers
func generateWeeklyPlan(c *gin.Context) {
	var req GeneratePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "details": err.Error()})
		return
	}

	// Get user
	var user User
	err := database.Collection("users").FindOne(context.Background(), bson.M{"authProviderId": req.AuthProviderID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	weekStartDate := getMonday(time.Now())

	// Get meals (liked meals or random sample)
	var meals []Meal
	if len(user.LikedMeals) > 0 {
		cursor, err := database.Collection("meals").Find(context.Background(), bson.M{"_id": bson.M{"$in": user.LikedMeals}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get liked meals"})
			return
		}
		defer cursor.Close(context.Background())
		cursor.All(context.Background(), &meals)
	} else {
		// Get random sample of 21 meals
		pipeline := []bson.M{{"$sample": bson.M{"size": 21}}}
		cursor, err := database.Collection("meals").Aggregate(context.Background(), pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get random meals"})
			return
		}
		defer cursor.Close(context.Background())
		cursor.All(context.Background(), &meals)
	}

	if len(meals) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No meals available"})
		return
	}

	// Generate 7 days
	days := make([]DayPlan, 7)
	mealTypes := []string{"breakfast", "lunch", "dinner"}
	mealIndex := 0

	for i := 0; i < 7; i++ {
		date := weekStartDate.AddDate(0, 0, i)
		dayMeals := make([]DailyMealEntry, 0, 3)

		for _, mealType := range mealTypes {
			meal := meals[mealIndex%len(meals)]
			dayMeals = append(dayMeals, DailyMealEntry{
				Meal:     meal.ID,
				MealType: mealType,
			})
			mealIndex++
		}

		days[i] = DayPlan{
			Date:  date,
			Meals: dayMeals,
		}
	}

	now := time.Now()

	// Upsert the plan
	filter := bson.M{"user": user.ID, "weekStartDate": weekStartDate}
	update := bson.M{"$set": bson.M{
		"days":      days,
		"updatedAt": now,
	}, "$setOnInsert": bson.M{
		"user":          user.ID,
		"weekStartDate": weekStartDate,
		"createdAt":     now,
	}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result WeeklyPlan
	err = database.Collection("weeklyplans").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate plan"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func getWeeklyPlan(c *gin.Context) {
	authProviderID := c.Query("authProviderId")
	if authProviderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authProviderId is required"})
		return
	}

	// Get user
	var user User
	err := database.Collection("users").FindOne(context.Background(), bson.M{"authProviderId": authProviderID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	weekStartDate := getMonday(time.Now())

	// Get plan with populated meals
	pipeline := []bson.M{
		{"$match": bson.M{"user": user.ID, "weekStartDate": weekStartDate}},
		{"$unwind": "$days"},
		{"$unwind": "$days.meals"},
		{"$lookup": bson.M{
			"from":         "meals",
			"localField":   "days.meals.meal",
			"foreignField": "_id",
			"as":           "days.meals.meal",
		}},
		{"$unwind": "$days.meals.meal"},
		{"$group": bson.M{
			"_id": "$_id",
			"user": bson.M{"$first": "$user"},
			"weekStartDate": bson.M{"$first": "$weekStartDate"},
			"createdAt": bson.M{"$first": "$createdAt"},
			"updatedAt": bson.M{"$first": "$updatedAt"},
			"days": bson.M{"$push": "$days"},
		}},
	}

	cursor, err := database.Collection("weeklyplans").Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get plan"})
		return
	}
	defer cursor.Close(context.Background())

	var plans []bson.M
	if err = cursor.All(context.Background(), &plans); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode plan"})
		return
	}

	if len(plans) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	c.JSON(http.StatusOK, plans[0])
}

// Shopping handler
func getShoppingList(c *gin.Context) {
	authProviderID := c.Query("authProviderId")
	if authProviderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authProviderId is required"})
		return
	}

	// Get user
	var user User
	err := database.Collection("users").FindOne(context.Background(), bson.M{"authProviderId": authProviderID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	weekStartDate := getMonday(time.Now())

	// Get plan with populated meals
	pipeline := []bson.M{
		{"$match": bson.M{"user": user.ID, "weekStartDate": weekStartDate}},
		{"$unwind": "$days"},
		{"$unwind": "$days.meals"},
		{"$lookup": bson.M{
			"from":         "meals",
			"localField":   "days.meals.meal",
			"foreignField": "_id",
			"as":           "mealData",
		}},
		{"$unwind": "$mealData"},
		{"$unwind": "$mealData.ingredients"},
		{"$group": bson.M{
			"_id": bson.M{"$toLower": "$mealData.ingredients.name"},
			"quantities": bson.M{"$push": "$mealData.ingredients.quantity"},
		}},
		{"$project": bson.M{
			"name":       "$_id",
			"quantities": 1,
			"_id":        0,
		}},
		{"$sort": bson.M{"name": 1}},
	}

	cursor, err := database.Collection("weeklyplans").Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build shopping list"})
		return
	}
	defer cursor.Close(context.Background())

	var items []ShoppingListItem
	if err = cursor.All(context.Background(), &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode shopping list"})
		return
	}

	response := ShoppingListResponse{
		WeekStartDate: weekStartDate,
		Items:         items,
	}

	c.JSON(http.StatusOK, response)
}

// Helper functions
func getMonday(date time.Time) time.Time {
	// Calculate days until Monday (1 = Monday, 0 = Sunday)
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday becomes 7
	}
	daysUntilMonday := 1 - weekday
	monday := date.AddDate(0, 0, daysUntilMonday)
	
	// Set to beginning of day
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
}
