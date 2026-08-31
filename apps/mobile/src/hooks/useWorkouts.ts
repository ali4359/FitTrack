import { useMutation, useQuery } from '@tanstack/react-query';
import type { CompleteWorkoutRequest } from '@fittrack/shared';
import { completeWorkout, getWorkoutDay, getWorkoutHistory } from '../api/endpoints';
import { queryClient, queryKeys } from '../api/queryClient';

export function useWorkoutDay(dayId: string) {
  return useQuery({
    queryKey: queryKeys.workoutDay(dayId),
    queryFn: () => getWorkoutDay(dayId),
    enabled: !!dayId,
  });
}

export function useWorkoutHistory() {
  return useQuery({
    queryKey: queryKeys.workoutHistory,
    queryFn: getWorkoutHistory,
  });
}

export function useCompleteWorkout() {
  return useMutation({
    mutationFn: (body: CompleteWorkoutRequest) => completeWorkout(body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.workoutHistory });
      queryClient.invalidateQueries({ queryKey: queryKeys.progressSummary });
    },
  });
}
