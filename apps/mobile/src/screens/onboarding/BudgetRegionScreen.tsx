import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import type { BudgetTier } from '@iron-and-spice/shared';
import { AppText, Button, Field, OptionCard, Screen } from '../../components';
import { spacing } from '../../theme';
import { useOnboardingStore } from '../../store/onboardingStore';
import type { OnboardingStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<OnboardingStackParamList, 'BudgetRegion'>;

const TIERS: { value: BudgetTier; title: string; description: string }[] = [
  { value: 'low', title: 'Low', description: 'Daal, chana, seasonal sabzi, eggs.' },
  { value: 'mid', title: 'Mid', description: 'Chicken a few times a week, some beef.' },
  { value: 'high', title: 'High', description: 'Meat daily, more variety, less price-sensitive.' },
];

export function BudgetRegionScreen({ navigation }: Props) {
  const { budgetTier, region } = useOnboardingStore();
  const set = useOnboardingStore((s) => s.set);

  const canContinue = !!budgetTier && region.trim().length > 1;

  return (
    <Screen
      scroll
      footer={
        <Button
          label="Continue"
          fullWidth
          disabled={!canContinue}
          onPress={() => navigation.navigate('Restrictions')}
        />
      }
    >
      <AppText variant="display">Budget &amp; region</AppText>
      <AppText variant="subtitle" style={{ marginBottom: spacing.sm }}>
        We only suggest food that&apos;s actually cheap and available where you are.
      </AppText>

      <Field
        label="City / region"
        value={region}
        onChangeText={(t) => set({ region: t })}
        placeholder="e.g. Lahore"
        autoCapitalize="words"
      />

      <AppText variant="label" style={{ textTransform: 'uppercase', marginTop: spacing.sm }}>
        Food budget
      </AppText>
      {TIERS.map((t) => (
        <OptionCard
          key={t.value}
          title={t.title}
          description={t.description}
          selected={budgetTier === t.value}
          onPress={() => set({ budgetTier: t.value })}
        />
      ))}
    </Screen>
  );
}
