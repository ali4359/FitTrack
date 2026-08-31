import { AppText, Button, Card, Screen, Tag } from '../../components';
import { useLogout } from '../../hooks/useAuth';
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
