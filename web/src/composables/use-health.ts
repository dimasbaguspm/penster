import { ref } from "vue";
import { useApi } from "@/composables/use-api";
import type { HandlerHealthResponse } from "@/api/types";

export function useHealth() {
  const { api, loading, error, wrap } = useApi();

  const status = ref<HandlerHealthResponse | null>(null);

  async function check() {
    return wrap(async () => {
      const res = await api.health.healthList();
      status.value = res.data || null;
    });
  }

  return { status, loading, error, check };
}