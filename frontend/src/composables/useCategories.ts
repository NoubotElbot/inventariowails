import { computed, readonly, ref, watch } from 'vue';
import type { DataTablePageEvent } from 'primevue/datatable';

import { usePagination } from '@/composables/usePagination';
import { categoryService } from '@/services/categoryService';
import type { Category } from '@/types/category';

export function useCategories() {
  const items = ref<Category[]>([]);
  const pending = ref(0);
  const error = ref<string | null>(null);
  const onlyInactive = ref(false);
  const pagination = usePagination();

  const loading = computed(() => pending.value > 0);

  async function withLoading<T>(action: () => Promise<T>): Promise<T> {
    pending.value++;

    try {
      return await action();
    } finally {
      pending.value--;
    }
  }

  function setError(cause: unknown): void {
    error.value = cause instanceof Error ? cause.message : String(cause);
  }

  function clearError(): void {
    error.value = null;
  }

  async function refresh(): Promise<void> {
    const res = await categoryService.list({ ...pagination.params.value, only_inactive: onlyInactive.value });

    // La página pedida ya no existe (p. ej. se borró el último elemento): retrocede.
    if (res.items.length === 0 && res.total > 0 && pagination.page.value > 1) {
      pagination.goTo(res.total_pages);
      return refresh();
    }

    items.value = res.items;
    pagination.apply(res);
  }

  async function fetchAll(): Promise<void> {
    clearError();

    try {
      await withLoading(refresh);
    } catch (cause) {
      setError(cause);
    }
  }

  async function onPage(event: DataTablePageEvent): Promise<void> {
    pagination.onPage(event);
    await fetchAll();
  }

  async function onSearch(value: string): Promise<void> {
    pagination.setSearch(value);
    await fetchAll();
  }

  /** Crea la categoría. El llamador es responsable de recargar la lista. */
  async function create(name: string): Promise<boolean> {
    clearError();

    try {
      await withLoading(async () => {
        await categoryService.create(name);
      });

      return true;
    } catch (cause) {
      setError(cause);
      return false;
    }
  }

  async function update(id: number, name: string): Promise<boolean> {
    clearError();

    try {
      await withLoading(async () => {
        await categoryService.update(id, name);
        await refresh();
      });

      return true;
    } catch (cause) {
      setError(cause);
      return false;
    }
  }

  async function remove(id: number): Promise<boolean> {
    clearError();

    try {
      await withLoading(async () => {
        await categoryService.remove(id);
        await refresh();
      });

      return true;
    } catch (cause) {
      setError(cause);
      return false;
    }
  }

  async function restore(id: number): Promise<boolean> {
    clearError();

    try {
      await withLoading(async () => {
        await categoryService.restore(id);
        await refresh();
      });

      return true;
    } catch (cause) {
      setError(cause);
      return false;
    }
  }

  watch(onlyInactive, () => {
    pagination.goTo(1);
    fetchAll();
  });

  return {
    items: readonly(items),
    loading,
    error: readonly(error),
    pagination: readonly(pagination),
    onlyInactive,
    fetchAll,
    onPage,
    onSearch,
    create,
    update,
    remove,
    restore,
    clearError,
  };
}
