"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useEffect, useRef } from "react";
import { toast } from "sonner";

import { useRouter } from "@/i18n/navigation";
import { ROUTES } from "@/lib/routes";

import {
  fetchSpecDocument,
  ForbiddenError,
  NoSubscriptionError,
  QuotaExceededError,
  requestSpecGeneration,
  UnauthorizedError,
} from "../api";
import { repoSpecViewKeys } from "./use-repo-spec-view";
import type {
  BehaviorCacheStats,
  SpecDocument,
  SpecDocumentResponse,
  SpecGenerationStatusEnum,
  SpecLanguage,
} from "../types";
import { isDocumentCompleted, isDocumentGenerating } from "../types";

const POLLING_INTERVAL_MS = 5000;
const MAX_POLLING_DURATION_MS = 5 * 60 * 1000; // 5 minutes
const DEFAULT_LANGUAGE: SpecLanguage = "Korean";

export const specViewKeys = {
  all: ["spec-view"] as const,
  document: (analysisId: string, language?: SpecLanguage, version?: number) =>
    version !== undefined
      ? ([...specViewKeys.all, "document", analysisId, language, version] as const)
      : language
        ? ([...specViewKeys.all, "document", analysisId, language] as const)
        : ([...specViewKeys.all, "document", analysisId] as const),
  versions: (analysisId: string, language: SpecLanguage) =>
    [...specViewKeys.all, "versions", analysisId, language] as const,
};

type UseSpecViewOptions = {
  language?: SpecLanguage;
  version?: number;
};

/**
 * Unified generation state for UI consumption.
 * - idle: No generation in progress, no document
 * - requesting: Mutation in progress (sending request to server)
 * - pending: Server queued the generation
 * - running: Server is actively generating
 * - completed: Document is ready
 * - failed: Generation failed
 */
export type GenerationState =
  | "idle"
  | "requesting"
  | "pending"
  | "running"
  | "completed"
  | "failed";

export type AccessErrorType = "unauthorized" | "forbidden" | null;

type UseSpecViewReturn = {
  accessError: AccessErrorType;
  behaviorCacheStats: BehaviorCacheStats | undefined;
  currentLanguage: SpecLanguage | undefined;
  data: SpecDocument | null;
  error: Error | null;
  generationState: GenerationState;
  isFetching: boolean;
  isLoading: boolean;
  isPollingTimeout: boolean;
  requestGenerate: (language?: SpecLanguage, isForceRegenerate?: boolean) => void;
  serverStatus: SpecGenerationStatusEnum | null;
};

export const useSpecView = (
  analysisId: string,
  options: UseSpecViewOptions = {}
): UseSpecViewReturn => {
  const { language, version } = options;
  const t = useTranslations("specView.toast");
  const queryClient = useQueryClient();
  const router = useRouter();
  const pollingStartTimeRef = useRef<number | null>(null);

  // Track previous status to detect completion transition
  const previousStatusRef = useRef<string | null>(null);

  useEffect(() => {
    pollingStartTimeRef.current = null;
    previousStatusRef.current = null;
  }, [analysisId, language, version]);

  const query = useQuery({
    enabled: Boolean(analysisId),
    queryFn: () => fetchSpecDocument(analysisId, { language, version }),
    queryKey: specViewKeys.document(analysisId, language, version),
    refetchInterval: (query) => {
      const response = query.state.data;

      if (!response) return false;

      if (isDocumentCompleted(response)) {
        pollingStartTimeRef.current = null;
        return false;
      }

      if (isDocumentGenerating(response)) {
        const status = response.generationStatus.status;

        // Detect completion transition and invalidate queries for fresh document data
        if (
          status === "completed" &&
          previousStatusRef.current !== "completed" &&
          previousStatusRef.current !== null
        ) {
          previousStatusRef.current = status;
          // Invalidate analysis-based spec view query
          queryClient.invalidateQueries({
            queryKey: specViewKeys.document(analysisId, language, version),
          });
          // Invalidate repo-based spec view queries for immediate document display
          queryClient.invalidateQueries({
            queryKey: repoSpecViewKeys.all,
          });
          pollingStartTimeRef.current = null;
          return false;
        }

        previousStatusRef.current = status;

        if (status === "completed" || status === "failed" || status === "not_found") {
          pollingStartTimeRef.current = null;
          return false;
        }

        // Start tracking polling time
        if (pollingStartTimeRef.current === null) {
          pollingStartTimeRef.current = Date.now();
        }

        // Stop polling after max duration
        if (Date.now() - pollingStartTimeRef.current > MAX_POLLING_DURATION_MS) {
          return false;
        }

        return POLLING_INTERVAL_MS;
      }

      return false;
    },
    retry: false,
    staleTime: 30000,
  });

  const generateMutation = useMutation({
    mutationFn: ({
      isForceRegenerate = false,
      language = DEFAULT_LANGUAGE,
    }: {
      isForceRegenerate?: boolean;
      language?: SpecLanguage;
    }) =>
      requestSpecGeneration({
        analysisId,
        isForceRegenerate,
        language,
      }),
    onError: (error) => {
      if (error instanceof NoSubscriptionError) {
        toast.error(t("noSubscription.title"), {
          action: {
            label: t("noSubscription.viewPlans"),
            onClick: () => {
              router.push(ROUTES.ACCOUNT);
            },
          },
          description: t("noSubscription.description"),
        });
        return;
      }

      if (error instanceof QuotaExceededError) {
        toast.error(t("quotaExceeded.title"), {
          action: {
            label: t("quotaExceeded.viewAccount"),
            onClick: () => {
              router.push(ROUTES.ACCOUNT);
            },
          },
          description: t("quotaExceeded.description", {
            limit: error.limit.toLocaleString(),
            used: error.used.toLocaleString(),
          }),
        });
        return;
      }

      toast.error(t("generateFailed.title"), {
        description: error instanceof Error ? error.message : String(error),
      });
    },
    onSuccess: (_data, variables) => {
      pollingStartTimeRef.current = null;
      previousStatusRef.current = null;
      // Invalidate language-specific query
      queryClient.invalidateQueries({
        queryKey: specViewKeys.document(analysisId, variables.language),
      });
      // Also invalidate language-less query to trigger polling in spec-panel
      queryClient.invalidateQueries({
        queryKey: specViewKeys.document(analysisId),
      });
    },
  });

  const requestGenerate = (
    language: SpecLanguage = DEFAULT_LANGUAGE,
    isForceRegenerate = false
  ) => {
    generateMutation.mutate({ isForceRegenerate, language });
  };

  const response: SpecDocumentResponse | undefined = query.data;
  const data = response && isDocumentCompleted(response) ? response.data : null;

  const serverStatus: SpecGenerationStatusEnum | null =
    response && isDocumentGenerating(response) ? response.generationStatus.status : null;

  const isRequesting = generateMutation.isPending;
  const isLoading = query.isPending || isRequesting;
  const isActiveGeneration = serverStatus === "pending" || serverStatus === "running";

  const isPollingTimeout =
    pollingStartTimeRef.current !== null &&
    Date.now() - pollingStartTimeRef.current > MAX_POLLING_DURATION_MS &&
    isActiveGeneration;

  // Extract language and cache stats from completed document data
  const currentLanguage = data?.language;
  const behaviorCacheStats = data?.behaviorCacheStats;

  // Compute unified generation state
  // Note: serverStatus takes priority over data existence for regeneration scenarios
  const generationState: GenerationState = (() => {
    if (isRequesting) return "requesting";
    if (serverStatus === "pending") return "pending";
    if (serverStatus === "running") return "running";
    if (serverStatus === "failed") return "failed";
    if (data) return "completed";
    return "idle";
  })();

  // Determine access error type from query error
  const accessError: AccessErrorType = (() => {
    if (query.error instanceof UnauthorizedError) return "unauthorized";
    if (query.error instanceof ForbiddenError) return "forbidden";
    return null;
  })();

  return {
    accessError,
    behaviorCacheStats,
    currentLanguage,
    data,
    error: query.error || generateMutation.error,
    generationState,
    isFetching: query.isFetching,
    isLoading,
    isPollingTimeout,
    requestGenerate,
    serverStatus,
  };
};
