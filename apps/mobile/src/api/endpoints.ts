import type {
  AuthResponse,
  CompleteWorkoutRequest,
  CompleteWorkoutResponse,
  LoginRequest,
  LogMealRequest,
  LogMealResponse,
  MealSuggestResponse,
  MealSuggestionType,
  ProgressSummaryResponse,
  RegisterRequest,
  UpdateProfileRequest,
  User,
  WorkoutDay,
  WorkoutHistoryResponse,
} from '@fittrack/shared';
import { api } from './client';

// ---- Auth ----
export async function register(body: RegisterRequest) {
  const { data } = await api.post<AuthResponse>('/auth/register', body);
  return data;
}

export async function login(body: LoginRequest) {
  const { data } = await api.post<AuthResponse>('/auth/login', body);
  return data;
}

// ---- Profile ----
export async function getProfile() {
  const { data } = await api.get<{ user: User }>('/profile');
  return data.user;
}

export async function updateProfile(body: UpdateProfileRequest) {
  const { data } = await api.patch<{ user: User }>('/profile', body);
  return data.user;
}

// ---- Workouts ----
export async function getWorkoutDay(dayId: string) {
  const { data } = await api.get<WorkoutDay>(`/workouts/${dayId}`);
  return data;
}

export async function completeWorkout(body: CompleteWorkoutRequest) {
  const { data } = await api.post<CompleteWorkoutResponse>('/workouts/complete', body);
  return data;
}

export async function getWorkoutHistory() {
  const { data } = await api.get<WorkoutHistoryResponse>('/workouts/history');
  return data.results;
}

// ---- Meals ----
export async function suggestMeals(type: MealSuggestionType) {
  const { data } = await api.get<MealSuggestResponse>('/meals/suggest', { params: { type } });
  return data.results;
}

export async function logMeal(body: LogMealRequest) {
  const { data } = await api.post<LogMealResponse>('/meals/log', body);
  return data;
}

// ---- Progress ----
export async function getProgressSummary() {
  const { data } = await api.get<ProgressSummaryResponse>('/progress/summary');
  return data;
}
