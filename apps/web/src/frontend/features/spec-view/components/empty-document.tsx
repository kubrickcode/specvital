"use client";

import { FileText, Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

import { QuotaIndicator } from "./quota-indicator";
import { isQuotaExceeded } from "../utils/quota";

type QuotaInfo = {
  limit: number | null;
  percentage: number | null;
  reserved: number;
  used: number;
};

type EmptyDocumentProps = {
  isLoading?: boolean;
  onGenerate: () => void;
  quota?: QuotaInfo | null;
};

export const EmptyDocument = ({ isLoading = false, onGenerate, quota }: EmptyDocumentProps) => {
  const t = useTranslations("specView.emptyDocument");

  const isDisabled = isLoading || Boolean(quota && isQuotaExceeded(quota.percentage));

  return (
    <Card className="border-dashed">
      <CardHeader className="text-center pb-2">
        <div className="mx-auto mb-4 rounded-full bg-primary/10 p-4">
          <FileText className="h-8 w-8 text-primary" />
        </div>
        <CardTitle>{t("title")}</CardTitle>
        <CardDescription className="max-w-sm mx-auto">{t("description")}</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col items-center gap-4 pt-2">
        {quota && (
          <div className="w-full max-w-sm">
            <QuotaIndicator
              limit={quota.limit}
              percentage={quota.percentage}
              reserved={quota.reserved}
              used={quota.used}
            />
          </div>
        )}
        <Button disabled={isDisabled} onClick={onGenerate} size="lg">
          <Sparkles className="h-4 w-4 mr-2" />
          {isLoading ? t("generating") : t("generateButton")}
        </Button>
      </CardContent>
    </Card>
  );
};
