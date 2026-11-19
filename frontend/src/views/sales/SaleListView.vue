<script setup lang="ts">
  import { onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import { Eye, FolderOpen, Plus, ShoppingCart } from '@lucide/vue';
  import DataTable from '@/components/common/DataTable.vue';
  import DataTablePagination from '@/components/common/DataTablePagination.vue';
  import EmptyState from '@/components/common/EmptyState.vue';
  import ErrorAlert from '@/components/common/ErrorAlert.vue';
  import IconButton from '@/components/common/IconButton.vue';
  import type { DataTableColumn } from '@/components/common/data-table';
  import PageHeader from '@/components/PageHeader.vue';
  import { Button } from '@/components/ui/button';
  import { useSales } from '@/composables/useSales';
  import type { Sale } from '@/types/sale';
  import { formatPrice } from '@/utils/currency';
  import { formatDateTime } from '@/utils/date';
  import { totalUnits } from '@/utils/sale';

  const router = useRouter();
  const { items, loading, error, pagination, fetchAll, onPage, clearError } = useSales();

  const columns: DataTableColumn<Sale>[] = [
    { key: 'id', header: 'N°', class: 'w-20' },
    { key: 'created_at', header: 'Fecha' },
    { key: 'units', header: 'Unidades', class: 'w-32', value: (sale) => totalUnits(sale.items) },
    { key: 'total', header: 'Total', class: 'w-40' },
    { key: 'actions', class: 'w-16' },
  ];

  onMounted(fetchAll);

  function openSale(sale: Sale): void {
    router.push(`/sales/${sale.id}`);
  }
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageHeader title="Ventas" :icon="ShoppingCart">
      <template #actions>
        <Button size="sm" @click="router.push('/sales/new')">
          <Plus />
          Nueva venta
        </Button>
      </template>
    </PageHeader>

    <ErrorAlert v-if="error" :message="error" closable @close="clearError" />

    <div class="flex flex-col gap-3">
      <DataTable
        :rows="items"
        :columns="columns"
        :row-key="(row) => row.id"
        :loading="loading"
        striped
        clickable
        table-class="min-w-[40rem]"
        @row-click="openSale"
      >
        <template #empty>
          <EmptyState :icon="FolderOpen" message="No hay ventas registradas" />
        </template>

        <template #cell-id="{ row }">
          <span class="font-mono">#{{ row.id }}</span>
        </template>

        <template #cell-created_at="{ row }">
          <span class="text-sm">{{ formatDateTime(row.created_at) }}</span>
        </template>

        <template #cell-total="{ row }">
          <span class="font-medium">{{ formatPrice(row.total) }}</span>
        </template>

        <template #cell-actions="{ row }">
          <IconButton
            :icon="Eye"
            tooltip="Ver detalle"
            :aria-label="`Ver venta #${row.id}`"
            @click.stop="openSale(row)"
          />
        </template>
      </DataTable>

      <DataTablePagination
        :page="pagination.page"
        :page-size="pagination.pageSize"
        :total="pagination.total"
        @update:page="onPage"
      />
    </div>
  </div>
</template>
