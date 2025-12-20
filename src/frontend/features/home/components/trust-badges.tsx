import { Gift, Layers, Zap } from "lucide-react";
import { getTranslations } from "next-intl/server";

export const TrustBadges = async () => {
  const t = await getTranslations("home.trustBadges");

  const badges = [
    {
      icon: Zap,
      label: t("zeroSetup"),
    },
    {
      icon: Layers,
      label: t("multiFramework"),
    },
    {
      icon: Gift,
      label: t("free"),
    },
  ];

  return (
    <div className="flex flex-wrap items-center justify-center gap-x-6 gap-y-2">
      {badges.map((badge) => {
        const Icon = badge.icon;
        return (
          <div key={badge.label} className="flex items-center gap-2 text-muted-foreground">
            <Icon className="h-4 w-4" aria-hidden="true" />
            <span className="text-sm">{badge.label}</span>
          </div>
        );
      })}
    </div>
  );
};
