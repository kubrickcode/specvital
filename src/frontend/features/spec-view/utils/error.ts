export type ErrorType = "generation" | "network" | "unknown";

export const detectErrorType = (error: Error | null): ErrorType => {
  if (!error) return "unknown";

  const message = error.message.toLowerCase();

  if (
    message.includes("network") ||
    message.includes("fetch") ||
    message.includes("connection") ||
    message.includes("timeout") ||
    message.includes("failed to fetch")
  ) {
    return "network";
  }

  if (message.includes("generation") || message.includes("generate")) {
    return "generation";
  }

  return "unknown";
};
