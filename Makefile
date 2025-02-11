run:
	CONFIG_PATH=./config/local.yaml go run ./cmd/url-shortener/main.go
test:
	go test ./internal/...
generate:
	go generate ./...
functional_tests:
	go test ./tests/url_shortener_test.go
up:
	docker compose -f docker-compose.yaml up -d

down:
	docker compose -f docker-compose.yaml down
	