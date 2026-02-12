import { describe, expect, it } from "vitest";

import { getStatusDisplayInfo } from "./generation-progress-utils";

describe("getStatusDisplayInfo", () => {
  it("should return pending display info for pending status", () => {
    const result = getStatusDisplayInfo("pending");

    expect(result).toEqual({
      descriptionKey: "status.pendingDescription",
      titleKey: "status.pending",
    });
  });

  it("should return running display info for running status", () => {
    const result = getStatusDisplayInfo("running");

    expect(result).toEqual({
      descriptionKey: "status.runningDescription",
      titleKey: "status.running",
    });
  });

  it("should return completed display info for completed status", () => {
    const result = getStatusDisplayInfo("completed");

    expect(result).toEqual({
      descriptionKey: "completed.description",
      titleKey: "completed.title",
    });
  });

  it("should return failed display info for failed status", () => {
    const result = getStatusDisplayInfo("failed");

    expect(result).toEqual({
      descriptionKey: "failed.description",
      titleKey: "failed.title",
    });
  });

  it("should return pending display info for null status", () => {
    const result = getStatusDisplayInfo(null);

    expect(result).toEqual({
      descriptionKey: "status.pendingDescription",
      titleKey: "status.pending",
    });
  });

  it("should return pending display info for not_found status", () => {
    const result = getStatusDisplayInfo("not_found");

    expect(result).toEqual({
      descriptionKey: "status.pendingDescription",
      titleKey: "status.pending",
    });
  });
});
