"use client";

import { Plus, Search } from "lucide-react";
import { useTranslations } from "next-intl";
import type { ReactNode } from "react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { UrlInputForm } from "./url-input-form";

type AnalyzeDialogProps = {
  children?: ReactNode;
  variant?: "header" | "empty-state";
};

export const AnalyzeDialog = ({ children, variant = "header" }: AnalyzeDialogProps) => {
  const t = useTranslations("analyzeDialog");
  const [isOpen, setIsOpen] = useState(false);

  const defaultTrigger =
    variant === "header" ? (
      <Button size="sm" variant="ghost">
        <Plus className="h-4 w-4" />
        <span className="hidden sm:inline">{t("trigger")}</span>
      </Button>
    ) : (
      <Button size="lg" variant="cta">
        <Search aria-hidden="true" className="size-4 mr-2" />
        {t("ctaTrigger")}
      </Button>
    );

  return (
    <Dialog onOpenChange={setIsOpen} open={isOpen}>
      <DialogTrigger asChild>{children ?? defaultTrigger}</DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>
        <UrlInputForm />
      </DialogContent>
    </Dialog>
  );
};
