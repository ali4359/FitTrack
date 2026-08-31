/**
 * Request/response shapes for the Go backend API.
 * Endpoint reference lives in the root README and apps/backend.
 */
import type {
  BudgetTier,
  Goal,
  MealEntry,
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
  weightKg: number;
  heightCm: number;
  region: string;
  budgetTier: BudgetTier;
  restrictions: string;
}>;

// ---- Workouts ----

export type WorkoutExerciseResult = {
  exerciseId: string;
  setsDone: number;
  repsDone: number;
  weightKg: number;
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
  results: MealEntry[];
};

export type LogMealRequest = {
  mealEntryId: string;
  mealType: MealSuggestionType;
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
