"use client";

import { createContext } from "react";

import type { DocumentNavigationContextValue } from "./document-navigation-types";

export const DocumentNavigationContext = createContext<DocumentNavigationContextValue | null>(null);
