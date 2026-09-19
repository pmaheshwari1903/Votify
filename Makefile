# ============================================
# Votify — Monorepo Makefile
# ============================================

.PHONY: help dev-frontend dev-gateway dev-auth dev-poll dev-vote dev-realtime dev-analytics dev-payment build-all clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ── Frontend ─────────────────────────────────

dev-frontend: ## Start frontend dev server
	cd frontend && npm run dev

build-frontend: ## Build frontend for production
	cd frontend && npm run build

# ── Backend Services ─────────────────────────

dev-gateway: ## Start API Gateway
	cd services/api-gateway && go run ./cmd/main.go

dev-auth: ## Start Auth Service
	cd services/auth-service && go run ./cmd/main.go

dev-poll: ## Start Poll Service
	cd services/poll-service && go run ./cmd/main.go

dev-vote: ## Start Vote Service
	cd services/vote-service && go run ./cmd/main.go

dev-realtime: ## Start Realtime Service
	cd services/realtime-service && go run ./cmd/main.go

dev-analytics: ## Start Analytics Service
	cd services/analytics-service && go run ./cmd/main.go

dev-payment: ## Start Payment Service
	cd services/payment-service && go run ./cmd/main.go

# ── Build ────────────────────────────────────

build-all: ## Build all Go services
	@echo "Building all services..."
	@for dir in services/*/; do \
		echo "Building $$dir..."; \
		cd $$dir && go build -o bin/server ./cmd/main.go && cd ../..; \
	done
	@echo "All services built."

# ── Utilities ────────────────────────────────

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@find . -name "bin" -type d -exec rm -rf {} + 2>/dev/null || true
	@rm -rf frontend/dist
	@echo "Clean complete."

lint: ## Run linters
	@echo "Linting Go services..."
	@for dir in services/*/; do \
		echo "Linting $$dir..."; \
		cd $$dir && go vet ./... && cd ../..; \
	done
