import { Events } from '@wailsio/runtime';
import { UpdateProgress } from 'bindings/events';
import type { UpdateInfoResponse } from 'bindings/inventario/internal/modules/system/adapters/handler/models';
import { computed, readonly, ref } from 'vue';

import { systemService } from '@/services/systemService';

const CHECK_INTERVAL_MS = 6 * 60 * 60 * 1000;

// Estado compartido: el sidebar y el diálogo de actualización ven lo mismo.
const info = ref<UpdateInfoResponse | null>(null);
const installing = ref(false);
const downloaded = ref(0);
const total = ref(-1);
const error = ref<string | null>(null);
let started = false;

const available = computed(() => info.value?.available ?? false);

/** Porcentaje de descarga, o null si el tamaño es desconocido. */
const percent = computed(() =>
  total.value > 0 ? Math.min(100, Math.round((downloaded.value / total.value) * 100)) : null
);

async function check(): Promise<void> {
  try {
    info.value = await systemService.checkForUpdate();
  } catch (cause) {
    // Sin red o límite de la API: no molestar al usuario, se reintenta luego.
    console.warn('No se pudo comprobar si hay actualizaciones:', cause);
  }
}

/**
 * Descarga e inicia el instalador. Si la plataforma no admite instalación
 * automática, abre la página del release. Devuelve false si falló (ver `error`).
 */
async function install(): Promise<boolean> {
  if (!info.value?.available || installing.value) return false;

  error.value = null;

  if (!info.value.auto_install) {
    try {
      await systemService.openReleasePage();
      return true;
    } catch (cause) {
      error.value = cause instanceof Error ? cause.message : String(cause);
      return false;
    }
  }

  installing.value = true;
  downloaded.value = 0;
  total.value = -1;

  const off = Events.On(UpdateProgress, (event) => {
    downloaded.value = event.data.downloaded;
    total.value = event.data.total;
  });

  try {
    await systemService.installUpdate();
    // La app se cerrará en breve; `installing` queda en true a propósito.
    return true;
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : String(cause);
    installing.value = false;
    return false;
  } finally {
    off();
  }
}

export function useUpdater() {
  if (!started) {
    started = true;
    void check();
    setInterval(check, CHECK_INTERVAL_MS);
  }

  return {
    info: readonly(info),
    available,
    installing: readonly(installing),
    percent,
    error: readonly(error),
    check,
    install,
  };
}
