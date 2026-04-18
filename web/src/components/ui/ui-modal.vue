<script setup lang="ts">
interface Props {
  open: boolean;
  title?: string;
}

defineProps<Props>();
defineEmits<{
  close: [];
}>();
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-[var(--ink)]/40 backdrop-blur-sm" @click="$emit('close')" />

        <!-- Dialog -->
        <div
          class="relative w-full max-w-md bg-[var(--paper)] border border-[var(--rule)] rounded-xl shadow-2xl animate-fade-up"
        >
          <div class="flex items-center justify-between px-6 py-4 border-b border-[var(--rule)]">
            <h3 v-if="title" class="font-display text-lg font-semibold text-[var(--ink)]">
              {{ title }}
            </h3>
            <button
              class="ml-auto p-1 rounded-md text-[var(--ink-soft)] hover:bg-[var(--paper-dark)] transition-colors"
              @click="$emit('close')"
            >
              <svg
                class="w-5 h-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="px-6 py-4">
            <slot />
          </div>

          <div
            v-if="$slots.footer"
            class="px-6 py-4 border-t border-[var(--rule)] flex items-center justify-end gap-3"
          >
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
