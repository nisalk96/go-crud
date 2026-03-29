# Movies API (Go + MongoDB)

JSON REST service for **movies** (title, rating, description, IMDb URL, YouTube trailer URL) with **cover art** uploads stored on disk and served under `/api/v1/files/covers/…`.

Author: `nisalk.dev`

## Requirements

- [Go](https://go.dev/dl/) 1.22 or newer  
- A MongoDB instance (local or [Atlas](https://www.mongodb.com/cloud/atlas)) and a connection URI  

## Configuration

1. Copy `.env.example` to `.env`.  
2. Set **`MONGODB_URI`** (required at startup).  
3. Optional: `MONGODB_DATABASE`, `MONGODB_COLLECTION_MOVIES`, **`UPLOAD_DIR`** (default `data/covers`), **`MAX_UPLOAD_MB`** (default `10`), `HTTP_ADDR` (default `:8080`).  

The app creates `UPLOAD_DIR` on startup. Environment variables can also be set without a `.env` file.

## Run

```bash
go run ./cmd/server
```

```bash
go build -o bin/server ./cmd/server
```

On Windows: `bin\server.exe` if you build with `-o bin/server.exe`.

## Docker

```bash
docker build -t restapi .
docker run --rm -p 8080:8080 \
  -e MONGODB_URI="your-connection-string" \
  -e UPLOAD_DIR=/data/covers \
  -v restapi_covers:/data/covers \
  restapi
```

Compose loads variables from `.env` and persists uploads on the `movie_uploads` volume:

```bash
export MONGODB_URI="your-connection-string"
docker compose up --build
```

Local MongoDB stack:

```bash
docker compose -f docker-compose.yml -f compose.mongo.yaml up --build
```

Host port: `HTTP_PORT` (defaults to `8080`).

## API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Liveness |
| `GET` | `/api/v1/movies` | List movies |
| `POST` | `/api/v1/movies` | Create (JSON **or** `multipart/form-data`; see below) |
| `GET` | `/api/v1/movies/{id}` | Get one |
| `PATCH` | `/api/v1/movies/{id}` | Partial update (JSON) |
| `DELETE` | `/api/v1/movies/{id}` | Delete (removes cover file if present) |
| `POST` | `/api/v1/movies/{id}/cover` | Upload or replace cover (`multipart` field **`cover`**) |
| `GET` | `/api/v1/files/covers/{filename}` | Download a stored cover image |

**Create JSON** (`Content-Type: application/json`): `title` (required), `rate` (0–10), `description`, `imdbLink`, `trailerYouTubeLink` (URLs must use `http://` or `https://` when set).

**Create multipart** (`Content-Type: multipart/form-data`): fields `title`, optional `rate`, `description`, `imdbLink`, `trailerYouTubeLink`, optional file field **`cover`** (jpeg/png/webp/gif).

Responses include a computed **`coverArtURL`** path (e.g. `/api/v1/files/covers/<filename>.jpg`) when a cover exists.

## Postman

See [docs/postman.md](docs/postman.md).

## Project layout

```
cmd/server/           # Entrypoint
internal/config/      # Configuration
internal/models/      # Movie types
internal/repository/  # MongoDB
internal/storage/     # Cover file storage
internal/handlers/    # HTTP handlers
internal/router/      # Routes and static cover files
```
