import { View, Text, Pressable } from "react-native";
import { useAuth } from "@/lib/auth";
import { theme } from "@/lib/theme";
import { Link } from "expo-router";
import { useFollowers, useFollowing, usePosts } from "@/lib/social";

export default function ProfileScreen() {
  const { session, user, signOut } = useAuth();
  const uid = user?.id ?? "";
  const { data: followers } = useFollowers(uid);
  const { data: following } = useFollowing(uid);
  const { data: posts } = usePosts(undefined, uid);

  if (!session) {
    return (
      <View style={{ flex: 1, backgroundColor: theme.colors.cream, alignItems: "center", justifyContent: "center", padding: 24 }}>
        <Text style={{ fontWeight: "600", color: theme.colors.charcoal }}>Sign in to see profile</Text>
        <Link href="/(auth)/login" asChild>
          <Pressable style={{ marginTop: 12, backgroundColor: theme.colors.terracotta, paddingHorizontal: 16, paddingVertical: 10, borderRadius: theme.radius.sm }}>
            <Text style={{ color: theme.colors.card }}>Sign in</Text>
          </Pressable>
        </Link>
      </View>
    );
  }

  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream, padding: 16, paddingTop: 48, gap: 16 }}>
      <View style={{ alignItems: "center", gap: 8 }}>
        <View style={{ width: 80, height: 80, borderRadius: 40, backgroundColor: theme.colors.border, alignItems: "center", justifyContent: "center" }}>
          <Text style={{ fontSize: 28, fontWeight: "700", color: theme.colors.charcoal }}>{user?.email?.[0]?.toUpperCase() ?? "U"}</Text>
        </View>
        <Text style={{ fontSize: theme.text.subtitle, fontWeight: "700", color: theme.colors.charcoal }}>{user?.email}</Text>
        <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12 }}>
          {followers?.length ?? 0} Followers • {following?.length ?? 0} Following • {posts?.length ?? 0} Made
        </Text>
        <Pressable onPress={() => signOut()} style={{ marginTop: 8, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.pill, paddingHorizontal: 16, paddingVertical: 8 }}>
          <Text style={{ color: theme.colors.charcoal }}>Sign out</Text>
        </Pressable>
      </View>

      <View style={{ backgroundColor: theme.colors.card, borderRadius: theme.radius.lg, padding: 16, borderWidth: 1, borderColor: theme.colors.border }}>
        <Text style={{ fontWeight: "600", color: theme.colors.charcoal, marginBottom: 8 }}>I made this</Text>
        {(posts ?? []).length === 0 ? (
          <Text style={{ color: theme.colors.charcoalMuted, fontSize: 13 }}>No posts yet — cook a recipe and share a photo.</Text>
        ) : (
          posts!.slice(0, 3).map((p: any) => (
            <View key={p.id} style={{ paddingVertical: 8, borderBottomWidth: 1, borderBottomColor: theme.colors.borderSoft }}>
              <Text style={{ color: theme.colors.charcoal, fontWeight: "500" }}>{p.caption || "Made it!"}</Text>
              <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12 }}>{p.photo_url}</Text>
            </View>
          ))
        )}
      </View>
    </View>
  );
}
