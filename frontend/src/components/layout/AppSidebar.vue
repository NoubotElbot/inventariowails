<script setup lang="ts">
  import { Box, CloudDownload, Info, ShoppingCart, Tag } from '@lucide/vue';
  import { computed, onMounted, ref } from 'vue';
  import { RouterLink, useRoute } from 'vue-router';

  import UpdateDialog from '@/components/UpdateDialog.vue';
  import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarRail,
  } from '@/components/ui/sidebar';
  import { useUpdater } from '@/composables/useUpdater';
  import { systemService } from '@/services/systemService';

  const route = useRoute();

  const nav = [
    { label: 'Productos', icon: Box, path: '/products' },
    { label: 'Categorías', icon: Tag, path: '/categories' },
    { label: 'Ventas', icon: ShoppingCart, path: '/sales' },
  ];

  const isActive = (path: string) => route.path.startsWith(path);

  const version = ref('');
  const versionLabel = computed(() =>
    !version.value || version.value === 'dev' ? version.value : `v${version.value}`
  );

  const { info: updateInfo, available: updateAvailable, installing, percent } = useUpdater();
  const updateDialogOpen = ref(false);
  const updateLabel = computed(() => {
    if (installing.value) {
      return percent.value === null ? 'Descargando…' : `Descargando… ${percent.value}%`;
    }
    return `Actualizar a v${updateInfo.value?.latest_version}`;
  });

  onMounted(async () => {
    try {
      version.value = await systemService.getVersion();
    } catch {
      version.value = '';
    }
  });
</script>

<template>
  <Sidebar collapsible="icon" aria-label="Navegación principal">
    <!-- Logo -->
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" class="pointer-events-none">
            <div
              class="flex size-8 shrink-0 items-center justify-center rounded-md bg-sidebar-primary text-sm font-bold text-sidebar-primary-foreground"
            >
              I
            </div>
            <span class="font-semibold">Inventario</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>

    <!-- Nav -->
    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupLabel>Menú</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in nav" :key="item.path">
              <SidebarMenuButton as-child :is-active="isActive(item.path)" :tooltip="item.label">
                <RouterLink :to="item.path">
                  <component :is="item.icon" />
                  <span>{{ item.label }}</span>
                </RouterLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>

    <!-- Footer -->
    <SidebarFooter>
      <SidebarMenu>
        <SidebarMenuItem v-if="updateAvailable">
          <SidebarMenuButton
            class="text-primary"
            :tooltip="updateLabel"
            @click="updateDialogOpen = true"
          >
            <span class="relative shrink-0">
              <CloudDownload />
              <span
                v-if="!installing"
                class="absolute -top-0.5 -right-0.5 size-2 rounded-full bg-destructive"
                aria-hidden="true"
              />
            </span>
            <span class="text-xs font-medium">{{ updateLabel }}</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton :tooltip="versionLabel || 'Versión'" class="cursor-default">
            <Info />
            <span class="text-xs">{{ versionLabel }}</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooter>

    <SidebarRail />
  </Sidebar>

  <UpdateDialog v-model:visible="updateDialogOpen" />
</template>
