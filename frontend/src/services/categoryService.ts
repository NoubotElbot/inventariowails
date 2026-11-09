import {
  CreateCategory,
  DeleteCategory,
  ListCategories,
  RestoreCategory,
  UpdateCategory,
} from 'bindings/inventario/internal/modules/category/handler/categoryhandler';

import type { ListCategoriesRequest } from '@/types/pagination';

export const categoryService = {
  list: (req: ListCategoriesRequest) => ListCategories(req),

  create: (name: string) => CreateCategory({ name }),

  update: (id: number, name: string) => UpdateCategory(id, { name }),

  remove: (id: number) => DeleteCategory(id),

  restore: (id: number) => RestoreCategory(id),
};
