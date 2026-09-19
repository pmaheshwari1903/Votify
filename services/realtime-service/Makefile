.PHONY: run build test clean

run: ## Run the Realtime Service
	go run ./cmd/main.go

build: ## Build the Realtime Service binary
	go build -o bin/realtime-service ./cmd/main.go

test: ## Run tests
	go test ./... -v

clean: ## Clean build artifacts
	rm -rf bin/
