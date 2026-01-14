import type {
  FilteredBehavior,
  FilteredDomain,
  FilteredFeature,
} from "../hooks/use-document-filter";

export type FlatSpecItemType = "domain-header" | "feature-header" | "behavior";

export type FlatSpecDomainItem = {
  depth: 0;
  domain: FilteredDomain;
  domainId: string;
  isExpanded: boolean;
  isLastInDomain: boolean; // true when domain is collapsed (header is the last item)
  type: "domain-header";
};

export type FlatSpecFeatureItem = {
  depth: 1;
  domainId: string;
  feature: FilteredFeature;
  featureId: string;
  isExpanded: boolean;
  isLastInDomain: boolean;
  type: "feature-header";
};

export type FlatSpecBehaviorItem = {
  behavior: FilteredBehavior;
  behaviorId: string;
  depth: 2;
  domainId: string;
  featureId: string;
  isLastInDomain: boolean;
  type: "behavior";
};

export type FlatSpecItem = FlatSpecDomainItem | FlatSpecFeatureItem | FlatSpecBehaviorItem;
