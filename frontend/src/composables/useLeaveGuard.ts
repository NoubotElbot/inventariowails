import { onBeforeRouteLeave } from 'vue-router';
import { useConfirm } from 'primevue/useconfirm';
import ExclamationTriangle from '@primeicons/vue/exclamation-triangle';

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

  const confirm = useConfirm();

  onBeforeRouteLeave(() => {
    if (!shouldConfirm()) return true;

    return new Promise<boolean>((resolve) => {
      confirm.require({
        header,
        message,
        icon: ExclamationTriangle,
        rejectProps: { label: rejectLabel, severity: 'secondary', outlined: true },
        acceptProps: { label: acceptLabel, severity: 'danger' },
        accept: () => resolve(true),
        reject: () => resolve(false),
        onHide: () => resolve(false),
      });
    });
  });
}
