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
  onOpenChange?: (open: boolean) => void;
  open?: boolean;
  variant?: "header" | "empty-state";
};

export const AnalyzeDialog = ({
  children,
  onOpenChange,
  open,
  variant = "header",
}: AnalyzeDialogProps) => {
  const t = useTranslations("analyzeDialog");
  const [internalOpen, setInternalOpen] = useState(false);

  const isControlled = open !== undefined;
  const isOpen = isControlled ? open : internalOpen;
  const setIsOpen = isControlled ? onOpenChange : setInternalOpen;

  const defaultTrigger =
    variant === "header" ? (
      <Button size="sm" variant="header-action">
        <Plus className="h-4 w-4" />
        <span className="hidden sm:inline">{t("trigger")}</span>
      </Button>
    ) : (
      <Button size="lg" variant="cta">
        <Search aria-hidden="true" className="size-4 mr-2" />
        {t("ctaTrigger")}
      </Button>
    );

  const handleOpenChange = (value: boolean) => {
    setIsOpen?.(value);
  };

  const handleSuccess = () => {
    setIsOpen?.(false);
  };

  return (
    <Dialog onOpenChange={handleOpenChange} open={isOpen}>
      {!isControlled && <DialogTrigger asChild>{children ?? defaultTrigger}</DialogTrigger>}
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>
        <UrlInputForm onSuccess={handleSuccess} />
      </DialogContent>
    </Dialog>
  );
};
