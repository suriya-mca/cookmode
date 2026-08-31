import { View, Text, Pressable, TextInput, ActivityIndicator } from "react-native";
import { FlashList } from "@shopify/flash-list";
import { Image } from "expo-image";
import { Link } from "expo-router";
import { useState } from "react";
import { useRecipes, useSearch, type Recipe } from "@/lib/recipes";
import { theme } from "@/lib/theme";

const chips = ["All", "Italian", "Asian", "Vegan", "Under 30m"];

function Chip({
  label,
  active,
  onPress,
}: {
  label: string;
  active: boolean;
  onPress: () => void;
}) {
  return (
    <Pressable
      onPress={onPress}
      style={{
        paddingHorizontal: 16,
        paddingVertical: 8,
        borderRadius: theme.radius.pill,
        backgroundColor: active ? theme.colors.terracotta : theme.colors.card,
        borderWidth: 1,
        borderColor: active ? theme.colors.terracotta : theme.colors.border,
        marginRight: 8,
      }}
    >
      <Text
        style={{
          color: active ? theme.colors.card : theme.colors.charcoal,
          fontSize: 13,
          fontWeight: "500",
        }}
      >
        {label}
      </Text>
    </Pressable>
  );
}

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
  const [q, setQ] = useState("");
  const [activeChip, setActiveChip] = useState("All");
  const [searchQ, setSearchQ] = useState("");

  const dietary = activeChip === "Vegan" ? "vegan" : undefined;
  const cuisine = ["Italian", "Asian"].includes(activeChip) ? activeChip : undefined;
  const maxTime = activeChip === "Under 30m" ? 30 : undefined;
  const useSearchMode = !!searchQ || !!dietary || !!cuisine || !!maxTime;

  const recipesQuery = useRecipes(20);
  const searchQuery = useSearch({
    q: searchQ || undefined,
    dietary,
    cuisine,
    max_total_time: maxTime,
  });

  const data: Recipe[] = useSearchMode
    ? (searchQuery.data?.hits ?? [])
    : (recipesQuery.data?.pages.flatMap((p) => p.data) ?? []);

  const isLoading = useSearchMode ? searchQuery.isLoading : recipesQuery.isLoading;

  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream }}>
      <View style={{ padding: theme.spacing.md, gap: 12, paddingTop: 48 }}>
        <TextInput
          placeholder="Search recipes, ingredients..."
          placeholderTextColor={theme.colors.charcoalMuted}
          value={q}
          onChangeText={setQ}
          onSubmitEditing={() => setSearchQ(q)}
          returnKeyType="search"
          style={{
            backgroundColor: theme.colors.card,
            borderRadius: theme.radius.pill,
            paddingHorizontal: 16,
            paddingVertical: 12,
            borderWidth: 1,
            borderColor: theme.colors.border,
            color: theme.colors.charcoal,
          }}
        />
        <View style={{ flexDirection: "row" }}>
          <FlashList
            data={chips}
            horizontal
            showsHorizontalScrollIndicator={false}
            renderItem={({ item }) => (
              <Chip label={item} active={activeChip === item} onPress={() => setActiveChip(item)} />
            )}
          />
        </View>
      </View>

      {isLoading ? (
        <View style={{ flex: 1, alignItems: "center", justifyContent: "center" }}>
          <ActivityIndicator color={theme.colors.terracotta} />
        </View>
      ) : (
        <FlashList
          data={data}
          contentContainerStyle={{ padding: 16, paddingTop: 0 }}
          renderItem={({ item }) => <RecipeCard r={item} />}
          onEndReached={() => {
            if (!useSearchMode && recipesQuery.hasNextPage) recipesQuery.fetchNextPage();
          }}
        />
      )}
    </View>
  );
}
