import {
  CreateSale,
  GetSale,
  ListSales,
} from 'bindings/inventario/internal/modules/sale/adapters/handler/salehandler';
import type {
  CreateSaleRequest,
  ListSalesRequest,
} from 'bindings/inventario/internal/modules/sale/adapters/handler/models';

export const saleService = {
  list: (req: ListSalesRequest) => ListSales(req),
  get: (id: number) => GetSale(id),
  create: (req: CreateSaleRequest) => CreateSale(req),
};
