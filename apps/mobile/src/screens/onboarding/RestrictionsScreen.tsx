import { Pressable, StyleSheet, View } from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { AppText, Button, Screen } from '../../components';
import { apiErrorMessage } from '../../api/client';
import { useUpdateProfile } from '../../hooks/useAuth';
import { colors, radius, spacing } from '../../theme';
import { useOnboardingStore } from '../../store/onboardingStore';
import type { OnboardingStackParamList } from '../../navigation/types';

type Props = NativeStackScreenProps<OnboardingStackParamList, 'Restrictions'>;

const OPTIONS = ['halal', 'vegetarian', 'no beef', 'lactose-free', 'no eggs'];

export function RestrictionsScreen(_: Props) {
  const { goal, budgetTier, region, restrictions } = useOnboardingStore();
  const set = useOnboardingStore((s) => s.set);
  const reset = useOnboardingStore((s) => s.reset);
  const updateProfile = useUpdateProfile();

  const toggle = (opt: string) => {
    set({
      restrictions: restrictions.includes(opt)
        ? restrictions.filter((r) => r !== opt)
        : [...restrictions, opt],
    });
  };

  const finish = () => {
    if (!goal || !budgetTier) return;
    updateProfile.mutate(
      {
        goal,
        budgetTier,
        region: region.trim(),
        restrictions: restrictions.join(','),
      },
      { onSuccess: () => reset() },
    );
    // On success, RootNavigator sees a complete profile and swaps to the tabs.
  };

  return (
    <Screen
      scroll
      footer={
        <Button
          label="Finish setup"
          fullWidth
          loading={updateProfile.isPending}
          onPress={finish}
        />
      }
    >
      <AppText variant="display">Anything we should avoid?</AppText>
      <AppText variant="subtitle" style={{ marginBottom: spacing.sm }}>
        Pick all that apply. You can change these later in your profile.
      </AppText>

      <View style={styles.chips}>
        {OPTIONS.map((opt) => {
          const on = restrictions.includes(opt);
          return (
            <Pressable
              key={opt}
              onPress={() => toggle(opt)}
              style={[styles.chip, { borderColor: on ? colors.turmeric : colors.border, backgroundColor: on ? 'rgba(232,163,61,0.12)' : colors.steel }]}
            >
              <AppText color={on ? colors.turmeric : colors.textDim} variant="bodyMedium">
                {opt}
              </AppText>
            </Pressable>
          );
        })}
      </View>

      {updateProfile.isError ? (
        <AppText variant="caption" color={colors.chili}>
          {apiErrorMessage(updateProfile.error, 'Could not save your profile')}
        </AppText>
      ) : null}
    </Screen>
  );
}

const styles = StyleSheet.create({
  chips: { flexDirection: 'row', flexWrap: 'wrap', gap: spacing.sm },
  chip: {
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderRadius: radius.pill,
    borderWidth: 1.5,
  },
});
