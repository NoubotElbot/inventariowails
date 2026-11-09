import { computed, readonly, ref } from 'vue';
import type { DataTablePageEvent } from 'primevue/datatable';

import { usePagination } from '@/composables/usePagination';
import { productService } from '@/services/productService';
import type {
  CategoryOptionResponse,
  CreateProductRequest,
  UpdateProductRequest,
} from 'bindings/inventario/internal/handler/models';
import type { Product } from '@/types/product';

export function useProducts() {
  const items = ref<Product[]>([]);
  const pending = ref(0);
  const error = ref<string | null>(null);
  const categoryOptions = ref<CategoryOptionResponse[]>([]);
  const selectedCategoryIds = ref<number[]>([]);
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

  async function loadCategoryOptions(): Promise<void> {
    try {
      const opts = await productService.categories();
      categoryOptions.value = opts ?? [];
    } catch {
      // non-critical — filter just won't show categories
    }
  }

  async function refresh(): Promise<void> {
    const req = {
      search: pagination.search.value,
      category_ids: selectedCategoryIds.value,
      page: pagination.page.value,
      page_size: pagination.pageSize.value,
    };

    const res = await productService.list(req);

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

  async function onCategoryFilter(ids: number[]): Promise<void> {
    selectedCategoryIds.value = ids;
    pagination.goTo(1);
    await fetchAll();
  }

  async function create(req: CreateProductRequest): Promise<boolean> {
    clearError();
    try {
      await withLoading(() => productService.create(req));
      return true;
    } catch (cause) {
      setError(cause);
      return false;
    }
  }

  async function update(id: number, req: UpdateProductRequest): Promise<boolean> {
    clearError();
    try {
      await withLoading(async () => {
        await productService.update(id, req);
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
        await productService.remove(id);
        await refresh();
      });
      return true;
    } catch (cause) {
      setError(cause);
      return false;
    }
  }

  return {
    items: readonly(items),
    loading,
    error: readonly(error),
    categoryOptions: categoryOptions,
    selectedCategoryIds: readonly(selectedCategoryIds),
    pagination: readonly(pagination),
    fetchAll,
    loadCategoryOptions,
    onPage,
    onSearch,
    onCategoryFilter,
    create,
    update,
    remove,
    clearError,
  };
}
