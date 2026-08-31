import { View } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { AppText, Button, Card, Screen, Tag } from '../../components';
import { apiErrorMessage } from '../../api/client';
import { useMealSuggestions } from '../../hooks/useMeals';
import { colors, spacing } from '../../theme';
import type { HomeStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<HomeStackParamList, 'MealSuggestion'>;

/** Placeholder — full build (swap / log / image) is screen 6 in the plan. */
export function MealSuggestionScreen({ route }: Props) {
  const { mealType } = route.params;
  const meals = useMealSuggestions(mealType);
  const top = meals.data?.[0];

  return (
    <Screen scroll refreshing={meals.isRefetching} onRefresh={() => meals.refetch()}>
      {meals.isLoading ? <AppText variant="subtitle">Finding something for you…</AppText> : null}

      {meals.isError ? (
        <AppText variant="caption" color={colors.chili}>
          {apiErrorMessage(meals.error)}
        </AppText>
      ) : null}

      {top ? (
        <Card padded style={{ gap: spacing.md }}>
          <View style={{ height: 140, borderRadius: 12, backgroundColor: colors.ink, alignItems: 'center', justifyContent: 'center' }}>
            <AppText variant="label">IMAGE</AppText>
          </View>
          <AppText variant="title">{top.dishName}</AppText>
          <View style={{ flexDirection: 'row', gap: spacing.sm, flexWrap: 'wrap' }}>
            <Tag label={top.region} tone="turmeric" />
            <Tag label={`${top.budgetTier} budget`} tone="neutral" />
            {top.halal ? <Tag label="halal" tone="cardamom" /> : null}
          </View>
          <View style={{ flexDirection: 'row', gap: spacing.lg }}>
            <AppText variant="mono">{Math.round(top.calories)} kcal</AppText>
            <AppText variant="mono" color={colors.textDim}>
              P{Math.round(top.proteinG)} · C{Math.round(top.carbsG)} · F{Math.round(top.fatG)}
            </AppText>
          </View>
        </Card>
      ) : null}

      <AppText variant="caption" center>
        Swap &amp; log actions come with the full Meal Suggestion screen.
      </AppText>
    </Screen>
  );
}
