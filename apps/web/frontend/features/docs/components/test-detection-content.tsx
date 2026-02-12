"use client";

import { AlertCircle, CheckCircle2, Zap } from "lucide-react";
import { useTranslations } from "next-intl";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export const TestDetectionContent = () => {
  const t = useTranslations("docs.testDetection");

  return (
    <div className="space-y-8">
      {/* What is Test Detection */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">{t("sections.what.title")}</h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.what.description")}</p>

        <Alert className="border-blue-200 bg-blue-50 dark:border-blue-900 dark:bg-blue-950/50">
          <Zap className="size-4 text-blue-600 dark:text-blue-400" />
          <AlertTitle className="text-blue-900 dark:text-blue-100">
            {t("sections.what.keyPoint.title")}
          </AlertTitle>
          <AlertDescription className="text-blue-800 dark:text-blue-200">
            {t("sections.what.keyPoint.description")}
          </AlertDescription>
        </Alert>
      </section>

      {/* How It Works */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">{t("sections.how.title")}</h2>
        <p className="leading-7 text-muted-foreground">{t("sections.how.description")}</p>
      </section>

      {/* Test Status Detection */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">{t("sections.status.title")}</h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.status.description")}</p>

        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <Card className="border-l-4 border-l-green-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-green-700 dark:text-green-400">
                {t("sections.status.active.label")}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.status.active.description")}
              </p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-yellow-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-yellow-700 dark:text-yellow-400">
                {t("sections.status.skipped.label")}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.status.skipped.description")}
              </p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-blue-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-blue-700 dark:text-blue-400">
                {t("sections.status.todo.label")}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.status.todo.description")}
              </p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-purple-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-purple-700 dark:text-purple-400">
                {t("sections.status.focused.label")}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.status.focused.description")}
              </p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-orange-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-orange-700 dark:text-orange-400">
                {t("sections.status.xfail.label")}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.status.xfail.description")}
              </p>
            </CardContent>
          </Card>
        </div>
      </section>

      {/* Known Limitations */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.limitations.title")}
        </h2>

        <Alert className="border-orange-200 bg-orange-50 dark:border-orange-900 dark:bg-orange-950/50">
          <AlertCircle className="size-4 text-orange-600 dark:text-orange-400" />
          <AlertTitle className="text-orange-900 dark:text-orange-100">
            {t("sections.limitations.staticAnalysis.title")}
          </AlertTitle>
          <AlertDescription className="text-orange-800 dark:text-orange-200">
            {t("sections.limitations.staticAnalysis.description")}
          </AlertDescription>
        </Alert>
      </section>

      {/* Accuracy */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.accuracy.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.accuracy.description")}</p>

        <Alert className="border-green-200 bg-green-50 dark:border-green-900 dark:bg-green-950/50">
          <CheckCircle2 className="size-4 text-green-600 dark:text-green-400" />
          <AlertTitle className="text-green-900 dark:text-green-100">
            {t("sections.accuracy.typical.title")}
          </AlertTitle>
          <AlertDescription className="text-green-800 dark:text-green-200">
            {t("sections.accuracy.typical.description")}
          </AlertDescription>
        </Alert>
      </section>

      {/* FAQ */}
      <section>
        <h2 className="mb-6 text-2xl font-semibold tracking-tight">{t("sections.faq.title")}</h2>

        <div className="space-y-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-base">{t("sections.faq.q1.question")}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {t("sections.faq.q1.answer")}
              </p>
            </CardContent>
          </Card>
        </div>
      </section>
    </div>
  );
};
