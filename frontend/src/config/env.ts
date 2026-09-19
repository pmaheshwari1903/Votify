// Votify Frontend Environment Configuration
// Encapsulates configuration values loaded from environment variables with sensible defaults.

const isLocalhost =
  typeof window !== 'undefined' &&
  (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1');

export const ENV = {
  // API Gateway base URL
  API_BASE_URL:
    import.meta.env.VITE_API_BASE_URL ||
    (isLocalhost
      ? 'http://localhost:8080/api/v1'
      : 'https://votify-production-2915.up.railway.app/api/v1'),

  // Realtime Service WebSocket URL
  WS_URL:
    import.meta.env.VITE_WS_URL ||
    (isLocalhost
      ? 'ws://localhost:8080/ws'
      : 'wss://votify-production-2915.up.railway.app/ws'),

  // Current environment stage
  APP_ENV: import.meta.env.VITE_APP_ENV || (isLocalhost ? 'development' : 'production'),

  // App version
  APP_VERSION: '0.1.0-foundation',
} as const;
