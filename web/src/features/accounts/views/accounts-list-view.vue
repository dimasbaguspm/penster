<script setup lang="ts">
import { ref, onMounted } from "vue";
import { Button, Badge, Card, Heading, Text, Icon } from "@/components/ui";
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
        <Heading as="h1" size="3xl">Accounts</Heading>
        <Text as="p" size="sm" muted mt="1">Manage your bank accounts and wallets</Text>
      </div>
      <Button>+ Add Account</Button>
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

    <Card v-else hover>
      <div v-if="accounts.length === 0" class="p-12 text-center">
        <Icon name="building-2" size="xl" class="mx-auto text-[var(--rule)] mb-4" />
        <Heading as="h3" size="lg" mb="1">No accounts yet</Heading>
        <Text as="p" size="sm" muted mb="6">
          Create your first account to start tracking your finances.
        </Text>
        <Button>Create Account</Button>
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
              <Badge :variant="getBadgeVariant(account.type)">
                {{ account.type }}
              </Badge>
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
          <Button
            variant="secondary"
            size="sm"
            :disabled="page <= 1"
            @click="
              page--;
              fetchAccounts();
            "
          >
            Previous
          </Button>
          <span class="text-xs text-[var(--ink-soft)] px-2">
            Page {{ page }} of {{ totalPages }}
          </span>
          <Button
            variant="secondary"
            size="sm"
            :disabled="page >= totalPages"
            @click="
              page++;
              fetchAccounts();
            "
          >
            Next
          </Button>
        </div>
      </div>
    </Card>
  </div>
</template>