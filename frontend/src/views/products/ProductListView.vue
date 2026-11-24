<script setup lang="ts">
  import { computed, onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import { Box, FolderOpen, Pencil, Plus, Trash2 } from '@lucide/vue';
  import ConfirmActionButton from '@/components/common/ConfirmActionButton.vue';
  import DataTable from '@/components/common/DataTable.vue';
  import DataTablePagination from '@/components/common/DataTablePagination.vue';
  import EmptyState from '@/components/common/EmptyState.vue';
  import ErrorAlert from '@/components/common/ErrorAlert.vue';
  import IconButton from '@/components/common/IconButton.vue';
  import MultiSelect from '@/components/common/MultiSelect.vue';
  import type { DataTableColumn } from '@/components/common/data-table';
  import { toSelectOptions } from '@/components/common/select';
  import PageHeader from '@/components/PageHeader.vue';
  import SearchInput from '@/components/SearchInput.vue';
  import { Badge } from '@/components/ui/badge';
  import { Button } from '@/components/ui/button';
  import { useProducts } from '@/composables/useProducts';
  import type { Product } from '@/types/product';
  import { formatPrice } from '@/utils/currency';
  import { formatDateTime } from '@/utils/date';

  const router = useRouter();

  const {
    items,
    loading,
    error,
    categoryOptions,
    selectedCategoryIds,
    pagination,
    fetchAll,
    loadCategoryOptions,
    onPage,
    onSearch,
    onCategoryFilter,
    remove,
    clearError,
  } = useProducts();

  const categorySelectOptions = computed(() => toSelectOptions(categoryOptions.value));

  const columns: DataTableColumn<Product>[] = [
    { key: 'code', header: 'Código', class: 'w-36' },
    { key: 'name', header: 'Nombre' },
    { key: 'price', header: 'Precio', class: 'w-36' },
    { key: 'stock', header: 'Stock', class: 'w-28' },
    { key: 'active', header: 'Estado', class: 'w-28' },
    { key: 'updated_at', header: 'Actualizado', class: 'w-48' },
    { key: 'actions', header: 'Acciones', class: 'w-28' },
  ];

  onMounted(async () => {
    await Promise.all([fetchAll(), loadCategoryOptions()]);
  });
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageHeader title="Productos" :icon="Box">
      <template #actions>
        <Button size="sm" @click="router.push('/products/new')">
          <Plus />
          Nuevo producto
        </Button>
      </template>
    </PageHeader>

    <div class="flex flex-wrap items-center gap-3">
      <SearchInput placeholder="Buscar por nombre, código o código de barras…" @search="onSearch" />
      <MultiSelect
        class="w-64"
        :model-value="selectedCategoryIds"
        :options="categorySelectOptions"
        placeholder="Todas las categorías"
        @update:model-value="onCategoryFilter"
      />
    </div>

    <ErrorAlert v-if="error" :message="error" closable @close="clearError" />

    <div class="flex flex-col gap-3">
      <DataTable
        :rows="items"
        :columns="columns"
        :row-key="(row) => row.id"
        :loading="loading"
        striped
        table-class="min-w-[50rem]"
      >
        <template #empty>
          <EmptyState :icon="FolderOpen" message="No hay productos registrados" />
        </template>

        <template #cell-name="{ row }">
          <div class="flex items-center gap-2">
            <Box class="size-3.5 text-muted-foreground" />
            <span class="font-medium">{{ row.name }}</span>
          </div>
        </template>

        <template #cell-code="{ row }">
          <span class="font-mono text-sm">{{ row.code }}</span>
        </template>

        <template #cell-barcode="{ row }">
          <span class="font-mono text-sm text-muted-foreground">{{ row.barcode || '—' }}</span>
        </template>

        <template #cell-price="{ row }">{{ formatPrice(row.price) }}</template>

        <template #cell-stock="{ row }">
          <span :class="{ 'font-semibold text-destructive': row.low_stock }">
            {{ row.stock }}
            <span v-if="row.low_stock" class="text-xs">(bajo)</span>
          </span>
        </template>

        <template #cell-active="{ row }">
          <Badge :variant="row.active ? 'success' : 'secondary'">
            {{ row.active ? 'Activo' : 'Inactivo' }}
          </Badge>
        </template>

        <template #cell-updated_at="{ row }">
          <span class="text-sm text-muted-foreground">{{ formatDateTime(row.updated_at) }}</span>
        </template>

        <template #cell-actions="{ row }">
          <div class="flex items-center gap-1">
            <IconButton
              :icon="Pencil"
              tooltip="Editar"
              @click="router.push(`/products/${row.id}/edit`)"
            />
            <ConfirmActionButton
              :icon="Trash2"
              tooltip="Eliminar"
              :confirm="{
                title: 'Confirmar eliminación',
                description: `¿Eliminar el producto «${row.name}»?`,
                confirmLabel: 'Eliminar',
              }"
              :action="() => remove(row.id)"
              success-message="Producto eliminado"
              @done="fetchAll"
            />
          </div>
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
