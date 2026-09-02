package nutrition

import (
	"math"
	"time"
)

// Resistance-training energy-burn model. The estimate is the sum of three
// buckets:
//
//  1. Mechanical work — moving the load through its range of motion, rep by
//     rep, converted to metabolic energy via muscular efficiency. This is the
//     part that scales with the weight and reps the user actually logged.
//  2. Baseline metabolism — the calories burned just being at the gym, charged
//     over the time under tension plus the rest taken between sets. Rest is
//     measured from the per-set completion timestamps the app captures, and
//     falls back to an even split of the leftover session time for any set
//     that has no timestamp.
//  3. EPOC — the post-exercise "afterburn", taken as a fixed fraction of 1+2.
//
// Every value here is approximate: without heart-rate or motion data a
// resistance-session estimate is good to maybe +/-25%. The model is pure so it
// can be unit-tested and lifted into the real backend unchanged.
const (
	gravity           = 9.81   // m/s^2
	eccentricFactor   = 0.30   // lowering the load costs ~30% of the energy of lifting it
	liftingEfficiency = 0.22   // net mechanical efficiency of weight training (20-25%)
	joulesPerKcal     = 4184.0 // J per dietary kcal
	secondsPerRep     = 3.0    // TUT estimate: ~2s concentric + slower eccentric
	restMET           = 2.5    // standing recovery between sets
	maxRestSecPerSet  = 150.0  // rest cap per set, so a forgotten timer can't inflate the burn
	epocFraction      = 0.10   // afterburn as a fraction of session energy
	defaultROMMeters  = 0.45   // used when the exercise has no catalogued range of motion
)

// BurnSetInput is one performed set.
type BurnSetInput struct {
	Reps        int
	WeightKg    float64
	CompletedAt *time.Time // when the user tapped "set done"; nil if not captured
}

// BurnExerciseInput is one exercise in the session plus its catalogue params.
type BurnExerciseInput struct {
	ROMMeters            float64 // range of motion; <=0 falls back to defaultROMMeters
	BodyweightLoadFactor float64 // fraction of body mass moved through the ROM (≈0 for most isolation work)
	Sets                 []BurnSetInput
}

// BurnSetResult is the resolved per-set timing, echoed back so it can be stored.
type BurnSetResult struct {
	TUTSeconds  int
	RestSeconds int // rest taken after this set
}

// BurnExerciseResult is one exercise's share of the session burn.
type BurnExerciseResult struct {
	Kcal float64
	Sets []BurnSetResult
}

// BurnResult is the whole-session estimate and its per-exercise breakdown.
type BurnResult struct {
	TotalKcal float64
	Exercises []BurnExerciseResult // aligned to the exercises passed in
}

// WorkoutBurn estimates the calories burned in one resistance-training session.
// exercises must be in the order they were performed. Returns a zeroed result
// when body weight is unknown (<= 0) — callers should require it upstream.
func WorkoutBurn(bodyWeightKg float64, durationMinutes int, exercises []BurnExerciseInput) BurnResult {
	res := BurnResult{Exercises: make([]BurnExerciseResult, len(exercises))}
	for i := range exercises {
		res.Exercises[i].Sets = make([]BurnSetResult, len(exercises[i].Sets))
	}
	if bodyWeightKg <= 0 {
		return res
	}

	tut, rest := resolveTiming(exercises, durationMinutes)

	// baseline metabolic rate in kcal per second, from the resting-recovery MET
	baselineRate := restMET * 3.5 * bodyWeightKg / 200 / 60

	active := make([]float64, len(exercises))
	baseline := make([]float64, len(exercises))
	var sumActive, sumBaseline float64

	for i, ex := range exercises {
		rom := ex.ROMMeters
		if rom <= 0 {
			rom = defaultROMMeters
		}
		for j, set := range ex.Sets {
			res.Exercises[i].Sets[j] = BurnSetResult{
				TUTSeconds:  int(math.Round(tut[i][j])),
				RestSeconds: int(math.Round(rest[i][j])),
			}
			if set.Reps <= 0 {
				continue
			}
			moved := set.WeightKg + ex.BodyweightLoadFactor*bodyWeightKg
			if moved > 0 {
				workJ := float64(set.Reps) * moved * gravity * rom * (1 + eccentricFactor)
				active[i] += workJ / liftingEfficiency / joulesPerKcal
			}
			baseline[i] += (tut[i][j] + rest[i][j]) * baselineRate
		}
		sumActive += active[i]
		sumBaseline += baseline[i]
	}

	work := sumActive + sumBaseline
	epoc := epocFraction * work

	for i := range exercises {
		kcal := active[i] + baseline[i]
		if work > 0 {
			kcal += epoc * (active[i] + baseline[i]) / work
		}
		res.Exercises[i].Kcal = round(kcal)
	}
	res.TotalKcal = round(work + epoc)
	return res
}

// resolveTiming returns time-under-tension and post-set rest (seconds) for every
// set, indexed [exercise][set]. TUT is always estimated from the rep count — one
// "set done" tap can't measure it. Rest is the gap between consecutive set
// completions minus the next set's TUT; any set without usable timestamps gets
// an even share of the session time the logged TUT doesn't account for (which is
// how the whole session behaves when no timestamps are captured at all).
func resolveTiming(exercises []BurnExerciseInput, durationMinutes int) (tut, rest [][]float64) {
	tut = make([][]float64, len(exercises))
	rest = make([][]float64, len(exercises))

	type coord struct{ i, j int }
	var order []coord
	var totalTUT float64
	for i, ex := range exercises {
		tut[i] = make([]float64, len(ex.Sets))
		rest[i] = make([]float64, len(ex.Sets))
		for j, set := range ex.Sets {
			tut[i][j] = math.Max(0, float64(set.Reps)) * secondsPerRep
			totalTUT += tut[i][j]
			order = append(order, coord{i, j})
		}
	}
	if len(order) == 0 {
		return tut, rest
	}

	fallbackRest := clamp((float64(durationMinutes)*60-totalTUT)/float64(len(order)), 0, maxRestSecPerSet)
	at := func(c coord) *time.Time { return exercises[c.i].Sets[c.j].CompletedAt }

	for k, c := range order {
		r := fallbackRest
		switch {
		case k+1 < len(order) && at(c) != nil && at(order[k+1]) != nil:
			next := order[k+1]
			r = clamp(at(next).Sub(*at(c)).Seconds()-tut[next.i][next.j], 0, maxRestSecPerSet)
		case k+1 == len(order) && at(c) != nil:
			r = 0 // last set: no measurable trailing rest
		}
		rest[c.i][c.j] = r
	}
	return tut, rest
}
