/**
 * User role enumeration defining access levels in the ELO system
 */
export type UserRole = 'student' | 'teacher' | 'admin';

/**
 * User interface representing system users
 */
export interface User {
  id: string;
  username: string;
  email: string;
  role: UserRole;
}

/**
 * Authentication credentials for login
 */
export interface LoginCredentials {
  email: string;
  password: string;
}

/**
 * Registration data with validation requirements
 */
export interface RegisterData {
  username: string;
  password: string;
  email: string;
}

/**
 * Authentication response from the API
 */
export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}

/**
 * Refresh token request
 */
export interface RefreshTokenRequest {
  refresh_token: string;
}

/**
 * Logout response
 */
export interface LogoutResponse {
  message: string;
}

/**
 * Current authentication state
 */
export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error?: string;
}
