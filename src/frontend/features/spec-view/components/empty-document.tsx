"use client";

import { FileText, Sparkles } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

type EmptyDocumentProps = {
  isLoading?: boolean;
  onGenerate: () => void;
};

export const EmptyDocument = ({ isLoading = false, onGenerate }: EmptyDocumentProps) => {
  const t = useTranslations("specView.emptyDocument");

  return (
    <Card className="border-dashed">
      <CardHeader className="text-center pb-2">
        <div className="mx-auto mb-4 rounded-full bg-primary/10 p-4">
          <FileText className="h-8 w-8 text-primary" />
        </div>
        <CardTitle>{t("title")}</CardTitle>
        <CardDescription className="max-w-sm mx-auto">{t("description")}</CardDescription>
      </CardHeader>
      <CardContent className="flex justify-center pt-2">
        <Button disabled={isLoading} onClick={onGenerate} size="lg">
          <Sparkles className="h-4 w-4 mr-2" />
          {isLoading ? t("generating") : t("generateButton")}
        </Button>
      </CardContent>
    </Card>
  );
};
