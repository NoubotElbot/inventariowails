/** Suma las unidades de una lista de líneas de venta. */
export function totalUnits(items: readonly { quantity: number }[]): number {
  return items.reduce((sum, item) => sum + item.quantity, 0);
}
