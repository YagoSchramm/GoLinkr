# GoLinkr

GoLinkr is a small URL shortener API built with Go, Gorilla Mux, PostgreSQL, and JWT-based authentication.

## Current Features

- User registration and login with password hashing.
- JWT-protected routes for creating links and reading analytics.
- Short link creation from a long URL.
- Public redirect by short code.
- Redirect click tracking through the analytics repository.
- PostgreSQL-backed repositories.
- JSON request and response DTOs.
- Route registration using the `module` and `router` pattern.
- Local configuration loading from `.env`.
- Health checks for the API and database connection.

## API Overview

Public routes do not require authentication. Protected routes require:

```http
Authorization: Bearer <token>
```

### Auth

#### Register

- `POST /auth/register`
- Public

Request body:

```json
{
  "name": "Ada Lovelace",
  "email": "ada@example.com",
  "password": "strong-password"
}
```

Response:

```json
"<jwt-token>"
```

The same token is also returned in the `Authorization` response header.

#### Login

- `POST /auth/login`
- Public

Request body:

```json
{
  "email": "ada@example.com",
  "password": "strong-password"
}
```

Response:

```json
"<jwt-token>"
```

The same token is also returned in the `Authorization` response header.

### Links

#### Create Link

- `POST /link`
- Protected

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

#### Redirect

- `GET /link/{code}`
- Public

Returns a `302 Found` redirect to the original URL and increments the link analytics counter.

### Analytics

#### Get Analytics By Link ID

- `GET /analytics/{link_id}`
- Protected

Response example:

```json
{
  "id": "8c88c9c7-35a3-45c5-84f9-39cf188fa9e7",
  "link_id": "b3d4b0f2-4e6a-4f6b-9f7b-1c2d3e4f5a6b",
  "clicks": 12
}
```

### Health

#### Liveness

- `GET /health`
- Public

Response example:

```json
{
  "code": "SUCCESS",
  "message": "service is healthy"
}
```

#### Database Readiness

- `GET /health/db`
- Public

Returns `200 OK` when the API can ping PostgreSQL:

```json
{
  "code": "SUCCESS",
  "message": "service is ready",
  "database": "up"
}
```

Returns `503 Service Unavailable` when PostgreSQL is not reachable.

## Project Structure

- `domain/entity` - core entities and domain errors
- `domain/rules` - validation rules
- `domain/usecase` - use case contracts
- `domain/usecase/impl` - use case implementations
- `domain/usecase/dtos` - request and response DTOs
- `domain/util` - domain helpers
- `infrastructure/datastore/repository` - repository contracts
- `infrastructure/datastore/repository/impl` - PostgreSQL repository implementations
- `infrastructure/datastore/repository/impl/_query` - SQL query files
- `infrastructure/router` - module/router primitives and HTTP helpers
- `infrastructure/router/modules` - HTTP modules
- `infrastructure/service/db` - PostgreSQL connection setup
- `infrastructure/service/hash` - password hashing service
- `infrastructure/service/jwt` - JWT generation and validation
- `infrastructure/script/migrate` - SQL bootstrap scripts

## How To Run Locally

### 1. Start PostgreSQL

The current `docker-compose.yml` starts only the database and mounts the initial SQL script.

```powershell
docker compose down -v
docker compose up -d db
```

The database is exposed on port `5433`.

### 2. Configure Environment Variables

The app reads configuration from `.env`.

```env
DATABASE_URL="postgres://golinkr:golinkr@localhost:5433/golinkr?sslmode=disable"
JWT_SECRET="change-this-secret"
```

### 3. Run The API

```powershell
go run .
```

The API listens on `http://localhost:8080`.

### 4. Try The Main Flow

1. Register or log in to receive a JWT.
2. Send the JWT as `Authorization: Bearer <token>`.
3. Create a link with `POST http://localhost:8080/link`.
4. Open `GET http://localhost:8080/link/{code}` to redirect.
5. Read analytics with `GET http://localhost:8080/analytics/{link_id}`.

## Database

The schema is created by:

- `infrastructure/script/migrate/001-create-tables.up.sql`

Tables currently included:

- `users`
- `links`
- `analytics`

## Known Gaps

- There is no migration runner yet; Docker mounts the initial SQL file directly.
- The API does not have a production Docker image yet.
- Automated tests are still missing.
- Validation and error responses can still be made more consistent.

## Notes

- The API currently runs locally with `go run .`.
- PostgreSQL is managed separately with Docker.
- If you change database credentials or the JWT secret, update `.env` as well.
