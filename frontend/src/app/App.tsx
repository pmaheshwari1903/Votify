import React, { useEffect, useState } from 'react';
import { Button } from '../components/ui/Button';
import { LoginView } from '../views/LoginView';
import { DashboardView } from '../views/DashboardView';
import { CreatePollView } from '../views/CreatePollView';
import { PublicVoteView } from '../views/PublicVoteView';
import { ResultsView } from '../views/ResultsView';
import { apiClient } from '../lib/api';
import type { User } from '../types';
import './App.css';

type ViewMode = 'dashboard' | 'create' | 'login' | 'public-vote' | 'results';

export const App: React.FC = () => {
  const [user, setUser] = useState<User | null>(null);
  const [view, setView] = useState<ViewMode>('dashboard');
  const [activePollId, setActivePollId] = useState<string>('');
  const [initializing, setInitializing] = useState(true);

  useEffect(() => {
    // Check URL pathname for direct links e.g. /poll/:id or /poll/:id/results
    const path = window.location.pathname;
    if (path.startsWith('/poll/')) {
      const parts = path.split('/');
      const id = parts[2];
      if (id) {
        setActivePollId(id);
        if (parts[3] === 'results') {
          setView('results');
        } else {
          setView('public-vote');
        }
      }
    }

    // Check auth via HttpOnly cookie (sent automatically with credentials)
    apiClient
      .get<User>('/auth/me')
      .then((u) => {
        setUser(u);
      })
      .catch(() => {
        // Not authenticated — cookie missing or expired
      })
      .finally(() => setInitializing(false));
  }, []);

  const handleLogout = () => {
    apiClient.post('/auth/logout').catch(() => {});
    setUser(null);
    setView('login');
  };

  const navigateToVote = (pollId: string) => {
    setActivePollId(pollId);
    window.history.pushState({}, '', `/poll/${pollId}`);
    setView('public-vote');
  };

  const navigateToResults = (pollId: string) => {
    setActivePollId(pollId);
    window.history.pushState({}, '', `/poll/${pollId}/results`);
    setView('results');
  };

  const navigateToDashboard = () => {
    window.history.pushState({}, '', '/dashboard');
    setView('dashboard');
  };

  if (initializing) {
    return (
      <div style={{ maxWidth: '400px', margin: '100px auto', textAlign: 'center', color: 'var(--v-color-warm-gray)' }}>
        Initializing Votify...
      </div>
    );
  }

  return (
    <div style={{ minHeight: '100vh', backgroundColor: 'var(--v-color-washi)' }}>
      {/* Navigation Header */}
      <nav className="v-navbar">
        <div className="v-navbar-inner">
          <div className="v-logo" onClick={navigateToDashboard}>
            <span>VOTIFY</span>
            <span className="v-logo-dot" />
          </div>

          <div className="v-nav-actions">
            {user ? (
              <>
                <div className="v-user-chip">
                  <span style={{ width: '6px', height: '6px', borderRadius: '50%', backgroundColor: 'var(--v-color-success)' }} />
                  {user.name}
                </div>
                <Button variant="outline" size="sm" onClick={navigateToDashboard}>
                  Dashboard
                </Button>
                <Button variant="ghost" size="sm" onClick={handleLogout}>
                  Logout
                </Button>
              </>
            ) : (
              <Button variant="primary" size="sm" onClick={() => setView('login')}>
                Sign In
              </Button>
            )}
          </div>
        </div>
      </nav>

      {/* Main View Router */}
      <main>
        {view === 'login' && <LoginView />}

        {view === 'dashboard' && (
          user ? (
            <DashboardView
              onNavigateCreate={() => setView('create')}
              onNavigateVote={navigateToVote}
              onNavigateResults={navigateToResults}
            />
          ) : (
            <LoginView />
          )
        )}

        {view === 'create' && (
          <CreatePollView
            onSuccess={(poll) => navigateToVote(poll.id)}
            onCancel={navigateToDashboard}
          />
        )}

        {view === 'public-vote' && (
          <PublicVoteView
            pollId={activePollId}
            onNavigateResults={navigateToResults}
            onNavigateDashboard={navigateToDashboard}
          />
        )}

        {view === 'results' && (
          <ResultsView
            pollId={activePollId}
            onNavigateVote={navigateToVote}
            onNavigateDashboard={navigateToDashboard}
          />
        )}
      </main>
    </div>
  );
};

export default App;
