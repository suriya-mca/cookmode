import { View, Text, TextInput, Pressable, ActivityIndicator, Alert } from "react-native";
import { FlashList } from "@shopify/flash-list";
import { useState } from "react";
import { useShoppingList, useToggleShoppingItem, useAddShoppingItems } from "@/lib/shopping";
import { theme } from "@/lib/theme";
import { useLocalSearchParams } from "expo-router";
import { RequireAuth } from "@/components/RequireAuth";

export default function ShoppingListDetail() {
  return (
    <RequireAuth title="Sign in to see this list" subtitle="Shopping lists are tied to your account.">
      <ShoppingListView />
    </RequireAuth>
  );
}

function ShoppingListView() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { data: list, isLoading } = useShoppingList(id!);
  const toggle = useToggleShoppingItem();
  const add = useAddShoppingItems();
  const [name, setName] = useState("");
  const [qty, setQty] = useState("");
  const [unit, setUnit] = useState("");
  if (isLoading) return <View style={{ flex: 1, alignItems: "center", justifyContent: "center" }}><ActivityIndicator /></View>;
  if (!list) return <View style={{ flex: 1, alignItems: "center", justifyContent: "center" }}><Text>Not found</Text></View>;
  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream, padding: 16, paddingTop: 48 }}>
      <Text style={{ fontSize: theme.text.title, fontWeight: "700", color: theme.colors.charcoal }}>{list.title}</Text>
      <View style={{ flexDirection: "row", gap: 8, marginTop: 12, marginBottom: 8 }}>
        <TextInput placeholder="Item" value={name} onChangeText={setName} style={{ flex: 1, backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 10, color: theme.colors.charcoal }} />
        <TextInput placeholder="Qty" value={qty} onChangeText={setQty} keyboardType="numeric" style={{ width: 60, backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 10, color: theme.colors.charcoal }} />
        <TextInput placeholder="Unit" value={unit} onChangeText={setUnit} style={{ width: 60, backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 10, color: theme.colors.charcoal }} />
        <Pressable
          onPress={() => {
            if (!name.trim()) return Alert.alert("Name required");
            const quantity = parseFloat(qty) || 1;
            add.mutate({ id: id!, items: [{ ingredient_name: name.trim(), quantity, unit }] }, { onSuccess: () => { setName(""); setQty(""); setUnit(""); } });
          }}
          style={{ backgroundColor: theme.colors.terracotta, borderRadius: theme.radius.sm, paddingHorizontal: 14, justifyContent: "center" }}
        >
          <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Add</Text>
        </Pressable>
      </View>
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
