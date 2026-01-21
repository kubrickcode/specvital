"use client";

import { FileSearch } from "lucide-react";
import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useRouter } from "@/i18n/navigation";
import {
  FRAMEWORK_CATEGORIES,
  HIGHLIGHTED_FRAMEWORKS,
  TOTAL_FRAMEWORK_COUNT,
} from "@/lib/constants/frameworks";

export const EmptyState = () => {
  const router = useRouter();
  const t = useTranslations("emptyState");

  const handleAnalyzeAnother = () => {
    router.push("/");
  };

  const remainingCount = TOTAL_FRAMEWORK_COUNT - HIGHLIGHTED_FRAMEWORKS.length;

  return (
    <div className="rounded-lg border bg-card p-8 shadow-xs md:p-12">
      <div className="flex flex-col items-center space-y-6 text-center">
        <div className="rounded-full bg-muted p-4">
          <FileSearch aria-hidden="true" className="h-8 w-8 text-muted-foreground" />
        </div>

        <div className="space-y-2">
          <h3 className="text-xl font-semibold text-foreground">{t("title")}</h3>
          <p className="max-w-md text-muted-foreground">{t("description")}</p>
        </div>

        <div className="w-full max-w-lg">
          <h4 className="mb-4 text-sm font-medium text-muted-foreground">
            {t("supportedFrameworks")}
          </h4>

          <div className="flex flex-wrap justify-center gap-2">
            {HIGHLIGHTED_FRAMEWORKS.map((framework) => (
              <Badge className="bg-black/5 text-foreground dark:bg-white/10" key={framework}>
                {framework}
              </Badge>
            ))}
            <Badge className="bg-primary/10 text-primary" variant="outline">
              +{remainingCount} more
            </Badge>
          </div>

          <div className="mt-4 flex flex-wrap justify-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
            {FRAMEWORK_CATEGORIES.map((category) => (
              <span key={category.category}>{category.frameworks.join(", ")}</span>
            )).reduce<React.ReactNode[]>((acc, curr, idx) => {
              if (idx === 0) return [curr];
              return [
                ...acc,
                <span aria-hidden="true" key={`sep-${idx}`}>
                  ·
                </span>,
                curr,
              ];
            }, [])}
          </div>
        </div>

        <Button className="mt-4" onClick={handleAnalyzeAnother} size="lg" variant="cta">
          {t("analyzeAnother")}
        </Button>
      </div>
    </div>
  );
};
