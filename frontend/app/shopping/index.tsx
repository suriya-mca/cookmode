import { View, Text, Pressable, TextInput, ActivityIndicator, Alert } from "react-native";
import { FlashList } from "@shopify/flash-list";
import { useState } from "react";
import { useShoppingLists, useCreateShoppingList, useDeleteShoppingList } from "@/lib/shopping";
import { theme } from "@/lib/theme";
import { Link } from "expo-router";
import { RequireAuth } from "@/components/RequireAuth";

export default function ShoppingListsScreen() {
  return (
    <RequireAuth title="Sign in to see shopping lists" subtitle="Shopping lists are tied to your account.">
      <ShoppingListsView />
    </RequireAuth>
  );
}

function ShoppingListsView() {
  const { data, isLoading } = useShoppingLists();
  const create = useCreateShoppingList();
  const del = useDeleteShoppingList();
  const [title, setTitle] = useState("");
  if (isLoading)
    return (
      <View style={{ flex: 1, alignItems: "center", justifyContent: "center", backgroundColor: theme.colors.cream }}>
        <ActivityIndicator color={theme.colors.terracotta} />
      </View>
    );
  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream, padding: 16, paddingTop: 48 }}>
      <Text style={{ fontSize: theme.text.title, fontWeight: "700", color: theme.colors.charcoal }}>Shopping Lists</Text>
      <View style={{ flexDirection: "row", gap: 8, marginTop: 12 }}>
        <TextInput
          placeholder="New list title"
          value={title}
          onChangeText={setTitle}
          style={{ flex: 1, backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 10, color: theme.colors.charcoal }}
        />
        <Pressable
          onPress={() => {
            if (!title.trim()) return Alert.alert("Title required");
            create.mutate({ title: title.trim() }, { onSuccess: () => setTitle("") });
          }}
          style={{ backgroundColor: theme.colors.terracotta, borderRadius: theme.radius.sm, paddingHorizontal: 14, justifyContent: "center" }}
        >
          <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Add</Text>
        </Pressable>
      </View>
      <FlashList
        data={data ?? []}
        contentContainerStyle={{ paddingTop: 12 }}
        renderItem={({ item }: any) => (
          <View style={{ backgroundColor: theme.colors.card, borderRadius: theme.radius.md, padding: 14, marginBottom: 10, borderWidth: 1, borderColor: theme.colors.border }}>
            <Link href={`/shopping/${item.id}`} asChild>
              <Pressable>
                <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>{item.title}</Text>
                <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12 }}>{item.items?.length ?? 0} items</Text>
              </Pressable>
            </Link>
            <Pressable
              onPress={() =>
                Alert.alert("Delete?", undefined, [
                  { text: "Cancel", style: "cancel" },
                  { text: "Delete", style: "destructive", onPress: () => del.mutate(item.id) },
                ])
              }
              style={{ marginTop: 8, alignSelf: "flex-start", paddingHorizontal: 10, paddingVertical: 6, backgroundColor: theme.colors.terracottaDark, borderRadius: theme.radius.sm }}
            >
              <Text style={{ fontSize: 12, color: theme.colors.card }}>Delete</Text>
            </Pressable>
          </View>
        )}
      />
    </View>
  );
}
