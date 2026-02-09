"use client";

import { useQuery } from "@tanstack/react-query";

import { fetchVersionHistory } from "../api";
import type { SpecLanguage, VersionHistoryResponse } from "../types";
import { specViewKeys } from "./use-spec-view";

type UseVersionHistoryOptions = {
  enabled?: boolean;
};

type UseVersionHistoryReturn = {
  data: VersionHistoryResponse | undefined;
  error: Error | null;
  isLoading: boolean;
};

export const useVersionHistory = (
  analysisId: string,
  language: SpecLanguage | undefined,
  options: UseVersionHistoryOptions = {}
): UseVersionHistoryReturn => {
  const { enabled = true } = options;

  const query = useQuery({
    enabled: enabled && Boolean(analysisId) && Boolean(language),
    queryFn: () => fetchVersionHistory(analysisId, language!),
    queryKey: specViewKeys.versions(analysisId, language!),
  });

  return {
    data: query.data,
    error: query.error,
    isLoading: query.isPending,
  };
};
