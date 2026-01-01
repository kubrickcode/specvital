"use client";

import { parseAsStringLiteral, useQueryState } from "nuqs";

export type ViewFilter = "all" | "mine" | "starred";

const VIEW_FILTER_OPTIONS: ViewFilter[] = ["all", "mine", "starred"];

const viewFilterParser = parseAsStringLiteral(VIEW_FILTER_OPTIONS).withDefault("all");

export const useViewFilter = () => {
  const [viewFilter, setViewFilter] = useQueryState("view", viewFilterParser);

  return { setViewFilter, viewFilter } as const;
};
