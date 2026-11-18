<script setup lang="ts">
  import { formatPrice } from '@/utils/currency';
  import type { SaleItem } from '@/types/sale';
  import SaleProductCell from './SaleProductCell.vue';

  defineProps<{
    items: SaleItem[];
    total: number;
  }>();
</script>

<template>
  <DataTable :value="items" data-key="id" striped-rows table-style="min-width: 40rem">
    <Column header="Producto">
      <template #body="{ data }">
        <SaleProductCell
          :name="data.product?.name ?? `Producto #${data.product_id}`"
          :code="data.product?.code ?? '—'"
        />
      </template>
    </Column>

    <Column header="Precio unit." style="width: 9rem">
      <template #body="{ data }">{{ formatPrice(data.unit_price) }}</template>
    </Column>

    <Column field="quantity" header="Cantidad" style="width: 7rem" />

    <Column header="Subtotal" style="width: 9rem">
      <template #body="{ data }">
        <span class="font-medium">{{ formatPrice(data.subtotal) }}</span>
      </template>
    </Column>

    <ColumnGroup type="footer">
      <Row>
        <Column footer="Total" :colspan="3" footer-style="text-align: right" />
        <Column :footer="formatPrice(total)" footer-class="font-semibold" />
      </Row>
    </ColumnGroup>
  </DataTable>
</template>
