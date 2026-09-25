import { View, Text, Pressable, TextInput, ActivityIndicator, Alert } from "react-native";
import { useLocalSearchParams, Link } from "expo-router";
import { FlashList } from "@shopify/flash-list";
import { useState } from "react";
import { useCollection, useCollectionRecipes, useAddToCollection, useRemoveFromCollection } from "@/lib/collections";
import { theme } from "@/lib/theme";
import { RequireAuth } from "@/components/RequireAuth";

export default function CollectionDetail() {
  return (
    <RequireAuth title="Sign in to see collections" subtitle="Collections are private to your account.">
      <CollectionView />
    </RequireAuth>
  );
}

function CollectionView() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { data: col, isLoading } = useCollection(id!);
  const { data: recipes } = useCollectionRecipes(id!);
  const [recipeId, setRecipeId] = useState("");
  const add = useAddToCollection();
  const remove = useRemoveFromCollection();

  if (isLoading) {
    return (
      <View style={{ flex: 1, alignItems: "center", justifyContent: "center", backgroundColor: theme.colors.cream }}>
        <ActivityIndicator color={theme.colors.terracotta} />
      </View>
    );
  }
  if (!col) {
    return (
      <View style={{ flex: 1, alignItems: "center", justifyContent: "center", backgroundColor: theme.colors.cream }}>
        <Text style={{ color: theme.colors.charcoalMuted }}>Collection not found</Text>
      </View>
    );
  }

  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream, padding: 16, paddingTop: 48 }}>
      <Text style={{ fontSize: theme.text.title, fontWeight: "700", color: theme.colors.charcoal }}>{col.title}</Text>
      <Text style={{ color: theme.colors.charcoalMuted, marginTop: 4 }}>{col.is_public ? "Public" : "Private"} • {col.description}</Text>

      <View style={{ flexDirection: "row", gap: 8, marginTop: 12 }}>
        <TextInput
          placeholder="Recipe ID to add"
          value={recipeId}
          onChangeText={setRecipeId}
          style={{ flex: 1, backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 10, color: theme.colors.charcoal }}
        />
        <Pressable
          onPress={() => {
            if (!recipeId.trim()) return Alert.alert("Recipe ID required");
            add.mutate({ id: id!, recipe_id: recipeId.trim() }, { onSuccess: () => setRecipeId("") });
          }}
          style={{ backgroundColor: theme.colors.terracotta, borderRadius: theme.radius.sm, paddingHorizontal: 14, justifyContent: "center" }}
        >
          <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Add</Text>
        </Pressable>
      </View>

      <Text style={{ marginTop: 16, fontWeight: "600", color: theme.colors.charcoal }}>Recipes • {recipes?.length ?? 0}</Text>

      <FlashList
        data={recipes ?? []}
        contentContainerStyle={{ paddingTop: 12 }}
        renderItem={({ item }: any) => (
          <View style={{ backgroundColor: theme.colors.card, borderRadius: theme.radius.md, padding: 14, marginBottom: 10, borderWidth: 1, borderColor: theme.colors.border }}>
            <Link href={`/recipe/${item.id}`} asChild>
              <Pressable>
                <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>{item.title}</Text>
                <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12, marginTop: 2 }}>{item.cuisine} • {item.difficulty}</Text>
              </Pressable>
            </Link>
            <Pressable
              onPress={() => remove.mutate({ id: id!, recipeId: item.id })}
              style={{ marginTop: 8, alignSelf: "flex-start", paddingHorizontal: 10, paddingVertical: 6, borderWidth: 1, borderColor: theme.colors.terracottaDark, borderRadius: theme.radius.sm }}
            >
              <Text style={{ fontSize: 12, color: theme.colors.terracottaDark }}>Remove</Text>
            </Pressable>
          </View>
        )}
      />
    </View>
  );
}
