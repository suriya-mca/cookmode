import { View, Text, Pressable, ScrollView, ActivityIndicator, Alert } from "react-native";
import { useLocalSearchParams, Link, useRouter } from "expo-router";
import { Image } from "expo-image";
import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRecipe } from "@/lib/recipes";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { theme } from "@/lib/theme";

export default function RecipeDetail() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const qc = useQueryClient();
  const { session } = useAuth();
  const { data: r, isLoading } = useRecipe(id!);
  const [tab, setTab] = useState<"ingredients" | "steps" | "nutrition">("ingredients");

  const save = useMutation({
    mutationFn: () => api.post(`/recipes/${id}/saves`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["recipe", id] }),
  });
  const fork = useMutation({
    mutationFn: () => api.post<{ id: string }>(`/recipes/${id}/fork`),
    onSuccess: (data) => router.push(`/recipe/${data.id}`),
  });
  const publish = useMutation({
    mutationFn: () => api.post(`/recipes/${id}/publish`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["recipe", id] }),
  });

  if (isLoading) {
    return (
      <View style={{ flex: 1, alignItems: "center", justifyContent: "center", backgroundColor: theme.colors.cream }}>
        <ActivityIndicator color={theme.colors.terracotta} />
      </View>
    );
  }
  if (!r) {
    return (
      <View style={{ flex: 1, alignItems: "center", justifyContent: "center", backgroundColor: theme.colors.cream }}>
        <Text style={{ color: theme.colors.charcoalMuted }}>Recipe not found</Text>
      </View>
    );
  }

  const authed = !!session;

  return (
    <ScrollView style={{ flex: 1, backgroundColor: theme.colors.cream }} contentContainerStyle={{ paddingBottom: 24 }}>
      <Image
        source={{ uri: r.video_thumbnail_url || undefined }}
        style={{ width: "100%", height: 220, backgroundColor: theme.colors.creamSoft }}
        contentFit="cover"
      />
      <View style={{ padding: 16, gap: 12 }}>
        <Text style={{ fontSize: theme.text.title, fontWeight: "700", color: theme.colors.charcoal }}>{r.title}</Text>
        <Text style={{ color: theme.colors.charcoalMuted }}>{r.description}</Text>
        <Text style={{ fontSize: 12, color: theme.colors.charcoalMuted }}>
          {r.cuisine} • {r.prep_time_min + r.cook_time_min}m • {r.difficulty} • {r.saves} saves
        </Text>

        <View style={{ flexDirection: "row", gap: 8, marginTop: 8 }}>
          <Pressable
            onPress={() => router.push(`/cook/${id}`)}
            style={{ flex: 1, backgroundColor: theme.colors.terracotta, borderRadius: theme.radius.sm, padding: 14, alignItems: "center" }}
          >
            <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Cook Mode</Text>
          </Pressable>
          <Pressable
            onPress={() => (authed ? save.mutate() : Alert.alert("Sign in required"))}
            style={{ borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, paddingHorizontal: 16, justifyContent: "center" }}
          >
            <Text style={{ color: theme.colors.charcoal }}>Save</Text>
          </Pressable>
          <Pressable
            onPress={() => (authed ? fork.mutate() : Alert.alert("Sign in required"))}
            style={{ borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, paddingHorizontal: 16, justifyContent: "center" }}
          >
            <Text style={{ color: theme.colors.charcoal }}>Fork</Text>
          </Pressable>
        </View>

        {r.user_id === session?.user?.id && r.status !== "published" && (
          <Pressable
            onPress={() => publish.mutate()}
            style={{ backgroundColor: theme.colors.sage, borderRadius: theme.radius.sm, padding: 12, alignItems: "center" }}
          >
            <Text style={{ color: theme.colors.card, fontWeight: "600" }}>
              {publish.isPending ? "Publishing..." : `Publish (${r.status})`}
            </Text>
          </Pressable>
        )}

        <View style={{ flexDirection: "row", gap: 8, marginTop: 8, borderBottomWidth: 1, borderBottomColor: theme.colors.border, paddingBottom: 8 }}>
          {(["ingredients", "steps", "nutrition"] as const).map((t) => (
            <Pressable key={t} onPress={() => setTab(t)} style={{ paddingVertical: 6, borderBottomWidth: tab === t ? 2 : 0, borderBottomColor: theme.colors.terracotta }}>
              <Text style={{ color: tab === t ? theme.colors.terracotta : theme.colors.charcoalMuted, fontWeight: tab === t ? "600" : "400", textTransform: "capitalize" }}>{t}</Text>
            </Pressable>
          ))}
        </View>

        {tab === "ingredients" && (
          <View style={{ gap: 8 }}>
            <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>Ingredients • {r.servings} servings</Text>
            {r.ingredients.map((ing, i) => (
              <View key={i} style={{ flexDirection: "row", gap: 8, paddingVertical: 6, borderBottomWidth: 1, borderBottomColor: theme.colors.borderSoft }}>
                <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>{ing.quantity} {ing.unit}</Text>
                <Text style={{ color: theme.colors.charcoal }}>{ing.name}</Text>
              </View>
            ))}
            {r.ingredients.length === 0 && <Text style={{ color: theme.colors.charcoalMuted }}>No ingredients</Text>}
          </View>
        )}

        {tab === "steps" && (
          <View style={{ gap: 12 }}>
            {r.steps.map((s, i) => (
              <View key={i} style={{ flexDirection: "row", gap: 12 }}>
                <View style={{ width: 28, height: 28, borderRadius: 14, backgroundColor: theme.colors.terracotta, alignItems: "center", justifyContent: "center" }}>
                  <Text style={{ color: theme.colors.card, fontSize: 12, fontWeight: "700" }}>{String(i + 1).padStart(2, "0")}</Text>
                </View>
                <View style={{ flex: 1 }}>
                  <Text style={{ color: theme.colors.charcoal }}>{s.text}</Text>
                  <Text style={{ color: theme.colors.sage, fontSize: 12, marginTop: 4 }}>▶ {s.anchor_seconds}s</Text>
                </View>
              </View>
            ))}
            <Link href={`/cook/${id}`} asChild>
              <Pressable style={{ marginTop: 8, backgroundColor: theme.colors.charcoal, borderRadius: theme.radius.sm, padding: 14, alignItems: "center" }}>
                <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Start Cook Mode</Text>
              </Pressable>
            </Link>
          </View>
        )}

        {tab === "nutrition" && (
          <View style={{ gap: 8 }}>
            <Text style={{ color: theme.colors.charcoalMuted }}>Nutrition from ingredients (Edamam/USDA)</Text>
            <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12 }}>{r.nutrition ? JSON.stringify(r.nutrition) : "Not calculated yet"}</Text>
          </View>
        )}
      </View>
    </ScrollView>
  );
}
