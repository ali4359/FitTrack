package store

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ali4359/fittrack/backend/internal/models"
)

// Open connects to Postgres when DATABASE_URL is set, otherwise falls back to a
// local SQLite file so the stub runs with zero setup.
func Open() *gorm.DB {
	cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}

	var db *gorm.DB
	var err error

	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		db, err = gorm.Open(postgres.Open(dsn), cfg)
		log.Println("store: using postgres")
	} else {
		path := os.Getenv("SQLITE_PATH")
		if path == "" {
			path = "fittrack.db"
		}
		db, err = gorm.Open(sqlite.Open(path), cfg)
		log.Printf("store: using sqlite (%s)", path)
	}
	if err != nil {
		log.Fatalf("store: connect: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Exercise{},
		&models.WorkoutDay{},
		&models.WorkoutDayExercise{},
		&models.WorkoutLog{},
		&models.WorkoutExerciseLog{},
		&models.WorkoutSetLog{},
		&models.MealEntry{},
		&models.MealLog{},
	); err != nil {
		log.Fatalf("store: migrate: %v", err)
	}

	seed(db)
	return db
}

// exercisePhysics holds the per-exercise terms the calorie model needs:
// {rangeOfMotionM, bodyweightLoadFactor}. Kept next to the seed so new
// exercises get values at the same time.
var exercisePhysics = map[string][2]float64{
	"ex-bench":       {0.40, 0},
	"ex-incline-db":  {0.40, 0},
	"ex-cable-fly":   {0.55, 0},
	"ex-tricep-push": {0.28, 0},
	"ex-pullup":      {0.55, 0.90},
	"ex-row":         {0.45, 0},
	"ex-squat":       {0.50, 0.65},
	"ex-rdl":         {0.40, 0.35},
}

// backfillExercisePhysics fills range-of-motion / bodyweight-load values on
// exercises that predate those columns. Idempotent — only touches rows still
// at the zero default.
func backfillExercisePhysics(db *gorm.DB) {
	for id, p := range exercisePhysics {
		db.Model(&models.Exercise{}).
			Where("id = ? AND range_of_motion_m = 0", id).
			Updates(map[string]any{"range_of_motion_m": p[0], "bodyweight_load_factor": p[1]})
	}
}

func seed(db *gorm.DB) {
	var n int64
	db.Model(&models.Exercise{}).Count(&n)
	if n > 0 {
		backfillExercisePhysics(db)
		return
	}
	log.Println("store: seeding demo data")

	exercises := []models.Exercise{
		{ID: "ex-bench", Name: "Barbell Bench Press", MuscleGroup: "chest", MetValue: 5.0},
		{ID: "ex-incline-db", Name: "Incline Dumbbell Press", MuscleGroup: "chest", MetValue: 5.0},
		{ID: "ex-cable-fly", Name: "Cable Fly", MuscleGroup: "chest", MetValue: 3.5},
		{ID: "ex-tricep-push", Name: "Tricep Pushdown", MuscleGroup: "triceps", MetValue: 3.5},
		{ID: "ex-pullup", Name: "Pull-up", MuscleGroup: "back", MetValue: 8.0},
		{ID: "ex-row", Name: "Barbell Row", MuscleGroup: "back", MetValue: 6.0},
		{ID: "ex-squat", Name: "Back Squat", MuscleGroup: "legs", MetValue: 6.0},
		{ID: "ex-rdl", Name: "Romanian Deadlift", MuscleGroup: "legs", MetValue: 6.0},
	}
	for i := range exercises {
		if p, ok := exercisePhysics[exercises[i].ID]; ok {
			exercises[i].RangeOfMotionM = p[0]
			exercises[i].BodyweightLoadFactor = p[1]
		}
	}
	db.Create(&exercises)

	chestDay := models.WorkoutDay{
		ID: "day-1", Name: "Day 1 — Chest & Triceps", DayOrder: 1,
		Exercises: []models.WorkoutDayExercise{
			{ID: uuid.NewString(), ExerciseID: "ex-bench", DefaultSets: 4, DefaultReps: 8, Position: 0},
			{ID: uuid.NewString(), ExerciseID: "ex-incline-db", DefaultSets: 3, DefaultReps: 10, Position: 1},
			{ID: uuid.NewString(), ExerciseID: "ex-cable-fly", DefaultSets: 3, DefaultReps: 12, Position: 2},
			{ID: uuid.NewString(), ExerciseID: "ex-tricep-push", DefaultSets: 3, DefaultReps: 12, Position: 3},
		},
	}
	backDay := models.WorkoutDay{
		ID: "day-2", Name: "Day 2 — Back & Biceps", DayOrder: 2,
		Exercises: []models.WorkoutDayExercise{
			{ID: uuid.NewString(), ExerciseID: "ex-pullup", DefaultSets: 4, DefaultReps: 8, Position: 0},
			{ID: uuid.NewString(), ExerciseID: "ex-row", DefaultSets: 4, DefaultReps: 10, Position: 1},
		},
	}
	legDay := models.WorkoutDay{
		ID: "day-3", Name: "Day 3 — Legs", DayOrder: 3,
		Exercises: []models.WorkoutDayExercise{
			{ID: uuid.NewString(), ExerciseID: "ex-squat", DefaultSets: 4, DefaultReps: 8, Position: 0},
			{ID: uuid.NewString(), ExerciseID: "ex-rdl", DefaultSets: 3, DefaultReps: 10, Position: 1},
		},
	}
	db.Create(&[]models.WorkoutDay{chestDay, backDay, legDay})

	meals := []models.MealEntry{
		{ID: "meal-chana", DishName: "Chana Chaat", Region: "Lahore", BudgetTier: "low", Calories: 420, ProteinG: 22, CarbsG: 58, FatG: 10, GoalTags: "cut,maintain,post-workout", Halal: true, Vegetarian: true},
		{ID: "meal-chapli", DishName: "Chapli Kebab with Roti", Region: "Peshawar", BudgetTier: "mid", Calories: 720, ProteinG: 41, CarbsG: 52, FatG: 38, GoalTags: "bulk,post-workout", Halal: true, Vegetarian: false},
		{ID: "meal-daal-roti", DishName: "Daal Chawal", Region: "Karachi", BudgetTier: "low", Calories: 560, ProteinG: 19, CarbsG: 92, FatG: 12, GoalTags: "bulk,maintain", Halal: true, Vegetarian: true},
		{ID: "meal-nihari", DishName: "Beef Nihari with Naan", Region: "Karachi", BudgetTier: "mid", Calories: 850, ProteinG: 45, CarbsG: 70, FatG: 45, GoalTags: "bulk", Halal: true, Vegetarian: false},
		{ID: "meal-anda-shami", DishName: "Anda Shami Wrap", Region: "Lahore", BudgetTier: "low", Calories: 480, ProteinG: 28, CarbsG: 40, FatG: 22, GoalTags: "maintain,post-workout,breakfast", Halal: true, Vegetarian: false},
		{ID: "meal-sabzi", DishName: "Aloo Palak with Roti", Region: "Lahore", BudgetTier: "low", Calories: 390, ProteinG: 12, CarbsG: 55, FatG: 14, GoalTags: "cut,maintain", Halal: true, Vegetarian: true},
		{ID: "meal-chicken-karahi", DishName: "Chicken Karahi with Roti", Region: "Lahore", BudgetTier: "high", Calories: 680, ProteinG: 52, CarbsG: 34, FatG: 36, GoalTags: "bulk,cut,post-workout,dinner", Halal: true, Vegetarian: false},
	}
	db.Create(&meals)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	demo := models.User{
		ID: uuid.NewString(), Email: "demo@fittrack.app", Name: "Demo Lifter",
		PasswordHash: string(hash), Goal: "bulk", Sex: "male", Age: 27,
		WeightKg: 78, HeightCm: 176,
		Region: "Lahore", BudgetTier: "mid", Restrictions: "halal",
	}
	db.Create(&demo)

	// a couple of past logs so Progress isn't empty
	db.Create(&[]models.WorkoutLog{
		{ID: uuid.NewString(), UserID: demo.ID, WorkoutDayID: "day-1", CompletedAt: time.Now().AddDate(0, 0, -3), DurationMinutes: 52, CaloriesBurned: 331},
		{ID: uuid.NewString(), UserID: demo.ID, WorkoutDayID: "day-2", CompletedAt: time.Now().AddDate(0, 0, -1), DurationMinutes: 47, CaloriesBurned: 402},
	})
}
