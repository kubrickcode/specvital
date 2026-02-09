import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { cleanup, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import * as hooksModule from "../hooks";
import * as useUserActiveTasksModule from "../use-user-active-tasks";
import { TaskBadge } from "./task-badge";

const messages = {
  backgroundTasks: {
    badge: {
      ariaLabel: "{count, plural, =1 {1 task in progress} other {# tasks in progress}}",
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

describe("TaskBadge", () => {
  beforeEach(() => {
    vi.spyOn(useUserActiveTasksModule, "useUserActiveTasks");
    vi.spyOn(hooksModule, "useActiveTasks").mockReturnValue([]);
  });

  afterEach(() => {
    cleanup();
    vi.restoreAllMocks();
  });

  it("should render nothing when no active tasks", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 0,
      data: null,
      isLoading: false,
    });

    const { container } = renderWithProviders(<TaskBadge />);

    expect(container.firstChild).toBeNull();
  });

  it("should render badge when tasks are active", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [] },
      isLoading: false,
    });

    renderWithProviders(<TaskBadge />);

    const badge = screen.getByLabelText("1 task in progress");
    expect(badge).toBeInTheDocument();
  });

  it("should render badge with multiple tasks", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 3,
      data: { tasks: [] },
      isLoading: false,
    });

    renderWithProviders(<TaskBadge />);

    const badge = screen.getByLabelText("3 tasks in progress");
    expect(badge).toBeInTheDocument();
  });

  it("should apply custom className", () => {
    vi.mocked(useUserActiveTasksModule.useUserActiveTasks).mockReturnValue({
      activeTaskCount: 1,
      data: { tasks: [] },
      isLoading: false,
    });

    renderWithProviders(<TaskBadge className="custom-class" />);

    const badge = screen.getByLabelText("1 task in progress");
    expect(badge).toHaveClass("custom-class");
  });
});
