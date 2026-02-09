"use client";

import {
  AlertCircle,
  BarChart3,
  CalendarClock,
  CheckCircle2,
  Database,
  FileText,
  HelpCircle,
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
import { Link } from "@/i18n/navigation";

const USAGE_EVENTS = [
  { countsUsage: true, eventKey: "analysisRun" },
  { countsUsage: true, eventKey: "specDocGeneration" },
  { countsUsage: false, eventKey: "cacheHit" },
  { countsUsage: false, eventKey: "viewResults" },
  { countsUsage: false, eventKey: "exportDownload" },
] as const;

const CACHE_SCENARIOS = [
  { cached: true, cacheKey: "sameCommit" },
  { cached: true, cacheKey: "newCommit" },
  { cached: true, cacheKey: "parserUpdate" },
  { cached: false, cacheKey: "manualReanalysis" },
] as const;

export const UsageBillingContent = () => {
  const t = useTranslations("docs.usageBilling");

  return (
    <div className="space-y-8">
      {/* What Counts as Usage */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.whatCounts.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.whatCounts.description")}
        </p>

        <Alert className="mb-6 border-blue-200 bg-blue-50 dark:border-blue-900 dark:bg-blue-950/50">
          <Zap className="size-4 text-blue-600 dark:text-blue-400" />
          <AlertTitle className="text-blue-900 dark:text-blue-100">
            {t("sections.whatCounts.keyPoint.title")}
          </AlertTitle>
          <AlertDescription className="text-blue-800 dark:text-blue-200">
            {t("sections.whatCounts.keyPoint.description")}
          </AlertDescription>
        </Alert>

        <Card className="py-0">
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.whatCounts.table.event")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.whatCounts.table.description")}
                  </TableHead>
                  <TableHead className="w-[100px] px-4 py-2 text-center font-semibold">
                    {t("sections.whatCounts.table.countsUsage")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {USAGE_EVENTS.map(({ countsUsage, eventKey }) => (
                  <TableRow key={eventKey}>
                    <TableCell className="px-4 py-2.5 font-medium">
                      {t(`sections.whatCounts.events.${eventKey}.name`)}
                    </TableCell>
                    <TableCell className="px-4 py-2.5 text-muted-foreground">
                      {t(`sections.whatCounts.events.${eventKey}.description`)}
                    </TableCell>
                    <TableCell className="px-4 py-2.5 text-center">
                      {countsUsage ? (
                        <Badge variant="default">{t("sections.whatCounts.table.yes")}</Badge>
                      ) : (
                        <Badge variant="secondary">{t("sections.whatCounts.table.no")}</Badge>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </section>

      {/* Two Types of Usage */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.usageTypes.title")}
        </h2>
        <p className="mb-6 leading-7 text-muted-foreground">
          {t("sections.usageTypes.description")}
        </p>

        <div className="grid gap-4 md:grid-cols-2">
          <Card className="border-l-4 border-l-green-500">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <BarChart3 className="size-5 text-green-600 dark:text-green-400" />
                <CardTitle className="text-base">
                  {t("sections.usageTypes.analysis.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-sm text-muted-foreground">
                {t("sections.usageTypes.analysis.description")}
              </p>
              <ul className="space-y-1 text-sm text-muted-foreground">
                <li className="flex items-start gap-2">
                  <CheckCircle2 className="mt-0.5 size-4 shrink-0 text-green-500" />
                  <span>{t("sections.usageTypes.analysis.item1")}</span>
                </li>
                <li className="flex items-start gap-2">
                  <CheckCircle2 className="mt-0.5 size-4 shrink-0 text-green-500" />
                  <span>{t("sections.usageTypes.analysis.item2")}</span>
                </li>
              </ul>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-purple-500">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <FileText className="size-5 text-purple-600 dark:text-purple-400" />
                <CardTitle className="text-base">
                  {t("sections.usageTypes.specview.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-sm text-muted-foreground">
                {t("sections.usageTypes.specview.description")}
              </p>
              <ul className="space-y-1 text-sm text-muted-foreground">
                <li className="flex items-start gap-2">
                  <CheckCircle2 className="mt-0.5 size-4 shrink-0 text-purple-500" />
                  <span>{t("sections.usageTypes.specview.item1")}</span>
                </li>
                <li className="flex items-start gap-2">
                  <CheckCircle2 className="mt-0.5 size-4 shrink-0 text-purple-500" />
                  <span>{t("sections.usageTypes.specview.item2")}</span>
                </li>
              </ul>
            </CardContent>
          </Card>
        </div>
      </section>

      {/* Cache Behavior */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">{t("sections.cache.title")}</h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.cache.description")}</p>

        <Card className="mb-4 py-0">
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.cache.table.scenario")}
                  </TableHead>
                  <TableHead className="w-[120px] px-4 py-2 text-center font-semibold">
                    {t("sections.cache.table.cached")}
                  </TableHead>
                  <TableHead className="w-[120px] px-4 py-2 text-center font-semibold">
                    {t("sections.cache.table.usageDeducted")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {CACHE_SCENARIOS.map(({ cached, cacheKey }) => (
                  <TableRow key={cacheKey}>
                    <TableCell className="px-4 py-2.5">
                      <div>
                        <p className="font-medium">
                          {t(`sections.cache.scenarios.${cacheKey}.name`)}
                        </p>
                        <p className="text-sm text-muted-foreground">
                          {t(`sections.cache.scenarios.${cacheKey}.description`)}
                        </p>
                      </div>
                    </TableCell>
                    <TableCell className="px-4 py-2.5 text-center">
                      {cached ? (
                        <Badge className="bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
                          <Database className="mr-1 size-3" />
                          {t("sections.cache.table.yes")}
                        </Badge>
                      ) : (
                        <Badge variant="secondary">{t("sections.cache.table.no")}</Badge>
                      )}
                    </TableCell>
                    <TableCell className="px-4 py-2.5 text-center">
                      {cached ? (
                        <Badge variant="secondary">{t("sections.cache.table.no")}</Badge>
                      ) : (
                        <Badge variant="default">{t("sections.cache.table.yes")}</Badge>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <div className="space-y-3">
          <Alert className="border-green-200 bg-green-50 dark:border-green-900 dark:bg-green-950/50">
            <CheckCircle2 className="size-4 text-green-600 dark:text-green-400" />
            <AlertTitle className="text-green-900 dark:text-green-100">
              {t("sections.cache.tip.title")}
            </AlertTitle>
            <AlertDescription className="text-green-800 dark:text-green-200">
              {t("sections.cache.tip.description")}
            </AlertDescription>
          </Alert>

          <Alert className="border-purple-200 bg-purple-50 dark:border-purple-900 dark:bg-purple-950/50">
            <FileText className="size-4 text-purple-600 dark:text-purple-400" />
            <AlertTitle className="text-purple-900 dark:text-purple-100">
              {t("sections.cache.specDocCache.title")}
            </AlertTitle>
            <AlertDescription className="text-purple-800 dark:text-purple-200">
              {t("sections.cache.specDocCache.description")}
            </AlertDescription>
          </Alert>
        </div>
      </section>

      {/* Billing Cycle */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">{t("sections.cycle.title")}</h2>
        <p className="mb-6 leading-7 text-muted-foreground">{t("sections.cycle.description")}</p>

        <Card>
          <CardHeader className="pb-3">
            <div className="flex items-center gap-2">
              <CalendarClock className="size-5 text-primary" />
              <CardTitle className="text-base">{t("sections.cycle.howItWorks.title")}</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-start gap-3 rounded-lg border bg-muted/30 p-3">
              <Badge className="mt-0.5 shrink-0" variant="secondary">
                1
              </Badge>
              <div>
                <p className="font-medium">{t("sections.cycle.howItWorks.step1Title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.cycle.howItWorks.step1Desc")}
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3 rounded-lg border bg-muted/30 p-3">
              <Badge className="mt-0.5 shrink-0" variant="secondary">
                2
              </Badge>
              <div>
                <p className="font-medium">{t("sections.cycle.howItWorks.step2Title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.cycle.howItWorks.step2Desc")}
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3 rounded-lg border bg-muted/30 p-3">
              <Badge className="mt-0.5 shrink-0" variant="secondary">
                3
              </Badge>
              <div>
                <p className="font-medium">{t("sections.cycle.howItWorks.step3Title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.cycle.howItWorks.step3Desc")}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </section>

      {/* Plan Limits */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">{t("sections.limits.title")}</h2>
        <p className="leading-7 text-muted-foreground">{t("sections.limits.description")}</p>
        <Link
          className="mt-2 inline-flex items-center text-sm font-medium text-primary underline-offset-4 hover:underline"
          href="/pricing"
        >
          {t("sections.limits.viewPricing")}
        </Link>
      </section>

      {/* Quota Warnings */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.warnings.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.warnings.description")}</p>

        <div className="space-y-3">
          <Alert className="border-amber-200 bg-amber-50 dark:border-amber-900 dark:bg-amber-950/50">
            <AlertCircle className="size-4 text-amber-600 dark:text-amber-400" />
            <AlertTitle className="text-amber-900 dark:text-amber-100">
              {t("sections.warnings.approaching.title")}
            </AlertTitle>
            <AlertDescription className="text-amber-800 dark:text-amber-200">
              {t("sections.warnings.approaching.description")}
            </AlertDescription>
          </Alert>

          <Alert className="border-red-200 bg-red-50 dark:border-red-900 dark:bg-red-950/50">
            <AlertCircle className="size-4 text-red-600 dark:text-red-400" />
            <AlertTitle className="text-red-900 dark:text-red-100">
              {t("sections.warnings.exceeded.title")}
            </AlertTitle>
            <AlertDescription className="text-red-800 dark:text-red-200">
              {t("sections.warnings.exceeded.description")}
            </AlertDescription>
          </Alert>
        </div>
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
