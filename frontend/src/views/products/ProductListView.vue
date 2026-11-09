<script setup lang="ts">
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Box from '@primeicons/vue/box';
import FolderOpen from '@primeicons/vue/folder-open';
import PencilLine from '@primeicons/vue/pencil';
import Plus from '@primeicons/vue/plus';
import { useProducts } from '@/composables/useProducts';
import { formatDateTime } from '@/utils/date';
import SearchInput from '@/components/SearchInput.vue';
import ProductDeleteButton from './ProductDeleteButton.vue';

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
  clearError,
} = useProducts();

onMounted(async () => {
  await Promise.all([fetchAll(), loadCategoryOptions()]);
});

function formatPrice(value: number): string {
  return new Intl.NumberFormat('es-CL', { style: 'currency', currency: 'CLP' }).format(value);
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Encabezado -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <Box :size="20" class="text-primary" />
        <h1 class="text-xl font-semibold text-surface-800 dark:text-surface-100">Productos</h1>
      </div>
      <Button size="small" @click="router.push('/products/new')">
        <Plus />
        Nuevo producto
      </Button>
    </div>

    <!-- Búsqueda y filtro por categoría -->
    <div class="flex flex-row items-center justify-start gap-3">
      <SearchInput placeholder="Buscar por nombre, código o código de barras…" @search="onSearch" />
      <Select
        :model-value="selectedCategoryIds"
        :options="categoryOptions"
        multiple
        checkmark
        showClear
        option-label="name"
        option-value="id"
        placeholder="Todas las categorías"
        :max-selected-labels="2"
        class="w-64"
        @update:model-value="onCategoryFilter"
      />
    </div>

    <!-- Error banner -->
    <Message v-if="error" severity="error" :closable="true" @close="clearError">
      {{ error }}
    </Message>

    <!-- Tabla -->
    <DataTable
      :value="items"
      :loading="loading"
      lazy
      paginator
      :rows="pagination.pageSize"
      :first="pagination.first"
      :total-records="pagination.total"
      paginator-template="PrevPageLink PageLinks NextPageLink"
      striped-rows
      table-style="min-width: 50rem"
      @page="onPage"
    >
      <template #empty>
        <div class="flex flex-col items-center gap-3 py-12 text-surface-400">
          <FolderOpen :size="40" />
          <span class="text-sm">No hay productos registrados</span>
        </div>
      </template>

      <Column field="id" header="ID" style="width: 4rem" />

      <Column field="name" header="Nombre">
        <template #body="{ data }">
          <div class="flex items-center gap-2">
            <Box :size="14" class="text-surface-400" />
            <span class="font-medium">{{ data.name }}</span>
          </div>
        </template>
      </Column>

      <Column field="code" header="Código" style="width: 9rem">
        <template #body="{ data }">
          <span class="font-mono text-sm">{{ data.code }}</span>
        </template>
      </Column>

      <Column field="barcode" header="C. barras" style="width: 10rem">
        <template #body="{ data }">
          <span class="font-mono text-sm text-surface-500">{{ data.barcode || '—' }}</span>
        </template>
      </Column>

      <Column field="price" header="Precio" style="width: 9rem">
        <template #body="{ data }">
          <span>{{ formatPrice(data.price) }}</span>
        </template>
      </Column>

      <Column field="stock" header="Stock" style="width: 7rem">
        <template #body="{ data }">
          <span :class="data.low_stock ? 'text-red-500 font-semibold' : ''">
            {{ data.stock }}
            <span v-if="data.low_stock" class="text-xs">(bajo)</span>
          </span>
        </template>
      </Column>

      <Column field="active" header="Estado" style="width: 7rem">
        <template #body="{ data }">
          <Tag
            :value="data.active ? 'Activo' : 'Inactivo'"
            :severity="data.active ? 'success' : 'secondary'"
          />
        </template>
      </Column>

      <Column field="updated_at" header="Actualizado" style="width: 12rem">
        <template #body="{ data }">
          <span class="text-sm text-surface-500">{{ formatDateTime(data.updated_at) }}</span>
        </template>
      </Column>

      <Column header="Acciones" style="width: 7rem">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button
              severity="secondary"
              text
              rounded
              size="small"
              v-tooltip.top="'Editar'"
              @click="router.push(`/products/${data.id}/edit`)"
            >
              <PencilLine :size="15" />
            </Button>
            <ProductDeleteButton :product="data" @deleted="fetchAll" />
          </div>
        </template>
      </Column>
    </DataTable>
  </div>
</template>
