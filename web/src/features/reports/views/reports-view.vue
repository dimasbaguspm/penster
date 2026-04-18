<script setup lang="ts">
import uiCard from "@/components/ui/ui-card.vue";
import uiButton from "@/components/ui/ui-button.vue";
import uiInput from "@/components/ui/ui-input.vue";
</script>

<template>
  <div class="max-w-7xl mx-auto px-6 lg:px-10 py-10">
    <div class="mb-8 animate-fade-up">
      <h1 class="font-display text-3xl font-semibold text-[var(--ink)]">Reports</h1>
      <p class="text-sm text-[var(--ink-soft)] mt-1">Spending trends and financial insights</p>
    </div>

    <!-- Date range filter -->
    <div class="flex flex-col sm:flex-row gap-4 mb-8 animate-fade-up delay-1 items-end">
      <uiInput label="Start Date" type="date" />
      <uiInput label="End Date" type="date" />
      <uiButton variant="secondary">Apply</uiButton>
    </div>

    <uiCard>
      <div class="p-12 text-center animate-fade-up delay-2">
        <svg
          class="w-10 h-10 mx-auto text-[var(--rule)] mb-4"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1"
        >
          <path
            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
        <h3 class="font-display text-lg font-medium text-[var(--ink)] mb-1">No report data</h3>
        <p class="text-sm text-[var(--ink-soft)]">
          Add transactions to see spending breakdowns and trends.
        </p>
      </div>
    </uiCard>
  </div>
</template>
