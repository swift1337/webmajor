build:
	go build -o bin/service cmd/service.go

run: build
	./bin/service -service-port 80 -source http://localhost:8000

lint: ## Run golangci-lint
	golangci-lint run -v ./...