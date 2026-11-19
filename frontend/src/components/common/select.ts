export interface SelectOption<V extends string | number = number> {
  value: V;
  label: string;
}

/** Maps `{ id, name }` records (e.g. category options) to select options. */
export function toSelectOptions(items: readonly { id: number; name: string }[]): SelectOption[] {
  return items.map((item) => ({ value: item.id, label: item.name }));
}
