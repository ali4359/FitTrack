package meals

import (
	"testing"

	"github.com/ali4359/fittrack/backend/internal/nutrition"
)

func dish(name string, kcal, p, c, f float64, ings ...string) Suggestion {
	return Suggestion{Name: name, Portion: "1 plate", Ingredients: ings, Calories: kcal, ProteinG: p, CarbsG: c, FatG: f}
}

func TestMacroSanity(t *testing.T) {
	// 40*4 + 50*4 + 20*9 = 160 + 200 + 180 = 540
	if !macroSanity(dish("ok", 540, 40, 50, 20, "chicken")) {
		t.Fatal("dish whose macros reconcile should pass")
	}
	// Same macros, claimed 900 kcal — 4/4/9 says 540, off by 40%.
	if macroSanity(dish("bad", 900, 40, 50, 20, "chicken")) {
		t.Fatal("dish with impossible calorie count should fail")
	}
	if macroSanity(dish("zero", 0, 0, 0, 0, "water")) {
		t.Fatal("zero-calorie dish should fail")
	}
}

func TestDietOK(t *testing.T) {
	veg := dish("Daal Chawal", 560, 19, 92, 12, "lentils", "rice", "onion")
	meat := dish("Chicken Karahi", 640, 52, 20, 40, "chicken", "tomato", "ghee")
	pork := dish("Pork Ribs", 700, 45, 5, 55, "pork ribs", "bbq sauce")

	if !dietOK(veg, DietVegetarian) {
		t.Fatal("daal chawal should be vegetarian-OK")
	}
	if dietOK(meat, DietVegetarian) {
		t.Fatal("chicken karahi must fail vegetarian even if the model claimed otherwise")
	}
	if dietOK(pork, DietHalal) {
		t.Fatal("pork must fail halal")
	}
	if !dietOK(meat, DietHalal) {
		t.Fatal("chicken karahi should pass halal")
	}
	if dietOK(dish("Anda Paratha", 480, 18, 40, 26, "egg", "wheat flour", "butter"), DietVegan) {
		t.Fatal("egg + butter must fail vegan")
	}
}

func TestModelLiedAboutFlags(t *testing.T) {
	// Model says isVegetarian:true but the ingredients contain mince.
	liar := Suggestion{
		Name: "Keema Matar", Portion: "1 cup", Ingredients: []string{"beef mince", "peas", "onion"},
		Calories: 400, ProteinG: 30, CarbsG: 15, FatG: 22, IsVegetarian: true,
	}
	got := guard([]Suggestion{liar}, Request{Diet: DietVegetarian, Target: nutrition.Macros{Calories: 500}})
	if len(got) != 0 {
		t.Fatal("guard must drop a meat dish regardless of the model's isVegetarian flag")
	}
}

func TestAllergensOK(t *testing.T) {
	d := dish("Chicken Karahi with Roti", 640, 52, 44, 28, "chicken", "tomato", "roti")
	if allergensOK(d, []string{"gluten"}) {
		t.Fatal("roti should trip the gluten allergen")
	}
	if !allergensOK(d, []string{"peanut"}) {
		t.Fatal("no peanuts here")
	}
}

func TestCalorieFit(t *testing.T) {
	d := dish("Nihari", 850, 45, 70, 45, "beef", "flour")
	if calorieFit(d, 500, "cut") {
		t.Fatal("850 kcal is way over a 500 kcal cut target")
	}
	if !calorieFit(d, 700, "bulk") {
		t.Fatal("850 kcal should be allowed against a 700 kcal bulk target")
	}
}

func TestGuardDiversifies(t *testing.T) {
	a := dish("Chicken Karahi", 600, 45, 30, 32, "chicken", "tomato")
	b := dish("Chicken Pulao", 620, 40, 60, 22, "chicken", "rice")
	c := dish("Daal Chawal", 560, 19, 92, 12, "lentils", "rice")
	req := Request{Diet: DietOmnivore, Target: nutrition.Macros{Calories: 650}, Goal: "maintain"}

	got := guard([]Suggestion{a, b, c}, req)
	if len(got) != 2 {
		t.Fatalf("expected the two chicken dishes to collapse to one, got %d", len(got))
	}
}

func TestGuardRespectsExclude(t *testing.T) {
	a := dish("Chicken Karahi", 600, 45, 30, 32, "chicken", "tomato")
	req := Request{Diet: DietOmnivore, Target: nutrition.Macros{Calories: 650}, Exclude: []string{"chicken karahi"}}
	if len(guard([]Suggestion{a}, req)) != 0 {
		t.Fatal("excluded dish should not come back")
	}
}
