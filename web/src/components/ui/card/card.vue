<script setup lang="ts">
import { cva } from "class-variance-authority";
import type { CardProps } from "./types";

const cardVariant = cva(
  "bg-[var(--paper)] border border-[var(--rule)] rounded-lg overflow-hidden",
  {
    variants: {
      hover: {
        true: "card-hover cursor-pointer",
      },
    },
    defaultVariants: {
      hover: false,
    },
  },
);

const props = withDefaults(defineProps<CardProps>(), {
  hover: false,
});
</script>

<template>
  <div :class="cardVariant({ hover: props.hover })">
    <slot />
  </div>
</template>
