import { Pressable, StyleSheet, View } from 'react-native';
import { colors, radius, spacing } from '../theme';
import { AppText } from './AppText';

type Props = {
  title: string;
  description?: string;
  selected: boolean;
  onPress: () => void;
  accent?: string;
};

export function OptionCard({ title, description, selected, onPress, accent = colors.turmeric }: Props) {
  return (
    <Pressable
      accessibilityRole="radio"
      accessibilityState={{ selected }}
      onPress={onPress}
      style={[
        styles.card,
        { borderColor: selected ? accent : colors.border, backgroundColor: selected ? 'rgba(232,163,61,0.08)' : colors.steel },
      ]}
    >
      <View style={styles.row}>
        <View style={styles.textCol}>
          <AppText variant="title" style={{ fontSize: 18 }}>
            {title}
          </AppText>
          {description ? (
            <AppText variant="caption" style={{ marginTop: 2 }}>
              {description}
            </AppText>
          ) : null}
        </View>
        <View style={[styles.dot, { borderColor: selected ? accent : colors.textSecondary }]}>
          {selected ? <View style={[styles.dotFill, { backgroundColor: accent }]} /> : null}
        </View>
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    borderRadius: radius.md,
    borderWidth: 1.5,
    padding: spacing.lg,
  },
  row: { flexDirection: 'row', alignItems: 'center', gap: spacing.md },
  textCol: { flex: 1 },
  dot: {
    width: 22,
    height: 22,
    borderRadius: radius.pill,
    borderWidth: 2,
    alignItems: 'center',
    justifyContent: 'center',
  },
  dotFill: { width: 10, height: 10, borderRadius: radius.pill },
});
