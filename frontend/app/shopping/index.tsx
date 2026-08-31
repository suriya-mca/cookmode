import { View, Text, Pressable, ActivityIndicator } from "react-native";
import { FlashList } from "@shopify/flash-list";
import { useShoppingLists } from "@/lib/shopping";
import { theme } from "@/lib/theme";
import { Link } from "expo-router";

export default function ShoppingListsScreen() {
  const { data, isLoading } = useShoppingLists();
  if (isLoading)
    return (
      <View style={{ flex: 1, alignItems: "center", justifyContent: "center", backgroundColor: theme.colors.cream }}>
        <ActivityIndicator color={theme.colors.terracotta} />
      </View>
    );
  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream, padding: 16, paddingTop: 48 }}>
      <Text style={{ fontSize: theme.text.title, fontWeight: "700", color: theme.colors.charcoal }}>Shopping Lists</Text>
      <FlashList
        data={data ?? []}
        contentContainerStyle={{ paddingTop: 12 }}
        renderItem={({ item }: any) => (
          <Link href={`/shopping/${item.id}`} asChild>
            <Pressable style={{ backgroundColor: theme.colors.card, borderRadius: theme.radius.md, padding: 14, marginBottom: 10, borderWidth: 1, borderColor: theme.colors.border }}>
              <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>{item.title}</Text>
              <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12 }}>{item.items?.length ?? 0} items</Text>
            </Pressable>
          </Link>
        )}
      />
    </View>
  );
}
