/**
 * Request/response shapes for the Go backend API.
 * Endpoint reference lives in the root README and apps/backend.
 */
import type {
  BudgetTier,
  Goal,
  Macros,
  MealSuggestion,
  Sex,
  User,
  WorkoutDay,
  WorkoutLog,
} from './models';

// ---- Auth ----

export type RegisterRequest = {
  email: string;
  password: string;
  name: string;
};

export type LoginRequest = {
  email: string;
  password: string;
};

export type AuthResponse = {
  token: string;
  user: User;
};

// ---- Profile (onboarding writes here) ----

export type UpdateProfileRequest = Partial<{
  name: string;
  goal: Goal;
  sex: Sex;
  age: number;
  weightKg: number;
  heightCm: number;
  region: string;
  budgetTier: BudgetTier;
  restrictions: string;
}>;

// ---- Workouts ----

/** One set the user logged for an exercise. */
export type WorkoutSetInput = {
  reps: number;
  weightKg: number;
};

export type WorkoutExerciseResult = {
  exerciseId: string;
  /** in performed order; empty means the exercise was skipped */
  sets: WorkoutSetInput[];
};

export type CompleteWorkoutRequest = {
  workoutDayId: string;
  durationMinutes: number;
  exercises: WorkoutExerciseResult[];
};

export type NextStep = {
  kind: 'meal-suggestion';
  mealType: MealSuggestionType;
  message: string;
};

export type CompleteWorkoutResponse = {
  workoutLog: WorkoutLog;
  nextStep: NextStep;
};

export type WorkoutHistoryResponse = {
  results: WorkoutLog[];
};

export type WorkoutDayResponse = WorkoutDay;

// ---- Meals ----

export type MealSuggestionType =
  | 'post-workout'
  | 'breakfast'
  | 'lunch'
  | 'dinner'
  | 'daily';

export type MealSuggestResponse = {
  /** macro target for this one meal */
  target: Macros;
  /** the full-day target, for context on the card */
  dailyTarget: Macros;
  results: MealSuggestion[];
  /** "llm" | "llm-cache" | "fallback" */
  source: string;
  /** true when constraints had to be relaxed to return three dishes */
  broadened: boolean;
};

export type LogMealRequest = {
  mealType: MealSuggestionType;
  /** portion multiplier the user picked; defaults to 1 */
  servings?: number;
  /** reference a seeded catalogue entry... */
  mealEntryId?: string;
  /** ...or pass an LLM suggestion inline (per-serving macros) */
  dishName?: string;
  calories?: number;
  proteinG?: number;
  carbsG?: number;
  fatG?: number;
};

export type LogMealResponse = {
  logged: true;
};

// ---- Progress ----

export type ProgressSummaryResponse = {
  workoutsThisMonth: number;
  avgCaloriesBurned: number;
};

// ---- Errors ----

export type ApiError = {
  error: string;
};
