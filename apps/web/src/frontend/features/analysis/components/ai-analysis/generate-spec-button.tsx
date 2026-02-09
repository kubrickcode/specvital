"use client";

import { Loader2, Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type GenerateSpecButtonProps = {
  className?: string;
  isActive: boolean;
  isDocumentAvailable: boolean;
  isGenerating: boolean;
  onClick: () => void;
};

export const GenerateSpecButton = ({
  className,
  isActive,
  isDocumentAvailable,
  isGenerating,
  onClick,
}: GenerateSpecButtonProps) => {
  const t = useTranslations("analyze.generateSpec");

  const buttonText = isGenerating
    ? t("generating")
    : isDocumentAvailable
      ? t("aiSpecDocument")
      : t("generate");

  return (
    <Button
      aria-label={t("ariaLabel")}
      aria-pressed={isActive}
      className={cn(
        "relative overflow-hidden",
        "bg-gradient-to-r from-violet-600 to-indigo-600",
        "dark:from-violet-500 dark:to-indigo-500",
        "text-white border-0",
        "hover:from-violet-500 hover:to-indigo-500",
        "dark:hover:from-violet-400 dark:hover:to-indigo-400",
        "shadow-sm hover:shadow-md hover:shadow-violet-500/20",
        "transition-all duration-200",
        // Active state when viewing document - subtle brightness increase
        isActive && [
          "from-violet-500 to-indigo-500",
          "dark:from-violet-400 dark:to-indigo-400",
          "shadow-lg shadow-violet-500/30",
        ],
        className
      )}
      disabled={isGenerating}
      onClick={onClick}
      size="default"
    >
      {isGenerating ? (
        <Loader2 className="h-4 w-4 animate-spin" />
      ) : (
        <Sparkles className="h-4 w-4" />
      )}
      <span>{buttonText}</span>
    </Button>
  );
};
