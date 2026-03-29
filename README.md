# Movies API (Go + MongoDB)

JSON REST service for **movies** (title, rating, description, IMDb URL, YouTube trailer URL) with **cover art** uploads stored on disk and served under `/api/v1/files/covers/…`.

Author: `nisalk.dev`

## Requirements

- [Go](https://go.dev/dl/) 1.22 or newer
- A MongoDB instance (local or [Atlas](https://www.mongodb.com/cloud/atlas)) and a connection URI

## Configuration

1. Copy `.env.example` to `.env`.
2. Set **`MONGODB_URI`** (required at startup).
3. Set **`API_TOKEN`** — **every** route requires `Authorization: Bearer <token>` or `X-API-Key: <token>` (including **`/health`** and cover images). Plain `<img src="...">` without credentials will not work unless the client sends the header (e.g. fetch with auth, or a front-end proxy).
4. Optional: `MONGODB_DATABASE`, `MONGODB_COLLECTION_MOVIES`, **`UPLOAD_DIR`** (default `data/covers`), **`MAX_UPLOAD_MB`** (default `10`), `HTTP_ADDR` (default `:8080`).

The app creates `UPLOAD_DIR` on startup. Environment variables can also be set without a `.env` file.

### API token (`API_TOKEN`)

The server does not issue tokens. **`API_TOKEN` is a shared secret you create** and paste into `.env`; clients must send the **exact same** value as `Authorization: Bearer <token>` or `X-API-Key: <token>`.

Generate a random secret (pick one):

```bash
openssl rand -hex 32
```

```bash
node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"
```

On **PowerShell** (.NET):

```powershell
$b = New-Object byte[] 32; [System.Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($b); [BitConverter]::ToString($b).Replace('-','').ToLower()
```

Then set `API_TOKEN=<output>` and use that same string in Postman’s **`apiToken`** variable.

## Run

```bash
go run ./cmd/server
```

```bash
go build -o bin/server ./cmd/server
```

On Windows: `bin\server.exe` if you build with `-o bin/server.exe`.

## Development (auto-restart)

[Air](https://github.com/air-verse/air) rebuilds and restarts the server when you change Go files.

```bash
go install github.com/air-verse/air@latest
```

From the project root (with `.env` configured):

```bash
air
```

Settings live in **`.air.toml`** (build output under `tmp/`, ignored by git). On Windows, ensure `air` is on your `PATH` (often `%USERPROFILE%\go\bin`).

### Formatting (Prettier)

JSON, Markdown, and YAML in this repo are formatted with [Prettier](https://prettier.io/) (`.prettierrc`). Go source uses `gofmt` / your editor’s Go formatter, not Prettier.

```bash
npm install
npm run format
```

`npm run format:check` exits with an error if anything would change (useful in CI).

## Docker

```bash
docker build -t restapi .
docker run --rm -p 8080:8080 \
  -e MONGODB_URI="your-connection-string" \
  -e API_TOKEN="your-secret-token" \
  -e UPLOAD_DIR=/data/covers \
  -v restapi_covers:/data/covers \
  restapi
```

Compose loads variables from `.env` (including **`API_TOKEN`**) and persists uploads on the `movie_uploads` volume:

```bash
export MONGODB_URI="your-connection-string"
export API_TOKEN="your-secret-token"
docker compose up --build
```

Local MongoDB stack:

```bash
docker compose -f docker-compose.yml -f compose.mongo.yaml up --build
```

Host port: `HTTP_PORT` (defaults to `8080`).

## API

| Method   | Path                              | Description                                             |
| -------- | --------------------------------- | ------------------------------------------------------- |
| `GET`    | `/health`                         | Liveness                                                |
| `GET`    | `/api/v1/movies`                  | List movies                                             |
| `POST`   | `/api/v1/movies`                  | Create (JSON **or** `multipart/form-data`; see below)   |
| `GET`    | `/api/v1/movies/{id}`             | Get one                                                 |
| `PATCH`  | `/api/v1/movies/{id}`             | Partial update (JSON)                                   |
| `DELETE` | `/api/v1/movies/{id}`             | Delete (removes cover file if present)                  |
| `POST`   | `/api/v1/movies/{id}/cover`       | Upload or replace cover (`multipart` field **`cover`**) |
| `GET`    | `/api/v1/files/covers/{filename}` | Download a stored cover image (requires token)          |

Send **`Authorization: Bearer <token>`** or **`X-API-Key: <token>`** on **every** request (same value as **`API_TOKEN`** in the server env).

**Create JSON** (`Content-Type: application/json`): `title` (required), `rate` (0–10), `description`, `imdbLink`, `trailerYouTubeLink` (URLs must use `http://` or `https://` when set).

**Create multipart** (`Content-Type: multipart/form-data`): fields `title`, optional `rate`, `description`, `imdbLink`, `trailerYouTubeLink`, optional file field **`cover`** (jpeg/png/webp/gif).

Responses include a computed **`coverArtURL`** path (e.g. `/api/v1/files/covers/<filename>.jpg`) when a cover exists.

## Postman

Import **`Movies_API.postman_collection.json`** from the project root (Postman → Import → File), or use the team workspace link in [docs/postman.md](docs/postman.md).

## Project layout

```
cmd/server/           # Entrypoint
internal/config/      # Configuration
internal/models/      # Movie types
internal/repository/  # MongoDB
internal/storage/     # Cover file storage
internal/handlers/    # HTTP handlers
internal/auth/        # API token middleware
internal/router/      # Routes and static cover files
```

## License

[MIT](LICENSE)
