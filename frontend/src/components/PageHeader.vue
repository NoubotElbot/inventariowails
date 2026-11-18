<script setup lang="ts">
  import type { Component } from 'vue';

  defineProps<{
    title: string;
    icon: Component;
  }>();
</script>

<template>
  <div class="flex items-center justify-between gap-3">
    <div class="flex items-center gap-3">
      <component :is="icon" :size="20" class="text-primary" />
      <h1 class="text-xl font-semibold text-surface-800 dark:text-surface-100">{{ title }}</h1>
    </div>
    <div v-if="$slots.actions" class="flex items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>
