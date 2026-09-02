# FitTrack — Backend (stub)

> **Note:** This is a minimal Gin + GORM stub scaffolded to give the mobile app a
> working API to build against. Replace it with the real Go backend when ready —
> keep the routes and JSON shapes identical (they mirror `packages/shared`).

## Run

```bash
# zero-setup: uses a local SQLite file, seeds demo data
go run ./cmd/api

# or against Postgres
docker compose up -d
DATABASE_URL='postgres://postgres:postgres@localhost:5432/fittrack?sslmode=disable' go run ./cmd/api
```

Server listens on `:8080` (override with `PORT`).

### Meal suggestions (LLM)

`GET /api/meals/suggest` computes a per-meal calorie/macro target from the
profile, today's workout burn, and what's been logged today, then asks Claude for
Pakistani dishes that fill the gap. Every candidate is run through code
guardrails (Atwater macro sanity, dietary rule + allergen scan, calorie ceiling,
diversity) before it is returned, and results are cached per profile bucket.

| Env var | Default | Purpose |
| --- | --- | --- |
| `ANTHROPIC_API_KEY` | – | when unset, the endpoint falls back to the seeded `MealEntry` catalogue |
| `MEALS_LLM_MODEL` | `claude-haiku-4-5` | override the model |

## Demo credentials

`demo@fittrack.app` / `password123` (goal: bulk, region: Lahore, budget: mid, halal)

## Endpoints

| Method | Path | Auth | Notes |
| --- | --- | --- | --- |
| POST | `/api/auth/register` | – | `{ email, password, name }` → `{ token, user }` |
| POST | `/api/auth/login` | – | `{ email, password }` → `{ token, user }` |
| GET | `/api/profile` | ✓ | `{ user }` |
| PATCH | `/api/profile` | ✓ | onboarding / profile edits write here |
| GET | `/api/workouts/:dayId` | ✓ | `WorkoutDay` with nested exercises |
| POST | `/api/workouts/complete` | ✓ | → `{ workoutLog, nextStep }` |
| GET | `/api/workouts/history` | ✓ | `{ results: WorkoutLog[] }` (last 30) |
| GET | `/api/meals/suggest?type=post-workout&exclude=Nihari,Daal` | ✓ | `{ target, dailyTarget, results: MealSuggestion[], source, broadened }` |
| POST | `/api/meals/log` | ✓ | `{ mealType, servings?, mealEntryId? }` or an inline `{ dishName, calories, proteinG, carbsG, fatG }` |
| GET | `/api/progress/summary` | ✓ | `{ workoutsThisMonth, avgCaloriesBurned }` |

Calorie burn is a MET-based estimate: `kcal = avgMET * 3.5 * weightKg / 200 * minutes`.

Seeded workout days: `day-1` (Chest & Triceps), `day-2` (Back & Biceps), `day-3` (Legs).
