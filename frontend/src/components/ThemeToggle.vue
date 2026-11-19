<script setup lang="ts">
  import { Moon, Sun } from '@lucide/vue';
  import IconButton from '@/components/common/IconButton.vue';
  import { useTheme } from '@/composables/useTheme';

  const { theme, toggleTheme } = useTheme();
</script>

<template>
  <IconButton
    :icon="theme === 'light' ? Sun : Moon"
    :tooltip="theme === 'light' ? 'Cambiar a tema oscuro' : 'Cambiar a tema claro'"
    @click="toggleTheme"
  />
</template>
