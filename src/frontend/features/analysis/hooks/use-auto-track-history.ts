"use client";

import { useEffect, useRef } from "react";

import { useAuth } from "@/features/auth";
import { invalidationEvents, useInvalidationTrigger } from "@/lib/query";

import { addToHistory } from "../api";

export const useAutoTrackHistory = (owner: string, repo: string, isReady: boolean): void => {
  const { isAuthenticated, isLoading: isAuthLoading } = useAuth();
  const triggerInvalidation = useInvalidationTrigger();
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
        triggerInvalidation(invalidationEvents.HISTORY_CHANGED);
      })
      .catch(() => {
        // Silent fail - history tracking is not critical
      });
  }, [isAuthLoading, isAuthenticated, isReady, owner, repo, triggerInvalidation]);
};
