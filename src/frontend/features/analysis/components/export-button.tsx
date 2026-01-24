"use client";

import { Check, ClipboardCopy, Download, FileText } from "lucide-react";
import type { Transition } from "motion/react";
import { AnimatePresence, motion } from "motion/react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import type { components } from "@/lib/api/generated-types";
import { useReducedMotion } from "@/lib/motion";

import { copyToClipboard, downloadMarkdown, exportToMarkdown } from "../utils/export-markdown";

type AnalysisResult = components["schemas"]["AnalysisResult"];

type ExportButtonProps = {
  data: AnalysisResult;
  size?: "default" | "sm" | "lg" | "icon";
  variant?: "default" | "outline" | "ghost";
};

type ExportState = "idle" | "copied" | "downloaded";

const iconTransition: Transition = {
  duration: 0.2,
  ease: "easeOut",
  type: "tween",
};

export const ExportButton = ({
  data,
  size = "default",
  variant = "outline",
}: ExportButtonProps) => {
  const t = useTranslations("export");
  const shouldReduceMotion = useReducedMotion();
  const [state, setState] = useState<ExportState>("idle");
  const isIconOnly = size === "icon";

  const handleDownload = () => {
    const markdown = exportToMarkdown(data);
    const filename = `${data.owner}-${data.repo}-tests.md`;
    downloadMarkdown(markdown, filename);
    setState("downloaded");
    setTimeout(() => setState("idle"), 2000);
  };

  const handleCopy = async () => {
    const markdown = exportToMarkdown(data);
    const success = await copyToClipboard(markdown);
    if (success) {
      setState("copied");
      setTimeout(() => setState("idle"), 2000);
    }
  };

  const getIcon = () => {
    if (state === "copied" || state === "downloaded") {
      return Check;
    }
    return FileText;
  };

  const IconComponent = getIcon();
  const iconKey = state === "idle" ? "file" : "check";

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          aria-label={t("ariaLabel")}
          className={isIconOnly ? "" : "gap-2"}
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
          {!isIconOnly && <span>{state !== "idle" ? t("success") : t("button")}</span>}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={handleDownload}>
          <Download className="mr-2 h-4 w-4" />
          {t("downloadMarkdown")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={handleCopy}>
          <ClipboardCopy className="mr-2 h-4 w-4" />
          {t("copyMarkdown")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
