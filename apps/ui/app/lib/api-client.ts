import type { ApiResponse } from "~/types";

/**
 * API configuration and base URL
 */
const API_URL = process.env.API_URL || "http://localhost:8080/api";

/**
 * API client class for making HTTP requests to the ELO backend
 */
export class ApiClient {
  private baseUrl: string;
  private token: string | null = null;

  constructor(baseUrl: string = API_URL) {
    this.baseUrl = baseUrl;

    // Try to get token from localStorage on client side
    if (typeof window !== "undefined") {
      this.token = localStorage.getItem("auth-token");
    }
  }

  /**
   * Set authentication token for subsequent requests
   */
  setToken(token: string | null): void {
    this.token = token;
    if (typeof window !== "undefined") {
      if (token) {
        localStorage.setItem("auth-token", token);
      } else {
        localStorage.removeItem("auth-token");
      }
    }
  }

  /**
   * Get current authentication token
   */
  getToken(): string | null {
    return this.token;
  }

  /**
   * Make a GET request to the API
   */
  async get<T>(
    endpoint: string,
    params?: Record<string, any>
  ): Promise<ApiResponse<T>> {
    const url = new URL(`${this.baseUrl}${endpoint}`);

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          url.searchParams.append(key, String(value));
        }
      });
    }

    return this.request<T>(url.toString(), {
      method: "GET",
    });
  }

  /**
   * Make a POST request to the API
   */
  async post<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>(`${this.baseUrl}${endpoint}`, {
      method: "POST",
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  /**
   * Make a PUT request to the API
   */
  async put<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>(`${this.baseUrl}${endpoint}`, {
      method: "PUT",
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  /**
   * Make a DELETE request to the API
   */
  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(`${this.baseUrl}${endpoint}`, {
      method: "DELETE",
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
    formData.append("file", file);

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
          reject(new Error("Invalid JSON response"));
        }
      };

      xhr.onerror = () => reject(new Error("Upload failed"));

      xhr.open("POST", `${this.baseUrl}${endpoint}`);

      if (this.token) {
        xhr.setRequestHeader("Authorization", `Bearer ${this.token}`);
      }

      xhr.send(formData);
    });
  }

  /**
   * Internal method for making HTTP requests
   */
  private async request<T>(
    url: string,
    options: RequestInit
  ): Promise<ApiResponse<T>> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...((options.headers as Record<string, string>) || {}),
    };

    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      const data = await response.json();

      if (!response.ok) {
        return {
          success: false,
          message:
            data.message || `HTTP ${response.status}: ${response.statusText}`,
          errors: data.errors,
        };
      }

      return {
        success: true,
        data: data.data || data,
        message: data.message,
      };
    } catch (error) {
      console.error("API request failed:", error);
      return {
        success: false,
        message:
          error instanceof Error ? error.message : "Onbekende fout opgetreden",
      };
    }
  }
}

// Global API client instance
export const apiClient = new ApiClient();
