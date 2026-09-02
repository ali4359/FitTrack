package meals

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

// poolCache holds a rotating pool of vetted dishes per profile bucket, so most
// requests never hit the LLM. It is in-memory and process-local — fine for the
// stub; a real deployment would back this with Redis or a table.
type poolCache struct {
	mu    sync.Mutex
	pools map[string]*pool
	ttl   time.Duration
}

type pool struct {
	dishes   []Suggestion
	cursor   int
	filledAt time.Time
}

func newPoolCache() *poolCache {
	return &poolCache{pools: map[string]*pool{}, ttl: 7 * 24 * time.Hour}
}

// bucketKey collapses a request to the axes that actually change the answer.
// Calorie target is bucketed to 250 kcal so near-identical requests share a pool.
func bucketKey(r Request) string {
	allergens := append([]string(nil), r.Allergens...)
	for i := range allergens {
		allergens[i] = strings.ToLower(strings.TrimSpace(allergens[i]))
	}
	sortStrings(allergens)
	kcalBucket := int(math.Round(r.Target.Calories/250) * 250)
	return fmt.Sprintf("%s|%s|%s|%s|%s|%d",
		strings.ToLower(r.MealType),
		strings.ToLower(r.Goal),
		strings.ToLower(string(r.Diet)),
		strings.Join(allergens, ","),
		strings.ToLower(strings.TrimSpace(r.Region)),
		kcalBucket,
	)
}

// take returns n dishes from the pool for this bucket, advancing the rotation
// cursor and skipping anything in req.Exclude. ok is false when the pool is
// missing, stale, or exhausted — the caller then goes to the LLM.
func (c *poolCache) take(req Request, n int) (out []Suggestion, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	p := c.pools[bucketKey(req)]
	if p == nil || time.Since(p.filledAt) > c.ttl || len(p.dishes) < n {
		return nil, false
	}

	excluded := map[string]bool{}
	for _, e := range req.Exclude {
		excluded[normalizeName(e)] = true
	}

	for i := 0; i < len(p.dishes) && len(out) < n; i++ {
		d := p.dishes[(p.cursor+i)%len(p.dishes)]
		if excluded[normalizeName(d.Name)] {
			continue
		}
		out = append(out, d)
	}
	if len(out) < n {
		return nil, false // exclude-list exhausted the pool; regenerate
	}
	p.cursor = (p.cursor + n) % len(p.dishes)
	return out, true
}

// put merges freshly vetted dishes into the bucket's pool (dedup by name).
func (c *poolCache) put(req Request, dishes []Suggestion) {
	if len(dishes) == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	key := bucketKey(req)
	p := c.pools[key]
	if p == nil {
		p = &pool{}
		c.pools[key] = p
	}
	have := map[string]bool{}
	for _, d := range p.dishes {
		have[normalizeName(d.Name)] = true
	}
	for _, d := range dishes {
		if !have[normalizeName(d.Name)] {
			p.dishes = append(p.dishes, d)
			have[normalizeName(d.Name)] = true
		}
	}
	// Keep the pool bounded.
	if len(p.dishes) > 20 {
		p.dishes = p.dishes[len(p.dishes)-20:]
	}
	p.filledAt = time.Now()
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}
