"use client";

import { parseAsString, useQueryState } from "nuqs";

const queryParser = parseAsString.withDefault("");

export const useFilterState = () => {
  const [query, setQuery] = useQueryState("q", queryParser);

  return { query, setQuery } as const;
};
