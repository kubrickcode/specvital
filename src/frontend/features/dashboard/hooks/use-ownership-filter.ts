"use client";

import { parseAsStringLiteral, useQueryState } from "nuqs";

import { OWNERSHIP_FILTER_VALUES } from "@/lib/api/types";

const ownershipParser = parseAsStringLiteral(OWNERSHIP_FILTER_VALUES).withDefault("all");

export const useOwnershipFilter = () => {
  const [ownership, setOwnership] = useQueryState("ownership", ownershipParser);

  return { ownership, setOwnership } as const;
};
