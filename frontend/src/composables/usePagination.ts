import { computed, ref } from 'vue';
import type { DataTablePageEvent } from 'primevue/datatable';

import type { PageParams, PageResponse } from '@/types/pagination';

export function usePagination(initialPageSize = 10) {
  const page = ref(1); // 1-based
  const pageSize = ref(initialPageSize);
  const total = ref(0);

  const first = computed(() => (page.value - 1) * pageSize.value);
  const params = computed<PageParams>(() => ({ page: page.value, page_size: pageSize.value }));

  /** Updates state from a PrimeVue DataTable `@page` event (0-based page). */
  function onPage(event: DataTablePageEvent): void {
    page.value = event.page + 1;
    pageSize.value = event.rows;
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

  return { page, pageSize, total, first, params, onPage, apply, goTo };
}
