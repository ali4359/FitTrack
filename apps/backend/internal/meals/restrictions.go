package meals

import "strings"

// ParseRestrictions turns the profile's free-form comma-separated restrictions
// string (e.g. "halal,no beef,lactose-free") into a hard diet style plus a list
// of allergen keywords the guardrails understand.
func ParseRestrictions(csv string) (DietStyle, []string) {
	diet := DietOmnivore
	var allergens []string
	seen := map[string]bool{}

	add := func(a string) {
		if a != "" && !seen[a] {
			seen[a] = true
			allergens = append(allergens, a)
		}
	}

	for _, raw := range strings.Split(csv, ",") {
		t := strings.ToLower(strings.TrimSpace(raw))
		if t == "" {
			continue
		}
		switch {
		case t == "vegan" || t == "plant-based":
			diet = DietVegan
		case t == "vegetarian" || t == "veg":
			if diet != DietVegan {
				diet = DietVegetarian
			}
		case t == "halal":
			if diet == DietOmnivore {
				diet = DietHalal
			}
		case strings.Contains(t, "beef"):
			add("beef")
		case strings.Contains(t, "lactose") || strings.Contains(t, "dairy") || t == "no milk":
			add("dairy")
		case strings.Contains(t, "egg"):
			add("egg")
		case strings.Contains(t, "gluten") || strings.Contains(t, "wheat"):
			add("gluten")
		case strings.Contains(t, "peanut"):
			add("peanut")
		case strings.Contains(t, "tree nut") || strings.Contains(t, "treenut") || t == "no nuts" || t == "nut-free":
			add("treenut")
		case strings.Contains(t, "soy"):
			add("soy")
		case strings.Contains(t, "shellfish") || strings.Contains(t, "prawn") || strings.Contains(t, "shrimp"):
			add("shellfish")
		case strings.HasPrefix(t, "no "):
			add(strings.TrimSpace(strings.TrimPrefix(t, "no ")))
		}
	}
	return diet, allergens
}
