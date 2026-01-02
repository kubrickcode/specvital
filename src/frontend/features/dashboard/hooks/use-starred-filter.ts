"use client";

import { parseAsBoolean, useQueryState } from "nuqs";

const starredFilterParser = parseAsBoolean.withDefault(false);

export const useStarredFilter = () => {
  const [starredOnly, setStarredOnly] = useQueryState("starred", starredFilterParser);

  return { setStarredOnly, starredOnly } as const;
};
