# Votify Single-Platform Cloud Deployment Guide (Zero Docker)

This guide details how to deploy the **entire Votify application** (Frontend React SPA + All 7 Go Backend Microservices) onto a **single unified cloud platform** (**Render** or **Railway**), paired with managed cloud services for MongoDB, Redis, and Kafka.

---

## Unified Stack Architecture

- **Single Host Platform**: **Render** (or **Railway**)
  - Frontend SPA (Static Web Site)
  - API Gateway (Web Service)
  - Auth Service (Web Service)
  - Poll Service (Web Service)
  - Vote Service (Web Service)
  - Realtime Service (Web Service with WebSockets)
  - Analytics Service (Web Service / Worker)
  - Payment Service (Web Service)
- **Managed Data Services**:
  - **MongoDB Atlas**: Managed MongoDB cluster
  - **Redis Cloud**: Managed Redis for caching & Pub/Sub
  - **Aiven for Apache Kafka**: Managed Kafka event broker

---

## Step-by-Step Deployment Blueprint

### Step 1: Managed Cloud Data Services Setup

1. **MongoDB Atlas**:
   - Create a free cluster on MongoDB Atlas.
   - Database user credentials and IP Access List (`0.0.0.0/0` for cloud services).
   - Connection String format: `mongodb+srv://<username>:<password>@cluster0.mongodb.net/`

2. **Redis Cloud**:
   - Create a database instance on Redis Cloud.
   - Note the endpoint address (`redis-XXXXX.c1.us-east-1-1.ec2.cloud.redislabs.com:PORT`) and default password.

3. **Aiven for Apache Kafka**:
   - Create an Apache Kafka service on Aiven.
   - Copy the Service URI (`kafka-XXXXX.aivencloud.com:PORT`) and SSL certificate / SASL credentials.

---

### Step 2: Deploying Everything on Render (Single Platform)

Connect your GitHub repository to Render and create the following services under a single Project:

#### 1. Frontend Web Service
- **Type**: Static Site
- **Build Command**: `npm run build`
- **Publish Directory**: `dist`
- **Root Directory**: `frontend`
- **Environment Variables**:
  - `VITE_API_BASE_URL`: `https://votify-gateway.onrender.com/api/v1`
  - `VITE_WS_URL`: `wss://votify-realtime.onrender.com/ws`

#### 2. API Gateway Web Service
- **Type**: Web Service
- **Root Directory**: `services/api-gateway`
- **Build Command**: `go build -o server ./cmd/main.go`
- **Start Command**: `./server`
- **Environment Variables**:
  - `PORT`: `8080`
  - `AUTH_SERVICE_URL`: `https://votify-auth.onrender.com`
  - `POLL_SERVICE_URL`: `https://votify-poll.onrender.com`
  - `VOTE_SERVICE_URL`: `https://votify-vote.onrender.com`
  - `REALTIME_SERVICE_URL`: `https://votify-realtime.onrender.com`

#### 3. Realtime Service (WebSockets)
- **Type**: Web Service (WebSockets enabled)
- **Root Directory**: `services/realtime-service`
- **Build Command**: `go build -o server ./cmd/main.go`
- **Start Command**: `./server`
- **Environment Variables**:
  - `REDIS_ADDR`: `<your-redis-cloud-endpoint>`
  - `REDIS_PASSWORD`: `<your-redis-cloud-password>`

#### 4. Auth, Poll, Vote, Analytics & Payment Services
- Deploy each service under `services/<service-name>` using the same native Go build command (`go build -o server ./cmd/main.go`).
- Supply environment variables connecting them to MongoDB Atlas, Redis Cloud, and Aiven Kafka.

---

## Why Single-Platform + Managed Cloud is Ideal

1. **Unified Dashboard & Billing**: Manage all frontend and backend services in one place (Render or Railway).
2. **Zero Container Maintenance**: No Dockerfiles to update, no container registry management, no Kubernetes YAML.
3. **Automated CI/CD**: Pushing code to `main` on GitHub automatically triggers build and deploy for affected services.
4. **Native WebSockets**: Web Services on Render/Railway natively support persistent WebSocket connections.
