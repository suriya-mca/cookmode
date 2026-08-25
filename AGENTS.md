# AGENTS.md

## Repo layout & branches

- Monorepo: `backend/` (Go API), `frontend/` (Expo app, separate owner), `API_CONTRACT.md` (REST spec — keep in sync when changing endpoints).
- Branches: `main` = production, `dev` = integration, `api` = backend work, `ui` = frontend work. Backend agents commit to `api`.
- All backend commands run from `backend/`.

## Backend commands

- `make run` — dev server on :8080
- `go test ./internal/<pkg> -run TestName` — single test
- `make migrate-up` / `make migrate-down` — runs `cmd/migrate` (embedded SQL from `backend/migrations/*.sql`)
- `go mod tidy` after changing imports; deps are not vendored
- New migrations need a matching `.up.sql`/`.down.sql` pair; the embed glob (`migrations/embed.go`) picks them up automatically

## Architecture

- Stack: Go + Gin, PostgreSQL, River queue (Postgres-backed — **no Redis**), Meilisearch, Cloudflare R2 (photos), Cloudflare Stream (video), Supabase Auth.
- Routes are registered in `internal/api/routes.go`; handlers in `internal/api/handlers/`, all under `/api/v1`.
- Auth is Supabase JWT (HS256) via `internal/api/middleware/auth.go`; handlers read `user_id` with `middleware.UserIDFrom(c)`.
- One pgx pool (`internal/db`) backs both queries and the River client. Initialize River via `jobs.NewClient(pool)` and `river.Start` it in `cmd/api/main.go`.
- Recipes are JSONB-first: `ingredients`, `steps`, `substitutions`, `nutrition` columns mirror structs in `internal/models/recipe.go`. Keep model tags and `migrations/000001_init.up.sql` in sync when changing fields.

## Async pipeline (River)

Upload flow enqueues jobs in strict order:
`process_video → transcribe → infer_anchors → fetch_nutrition + index_recipe`

- Video is never self-hosted: use Cloudflare Stream direct-upload URLs; HLS URL/thumbnail/duration come from Stream after processing.
- Recipe status lifecycle: `draft → processing → published → archived`. Only published recipes go into Meilisearch.
- Nutrition comes from Edamam (primary) / USDA API fallback by ingredient list — never computer vision.
- Anchor points (`steps[].anchor_seconds`) are AI-inferred by aligning step text with transcription segments; they power cook mode ("next step", "repeat", "set timer").

## Frontend

- `frontend/` is Expo/React Native owned by another person. Don't modify it from backend work; coordinate endpoint changes through `API_CONTRACT.md`.
