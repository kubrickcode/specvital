"use client";

import { Lock } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useLoginModal } from "@/features/auth";

type LoginRequiredStateProps = {
  descriptionKey: string;
  titleKey: string;
};

export const LoginRequiredState = ({ descriptionKey, titleKey }: LoginRequiredStateProps) => {
  const t = useTranslations("explore.loginRequired");
  const { open } = useLoginModal();

  return (
    <Card className="flex flex-col items-center justify-center py-16 px-6 text-center">
      <div className="rounded-full bg-muted p-4 mb-6">
        <Lock aria-hidden="true" className="size-8 text-muted-foreground" />
      </div>

      <h2 className="text-xl font-semibold mb-2">{t(titleKey)}</h2>
      <p className="text-muted-foreground mb-6 max-w-md">{t(descriptionKey)}</p>

      <Button onClick={open} variant="cta">
        {t("signIn")}
      </Button>
    </Card>
  );
};
