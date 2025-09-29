import type { ApiResponse } from '~/types';
import { env } from './env';

/**
 * API configuration and base URL
 */
const API_URL = env.apiUrl;

/**
 * API client class for making HTTP requests to the ELO backend with cookie-based auth
 */
export class ApiClient {
  private baseUrl: string;
  private isRefreshing = false;
  private refreshPromise: Promise<boolean> | null = null;

  constructor(baseUrl: string = API_URL) {
    this.baseUrl = baseUrl;
  }

  /**
   * Refresh the access token using the refresh token from cookies
   */
  private async refreshToken(): Promise<boolean> {
    try {
      const response = await fetch(`${this.baseUrl}/auth/refresh`, {
        method: 'POST',
        credentials: 'include', // Include cookies
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        // New tokens are automatically set as cookies by the server
        return true;
      }

      return false;
    } catch (error) {
      if (import.meta.env.DEV) {
        console.error('Token refresh failed:', error);
      }
      return false;
    }
  }

  /**
   * Handle token refresh for 401 responses
   */
  private async handleTokenRefresh(): Promise<boolean> {
    if (this.isRefreshing) {
      // If already refreshing, wait for the existing refresh
      return this.refreshPromise || Promise.resolve(false);
    }

    this.isRefreshing = true;
    this.refreshPromise = this.refreshToken();

    try {
      const success = await this.refreshPromise;
      return success;
    } finally {
      this.isRefreshing = false;
      this.refreshPromise = null;
    }
  }

  /**
   * Internal method for making HTTP requests with automatic token refresh
   */
  private async request<T>(
    url: string,
    options: RequestInit,
    retryOnAuth = true
  ): Promise<ApiResponse<T>> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...((options.headers as Record<string, string>) || {}),
    };

    const requestOptions: RequestInit = {
      ...options,
      headers,
      credentials: 'include', // Always include cookies
    };

    try {
      const response = await fetch(url, requestOptions);
      const data = await response.json();

      // Handle 401 Unauthorized responses
      if (response.status === 401 && retryOnAuth) {
        // Try to refresh tokens since we can't check HttpOnly cookies
        const refreshSuccess = await this.handleTokenRefresh();

        if (refreshSuccess) {
          // Retry the original request once
          return this.request<T>(url, options, false);
        }

        // Only trigger auth failure for non-auth endpoints
        // /auth/me returning 401 is expected for unauthenticated users
        if (!url.includes('/auth/me')) {
          this.handleAuthFailure();
        }

        return {
          success: false,
          message: 'Authentication failed',
        };
      }

      if (!response.ok) {
        return {
          success: false,
          message: data.message || `HTTP ${response.status}: ${response.statusText}`,
          errors: data.errors,
        };
      }

      return {
        success: true,
        data: data.data || data,
        message: data.message,
      };
    } catch (error) {
      if (import.meta.env.DEV) {
        console.error('API request failed:', error);
      }
      return {
        success: false,
        message: error instanceof Error ? error.message : 'Onbekende fout opgetreden',
      };
    }
  }

  /**
   * Handle authentication failure (redirect to login or emit event)
   */
  private handleAuthFailure(): void {
    // Emit custom event for auth failure
    window.dispatchEvent(new CustomEvent('auth:failure'));
  }

  /**
   * Make a GET request to the API
   */
  async get<T>(endpoint: string, params?: Record<string, any>): Promise<ApiResponse<T>> {
    const url = new URL(`${this.baseUrl}${endpoint}`);

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          url.searchParams.append(key, String(value));
        }
      });
    }

    return this.request<T>(url.toString(), {
      method: 'GET',
    });
  }

  /**
   * Make a POST request to the API
   */
  async post<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>(`${this.baseUrl}${endpoint}`, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  /**
   * Make a PUT request to the API
   */
  async put<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>(`${this.baseUrl}${endpoint}`, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  /**
   * Make a DELETE request to the API
   */
  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(`${this.baseUrl}${endpoint}`, {
      method: 'DELETE',
    });
  }

  /**
   * Upload a file to the API
   */
  async uploadFile<T>(
    endpoint: string,
    file: File,
    additionalData?: Record<string, any>,
    onProgress?: (progress: number) => void
  ): Promise<ApiResponse<T>> {
    const formData = new FormData();
    formData.append('file', file);

    if (additionalData) {
      Object.entries(additionalData).forEach(([key, value]) => {
        formData.append(key, String(value));
      });
    }

    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();

      if (onProgress) {
        xhr.upload.onprogress = (event) => {
          if (event.lengthComputable) {
            const progress = (event.loaded / event.total) * 100;
            onProgress(progress);
          }
        };
      }

      xhr.onload = () => {
        try {
          const response = JSON.parse(xhr.responseText);
          resolve(response);
        } catch (error) {
          reject(new Error('Invalid JSON response'));
        }
      };

      xhr.onerror = () => reject(new Error('Upload failed'));

      xhr.open('POST', `${this.baseUrl}${endpoint}`);
      xhr.withCredentials = true; // Include cookies for uploads
      xhr.send(formData);
    });
  }

  /**
   * Check if user is authenticated by making a request to /auth/me
   */
  async checkAuth(): Promise<boolean> {
    try {
      const response = await this.get('/auth/me');
      return response.success;
    } catch (error) {
      return false;
    }
  }

  /**
   * Check if authentication cookies are present (for HttpOnly cookies, this always returns false)
   * This method is kept for compatibility but will always return false for HttpOnly cookies
   */
  hasAuthCookies(): boolean {
    if (typeof document === 'undefined') return false;

    // Note: HttpOnly cookies cannot be read by JavaScript
    // This method will always return false in production where cookies are HttpOnly
    // We rely on server-side authentication validation instead
    return false;
  }

  /**
   * Logout user by calling the logout endpoint
   */
  async logout(): Promise<boolean> {
    try {
      const response = await this.post('/auth/logout');
      return response.success;
    } catch {
      return false;
    }
  }
}

// Global API client instance
export const apiClient = new ApiClient();
