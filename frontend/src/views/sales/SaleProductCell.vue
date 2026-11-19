<script setup lang="ts">
  import { Box } from '@lucide/vue';

  defineProps<{
    name: string;
    code: string;
  }>();
</script>

<template>
  <div class="flex items-center gap-2">
    <Box class="size-3.5 shrink-0 text-muted-foreground" />
    <div class="flex min-w-0 flex-col">
      <span class="truncate font-medium">{{ name }}</span>
      <span class="font-mono text-xs text-muted-foreground">{{ code }}</span>
    </div>
  </div>
</template>
