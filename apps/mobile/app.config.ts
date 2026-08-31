import type { ExpoConfig } from 'expo/config';

/**
 * Expo loads `.env` into process.env before evaluating this file, so
 * API_BASE_URL from apps/mobile/.env lands in `extra.apiBaseUrl` and is
 * read at runtime via expo-constants (see src/config.ts).
 */
const apiBaseUrl =
  process.env.API_BASE_URL ??
  process.env.EXPO_PUBLIC_API_BASE_URL ??
  'http://localhost:8080';

const config: ExpoConfig = {
  name: 'Iron & Spice',
  slug: 'iron-and-spice',
  version: '1.0.0',
  orientation: 'portrait',
  icon: './assets/icon.png',
  scheme: 'ironandspice',
  userInterfaceStyle: 'dark',
  backgroundColor: '#17181B',
  ios: {
    supportsTablet: true,
    bundleIdentifier: 'app.ironandspice.mobile',
  },
  android: {
    package: 'app.ironandspice.mobile',
    adaptiveIcon: {
      foregroundImage: './assets/android-icon-foreground.png',
      backgroundColor: '#17181B',
    },
  },
  web: {
    favicon: './assets/favicon.png',
  },
  plugins: [
    'expo-secure-store',
    'expo-font',
    [
      'expo-splash-screen',
      {
        image: './assets/splash-icon.png',
        resizeMode: 'contain',
        backgroundColor: '#17181B',
      },
    ],
  ],
  extra: {
    apiBaseUrl,
  },
};

export default config;
