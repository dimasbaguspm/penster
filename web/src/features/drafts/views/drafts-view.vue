<script setup lang="ts">
import { ref, onMounted } from "vue";
import uiCard from "@/components/ui/ui-card.vue";
import uiBadge from "@/components/ui/ui-badge.vue";
import uiButton from "@/components/ui/ui-button.vue";
import { useApi } from "@/composables/use-api";
import type { ModelsDraft } from "@/api/types";

const { api, loading, error, wrap } = useApi();

const drafts = ref<ModelsDraft[]>([]);
const totalItems = ref(0);
const pageSize = ref(10);

function formatCurrency(amount?: number) {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format((amount || 0) / 100);
}

function getStatusVariant(status?: string) {
  if (status === "pending") return "gold";
  if (status === "confirmed") return "teal";
  if (status === "rejected") return "rust";
  return "default";
}

async function fetchDrafts() {
  await wrap(async () => {
    const res = await api.drafts.draftsList({
      status: "pending",
      page_size: pageSize.value,
    });
    drafts.value = res.data.items || [];
    totalItems.value = res.data.total_items || 0;
  });
}

async function confirmDraft(id: string) {
  await wrap(async () => {
    await api.drafts.confirmCreate(id);
    await fetchDrafts();
  });
}

async function rejectDraft(id: string) {
  await wrap(async () => {
    await api.drafts.rejectCreate(id);
    await fetchDrafts();
  });
}

onMounted(fetchDrafts);
</script>

<template>
  <div class="max-w-7xl mx-auto px-6 lg:px-10 py-10">
    <div class="flex items-center justify-between mb-8 animate-fade-up">
      <div>
        <h1 class="font-display text-3xl font-semibold text-[var(--ink)]">Drafts</h1>
        <p class="text-sm text-[var(--ink-soft)] mt-1">Pending transactions awaiting your review</p>
      </div>
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

    <uiCard v-else>
      <div v-if="drafts.length === 0" class="p-12 text-center">
        <svg
          class="w-10 h-10 mx-auto text-[var(--rule)] mb-4"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1"
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
        <h3 class="font-display text-lg font-medium text-[var(--ink)] mb-1">No pending drafts</h3>
        <p class="text-sm text-[var(--ink-soft)]">
          All drafts have been reviewed. You're all caught up.
        </p>
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
              class="text-left text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Status
            </th>
            <th
              class="text-left text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Source
            </th>
            <th
              class="text-right text-xs font-medium uppercase tracking-widest text-[var(--ink-soft)] px-5 py-3"
            >
              Actions
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="draft in drafts"
            :key="draft.id"
            class="border-b border-[var(--rule)] last:border-0 hover:bg-[var(--paper-dark)]/40 transition-colors duration-150"
          >
            <td class="px-5 py-4">
              <span class="text-sm font-medium text-[var(--ink)]">{{ draft.title }}</span>
              <p v-if="draft.notes" class="text-xs text-[var(--ink-soft)] mt-0.5 truncate max-w-xs">
                {{ draft.notes }}
              </p>
            </td>
            <td class="px-5 py-4 text-right">
              <span class="font-mono text-sm font-medium text-[var(--ink)]">
                {{ formatCurrency(draft.amount) }}
              </span>
            </td>
            <td class="px-5 py-4">
              <uiBadge :variant="getStatusVariant(draft.status)">
                {{ draft.status || "pending" }}
              </uiBadge>
            </td>
            <td class="px-5 py-4">
              <span class="text-xs text-[var(--ink-soft)] capitalize">
                {{ draft.source }}
              </span>
            </td>
            <td class="px-5 py-4 text-right">
              <div class="flex items-center justify-end gap-2">
                <uiButton
                  variant="secondary"
                  size="sm"
                  :loading="loading"
                  @click="confirmDraft(draft.id!)"
                >
                  Confirm
                </uiButton>
                <uiButton
                  variant="ghost"
                  size="sm"
                  :loading="loading"
                  @click="rejectDraft(draft.id!)"
                >
                  Reject
                </uiButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </uiCard>
  </div>
</template>
