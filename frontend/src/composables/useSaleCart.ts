import { computed, readonly, ref } from 'vue';

import { totalUnits } from '@/utils/sale';
import type { CreateSaleRequest } from 'bindings/inventario/internal/modules/sale/adapters/handler/models';
import type { Product } from '@/types/product';
import type { CartItem } from '@/types/sale';

const clamp = (value: number, min: number, max: number): number =>
  Math.min(Math.max(value, min), max);

/** Products picked for a sale that hasn't been registered yet. */
export function useSaleCart() {
  const items = ref<CartItem[]>([]);

  const total = computed(() =>
    items.value.reduce((sum, item) => sum + item.product.price * item.quantity, 0)
  );
  const count = computed(() => totalUnits(items.value));
  const isEmpty = computed(() => items.value.length === 0);

  function find(productId: number): CartItem | undefined {
    return items.value.find((item) => item.product.id === productId);
  }

  /**
   * Adds one unit of the product (a new line, or +1 on its existing line).
   * Returns false when the product's stock doesn't allow another unit.
   */
  function add(product: Product): boolean {
    const existing = find(product.id);
    const quantity = (existing?.quantity ?? 0) + 1;
    if (quantity > product.stock) return false;

    if (existing) existing.quantity = quantity;
    else items.value.push({ product, quantity });
    return true;
  }

  /** Sets a line's quantity, kept between 1 and the product's stock. */
  function setQuantity(productId: number, quantity: number): void {
    const item = find(productId);
    if (item) item.quantity = clamp(quantity, 1, item.product.stock);
  }

  function remove(productId: number): void {
    items.value = items.value.filter((item) => item.product.id !== productId);
  }

  function clear(): void {
    items.value = [];
  }

  function toRequest(): CreateSaleRequest {
    return {
      items: items.value.map((item) => ({
        product_id: item.product.id,
        quantity: item.quantity,
      })),
    };
  }

  return {
    items: readonly(items),
    total,
    count,
    isEmpty,
    add,
    setQuantity,
    remove,
    clear,
    toRequest,
  };
}
