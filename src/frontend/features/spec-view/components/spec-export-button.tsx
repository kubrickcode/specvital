"use client";

import { Check, ChevronDown, ClipboardCopy, Download, FileText } from "lucide-react";
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
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useReducedMotion } from "@/lib/motion";
import { cn } from "@/lib/utils";

import type { RepoSpecDocument, SpecDocument } from "../types";
import {
  copySpecToClipboard,
  downloadSpecMarkdown,
  exportSpecToMarkdown,
  generateSpecFilename,
} from "../utils/export-spec-markdown";

type ExportableDocument = SpecDocument | RepoSpecDocument;

type SpecExportButtonProps = {
  className?: string;
  disabled?: boolean;
  document: ExportableDocument;
  owner: string;
  repo: string;
  variant?: "default" | "ghost" | "outline";
};

type ExportState = "copied" | "downloaded" | "idle";

const iconTransition: Transition = {
  duration: 0.2,
  ease: "easeOut",
  type: "tween",
};

export const SpecExportButton = ({
  className,
  disabled = false,
  document,
  owner,
  repo,
  variant = "outline",
}: SpecExportButtonProps) => {
  const t = useTranslations("specView.export");
  const shouldReduceMotion = useReducedMotion();
  const [state, setState] = useState<ExportState>("idle");

  const handleDownload = () => {
    const markdown = exportSpecToMarkdown(document, { owner, repo });
    const filename = generateSpecFilename(owner, repo, document.language);
    downloadSpecMarkdown(markdown, filename);
    setState("downloaded");
    setTimeout(() => setState("idle"), 2000);
  };

  const handleCopy = async () => {
    const markdown = exportSpecToMarkdown(document, { owner, repo });
    const success = await copySpecToClipboard(markdown);
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
      <Tooltip>
        <TooltipTrigger asChild>
          <DropdownMenuTrigger asChild>
            <Button
              aria-label={t("ariaLabel")}
              className={cn("h-auto px-2.5 py-1 gap-1.5 text-xs font-normal", className)}
              disabled={disabled}
              variant={variant}
            >
              <span className="relative flex h-3 w-3 items-center justify-center">
                <AnimatePresence initial={false} mode="wait">
                  <motion.span
                    animate={shouldReduceMotion ? {} : { opacity: 1, scale: 1 }}
                    className="absolute inset-0 flex items-center justify-center"
                    exit={shouldReduceMotion ? {} : { opacity: 0, scale: 0.8 }}
                    initial={shouldReduceMotion ? false : { opacity: 0, scale: 0.8 }}
                    key={iconKey}
                    transition={shouldReduceMotion ? undefined : iconTransition}
                  >
                    <IconComponent className="h-3 w-3" />
                  </motion.span>
                </AnimatePresence>
              </span>
              <span>{state !== "idle" ? t("success") : t("button")}</span>
              <ChevronDown className="h-3 w-3 opacity-50" />
            </Button>
          </DropdownMenuTrigger>
        </TooltipTrigger>
        <TooltipContent>{t("tooltip")}</TooltipContent>
      </Tooltip>
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
