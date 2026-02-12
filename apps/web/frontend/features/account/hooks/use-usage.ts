"use client";

import { useQuery } from "@tanstack/react-query";

import { fetchCurrentUsage } from "../api/usage";

export const useUsage = (enabled = true) => {
  return useQuery({
    enabled,
    queryFn: fetchCurrentUsage,
    queryKey: ["user", "usage"],
  });
};
