import React from 'react';
import { Button } from '../components/ui/Button';
import { ENV } from '../config/env';
import { apiClient } from '../lib/api';
import type { User } from '../types';

export const LoginView: React.FC = () => {
  const handleSignIn = () => {
    const loginUrl = `${ENV.API_BASE_URL}/auth/oauth/login?popup=true`;

    const width = 500;
    const height = 650;
    const left = window.screen.width / 2 - width / 2;
    const top = window.screen.height / 2 - height / 2;

    /*
     * Register the message listener BEFORE opening the popup.
     * This prevents a race condition where the popup finishes
     * authentication before the listener is attached.
     */
    const handleMessage = async (event: MessageEvent) => {
      if (event.data?.type !== 'OAUTH_SUCCESS') {
        return;
      }

      window.removeEventListener('message', handleMessage);

      if (event.data.token) {
        localStorage.setItem('votify_token', event.data.token);
      }

      try {
        const user = await apiClient.get<User>('/auth/me');
        if (user) {
          window.location.reload();
        }
      } catch (error) {
        console.error('OAuth succeeded but Votify session fetch failed:', error);
        window.location.reload();
      }
    };

    window.addEventListener('message', handleMessage);

    // Open OIDC login popup AFTER the listener is ready.
    const popup = window.open(
      loginUrl,
      'MaheshwariOIDCAuth',
      `width=${width},height=${height},top=${top},left=${left},scrollbars=yes,status=yes`
    );

    // Popup blocked → use normal navigation.
    if (!popup || popup.closed || typeof popup.closed === 'undefined') {
      window.removeEventListener('message', handleMessage);
      window.location.href = loginUrl;
      return;
    }

    /*
     * Fallback:
     * If the popup closes without sending the postMessage,
     * check the Votify session automatically.
     */
    const checkPopupClosed = window.setInterval(async () => {
      if (popup.closed) {
        window.clearInterval(checkPopupClosed);

        try {
          const user = await apiClient.get<User>('/auth/me');

          if (user) {
            window.removeEventListener('message', handleMessage);
            window.location.reload();
          }
        } catch (error) {
          console.error('Popup closed but authentication was not detected:', error);
        }
      }
    }, 1000);
  };

  return (
    <div style={{ maxWidth: '460px', margin: '60px auto', padding: '0 20px' }}>
      <div className="card" style={{ padding: '36px 32px' }}>
        <div style={{ marginBottom: '28px', textAlign: 'center' }}>
          <span
            className="brand-badge"
            style={{
              backgroundColor: 'rgba(196, 52, 45, 0.1)',
              color: 'var(--v-color-vermilion)',
            }}
          >
            SINGLE SIGN-ON (OIDC)
          </span>

          <h2
            style={{
              fontFamily: 'var(--v-font-heading)',
              fontSize: '28px',
              color: 'var(--v-color-deep-ink)',
              marginTop: '10px',
            }}
          >
            Personal OAuth 2.0 Login
          </h2>

          <p
            style={{
              color: 'var(--v-color-warm-gray)',
              fontSize: '14px',
              marginTop: '6px',
              lineHeight: 1.5,
            }}
          >
            Votify uses unified single user authentication via your personal
            OIDC provider at <strong>Maheshwari.com</strong>.
          </p>
        </div>

        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            gap: '16px',
            marginTop: '24px',
          }}
        >
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

          <p
            style={{
              textAlign: 'center',
              fontSize: '12px',
              color: 'var(--v-color-warm-gray)',
              marginTop: '8px',
            }}
          >
            Protected by OpenID Connect (OIDC) + PKCE standard
          </p>
        </div>
      </div>
    </div>
  );
};