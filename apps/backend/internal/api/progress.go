package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ali4359/fittrack/backend/internal/models"
)

func (s *Server) handleProgressSummary(c *gin.Context) {
	userID := currentUserID(c)
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)

	var logs []models.WorkoutLog
	s.db.Where("user_id = ? AND completed_at >= ?", userID, monthStart).Find(&logs)

	var totalCals float64
	for _, l := range logs {
		totalCals += l.CaloriesBurned
	}
	avg := 0.0
	if len(logs) > 0 {
		avg = totalCals / float64(len(logs))
	}

	c.JSON(http.StatusOK, gin.H{
		"workoutsThisMonth": len(logs),
		"avgCaloriesBurned": avg,
	})
}
