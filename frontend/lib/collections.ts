import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "./api";

export type Collection = {
  id: string;
  user_id: string;
  title: string;
  description: string;
  is_public: boolean;
};

export function useCollections() {
  return useQuery({
    queryKey: ["collections"],
    queryFn: () => api.get<{ data: Collection[] }>("/collections").then((r) => r.data),
  });
}

export function useCollection(id: string) {
  return useQuery({
    queryKey: ["collection", id],
    queryFn: () => api.get<Collection>(`/collections/${id}`),
    enabled: !!id,
  });
}

export function useCollectionRecipes(id: string) {
  return useQuery({
    queryKey: ["collection-recipes", id],
    queryFn: () => api.get<{ data: any[] }>(`/collections/${id}/recipes`).then((r) => r.data),
    enabled: !!id,
  });
}

export function useCreateCollection() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { title: string; description?: string; is_public?: boolean }) =>
      api.post<Collection>("/collections", body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["collections"] }),
  });
}

export function useAddToCollection() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, recipe_id }: { id: string; recipe_id: string }) =>
      api.post(`/collections/${id}/recipes`, { recipe_id }),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: ["collections"] });
      qc.invalidateQueries({ queryKey: ["collection-recipes", vars.id] });
    },
  });
}

export function useRemoveFromCollection() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, recipeId }: { id: string; recipeId: string }) =>
      api.del(`/collections/${id}/recipes/${recipeId}`),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: ["collection-recipes", vars.id] });
    },
  });
}
