package meals

import (
	"context"
	"testing"

	"github.com/ali4359/fittrack/backend/internal/nutrition"
)

type fakeSuggester struct {
	batches [][]Suggestion // one per call
	calls   int
}

func (f *fakeSuggester) Suggest(_ context.Context, _ Request, _ int) ([]Suggestion, error) {
	i := f.calls
	f.calls++
	if i < len(f.batches) {
		return f.batches[i], nil
	}
	return nil, nil
}

func baseReq() Request {
	return Request{
		Target:   nutrition.Macros{Calories: 650, ProteinG: 45, CarbsG: 60, FatG: 20},
		MealType: "post-workout",
		Goal:     "maintain",
		Diet:     DietHalal,
		Region:   "Lahore",
	}
}

func TestServiceHappyPath(t *testing.T) {
	f := &fakeSuggester{batches: [][]Suggestion{{
		dish("Chicken Karahi", 640, 48, 55, 24, "chicken", "tomato"),
		dish("Daal Chawal", 600, 20, 95, 12, "lentils", "rice"),
		dish("Chapli Kebab", 660, 40, 40, 34, "beef mince", "onion"),
		dish("Pork Chop", 600, 40, 5, 45, "pork"), // must be dropped (halal)
	}}}
	svc := NewService(f)

	resp, err := svc.Suggest(context.Background(), baseReq())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Source != "llm" {
		t.Fatalf("source = %q, want llm", resp.Source)
	}
	if len(resp.Results) != 3 {
		t.Fatalf("got %d results, want 3", len(resp.Results))
	}
	for _, r := range resp.Results {
		if r.Name == "Pork Chop" {
			t.Fatal("pork dish leaked past the guardrails")
		}
		if !r.Estimated {
			t.Fatal("results should be flagged Estimated")
		}
	}
	roles := map[string]bool{}
	for _, r := range resp.Results {
		roles[r.Role] = true
	}
	if !roles["best-fit"] || !roles["higher-protein"] || !roles["lighter"] {
		t.Fatalf("missing role labels: %v", roles)
	}
}

func TestServiceServesFromCache(t *testing.T) {
	f := &fakeSuggester{batches: [][]Suggestion{{
		dish("Chicken Karahi", 640, 48, 55, 24, "chicken", "tomato"),
		dish("Daal Chawal", 600, 20, 95, 12, "lentils", "rice"),
		dish("Chapli Kebab", 660, 40, 40, 34, "beef mince", "onion"),
		dish("Anda Curry", 620, 30, 30, 40, "egg", "onion"),
		dish("Fish Curry", 610, 44, 30, 32, "fish", "spices"),
	}}}
	svc := NewService(f)

	if _, err := svc.Suggest(context.Background(), baseReq()); err != nil {
		t.Fatal(err)
	}
	callsAfterFirst := f.calls
	if _, err := svc.Suggest(context.Background(), baseReq()); err != nil {
		t.Fatal(err)
	}
	if f.calls != callsAfterFirst {
		t.Fatalf("second request should have hit the cache, but the LLM was called again (%d -> %d)", callsAfterFirst, f.calls)
	}
}

func TestServiceRetryBroadens(t *testing.T) {
	// First batch: everything overshoots a tight cut target -> all dropped.
	// Second (relaxed) batch: acceptable dishes.
	f := &fakeSuggester{batches: [][]Suggestion{
		{dish("Nihari", 1200, 60, 90, 65, "beef", "flour")},
		{
			dish("Grilled Chicken + Salad", 520, 55, 20, 22, "chicken", "lettuce"),
			dish("Daal + Rice", 540, 22, 80, 12, "lentils", "rice"),
			dish("Fish Tikka", 500, 48, 8, 28, "fish", "yogurt"),
		},
	}}
	svc := NewService(f)
	req := baseReq()
	req.Goal = "cut"
	req.Target = nutrition.Macros{Calories: 500, ProteinG: 45, CarbsG: 40, FatG: 15}

	resp, err := svc.Suggest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Broadened {
		t.Fatal("expected Broadened=true after the retry")
	}
	if len(resp.Results) == 0 {
		t.Fatal("expected results from the relaxed retry")
	}
}

func TestServiceNoBackend(t *testing.T) {
	svc := NewService(nil)
	if svc.Enabled() {
		t.Fatal("service with nil suggester should not be Enabled")
	}
	if _, err := svc.Suggest(context.Background(), baseReq()); err != ErrNoBackend {
		t.Fatalf("err = %v, want ErrNoBackend", err)
	}
}
