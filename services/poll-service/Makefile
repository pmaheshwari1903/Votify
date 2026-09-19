.PHONY: run build test clean

run: ## Run the Poll Service
	go run ./cmd/main.go

build: ## Build the Poll Service binary
	go build -o bin/poll-service ./cmd/main.go

test: ## Run tests
	go test ./... -v

clean: ## Clean build artifacts
	rm -rf bin/
