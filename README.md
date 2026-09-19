# Votify

**Real-time polling, beautifully simple.**

Votify is a production-quality SaaS platform for creating, sharing, and analyzing live polls. Built with a cloud-native microservice architecture, it provides real-time voting updates, shareable poll links, and rich analytics — all wrapped in a refined Japanese editorial-inspired design.

## Architecture

```
┌─────────────┐
│   Frontend  │    React + Vite + TypeScript (Static Site on Render/Railway)
│  (SPA)      │    Port 5173
└──────┬──────┘
       │ HTTP / WebSockets
       ▼
┌──────────────┐
│ API Gateway  │    Go + Gin — Port 8080 (Web Service)
│              │    Routes requests to services
└──┬──┬──┬──┬──┘
   │  │  │  │
   ▼  ▼  ▼  ▼
┌──────┐ ┌──────┐ ┌──────┐ ┌──────────┐
│ Auth │ │ Poll │ │ Vote │ │ Realtime │
│ 8081 │ │ 8082 │ │ 8083 │ │  8084    │
└──────┘ └──────┘ └──┬───┘ └────┬─────┘
                     │          │
                     ▼          ▼
             ┌──────────────┐ ┌─────────────┐
             │ Aiven Kafka  │ │ Redis Cloud │
             └──────┬───────┘ └─────────────┘
                    │
           ┌────────┴─────────┐
           ▼                  ▼
      ┌──────────┐       ┌─────────┐
      │Analytics │       │ Payment │
      │  8085    │       │  8086   │
      └──────────┘       └─────────┘
```

Each service owns its own database on MongoDB Atlas. Services communicate asynchronously via Aiven Kafka domain events and in-memory Redis Cloud channels.

## Services

| Service | Port | Responsibility |
|---|---|---|
| **API Gateway** | 8080 | Public entry point. Routes, auth forwarding, rate limiting |
| **Auth Service** | 8081 | Authentication, user management, JWT tokens |
| **Poll Service** | 8082 | Poll CRUD, lifecycle management |
| **Vote Service** | 8083 | Vote ingestion, tallying, Kafka publishing |
| **Realtime Service** | 8084 | WebSocket connections, Redis pub/sub, live updates |
| **Analytics Service** | 8085 | Read models, poll analytics, Kafka consumption |
| **Payment Service** | 8086 | Subscriptions, payments, plan management |

## Cloud Tech Stack (Zero Docker)

| Layer | Technology & Host |
|---|---|
| Platform Host | **Render** / **Railway** (Unified platform for Frontend & Go Services) |
| Frontend | React, Vite, TypeScript, Vanilla CSS |
| Backend | Go, Gin |
| Database | **MongoDB Atlas** (Managed MongoDB per service) |
| Realtime | **Redis Cloud** (Managed Redis pub/sub & caching) |
| Events | **Aiven Kafka** (Managed event streaming) |
| Transport | WebSocket (browser ↔ Realtime Service) |

## Quickstart

```bash
# 1. Configure environment
cp .env.example .env

# 2. Run Frontend
cd frontend
npm install
npm run dev

# 3. Build & Run Go Service (Native Go)
cd services/api-gateway
go run ./cmd/main.go
```

## Documentation

- [Architecture](docs/architecture.md) — System topology, service responsibilities, data flow
- [Conventions](docs/conventions.md) — Code organization, naming, error handling patterns
- [Single-Platform Deployment Guide](docs/deployment.md) — Step-by-step setup for Render/Railway + MongoDB Atlas + Redis Cloud + Aiven Kafka

## License

Private — All rights reserved.
