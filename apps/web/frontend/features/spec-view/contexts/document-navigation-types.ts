import type { Virtualizer } from "@tanstack/react-virtual";

import type { FlatSpecItem } from "../types";

export type VirtualizerRef = Virtualizer<Window, Element>;

export type DocumentNavigationContextValue = {
  // State
  activeSection: string | null;
  expandedDomains: Set<string>;
  expandedFeatures: Set<string>;

  // Virtualizer registration
  registerVirtualizer: (virtualizer: VirtualizerRef, flatItems: FlatSpecItem[]) => void;
  // Actions
  scrollToSection: (sectionId: string) => void;
  setActiveSection: (sectionId: string | null) => void;
  toggleDomain: (domainId: string) => void;

  toggleFeature: (featureId: string) => void;
};
