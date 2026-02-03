"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { GitCommit, Loader2, RefreshCw, X } from "lucide-react";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";
import { toast } from "sonner";

import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { paginatedRepositoriesKeys } from "@/features/dashboard";
import { addTask, getTask, removeTask } from "@/lib/background-tasks";

import { fetchAnalysisStatus, triggerReanalyze } from "../api";
import { analysisKeys } from "../hooks/use-analysis";
import { updateStatusKeys, useUpdateStatus } from "../hooks/use-update-status";

type UpdateBannerProps = {
  owner: string;
  repo: string;
};

type BannerMessageKey = "analyzing" | "newCommitsDetected" | "parserUpdated";

const POLL_INTERVAL_MS = 1000;
const createTaskId = (owner: string, repo: string): string => `update-banner:${owner}/${repo}`;

const isTerminalStatus = (status: string): boolean => status === "completed" || status === "failed";

export const UpdateBanner = ({ owner, repo }: UpdateBannerProps) => {
  const t = useTranslations("analyze.updateBanner");
  const queryClient = useQueryClient();
  const [isDismissed, setIsDismissed] = useState(false);
  const [isPolling, setIsPolling] = useState(false);

  const {
    isLoading: isCheckingStatus,
    parserOutdated,
    status,
  } = useUpdateStatus(
    owner,
    repo,
    !isDismissed && !isPolling // Disable when dismissed or polling
  );

  // Polling query for reanalysis completion
  const pollingQuery = useQuery({
    enabled: isPolling,
    queryFn: () => fetchAnalysisStatus(owner, repo),
    queryKey: ["updateBannerPolling", owner, repo],
    refetchInterval: (query) => {
      const data = query.state.data;
      if (!data) return POLL_INTERVAL_MS;
      return isTerminalStatus(data.status) ? false : POLL_INTERVAL_MS;
    },
    // Prevent returning cached "completed" status from previous polling sessions
    gcTime: 0,
    staleTime: 0,
  });

  // Clean up polling state and dismiss banner on completion
  // Task removal is handled by AnalysisMonitor (global) for deduplication
  useEffect(() => {
    if (!isPolling || !pollingQuery.data) return;

    const { status: pollingStatus } = pollingQuery.data;
    if (!isTerminalStatus(pollingStatus)) return;

    const taskId = createTaskId(owner, repo);
    const task = getTask(taskId);
    if (task) {
      removeTask(taskId);
      if (pollingStatus === "failed") {
        toast.error(t("reanalyzeFailed"));
      }
    }

    queryClient.invalidateQueries({ queryKey: analysisKeys.detail(owner, repo) });
    queryClient.invalidateQueries({ queryKey: updateStatusKeys.detail(owner, repo) });
    queryClient.invalidateQueries({ queryKey: paginatedRepositoriesKeys.all });

    setIsPolling(false);
    setIsDismissed(true);
  }, [isPolling, pollingQuery.data, owner, repo, queryClient, t]);

  const reanalyzeMutation = useMutation({
    mutationFn: () => triggerReanalyze(owner, repo),
    onError: (error) => {
      toast.error(t("reanalyzeFailed"), {
        description: error instanceof Error ? error.message : String(error),
      });
    },
    onSuccess: () => {
      toast.success(t("reanalyzeQueued"));

      // Register task in TaskStore for Account badge
      const taskId = createTaskId(owner, repo);
      addTask({
        id: taskId,
        metadata: { owner, repo },
        startedAt: new Date().toISOString(),
        status: "processing",
        type: "analysis",
      });

      // Clear cached polling data before starting new polling session
      // Prevents refetchInterval from seeing stale "completed" status
      queryClient.removeQueries({ queryKey: ["updateBannerPolling", owner, repo] });

      // Start polling for completion
      setIsPolling(true);
    },
  });

  const handleDismiss = () => {
    setIsDismissed(true);
  };

  const handleReanalyze = () => {
    reanalyzeMutation.mutate();
  };

  const hasNewCommits = status === "new-commits";
  const isReanalyzing = reanalyzeMutation.isPending || isPolling;
  const shouldShowBanner =
    !isDismissed && !isCheckingStatus && (hasNewCommits || parserOutdated || isPolling);

  // Don't show banner if dismissed, loading, or no updates available
  if (!shouldShowBanner) {
    return null;
  }

  // Determine message based on state
  let messageKey: BannerMessageKey;
  if (isPolling) {
    messageKey = "analyzing";
  } else if (hasNewCommits) {
    messageKey = "newCommitsDetected";
  } else {
    messageKey = "parserUpdated";
  }

  return (
    <Alert className="mb-4 border-amber-200 bg-amber-50 dark:border-amber-900 dark:bg-amber-950/30">
      <GitCommit className="size-4 text-amber-600 dark:text-amber-400" />
      <AlertDescription className="flex items-center justify-between gap-4">
        <span className="text-amber-800 dark:text-amber-200">{t(messageKey)}</span>
        <div className="flex shrink-0 items-center gap-2">
          <Button
            className="h-7 gap-1.5 px-2.5 text-xs"
            disabled={isReanalyzing}
            onClick={handleReanalyze}
            size="sm"
            variant="default"
          >
            {isReanalyzing ? (
              <Loader2 className="size-3 animate-spin" />
            ) : (
              <RefreshCw className="size-3" />
            )}
            {t("updateNow")}
          </Button>
          {!isPolling && (
            <Button
              aria-label={t("dismissLabel")}
              className="size-7 p-0"
              onClick={handleDismiss}
              size="sm"
              variant="ghost"
            >
              <X className="size-4" />
            </Button>
          )}
        </div>
      </AlertDescription>
    </Alert>
  );
};
