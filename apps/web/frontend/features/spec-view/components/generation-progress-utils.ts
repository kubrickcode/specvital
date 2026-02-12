import type { SpecGenerationStatusEnum } from "../types";

export type StatusDisplayInfo = {
  descriptionKey: string;
  titleKey: string;
};

export const getStatusDisplayInfo = (
  status: SpecGenerationStatusEnum | null
): StatusDisplayInfo => {
  switch (status) {
    case "pending":
      return { descriptionKey: "status.pendingDescription", titleKey: "status.pending" };
    case "running":
      return { descriptionKey: "status.runningDescription", titleKey: "status.running" };
    case "completed":
      return { descriptionKey: "completed.description", titleKey: "completed.title" };
    case "failed":
      return { descriptionKey: "failed.description", titleKey: "failed.title" };
    default:
      return { descriptionKey: "status.pendingDescription", titleKey: "status.pending" };
  }
};
