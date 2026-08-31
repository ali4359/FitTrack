import { StyleSheet, View } from 'react-native';
import { colors, fonts, fontSize, radius, spacing } from '../theme';
import { AppText } from './AppText';

type Props = {
  label: string;
  tone?: 'neutral' | 'turmeric' | 'chili' | 'cardamom';
};

const TONE: Record<NonNullable<Props['tone']>, { bg: string; fg: string }> = {
  neutral: { bg: colors.ink, fg: colors.textDim },
  turmeric: { bg: 'rgba(232,163,61,0.16)', fg: colors.turmeric },
  chili: { bg: 'rgba(214,73,51,0.16)', fg: colors.chili },
  cardamom: { bg: 'rgba(92,138,102,0.18)', fg: colors.cardamom },
};

export function Tag({ label, tone = 'neutral' }: Props) {
  const t = TONE[tone];
  return (
    <View style={[styles.pill, { backgroundColor: t.bg }]}>
      <AppText
        style={{
          fontFamily: fonts.bodySemibold,
          fontSize: fontSize.xs,
          color: t.fg,
          letterSpacing: 0.4,
        }}
      >
        {label}
      </AppText>
    </View>
  );
}

const styles = StyleSheet.create({
  pill: {
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.xs + 1,
    borderRadius: radius.pill,
    alignSelf: 'flex-start',
  },
});
