<script setup lang="ts">
  import { onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import type { DataTableRowClickEvent } from 'primevue/datatable';
  import Eye from '@primeicons/vue/eye';
  import FolderOpen from '@primeicons/vue/folder-open';
  import Plus from '@primeicons/vue/plus';
  import ShoppingCart from '@primeicons/vue/shopping-cart';
  import PageHeader from '@/components/PageHeader.vue';
  import { useSales } from '@/composables/useSales';
  import { formatPrice } from '@/utils/currency';
  import { formatDateTime } from '@/utils/date';
  import { totalUnits } from '@/utils/sale';

  const router = useRouter();
  const { items, loading, error, pagination, fetchAll, onPage, clearError } = useSales();

  onMounted(fetchAll);

  function openSale(id: number): void {
    router.push(`/sales/${id}`);
  }

  function onRowClick(event: DataTableRowClickEvent): void {
    openSale(event.data.id);
  }
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageHeader title="Ventas" :icon="ShoppingCart">
      <template #actions>
        <Button size="small" @click="router.push('/sales/new')">
          <Plus />
          Nueva venta
        </Button>
      </template>
    </PageHeader>

    <Message v-if="error" severity="error" :closable="true" @close="clearError">
      {{ error }}
    </Message>

    <DataTable
      :value="items"
      :loading="loading"
      lazy
      paginator
      :rows="pagination.pageSize"
      :first="pagination.first"
      :total-records="pagination.total"
      paginator-template="PrevPageLink PageLinks NextPageLink"
      data-key="id"
      striped-rows
      row-hover
      class="cursor-pointer"
      table-style="min-width: 40rem"
      @page="onPage"
      @row-click="onRowClick"
    >
      <template #empty>
        <div class="flex flex-col items-center gap-3 py-12 text-surface-400">
          <FolderOpen :size="40" />
          <span class="text-sm">No hay ventas registradas</span>
        </div>
      </template>

      <Column field="id" header="N°" style="width: 5rem">
        <template #body="{ data }">
          <span class="font-mono">#{{ data.id }}</span>
        </template>
      </Column>

      <Column header="Fecha">
        <template #body="{ data }">
          <span class="text-sm">{{ formatDateTime(data.created_at) }}</span>
        </template>
      </Column>

      <Column header="Unidades" style="width: 8rem">
        <template #body="{ data }">{{ totalUnits(data.items) }}</template>
      </Column>

      <Column header="Total" style="width: 10rem">
        <template #body="{ data }">
          <span class="font-medium">{{ formatPrice(data.total) }}</span>
        </template>
      </Column>

      <Column style="width: 4rem">
        <template #body="{ data }">
          <Button
            severity="secondary"
            text
            rounded
            size="small"
            v-tooltip.top="'Ver detalle'"
            :aria-label="`Ver venta #${data.id}`"
            @click.stop="openSale(data.id)"
          >
            <Eye :size="15" />
          </Button>
        </template>
      </Column>
    </DataTable>
  </div>
</template>
