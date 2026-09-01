import { useState } from "react";
import { View, Text, TextInput, Pressable, Alert } from "react-native";
import { useRouter } from "expo-router";
import { supabase } from "@/lib/api";
import { theme } from "@/lib/theme";

export default function LoginScreen() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const signIn = async () => {
    setLoading(true);
    const { error } = await supabase.auth.signInWithPassword({
      email: email.trim(),
      password,
    });
    setLoading(false);
    if (error) {
      Alert.alert("Login failed", error.message);
      return;
    }
    router.replace("/(tabs)");
  };

  const signUp = async () => {
    setLoading(true);
    const { error } = await supabase.auth.signUp({
      email: email.trim(),
      password,
    });
    setLoading(false);
    if (error) {
      Alert.alert("Sign up failed", error.message);
      return;
    }
    Alert.alert("Check email", "Confirm your email to continue.");
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
        placeholder="Email"
        placeholderTextColor={theme.colors.charcoalMuted}
        autoCapitalize="none"
        keyboardType="email-address"
        value={email}
        onChangeText={setEmail}
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
