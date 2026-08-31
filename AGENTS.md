# AGENTS.md

## Repo layout & branches

- Monorepo: `backend/` (Go API, `pnpm` not used), `frontend/` (Expo `~52` + `expo-router` + `gluestack-ui` + `UniWind`, `pnpm`), `API_CONTRACT.md` (REST spec — keep in sync when changing endpoints).
- Branches: `main` = production, `dev` = integration, `api` = backend work, `ui` = frontend work. Backend agents commit to `api`, frontend to `ui`. Both PR into `dev`, `dev` → `main` for releases.
- Serena has **two** projects: `backend/.serena/project.yml` (`language_servers: [go]` → `gopls` v0.23) and `frontend/.serena/project.yml` (`language_servers: [typescript_vts]` → `vtsls` via `pnpm`). Keep `ls_workspace_folders` isolated.

## Environment

- Go 1.27 lives in `/usr/local/go/bin` (official tarball, shadows Fedora's package). In fresh/non-login shells: `export PATH=/usr/local/go/bin:$HOME/go/bin:$PATH`. `gopls` v0.23 at `~/go/bin/gopls`.
- Node `24` + `pnpm 11` for frontend. No `npm` — use `pnpm install` / `pnpm exec tsc`.
- Local infra via `docker compose up -d` from `backend/`: Postgres 17 on **host port 5433** (5432 is taken by an unrelated project — never use it), Meilisearch v1.53 on 7700. `backend/.env` already points at this stack.

## Backend commands

- `make run` — dev server on :8080 (from `backend/`)
- `go test ./internal/<pkg> -run TestName` — single test
- `make migrate-up` / `make migrate-down` — runs `cmd/migrate` (embedded SQL from `backend/migrations/*.sql`)
- `go mod tidy` after changing imports; deps are not vendored
- New migrations need a matching `.up.sql`/`.down.sql` pair; the embed glob (`migrations/embed.go`) picks them up automatically

## Frontend commands

- `pnpm install` — install deps (generates `pnpm-lock.yaml`, do not use `npm`)
- `pnpm start` / `pnpm exec expo start -- --web` — Expo dev server (web + native)
- `pnpm exec tsc --noEmit --skipLibCheck` — typecheck (`typescript` 5.9, `moduleResolution: bundler` for `uniwind`)
- `npx gluestack-ui add <component>` — add gluestack copy-paste components to `components/ui/`

## Architecture

- Stack: Go + Gin, PostgreSQL, River queue (Postgres-backed — **no Redis**), Meilisearch, Cloudflare R2 (photos), Cloudflare Stream (video), Supabase Auth; Expo `~52` + `expo-router` + `gluestack-ui` (UniWind/Tailwind v4) + `react-native-reanimated` + `uniwind`.
- Routes are registered in `internal/api/routes.go`; handlers in `internal/api/handlers/`, all under `/api/v1`. `POST /recipes/:id/publish` flips `draft|processing → published` and sync-indexes to Meilisearch.
- Auth is Supabase JWT (HS256) via `internal/api/middleware/auth.go`; handlers read `user_id` with `middleware.UserIDFrom(c)`; `GET /recipes/:id` uses `auth.Optional()` for draft visibility.
- One pgx pool (`internal/db`) backs both queries and the River client. River is started in `cmd/api/main.go` via `jobs.NewClient(pool, indexer)` → `rivermigrate` → `river.Start`; graceful `Stop` on SIGTERM.
- Recipes are JSONB-first: `ingredients`, `steps`, `substitutions`, `nutrition` columns mirror structs in `internal/models/recipe.go`. Keep model tags and `migrations/000001_init.up.sql` in sync when changing fields. `steps[].anchor_seconds` is `float64` (sub-second).
- Frontend theme tokens live in `frontend/lib/theme.ts` (extracted from `ui/*.jpeg`: cream/terracotta/sage palette, `Playfair`/`Inter`, `radius 20`). Gluestack config in `components/ui/gluestack-ui-provider`, `global.css` (`@import tailwindcss; @import 'uniwind'`), `babel/metro` `withUniwindConfig`.

## Async pipeline (River)

Upload flow enqueues jobs in strict order:
`process_video → transcribe → infer_anchors → fetch_nutrition + index_recipe`

- `SearchIndexer` uses `sync.Mutex` + `initialized` (not `sync.Once`) so transient Meilisearch failures retry; `EnsureIndex` is called once at startup in `main.go` and again via `ensureOnce` on first `IndexRecipe`. `Search` no longer calls `EnsureIndex` per request. Filter values are escaped via `quoteFilterValue` (`\`→`\\`, `"`→`\"`).
- Recipe status lifecycle: `draft → processing → published → archived`. Only `published` goes into Meilisearch (via `index_recipe` worker or `POST :id/publish` sync path).
- Video is never self-hosted: use Cloudflare Stream direct-upload URLs; HLS URL/thumbnail/duration come from Stream after processing. Dev stub `https://dev.local/upload/{uuid}` when `CF_*` empty.
- Anchor points (`steps[].anchor_seconds`) are AI-inferred by aligning step text with transcription segments; they power cook mode ("next step", "repeat", "set timer").

## CI

- `.github/workflows/ci.yml` (Go 1.27 + Postgres 17 + Meilisearch + Node 24 via `actions/checkout@v6`/`setup-go@v6`/`setup-node@v6`, `persist-credentials: false`, `permissions: contents: read`) runs `go vet`/`gofmt`/`go test`/`migrate-up`/`build` and `pnpm tsc` on PRs to `dev`/`main`.

## Frontend

- `frontend/` is `expo-router` file-based routing: `app/_layout.tsx` (Gluestack + SafeArea + Gesture + UniWind), `app/(tabs)/[index,search,add,collections,profile]` (5-tab `theme.colors.terracotta` active), `app/recipe/[id].tsx`, `app/cook/[id].tsx`. `package.json` `main: expo-router/entry`, `app.json` `scheme: cookmode` + `plugins: [expo-video, expo-router]`.
- `pnpm` only — don't use `npm` in `frontend/` (lockfile is `pnpm-lock.yaml`). Coordinate endpoint changes through `API_CONTRACT.md`.
