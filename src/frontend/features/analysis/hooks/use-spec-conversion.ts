"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";

import { convertSpecView } from "../api";
import type { ConversionLanguage, ConvertSpecViewResponse } from "../types";

export const specConversionKeys = {
  all: ["spec-conversion"] as const,
  detail: (owner: string, repo: string, commitSha: string, language: ConversionLanguage) =>
    [...specConversionKeys.all, owner, repo, commitSha, language] as const,
};

type UseSpecConversionParams = {
  commitSha: string;
  enabled?: boolean;
  language: ConversionLanguage;
  owner: string;
  repo: string;
};

type UseSpecConversionReturn = {
  data: ConvertSpecViewResponse | null;
  error: Error | null;
  isLoading: boolean;
  isRegenerating: boolean;
  regenerate: () => void;
};

export const useSpecConversion = ({
  commitSha,
  enabled = true,
  language,
  owner,
  repo,
}: UseSpecConversionParams): UseSpecConversionReturn => {
  const queryClient = useQueryClient();

  const queryKey = specConversionKeys.detail(owner, repo, commitSha, language);

  const query = useQuery({
    enabled,
    queryFn: () => convertSpecView({ commitSha, language, owner, repo }),
    queryKey,
    staleTime: 5 * 60 * 1000,
  });

  const regenerateMutation = useMutation({
    mutationFn: () => convertSpecView({ commitSha, isForceRefresh: true, language, owner, repo }),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKey, data);
    },
  });

  const regenerate = useCallback(() => {
    regenerateMutation.mutate();
  }, [regenerateMutation]);

  return {
    data: query.data ?? null,
    error: query.error ?? regenerateMutation.error ?? null,
    isLoading: query.isPending,
    isRegenerating: regenerateMutation.isPending,
    regenerate,
  };
};
