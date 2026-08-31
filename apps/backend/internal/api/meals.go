package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ali4359/iron-and-spice/backend/internal/models"
)

// handleSuggestMeals reads goal/region/budget/restrictions from the saved profile
// (the client does not pass them) and ranks the seeded meal catalogue.
func (s *Server) handleSuggestMeals(c *gin.Context) {
	user, err := s.currentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	mealType := c.DefaultQuery("type", "daily")

	var all []models.MealEntry
	s.db.Find(&all)

	restrictions := splitCSV(user.Restrictions)
	wantHalal := contains(restrictions, "halal")
	wantVeg := contains(restrictions, "vegetarian")

	type scored struct {
		meal  models.MealEntry
		score int
	}
	var ranked []scored
	for _, m := range all {
		if wantHalal && !m.Halal {
			continue
		}
		if wantVeg && !m.Vegetarian {
			continue
		}

		score := 0
		if strings.EqualFold(m.Region, user.Region) {
			score += 3
		}
		if strings.EqualFold(m.BudgetTier, user.BudgetTier) {
			score += 2
		} else if budgetRank(m.BudgetTier) <= budgetRank(user.BudgetTier) {
			score += 1 // cheaper than the user's ceiling is still fine
		}
		tags := splitCSV(m.GoalTags)
		if contains(tags, user.Goal) {
			score += 2
		}
		if mealType != "daily" && contains(tags, mealType) {
			score += 2
		}
		ranked = append(ranked, scored{m, score})
	}

	// simple insertion sort by score desc — small catalogue
	for i := 1; i < len(ranked); i++ {
		for j := i; j > 0 && ranked[j-1].score < ranked[j].score; j-- {
			ranked[j-1], ranked[j] = ranked[j], ranked[j-1]
		}
	}

	results := make([]models.MealEntry, 0, len(ranked))
	for _, r := range ranked {
		results = append(results, r.meal)
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

type logMealBody struct {
	MealEntryID string `json:"mealEntryId" binding:"required"`
	MealType    string `json:"mealType" binding:"required"`
}

func (s *Server) handleLogMeal(c *gin.Context) {
	var body logMealBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entry := models.MealLog{
		ID:          uuid.NewString(),
		UserID:      currentUserID(c),
		MealEntryID: body.MealEntryID,
		MealType:    body.MealType,
		EatenAt:     time.Now(),
	}
	if err := s.db.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not log meal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logged": true})
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(strings.ToLower(p)); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func contains(xs []string, want string) bool {
	for _, x := range xs {
		if x == strings.ToLower(want) {
			return true
		}
	}
	return false
}

func budgetRank(tier string) int {
	switch strings.ToLower(tier) {
	case "low":
		return 0
	case "mid":
		return 1
	case "high":
		return 2
	default:
		return 1
	}
}
