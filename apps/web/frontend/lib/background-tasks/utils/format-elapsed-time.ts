export const formatElapsedTime = (startedAt: string | null): string => {
  if (!startedAt) return "";

  const elapsed = Math.max(0, Math.floor((Date.now() - new Date(startedAt).getTime()) / 1000));

  if (elapsed < 60) {
    return `${elapsed}s`;
  }

  const minutes = Math.floor(elapsed / 60);
  const seconds = elapsed % 60;
  return `${minutes}m ${seconds}s`;
};
