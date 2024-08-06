# Gossip

A chat room application.

## Demo

<video src='https://github.com/user-attachments/assets/2e973162-2a8a-4085-93ad-4e97575e4e66' width=200></video>

## Architecture

![architecture image](https://github.com/user-attachments/assets/296d77b5-6ac6-4bf3-9611-8cd5bb5e73c5)

## How to run

### Dependencies

- Docker
- Go
- Node

### Pre-requisites

Set the correct permissions for the PGAdmin configuration files.

```bash
chmod -R 600 ./docker/pg_config
```

### Full suite

- Spin up PostgreSQL database container
- Run migration scripts against database
- Spin up accompanying PGAdmin container
- Spin up built image of the server

```bash
make start
```

### Development environment

- Spin up PostgreSQL database container
- Run migration scripts against database
- Spin up accompanying PGAdmin container

```bash
make dev
```

- Run the server locally

```bash
make server
```

### Run migration scripts

Run the migration scripts against the PostgreSQL database specified.

```bash
POSTGRES_URL="<postgres-uri-here>" make migrate
```

### Stopping services

Stops any of the services spun up by running `make start` or `make dev`

```bash
make stop
```

### Clean up volumes

Delete persisted volumes for PostgreSQL and PGAdmin containers.

```bash
make clean
```
