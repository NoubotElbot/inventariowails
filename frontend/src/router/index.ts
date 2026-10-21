import { createRouter, createMemoryHistory } from 'vue-router'

export const router = createRouter({
  history: createMemoryHistory(),
  routes: [
    {
      path: '/',
      redirect: '/products',
    },
    // {
    //   path: '/categories',
    //   component: () => import('@/views/categories/CategoryListView.vue'),
    // },
    // {
    //   path: '/products',
    //   component: () => import('@/views/products/ProductListView.vue'),
    // },
    // {
    //   path: '/products/new',
    //   component: () => import('@/views/products/ProductFormView.vue'),
    // },
    // {
    //   path: '/products/:id/edit',
    //   component: () => import('@/views/products/ProductFormView.vue'),
    // },
    // {
    //   path: '/products/:id/stock',
    //   component: () => import('@/views/stock/StockMovementsView.vue'),
    // },
    // {
    //   path: '/sales',
    //   component: () => import('@/views/sales/SaleListView.vue'),
    // },
    // {
    //   path: '/sales/new',
    //   component: () => import('@/views/sales/SaleCreateView.vue'),
    // },
    // {
    //   path: '/sales/:id',
    //   component: () => import('@/views/sales/SaleDetailView.vue'),
    // },
  ],
})
