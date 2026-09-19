// Votify API Client Base
// Centralized fetch wrapper with HttpOnly cookie authentication.

import { ENV } from '../config/env';
import type { ApiResponse } from '../types';

export class ApiError extends Error {
  code: string;
  details?: Record<string, string>;

  constructor(message: string, code: string = 'UNKNOWN_ERROR', details?: Record<string, string>) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.details = details;
  }
}

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${ENV.API_BASE_URL}${endpoint}`;
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  const response = await fetch(url, {
    ...options,
    headers,
    credentials: 'include', // Send HttpOnly cookies automatically
  });

  const body: ApiResponse<T> = await response.json().catch(() => ({
    success: false,
    data: null as unknown as T,
    error: {
      code: 'INVALID_JSON',
      message: 'Server returned a non-JSON response.',
    },
  }));

  if (!response.ok || !body.success) {
    throw new ApiError(
      body.error?.message || `Request failed with status ${response.status}`,
      body.error?.code || `HTTP_${response.status}`,
      body.error?.details
    );
  }

  return body.data;
}

export const apiClient = {
  get: <T>(endpoint: string, options?: RequestInit) =>
    request<T>(endpoint, { ...options, method: 'GET' }),

  post: <T>(endpoint: string, payload?: unknown, options?: RequestInit) =>
    request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: payload ? JSON.stringify(payload) : undefined,
    }),

  put: <T>(endpoint: string, payload?: unknown, options?: RequestInit) =>
    request<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: payload ? JSON.stringify(payload) : undefined,
    }),

  delete: <T>(endpoint: string, options?: RequestInit) =>
    request<T>(endpoint, { ...options, method: 'DELETE' }),
};
