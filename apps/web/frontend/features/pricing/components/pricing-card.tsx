"use client";

import { Check } from "lucide-react";
import { useTranslations } from "next-intl";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { PricingCta } from "./pricing-cta";
import type { PlanTier, PricingPlan } from "../types";

type PricingCardProps = {
  currentTier?: PlanTier | null;
  plan: PricingPlan;
};

const formatFeatureValue = (value: string | number | null, unlimitedLabel: string): string => {
  if (value === null) return unlimitedLabel;
  if (typeof value === "number") return value.toLocaleString();
  return value;
};

export const PricingCard = ({ currentTier, plan }: PricingCardProps) => {
  const t = useTranslations("pricing");
  const tTiers = useTranslations("pricing.tiers");
  const tFeatures = useTranslations("pricing.features");

  return (
    <Card className="relative flex flex-col transition-all duration-200">
      <CardHeader className="text-center">
        <CardTitle className="text-xl">{tTiers(`${plan.tier}.name`)}</CardTitle>
        <CardDescription className="mt-1 min-h-[40px]">
          {tTiers(`${plan.tier}.description`)}
        </CardDescription>
      </CardHeader>

      <CardContent className="flex-1">
        <div className="mb-6 text-center">
          {plan.monthlyPrice === null ? (
            <div className="text-3xl font-bold">{t("customPricing")}</div>
          ) : (
            <div className="flex items-baseline justify-center gap-1">
              <span className="text-4xl font-bold">${plan.monthlyPrice}</span>
              <span className="text-muted-foreground">/{t("month")}</span>
            </div>
          )}
        </div>

        <ul className="space-y-3">
          {plan.features.map((feature) => (
            <li className="flex items-center gap-2" key={feature.label}>
              <Check className="size-4 shrink-0 text-green-500" />
              <span className="text-sm">
                <span className="text-muted-foreground">
                  {tFeatures(`${feature.label}.label`)}:{" "}
                </span>
                <span className="font-medium">
                  {formatFeatureValue(feature.value, t("unlimited"))}
                  {feature.label !== "retention" &&
                    feature.value !== null &&
                    ` ${tFeatures(`${feature.label}.unit`)}/${t("month")}`}
                  {feature.label === "retention" && feature.value !== null && ` ${t("days")}`}
                </span>
              </span>
            </li>
          ))}
        </ul>
      </CardContent>

      <CardFooter>
        <PricingCta currentTier={currentTier} plan={plan} />
      </CardFooter>
    </Card>
  );
};
