"use client";

import { useTranslations } from "next-intl";

export const PricingHeader = () => {
  const t = useTranslations("pricing.header");

  return (
    <div className="mb-12 text-center">
      <h1 className="text-3xl font-bold tracking-tight sm:text-4xl">{t("title")}</h1>
      <p className="mt-4 text-lg text-muted-foreground">{t("subtitle")}</p>
    </div>
  );
};
