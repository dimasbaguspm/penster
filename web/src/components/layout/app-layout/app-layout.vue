<script setup lang="ts">
import { useIsMobile } from "@/composables/use-media-query";
import { cx } from "@/lib/cx";
import { computed, ref } from "vue";
import { defaultNavItems } from "../config";
import type { AppLayoutProps } from "./types";
import SideBar from "../side-bar/side-bar.vue";
import TopBar from "../top-bar/top-bar.vue";

const props = withDefaults(defineProps<AppLayoutProps>(), {
  topBar: () => ({}) as any,
  sideBar: () => ({ enabled: true }),
  bottomBar: () => ({ enabled: true }),
  contentMaxWidth: "7xl",
  footerEnabled: true,
});

const isMobile = useIsMobile();
const mobileMenuOpen = ref(false);

const contentMaxWidthClass = computed(
  () =>
    ({
      sm: "max-w-screen-sm",
      md: "max-w-screen-md",
      lg: "max-w-screen-lg",
      xl: "max-w-screen-xl",
      "2xl": "max-w-screen-2xl",
      "7xl": "max-w-7xl",
      full: "max-w-full",
    })[props.contentMaxWidth || "7xl"],
);
</script>

<template>
  <div class="min-h-screen flex flex-col bg-[var(--paper)]">
    <!-- TopBar -->
    <TopBar
      :logo="props.topBar?.logo"
      :nav-items="props.topBar?.navItems || defaultNavItems"
      :show-mobile-menu="mobileMenuOpen"
      :sticky="props.topBar?.sticky ?? true"
      @toggle-mobile-menu="mobileMenuOpen = !mobileMenuOpen"
    >
      <template v-if="$slots['top-bar-logo']" #logo>
        <slot name="top-bar-logo" />
      </template>
      <template v-if="$slots['top-bar-nav']" #nav>
        <slot name="top-bar-nav" />
      </template>
      <template v-if="$slots['top-bar-menu-icon']" #menu-icon>
        <slot name="top-bar-menu-icon" />
      </template>
    </TopBar>

    <!-- Sidebar (mobile drawer) -->
    <SideBar
      v-if="props.sideBar?.enabled"
      :open="mobileMenuOpen"
      :nav-items="props.sideBar?.navItems || defaultNavItems"
      :sections="props.sideBar?.sections"
      @update:open="mobileMenuOpen = $event"
    />

    <!-- Main container -->
    <div class="flex flex-1 relative">
      <main
        :class="cx('flex-1 flex flex-col min-w-0', props.bottomBar?.enabled && isMobile && 'pb-20')"
      >
        <div :class="cx('w-full mx-auto px-6 lg:px-10 py-8', contentMaxWidthClass)">
          <slot />
        </div>

        <!-- Footer slot -->
        <div v-if="$slots.footer || props.footerEnabled" class="mt-auto">
          <slot name="footer">
            <footer class="border-t border-[var(--rule)] mt-16">
              <div class="max-w-7xl mx-auto px-6 lg:px-10 py-6 flex items-center justify-between">
                <p class="text-xs text-[var(--ink-soft)]">Penster — Personal Finance Ledger</p>
                <p class="text-xs font-mono text-[var(--ink-soft)]">All figures in USD</p>
              </div>
            </footer>
          </slot>
        </div>
      </main>
    </div>
  </div>
</template>