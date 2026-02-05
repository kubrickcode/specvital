"use client";

import { useQuery } from "@tanstack/react-query";

import { fetchSubscription } from "../api/subscription";

export const useSubscription = (enabled = true) => {
  return useQuery({
    enabled,
    queryFn: fetchSubscription,
    queryKey: ["user", "subscription"],
  });
};
