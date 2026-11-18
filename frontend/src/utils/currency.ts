const clp = new Intl.NumberFormat('es-CL', { style: 'currency', currency: 'CLP' });

/** Formatea un monto entero como pesos chilenos (ej. "$1.990"). */
export function formatPrice(value: number): string {
  return clp.format(value);
}
