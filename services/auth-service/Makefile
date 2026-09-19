.PHONY: run build test clean

run: ## Run the Auth Service
	go run ./cmd/main.go

build: ## Build the Auth Service binary
	go build -o bin/auth-service ./cmd/main.go

test: ## Run tests
	go test ./... -v

clean: ## Clean build artifacts
	rm -rf bin/
