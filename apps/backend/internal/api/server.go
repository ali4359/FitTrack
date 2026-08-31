package api

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db        *gorm.DB
	jwtSecret []byte
}

func New(db *gorm.DB) *Server {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-only-insecure-secret-change-me"
	}
	return &Server{db: db, jwtSecret: []byte(secret)}
}

func (s *Server) Router() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:          12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	base := r.Group("/api")
	{
		base.POST("/auth/register", s.handleRegister)
		base.POST("/auth/login", s.handleLogin)

		auth := base.Group("")
		auth.Use(s.authRequired())
		{
			auth.GET("/profile", s.handleGetProfile)
			auth.PATCH("/profile", s.handleUpdateProfile)

			auth.GET("/workouts/history", s.handleWorkoutHistory)
			auth.POST("/workouts/complete", s.handleCompleteWorkout)
			auth.GET("/workouts/:dayId", s.handleGetWorkoutDay)

			auth.GET("/meals/suggest", s.handleSuggestMeals)
			auth.POST("/meals/log", s.handleLogMeal)

			auth.GET("/progress/summary", s.handleProgressSummary)
		}
	}

	return r
}
