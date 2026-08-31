import { useMutation } from '@tanstack/react-query';
import type { LoginRequest, RegisterRequest, UpdateProfileRequest } from '@fittrack/shared';
import { login, register, updateProfile } from '../api/endpoints';
import { queryClient, queryKeys } from '../api/queryClient';
import { useAuthStore } from '../store/authStore';

export function useLogin() {
  const signIn = useAuthStore((s) => s.signIn);
  return useMutation({
    mutationFn: (body: LoginRequest) => login(body),
    onSuccess: async ({ token, user }) => {
      await signIn(token, user);
    },
  });
}

export function useRegister() {
  const signIn = useAuthStore((s) => s.signIn);
  return useMutation({
    mutationFn: (body: RegisterRequest) => register(body),
    onSuccess: async ({ token, user }) => {
      await signIn(token, user);
    },
  });
}

export function useUpdateProfile() {
  const setUser = useAuthStore((s) => s.setUser);
  return useMutation({
    mutationFn: (body: UpdateProfileRequest) => updateProfile(body),
    onSuccess: async (user) => {
      await setUser(user);
      queryClient.setQueryData(queryKeys.profile, user);
    },
  });
}

export function useLogout() {
  const signOut = useAuthStore((s) => s.signOut);
  return async () => {
    await signOut();
    queryClient.clear();
  };
}
