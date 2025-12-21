"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";

import { ErrorFallback } from "@/components/feedback";
import { Button } from "@/components/ui/button";

type ErrorProps = {
  error: Error & { digest?: string };
  reset: () => void;
};

const Error = ({ error, reset }: ErrorProps) => {
  const t = useTranslations("errors");
  const tCommon = useTranslations("common");

  return (
    <ErrorFallback
      action={
        <div className="flex flex-col gap-3 sm:flex-row">
          <Button onClick={reset} variant="default">
            {tCommon("tryAgain")}
          </Button>
          <Button asChild variant="outline">
            <Link href="/">{tCommon("goHome")}</Link>
          </Button>
        </div>
      }
      description={error.message}
      title={t("somethingWentWrong")}
    />
  );
};

export default Error;
