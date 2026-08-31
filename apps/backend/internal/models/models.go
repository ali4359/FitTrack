package models

import "time"

// User mirrors packages/shared/src/models.ts `User`.
type User struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Goal         string    `json:"goal"`       // bulk | cut | maintain
	WeightKg     float64   `json:"weightKg"`
	HeightCm     float64   `json:"heightCm"`
	Region       string    `json:"region"`
	BudgetTier   string    `json:"budgetTier"` // low | mid | high
	Restrictions string    `json:"restrictions"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

// Exercise mirrors `Exercise`.
type Exercise struct {
	ID          string  `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	MuscleGroup string  `json:"muscleGroup"`
	MetValue    float64 `json:"metValue"`
}

// WorkoutDay mirrors `WorkoutDay`.
type WorkoutDay struct {
	ID        string               `gorm:"primaryKey" json:"id"`
	Name      string               `json:"name"`
	DayOrder  int                  `json:"dayOrder"`
	Exercises []WorkoutDayExercise `gorm:"foreignKey:WorkoutDayID" json:"exercises"`
}

// WorkoutDayExercise is the join row; serialized as { exercise, defaultSets, defaultReps }.
type WorkoutDayExercise struct {
	ID           string   `gorm:"primaryKey" json:"-"`
	WorkoutDayID string   `json:"-"`
	ExerciseID   string   `json:"-"`
	Exercise     Exercise `gorm:"foreignKey:ExerciseID" json:"exercise"`
	DefaultSets  int      `json:"defaultSets"`
	DefaultReps  int      `json:"defaultReps"`
	Position     int      `json:"-"`
}

// WorkoutLog mirrors `WorkoutLog`.
type WorkoutLog struct {
	ID              string    `gorm:"primaryKey" json:"id"`
	UserID          string    `gorm:"index" json:"-"`
	WorkoutDayID    string    `json:"workoutDayId"`
	CompletedAt     time.Time `json:"completedAt"`
	DurationMinutes int       `json:"durationMinutes"`
	CaloriesBurned  float64   `json:"caloriesBurned"`
}

// MealEntry mirrors `MealEntry`.
type MealEntry struct {
	ID         string  `gorm:"primaryKey" json:"id"`
	DishName   string  `json:"dishName"`
	Region     string  `json:"region"`
	BudgetTier string  `json:"budgetTier"`
	Calories   float64 `json:"calories"`
	ProteinG   float64 `json:"proteinG"`
	CarbsG     float64 `json:"carbsG"`
	FatG       float64 `json:"fatG"`
	GoalTags   string  `json:"goalTags"`
	Halal      bool    `json:"halal"`
	Vegetarian bool    `json:"vegetarian"`
}

// MealLog records that a user ate a suggested meal.
type MealLog struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	UserID      string    `gorm:"index" json:"-"`
	MealEntryID string    `json:"mealEntryId"`
	MealType    string    `json:"mealType"`
	EatenAt     time.Time `json:"eatenAt"`
}
