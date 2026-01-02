"use client";

import { parseAsStringLiteral, useQueryState } from "nuqs";

// TODO: Add "organization" option when backend API supports ownership parameter
export type OwnershipFilter = "all" | "mine";

export const OWNERSHIP_FILTER_OPTIONS: OwnershipFilter[] = ["all", "mine"];

const ownershipFilterParser = parseAsStringLiteral(OWNERSHIP_FILTER_OPTIONS).withDefault("all");

export const useOwnershipFilter = () => {
  const [ownershipFilter, setOwnershipFilter] = useQueryState("ownership", ownershipFilterParser);

  return { ownershipFilter, setOwnershipFilter } as const;
};
