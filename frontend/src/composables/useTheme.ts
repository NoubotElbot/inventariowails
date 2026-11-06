import { ref, watch } from 'vue';

export type Theme = 'light' | 'dark';

const STORAGE_KEY = 'app-theme';

function readTheme(): Theme {
  return localStorage.getItem(STORAGE_KEY) === 'dark' ? 'dark' : 'light';
}

const theme = ref<Theme>(readTheme());

function applyTheme(value: Theme): void {
  document.documentElement.classList.toggle('app-dark', value === 'dark');

  document.documentElement.style.colorScheme = value;
}

// Aplicar inmediatamente el tema guardado.
applyTheme(theme.value);

watch(theme, (value) => {
  applyTheme(value);
  localStorage.setItem(STORAGE_KEY, value);
});

export function useTheme() {
  function setTheme(value: Theme): void {
    theme.value = value;
  }

  function toggleTheme(): void {
    setTheme(theme.value === 'dark' ? 'light' : 'dark');
  }

  return {
    theme,
    setTheme,
    toggleTheme,
  };
}
