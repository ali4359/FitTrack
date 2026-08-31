import type { NavigatorScreenParams } from '@react-navigation/native';
import type { MealSuggestionType } from '@iron-and-spice/shared';

export type AuthStackParamList = {
  Login: undefined;
  Register: undefined;
};

export type OnboardingStackParamList = {
  Goal: undefined;
  BudgetRegion: undefined;
  Restrictions: undefined;
};

export type HomeStackParamList = {
  Today: undefined;
  WorkoutSession: { workoutDayId: string };
  SessionComplete: {
    workoutLogId: string;
    caloriesBurned: number;
    durationMinutes: number;
    workoutDayId: string;
  };
  MealSuggestion: { mealType: MealSuggestionType };
};

export type MainTabParamList = {
  HomeTab: NavigatorScreenParams<HomeStackParamList>;
  ProgressTab: undefined;
  ProfileTab: undefined;
};
