<script setup lang="ts">
interface Props {
  label?: string;
  type?: "text" | "email" | "password" | "number" | "search" | "date";
  placeholder?: string;
  modelValue?: string | number;
  disabled?: boolean;
  error?: string;
}

defineProps<Props>();
defineEmits<{
  "update:modelValue": [value: string];
}>();
</script>

<template>
  <div class="flex flex-col gap-1">
    <label v-if="label" class="text-xs font-medium text-[var(--ink-soft)]">
      {{ label }}
    </label>
    <input
      :type="type || 'text'"
      :placeholder="placeholder"
      :disabled="disabled"
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      :class="[
        'w-full px-3 py-2 text-sm text-[var(--ink)] bg-[var(--paper)] border rounded-lg transition-colors duration-200',
        'placeholder:text-[var(--ink-soft)]/50',
        'focus:outline-none focus:border-[var(--gold)] focus:ring-1 focus:ring-[var(--gold)]/30',
        'disabled:opacity-50 disabled:cursor-not-allowed',
        error ? 'border-[var(--rust)]' : 'border-[var(--rule)] hover:border-[var(--ink-soft)]/50',
      ]"
    />
    <p v-if="error" class="text-xs text-[var(--rust)]">{{ error }}</p>
  </div>
</template>
