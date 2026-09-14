import * as SecureStore from "expo-secure-store";

const API_BASE =
  process.env.EXPO_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

const TOKEN_KEY = "cookmode_token";
const USER_KEY = "cookmode_user";

let cachedToken: string | null = null;

export const tokenStore = {
  async getToken(): Promise<string | null> {
    if (cachedToken !== null) return cachedToken;
    cachedToken = await SecureStore.getItemAsync(TOKEN_KEY);
    return cachedToken;
  },
  async setToken(token: string, user?: unknown) {
    cachedToken = token;
    await SecureStore.setItemAsync(TOKEN_KEY, token);
    if (user !== undefined) {
      await SecureStore.setItemAsync(USER_KEY, JSON.stringify(user));
    }
  },
  async clear() {
    cachedToken = null;
    await SecureStore.deleteItemAsync(TOKEN_KEY);
    await SecureStore.deleteItemAsync(USER_KEY);
  },
  async getUser<T>(): Promise<T | null> {
    const raw = await SecureStore.getItemAsync(USER_KEY);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as T;
    } catch {
      return null;
    }
  },
};

async function authHeaders(): Promise<Record<string, string>> {
  const token = await tokenStore.getToken();
  return token ? { Authorization: `Bearer ${token}` } : {};
}

export async function apiFetch<T>(
  path: string,
  opts: RequestInit = {}
): Promise<T> {
  const headers = {
    "Content-Type": "application/json",
    ...(opts.headers as Record<string, string> | undefined),
    ...(await authHeaders()),
  };
  const res = await fetch(`${API_BASE}${path}`, { ...opts, headers });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error ?? `Request failed: ${res.status}`);
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
