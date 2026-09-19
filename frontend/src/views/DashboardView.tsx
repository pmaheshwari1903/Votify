import React, { useEffect, useState } from 'react';
import { Button } from '../components/ui/Button';
import { apiClient } from '../lib/api';
import type { Poll } from '../types';

interface DashboardViewProps {
  onNavigateCreate: () => void;
  onNavigateVote: (pollId: string) => void;
  onNavigateResults: (pollId: string) => void;
}

export const DashboardView: React.FC<DashboardViewProps> = ({
  onNavigateCreate,
  onNavigateVote,
  onNavigateResults,
}) => {
  const [polls, setPolls] = useState<Poll[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [copiedId, setCopiedId] = useState<string | null>(null);

  const fetchPolls = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await apiClient.get<Poll[]>('/polls');
      setPolls(data || []);
    } catch (err: any) {
      setError(err.message || 'Failed to load polls.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPolls();
  }, []);

  const handleToggleStatus = async (poll: Poll) => {
    const action = poll.status === 'open' ? 'close' : 'open';
    try {
      await apiClient.post(`/polls/${poll.id}/${action}`);
      fetchPolls();
    } catch (err: any) {
      alert(err.message || `Failed to ${action} poll.`);
    }
  };

  const handleDelete = async (pollId: string) => {
    if (!confirm('Are you sure you want to delete this poll?')) return;
    try {
      await apiClient.delete(`/polls/${pollId}`);
      fetchPolls();
    } catch (err: any) {
      alert(err.message || 'Failed to delete poll.');
    }
  };

  const handleCopyLink = (pollId: string) => {
    const shareUrl = `${window.location.origin}/poll/${pollId}`;
    navigator.clipboard.writeText(shareUrl);
    setCopiedId(pollId);
    setTimeout(() => setCopiedId(null), 2500);
  };

  return (
    <div style={{ maxWidth: '1000px', margin: '40px auto', padding: '0 20px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
        <div>
          <span className="brand-badge">POLL DASHBOARD</span>
          <h2 style={{ fontFamily: 'var(--v-font-heading)', fontSize: '32px', color: 'var(--v-color-deep-ink)' }}>
            Your Polls
          </h2>
        </div>
        <Button variant="primary" size="md" onClick={onNavigateCreate}>
          + Create New Poll
        </Button>
      </div>

      {error && (
        <div style={{
          padding: '12px 16px',
          backgroundColor: 'rgba(196, 52, 45, 0.08)',
          borderLeft: '3px solid var(--v-color-error)',
          color: 'var(--v-color-error)',
          fontSize: '14px',
          marginBottom: '24px'
        }}>
          {error}
        </div>
      )}

      {loading ? (
        <div style={{ textAlign: 'center', padding: '48px 0', color: 'var(--v-color-warm-gray)' }}>
          Loading your polls...
        </div>
      ) : polls.length === 0 ? (
        <div className="card" style={{ textAlign: 'center', padding: '60px 20px' }}>
          <h3 style={{ fontFamily: 'var(--v-font-heading)', fontSize: '20px', color: 'var(--v-color-deep-ink)' }}>
            No Polls Created Yet
          </h3>
          <p style={{ color: 'var(--v-color-warm-gray)', marginTop: '8px', marginBottom: '24px' }}>
            Create your first live poll in seconds and share it with your audience.
          </p>
          <Button variant="primary" size="md" onClick={onNavigateCreate}>
            Create Your First Poll
          </Button>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          {polls.map((poll) => (
            <div key={poll.id} className="card" style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                <div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '6px' }}>
                    <span style={{
                      display: 'inline-block',
                      padding: '2px 8px',
                      borderRadius: 'var(--v-radius-sm)',
                      fontSize: '11px',
                      fontWeight: 700,
                      letterSpacing: '0.05em',
                      backgroundColor: poll.status === 'open' ? 'rgba(45, 122, 79, 0.12)' : 'rgba(107, 100, 96, 0.12)',
                      color: poll.status === 'open' ? 'var(--v-color-success)' : 'var(--v-color-warm-gray)',
                      textTransform: 'uppercase'
                    }}>
                      {poll.status}
                    </span>
                    <span style={{ fontSize: '12px', color: 'var(--v-color-warm-gray)' }}>
                      {poll.options.length} Options • {poll.totalVotes} Total Votes
                    </span>
                  </div>
                  <h3 style={{ fontFamily: 'var(--v-font-heading)', fontSize: '20px', color: 'var(--v-color-deep-ink)' }}>
                    {poll.question}
                  </h3>
                  {poll.description && (
                    <p style={{ fontSize: '14px', color: 'var(--v-color-warm-gray)', marginTop: '4px' }}>
                      {poll.description}
                    </p>
                  )}
                </div>
              </div>

              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', paddingTop: '12px', borderTop: '1px solid var(--v-color-stone)' }}>
                <Button variant="outline" size="sm" onClick={() => handleCopyLink(poll.id)}>
                  {copiedId === poll.id ? '✓ Link Copied!' : 'Copy Share Link'}
                </Button>
                <Button variant="secondary" size="sm" onClick={() => onNavigateVote(poll.id)}>
                  Public Voting View
                </Button>
                <Button variant="outline" size="sm" onClick={() => onNavigateResults(poll.id)}>
                  View Results ({poll.totalVotes})
                </Button>
                <Button variant="ghost" size="sm" onClick={() => handleToggleStatus(poll)}>
                  {poll.status === 'open' ? 'Close Poll' : 'Open Poll'}
                </Button>
                <Button variant="ghost" size="sm" onClick={() => handleDelete(poll.id)} style={{ color: 'var(--v-color-error)' }}>
                  Delete
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
