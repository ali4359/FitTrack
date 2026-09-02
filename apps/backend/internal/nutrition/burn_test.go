package nutrition

import (
	"testing"
	"time"
)

// a realistic push day for a 75 kg lifter, no timestamps captured.
func pushDay() []BurnExerciseInput {
	return []BurnExerciseInput{
		{ROMMeters: 0.40, Sets: reps(4, 8, 60)},  // bench
		{ROMMeters: 0.40, Sets: reps(3, 10, 40)}, // incline db
		{ROMMeters: 0.55, Sets: reps(3, 12, 15)}, // cable fly
		{ROMMeters: 0.28, Sets: reps(3, 12, 25)}, // pushdown
	}
}

func reps(sets, n int, kg float64) []BurnSetInput {
	out := make([]BurnSetInput, sets)
	for i := range out {
		out[i] = BurnSetInput{Reps: n, WeightKg: kg}
	}
	return out
}

func TestWorkoutBurn_NoBodyWeight(t *testing.T) {
	got := WorkoutBurn(0, 50, pushDay())
	if got.TotalKcal != 0 {
		t.Fatalf("want 0 kcal with unknown body weight, got %v", got.TotalKcal)
	}
	if len(got.Exercises) != 4 {
		t.Fatalf("want per-exercise slice aligned to input, got %d", len(got.Exercises))
	}
}

func TestWorkoutBurn_PlausibleRange(t *testing.T) {
	got := WorkoutBurn(75, 50, pushDay())
	if got.TotalKcal < 140 || got.TotalKcal > 260 {
		t.Fatalf("50-min push day should land ~150-240 kcal, got %v", got.TotalKcal)
	}

	var sum float64
	for _, e := range got.Exercises {
		if e.Kcal <= 0 {
			t.Fatalf("every exercise should carry some burn, got %v", e.Kcal)
		}
		sum += e.Kcal
	}
	if diff := sum - got.TotalKcal; diff < -2 || diff > 2 {
		t.Fatalf("per-exercise kcal should sum to total (±rounding), sum=%v total=%v", sum, got.TotalKcal)
	}
}

func TestWorkoutBurn_HeavierMeansMore(t *testing.T) {
	light := WorkoutBurn(75, 50, []BurnExerciseInput{{ROMMeters: 0.40, Sets: reps(4, 8, 40)}})
	heavy := WorkoutBurn(75, 50, []BurnExerciseInput{{ROMMeters: 0.40, Sets: reps(4, 8, 100)}})
	if heavy.TotalKcal <= light.TotalKcal {
		t.Fatalf("heavier load should burn more: light=%v heavy=%v", light.TotalKcal, heavy.TotalKcal)
	}
}

func TestWorkoutBurn_BodyweightExercise(t *testing.T) {
	// pull-ups: no external weight, but body mass moves through the ROM.
	got := WorkoutBurn(80, 20, []BurnExerciseInput{
		{ROMMeters: 0.55, BodyweightLoadFactor: 0.90, Sets: reps(4, 8, 0)},
	})
	if got.Exercises[0].Kcal <= 0 {
		t.Fatalf("bodyweight work should still register a burn, got %v", got.Exercises[0].Kcal)
	}
}

func TestWorkoutBurn_MeasuredRestUsedOverFallback(t *testing.T) {
	base := time.Date(2026, 9, 2, 18, 0, 0, 0, time.UTC)
	// 3 sets of 10, completed 30s TUT + 30s rest apart -> dense session
	dense := []BurnExerciseInput{{ROMMeters: 0.40, Sets: []BurnSetInput{
		{Reps: 10, WeightKg: 60, CompletedAt: at(base.Add(30 * time.Second))},
		{Reps: 10, WeightKg: 60, CompletedAt: at(base.Add(90 * time.Second))},
		{Reps: 10, WeightKg: 60, CompletedAt: at(base.Add(150 * time.Second))},
	}}}
	// same sets, 3 min apart -> long rests, capped
	slow := []BurnExerciseInput{{ROMMeters: 0.40, Sets: []BurnSetInput{
		{Reps: 10, WeightKg: 60, CompletedAt: at(base.Add(30 * time.Second))},
		{Reps: 10, WeightKg: 60, CompletedAt: at(base.Add(210 * time.Second))},
		{Reps: 10, WeightKg: 60, CompletedAt: at(base.Add(390 * time.Second))},
	}}}

	d := WorkoutBurn(75, 30, dense)
	s := WorkoutBurn(75, 30, slow)
	if s.TotalKcal <= d.TotalKcal {
		t.Fatalf("longer measured rests should raise the baseline burn: dense=%v slow=%v", d.TotalKcal, s.TotalKcal)
	}
	if d.Exercises[0].Sets[0].RestSeconds != 30 {
		t.Fatalf("want measured 30s rest after set 1, got %d", d.Exercises[0].Sets[0].RestSeconds)
	}
}

func at(t time.Time) *time.Time { return &t }
