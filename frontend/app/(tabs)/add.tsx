import { useState } from "react";
import { View, Text, TextInput, Pressable, ActivityIndicator, Alert, ScrollView } from "react-native";
import { useRouter } from "expo-router";
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { theme } from "@/lib/theme";
import type { Recipe } from "@/lib/recipes";

export default function AddScreen() {
  const router = useRouter();
  const { session } = useAuth();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [cuisine, setCuisine] = useState("");
  const [servings, setServings] = useState("2");
  const [difficulty, setDifficulty] = useState("easy");

  const create = useMutation({
    mutationFn: (body: unknown) => api.post<Recipe>("/recipes", body),
    onSuccess: (data) => router.push(`/recipe/${data.id}`),
    onError: (e: Error) => Alert.alert("Failed", e.message),
  });

  const requestUpload = useMutation({
    mutationFn: (recipe_id: string) =>
      api.post<{ upload_url: string; video_uid: string }>("/recipes/upload-url", { recipe_id }),
    onSuccess: (data) => Alert.alert("Upload URL", data.upload_url),
    onError: (e: Error) => Alert.alert("Upload failed", e.message),
  });

  if (!session) {
    return (
      <View style={{ flex: 1, backgroundColor: theme.colors.cream, alignItems: "center", justifyContent: "center", padding: 24 }}>
        <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>Sign in to create recipes</Text>
      </View>
    );
  }

  const onCreate = () => {
    if (!title.trim()) return Alert.alert("Title required");
    const servingsNum = parseInt(servings, 10) || 2;
    create.mutate({
      title: title.trim(),
      description,
      cuisine,
      servings: servingsNum,
      difficulty,
      ingredients: [],
      steps: [],
      dietary_tags: [],
    });
  };

  return (
    <ScrollView style={{ flex: 1, backgroundColor: theme.colors.cream }} contentContainerStyle={{ padding: 16, paddingTop: 48, gap: 12 }}>
      <Text style={{ fontSize: theme.text.title, fontWeight: "700", color: theme.colors.charcoal }}>Add Recipe</Text>
      <Text style={{ color: theme.colors.charcoalMuted }}>Creates a draft (POST /recipes). Then request a video upload URL.</Text>

      <TextInput
        placeholder="Title *"
        placeholderTextColor={theme.colors.charcoalMuted}
        value={title}
        onChangeText={setTitle}
        style={{ backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 12, color: theme.colors.charcoal }}
      />
      <TextInput
        placeholder="Description"
        placeholderTextColor={theme.colors.charcoalMuted}
        value={description}
        onChangeText={setDescription}
        multiline
        style={{ backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 12, color: theme.colors.charcoal, minHeight: 60 }}
      />
      <TextInput
        placeholder="Cuisine (e.g. Italian)"
        placeholderTextColor={theme.colors.charcoalMuted}
        value={cuisine}
        onChangeText={setCuisine}
        style={{ backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 12, color: theme.colors.charcoal }}
      />
      <View style={{ flexDirection: "row", gap: 8 }}>
        <TextInput
          placeholder="Servings"
          placeholderTextColor={theme.colors.charcoalMuted}
          value={servings}
          onChangeText={setServings}
          keyboardType="number-pad"
          style={{ flex: 1, backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 12, color: theme.colors.charcoal }}
        />
        <TextInput
          placeholder="Difficulty (easy/medium/hard)"
          placeholderTextColor={theme.colors.charcoalMuted}
          value={difficulty}
          onChangeText={setDifficulty}
          style={{ flex: 1, backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 12, color: theme.colors.charcoal }}
        />
      </View>

      <Pressable
        onPress={onCreate}
        disabled={create.isPending}
        style={{ backgroundColor: theme.colors.terracotta, borderRadius: theme.radius.sm, padding: 14, alignItems: "center", opacity: create.isPending ? 0.6 : 1 }}
      >
        {create.isPending ? <ActivityIndicator color={theme.colors.card} /> : <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Create draft</Text>}
      </Pressable>

      {create.data && (
        <View style={{ gap: 8, marginTop: 8, backgroundColor: theme.colors.card, borderRadius: theme.radius.md, padding: 12, borderWidth: 1, borderColor: theme.colors.border }}>
          <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>Draft created: {create.data.title}</Text>
          <Pressable
            onPress={() => requestUpload.mutate(create.data!.id)}
            style={{ backgroundColor: theme.colors.sage, borderRadius: theme.radius.sm, padding: 12, alignItems: "center" }}
          >
            <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Get upload URL (POST /recipes/upload-url)</Text>
          </Pressable>
          {requestUpload.data && <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12 }}>{requestUpload.data.upload_url}</Text>}
          <Pressable onPress={() => create.reset()} style={{ alignItems: "center", padding: 8 }}>
            <Text style={{ color: theme.colors.charcoalMuted }}>Create another</Text>
          </Pressable>
        </View>
      )}
    </ScrollView>
  );
}
