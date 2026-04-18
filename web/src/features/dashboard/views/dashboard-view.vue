<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { Line, Doughnut } from "vue-chartjs";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from "chart.js";
import { Button, Badge, Card, Heading, Text } from "@/components/ui";
import { useApi } from "@/composables/use-api";
import type {
  ModelsAccount,
  ModelsTransaction,
  ModelsDraft,
  ModelsReportSummary,
  ModelsReportTrends,
  ModelsReportByCategory,
} from "@/api/types";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

const { api, loading, error, wrap } = useApi();

const accounts = ref<ModelsAccount[]>([]);
const transactions = ref<ModelsTransaction[]>([]);
const drafts = ref<ModelsDraft[]>([]);
const reportSummary = ref<ModelsReportSummary | null>(null);
const reportTrends = ref<ModelsReportTrends | null>(null);
const reportByCategory = ref<ModelsReportByCategory | null>(null);
const draftActionLoading = ref<string | null>(null);

// Budget thresholds — hardcoded mock; no real budget API exists
const MOCK_BUDGETS: Record<string, number> = {
  "cat-001": 800,   // Groceries
  "cat-002": 2500,  // Rent
  "cat-003": 10000, // Salary
};

const netBalance = computed(() =>
  accounts.value.reduce((sum, a) => sum + (a.balance || 0), 0)
);

function formatCurrency(amount?: number) {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format((amount || 0) / 100);
}

function getBadgeVariant(type?: string) {
  if (type === "income") return "teal";
  if (type === "expense") return "rust";
  if (type === "transfer") return "gold";
  return "default";
}

function getAmountColor(type?: string) {
  if (type === "income") return "text-[var(--teal)]";
  if (type === "expense") return "text-[var(--rust)]";
  return "text-[var(--ink)]";
}

function formatRelativeDate(dateStr?: string) {
  if (!dateStr) return "";
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  if (diffDays === 0) return "Today";
  if (diffDays === 1) return "Yesterday";
  if (diffDays < 7) return `${diffDays}d ago`;
  return date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

// Line chart — monthly spending trend
const trendChartData = computed(() => {
  const pts = reportTrends.value?.data_points || [];
  return {
    labels: pts.map((p) =>
      new Date(p.date || "").toLocaleDateString("en-US", { month: "short", day: "numeric" })
    ),
    datasets: [
      {
        label: "Spending",
        data: pts.map((p) => Math.abs(p.total || 0)),
        borderColor: "#3d9a8b",
        backgroundColor: "rgba(61,154,139,0.1)",
        fill: true,
        tension: 0.4,
      },
    ],
  };
});

const trendChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  scales: {
    x: { grid: { display: false }, ticks: { color: "var(--ink-soft)" } },
    y: { grid: { color: "var(--rule)" }, ticks: { color: "var(--ink-soft)" } },
  },
};

// Donut chart — spending by category
const categoryChartData = computed(() => {
  const cats = reportByCategory.value?.categories || [];
  return {
    labels: cats.map((c) => c.category_name || ""),
    datasets: [
      {
        data: cats.map((c) => c.total || 0),
        backgroundColor: ["#b5534a", "#3d9a8b", "#e8b05b", "#7a6f5d"],
        borderWidth: 0,
      },
    ],
  };
});

const categoryChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  cutout: "65%",
  plugins: {
    legend: {
      position: "bottom" as const,
      labels: { color: "var(--ink-soft)", padding: 16, boxWidth: 12 },
    },
  },
};

async function confirmDraft(id: string) {
  draftActionLoading.value = id;
  try {
    await api.drafts.confirmCreate(id);
    drafts.value = drafts.value.filter((d) => d.id !== id);
  } finally {
    draftActionLoading.value = null;
  }
}

async function rejectDraft(id: string) {
  draftActionLoading.value = id;
  try {
    await api.drafts.rejectCreate(id);
    drafts.value = drafts.value.filter((d) => d.id !== id);
  } finally {
    draftActionLoading.value = null;
  }
}

onMounted(async () => {
  const start = new Date(new Date().getFullYear(), new Date().getMonth(), 1)
    .toISOString()
    .split("T")[0];
  const end = new Date().toISOString().split("T")[0];

  await Promise.allSettled([
    wrap(async () => {
      const res = await api.accounts.accountsList({ page_size: 10 });
      accounts.value = res.data.items || [];
    }),
    wrap(async () => {
      const res = await api.transactions.transactionsList({ page_size: 5 });
      transactions.value = res.data.items || [];
    }),
    wrap(async () => {
      const res = await api.drafts.draftsList({ status: "pending" });
      drafts.value = res.data.items || [];
    }),
    wrap(async () => {
      const res = await api.reports.summaryList({ start_date: start, end_date: end });
      reportSummary.value = res.data.data || null;
    }),
    wrap(async () => {
      const res = await api.reports.trendsList({ start_date: start, end_date: end });
      reportTrends.value = res.data || null;
    }),
    wrap(async () => {
      const res = await api.reports.byCategoryList({ start_date: start, end_date: end });
      reportByCategory.value = res.data || null;
    }),
  ]);
});
</script>

<template>
  <div>
    <!-- Section A — Hero greeting strip -->
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
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        Loading...
      </div>

      <template v-else>
        <!-- Section B — Metric strip (3 tiles: Net Balance | Income | Expenses) -->
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
            <Text as="p" size="xs" muted mt="1">This month</Text>
          </div>
          <div class="bg-[var(--paper)] p-6">
            <Text as="p" size="xs" muted class="uppercase tracking-widest mb-3">
              Income
            </Text>
            <p class="font-mono text-3xl font-medium text-[var(--teal)]">
              {{ formatCurrency(reportSummary?.total_income) }}
            </p>
            <Text as="p" size="xs" muted mt="1">This month</Text>
          </div>
          <div class="bg-[var(--paper)] p-6">
            <Text as="p" size="xs" muted class="uppercase tracking-widest mb-3">
              Expenses
            </Text>
            <p class="font-mono text-3xl font-medium text-[var(--rust)]">
              {{ formatCurrency(reportSummary?.total_expenses) }}
            </p>
            <Text as="p" size="xs" muted mt="1">This month</Text>
          </div>
        </div>

        <!-- Section C — Charts row -->
        <div class="grid grid-cols-1 lg:grid-cols-5 gap-8 mb-10">
          <Card class="lg:col-span-3 animate-fade-up delay-2">
            <Heading as="h3" size="lg" mb="4">Spending Trend</Heading>
            <div class="h-64">
              <Line :data="trendChartData" :options="trendChartOptions" />
            </div>
          </Card>
          <Card class="lg:col-span-2 animate-fade-up delay-3">
            <Heading as="h3" size="lg" mb="4">By Category</Heading>
            <div class="h-64">
              <Doughnut :data="categoryChartData" :options="categoryChartOptions" />
            </div>
          </Card>
        </div>

        <!-- Section D — Budget progress bars -->
        <Card class="mb-10 animate-fade-up delay-4">
          <Heading as="h3" size="lg" mb="4">Spending vs Budget</Heading>
          <div class="space-y-6">
            <div
              v-for="cat in (reportByCategory?.categories || [])"
              :key="cat.category_id"
            >
              <div class="flex items-center justify-between mb-2">
                <span class="text-sm font-medium text-[var(--ink)]">{{ cat.category_name }}</span>
                <span class="text-xs text-[var(--ink-soft)]">
                  {{ formatCurrency(cat.total) }} / {{ formatCurrency(MOCK_BUDGETS[cat.category_id!] || 0) }}
                </span>
              </div>
              <div class="h-2 bg-[var(--rule)] rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-300"
                  :class="cat.type === 'income' ? 'bg-[var(--teal)]' : 'bg-[var(--rust)]'"
                  :style="{ width: Math.min(Math.round(((cat.total ?? 0) / (MOCK_BUDGETS[cat.category_id!] || 1)) * 100), 100) + '%' }"
                />
              </div>
            </div>
            <div v-if="!reportByCategory?.categories?.length" class="text-sm text-[var(--ink-soft)]">
              No category data available.
            </div>
          </div>
        </Card>

        <!-- Section E — Recent Transactions (bottom-left) and Section F — Pending Drafts (bottom-right) -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-10">
          <!-- Section E — Recent Transactions -->
          <Card class="animate-fade-up delay-5">
            <div class="flex items-center justify-between mb-4">
              <Heading as="h3" size="lg">Recent Transactions</Heading>
              <RouterLink
                to="/transactions"
                class="text-xs font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] underline transition-colors"
              >
                View all
              </RouterLink>
            </div>
            <div v-if="transactions.length === 0" class="p-6 text-center">
              <Text as="p" size="sm" muted>No recent transactions.</Text>
            </div>
            <div v-else class="divide-y divide-[var(--rule)]">
              <div
                v-for="tx in transactions.slice(0, 5)"
                :key="tx.id"
                class="flex items-center justify-between py-3"
              >
                <div>
                  <p class="text-sm font-medium text-[var(--ink)]">{{ tx.title || 'Untitled Transaction' }}</p>
                  <Text as="p" size="xs" muted>{{ formatRelativeDate(tx.created_at) }}</Text>
                </div>
                <span class="font-mono text-sm" :class="getAmountColor(tx.transaction_type)">
                  {{ tx.transaction_type === 'income' ? '+' : tx.transaction_type === 'expense' ? '-' : '~' }}{{ formatCurrency(tx.amount) }}
                </span>
              </div>
            </div>
          </Card>

          <!-- Section F — Pending Drafts -->
          <Card class="animate-fade-up delay-6">
            <div class="flex items-center justify-between mb-4">
              <Heading as="h3" size="lg">Pending Drafts</Heading>
              <RouterLink
                to="/drafts"
                class="text-xs font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] underline transition-colors"
              >
                View all
              </RouterLink>
            </div>
            <div v-if="drafts.length === 0" class="p-6 text-center">
              <Text as="p" size="sm" muted>No pending drafts.</Text>
            </div>
            <div v-else class="space-y-4">
              <div
                v-for="draft in drafts"
                :key="draft.id"
                class="p-4 bg-[var(--paper-dark)]/40 rounded-lg"
              >
                <div class="flex items-center justify-between mb-3">
                  <div>
                    <p class="text-sm font-medium text-[var(--ink)]">{{ draft.title || 'Untitled Draft' }}</p>
                    <Text as="p" size="xs" muted>{{ formatCurrency(draft.amount) }} · {{ draft.source }}</Text>
                  </div>
                </div>
                <div class="flex gap-2">
                  <Button variant="secondary" size="sm" :disabled="draftActionLoading === (draft.id ?? '')" @click="rejectDraft(draft.id ?? '')">Reject</Button>
                  <Button size="sm" :disabled="draftActionLoading === (draft.id ?? '')" @click="confirmDraft(draft.id ?? '')">Confirm</Button>
                </div>
              </div>
            </div>
          </Card>
        </div>

        <!-- Section G — Accounts table (collapsible) -->
        <details class="group">
          <summary class="flex items-center justify-between cursor-pointer list-none py-4">
            <Heading as="h3" size="lg">Accounts</Heading>
            <span class="text-xs text-[var(--ink-soft)] group-open:rotate-180 transition-transform">▼</span>
          </summary>
          <Card>
            <div v-if="accounts.length === 0" class="p-8 text-center">
              <Text as="p" size="sm" muted>
                No accounts yet. Create one to get started.
              </Text>
            </div>
            <table v-else class="w-full">
              <thead>
                <tr class="border-b border-[var(--rule)]">
                  <th class="text-left text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3">
                    Account
                  </th>
                  <th class="text-right text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3">
                    Balance
                  </th>
                  <th class="text-right text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3">
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
                    <span class="font-mono text-sm text-[var(--ink)]">
                      {{ formatCurrency(account.balance) }}
                    </span>
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
        </details>
      </template>
    </main>
  </div>
</template>