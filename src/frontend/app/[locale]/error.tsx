"use client";

import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";

type ErrorProps = {
  error: Error & { digest?: string };
  reset: () => void;
};

export const ErrorBoundary = ({ error, reset }: ErrorProps) => {
  const t = useTranslations("errors");

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8">
      <div className="w-full max-w-md space-y-6 text-center">
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

export default ErrorBoundary;
