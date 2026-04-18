<script setup lang="ts">
import { Icon } from "@/components/ui";
import type { ModalProps, ModalEmits } from './types';

defineProps<ModalProps>();
defineEmits<ModalEmits>();
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
      <div
        v-if="$props.open"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="$props.title ? 'modal-title' : undefined"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-[var(--ink)]/40 backdrop-blur-sm" @click="$emit('close')" />

        <!-- Dialog -->
        <div
          class="relative w-full max-w-md bg-[var(--paper)] border border-[var(--rule)] rounded-xl shadow-2xl animate-fade-up"
        >
          <div class="flex items-center justify-between px-6 py-4 border-b border-[var(--rule)]">
            <h3
              v-if="$props.title"
              id="modal-title"
              class="font-display text-lg font-semibold text-[var(--ink)]"
            >
              {{ $props.title }}
            </h3>
            <button
              class="ml-auto p-1 rounded-md text-[var(--ink-soft)] hover:bg-[var(--paper-dark)] transition-colors"
              @click="$emit('close')"
              aria-label="Close dialog"
            >
              <Icon name="x" size="md" />
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