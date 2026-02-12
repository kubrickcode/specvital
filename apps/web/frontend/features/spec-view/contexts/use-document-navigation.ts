"use client";

import { useContext } from "react";

import { DocumentNavigationContext } from "./document-navigation-context-value";
import type { DocumentNavigationContextValue } from "./document-navigation-types";

export const useDocumentNavigation = (): DocumentNavigationContextValue => {
  const context = useContext(DocumentNavigationContext);
  if (!context) {
    throw new Error("useDocumentNavigation must be used within DocumentNavigationProvider");
  }
  return context;
};
