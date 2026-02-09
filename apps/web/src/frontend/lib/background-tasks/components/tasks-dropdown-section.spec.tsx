import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@radix-ui/react-dropdown-menu";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { cleanup, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import type { ActiveTask } from "@/lib/api/types";

import * as hooksModule from "../hooks";
import * as useUserActiveTasksModule from "../use-user-active-tasks";
import { TasksDropdownSection } from "./tasks-dropdown-section";

vi.mock("@/i18n/navigation", () => ({
  Link: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}));

const messages = {
  backgroundTasks: {
    dropdown: {
      title: "Tasks",
      viewAll: "+{count} more - View all",
      viewPage: "View",
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
        <DropdownMenu defaultOpen>
          <DropdownMenuTrigger>Trigger</DropdownMenuTrigger>
          <DropdownMenuContent>{ui}</DropdownMenuContent>
        </DropdownMenu>
      </NextIntlClientProvider>
    </QueryClientProvider>
  );
};

const renderWithProvidersNoWrapper = (ui: React.ReactElement) => {
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

describe("TasksDropdownSection", () => {
  beforeEach(() => {
    vi.spyOn(useUserActiveTasksModule, "useUserActiveTasks");
    vi.spyOn(hooksModule, "useActiveTasks").mockReturnValue([]);
    vi.useFakeTimers({ shouldAdvanceTime: true });
  });

  afterEach(() => {
    cleanup();
    vi.restoreAllMocks();
    vi.useRealTimers();
  });

  it("should render nothing when no active tasks", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 0,
      data: { tasks: [] },
      isLoading: false,
    });

    const { container } = renderWithProvidersNoWrapper(<TasksDropdownSection />);

    expect(container.firstChild).toBeNull();
  });

  it("should render section title when tasks exist", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [createTask()] },
      isLoading: false,
    });

    renderWithProviders(<TasksDropdownSection />);

    expect(screen.getByText("Tasks")).toBeInTheDocument();
  });

  it("should render task with owner/repo as display name", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [createTask()] },
      isLoading: false,
    });

    renderWithProviders(<TasksDropdownSection />);

    expect(screen.getByText("test-owner/test-repo")).toBeInTheDocument();
  });

  it("should render task type label", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [createTask()] },
      isLoading: false,
    });

    renderWithProviders(<TasksDropdownSection />);

    expect(screen.getByText("Analysis")).toBeInTheDocument();
  });

  it("should render queued status for queued tasks", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [createTask({ startedAt: undefined, status: "queued" })] },
      isLoading: false,
    });

    renderWithProviders(<TasksDropdownSection />);

    expect(screen.getByText("Queued")).toBeInTheDocument();
  });

  it("should render processing status with elapsed time", () => {
    const startedAt = new Date(Date.now() - 45000).toISOString();
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [createTask({ startedAt, status: "analyzing" })] },
      isLoading: false,
    });

    renderWithProviders(<TasksDropdownSection />);

    expect(screen.getByText("Processing 45s")).toBeInTheDocument();
  });

  it("should render View link for analysis tasks", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [createTask()] },
      isLoading: false,
    });

    renderWithProviders(<TasksDropdownSection />);

    expect(screen.getByText("View")).toBeInTheDocument();
    expect(screen.getByRole("link")).toHaveAttribute("href", "/analyze/test-owner/test-repo");
  });

  it("should render multiple tasks", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
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

    renderWithProviders(<TasksDropdownSection />);

    expect(screen.getByText("test-owner/test-repo")).toBeInTheDocument();
    expect(screen.getByText("other/repo2")).toBeInTheDocument();
  });
});
