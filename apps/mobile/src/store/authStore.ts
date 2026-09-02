import AsyncStorage from '@react-native-async-storage/async-storage';
import * as SecureStore from 'expo-secure-store';
import { create } from 'zustand';
import type { User } from '@fittrack/shared';

const TOKEN_KEY = 'fittrack.jwt';
const USER_KEY = 'fittrack.user';

type AuthState = {
  token: string | null;
  user: User | null;
  /** false until we've tried to restore a session from storage */
  hydrated: boolean;

  hydrate: () => Promise<void>;
  signIn: (token: string, user: User) => Promise<void>;
  setUser: (user: User) => Promise<void>;
  signOut: () => Promise<void>;
};

/** Whether onboarding still needs to run for the signed-in user. */
export function needsOnboarding(user: User | null): boolean {
  if (!user) return false;
  return !user.goal || !user.region || !user.budgetTier || !user.weightKg;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  token: null,
  user: null,
  hydrated: false,

  hydrate: async () => {
    try {
      const [token, rawUser] = await Promise.all([
        SecureStore.getItemAsync(TOKEN_KEY),
        AsyncStorage.getItem(USER_KEY),
      ]);
      set({
        token: token ?? null,
        user: rawUser ? (JSON.parse(rawUser) as User) : null,
        hydrated: true,
      });
    } catch {
      set({ hydrated: true });
    }
  },

  signIn: async (token, user) => {
    await Promise.all([
      SecureStore.setItemAsync(TOKEN_KEY, token),
      AsyncStorage.setItem(USER_KEY, JSON.stringify(user)),
    ]);
    set({ token, user });
  },

  setUser: async (user) => {
    await AsyncStorage.setItem(USER_KEY, JSON.stringify(user));
    set({ user });
  },

  signOut: async () => {
    await Promise.all([
      SecureStore.deleteItemAsync(TOKEN_KEY),
      AsyncStorage.removeItem(USER_KEY),
    ]);
    set({ token: null, user: null });
  },
}));

/** Non-hook accessor for the axios interceptor. */
export const getToken = () => useAuthStore.getState().token;
export const forceSignOut = () => useAuthStore.getState().signOut();
