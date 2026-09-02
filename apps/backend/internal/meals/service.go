package meals

import (
	"context"
	"errors"
	"math"
	"sort"

	"github.com/ali4359/fittrack/backend/internal/nutrition"
)

var (
	// ErrNoBackend means no LLM API key is configured.
	ErrNoBackend = errors.New("meals: no suggester backend configured")
	// ErrNoneSurvived means every candidate the model returned failed the guardrails.
	ErrNoneSurvived = errors.New("meals: no suggestions passed the guardrails")
)

// Suggester produces raw (un-guarded) candidate dishes for a request.
type Suggester interface {
	Suggest(ctx context.Context, req Request, n int) ([]Suggestion, error)
}

// Service runs the full pipeline: cache -> Suggester -> guardrails -> top 3.
type Service struct {
	llm   Suggester
	cache *poolCache
}

func NewService(llm Suggester) *Service {
	return &Service{llm: llm, cache: newPoolCache()}
}

// Enabled reports whether a real LLM backend is wired up.
func (s *Service) Enabled() bool { return s != nil && s.llm != nil }

const wantResults = 3

// Suggest returns up to three guarded dishes for the request.
func (s *Service) Suggest(ctx context.Context, req Request) (Response, error) {
	if !s.Enabled() {
		return Response{}, ErrNoBackend
	}

	// 1. Serve from the cached pool for this profile bucket when we can.
	if pooled, ok := s.cache.take(req, wantResults); ok {
		return Response{Target: req.Target, Results: assignRoles(pooled, req), Source: "llm-cache"}, nil
	}

	// 2. Ask for a generous batch so the pool has depth and the guardrails have
	//    room to drop the bad ones.
	raw, err := s.llm.Suggest(ctx, req, 8)
	if err != nil {
		return Response{}, err
	}
	kept := guard(raw, req)

	broadened := false
	if len(kept) < wantResults {
		// 3. One retry with the calorie ceiling relaxed and budget dropped —
		//    never diet or allergens.
		relaxed := req
		relaxed.Budget = ""
		relaxed.Goal = "maintain"
		if more, rerr := s.llm.Suggest(ctx, relaxed, 8); rerr == nil {
			if wider := guard(append(raw, more...), relaxed); len(wider) > len(kept) {
				kept, broadened = wider, true
			}
		}
	}

	if len(kept) == 0 {
		return Response{}, ErrNoneSurvived
	}

	s.cache.put(req, kept)
	top := kept
	if len(top) > wantResults {
		top = top[:wantResults]
	}
	return Response{
		Target:    req.Target,
		Results:   assignRoles(top, req),
		Source:    "llm",
		Broadened: broadened,
	}, nil
}

// score is the weighted distance from the meal target. Lower is better.
// Protein shortfall dominates because that is the point of the meal.
func score(s Suggestion, req Request) float64 {
	tgt := req.Target
	if tgt.Calories <= 0 {
		return s.Calories // "already at target": prefer the lightest
	}

	calErr := (s.Calories - tgt.Calories) / tgt.Calories
	switch req.Goal {
	case "cut":
		if calErr > 0 {
			calErr *= 2 // overshoot hurts a cut
		}
	case "bulk":
		if calErr < 0 {
			calErr *= 2 // undershoot hurts a bulk
		}
	}

	var protShort, carbErr, fatOver float64
	if tgt.ProteinG > 0 {
		protShort = math.Max(0, tgt.ProteinG-s.ProteinG) / tgt.ProteinG
	}
	if tgt.CarbsG > 0 {
		carbErr = math.Abs(s.CarbsG-tgt.CarbsG) / tgt.CarbsG
	}
	if tgt.FatG > 0 {
		fatOver = math.Max(0, s.FatG-tgt.FatG*1.3) / tgt.FatG
	}

	return 3*protShort + 2*math.Abs(calErr) + 1.5*carbErr + fatOver
}

// assignRoles orders the (already guarded) dishes and tags each with the slot
// it fills: best-fit, higher-protein, lighter.
func assignRoles(in []Suggestion, req Request) []Suggestion {
	out := make([]Suggestion, len(in))
	copy(out, in)
	if len(out) == 0 {
		return out
	}

	sort.SliceStable(out, func(i, j int) bool { return score(out[i], req) < score(out[j], req) })
	out[0].Role = "best-fit"

	if len(out) > 1 {
		hp := 1
		for i := 2; i < len(out); i++ {
			if out[i].ProteinG > out[hp].ProteinG {
				hp = i
			}
		}
		out[hp].Role = "higher-protein"
	}
	if len(out) > 2 {
		lt := -1
		for i := 1; i < len(out); i++ {
			if out[i].Role != "" {
				continue
			}
			if lt == -1 || out[i].Calories < out[lt].Calories {
				lt = i
			}
		}
		if lt != -1 {
			out[lt].Role = "lighter"
		}
	}
	return out
}

// FallbackResponse shapes an already-ranked catalogue result for when no LLM
// backend is configured.
func FallbackResponse(target nutrition.Macros, dishes []Suggestion) Response {
	if len(dishes) > wantResults {
		dishes = dishes[:wantResults]
	}
	return Response{Target: target, Results: dishes, Source: "fallback"}
}
