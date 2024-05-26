.PHONY: start
start:
	docker compose --file ./docker/compose.yaml --env-file ./docker/docker.env up --build -d

.PHONY: stop
stop:
	docker compose --file ./docker/compose.yaml --env-file ./docker/docker.env down --remove-orphans --volumes && docker image prune

.PHONY: clean
clean:
	rm -rf ./docker/pg_data/ ./docker/pg_admin/
