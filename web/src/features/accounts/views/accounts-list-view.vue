<script setup lang="ts">
import { ref, onMounted } from "vue";
import uiCard from "@/components/ui/ui-card.vue";
import uiButton from "@/components/ui/ui-button.vue";
import uiBadge from "@/components/ui/ui-badge.vue";
import { useApi } from "@/composables/use-api";
import type { ModelsAccount } from "@/api/types";

const { api, loading, error, wrap } = useApi();

const accounts = ref<ModelsAccount[]>([]);
const totalItems = ref(0);
const totalPages = ref(0);
const page = ref(1);
const pageSize = ref(10);

function getBadgeVariant(type?: string) {
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

async function fetchAccounts() {
  await wrap(async () => {
    const res = await api.accounts.accountsList({
      page: page.value,
      page_size: pageSize.value,
    });
    accounts.value = res.data.items || [];
    totalItems.value = res.data.total_items || 0;
    totalPages.value = res.data.total_pages || 0;
  });
}

onMounted(fetchAccounts);
</script>

<template>
  <div class="max-w-7xl mx-auto px-6 lg:px-10 py-10">
    <div class="flex items-center justify-between mb-8 animate-fade-up">
      <div>
        <h1 class="font-display text-3xl font-semibold text-[var(--ink)]">Accounts</h1>
        <p class="text-sm text-[var(--ink-soft)] mt-1">Manage your bank accounts and wallets</p>
      </div>
      <uiButton>+ Add Account</uiButton>
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
      <div v-if="accounts.length === 0" class="p-12 text-center">
        <svg
          class="w-10 h-10 mx-auto text-[var(--rule)] mb-4"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1"
        >
          <path
            d="M3 21h18M3 7v1a3 3 0 003 3h12a3 3 0 003-3V7M3 7l9-4 9 4M9 21V11h6v10"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
        <h3 class="font-display text-lg font-medium text-[var(--ink)] mb-1">No accounts yet</h3>
        <p class="text-sm text-[var(--ink-soft)] mb-6">
          Create your first account to start tracking your finances.
        </p>
        <uiButton>Create Account</uiButton>
      </div>

      <table v-else class="w-full">
        <thead>
          <tr class="border-b border-[var(--rule)]">
            <th
              class="text-left text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Name
            </th>
            <th
              class="text-right text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Balance
            </th>
            <th
              class="text-right text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Type
            </th>
            <th
              class="text-right text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Created
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="account in accounts"
            :key="account.id"
            class="border-b border-[var(--rule)] last:border-0 hover:bg-[var(--paper-dark)]/40 transition-colors duration-150 cursor-pointer"
          >
            <td class="px-5 py-4">
              <RouterLink
                :to="`/accounts/${account.id}`"
                class="text-sm font-medium text-[var(--ink)] hover:text-[var(--gold)] transition-colors"
              >
                {{ account.name }}
              </RouterLink>
            </td>
            <td class="px-5 py-4 text-right">
              <span class="font-mono text-sm text-[var(--ink)]">{{
                formatCurrency(account.balance)
              }}</span>
            </td>
            <td class="px-5 py-4 text-right">
              <uiBadge :variant="getBadgeVariant(account.type)">
                {{ account.type }}
              </uiBadge>
            </td>
            <td class="px-5 py-4 text-right">
              <span class="text-xs text-[var(--ink-soft)]">
                {{ account.created_at ? new Date(account.created_at).toLocaleDateString() : "—" }}
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
          Showing {{ accounts.length }} of {{ totalItems }} accounts
        </p>
        <div class="flex items-center gap-2">
          <uiButton
            variant="secondary"
            size="sm"
            :disabled="page <= 1"
            @click="
              page--;
              fetchAccounts();
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
              fetchAccounts();
            "
          >
            Next
          </uiButton>
        </div>
      </div>
    </uiCard>
  </div>
</template>
