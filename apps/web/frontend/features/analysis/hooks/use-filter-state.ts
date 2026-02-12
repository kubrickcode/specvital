"use client";

import { parseAsArrayOf, parseAsString, useQueryState } from "nuqs";
import { useTransition } from "react";

import type { TestStatus } from "@/lib/api";

const VALID_STATUSES: TestStatus[] = ["active", "focused", "skipped", "todo", "xfail"];

const queryParser = parseAsString.withDefault("");
const arrayParser = parseAsArrayOf(parseAsString, ",").withDefault([]);

export const useFilterState = () => {
  const [isPending, startTransition] = useTransition();
  const [frameworks, setFrameworks] = useQueryState(
    "frameworks",
    arrayParser.withOptions({ startTransition })
  );
  const [query, setQuery] = useQueryState("q", queryParser.withOptions({ startTransition }));
  const [rawStatuses, setRawStatuses] = useQueryState(
    "statuses",
    arrayParser.withOptions({ startTransition })
  );

  const statuses = rawStatuses.filter((s): s is TestStatus =>
    VALID_STATUSES.includes(s as TestStatus)
  );

  const setStatuses = (value: TestStatus[] | null) => {
    setRawStatuses(value);
  };

  return {
    frameworks,
    isPending,
    query,
    setFrameworks,
    setQuery,
    setStatuses,
    statuses,
  } as const;
};
