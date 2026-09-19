import React, { useEffect, useState } from 'react';
import { Button } from '../components/ui/Button';
import { apiClient } from '../lib/api';
import type { Poll, PollOption } from '../types';

interface PublicVoteViewProps {
  pollId: string;
  onNavigateResults: (pollId: string) => void;
  onNavigateDashboard?: () => void;
}

export const PublicVoteView: React.FC<PublicVoteViewProps> = ({
  pollId,
  onNavigateResults,
  onNavigateDashboard,
}) => {
  const [poll, setPoll] = useState<Poll | null>(null);
  const [selectedOptionId, setSelectedOptionId] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [votedSuccess, setVotedSuccess] = useState(false);

  useEffect(() => {
    const loadPoll = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await apiClient.get<Poll>(`/polls/public/${pollId}`);
        setPoll(data);
      } catch (err: any) {
        setError(err.message || 'Poll not found or not open for public voting.');
      } finally {
        setLoading(false);
      }
    };

    loadPoll();
  }, [pollId]);

  const handleVote = async () => {
    if (!selectedOptionId || !poll) return;
    setError(null);
    setSubmitting(true);

    try {
      await apiClient.post('/votes', {
        pollId: poll.id,
        optionId: selectedOptionId,
      });
      setVotedSuccess(true);
    } catch (err: any) {
      setError(err.message || 'Failed to submit vote.');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div style={{ maxWidth: '600px', margin: '60px auto', padding: '0 20px', textAlign: 'center', color: 'var(--v-color-warm-gray)' }}>
        Loading poll...
      </div>
    );
  }

  if (error && !poll) {
    return (
      <div style={{ maxWidth: '600px', margin: '60px auto', padding: '0 20px' }}>
        <div className="card" style={{ padding: '32px', textAlign: 'center' }}>
          <span className="brand-badge" style={{ backgroundColor: 'rgba(196, 52, 45, 0.12)', color: 'var(--v-color-error)' }}>
            POLL UNAVAILABLE
          </span>
          <h2 style={{ fontFamily: 'var(--v-font-heading)', fontSize: '24px', color: 'var(--v-color-deep-ink)', marginTop: '8px' }}>
            Poll Not Found
          </h2>
          <p style={{ color: 'var(--v-color-warm-gray)', marginTop: '8px', marginBottom: '24px' }}>
            {error}
          </p>
          {onNavigateDashboard && (
            <Button variant="primary" size="md" onClick={onNavigateDashboard}>
              Return to Dashboard
            </Button>
          )}
        </div>
      </div>
    );
  }

  if (votedSuccess && poll) {
    const chosenOption = poll.options.find((o) => o.id === selectedOptionId);
    return (
      <div style={{ maxWidth: '600px', margin: '60px auto', padding: '0 20px' }}>
        <div className="card" style={{ padding: '40px', textAlign: 'center' }}>
          <div style={{
            width: '48px',
            height: '48px',
            borderRadius: '50%',
            backgroundColor: 'rgba(45, 122, 79, 0.12)',
            color: 'var(--v-color-success)',
            fontSize: '24px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            margin: '0 auto 16px auto'
          }}>
            ✓
          </div>
          <span className="brand-badge" style={{ backgroundColor: 'rgba(45, 122, 79, 0.12)', color: 'var(--v-color-success)' }}>
            VOTE RECORDED
          </span>
          <h2 style={{ fontFamily: 'var(--v-font-heading)', fontSize: '28px', color: 'var(--v-color-deep-ink)', marginTop: '8px' }}>
            Thank You for Voting!
          </h2>
          <p style={{ color: 'var(--v-color-warm-gray)', marginTop: '6px', fontSize: '15px' }}>
            Your choice <strong>"{chosenOption?.text}"</strong> has been submitted.
          </p>
          <div style={{ marginTop: '32px' }}>
            <Button variant="primary" size="lg" onClick={() => onNavigateResults(poll.id)}>
              View Live Results
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '600px', margin: '60px auto', padding: '0 20px' }}>
      <div className="card" style={{ padding: '32px' }}>
        <div style={{ marginBottom: '24px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
            <span className="brand-badge">PUBLIC POLL</span>
            <span style={{
              fontSize: '11px',
              fontWeight: 700,
              padding: '2px 8px',
              borderRadius: 'var(--v-radius-sm)',
              backgroundColor: poll?.status === 'open' ? 'rgba(45, 122, 79, 0.12)' : 'rgba(196, 52, 45, 0.12)',
              color: poll?.status === 'open' ? 'var(--v-color-success)' : 'var(--v-color-error)',
              textTransform: 'uppercase'
            }}>
              {poll?.status}
            </span>
          </div>

          <h1 style={{ fontFamily: 'var(--v-font-heading)', fontSize: '28px', color: 'var(--v-color-deep-ink)', lineHeight: 1.2 }}>
            {poll?.question}
          </h1>

          {poll?.description && (
            <p style={{ color: 'var(--v-color-warm-gray)', fontSize: '15px', marginTop: '8px' }}>
              {poll.description}
            </p>
          )}
        </div>

        {error && (
          <div style={{
            padding: '12px 16px',
            backgroundColor: 'rgba(196, 52, 45, 0.08)',
            borderLeft: '3px solid var(--v-color-error)',
            color: 'var(--v-color-error)',
            fontSize: '14px',
            marginBottom: '20px'
          }}>
            {error}
          </div>
        )}

        {poll?.status !== 'open' ? (
          <div style={{
            padding: '20px',
            backgroundColor: 'var(--v-color-washi-light)',
            borderRadius: 'var(--v-radius-md)',
            textAlign: 'center',
            color: 'var(--v-color-warm-gray)',
            marginBottom: '24px'
          }}>
            This poll is currently closed for voting. You can still view the final results.
            <div style={{ marginTop: '16px' }}>
              <Button variant="secondary" size="md" onClick={() => onNavigateResults(poll!.id)}>
                View Results
              </Button>
            </div>
          </div>
        ) : (
          <div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', marginBottom: '28px' }}>
              {poll?.options.map((opt: PollOption) => {
                const isSelected = selectedOptionId === opt.id;
                return (
                  <label
                    key={opt.id}
                    onClick={() => setSelectedOptionId(opt.id)}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: '14px',
                      padding: '16px 20px',
                      borderRadius: 'var(--v-radius-md)',
                      border: isSelected
                        ? '2px solid var(--v-color-vermilion)'
                        : '1px solid var(--v-color-stone)',
                      backgroundColor: isSelected ? 'var(--v-color-sakura-light)' : 'var(--v-color-white)',
                      cursor: 'pointer',
                      transition: 'all 0.15s ease-in-out',
                    }}
                  >
                    <input
                      type="radio"
                      name="poll-option"
                      checked={isSelected}
                      onChange={() => setSelectedOptionId(opt.id)}
                      style={{ accentColor: 'var(--v-color-vermilion)', width: '18px', height: '18px' }}
                    />
                    <span style={{ fontSize: '16px', fontWeight: 500, color: 'var(--v-color-deep-ink)' }}>
                      {opt.text}
                    </span>
                  </label>
                );
              })}
            </div>

            <div style={{ display: 'flex', gap: '12px' }}>
              <Button
                variant="primary"
                size="lg"
                fullWidth
                disabled={!selectedOptionId || submitting}
                onClick={handleVote}
              >
                {submitting ? 'Submitting Vote...' : 'Submit Vote'}
              </Button>
              <Button
                variant="outline"
                size="lg"
                onClick={() => onNavigateResults(poll!.id)}
              >
                View Results
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
