<script setup lang="ts" generic="V extends string | number">
  import { X } from '@lucide/vue';
  import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
  } from '@/components/ui/select';
  import type { SelectOption } from './select';

  /** Single-value select. With `clearable`, an × button resets the value to `null`. */
  const model = defineModel<V | null>({ required: true });

  const {
    id,
    options,
    placeholder,
    clearable = false,
  } = defineProps<{
    id?: string;
    options: readonly SelectOption<V>[];
    placeholder?: string;
    clearable?: boolean;
  }>();
</script>

<template>
  <div class="relative">
    <Select v-model="model">
      <SelectTrigger :id="id" class="w-full" :class="{ 'pr-14': clearable && model !== null }">
        <SelectValue :placeholder="placeholder" />
      </SelectTrigger>
      <SelectContent position="popper">
        <SelectItem v-for="option in options" :key="option.value" :value="option.value">
          {{ option.label }}
        </SelectItem>
      </SelectContent>
    </Select>

    <button
      v-if="clearable && model !== null"
      type="button"
      class="absolute top-1/2 right-8 -translate-y-1/2 rounded-sm p-0.5 text-muted-foreground hover:text-foreground"
      aria-label="Limpiar selección"
      @click="model = null"
    >
      <X class="size-3.5" />
    </button>
  </div>
</template>
