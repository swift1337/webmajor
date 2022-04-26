build:
	go build -o bin/service cmd/service.go

run: build
	./bin/service -service-port 80 -proxy-port 8001