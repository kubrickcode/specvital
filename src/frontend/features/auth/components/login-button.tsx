"use client";

import { Github, Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";

import { useAuth } from "../hooks";

export const LoginButton = () => {
  const t = useTranslations("auth");
  const { login, loginPending } = useAuth();

  return (
    <Button disabled={loginPending} onClick={login} size="lg" variant="cta">
      {loginPending ? (
        <Loader2 className="size-4 shrink-0 animate-spin text-primary-foreground/70 sm:mr-1.5" />
      ) : (
        <Github className="size-4 shrink-0 sm:mr-1.5" />
      )}
      <span className="hidden sm:inline">{t("login")}</span>
    </Button>
  );
};
