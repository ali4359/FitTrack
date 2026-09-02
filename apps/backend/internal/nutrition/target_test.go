package nutrition

import (
	"math"
	"testing"
)

func approx(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

func TestBMR(t *testing.T) {
	// Mifflin–St Jeor, 80kg / 180cm / 30y male:
	// 10*80 + 6.25*180 - 5*30 + 5 = 800 + 1125 - 150 + 5 = 1780
	got := BMR(Profile{Sex: "male", Age: 30, WeightKg: 80, HeightCm: 180})
	if !approx(got, 1780, 0.5) {
		t.Fatalf("male BMR = %v, want 1780", got)
	}
	// Female constant is -161: 1780 - 5 - 161 = 1614
	got = BMR(Profile{Sex: "female", Age: 30, WeightKg: 80, HeightCm: 180})
	if !approx(got, 1614, 0.5) {
		t.Fatalf("female BMR = %v, want 1614", got)
	}
	// Unknown sex averages the constants (-78): 1780 - 5 - 78 = 1697
	got = BMR(Profile{Age: 30, WeightKg: 80, HeightCm: 180})
	if !approx(got, 1697, 0.5) {
		t.Fatalf("unknown-sex BMR = %v, want 1697", got)
	}
}

func TestBMRDefaults(t *testing.T) {
	// Empty profile must not panic or return NaN/0.
	got := BMR(Profile{})
	if got <= 0 || math.IsNaN(got) {
		t.Fatalf("default BMR = %v, want a positive number", got)
	}
}

func TestDailyTargetGoalDirection(t *testing.T) {
	p := Profile{Sex: "male", Age: 30, WeightKg: 80, HeightCm: 180}
	cut := DailyTarget(withGoal(p, "cut"), 0)
	maintain := DailyTarget(withGoal(p, "maintain"), 0)
	bulk := DailyTarget(withGoal(p, "bulk"), 0)

	if !(cut.Calories < maintain.Calories && maintain.Calories < bulk.Calories) {
		t.Fatalf("expected cut < maintain < bulk, got %v / %v / %v", cut.Calories, maintain.Calories, bulk.Calories)
	}
	// Cut should raise protein above maintain despite fewer calories.
	if cut.ProteinG <= maintain.ProteinG {
		t.Fatalf("cut protein %v should exceed maintain protein %v", cut.ProteinG, maintain.ProteinG)
	}
}

func TestDailyTargetMacrosReconcile(t *testing.T) {
	p := Profile{Sex: "female", Age: 28, WeightKg: 62, HeightCm: 165, Goal: "maintain"}
	m := DailyTarget(p, 300)
	sum := 4*m.ProteinG + 4*m.CarbsG + 9*m.FatG
	if !approx(sum, m.Calories, m.Calories*0.02) {
		t.Fatalf("4/4/9 sum %v should match calories %v", sum, m.Calories)
	}
}

func TestDailyTargetAddsWorkoutBurn(t *testing.T) {
	p := Profile{Sex: "male", Age: 30, WeightKg: 80, HeightCm: 180, Goal: "maintain"}
	rest := DailyTarget(p, 0)
	trained := DailyTarget(p, 400)
	if trained.Calories <= rest.Calories {
		t.Fatalf("workout burn should raise the target: %v vs %v", trained.Calories, rest.Calories)
	}
}

func TestRemainingClampsAtZero(t *testing.T) {
	target := Macros{Calories: 2000, ProteinG: 150, CarbsG: 200, FatG: 60}
	eaten := Macros{Calories: 2500, ProteinG: 100, CarbsG: 300, FatG: 90}
	r := Remaining(target, eaten)
	if r.Calories != 0 || r.CarbsG != 0 || r.FatG != 0 {
		t.Fatalf("over-target values should clamp to 0, got %+v", r)
	}
	if r.ProteinG != 50 {
		t.Fatalf("protein remaining = %v, want 50", r.ProteinG)
	}
}

func TestMealTargetPostWorkoutFloors(t *testing.T) {
	p := Profile{Sex: "male", Age: 30, WeightKg: 80, HeightCm: 180, Goal: "maintain"}
	remaining := Macros{Calories: 900, ProteinG: 20, CarbsG: 30, FatG: 20}

	plain := MealTarget(p, remaining, 2, false)
	pw := MealTarget(p, remaining, 2, true)

	if pw.ProteinG <= plain.ProteinG {
		t.Fatalf("post-workout protein floor not applied: %v vs %v", pw.ProteinG, plain.ProteinG)
	}
	if pw.ProteinG < 0.4*p.WeightKg {
		t.Fatalf("post-workout protein %v below 0.4 g/kg floor", pw.ProteinG)
	}
	if pw.CarbsG < 0.8*p.WeightKg {
		t.Fatalf("post-workout carbs %v below 0.8 g/kg floor", pw.CarbsG)
	}
}

func TestMealTargetCapsShare(t *testing.T) {
	p := Profile{Sex: "male", Age: 30, WeightKg: 80, HeightCm: 180, Goal: "bulk"}
	// Huge remaining, only one meal left — must not dump the whole day into it.
	remaining := Macros{Calories: 5000, ProteinG: 300, CarbsG: 600, FatG: 150}
	m := MealTarget(p, remaining, 1, false)
	daily := DailyTarget(p, 0).Calories
	if m.Calories > 0.45*daily+1 {
		t.Fatalf("meal share %v exceeds 45%% of daily %v", m.Calories, daily)
	}
}

func withGoal(p Profile, g string) Profile { p.Goal = g; return p }
