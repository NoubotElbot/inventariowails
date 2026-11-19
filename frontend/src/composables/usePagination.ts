import { computed, ref } from 'vue';
import type { ComputedRef } from 'vue';

import type { PageResponse } from '@/types/pagination';

export function usePagination<TFilters extends object = Record<string, never>>(
  initialPageSize = 10,
  extraFilters?: ComputedRef<TFilters>
) {
  const page = ref(1); // 1-based
  const pageSize = ref(initialPageSize);
  const total = ref(0);
  const search = ref('');

  const params = computed(() => ({
    search: search.value,
    page: page.value,
    page_size: pageSize.value,
    ...(extraFilters?.value ?? ({} as TFilters)),
  }));

  /** Moves to another (1-based) page, e.g. from DataTablePagination. */
  function onPage(target: number): void {
    page.value = target;
  }

  /** Syncs state with the server's response. */
  function apply(res: PageResponse<unknown>): void {
    total.value = res.total;
    page.value = res.page;
    pageSize.value = res.page_size;
  }

  function goTo(target: number): void {
    page.value = Math.max(1, target);
  }

  function setSearch(value: string): void {
    search.value = value;
    page.value = 1;
  }

  return { page, pageSize, total, search, params, onPage, apply, goTo, setSearch };
}
