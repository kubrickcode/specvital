"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
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
import type { RequestSpecGenerationResponse, SpecGenerationMode, SpecLanguage } from "../types";

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

export type AccessErrorType = "unauthorized" | "forbidden" | null;

type UseSpecViewReturn = {
  accessError: AccessErrorType;
  isRequesting: boolean;
  requestGenerate: (
    language?: SpecLanguage,
    mode?: SpecGenerationMode
  ) => Promise<RequestSpecGenerationResponse>;
};

/**
 * Hook for spec generation mutation and access error detection.
 * Document data is provided by useRepoSpecView; this hook's internal query
 * exists solely to detect auth errors (unauthorized/forbidden).
 * Polling is handled separately by useSpecGenerationStatus.
 */
export const useSpecView = (
  analysisId: string,
  options: UseSpecViewOptions = {}
): UseSpecViewReturn => {
  const { language, version } = options;
  const t = useTranslations("specView.toast");
  const queryClient = useQueryClient();
  const router = useRouter();

  const query = useQuery({
    enabled: Boolean(analysisId),
    queryFn: () => fetchSpecDocument(analysisId, { language, version }),
    queryKey: specViewKeys.document(analysisId, language, version),
    retry: false,
    staleTime: 30000,
  });

  const generateMutation = useMutation({
    mutationFn: ({
      generationMode = "initial",
      language = DEFAULT_LANGUAGE,
    }: {
      generationMode?: SpecGenerationMode;
      language?: SpecLanguage;
    }) =>
      requestSpecGeneration({
        analysisId,
        generationMode,
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
      // Invalidate queries to trigger refetch
      queryClient.invalidateQueries({
        queryKey: specViewKeys.document(analysisId, variables.language),
      });
      queryClient.invalidateQueries({
        queryKey: specViewKeys.document(analysisId),
      });
      // Also invalidate repo-based queries
      queryClient.invalidateQueries({
        queryKey: repoSpecViewKeys.all,
      });
    },
  });

  const requestGenerate = (
    language: SpecLanguage = DEFAULT_LANGUAGE,
    mode: SpecGenerationMode = "initial"
  ): Promise<RequestSpecGenerationResponse> => {
    return generateMutation.mutateAsync({ generationMode: mode, language });
  };

  const isRequesting = generateMutation.isPending;

  // Determine access error type from query error
  const accessError: AccessErrorType = (() => {
    if (query.error instanceof UnauthorizedError) return "unauthorized";
    if (query.error instanceof ForbiddenError) return "forbidden";
    return null;
  })();

  return {
    accessError,
    isRequesting,
    requestGenerate,
  };
};
