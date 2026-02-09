"use client";

import { useQuery } from "@tanstack/react-query";

import { apiFetch, parseJsonResponse } from "@/lib/api/client";
import type { UserActiveTasksResponse } from "@/lib/api/types";

const POLL_INTERVAL_MS = 3000;

export const userActiveTasksKeys = {
  all: ["user", "active-tasks"] as const,
};

const fetchUserActiveTasks = async (): Promise<UserActiveTasksResponse> => {
  const response = await apiFetch("/api/user/active-tasks");
  return parseJsonResponse<UserActiveTasksResponse>(response);
};

type UseUserActiveTasksOptions = {
  enabled?: boolean;
};

type UseUserActiveTasksReturn = {
  activeTaskCount: number;
  data: UserActiveTasksResponse | null;
  isLoading: boolean;
};

/**
 * Fetches active background tasks for authenticated users from the server.
 * Polls every 3 seconds while there are active tasks.
 *
 * IMPORTANT: Disabled by default. Callers must explicitly pass `enabled: true`
 * (typically `enabled: isAuthenticated`) to prevent 401 errors for unauthenticated users.
 */
export const useUserActiveTasks = (
  options: UseUserActiveTasksOptions = {}
): UseUserActiveTasksReturn => {
  const { enabled = false } = options;

  const query = useQuery({
    enabled,
    queryFn: fetchUserActiveTasks,
    queryKey: userActiveTasksKeys.all,
    // Always poll - stopping causes race conditions when tasks complete
    // Server returns empty array when no active tasks (minimal overhead)
    refetchInterval: POLL_INTERVAL_MS,
  });

  const tasks = query.data?.tasks ?? [];

  return {
    activeTaskCount: tasks.length,
    data: query.data ?? null,
    isLoading: query.isLoading,
  };
};
