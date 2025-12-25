"use client";

import { parseAsStringLiteral, useQueryState } from "nuqs";

const TAB_VALUES = ["bookmarked", "my-analyses"] as const;
type TabValue = (typeof TAB_VALUES)[number];

const tabParser = parseAsStringLiteral(TAB_VALUES).withDefault("bookmarked");

export const useTabState = () => {
  const [tab, setTab] = useQueryState("tab", tabParser);

  return { setTab, tab } as const;
};

export type { TabValue };
