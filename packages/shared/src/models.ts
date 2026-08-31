/**
 * Core domain models — mirror the Go structs in apps/backend.
 * Keep field names and unions in sync with the GORM models so the
 * mobile app's API types don't drift from what the backend returns.
 */

export type Goal = 'bulk' | 'cut' | 'maintain';
export type BudgetTier = 'low' | 'mid' | 'high';

export type User = {
  id: string;
  email: string;
  name: string;
  goal: Goal;
  weightKg: number;
  heightCm: number;
  region: string;
  budgetTier: BudgetTier;
  /** comma-separated, e.g. "halal,vegetarian" */
  restrictions: string;
};

export type Exercise = {
  id: string;
  name: string;
  muscleGroup: string;
  metValue: number;
};

export type WorkoutDayExercise = {
  exercise: Exercise;
  defaultSets: number;
  defaultReps: number;
};

export type WorkoutDay = {
  id: string;
  name: string;
  dayOrder: number;
  exercises: WorkoutDayExercise[];
};

export type WorkoutLog = {
  id: string;
  workoutDayId: string;
  /** ISO 8601 timestamp */
  completedAt: string;
  durationMinutes: number;
  caloriesBurned: number;
};

export type MealEntry = {
  id: string;
  dishName: string;
  region: string;
  budgetTier: string;
  calories: number;
  proteinG: number;
  carbsG: number;
  fatG: number;
  /** comma-separated tags, e.g. "bulk,post-workout" */
  goalTags: string;
  halal: boolean;
  vegetarian: boolean;
};
