.PHONY: run build test clean

run: ## Run the API Gateway
	go run ./cmd/main.go

build: ## Build the API Gateway binary
	go build -o bin/api-gateway ./cmd/main.go

test: ## Run tests
	go test ./... -v

clean: ## Clean build artifacts
	rm -rf bin/
