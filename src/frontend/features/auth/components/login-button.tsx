"use client";

import { Github, Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";

import { useAuth } from "../hooks";

export const LoginButton = () => {
  const t = useTranslations("auth");
  const { login, loginPending } = useAuth();

  return (
    <Button disabled={loginPending} onClick={login} size="sm" variant="outline">
      {loginPending ? (
        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
      ) : (
        <Github className="mr-2 h-4 w-4" />
      )}
      {t("login")}
    </Button>
  );
};
