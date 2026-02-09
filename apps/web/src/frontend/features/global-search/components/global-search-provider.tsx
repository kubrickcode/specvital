"use client";

import type { ReactNode } from "react";

import { GlobalSearchDialog } from "./global-search-dialog";
import { useKeyboardShortcut } from "../hooks/use-keyboard-shortcut";

type GlobalSearchProviderProps = {
  children: ReactNode;
};

export const GlobalSearchProvider = ({ children }: GlobalSearchProviderProps) => {
  useKeyboardShortcut();

  return (
    <>
      {children}
      <GlobalSearchDialog />
    </>
  );
};
