import { View, Text, Pressable, TextInput, Alert } from "react-native";
import { useState } from "react";
import { useAuth } from "@/lib/auth";
import { theme } from "@/lib/theme";
import { RequireAuth } from "@/components/RequireAuth";
import { useFollowers, useFollowing, usePosts, useFollow, useUnfollow, useCreatePost, useDeletePost } from "@/lib/social";

export default function ProfileScreen() {
  return (
    <RequireAuth title="Sign in to see your profile">
      <ProfileView />
    </RequireAuth>
  );
}

function ProfileView() {
  const { user, signOut } = useAuth();
  const uid = user?.id ?? "";
  const { data: followers } = useFollowers(uid);
  const { data: following } = useFollowing(uid);
  const { data: posts } = usePosts(undefined, uid);
  const follow = useFollow();
  const unfollow = useUnfollow();
  const createPost = useCreatePost();
  const delPost = useDeletePost();
  const [followId, setFollowId] = useState("");
  const [postRecipeId, setPostRecipeId] = useState("");
  const [photoUrl, setPhotoUrl] = useState("");

  return (
    <View style={{ flex: 1, backgroundColor: theme.colors.cream, padding: 16, paddingTop: 48, gap: 16 }}>
      <View style={{ alignItems: "center", gap: 8 }}>
        <View style={{ width: 80, height: 80, borderRadius: 40, backgroundColor: theme.colors.border, alignItems: "center", justifyContent: "center" }}>
          <Text style={{ fontSize: 28, fontWeight: "700", color: theme.colors.charcoal }}>{(user?.username?.[0] ?? user?.display_name?.[0] ?? "U").toUpperCase()}</Text>
        </View>
        <Text style={{ fontSize: theme.text.subtitle, fontWeight: "700", color: theme.colors.charcoal }}>{user?.username ?? user?.display_name ?? "User"}</Text>
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
            <View key={p.id} style={{ flexDirection: "row", justifyContent: "space-between", alignItems: "center", paddingVertical: 8, borderBottomWidth: 1, borderBottomColor: theme.colors.borderSoft }}>
              <View style={{ flex: 1 }}>
                <Text style={{ color: theme.colors.charcoal, fontWeight: "500" }}>{p.caption || "Made it!"}</Text>
                <Text style={{ color: theme.colors.charcoalMuted, fontSize: 12 }}>{p.photo_url}</Text>
              </View>
              <Pressable onPress={() => delPost.mutate(p.id)} style={{ paddingHorizontal: 10, paddingVertical: 6, borderWidth: 1, borderColor: theme.colors.terracottaDark, borderRadius: theme.radius.sm }}>
                <Text style={{ fontSize: 12, color: theme.colors.terracottaDark }}>Delete</Text>
              </Pressable>
            </View>
          ))
        )}
        <View style={{ flexDirection: "row", gap: 8, marginTop: 12 }}>
          <TextInput placeholder="Recipe ID" value={postRecipeId} onChangeText={setPostRecipeId} style={{ flex: 1, backgroundColor: theme.colors.cream, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 8, color: theme.colors.charcoal, fontSize: 12 }} />
          <TextInput placeholder="Photo URL" value={photoUrl} onChangeText={setPhotoUrl} style={{ flex: 1, backgroundColor: theme.colors.cream, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 8, color: theme.colors.charcoal, fontSize: 12 }} />
          <Pressable
            onPress={() => {
              if (!postRecipeId.trim() || !photoUrl.trim()) return Alert.alert("Recipe ID and Photo URL required");
              createPost.mutate({ recipe_id: postRecipeId.trim(), photo_url: photoUrl.trim() });
            }}
            style={{ backgroundColor: theme.colors.sage, borderRadius: theme.radius.sm, paddingHorizontal: 12, justifyContent: "center" }}
          >
            <Text style={{ color: theme.colors.card, fontWeight: "600", fontSize: 12 }}>Post</Text>
          </Pressable>
        </View>
      </View>

      <View style={{ backgroundColor: theme.colors.card, borderRadius: theme.radius.lg, padding: 16, borderWidth: 1, borderColor: theme.colors.border }}>
        <Text style={{ fontWeight: "600", color: theme.colors.charcoal, marginBottom: 8 }}>Follow</Text>
        <View style={{ flexDirection: "row", gap: 8 }}>
          <TextInput placeholder="User ID to follow" value={followId} onChangeText={setFollowId} style={{ flex: 1, backgroundColor: theme.colors.cream, borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, padding: 8, color: theme.colors.charcoal, fontSize: 12 }} />
          <Pressable onPress={() => followId && follow.mutate(followId.trim())} style={{ backgroundColor: theme.colors.terracotta, borderRadius: theme.radius.sm, paddingHorizontal: 12, justifyContent: "center" }}>
            <Text style={{ color: theme.colors.card, fontWeight: "600", fontSize: 12 }}>Follow</Text>
          </Pressable>
          <Pressable onPress={() => followId && unfollow.mutate(followId.trim())} style={{ borderWidth: 1, borderColor: theme.colors.border, borderRadius: theme.radius.sm, paddingHorizontal: 12, justifyContent: "center" }}>
            <Text style={{ color: theme.colors.charcoal, fontSize: 12 }}>Unfollow</Text>
          </Pressable>
        </View>
      </View>
    </View>
  );
}
