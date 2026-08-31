import { useMutation, useQuery } from '@tanstack/react-query';
import type { LogMealRequest, MealSuggestionType } from '@fittrack/shared';
import { logMeal, suggestMeals } from '../api/endpoints';
import { queryKeys } from '../api/queryClient';

export function useMealSuggestions(type: MealSuggestionType) {
  return useQuery({
    queryKey: queryKeys.mealSuggestions(type),
    queryFn: () => suggestMeals(type),
  });
}

export function useLogMeal() {
  return useMutation({
    mutationFn: (body: LogMealRequest) => logMeal(body),
  });
}
