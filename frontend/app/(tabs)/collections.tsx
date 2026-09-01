import { View, Text, Pressable, TextInput, ActivityIndicator, Alert } from "react-native";
import { FlashList } from "@shopify/flash-list";
import { useState } from "react";
import { Link } from "expo-router";
import { useCollections, useCreateCollection, useDeleteCollection, useUpdateCollection } from "@/lib/collections";
import { theme } from "@/lib/theme";
import { useAuth } from "@/lib/auth";

export default function CollectionsScreen() {
  const { session } = useAuth();
  const { data: cols, isLoading } = useCollections();
  const create = useCreateCollection();
  const del = useDeleteCollection();
  const update = useUpdateCollection();
  const [title, setTitle] = useState("");

  if (!session) {
    return (
      <View style={{ flex: 1, backgroundColor: theme.colors.cream, alignItems: "center", justifyContent: "center", padding: 24 }}>
        <Text style={{ color: theme.colors.charcoal, fontWeight: "600" }}>Sign in to see collections</Text>
        <Link href="/(auth)/login" asChild>
          <Pressable style={{ marginTop: 12, backgroundColor: theme.colors.terracotta, paddingHorizontal: 16, paddingVertical: 10, borderRadius: theme.radius.sm }}>
            <Text style={{ color: theme.colors.card }}>Sign in</Text>
          </Pressable>
        </Link>
      </View>
    );
  }

  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream, padding: 16, paddingTop: 48 }}>
      <Text style={{ fontSize: theme.text.title, fontWeight: "700", color: theme.colors.charcoal }}>Collections</Text>

      <View style={{ flexDirection: "row", gap: 8, marginTop: 12, marginBottom: 16 }}>
        <TextInput
          placeholder="New collection title"
          value={title}
          onChangeText={setTitle}
          style={{ flex: 1, backgroundColor: theme.colors.card, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 12, color: theme.colors.charcoal }}
        />
        <Pressable
          onPress={() => {
            if (!title.trim()) return Alert.alert("Title required");
            create.mutate({ title: title.trim() }, { onSuccess: () => setTitle("") });
          }}
          style={{ backgroundColor: theme.colors.terracotta, borderRadius: theme.radius.sm, paddingHorizontal: 16, justifyContent: "center" }}
        >
          <Text style={{ color: theme.colors.card, fontWeight: "600" }}>Add</Text>
        </Pressable>
      </View>

      {isLoading ? (
        <ActivityIndicator color={theme.colors.terracotta} />
      ) : (
        <FlashList
          data={cols ?? []}
          contentContainerStyle={{ paddingBottom: 16 }}
          renderItem={({ item }) => (
            <View style={{ backgroundColor: theme.colors.card, borderRadius: theme.radius.lg, padding: 16, marginBottom: 12, borderWidth: 1, borderColor: theme.colors.border }}>
              <Link href={`/collection/${item.id}`} asChild>
                <Pressable>
                  <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>{item.title}</Text>
                  <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12, marginTop: 4 }}>
                    {item.is_public ? "Public" : "Private"} • {item.description || "No description"}
                  </Text>
                </Pressable>
              </Link>
              <View style={{ flexDirection: "row", gap: 8, marginTop: 8 }}>
                <Pressable
                  onPress={() =>
                    Alert.prompt("Edit title", undefined, (text) => {
                      if (text) update.mutate({ id: item.id, title: text });
                    })
                  }
                  style={{ paddingHorizontal: 12, paddingVertical: 6, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm }}
                >
                  <Text style={{ fontSize: 12, color: theme.colors.charcoal }}>Edit</Text>
                </Pressable>
                <Pressable
                  onPress={() =>
                    Alert.alert("Delete?", undefined, [
                      { text: "Cancel", style: "cancel" },
                      { text: "Delete", style: "destructive", onPress: () => del.mutate(item.id) },
                    ])
                  }
                  style={{ paddingHorizontal: 12, paddingVertical: 6, backgroundColor: theme.colors.terracottaDark, borderRadius: theme.radius.sm }}
                >
                  <Text style={{ fontSize: 12, color: theme.colors.card }}>Delete</Text>
                </Pressable>
              </View>
            </View>
          )}
        />
      )}
    </View>
  );
}
