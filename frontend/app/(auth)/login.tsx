import { useState } from "react";
import { View, Text, TextInput, Pressable, Alert } from "react-native";
import { useRouter } from "expo-router";
import { api } from "@/lib/api";
import { useAuth, type User } from "@/lib/auth";
import { theme } from "@/lib/theme";

export default function LoginScreen() {
  const router = useRouter();
  const { signIn } = useAuth();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const signInWith = async (path: "/auth/login" | "/auth/signup", title: string) => {
    if (!username.trim() || !password) {
      Alert.alert("Missing fields", "Username and password are required");
      return;
    }
    if (path === "/auth/signup") {
      if (username.trim().length < 3 || username.trim().length > 30) {
        Alert.alert("Invalid username", "Username must be 3-30 characters");
        return;
      }
      if (!/^[A-Za-z0-9_]{3,30}$/.test(username.trim())) {
        Alert.alert("Invalid username", "Use letters, numbers and underscores only");
        return;
      }
      if (password.length < 8 || password.length > 72) {
        Alert.alert("Invalid password", "Password must be 8-72 characters");
        return;
      }
    }
    setLoading(true);
    try {
      const res = await api.post<{ token: string; user: User }>(path, {
        username: username.trim(),
        password,
      });
      await signIn(res.token, res.user);
      router.replace("/(tabs)");
    } catch (e: any) {
      Alert.alert(title, e.message ?? "Unknown error");
    } finally {
      setLoading(false);
    }
  };

  const signInWithPassword = () => signInWith("/auth/login", "Login failed");
  const signUp = () => signInWith("/auth/signup", "Sign up failed");

  return (
    <View
      style={{
        flex: 1,
        backgroundColor: theme.colors.cream,
        padding: theme.spacing.lg,
        justifyContent: "center",
        gap: 12,
      }}
    >
      <Text
        style={{
          fontSize: theme.text.title,
          color: theme.colors.charcoal,
          fontWeight: "700",
          textAlign: "center",
          marginBottom: 8,
        }}
      >
        CookMode
      </Text>
      <Text
        style={{
          color: theme.colors.charcoalMuted,
          textAlign: "center",
          marginBottom: 16,
        }}
      >
        Sign in to save recipes and cook
      </Text>

      <TextInput
        placeholder="Username"
        placeholderTextColor={theme.colors.charcoalMuted}
        autoCapitalize="none"
        value={username}
        onChangeText={setUsername}
        style={{
          backgroundColor: theme.colors.card,
          borderWidth: 1,
          borderColor: theme.colors.border,
          borderRadius: theme.radius.sm,
          padding: 14,
          color: theme.colors.charcoal,
        }}
      />
      <TextInput
        placeholder="Password"
        placeholderTextColor={theme.colors.charcoalMuted}
        secureTextEntry
        value={password}
        onChangeText={setPassword}
        style={{
          backgroundColor: theme.colors.card,
          borderWidth: 1,
          borderColor: theme.colors.border,
          borderRadius: theme.radius.sm,
          padding: 14,
          color: theme.colors.charcoal,
        }}
      />

      <Pressable
        onPress={signInWithPassword}
        disabled={loading}
        style={{
          backgroundColor: theme.colors.terracotta,
          borderRadius: theme.radius.sm,
          padding: 14,
          alignItems: "center",
          marginTop: 8,
          opacity: loading ? 0.6 : 1,
        }}
      >
        <Text style={{ color: theme.colors.card, fontWeight: "600" }}>
          {loading ? "Please wait..." : "Sign in"}
        </Text>
      </Pressable>

      <Pressable onPress={signUp} disabled={loading} style={{ alignItems: "center", padding: 8 }}>
        <Text style={{ color: theme.colors.sage, fontWeight: "500" }}>
          Create account
        </Text>
      </Pressable>

      <Pressable onPress={() => router.replace("/(tabs)")} style={{ alignItems: "center", marginTop: 8 }}>
        <Text style={{ color: theme.colors.charcoalMuted }}>Continue as guest — browsing only</Text>
      </Pressable>
    </View>
  );
}
