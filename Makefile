local.up:
	docker compose up -d

local.down:
	docker compose down

docs.generate:
	swag init -d "./cmd/tuiter,./internal/application/handlers" -o "./cmd/tuiter/docs" --parseInternal --parseDependency