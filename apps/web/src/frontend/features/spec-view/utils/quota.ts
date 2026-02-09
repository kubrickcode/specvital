export type QuotaLevel = "danger" | "normal" | "unlimited" | "warning";

const QUOTA_WARNING_THRESHOLD = 70;
const QUOTA_DANGER_THRESHOLD = 90;
const QUOTA_EXCEEDED_THRESHOLD = 100;

export const getQuotaLevel = (percentage: number | null): QuotaLevel => {
  if (percentage === null) return "unlimited";
  if (percentage >= QUOTA_DANGER_THRESHOLD) return "danger";
  if (percentage >= QUOTA_WARNING_THRESHOLD) return "warning";
  return "normal";
};

export const formatQuotaNumber = (num: number): string => {
  return num.toLocaleString();
};

export const isQuotaExceeded = (percentage: number | null): boolean => {
  return percentage !== null && percentage >= QUOTA_EXCEEDED_THRESHOLD;
};
