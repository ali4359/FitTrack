import { useQuery } from '@tanstack/react-query';
import { getProgressSummary } from '../api/endpoints';
import { queryKeys } from '../api/queryClient';

export function useProgressSummary() {
  return useQuery({
    queryKey: queryKeys.progressSummary,
    queryFn: getProgressSummary,
  });
}
