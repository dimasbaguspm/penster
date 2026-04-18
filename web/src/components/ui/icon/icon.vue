<script setup lang="ts">
import { computed } from "vue";
import { cva } from "class-variance-authority";
import * as LucideIcons from "lucide-vue-next";

const iconVariant = cva("inline-flex items-center justify-center flex-shrink-0", {
  variants: {
    size: {
      xs: "w-3 h-3",
      sm: "w-4 h-4",
      md: "w-5 h-5",
      lg: "w-6 h-6",
      xl: "w-8 h-8",
    },
    strokeWidth: {
      thin: 1,
      normal: 1.5,
      bold: 2,
    },
  },
  defaultVariants: {
    size: "md",
    strokeWidth: "normal",
  },
});

const props = withDefaults(defineProps<{
  name: string;
  size?: "xs" | "sm" | "md" | "lg" | "xl";
  strokeWidth?: "thin" | "normal" | "bold";
}>(), {
  size: "md",
  strokeWidth: "normal",
});

// Convert kebab-case or any case to PascalCase for lucide
const iconName = computed(() => {
  const name = props.name;
  // Handle kebab-case (icon-name) or already PascalCase
  const PascalCase = name.split(/[-_\s]/).map(s =>
    s.charAt(0).toUpperCase() + s.slice(1).toLowerCase()
  ).join('');
  return PascalCase as keyof typeof LucideIcons;
});

const IconComponent = computed(() => {
  return (LucideIcons as any)[iconName.value] || LucideIcons.Circle;
});

const iconClass = computed(() => iconVariant({ size: props.size }));
const strokeWidthValue = computed(() => iconVariant({ strokeWidth: props.strokeWidth }).replace(/.*stroke-\[?/, '').replace(/\]/, '') || '1.5');
</script>

<template>
  <component :is="IconComponent" :size="props.size === 'xs' ? 12 : props.size === 'sm' ? 16 : props.size === 'lg' ? 24 : props.size === 'xl' ? 32 : 20" :stroke-width="strokeWidthValue" :class="iconClass" />
</template>