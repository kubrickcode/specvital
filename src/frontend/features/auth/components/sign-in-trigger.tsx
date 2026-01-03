"use client";

import { LogIn } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";

import { useLoginModal } from "../hooks";

export const SignInTrigger = () => {
  const t = useTranslations("auth");
  const { open } = useLoginModal();

  return (
    <Button onClick={open} size="lg" variant="cta">
      <LogIn className="size-4 shrink-0 sm:mr-1.5" />
      <span className="hidden sm:inline">{t("login")}</span>
    </Button>
  );
};
