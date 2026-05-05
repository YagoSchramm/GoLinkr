# GoLinkr

GoLinkr is a simple URL shortener built with Go, Gorilla Mux, and PostgreSQL.

## What is working now

- Create a short link from a long URL.
- Redirect by short code.
- PostgreSQL-backed repository.
- JSON request/response DTOs.
- Route registration using the `module` and `router` pattern.
- Local config loading from `.env`.

## Current API

### Create link

- `POST /link`
- `GET /health`
- `GET /health/db`

Request body:

```json
{
  "original_url": "https://google.com"
}
```

Response example:

```json
{
  "id": "b3d4b0f2-4e6a-4f6b-9f7b-1c2d3e4f5a6b",
  "code": "aZ3kQ9",
  "original_url": "https://google.com",
  "created_at": "2026-05-04T21:00:00Z",
  "short_url": "http://localhost:8080/link/aZ3kQ9"
}
```

### Redirect

- `GET /link/{code}`

This returns a `302 Found` redirect to the original URL.

### Health

- `GET /health`

Returns a simple liveness response.

- `GET /health/db`

Returns `200 OK` when the API can ping PostgreSQL and `503 Service Unavailable` when it cannot.

## Project Structure

- `domain/entity` - core entities
- `domain/usecase` - use case contracts
- `domain/usecase/impl` - use case implementations
- `domain/usecase/dtos` - request and response DTOs
- `domain/util` - domain helpers
- `infrastructure/datastore/repository` - repository contracts
- `infrastructure/datastore/repository/impl` - PostgreSQL implementation
- `infrastructure/router` - module/router primitives
- `infrastructure/router/modules` - application modules
- `infrastructure/script/migrate` - SQL bootstrap scripts

## How to run locally

### 1. Start PostgreSQL

The current `docker-compose.yml` starts only the database and mounts the init SQL.

```powershell
docker compose down -v
docker compose up -d db
```

The database is exposed on port `5433`.

### 2. Run the API

The app reads `DATABASE_URL` from `.env`.

```powershell
go run .
```

Current `.env` example:

```env
DATABASE_URL="postgres://golinkr:golinkr@localhost:5433/golinkr?sslmode=disable"
```

### 3. Test with Insomnia

- `POST http://localhost:8080/link`
- `GET http://localhost:8080/link/{code}`
- `GET http://localhost:8080/health`
- `GET http://localhost:8080/health/db`

## Database

The schema is created by:

- `infrastructure/script/migrate/001-create-tables.up.sql`

Tables currently included:

- `users`
- `links`
- `analytics`

## What is still missing

- A proper migration runner in the repo.
- A production Docker image for the API.
- Tests for HTTP routes and use cases.
- Better error mapping and validation responses.
- Pagination and richer CRUD for users and analytics.
- Click tracking for link redirects.

## Notes

- The API currently runs locally with `go run .`.
- The database is managed separately with Docker.
- If you change the database credentials, update `.env` as well.
