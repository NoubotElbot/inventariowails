import { onBeforeRouteLeave } from 'vue-router';
import { useConfirm } from '@/composables/useConfirm';

interface LeaveGuardOptions {
  header?: string;
  message?: string;
  acceptLabel?: string;
  rejectLabel?: string;
}

/**
 * Asks for confirmation before leaving the current route while
 * `shouldConfirm()` is true (e.g. a form with unsaved changes).
 */
export function useLeaveGuard(shouldConfirm: () => boolean, options: LeaveGuardOptions = {}) {
  const {
    header = 'Salir sin guardar',
    message = 'Hay cambios sin guardar. ¿Salir y descartarlos?',
    acceptLabel = 'Salir',
    rejectLabel = 'Quedarme',
  } = options;

  const { confirm } = useConfirm();

  onBeforeRouteLeave(() => {
    if (!shouldConfirm()) return true;

    // Cancelling or dismissing the dialog keeps the user on the page.
    return confirm({
      title: header,
      description: message,
      confirmLabel: acceptLabel,
      cancelLabel: rejectLabel,
      variant: 'destructive',
    });
  });
}
