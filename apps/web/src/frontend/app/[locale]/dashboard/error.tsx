"use client";

import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";

type DashboardErrorProps = {
  error: Error & { digest?: string };
  reset: () => void;
};

const DashboardError = ({ error, reset }: DashboardErrorProps) => {
  const t = useTranslations("errors");

  return (
    <main className="container mx-auto px-4 py-8">
      <div className="flex flex-col items-center justify-center space-y-6 py-12 text-center">
        <div className="space-y-2">
          <h1 className="text-2xl font-bold">{t("somethingWentWrong")}</h1>
          <p className="text-muted-foreground">{error.message || t("unexpectedError")}</p>
        </div>

        <Button onClick={reset} variant="outline">
          {t("tryAgain")}
        </Button>
      </div>
    </main>
  );
};

export default DashboardError;
