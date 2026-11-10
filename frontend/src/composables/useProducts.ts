import { computed, readonly, ref } from 'vue';
import type { DataTablePageEvent } from 'primevue/datatable';

import { useLoading } from '@/composables/useLoading';
import { usePagination } from '@/composables/usePagination';
import { productService } from '@/services/productService';
import type { CategoryOptionResponse } from 'bindings/inventario/internal/modules/category/handler/models';
import type {
  CreateProductRequest,
  UpdateProductRequest,
} from 'bindings/inventario/internal/modules/product/handler/models';
import type { Product } from '@/types/product';

export function useProducts() {
  const items = ref<Product[]>([]);
  const { loading, error, withLoading, setError, clearError } = useLoading();
  const categoryOptions = ref<CategoryOptionResponse[]>([]);
  const selectedCategoryIds = ref<number[]>([]);
  const extraFilters = computed(() => ({ category_ids: selectedCategoryIds.value }));
  const pagination = usePagination(10, extraFilters);

  async function loadCategoryOptions(): Promise<void> {
    try {
      const opts = await productService.categories();
      categoryOptions.value = opts ?? [];
    } catch {
      // non-critical — filter just won't show categories
    }
  }

  async function refresh(): Promise<void> {
    const res = await productService.list(pagination.params.value);

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
      await withLoading(() => productService.update(id, req));
      return true;
    } catch (cause) {
      setError(cause);
      return false;
    }
  }

  async function remove(id: number): Promise<boolean> {
    clearError();
    try {
      await withLoading(() => productService.remove(id));
      return true;
    } catch (cause) {
      setError(cause);
      return false;
    }
  }

  return {
    items: readonly(items),
    loading,
    error,
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
