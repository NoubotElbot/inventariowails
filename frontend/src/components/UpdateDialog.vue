<script setup lang="ts">
  import { computed } from 'vue';
  import { toast } from 'vue-sonner';
  import { Button } from '@/components/ui/button';
  import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
  } from '@/components/ui/dialog';
  import { Progress } from '@/components/ui/progress';
  import { Spinner } from '@/components/ui/spinner';
  import { useUpdater } from '@/composables/useUpdater';
  import { formatDateTime } from '@/utils/date';

  const open = defineModel<boolean>('visible', { default: false });

  const { info, installing, percent, error, install } = useUpdater();

  const publishedOn = computed(() => formatDateTime(info.value?.published_at).split(' ')[0]);

  /** While installing the dialog can't be dismissed. */
  function preventWhileInstalling(event: Event): void {
    if (installing.value) event.preventDefault();
  }

  async function handleInstall() {
    const autoInstall = info.value?.auto_install;
    const ok = await install();

    if (!ok) {
      toast.error('No se pudo actualizar', {
        description: error.value ?? undefined,
        duration: 6000,
      });
      return;
    }

    // Se abrió la página de descargas en el navegador.
    if (!autoInstall) open.value = false;
  }
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent
      class="sm:max-w-lg"
      :show-close-button="!installing"
      @escape-key-down="preventWhileInstalling"
      @interact-outside="preventWhileInstalling"
    >
      <DialogHeader>
        <DialogTitle>Nueva versión disponible</DialogTitle>
        <DialogDescription v-if="info">
          Versión actual <strong>v{{ info.current_version }}</strong> → nueva
          <strong>v{{ info.latest_version }}</strong>
          <span v-if="publishedOn"> (publicada el {{ publishedOn }})</span>
        </DialogDescription>
      </DialogHeader>

      <div v-if="info" class="flex flex-col gap-4">
        <div v-if="info.notes" class="max-h-64 overflow-auto rounded-md border p-3">
          <pre class="font-sans text-xs whitespace-pre-wrap">{{ info.notes }}</pre>
        </div>

        <p v-if="!info.auto_install" class="text-xs text-muted-foreground">
          Esta versión no se puede instalar automáticamente. Se abrirá la página de descargas para
          que instales la nueva versión manualmente.
        </p>

        <div v-if="installing" class="flex flex-col gap-2">
          <Progress :model-value="percent" class="h-2" />
          <span class="text-xs text-muted-foreground">
            {{ percent === null ? 'Descargando…' : `Descargando… ${percent}%` }}
            La aplicación se cerrará para instalar la nueva versión.
          </span>
        </div>
      </div>

      <DialogFooter>
        <Button variant="outline" :disabled="installing" @click="open = false">Más tarde</Button>
        <Button :disabled="installing" @click="handleInstall">
          <Spinner v-if="installing" />
          {{ info?.auto_install ? 'Actualizar ahora' : 'Ir a descargas' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
