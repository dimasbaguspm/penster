<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router";
</script>

<template>
  <div class="min-h-screen">
    <header
      class="border-b border-[var(--rule)] bg-[var(--paper)]/80 backdrop-blur-sm sticky top-0 z-20 animate-fade-in"
    >
      <div class="max-w-7xl mx-auto px-6 lg:px-10">
        <div class="flex items-center justify-between h-14">
          <div class="flex items-center gap-3">
            <svg
              class="w-5 h-5 text-[var(--gold)]"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
            >
              <path d="M12 2L2 7l10 5 10-5-10-5z" stroke-linecap="round" stroke-linejoin="round" />
              <path d="M2 17l10 5 10-5" stroke-linecap="round" stroke-linejoin="round" />
              <path d="M2 12l10 5 10-5" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
            <h1 class="font-display text-2xl font-semibold tracking-tight text-[var(--ink)]">
              Penster
            </h1>
          </div>

          <nav class="hidden md:flex items-center gap-1">
            <RouterLink
              to="/"
              class="px-4 py-1.5 text-sm font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] transition-colors duration-200 rounded-md hover:bg-[var(--paper-dark)] ink-underline"
              active-class="!text-[var(--ink)] bg-[var(--paper-dark)]"
            >
              Dashboard
            </RouterLink>
            <RouterLink
              to="/accounts"
              class="px-4 py-1.5 text-sm font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] transition-colors duration-200 rounded-md hover:bg-[var(--paper-dark)] ink-underline"
              active-class="!text-[var(--ink)] bg-[var(--paper-dark)]"
            >
              Accounts
            </RouterLink>
            <RouterLink
              to="/transactions"
              class="px-4 py-1.5 text-sm font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] transition-colors duration-200 rounded-md hover:bg-[var(--paper-dark)] ink-underline"
              active-class="!text-[var(--ink)] bg-[var(--paper-dark)]"
            >
              Transactions
            </RouterLink>
            <RouterLink
              to="/drafts"
              class="px-4 py-1.5 text-sm font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] transition-colors duration-200 rounded-md hover:bg-[var(--paper-dark)] ink-underline"
              active-class="!text-[var(--ink)] bg-[var(--paper-dark)]"
            >
              Drafts
            </RouterLink>
            <RouterLink
              to="/reports"
              class="px-4 py-1.5 text-sm font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] transition-colors duration-200 rounded-md hover:bg-[var(--paper-dark)] ink-underline"
              active-class="!text-[var(--ink)] bg-[var(--paper-dark)]"
            >
              Reports
            </RouterLink>
          </nav>
        </div>
      </div>
    </header>

    <RouterView />

    <footer class="border-t border-[var(--rule)] mt-16">
      <div class="max-w-7xl mx-auto px-6 lg:px-10 py-6 flex items-center justify-between">
        <p class="text-xs text-[var(--ink-soft)]">Penster — Personal Finance Ledger</p>
        <p class="text-xs font-mono text-[var(--ink-soft)]">All figures in USD</p>
      </div>
    </footer>
  </div>
</template>
