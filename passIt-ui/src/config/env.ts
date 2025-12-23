/**
 * Environment configuration
 * All environment variables must be prefixed with VITE_ to be exposed to the client
 */

export const config = {
  // Backend API URL
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  
  // App configuration
  appPort: import.meta.env.VITE_APP_PORT || 3000,
  
  // Helper to check if running in production
  isProd: import.meta.env.PROD,
  isDev: import.meta.env.DEV,
} as const;

// Validate required environment variables
if (!config.apiBaseUrl) {
  console.error('Missing required environment variable: VITE_API_BASE_URL');
}

export default config;
