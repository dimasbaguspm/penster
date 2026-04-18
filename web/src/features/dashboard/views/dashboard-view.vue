<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import uiButton from "@/components/ui/ui-button.vue";
import uiBadge from "@/components/ui/ui-badge.vue";
import uiCard from "@/components/ui/ui-card.vue";
import { useApi } from "@/composables/use-api";
import type { ModelsAccount, ModelsReportSummary } from "@/api/types";

const { api, loading, error, wrap } = useApi();

const accounts = ref<ModelsAccount[]>([]);
const draftsCount = ref(0);
const report = ref<ModelsReportSummary | null>(null);

const netBalance = computed(() => {
  return accounts.value.reduce((sum, a) => sum + (a.balance || 0), 0);
});

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

onMounted(async () => {
  await Promise.allSettled([
    wrap(async () => {
      const res = await api.accounts.accountsList({ page_size: 10 });
      accounts.value = res.data.items || [];
    }),
    wrap(async () => {
      const res = await api.drafts.draftsList({ status: "pending" });
      draftsCount.value = res.data.total_items || 0;
    }),
    wrap(async () => {
      const now = new Date();
      const start = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split("T")[0];
      const end = now.toISOString().split("T")[0];
      const res = await api.reports.summaryList({
        start_date: start,
        end_date: end,
      });
      report.value = res.data.data || null;
    }),
  ]);
});
</script>

<template>
  <div>
    <!-- Hero greeting -->
    <section class="border-b border-[var(--rule)] bg-[var(--paper-dark)]/40">
      <div class="max-w-7xl mx-auto px-6 lg:px-10 py-10">
        <div class="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-6">
          <div class="animate-fade-up">
            <p class="text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] mb-1">
              Welcome back
            </p>
            <h2
              class="font-display text-4xl lg:text-5xl font-semibold text-[var(--ink)] leading-tight"
            >
              Your Ledger
            </h2>
          </div>
          <div class="animate-fade-up delay-2 flex items-center gap-3">
            <uiButton variant="secondary">Filter</uiButton>
            <uiButton>+ New Transaction</uiButton>
          </div>
        </div>
      </div>
    </section>

    <!-- Main content -->
    <main class="max-w-7xl mx-auto px-6 lg:px-10 py-10">
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
          <circle
            class="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            stroke-width="4"
          />
          <path
            class="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
          />
        </svg>
        Loading...
      </div>

      <template v-else>
        <!-- Summary strip -->
        <div
          class="grid grid-cols-1 sm:grid-cols-3 gap-px bg-[var(--rule)] border border-[var(--rule)] rounded-lg overflow-hidden mb-10 animate-fade-up delay-1"
        >
          <div class="bg-[var(--paper)] p-6">
            <p class="text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] mb-3">
              Net Balance
            </p>
            <p class="font-mono text-3xl font-medium text-[var(--ink)]">
              {{ formatCurrency(netBalance) }}
            </p>
            <p class="text-xs text-[var(--ink-soft)] mt-1">
              Across {{ accounts.length }}
              {{ accounts.length === 1 ? "account" : "accounts" }}
            </p>
          </div>
          <div class="bg-[var(--paper)] p-6">
            <p class="text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] mb-3">
              Income
            </p>
            <p class="font-mono text-3xl font-medium text-[var(--teal)]">
              {{ formatCurrency(report?.total_income) }}
            </p>
            <p class="text-xs text-[var(--ink-soft)] mt-1">This month</p>
          </div>
          <div class="bg-[var(--paper)] p-6">
            <p class="text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] mb-3">
              Expenses
            </p>
            <p class="font-mono text-3xl font-medium text-[var(--rust)]">
              {{ formatCurrency(report?.total_expenses) }}
            </p>
            <p class="text-xs text-[var(--ink-soft)] mt-1">This month</p>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-5 gap-8">
          <!-- Accounts list -->
          <div class="lg:col-span-3 animate-fade-up delay-2">
            <div class="flex items-center justify-between mb-4">
              <h3 class="font-display text-xl font-semibold text-[var(--ink)]">Accounts</h3>
              <RouterLink
                to="/accounts"
                class="text-xs font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] ink-underline transition-colors"
              >
                View all
              </RouterLink>
            </div>

            <uiCard>
              <div v-if="accounts.length === 0" class="p-8 text-center">
                <p class="text-sm text-[var(--ink-soft)]">
                  No accounts yet. Create one to get started.
                </p>
              </div>
              <table v-else class="w-full">
                <thead>
                  <tr class="border-b border-[var(--rule)]">
                    <th
                      class="text-left text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
                    >
                      Account
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
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="account in accounts"
                    :key="account.id"
                    class="border-b border-[var(--rule)] last:border-0 hover:bg-[var(--paper-dark)]/40 transition-colors duration-150"
                  >
                    <td class="px-5 py-4">
                      <span class="text-sm font-medium text-[var(--ink)]">{{ account.name }}</span>
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
                  </tr>
                </tbody>
              </table>
            </uiCard>
          </div>

          <!-- Quick actions -->
          <div class="lg:col-span-2 animate-fade-up delay-3">
            <h3 class="font-display text-xl font-semibold text-[var(--ink)] mb-4">Quick Actions</h3>
            <div class="space-y-3">
              <RouterLink
                to="/accounts"
                class="w-full group flex items-center gap-4 p-4 bg-[var(--paper)] border border-[var(--rule)] rounded-lg hover:border-[var(--gold)] hover:shadow-md transition-all duration-200 text-left card-hover"
              >
                <div
                  class="w-9 h-9 rounded-full bg-[var(--paper-dark)] flex items-center justify-center flex-shrink-0 group-hover:bg-[var(--gold-light)]/20 transition-colors"
                >
                  <svg
                    class="w-4 h-4 text-[var(--gold)]"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path d="M12 5v14M5 12h14" stroke-linecap="round" />
                  </svg>
                </div>
                <div>
                  <p class="text-sm font-medium text-[var(--ink)]">Add Account</p>
                  <p class="text-xs text-[var(--ink-soft)]">Track a new bank or wallet</p>
                </div>
              </RouterLink>

              <RouterLink
                to="/transactions"
                class="w-full group flex items-center gap-4 p-4 bg-[var(--paper)] border border-[var(--rule)] rounded-lg hover:border-[var(--teal)] hover:shadow-md transition-all duration-200 text-left card-hover"
              >
                <div
                  class="w-9 h-9 rounded-full bg-[var(--paper-dark)] flex items-center justify-center flex-shrink-0 group-hover:bg-[var(--teal)]/10 transition-colors"
                >
                  <svg
                    class="w-4 h-4 text-[var(--teal)]"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
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
                </div>
                <div>
                  <p class="text-sm font-medium text-[var(--ink)]">Record Transaction</p>
                  <p class="text-xs text-[var(--ink-soft)]">Log income or an expense</p>
                </div>
              </RouterLink>

              <RouterLink
                to="/drafts"
                class="w-full group flex items-center gap-4 p-4 bg-[var(--paper)] border border-[var(--rule)] rounded-lg hover:border-[var(--rust-muted)] hover:shadow-md transition-all duration-200 text-left card-hover"
              >
                <div
                  class="w-9 h-9 rounded-full bg-[var(--paper-dark)] flex items-center justify-center flex-shrink-0 group-hover:bg-[var(--rust)]/10 transition-colors"
                >
                  <svg
                    class="w-4 h-4 text-[var(--rust)]"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path
                      d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                    <path
                      d="M14 2v6h6M16 13H8M16 17H8M10 9H8"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </div>
                <div>
                  <p class="text-sm font-medium text-[var(--ink)]">Review Drafts</p>
                  <p class="text-xs text-[var(--ink-soft)]">Confirm or reject pending</p>
                </div>
                <span
                  v-if="draftsCount > 0"
                  class="ml-auto inline-flex items-center justify-center w-5 h-5 rounded-full bg-[var(--gold-light)]/20 text-xs font-mono font-medium text-[var(--gold)]"
                >
                  {{ draftsCount }}
                </span>
              </RouterLink>

              <RouterLink
                to="/reports"
                class="w-full group flex items-center gap-4 p-4 bg-[var(--paper)] border border-[var(--rule)] rounded-lg hover:border-[var(--ink-soft)] hover:shadow-md transition-all duration-200 text-left card-hover"
              >
                <div
                  class="w-9 h-9 rounded-full bg-[var(--paper-dark)] flex items-center justify-center flex-shrink-0 group-hover:bg-[var(--ink)]/10 transition-colors"
                >
                  <svg
                    class="w-4 h-4 text-[var(--ink-soft)]"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                  >
                    <path
                      d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </div>
                <div>
                  <p class="text-sm font-medium text-[var(--ink)]">View Reports</p>
                  <p class="text-xs text-[var(--ink-soft)]">Spending trends & insights</p>
                </div>
              </RouterLink>
            </div>
          </div>
        </div>
      </template>
    </main>
  </div>
</template>
