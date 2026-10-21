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
  <aside
    class="flex w-56 flex-col border-r border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800"
  >
    <!-- Logo -->
    <div class="flex h-14 items-center border-b border-gray-200 px-5 dark:border-gray-700">
      <span class="font-semibold tracking-tight text-gray-800 dark:text-gray-100">
        Inventario
      </span>
    </div>

    <!-- Nav -->
    <nav class="flex-1 space-y-1 px-3 py-4">
      <button
        v-for="item in nav"
        :key="item.path"
        class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors"
        :class="
          isActive(item.path)
            ? 'bg-primary-50 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300 font-medium'
            : 'text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700'
        "
        @click="router.push(item.path)"
      >
        <i :class="item.icon" />
        {{ item.label }}
      </button>
    </nav>

    <!-- Footer -->
    <div class="border-t border-gray-200 p-4 dark:border-gray-700">
      <span class="text-xs text-gray-400">v0.1.0</span>
    </div>
  </aside>
</template>
