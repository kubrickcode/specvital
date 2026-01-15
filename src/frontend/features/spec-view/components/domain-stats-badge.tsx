import { useTranslations } from "next-intl";

import type { SpecDocument, SpecDomain } from "../types";
import { calculateDocumentStats, calculateDomainStats } from "../utils/stats";

type DomainStatsBadgeProps =
  | { document: SpecDocument; domain?: never }
  | { document?: never; domain: SpecDomain };

export const DomainStatsBadge = ({ document, domain }: DomainStatsBadgeProps) => {
  const t = useTranslations("specView");

  if (document) {
    const { behaviorCount, domainCount, featureCount } = calculateDocumentStats(document);

    return (
      <span className="text-muted-foreground">
        {t("statsBadge.full", { behaviorCount, domainCount, featureCount })}
      </span>
    );
  }

  if (domain) {
    const { behaviorCount, featureCount } = calculateDomainStats(domain);

    return (
      <span className="text-muted-foreground">
        {t("statsBadge.domainOnly", { behaviorCount, featureCount })}
      </span>
    );
  }

  return null;
};
