import { View, Text, Pressable, ActivityIndicator } from "react-native";
import { FlashList } from "@shopify/flash-list";
import { useShoppingList, useToggleShoppingItem } from "@/lib/shopping";
import { theme } from "@/lib/theme";
import { useLocalSearchParams } from "expo-router";

export default function ShoppingListDetail() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { data: list, isLoading } = useShoppingList(id!);
  const toggle = useToggleShoppingItem();
  if (isLoading) return <View style={{ flex: 1, alignItems: "center", justifyContent: "center" }}><ActivityIndicator /></View>;
  if (!list) return <View style={{ flex: 1, alignItems: "center", justifyContent: "center" }}><Text>Not found</Text></View>;
  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream, padding: 16, paddingTop: 48 }}>
      <Text style={{ fontSize: theme.text.title, fontWeight: "700", color: theme.colors.charcoal }}>{list.title}</Text>
      <FlashList
        data={list.items}
        contentContainerStyle={{ paddingTop: 12 }}
        renderItem={({ item, index }: any) => (
          <Pressable
            onPress={() => toggle.mutate({ id: id!, index, checked: !item.checked })}
            style={{ flexDirection: "row", alignItems: "center", gap: 12, backgroundColor: theme.colors.card, padding: 12, borderRadius: theme.radius.sm, marginBottom: 8, borderWidth: 1, borderColor: theme.colors.border, opacity: item.checked ? 0.6 : 1 }}
          >
            <View style={{ width: 22, height: 22, borderRadius: 11, borderWidth: 1, borderColor: theme.colors.sage, backgroundColor: item.checked ? theme.colors.sage : "transparent", alignItems: "center", justifyContent: "center" }}>
              {item.checked && <Text style={{ color: theme.colors.card, fontSize: 12 }}>✓</Text>}
            </View>
            <Text style={{ flex: 1, color: theme.colors.charcoal, textDecorationLine: item.checked ? "line-through" : "none" }}>
              {item.ingredient_name} — {item.quantity} {item.unit}
            </Text>
          </Pressable>
        )}
      />
    </View>
  );
}
