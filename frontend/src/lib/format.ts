export function formatCurrency(amount: number | undefined | null): string {
  if (amount == null) return "";
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    maximumFractionDigits: 0,
  }).format(amount);
}

export function formatOdometer(km: number | undefined | null): string {
  if (km == null) return "";
  return `${km.toLocaleString()} km`;
}
