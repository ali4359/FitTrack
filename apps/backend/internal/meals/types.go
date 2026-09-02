// Package meals produces post-workout / mealtime dish suggestions.
//
// The pipeline is: compute a macro target in code (internal/nutrition) -> ask an
// LLM for candidate dishes -> run every candidate through code guardrails
// (macro sanity, dietary restrictions, calorie fit) -> cache the survivors.
// The LLM never computes the target and is never the only thing enforcing a
// dietary restriction.
package meals

import "github.com/ali4359/fittrack/backend/internal/nutrition"

// DietStyle is the user's hard dietary rule.
type DietStyle string

const (
	DietOmnivore   DietStyle = "omnivore"
	DietHalal      DietStyle = "halal"
	DietVegetarian DietStyle = "vegetarian"
	DietVegan      DietStyle = "vegan"
)

// Request is everything the suggester needs for one call.
type Request struct {
	// Target is the macro vector this single meal should aim for.
	Target nutrition.Macros
	// MealType is "post-workout" | "breakfast" | "lunch" | "dinner" | "daily".
	MealType string
	// Goal is the user's training goal ("cut" | "bulk" | "maintain"); it tunes
	// how overshoot vs. undershoot is scored.
	Goal string

	Diet      DietStyle
	Allergens []string // lowercase: "egg", "dairy", "beef", "gluten", "peanut", ...
	Region    string   // free text, e.g. "Lahore" / "Punjab" / "Pakistan"
	Budget    string   // "low" | "mid" | "high"

	// Exclude lists dish names already shown recently, so we don't repeat.
	Exclude []string
}

// Suggestion is one dish returned to the client.
type Suggestion struct {
	Name        string   `json:"name"`
	Portion     string   `json:"portion"`     // human portion, e.g. "1.5 plates (~420g)"
	Ingredients []string `json:"ingredients"` // used by the guardrails, also shown

	Calories float64 `json:"calories"`
	ProteinG float64 `json:"proteinG"`
	CarbsG   float64 `json:"carbsG"`
	FatG     float64 `json:"fatG"`

	IsVegetarian bool `json:"isVegetarian"`
	IsVegan      bool `json:"isVegan"`
	IsHalal      bool `json:"isHalal"`

	// Role labels the slot this dish fills in the returned set of three:
	// "best-fit" | "higher-protein" | "lighter".
	Role string `json:"role"`
	// WhyItFits is a one-line rationale from the model.
	WhyItFits string `json:"whyItFits"`
	// Estimated is always true today — the macros are model estimates, not
	// measured. The client should label them as such.
	Estimated bool `json:"estimated"`
}

// Response is the handler payload.
type Response struct {
	Target    nutrition.Macros `json:"target"`
	Results   []Suggestion     `json:"results"`
	Source    string           `json:"source"`    // "llm" | "llm-cache" | "fallback"
	Broadened bool             `json:"broadened"` // true if constraints were relaxed to fill 3
}
