export type AnalysisStatus = "analyzing" | "loading" | "queued";

export type WaitingStatus = Extract<AnalysisStatus, "analyzing" | "queued">;
