"use client";

import { FileSearch } from "lucide-react";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { useRouter } from "@/i18n/navigation";

const SUPPORTED_FRAMEWORKS = [
  { name: "Jest", patterns: ["*.test.ts", "*.test.tsx", "*.test.js", "__tests__/*"] },
  { name: "Vitest", patterns: ["*.test.ts", "*.test.tsx", "*.spec.ts"] },
  { name: "Playwright", patterns: ["*.spec.ts", "*.spec.tsx"] },
  { name: "Go", patterns: ["*_test.go"] },
];

export const EmptyState = () => {
  const router = useRouter();
  const t = useTranslations("emptyState");

  const handleAnalyzeAnother = () => {
    router.push("/");
  };

  return (
    <div className="rounded-lg border bg-card p-8 md:p-12 shadow-xs">
      <div className="flex flex-col items-center text-center space-y-6">
        <div className="rounded-full bg-muted p-4">
          <FileSearch className="h-8 w-8 text-muted-foreground" aria-hidden="true" />
        </div>

        <div className="space-y-2">
          <h3 className="text-xl font-semibold text-foreground">{t("title")}</h3>
          <p className="text-muted-foreground max-w-md">{t("description")}</p>
        </div>

        <div className="w-full max-w-md">
          <h4 className="text-sm font-medium text-muted-foreground mb-3">
            {t("supportedFrameworks")}
          </h4>
          <ul className="space-y-2 text-sm text-left">
            {SUPPORTED_FRAMEWORKS.map((framework) => (
              <li key={framework.name} className="flex items-start gap-2">
                <span className="font-medium text-foreground min-w-[80px]">
                  {framework.name}
                </span>
                <span className="text-muted-foreground">{framework.patterns.join(", ")}</span>
              </li>
            ))}
          </ul>
        </div>

        <Button onClick={handleAnalyzeAnother} size="lg" className="mt-4">
          {t("analyzeAnother")}
        </Button>
      </div>
    </div>
  );
};
