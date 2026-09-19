.PHONY: run build test clean

run: ## Run the Payment Service
	go run ./cmd/main.go

build: ## Build the Payment Service binary
	go build -o bin/payment-service ./cmd/main.go

test: ## Run tests
	go test ./... -v

clean: ## Clean build artifacts
	rm -rf bin/
