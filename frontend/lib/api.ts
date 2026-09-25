import { secureStorage } from "./storage";

const API_BASE =
  process.env.EXPO_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

const TOKEN_KEY = "cookmode_token";
const USER_KEY = "cookmode_user";

let cachedToken: string | null = null;

export const tokenStore = {
  async getToken(): Promise<string | null> {
    if (cachedToken !== null) return cachedToken;
    cachedToken = await secureStorage.getItem(TOKEN_KEY);
    return cachedToken;
  },
  async setToken(token: string, user?: unknown) {
    cachedToken = token;
    await secureStorage.setItem(TOKEN_KEY, token);
    if (user !== undefined) {
      await secureStorage.setItem(USER_KEY, JSON.stringify(user));
    }
  },
  async clear() {
    cachedToken = null;
    await secureStorage.removeItem(TOKEN_KEY);
    await secureStorage.removeItem(USER_KEY);
  },
  async getUser<T>(): Promise<T | null> {
    const raw = await secureStorage.getItem(USER_KEY);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as T;
    } catch {
      return null;
    }
  },
};

/** Error carrying the HTTP status so callers can branch on 401/403/404. */
export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

type UnauthorizedListener = () => void;
const unauthorizedListeners = new Set<UnauthorizedListener>();

/**
 * Subscribe to "the token we sent was rejected" events. The auth provider
 * uses this to drop the session; api.ts stays free of imports from it.
 * Returns an unsubscribe function.
 */
export function onUnauthorized(listener: UnauthorizedListener): () => void {
  unauthorizedListeners.add(listener);
  return () => {
    unauthorizedListeners.delete(listener);
  };
}

function notifyUnauthorized() {
  for (const listener of unauthorizedListeners) listener();
}

// A 401 from the auth endpoints is a credential error, not a stale session.
const AUTH_PATHS = ["/auth/login", "/auth/signup"];

async function authHeaders(): Promise<Record<string, string>> {
  const token = await tokenStore.getToken();
  return token ? { Authorization: `Bearer ${token}` } : {};
}

export async function apiFetch<T>(
  path: string,
  opts: RequestInit = {}
): Promise<T> {
  const auth = await authHeaders();
  const headers = {
    "Content-Type": "application/json",
    ...(opts.headers as Record<string, string> | undefined),
    ...auth,
  };
  const res = await fetch(`${API_BASE}${path}`, { ...opts, headers });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    if (res.status === 401 && auth.Authorization && !AUTH_PATHS.includes(path)) {
      // Stale or revoked token: drop it so the UI can fall back to signed-out.
      await tokenStore.clear();
      notifyUnauthorized();
    }
    throw new ApiError(res.status, body.error ?? `Request failed: ${res.status}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export const api = {
  get: <T>(path: string) => apiFetch<T>(path),
  post: <T>(path: string, body?: unknown) =>
    apiFetch<T>(path, { method: "POST", body: JSON.stringify(body ?? {}) }),
  patch: <T>(path: string, body?: unknown) =>
    apiFetch<T>(path, { method: "PATCH", body: JSON.stringify(body) }),
  del: <T>(path: string) => apiFetch<T>(path, { method: "DELETE" }),
};
