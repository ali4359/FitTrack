import { Text, type TextProps, type TextStyle } from 'react-native';
import { colors, fonts, fontSize } from '../theme';

type Variant =
  | 'hero' // giant numbers (calories burned)
  | 'display' // screen titles
  | 'title' // section / card titles
  | 'subtitle'
  | 'body'
  | 'bodyMedium'
  | 'label' // pill text, small caps-ish labels
  | 'caption'
  | 'mono' // numeric data
  | 'monoLarge';

const VARIANT: Record<Variant, TextStyle> = {
  hero: { fontFamily: fonts.monoBold, fontSize: fontSize.hero, color: colors.textLight },
  display: { fontFamily: fonts.displayBold, fontSize: fontSize.display, color: colors.textLight },
  title: { fontFamily: fonts.displaySemibold, fontSize: fontSize.xl, color: colors.textLight },
  subtitle: { fontFamily: fonts.bodyMedium, fontSize: fontSize.md, color: colors.textDim },
  body: { fontFamily: fonts.body, fontSize: fontSize.md, color: colors.textLight },
  bodyMedium: { fontFamily: fonts.bodyMedium, fontSize: fontSize.md, color: colors.textLight },
  label: { fontFamily: fonts.bodySemibold, fontSize: fontSize.xs, color: colors.textDim, letterSpacing: 0.5 },
  caption: { fontFamily: fonts.body, fontSize: fontSize.sm, color: colors.textSecondary },
  mono: { fontFamily: fonts.mono, fontSize: fontSize.md, color: colors.textLight },
  monoLarge: { fontFamily: fonts.monoBold, fontSize: fontSize.xxl, color: colors.textLight },
};

type Props = TextProps & {
  variant?: Variant;
  color?: string;
  center?: boolean;
};

export function AppText({ variant = 'body', color, center, style, ...rest }: Props) {
  return (
    <Text
      style={[VARIANT[variant], color ? { color } : null, center ? { textAlign: 'center' } : null, style]}
      {...rest}
    />
  );
}
