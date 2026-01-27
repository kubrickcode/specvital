"use client";

import {
  AlertTriangle,
  ArrowRight,
  FileText,
  HelpCircle,
  Layers,
  Lightbulb,
  Sparkles,
  Zap,
} from "lucide-react";
import Link from "next/link";
import { useLocale, useTranslations } from "next-intl";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const LIMITATIONS = ["namesBased", "nameQuality", "languageSupport"] as const;

export const SpecviewGenerationContent = () => {
  const t = useTranslations("docs.specviewGeneration");
  const locale = useLocale();

  return (
    <div className="space-y-8">
      {/* Overview */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.overview.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.overview.description")}</p>

        <Alert className="border-purple-200 bg-purple-50 dark:border-purple-900 dark:bg-purple-950/50">
          <Sparkles className="size-4 text-purple-600 dark:text-purple-400" />
          <AlertTitle className="text-purple-900 dark:text-purple-100">
            {t("sections.overview.keyPoint.title")}
          </AlertTitle>
          <AlertDescription className="text-purple-800 dark:text-purple-200">
            {t("sections.overview.keyPoint.description")}
          </AlertDescription>
        </Alert>
      </section>

      {/* What You Get */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">{t("sections.output.title")}</h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.output.description")}</p>

        <div className="grid gap-3 sm:grid-cols-2">
          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Layers className="size-4 text-primary" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.output.features.domains.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.output.features.domains.description")}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <FileText className="size-4 text-primary" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.output.features.descriptions.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.output.features.descriptions.description")}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Sparkles className="size-4 text-primary" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.output.features.summary.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.output.features.summary.description")}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Zap className="size-4 text-primary" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.output.features.caching.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.output.features.caching.description")}
              </p>
            </CardContent>
          </Card>
        </div>
      </section>

      {/* Things to Keep in Mind */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.limitations.title")}
        </h2>

        <Card>
          <CardContent className="pt-6">
            <ul className="space-y-3 text-sm text-muted-foreground">
              {LIMITATIONS.map((limitKey) => (
                <li className="flex items-start gap-3" key={limitKey}>
                  <AlertTriangle className="mt-0.5 size-5 shrink-0 text-amber-500" />
                  <div>
                    <p className="font-medium text-foreground">
                      {t(`sections.limitations.items.${limitKey}.title`)}
                    </p>
                    <p>{t(`sections.limitations.items.${limitKey}.description`)}</p>
                  </div>
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>
      </section>

      {/* Related Guide */}
      <section>
        <Link className="block" href={`/${locale}/docs/writing-guide`}>
          <Card className="transition-colors hover:bg-muted/50">
            <CardContent className="flex items-center gap-3 py-4">
              <Lightbulb className="size-5 shrink-0 text-amber-500" />
              <div className="flex-1">
                <p className="font-medium">{t("sections.relatedGuide.title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.relatedGuide.description")}
                </p>
              </div>
              <ArrowRight className="size-4 shrink-0 text-muted-foreground" />
            </CardContent>
          </Card>
        </Link>
      </section>

      {/* FAQ */}
      <section>
        <h2 className="mb-6 text-2xl font-semibold tracking-tight">{t("sections.faq.title")}</h2>

        <div className="space-y-4">
          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-start gap-2">
                <HelpCircle className="mt-0.5 size-4 shrink-0 text-muted-foreground" />
                <CardTitle className="text-base">{t("sections.faq.q1.question")}</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {t("sections.faq.q1.answer")}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-start gap-2">
                <HelpCircle className="mt-0.5 size-4 shrink-0 text-muted-foreground" />
                <CardTitle className="text-base">{t("sections.faq.q2.question")}</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {t("sections.faq.q2.answer")}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-start gap-2">
                <HelpCircle className="mt-0.5 size-4 shrink-0 text-muted-foreground" />
                <CardTitle className="text-base">{t("sections.faq.q3.question")}</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {t("sections.faq.q3.answer")}
              </p>
            </CardContent>
          </Card>
        </div>
      </section>
    </div>
  );
};
