// Votify Frontend Domain Type Definitions

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

export interface Poll {
  id: string;
  ownerId?: string;
  question: string;
  title?: string;
  description?: string;
  options: PollOption[];
  status: 'draft' | 'open' | 'closed';
  totalVotes: number;
  createdAt: string;
  updatedAt?: string;
}

export interface CreatePollPayload {
  question: string;
  description?: string;
  options: string[];
}

export interface CastVotePayload {
  pollId: string;
  optionId: string;
}

export interface PollResults {
  pollId: string;
  question: string;
  status: 'draft' | 'open' | 'closed';
  totalVotes: number;
  options: PollOption[];
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}