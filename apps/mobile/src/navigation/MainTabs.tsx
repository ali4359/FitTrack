import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { StyleSheet, View } from 'react-native';
import { colors, fonts } from '../theme';
import { AppText } from '../components/AppText';
import { HomeStack } from './HomeStack';
import { ProgressScreen } from '../screens/progress/ProgressScreen';
import { ProfileScreen } from '../screens/profile/ProfileScreen';
import type { MainTabParamList } from './types';

const Tab = createBottomTabNavigator<MainTabParamList>();

// Text-glyph tab icons keep us dependency-free until an icon set is wired up.
function TabIcon({ glyph, focused }: { glyph: string; focused: boolean }) {
  return (
    <View style={styles.icon}>
      <AppText
        style={{
          fontFamily: fonts.displaySemibold,
          fontSize: 18,
          color: focused ? colors.turmeric : colors.textSecondary,
        }}
      >
        {glyph}
      </AppText>
    </View>
  );
}

export function MainTabs() {
  return (
    <Tab.Navigator
      screenOptions={{
        headerShown: false,
        tabBarStyle: {
          backgroundColor: colors.ink,
          borderTopColor: colors.border,
          borderTopWidth: StyleSheet.hairlineWidth,
          height: 64,
          paddingTop: 6,
        },
        tabBarLabelStyle: { fontFamily: fonts.bodySemibold, fontSize: 11 },
        tabBarActiveTintColor: colors.turmeric,
        tabBarInactiveTintColor: colors.textSecondary,
      }}
    >
      <Tab.Screen
        name="HomeTab"
        component={HomeStack}
        options={{
          title: 'Today',
          tabBarIcon: ({ focused }) => <TabIcon glyph="◈" focused={focused} />,
        }}
      />
      <Tab.Screen
        name="ProgressTab"
        component={ProgressScreen}
        options={{
          title: 'Progress',
          tabBarIcon: ({ focused }) => <TabIcon glyph="▲" focused={focused} />,
        }}
      />
      <Tab.Screen
        name="ProfileTab"
        component={ProfileScreen}
        options={{
          title: 'Profile',
          tabBarIcon: ({ focused }) => <TabIcon glyph="●" focused={focused} />,
        }}
      />
    </Tab.Navigator>
  );
}

const styles = StyleSheet.create({
  icon: { height: 20, alignItems: 'center', justifyContent: 'center' },
});
