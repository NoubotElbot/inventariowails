<script setup lang="ts">
  import { ShoppingBasket, Trash2 } from '@lucide/vue';
  import DataTable from '@/components/common/DataTable.vue';
  import EmptyState from '@/components/common/EmptyState.vue';
  import IconButton from '@/components/common/IconButton.vue';
  import NumberInput from '@/components/common/NumberInput.vue';
  import type { DataTableColumn } from '@/components/common/data-table';
  import type { CartItem } from '@/types/sale';
  import { formatPrice } from '@/utils/currency';
  import SaleProductCell from './SaleProductCell.vue';

  defineProps<{
    items: readonly CartItem[];
  }>();

  const emit = defineEmits<{
    updateQuantity: [productId: number, quantity: number];
    remove: [productId: number];
  }>();

  const columns: DataTableColumn<CartItem>[] = [
    { key: 'product', header: 'Producto' },
    {
      key: 'price',
      header: 'Precio unit.',
      class: 'w-36',
      value: (i) => formatPrice(i.product.price),
    },
    { key: 'quantity', header: 'Cantidad', class: 'w-40' },
    { key: 'subtotal', header: 'Subtotal', class: 'w-36' },
    { key: 'actions', class: 'w-16' },
  ];
</script>

<template>
  <DataTable
    :rows="items"
    :columns="columns"
    :row-key="(item) => item.product.id"
    table-class="min-w-[44rem]"
  >
    <template #empty>
      <EmptyState :icon="ShoppingBasket" message="Busca un producto para agregarlo a la venta" />
    </template>

    <template #cell-product="{ row }">
      <SaleProductCell :name="row.product.name" :code="row.product.code" />
    </template>

    <template #cell-quantity="{ row }">
      <NumberInput
        class="w-32"
        :model-value="row.quantity"
        :min="1"
        :max="row.product.stock"
        show-buttons
        size="sm"
        :aria-label="`Cantidad de ${row.product.name}`"
        @update:model-value="emit('updateQuantity', row.product.id, $event)"
      />
    </template>

    <template #cell-subtotal="{ row }">
      <span class="font-medium">{{ formatPrice(row.product.price * row.quantity) }}</span>
    </template>

    <template #cell-actions="{ row }">
      <IconButton
        :icon="Trash2"
        tooltip="Quitar"
        :aria-label="`Quitar ${row.product.name}`"
        class="text-destructive hover:text-destructive"
        @click="emit('remove', row.product.id)"
      />
    </template>
  </DataTable>
</template>
