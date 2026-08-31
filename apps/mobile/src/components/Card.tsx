import { StyleSheet, View, type ViewProps } from 'react-native';
import { colors, radius, spacing } from '../theme';

type Props = ViewProps & {
  padded?: boolean;
};

export function Card({ padded = true, style, ...rest }: Props) {
  return <View style={[styles.card, padded && styles.padded, style]} {...rest} />;
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.steel,
    borderRadius: radius.lg,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: colors.border,
  },
  padded: {
    padding: spacing.lg,
  },
});
