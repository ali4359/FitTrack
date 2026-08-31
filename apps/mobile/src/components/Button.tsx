import {
  ActivityIndicator,
  Pressable,
  StyleSheet,
  View,
  type PressableProps,
  type ViewStyle,
} from 'react-native';
import { colors, fonts, fontSize, radius, spacing } from '../theme';
import { AppText } from './AppText';

type Variant = 'primary' | 'secondary' | 'ghost' | 'danger';

type Props = Omit<PressableProps, 'style'> & {
  label: string;
  variant?: Variant;
  loading?: boolean;
  fullWidth?: boolean;
  style?: ViewStyle;
};

const FILL: Record<Variant, string> = {
  primary: colors.turmeric,
  secondary: colors.steel,
  ghost: 'transparent',
  danger: 'transparent',
};

const TEXT_COLOR: Record<Variant, string> = {
  primary: colors.onAccent,
  secondary: colors.textLight,
  ghost: colors.textDim,
  danger: colors.chili,
};

export function Button({
  label,
  variant = 'primary',
  loading,
  fullWidth,
  disabled,
  style,
  ...rest
}: Props) {
  const isDisabled = disabled || loading;
  return (
    <Pressable
      accessibilityRole="button"
      disabled={isDisabled}
      style={({ pressed }) => [
        styles.base,
        {
          backgroundColor: FILL[variant],
          borderColor: variant === 'secondary' ? colors.border : 'transparent',
          borderWidth: variant === 'secondary' ? StyleSheet.hairlineWidth : 0,
          opacity: isDisabled ? 0.5 : pressed ? 0.85 : 1,
          alignSelf: fullWidth ? 'stretch' : 'flex-start',
        },
        style,
      ]}
      {...rest}
    >
      <View style={styles.inner}>
        {loading ? <ActivityIndicator color={TEXT_COLOR[variant]} /> : null}
        <AppText
          style={{
            fontFamily: fonts.displaySemibold,
            fontSize: fontSize.md,
            color: TEXT_COLOR[variant],
            letterSpacing: 0.5,
            textTransform: 'uppercase',
          }}
        >
          {label}
        </AppText>
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  base: {
    borderRadius: radius.pill,
    paddingVertical: spacing.md + 2,
    paddingHorizontal: spacing.xl,
    minHeight: 52,
    justifyContent: 'center',
  },
  inner: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing.sm,
  },
});
