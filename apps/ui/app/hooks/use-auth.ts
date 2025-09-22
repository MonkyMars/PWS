import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '~/lib/api-client';
import type { User, LoginCredentials, RegisterData, AuthResponse } from '~/types';

/**
 * Hook for user login
 */
export function useLogin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (credentials: LoginCredentials): Promise<AuthResponse> => {
      const response = await apiClient.post<AuthResponse>('/auth/login', credentials);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Login mislukt');
      }

      // Set token in API client
      apiClient.setToken(response.data.token);

      return response.data;
    },
    onSuccess: (data) => {
      // Cache user data
      queryClient.setQueryData(['auth', 'user'], data.user);
      queryClient.setQueryData(['auth', 'token'], data.token);
    },
    onError: () => {
      // Clear any existing auth data on login failure
      queryClient.removeQueries({ queryKey: ['auth'] });
      apiClient.setToken(null);
    },
  });
}

/**
 * Hook for user registration
 */
export function useRegister() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (userData: RegisterData): Promise<AuthResponse> => {
      const response = await apiClient.post<AuthResponse>('/auth/register', userData);

      if (!response.success || !response.data) {
        throw new Error(response.message || 'Registratie mislukt');
      }

      // Set token in API client
      apiClient.setToken(response.data.token);

      return response.data;
    },
    onSuccess: (data) => {
      // Cache user data
      queryClient.setQueryData(['auth', 'user'], data.user);
      queryClient.setQueryData(['auth', 'token'], data.token);
    },
  });
}

/**
 * Hook for user logout
 */
export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (): Promise<void> => {
      const response = await apiClient.post('/auth/logout');

      if (!response.success) {
        throw new Error(response.message || 'Uitloggen mislukt');
      }
    },
    onSuccess: () => {
      // Clear all cached data
      queryClient.clear();
      apiClient.setToken(null);
    },
    onError: () => {
      // Even if logout API fails, clear local data
      queryClient.clear();
      apiClient.setToken(null);
    },
  });
}

/**
 * Hook to get current user data
 */
export function useCurrentUser() {
  return useQuery({
    queryKey: ['auth', 'user'],
    queryFn: async (): Promise<User | null> => {
      const token = apiClient.getToken();

      if (!token) {
        return null;
      }

      const response = await apiClient.get<User>('/auth/me');

      if (!response.success || !response.data) {
        // Invalid token, clear auth data
        apiClient.setToken(null);
        return null;
      }

      return response.data;
    },
    retry: false,
    staleTime: 1000 * 60 * 10, // 10 minutes
  });
}

/**
 * Hook to check if user is authenticated
 */
export function useIsAuthenticated(): boolean {
  const { data: user } = useCurrentUser();
  return !!user;
}
