import type { HTMLAttributes } from 'vue';

export interface DataTableColumn<T> {
  /** Unique id; also names the cell slot (`#cell-<key>`). */
  key: string;
  header?: string;
  /** Rendered as plain text when there is no `#cell-<key>` slot. */
  value?: (row: T) => string | number | null | undefined;
  /** Classes for both the header and the body cells (width, alignment…). */
  class?: HTMLAttributes['class'];
}
