import type { FilteredDocument } from "../hooks/use-document-filter";
import type {
  FlatSpecBehaviorItem,
  FlatSpecDomainItem,
  FlatSpecFeatureItem,
  FlatSpecItem,
} from "../types";

/**
 * Flatten hierarchical FilteredDocument into flat array for window-level virtualization
 * Each item includes isLastInDomain flag to enable CSS-based card grouping
 */
export const flattenSpecDocument = (
  document: FilteredDocument,
  expandedDomains: Set<string>,
  expandedFeatures: Set<string>,
  hasFilter: boolean
): FlatSpecItem[] => {
  const result: FlatSpecItem[] = [];

  for (const domain of document.domains) {
    const isDomainExpanded = expandedDomains.has(domain.id);

    // Collect all items for this domain first to determine last item
    const domainItems: FlatSpecItem[] = [];

    const domainItem: FlatSpecDomainItem = {
      depth: 0,
      domain,
      domainId: domain.id,
      isExpanded: isDomainExpanded,
      isLastInDomain: false, // will be updated if collapsed
      type: "domain-header",
    };
    domainItems.push(domainItem);

    if (isDomainExpanded) {
      for (const feature of domain.features) {
        const isFeatureExpanded = expandedFeatures.has(feature.id);

        const featureItem: FlatSpecFeatureItem = {
          depth: 1,
          domainId: domain.id,
          feature,
          featureId: feature.id,
          isExpanded: isFeatureExpanded,
          isLastInDomain: false,
          type: "feature-header",
        };
        domainItems.push(featureItem);

        if (isFeatureExpanded) {
          const visibleBehaviors = hasFilter
            ? feature.behaviors.filter((b) => b.hasMatch)
            : feature.behaviors;

          for (const behavior of visibleBehaviors) {
            const behaviorItem: FlatSpecBehaviorItem = {
              behavior,
              behaviorId: behavior.id,
              depth: 2,
              domainId: domain.id,
              featureId: feature.id,
              isLastInDomain: false,
              type: "behavior",
            };
            domainItems.push(behaviorItem);
          }
        }
      }
    }

    // Mark the last item in this domain
    if (domainItems.length > 0) {
      const lastItem = domainItems[domainItems.length - 1];
      if (lastItem) {
        lastItem.isLastInDomain = true;
      }
    }

    result.push(...domainItems);
  }

  return result;
};

// Estimated heights per item type (pixels)
export const DOMAIN_HEADER_HEIGHT = 80;
export const FEATURE_HEADER_HEIGHT = 56;
export const BEHAVIOR_ITEM_HEIGHT = 72;

// Gap between domain cards (applied as margin-bottom on last item in domain)
export const DOMAIN_GAP = 24; // space-y-6 equivalent

/**
 * Returns estimated height for a flat spec item
 * Domain gap is only added for the last item in a domain
 */
export const getItemHeight = (item: FlatSpecItem): number => {
  const baseHeight =
    item.type === "domain-header"
      ? DOMAIN_HEADER_HEIGHT
      : item.type === "feature-header"
        ? FEATURE_HEADER_HEIGHT
        : BEHAVIOR_ITEM_HEIGHT;

  // Add domain gap only for last item in domain
  return baseHeight + (item.isLastInDomain ? DOMAIN_GAP : 0);
};
