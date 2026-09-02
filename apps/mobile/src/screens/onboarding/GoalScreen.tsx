import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import type { Goal } from '@fittrack/shared';
import { AppText, Button, OptionCard, Screen } from '../../components';
import { goalAccent, spacing } from '../../theme';
import { useOnboardingStore } from '../../store/onboardingStore';
import type { OnboardingStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<OnboardingStackParamList, 'Goal'>;

const GOALS: { value: Goal; title: string; description: string }[] = [
  { value: 'bulk', title: 'Bulk', description: 'Add size and strength. Eat in a surplus.' },
  { value: 'cut', title: 'Cut', description: 'Lean out while keeping muscle. Eat in a deficit.' },
  { value: 'maintain', title: 'Maintain', description: 'Hold your weight, train for performance.' },
];

export function GoalScreen({ navigation }: Props) {
  const goal = useOnboardingStore((s) => s.goal);
  const set = useOnboardingStore((s) => s.set);

  return (
    <Screen
      scroll
      footer={
        <Button
          label="Continue"
          fullWidth
          disabled={!goal}
          onPress={() => navigation.navigate('BodyStats')}
        />
      }
    >
      <AppText variant="display">What are you training for?</AppText>
      <AppText variant="subtitle" style={{ marginBottom: spacing.sm }}>
        This shapes your workouts and the meals we suggest.
      </AppText>

      {GOALS.map((g) => (
        <OptionCard
          key={g.value}
          title={g.title}
          description={g.description}
          selected={goal === g.value}
          accent={goalAccent[g.value]}
          onPress={() => set({ goal: g.value })}
        />
      ))}
    </Screen>
  );
}
