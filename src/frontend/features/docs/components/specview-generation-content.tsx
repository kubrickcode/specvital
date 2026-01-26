"use client";

import {
  AlertTriangle,
  Brain,
  CheckCircle2,
  FileText,
  HelpCircle,
  Layers,
  Sparkles,
  Zap,
} from "lucide-react";
import { useTranslations } from "next-intl";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

const CLASSIFICATION_TYPES = ["functional", "edge", "integration", "unit", "performance"] as const;

const LIMITATIONS = ["accuracy", "language", "context", "complexity"] as const;

export const SpecviewGenerationContent = () => {
  const t = useTranslations("docs.specviewGeneration");

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

      {/* Classification Logic */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.classification.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.classification.description")}
        </p>

        <Card className="py-0">
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.classification.table.type")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.classification.table.description")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.classification.table.example")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {CLASSIFICATION_TYPES.map((typeKey) => (
                  <TableRow key={typeKey}>
                    <TableCell className="px-4 py-2.5">
                      <Badge
                        className={
                          typeKey === "functional"
                            ? "bg-green-500 hover:bg-green-600"
                            : typeKey === "edge"
                              ? "bg-amber-500 hover:bg-amber-600"
                              : typeKey === "integration"
                                ? "bg-blue-500 hover:bg-blue-600"
                                : typeKey === "unit"
                                  ? "bg-slate-500 hover:bg-slate-600"
                                  : "bg-purple-500 hover:bg-purple-600"
                        }
                      >
                        {t(`sections.classification.types.${typeKey}.name`)}
                      </Badge>
                    </TableCell>
                    <TableCell className="px-4 py-2.5 text-muted-foreground">
                      {t(`sections.classification.types.${typeKey}.description`)}
                    </TableCell>
                    <TableCell className="px-4 py-2.5 font-mono text-xs text-muted-foreground">
                      {t(`sections.classification.types.${typeKey}.example`)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </section>

      {/* Output Format */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">{t("sections.output.title")}</h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.output.description")}</p>

        <div className="grid gap-3 sm:grid-cols-2">
          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Layers className="size-4 text-primary" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.output.features.structure.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.output.features.structure.description")}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Brain className="size-4 text-primary" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.output.features.context.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.output.features.context.description")}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <FileText className="size-4 text-primary" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.output.features.markdown.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.output.features.markdown.description")}
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

      {/* Limitations */}
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

      {/* Best Practices */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.bestPractices.title")}
        </h2>

        <Card>
          <CardContent className="pt-6">
            <ul className="space-y-3 text-sm text-muted-foreground">
              <li className="flex items-start gap-3">
                <CheckCircle2 className="mt-0.5 size-5 shrink-0 text-green-500" />
                <div>
                  <p className="font-medium text-foreground">
                    {t("sections.bestPractices.tips.naming.title")}
                  </p>
                  <p>{t("sections.bestPractices.tips.naming.description")}</p>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <CheckCircle2 className="mt-0.5 size-5 shrink-0 text-green-500" />
                <div>
                  <p className="font-medium text-foreground">
                    {t("sections.bestPractices.tips.structure.title")}
                  </p>
                  <p>{t("sections.bestPractices.tips.structure.description")}</p>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <CheckCircle2 className="mt-0.5 size-5 shrink-0 text-green-500" />
                <div>
                  <p className="font-medium text-foreground">
                    {t("sections.bestPractices.tips.coverage.title")}
                  </p>
                  <p>{t("sections.bestPractices.tips.coverage.description")}</p>
                </div>
              </li>
            </ul>
          </CardContent>
        </Card>
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
                <CardTitle className="text-base">{t("sections.faq.q4.question")}</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {t("sections.faq.q4.answer")}
              </p>
            </CardContent>
          </Card>
        </div>
      </section>
    </div>
  );
};
