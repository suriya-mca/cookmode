# CookMode API Contract

Base URL: `/api/v1`
Auth: Supabase JWT (HS256) — `Authorization: Bearer <token>`. Endpoints marked 🔒 require it; the token's `sub` claim identifies the user.

Status codes: `200` success, `201` created, `400` validation error, `401` unauthorized, `403/404` ownership/not found, `501` not yet implemented.

## Health

### `GET /health` — public
Returns `{ "status": "ok" }`.

## Recipes

Recipe object shape mirrors `backend/internal/models/recipe.go`:
`id, user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags[], ingredients[], steps[] (with anchor_seconds), substitutions[], nutrition, video_uid, video_hls_url, video_thumbnail_url, video_duration_sec, views, saves, status (draft|processing|published|archived), timestamps`

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/recipes` | public | List published recipes (paginated: `?limit=&cursor=`) |
| GET | `/recipes/:id` | public | Single recipe |
| POST | `/recipes` | 🔒 | Create draft. Response includes Cloudflare Stream direct-upload URL for the video |
| PATCH | `/recipes/:id` | 🔒 | Update own recipe |
| DELETE | `/recipes/:id` | 🔒 | Archive own recipe (`status → archived`) |
| POST | `/recipes/:id/saves` | 🔒 | Save/bookmark recipe |
| POST | `/recipes/upload-url` | 🔒 | Get a fresh Stream direct-upload URL |

Async pipeline after upload completes:
`process_video → transcribe → infer_anchors → fetch_nutrition + index_recipe`.
While processing, `status = "processing"` and `steps[].anchor_seconds` may be absent.

## Search

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/search?q=` | public | Meilisearch-backed. Filters: `?ingredients=a,b`, `?max_total_time=`, `?dietary=vegan,gluten-free`, `?difficulty=` |

## Posts ("I made this")

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/posts?recipe_id=` or `?user_id=` | 🔒 | List posts |
| POST | `/posts` | 🔒 | Body: `{ recipe_id, photo_url (R2), caption?, rating_value? }` |
| DELETE | `/posts/:id` | 🔒 | Delete own post |

## Shopping lists

List items: `[ { ingredient_name, quantity, unit, checked } ]`.

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/shopping-lists` | 🔒 | Own lists |
| POST | `/shopping-lists` | 🔒 | Create from recipes: `{ title?, recipe_ids[], servings }` — quantities scale to servings |
| POST | `/shopping-lists/:id/items` | 🔒 | Append items |
| PATCH | `/shopping-lists/:id/items/:itemIndex` | 🔒 | Toggle `checked` |

## Error shape

```json
{ "error": "human-readable message" }
```
