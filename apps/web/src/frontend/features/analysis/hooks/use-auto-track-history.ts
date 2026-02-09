"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useEffect, useRef } from "react";

import { useAuth } from "@/features/auth";
import { paginatedRepositoriesKeys } from "@/features/dashboard";

import { addToHistory } from "../api";

export const useAutoTrackHistory = (owner: string, repo: string, isReady: boolean): void => {
  const queryClient = useQueryClient();
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const hasTracked = useRef(false);

  useEffect(() => {
    hasTracked.current = false;
  }, [owner, repo]);

  useEffect(() => {
    if (isAuthLoading || !isReady || !isAuthenticated || hasTracked.current) {
      return;
    }

    hasTracked.current = true;

    addToHistory(owner, repo)
      .then(() => {
        queryClient.invalidateQueries({ queryKey: paginatedRepositoriesKeys.all });
      })
      .catch(() => {
        // Silent fail - history tracking is not critical
      });
  }, [isAuthLoading, isAuthenticated, isReady, owner, repo, queryClient]);
};
