import { GetVersion } from 'bindings/inventario/internal/modules/system/adapters/handler/systemhandler';

export const systemService = {
  getVersion: () => GetVersion(),
};
