import { useEffect, useState, useRef } from 'react';
import { ENV } from '../config/env';
import type { PollOption } from '../types';

export type RealtimeStatus = 'connected' | 'reconnecting' | 'disconnected';

export interface RealtimeResultsPayload {
  type: string;
  pollId: string;
  totalVotes: number;
  options: Array<{ id: string; voteCount: number }>;
}

export function useRealtimeResults(pollId: string, initialOptions?: PollOption[], initialTotalVotes?: number) {
  const [options, setOptions] = useState<PollOption[]>(initialOptions || []);
  const [totalVotes, setTotalVotes] = useState<number>(initialTotalVotes || 0);
  const [status, setStatus] = useState<RealtimeStatus>('disconnected');
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (initialOptions && initialOptions.length > 0) {
      setOptions(initialOptions);
    }
    if (initialTotalVotes !== undefined) {
      setTotalVotes(initialTotalVotes);
    }
  }, [initialOptions, initialTotalVotes]);

  useEffect(() => {
    if (!pollId) return;

    let isUnmounted = false;
    const wsUrl = `${ENV.WS_URL}/${pollId}`;

    const connect = () => {
      setStatus('reconnecting');
      const ws = new WebSocket(wsUrl);
      socketRef.current = ws;

      ws.onopen = () => {
        if (!isUnmounted) setStatus('connected');
      };

      ws.onmessage = (event) => {
        if (isUnmounted) return;
        try {
          const data: RealtimeResultsPayload = JSON.parse(event.data);
          if (data.type === 'POLL_RESULTS_UPDATED' && data.pollId === pollId) {
            setTotalVotes(data.totalVotes);
            setOptions((prevOptions) =>
              prevOptions.map((opt) => {
                const updated = data.options.find((o) => o.id === opt.id);
                return updated ? { ...opt, voteCount: updated.voteCount } : opt;
              })
            );
          }
        } catch {
          // Ignore malformed frames
        }
      };

      ws.onerror = () => {
        if (!isUnmounted) setStatus('disconnected');
      };

      ws.onclose = () => {
        if (!isUnmounted) {
          setStatus('disconnected');
          // Auto-reconnect after 3 seconds
          setTimeout(() => {
            if (!isUnmounted) connect();
          }, 3000);
        }
      };
    };

    connect();

    return () => {
      isUnmounted = true;
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, [pollId]);

  return { options, totalVotes, status };
}
