import React from 'react';
import { Button } from '../components/ui/Button';
import { APP_ROUTES } from './router';
import { ENV } from '../config/env';
import './App.css';

export const App: React.FC = () => {
  const services = [
    { name: 'API Gateway', port: '8080', route: '/api/v1/*' },
    { name: 'Auth Service', port: '8081', route: '/api/v1/auth/*' },
    { name: 'Poll Service', port: '8082', route: '/api/v1/polls/*' },
    { name: 'Vote Service', port: '8083', route: '/api/v1/votes/*' },
    { name: 'Realtime Service', port: '8084', route: '/ws/*' },
    { name: 'Analytics Service', port: '8085', route: '/api/v1/analytics/*' },
    { name: 'Payment Service', port: '8086', route: '/api/v1/payments/*' },
  ];

  const colorSwatches = [
    { name: 'Vermilion', code: '#D94032', textLight: true, var: '--v-color-vermilion' },
    { name: 'Ink Black', code: '#1A1A1A', textLight: true, var: '--v-color-ink-black' },
    { name: 'Deep Ink', code: '#0D0D0D', textLight: true, var: '--v-color-deep-ink' },
    { name: 'Washi Cream', code: '#FAF7F2', textLight: false, var: '--v-color-washi' },
    { name: 'Warm Gray', code: '#6B6460', textLight: true, var: '--v-color-warm-gray' },
    { name: 'Sakura', code: '#F0C4C4', textLight: false, var: '--v-color-sakura' },
  ];

  return (
    <div className="app-container">
      <header className="hero-header">
        <span className="brand-badge">VOTIFY ARCHITECTURE FOUNDATION</span>
        <h1 className="hero-title">VOTIFY</h1>
        <p className="hero-tagline">Real-time polling, beautifully simple.</p>
      </header>

      <main className="section-grid">
        {/* Japanese Editorial Design Tokens */}
        <section className="card">
          <div className="card-title">
            <span>Design System Tokens</span>
            <span style={{ fontSize: '12px', color: 'var(--v-color-vermilion)' }}>Japanese Editorial</span>
          </div>
          <p className="card-desc">
            Built using custom CSS custom properties (variables) featuring warm washi backgrounds, deep ink typography, and signature vermilion accents.
          </p>

          <div className="palette-grid">
            {colorSwatches.map((s) => (
              <div
                key={s.name}
                className="swatch"
                style={{
                  backgroundColor: `var(${s.var})`,
                  color: s.textLight ? '#FFFFFF' : '#1A1A1A',
                }}
              >
                <div className="swatch-label">
                  <div>{s.name}</div>
                  <div style={{ opacity: 0.8, fontSize: '10px' }}>{s.code}</div>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* UI Components Showcase */}
        <section className="card">
          <div className="card-title">
            <span>Button Component</span>
            <span style={{ fontSize: '12px', color: 'var(--v-color-warm-gray)' }}>Reusable UI</span>
          </div>
          <p className="card-desc">
            Accessible, state-aware button components adhering to Votify spacing, typography, and hover micro-interactions.
          </p>

          <div className="button-showcase" style={{ marginBottom: '16px' }}>
            <Button variant="primary" size="md">Primary Button</Button>
            <Button variant="secondary" size="md">Secondary</Button>
            <Button variant="outline" size="md">Outline</Button>
            <Button variant="ghost" size="md">Ghost</Button>
          </div>

          <div className="button-showcase">
            <Button variant="primary" size="sm">Small</Button>
            <Button variant="primary" size="md">Medium</Button>
            <Button variant="primary" size="lg">Large</Button>
            <Button variant="primary" size="md" disabled>Disabled</Button>
          </div>
        </section>

        {/* Backend Microservices Architecture */}
        <section className="card" style={{ gridColumn: '1 / -1' }}>
          <div className="card-title">
            <span>Backend Microservices Topology</span>
            <span style={{ fontSize: '12px', color: 'var(--v-color-success)' }}>7 Go Services Ready</span>
          </div>
          <p className="card-desc">
            Independent Go modules communicating via Gin HTTP routes, Kafka domain events, and Redis pub/sub channels.
          </p>

          <div className="service-list">
            {services.map((svc) => (
              <div key={svc.name} className="service-item">
                <div className="service-info">
                  <span className="dot"></span>
                  <span className="service-name">{svc.name}</span>
                </div>
                <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                  <span style={{ fontSize: '12px', color: 'var(--v-color-warm-gray)' }}>{svc.route}</span>
                  <span className="service-port">:{svc.port}</span>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Application Route Stubs */}
        <section className="card" style={{ gridColumn: '1 / -1' }}>
          <div className="card-title">
            <span>Application Routing Blueprint</span>
            <span style={{ fontSize: '12px', color: 'var(--v-color-warm-gray)' }}>{APP_ROUTES.length} Routes Defined</span>
          </div>
          <p className="card-desc">
            Navigation structure and API service boundaries ready for feature development.
          </p>

          <div className="service-list">
            {APP_ROUTES.map((route) => (
              <div key={route.path} className="service-item">
                <div className="service-info">
                  <span className="service-name" style={{ fontFamily: 'monospace' }}>{route.path}</span>
                  <span style={{ fontSize: '13px', color: 'var(--v-color-ink-black)' }}>— {route.name}</span>
                </div>
                <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
                  <span style={{ fontSize: '12px', color: 'var(--v-color-warm-gray)' }}>{route.description}</span>
                  <span className="service-port">{route.service}</span>
                </div>
              </div>
            ))}
          </div>
        </section>
      </main>

      <footer className="footer-note">
        <p>Votify Foundation v{ENV.APP_VERSION} | Environment: {ENV.APP_ENV} | Monorepo Structure Established</p>
      </footer>
    </div>
  );
};

export default App;
