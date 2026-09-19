# Votify System Architecture

> "Real-time polling, beautifully simple."

## System Topology & Microservices Overview

Votify is designed as an event-driven microservices architecture built around Go backend services, a React + Vite frontend, and high-performance infrastructure components (MongoDB, Redis, Apache Kafka, WebSockets).

```
                            ┌────────────────────────┐
                            │    React + Vite SPA    │
                            └───────────┬────────────┘
                                        │
                         ┌──────────────┴──────────────┐
                         │                             │ (HTTP / WebSocket)
                         ▼                             ▼
              ┌────────────────────┐       ┌──────────────────────┐
              │    API Gateway     │       │   Realtime Service   │
              │     (Port 8080)    │       │     (Port 8084)      │
              └──────────┬─────────┘       └──────────▲───────────┘
                         │                            │ (Redis PubSub)
      ┌──────────────────┼──────────────────┐         │
      ▼                  ▼                  ▼         │
┌───────────┐      ┌───────────┐      ┌───────────┐   │
│   Auth    │      │   Poll    │      │   Vote    ├───┴──────────┐
│ Service   │      │ Service   │      │ Service   │              │
│  (:8081)  │      │  (:8082)  │      │  (:8083)  │              │
└─────┬─────┘      └─────┬─────┘      └─────┬─────┘              │
      │                  │                  │                    │
      ▼                  ▼                  ▼                    │
┌───────────┐      ┌───────────┐      ┌───────────┐              │
│ Auth DB   │      │  Poll DB  │      │   Redis   │              │
│ (MongoDB) │      │ (MongoDB) │      │  (Cache)  │              │
└───────────┘      └───────────┘      └───────────┘              │
                                            ▲                    │
                                            │                    ▼
                                 ┌──────────┴────────────────────────┐
                                 │       Apache Kafka Event Bus      │
                                 └──────────┬──────────────────┬─────┘
                                            │                  │
                                            ▼                  ▼
                                     ┌───────────┐      ┌───────────┐
                                     │ Analytics │      │  Payment  │
                                     │  Service  │      │  Service  │
                                     │  (:8085)  │      │  (:8086)  │
                                     └───────────┘      └───────────┘
```

---

## Service Responsibilities & Ports

| Service | Port | Description & Responsibilities | Data Store |
|---|---|---|---|
| **API Gateway** | `8080` | Entry point for HTTP traffic. Route forwarding, rate limiting, JWT validation, CORS header enforcement. | None |
| **Auth Service** | `8081` | User registration, login, bcrypt password hashing, JWT token issuance & validation. | MongoDB (`auth_db`) |
| **Poll Service** | `8082` | Poll creation, update, deletion, listing, status toggling, configuration enforcement. | MongoDB (`poll_db`) |
| **Vote Service** | `8083` | Ingestion of high-throughput votes, atomic count updates in Redis, publishing raw vote events to Kafka. | Redis |
| **Realtime Service** | `8084` | WebSocket connections with browser clients, listening to Redis Pub/Sub for immediate vote updates. | Redis (Pub/Sub) |
| **Analytics Service** | `8085` | Consumes Kafka vote events to build time-series analytics, voter demographics, and aggregate reports. | MongoDB (`analytics_db`) |
| **Payment Service** | `8086` | Subscription management, Stripe integration / webhooks, plan tier enforcement. | MongoDB (`payment_db`) |

---

## Data Ownership Principles

1. **Database per Service**: Microservices MUST NOT query each other's databases directly. Cross-service data requests are satisfied via API Gateway routing or Kafka domain events.
2. **Redis as Hot State**: Vote tallies are cached in Redis Hash structures (`poll:{id}:counts`) to eliminate DB bottlenecks during live polling bursts.
3. **Kafka as Event Log**: Every vote cast produces an immutable `VoteCreated` event published to `votify.votes.raw`. This enables asynchronous processing without blocking the voting HTTP request.

---

## Authentication & Security Flow

- Auth Service generates RS256 / HS256 JWT tokens containing `user_id`, `email`, and `role`.
- API Gateway verifies JWT headers on protected routes (`/api/v1/polls/new`, `/api/v1/payments/*`) before forwarding requests downstream.
- Public voting endpoints (`POST /api/v1/votes`) are unauthenticated by default, but protected by IP rate-limiting in API Gateway and vote fingerprinting in Vote Service.

---

## Realtime Architecture Flow

1. Browser opens WebSocket connection to `ws://localhost:8084/ws/polls/{id}`.
2. Realtime Service registers connection in connection manager and subscribes to Redis channel `channel:poll:{id}`.
3. Voter casts vote → `POST /api/v1/votes` → Vote Service.
4. Vote Service updates Redis counter and publishes vote event to Redis Pub/Sub `channel:poll:{id}`.
5. Realtime Service receives Redis message and immediately broadcasts JSON update frame to all connected WebSocket clients for that poll.
