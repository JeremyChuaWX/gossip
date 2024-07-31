include ./.env.dev
export

.PHONY: start
start:
	docker compose --file ./docker/compose.yaml --env-file ./docker/docker.env --profile development --profile server up --build -d

.PHONY: dev
dev:
	docker compose --file ./docker/compose.yaml --env-file ./docker/docker.env --profile development up --build

.PHONY: server
server:
	go run ./cmd/server/main.go

.PHONY: stop
stop:
	docker compose --file ./docker/compose.yaml --env-file ./docker/docker.env --profile development --profile server down --remove-orphans --volumes
	docker image prune --force

.PHONY: clean
clean:
	rm -rf ./docker/pg_data/ ./docker/pg_admin/
