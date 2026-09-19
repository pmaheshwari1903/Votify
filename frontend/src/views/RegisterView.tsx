import React, { useState } from 'react';
import { Button } from '../components/ui/Button';
import { apiClient } from '../lib/api';
import type { AuthResponse } from '../types';

interface RegisterViewProps {
  onSuccess: (user: AuthResponse['user'], token: string) => void;
  onNavigateLogin: () => void;
}

export const RegisterView: React.FC<RegisterViewProps> = ({ onSuccess, onNavigateLogin }) => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const res = await apiClient.post<AuthResponse>('/auth/register', { name, email, password });
      onSuccess(res.user, res.token);
    } catch (err: any) {
      setError(err.message || 'Registration failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '440px', margin: '60px auto', padding: '0 20px' }}>
      <div className="card" style={{ padding: '32px' }}>
        <div style={{ marginBottom: '24px', textAlign: 'center' }}>
          <span className="brand-badge">NEW ACCOUNT</span>
          <h2 style={{ fontFamily: 'var(--v-font-heading)', fontSize: '28px', color: 'var(--v-color-deep-ink)', marginTop: '8px' }}>
            Create Account
          </h2>
          <p style={{ color: 'var(--v-color-warm-gray)', fontSize: '14px', marginTop: '4px' }}>
            Join Votify for simple, elegant live polling.
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
            <label style={{ display: 'block', fontSize: '13px', fontWeight: 600, marginBottom: '6px' }}>Full Name</label>
            <input
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Sora Takahashi"
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
            <label style={{ display: 'block', fontSize: '13px', fontWeight: 600, marginBottom: '6px' }}>Password (min 8 chars)</label>
            <input
              type="password"
              required
              minLength={8}
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
            {loading ? 'Creating Account...' : 'Register'}
          </Button>
        </form>

        <div style={{ marginTop: '24px', textAlign: 'center', fontSize: '13px', color: 'var(--v-color-warm-gray)' }}>
          Already have an account?{' '}
          <button
            onClick={onNavigateLogin}
            style={{ background: 'none', border: 'none', color: 'var(--v-color-vermilion)', fontWeight: 600, cursor: 'pointer' }}
          >
            Sign In
          </button>
        </div>
      </div>
    </div>
  );
};
