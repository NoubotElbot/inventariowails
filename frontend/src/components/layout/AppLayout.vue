<script setup lang="ts">
import AppSidebar from './AppSidebar.vue'
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-gray-50 dark:bg-gray-900">
    <AppSidebar />
    <main class="flex-1 overflow-y-auto">
      <div class="p-6">
        <RouterView />
      </div>
    </main>
  </div>
</template>
