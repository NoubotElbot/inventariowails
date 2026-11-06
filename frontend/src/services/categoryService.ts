import {
  CreateCategory,
  DeleteCategory,
  ListCategories,
  UpdateCategory,
} from 'bindings/inventario/internal/handler/categoryhandler';

import type { PageParams } from '@/types/pagination';

export const categoryService = {
  list: (params: PageParams) => ListCategories(params),

  create: (name: string) => CreateCategory({ name }),

  update: (id: number, name: string) => UpdateCategory(id, { name }),

  remove: (id: number) => DeleteCategory(id),
};
