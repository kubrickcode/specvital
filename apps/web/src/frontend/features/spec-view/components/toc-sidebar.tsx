"use client";

import { ChevronDown, ChevronRight, List } from "lucide-react";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { useTruncateDetection } from "@/lib/hooks";
import { cn } from "@/lib/utils";

import { useDocumentNavigation } from "../contexts";
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
};

type TocItemProps = {
  activeId: string | null;
  domain: FilteredDomain;
  expandedDomainId: string | null;
  hasFilter?: boolean;
  onExpandToggle: (domainId: string) => void;
  onNavigate: (sectionId: string) => void;
};

const TocDomainName = ({ name }: { name: string }) => {
  const { isTruncated, ref } = useTruncateDetection<HTMLSpanElement>();

  const nameElement = (
    <span className="flex-1 min-w-0 lg:truncate" ref={ref}>
      {name}
    </span>
  );

  return isTruncated ? (
    <ResponsiveTooltip content={name} side="right">
      {nameElement}
    </ResponsiveTooltip>
  ) : (
    nameElement
  );
};

const TocFeatureName = ({ name }: { name: string }) => {
  const { isTruncated, ref } = useTruncateDetection<HTMLSpanElement>();

  const nameElement = (
    <span className="flex-1 min-w-0 lg:truncate" ref={ref}>
      {name}
    </span>
  );

  return isTruncated ? (
    <ResponsiveTooltip content={name} side="right">
      {nameElement}
    </ResponsiveTooltip>
  ) : (
    nameElement
  );
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
          "w-full max-w-full flex items-center gap-2 px-3 py-2 text-sm rounded-md transition-colors text-left overflow-hidden",
          isDomainActive
            ? "bg-gradient-to-r from-foreground/[0.10] to-transparent font-medium text-foreground"
            : "hover:bg-muted/50 text-muted-foreground"
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
        <TocDomainName name={domain.name} />
        <span className="text-xs text-muted-foreground flex-shrink-0">{displayCount}</span>
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
                  "w-full max-w-full flex items-center gap-2 px-3 py-1.5 text-xs rounded-md transition-colors text-left overflow-hidden",
                  isActive
                    ? "bg-gradient-to-r from-foreground/[0.10] to-transparent font-medium text-foreground"
                    : "hover:bg-muted/50 text-muted-foreground"
                )}
                key={feature.id}
                onClick={() => onNavigate(featureId)}
                type="button"
              >
                <TocFeatureName name={feature.name} />
                <span className="text-muted-foreground ml-auto flex-shrink-0">
                  {featureDisplayCount}
                </span>
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
    <div className="flex flex-col h-full overflow-hidden">
      <div className="px-3 py-2 border-b flex-shrink-0">
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
      <ScrollArea className="flex-1 min-h-0 [&>div>div]:!block [&>div>div]:!w-full">
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

export const TocSidebar = ({ document, filteredDocument, hasFilter = false }: TocSidebarProps) => {
  const t = useTranslations("specView");
  const { activeSection, scrollToSection } = useDocumentNavigation();

  // Track expanded domain in TOC (independent of main document view)
  const [tocExpandedDomainId, setTocExpandedDomainId] = useState<string | null>(null);

  // Auto-expand domain when active section changes
  useEffect(() => {
    if (activeSection) {
      const domainId = extractDomainIdFromSectionId(activeSection, document.domains);
      if (domainId) {
        setTocExpandedDomainId(domainId);
      }
    }
  }, [activeSection, document.domains]);

  const handleNavigate = (sectionId: string) => {
    scrollToSection(sectionId);
  };

  return (
    <>
      {/* Desktop: Fixed sidebar */}
      <div className="hidden lg:block w-64 flex-shrink-0">
        <div className="sticky top-20 h-[calc(100vh-6rem)] border rounded-lg overflow-hidden bg-card">
          <TocContent
            activeId={activeSection}
            document={document}
            expandedDomainId={tocExpandedDomainId}
            filteredDocument={filteredDocument}
            hasFilter={hasFilter}
            onExpandToggle={setTocExpandedDomainId}
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
              className="fixed bottom-20 right-4 z-40 shadow-lg rounded-full h-12 w-12"
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
              activeId={activeSection}
              document={document}
              expandedDomainId={tocExpandedDomainId}
              filteredDocument={filteredDocument}
              hasFilter={hasFilter}
              onExpandToggle={setTocExpandedDomainId}
              onNavigate={handleNavigate}
              t={t}
            />
          </SheetContent>
        </Sheet>
      </div>
    </>
  );
};
