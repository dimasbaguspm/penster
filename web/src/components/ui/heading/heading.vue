<script setup lang="ts">
import { cva } from "class-variance-authority";
import type { HeadingProps } from "./types";

const headingVariant = cva("font-display font-semibold text-[var(--ink)] leading-tight", {
  variants: {
    size: {
      xs: "text-xs",
      sm: "text-sm",
      md: "text-base",
      lg: "text-xl",
      xl: "text-2xl",
      "2xl": "text-3xl",
      "3xl": "text-4xl",
      "4xl": "text-5xl",
    },
    mb: {
      "0": "mb-0",
      "1": "mb-1",
      "2": "mb-2",
      "3": "mb-3",
      "4": "mb-4",
      "6": "mb-6",
      "8": "mb-8",
    },
  },
  defaultVariants: {
    size: "lg",
  },
});

const props = withDefaults(defineProps<HeadingProps>(), {
  as: "h2",
  size: "lg",
  mb: undefined,
});
</script>

<template>
  <component
    :is="props.as"
    :class="headingVariant({ size: props.size, mb: props.mb })"
  >
    <slot />
  </component>
</template>
