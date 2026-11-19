<script setup lang="ts">
  import { computed } from 'vue';
  import { Minus, Plus } from '@lucide/vue';
  import { cn } from '@/lib/utils';
  import {
    NumberField,
    NumberFieldContent,
    NumberFieldDecrement,
    NumberFieldIncrement,
    NumberFieldInput,
  } from '@/components/ui/number-field';

  /**
   * Numeric input with optional +/- buttons. `currency` formats the value
   * as Chilean pesos (no decimals) using the es-CL locale.
   * Extra attributes (`aria-label`, `placeholder`…) go to the <input>.
   */
  defineOptions({ inheritAttrs: false });

  const model = defineModel<number>({ required: true });

  const {
    id,
    min,
    max,
    step = 1,
    currency = false,
    showButtons = false,
    size = 'default',
  } = defineProps<{
    id?: string;
    min?: number;
    max?: number;
    step?: number;
    currency?: boolean;
    showButtons?: boolean;
    size?: 'sm' | 'default';
  }>();

  const formatOptions = computed<Intl.NumberFormatOptions>(() =>
    currency
      ? { style: 'currency', currency: 'CLP', maximumFractionDigits: 0 }
      : { maximumFractionDigits: 0, useGrouping: false }
  );

  // Clearing the input emits undefined/NaN; keep the last valid value instead.
  function onUpdate(value: number | undefined): void {
    if (value !== undefined && !Number.isNaN(value)) model.value = value;
  }
</script>

<template>
  <NumberField
    :id="id"
    :model-value="model"
    :min="min"
    :max="max"
    :step="step"
    :format-options="formatOptions"
    locale="es-CL"
    @update:model-value="onUpdate"
  >
    <NumberFieldContent>
      <NumberFieldDecrement v-if="showButtons" :class="cn(size === 'sm' && 'p-2')">
        <Minus class="size-3.5" />
      </NumberFieldDecrement>
      <NumberFieldInput
        v-bind="$attrs"
        :class="cn(!showButtons && 'text-left', size === 'sm' && 'h-8')"
      />
      <NumberFieldIncrement v-if="showButtons" :class="cn(size === 'sm' && 'p-2')">
        <Plus class="size-3.5" />
      </NumberFieldIncrement>
    </NumberFieldContent>
  </NumberField>
</template>
