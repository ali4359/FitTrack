export { colors, goalAccent, type ColorName } from './colors';
export { fonts, fontSize } from './typography';
export { spacing, radius } from './spacing';

import { colors } from './colors';
import type { Theme } from '@react-navigation/native';

/** React Navigation theme so headers / containers match the dark UI. */
export const navTheme: Theme = {
  dark: true,
  colors: {
    primary: colors.turmeric,
    background: colors.ink,
    card: colors.ink,
    text: colors.textLight,
    border: colors.border,
    notification: colors.chili,
  },
  fonts: {
    regular: { fontFamily: 'Inter_400Regular', fontWeight: '400' },
    medium: { fontFamily: 'Inter_500Medium', fontWeight: '500' },
    bold: { fontFamily: 'Inter_600SemiBold', fontWeight: '600' },
    heavy: { fontFamily: 'Oswald_700Bold', fontWeight: '700' },
  },
};
