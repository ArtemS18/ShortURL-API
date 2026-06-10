run:
	go run ./cmd/main

build:
	docker compose up -d --build

swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --parseDependency --parseInternal -g app/app.go -d ./internal -o ./docs
