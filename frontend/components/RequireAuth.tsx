import type { ReactNode } from "react";
import { View, Text, Pressable, ActivityIndicator } from "react-native";
import { useRouter } from "expo-router";
import { useAuth } from "@/lib/auth";
import { theme } from "@/lib/theme";

/** Centered spinner shown while the stored session is being validated. */
export function AuthLoading() {
  return (
    <View
      style={{
        flex: 1,
        backgroundColor: theme.colors.cream,
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <ActivityIndicator color={theme.colors.terracotta} />
    </View>
  );
}

/** Sign-in prompt with a working link to the login screen. */
export function SignInPrompt({
  title = "Sign in to continue",
  subtitle,
}: {
  title?: string;
  subtitle?: string;
}) {
  const router = useRouter();
  return (
    <View
      style={{
        flex: 1,
        backgroundColor: theme.colors.cream,
        alignItems: "center",
        justifyContent: "center",
        padding: 24,
        gap: 8,
      }}
    >
      <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>{title}</Text>
      {subtitle ? (
        <Text
          style={{
            color: theme.colors.charcoalMuted,
            fontSize: 13,
            textAlign: "center",
          }}
        >
          {subtitle}
        </Text>
      ) : null}
      <Pressable
        onPress={() => router.replace("/(auth)/login")}
        style={{
          marginTop: 8,
          backgroundColor: theme.colors.terracotta,
          paddingHorizontal: 16,
          paddingVertical: 10,
          borderRadius: theme.radius.sm,
        }}
      >
        <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Sign in</Text>
      </Pressable>
    </View>
  );
}

/**
 * Wraps screens that call auth-required endpoints. Hooks inside the screen
 * still run, so the data hooks in lib/* must also gate on `signedIn`; this
 * component only decides what the user sees.
 */
export function RequireAuth({
  children,
  title,
  subtitle,
}: {
  children: ReactNode;
  title?: string;
  subtitle?: string;
}) {
  const { signedIn, loading } = useAuth();
  if (loading) return <AuthLoading />;
  if (!signedIn) return <SignInPrompt title={title} subtitle={subtitle} />;
  return <>{children}</>;
}
