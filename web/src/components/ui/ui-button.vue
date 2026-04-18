<script setup lang="ts">
interface Props {
  variant?: "primary" | "secondary" | "ghost" | "danger";
  size?: "sm" | "md" | "lg";
  loading?: boolean;
  disabled?: boolean;
  type?: "button" | "submit" | "reset";
}

withDefaults(defineProps<Props>(), {
  variant: "primary",
  size: "md",
  loading: false,
  disabled: false,
  type: "button",
});

defineEmits<{
  click: [event: MouseEvent];
}>();
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="[
      'inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-all duration-200 btn-press',
      'disabled:opacity-50 disabled:cursor-not-allowed disabled:pointer-events-none',
      {
        'bg-[var(--ink)] text-[var(--paper)] hover:bg-[var(--ink-soft)]': variant === 'primary',
        'border border-[var(--rule)] text-[var(--ink-soft)] hover:bg-[var(--paper-dark)] hover:text-[var(--ink)]':
          variant === 'secondary',
        'text-[var(--ink-soft)] hover:text-[var(--ink)] hover:bg-[var(--paper-dark)]':
          variant === 'ghost',
        'bg-[var(--rust)] text-white hover:bg-[var(--rust-muted)]': variant === 'danger',
      },
      {
        'px-3 py-1.5 text-xs': size === 'sm',
        'px-4 py-2 text-sm': size === 'md',
        'px-6 py-3 text-base': size === 'lg',
      },
    ]"
  >
    <svg v-if="loading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
      />
    </svg>
    <slot />
  </button>
</template>
