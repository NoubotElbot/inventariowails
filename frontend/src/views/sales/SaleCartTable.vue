<script setup lang="ts">
  import CartPlus from '@primeicons/vue/cart-plus';
  import Minus from '@primeicons/vue/minus';
  import Plus from '@primeicons/vue/plus';
  import Trash from '@primeicons/vue/trash';
  import { formatPrice } from '@/utils/currency';
  import type { CartItem } from '@/types/sale';
  import SaleProductCell from './SaleProductCell.vue';

  defineProps<{
    items: readonly CartItem[];
  }>();

  const emit = defineEmits<{
    updateQuantity: [productId: number, quantity: number];
    remove: [productId: number];
  }>();

  function onQuantityChange(productId: number, quantity: number | null): void {
    if (quantity !== null) emit('updateQuantity', productId, quantity);
  }
</script>

<template>
  <DataTable :value="items" data-key="product.id" table-style="min-width: 44rem">
    <template #empty>
      <div class="flex flex-col items-center gap-3 py-10 text-surface-400">
        <CartPlus :size="36" />
        <span class="text-sm">Busca un producto para agregarlo a la venta</span>
      </div>
    </template>

    <Column header="Producto">
      <template #body="{ data }">
        <SaleProductCell :name="data.product.name" :code="data.product.code" />
      </template>
    </Column>

    <Column header="Precio unit." style="width: 9rem">
      <template #body="{ data }">{{ formatPrice(data.product.price) }}</template>
    </Column>

    <Column header="Cantidad" style="width: 11rem">
      <template #body="{ data }">
        <InputNumber
          :model-value="data.quantity"
          show-buttons
          button-layout="horizontal"
          :min="1"
          :max="data.product.stock"
          :allow-empty="false"
          size="small"
          input-class="w-14 text-center"
          :aria-label="`Cantidad de ${data.product.name}`"
          @update:model-value="onQuantityChange(data.product.id, $event)"
        >
          <template #incrementicon><Plus :size="12" /></template>
          <template #decrementicon><Minus :size="12" /></template>
        </InputNumber>
      </template>
    </Column>

    <Column header="Subtotal" style="width: 9rem">
      <template #body="{ data }">
        <span class="font-medium">{{ formatPrice(data.product.price * data.quantity) }}</span>
      </template>
    </Column>

    <Column style="width: 4rem">
      <template #body="{ data }">
        <Button
          severity="danger"
          text
          rounded
          size="small"
          v-tooltip.top="'Quitar'"
          :aria-label="`Quitar ${data.product.name}`"
          @click="emit('remove', data.product.id)"
        >
          <Trash :size="15" />
        </Button>
      </template>
    </Column>
  </DataTable>
</template>
