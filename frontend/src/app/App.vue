<script setup lang="ts">
  import AppLayout from '@/components/layout/AppLayout.vue';
  import ConfirmDialogHost from '@/components/common/ConfirmDialogHost.vue';
  import { Toaster } from '@/components/ui/sonner';
  import { TooltipProvider } from '@/components/ui/tooltip';
</script>

<template>
  <TooltipProvider :delay-duration="300">
    <AppLayout />
    <Toaster position="bottom-right" rich-colors close-button />
    <ConfirmDialogHost />
  </TooltipProvider>
</template>
