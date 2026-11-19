import {
  CheckForUpdate,
  GetVersion,
  InstallUpdate,
  OpenReleasePage,
} from 'bindings/inventario/internal/modules/system/adapters/handler/systemhandler';

export const systemService = {
  getVersion: () => GetVersion(),

  checkForUpdate: () => CheckForUpdate(),

  /** Descarga e inicia el instalador; si tiene éxito, la app se cierra sola. */
  installUpdate: () => InstallUpdate(),

  openReleasePage: () => OpenReleasePage(),
};
