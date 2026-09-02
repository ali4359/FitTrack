package api

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ali4359/fittrack/backend/internal/models"
	"github.com/ali4359/fittrack/backend/internal/nutrition"
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
		ExerciseID string `json:"exerciseId"`
		Sets       []struct {
			Reps        int        `json:"reps"`
			WeightKg    float64    `json:"weightKg"`
			CompletedAt *time.Time `json:"completedAt"` // when "set done" was tapped; optional
		} `json:"sets"`
	} `json:"exercises"`
}

func (s *Server) handleCompleteWorkout(c *gin.Context) {
	user, err := s.currentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if user.WeightKg <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "set your body weight in your profile before logging a workout",
			"code":  "weight_required",
		})
		return
	}

	var body completeWorkoutBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logID := uuid.NewString()

	exerciseLogs := []models.WorkoutExerciseLog{}
	burnInputs := []nutrition.BurnExerciseInput{}
	for i, e := range body.Exercises {
		sets := []models.WorkoutSetLog{}
		burnSets := []nutrition.BurnSetInput{}
		for _, st := range e.Sets {
			if st.Reps <= 0 {
				continue
			}
			sets = append(sets, models.WorkoutSetLog{
				ID:          uuid.NewString(),
				SetNumber:   len(sets) + 1,
				Reps:        st.Reps,
				WeightKg:    st.WeightKg,
				CompletedAt: st.CompletedAt,
			})
			burnSets = append(burnSets, nutrition.BurnSetInput{
				Reps:        st.Reps,
				WeightKg:    st.WeightKg,
				CompletedAt: st.CompletedAt,
			})
		}
		if len(sets) == 0 {
			continue
		}

		var ex models.Exercise
		s.db.First(&ex, "id = ?", e.ExerciseID)
		exerciseLogs = append(exerciseLogs, models.WorkoutExerciseLog{
			ID:           uuid.NewString(),
			WorkoutLogID: logID,
			ExerciseID:   e.ExerciseID,
			Position:     i,
			Sets:         sets,
		})
		burnInputs = append(burnInputs, nutrition.BurnExerciseInput{
			ROMMeters:            ex.RangeOfMotionM,
			BodyweightLoadFactor: ex.BodyweightLoadFactor,
			Sets:                 burnSets,
		})
	}

	burn := nutrition.WorkoutBurn(user.WeightKg, body.DurationMinutes, burnInputs)
	for i := range exerciseLogs {
		exerciseLogs[i].CaloriesBurned = burn.Exercises[i].Kcal
		for j := range exerciseLogs[i].Sets {
			exerciseLogs[i].Sets[j].TUTSeconds = burn.Exercises[i].Sets[j].TUTSeconds
			exerciseLogs[i].Sets[j].RestSeconds = burn.Exercises[i].Sets[j].RestSeconds
		}
	}

	logEntry := models.WorkoutLog{
		ID:              logID,
		UserID:          user.ID,
		WorkoutDayID:    body.WorkoutDayID,
		CompletedAt:     time.Now(),
		DurationMinutes: body.DurationMinutes,
		CaloriesBurned:  burn.TotalKcal,
		Exercises:       exerciseLogs,
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
		Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("position asc") }).
		Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("set_number asc") }).
		Where("user_id = ?", currentUserID(c)).
		Order("completed_at desc").
		Limit(30).
		Find(&logs)
	c.JSON(http.StatusOK, gin.H{"results": logs})
}
