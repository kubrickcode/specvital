"use client";

import { Construction } from "lucide-react";
import { useTranslations } from "next-intl";

import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";

import { PulseRing } from "./pulse-ring";

type MaintenancePageProps = {
  className?: string;
  description?: string;
  title?: string;
};

export const MaintenancePage = ({ className, description, title }: MaintenancePageProps) => {
  const t = useTranslations("maintenance");

  const displayTitle = title ?? t("title");
  const displayDescription = description ?? t("description");

  return (
    <main
      aria-labelledby="maintenance-title"
      className={cn(
        "flex min-h-screen flex-col items-center justify-center p-6",
        "bg-background",
        className
      )}
      role="alert"
    >
      <Card className="w-full max-w-md border-border/50">
        <CardContent className="flex flex-col items-center gap-6 pt-8 text-center">
          {/* 상태 아이콘 */}
          <div className="relative">
            <div className="rounded-full bg-muted p-4">
              <Construction aria-hidden="true" className="size-8 text-muted-foreground" />
            </div>
            <span className="absolute -right-1 -top-1 text-status-skipped">
              <PulseRing size="md" />
            </span>
          </div>

          {/* 헤드라인 */}
          <div className="space-y-2">
            <h1 className="text-xl font-semibold tracking-tight" id="maintenance-title">
              {displayTitle}
            </h1>
            <p className="text-sm leading-relaxed text-muted-foreground">{displayDescription}</p>
          </div>

          {/* 브랜드 */}
          <p className="pt-4 text-xs text-muted-foreground/60">SpecVital</p>
        </CardContent>
      </Card>
    </main>
  );
};
