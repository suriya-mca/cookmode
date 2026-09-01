import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "./api";

export type ShoppingList = {
  id: string;
  title: string;
  items: { ingredient_name: string; quantity: number; unit: string; checked: boolean }[];
};

export function useShoppingLists() {
  return useQuery({
    queryKey: ["shopping-lists"],
    queryFn: () => api.get<{ data: ShoppingList[] }>("/shopping-lists").then((r) => r.data),
  });
}

export function useShoppingList(id: string) {
  return useQuery({
    queryKey: ["shopping-list", id],
    queryFn: () => api.get<ShoppingList>(`/shopping-lists/${id}`),
    enabled: !!id,
  });
}

export function useCreateShoppingList() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { title?: string; recipe_ids?: string[]; servings?: number }) =>
      api.post<ShoppingList>("/shopping-lists", body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["shopping-lists"] }),
  });
}

export function useToggleShoppingItem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, index, checked }: { id: string; index: number; checked: boolean }) =>
      api.patch<ShoppingList>(`/shopping-lists/${id}/items/${index}`, { checked }),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: ["shopping-list", vars.id] });
      qc.invalidateQueries({ queryKey: ["shopping-lists"] });
    },
  });
}

export function useDeleteShoppingList() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.del(`/shopping-lists/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["shopping-lists"] }),
  });
}

export function useAddShoppingItems() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, items }: { id: string; items: { ingredient_name: string; quantity: number; unit: string }[] }) =>
      api.post<ShoppingList>(`/shopping-lists/${id}/items`, { items }),
    onSuccess: (_data, vars) => {
      qc.invalidateQueries({ queryKey: ["shopping-list", vars.id] });
      qc.invalidateQueries({ queryKey: ["shopping-lists"] });
    },
  });
}
