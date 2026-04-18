<script setup lang="ts">
interface Props {
  hover?: boolean;
}

withDefaults(defineProps<Props>(), {
  hover: false,
});
</script>

<template>
  <div
    :class="[
      'bg-[var(--paper)] border border-[var(--rule)] rounded-lg overflow-hidden',
      hover ? 'card-hover cursor-pointer' : '',
    ]"
  >
    <slot />
  </div>
</template>
