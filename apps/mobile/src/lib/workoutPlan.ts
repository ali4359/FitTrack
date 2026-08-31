import type { WorkoutLog } from '@fittrack/shared';

/** The seeded rotating split. Replace when the backend serves a real plan. */
export const PLAN_DAYS = ['day-1', 'day-2', 'day-3'] as const;

export const DAY_LABELS: Record<string, string> = {
  'day-1': 'Day 1 — Chest & Triceps',
  'day-2': 'Day 2 — Back & Biceps',
  'day-3': 'Day 3 — Legs',
};

function isSameCalendarDay(a: Date, b: Date) {
  return (
    a.getFullYear() === b.getFullYear() &&
    a.getMonth() === b.getMonth() &&
    a.getDate() === b.getDate()
  );
}

/** Next day in the rotation after the most recent completed one. */
export function todaysWorkoutDayId(history: WorkoutLog[] | undefined): string {
  if (!history || history.length === 0) return PLAN_DAYS[0];
  const last = [...history].sort(
    (a, b) => +new Date(b.completedAt) - +new Date(a.completedAt),
  )[0];

  // Already trained today → keep showing the same day.
  if (isSameCalendarDay(new Date(last.completedAt), new Date())) {
    return last.workoutDayId;
  }
  const idx = PLAN_DAYS.indexOf(last.workoutDayId as (typeof PLAN_DAYS)[number]);
  return PLAN_DAYS[(idx + 1) % PLAN_DAYS.length] ?? PLAN_DAYS[0];
}

export function trainedToday(history: WorkoutLog[] | undefined): boolean {
  if (!history) return false;
  const now = new Date();
  return history.some((l) => isSameCalendarDay(new Date(l.completedAt), now));
}

/** Consecutive-day streak ending today or yesterday. */
export function computeStreak(history: WorkoutLog[] | undefined): number {
  if (!history || history.length === 0) return 0;
  const days = new Set(
    history.map((l) => {
      const d = new Date(l.completedAt);
      return `${d.getFullYear()}-${d.getMonth()}-${d.getDate()}`;
    }),
  );
  const key = (d: Date) => `${d.getFullYear()}-${d.getMonth()}-${d.getDate()}`;

  const cursor = new Date();
  if (!days.has(key(cursor))) {
    cursor.setDate(cursor.getDate() - 1);
    if (!days.has(key(cursor))) return 0;
  }
  let streak = 0;
  while (days.has(key(cursor))) {
    streak += 1;
    cursor.setDate(cursor.getDate() - 1);
  }
  return streak;
}

/** Mon-first 7-slot view of the current week. */
export function weekStrip(history: WorkoutLog[] | undefined) {
  const now = new Date();
  const monday = new Date(now);
  const dow = (now.getDay() + 6) % 7; // 0 = Monday
  monday.setDate(now.getDate() - dow);
  monday.setHours(0, 0, 0, 0);

  const trained = new Set(
    (history ?? []).map((l) => {
      const d = new Date(l.completedAt);
      return `${d.getFullYear()}-${d.getMonth()}-${d.getDate()}`;
    }),
  );
  const key = (d: Date) => `${d.getFullYear()}-${d.getMonth()}-${d.getDate()}`;

  return ['M', 'T', 'W', 'T', 'F', 'S', 'S'].map((label, i) => {
    const d = new Date(monday);
    d.setDate(monday.getDate() + i);
    return {
      label,
      done: trained.has(key(d)),
      isToday: isSameCalendarDay(d, now),
      future: d > now && !isSameCalendarDay(d, now),
    };
  });
}
