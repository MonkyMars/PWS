/**
 * Environment configuration for feature flags and application settings.
 *
 * This module provides type-safe access to environment variables and feature flags.
 * Variables prefixed with VITE_ are available in the browser.
 */

export interface FeatureFlags {
  /** Enable/disable login functionality */
  enableLogin: boolean;
  /** Enable/disable registration functionality */
  enableRegister: boolean;
}

/**
 * Parse a boolean environment variable with a default value.
 * Supports various truthy values: 'true', '1', 'yes', 'on', 'enabled'
 */
function parseBooleanEnv(value: string | undefined, defaultValue: boolean): boolean {
  if (value === undefined || value === '') {
    return defaultValue;
  }

  const normalizedValue = value.toLowerCase().trim();
  const truthyValues = ['true', '1', 'yes', 'on', 'enabled'];
  const falsyValues = ['false', '0', 'no', 'off', 'disabled'];

  if (truthyValues.includes(normalizedValue)) {
    return true;
  }

  if (falsyValues.includes(normalizedValue)) {
    return false;
  }

  return defaultValue;
}

/**
 * Get the current feature flags configuration.
 * When environment variables are not provided, both login and register are enabled by default.
 */
export function getFeatureFlags(): FeatureFlags {
  return {
    enableLogin: parseBooleanEnv(import.meta.env.VITE_ENABLE_LOGIN, true),
    enableRegister: parseBooleanEnv(import.meta.env.VITE_ENABLE_REGISTER, true),
  };
}

/**
 * Environment configuration object.
 */
export const env = {
  /** API base URL */
  apiUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',

  /** Node environment */
  nodeEnv: import.meta.env.NODE_ENV || 'development',

  /** Feature flags */
  features: getFeatureFlags(),
} as const;

export type Env = typeof env;
