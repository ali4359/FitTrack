import { useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
import { Alert, Pressable, StyleSheet, TextInput, View } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import type { WorkoutExerciseResult } from '@fittrack/shared';
import { AppText, Button, Card, Screen, Tag } from '../../components';
import { apiErrorMessage } from '../../api/client';
import { useCompleteWorkout, useWorkoutDay } from '../../hooks/useWorkouts';
import { colors, fonts, fontSize, radius, spacing } from '../../theme';
import type { HomeStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<HomeStackParamList, 'WorkoutSession'>;

type ExerciseState = {
  setsDone: number;
  repsDone: string;
  weightKg: string;
};

function fmtClock(totalSeconds: number) {
  const m = Math.floor(totalSeconds / 60);
  const s = totalSeconds % 60;
  return `${m}:${String(s).padStart(2, '0')}`;
}

export function WorkoutSessionScreen({ navigation, route }: Props) {
  const { workoutDayId } = route.params;
  const day = useWorkoutDay(workoutDayId);
  const complete = useCompleteWorkout();

  const startedAt = useRef(Date.now());
  const [elapsed, setElapsed] = useState(0);
  useEffect(() => {
    const id = setInterval(() => setElapsed(Math.floor((Date.now() - startedAt.current) / 1000)), 1000);
    return () => clearInterval(id);
  }, []);

  const [state, setState] = useState<Record<string, ExerciseState>>({});

  // seed per-exercise state once the plan loads
  useEffect(() => {
    if (!day.data) return;
    setState((prev) => {
      if (Object.keys(prev).length) return prev;
      const seeded: Record<string, ExerciseState> = {};
      for (const row of day.data.exercises) {
        seeded[row.exercise.id] = {
          setsDone: 0,
          repsDone: String(row.defaultReps),
          weightKg: '',
        };
      }
      return seeded;
    });
  }, [day.data]);

  useLayoutEffect(() => {
    navigation.setOptions({
      title: day.data?.name ?? 'Workout',
      headerRight: () => (
        <AppText variant="mono" color={colors.turmeric}>
          {fmtClock(elapsed)}
        </AppText>
      ),
    });
  }, [navigation, day.data?.name, elapsed]);

  const touched = useMemo(
    () => Object.values(state).filter((s) => s.setsDone > 0).length,
    [state],
  );

  const patch = (id: string, p: Partial<ExerciseState>) =>
    setState((prev) => ({ ...prev, [id]: { ...prev[id], ...p } }));

  const finish = () => {
    if (!day.data) return;
    const exercises: WorkoutExerciseResult[] = day.data.exercises.map((row) => {
      const s = state[row.exercise.id];
      return {
        exerciseId: row.exercise.id,
        setsDone: s?.setsDone ?? 0,
        repsDone: Number(s?.repsDone) || 0,
        weightKg: Number(s?.weightKg) || 0,
      };
    });
    const durationMinutes = Math.max(1, Math.round(elapsed / 60));

    complete.mutate(
      { workoutDayId, durationMinutes, exercises },
      {
        onSuccess: (res) => {
          navigation.replace('SessionComplete', {
            workoutLogId: res.workoutLog.id,
            caloriesBurned: res.workoutLog.caloriesBurned,
            durationMinutes: res.workoutLog.durationMinutes,
            workoutDayId,
          });
        },
        onError: (err) => Alert.alert('Could not save', apiErrorMessage(err)),
      },
    );
  };

  const confirmFinish = () => {
    const msg =
      touched === 0
        ? "You haven't logged any sets. Finish anyway?"
        : `Logging ${touched} of ${day.data?.exercises.length} exercises.`;
    Alert.alert('Finish workout', msg, [
      { text: 'Keep going', style: 'cancel' },
      { text: 'Finish', style: 'default', onPress: finish },
    ]);
  };

  const confirmQuit = () => {
    Alert.alert('Discard workout?', 'Your progress in this session will be lost.', [
      { text: 'Cancel', style: 'cancel' },
      { text: 'Discard', style: 'destructive', onPress: () => navigation.goBack() },
    ]);
  };

  if (day.isLoading) {
    return (
      <Screen>
        <AppText variant="subtitle">Loading your session…</AppText>
      </Screen>
    );
  }

  if (day.isError || !day.data) {
    return (
      <Screen footer={<Button label="Back" variant="secondary" fullWidth onPress={() => navigation.goBack()} />}>
        <AppText variant="title">Couldn&apos;t load this workout</AppText>
        <AppText variant="caption">{apiErrorMessage(day.error)}</AppText>
        <Button label="Retry" onPress={() => day.refetch()} />
      </Screen>
    );
  }

  return (
    <Screen
      scroll
      footer={
        <View style={{ gap: spacing.sm }}>
          <Button
            label={complete.isPending ? 'Saving…' : 'Finish workout'}
            fullWidth
            loading={complete.isPending}
            onPress={confirmFinish}
          />
          <Button label="Discard" variant="danger" fullWidth onPress={confirmQuit} />
        </View>
      }
    >
      {day.data.exercises.map((row) => {
        const s = state[row.exercise.id] ?? { setsDone: 0, repsDone: '', weightKg: '' };
        const complete_ = s.setsDone >= row.defaultSets;
        return (
          <Card key={row.exercise.id} style={{ gap: spacing.md }}>
            <View style={styles.exHeader}>
              <View style={{ flex: 1, gap: 4 }}>
                <AppText variant="title" style={{ fontSize: 18 }}>
                  {row.exercise.name}
                </AppText>
                <Tag label={row.exercise.muscleGroup.toUpperCase()} tone={complete_ ? 'cardamom' : 'neutral'} />
              </View>
              <View style={{ alignItems: 'flex-end' }}>
                <AppText variant="label" style={{ textTransform: 'uppercase' }}>
                  target
                </AppText>
                <AppText variant="mono" color={colors.textDim}>
                  {row.defaultSets} × {row.defaultReps}
                </AppText>
              </View>
            </View>

            <View>
              <AppText variant="label" style={{ textTransform: 'uppercase', marginBottom: spacing.sm }}>
                sets done
              </AppText>
              <View style={styles.setRow}>
                {Array.from({ length: row.defaultSets }).map((_, i) => {
                  const filled = i < s.setsDone;
                  return (
                    <Pressable
                      key={i}
                      onPress={() => patch(row.exercise.id, { setsDone: filled ? i : i + 1 })}
                      style={[styles.setPill, filled && { backgroundColor: colors.turmeric, borderColor: colors.turmeric }]}
                    >
                      <AppText
                        style={{
                          fontFamily: fonts.monoBold,
                          fontSize: fontSize.md,
                          color: filled ? colors.onAccent : colors.textSecondary,
                        }}
                      >
                        {i + 1}
                      </AppText>
                    </Pressable>
                  );
                })}
              </View>
            </View>

            <View style={styles.inputsRow}>
              <NumField
                label="reps / set"
                value={s.repsDone}
                onChangeText={(t) => patch(row.exercise.id, { repsDone: t })}
              />
              <NumField
                label="weight (kg)"
                value={s.weightKg}
                placeholder="0"
                onChangeText={(t) => patch(row.exercise.id, { weightKg: t })}
              />
            </View>
          </Card>
        );
      })}

      <AppText variant="caption" center>
        Calorie burn is estimated from these numbers — log what you actually did.
      </AppText>
    </Screen>
  );
}

function NumField({
  label,
  value,
  onChangeText,
  placeholder,
}: {
  label: string;
  value: string;
  onChangeText: (t: string) => void;
  placeholder?: string;
}) {
  return (
    <View style={{ flex: 1, gap: spacing.sm }}>
      <AppText variant="label" style={{ textTransform: 'uppercase' }}>
        {label}
      </AppText>
      <TextInput
        value={value}
        onChangeText={(t) => onChangeText(t.replace(/[^0-9.]/g, ''))}
        keyboardType="numeric"
        placeholder={placeholder}
        placeholderTextColor={colors.textSecondary}
        style={styles.numInput}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  exHeader: { flexDirection: 'row', gap: spacing.md },
  setRow: { flexDirection: 'row', flexWrap: 'wrap', gap: spacing.sm },
  setPill: {
    width: 44,
    height: 44,
    borderRadius: radius.md,
    borderWidth: 1.5,
    borderColor: colors.border,
    alignItems: 'center',
    justifyContent: 'center',
  },
  inputsRow: { flexDirection: 'row', gap: spacing.md },
  numInput: {
    backgroundColor: colors.ink,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border,
    borderRadius: radius.md,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    color: colors.textLight,
    fontFamily: fonts.mono,
    fontSize: fontSize.lg,
  },
});
