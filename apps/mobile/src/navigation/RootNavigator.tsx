import { useEffect } from 'react';
import { ActivityIndicator, View } from 'react-native';
import { NavigationContainer } from '@react-navigation/native';
import { colors, navTheme } from '../theme';
import { needsOnboarding, useAuthStore } from '../store/authStore';
import { AuthStack } from './AuthStack';
import { OnboardingStack } from './OnboardingStack';
import { MainTabs } from './MainTabs';

export function RootNavigator() {
  const hydrated = useAuthStore((s) => s.hydrated);
  const hydrate = useAuthStore((s) => s.hydrate);
  const token = useAuthStore((s) => s.token);
  const user = useAuthStore((s) => s.user);

  useEffect(() => {
    void hydrate();
  }, [hydrate]);

  if (!hydrated) {
    return (
      <View style={{ flex: 1, backgroundColor: colors.ink, alignItems: 'center', justifyContent: 'center' }}>
        <ActivityIndicator color={colors.turmeric} />
      </View>
    );
  }

  return (
    <NavigationContainer theme={navTheme}>
      {!token ? <AuthStack /> : needsOnboarding(user) ? <OnboardingStack /> : <MainTabs />}
    </NavigationContainer>
  );
}
