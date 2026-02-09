"use client";

import { Building2, List, User, Users } from "lucide-react";
import { parseAsStringLiteral, useQueryState } from "nuqs";

import type { OwnershipFilterParam } from "@/lib/api/types";

export type OwnershipFilter = OwnershipFilterParam;

export const OWNERSHIP_FILTER_OPTIONS: OwnershipFilter[] = [
  "all",
  "mine",
  "organization",
  "others",
];

export const OWNERSHIP_FILTER_ICONS: Record<
  OwnershipFilter,
  React.ComponentType<{ className?: string }>
> = {
  all: List,
  mine: User,
  organization: Building2,
  others: Users,
};

const ownershipFilterParser = parseAsStringLiteral(OWNERSHIP_FILTER_OPTIONS).withDefault("all");

export const useOwnershipFilter = () => {
  const [ownershipFilter, setOwnershipFilter] = useQueryState("ownership", ownershipFilterParser);

  return { ownershipFilter, setOwnershipFilter } as const;
};
