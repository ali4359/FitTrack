import Constants from 'expo-constants';

const extra = (Constants.expoConfig?.extra ?? {}) as { apiBaseUrl?: string };

/** Backend root, e.g. http://localhost:8080 (no trailing slash). */
export const API_BASE_URL = (extra.apiBaseUrl ?? 'http://localhost:8080').replace(/\/$/, '');

/** Full API prefix the client talks to. */
export const API_URL = `${API_BASE_URL}/api`;
