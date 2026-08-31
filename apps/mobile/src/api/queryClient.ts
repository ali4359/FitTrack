import { QueryClient } from '@tanstack/react-query';

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 30_000,
      refetchOnWindowFocus: false,
    },
  },
});

export const queryKeys = {
  profile: ['profile'] as const,
  workoutDay: (dayId: string) => ['workout-day', dayId] as const,
  workoutHistory: ['workout-history'] as const,
  mealSuggestions: (type: string) => ['meal-suggestions', type] as const,
  progressSummary: ['progress-summary'] as const,
};
