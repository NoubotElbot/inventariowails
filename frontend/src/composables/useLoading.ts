import { computed, readonly, ref } from 'vue';

export function useLoading() {
  const pending = ref(0);
  const error = ref<string | null>(null);
  const loading = computed(() => pending.value > 0);

  async function withLoading<T>(action: () => Promise<T>): Promise<T> {
    pending.value++;
    try {
      return await action();
    } finally {
      pending.value--;
    }
  }

  function setError(cause: unknown): void {
    error.value = cause instanceof Error ? cause.message : String(cause);
  }

  function clearError(): void {
    error.value = null;
  }

  return { loading, error: readonly(error), withLoading, setError, clearError };
}
