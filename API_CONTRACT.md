# CookMode API Contract

Base URL: `/api/v1`
Auth: Supabase JWT (HS256) — `Authorization: Bearer <token>`. Endpoints marked 🔒 require it; the token's `sub` claim identifies the user.

Status codes: `200` success, `201` created, `204` no content, `400` validation error, `401` unauthorized, `403` forbidden, `404` not found.

Recipe object: `id, user_id, title, description, cuisine, prep_time_min, cook_time_min, servings, difficulty, dietary_tags[], ingredients[{name,quantity,unit,optional}], steps[{order,text,anchor_seconds,duration_hint}], substitutions[{for_ingredient,alternatives[]}], nutrition{calories,protein_g,...}, video_uid, video_hls_url, video_thumbnail_url, video_duration_sec, views, saves, status(draft|processing|published|archived), created_at, updated_at`

## Health

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/health` | public | `{ "status": "ok" }` |

## Recipes

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/recipes` | public | List **published** only. Paginated: `?limit=&cursor=` (cursor is base64 `created_at\|id`) |
| GET | `/recipes/:id` | public* | Single recipe. Draft/processing/archived visible only to owner (otherwise 404). Uses `Auth.Optional()` |
| POST | `/recipes` | 🔒 | Create draft (`status=draft`). Body is `Recipe` JSON without `id/user_id/status/video_*`. Validates title 1-200, servings ≥0, difficulty enum, ingredients/steps |
| PATCH | `/recipes/:id` | 🔒 | Partial update (only supplied JSON keys are merged; omitted keys keep existing). Owner only, archived → 400 |
| DELETE | `/recipes/:id` | 🔒 | Archive own recipe (`status→archived`). Idempotent if already archived |
| POST | `/recipes/:id/saves` | 🔒 | Bookmark recipe (idempotent). Draft saves allowed only for owner; archived → 400; otherwise 404 if not readable |
| DELETE | `/recipes/:id/saves` | 🔒 | Remove bookmark |
| POST | `/recipes/:id/fork` | 🔒 | Clone recipe as new draft owned by caller; inserts `recipe_forks` link. Published always forkable; draft forkable only by owner; archived → 400 |
| POST | `/recipes/upload-url` | 🔒 | Body `{ recipe_id }`. Generates Cloudflare Stream direct-upload URL (dev stub `https://dev.local/upload/{uid}` when `CF_*` empty). Sets `video_uid` + `status=processing`. Owner only, archived → 400 |
| POST | `/recipes/:id/publish` | 🔒 | Publish own recipe (`draft|processing → published`). Triggers Meilisearch indexing. Owner only |

Async pipeline (when `CF_*` configured): `process_video → transcribe → infer_anchors → fetch_nutrition + index_recipe`. While `processing`, `steps[].anchor_seconds` may be 0/absent.

## Search

Meilisearch `recipes` index. Typo-tolerant, filterable on `ingredients`, `dietary_tags`, `difficulty`, `cuisine`, `total_time`.

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/search?q=&ingredients=a,b&dietary=vegan&difficulty=easy&cuisine=Asian&max_total_time=30&limit=20&offset=0` | public | `q` is full-text query (title/description/ingredients/cuisine). Filters are `AND`ed. Returns `{ hits[], query, processingTimeMs, limit, offset, estimatedTotalHits }` |

`ingredients`/`dietary` are comma-separated. `max_total_time` = `prep_time_min+cook_time_min` upper bound.

## Collections (curated lists)

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/collections` | 🔒 | List own collections |
| POST | `/collections` | 🔒 | Body `{ title (1-100), description?, is_public? default true }` |
| GET | `/collections/:id` | 🔒 | Single collection. Private visible only to owner (otherwise 404) |
| PATCH | `/collections/:id` | 🔒 | Owner only. Partial `{ title?, description?, is_public? }` |
| DELETE | `/collections/:id` | 🔒 | Owner only, cascades `collection_recipes` |
| GET | `/collections/:id/recipes` | 🔒 | List recipes in collection (respects collection visibility) |
| POST | `/collections/:id/recipes` | 🔒 | Body `{ recipe_id }`. Owner only, idempotent. Validates recipe exists |
| DELETE | `/collections/:id/recipes/:recipeId` | 🔒 | Owner only |

## Follows

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| POST | `/follows` | 🔒 | Body `{ following_id }`. Idempotent. `400` if self-follow |
| DELETE | `/follows/:id` | 🔒 | Unfollow `following_id` |
| GET | `/users/:id/followers` | public | List follower user IDs |
| GET | `/users/:id/following` | public | List following user IDs |

## Posts ("I made this")

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/posts?recipe_id=&user_id=` | 🔒 | List visible posts. Requires `recipe_id` or `user_id` (400 otherwise) |
| POST | `/posts` | 🔒 | Body `{ recipe_id, photo_url (R2), caption? ≤500, rating_value? 0-5 }`. Validates recipe exists |
| DELETE | `/posts/:id` | 🔒 | Delete own post (owner only, 404 for others) |

## Shopping lists

List shape: `{ id, user_id, title, items[{ingredient_name,quantity,unit,checked}], created_at, updated_at }`

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/shopping-lists` | 🔒 | List own lists |
| GET | `/shopping-lists/:id` | 🔒 | Single list (owner only) |
| POST | `/shopping-lists` | 🔒 | Body `{ title? default "Shopping List", recipe_ids[]?, servings? }`. For each recipe, ingredients are scaled `qty * servings/targetServings` (target defaults to recipe servings) |
| DELETE | `/shopping-lists/:id` | 🔒 | Owner only |
| POST | `/shopping-lists/:id/items` | 🔒 | Body `{ items[{ingredient_name,quantity,unit}] }`. Appends to list |
| PATCH | `/shopping-lists/:id/items/:itemIndex` | 🔒 | Body `{ checked: bool }`. Toggles `checked` at index (400 if out of range) |

## Error shape

```json
{ "error": "human-readable message" }
```

Pagination: `GET /recipes` returns `{ data[], next_cursor }` where `next_cursor` is base64 `created_at|id` or `""` if no more pages.
