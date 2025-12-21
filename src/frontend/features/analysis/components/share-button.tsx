"use client";

import { Check, Link2, Share2 } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";

import { useShare } from "../hooks/use-share";

type ShareButtonProps = {
  showLabel?: boolean;
  size?: "default" | "sm" | "lg" | "icon";
  variant?: "default" | "outline" | "ghost";
};

export const ShareButton = ({
  showLabel = true,
  size = "sm",
  variant = "outline",
}: ShareButtonProps) => {
  const t = useTranslations("share");
  const { copyToClipboard, isNativeShareSupported, shareNative, state } = useShare();

  const handleClick = isNativeShareSupported ? shareNative : copyToClipboard;

  const getIcon = () => {
    if (state === "success") {
      return <Check className="h-4 w-4" />;
    }
    return isNativeShareSupported ? <Share2 className="h-4 w-4" /> : <Link2 className="h-4 w-4" />;
  };

  return (
    <Button
      aria-label={t("ariaLabel")}
      className="gap-2"
      onClick={handleClick}
      size={size}
      variant={variant}
    >
      {getIcon()}
      {showLabel && (
        <span className="hidden sm:inline">
          {state === "success" ? t("linkCopied") : t("button")}
        </span>
      )}
    </Button>
  );
};
