# Frontend

React Native + Expo app. **Owned by the UI team** (see branch `ui`).

## Setup

```sh
npm install
npm start          # Expo dev server
```

## Conventions

- All backend calls go through the REST API documented in `../API_CONTRACT.md`
- Base URL for local dev: `http://localhost:8080/api/v1`
- Auth uses Supabase Auth; send the JWT as `Authorization: Bearer <token>`
