import axios, { AxiosError, AxiosHeaders } from 'axios';
import type { ApiError } from '@iron-and-spice/shared';
import { API_URL } from '../config';
import { forceSignOut, getToken } from '../store/authStore';

export const api = axios.create({
  baseURL: API_URL,
  timeout: 15000,
});

// Attach the JWT to every request.
api.interceptors.request.use((config) => {
  const token = getToken();
  if (token) {
    const headers = AxiosHeaders.from(config.headers);
    headers.set('Authorization', `Bearer ${token}`);
    config.headers = headers;
  }
  return config;
});

// On 401, drop the session — RootNavigator swaps to the auth stack reactively.
api.interceptors.response.use(
  (res) => res,
  (error: AxiosError<ApiError>) => {
    if (error.response?.status === 401) {
      void forceSignOut();
    }
    return Promise.reject(error);
  },
);

/** Pull a human-readable message out of an axios error. */
export function apiErrorMessage(error: unknown, fallback = 'Something went wrong'): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as ApiError | undefined;
    if (data?.error) return data.error;
    if (error.code === 'ECONNABORTED') return 'Request timed out';
    if (!error.response) return 'Cannot reach the server';
  }
  return fallback;
}
