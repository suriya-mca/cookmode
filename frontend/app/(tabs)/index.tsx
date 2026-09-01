import { View, Text, Pressable, ActivityIndicator } from "react-native";
import { FlashList } from "@shopify/flash-list";
import { Image } from "expo-image";
import { Link } from "expo-router";
import { useRecipes, type Recipe } from "@/lib/recipes";
import { theme } from "@/lib/theme";

function RecipeCard({ r }: { r: Recipe }) {
  return (
    <Link href={`/recipe/${r.id}`} asChild>
      <Pressable
        style={{
          backgroundColor: theme.colors.card,
          borderRadius: theme.radius.lg,
          overflow: "hidden",
          marginBottom: 16,
          borderWidth: 1,
          borderColor: theme.colors.border,
          shadowColor: theme.colors.shadow,
          shadowOpacity: 0.08,
          shadowRadius: 20,
        }}
      >
        {r.video_thumbnail_url ? (
          <Image
            source={{ uri: r.video_thumbnail_url }}
            style={{ width: "100%", height: 180, backgroundColor: theme.colors.creamSoft }}
            contentFit="cover"
            transition={200}
          />
        ) : (
          <View
            style={{ width: "100%", height: 180, backgroundColor: theme.colors.creamSoft }}
          />
        )}
        <View style={{ padding: 14, gap: 4 }}>
          <Text numberOfLines={2} style={{ fontSize: 16, fontWeight: "700", color: theme.colors.charcoal }}>
            {r.title}
          </Text>
          <Text style={{ fontSize: 12, color: theme.colors.charcoalMuted }}>
            {r.cuisine} • {r.prep_time_min + r.cook_time_min}m • {r.difficulty} • ★ {r.saves}
          </Text>
        </View>
      </Pressable>
    </Link>
  );
}

export default function FeedScreen() {
  const recipesQuery = useRecipes(20);
  const data: Recipe[] = recipesQuery.data?.pages.flatMap((p) => p.data) ?? [];

  if (recipesQuery.isLoading) {
    return (
      <View style={{ flex: 1, backgroundColor: theme.colors.cream, alignItems: "center", justifyContent: "center" }}>
        <ActivityIndicator color={theme.colors.terracotta} />
      </View>
    );
  }

  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream, paddingTop: 48 }}>
      <View style={{ padding: theme.spacing.md, paddingBottom: 8 }}>
        <Text style={{ fontSize: theme.text.title, fontWeight: "700", color: theme.colors.charcoal }}>Discover</Text>
        <Text style={{ color: theme.colors.charcoalMuted, marginTop: 4 }}>Published recipes • pull to refresh</Text>
      </View>
      <FlashList
        data={data}
        contentContainerStyle={{ padding: 16, paddingTop: 0 }}
        renderItem={({ item }) => <RecipeCard r={item} />}
        onEndReached={() => {
          if (recipesQuery.hasNextPage) recipesQuery.fetchNextPage();
        }}
      />
    </View>
  );
}
