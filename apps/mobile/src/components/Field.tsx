import { forwardRef } from 'react';
import { StyleSheet, TextInput, View, type TextInputProps } from 'react-native';
import { colors, fonts, fontSize, radius, spacing } from '../theme';
import { AppText } from './AppText';

type Props = TextInputProps & {
  label: string;
  error?: string | null;
};

export const Field = forwardRef<TextInput, Props>(function Field(
  { label, error, style, ...rest },
  ref,
) {
  return (
    <View style={styles.wrap}>
      <AppText variant="label" style={{ textTransform: 'uppercase' }}>
        {label}
      </AppText>
      <TextInput
        ref={ref}
        placeholderTextColor={colors.textSecondary}
        style={[styles.input, !!error && styles.inputError, style]}
        {...rest}
      />
      {error ? (
        <AppText variant="caption" color={colors.chili}>
          {error}
        </AppText>
      ) : null}
    </View>
  );
});

const styles = StyleSheet.create({
  wrap: { gap: spacing.sm },
  input: {
    backgroundColor: colors.steel,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border,
    borderRadius: radius.md,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    color: colors.textLight,
    fontFamily: fonts.body,
    fontSize: fontSize.md,
    minHeight: 50,
  },
  inputError: {
    borderColor: colors.chili,
  },
});
