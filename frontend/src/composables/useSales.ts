import { readonly, shallowReadonly, shallowRef } from 'vue';

import { useLoading } from '@/composables/useLoading';
import { usePagination } from '@/composables/usePagination';
import { saleService } from '@/services/saleService';
import type { CreateSaleRequest } from 'bindings/inventario/internal/modules/sale/adapters/handler/models';
import type { Sale } from '@/types/sale';

export function useSales() {
  const items = shallowRef<Sale[]>([]);
  const { loading, error, withLoading, setError, clearError } = useLoading();
  const pagination = usePagination(10);

  async function refresh(): Promise<void> {
    const res = await saleService.list(pagination.params.value);
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

  async function onPage(page: number): Promise<void> {
    pagination.onPage(page);
    await fetchAll();
  }

  /** Registers a sale. Returns the created sale, or null if it failed (see `error`). */
  async function create(req: CreateSaleRequest): Promise<Sale | null> {
    clearError();
    try {
      return await withLoading(() => saleService.create(req));
    } catch (cause) {
      setError(cause);
      return null;
    }
  }

  return {
    items: shallowReadonly(items),
    loading,
    error,
    pagination: readonly(pagination),
    fetchAll,
    onPage,
    create,
    clearError,
  };
}
