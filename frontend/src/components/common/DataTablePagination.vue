<script setup lang="ts">
  import { ChevronLeft, ChevronRight } from '@lucide/vue';
  import {
    Pagination,
    PaginationContent,
    PaginationEllipsis,
    PaginationItem,
    PaginationNext,
    PaginationPrevious,
  } from '@/components/ui/pagination';

  /** Server-side pagination controls (1-based pages). Hidden when everything fits in one page. */
  const { page, pageSize, total } = defineProps<{
    page: number;
    pageSize: number;
    total: number;
  }>();

  const emit = defineEmits<{
    'update:page': [page: number];
  }>();
</script>

<template>
  <Pagination
    v-if="total > pageSize"
    :page="page"
    :items-per-page="pageSize"
    :total="total"
    :sibling-count="1"
    show-edges
    @update:page="emit('update:page', $event)"
  >
    <PaginationContent v-slot="{ items }">
      <PaginationPrevious>
        <ChevronLeft data-icon="inline-start" />
        <span class="hidden sm:block">Anterior</span>
      </PaginationPrevious>

      <template v-for="(item, index) in items" :key="index">
        <PaginationItem
          v-if="item.type === 'page'"
          :value="item.value"
          :is-active="item.value === page"
        >
          {{ item.value }}
        </PaginationItem>
        <PaginationEllipsis v-else :index="index" />
      </template>

      <PaginationNext>
        <span class="hidden sm:block">Siguiente</span>
        <ChevronRight data-icon="inline-end" />
      </PaginationNext>
    </PaginationContent>
  </Pagination>
</template>
