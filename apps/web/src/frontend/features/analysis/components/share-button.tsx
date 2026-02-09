"use client";

import { Check, Link2, Share2 } from "lucide-react";
import type { Transition } from "motion/react";
import { AnimatePresence, motion } from "motion/react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { useReducedMotion } from "@/lib/motion";

import { useShare } from "../hooks/use-share";

type ShareButtonProps = {
  showLabel?: boolean;
  size?: "default" | "sm" | "lg" | "icon";
  variant?: "default" | "outline" | "ghost";
};

const iconTransition: Transition = {
  duration: 0.2,
  ease: "easeOut",
  type: "tween",
};

export const ShareButton = ({
  showLabel = true,
  size = "sm",
  variant = "outline",
}: ShareButtonProps) => {
  const t = useTranslations("share");
  const { copyToClipboard, isNativeShareSupported, shareNative, state } = useShare();
  const shouldReduceMotion = useReducedMotion();

  const handleClick = isNativeShareSupported ? shareNative : copyToClipboard;

  const iconKey = state === "success" ? "check" : isNativeShareSupported ? "share" : "link";
  const IconComponent = state === "success" ? Check : isNativeShareSupported ? Share2 : Link2;

  return (
    <Button
      aria-label={t("ariaLabel")}
      className="gap-2"
      onClick={handleClick}
      size={size}
      variant={variant}
    >
      <span className="relative flex h-4 w-4 items-center justify-center">
        <AnimatePresence initial={false} mode="wait">
          <motion.span
            animate={shouldReduceMotion ? {} : { opacity: 1, scale: 1 }}
            className="absolute inset-0 flex items-center justify-center"
            exit={shouldReduceMotion ? {} : { opacity: 0, scale: 0.8 }}
            initial={shouldReduceMotion ? false : { opacity: 0, scale: 0.8 }}
            key={iconKey}
            transition={shouldReduceMotion ? undefined : iconTransition}
          >
            <IconComponent className="h-4 w-4" />
          </motion.span>
        </AnimatePresence>
      </span>
      {showLabel && <span>{state === "success" ? t("linkCopied") : t("button")}</span>}
    </Button>
  );
};
