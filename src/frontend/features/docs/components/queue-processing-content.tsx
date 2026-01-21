"use client";

import { CheckCircle2, Clock, HelpCircle, Layers, Rocket, Timer, Users, Zap } from "lucide-react";
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

const QUEUE_TIERS = [
  { priorityKey: "free", queueKey: "standard" },
  { priorityKey: "pro", queueKey: "priority" },
  { priorityKey: "proPlusEnterprise", queueKey: "dedicated" },
] as const;

const WAIT_TIME_FACTORS = ["queueLength", "jobSize", "planTier", "timeOfDay"] as const;

export const QueueProcessingContent = () => {
  const t = useTranslations("docs.queueProcessing");

  return (
    <div className="space-y-8">
      {/* Overview */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.overview.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.overview.description")}</p>

        <Alert className="border-blue-200 bg-blue-50 dark:border-blue-900 dark:bg-blue-950/50">
          <Zap className="size-4 text-blue-600 dark:text-blue-400" />
          <AlertTitle className="text-blue-900 dark:text-blue-100">
            {t("sections.overview.keyPoint.title")}
          </AlertTitle>
          <AlertDescription className="text-blue-800 dark:text-blue-200">
            {t("sections.overview.keyPoint.description")}
          </AlertDescription>
        </Alert>
      </section>

      {/* Multi-Queue Architecture */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.architecture.title")}
        </h2>
        <p className="mb-6 leading-7 text-muted-foreground">
          {t("sections.architecture.description")}
        </p>

        <div className="grid gap-4 md:grid-cols-3">
          <Card className="border-l-4 border-l-slate-400">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Users className="size-5 text-slate-500" />
                <CardTitle className="text-base">
                  {t("sections.architecture.standard.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-sm text-muted-foreground">
                {t("sections.architecture.standard.description")}
              </p>
              <Badge variant="secondary">{t("sections.architecture.standard.plans")}</Badge>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-blue-500">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Rocket className="size-5 text-blue-500" />
                <CardTitle className="text-base">
                  {t("sections.architecture.priority.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-sm text-muted-foreground">
                {t("sections.architecture.priority.description")}
              </p>
              <Badge className="bg-blue-500 hover:bg-blue-600">
                {t("sections.architecture.priority.plans")}
              </Badge>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-purple-500">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Zap className="size-5 text-purple-500" />
                <CardTitle className="text-base">
                  {t("sections.architecture.dedicated.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-sm text-muted-foreground">
                {t("sections.architecture.dedicated.description")}
              </p>
              <Badge className="bg-purple-500 hover:bg-purple-600">
                {t("sections.architecture.dedicated.plans")}
              </Badge>
            </CardContent>
          </Card>
        </div>
      </section>

      {/* Priority Tiers Table */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">{t("sections.tiers.title")}</h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.tiers.description")}</p>

        <Card className="py-0">
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.tiers.table.plan")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.tiers.table.queue")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.tiers.table.typicalWait")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.tiers.table.concurrency")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {QUEUE_TIERS.map(({ priorityKey, queueKey }) => (
                  <TableRow key={priorityKey}>
                    <TableCell className="px-4 py-2.5 font-medium">
                      {t(`sections.tiers.priorities.${priorityKey}.plan`)}
                    </TableCell>
                    <TableCell className="px-4 py-2.5">
                      <Badge
                        className={
                          queueKey === "priority"
                            ? "bg-blue-500 hover:bg-blue-600"
                            : queueKey === "dedicated"
                              ? "bg-purple-500 hover:bg-purple-600"
                              : undefined
                        }
                        variant={
                          queueKey === "standard"
                            ? "secondary"
                            : queueKey === "priority"
                              ? "default"
                              : "default"
                        }
                      >
                        {t(`sections.tiers.priorities.${priorityKey}.queue`)}
                      </Badge>
                    </TableCell>
                    <TableCell className="px-4 py-2.5 text-muted-foreground">
                      {t(`sections.tiers.priorities.${priorityKey}.wait`)}
                    </TableCell>
                    <TableCell className="px-4 py-2.5 text-muted-foreground">
                      {t(`sections.tiers.priorities.${priorityKey}.concurrency`)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </section>

      {/* How Queue Processing Works */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.howItWorks.title")}
        </h2>
        <p className="mb-6 leading-7 text-muted-foreground">
          {t("sections.howItWorks.description")}
        </p>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base">{t("sections.howItWorks.steps.title")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-start gap-3 rounded-lg border bg-muted/30 p-3">
              <Badge className="mt-0.5 shrink-0" variant="secondary">
                1
              </Badge>
              <div>
                <p className="font-medium">{t("sections.howItWorks.steps.submit.title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.howItWorks.steps.submit.description")}
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3 rounded-lg border bg-muted/30 p-3">
              <Badge className="mt-0.5 shrink-0" variant="secondary">
                2
              </Badge>
              <div>
                <p className="font-medium">{t("sections.howItWorks.steps.route.title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.howItWorks.steps.route.description")}
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3 rounded-lg border bg-muted/30 p-3">
              <Badge className="mt-0.5 shrink-0" variant="secondary">
                3
              </Badge>
              <div>
                <p className="font-medium">{t("sections.howItWorks.steps.process.title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.howItWorks.steps.process.description")}
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3 rounded-lg border bg-muted/30 p-3">
              <Badge className="mt-0.5 shrink-0" variant="secondary">
                4
              </Badge>
              <div>
                <p className="font-medium">{t("sections.howItWorks.steps.complete.title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.howItWorks.steps.complete.description")}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </section>

      {/* Wait Time Factors */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.waitFactors.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.waitFactors.description")}
        </p>

        <div className="grid gap-3 sm:grid-cols-2">
          {WAIT_TIME_FACTORS.map((factorKey) => (
            <Card key={factorKey}>
              <CardHeader className="pb-2">
                <div className="flex items-center gap-2">
                  {factorKey === "queueLength" && <Layers className="size-4 text-primary" />}
                  {factorKey === "jobSize" && <Timer className="size-4 text-primary" />}
                  {factorKey === "planTier" && <Rocket className="size-4 text-primary" />}
                  {factorKey === "timeOfDay" && <Clock className="size-4 text-primary" />}
                  <CardTitle className="text-sm font-medium">
                    {t(`sections.waitFactors.factors.${factorKey}.title`)}
                  </CardTitle>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground">
                  {t(`sections.waitFactors.factors.${factorKey}.description`)}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      {/* Why Priority Queues Matter */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.whyPriority.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.whyPriority.description")}
        </p>

        <Card>
          <CardContent className="pt-6">
            <ul className="space-y-3 text-sm text-muted-foreground">
              <li className="flex items-start gap-3">
                <CheckCircle2 className="mt-0.5 size-5 shrink-0 text-green-500" />
                <div>
                  <p className="font-medium text-foreground">
                    {t("sections.whyPriority.benefits.fairness.title")}
                  </p>
                  <p>{t("sections.whyPriority.benefits.fairness.description")}</p>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <CheckCircle2 className="mt-0.5 size-5 shrink-0 text-green-500" />
                <div>
                  <p className="font-medium text-foreground">
                    {t("sections.whyPriority.benefits.predictability.title")}
                  </p>
                  <p>{t("sections.whyPriority.benefits.predictability.description")}</p>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <CheckCircle2 className="mt-0.5 size-5 shrink-0 text-green-500" />
                <div>
                  <p className="font-medium text-foreground">
                    {t("sections.whyPriority.benefits.scalability.title")}
                  </p>
                  <p>{t("sections.whyPriority.benefits.scalability.description")}</p>
                </div>
              </li>
            </ul>
          </CardContent>
        </Card>
      </section>

      {/* Real-Time Status */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">{t("sections.status.title")}</h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.status.description")}</p>

        <div className="grid gap-3 sm:grid-cols-3">
          <Card className="border-l-4 border-l-amber-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-amber-700 dark:text-amber-400">
                {t("sections.status.states.queued.label")}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.status.states.queued.description")}
              </p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-blue-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-blue-700 dark:text-blue-400">
                {t("sections.status.states.processing.label")}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.status.states.processing.description")}
              </p>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-green-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-green-700 dark:text-green-400">
                {t("sections.status.states.complete.label")}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.status.states.complete.description")}
              </p>
            </CardContent>
          </Card>
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
