/**
 * Authentication utility functions for token management and API calls
 */

import type { AuthResponse, RefreshTokenRequest, LogoutResponse, LoginCredentials, RegisterData } from '../types/auth';

const API_BASE_URL = 'http://localhost:8082';

/**
 * Token storage keys
 */
const TOKEN_KEYS = {
  ACCESS_TOKEN: 'access_token',
  REFRESH_TOKEN: 'refresh_token',
} as const;

/**
 * Store tokens in localStorage
 */
export function storeTokens(accessToken: string, refreshToken: string): void {
  localStorage.setItem(TOKEN_KEYS.ACCESS_TOKEN, accessToken);
  localStorage.setItem(TOKEN_KEYS.REFRESH_TOKEN, refreshToken);
}

/**
 * Get access token from localStorage
 */
export function getAccessToken(): string | null {
  return localStorage.getItem(TOKEN_KEYS.ACCESS_TOKEN);
}

/**
 * Get refresh token from localStorage
 */
export function getRefreshToken(): string | null {
  return localStorage.getItem(TOKEN_KEYS.REFRESH_TOKEN);
}

/**
 * Clear all tokens from localStorage
 */
export function clearTokens(): void {
  localStorage.removeItem(TOKEN_KEYS.ACCESS_TOKEN);
  localStorage.removeItem(TOKEN_KEYS.REFRESH_TOKEN);
}

/**
 * Check if user is authenticated (has valid access token)
 */
export function isAuthenticated(): boolean {
  return !!getAccessToken();
}

/**
 * Create Authorization header for API requests
 */
export function getAuthHeader(): Record<string, string> {
  const token = getAccessToken();
  return token ? { Authorization: `Bearer ${token}` } : {};
}

/**
 * Login user with email and password
 */
export async function login(credentials: LoginCredentials): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Login failed');
  }

  const authResponse: AuthResponse = await response.json();

  // Store tokens
  storeTokens(authResponse.access_token, authResponse.refresh_token);

  return authResponse;
}

/**
 * Register new user
 */
export async function register(userData: RegisterData): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Registration failed');
  }

  const authResponse: AuthResponse = await response.json();

  // Store tokens
  storeTokens(authResponse.access_token, authResponse.refresh_token);

  return authResponse;
}

/**
 * Refresh access token using refresh token
 */
export async function refreshAccessToken(): Promise<AuthResponse> {
  const refreshToken = getRefreshToken();

  if (!refreshToken) {
    throw new Error('No refresh token available');
  }

  const refreshRequest: RefreshTokenRequest = {
    refresh_token: refreshToken,
  };

  const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(refreshRequest),
  });

  if (!response.ok) {
    // If refresh fails, clear tokens and redirect to login
    clearTokens();
    throw new Error('Token refresh failed');
  }

  const authResponse: AuthResponse = await response.json();

  // Store new tokens
  storeTokens(authResponse.access_token, authResponse.refresh_token);

  return authResponse;
}

/**
 * Get current user information
 */
export async function getCurrentUser() {
  const response = await fetch(`${API_BASE_URL}/auth/me`, {
    method: 'GET',
    headers: {
      ...getAuthHeader(),
    },
  });

  if (!response.ok) {
    if (response.status === 401) {
      // Try to refresh token
      try {
        await refreshAccessToken();
        // Retry the request with new token
        return getCurrentUser();
      } catch {
        clearTokens();
        throw new Error('Authentication failed');
      }
    }
    throw new Error('Failed to get user information');
  }

  return response.json();
}

/**
 * Logout user
 */
export async function logout(): Promise<LogoutResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/logout`, {
    method: 'POST',
    headers: {
      ...getAuthHeader(),
    },
  });

  // Clear tokens regardless of response status
  clearTokens();

  if (!response.ok) {
    throw new Error('Logout request failed');
  }

  return response.json();
}

/**
 * Make authenticated API request with automatic token refresh
 */
export async function authenticatedRequest(
  url: string,
  options: RequestInit = {}
): Promise<Response> {
  const requestOptions: RequestInit = {
    ...options,
    headers: {
      ...options.headers,
      ...getAuthHeader(),
    },
  };

  let response = await fetch(url, requestOptions);

  // If unauthorized, try to refresh token and retry
  if (response.status === 401) {
    try {
      await refreshAccessToken();

      // Retry with new token
      requestOptions.headers = {
        ...options.headers,
        ...getAuthHeader(),
      };

      response = await fetch(url, requestOptions);
    } catch {
      clearTokens();
      throw new Error('Authentication failed');
    }
  }

  return response;
}

/**
 * Decode JWT token payload (client-side only for display purposes)
 * Note: Never trust client-side JWT decoding for security decisions
 */
export function decodeJWT(token: string): any {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch {
    return null;
  }
}

/**
 * Check if token is expired (client-side check only)
 */
export function isTokenExpired(token: string): boolean {
  const decoded = decodeJWT(token);
  if (!decoded || !decoded.exp) return true;

  const currentTime = Math.floor(Date.now() / 1000);
  return decoded.exp < currentTime;
}

/**
 * Check if access token needs refresh (expires in less than 5 minutes)
 */
export function needsRefresh(): boolean {
  const token = getAccessToken();
  if (!token) return false;

  const decoded = decodeJWT(token);
  if (!decoded || !decoded.exp) return true;

  const currentTime = Math.floor(Date.now() / 1000);
  const fiveMinutes = 5 * 60;

  return decoded.exp - currentTime < fiveMinutes;
}
