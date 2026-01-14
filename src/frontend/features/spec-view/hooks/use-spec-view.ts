"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useRef } from "react";

import { fetchSpecDocument, requestSpecGeneration } from "../api";
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
  document: (analysisId: string) => [...specViewKeys.all, "document", analysisId] as const,
};

type UseSpecViewReturn = {
  data: SpecDocument | null;
  error: Error | null;
  generationStatus: SpecGenerationStatusEnum | null;
  isGenerating: boolean;
  isLoading: boolean;
  isPollingTimeout: boolean;
  requestGenerate: (language?: SpecLanguage, isForceRegenerate?: boolean) => void;
};

export const useSpecView = (analysisId: string): UseSpecViewReturn => {
  const queryClient = useQueryClient();
  const pollingStartTimeRef = useRef<number | null>(null);

  useEffect(() => {
    pollingStartTimeRef.current = null;
  }, [analysisId]);

  const query = useQuery({
    enabled: Boolean(analysisId),
    queryFn: () => fetchSpecDocument(analysisId),
    queryKey: specViewKeys.document(analysisId),
    refetchInterval: (query) => {
      const response = query.state.data;

      if (!response) return false;

      if (isDocumentCompleted(response)) {
        pollingStartTimeRef.current = null;
        return false;
      }

      if (isDocumentGenerating(response)) {
        const status = response.generationStatus.status;
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
    onSuccess: () => {
      pollingStartTimeRef.current = null;
      queryClient.invalidateQueries({ queryKey: specViewKeys.document(analysisId) });
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

  return {
    data,
    error: query.error || generateMutation.error,
    generationStatus,
    isGenerating,
    isLoading,
    isPollingTimeout,
    requestGenerate,
  };
};
