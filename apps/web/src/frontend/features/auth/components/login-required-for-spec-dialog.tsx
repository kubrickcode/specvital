"use client";

import { Check, Github, Loader2, Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { useAuth, useSpecLoginDialog } from "../hooks";

type BenefitItemProps = {
  description: string;
  title: string;
};

const BenefitItem = ({ description, title }: BenefitItemProps) => (
  <li className="flex items-start gap-3">
    <div className="bg-green-500/10 flex h-5 w-5 shrink-0 items-center justify-center rounded-full">
      <Check className="h-3 w-3 text-green-600" />
    </div>
    <div>
      <p className="text-sm font-medium">{title}</p>
      <p className="text-muted-foreground text-xs">{description}</p>
    </div>
  </li>
);

export const LoginRequiredForSpecDialog = () => {
  const t = useTranslations("specLoginRequired");
  const { login, loginPending } = useAuth();
  const { close, isOpen, onOpenChange } = useSpecLoginDialog();

  const handleLogin = () => {
    close();
    login();
  };

  return (
    <Dialog onOpenChange={onOpenChange} open={isOpen}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader className="text-center sm:text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gradient-to-br from-violet-500/20 to-indigo-500/20">
            <Sparkles className="h-6 w-6 text-violet-500" />
          </div>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>

        <div className="mt-2">
          <p className="text-muted-foreground mb-3 text-sm font-medium">{t("whySignIn")}</p>
          <ul className="space-y-3">
            <BenefitItem description={t("benefit1Description")} title={t("benefit1")} />
            <BenefitItem description={t("benefit2Description")} title={t("benefit2")} />
            <BenefitItem description={t("benefit3Description")} title={t("benefit3")} />
          </ul>
        </div>

        <div className="mt-4 flex flex-col gap-3">
          <Button
            className="w-full"
            disabled={loginPending}
            onClick={handleLogin}
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

        <p className="text-muted-foreground mt-2 text-center text-xs">{t("freeNote")}</p>
      </DialogContent>
    </Dialog>
  );
};
