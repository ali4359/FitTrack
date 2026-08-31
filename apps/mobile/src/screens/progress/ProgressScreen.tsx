import { View } from 'react-native';
import { AppText, Card, Screen } from '../../components';
import { useProgressSummary } from '../../hooks/useProgress';
import { useWorkoutHistory } from '../../hooks/useWorkouts';
import { colors, spacing } from '../../theme';

/** Placeholder — full build (chart + sessions list) is screen 7 in the plan. */
export function ProgressScreen() {
  const summary = useProgressSummary();
  const history = useWorkoutHistory();

  return (
    <Screen
      scroll
      contentStyle={{ paddingTop: spacing.xxl }}
      refreshing={summary.isRefetching}
      onRefresh={() => {
        void summary.refetch();
        void history.refetch();
      }}
    >
      <AppText variant="display">Progress</AppText>

      <View style={{ flexDirection: 'row', gap: spacing.md }}>
        <Card padded style={{ flex: 1, gap: spacing.xs }}>
          <AppText variant="label" style={{ textTransform: 'uppercase' }}>
            workouts / month
          </AppText>
          <AppText variant="monoLarge" color={colors.turmeric}>
            {summary.data?.workoutsThisMonth ?? '—'}
          </AppText>
        </Card>
        <Card padded style={{ flex: 1, gap: spacing.xs }}>
          <AppText variant="label" style={{ textTransform: 'uppercase' }}>
            avg burn
          </AppText>
          <AppText variant="monoLarge" color={colors.turmeric}>
            {summary.data ? Math.round(summary.data.avgCaloriesBurned) : '—'}
          </AppText>
        </Card>
      </View>

      <AppText variant="caption">
        Recent sessions: {history.data?.length ?? 0}. Chart and session list come with the full
        Progress screen.
      </AppText>
    </Screen>
  );
}
