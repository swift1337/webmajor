build:
	go build -o bin/service cmd/service.go

run: build
	./bin/service -service-port 8801 -proxy-port 8001

lint: ## Run golangci-lint
	golangci-lint run -v ./...