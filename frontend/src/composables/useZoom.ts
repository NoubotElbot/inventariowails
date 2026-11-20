import { computed, ref, watch } from 'vue';
import { useEventListener } from '@vueuse/core';
import { Window } from '@wailsio/runtime';

const STORAGE_KEY = 'app-zoom';

/** Niveles de zoom disponibles, al estilo de los navegadores. */
export const ZOOM_LEVELS = [0.5, 0.67, 0.75, 0.8, 0.9, 1, 1.1, 1.25, 1.5, 1.75, 2] as const;

const MIN_ZOOM = ZOOM_LEVELS[0];
const MAX_ZOOM = ZOOM_LEVELS[ZOOM_LEVELS.length - 1];

function readZoom(): number {
  try {
    const value = Number(localStorage.getItem(STORAGE_KEY));
    return Number.isFinite(value) && value >= MIN_ZOOM && value <= MAX_ZOOM ? value : 1;
  } catch {
    return 1;
  }
}

const zoom = ref<number>(readZoom());

function applyZoom(value: number): void {
  // Zoom nativo del webview: escala todo (incluidos diálogos y popovers).
  // Falla silenciosamente si no se ejecuta dentro de Wails.
  Window.SetZoom(value).catch(() => {});
}

// Aplicar inmediatamente el zoom guardado.
applyZoom(zoom.value);

watch(zoom, (value) => {
  applyZoom(value);
  try {
    localStorage.setItem(STORAGE_KEY, String(value));
  } catch {
    // Sin almacenamiento disponible: el zoom no se recuerda.
  }
});

const canZoomIn = computed(() => zoom.value < MAX_ZOOM);
const canZoomOut = computed(() => zoom.value > MIN_ZOOM);
const percent = computed(() => Math.round(zoom.value * 100));

function zoomIn(): void {
  const next = ZOOM_LEVELS.find((level) => level > zoom.value + 0.001);
  if (next !== undefined) zoom.value = next;
}

function zoomOut(): void {
  const prev = [...ZOOM_LEVELS].reverse().find((level) => level < zoom.value - 0.001);
  if (prev !== undefined) zoom.value = prev;
}

function resetZoom(): void {
  zoom.value = 1;
}

// Atajos globales: Ctrl/⌘ + '+', '-', '0' y Ctrl/⌘ + rueda.
useEventListener(window, 'keydown', (event: KeyboardEvent) => {
  if (!(event.ctrlKey || event.metaKey) || event.altKey) return;

  switch (event.key) {
    case '+':
    case '=':
      event.preventDefault();
      zoomIn();
      break;
    case '-':
    case '_':
      event.preventDefault();
      zoomOut();
      break;
    case '0':
      event.preventDefault();
      resetZoom();
      break;
  }
});

useEventListener(
  window,
  'wheel',
  (event: WheelEvent) => {
    if (!(event.ctrlKey || event.metaKey)) return;
    event.preventDefault();
    if (event.deltaY < 0) zoomIn();
    else if (event.deltaY > 0) zoomOut();
  },
  { passive: false }
);

export function useZoom() {
  return {
    zoom,
    percent,
    canZoomIn,
    canZoomOut,
    zoomIn,
    zoomOut,
    resetZoom,
  };
}
