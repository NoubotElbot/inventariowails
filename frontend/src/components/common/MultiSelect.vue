<script setup lang="ts" generic="V extends string | number">
  import { computed } from 'vue';
  import { X } from '@lucide/vue';
  import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
  } from '@/components/ui/select';
  import type { SelectOption } from './select';

  /**
   * Multi-value select with checkmarks and a clear button. Past `maxLabels`
   * selected items the trigger shows a count instead of the labels.
   */
  const {
    modelValue,
    options,
    placeholder,
    maxLabels = 2,
  } = defineProps<{
    modelValue: readonly V[];
    options: readonly SelectOption<V>[];
    placeholder?: string;
    maxLabels?: number;
  }>();

  const emit = defineEmits<{
    'update:modelValue': [value: V[]];
  }>();

  const summary = computed(() => {
    if (modelValue.length === 0) return '';
    if (modelValue.length > maxLabels) return `${modelValue.length} seleccionadas`;
    return options
      .filter((option) => modelValue.includes(option.value))
      .map((option) => option.label)
      .join(', ');
  });

  function update(value: unknown): void {
    emit('update:modelValue', [...(value as V[])]);
  }
</script>

<template>
  <div class="relative">
    <Select :model-value="[...modelValue]" multiple @update:model-value="update">
      <SelectTrigger class="w-full" :class="{ 'pr-14': modelValue.length > 0 }">
        <SelectValue :placeholder="placeholder">{{ summary || placeholder }}</SelectValue>
      </SelectTrigger>
      <SelectContent position="popper">
        <SelectItem v-for="option in options" :key="option.value" :value="option.value">
          {{ option.label }}
        </SelectItem>
      </SelectContent>
    </Select>

    <button
      v-if="modelValue.length > 0"
      type="button"
      class="absolute top-1/2 right-8 -translate-y-1/2 rounded-sm p-0.5 text-muted-foreground hover:text-foreground"
      aria-label="Limpiar selección"
      @click="emit('update:modelValue', [])"
    >
      <X class="size-3.5" />
    </button>
  </div>
</template>
