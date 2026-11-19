<script setup lang="ts">
  import { computed } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { ArrowLeft, Receipt } from '@lucide/vue';
  import ErrorAlert from '@/components/common/ErrorAlert.vue';
  import PageHeader from '@/components/PageHeader.vue';
  import { Button } from '@/components/ui/button';
  import { Skeleton } from '@/components/ui/skeleton';
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
        <Button size="sm" variant="secondary" @click="router.push('/sales')">
          <ArrowLeft />
          Volver
        </Button>
      </template>
    </PageHeader>

    <ErrorAlert v-if="error" :message="error" />

    <div v-if="loading" class="flex flex-col gap-3">
      <Skeleton class="h-16" />
      <Skeleton class="h-48" />
    </div>

    <template v-else-if="sale">
      <SaleSummary :sale="sale" />
      <SaleItemsTable :items="sale.items" :total="sale.total" />
    </template>
  </div>
</template>
