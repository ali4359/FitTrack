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
  const data = meals.data;
  const results = data?.results ?? [];

  return (
    <Screen scroll refreshing={meals.isRefetching} onRefresh={() => meals.refetch()}>
      {meals.isLoading ? <AppText variant="subtitle">Finding something for you…</AppText> : null}

      {meals.isError ? (
        <AppText variant="caption" color={colors.chili}>
          {apiErrorMessage(meals.error)}
        </AppText>
      ) : null}

      {data ? (
        <AppText variant="caption" color={colors.textDim}>
          This meal ≈ {Math.round(data.target.calories)} kcal · P{Math.round(data.target.proteinG)} ·
          C{Math.round(data.target.carbsG)} · F{Math.round(data.target.fatG)}
          {data.broadened ? '  (widened to fill 3)' : ''}
        </AppText>
      ) : null}

      {results.map((dish) => (
        <Card key={dish.name} padded style={{ gap: spacing.sm }}>
          <View style={{ flexDirection: 'row', gap: spacing.sm, flexWrap: 'wrap' }}>
            {dish.role ? <Tag label={dish.role} tone="turmeric" /> : null}
            {dish.isHalal ? <Tag label="halal" tone="cardamom" /> : null}
            {dish.isVegetarian ? <Tag label="veg" tone="neutral" /> : null}
          </View>
          <AppText variant="title">{dish.name}</AppText>
          <AppText variant="caption" color={colors.textDim}>
            {dish.portion}
          </AppText>
          <View style={{ flexDirection: 'row', gap: spacing.lg }}>
            <AppText variant="mono">{Math.round(dish.calories)} kcal*</AppText>
            <AppText variant="mono" color={colors.textDim}>
              P{Math.round(dish.proteinG)} · C{Math.round(dish.carbsG)} · F{Math.round(dish.fatG)}
            </AppText>
          </View>
          <AppText variant="caption">{dish.whyItFits}</AppText>
        </Card>
      ))}

      <AppText variant="caption" center color={colors.textDim}>
        * estimated macros. Swap &amp; log actions come with the full Meal Suggestion screen.
      </AppText>
    </Screen>
  );
}
