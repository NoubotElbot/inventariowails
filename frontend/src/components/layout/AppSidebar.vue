<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const nav = [
  { label: 'Productos', icon: 'pi pi-box', path: '/products' },
  { label: 'Categorías', icon: 'pi pi-tag', path: '/categories' },
  { label: 'Ventas', icon: 'pi pi-shopping-cart', path: '/sales' },
]

const isActive = (path: string) => route.path.startsWith(path)
</script>

<template>
  <aside class="w-56 flex flex-col bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700">
    <!-- Logo -->
    <div class="h-14 flex items-center px-5 border-b border-gray-200 dark:border-gray-700">
      <span class="font-semibold text-gray-800 dark:text-gray-100 tracking-tight">
        Inventario
      </span>
    </div>

    <!-- Nav -->
    <nav class="flex-1 py-4 space-y-1 px-3">
      <button v-for="item in nav" :key="item.path"
        class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors" :class="isActive(item.path)
          ? 'bg-primary-50 text-primary-700 font-medium dark:bg-primary-900/30 dark:text-primary-300'
          : 'text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'"
        @click="router.push(item.path)">
        <i :class="item.icon" />
        {{ item.label }}
      </button>
    </nav>

    <!-- Footer -->
    <div class="p-4 border-t border-gray-200 dark:border-gray-700">
      <span class="text-xs text-gray-400">v0.1.0</span>
    </div>
  </aside>
</template>