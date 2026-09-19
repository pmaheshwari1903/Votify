# Votify Backend & Architecture Conventions

Guidelines and rules for maintaining code consistency, security, and architectural integrity across all Votify microservices and frontend modules.

---

## 1. Request DTO & Transport Boundaries

### Structure & Location
- DTOs (Data Transfer Objects) represent request/response structures at the transport boundary (`internal/dto/`).
- DTOs use Gin struct tags for binding and transport validation:
  ```go
  type CreatePollRequest struct {
      Title    string   `json:"title" binding:"required,min=3,max=200"`
      Options  []string `json:"options" binding:"required,min=2,max=10"`
  }
  ```
- DTOs are **not** domain models and **not** database entities. Handlers parse JSON into DTOs and pass them to the service layer.

---

## 2. Validation Architecture

### Transport Validation vs Business Validation
- **Transport Validation (Handler/DTO level)**:
  - Validates JSON shape, field types, required fields, and basic length/range constraints.
  - Executed via Gin binding tags or `pkg/validator`.
- **Business Validation (Service layer)**:
  - Validates business rules (e.g., poll expiration, unique voter IP/fingerprint, user role permissions, subscription tier limits).
  - Executed inside `internal/service/`.
- **Frontend Validation**:
  - React/Zod validation is purely for user experience (instant feedback).
  - **Backend validation is authoritative for security**. The backend NEVER assumes incoming requests are valid.

---

## 3. Response & Error Handling Standard

### Response Envelope (`pkg/response`)
All endpoints return responses in a standard JSON envelope:
```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "meta": null
}
```

### Domain Errors (`pkg/errors`)
- Service layers return domain errors (`*errors.AppError`) using type constants (`ErrorTypeValidation`, `ErrorTypeNotFound`, `ErrorTypeUnauthorized`, `ErrorTypeForbidden`, `ErrorTypeConflict`, `ErrorTypeBadRequest`, `ErrorTypeInternal`).
- Handlers invoke `errors.RespondWithError(c, err)` to convert domain errors into HTTP status codes and `pkg/response` envelopes.
- Internal errors (stack traces, MongoDB connection errors, SQL/driver details) are **never** exposed to clients.

---

## 4. Configuration & Security Baseline

- **Zero Hardcoded Secrets**: All configuration is loaded from environment variables using `internal/config/`.
- Connection strings for **MongoDB Atlas** (`MONGO_URI`), **Redis Cloud** (`REDIS_ADDR`), and **Aiven Kafka** (`KAFKA_BROKERS`) are passed via environment variables.
- `.env.example` contains safe placeholders only.
- CORS policy is configured at the API Gateway (`services/api-gateway/internal/middleware/cors.go`).

---

## 5. Inter-Service Communication Boundaries

| Transport | Technology | Use Case |
|---|---|---|
| **Synchronous HTTP** | Gin / net/http | API Gateway → Microservices for immediate request/response. |
| **Asynchronous Domain Events** | Aiven Kafka | Cross-service state synchronization (e.g. `VoteCreated` → Analytics). |
| **Realtime Pub/Sub** | Redis Cloud | Low-latency live vote distribution (Vote Service → Realtime Service). |
| **Browser Live Stream** | Gorilla WebSocket | Realtime Service → Browser client. |

---

## 6. Operational Health & Readiness

Every Go microservice exposes operational endpoints:
- `GET /health`: Liveness check confirming the process is alive (returns HTTP 200).
- `GET /ready`: Readiness check verifying dependency connectivity (MongoDB Atlas, Redis Cloud, Aiven Kafka) before accepting traffic.

---

## 7. Service Isolation Boundaries

- **Auth Service**: Isolates user identity and token issuance. Future Myshri.com authentication integration will live strictly behind `services/auth-service/internal/service/` without altering other microservices.
- **Payment Service**: Isolates plan definitions, checkout stubs, and webhook handlers (`services/payment-service`).
- **Database Isolation**: No microservice may directly query or connect to another service's MongoDB collection.
