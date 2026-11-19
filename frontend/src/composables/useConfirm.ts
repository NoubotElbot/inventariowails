import { readonly, shallowRef } from 'vue';

export interface ConfirmOptions {
  title: string;
  description: string;
  confirmLabel?: string;
  cancelLabel?: string;
  /** `destructive` paints the confirm button red. */
  variant?: 'default' | 'destructive' | 'success';
}

interface ConfirmRequest extends Required<ConfirmOptions> {
  resolve: (accepted: boolean) => void;
}

// One app-wide confirmation at a time, rendered by <ConfirmDialogHost />.
const current = shallowRef<ConfirmRequest | null>(null);

function settle(accepted: boolean): void {
  current.value?.resolve(accepted);
  current.value = null;
}

/**
 * Promise-based confirmation dialog. `confirm()` resolves to `true` when the
 * user accepts and `false` when they cancel or dismiss the dialog.
 */
export function useConfirm() {
  function confirm(options: ConfirmOptions): Promise<boolean> {
    // A new request supersedes (and rejects) a pending one.
    settle(false);

    return new Promise<boolean>((resolve) => {
      current.value = {
        confirmLabel: 'Aceptar',
        cancelLabel: 'Cancelar',
        variant: 'default',
        ...options,
        resolve,
      };
    });
  }

  return { confirm };
}

/** Internal API for the dialog host. */
export function useConfirmHost() {
  return { request: readonly(current), settle };
}
