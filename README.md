# CookMode

Video-first cooking platform: structured recipe data, AI-inferred video anchor
points, cook mode voice control, nutrition from ingredients, "I made this"
posts, shopping lists, and ingredient-aware search.

## Layout (monorepo)

```
backend/     Go API — owned by @suriya (branch: api)
  cmd/api/       main server
  cmd/migrate/   migration runner (embedded SQL)
  internal/      api/, db/, jobs/, models/, services/, config/
  migrations/    golang-migrate SQL files
  pkg/           thin clients (Meilisearch, R2, Cloudflare Stream)
frontend/    React Native + Expo app — owned by UI team (branch: ui)
API_CONTRACT.md   REST API spec shared between backend and frontend
```

## Backend quick start

```sh
cd backend
cp .env.example .env          # fill in real values
go mod tidy
make migrate-up
make run                      # serves on :8080
```

Commands: `make run | test | build | migrate-up | migrate-down`
Single test: `go test ./internal/services -run TestName`

## Stack

Go + Gin · PostgreSQL (JSONB) · River queue (Postgres-backed) · Meilisearch ·
Cloudflare R2 (photos) · Cloudflare Stream (video) · Supabase Auth ·
Edamam/USDA (nutrition)

## Branches

`main` = production · `dev` = integration · `api` = backend · `ui` = frontend
