import React, { useState } from 'react';
import { Button } from '../components/ui/Button';
import { apiClient } from '../lib/api';
import type { AuthResponse } from '../types';

interface LoginViewProps {
  onSuccess: (user: AuthResponse['user'], token: string) => void;
  onNavigateRegister: () => void;
}

export const LoginView: React.FC<LoginViewProps> = ({ onSuccess, onNavigateRegister }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res = await apiClient.post<AuthResponse>('/auth/login', { email, password });
      onSuccess(res.user, res.token);
    } catch (err: any) {
      setError(err.message || 'Login failed. Please check your credentials.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '440px', margin: '60px auto', padding: '0 20px' }}>
      <div className="card" style={{ padding: '32px' }}>
        <div style={{ marginBottom: '24px', textAlign: 'center' }}>
          <span className="brand-badge">AUTHENTICATION</span>
          <h2 style={{ fontFamily: 'var(--v-font-heading)', fontSize: '28px', color: 'var(--v-color-deep-ink)', marginTop: '8px' }}>
            Welcome Back
          </h2>
          <p style={{ color: 'var(--v-color-warm-gray)', fontSize: '14px', marginTop: '4px' }}>
            Sign in to manage your real-time polls.
          </p>
        </div>

        {error && (
          <div style={{
            padding: '12px 16px',
            backgroundColor: 'rgba(196, 52, 45, 0.08)',
            borderLeft: '3px solid var(--v-color-error)',
            borderRadius: 'var(--v-radius-sm)',
            color: 'var(--v-color-error)',
            fontSize: '13px',
            marginBottom: '20px'
          }}>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <div>
            <label style={{ display: 'block', fontSize: '13px', fontWeight: 600, marginBottom: '6px' }}>Email Address</label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              style={{
                width: '100%',
                padding: '10px 14px',
                borderRadius: 'var(--v-radius-md)',
                border: '1px solid var(--v-color-stone)',
                fontFamily: 'var(--v-font-body)',
                fontSize: '14px',
                outline: 'none',
              }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: '13px', fontWeight: 600, marginBottom: '6px' }}>Password</label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              style={{
                width: '100%',
                padding: '10px 14px',
                borderRadius: 'var(--v-radius-md)',
                border: '1px solid var(--v-color-stone)',
                fontFamily: 'var(--v-font-body)',
                fontSize: '14px',
                outline: 'none',
              }}
            />
          </div>

          <Button type="submit" variant="primary" size="lg" fullWidth disabled={loading}>
            {loading ? 'Signing in...' : 'Sign In'}
          </Button>
        </form>

        <div style={{ marginTop: '24px', textAlign: 'center', fontSize: '13px', color: 'var(--v-color-warm-gray)' }}>
          Don't have an account?{' '}
          <button
            onClick={onNavigateRegister}
            style={{ background: 'none', border: 'none', color: 'var(--v-color-vermilion)', fontWeight: 600, cursor: 'pointer' }}
          >
            Create an Account
          </button>
        </div>
      </div>
    </div>
  );
};
