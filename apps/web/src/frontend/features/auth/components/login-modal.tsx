"use client";

import { Github, Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { useAuth, useLoginModal } from "../hooks";

export const LoginModal = () => {
  const t = useTranslations("loginModal");
  const { login, loginPending } = useAuth();
  const { isOpen, onOpenChange } = useLoginModal();

  return (
    <Dialog onOpenChange={onOpenChange} open={isOpen}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader className="text-center sm:text-center">
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>

        <div className="mt-4 flex flex-col gap-3">
          <Button
            className="w-full"
            disabled={loginPending}
            onClick={login}
            size="lg"
            variant="cta"
          >
            {loginPending ? (
              <Loader2 className="mr-2 size-5 animate-spin" />
            ) : (
              <Github className="mr-2 size-5" />
            )}
            {t("continueWithGitHub")}
          </Button>
        </div>

        <p className="text-muted-foreground mt-4 text-center text-xs">{t("terms")}</p>
      </DialogContent>
    </Dialog>
  );
};
