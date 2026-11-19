<script setup lang="ts">
  import { useRouter } from 'vue-router';
  import { toast } from 'vue-sonner';
  import { ShoppingBasket } from '@lucide/vue';
  import ErrorAlert from '@/components/common/ErrorAlert.vue';
  import PageHeader from '@/components/PageHeader.vue';
  import ProductSearchSelect from '@/components/ProductSearchSelect.vue';
  import { useLeaveGuard } from '@/composables/useLeaveGuard';
  import { useSaleCart } from '@/composables/useSaleCart';
  import { useSales } from '@/composables/useSales';
  import type { Product } from '@/types/product';
  import SaleCartTable from './SaleCartTable.vue';
  import SaleCheckoutBar from './SaleCheckoutBar.vue';

  const router = useRouter();

  const { items, total, count, isEmpty, add, setQuantity, remove, clear, toRequest } =
    useSaleCart();
  const { create, loading, error, clearError } = useSales();

  useLeaveGuard(() => !isEmpty.value, {
    header: 'Venta sin confirmar',
    message: 'Hay productos en la venta. ¿Salir y descartarlos?',
  });

  function addProduct(product: Product): void {
    if (add(product)) return;
    toast.warning('Stock insuficiente', {
      description: `Solo hay ${product.stock} unidades de "${product.name}".`,
    });
  }

  async function confirmSale(): Promise<void> {
    const sale = await create(toRequest());
    if (!sale) return;

    // Empty the cart first so the leave guard lets us go to the detail.
    clear();
    toast.success(`Venta #${sale.id} registrada`);
    await router.push(`/sales/${sale.id}`);
  }
</script>

<template>
  <div class="flex max-w-5xl flex-col gap-6">
    <PageHeader title="Nueva venta" :icon="ShoppingBasket" />

    <ErrorAlert v-if="error" :message="error" closable @close="clearError" />

    <ProductSearchSelect @select="addProduct" />

    <SaleCartTable :items="items" @update-quantity="setQuantity" @remove="remove" />

    <SaleCheckoutBar
      :units="count"
      :total="total"
      :loading="loading"
      :disabled="isEmpty"
      @cancel="router.push('/sales')"
      @confirm="confirmSale"
    />
  </div>
</template>
