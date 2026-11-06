<script setup lang="ts">
import Button from 'primevue/button';
import { useTheme } from '@/composables/useTheme';
import Sun from '@primeicons/vue/sun';
import Moon from '@primeicons/vue/moon';
const { theme, toggleTheme } = useTheme();
</script>

<template>
  <Button severity="secondary" text @click="toggleTheme">
    <Sun v-if="theme === 'light'" />
    <Moon v-else />
  </Button>
</template>
