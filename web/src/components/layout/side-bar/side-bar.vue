<script setup lang="ts">
import { RouterLink } from "vue-router";
import { cx } from "@/lib/cx";
import { Icon } from "@/components/ui";
import type { SideBarProps } from "./types";

const props = withDefaults(defineProps<SideBarProps>(), {
  open: false,
  navItems: () => [],
  sections: () => [],
});

const emit = defineEmits<{
  "update:open": [value: boolean];
}>();
</script>

<template>
  <!-- Mobile Backdrop -->
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-300"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-300"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="props.open"
        class="fixed inset-0 bg-[var(--ink)]/50 backdrop-blur-sm z-40 md:hidden"
        @click="emit('update:open', false)"
      />
    </Transition>
  </Teleport>

  <!-- Sidebar Drawer -->
  <aside
    :class="
      cx(
        'bg-[var(--paper)] border-r border-[var(--rule)] flex flex-col',
        'fixed top-0 left-0 h-full z-50 w-80',
        'transform transition-transform duration-300 ease-out',
        props.open ? 'translate-x-0' : '-translate-x-full',
        'md:hidden'
      )
    "
  >
    <!-- Header -->
    <div
      class="flex items-center justify-between h-16 px-6 border-b border-[var(--rule)] bg-[var(--paper)]"
    >
      <div class="flex items-center gap-3">
        <div
          class="w-8 h-8 rounded-lg bg-gradient-to-br from-[var(--gold)]/20 to-[var(--gold)]/5 border border-[var(--gold)]/20 flex items-center justify-center"
        >
          <Icon name="layers" size="sm" class="text-[var(--gold)]" />
        </div>
        <span class="font-display text-lg font-semibold text-[var(--ink)]">Menu</span>
      </div>
      <button
        class="p-2 rounded-xl text-[var(--ink-soft)] hover:bg-[var(--paper-dark)] hover:text-[var(--ink)] transition-all duration-200 active:scale-95"
        @click="emit('update:open', false)"
      >
        <Icon name="x" size="md" />
      </button>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto py-4 px-4">
      <div class="space-y-1">
        <RouterLink
          v-for="item in props.navItems"
          :key="item.to"
          :to="item.to"
          class="group flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] hover:bg-[var(--paper-dark)] transition-all duration-200"
          active-class="!text-[var(--ink)] bg-[var(--paper-dark)] shadow-sm"
          @click="emit('update:open', false)"
        >
          <div
            :class="
              cx(
                'w-9 h-9 rounded-lg bg-[var(--paper-dark)] flex items-center justify-center transition-all duration-200',
                'group-hover:bg-[var(--gold)]/10'
              )
            "
          >
            <Icon
              :name="item.icon ?? 'circle'"
              size="sm"
              :class="
                cx(
                  'text-[var(--ink-soft)] transition-colors duration-200',
                  'group-hover:text-[var(--gold)]'
                )
              "
            />
          </div>
          <span class="flex-1">{{ item.label }}</span>
          <span
            v-if="item.badge"
            class="inline-flex items-center justify-center min-w-[1.5rem] h-6 px-2 rounded-full bg-[var(--gold)]/15 text-xs font-mono font-semibold text-[var(--gold)]"
          >
            {{ item.badge }}
          </span>
        </RouterLink>
      </div>

      <!-- Sections -->
      <template v-for="section in props.sections" :key="section.title">
        <div class="pt-6 mt-6 border-t border-[var(--rule)]">
          <h3
            v-if="section.title"
            class="px-4 mb-3 text-xs font-semibold uppercase tracking-wider text-[var(--ink-soft)]/60"
          >
            {{ section.title }}
          </h3>
          <div class="space-y-1">
            <RouterLink
              v-for="item in section.items"
              :key="item.to"
              :to="item.to"
              class="group flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] hover:bg-[var(--paper-dark)] transition-all duration-200"
              active-class="!text-[var(--ink)] bg-[var(--paper-dark)] shadow-sm"
              @click="emit('update:open', false)"
            >
              <div
                class="w-9 h-9 rounded-lg bg-[var(--paper-dark)] flex items-center justify-center transition-all duration-200 group-hover:bg-[var(--gold)]/10"
              >
                <Icon
                  :name="item.icon ?? 'circle'"
                  size="sm"
                  class="text-[var(--ink-soft)] group-hover:text-[var(--gold)] transition-colors duration-200"
                />
              </div>
              <span class="flex-1">{{ item.label }}</span>
            </RouterLink>
          </div>
        </div>
      </template>
    </nav>

    <!-- Footer -->
    <div class="px-4 py-6 border-t border-[var(--rule)]">
      <p class="text-xs text-center text-[var(--ink-soft)]/60">Penster v1.0</p>
    </div>
  </aside>
</template>