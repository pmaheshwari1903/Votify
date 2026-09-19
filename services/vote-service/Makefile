.PHONY: run build test clean

run: ## Run the Vote Service
	go run ./cmd/main.go

build: ## Build the Vote Service binary
	go build -o bin/vote-service ./cmd/main.go

test: ## Run tests
	go test ./... -v

clean: ## Clean build artifacts
	rm -rf bin/
