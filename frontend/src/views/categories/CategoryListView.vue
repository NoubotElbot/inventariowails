<script setup lang="ts">
  import { onMounted, ref } from 'vue';
  import { FolderOpen, Pencil, Plus, RotateCcw, Tag, Trash2 } from '@lucide/vue';
  import ConfirmActionButton from '@/components/common/ConfirmActionButton.vue';
  import DataTable from '@/components/common/DataTable.vue';
  import DataTablePagination from '@/components/common/DataTablePagination.vue';
  import EmptyState from '@/components/common/EmptyState.vue';
  import ErrorAlert from '@/components/common/ErrorAlert.vue';
  import IconButton from '@/components/common/IconButton.vue';
  import type { DataTableColumn } from '@/components/common/data-table';
  import PageHeader from '@/components/PageHeader.vue';
  import SearchInput from '@/components/SearchInput.vue';
  import { Button } from '@/components/ui/button';
  import { Checkbox } from '@/components/ui/checkbox';
  import { Label } from '@/components/ui/label';
  import { useCategories } from '@/composables/useCategories';
  import type { Category } from '@/types/category';
  import { formatDateTime } from '@/utils/date';
  import CategoryFormDialog from './CategoryFormDialog.vue';

  const {
    items,
    loading,
    error,
    pagination,
    fetchAll,
    onPage,
    onSearch,
    remove,
    restore,
    clearError,
    onlyInactive,
  } = useCategories();

  const columns: DataTableColumn<Category>[] = [
    { key: 'id', header: 'ID', class: 'w-16' },
    { key: 'name', header: 'Nombre' },
    { key: 'created_at', header: 'Creada', class: 'w-48' },
    { key: 'actions', header: 'Acciones', class: 'w-28' },
  ];

  // One dialog for the page: `editing` null means "create".
  const dialogOpen = ref(false);
  const editing = ref<Category | null>(null);

  function openDialog(category: Category | null): void {
    editing.value = category;
    dialogOpen.value = true;
  }

  onMounted(fetchAll);
</script>

<template>
  <div class="flex flex-col gap-6">
    <PageHeader title="Categorías" :icon="Tag">
      <template #actions>
        <Button size="sm" @click="openDialog(null)">
          <Plus />
          Nueva categoría
        </Button>
      </template>
    </PageHeader>

    <div class="flex flex-wrap items-center gap-4">
      <SearchInput placeholder="Buscar por nombre…" @search="onSearch" />
      <div class="flex items-center gap-2">
        <Checkbox id="show-deleted" v-model="onlyInactive" />
        <Label for="show-deleted" class="cursor-pointer">Mostrar eliminadas</Label>
      </div>
    </div>

    <ErrorAlert v-if="error" :message="error" closable @close="clearError" />

    <div class="flex flex-col gap-3">
      <DataTable
        :rows="items"
        :columns="columns"
        :row-key="(row) => row.id"
        :loading="loading"
        striped
        table-class="min-w-[30rem]"
      >
        <template #empty>
          <EmptyState :icon="FolderOpen" message="No hay categorías registradas" />
        </template>

        <template #cell-name="{ row }">
          <div class="flex items-center gap-2">
            <Tag class="size-3.5 text-muted-foreground" />
            <span class="font-medium">{{ row.name }}</span>
          </div>
        </template>

        <template #cell-created_at="{ row }">
          <span class="text-sm text-muted-foreground">{{ formatDateTime(row.created_at) }}</span>
        </template>

        <template #cell-actions="{ row }">
          <div class="flex items-center gap-1">
            <IconButton :icon="Pencil" tooltip="Editar" @click="openDialog(row)" />
            <ConfirmActionButton
              v-if="!row.deleted_at"
              :icon="Trash2"
              tooltip="Eliminar"
              :confirm="{
                title: 'Confirmar eliminación',
                description: `¿Eliminar la categoría «${row.name}»?`,
                confirmLabel: 'Eliminar',
              }"
              :action="() => remove(row.id)"
              success-message="Categoría eliminada"
              @done="fetchAll"
            />
            <ConfirmActionButton
              v-else
              :icon="RotateCcw"
              tooltip="Restaurar"
              variant="success"
              :confirm="{
                title: 'Confirmar restauración',
                description: `¿Restaurar la categoría «${row.name}»?`,
                confirmLabel: 'Restaurar',
              }"
              :action="() => restore(row.id)"
              success-message="Categoría restaurada"
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

    <CategoryFormDialog v-model:open="dialogOpen" :category="editing" @saved="fetchAll" />
  </div>
</template>
