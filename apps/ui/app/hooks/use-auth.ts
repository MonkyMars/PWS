import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "~/lib/api-client";
import type { User, LoginCredentials, RegisterData } from "~/types";

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
      // Cache user data - tokens are now handled via cookies
      queryClient.setQueryData(["auth", "user"], user);
      queryClient.invalidateQueries({ queryKey: ["auth"] });
      // Clear logged out flag on successful login
      if (typeof sessionStorage !== "undefined") {
        sessionStorage.removeItem("auth_logged_out");
      }
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
      // Clear all cached data
      queryClient.clear();
      // Mark as logged out to prevent further auth requests
      if (typeof sessionStorage !== "undefined") {
        sessionStorage.setItem("auth_logged_out", "true");
      }
    },
    onError: () => {
      // Even if logout API fails, clear local data
      queryClient.clear();
    },
  });
}

/**
 * Hook to get current user data
 */
export function useCurrentUser() {
  return useQuery({
    queryKey: ["auth", "user"],
    queryFn: async (): Promise<User | null> => {
      const response = await apiClient.get<User>("/auth/me");

      if (!response.success || !response.data) {
        return null;
      }

      return response.data;
    },
    enabled: () => {
      // SSR check - don't run on server
      if (typeof document === "undefined") return false;

      // Check if we've explicitly logged out
      const hasLoggedOut = sessionStorage.getItem("auth_logged_out") === "true";
      if (hasLoggedOut) return false;

      // Always attempt auth check since HttpOnly cookies can't be read by JS
      // but only if we haven't explicitly logged out
      return true;
    },
    retry: (failureCount: number, error: any) => {
      // Don't retry on 401 errors (handled by API client)
      if (error?.message?.includes("Authentication failed")) {
        return false;
      }
      return failureCount < 3;
    },

    staleTime: 1000 * 60 * 5, // 5 minutes
    gcTime: 1000 * 60 * 10, // 10 minutes (formerly cacheTime)
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
