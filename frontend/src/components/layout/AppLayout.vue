<script setup lang="ts">
import { ref } from 'vue';
import AppSidebar from './AppSidebar.vue';
import ThemeToggle from '@/components/ThemeToggle.vue';
import { SidebarLayout } from 'primevue';
import SidebarIcon from '@primeicons/vue/sidebar';

const open = ref(true);
</script>

<template>
  <SidebarLayout class="h-screen">
    <AppSidebar v-model:open="open" />

    <SidebarMain>
      <!-- Header -->
      <header
        class="flex h-12 items-center gap-2 border-b border-surface-200 px-4 dark:border-surface-700"
      >
        <!-- Botón colapsar -->
        <SidebarTrigger target="main-nav" severity="secondary" :text="true" size="small">
          <SidebarIcon />
        </SidebarTrigger>

        <span class="flex-1" />

        <!-- ThemeToggle arriba a la derecha -->
        <ThemeToggle />
      </header>

      <div class="flex-1 overflow-auto p-4">
        <router-view />
      </div>
    </SidebarMain>
  </SidebarLayout>

  <Toast position="bottom-right" />
</template>
