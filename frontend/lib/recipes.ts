import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import { apiFetch } from "./api";

export type Recipe = {
  id: string;
  user_id: string;
  title: string;
  description: string;
  cuisine: string;
  prep_time_min: number;
  cook_time_min: number;
  servings: number;
  difficulty: string;
  dietary_tags: string[];
  ingredients: { name: string; quantity: number; unit: string }[];
  steps: { order: number; text: string; anchor_seconds: number }[];
  video_hls_url: string;
  video_thumbnail_url: string;
  saves: number;
  views: number;
  status: string;
  nutrition?: unknown;
};

export function useRecipes(limit = 20) {
  return useInfiniteQuery({
    queryKey: ["recipes", limit],
    queryFn: ({ pageParam }) =>
      apiFetch<{ data: Recipe[]; next_cursor: string }>(
        `/recipes?limit=${limit}${pageParam ? `&cursor=${pageParam}` : ""}`
      ),
    initialPageParam: "" as string,
    getNextPageParam: (last) => last.next_cursor || undefined,
  });
}

export function useRecipe(id: string) {
  return useQuery({
    queryKey: ["recipe", id],
    queryFn: () => apiFetch<Recipe>(`/recipes/${id}`),
    enabled: !!id,
  });
}

export function useSearch(params: {
  q?: string;
  ingredients?: string;
  dietary?: string;
  difficulty?: string;
  cuisine?: string;
  max_total_time?: number;
}) {
  const qs = new URLSearchParams();
  if (params.q) qs.set("q", params.q);
  if (params.ingredients) qs.set("ingredients", params.ingredients);
  if (params.dietary) qs.set("dietary", params.dietary);
  if (params.difficulty) qs.set("difficulty", params.difficulty);
  if (params.cuisine) qs.set("cuisine", params.cuisine);
  if (params.max_total_time) qs.set("max_total_time", String(params.max_total_time));
  const qstr = qs.toString();
  return useQuery({
    queryKey: ["search", qstr],
    queryFn: () =>
      apiFetch<{ hits: Recipe[] }>(`/search?${qstr}`),
    enabled: !!params.q || !!params.ingredients || !!params.dietary,
  });
}
