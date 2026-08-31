import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { colors, fonts } from '../theme';
import { GoalScreen } from '../screens/onboarding/GoalScreen';
import { BudgetRegionScreen } from '../screens/onboarding/BudgetRegionScreen';
import { RestrictionsScreen } from '../screens/onboarding/RestrictionsScreen';
import type { OnboardingStackParamList } from './types';

const Stack = createNativeStackNavigator<OnboardingStackParamList>();

export function OnboardingStack() {
  return (
    <Stack.Navigator
      screenOptions={{
        headerStyle: { backgroundColor: colors.ink },
        headerTitleStyle: { fontFamily: fonts.displaySemibold, color: colors.textLight },
        headerTintColor: colors.turmeric,
        headerShadowVisible: false,
        contentStyle: { backgroundColor: colors.ink },
      }}
    >
      <Stack.Screen name="Goal" component={GoalScreen} options={{ title: 'Your goal' }} />
      <Stack.Screen
        name="BudgetRegion"
        component={BudgetRegionScreen}
        options={{ title: 'Budget & region' }}
      />
      <Stack.Screen
        name="Restrictions"
        component={RestrictionsScreen}
        options={{ title: 'Dietary needs' }}
      />
    </Stack.Navigator>
  );
}
