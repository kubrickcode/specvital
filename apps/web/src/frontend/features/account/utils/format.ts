export const formatNumber = (value: number | null | undefined, unlimitedSymbol = "∞"): string => {
  if (value === null || value === undefined) return unlimitedSymbol;
  return value.toLocaleString();
};

type ResetInfo = {
  date: string;
  days: number;
};

export const getResetInfo = (resetAt: string, locale: string): ResetInfo => {
  const resetDate = new Date(resetAt);
  const now = new Date();
  const diffMs = resetDate.getTime() - now.getTime();
  const diffDays = Math.ceil(diffMs / (1000 * 60 * 60 * 24));

  const dateStr = resetDate.toLocaleDateString(locale, {
    day: "numeric",
    month: "short",
  });

  return { date: dateStr, days: diffDays };
};
