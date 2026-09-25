import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "./api";
import { useAuth } from "./auth";

export function useFollowers(userId: string) {
  return useQuery({
    queryKey: ["followers", userId],
    queryFn: () => api.get<{ data: string[] }>(`/users/${userId}/followers`).then((r) => r.data),
    enabled: !!userId,
  });
}

export function useFollowing(userId: string) {
  return useQuery({
    queryKey: ["following", userId],
    queryFn: () => api.get<{ data: string[] }>(`/users/${userId}/following`).then((r) => r.data),
    enabled: !!userId,
  });
}

export function useFollow() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (following_id: string) => api.post("/follows", { following_id }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["followers"] });
      qc.invalidateQueries({ queryKey: ["following"] });
    },
  });
}

export function useUnfollow() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (followingId: string) => api.del(`/follows/${followingId}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["followers"] });
      qc.invalidateQueries({ queryKey: ["following"] });
    },
  });
}

export function usePosts(recipeId?: string, userId?: string) {
  const { signedIn } = useAuth();
  const qs = recipeId ? `recipe_id=${recipeId}` : userId ? `user_id=${userId}` : "";
  return useQuery({
    queryKey: ["posts", recipeId, userId],
    queryFn: () => api.get<{ data: any[] }>(`/posts?${qs}`).then((r) => r.data),
    // GET /posts is auth-required, so guests never trigger a 401.
    enabled: signedIn && !!qs,
  });
}

export function useCreatePost() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { recipe_id: string; photo_url: string; caption?: string; rating_value?: number }) =>
      api.post("/posts", body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["posts"] }),
  });
}

export function useDeletePost() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.del(`/posts/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["posts"] }),
  });
}
