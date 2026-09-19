// Votify Route Stubs
// Outlines application views for future feature implementation.

export interface RouteDefinition {
  path: string;
  name: string;
  description: string;
  service: string;
}

export const APP_ROUTES: RouteDefinition[] = [
  {
    path: '/',
    name: 'Home / Showcase',
    description: 'Landing page and Design System foundation showcase',
    service: 'Frontend',
  },
  {
    path: '/auth/login',
    name: 'Login',
    description: 'User authentication & JWT token issuance',
    service: 'Auth Service (:8081)',
  },
  {
    path: '/auth/register',
    name: 'Register',
    description: 'Account registration and password hashing',
    service: 'Auth Service (:8081)',
  },
  {
    path: '/polls',
    name: 'Poll Dashboard',
    description: 'List user polls and quick status stats',
    service: 'Poll Service (:8082)',
  },
  {
    path: '/polls/new',
    name: 'Create Poll',
    description: 'Interactive poll builder with real-time preview',
    service: 'Poll Service (:8082)',
  },
  {
    path: '/poll/:id',
    name: 'Public Vote Screen',
    description: 'Public voting page with live WebSocket stream',
    service: 'Vote Service (:8083) & Realtime (:8084)',
  },
  {
    path: '/poll/:id/analytics',
    name: 'Poll Analytics',
    description: 'Live charts, vote totals, and demographic breakdown',
    service: 'Analytics Service (:8085)',
  },
  {
    path: '/pricing',
    name: 'Pricing & Subscription',
    description: 'Plan comparison and payment checkout',
    service: 'Payment Service (:8086)',
  },
];
