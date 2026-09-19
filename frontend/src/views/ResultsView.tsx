import React, { useEffect, useState } from 'react';
import { Button } from '../components/ui/Button';
import { apiClient } from '../lib/api';
import { useRealtimeResults } from '../hooks/useRealtimeResults';
import type { PollOption, PollResults } from '../types';

interface ResultsViewProps {
  pollId: string;
  onNavigateVote: (pollId: string) => void;
  onNavigateDashboard?: () => void;
}

export const ResultsView: React.FC<ResultsViewProps> = ({
  pollId,
  onNavigateVote,
  onNavigateDashboard,
}) => {
  const [initialData, setInitialData] = useState<PollResults | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchResults = async () => {
    setLoading(true);
    setError(null);

    try {
      const data = await apiClient.get<PollResults>(
        `/polls/${pollId}/results`
      );
      setInitialData(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load poll results.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchResults();
  }, [pollId]);

  // Connect WebSocket for live updates without page refresh
  const { options, totalVotes, status: wsStatus } = useRealtimeResults(
    pollId,
    initialData?.options,
    initialData?.totalVotes
  );

  if (loading) {
    return (
      <div
        style={{
          maxWidth: '640px',
          margin: '60px auto',
          padding: '0 20px',
          textAlign: 'center',
          color: 'var(--v-color-warm-gray)',
        }}
      >
        Loading poll results...
      </div>
    );
  }

  if (error || !initialData) {
    return (
      <div
        style={{
          maxWidth: '640px',
          margin: '60px auto',
          padding: '0 20px',
        }}
      >
        <div
          className="card"
          style={{
            padding: '32px',
            textAlign: 'center',
          }}
        >
          <span
            className="brand-badge"
            style={{
              backgroundColor: 'rgba(196, 52, 45, 0.12)',
              color: 'var(--v-color-error)',
            }}
          >
            ERROR
          </span>

          <h2
            style={{
              fontFamily: 'var(--v-font-heading)',
              fontSize: '24px',
              color: 'var(--v-color-deep-ink)',
              marginTop: '8px',
            }}
          >
            Results Unavailable
          </h2>

          <p
            style={{
              color: 'var(--v-color-warm-gray)',
              marginTop: '8px',
              marginBottom: '24px',
            }}
          >
            {error || 'Could not load results.'}
          </p>

          {onNavigateDashboard && (
            <Button
              variant="primary"
              size="md"
              onClick={onNavigateDashboard}
            >
              Back to Dashboard
            </Button>
          )}
        </div>
      </div>
    );
  }

  return (
    <div
      style={{
        maxWidth: '640px',
        margin: '60px auto',
        padding: '0 20px',
      }}
    >
      <div className="card" style={{ padding: '32px' }}>
        <div style={{ marginBottom: '28px' }}>
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              marginBottom: '8px',
            }}
          >
            <span className="brand-badge">LIVE RESULTS</span>

            {/* Realtime Connection Status Indicator */}
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '8px',
              }}
            >
              <span
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  gap: '6px',
                  fontSize: '11px',
                  fontWeight: 700,
                  letterSpacing: '0.05em',
                  padding: '2px 8px',
                  borderRadius: 'var(--v-radius-sm)',
                  backgroundColor:
                    wsStatus === 'connected'
                      ? 'rgba(45, 122, 79, 0.12)'
                      : 'rgba(107, 100, 96, 0.12)',
                  color:
                    wsStatus === 'connected'
                      ? 'var(--v-color-success)'
                      : 'var(--v-color-warm-gray)',
                  textTransform: 'uppercase',
                }}
              >
                <span
                  style={{
                    width: '6px',
                    height: '6px',
                    borderRadius: '50%',
                    backgroundColor:
                      wsStatus === 'connected'
                        ? 'var(--v-color-success)'
                        : 'var(--v-color-warm-gray)',
                  }}
                />

                {wsStatus === 'connected'
                  ? 'LIVE'
                  : wsStatus === 'reconnecting'
                    ? 'RECONNECTING'
                    : 'OFFLINE'}
              </span>

              <span
                style={{
                  fontSize: '13px',
                  color: 'var(--v-color-warm-gray)',
                  fontWeight: 600,
                }}
              >
                {totalVotes} {totalVotes === 1 ? 'Total Vote' : 'Total Votes'}
              </span>
            </div>
          </div>

          <h1
            style={{
              fontFamily: 'var(--v-font-heading)',
              fontSize: '28px',
              color: 'var(--v-color-deep-ink)',
              lineHeight: 1.2,
            }}
          >
            {initialData.question}
          </h1>
        </div>

        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            gap: '20px',
            marginBottom: '32px',
          }}
        >
          {options.map((opt: PollOption) => {
            const votes = opt.voteCount || 0;
            const percentage =
              totalVotes > 0
                ? Math.round((votes / totalVotes) * 100)
                : 0;

            return (
              <div
                key={opt.id}
                style={{
                  display: 'flex',
                  flexDirection: 'column',
                  gap: '6px',
                }}
              >
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    fontSize: '14px',
                  }}
                >
                  <span
                    style={{
                      fontWeight: 600,
                      color: 'var(--v-color-deep-ink)',
                    }}
                  >
                    {opt.text}
                  </span>

                  <span
                    style={{
                      fontSize: '13px',
                      color: 'var(--v-color-warm-gray)',
                    }}
                  >
                    {votes} {votes === 1 ? 'vote' : 'votes'} ({percentage}%)
                  </span>
                </div>

                {/* Animated progress bar */}
                <div
                  style={{
                    width: '100%',
                    height: '14px',
                    backgroundColor: 'var(--v-color-washi-light)',
                    border: '1px solid var(--v-color-stone)',
                    borderRadius: 'var(--v-radius-sm)',
                    overflow: 'hidden',
                    position: 'relative',
                  }}
                >
                  <div
                    style={{
                      width: `${percentage}%`,
                      height: '100%',
                      backgroundColor:
                        percentage > 0
                          ? 'var(--v-color-vermilion)'
                          : 'transparent',
                      transition: 'width 0.4s ease-out',
                    }}
                  />
                </div>
              </div>
            );
          })}
        </div>

        <div
          style={{
            display: 'flex',
            gap: '12px',
            justifyContent: 'space-between',
          }}
        >
          <Button variant="outline" size="md" onClick={fetchResults}>
            ↻ Refresh Results
          </Button>

          <Button
            variant="primary"
            size="md"
            onClick={() => onNavigateVote(initialData.pollId)}
          >
            Go to Vote Screen
          </Button>
        </div>
      </div>
    </div>
  );
};