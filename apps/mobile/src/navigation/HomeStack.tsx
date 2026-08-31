import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { colors, fonts } from '../theme';
import { HomeScreen } from '../screens/home/HomeScreen';
import { WorkoutSessionScreen } from '../screens/workout/WorkoutSessionScreen';
import { SessionCompleteScreen } from '../screens/workout/SessionCompleteScreen';
import { MealSuggestionScreen } from '../screens/meals/MealSuggestionScreen';
import type { HomeStackParamList } from './types';

const Stack = createNativeStackNavigator<HomeStackParamList>();

export function HomeStack() {
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
      <Stack.Screen name="Today" component={HomeScreen} options={{ headerShown: false }} />
      <Stack.Screen
        name="WorkoutSession"
        component={WorkoutSessionScreen}
        options={{ title: 'Workout', headerBackVisible: false }}
      />
      <Stack.Screen
        name="SessionComplete"
        component={SessionCompleteScreen}
        options={{ headerShown: false }}
      />
      <Stack.Screen
        name="MealSuggestion"
        component={MealSuggestionScreen}
        options={{ title: 'Post-workout meal' }}
      />
    </Stack.Navigator>
  );
}
