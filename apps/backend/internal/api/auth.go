package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/ali4359/fittrack/backend/internal/models"
)

const userCtxKey = "userID"

func (s *Server) issueToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

// authRequired validates the Bearer token and stashes the user id on the context.
func (s *Server) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		raw, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token, err := jwt.ParseWithClaims(raw, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			return s.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims := token.Claims.(*jwt.RegisteredClaims)
		c.Set(userCtxKey, claims.Subject)
		c.Next()
	}
}

func currentUserID(c *gin.Context) string { return c.GetString(userCtxKey) }

func (s *Server) currentUser(c *gin.Context) (*models.User, error) {
	var u models.User
	if err := s.db.First(&u, "id = ?", currentUserID(c)).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

type registerBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

func (s *Server) handleRegister(c *gin.Context) {
	var body registerBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing int64
	s.db.Model(&models.User{}).Where("email = ?", strings.ToLower(body.Email)).Count(&existing)
	if existing > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	user := models.User{
		ID:           uuid.NewString(),
		Email:        strings.ToLower(body.Email),
		Name:         body.Name,
		PasswordHash: string(hash),
		Goal:         "maintain",
		BudgetTier:   "mid",
	}
	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	token, _ := s.issueToken(user.ID)
	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

type loginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) handleLogin(c *gin.Context) {
	var body loginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := s.db.First(&user, "email = ?", strings.ToLower(body.Email)).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, _ := s.issueToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

type updateProfileBody struct {
	Name         *string  `json:"name"`
	Goal         *string  `json:"goal"`
	Sex          *string  `json:"sex"`
	Age          *int     `json:"age"`
	WeightKg     *float64 `json:"weightKg"`
	HeightCm     *float64 `json:"heightCm"`
	Region       *string  `json:"region"`
	BudgetTier   *string  `json:"budgetTier"`
	Restrictions *string  `json:"restrictions"`
}

func (s *Server) handleGetProfile(c *gin.Context) {
	user, err := s.currentUser(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (s *Server) handleUpdateProfile(c *gin.Context) {
	user, err := s.currentUser(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var body updateProfileBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]any{}
	if body.Name != nil {
		updates["name"] = *body.Name
	}
	if body.Goal != nil {
		updates["goal"] = *body.Goal
	}
	if body.Sex != nil {
		updates["sex"] = *body.Sex
	}
	if body.Age != nil {
		updates["age"] = *body.Age
	}
	if body.WeightKg != nil {
		updates["weight_kg"] = *body.WeightKg
	}
	if body.HeightCm != nil {
		updates["height_cm"] = *body.HeightCm
	}
	if body.Region != nil {
		updates["region"] = *body.Region
	}
	if body.BudgetTier != nil {
		updates["budget_tier"] = *body.BudgetTier
	}
	if body.Restrictions != nil {
		updates["restrictions"] = *body.Restrictions
	}

	if len(updates) > 0 {
		if err := s.db.Model(user).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update profile"})
			return
		}
	}
	s.db.First(user, "id = ?", user.ID)
	c.JSON(http.StatusOK, gin.H{"user": user})
}
