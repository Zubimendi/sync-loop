.PHONY: dev install build down
install:  ## first-time setup
cd api && go mod tidy
cd web && pnpm install
dev:      ## run everything concurrently
@make -j 2 dev-api dev-web
dev-api:
cd api && make dev
dev-web:
cd web && pnpm dev
build:
cd api && make build
cd web && pnpm build
down:
docker compose down
migrate:  ## create goose migrations
cd api && goose create $(name) sql
