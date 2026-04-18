<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { Button, Badge, Card, Heading, Text, Icon } from "@/components/ui";
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
            <Text as="p" size="xs" muted class="uppercase tracking-widest mb-1">
              Welcome back
            </Text>
            <Heading as="h2" size="4xl" class="lg:text-5xl">
              Your Ledger
            </Heading>
          </div>
          <div class="animate-fade-up delay-2 flex items-center gap-3">
            <Button variant="secondary">Filter</Button>
            <Button>+ New Transaction</Button>
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
            <Text as="p" size="xs" muted class="uppercase tracking-widest mb-3">
              Net Balance
            </Text>
            <p class="font-mono text-3xl font-medium text-[var(--ink)]">
              {{ formatCurrency(netBalance) }}
            </p>
            <Text as="p" size="xs" muted mt="1">
              Across {{ accounts.length }}
              {{ accounts.length === 1 ? "account" : "accounts" }}
            </Text>
          </div>
          <div class="bg-[var(--paper)] p-6">
            <Text as="p" size="xs" muted class="uppercase tracking-widest mb-3">
              Income
            </Text>
            <p class="font-mono text-3xl font-medium text-[var(--teal)]">
              {{ formatCurrency(report?.total_income) }}
            </p>
            <Text as="p" size="xs" muted mt="1">This month</Text>
          </div>
          <div class="bg-[var(--paper)] p-6">
            <Text as="p" size="xs" muted class="uppercase tracking-widest mb-3">
              Expenses
            </Text>
            <p class="font-mono text-3xl font-medium text-[var(--rust)]">
              {{ formatCurrency(report?.total_expenses) }}
            </p>
            <Text as="p" size="xs" muted mt="1">This month</Text>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-5 gap-8">
          <!-- Accounts list -->
          <div class="lg:col-span-3 animate-fade-up delay-2">
            <div class="flex items-center justify-between mb-4">
              <Heading as="h3" size="xl">Accounts</Heading>
              <RouterLink
                to="/accounts"
                class="text-xs font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] ink-underline transition-colors"
              >
                View all
              </RouterLink>
            </div>

            <Card>
              <div v-if="accounts.length === 0" class="p-8 text-center">
                <Text as="p" size="sm" muted>
                  No accounts yet. Create one to get started.
                </Text>
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
                      <Badge :variant="getBadgeVariant(account.type)">
                        {{ account.type }}
                      </Badge>
                    </td>
                  </tr>
                </tbody>
              </table>
            </Card>
          </div>

          <!-- Quick actions -->
          <div class="lg:col-span-2 animate-fade-up delay-3">
            <Heading as="h3" size="xl" mb="4">Quick Actions</Heading>
            <div class="space-y-3">
              <RouterLink
                to="/accounts"
                class="w-full group flex items-center gap-4 p-4 bg-[var(--paper)] border border-[var(--rule)] rounded-lg hover:border-[var(--gold)] hover:shadow-md transition-all duration-200 text-left card-hover"
              >
                <div
                  class="w-9 h-9 rounded-full bg-[var(--paper-dark)] flex items-center justify-center flex-shrink-0 group-hover:bg-[var(--gold-light)]/20 transition-colors"
                >
                  <Icon name="plus" size="sm" class="text-[var(--gold)]" />
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
                  <Icon name="file-plus" size="sm" class="text-[var(--teal)]" />
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
                  <Icon name="file-text" size="sm" class="text-[var(--rust)]" />
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
                  <Icon name="bar-chart-2" size="sm" class="text-[var(--ink-soft)]" />
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