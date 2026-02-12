"use client";

import { useCallback, useEffect, useRef, useState, type ReactNode } from "react";

import { DocumentNavigationContext } from "./document-navigation-context-value";
import type { DocumentNavigationContextValue, VirtualizerRef } from "./document-navigation-types";
import type { FlatSpecItem, SpecDocument } from "../types";

type DocumentNavigationProviderProps = {
  children: ReactNode;
  document: SpecDocument;
};

export const DocumentNavigationProvider = ({
  children,
  document,
}: DocumentNavigationProviderProps) => {
  // State
  const [activeSection, setActiveSection] = useState<string | null>(null);
  const [expandedDomains, setExpandedDomains] = useState<Set<string>>(
    () => new Set(document.domains.map((d) => d.id))
  );
  const [expandedFeatures, setExpandedFeatures] = useState<Set<string>>(
    () => new Set(document.domains.flatMap((d) => d.features.map((f) => f.id)))
  );

  // Trigger for pending scroll effect
  const [flatItemsVersion, setFlatItemsVersion] = useState(0);

  // Virtualizer refs
  const virtualizerRef = useRef<VirtualizerRef | null>(null);
  const flatItemsRef = useRef<FlatSpecItem[]>([]);

  // Pending scroll target (for when section needs expansion first)
  const pendingScrollRef = useRef<string | null>(null);

  // Reset expanded state when document structure changes
  useEffect(() => {
    setExpandedDomains(new Set(document.domains.map((d) => d.id)));
    setExpandedFeatures(new Set(document.domains.flatMap((d) => d.features.map((f) => f.id))));
  }, [document.domains]);

  // Handle initial hash on mount
  useEffect(() => {
    const hash = window.location.hash.slice(1);
    if (hash) {
      try {
        const targetId = decodeURIComponent(hash);
        setActiveSection(targetId);
        pendingScrollRef.current = targetId;
      } catch {
        // Invalid hash, ignore
      }
    }
  }, []);

  // Process pending scroll when flatItems update
  useEffect(() => {
    if (!pendingScrollRef.current || !virtualizerRef.current) return;

    const targetId = pendingScrollRef.current;
    const index = flatItemsRef.current.findIndex((item) => {
      if (item.type === "domain-header") return `domain-${item.domainId}` === targetId;
      if (item.type === "feature-header") return `feature-${item.featureId}` === targetId;
      return false;
    });

    if (index >= 0) {
      requestAnimationFrame(() => {
        virtualizerRef.current?.scrollToIndex(index, {
          align: "start",
        });
      });
      pendingScrollRef.current = null;
    }
  }, [flatItemsVersion]);

  const registerVirtualizer = useCallback(
    (virtualizer: VirtualizerRef, flatItems: FlatSpecItem[]) => {
      const lengthChanged = flatItemsRef.current.length !== flatItems.length;
      virtualizerRef.current = virtualizer;
      flatItemsRef.current = flatItems;

      // Only trigger effect when flatItems actually changed
      if (lengthChanged) {
        setFlatItemsVersion((v) => v + 1);
      }
    },
    []
  );

  const toggleDomain = useCallback((domainId: string) => {
    setExpandedDomains((prev) => {
      const next = new Set(prev);
      if (next.has(domainId)) {
        next.delete(domainId);
      } else {
        next.add(domainId);
      }
      return next;
    });
  }, []);

  const toggleFeature = useCallback((featureId: string) => {
    setExpandedFeatures((prev) => {
      const next = new Set(prev);
      if (next.has(featureId)) {
        next.delete(featureId);
      } else {
        next.add(featureId);
      }
      return next;
    });
  }, []);

  const scrollToSection = useCallback(
    (sectionId: string) => {
      // Update URL hash
      window.history.replaceState(null, "", `#${sectionId}`);
      setActiveSection(sectionId);

      // Determine if we need to expand parent sections
      let needsExpand = false;
      let targetDomainId: string | null = null;

      if (sectionId.startsWith("domain-")) {
        targetDomainId = sectionId.replace("domain-", "");
        if (!expandedDomains.has(targetDomainId)) {
          needsExpand = true;
        }
      } else if (sectionId.startsWith("feature-")) {
        const targetFeatureId = sectionId.replace("feature-", "");

        // Find parent domain
        const parentDomain = document.domains.find((d) =>
          d.features.some((f) => f.id === targetFeatureId)
        );

        if (parentDomain) {
          targetDomainId = parentDomain.id;

          // Check if parent domain needs to be expanded
          if (!expandedDomains.has(targetDomainId)) {
            needsExpand = true;
          }
        }
      }

      // Expand sections if needed
      if (needsExpand && targetDomainId) {
        const domainIdToExpand = targetDomainId;
        setExpandedDomains((prev) => new Set(prev).add(domainIdToExpand));

        // Set pending scroll - will be processed after re-render
        pendingScrollRef.current = sectionId;
        return;
      }

      // Find index in flatItems and scroll
      const index = flatItemsRef.current.findIndex((item) => {
        if (item.type === "domain-header") return `domain-${item.domainId}` === sectionId;
        if (item.type === "feature-header") return `feature-${item.featureId}` === sectionId;
        return false;
      });

      if (index >= 0 && virtualizerRef.current) {
        virtualizerRef.current.scrollToIndex(index, {
          align: "start",
        });
      }
    },
    [document.domains, expandedDomains]
  );

  const value: DocumentNavigationContextValue = {
    activeSection,
    expandedDomains,
    expandedFeatures,
    registerVirtualizer,
    scrollToSection,
    setActiveSection,
    toggleDomain,
    toggleFeature,
  };

  return (
    <DocumentNavigationContext.Provider value={value}>
      {children}
    </DocumentNavigationContext.Provider>
  );
};
