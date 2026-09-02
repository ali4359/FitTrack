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

type SetEntry = { reps: string; weightKg: string; doneAt: number | null };
type ExerciseState = { sets: SetEntry[] };

function fmtClock(totalSeconds: number) {
  const m = Math.floor(totalSeconds / 60);
  const s = totalSeconds % 60;
  return `${m}:${String(s).padStart(2, '0')}`;
}

const isLogged = (set: SetEntry) => (Number(set.reps) || 0) > 0;

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

  // seed one row per planned set once the plan loads
  useEffect(() => {
    if (!day.data) return;
    setState((prev) => {
      if (Object.keys(prev).length) return prev;
      const seeded: Record<string, ExerciseState> = {};
      for (const row of day.data.exercises) {
        seeded[row.exercise.id] = {
          sets: Array.from({ length: Math.max(1, row.defaultSets) }, () => ({
            reps: String(row.defaultReps),
            weightKg: '',
            doneAt: null,
          })),
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
    () => Object.values(state).filter((s) => s.sets.some(isLogged)).length,
    [state],
  );

  const setSets = (id: string, next: (sets: SetEntry[]) => SetEntry[]) =>
    setState((prev) => ({ ...prev, [id]: { sets: next(prev[id]?.sets ?? []) } }));

  const patchSet = (id: string, index: number, p: Partial<SetEntry>) =>
    setSets(id, (sets) => sets.map((s, i) => (i === index ? { ...s, ...p } : s)));

  const addSet = (id: string) =>
    setSets(id, (sets) => [
      ...sets,
      { reps: sets.at(-1)?.reps ?? '', weightKg: sets.at(-1)?.weightKg ?? '', doneAt: null },
    ]);

  const removeSet = (id: string, index: number) =>
    setSets(id, (sets) => sets.filter((_, i) => i !== index));

  // One tap marks the set done and timestamps it; tapping again clears it.
  // The timestamps let the backend measure real rest between sets.
  const toggleDone = (id: string, index: number) =>
    setSets(id, (sets) =>
      sets.map((s, i) => (i === index ? { ...s, doneAt: s.doneAt ? null : Date.now() } : s)),
    );

  const finish = () => {
    if (!day.data) return;
    const exercises: WorkoutExerciseResult[] = day.data.exercises.map((row) => ({
      exerciseId: row.exercise.id,
      sets: (state[row.exercise.id]?.sets ?? [])
        .filter(isLogged)
        .map((s) => ({
          reps: Number(s.reps) || 0,
          weightKg: Number(s.weightKg) || 0,
          completedAt: s.doneAt ? new Date(s.doneAt).toISOString() : undefined,
        })),
    }));
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
        const sets = state[row.exercise.id]?.sets ?? [];
        const loggedCount = sets.filter(isLogged).length;
        const complete_ = loggedCount >= row.defaultSets;
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

            <View style={{ gap: spacing.sm }}>
              <View style={styles.setHeaderRow}>
                <AppText style={[styles.colSet, styles.colLabel]}>SET</AppText>
                <AppText style={[styles.colNum, styles.colLabel]}>REPS</AppText>
                <AppText style={[styles.colNum, styles.colLabel]}>WEIGHT (KG)</AppText>
                <AppText style={[styles.colDone, styles.colLabel]}>DONE</AppText>
                <View style={styles.colRemove} />
              </View>

              {sets.map((set, i) => (
                <View key={i} style={styles.setRow}>
                  <AppText style={[styles.colSet, styles.setIndex, isLogged(set) && { color: colors.turmeric }]}>
                    {i + 1}
                  </AppText>
                  <TextInput
                    value={set.reps}
                    onChangeText={(t) => patchSet(row.exercise.id, i, { reps: t.replace(/[^0-9]/g, '') })}
                    keyboardType="numeric"
                    placeholder="0"
                    placeholderTextColor={colors.textSecondary}
                    style={[styles.numInput, styles.colNum]}
                  />
                  <TextInput
                    value={set.weightKg}
                    onChangeText={(t) => patchSet(row.exercise.id, i, { weightKg: t.replace(/[^0-9.]/g, '') })}
                    keyboardType="numeric"
                    placeholder="0"
                    placeholderTextColor={colors.textSecondary}
                    style={[styles.numInput, styles.colNum]}
                  />
                  <Pressable
                    onPress={() => toggleDone(row.exercise.id, i)}
                    hitSlop={8}
                    style={[styles.colDone, styles.doneBox, !!set.doneAt && styles.doneBoxOn]}
                    accessibilityRole="checkbox"
                    accessibilityState={{ checked: !!set.doneAt }}
                    accessibilityLabel={`Mark set ${i + 1} done`}
                  >
                    <AppText style={{ fontFamily: fonts.monoBold, fontSize: fontSize.md, color: set.doneAt ? colors.onAccent : colors.textSecondary }}>
                      {set.doneAt ? '✓' : '○'}
                    </AppText>
                  </Pressable>
                  <Pressable
                    onPress={() => removeSet(row.exercise.id, i)}
                    disabled={sets.length === 1}
                    hitSlop={8}
                    style={styles.colRemove}
                  >
                    <AppText
                      style={{
                        fontFamily: fonts.monoBold,
                        fontSize: fontSize.lg,
                        color: sets.length === 1 ? colors.border : colors.textSecondary,
                      }}
                    >
                      ✕
                    </AppText>
                  </Pressable>
                </View>
              ))}

              <Pressable onPress={() => addSet(row.exercise.id)} hitSlop={8} style={styles.addSet}>
                <AppText variant="label" color={colors.turmeric} style={{ textTransform: 'uppercase' }}>
                  + Add set
                </AppText>
              </Pressable>
            </View>
          </Card>
        );
      })}

      <AppText variant="caption" center>
        Tap ✓ as you finish each set — the timing sharpens the calorie estimate. Reps and
        weight still count if you skip it.
      </AppText>
    </Screen>
  );
}

const styles = StyleSheet.create({
  exHeader: { flexDirection: 'row', gap: spacing.md },
  setHeaderRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm },
  setRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm },
  colLabel: {
    fontFamily: fonts.bodySemibold,
    fontSize: fontSize.xs,
    color: colors.textDim,
    letterSpacing: 0.5,
  },
  colSet: { width: 24, textAlign: 'center' },
  colNum: { flex: 1 },
  colDone: { width: 34, textAlign: 'center' },
  colRemove: { width: 24, alignItems: 'center', justifyContent: 'center' },
  doneBox: {
    alignItems: 'center',
    justifyContent: 'center',
    height: 34,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.ink,
  },
  doneBoxOn: {
    backgroundColor: colors.cardamom,
    borderColor: colors.cardamom,
  },
  setIndex: { fontFamily: fonts.monoBold, fontSize: fontSize.md, color: colors.textSecondary },
  addSet: { paddingVertical: spacing.sm, alignSelf: 'flex-start' },
  numInput: {
    backgroundColor: colors.ink,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border,
    borderRadius: radius.md,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    color: colors.textLight,
    fontFamily: fonts.mono,
    fontSize: fontSize.lg,
    textAlign: 'center',
  },
});
