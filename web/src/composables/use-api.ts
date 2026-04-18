import { ref } from "vue";
import { Api } from "@/api/types";

const apiClient = new Api({
  baseUrl: import.meta.env.VITE_API_BASE_URL || "http://localhost:8080",
});

export function useApi() {
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function wrap<T>(fn: () => Promise<T>) {
    loading.value = true;
    error.value = null;
    try {
      return await fn();
    } catch (e: any) {
      error.value = e?.error?.error || e?.message || "An error occurred";
      throw e;
    } finally {
      loading.value = false;
    }
  }

  return {
    api: apiClient,
    loading,
    error,
    wrap,
  };
}
