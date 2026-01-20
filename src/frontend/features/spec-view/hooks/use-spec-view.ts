"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useEffect, useRef } from "react";
import { toast } from "sonner";

import { useRouter } from "@/i18n/navigation";
import { ROUTES } from "@/lib/routes";

import {
  fetchSpecDocument,
  NoSubscriptionError,
  QuotaExceededError,
  requestSpecGeneration,
} from "../api";
import type {
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
  document: (analysisId: string, language?: SpecLanguage) =>
    language
      ? ([...specViewKeys.all, "document", analysisId, language] as const)
      : ([...specViewKeys.all, "document", analysisId] as const),
};

type UseSpecViewOptions = {
  language?: SpecLanguage;
};

type UseSpecViewReturn = {
  currentLanguage: SpecLanguage | undefined;
  data: SpecDocument | null;
  error: Error | null;
  generationStatus: SpecGenerationStatusEnum | null;
  isGenerating: boolean;
  isLoading: boolean;
  isPollingTimeout: boolean;
  requestGenerate: (language?: SpecLanguage, isForceRegenerate?: boolean) => void;
};

export const useSpecView = (
  analysisId: string,
  options: UseSpecViewOptions = {}
): UseSpecViewReturn => {
  const { language } = options;
  const t = useTranslations("specView.toast");
  const queryClient = useQueryClient();
  const router = useRouter();
  const pollingStartTimeRef = useRef<number | null>(null);

  // Track previous status to detect completion transition
  const previousStatusRef = useRef<string | null>(null);

  useEffect(() => {
    pollingStartTimeRef.current = null;
    previousStatusRef.current = null;
  }, [analysisId, language]);

  const query = useQuery({
    enabled: Boolean(analysisId),
    queryFn: () => fetchSpecDocument(analysisId, language),
    queryKey: specViewKeys.document(analysisId, language),
    refetchInterval: (query) => {
      const response = query.state.data;

      if (!response) return false;

      if (isDocumentCompleted(response)) {
        pollingStartTimeRef.current = null;
        return false;
      }

      if (isDocumentGenerating(response)) {
        const status = response.generationStatus.status;

        // Detect completion transition and invalidate query for fresh document data
        if (
          status === "completed" &&
          previousStatusRef.current !== "completed" &&
          previousStatusRef.current !== null
        ) {
          previousStatusRef.current = status;
          // Invalidate to fetch completed document data
          queryClient.invalidateQueries({
            queryKey: specViewKeys.document(analysisId, language),
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

  const generationStatus: SpecGenerationStatusEnum | null =
    response && isDocumentGenerating(response) ? response.generationStatus.status : null;

  const isGenerating = generationStatus === "pending" || generationStatus === "running";
  const isLoading = query.isPending || generateMutation.isPending;

  const isPollingTimeout =
    pollingStartTimeRef.current !== null &&
    Date.now() - pollingStartTimeRef.current > MAX_POLLING_DURATION_MS &&
    isGenerating;

  // Extract language from completed document data
  const currentLanguage = data?.language;

  return {
    currentLanguage,
    data,
    error: query.error || generateMutation.error,
    generationStatus,
    isGenerating,
    isLoading,
    isPollingTimeout,
    requestGenerate,
  };
};
