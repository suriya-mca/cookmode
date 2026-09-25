# Frontend

React Native + Expo app. **Owned by the UI team** (see branch `ui`).

## Setup

```sh
pnpm install
pnpm start        # Expo dev server
```

## Conventions

- All backend calls go through the REST API documented in `../API_CONTRACT.md`
- Base URL for local dev: `http://localhost:8080/api/v1`
- Auth is local: `POST /auth/login` returns an HS256 JWT, sent as `Authorization: Bearer <token>` (7-day expiry, no refresh)
