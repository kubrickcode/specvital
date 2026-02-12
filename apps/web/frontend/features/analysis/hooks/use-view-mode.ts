"use client";

import { parseAsStringLiteral, useQueryState } from "nuqs";

import { DEFAULT_VIEW_MODE, VIEW_MODES } from "../types";

const viewModeParser = parseAsStringLiteral(VIEW_MODES).withDefault(DEFAULT_VIEW_MODE);

export const useViewMode = () => {
  const [viewMode, setViewMode] = useQueryState("view", viewModeParser);

  return { setViewMode, viewMode } as const;
};
