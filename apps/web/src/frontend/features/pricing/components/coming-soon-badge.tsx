"use client";

import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";

export const ComingSoonBadge = () => {
  const t = useTranslations("pricing");

  return (
    <Badge className="text-xs" variant="outline">
      {t("comingSoon")}
    </Badge>
  );
};
