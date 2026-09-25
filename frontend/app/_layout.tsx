import { Stack } from "expo-router";
import { GluestackUIProvider } from "@/components/ui/gluestack-ui-provider";
import "@/global.css";
import { SafeAreaListener } from "react-native-safe-area-context";
import { GestureHandlerRootView } from "react-native-gesture-handler";
import { Uniwind } from "uniwind";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/query";
import { AuthProvider, useAuth } from "@/lib/auth";
import { AuthLoading } from "@/components/RequireAuth";

export default function RootLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <SafeAreaListener
          onChange={({ insets }) => {
            Uniwind.updateInsets(insets);
          }}
        >
          <GestureHandlerRootView style={{ flex: 1 }}>
            <GluestackUIProvider mode="light">
              <Router />
            </GluestackUIProvider>
          </GestureHandlerRootView>
        </SafeAreaListener>
      </AuthProvider>
    </QueryClientProvider>
  );
}

// Hold the navigator until the stored token has been validated, so a returning
// user never sees the signed-out UI flash before hydration completes.
function Router() {
  const { loading } = useAuth();
  if (loading) return <AuthLoading />;
  return (
    <Stack
      screenOptions={{
        headerShown: false,
        contentStyle: { backgroundColor: "#FFFBF5" },
      }}
    >
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
      <Stack.Screen name="(auth)" options={{ headerShown: false }} />
      <Stack.Screen
        name="recipe/[id]"
        options={{ title: "Recipe", presentation: "card" }}
      />
      <Stack.Screen
        name="cook/[id]"
        options={{ title: "Cook Mode", presentation: "fullScreenModal" }}
      />
    </Stack>
  );
}
