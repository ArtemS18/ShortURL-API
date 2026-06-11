run:
	go run ./cmd/main

prod:
	docker compose -f 'docker-compose.prod.yaml' up -d --build

dev:
	docker compose -f 'docker-compose.dev.yaml' up -d --build

swag:
	swag init --parseDependency --parseInternal -g app/app.go -d ./internal -o ./docs

test:
	go clean -testcache
	go test -v ./...

cover:
	go clean -testcache
	-go test ./internal/delivery/... ./internal/repository/... ./internal/usecase/... -coverprofile cover.out -covermode=count 
	go tool cover -func cover.out
