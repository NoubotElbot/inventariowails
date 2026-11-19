<script setup lang="ts" generic="T">
  import type { HTMLAttributes, VNode } from 'vue';
  import { computed } from 'vue';
  import { cn } from '@/lib/utils';
  import { Skeleton } from '@/components/ui/skeleton';
  import {
    Table,
    TableBody,
    TableCell,
    TableEmpty,
    TableFooter,
    TableHead,
    TableHeader,
    TableRow,
  } from '@/components/ui/table';
  import type { DataTableColumn } from './data-table';

  /**
   * Presentational table over shadcn's Table: columns are declared as data and
   * each cell can be customized with a `#cell-<key>="{ row }"` slot.
   * Pagination is left to the caller (see DataTablePagination).
   */
  const {
    rows,
    columns,
    rowKey,
    loading = false,
    striped = false,
    clickable = false,
    skeletonRows = 5,
    tableClass,
  } = defineProps<{
    rows: readonly T[];
    columns: readonly DataTableColumn<T>[];
    rowKey: (row: T) => PropertyKey;
    loading?: boolean;
    striped?: boolean;
    /** Rows react to clicks (`row-click`) and show a pointer cursor. */
    clickable?: boolean;
    skeletonRows?: number;
    tableClass?: HTMLAttributes['class'];
  }>();

  const emit = defineEmits<{
    'row-click': [row: T];
  }>();

  defineSlots<
    {
      [K in `cell-${string}`]?: (props: { row: T; index: number }) => VNode[];
    } & {
      empty?: () => VNode[];
      footer?: () => VNode[];
    }
  >();

  // First load: no rows yet, show placeholders instead of the empty state.
  const showSkeleton = computed(() => loading && rows.length === 0);

  function cellText(column: DataTableColumn<T>, row: T): string {
    const value = column.value ? column.value(row) : (row as Record<string, unknown>)[column.key];
    return value === null || value === undefined ? '' : String(value);
  }
</script>

<template>
  <div class="overflow-hidden rounded-lg border">
    <Table :class="tableClass">
      <TableHeader class="bg-muted/50">
        <TableRow class="hover:bg-transparent">
          <TableHead v-for="column in columns" :key="column.key" :class="cn('px-3', column.class)">
            {{ column.header }}
          </TableHead>
        </TableRow>
      </TableHeader>

      <TableBody :class="{ 'opacity-60 transition-opacity': loading && !showSkeleton }">
        <template v-if="showSkeleton">
          <TableRow v-for="n in skeletonRows" :key="n" class="hover:bg-transparent">
            <TableCell v-for="column in columns" :key="column.key" class="px-3 py-3">
              <Skeleton class="h-4 w-full max-w-32" />
            </TableCell>
          </TableRow>
        </template>

        <TableEmpty v-else-if="rows.length === 0" :colspan="columns.length">
          <slot name="empty" />
        </TableEmpty>

        <template v-else>
          <TableRow
            v-for="(row, index) in rows"
            :key="rowKey(row)"
            :class="cn(striped && 'even:bg-muted/30', clickable && 'cursor-pointer')"
            @click="clickable && emit('row-click', row)"
          >
            <TableCell
              v-for="column in columns"
              :key="column.key"
              :class="cn('px-3 py-2.5', column.class)"
            >
              <slot :name="`cell-${column.key}`" :row="row" :index="index">
                {{ cellText(column, row) }}
              </slot>
            </TableCell>
          </TableRow>
        </template>
      </TableBody>

      <TableFooter v-if="$slots.footer">
        <slot name="footer" />
      </TableFooter>
    </Table>
  </div>
</template>
