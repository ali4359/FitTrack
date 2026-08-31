/**
 * Iron & Spice palette — from the approved mockup. Dark UI only.
 * Do not introduce new colors without updating the mockup.
 */
export const colors = {
  // surfaces
  ink: '#17181B', // app background
  steel: '#2A2D32', // card surface
  border: '#3A3E45',

  // text
  textLight: '#F3F0E7',
  textDim: '#B8B4A8',
  textSecondary: '#8A8F98',

  // accents
  turmeric: '#E8A33D', // primary / CTAs
  chili: '#D64933', // bulk / intensity
  cardamom: '#5C8A66', // cut / health / success

  // fixed
  onAccent: '#17181B', // text/icons sitting on a turmeric fill
  danger: '#D64933',
} as const;

export type ColorName = keyof typeof colors;

/** Map a user goal to its accent color. */
export const goalAccent: Record<'bulk' | 'cut' | 'maintain', string> = {
  bulk: colors.chili,
  cut: colors.cardamom,
  maintain: colors.turmeric,
};
