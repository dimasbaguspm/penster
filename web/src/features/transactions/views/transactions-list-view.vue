<script setup lang="ts">
import { ref, onMounted } from "vue";
import uiCard from "@/components/ui/ui-card.vue";
import uiButton from "@/components/ui/ui-button.vue";
import uiBadge from "@/components/ui/ui-badge.vue";
import { useApi } from "@/composables/use-api";
import type { ModelsTransaction, ModelsTransactionType } from "@/api/types";

const { api, loading, error, wrap } = useApi();

const transactions = ref<ModelsTransaction[]>([]);
const totalItems = ref(0);
const totalPages = ref(0);
const page = ref(1);
const pageSize = ref(10);

function getBadgeVariant(type?: ModelsTransactionType) {
  if (type === "income") return "teal";
  if (type === "expense") return "rust";
  return "default";
}

function formatCurrency(amount?: number) {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format((amount || 0) / 100);
}

async function fetchTransactions() {
  await wrap(async () => {
    const res = await api.transactions.transactionsList({
      page: page.value,
      page_size: pageSize.value,
    });
    transactions.value = res.data.items || [];
    totalItems.value = res.data.total_items || 0;
    totalPages.value = res.data.total_pages || 0;
  });
}

onMounted(fetchTransactions);
</script>

<template>
  <div class="max-w-7xl mx-auto px-6 lg:px-10 py-10">
    <div class="flex items-center justify-between mb-8 animate-fade-up">
      <div>
        <h1 class="font-display text-3xl font-semibold text-[var(--ink)]">Transactions</h1>
        <p class="text-sm text-[var(--ink-soft)] mt-1">A record of all your income and expenses</p>
      </div>
      <uiButton>+ New Transaction</uiButton>
    </div>

    <!-- Error banner -->
    <div
      v-if="error"
      class="mb-6 px-4 py-3 bg-[var(--rust)]/10 border border-[var(--rust)]/30 rounded-lg text-sm text-[var(--rust)]"
    >
      {{ error }}
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex items-center justify-center py-16 text-[var(--ink-soft)]">
      <svg class="w-5 h-5 animate-spin mr-3" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path
          class="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
        />
      </svg>
      Loading...
    </div>

    <uiCard v-else hover>
      <div v-if="transactions.length === 0" class="p-12 text-center">
        <svg
          class="w-10 h-10 mx-auto text-[var(--rule)] mb-4"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1"
        >
          <path
            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
          <rect
            x="9"
            y="3"
            width="6"
            height="4"
            rx="1"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
          <path d="M9 12h6M9 16h4" stroke-linecap="round" />
        </svg>
        <h3 class="font-display text-lg font-medium text-[var(--ink)] mb-1">No transactions yet</h3>
        <p class="text-sm text-[var(--ink-soft)] mb-6">
          Record your first transaction to begin tracking.
        </p>
        <uiButton>Record Transaction</uiButton>
      </div>

      <table v-else class="w-full">
        <thead>
          <tr class="border-b border-[var(--rule)]">
            <th
              class="text-left text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Title
            </th>
            <th
              class="text-right text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Amount
            </th>
            <th
              class="text-right text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Type
            </th>
            <th
              class="text-right text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Date
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="tx in transactions"
            :key="tx.id"
            class="border-b border-[var(--rule)] last:border-0 hover:bg-[var(--paper-dark)]/40 transition-colors duration-150 cursor-pointer"
          >
            <td class="px-5 py-4">
              <RouterLink
                :to="`/transactions/${tx.id}`"
                class="text-sm font-medium text-[var(--ink)] hover:text-[var(--gold)] transition-colors"
              >
                {{ tx.title }}
              </RouterLink>
              <p v-if="tx.notes" class="text-xs text-[var(--ink-soft)] mt-0.5 truncate max-w-xs">
                {{ tx.notes }}
              </p>
            </td>
            <td class="px-5 py-4 text-right">
              <span
                :class="[
                  'font-mono text-sm font-medium',
                  tx.transaction_type === 'income'
                    ? 'text-[var(--teal)]'
                    : tx.transaction_type === 'expense'
                      ? 'text-[var(--rust)]'
                      : 'text-[var(--ink)]',
                ]"
              >
                {{ formatCurrency(tx.amount) }}
              </span>
            </td>
            <td class="px-5 py-4 text-right">
              <uiBadge :variant="getBadgeVariant(tx.transaction_type)">
                {{ tx.transaction_type }}
              </uiBadge>
            </td>
            <td class="px-5 py-4 text-right">
              <span class="text-xs text-[var(--ink-soft)]">
                {{ tx.created_at ? new Date(tx.created_at).toLocaleDateString() : "—" }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div
        v-if="totalPages > 1"
        class="flex items-center justify-between px-5 py-4 border-t border-[var(--rule)]"
      >
        <p class="text-xs text-[var(--ink-soft)]">
          Showing {{ transactions.length }} of {{ totalItems }} transactions
        </p>
        <div class="flex items-center gap-2">
          <uiButton
            variant="secondary"
            size="sm"
            :disabled="page <= 1"
            @click="
              page--;
              fetchTransactions();
            "
          >
            Previous
          </uiButton>
          <span class="text-xs text-[var(--ink-soft)] px-2">
            Page {{ page }} of {{ totalPages }}
          </span>
          <uiButton
            variant="secondary"
            size="sm"
            :disabled="page >= totalPages"
            @click="
              page++;
              fetchTransactions();
            "
          >
            Next
          </uiButton>
        </div>
      </div>
    </uiCard>
  </div>
</template>
