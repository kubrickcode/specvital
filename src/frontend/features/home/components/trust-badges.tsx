"use client";

import { Gift, Layers, Zap } from "lucide-react";
import { motion } from "motion/react";
import { useTranslations } from "next-intl";

import { staggerContainer, staggerItem, useReducedMotion } from "@/lib/motion";
import { cn } from "@/lib/utils";

type TrustBadgesProps = {
  onFrameworkClick?: () => void;
};

export const TrustBadges = ({ onFrameworkClick }: TrustBadgesProps) => {
  const t = useTranslations("home.trustBadges");
  const shouldReduceMotion = useReducedMotion();

  const badges = [
    {
      icon: Zap,
      isClickable: false,
      key: "instant",
      label: t("instant"),
    },
    {
      icon: Layers,
      isClickable: true,
      key: "multiFramework",
      label: t("multiFramework"),
    },
    {
      icon: Gift,
      isClickable: false,
      key: "free",
      label: t("free"),
    },
  ];

  const containerVariants = shouldReduceMotion ? {} : staggerContainer;
  const itemVariants = shouldReduceMotion ? {} : staggerItem;

  return (
    <motion.div
      animate="visible"
      className="flex flex-wrap items-center justify-center gap-x-4 gap-y-2"
      initial={shouldReduceMotion ? false : "hidden"}
      variants={containerVariants}
    >
      {badges.map((badge) => {
        const Icon = badge.icon;
        const isClickable = badge.isClickable && onFrameworkClick;

        const badgeContent = (
          <>
            <Icon aria-hidden="true" className="size-4" />
            <span className="text-sm">{badge.label}</span>
          </>
        );

        const badgeClassName = cn(
          "flex items-center gap-2 rounded-full px-3 py-1.5 text-muted-foreground transition-all",
          "bg-gradient-to-b from-white/80 to-secondary/60 border border-border/40 shadow-sm",
          "dark:from-white/[0.08] dark:to-white/[0.03] dark:border-white/10",
          isClickable
            ? "cursor-pointer hover:border-border hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            : "hover:border-border hover:shadow-md"
        );

        if (isClickable) {
          return (
            <motion.button
              aria-label={badge.label}
              className={badgeClassName}
              key={badge.key}
              onClick={onFrameworkClick}
              type="button"
              variants={itemVariants}
            >
              {badgeContent}
            </motion.button>
          );
        }

        return (
          <motion.div className={badgeClassName} key={badge.key} variants={itemVariants}>
            {badgeContent}
          </motion.div>
        );
      })}
    </motion.div>
  );
};
