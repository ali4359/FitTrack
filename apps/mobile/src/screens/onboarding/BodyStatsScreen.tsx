import { useState } from 'react';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { AppText, Button, Field, Screen } from '../../components';
import { spacing } from '../../theme';
import { useOnboardingStore } from '../../store/onboardingStore';
import type { OnboardingStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<OnboardingStackParamList, 'BodyStats'>;

export function BodyStatsScreen({ navigation }: Props) {
  const weightKg = useOnboardingStore((s) => s.weightKg);
  const set = useOnboardingStore((s) => s.set);
  const [text, setText] = useState(weightKg ? String(weightKg) : '');

  const parsed = Number(text);
  const valid = text.trim() !== '' && Number.isFinite(parsed) && parsed >= 30 && parsed <= 300;

  return (
    <Screen
      scroll
      footer={
        <Button
          label="Continue"
          fullWidth
          disabled={!valid}
          onPress={() => {
            set({ weightKg: parsed });
            navigation.navigate('BudgetRegion');
          }}
        />
      }
    >
      <AppText variant="display">Your body weight</AppText>
      <AppText variant="subtitle" style={{ marginBottom: spacing.sm }}>
        We use this to estimate the calories you burn each workout and to set your
        daily targets. You can change it any time in your profile.
      </AppText>

      <Field
        label="Weight (kg)"
        keyboardType="decimal-pad"
        value={text}
        onChangeText={(t) => setText(t.replace(/[^0-9.]/g, ''))}
        placeholder="e.g. 74.5"
        error={text.trim() !== '' && !valid ? 'Enter a weight between 30 and 300 kg' : null}
      />
    </Screen>
  );
}
