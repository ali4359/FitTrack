package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ali4359/fittrack/backend/internal/meals"
	"github.com/ali4359/fittrack/backend/internal/models"
	"github.com/ali4359/fittrack/backend/internal/nutrition"
)

// handleSuggestMeals builds a per-meal macro target from the user's profile,
// today's workout burn, and what they've eaten so far, then asks the meal
// service for dishes that fill the gap. Falls back to the seeded catalogue when
// no LLM backend is configured.
func (s *Server) handleSuggestMeals(c *gin.Context) {
	user, err := s.currentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	mealType := c.DefaultQuery("type", "daily")
	exclude := splitCSV(c.Query("exclude"))

	target, dailyTarget := s.mealTarget(user, mealType)

	diet, allergens := meals.ParseRestrictions(user.Restrictions)
	req := meals.Request{
		Target:    target,
		MealType:  mealType,
		Goal:      user.Goal,
		Diet:      diet,
		Allergens: allergens,
		Region:    user.Region,
		Budget:    user.BudgetTier,
		Exclude:   exclude,
	}

	if s.meals.Enabled() {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
		defer cancel()
		resp, serr := s.meals.Suggest(ctx, req)
		if serr == nil {
			resp.Target = target
			c.JSON(http.StatusOK, withDaily(resp, dailyTarget))
			return
		}
		// fall through to the catalogue on any LLM error
	}

	c.JSON(http.StatusOK, withDaily(s.fallbackSuggest(req), dailyTarget))
}

// mealTarget returns the target for this one meal plus the full-day target.
func (s *Server) mealTarget(user *models.User, mealType string) (nutrition.Macros, nutrition.Macros) {
	prof := nutrition.Profile{
		Sex: user.Sex, Age: user.Age,
		WeightKg: user.WeightKg, HeightCm: user.HeightCm,
		Goal: user.Goal,
	}

	startOfDay := time.Now().Truncate(24 * time.Hour)

	var burn struct{ Total float64 }
	s.db.Model(&models.WorkoutLog{}).
		Select("COALESCE(SUM(calories_burned),0) as total").
		Where("user_id = ? AND completed_at >= ?", user.ID, startOfDay).
		Scan(&burn)

	var eaten nutrition.Macros
	s.db.Model(&models.MealLog{}).
		Select("COALESCE(SUM(calories),0) as calories, COALESCE(SUM(protein_g),0) as protein_g, COALESCE(SUM(carbs_g),0) as carbs_g, COALESCE(SUM(fat_g),0) as fat_g").
		Where("user_id = ? AND eaten_at >= ?", user.ID, startOfDay).
		Scan(&eaten)

	daily := nutrition.DailyTarget(prof, burn.Total)
	remaining := nutrition.Remaining(daily, eaten)
	meal := nutrition.MealTarget(prof, remaining, mealsLeft(mealType), mealType == "post-workout")
	return meal, daily
}

// mealsLeft estimates how many eating occasions remain today, counting the
// current one. A post-workout meal is always "now".
func mealsLeft(mealType string) int {
	switch mealType {
	case "breakfast":
		return 3
	case "lunch":
		return 2
	case "dinner", "post-workout":
		return 1
	}
	switch h := time.Now().Hour(); {
	case h < 11:
		return 3
	case h < 16:
		return 2
	default:
		return 1
	}
}

func withDaily(resp meals.Response, daily nutrition.Macros) gin.H {
	return gin.H{
		"target":      resp.Target,
		"dailyTarget": daily,
		"results":     resp.Results,
		"source":      resp.Source,
		"broadened":   resp.Broadened,
	}
}

// fallbackSuggest ranks the seeded MealEntry catalogue when no LLM is available.
func (s *Server) fallbackSuggest(req meals.Request) meals.Response {
	var all []models.MealEntry
	s.db.Find(&all)

	wantVeg := req.Diet == meals.DietVegetarian || req.Diet == meals.DietVegan
	wantHalal := req.Diet == meals.DietHalal

	type scored struct {
		s     meals.Suggestion
		score float64
	}
	var ranked []scored
	for _, m := range all {
		if wantVeg && !m.Vegetarian {
			continue
		}
		if wantHalal && !m.Halal {
			continue
		}
		sug := meals.Suggestion{
			Name:         m.DishName,
			Portion:      "1 serving",
			Ingredients:  []string{m.DishName},
			Calories:     m.Calories,
			ProteinG:     m.ProteinG,
			CarbsG:       m.CarbsG,
			FatG:         m.FatG,
			IsVegetarian: m.Vegetarian,
			IsHalal:      m.Halal,
			Estimated:    true,
			WhyItFits:    "From the built-in catalogue",
		}
		sc := 0.0
		if req.Target.Calories > 0 {
			sc = abs(m.Calories-req.Target.Calories) / req.Target.Calories
		}
		if strings.EqualFold(m.Region, req.Region) {
			sc -= 0.3
		}
		if strings.EqualFold(m.BudgetTier, req.Budget) {
			sc -= 0.2
		}
		ranked = append(ranked, scored{sug, sc})
	}
	for i := 1; i < len(ranked); i++ {
		for j := i; j > 0 && ranked[j-1].score > ranked[j].score; j-- {
			ranked[j-1], ranked[j] = ranked[j], ranked[j-1]
		}
	}

	out := make([]meals.Suggestion, 0, 3)
	for i, r := range ranked {
		if i == 3 {
			break
		}
		out = append(out, r.s)
	}
	return meals.FallbackResponse(req.Target, out)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

type logMealBody struct {
	// Either reference a catalogue entry...
	MealEntryID string `json:"mealEntryId"`
	// ...or pass an LLM suggestion inline.
	DishName string  `json:"dishName"`
	Calories float64 `json:"calories"`
	ProteinG float64 `json:"proteinG"`
	CarbsG   float64 `json:"carbsG"`
	FatG     float64 `json:"fatG"`

	MealType string  `json:"mealType" binding:"required"`
	Servings float64 `json:"servings"`
}

func (s *Server) handleLogMeal(c *gin.Context) {
	var body logMealBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	servings := body.Servings
	if servings <= 0 {
		servings = 1
	}

	entry := models.MealLog{
		ID:          uuid.NewString(),
		UserID:      currentUserID(c),
		MealEntryID: body.MealEntryID,
		DishName:    body.DishName,
		MealType:    body.MealType,
		Servings:    servings,
		Calories:    body.Calories * servings,
		ProteinG:    body.ProteinG * servings,
		CarbsG:      body.CarbsG * servings,
		FatG:        body.FatG * servings,
		EatenAt:     time.Now(),
	}

	// If a catalogue entry was referenced without inline macros, snapshot them.
	if body.MealEntryID != "" && body.Calories == 0 {
		var m models.MealEntry
		if s.db.First(&m, "id = ?", body.MealEntryID).Error == nil {
			entry.DishName = m.DishName
			entry.Calories = m.Calories * servings
			entry.ProteinG = m.ProteinG * servings
			entry.CarbsG = m.CarbsG * servings
			entry.FatG = m.FatG * servings
		}
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
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
