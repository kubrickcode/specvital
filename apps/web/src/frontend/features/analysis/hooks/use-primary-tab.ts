"use client";

import { parseAsStringLiteral, useQueryState } from "nuqs";
import { useEffect, useTransition } from "react";

import { DEFAULT_PRIMARY_TAB, PRIMARY_TABS } from "../types/primary-tab";

const primaryTabParser = parseAsStringLiteral(PRIMARY_TABS).withDefault(DEFAULT_PRIMARY_TAB);

export const usePrimaryTab = () => {
  const [isPending, startTransition] = useTransition();
  const [tab, setTab] = useQueryState("tab", primaryTabParser.withOptions({ startTransition }));
  const [viewParam, setViewParam] = useQueryState("view");

  // URL backward compatibility: migrate ?view=document to ?tab=spec
  useEffect(() => {
    if (viewParam === "document") {
      setTab("spec");
      setViewParam(null);
    }
  }, [viewParam, setTab, setViewParam]);

  return { isPending, setTab, tab } as const;
};
