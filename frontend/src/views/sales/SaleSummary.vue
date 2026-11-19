<script setup lang="ts">
  import { computed } from 'vue';
  import { Box, Calendar, Wallet } from '@lucide/vue';
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
      class="flex items-center gap-3 rounded-lg border px-4 py-3"
    >
      <component :is="stat.icon" class="size-4.5 text-muted-foreground" />
      <div class="flex flex-col">
        <span class="text-xs text-muted-foreground">{{ stat.label }}</span>
        <span class="font-semibold text-foreground">{{ stat.value }}</span>
      </div>
    </div>
  </div>
</template>
