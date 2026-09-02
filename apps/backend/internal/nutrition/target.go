// Package nutrition computes calorie and macro targets from a user's profile.
// All functions are pure so they can be unit-tested and lifted into the real
// backend unchanged.
package nutrition

import "math"

// Macros is a calorie + macronutrient vector (grams).
type Macros struct {
	Calories float64 `json:"calories"`
	ProteinG float64 `json:"proteinG"`
	CarbsG   float64 `json:"carbsG"`
	FatG     float64 `json:"fatG"`
}

// Profile is the subset of the user record the math needs.
type Profile struct {
	Sex      string  // "male" | "female" | "" (unknown -> midpoint estimate)
	Age      int     // years; 0 -> assume 30
	WeightKg float64 // 0 -> assume 75
	HeightCm float64 // 0 -> assume 170
	Goal     string  // "cut" | "bulk" | "maintain"
}

func (p Profile) weight() float64 {
	if p.WeightKg > 0 {
		return p.WeightKg
	}
	return 75
}

func (p Profile) height() float64 {
	if p.HeightCm > 0 {
		return p.HeightCm
	}
	return 170
}

func (p Profile) age() float64 {
	if p.Age > 0 {
		return float64(p.Age)
	}
	return 30
}

// BMR is the Mifflin–St Jeor basal metabolic rate. With unknown sex we average
// the male (+5) and female (−161) constants, i.e. −78.
func BMR(p Profile) float64 {
	base := 10*p.weight() + 6.25*p.height() - 5*p.age()
	switch p.Sex {
	case "male":
		return base + 5
	case "female":
		return base - 161
	default:
		return base - 78
	}
}

// sedentaryFactor is the baseline activity multiplier applied to BMR. Workout
// burn is added explicitly on top (see DailyTarget) rather than baked in here,
// because the app already measures it per session.
const sedentaryFactor = 1.2

// goalDelta returns the fractional calorie adjustment for a goal.
func goalDelta(goal string) float64 {
	switch goal {
	case "cut":
		return -0.20
	case "bulk":
		return 0.12
	default:
		return 0
	}
}

// DailyTarget is the calorie + macro goal for the whole day, including the
// calories burned in workouts so far today.
func DailyTarget(p Profile, workoutBurnToday float64) Macros {
	maintenance := BMR(p)*sedentaryFactor + workoutBurnToday
	kcal := maintenance * (1 + goalDelta(p.Goal))

	// Protein: g/kg bodyweight. Higher on a cut to protect lean mass.
	proteinPerKg := 1.8
	if p.Goal == "cut" {
		proteinPerKg = 2.2
	}
	protein := proteinPerKg * p.weight()

	// Fat: ~0.9 g/kg, floored so very light users still get enough.
	fat := math.Max(0.9*p.weight(), 40)

	// Carbs fill whatever calories remain.
	carbs := math.Max((kcal-4*protein-9*fat)/4, 0)

	return Macros{Calories: round(kcal), ProteinG: round(protein), CarbsG: round(carbs), FatG: round(fat)}
}

// Remaining subtracts what's already been eaten from the daily target. Values
// are clamped at zero — a negative "remaining" just means the user is at or over
// target, which the meal logic handles as a "light option" case.
func Remaining(target, eaten Macros) Macros {
	return Macros{
		Calories: math.Max(target.Calories-eaten.Calories, 0),
		ProteinG: math.Max(target.ProteinG-eaten.ProteinG, 0),
		CarbsG:   math.Max(target.CarbsG-eaten.CarbsG, 0),
		FatG:     math.Max(target.FatG-eaten.FatG, 0),
	}
}

// MealTarget is the macro vector a single suggested meal should aim for.
//
// mealsLeft is how many eating occasions remain today (>=1). isPostWorkout
// raises the protein and carb floors for glycogen refill and muscle protein
// synthesis, and does not push fat up.
func MealTarget(p Profile, remaining Macros, mealsLeft int, isPostWorkout bool) Macros {
	if mealsLeft < 1 {
		mealsLeft = 1
	}
	daily := DailyTarget(p, 0).Calories

	share := remaining.Calories / float64(mealsLeft)
	// Keep a single meal within a sane band of the day's intake.
	share = clamp(share, 0.25*daily, 0.45*daily)
	if remaining.Calories == 0 {
		// Already at target: a small, high-protein "top up" only.
		share = math.Min(400, 0.2*daily)
	}

	// Scale the macros to the same fraction of what's left as the calorie share,
	// so the meal's macros stay consistent with its calorie budget even when the
	// share was clamped (e.g. one meal left, nothing eaten yet).
	frac := 1.0 / float64(mealsLeft)
	if remaining.Calories > 0 {
		frac = share / remaining.Calories
	}
	protein := remaining.ProteinG * frac
	carbs := remaining.CarbsG * frac
	fat := remaining.FatG * frac

	if isPostWorkout {
		protein = math.Max(protein, 0.4*p.weight())
		carbs = math.Max(carbs, 0.8*p.weight())
	}

	return Macros{Calories: round(share), ProteinG: round(protein), CarbsG: round(carbs), FatG: round(fat)}
}

func round(x float64) float64 { return math.Round(x) }
func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
