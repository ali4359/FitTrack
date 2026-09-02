package meals

import (
	"math"
	"strings"
)

// macroSanityTolerance is how far a dish's stated calories may drift from the
// Atwater (4/4/9) sum of its macros before we treat the numbers as unreliable.
const macroSanityTolerance = 0.15

// forbiddenIngredients maps a diet style / allergen keyword to substrings that
// must NOT appear in an ingredient list. This is the code-side backstop — we do
// not trust the model's own isVegetarian / isHalal flags on their own.
var haramSubstrings = []string{"pork", "bacon", "ham", "lard", "gelatin", "wine", "beer", "rum", "alcohol", "prosciutto", "pancetta", "chorizo"}
var meatSubstrings = []string{"chicken", "beef", "mutton", "lamb", "goat", "fish", "prawn", "shrimp", "keema", "qeema", "mince", "kebab", "tikka", "nihari", "haleem", "gosht", "boti", "seekh"}
var animalSubstrings = []string{"egg", "milk", "yogurt", "yoghurt", "dahi", "cheese", "paneer", "butter", "ghee", "cream", "honey", "khoya", "malai"}

var allergenSubstrings = map[string][]string{
	"egg":       {"egg", "anda", "mayonnaise", "mayo"},
	"dairy":     {"milk", "yogurt", "yoghurt", "dahi", "cheese", "paneer", "butter", "ghee", "cream", "khoya", "malai", "lassi"},
	"beef":      {"beef", "steak", "veal"},
	"gluten":    {"wheat", "roti", "naan", "chapati", "paratha", "bread", "flour", "atta", "maida", "pasta", "noodle", "suji", "semolina"},
	"peanut":    {"peanut", "groundnut", "moongphali"},
	"treenut":   {"almond", "cashew", "pistachio", "walnut", "badam", "kaju", "pista"},
	"soy":       {"soy", "soya", "tofu", "edamame"},
	"shellfish": {"prawn", "shrimp", "crab", "lobster", "shellfish"},
}

// macroSanity reports whether 4·protein + 4·carbs + 9·fat is within tolerance of
// the stated calories. Catches invented numbers that don't add up.
func macroSanity(s Suggestion) bool {
	if s.Calories <= 0 {
		return false
	}
	atwater := 4*s.ProteinG + 4*s.CarbsG + 9*s.FatG
	return math.Abs(atwater-s.Calories)/s.Calories <= macroSanityTolerance
}

func ingredientsBlob(s Suggestion) string {
	return strings.ToLower(s.Name + " " + strings.Join(s.Ingredients, " "))
}

func containsAny(blob string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(blob, n) {
			return true
		}
	}
	return false
}

// dietOK enforces the user's hard dietary rule against the ingredient list,
// independent of the model's self-reported flags.
func dietOK(s Suggestion, diet DietStyle) bool {
	blob := ingredientsBlob(s)
	switch diet {
	case DietHalal:
		return !containsAny(blob, haramSubstrings)
	case DietVegetarian:
		return !containsAny(blob, meatSubstrings)
	case DietVegan:
		return !containsAny(blob, meatSubstrings) && !containsAny(blob, animalSubstrings)
	default: // omnivore
		return true
	}
}

// allergensOK reports whether the dish avoids every listed allergen.
func allergensOK(s Suggestion, allergens []string) bool {
	blob := ingredientsBlob(s)
	for _, a := range allergens {
		needles, known := allergenSubstrings[strings.ToLower(strings.TrimSpace(a))]
		if !known {
			// Unknown allergen: fall back to a literal match on the word.
			needles = []string{strings.ToLower(strings.TrimSpace(a))}
		}
		if containsAny(blob, needles) {
			return false
		}
	}
	return true
}

// calorieFit rejects a dish that blows past the meal's calorie target. A cut
// user gets a tight ceiling; a bulk user is allowed to overshoot more.
func calorieFit(s Suggestion, target float64, goal string) bool {
	if target <= 0 {
		return s.Calories <= 500 // "already at target" case: light only
	}
	ceil := 1.30
	switch goal {
	case "cut":
		ceil = 1.15
	case "bulk":
		ceil = 1.60
	}
	return s.Calories <= target*ceil
}

// guard runs every check and returns the dishes that pass, de-duplicated and
// diversified (no two sharing the same primary protein + starch).
func guard(in []Suggestion, req Request) []Suggestion {
	seen := map[string]bool{}
	for _, e := range req.Exclude {
		seen[normalizeName(e)] = true
	}

	var kept []Suggestion
	var primaries []string
	for _, s := range in {
		if s.Name == "" || len(s.Ingredients) == 0 {
			continue
		}
		if !macroSanity(s) {
			continue
		}
		if !dietOK(s, req.Diet) {
			continue
		}
		if !allergensOK(s, req.Allergens) {
			continue
		}
		if !calorieFit(s, req.Target.Calories, req.Goal) {
			continue
		}
		key := normalizeName(s.Name)
		if seen[key] {
			continue
		}
		if p := primaryIngredient(s); p != "" && contains(primaries, p) {
			continue // already have a dish built on this protein/starch
		}

		seen[key] = true
		if p := primaryIngredient(s); p != "" {
			primaries = append(primaries, p)
		}
		s.Estimated = true
		kept = append(kept, s)
	}
	return kept
}

func normalizeName(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(s))), " ")
}

// primaryIngredient returns the first recognised protein or starch keyword in
// the dish, used only for the diversity check.
func primaryIngredient(s Suggestion) string {
	blob := ingredientsBlob(s)
	for _, k := range []string{"chicken", "beef", "mutton", "lamb", "goat", "fish", "prawn", "egg", "paneer", "lentil", "daal", "dal", "chana", "chickpea", "rajma", "kidney bean"} {
		if strings.Contains(blob, k) {
			return k
		}
	}
	return ""
}

func contains(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
