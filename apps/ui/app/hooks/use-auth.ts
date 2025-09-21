import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "~/lib/api-client";
import type { User, LoginCredentials, RegisterData } from "~/types";

/**
 * Track authentication state across app lifecycle
 * This helps minimize unnecessary API calls for unauthenticated users
 */
let hasAttemptedAuth = false;
let lastAuthResult: boolean | null = null;

/**
 * Hook for user login
 */
export function useLogin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (credentials: LoginCredentials): Promise<User> => {
      const response = await apiClient.post<User>("/auth/login", credentials);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Login mislukt");
      }

      return response.data;
    },
    onSuccess: (user: User) => {
      // Reset auth tracking on successful login
      hasAttemptedAuth = true;
      lastAuthResult = true;
      // Cache user data - tokens are now handled via cookies
      queryClient.setQueryData(["auth", "user"], user);
      queryClient.invalidateQueries({ queryKey: ["auth"] });
    },
    onError: () => {
      // Clear any existing auth data on login failure
      queryClient.removeQueries({ queryKey: ["auth"] });
    },
  });
}

/**
 * Hook for user registration
 */
export function useRegister() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (userData: RegisterData): Promise<User> => {
      const response = await apiClient.post<User>("/auth/register", userData);

      if (!response.success || !response.data) {
        throw new Error(response.message || "Registratie mislukt");
      }

      return response.data;
    },
    onSuccess: (user: User) => {
      // Reset auth tracking on successful registration
      hasAttemptedAuth = true;
      lastAuthResult = true;
      // Cache user data - tokens are now handled via cookies
      queryClient.setQueryData(["auth", "user"], user);
      queryClient.invalidateQueries({ queryKey: ["auth"] });
    },
    onError: () => {
      // Clear any existing auth data on registration failure
      queryClient.removeQueries({ queryKey: ["auth"] });
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
      const success = await apiClient.logout();

      if (!success) {
        throw new Error("Uitloggen mislukt");
      }
    },
    onSuccess: () => {
      // Reset auth tracking on logout
      hasAttemptedAuth = true;
      lastAuthResult = false;
      // Clear all cached data
      queryClient.clear();
    },
    onError: () => {
      // Reset auth tracking even if logout fails
      hasAttemptedAuth = true;
      lastAuthResult = false;
      // Even if logout API fails, clear local data
      queryClient.clear();
    },
  });
}

/**
 * Hook to get current user data
 */
export function useCurrentUser() {
  const queryClient = useQueryClient();

  return useQuery({
    queryKey: ["auth", "user"],
    queryFn: async (): Promise<User | null> => {
      const response = await apiClient.get<User>("/auth/me");

      if (!response.success || !response.data) {
        // Track failed auth attempt
        hasAttemptedAuth = true;
        lastAuthResult = false;
        return null;
      }

      // Track successful auth
      hasAttemptedAuth = true;
      lastAuthResult = true;
      return response.data;
    },
    enabled: () => {
      // SSR check - don't run on server
      if (typeof document === "undefined") return false;

      // Always allow first attempt since we can't read HttpOnly cookies
      if (!hasAttemptedAuth) return true;

      // If last attempt was successful, allow refetch (user might still be authenticated)
      if (lastAuthResult === true) return true;

      // If we have cached user data, allow refetch to check if still valid
      const cachedUser = queryClient.getQueryData(["auth", "user"]);
      if (cachedUser) return true;

      // If last attempt failed and no cached user, don't spam the API
      return false;
    },
    retry: (failureCount: number, error: any) => {
      // Don't retry on 401 errors (handled by API client)
      if (error?.message?.includes("Authentication failed")) {
        hasAttemptedAuth = true;
        lastAuthResult = false;
        return false;
      }
      return failureCount < 3;
    },

    staleTime: 1000 * 60 * 5, // 5 minutes
    gcTime: 1000 * 60 * 10, // 10 minutes (formerly cacheTime)
    refetchOnWindowFocus: false, // Prevent refetch on window focus
    refetchOnMount: true, // Always refetch on mount to get fresh data
    // Don't refetch on reconnect for auth queries to avoid spam
    refetchOnReconnect: false,
  });
}

/**
 * Hook to check if user is authenticated
 */
export function useIsAuthenticated(): boolean {
  const { data: user, isLoading } = useCurrentUser();

  // Don't consider user authenticated while loading
  if (isLoading) return false;

  return !!user;
}

/**
 * Hook to check authentication status without triggering queries
 */
export function useAuthStatus() {
  const { data: user, isLoading, error } = useCurrentUser();

  return {
    user,
    isAuthenticated: !!user && !isLoading,
    isLoading,
    error,
  };
}

/**
 * Hook for token refresh (manual trigger if needed)
 */
export function useRefreshToken() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (): Promise<boolean> => {
      // The API client handles token refresh automatically,
      // but this can be used to manually trigger a refresh
      return await apiClient.checkAuth();
    },
    onSuccess: (success: boolean) => {
      if (success) {
        // Invalidate user query to refetch user data
        queryClient.invalidateQueries({ queryKey: ["auth", "user"] });
      } else {
        // Clear auth data if refresh failed
        queryClient.removeQueries({ queryKey: ["auth"] });
      }
    },
  });
}
