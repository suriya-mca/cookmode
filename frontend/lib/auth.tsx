import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { tokenStore, api, onUnauthorized } from "./api";
import { queryClient } from "./query";

export type User = {
  id: string;
  username: string;
  display_name?: string;
  avatar_url?: string;
  bio?: string;
};

type AuthState = {
  user: User | null;
  signedIn: boolean;
  loading: boolean;
  signIn: (token: string, user: User) => Promise<void>;
  signOut: () => Promise<void>;
  refresh: () => Promise<void>;
};

const AuthContext = createContext<AuthState>({
  user: null,
  signedIn: false,
  loading: true,
  signIn: async () => {},
  signOut: async () => {},
  refresh: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // Always validates the stored token against /users/me: the backend rejects
  // expired tokens (7-day lifetime) and now also rejects malformed ones, so a
  // cached user object cannot be trusted on its own.
  const refresh = useCallback(async () => {
    try {
      const token = await tokenStore.getToken();
      if (!token) {
        setUser(null);
        return;
      }
      const me = await api.get<User>("/users/me");
      setUser(me);
      await tokenStore.setToken(token, me);
    } catch {
      await tokenStore.clear();
      setUser(null);
    } finally {
      setLoading(false);
    }
  }, []);

  const signIn = useCallback(async (token: string, me: User) => {
    await tokenStore.setToken(token, me);
    // Drop anything cached for a previous account.
    queryClient.clear();
    setUser(me);
    setLoading(false);
  }, []);

  const signOut = useCallback(async () => {
    await tokenStore.clear();
    queryClient.clear();
    setUser(null);
  }, []);

  // Registered before hydration so a rejected token on the very first
  // /users/me call is still observed.
  useEffect(
    () =>
      onUnauthorized(() => {
        queryClient.clear();
        setUser(null);
      }),
    []
  );

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return (
    <AuthContext.Provider
      value={{ user, signedIn: !!user, loading, signIn, signOut, refresh }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}

export function useRequireAuth() {
  const { signedIn, loading } = useAuth();
  return { authed: signedIn, loading };
}
