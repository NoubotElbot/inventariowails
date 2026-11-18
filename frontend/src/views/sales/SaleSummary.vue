<script setup lang="ts">
  import { computed } from 'vue';
  import Calendar from '@primeicons/vue/calendar';
  import Box from '@primeicons/vue/box';
  import Wallet from '@primeicons/vue/wallet';
  import { formatPrice } from '@/utils/currency';
  import { formatDateTime } from '@/utils/date';
  import { totalUnits } from '@/utils/sale';
  import type { Sale } from '@/types/sale';

  const { sale } = defineProps<{
    sale: Sale;
  }>();

  const stats = computed(() => [
    { label: 'Fecha', value: formatDateTime(sale.created_at), icon: Calendar },
    { label: 'Unidades', value: String(totalUnits(sale.items)), icon: Box },
    { label: 'Total', value: formatPrice(sale.total), icon: Wallet },
  ]);
</script>

<template>
  <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
    <div
      v-for="stat in stats"
      :key="stat.label"
      class="flex items-center gap-3 rounded-lg border border-surface-200 px-4 py-3 dark:border-surface-700"
    >
      <component :is="stat.icon" :size="18" class="text-surface-400" />
      <div class="flex flex-col">
        <span class="text-xs text-surface-500">{{ stat.label }}</span>
        <span class="font-semibold text-surface-800 dark:text-surface-100">{{ stat.value }}</span>
      </div>
    </div>
  </div>
</template>
