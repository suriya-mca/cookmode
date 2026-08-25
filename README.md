# CookMode

Video-first cooking platform: structured recipe data, AI-inferred video anchor
points, cook mode voice control, nutrition from ingredients, "I made this"
posts, shopping lists, and ingredient-aware search.

## Tech stack

| Layer     | Tech                                          |
|-----------|-----------------------------------------------|
| API       | Go + Gin                                      |
| Database  | PostgreSQL 17 (JSONB-first)                   |
| Queue     | River (Postgres-backed — no Redis)            |
| Search    | Meilisearch                                   |
| Storage   | Cloudflare R2 (photos)                        |
| Video     | Cloudflare Stream                             |
| Auth      | Supabase Auth (JWT)                           |
| Nutrition | Edamam (primary), USDA fallback               |

## Monorepo layout

```
backend/              Go API (branch: api)
├── cmd/
│   ├── api/          main server
│   └── migrate/      migration runner (embedded SQL)
├── internal/
│   ├── api/          routes.go, handlers/, middleware/
│   ├── config/       env loading
│   ├── db/           pgx pool
│   ├── jobs/         River workers (async pipeline)
│   ├── models/       domain structs (mirror migrations)
│   └── services/     business logic (video, nutrition, search)
├── migrations/       SQL files (.up.sql / .down.sql pairs)
├── pkg/              thin clients (Meilisearch, R2, Stream)
├── docker-compose.yml  local Postgres + Meilisearch
└── .env.example      template for backend/.env

frontend/             React Native + Expo app (branch: ui)
API_CONTRACT.md       REST API spec — shared with frontend team
```

## Local development

Prerequisites: [Docker](https://docs.docker.com/get-docker/) and Go 1.27+.

```sh
cd backend

# 1. Start Postgres (:5433) and Meilisearch (:7700)
docker compose up -d

# 2. Configure environment
cp .env.example .env   # dev defaults already point at the compose stack

# 3. Create tables
make migrate-up

# 4. Run the server on :8080
make run
```

Verify:

```sh
curl http://localhost:8080/api/v1/health
# → {"status":"ok"}
```

Notes:
- CookMode's Postgres maps to host port **5433** (not the default 5432) to
  avoid conflicts with other local projects.
- Stop everything with `docker compose down` (add `-v` to wipe data).
- API endpoints can be tested interactively with
  [Hoppscotch](https://hoppscotch.io) or any REST client.

## Commands

| Command             | Action                                    |
|---------------------|-------------------------------------------|
| `make run`          | Start the API server on :8080             |
| `make test`         | Run all tests                             |
| `make build`        | Build binary to `bin/api`                 |
| `make migrate-up`   | Apply pending migrations                  |
| `make migrate-down` | Roll back last migration                  |
| `make tidy`         | Run `go mod tidy`                         |

Single test: `go test ./internal/<pkg> -run TestName`

## Branches

| Branch | Purpose                    |
|--------|----------------------------|
| `main` | Production                 |
| `dev`  | Integration                |
| `api`  | Backend work               |
| `ui`   | Frontend work              |
