"use client";

import { AlertTriangle, RefreshCw } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";

type SpecViewErrorProps = {
  error: Error;
  onRetry: () => void;
};

export const SpecViewError = ({ error, onRetry }: SpecViewErrorProps) => {
  const t = useTranslations("analyze.specView");

  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      <AlertTriangle className="h-12 w-12 text-destructive mb-4" />
      <h3 className="text-lg font-semibold mb-2">{t("errorTitle")}</h3>
      <p className="text-sm text-muted-foreground mb-4 max-w-md">{error.message}</p>
      <Button onClick={onRetry} variant="outline">
        <RefreshCw className="mr-2 h-4 w-4" />
        {t("retry")}
      </Button>
    </div>
  );
};
