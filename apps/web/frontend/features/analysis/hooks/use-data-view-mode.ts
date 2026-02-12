"use client";

import { parseAsStringLiteral, useQueryState } from "nuqs";

import { DATA_VIEW_MODES, DEFAULT_DATA_VIEW_MODE } from "../types/data-view-mode";

const dataViewModeParser =
  parseAsStringLiteral(DATA_VIEW_MODES).withDefault(DEFAULT_DATA_VIEW_MODE);

export const useDataViewMode = () => {
  const [dataViewMode, setDataViewMode] = useQueryState("view", dataViewModeParser);

  return { dataViewMode, setDataViewMode } as const;
};
