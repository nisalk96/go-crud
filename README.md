# REST API (Go + MongoDB)

A small JSON REST service for managing **items** with MongoDB persistence, structured as a standard Go layout (`cmd/`, `internal/`).

Author: `nisalk.dev`

## Requirements

- [Go](https://go.dev/dl/) 1.22 or newer  
- A MongoDB instance (local or [Atlas](https://www.mongodb.com/cloud/atlas)) and a connection URI  

## Configuration

1. Copy `.env.example` to `.env`.  
2. Set **`MONGODB_URI`** to your MongoDB connection string (required at startup).  
3. Optional: adjust `MONGODB_DATABASE`, `MONGODB_COLLECTION_ITEMS`, and `HTTP_ADDR` (default `:8080`).  

Variables are loaded from the environment; if present, a `.env` file in the project root is loaded automatically.

## Run

```bash
go run ./cmd/server
```

Build a binary:

```bash
go build -o bin/server ./cmd/server
```

On Windows, you can use `bin\server.exe` after building with `-o bin/server.exe`.

The server performs a MongoDB ping on startup and shuts down gracefully on `SIGINT` / `SIGTERM`.

## Docker

Build and run the API image (set `MONGODB_URI` in your environment or in a `.env` file next to the compose file; Compose reads `.env` for variable substitution):

```bash
docker build -t restapi .
docker run --rm -p 8080:8080 -e MONGODB_URI="your-connection-string" restapi
```

Or with Compose (Atlas or any reachable URI):

```bash
export MONGODB_URI="your-connection-string"
docker compose up --build
```

Optional local MongoDB in Docker (API uses `mongodb://mongo:27017`):

```bash
docker compose -f docker-compose.yml -f compose.mongo.yaml up --build
```

Override the host port with `HTTP_PORT` (maps to container `8080`), for example `HTTP_PORT=3000 docker compose up`.

## API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Liveness |
| `GET` | `/api/v1/items` | List items |
| `POST` | `/api/v1/items` | Create (`name` required; `notes` optional) |
| `GET` | `/api/v1/items/{id}` | Get by MongoDB ObjectID (hex string) |
| `PATCH` | `/api/v1/items/{id}` | Partial update (`name` and/or `notes`) |
| `DELETE` | `/api/v1/items/{id}` | Delete |

Responses are JSON. Create/update bodies use `Content-Type: application/json`.

## Postman

Shared collection and team invite are documented in [docs/postman.md](docs/postman.md).

## Project layout

```
cmd/server/          # Application entrypoint
internal/config/     # Environment configuration
internal/models/     # Domain types
internal/repository/ # MongoDB access
internal/handlers/   # HTTP handlers
internal/router/     # Routes and middleware
```
