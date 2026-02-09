import type { BackgroundTask } from "../task-store";

export const getTaskPageUrl = (task: BackgroundTask): string | null => {
  const { analysisId, owner, repo } = task.metadata;

  if (task.type === "spec-generation" && analysisId && owner && repo) {
    return `/analyze/${owner}/${repo}?tab=spec`;
  }

  if (task.type === "analysis" && owner && repo) {
    return `/analyze/${owner}/${repo}`;
  }

  return null;
};
