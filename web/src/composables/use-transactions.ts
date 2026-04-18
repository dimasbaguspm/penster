import { ref } from "vue";
import { useApi } from "@/composables/use-api";
import type {
  ModelsTransaction,
  ModelsTransactionPagedResponse,
  ModelsTransactionResponse,
  ModelsCreateTransactionRequest,
  ModelsUpdateTransactionRequest,
} from "@/api/types";

export function useTransactions() {
  const { api, loading, error, wrap } = useApi();

  const transactions = ref<ModelsTransaction[]>([]);
  const transaction = ref<ModelsTransaction | null>(null);
  const totalItems = ref(0);
  const totalPages = ref(0);

  async function list(params?: {
    q?: string;
    account_id?: string;
    category_id?: string;
    transaction_type?: string;
    page?: number;
    page_size?: number;
  }) {
    return wrap(async () => {
      const res = await api.transactions.transactionsList(params);
      const data = res.data as ModelsTransactionPagedResponse;
      transactions.value = data.items || [];
      totalItems.value = data.total_items || 0;
      totalPages.value = data.total_pages || 0;
    });
  }

  async function get(id: string) {
    return wrap(async () => {
      const res = await api.transactions.transactionsDetail(id);
      transaction.value = (res.data as ModelsTransactionResponse).data || null;
    });
  }

  async function create(payload: ModelsCreateTransactionRequest) {
    return wrap(async () => {
      const res = await api.transactions.transactionsCreate(payload);
      return (res.data as ModelsTransactionResponse).data;
    });
  }

  async function update(id: string, payload: ModelsUpdateTransactionRequest) {
    return wrap(async () => {
      const res = await api.transactions.transactionsUpdate(id, payload);
      return (res.data as ModelsTransactionResponse).data;
    });
  }

  async function remove(id: string) {
    return wrap(async () => {
      await api.transactions.transactionsDelete(id);
    });
  }

  return {
    transactions,
    transaction,
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