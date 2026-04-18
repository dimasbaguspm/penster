<script setup lang="ts">
import { RouterLink } from "vue-router";
import { cx } from "@/lib/cx";
import { Icon } from "@/components/ui";
import type { TopBarProps } from "./types";

const props = withDefaults(defineProps<TopBarProps>(), {
  logo: () => ({ text: "Penster", href: "/" }),
  navItems: () => [],
  showMobileMenu: false,
  sticky: true,
});

defineEmits<{ (e: "toggle-mobile-menu"): void }>();
</script>

<template>
  <header
    :class="
      cx(
        'w-full border-b border-[var(--rule)] bg-[var(--paper)]/95 backdrop-blur-xl',
        'transition-all duration-300',
        props.sticky && 'sticky top-0 z-20'
      )
    "
  >
    <div class="w-full px-6 lg:px-10">
      <div class="flex items-center justify-between h-16">
        <!-- Logo -->
        <slot name="logo">
          <RouterLink :to="props.logo?.href || '/'" class="flex items-center gap-3 group">
            <div
              class="w-9 h-9 rounded-xl bg-gradient-to-br from-[var(--gold)]/20 to-[var(--gold)]/5 border border-[var(--gold)]/20 flex items-center justify-center shadow-sm"
            >
              <Icon name="layers" size="sm" class="text-[var(--gold)]" />
            </div>
            <span
              class="font-display text-xl font-semibold tracking-tight text-[var(--ink)] group-hover:text-[var(--gold)] transition-colors duration-200"
            >
              {{ props.logo?.text }}
            </span>
          </RouterLink>
        </slot>

        <!-- Desktop Navigation -->
        <nav class="hidden md:flex items-center gap-2">
          <slot name="nav">
            <RouterLink
              v-for="item in props.navItems"
              :key="item.to"
              :to="item.to"
              class="group flex items-center gap-2 px-4 py-2 text-sm font-medium text-[var(--ink-soft)] hover:text-[var(--ink)] transition-all duration-200 rounded-xl hover:bg-[var(--paper-dark)]"
              active-class="!text-[var(--ink)] bg-[var(--paper-dark)]"
            >
              <Icon
                :name="item.icon ?? 'layout-dashboard'"
                size="sm"
                class="opacity-60 group-hover:opacity-100 transition-opacity"
              />
              <span>{{ item.label }}</span>
            </RouterLink>
          </slot>
        </nav>

        <!-- Right side actions -->
        <div class="flex items-center gap-2">
          <!-- Mobile menu button -->
          <button
            class="md:hidden p-2.5 rounded-xl text-[var(--ink-soft)] hover:bg-[var(--paper-dark)] hover:text-[var(--ink)] transition-all duration-200 active:scale-95"
            @click="$emit('toggle-mobile-menu')"
            aria-label="Toggle navigation menu"
          >
            <slot name="menu-icon">
              <Icon :name="props.showMobileMenu ? 'x' : 'menu'" size="md" />
            </slot>
          </button>
        </div>
      </div>
    </div>
  </header>
</template>