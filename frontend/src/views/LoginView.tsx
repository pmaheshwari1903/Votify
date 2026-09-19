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
    <div style={{ maxWidth: '460px', margin: '60px auto', padding: '0 20px' }}>
      <div className="card" style={{ padding: '36px 32px' }}>
        <div style={{ marginBottom: '28px', textAlign: 'center' }}>
          <span className="brand-badge" style={{ backgroundColor: 'rgba(196, 52, 45, 0.1)', color: 'var(--v-color-vermilion)' }}>
            SINGLE SIGN-ON (OIDC)
          </span>
          <h2 style={{ fontFamily: 'var(--v-font-heading)', fontSize: '28px', color: 'var(--v-color-deep-ink)', marginTop: '10px' }}>
            Personal OAuth 2.0 Login
          </h2>
          <p style={{ color: 'var(--v-color-warm-gray)', fontSize: '14px', marginTop: '6px', lineHeight: 1.5 }}>
            Votify uses unified single user authentication via your personal OIDC provider at <strong>Maheshwari.com</strong>.
          </p>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', marginTop: '24px' }}>
          <Button
            type="button"
            variant="primary"
            size="lg"
            fullWidth
            onClick={handleSignIn}
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              gap: '10px',
              padding: '14px',
              fontSize: '15px',
              fontWeight: 600,
            }}
          >
            <span>🔐</span>
            <span>Sign in with Maheshwari.com</span>
          </Button>

          <p style={{ textAlign: 'center', fontSize: '12px', color: 'var(--v-color-warm-gray)', marginTop: '8px' }}>
            Protected by OpenID Connect (OIDC) + PKCE standard
          </p>
        </div>
      </div>
    </div>
  );
};
