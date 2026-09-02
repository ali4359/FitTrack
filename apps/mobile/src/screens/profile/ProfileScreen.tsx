import { useState } from 'react';
import { View } from 'react-native';
import { AppText, Button, Card, Field, Screen, Tag } from '../../components';
import { apiErrorMessage } from '../../api/client';
import { useLogout, useUpdateProfile } from '../../hooks/useAuth';
import { useAuthStore } from '../../store/authStore';
import { colors, spacing } from '../../theme';

/** Placeholder — full build (edit goal/budget/restrictions) is screen 8 in the plan. */
export function ProfileScreen() {
  const user = useAuthStore((s) => s.user);
  const logout = useLogout();

  return (
    <Screen
      scroll
      contentStyle={{ paddingTop: spacing.xxl }}
      footer={<Button label="Log out" variant="danger" fullWidth onPress={() => void logout()} />}
    >
      <AppText variant="display">{user?.name ?? 'Profile'}</AppText>
      <AppText variant="subtitle">{user?.email}</AppText>

      <WeightCard currentKg={user?.weightKg ?? 0} />

      <Card padded style={{ gap: spacing.md }}>
        <Row label="Goal" value={user?.goal} />
        <Row label="Region" value={user?.region} />
        <Row label="Budget" value={user?.budgetTier} />
        <Row label="Restrictions" value={user?.restrictions || 'none'} />
      </Card>

      <AppText variant="caption">Editing these fields comes with the full Profile screen.</AppText>
    </Screen>
  );
}

function WeightCard({ currentKg }: { currentKg: number }) {
  const updateProfile = useUpdateProfile();
  const [value, setValue] = useState(currentKg > 0 ? String(currentKg) : '');
  const [saved, setSaved] = useState(false);

  const parsed = Number(value);
  const valid = value.trim() !== '' && Number.isFinite(parsed) && parsed > 0 && parsed < 500;
  const changed = valid && parsed !== currentKg;

  const save = () => {
    if (!changed) return;
    setSaved(false);
    updateProfile.mutate(
      { weightKg: parsed },
      { onSuccess: () => setSaved(true) },
    );
  };

  return (
    <Card padded style={{ gap: spacing.md }}>
      <Field
        label="Body weight (kg)"
        keyboardType="decimal-pad"
        value={value}
        onChangeText={(t) => {
          setValue(t.replace(/[^0-9.]/g, ''));
          setSaved(false);
        }}
        placeholder="e.g. 74.5"
        error={value.trim() !== '' && !valid ? 'Enter a weight between 1 and 500 kg' : null}
      />
      <AppText variant="caption">
        Used to estimate the calories you burn each workout.
      </AppText>

      <View style={{ flexDirection: 'row', alignItems: 'center', gap: spacing.md }}>
        <Button
          label="Save weight"
          disabled={!changed || updateProfile.isPending}
          loading={updateProfile.isPending}
          onPress={save}
        />
        {saved && !changed ? (
          <AppText variant="caption" color={colors.cardamom}>
            Saved
          </AppText>
        ) : null}
      </View>

      {updateProfile.isError ? (
        <AppText variant="caption" color={colors.chili}>
          {apiErrorMessage(updateProfile.error, 'Could not save your weight')}
        </AppText>
      ) : null}
    </Card>
  );
}

function Row({ label, value }: { label: string; value?: string }) {
  return (
    <>
      <AppText variant="label" style={{ textTransform: 'uppercase' }}>
        {label}
      </AppText>
      <Tag label={value ?? '—'} tone="neutral" />
    </>
  );
}
