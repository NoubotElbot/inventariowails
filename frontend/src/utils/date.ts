const pad = (n: number): string => String(n).padStart(2, '0');

/**
 * Formatea un timestamp RFC3339 (como lo envía el backend) a
 * "dd/mm/YYYY HH:MM:SS" en hora local. Devuelve '' si el valor es vacío o inválido.
 */
export function formatDateTime(value: string | null | undefined): string {
  if (!value) return '';
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return '';
  const date = `${pad(d.getDate())}/${pad(d.getMonth() + 1)}/${d.getFullYear()}`;
  const time = `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
  return `${date} ${time}`;
}
