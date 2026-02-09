import type { SpecDocument, SpecDomain } from "../types";

type DocumentStats = {
  behaviorCount: number;
  domainCount: number;
  featureCount: number;
};

type DomainStats = {
  behaviorCount: number;
  featureCount: number;
};

export const calculateDocumentStats = (document: SpecDocument): DocumentStats => ({
  behaviorCount: document.domains.reduce(
    (sum, d) => sum + d.features.reduce((fSum, f) => fSum + f.behaviors.length, 0),
    0
  ),
  domainCount: document.domains.length,
  featureCount: document.domains.reduce((sum, d) => sum + d.features.length, 0),
});

export const calculateDomainStats = (domain: SpecDomain): DomainStats => ({
  behaviorCount: domain.features.reduce((sum, f) => sum + f.behaviors.length, 0),
  featureCount: domain.features.length,
});
