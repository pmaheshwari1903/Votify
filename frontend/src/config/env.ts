// Votify Frontend Environment Configuration
// Encapsulates configuration values loaded from environment variables with sensible defaults.

export const ENV = {
  // API Gateway base URL
  API_BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',

  // Realtime Service WebSocket URL
  WS_URL: import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws',

  // Current environment stage
  APP_ENV: import.meta.env.VITE_APP_ENV || 'development',

  // App version
  APP_VERSION: '0.1.0-foundation',
} as const;
