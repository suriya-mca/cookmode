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
| Auth      | Local auth (bcrypt passwords + HS256 JWT, `/auth/signup` + `/auth/login`) |
| Nutrition | Edamam (primary), USDA fallback               |
| App       | Expo `~57` + `expo-router` + `gluestack-ui` + `UniWind` (Tailwind v4) |

## Monorepo layout

```
backend/              Go API (branch: api)
├── cmd/
│   ├── api/          main server (River start + EnsureIndex)
│   └── migrate/      migration runner (embedded SQL)
├── internal/
│   ├── api/          routes.go, handlers/, middleware/ (auth.Optional)
│   ├── config/       env loading (`JWT_SECRET` signs local HS256 auth tokens, min 32 chars)
│   ├── db/           pgx pool + queries/ (recipes/saves/collections/follows/posts/shopping/forks)
│   ├── httpx/        respond/validate/paginate + visibility helper
│   ├── jobs/         River workers (process_video → index_recipe)
│   ├── models/       domain structs (Post.RatingValue *int, Step.anchor_seconds float64)
│   └── services/     video (10s timeout, context), search (mutex+quoteFilterValue), nutrition
├── migrations/       SQL files (.up.sql / .down.sql pairs, 000001-000004 + pg_trgm)
├── pkg/              thin clients (Meilisearch, R2, Stream)
├── docker-compose.yml  local Postgres :5433 + Meilisearch :7700
└── .env.example      template for backend/.env

frontend/             Expo app (branch: ui, pnpm)
├── app/
│   ├── _layout.tsx          Root (Gluestack + SafeArea + Gesture + UniWind)
│   ├── (tabs)/              Feed/Search/Add/Collections/Profile (5-tab)
│   ├── recipe/[id].tsx      Detail (hero video + tabs)
│   └── cook/[id].tsx        Cook Mode (full-screen anchors)
├── components/ui/    gluestack-ui-provider (copy-paste)
├── lib/theme.ts      palette from ui/*.jpeg (cream/terracotta/sage, Playfair/Inter, radius 20)
├── global.css        @import tailwindcss; @import 'uniwind'
├── pnpm-lock.yaml    (do not use npm)
└── .serena/          language_servers: [typescript_vts] (vtsls)

API_CONTRACT.md       REST API spec — shared with frontend team
```

## Local development

Prerequisites: [Docker](https://docs.docker.com/get-docker/), Go 1.27+, Node 24 + `pnpm`.

```sh
# Backend
cd backend
docker compose up -d          # Postgres :5433 + Meilisearch :7700
cp .env.example .env          # dev defaults already point at compose stack
# one-time: generate the JWT signing secret (server refuses to start without it)
sed -i "s|^JWT_SECRET=.*|JWT_SECRET=$(openssl rand -base64 48)|" .env
make migrate-up               # also runs rivermigrate for river_* tables
make run                      # :8080 (logs "listening on :8080")

# Frontend (separate terminal)
cd frontend
pnpm install
pnpm start                    # or pnpm exec expo start -- --web
pnpm exec tsc --noEmit --skipLibCheck  # typecheck
```

Verify:

```sh
curl http://localhost:8080/api/v1/health
# → {"status":"ok"}
curl "http://localhost:8080/api/v1/search?q=garlic"
```

### Auth (local, no external provider)

```sh
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"chef_amy","password":"correct-horse-9"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")
curl http://localhost:8080/api/v1/users/me -H "Authorization: Bearer $TOKEN"
```

Login works the same via `POST /api/v1/auth/login`. Tokens are HS256, 7-day
expiry, no refresh — re-login when they expire. Usernames are lowercase
`a-z0-9_`, 3-30 chars; passwords 8-72 **bytes** (bcrypt's input limit).

### Publishing and search

Recipes are created as `draft`. Publish to make them visible in `GET /recipes` and search:

```sh
curl -X POST http://localhost:8080/api/v1/recipes/<id>/publish \
  -H "Authorization: Bearer <token>"
```

Publishing also indexes the recipe in Meilisearch (`GET /search?q=...`). Video upload (`POST /recipes/upload-url` → `https://dev.local/upload/{uuid}` stub when `CF_*` empty) sets `status=processing`; the publish step is manual until the Stream → transcribe → anchors pipeline is fully wired.

Notes:
- CookMode's Postgres maps to host port **5433** (not the default 5432) to
  avoid conflicts with other local projects.
- Stop everything with `docker compose down` (add `-v` to wipe data).
- API endpoints can be tested interactively with
  [Hoppscotch](https://hoppscotch.io) or any REST client.

## Testing

```sh
cd backend
go test ./...                          # unit tests; DB-backed auth tests skip unless DATABASE_URL is set
DATABASE_URL=postgres://cookmode:cookmode@localhost:5433/cookmode?sslmode=disable \
  go test ./...                        # includes the DB-backed auth handler tests
go test ./internal/httpx -run TestValidateRecipe
cd ../frontend
pnpm exec tsc --noEmit --skipLibCheck
```

## Commands

| Command             | Where     | Action                                    |
|---------------------|-----------|-------------------------------------------|
| `make run`          | `backend` | Start the API server on :8080             |
| `make test`         | `backend` | `go test ./...`                           |
| `make build`        | `backend` | Build binary to `bin/api`                 |
| `make migrate-up`   | `backend` | Apply pending migrations                  |
| `make migrate-down` | `backend` | Roll back last migration                  |
| `make tidy`         | `backend` | Run `go mod tidy`                         |
| `pnpm start`        | `frontend`| Expo dev server                           |
| `pnpm exec tsc`     | `frontend`| Typecheck                                 |

Single test: `go test ./internal/<pkg> -run TestName`

## Branches

| Branch | Purpose                    |
|--------|----------------------------|
| `main` | Production                 |
| `dev`  | Integration                |
| `api`  | Backend work               |
| `ui`   | Frontend work              |
```

Work on `ui` for frontend (`pnpm`), `api` for backend (`go`). Both PR into `dev`, `dev` → `main` for releases:
```sh
git switch ui; git add frontend/; git commit -m "feat(ui): ..."; git push origin ui   # → PR ui → dev
git switch api; git merge origin/dev; git add backend/; git commit; git push origin api # → PR api → dev
```
