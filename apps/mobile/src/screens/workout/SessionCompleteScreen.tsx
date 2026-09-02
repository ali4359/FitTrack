import { View } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { AppText, Button, Card, Screen, Tag } from '../../components';
import { colors, spacing } from '../../theme';
import { DAY_LABELS } from '../../lib/workoutPlan';
import type { HomeStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<HomeStackParamList, 'SessionComplete'>;

/**
 * Placeholder — full build is screen 5 in the plan. Shows the estimated burn and
 * routes into the post-workout meal so the end-to-end flow is testable now.
 */
export function SessionCompleteScreen({ navigation, route }: Props) {
  const { caloriesBurned, durationMinutes, workoutDayId } = route.params;

  return (
    <Screen
      contentStyle={{ flexGrow: 1, justifyContent: 'center', alignItems: 'center' }}
      footer={
        <View style={{ gap: spacing.sm }}>
          <Button
            label="See post-workout meal"
            fullWidth
            onPress={() => navigation.replace('MealSuggestion', { mealType: 'post-workout' })}
          />
          <Button
            label="Back to today"
            variant="secondary"
            fullWidth
            onPress={() => navigation.navigate('Today')}
          />
        </View>
      }
    >
      <Tag label="SESSION COMPLETE" tone="cardamom" />
      <AppText variant="subtitle" center>
        {DAY_LABELS[workoutDayId] ?? 'Workout logged'}
      </AppText>

      <Card padded style={{ alignItems: 'center', gap: spacing.xs, alignSelf: 'stretch' }}>
        <AppText variant="label" style={{ textTransform: 'uppercase' }}>
          estimated burn
        </AppText>
        <AppText variant="hero" color={colors.turmeric}>
          {Math.round(caloriesBurned)}
        </AppText>
        <AppText variant="mono" color={colors.textDim}>
          kcal · {durationMinutes} min
        </AppText>
      </Card>

      <AppText variant="caption" center>
        Estimated from the work you logged and your set timing — expect roughly ±25%.
      </AppText>
    </Screen>
  );
}
