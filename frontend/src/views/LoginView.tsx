import React from 'react';
import { Button } from '../components/ui/Button';
import { ENV } from '../config/env';

export const LoginView: React.FC = () => {
  const handleSignIn = () => {
    // Navigate to OAuth Login endpoint on Votify Backend
    const loginUrl = `${ENV.API_BASE_URL}/auth/oauth/login`;
    window.location.href = loginUrl;
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

        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', marginTop: '32px' }}>
          <Button
            type="button"
            variant="primary"
            size="lg"
            fullWidth
            onClick={handleSignIn}
          >
            Sign in with Maheshwari.com
          </Button>
        </div>
      </div>
    </div>
  );
};
