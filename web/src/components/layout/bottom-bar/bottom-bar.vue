<script setup lang="ts">
import { RouterLink } from "vue-router";
import { cx } from "@/lib/cx";
import { Icon } from "@/components/ui";
import type { BottomBarProps } from "./types";

const props = withDefaults(defineProps<BottomBarProps>(), {
  items: () => [
    { label: "Dashboard", to: "/", icon: "layout-dashboard" },
    { label: "Accounts", to: "/accounts", icon: "wallet" },
    { label: "Add", to: "/transactions/new", icon: "plus" },
    { label: "Drafts", to: "/drafts", icon: "file-text" },
    { label: "Reports", to: "/reports", icon: "bar-chart-2" },
  ],
});
</script>

<template>
  <nav class="md:hidden fixed bottom-0 left-0 right-0 z-20 bg-[var(--paper)] border-t border-[var(--rule)] safe-area-pb">
    <div class="flex items-center justify-around h-16 px-1">
      <RouterLink
        v-for="(item, index) in props.items"
        :key="item.to"
        :to="item.to"
        :class="
          cx(
            'flex flex-col items-center justify-center gap-0.5 rounded-xl transition-all duration-200',
            index === 2
              ? 'text-[var(--paper)] bg-[var(--ink)] shadow-lg shadow-[var(--ink)]/20 rounded-full w-12 h-12 mx-2 hover:scale-105 active:scale-95'
              : 'text-[var(--ink-soft)] hover:text-[var(--ink)] w-14 h-14 hover:bg-[var(--paper-dark)]'
          )
        "
        active-class="!text-[var(--ink)]"
      >
        <!-- Center FAB for "Add" -->
        <template v-if="index === 2">
          <Icon :name="item.icon" :size="index === 2 ? 'lg' : 'sm'" />
        </template>
        <template v-else>
          <Icon :name="item.icon" size="sm" />
          <span class="text-[10px] font-medium">{{ item.label }}</span>
        </template>
      </RouterLink>
    </div>
  </nav>
</template>