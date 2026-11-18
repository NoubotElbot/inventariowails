import { computed, readonly, shallowRef } from 'vue';

import { productService } from '@/services/productService';
import type { Product } from '@/types/product';

interface UseProductSearchOptions {
  pageSize?: number;
}

/**
 * Paginated product search (by name, code or barcode) meant for an
 * infinite-scroll list: `search` starts over from page 1 and `loadMore`
 * appends the next page. Responses to superseded requests are discarded.
 */
export function useProductSearch(options: UseProductSearchOptions = {}) {
  const { pageSize = 20 } = options;

  const suggestions = shallowRef<Product[]>([]);
  const term = shallowRef('');
  const page = shallowRef(0);
  const totalPages = shallowRef(0);
  const loading = shallowRef(false);
  /** True while the current input hasn't been searched yet (e.g. during the debounce). */
  const pending = shallowRef(false);
  const error = shallowRef<string | null>(null);

  const hasMore = computed(() => page.value < totalPages.value);

  // Every new request (or reset) bumps the id, so late responses are ignored.
  let latestRequest = 0;

  async function fetchPage(target: number): Promise<void> {
    const requestId = ++latestRequest;
    loading.value = true;
    error.value = null;
    try {
      const res = await productService.list({
        search: term.value,
        category_ids: [],
        page: target,
        page_size: pageSize,
      });
      if (requestId !== latestRequest) return;

      suggestions.value = target === 1 ? res.items : [...suggestions.value, ...res.items];
      page.value = res.page;
      totalPages.value = res.total_pages;
    } catch (cause) {
      if (requestId !== latestRequest) return;
      error.value = cause instanceof Error ? cause.message : String(cause);
    } finally {
      if (requestId === latestRequest) {
        loading.value = false;
        pending.value = false;
      }
    }
  }

  function search(value: string): Promise<void> {
    term.value = value.trim();
    return fetchPage(1);
  }

  function loadMore(): void {
    if (loading.value || !hasMore.value) return;
    void fetchPage(page.value + 1);
  }

  function clearResults(): void {
    latestRequest++;
    suggestions.value = [];
    page.value = 0;
    totalPages.value = 0;
    loading.value = false;
    error.value = null;
  }

  /** The input changed: drop the outdated results until the next search runs. */
  function markStale(): void {
    clearResults();
    pending.value = true;
  }

  /** Back to the initial, empty state. */
  function reset(): void {
    clearResults();
    term.value = '';
    pending.value = false;
  }

  return {
    // A computed (not readonly()) so it stays assignable to AutoComplete's `any[]` prop.
    suggestions: computed(() => suggestions.value),
    term: readonly(term),
    loading: readonly(loading),
    pending: readonly(pending),
    error: readonly(error),
    hasMore,
    search,
    loadMore,
    markStale,
    reset,
  };
}
