<script setup lang="ts">
import { cva } from "class-variance-authority";
import type { InputProps, InputEmits } from "./types";

const inputVariant = cva(
  "w-full px-3 py-2 text-sm text-[var(--ink)] bg-[var(--paper)] border rounded-lg transition-colors duration-200 placeholder:text-[var(--ink-soft)]/50 focus:outline-none focus:border-[var(--gold)] focus:ring-1 focus:ring-[var(--gold)]/30 disabled:opacity-50 disabled:cursor-not-allowed",
  {
    variants: {
      error: {
        true: "border-[var(--rust)]",
        false: "border-[var(--rule)] hover:border-[var(--ink-soft)]/50",
      },
    },
    defaultVariants: {
      error: false,
    },
  },
);

const props = defineProps<InputProps>();
defineEmits<InputEmits>();
</script>

<template>
  <div class="flex flex-col gap-1">
    <label v-if="props.label" class="text-xs font-medium text-[var(--ink-soft)]">
      {{ props.label }}
    </label>
    <input
      :type="props.type || 'text'"
      :placeholder="props.placeholder"
      :disabled="props.disabled"
      :value="props.modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      :class="inputVariant({ error: !!props.error })"
    />
    <p v-if="props.error" class="text-xs text-[var(--rust)]">{{ props.error }}</p>
  </div>
</template>
