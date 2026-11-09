<script setup lang="ts">
import { onMounted } from 'vue';
import FolderOpen from '@primeicons/vue/folder-open';
import Tag from '@primeicons/vue/tag';
import { useToast } from 'primevue/usetoast';
import { useCategories } from '@/composables/useCategories';
import { formatDateTime } from '@/utils/date';
import CategoryCreateDialog from './CategoryCreateDialog.vue';
import CategoryUpdateDialog from './CategoryUpdateDialog.vue';
import CategoryDeleteButton from './CategoryDeleteButton.vue';
import SearchInput from '@/components/SearchInput.vue';
import CategoryRestoreButton from './CategoryRestoreButton.vue';

// ---------------------------------------------------------------------------
// Composables
// ---------------------------------------------------------------------------
const {
  items,
  loading,
  error,
  pagination,
  fetchAll,
  onPage,
  onSearch,
  clearError,
  onlyInactive,
  restore,
} = useCategories();
const toast = useToast();

// ---------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------
async function handleRestore(id: number, name: string): Promise<void> {
  const ok = await restore(id);
  if (ok) {
    toast.add({ severity: 'success', summary: `Categoría "${name}" restaurada`, life: 3000 });
  }
}

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------
onMounted(fetchAll);
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Encabezado ---------------------------------------------------------->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <Tag :size="20" class="text-primary" />
        <div>
          <h1 class="text-xl font-semibold text-surface-800 dark:text-surface-100">Categorías</h1>
        </div>
      </div>
      <CategoryCreateDialog @created="fetchAll" />
    </div>

    <!-- Búsqueda + filtro eliminadas ---------------------------------------->
    <div class="flex flex-wrap items-center gap-4">
      <SearchInput placeholder="Buscar por nombre…" @search="onSearch" />
      <div class="flex items-center gap-2">
        <Checkbox v-model="onlyInactive" :binary="true" inputId="show-deleted" />
        <label for="show-deleted" class="text-sm cursor-pointer select-none"
          >Mostrar eliminadas</label
        >
      </div>
    </div>

    <!-- Error banner -------------------------------------------------------->
    <Message v-if="error" severity="error" :closable="true" @close="clearError">
      {{ error }}
    </Message>

    <!-- Tabla (paginación en servidor) -------------------------------------->
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
      table-style="min-width: 30rem"
      @page="onPage"
    >
      <template #empty>
        <div class="flex flex-col items-center gap-3 py-12 text-surface-400">
          <FolderOpen :size="40" />
          <span class="text-sm">No hay categorías registradas</span>
        </div>
      </template>

      <Column field="id" header="ID" style="width: 4rem" />

      <Column field="name" header="Nombre">
        <template #body="{ data }">
          <div class="flex items-center gap-2">
            <Tag :size="14" class="text-surface-400" />
            <span class="font-medium">{{ data.name }}</span>
          </div>
        </template>
      </Column>

      <Column field="created_at" header="Creada" style="width: 12rem">
        <template #body="{ data }">
          <span class="text-sm text-surface-500">{{ formatDateTime(data.created_at) }}</span>
        </template>
      </Column>

      <Column header="Acciones" style="width: 7rem">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <CategoryUpdateDialog :category="data" @updated="fetchAll" />
            <CategoryDeleteButton v-if="!data.deleted_at" :category="data" @deleted="fetchAll" />
            <CategoryRestoreButton v-else :category="data" @restored="fetchAll" />
          </div>
        </template>
      </Column>
    </DataTable>
  </div>
</template>
