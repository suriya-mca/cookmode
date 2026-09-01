import { useState } from "react";
import { View, Text, TextInput, Pressable, ActivityIndicator, Alert } from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";
import { useQueryClient } from "@tanstack/react-query";
import { useRecipe } from "@/lib/recipes";
import { api } from "@/lib/api";
import { theme } from "@/lib/theme";

export default function EditRecipe() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const qc = useQueryClient();
  const { data: r, isLoading } = useRecipe(id!);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [saving, setSaving] = useState(false);

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
        <Text>Not found</Text>
      </View>
    );
  }

  // init once
  if (title === "" && description === "" && r) {
    if (title === "") setTitle(r.title);
    if (description === "") setDescription(r.description ?? "");
  }

  const onSave = async () => {
    if (!title.trim()) return Alert.alert("Title required");
    setSaving(true);
    try {
      await api.patch(`/recipes/${id}`, { title: title.trim(), description });
      qc.invalidateQueries({ queryKey: ["recipe", id] });
      qc.invalidateQueries({ queryKey: ["recipes"] });
      router.back();
    } catch (e: any) {
      Alert.alert("Failed", e.message);
    } finally {
      setSaving(false);
    }
  };

  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream, padding: 16, paddingTop: 48, gap: 12 }}>
      <Text style={{ fontSize: theme.text.title, fontWeight: "700", color: theme.colors.charcoal }}>Edit Recipe</Text>
      <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12 }}>PATCH /recipes/:id (only supplied keys merged)</Text>
      <TextInput
        value={title}
        onChangeText={setTitle}
        placeholder="Title"
        style={{ backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 12, color: theme.colors.charcoal }}
      />
      <TextInput
        value={description}
        onChangeText={setDescription}
        placeholder="Description"
        multiline
        style={{ backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 12, minHeight: 80, color: theme.colors.charcoal }}
      />
      <Pressable onPress={onSave} disabled={saving} style={{ backgroundColor: theme.colors.terracotta, borderRadius: theme.radius.sm, padding: 14, alignItems: "center", opacity: saving ? 0.6 : 1 }}>
        {saving ? <ActivityIndicator color={theme.colors.card} /> : <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Save (PATCH)</Text>}
      </Pressable>
    </View>
  );
}
