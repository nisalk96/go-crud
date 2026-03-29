# Postman collection

This API is covered by a shared Postman collection on the team workspace. Use the link below to join the workspace and open the collection (you need a Postman account).

**Join the Postman team / workspace**

https://app.getpostman.com/join-team?invite_code=5411d9453b92d79ca920ef1171176d05e83c0a37ec3a9753ee68ee88f87a808a&target_code=c0a24692ac36d14f700e24f619f1f19d

After you join, set the collection or environment **base URL** to match your server (for example `http://localhost:8080` if you use the default `HTTP_ADDR`).

## Endpoints (reference)

| Method | Path | Notes |
|--------|------|--------|
| GET | `/health` | Health check |
| GET | `/api/v1/movies` | List movies |
| POST | `/api/v1/movies` | Create: JSON body **or** multipart with `title`, optional `rate`, `description`, `imdbLink`, `trailerYouTubeLink`, file `cover` |
| GET | `/api/v1/movies/{id}` | Get by MongoDB ObjectID hex |
| PATCH | `/api/v1/movies/{id}` | JSON partial update; set `"coverArt": ""` to clear cover (deletes file) |
| DELETE | `/api/v1/movies/{id}` | Delete movie |
| POST | `/api/v1/movies/{id}/cover` | Multipart field `cover` — upload or replace poster |
| GET | `/api/v1/files/covers/{filename}` | Serve stored cover image |

Ensure `MONGODB_URI` is set (see `.env.example`) before running the server locally.
