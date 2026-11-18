import { shallowRef, toValue, watch } from 'vue';
import type { MaybeRefOrGetter } from 'vue';

import { useLoading } from '@/composables/useLoading';
import { saleService } from '@/services/saleService';
import type { Sale } from '@/types/sale';

/** Loads a sale by id, and reloads it whenever the id changes. */
export function useSaleDetail(id: MaybeRefOrGetter<number>) {
  const sale = shallowRef<Sale | null>(null);
  const { loading, error, withLoading, setError, clearError } = useLoading();

  async function load(saleId: number): Promise<void> {
    clearError();
    try {
      sale.value = await withLoading(() => saleService.get(saleId));
    } catch (cause) {
      sale.value = null;
      setError(cause);
    }
  }

  watch(() => toValue(id), load, { immediate: true });

  return { sale, loading, error };
}
