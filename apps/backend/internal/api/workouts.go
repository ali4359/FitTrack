package api

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ali4359/fittrack/backend/internal/models"
)

func (s *Server) loadWorkoutDay(id string) (*models.WorkoutDay, error) {
	var day models.WorkoutDay
	err := s.db.
		Preload("Exercises.Exercise").
		First(&day, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	sort.Slice(day.Exercises, func(i, j int) bool {
		return day.Exercises[i].Position < day.Exercises[j].Position
	})
	return &day, nil
}

func (s *Server) handleGetWorkoutDay(c *gin.Context) {
	day, err := s.loadWorkoutDay(c.Param("dayId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workout day not found"})
		return
	}
	c.JSON(http.StatusOK, day)
}

type completeWorkoutBody struct {
	WorkoutDayID    string `json:"workoutDayId" binding:"required"`
	DurationMinutes int    `json:"durationMinutes" binding:"required"`
	Exercises       []struct {
		ExerciseID string  `json:"exerciseId"`
		SetsDone   int     `json:"setsDone"`
		RepsDone   int     `json:"repsDone"`
		WeightKg   float64 `json:"weightKg"`
	} `json:"exercises"`
}

// estimateCalories: MET formula — kcal = MET * 3.5 * kg / 200 * minutes,
// averaged across the exercises actually performed.
func estimateCalories(weightKg float64, minutes int, mets []float64) float64 {
	if len(mets) == 0 || weightKg <= 0 {
		return 0
	}
	var sum float64
	for _, m := range mets {
		sum += m
	}
	avgMet := sum / float64(len(mets))
	return avgMet * 3.5 * weightKg / 200 * float64(minutes)
}

func (s *Server) handleCompleteWorkout(c *gin.Context) {
	user, err := s.currentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	var body completeWorkoutBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mets := []float64{}
	for _, e := range body.Exercises {
		if e.SetsDone == 0 {
			continue
		}
		var ex models.Exercise
		if s.db.First(&ex, "id = ?", e.ExerciseID).Error == nil {
			mets = append(mets, ex.MetValue)
		}
	}

	weight := user.WeightKg
	if weight <= 0 {
		weight = 75
	}
	calories := estimateCalories(weight, body.DurationMinutes, mets)

	logEntry := models.WorkoutLog{
		ID:              uuid.NewString(),
		UserID:          user.ID,
		WorkoutDayID:    body.WorkoutDayID,
		CompletedAt:     time.Now(),
		DurationMinutes: body.DurationMinutes,
		CaloriesBurned:  calories,
	}
	if err := s.db.Create(&logEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save workout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workoutLog": logEntry,
		"nextStep": gin.H{
			"kind":     "meal-suggestion",
			"mealType": "post-workout",
			"message":  "Nice work. Here's a post-workout meal picked for your goal and budget.",
		},
	})
}

func (s *Server) handleWorkoutHistory(c *gin.Context) {
	var logs []models.WorkoutLog
	s.db.
		Where("user_id = ?", currentUserID(c)).
		Order("completed_at desc").
		Limit(30).
		Find(&logs)
	c.JSON(http.StatusOK, gin.H{"results": logs})
}
