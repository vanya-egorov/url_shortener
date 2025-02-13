run:
	CONFIG_PATH=./config/local.yaml go run ./cmd/url-shortener/main.go
test:
	go test ./internal/...
generate:
	go generate ./...
up:
	docker compose -f docker-compose.yaml up -d

down:
	docker compose -f docker-compose.yaml down
