<script setup lang="ts">
import { cva } from "class-variance-authority";
import type { TextProps } from "./types";

const textVariant = cva("text-[var(--ink)]", {
  variants: {
    size: {
      xs: "text-xs",
      sm: "text-sm",
      base: "text-base",
      lg: "text-lg",
    },
    muted: {
      true: "text-[var(--ink-soft)]",
    },
    mt: {
      "0": "mt-0",
      "1": "mt-1",
      "2": "mt-2",
      "3": "mt-3",
      "4": "mt-4",
      "6": "mt-6",
      "8": "mt-8",
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
    size: "base",
    muted: false,
  },
});

const props = withDefaults(defineProps<TextProps>(), {
  as: "p",
  size: "base",
  muted: false,
  mt: undefined,
  mb: undefined,
});
</script>

<template>
  <component
    :is="props.as"
    :class="textVariant({ size: props.size, muted: props.muted, mt: props.mt, mb: props.mb })"
  >
    <slot />
  </component>
</template>
