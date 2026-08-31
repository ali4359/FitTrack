import { useMemo } from 'react';
import { StyleSheet, View } from 'react-native';
import type { BottomTabScreenProps } from '@react-navigation/bottom-tabs';
import type { CompositeScreenProps } from '@react-navigation/native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { AppText, Button, Card, Screen, Tag } from '../../components';
import { useWorkoutDay, useWorkoutHistory } from '../../hooks/useWorkouts';
import { useAuthStore } from '../../store/authStore';
import { colors, goalAccent, radius, spacing } from '../../theme';
import {
  DAY_LABELS,
  computeStreak,
  todaysWorkoutDayId,
  trainedToday,
  weekStrip,
} from '../../lib/workoutPlan';
import type { HomeStackParamList, MainTabParamList } from '../../navigation/types';

type Props = CompositeScreenProps<
  NativeStackScreenProps<HomeStackParamList, 'Today'>,
  BottomTabScreenProps<MainTabParamList>
>;

function greeting() {
  const h = new Date().getHours();
  if (h < 12) return 'Good morning';
  if (h < 18) return 'Good afternoon';
  return 'Good evening';
}

export function HomeScreen({ navigation }: Props) {
  const user = useAuthStore((s) => s.user);
  const history = useWorkoutHistory();

  const dayId = useMemo(() => todaysWorkoutDayId(history.data), [history.data]);
  const day = useWorkoutDay(dayId);
  const streak = useMemo(() => computeStreak(history.data), [history.data]);
  const done = trainedToday(history.data);
  const week = useMemo(() => weekStrip(history.data), [history.data]);

  const accent = user ? goalAccent[user.goal] ?? colors.turmeric : colors.turmeric;
  const exerciseCount = day.data?.exercises.length ?? 0;

  return (
    <Screen
      scroll
      contentStyle={{ paddingTop: spacing.xxl }}
      refreshing={history.isRefetching || day.isRefetching}
      onRefresh={() => {
        void history.refetch();
        void day.refetch();
      }}
    >
      <View style={styles.headerRow}>
        <View style={{ flex: 1 }}>
          <AppText variant="subtitle">{greeting()},</AppText>
          <AppText variant="display">{user?.name?.split(' ')[0] ?? 'Athlete'}</AppText>
        </View>
        <View style={styles.streak}>
          <AppText variant="monoLarge" color={colors.turmeric}>
            {streak}
          </AppText>
          <AppText variant="label" style={{ textTransform: 'uppercase' }}>
            day streak
          </AppText>
        </View>
      </View>

      {/* Week at a glance */}
      <View style={styles.week}>
        {week.map((d, i) => (
          <View
            key={i}
            style={[
              styles.weekDot,
              d.done && { backgroundColor: accent, borderColor: accent },
              d.isToday && !d.done && { borderColor: colors.turmeric },
              d.future && { opacity: 0.4 },
            ]}
          >
            <AppText
              variant="label"
              color={d.done ? colors.onAccent : d.isToday ? colors.turmeric : colors.textSecondary}
            >
              {d.label}
            </AppText>
          </View>
        ))}
      </View>

      {/* Today's workout */}
      <Card padded style={{ gap: spacing.md }}>
        <View style={styles.cardTop}>
          <Tag label={done ? 'COMPLETED TODAY' : "TODAY'S SESSION"} tone={done ? 'cardamom' : 'turmeric'} />
          {user ? <Tag label={user.goal.toUpperCase()} tone="neutral" /> : null}
        </View>

        <AppText variant="title">
          {day.data?.name ?? DAY_LABELS[dayId] ?? 'Loading…'}
        </AppText>

        <View style={styles.metaRow}>
          <AppText variant="mono" color={colors.textDim}>
            {exerciseCount > 0 ? `${exerciseCount} exercises` : '—'}
          </AppText>
          {day.data ? (
            <AppText variant="mono" color={colors.textDim}>
              ~{Math.max(30, exerciseCount * 12)} min
            </AppText>
          ) : null}
        </View>

        <Button
          label={done ? 'Train again' : 'Start workout'}
          variant={done ? 'secondary' : 'primary'}
          fullWidth
          disabled={!day.data}
          onPress={() => navigation.navigate('WorkoutSession', { workoutDayId: dayId })}
        />
      </Card>

      {done ? (
        <Card padded style={{ gap: spacing.sm }}>
          <AppText variant="title" style={{ fontSize: 18 }}>
            Nice work today
          </AppText>
          <AppText variant="caption">
            Check a post-workout meal picked for your goal and budget.
          </AppText>
          <Button
            label="See meal"
            variant="ghost"
            onPress={() => navigation.navigate('MealSuggestion', { mealType: 'post-workout' })}
          />
        </Card>
      ) : null}

      {history.isError ? (
        <AppText variant="caption" color={colors.chili}>
          Couldn&apos;t load your history. Pull the button above to retry.
        </AppText>
      ) : null}
    </Screen>
  );
}

const styles = StyleSheet.create({
  headerRow: { flexDirection: 'row', alignItems: 'flex-start', gap: spacing.md },
  streak: { alignItems: 'flex-end' },
  week: { flexDirection: 'row', justifyContent: 'space-between', gap: spacing.xs },
  weekDot: {
    flex: 1,
    aspectRatio: 1,
    borderRadius: radius.md,
    borderWidth: 1.5,
    borderColor: colors.border,
    alignItems: 'center',
    justifyContent: 'center',
  },
  cardTop: { flexDirection: 'row', justifyContent: 'space-between' },
  metaRow: { flexDirection: 'row', gap: spacing.lg },
});
