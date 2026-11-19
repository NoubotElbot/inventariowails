<script setup lang="ts">
  import { useToast } from 'primevue/usetoast';
  import { computed } from 'vue';

  import { useUpdater } from '@/composables/useUpdater';
  import { formatDateTime } from '@/utils/date';

  const open = defineModel<boolean>('visible', { default: false });

  const { info, installing, percent, error, install } = useUpdater();
  const toast = useToast();

  const publishedOn = computed(() => formatDateTime(info.value?.published_at).split(' ')[0]);

  async function handleInstall() {
    const autoInstall = info.value?.auto_install;
    const ok = await install();

    if (!ok) {
      toast.add({
        severity: 'error',
        summary: 'No se pudo actualizar',
        detail: error.value ?? undefined,
        life: 6000,
      });
      return;
    }

    // Se abrió la página de descargas en el navegador.
    if (!autoInstall) open.value = false;
  }
</script>

<template>
  <Dialog
    v-model:visible="open"
    header="Nueva versión disponible"
    modal
    :closable="!installing"
    :dismissableMask="!installing"
    :style="{ width: '32rem' }"
  >
    <div v-if="info" class="flex flex-col gap-4">
      <p class="text-sm">
        Versión actual <strong>v{{ info.current_version }}</strong> → nueva
        <strong>v{{ info.latest_version }}</strong>
        <span v-if="publishedOn" class="text-surface-500"> (publicada el {{ publishedOn }})</span>
      </p>

      <div
        v-if="info.notes"
        class="max-h-64 overflow-auto rounded-md border border-surface-200 p-3 dark:border-surface-700"
      >
        <pre class="font-sans text-xs whitespace-pre-wrap">{{ info.notes }}</pre>
      </div>

      <p v-if="!info.auto_install" class="text-xs text-surface-500">
        Esta instalación no admite actualización automática. Se abrirá la página de descargas para
        que instales la nueva versión manualmente.
      </p>

      <div v-if="installing" class="flex flex-col gap-2">
        <ProgressBar
          :value="percent ?? 0"
          :mode="percent === null ? 'indeterminate' : 'determinate'"
          :showValue="false"
          style="height: 0.5rem"
        />
        <span class="text-xs text-surface-500">
          {{ percent === null ? 'Descargando…' : `Descargando… ${percent}%` }}
          La aplicación se cerrará para instalar la nueva versión.
        </span>
      </div>
    </div>

    <template #footer>
      <Button
        label="Más tarde"
        severity="secondary"
        outlined
        :disabled="installing"
        @click="open = false"
      />
      <Button
        :label="info?.auto_install ? 'Actualizar ahora' : 'Ir a descargas'"
        :loading="installing"
        @click="handleInstall"
      />
    </template>
  </Dialog>
</template>
