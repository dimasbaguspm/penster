<script setup lang="ts">
interface Props {
  variant?: "default" | "teal" | "rust" | "gold";
}

withDefaults(defineProps<Props>(), {
  variant: "default",
});
</script>

<template>
  <span
    :class="[
      'inline-flex items-center px-2 py-0.5 rounded text-xs font-medium border border-[var(--rule)] capitalize',
      {
        'text-[var(--ink-soft)]': variant === 'default',
        'text-[var(--teal)] border-[var(--teal-muted)]/30 bg-[var(--teal)]/5': variant === 'teal',
        'text-[var(--rust)] border-[var(--rust-muted)]/30 bg-[var(--rust)]/5': variant === 'rust',
        'text-[var(--gold)] border-[var(--gold-light)]/30 bg-[var(--gold-light)]/10':
          variant === 'gold',
      },
    ]"
  >
    <slot />
  </span>
</template>
