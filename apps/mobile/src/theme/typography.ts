/**
 * Font families. Keys match the names registered with useFonts() in App.tsx.
 *
 * - Oswald  → display / headers (condensed, high energy)
 * - Inter   → body / UI
 * - JetBrainsMono → numeric data ONLY (calories, macros, reps, weights, timers)
 */
export const fonts = {
  displayBold: 'Oswald_700Bold',
  displaySemibold: 'Oswald_600SemiBold',
  display: 'Oswald_500Medium',

  body: 'Inter_400Regular',
  bodyMedium: 'Inter_500Medium',
  bodySemibold: 'Inter_600SemiBold',

  mono: 'JetBrainsMono_500Medium',
  monoBold: 'JetBrainsMono_700Bold',
} as const;

export const fontSize = {
  xs: 12,
  sm: 14,
  md: 16,
  lg: 18,
  xl: 22,
  xxl: 28,
  display: 34,
  hero: 52,
} as const;
