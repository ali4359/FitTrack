# Workout calorie model

How FitTrack estimates the calories burned in a resistance-training session, and
why it is built this way.

- **Code:** [`internal/nutrition/burn.go`](../internal/nutrition/burn.go) — pure, unit-tested in `burn_test.go`
- **Caller:** [`internal/api/workouts.go`](../internal/api/workouts.go) `handleCompleteWorkout`
- **Catalogue data:** [`internal/store/store.go`](../internal/store/store.go) `exercisePhysics`

---

## TL;DR

```
session kcal = mechanical work + baseline metabolism + EPOC
```

Every set the user logs (reps × weight) contributes **mechanical work**. The time
spent training — measured from per-set "set done" timestamps where the app
captured them — contributes **baseline metabolism**. A flat 10 % is added on top
for the post-exercise afterburn (**EPOC**).

The estimate is good to roughly **±25 %**. Without heart-rate or motion data that
is the realistic ceiling for resistance training, so the number is always
presented as an estimate, never a measurement.

---

## Why not just MET × duration?

The original stub used the ACSM MET equation:

```
kcal = avgMET × 3.5 × weightKg / 200 × durationMinutes
```

It has one fatal flaw for this app: **the sets, reps and weights the user
carefully logs do not change the result at all.** Only the average MET of the
exercises performed and the wall-clock duration matter. Logging 3 sets or 10
sets of bench press gives the same number. That makes the core "log your
workout" interaction feel pointless.

We also can't go the other way and use a pure physics model (`work = force ×
distance`), because in a real session the mechanical work of lifting is only
~15–25 % of the energy cost. The rest is the elevated baseline metabolism
between sets and the afterburn. A physics-only model would report ~30 kcal for a
50-minute session.

So the model keeps **both**: the mechanical term (which the logged data drives)
plus a time term (which duration and rest drive).

---

## The three buckets

### 1. Mechanical work

For each set:

```
moved_kg   = weightKg + bodyweightLoadFactor × userWeightKg
work_J     = reps × moved_kg × g × romMeters × (1 + ECCENTRIC_FACTOR)
active_kcal = work_J / LIFTING_EFFICIENCY / JOULES_PER_KCAL
```

| Term | Value | Rationale |
| --- | --- | --- |
| `g` | 9.81 m/s² | gravity |
| `romMeters` | per exercise (see below) | how far the load travels per rep |
| `bodyweightLoadFactor` | per exercise | fraction of body mass also moved through the ROM — without this a bodyweight pull-up computes **zero** work |
| `ECCENTRIC_FACTOR` | 0.30 | lowering the load (negative work) costs ~⅓ of the metabolic energy of lifting it |
| `LIFTING_EFFICIENCY` | 0.22 | net mechanical efficiency of weight training is ~20–25 % |
| `JOULES_PER_KCAL` | 4184 | unit conversion |

This is the only bucket the logged **reps and weight** feed directly. Heavier
load, more reps, or a longer-ROM movement → more calories, monotonically.

### 2. Baseline metabolism

The body burns well above rest for the whole time you're training — breathing
hard, bracing, standing between sets. Charged at a fixed rate over
`time_under_tension + rest` for every set:

```
baseline_rate = REST_MET × 3.5 × userWeightKg / 200 / 60      // kcal per second
baseline_kcal = Σ (tut_seconds + rest_seconds) × baseline_rate
```

| Term | Value | Rationale |
| --- | --- | --- |
| `REST_MET` | 2.5 | standing recovery with elevated breathing |
| `tut_seconds` | `reps × 3` | one "set done" tap can't measure tempo; 3 s/rep ≈ 2 s up + slower lower |
| `rest_seconds` | **measured** (see next section) | the big lever — capped at 150 s/set |

This is the largest bucket (~55–70 % of the total) and the reason session
duration still matters.

### 3. EPOC (afterburn)

Resistance training keeps metabolism elevated after you stop. Modelled as a flat
fraction of buckets 1 + 2:

```
epoc_kcal = 0.10 × (Σ active_kcal + Σ baseline_kcal)
```

Published EPOC for resistance work ranges ~6–15 % of session cost; 10 % is a
middle estimate that avoids per-user physiology we can't observe.

---

## Rest: measured, not guessed

This is what the **"✓ set done"** tap in the mobile session screen buys us.

Each tap records a `completedAt` timestamp for that set. From two consecutive
timestamps:

```
rest_after_set_i = completedAt[i+1] − completedAt[i] − tut[i+1]
```

i.e. the gap between finishing set *i* and finishing set *i+1*, minus the time
the next set itself took. Clamped to `[0, 150] s` so a forgotten running timer
(user drives home, timer still going) can't inflate the burn.

**Fallback cascade** (`resolveTiming` in `burn.go`):

| Situation | Rest used |
| --- | --- |
| Both this set and the next have a timestamp | measured gap (as above) |
| Last set of the session, timestamp present | 0 (no measurable trailing rest) |
| Any set with a missing timestamp | even split of `durationMinutes − Σ tut` across all sets, clamped to 150 s |

If the user never taps ✓, every set falls to the last row and the model
degrades gracefully to roughly `REST_MET × duration` — the same ballpark as the
old MET formula, but still with the mechanical term on top.

### Why per-set timestamps and not per-rep timers

Considered and rejected:

- **Per-rep time entry** — 8–12 inputs per set. Nobody fills it in, and
  rep-to-rep variance is noise for a calorie estimate.
- **Manual per-set TUT field** — better, but still typing during a workout.
- **Start + stop tap per set** — gives true TUT, but doubles the taps.

One "set done" tap gives us the measured **rest intervals** (the largest,
previously-guessed bucket) for the least possible friction, and TUT is cheap to
estimate from rep count. If per-set TUT is ever wanted, `BurnSetInput` can carry
an explicit value and `resolveTiming` can prefer it — the rest of the model
doesn't change.

---

## Per-exercise breakdown

`WorkoutBurn` returns each exercise's share, stored on
`WorkoutExerciseLog.CaloriesBurned`:

```
exercise_kcal = active[i] + baseline[i] + epoc × (active[i] + baseline[i]) / total_work
```

`active[i]` is exact per exercise (it's just that exercise's sets). `baseline[i]`
follows from each set's resolved `tut + rest`, so an exercise with more sets and
longer rests carries a bigger slice. EPOC is distributed in the same proportion.
Rounded per-exercise values may differ from the rounded session total by ~1 kcal.

---

## Exercise catalogue parameters

Every exercise needs `RangeOfMotionM` and `BodyweightLoadFactor`. Seeded in
`store.go` (`exercisePhysics`), backfilled onto pre-existing rows by
`backfillExercisePhysics` on boot.

| Exercise | ROM (m) | Bodyweight load factor |
| --- | --- | --- |
| Barbell Bench Press | 0.40 | 0 |
| Incline Dumbbell Press | 0.40 | 0 |
| Cable Fly | 0.55 | 0 |
| Tricep Pushdown | 0.28 | 0 |
| Pull-up | 0.55 | 0.90 |
| Barbell Row | 0.45 | 0 |
| Back Squat | 0.50 | 0.65 |
| Romanian Deadlift | 0.40 | 0.35 |

Unknown / custom exercises fall back to `defaultROMMeters = 0.45` and a factor of
0. ROM values are averages for a mid-size lifter; they are not personalised.

---

## Constants — one place to tune

All in `burn.go`:

| Constant | Value | Bucket |
| --- | --- | --- |
| `gravity` | 9.81 | mechanical |
| `eccentricFactor` | 0.30 | mechanical |
| `liftingEfficiency` | 0.22 | mechanical |
| `secondsPerRep` | 3.0 | baseline (TUT estimate) |
| `restMET` | 2.5 | baseline |
| `maxRestSecPerSet` | 150 | baseline (anti-gaming clamp) |
| `epocFraction` | 0.10 | EPOC |
| `defaultROMMeters` | 0.45 | mechanical (fallback) |

If these ever need to vary by user or be A/B-tested, promote them to a config
struct passed into `WorkoutBurn`; the function is already pure.

---

## Body weight is required

Every bucket scales with the lifter's mass, and there is no sensible default that
isn't wrong by 20+ kg for a lot of people. `handleCompleteWorkout` returns
**HTTP 422** with `{ "error": ..., "code": "weight_required" }` when
`user.WeightKg <= 0`. Onboarding collects it (the "Body weight" step, between
Goal and Budget & region) and it is editable afterwards on the profile screen;
`needsOnboarding` treats a missing weight as incomplete onboarding. `WorkoutBurn`
also returns a zeroed result for non-positive body weight as a defensive
fallback.

---

## Worked example

50-minute push day, 75 kg lifter, no ✓ taps (pure fallback timing):

| Exercise | Sets | kcal |
| --- | --- | --- |
| Bench 4×8 @ 60 kg | 4 | 54 |
| Incline DB 3×10 @ 40 kg | 3 | 40 |
| Cable Fly 3×12 @ 15 kg | 3 | 38 |
| Tricep Pushdown 3×12 @ 25 kg | 3 | 37 |
| **Total** | | **169** |

In line with published figures (~150–250 kcal) for that session. A denser
session with the same volume but measured 45-second rests lands lower (less time
→ fewer baseline calories) but higher *per minute* — which is the physiologically
correct direction.

---

## Accuracy and honesty

- No heart-rate, no accelerometer → **±25 %** is the floor.
- Isometric bracing, grip, setup/un-racking, and individual metabolic variation
  are not modelled.
- The value's job is to be **directional and consistent**: more work, heavier
  load, denser session → bigger number, every time. The old MET formula failed
  that basic test.
- UI copy always says "estimated", never "burned exactly".
