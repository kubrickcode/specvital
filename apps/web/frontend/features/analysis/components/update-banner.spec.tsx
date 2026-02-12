import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import { UpdateBanner } from "./update-banner";

// Mock hooks and API calls
const mockMutate = vi.fn();
let mockMutationState = { isPending: false };

vi.mock("@tanstack/react-query", () => ({
  useMutation: ({ onSuccess }: { onSuccess?: () => void }) => ({
    ...mockMutationState,
    mutate: () => {
      mockMutate();
      onSuccess?.();
    },
  }),
  useQuery: () => ({
    data: null,
    isLoading: false,
  }),
  useQueryClient: () => ({
    invalidateQueries: vi.fn(),
    removeQueries: vi.fn(),
  }),
}));

let mockUpdateStatus = { isLoading: false, parserOutdated: false, status: "new-commits" as string };

vi.mock("../hooks/use-update-status", () => ({
  updateStatusKeys: { detail: (owner: string, repo: string) => ["updateStatus", owner, repo] },
  useUpdateStatus: () => mockUpdateStatus,
}));

vi.mock("../api", () => ({
  fetchAnalysisStatus: vi.fn().mockResolvedValue({ status: "completed" }),
  triggerReanalyze: vi.fn().mockResolvedValue({}),
}));

vi.mock("../hooks/use-analysis", () => ({
  analysisKeys: { detail: (owner: string, repo: string) => ["analysis", owner, repo] },
}));

vi.mock("@/features/dashboard", () => ({
  paginatedRepositoriesKeys: { all: ["paginatedRepositories"] },
  repositoryStatsKeys: { all: ["repositoryStats"] },
}));

vi.mock("@/lib/background-tasks", () => ({
  userActiveTasksKeys: { all: ["userActiveTasks"] },
}));

vi.mock("sonner", () => ({
  toast: { error: vi.fn(), success: vi.fn() },
}));

const messages = {
  analyze: {
    updateBanner: {
      analyzing: "Analyzing...",
      dismissLabel: "Dismiss",
      newCommitsDetected: "New commits detected",
      parserUpdated: "Parser updated",
      reanalyzeFailed: "Reanalysis failed",
      reanalyzeQueued: "Reanalysis queued",
      updateNow: "Update Now",
    },
  },
};

const renderUpdateBanner = () => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <UpdateBanner owner="facebook" repo="react" />
    </NextIntlClientProvider>
  );
};

describe("UpdateBanner", () => {
  it("shows spinner on Update Now button when reanalysis is pending", () => {
    // First render: normal state with new commits
    mockUpdateStatus = { isLoading: false, parserOutdated: false, status: "new-commits" };
    mockMutationState = { isPending: true };

    renderUpdateBanner();

    const updateButton = screen.getByRole("button", { name: /Update Now/i });
    expect(updateButton).toBeDisabled();

    // The Loader2 spinner should be rendered (has animate-spin class)
    const spinner = updateButton.querySelector(".animate-spin");
    expect(spinner).toBeInTheDocument();
  });
});
