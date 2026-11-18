<script setup lang="ts">
  import Check from '@primeicons/vue/check';
  import { formatPrice } from '@/utils/currency';

  defineProps<{
    units: number;
    total: number;
    loading: boolean;
    disabled: boolean;
  }>();

  const emit = defineEmits<{
    cancel: [];
    confirm: [];
  }>();
</script>

<template>
  <div
    class="flex flex-wrap items-center justify-between gap-4 rounded-lg border border-surface-200 px-4 py-3 dark:border-surface-700"
  >
    <div class="flex flex-col">
      <span class="text-sm text-surface-500">{{ units }} unidades</span>
      <span class="text-2xl font-semibold text-surface-800 dark:text-surface-100">
        {{ formatPrice(total) }}
      </span>
    </div>

    <div class="flex gap-2">
      <Button label="Cancelar" severity="secondary" :disabled="loading" @click="emit('cancel')" />
      <Button :loading="loading" :disabled="disabled" @click="emit('confirm')">
        <Check v-if="!loading" :size="16" />
        Confirmar venta
      </Button>
    </div>
  </div>
</template>
