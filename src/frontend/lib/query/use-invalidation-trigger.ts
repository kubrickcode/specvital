"use client";

import { useQueryClient } from "@tanstack/react-query";

import { invalidationRegistry, type InvalidationEvent } from "./invalidation-registry";

export const useInvalidationTrigger = () => {
  const queryClient = useQueryClient();

  const trigger = (event: InvalidationEvent) => {
    const queryKeys = invalidationRegistry[event];
    if (!queryKeys) return;

    queryKeys.forEach((queryKey) => {
      queryClient.invalidateQueries({ exact: false, queryKey });
    });
  };

  return trigger;
};
