"use client";

import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useAuth, useLoginModal } from "@/features/auth";
import { Link } from "@/i18n/navigation";

import type { PlanTier, PricingPlan } from "../types";

type PricingCtaProps = {
  currentTier?: PlanTier | null;
  plan: PricingPlan;
};

export const PricingCta = ({ currentTier, plan }: PricingCtaProps) => {
  const t = useTranslations("pricing.cta");
  const { isAuthenticated } = useAuth();
  const { open: openLogin } = useLoginModal();

  const isCurrentPlan = currentTier === plan.tier;
  const isPaidPlan = plan.monthlyPrice !== null && plan.monthlyPrice > 0;
  const isEnterprise = plan.tier === "enterprise";

  if (!isAuthenticated) {
    if (isEnterprise) {
      return (
        <Button asChild className="w-full" variant="outline">
          <Link href="mailto:support@specvital.com">{t("contactUs")}</Link>
        </Button>
      );
    }
    return (
      <Button className="w-full" onClick={openLogin} variant={plan.highlighted ? "cta" : "default"}>
        {t("getStartedFree")}
      </Button>
    );
  }

  if (isCurrentPlan) {
    return (
      <Button className="w-full" disabled variant="secondary">
        {t("currentPlan")}
      </Button>
    );
  }

  if (isEnterprise) {
    return (
      <Button asChild className="w-full" variant="outline">
        <Link href="mailto:support@specvital.com">{t("contactUs")}</Link>
      </Button>
    );
  }

  if (isPaidPlan) {
    return (
      <Tooltip>
        <TooltipTrigger asChild>
          <span className="w-full">
            <Button className="w-full" disabled variant={plan.highlighted ? "cta" : "default"}>
              {t("upgrade")}
            </Button>
          </span>
        </TooltipTrigger>
        <TooltipContent>
          <p>{t("upgradeComingSoon")}</p>
        </TooltipContent>
      </Tooltip>
    );
  }

  // Free plan but not current (e.g., Pro user viewing Free card)
  // No downgrade functionality, so hide button
  return null;
};
