import { ref } from "vue";
import { useApi } from "@/composables/use-api";
import type {
  ModelsAccount,
  ModelsAccountPagedResponse,
  ModelsAccountResponse,
  ModelsCreateAccountRequest,
  ModelsUpdateAccountRequest,
} from "@/api/types";

export function useAccounts() {
  const { api, loading, error, wrap } = useApi();

  const accounts = ref<ModelsAccount[]>([]);
  const account = ref<ModelsAccount | null>(null);
  const totalItems = ref(0);
  const totalPages = ref(0);

  async function list(params?: { q?: string; page?: number; page_size?: number }) {
    return wrap(async () => {
      const res = await api.accounts.accountsList(params);
      const data = res.data as ModelsAccountPagedResponse;
      accounts.value = data.items || [];
      totalItems.value = data.total_items || 0;
      totalPages.value = data.total_pages || 0;
    });
  }

  async function get(id: string) {
    return wrap(async () => {
      const res = await api.accounts.accountsDetail(id);
      account.value = (res.data as ModelsAccountResponse).data || null;
    });
  }

  async function create(payload: ModelsCreateAccountRequest) {
    return wrap(async () => {
      const res = await api.accounts.accountsCreate(payload);
      return (res.data as ModelsAccountResponse).data;
    });
  }

  async function update(id: string, payload: ModelsUpdateAccountRequest) {
    return wrap(async () => {
      const res = await api.accounts.accountsUpdate(id, payload);
      return (res.data as ModelsAccountResponse).data;
    });
  }

  async function remove(id: string) {
    return wrap(async () => {
      await api.accounts.accountsDelete(id);
    });
  }

  return {
    accounts,
    account,
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