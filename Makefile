client:
	@go build -o bin/client cmd/client/main.go
	@./bin/client

server:
	@go build -o bin/server cmd/server/main.go
	@./bin/server