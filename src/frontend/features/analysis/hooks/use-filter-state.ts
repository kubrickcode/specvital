"use client";

import { parseAsArrayOf, parseAsString, useQueryState } from "nuqs";
import { useMemo } from "react";

import type { TestStatus } from "@/lib/api";

const VALID_STATUSES: TestStatus[] = ["active", "focused", "skipped", "todo", "xfail"];

const queryParser = parseAsString.withDefault("");
const arrayParser = parseAsArrayOf(parseAsString, ",").withDefault([]);

export const useFilterState = () => {
  const [frameworks, setFrameworks] = useQueryState("frameworks", arrayParser);
  const [query, setQuery] = useQueryState("q", queryParser);
  const [rawStatuses, setRawStatuses] = useQueryState("statuses", arrayParser);

  const statuses = useMemo(
    () => rawStatuses.filter((s): s is TestStatus => VALID_STATUSES.includes(s as TestStatus)),
    [rawStatuses]
  );

  const setStatuses = (value: TestStatus[] | null) => {
    setRawStatuses(value);
  };

  return {
    frameworks,
    query,
    setFrameworks,
    setQuery,
    setStatuses,
    statuses,
  } as const;
};
