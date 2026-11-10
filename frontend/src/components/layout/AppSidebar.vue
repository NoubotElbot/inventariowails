<script setup lang="ts">
  import { Box, InfoCircle, ShoppingCart, Tag } from '@primeicons/vue';
  import { useRoute, useRouter } from 'vue-router';

  const open = defineModel<boolean>('open', { default: true });

  const router = useRouter();
  const route = useRoute();

  const nav = [
    { label: 'Productos', icon: Box, path: '/products' },
    { label: 'Categorías', icon: Tag, path: '/categories' },
    { label: 'Ventas', icon: ShoppingCart, path: '/sales' },
  ];

  const isActive = (path: string) => route.path.startsWith(path);
</script>

<template>
  <Sidebar
    id="main-nav"
    v-model:open="open"
    collapsible="icon"
    as="nav"
    aria-label="Navegación principal"
  >
    <SidebarSpacer />
    <SidebarAside>
      <SidebarPanel>
        <!-- Logo -->
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton class="p-1!">
                <div
                  class="flex size-6 shrink-0 items-center justify-center rounded-md bg-primary text-primary-contrast text-xs font-bold leading-none"
                >
                  I
                </div>
                <span class="font-semibold text-sm">Inventario</span>
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
                  <SidebarMenuButton
                    :isActive="isActive(item.path)"
                    @click="router.push(item.path)"
                  >
                    <component :is="item.icon" />
                    <span>{{ item.label }}</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        <!-- Footer -->
        <SidebarFooter>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton>
                <InfoCircle class="size-6 shrink-0 text-xs" />
                <span class="text-xs">v0.1.0</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarFooter>

        <SidebarRail />
      </SidebarPanel>
    </SidebarAside>
  </Sidebar>
</template>
