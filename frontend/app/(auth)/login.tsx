import { useState } from "react";
import { View, Text, TextInput, Pressable, Alert } from "react-native";
import { useRouter } from "expo-router";
import { api, tokenStore } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { theme } from "@/lib/theme";

export default function LoginScreen() {
  const router = useRouter();
  const { refresh } = useAuth();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const signIn = async () => {
    if (!username.trim() || !password) {
      Alert.alert("Missing fields", "Username and password are required");
      return;
    }
    setLoading(true);
    try {
      const res = await api.post<{ token: string; user: unknown }>("/auth/login", {
        username: username.trim(),
        password,
      });
      await tokenStore.setToken(res.token, res.user);
      await refresh();
      router.replace("/(tabs)");
    } catch (e: any) {
      Alert.alert("Login failed", e.message ?? "Unknown error");
    } finally {
      setLoading(false);
    }
  };

  const signUp = async () => {
    if (!username.trim() || !password) {
      Alert.alert("Missing fields", "Username and password are required");
      return;
    }
    if (username.trim().length < 3 || username.trim().length > 30) {
      Alert.alert("Invalid username", "Username must be 3-30 characters");
      return;
    }
    if (password.length < 8) {
      Alert.alert("Invalid password", "Password must be at least 8 characters");
      return;
    }
    setLoading(true);
    try {
      const res = await api.post<{ token: string; user: unknown }>("/auth/signup", {
        username: username.trim(),
        password,
      });
      await tokenStore.setToken(res.token, res.user);
      await refresh();
      router.replace("/(tabs)");
    } catch (e: any) {
      Alert.alert("Sign up failed", e.message ?? "Unknown error");
    } finally {
      setLoading(false);
    }
  };

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
        onPress={signIn}
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
        <Text style={{ color: theme.colors.charcoalMuted }}>Continue as guest</Text>
      </Pressable>
    </View>
  );
}
