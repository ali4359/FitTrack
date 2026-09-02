package models

import "time"

// User mirrors packages/shared/src/models.ts `User`.
type User struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Goal         string    `json:"goal"` // bulk | cut | maintain
	Sex          string    `json:"sex"`  // male | female | "" (unknown)
	Age          int       `json:"age"`  // years; 0 = unknown
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
	// RangeOfMotionM is the bar/handle travel per rep, in metres — the "distance"
	// term of the mechanical-work calorie estimate.
	RangeOfMotionM float64 `json:"rangeOfMotionM"`
	// BodyweightLoadFactor is the fraction of the lifter's body mass that also
	// moves through the ROM (≈0 for most isolation work, ~0.65 for a squat).
	BodyweightLoadFactor float64 `json:"bodyweightLoadFactor"`
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
	ID              string               `gorm:"primaryKey" json:"id"`
	UserID          string               `gorm:"index" json:"-"`
	WorkoutDayID    string               `json:"workoutDayId"`
	CompletedAt     time.Time            `json:"completedAt"`
	DurationMinutes int                  `json:"durationMinutes"`
	CaloriesBurned  float64              `json:"caloriesBurned"`
	Exercises       []WorkoutExerciseLog `gorm:"foreignKey:WorkoutLogID;constraint:OnDelete:CASCADE" json:"exercises,omitempty"`
}

// WorkoutExerciseLog is one exercise within a logged session; serialized as
// { exerciseId, caloriesBurned, sets }.
type WorkoutExerciseLog struct {
	ID           string `gorm:"primaryKey" json:"-"`
	WorkoutLogID string `gorm:"index" json:"-"`
	ExerciseID   string `json:"exerciseId"`
	Position     int    `json:"-"`
	// CaloriesBurned is this exercise's share of the session estimate.
	CaloriesBurned float64         `json:"caloriesBurned"`
	Sets           []WorkoutSetLog `gorm:"foreignKey:WorkoutExerciseLogID;constraint:OnDelete:CASCADE" json:"sets"`
}

// WorkoutSetLog is a single recorded set: reps and weight actually performed,
// plus the timing the app captured or the model estimated.
type WorkoutSetLog struct {
	ID                   string  `gorm:"primaryKey" json:"-"`
	WorkoutExerciseLogID string  `gorm:"index" json:"-"`
	SetNumber            int     `json:"setNumber"`
	Reps                 int     `json:"reps"`
	WeightKg             float64 `json:"weightKg"`
	// CompletedAt is when the user tapped "set done"; nil for older logs or when
	// the set was never marked done.
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	// TUTSeconds / RestSeconds are the resolved timing the burn model used:
	// time under tension, and rest taken after this set.
	TUTSeconds  int `json:"tutSeconds"`
	RestSeconds int `json:"restSeconds"`
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

// MealLog records that a user ate a meal. Macros are snapshotted onto the row so
// LLM-generated suggestions (which never become MealEntry rows) can still be
// logged and summed for "eaten today". MealEntryID stays for seeded-catalogue
// logs; it is empty for LLM suggestions.
type MealLog struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	UserID      string    `gorm:"index" json:"-"`
	MealEntryID string    `json:"mealEntryId"`
	DishName    string    `json:"dishName"`
	MealType    string    `json:"mealType"`
	Servings    float64   `json:"servings"` // portion multiplier; defaults to 1
	Calories    float64   `json:"calories"` // per the logged servings, already multiplied
	ProteinG    float64   `json:"proteinG"`
	CarbsG      float64   `json:"carbsG"`
	FatG        float64   `json:"fatG"`
	EatenAt     time.Time `json:"eatenAt"`
}
