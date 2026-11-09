import {
  CreateProduct,
  DeleteProduct,
  GetProduct,
  ListProducts,
  UpdateProduct,
} from 'bindings/inventario/internal/handler/producthandler';
import { ListCategoryOptions } from 'bindings/inventario/internal/handler/categoryhandler';
import type {
  CreateProductRequest,
  ListProductsRequest,
  UpdateProductRequest,
} from 'bindings/inventario/internal/handler/models';

export const productService = {
  list: (req: ListProductsRequest) => ListProducts(req),
  get: (id: number) => GetProduct(id),
  create: (req: CreateProductRequest) => CreateProduct(req),
  update: (id: number, req: UpdateProductRequest) => UpdateProduct(id, req),
  remove: (id: number) => DeleteProduct(id),
  categories: () => ListCategoryOptions(),
};
