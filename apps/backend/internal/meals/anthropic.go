package meals

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// anthropicSuggester asks Claude for candidate dishes via a single forced tool
// call, which gives us a strict JSON shape without prose to parse around.
type anthropicSuggester struct {
	client anthropic.Client
	model  string
}

// NewAnthropicSuggester returns a Suggester, or nil if ANTHROPIC_API_KEY is
// unset (the handler then falls back to the seeded catalogue).
func NewAnthropicSuggester() Suggester {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		return nil
	}
	model := os.Getenv("MEALS_LLM_MODEL")
	if model == "" {
		// Cheap + fast: this is a constrained generate-and-format task, not
		// open-ended reasoning. Override with MEALS_LLM_MODEL if quality lags.
		model = string(anthropic.ModelClaudeHaiku4_5)
	}
	return &anthropicSuggester{
		client: anthropic.NewClient(option.WithAPIKey(key)),
		model:  model,
	}
}

const suggestToolName = "submit_dishes"

func suggestTool() anthropic.ToolUnionParam {
	dish := map[string]any{
		"type":     "object",
		"required": []string{"name", "portion", "ingredients", "calories", "proteinG", "carbsG", "fatG", "isVegetarian", "isVegan", "isHalal", "whyItFits"},
		"properties": map[string]any{
			"name":         map[string]any{"type": "string", "description": "Common dish name, English (add Urdu in parens if helpful)"},
			"portion":      map[string]any{"type": "string", "description": "Realistic home portion, e.g. '1.5 plates (~420g)' or '2 rotis + 1 cup curry'"},
			"ingredients":  map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Main ingredients, lowercase, one per item"},
			"calories":     map[string]any{"type": "number", "description": "kcal for the stated portion"},
			"proteinG":     map[string]any{"type": "number"},
			"carbsG":       map[string]any{"type": "number"},
			"fatG":         map[string]any{"type": "number"},
			"isVegetarian": map[string]any{"type": "boolean"},
			"isVegan":      map[string]any{"type": "boolean"},
			"isHalal":      map[string]any{"type": "boolean"},
			"whyItFits":    map[string]any{"type": "string", "description": "One short sentence tying the dish to the calorie/macro target"},
		},
	}
	return anthropic.ToolUnionParam{OfTool: &anthropic.ToolParam{
		Name:        suggestToolName,
		Description: anthropic.String("Return the list of suggested dishes."),
		InputSchema: anthropic.ToolInputSchemaParam{
			Properties: map[string]any{
				"dishes": map[string]any{"type": "array", "items": dish},
			},
			Required: []string{"dishes"},
		},
	}}
}

func (a *anthropicSuggester) Suggest(ctx context.Context, req Request, n int) ([]Suggestion, error) {
	resp, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(a.model),
		MaxTokens: 2048,
		System: []anthropic.TextBlockParam{{
			Text: systemPrompt,
		}},
		Tools:      []anthropic.ToolUnionParam{suggestTool()},
		ToolChoice: anthropic.ToolChoiceParamOfTool(suggestToolName),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt(req, n))),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("meals: anthropic call: %w", err)
	}

	for _, block := range resp.Content {
		if tu, ok := block.AsAny().(anthropic.ToolUseBlock); ok && tu.Name == suggestToolName {
			var payload struct {
				Dishes []Suggestion `json:"dishes"`
			}
			if err := json.Unmarshal([]byte(tu.JSON.Input.Raw()), &payload); err != nil {
				return nil, fmt.Errorf("meals: decode tool input: %w", err)
			}
			return payload.Dishes, nil
		}
	}
	return nil, fmt.Errorf("meals: model did not call %s", suggestToolName)
}

const systemPrompt = `You are a nutrition assistant for a fitness app used in Pakistan.
You suggest common, home-cooked Pakistani dishes for a single meal.

Rules:
- Only real, widely-eaten Pakistani/South Asian dishes. No invented fusion.
- Macros must be realistic for the portion you state and must satisfy
  calories ≈ 4·protein + 4·carbs + 9·fat (within ~10%).
- Respect the dietary rule and allergens ABSOLUTELY. If asked for vegetarian,
  never include meat/fish. If halal, never include pork or alcohol.
- Prefer dishes that land near the calorie and protein target for the meal.
- Vary the primary protein/starch across your suggestions.
- Portions should be realistic (0.5x–2x a normal serving).`

func userPrompt(req Request, n int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Suggest %d dishes for a %s meal.\n\n", n, orDefault(req.MealType, "daily"))
	fmt.Fprintf(&b, "Target for THIS meal: ~%.0f kcal, ~%.0fg protein, ~%.0fg carbs, ~%.0fg fat.\n",
		req.Target.Calories, req.Target.ProteinG, req.Target.CarbsG, req.Target.FatG)
	fmt.Fprintf(&b, "Training goal: %s.\n", orDefault(req.Goal, "maintain"))
	fmt.Fprintf(&b, "Dietary rule: %s.\n", orDefault(string(req.Diet), string(DietOmnivore)))
	if len(req.Allergens) > 0 {
		fmt.Fprintf(&b, "Avoid entirely (allergens): %s.\n", strings.Join(req.Allergens, ", "))
	}
	if req.Region != "" {
		fmt.Fprintf(&b, "Region: %s — favour dishes eaten there.\n", req.Region)
	}
	if req.Budget != "" {
		fmt.Fprintf(&b, "Budget: %s.\n", req.Budget)
	}
	if len(req.Exclude) > 0 {
		fmt.Fprintf(&b, "Do NOT suggest (shown recently): %s.\n", strings.Join(req.Exclude, ", "))
	}
	if req.MealType == "post-workout" {
		b.WriteString("This is a post-workout meal: prioritise protein and carbs for recovery; keep fat moderate.\n")
	}
	return b.String()
}

func orDefault(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}
