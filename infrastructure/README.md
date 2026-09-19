# Votify Cloud Infrastructure Blueprint (Zero Docker)

Votify is designed as a cloud-native platform deployed on a single cloud platform (**Render** / **Railway**) with managed cloud data providers. No Docker containers or server maintenance are required.

## Cloud Architecture Blueprint

```
                      ┌────────────────────────────────────────────────────────┐
                      │              SINGLE CLOUD PLATFORM                     │
                      │               (Render / Railway)                       │
                      │                                                        │
                      │   ┌────────────────────────────────────────────────┐   │
                      │   │       Frontend (Vite + React SPA)              │   │
                      │   └───────────────────────┬────────────────────────┘   │
                      │                           │                            │
                      │          ┌────────────────┴──────────────┐             │
                      │          │ HTTP                          │ WebSockets  │
                      │          ▼                               ▼             │
                      │   ┌──────────────┐             ┌───────────────────┐   │
                      │   │ API Gateway  │             │ Realtime Service  │   │
                      │   │  (Port 8080) │             │    (Port 8084)    │   │
                      │   └──────┬───────┘             └─────────▲─────────┘   │
                      │          │                               │ (PubSub)    │
                      │  ┌───────┼───────┐                       │             │
                      │  ▼       ▼       ▼                       │             │
                      │ Auth   Poll    Vote ─────────────────────┼─────────┐   │
                      │ (:8081)(:8082) (:8083)                   │         │   │
                      └──┬───────┬───────┬───────────────────────┼─────────┼───┘
                         │       │       │                       │         │
                         ▼       ▼       ▼                       │         ▼
                     ┌──────┐ ┌──────┐ ┌───────────────────┐     │     ┌─────────┐
                     │ Mongo│ │ Mongo│ │    Redis Cloud    ├─────┘     │  Aiven  │
                     │Atlas │ │Atlas │ │ (Managed Cache &  │           │  Kafka  │
                     │ (DB) │ │ (DB) │ │     Pub/Sub)      │           │(EventBus│
                     └──────┘ └──────┘ └───────────────────┘           └─────────┘
```

## Infrastructure Providers Overview

| Infrastructure Component | Provider | Description |
|---|---|---|
| **Application Platform** | **Render** / **Railway** | Single platform hosting Frontend SPA and all Go backend microservices directly from GitHub. |
| **Database** | **MongoDB Atlas** | Managed MongoDB cloud cluster for Auth (`auth_db`), Poll (`poll_db`), and Payment (`payment_db`). |
| **Cache & Realtime Pub/Sub** | **Redis Cloud** | Managed Redis instance for vote caching and Redis Pub/Sub channels. |
| **Event Bus** | **Aiven Kafka** | Managed Apache Kafka cluster for domain event streaming (`votify.votes.raw`). |
