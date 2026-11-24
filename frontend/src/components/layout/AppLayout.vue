<script setup lang="ts">
  import AppSidebar from './AppSidebar.vue';
  import ThemeToggle from '@/components/ThemeToggle.vue';
  import ZoomControl from '@/components/ZoomControl.vue';
  import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar';
</script>

<template>
  <SidebarProvider class="h-full min-h-0">
    <AppSidebar />

    <SidebarInset class="min-h-0 min-w-0">
      <header class="flex h-12 shrink-0 items-center gap-2 border-b px-4">
        <SidebarTrigger class="-ml-1" />
        <span class="flex-1" />
        <ZoomControl />
        <ThemeToggle />
      </header>

      <main class="flex-1 overflow-auto p-4">
        <router-view />
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
