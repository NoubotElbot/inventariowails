<script setup lang="ts">
  import { computed } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import ArrowLeft from '@primeicons/vue/arrow-left';
  import Receipt from '@primeicons/vue/receipt';
  import PageHeader from '@/components/PageHeader.vue';
  import { useSaleDetail } from '@/composables/useSaleDetail';
  import SaleItemsTable from './SaleItemsTable.vue';
  import SaleSummary from './SaleSummary.vue';

  const route = useRoute();
  const router = useRouter();

  const saleId = computed(() => Number(route.params.id));
  const { sale, loading, error } = useSaleDetail(saleId);
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageHeader :title="`Venta #${saleId}`" :icon="Receipt">
      <template #actions>
        <Button size="small" severity="secondary" @click="router.push('/sales')">
          <ArrowLeft :size="16" />
          Volver
        </Button>
      </template>
    </PageHeader>

    <Message v-if="error" severity="error">{{ error }}</Message>

    <div v-if="loading" class="flex flex-col gap-3">
      <Skeleton height="4rem" />
      <Skeleton height="12rem" />
    </div>

    <template v-else-if="sale">
      <SaleSummary :sale="sale" />
      <SaleItemsTable :items="sale.items" :total="sale.total" />
    </template>
  </div>
</template>
