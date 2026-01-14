"use client";

import { ChevronDown, ChevronRight, List } from "lucide-react";
import { useTranslations } from "next-intl";
import { useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { cn } from "@/lib/utils";

import type {
  FilteredDocument,
  FilteredDomain,
  FilteredFeature,
} from "../hooks/use-document-filter";
import type { SpecDocument, SpecDomain } from "../types";
import { calculateDocumentStats } from "../utils/stats";

/**
 * Convert SpecDomain to FilteredDomain with default filter values
 */
const toFilteredDomain = (domain: SpecDomain): FilteredDomain => ({
  ...domain,
  features: domain.features.map((feature) => ({
    ...feature,
    behaviors: feature.behaviors.map((behavior) => ({
      ...behavior,
      hasMatch: true,
      highlightRanges: [],
    })),
    hasMatch: true,
    matchCount: feature.behaviors.length,
  })),
  hasMatch: true,
  matchCount: domain.features.reduce((sum, f) => sum + f.behaviors.length, 0),
});

type TocSidebarProps = {
  document: SpecDocument;
  filteredDocument?: FilteredDocument | null;
  hasFilter?: boolean;
  onNavigate?: (sectionId: string) => void;
};

type TocItemProps = {
  activeId: string | null;
  domain: FilteredDomain;
  expandedDomainId: string | null;
  hasFilter?: boolean;
  onExpandToggle: (domainId: string) => void;
  onNavigate: (sectionId: string) => void;
};

const TocItem = ({
  activeId,
  domain,
  expandedDomainId,
  hasFilter = false,
  onExpandToggle,
  onNavigate,
}: TocItemProps) => {
  const domainId = `domain-${domain.id}`;
  const isDomainActive = activeId === domainId;
  const isExpanded = expandedDomainId === domainId;

  const featureCount = domain.features.length;
  const behaviorCount = domain.features.reduce((sum, f) => sum + f.behaviors.length, 0);

  // Display count: show match count when filtering
  const displayCount = hasFilter
    ? `${domain.matchCount}/${behaviorCount}`
    : `${featureCount}/${behaviorCount}`;

  return (
    <div>
      <button
        aria-current={isDomainActive ? "location" : undefined}
        className={cn(
          "w-full flex items-center gap-2 px-3 py-2 text-sm rounded-md transition-colors text-left",
          "hover:bg-muted/50",
          isDomainActive && "bg-primary/10 text-primary font-medium"
        )}
        onClick={() => {
          onNavigate(domainId);
          onExpandToggle(isExpanded ? "" : domainId);
        }}
        type="button"
      >
        {domain.features.length > 0 ? (
          isExpanded ? (
            <ChevronDown className="h-3 w-3 flex-shrink-0 text-muted-foreground" />
          ) : (
            <ChevronRight className="h-3 w-3 flex-shrink-0 text-muted-foreground" />
          )
        ) : (
          <span className="w-3" />
        )}
        <span className="flex-1 truncate">{domain.name}</span>
        <span className="text-xs text-muted-foreground">{displayCount}</span>
      </button>

      {isExpanded && domain.features.length > 0 && (
        <div className="ml-5 mt-1 space-y-0.5">
          {domain.features.map((feature: FilteredFeature) => {
            const featureId = `feature-${feature.id}`;
            const isActive = activeId === featureId;
            const featureDisplayCount = hasFilter
              ? `${feature.matchCount}/${feature.behaviors.length}`
              : feature.behaviors.length;

            return (
              <button
                aria-current={isActive ? "location" : undefined}
                className={cn(
                  "w-full flex items-center gap-2 px-3 py-1.5 text-xs rounded-md transition-colors text-left",
                  "hover:bg-muted/50",
                  isActive && "bg-primary/10 text-primary font-medium"
                )}
                key={feature.id}
                onClick={() => onNavigate(featureId)}
                type="button"
              >
                <span className="truncate">{feature.name}</span>
                <span className="text-muted-foreground ml-auto">{featureDisplayCount}</span>
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
};

type TocContentProps = {
  activeId: string | null;
  document: SpecDocument;
  expandedDomainId: string | null;
  filteredDocument?: FilteredDocument | null;
  hasFilter?: boolean;
  onExpandToggle: (domainId: string) => void;
  onNavigate: (sectionId: string) => void;
  t: ReturnType<typeof useTranslations<"specView">>;
};

const TocContent = ({
  activeId,
  document,
  expandedDomainId,
  filteredDocument,
  hasFilter = false,
  onExpandToggle,
  onNavigate,
  t,
}: TocContentProps) => {
  const { behaviorCount, domainCount, featureCount } = calculateDocumentStats(document);

  // Use filtered domains when available, or convert to FilteredDomain for consistent types
  const displayDomains: FilteredDomain[] =
    filteredDocument?.domains ?? document.domains.map(toFilteredDomain);

  return (
    <div className="flex flex-col h-full">
      <div className="px-3 py-2 border-b">
        <div className="text-sm font-medium">{t("toc.title")}</div>
        <div className="text-xs text-muted-foreground mt-0.5">
          {hasFilter ? (
            <>
              {t("toc.filterMatchCount", {
                matched: filteredDocument?.matchCount ?? 0,
                total: behaviorCount,
              })}
            </>
          ) : (
            <>
              {t("stats.domains", { count: domainCount })} ·{" "}
              {t("stats.features", { count: featureCount })} ·{" "}
              {t("stats.behaviors", { count: behaviorCount })}
            </>
          )}
        </div>
      </div>
      <ScrollArea className="flex-1">
        <div className="p-2 space-y-1">
          {displayDomains.map((domain) => (
            <TocItem
              activeId={activeId}
              domain={domain}
              expandedDomainId={expandedDomainId}
              hasFilter={hasFilter}
              key={domain.id}
              onExpandToggle={onExpandToggle}
              onNavigate={onNavigate}
            />
          ))}
        </div>
      </ScrollArea>
    </div>
  );
};

/**
 * Extract domain ID from section ID
 * - "domain-123" → "domain-123"
 * - "feature-456" → null (need to look up parent domain)
 */
const extractDomainIdFromSectionId = (sectionId: string, domains: SpecDomain[]): string | null => {
  if (sectionId.startsWith("domain-")) {
    return sectionId;
  }
  if (sectionId.startsWith("feature-")) {
    const featureId = sectionId.replace("feature-", "");
    const parentDomain = domains.find((d) => d.features.some((f) => f.id === featureId));
    return parentDomain ? `domain-${parentDomain.id}` : null;
  }
  return null;
};

export const TocSidebar = ({
  document,
  filteredDocument,
  hasFilter = false,
  onNavigate,
}: TocSidebarProps) => {
  const t = useTranslations("specView");
  const [activeId, setActiveId] = useState<string | null>(null);
  const [expandedDomainId, setExpandedDomainId] = useState<string | null>(null);
  const observerRef = useRef<IntersectionObserver | null>(null);

  useEffect(() => {
    const handleHashChange = () => {
      const hash = window.location.hash.slice(1);
      if (hash) {
        setActiveId(hash);
      }
    };

    handleHashChange();
    window.addEventListener("hashchange", handleHashChange);
    return () => window.removeEventListener("hashchange", handleHashChange);
  }, []);

  // Auto-expand domain when active section changes
  useEffect(() => {
    if (activeId) {
      const domainId = extractDomainIdFromSectionId(activeId, document.domains);
      if (domainId) {
        setExpandedDomainId(domainId);
      }
    }
  }, [activeId, document.domains]);

  useEffect(() => {
    observerRef.current = new IntersectionObserver(
      (entries) => {
        const visibleEntries = entries
          .filter((entry) => entry.isIntersecting)
          .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top);

        if (visibleEntries.length > 0) {
          const topEntry = visibleEntries[0];
          if (topEntry && topEntry.target.id) {
            setActiveId(topEntry.target.id);
          }
        }
      },
      {
        rootMargin: "-80px 0px -60% 0px",
        threshold: 0,
      }
    );

    document.domains.forEach((domain) => {
      const domainEl = globalThis.document.getElementById(`domain-${domain.id}`);
      if (domainEl) {
        observerRef.current?.observe(domainEl);
      }

      domain.features.forEach((feature) => {
        const featureEl = globalThis.document.getElementById(`feature-${feature.id}`);
        if (featureEl) {
          observerRef.current?.observe(featureEl);
        }
      });
    });

    return () => {
      observerRef.current?.disconnect();
    };
  }, [document]);

  const handleNavigate = (sectionId: string) => {
    const element = globalThis.document.getElementById(sectionId);
    if (element) {
      element.scrollIntoView({ behavior: "smooth", block: "start" });
      window.history.replaceState(null, "", `#${sectionId}`);
      setActiveId(sectionId);
      onNavigate?.(sectionId);
    }
  };

  return (
    <>
      {/* Desktop: Fixed sidebar */}
      <div className="hidden lg:block w-64 flex-shrink-0">
        <div className="sticky top-20 h-[calc(100vh-6rem)] border rounded-lg overflow-hidden bg-card">
          <TocContent
            activeId={activeId}
            document={document}
            expandedDomainId={expandedDomainId}
            filteredDocument={filteredDocument}
            hasFilter={hasFilter}
            onExpandToggle={setExpandedDomainId}
            onNavigate={handleNavigate}
            t={t}
          />
        </div>
      </div>

      {/* Mobile: Sheet */}
      <div className="lg:hidden">
        <Sheet>
          <SheetTrigger asChild>
            <Button
              className="fixed bottom-20 right-4 z-[100] shadow-lg rounded-full h-12 w-12"
              size="icon"
            >
              <List className="h-5 w-5" />
              <span className="sr-only">{t("toc.openToc")}</span>
            </Button>
          </SheetTrigger>
          <SheetContent className="w-80 p-0" side="left">
            <SheetHeader className="sr-only">
              <SheetTitle>{t("toc.title")}</SheetTitle>
              <SheetDescription>{t("toc.description")}</SheetDescription>
            </SheetHeader>
            <TocContent
              activeId={activeId}
              document={document}
              expandedDomainId={expandedDomainId}
              filteredDocument={filteredDocument}
              hasFilter={hasFilter}
              onExpandToggle={setExpandedDomainId}
              onNavigate={handleNavigate}
              t={t}
            />
          </SheetContent>
        </Sheet>
      </div>
    </>
  );
};
