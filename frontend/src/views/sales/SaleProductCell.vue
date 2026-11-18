<script setup lang="ts">
  import Box from '@primeicons/vue/box';

  defineProps<{
    name: string;
    code: string;
  }>();
</script>

<template>
  <div class="flex items-center gap-2">
    <Box :size="14" class="shrink-0 text-surface-400" />
    <div class="flex min-w-0 flex-col">
      <span class="truncate font-medium">{{ name }}</span>
      <span class="font-mono text-xs text-surface-500">{{ code }}</span>
    </div>
  </div>
</template>
