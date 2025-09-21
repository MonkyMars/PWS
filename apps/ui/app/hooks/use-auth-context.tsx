import React, { createContext, useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { useQueryClient } from "@tanstack/react-query";
import { useCurrentUser } from "./use-auth";
import type { User, AuthState } from "~/types";

interface AuthContextType extends AuthState {
  login: () => void;
  logout: () => void;
  refreshAuth: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: React.ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { data: user, isLoading, error } = useCurrentUser();
  const [authError, setAuthError] = useState<string | undefined>();

  // Handle authentication failures from API client
  useEffect(() => {
    const handleAuthFailure = () => {
      // Clear all auth data
      queryClient.removeQueries({ queryKey: ["auth"] });
      setAuthError("Your session has expired. Please log in again.");

      // Redirect to login page
      navigate("/login", { replace: true });
    };

    // Listen for auth failure events from API client
    window.addEventListener("auth:failure", handleAuthFailure);

    return () => {
      window.removeEventListener("auth:failure", handleAuthFailure);
    };
  }, [navigate, queryClient]);

  // Clear auth error when user changes (successful login)
  useEffect(() => {
    if (user) {
      setAuthError(undefined);
    }
  }, [user]);

  const login = () => {
    // This is called after successful login to trigger auth state update
    queryClient.invalidateQueries({ queryKey: ["auth"] });
    setAuthError(undefined);
  };

  const logout = () => {
    // Clear all auth data
    queryClient.removeQueries({ queryKey: ["auth"] });
    setAuthError(undefined);
    navigate("/login", { replace: true });
  };

  const refreshAuth = () => {
    // Manually refresh authentication state
    queryClient.invalidateQueries({ queryKey: ["auth", "user"] });
  };

  const value: AuthContextType = {
    user: user || null,
    isAuthenticated: !!user && !isLoading,
    isLoading,
    error: authError || (error as any)?.message,
    login,
    logout,
    refreshAuth,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}

export function useRequireAuth(): AuthContextType {
  const auth = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!auth.isLoading && !auth.isAuthenticated) {
      navigate("/login", { replace: true });
    }
  }, [auth.isAuthenticated, auth.isLoading, navigate]);

  return auth;
}
