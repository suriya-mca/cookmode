import { View, Text, Pressable, ActivityIndicator } from "react-native";
import { useLocalSearchParams, Link } from "expo-router";
import { FlashList } from "@shopify/flash-list";
import { useCollection, useCollectionRecipes } from "@/lib/collections";
import { theme } from "@/lib/theme";

export default function CollectionDetail() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { data: col, isLoading } = useCollection(id!);
  const { data: recipes } = useCollectionRecipes(id!);

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

      <Text style={{ marginTop: 16, fontWeight: "600", color: theme.colors.charcoal }}>Recipes • {recipes?.length ?? 0}</Text>

      <FlashList
        data={recipes ?? []}
        contentContainerStyle={{ paddingTop: 12 }}
        renderItem={({ item }: any) => (
          <Link href={`/recipe/${item.id}`} asChild>
            <Pressable style={{ backgroundColor: theme.colors.card, borderRadius: theme.radius.md, padding: 14, marginBottom: 10, borderWidth: 1, borderColor: theme.colors.border }}>
              <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>{item.title}</Text>
              <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12, marginTop: 2 }}>{item.cuisine} • {item.difficulty}</Text>
            </Pressable>
          </Link>
        )}
      />
    </View>
  );
}
