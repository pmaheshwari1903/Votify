.PHONY: run build test clean

run: ## Run the Analytics Service
	go run ./cmd/main.go

build: ## Build the Analytics Service binary
	go build -o bin/analytics-service ./cmd/main.go

test: ## Run tests
	go test ./... -v

clean: ## Clean build artifacts
	rm -rf bin/
