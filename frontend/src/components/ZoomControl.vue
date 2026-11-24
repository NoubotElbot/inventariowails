<script setup lang="ts">
  import { ZoomIn, ZoomOut } from '@lucide/vue';
  import IconButton from '@/components/common/IconButton.vue';
  import { Button } from '@/components/ui/button';
  import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
  import { useZoom } from '@/composables/useZoom';

  const { percent, canZoomIn, canZoomOut, zoomIn, zoomOut, resetZoom } = useZoom();
</script>

<template>
  <div class="flex items-center" role="group" aria-label="Zoom">
    <IconButton
      :icon="ZoomOut"
      tooltip="Alejar (Ctrl -)"
      :disabled="!canZoomOut"
      @click="zoomOut"
    />

    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          variant="ghost"
          size="sm"
          class="w-14 tabular-nums"
          aria-label="Restablecer zoom"
          @click="resetZoom"
        >
          {{ percent }}%
        </Button>
      </TooltipTrigger>
      <TooltipContent>Restablecer zoom (Ctrl 0)</TooltipContent>
    </Tooltip>

    <IconButton :icon="ZoomIn" tooltip="Acercar (Ctrl +)" :disabled="!canZoomIn" @click="zoomIn" />
  </div>
</template>
