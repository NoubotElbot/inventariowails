import type { ProductResponse } from 'bindings/inventario/internal/modules/product/adapters/handler/models';

export type Product = ProductResponse;

/** Values edited by ProductForm. `stock` only exists when creating a product. */
export interface ProductFormValues {
  name: string;
  code: string;
  barcode: string;
  category_id: number | null;
  price: number;
  min_stock: number;
  stock?: number;
}
