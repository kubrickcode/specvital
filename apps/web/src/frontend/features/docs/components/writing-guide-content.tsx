"use client";

import {
  ArrowRight,
  Bot,
  CheckCircle2,
  Code2,
  Copy,
  Lightbulb,
  Rocket,
  Type,
  XCircle,
} from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

const PIPELINE_STEPS = ["detect", "parse", "analyze", "specDocument"] as const;

const FILE_ORG_TIPS = ["naming", "colocation", "separation"] as const;

const ROADMAP_FEATURES = ["naming", "structure", "preview", "architecture"] as const;

const NAMING_EXAMPLES = [
  { context: "-", level: "low" as const, name: "tc_01", result: "tc_01" },
  {
    context: "AuthService > Login",
    level: "medium" as const,
    name: "err",
    result: "Login error handling",
  },
  {
    context: "OrderService > Sorting",
    level: "high" as const,
    name: "created desc",
    result: "Orders sorted by creation date in descending order",
  },
  {
    context: "AuthService > Login",
    level: "high" as const,
    name: "should reject with expired token",
    result: "Rejects login when token is expired",
  },
] as const;

const CONFIDENCE_COLORS: Record<string, string> = {
  high: "bg-green-500 hover:bg-green-600",
  low: "bg-red-500 hover:bg-red-600",
  medium: "bg-amber-500 hover:bg-amber-600",
};

const SIGNAL_LEVELS = [
  { color: "bg-green-500", key: "strongest" },
  { color: "bg-blue-500", key: "strong" },
  { color: "bg-amber-500", key: "moderate" },
  { color: "bg-slate-400", key: "weak" },
] as const;

const FILE_ORG_BAD = `tests/
  test1.test.ts
  test2.test.ts
  helpers.test.ts`;

const FILE_ORG_GOOD = `src/
  auth/
    __tests__/
      login.test.ts
      registration.test.ts
  payment/
    __tests__/
      checkout.test.ts
      refund.test.ts`;

const STRUCTURE_BAD = `test('login works', () => { ... })
test('logout works', () => { ... })
test('token refreshes', () => { ... })`;

const STRUCTURE_GOOD = `describe('AuthService', () => {
  describe('Login', () => {
    it('should authenticate with valid credentials', ...)
    it('should reject invalid password', ...)
  })
  describe('Token', () => {
    it('should refresh expired token', ...)
  })
})`;

const ARCH_BAD = `import { render } from '@/utils'
import { mockFn } from '@/test-helpers'`;

const ARCH_GOOD = `import { authService } from '@/services/auth'
import { loginValidator } from '@/validators/auth'
import { TokenManager } from '@/modules/auth/token'`;

const PROMPT_PRINCIPLES = [
  "suiteHierarchy",
  "behaviorNames",
  "domainPaths",
  "domainImports",
] as const;

const PROMPT_TEMPLATE = `# Test Writing Rules for SpecVital Compatibility

Follow these rules so generated tests produce high-quality SpecVital analysis results.

## 1. Suite Hierarchy Provides Domain Context
Wrap tests in nested describe/context blocks (or language equivalent)
that name the domain and feature.
The suite path is the strongest structural signal for domain classification.

  GOOD: describe('OrderService') > describe('Sorting') > it('created desc')
  BAD:  it('OrderService returns items sorted by creation date')

Short test names are fine when the suite context is rich.

## 2. Test Names Describe Behavior
Write test names as behavior descriptions, not implementation details.

  GOOD: 'rejects expired tokens with 401 status'
  BAD:  'test token validation'
  BAD:  'works correctly'

## 3. File Paths Mirror Domain Structure
Organize test files under directories that reflect business domains.

  GOOD: tests/auth/session-management.test.ts
  BAD:  tests/utils.test.ts

## 4. Import Real Domain Modules
Import actual services/modules under test by name.
Import statements are the strongest signal for domain classification.

  GOOD: import { OrderService } from '@/modules/order'
  BAD:  (inline all logic without imports)`;

type TranslationFn = ReturnType<typeof useTranslations<"docs.writingGuide">>;

const AiToolConfigSection = ({ t }: { t: TranslationFn }) => {
  const [copied, setCopied] = useState(false);

  const copyPrompt = async () => {
    await navigator.clipboard.writeText(PROMPT_TEMPLATE);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <section>
      <h2 className="mb-4 text-2xl font-semibold tracking-tight">
        {t("sections.aiToolConfig.title")}
      </h2>
      <p className="mb-4 leading-7 text-muted-foreground">
        {t("sections.aiToolConfig.description")}
      </p>

      <div className="grid gap-3 sm:grid-cols-2">
        {PROMPT_PRINCIPLES.map((principle) => (
          <Card key={principle}>
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Bot className="size-4 text-primary" />
                <CardTitle className="text-sm font-medium">
                  {t(`sections.aiToolConfig.principles.${principle}.title`)}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t(`sections.aiToolConfig.principles.${principle}.description`)}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="mt-4">
        <Dialog>
          <DialogTrigger asChild>
            <Button className="gap-2" variant="outline">
              <Copy className="size-4" />
              {t("sections.aiToolConfig.viewTemplate")}
            </Button>
          </DialogTrigger>
          <DialogContent className="max-h-[80vh] max-w-2xl">
            <DialogHeader>
              <DialogTitle>{t("sections.aiToolConfig.dialog.title")}</DialogTitle>
              <DialogDescription>{t("sections.aiToolConfig.dialog.description")}</DialogDescription>
            </DialogHeader>
            <pre className="max-h-[50vh] overflow-auto rounded-md bg-muted p-4 text-xs leading-relaxed">
              <code>{PROMPT_TEMPLATE}</code>
            </pre>
            <DialogFooter>
              <Button className="gap-1.5" onClick={copyPrompt} variant="outline">
                <Copy className="size-3.5" />
                {copied
                  ? t("sections.aiToolConfig.dialog.copied")
                  : t("sections.aiToolConfig.dialog.copy")}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </section>
  );
};

export const WritingGuideContent = () => {
  const t = useTranslations("docs.writingGuide");

  return (
    <div className="space-y-8">
      {/* Pipeline Overview */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.pipeline.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.pipeline.description")}</p>

        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
          {PIPELINE_STEPS.map((step, index) => (
            <Card className="relative" key={step}>
              <CardHeader className="pb-2">
                <div className="flex items-center gap-2">
                  <div className="flex size-6 items-center justify-center rounded-full bg-primary text-xs font-bold text-primary-foreground">
                    {index + 1}
                  </div>
                  <CardTitle className="text-sm font-medium">
                    {t(`sections.pipeline.steps.${step}.title`)}
                  </CardTitle>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-xs text-muted-foreground">
                  {t(`sections.pipeline.steps.${step}.description`)}
                </p>
              </CardContent>
              {index < PIPELINE_STEPS.length - 1 && (
                <ArrowRight className="absolute -right-2.5 top-1/2 z-10 hidden size-4 -translate-y-1/2 text-muted-foreground lg:block" />
              )}
            </Card>
          ))}
        </div>

        <Alert className="mt-4 border-blue-200 bg-blue-50 dark:border-blue-900 dark:bg-blue-950/50">
          <Lightbulb className="size-4 text-blue-600 dark:text-blue-400" />
          <AlertTitle className="text-blue-900 dark:text-blue-100">
            {t("sections.pipeline.keyPoint.title")}
          </AlertTitle>
          <AlertDescription className="text-blue-800 dark:text-blue-200">
            {t("sections.pipeline.keyPoint.description")}
          </AlertDescription>
        </Alert>
      </section>

      {/* Test Naming (highest priority) */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.testNaming.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.testNaming.description")}
        </p>

        <Alert className="mb-4 border-purple-200 bg-purple-50 dark:border-purple-900 dark:bg-purple-950/50">
          <Lightbulb className="size-4 text-purple-600 dark:text-purple-400" />
          <AlertTitle className="text-purple-900 dark:text-purple-100">
            {t("sections.testNaming.contextNote.title")}
          </AlertTitle>
          <AlertDescription className="text-purple-800 dark:text-purple-200">
            {t("sections.testNaming.contextNote.description")}
          </AlertDescription>
        </Alert>

        <Card className="py-0">
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.testNaming.table.testName")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.testNaming.table.context")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.testNaming.table.result")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.testNaming.table.confidence")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {NAMING_EXAMPLES.map((example) => (
                  <TableRow key={example.name}>
                    <TableCell className="px-4 py-2.5 font-mono text-xs">{example.name}</TableCell>
                    <TableCell className="px-4 py-2.5 text-xs text-muted-foreground">
                      {example.context}
                    </TableCell>
                    <TableCell className="px-4 py-2.5 text-sm text-muted-foreground">
                      {example.result}
                    </TableCell>
                    <TableCell className="px-4 py-2.5">
                      <Badge className={CONFIDENCE_COLORS[example.level]}>
                        {t(`sections.testNaming.confidence.${example.level}`)}
                      </Badge>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <div className="mt-4 grid gap-3 sm:grid-cols-2">
          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Type className="size-4 text-primary" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.testNaming.tips.pattern.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.testNaming.tips.pattern.description")}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <Lightbulb className="size-4 text-primary" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.testNaming.tips.specific.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                {t("sections.testNaming.tips.specific.description")}
              </p>
            </CardContent>
          </Card>
        </div>
      </section>

      {/* Test Structure */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.testStructure.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.testStructure.description")}
        </p>

        <div className="grid gap-4 sm:grid-cols-2">
          <Card className="overflow-hidden border-l-4 border-l-red-400">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <XCircle className="size-4 text-red-500" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.testStructure.bad.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-xs text-muted-foreground">
                {t("sections.testStructure.bad.description")}
              </p>
              <pre className="overflow-x-auto rounded-md bg-muted p-3 text-xs leading-relaxed">
                <code>{STRUCTURE_BAD}</code>
              </pre>
            </CardContent>
          </Card>

          <Card className="overflow-hidden border-l-4 border-l-green-400">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <CheckCircle2 className="size-4 text-green-500" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.testStructure.good.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-xs text-muted-foreground">
                {t("sections.testStructure.good.description")}
              </p>
              <pre className="overflow-x-auto rounded-md bg-muted p-3 text-xs leading-relaxed">
                <code>{STRUCTURE_GOOD}</code>
              </pre>
            </CardContent>
          </Card>
        </div>

        <Alert className="mt-4 border-amber-200 bg-amber-50 dark:border-amber-900 dark:bg-amber-950/50">
          <Code2 className="size-4 text-amber-600 dark:text-amber-400" />
          <AlertTitle className="text-amber-900 dark:text-amber-100">
            {t("sections.testStructure.tip.title")}
          </AlertTitle>
          <AlertDescription className="text-amber-800 dark:text-amber-200">
            {t("sections.testStructure.tip.description")}
          </AlertDescription>
        </Alert>
      </section>

      {/* File Organization */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.fileOrganization.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.fileOrganization.description")}
        </p>

        <div className="grid gap-4 sm:grid-cols-2">
          <Card className="overflow-hidden border-l-4 border-l-red-400">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <XCircle className="size-4 text-red-500" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.fileOrganization.bad.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-xs text-muted-foreground">
                {t("sections.fileOrganization.bad.description")}
              </p>
              <pre className="overflow-x-auto rounded-md bg-muted p-3 text-xs leading-relaxed">
                <code>{FILE_ORG_BAD}</code>
              </pre>
            </CardContent>
          </Card>

          <Card className="overflow-hidden border-l-4 border-l-green-400">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <CheckCircle2 className="size-4 text-green-500" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.fileOrganization.good.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-xs text-muted-foreground">
                {t("sections.fileOrganization.good.description")}
              </p>
              <pre className="overflow-x-auto rounded-md bg-muted p-3 text-xs leading-relaxed">
                <code>{FILE_ORG_GOOD}</code>
              </pre>
            </CardContent>
          </Card>
        </div>

        <Card className="mt-4">
          <CardContent className="pt-6">
            <ul className="space-y-3 text-sm text-muted-foreground">
              {FILE_ORG_TIPS.map((tip) => (
                <li className="flex items-start gap-3" key={tip}>
                  <CheckCircle2 className="mt-0.5 size-5 shrink-0 text-green-500" />
                  <div>
                    <p className="font-medium text-foreground">
                      {t(`sections.fileOrganization.tips.${tip}.title`)}
                    </p>
                    <p>{t(`sections.fileOrganization.tips.${tip}.description`)}</p>
                  </div>
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>
      </section>

      {/* Project Architecture */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.architecture.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.architecture.description")}
        </p>

        <Card className="mb-4">
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">
              {t("sections.architecture.signalStrength.title")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {SIGNAL_LEVELS.map(({ color, key }) => (
                <div className="flex items-center gap-3" key={key}>
                  <div className={`size-2.5 rounded-full ${color}`} />
                  <span className="text-sm text-muted-foreground">
                    {t(`sections.architecture.signalStrength.${key}`)}
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <div className="grid gap-4 sm:grid-cols-2">
          <Card className="overflow-hidden border-l-4 border-l-red-400">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <XCircle className="size-4 text-red-500" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.architecture.bad.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-xs text-muted-foreground">
                {t("sections.architecture.bad.description")}
              </p>
              <pre className="overflow-x-auto rounded-md bg-muted p-3 text-xs leading-relaxed">
                <code>{ARCH_BAD}</code>
              </pre>
            </CardContent>
          </Card>

          <Card className="overflow-hidden border-l-4 border-l-green-400">
            <CardHeader className="pb-2">
              <div className="flex items-center gap-2">
                <CheckCircle2 className="size-4 text-green-500" />
                <CardTitle className="text-sm font-medium">
                  {t("sections.architecture.good.title")}
                </CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="mb-3 text-xs text-muted-foreground">
                {t("sections.architecture.good.description")}
              </p>
              <pre className="overflow-x-auto rounded-md bg-muted p-3 text-xs leading-relaxed">
                <code>{ARCH_GOOD}</code>
              </pre>
            </CardContent>
          </Card>
        </div>

        <Alert className="mt-4 border-green-200 bg-green-50 dark:border-green-900 dark:bg-green-950/50">
          <Lightbulb className="size-4 text-green-600 dark:text-green-400" />
          <AlertTitle className="text-green-900 dark:text-green-100">
            {t("sections.architecture.tip.title")}
          </AlertTitle>
          <AlertDescription className="text-green-800 dark:text-green-200">
            {t("sections.architecture.tip.description")}
          </AlertDescription>
        </Alert>
      </section>

      {/* AI Tool Configuration */}
      <AiToolConfigSection t={t} />

      {/* Roadmap */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.roadmap.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">{t("sections.roadmap.description")}</p>

        <div className="grid gap-3 sm:grid-cols-2">
          {ROADMAP_FEATURES.map((feature) => (
            <Card key={feature}>
              <CardHeader className="pb-2">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Rocket className="size-4 text-primary" />
                    <CardTitle className="text-sm font-medium">
                      {t(`sections.roadmap.features.${feature}.title`)}
                    </CardTitle>
                  </div>
                  <Badge className="text-xs" variant="secondary">
                    {t("sections.roadmap.comingSoon")}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground">
                  {t(`sections.roadmap.features.${feature}.description`)}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>
    </div>
  );
};
