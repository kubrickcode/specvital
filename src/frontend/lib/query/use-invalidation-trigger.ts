"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";

import { invalidationRegistry, type InvalidationEvent } from "./invalidation-registry";

export const useInvalidationTrigger = () => {
  const queryClient = useQueryClient();

  const trigger = useCallback(
    (event: InvalidationEvent) => {
      const queryKeys = invalidationRegistry[event];
      if (!queryKeys) return;

      queryKeys.forEach((queryKey) => {
        queryClient.invalidateQueries({ exact: false, queryKey });
      });
    },
    [queryClient]
  );

  return trigger;
};
