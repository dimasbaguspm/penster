<script setup lang="ts">
import { cva } from "class-variance-authority";
import { Icon } from "@/components/ui";
import type { ButtonEmits, ButtonProps } from "./types";

const buttonVariant = cva(
  "inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-all duration-200 btn-press disabled:opacity-50 disabled:cursor-not-allowed disabled:pointer-events-none",
  {
    variants: {
      variant: {
        primary: "bg-[var(--ink)] text-[var(--paper)] hover:bg-[var(--ink-soft)]",
        secondary:
          "border border-[var(--rule)] text-[var(--ink-soft)] hover:bg-[var(--paper-dark)] hover:text-[var(--ink)]",
        ghost: "text-[var(--ink-soft)] hover:text-[var(--ink)] hover:bg-[var(--paper-dark)]",
        danger: "bg-[var(--rust)] text-white hover:bg-[var(--rust-muted)]",
      },
      size: {
        sm: "px-3 py-1.5 text-xs",
        md: "px-4 py-2 text-sm",
        lg: "px-6 py-3 text-base",
      },
    },
    defaultVariants: {
      variant: "primary",
      size: "md",
    },
  },
);

const props = withDefaults(defineProps<ButtonProps>(), {
  variant: "primary",
  size: "md",
  loading: false,
  disabled: false,
  type: "button",
});

defineEmits<ButtonEmits>();
</script>

<template>
  <button
    :type="props.type"
    :disabled="props.disabled || props.loading"
    :class="buttonVariant({ variant: props.variant, size: props.size })"
  >
    <Icon v-if="props.loading" name="loader-2" size="sm" class="animate-spin" />
    <slot />
  </button>
</template>