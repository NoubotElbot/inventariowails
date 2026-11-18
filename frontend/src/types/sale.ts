import type {
  SaleItemResponse,
  SaleResponse,
} from 'bindings/inventario/internal/modules/sale/adapters/handler/models';
import type { Product } from '@/types/product';

export type Sale = SaleResponse;
export type SaleItem = SaleItemResponse;

/** A product added to a sale that hasn't been registered yet. */
export interface CartItem {
  product: Product;
  quantity: number;
}
