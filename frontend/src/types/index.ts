// Votify Frontend Domain Type Definitions
// Defines core entities and API schemas aligned with microservice contracts.

export interface User {
  id: string;
  email: string;
  name: string;
  role: 'user' | 'admin';
  createdAt: string;
}

export interface PollOption {
  id: string;
  text: string;
  voteCount: number;
}

export interface PollSettings {
  isAnonymous: boolean;
  allowMultipleVotes: boolean;
  requireAuth: boolean;
  expiresAt?: string;
}

export interface Poll {
  id: string;
  title: string;
  description?: string;
  creatorId: string;
  options: PollOption[];
  settings: PollSettings;
  isActive: boolean;
  totalVotes: number;
  createdAt: string;
  updatedAt: string;
}

export interface Vote {
  pollId: string;
  optionId: string;
  voterIp?: string;
  userId?: string;
  timestamp: string;
}

export interface PollAnalytics {
  pollId: string;
  totalVotes: number;
  votesByOption: Record<string, number>;
  votesOverTime: Array<{ timestamp: string; count: number }>;
  uniqueVoters: number;
}

export interface SubscriptionPlan {
  id: string;
  name: string;
  priceMonthly: number;
  maxPolls: number;
  maxVotesPerPoll: number;
  customBranding: boolean;
  analyticsLevel: 'basic' | 'advanced' | 'enterprise';
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
  meta?: {
    timestamp: string;
    requestId?: string;
  };
}
