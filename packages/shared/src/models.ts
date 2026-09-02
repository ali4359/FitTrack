/**
 * Core domain models — mirror the Go structs in apps/backend.
 * Keep field names and unions in sync with the GORM models so the
 * mobile app's API types don't drift from what the backend returns.
 */

export type Goal = 'bulk' | 'cut' | 'maintain';
export type BudgetTier = 'low' | 'mid' | 'high';

export type Sex = 'male' | 'female' | '';

export type User = {
  id: string;
  email: string;
  name: string;
  goal: Goal;
  /** used for the BMR calculation; '' when not yet asked */
  sex: Sex;
  /** years; 0 when not yet asked */
  age: number;
  weightKg: number;
  heightCm: number;
  region: string;
  budgetTier: BudgetTier;
  /** comma-separated, e.g. "halal,vegetarian,no beef" */
  restrictions: string;
};

/** A calorie + macronutrient vector (grams). */
export type Macros = {
  calories: number;
  proteinG: number;
  carbsG: number;
  fatG: number;
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

/** MealEntry is the seeded catalogue row, used only by the no-LLM fallback. */
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

/** MealSuggestion is one dish returned by GET /meals/suggest. */
export type MealSuggestion = {
  name: string;
  /** realistic home portion, e.g. "1.5 plates (~420g)" */
  portion: string;
  ingredients: string[];
  calories: number;
  proteinG: number;
  carbsG: number;
  fatG: number;
  isVegetarian: boolean;
  isVegan: boolean;
  isHalal: boolean;
  /** which slot this dish fills in the set of three */
  role: 'best-fit' | 'higher-protein' | 'lighter' | '';
  whyItFits: string;
  /** macros are model estimates, not measured — show them as approximate */
  estimated: boolean;
};
