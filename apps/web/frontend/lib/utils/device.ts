const isMobileDevice = (): boolean => {
  if (typeof navigator === "undefined") return false;
  return /Android|iPhone|iPad|iPod/i.test(navigator.userAgent);
};

export const shouldUseNativeShare = (): boolean => {
  return (
    typeof navigator !== "undefined" && typeof navigator.share === "function" && isMobileDevice()
  );
};
