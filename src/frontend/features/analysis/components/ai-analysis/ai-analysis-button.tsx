"use client";

import { Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

import { AIComingSoonModal } from "./ai-coming-soon-modal";

type AIAnalysisButtonProps = {
  className?: string;
};

export const AIAnalysisButton = ({ className }: AIAnalysisButtonProps = {}) => {
  const t = useTranslations("analyze.aiAnalysis");
  const [isModalOpen, setIsModalOpen] = useState(false);

  return (
    <>
      <Button
        className={cn(
          "relative overflow-hidden",
          "bg-gradient-to-r from-violet-600 to-indigo-600",
          "dark:from-violet-500 dark:to-indigo-500",
          "text-white border-0",
          "hover:from-violet-500 hover:to-indigo-500",
          "dark:hover:from-violet-400 dark:hover:to-indigo-400",
          "shadow-sm hover:shadow-md hover:shadow-violet-500/20",
          "transition-all duration-200",
          className
        )}
        onClick={() => setIsModalOpen(true)}
        size="sm"
      >
        <Sparkles className="h-4 w-4" />
        <span>{t("button")}</span>
        <span
          className={cn(
            "absolute -top-1 -right-1",
            "px-1.5 py-0.5 text-[10px] font-medium",
            "bg-amber-500 text-white rounded-full",
            "shadow-sm"
          )}
        >
          {t("badge")}
        </span>
      </Button>

      <AIComingSoonModal onClose={() => setIsModalOpen(false)} open={isModalOpen} />
    </>
  );
};
