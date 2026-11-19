<script setup lang="ts">
  import type { Component } from 'vue';
  import { Button, type ButtonVariants } from '@/components/ui/button';
  import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';

  /**
   * Compact icon-only button with a tooltip. The tooltip text doubles as the
   * accessible name unless `ariaLabel` gives a more specific one. Attributes
   * and listeners (`class`, `@click`, `disabled`…) go to the button.
   */
  defineOptions({ inheritAttrs: false });

  const {
    icon,
    tooltip,
    ariaLabel,
    variant = 'ghost',
  } = defineProps<{
    icon: Component;
    tooltip: string;
    ariaLabel?: string;
    variant?: ButtonVariants['variant'];
  }>();
</script>

<template>
  <Tooltip>
    <TooltipTrigger as-child>
      <Button v-bind="$attrs" :variant="variant" size="icon-sm" :aria-label="ariaLabel ?? tooltip">
        <component :is="icon" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>{{ tooltip }}</TooltipContent>
  </Tooltip>
</template>
