import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("next/navigation", () => ({
  useSearchParams: vi.fn(() => ({
    get: vi.fn(() => null),
  })),
}));

vi.mock("@/i18n/navigation", () => ({
  Link: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}));

import type { ActiveTask } from "@/lib/api/types";
import * as bgTasks from "@/lib/background-tasks";

import { ActiveTasksSection } from "./active-tasks-section";

const messages = {
  backgroundTasks: {
    activeTasks: {
      countBadge: "{count}",
      title: "Active Tasks",
      viewTask: "View",
    },
    status: {
      processing: "Processing {time}",
      queued: "Queued",
    },
    taskType: {
      analysis: "Analysis",
      "spec-generation": "Spec Generation",
    },
  },
};

const createQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

const renderWithProviders = (ui: React.ReactElement) => {
  const queryClient = createQueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <NextIntlClientProvider locale="en" messages={messages}>
        {ui}
      </NextIntlClientProvider>
    </QueryClientProvider>
  );
};

const createTask = (overrides: Partial<ActiveTask> = {}): ActiveTask => ({
  createdAt: new Date().toISOString(),
  id: "task-1",
  owner: "test-owner",
  repo: "test-repo",
  startedAt: new Date(Date.now() - 45000).toISOString(),
  status: "analyzing",
  type: "analysis",
  ...overrides,
});

describe("ActiveTasksSection", () => {
  beforeEach(() => {
    vi.spyOn(bgTasks, "useUserActiveTasks");
    vi.spyOn(bgTasks, "useActiveTasks").mockReturnValue([]);
    vi.useFakeTimers({ shouldAdvanceTime: true });
  });

  afterEach(() => {
    cleanup();
    vi.restoreAllMocks();
    vi.useRealTimers();
  });

  it("should render nothing when no active tasks", () => {
    vi.mocked(bgTasks.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 0,
      data: { tasks: [] },
      isLoading: false,
    });

    const { container } = renderWithProviders(<ActiveTasksSection />);

    expect(container.firstChild).toBeNull();
  });

  it("should render section title with count badge when tasks exist", () => {
    vi.mocked(bgTasks.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 2,
      data: {
        tasks: [createTask(), createTask({ id: "task-2", status: "queued" })],
      },
      isLoading: false,
    });

    renderWithProviders(<ActiveTasksSection />);

    expect(screen.getByText("Active Tasks")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();
  });

  it("should expand/collapse on click", () => {
    vi.mocked(bgTasks.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [createTask()] },
      isLoading: false,
    });

    renderWithProviders(<ActiveTasksSection />);

    // Initially collapsed
    expect(screen.queryByText("test-owner/test-repo")).not.toBeInTheDocument();

    // Click to expand
    fireEvent.click(screen.getByText("Active Tasks"));

    // Now visible
    expect(screen.getByText("test-owner/test-repo")).toBeInTheDocument();
  });

  it("should render task with owner/repo as display name", () => {
    vi.mocked(bgTasks.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [createTask()] },
      isLoading: false,
    });

    renderWithProviders(<ActiveTasksSection />);
    fireEvent.click(screen.getByText("Active Tasks"));

    expect(screen.getByText("test-owner/test-repo")).toBeInTheDocument();
    expect(screen.getByText("Analysis")).toBeInTheDocument();
  });

  it("should render View link for tasks", () => {
    vi.mocked(bgTasks.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [createTask()] },
      isLoading: false,
    });

    renderWithProviders(<ActiveTasksSection />);
    fireEvent.click(screen.getByText("Active Tasks"));

    expect(screen.getByText("View")).toBeInTheDocument();
    expect(screen.getByRole("link")).toHaveAttribute("href", "/analyze/test-owner/test-repo");
  });

  it("should render multiple tasks", () => {
    vi.mocked(bgTasks.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 2,
      data: {
        tasks: [
          createTask({ id: "task-1", status: "analyzing" }),
          createTask({
            id: "task-2",
            owner: "other",
            repo: "repo2",
            startedAt: undefined,
            status: "queued",
          }),
        ],
      },
      isLoading: false,
    });

    renderWithProviders(<ActiveTasksSection />);
    fireEvent.click(screen.getByText("Active Tasks"));

    expect(screen.getByText("test-owner/test-repo")).toBeInTheDocument();
    expect(screen.getByText("other/repo2")).toBeInTheDocument();
  });
});
