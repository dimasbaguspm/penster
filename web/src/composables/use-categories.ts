import { ref } from "vue";
import { useApi } from "@/composables/use-api";
import type {
  ModelsCategory,
  ModelsCategoryPagedResponse,
  ModelsCategoryResponse,
  ModelsCreateCategoryRequest,
  ModelsUpdateCategoryRequest,
} from "@/api/types";

export function useCategories() {
  const { api, loading, error, wrap } = useApi();

  const categories = ref<ModelsCategory[]>([]);
  const category = ref<ModelsCategory | null>(null);
  const totalItems = ref(0);
  const totalPages = ref(0);

  async function list(params?: { q?: string; page?: number; page_size?: number }) {
    return wrap(async () => {
      const res = await api.categories.categoriesList(params);
      const data = res.data as ModelsCategoryPagedResponse;
      categories.value = data.items || [];
      totalItems.value = data.total_items || 0;
      totalPages.value = data.total_pages || 0;
    });
  }

  async function get(id: string) {
    return wrap(async () => {
      const res = await api.categories.categoriesDetail(id);
      category.value = (res.data as ModelsCategoryResponse).data || null;
    });
  }

  async function create(payload: ModelsCreateCategoryRequest) {
    return wrap(async () => {
      const res = await api.categories.categoriesCreate(payload);
      return (res.data as ModelsCategoryResponse).data;
    });
  }

  async function update(id: string, payload: ModelsUpdateCategoryRequest) {
    return wrap(async () => {
      const res = await api.categories.categoriesUpdate(id, payload);
      return (res.data as ModelsCategoryResponse).data;
    });
  }

  async function remove(id: string) {
    return wrap(async () => {
      await api.categories.categoriesDelete(id);
    });
  }

  return {
    categories,
    category,
    totalItems,
    totalPages,
    loading,
    error,
    list,
    get,
    create,
    update,
    remove,
  };
}