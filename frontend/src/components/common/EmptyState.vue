<script setup lang="ts">
  import type { Component } from 'vue';
  import { Empty, EmptyDescription, EmptyHeader, EmptyMedia } from '@/components/ui/empty';

  /** Icon + message placeholder for empty lists and tables. */
  defineProps<{
    icon: Component;
    message: string;
  }>();
</script>

<template>
  <Empty class="p-6">
    <EmptyHeader>
      <EmptyMedia variant="icon">
        <component :is="icon" />
      </EmptyMedia>
      <EmptyDescription>{{ message }}</EmptyDescription>
    </EmptyHeader>
    <slot />
  </Empty>
</template>
