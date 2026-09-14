import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { tokenStore, api } from "./api";

type User = {
  id: string;
  username: string;
  display_name?: string;
};

type AuthState = {
  user: User | null;
  session: { access_token: string } | null;
  loading: boolean;
  signOut: () => Promise<void>;
  refresh: () => Promise<void>;
};

const AuthContext = createContext<AuthState>({
  user: null,
  session: null,
  loading: true,
  signOut: async () => {},
  refresh: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [session, setSession] = useState<{ access_token: string } | null>(null);
  const [loading, setLoading] = useState(true);

  const refresh = async () => {
    const token = await tokenStore.getToken();
    if (!token) {
      setUser(null);
      setSession(null);
      setLoading(false);
      return;
    }
    const storedUser = await tokenStore.getUser<User>();
    if (storedUser) {
      setUser(storedUser);
      setSession({ access_token: token });
    } else {
      try {
        const fetched = await api.get<User>("/users/me");
        setUser(fetched);
        setSession({ access_token: token });
      } catch {
        await tokenStore.clear();
        setUser(null);
        setSession(null);
      }
    }
    setLoading(false);
  };

  useEffect(() => {
    refresh();
  }, []);

  const signOut = async () => {
    await tokenStore.clear();
    setUser(null);
    setSession(null);
  };

  return (
    <AuthContext.Provider value={{ user, session, loading, signOut, refresh }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}

export function useRequireAuth() {
  const { session, loading } = useAuth();
  return { authed: !!session, loading };
}
