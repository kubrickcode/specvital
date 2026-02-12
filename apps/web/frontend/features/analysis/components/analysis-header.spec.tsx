import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import { AnalysisHeader } from "./analysis-header";

vi.mock("@/lib/hooks", () => ({
  useTruncateDetection: () => ({ isTruncated: false, ref: { current: null } }),
}));

vi.mock("@/lib/motion", () => ({
  fadeInUp: {},
}));

vi.mock("./commit-selector", () => ({
  CommitSelector: ({ currentCommitSha }: { currentCommitSha: string }) => (
    <span data-testid="commit-selector">{currentCommitSha.slice(0, 7)}</span>
  ),
}));

vi.mock("./share-button", () => ({
  ShareButton: () => <button type="button">Share</button>,
}));

const messages = {
  analyze: {
    analyzedAt: "Analyzed: {date}",
    branch: "Branch",
    commit: "Commit",
    committedAt: "Committed: {date}",
    details: "Details",
    parserVersion: "Parser: {version}",
    viewOnGitHub: "View on GitHub",
  },
};

const defaultProps = {
  analyzedAt: "2024-06-01T12:00:00Z",
  branchName: "main",
  commitSha: "abc123def456",
  committedAt: "2024-05-31T10:00:00Z",
  onCommitSelect: vi.fn(),
  owner: "facebook",
  parserVersion: "v1.5.1",
  repo: "react",
};

const renderAnalysisHeader = (props = defaultProps) => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <AnalysisHeader {...props} />
    </NextIntlClientProvider>
  );
};

describe("AnalysisHeader", () => {
  it("renders the repository full name as heading", () => {
    renderAnalysisHeader();

    expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("facebook/react");
  });

  it("renders the GitHub link with correct href", () => {
    renderAnalysisHeader();

    const link = screen.getByRole("link");
    expect(link).toHaveAttribute("href", "https://github.com/facebook/react");
    expect(link).toHaveAttribute("target", "_blank");
    expect(link).toHaveAttribute("rel", "noopener noreferrer");
  });

  it("renders the details toggle button", () => {
    renderAnalysisHeader();

    expect(screen.getByRole("button", { name: "Details" })).toBeInTheDocument();
  });

  it("shows branch name, commit, and parser version in the collapsible details panel", () => {
    renderAnalysisHeader();

    fireEvent.click(screen.getByRole("button", { name: "Details" }));

    expect(screen.getByText("main")).toBeInTheDocument();
    expect(screen.getByTestId("commit-selector")).toBeInTheDocument();
    expect(screen.getByText(/Parser: v1.5.1/)).toBeInTheDocument();
  });

  it("omits branch name when not provided", () => {
    renderAnalysisHeader({ ...defaultProps, branchName: undefined });

    fireEvent.click(screen.getByRole("button", { name: "Details" }));

    expect(screen.queryByText("Branch")).not.toBeInTheDocument();
  });

  it("omits parser version when not provided", () => {
    renderAnalysisHeader({ ...defaultProps, parserVersion: undefined });

    fireEvent.click(screen.getByRole("button", { name: "Details" }));

    expect(screen.queryByText(/Parser:/)).not.toBeInTheDocument();
  });
});
