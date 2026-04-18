import { ref } from "vue";
import { useApi } from "@/composables/use-api";
import type {
  ModelsDraft,
  ModelsDraftPagedResponse,
  ModelsDraftResponse,
  ModelsCreateDraftRequest,
  ModelsUpdateDraftRequest,
  ModelsTransaction,
  ModelsTransactionResponse,
} from "@/api/types";

export function useDrafts() {
  const { api, loading, error, wrap } = useApi();

  const drafts = ref<ModelsDraft[]>([]);
  const draft = ref<ModelsDraft | null>(null);
  const totalItems = ref(0);
  const totalPages = ref(0);

  async function list(params?: { source?: string; status?: string; page_size?: number }) {
    return wrap(async () => {
      const res = await api.drafts.draftsList(params);
      const data = res.data as ModelsDraftPagedResponse;
      drafts.value = data.items || [];
      totalItems.value = data.total_items || 0;
      totalPages.value = data.total_pages || 0;
    });
  }

  async function get(id: string) {
    return wrap(async () => {
      const res = await api.drafts.draftsDetail(id);
      draft.value = (res.data as ModelsDraftResponse).data || null;
    });
  }

  async function create(payload: ModelsCreateDraftRequest) {
    return wrap(async () => {
      const res = await api.drafts.draftsCreate(payload);
      return (res.data as ModelsDraftResponse).data;
    });
  }

  async function update(id: string, payload: ModelsUpdateDraftRequest) {
    return wrap(async () => {
      const res = await api.drafts.draftsPartialUpdate(id, payload);
      return (res.data as ModelsDraftResponse).data;
    });
  }

  async function remove(id: string) {
    return wrap(async () => {
      await api.drafts.draftsDelete(id);
    });
  }

  async function confirm(id: string): Promise<ModelsTransaction | undefined> {
    return wrap(async () => {
      const res = await api.drafts.confirmCreate(id);
      return (res.data as ModelsTransactionResponse).data;
    });
  }

  async function reject(id: string) {
    return wrap(async () => {
      await api.drafts.rejectCreate(id);
    });
  }

  return {
    drafts,
    draft,
    totalItems,
    totalPages,
    loading,
    error,
    list,
    get,
    create,
    update,
    remove,
    confirm,
    reject,
  };
}