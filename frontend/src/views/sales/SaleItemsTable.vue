<script setup lang="ts">
  import DataTable from '@/components/common/DataTable.vue';
  import type { DataTableColumn } from '@/components/common/data-table';
  import { TableCell, TableRow } from '@/components/ui/table';
  import type { SaleItem } from '@/types/sale';
  import { formatPrice } from '@/utils/currency';
  import SaleProductCell from './SaleProductCell.vue';

  defineProps<{
    items: SaleItem[];
    total: number;
  }>();

  const columns: DataTableColumn<SaleItem>[] = [
    { key: 'product', header: 'Producto' },
    {
      key: 'unit_price',
      header: 'Precio unit.',
      class: 'w-36',
      value: (i) => formatPrice(i.unit_price),
    },
    { key: 'quantity', header: 'Cantidad', class: 'w-28' },
    { key: 'subtotal', header: 'Subtotal', class: 'w-36' },
  ];
</script>

<template>
  <DataTable
    :rows="items"
    :columns="columns"
    :row-key="(item) => item.id"
    striped
    table-class="min-w-[40rem]"
  >
    <template #cell-product="{ row }">
      <SaleProductCell
        :name="row.product?.name ?? `Producto #${row.product_id}`"
        :code="row.product?.code ?? '—'"
      />
    </template>

    <template #cell-subtotal="{ row }">
      <span class="font-medium">{{ formatPrice(row.subtotal) }}</span>
    </template>

    <template #footer>
      <TableRow class="hover:bg-transparent">
        <TableCell :colspan="3" class="px-3 py-2.5 text-right">Total</TableCell>
        <TableCell class="px-3 py-2.5 font-semibold">{{ formatPrice(total) }}</TableCell>
      </TableRow>
    </template>
  </DataTable>
</template>
