<script setup lang="ts">
import { cva } from "class-variance-authority";
import type { BadgeProps } from "./types";

const badgeVariant = cva(
  "inline-flex items-center px-2 py-0.5 rounded text-xs font-medium border border-[var(--rule)] capitalize",
  {
    variants: {
      variant: {
        default: "text-[var(--ink-soft)]",
        teal: "text-[var(--teal)] border-[var(--teal-muted)]/30 bg-[var(--teal)]/5",
        rust: "text-[var(--rust)] border-[var(--rust-muted)]/30 bg-[var(--rust)]/5",
        gold: "text-[var(--gold)] border-[var(--gold-light)]/30 bg-[var(--gold-light)]/10",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  },
);

const props = withDefaults(defineProps<BadgeProps>(), {
  variant: "default",
});
</script>

<template>
  <span :class="badgeVariant({ variant: props.variant })">
    <slot />
  </span>
</template>
