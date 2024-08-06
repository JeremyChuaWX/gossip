.PHONY: start
start:
	docker compose --file ./docker/compose.yaml --env-file ./docker/docker.env --profile "*" up --build

.PHONY: dev
dev:
	docker compose --file ./docker/compose.yaml --env-file ./docker/docker.env --profile development up --build

.PHONY: server
server:
	export POSTGRES_URL=postgresql://admin:password123@127.0.0.1:5432/my_db?sslmode=disable; \
	export SERVER_ADDRESS=127.0.0.1:3000; \
	go run ./cmd/server/main.go

.PHONY: stop
stop:
	docker compose --file ./docker/compose.yaml --env-file ./docker/docker.env --profile "*" down --remove-orphans --volumes
	docker image prune --force

.PHONY: clean
clean:
	rm -rf ./docker/pg_data/ ./docker/pg_admin/

.PHONY: migrate
migrate:
	export $$POSTGRES_URL; \
	docker compose --file ./docker/compose.yaml up --no-deps postgres.migrate
