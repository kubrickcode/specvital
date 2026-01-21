"use client";

import { AlertCircle, CheckCircle2, FileCode2, Lightbulb, Zap } from "lucide-react";
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

const SUPPORTED_FRAMEWORKS = [
  { frameworks: "Jest, Vitest, Playwright, Cypress, Mocha", language: "JavaScript / TypeScript" },
  { frameworks: "go testing", language: "Go" },
  { frameworks: "pytest, unittest", language: "Python" },
  { frameworks: "JUnit 4, JUnit 5, TestNG", language: "Java" },
  { frameworks: "Kotest", language: "Kotlin" },
  { frameworks: "NUnit, xUnit, MSTest", language: "C#" },
  { frameworks: "RSpec, Minitest", language: "Ruby" },
  { frameworks: "PHPUnit", language: "PHP" },
  { frameworks: "cargo test", language: "Rust" },
  { frameworks: "Google Test", language: "C++" },
  { frameworks: "XCTest, Swift Testing", language: "Swift" },
] as const;

const FILE_PATTERNS = [
  { language: "Go", patterns: "*_test.go" },
  { language: "JavaScript / TypeScript", patterns: "*.test.ts, *.spec.ts, __tests__/*" },
  { language: "Python", patterns: "test_*.py, *_test.py" },
  { language: "Java / Kotlin", patterns: "*Test.java, *Tests.java, src/test/*" },
  { language: "C#", patterns: "*Test.cs, *Tests.cs, *.Tests/*" },
  { language: "Ruby", patterns: "*_spec.rb, *_test.rb" },
  { language: "Rust", patterns: "*_test.rs, tests/*, #[test]" },
] as const;

const DYNAMIC_TEST_EXAMPLES = [
  { cli: "3", parser: "1", pattern: "it.each([1,2,3])('test %i', ...)" },
  { cli: "N", parser: "1", pattern: "@pytest.mark.parametrize" },
  { cli: "N", parser: "1", pattern: "@ParameterizedTest (JUnit 5)" },
  { cli: "N", parser: "1 per t.Run", pattern: "t.Run in loop (Go)" },
  { cli: "N", parser: "N", pattern: "[TestCase] x N (C#)" },
] as const;

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
        <p className="mb-6 leading-7 text-muted-foreground">{t("sections.how.description")}</p>

        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base">{t("sections.how.detectionOrder.title")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-start gap-3 rounded-lg border bg-muted/30 p-3">
              <Badge className="mt-0.5 shrink-0" variant="secondary">
                1
              </Badge>
              <div>
                <p className="font-medium">{t("sections.how.detectionOrder.step1Title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.how.detectionOrder.step1Desc")}
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3 rounded-lg border bg-muted/30 p-3">
              <Badge className="mt-0.5 shrink-0" variant="secondary">
                2
              </Badge>
              <div>
                <p className="font-medium">{t("sections.how.detectionOrder.step2Title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.how.detectionOrder.step2Desc")}
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3 rounded-lg border bg-muted/30 p-3">
              <Badge className="mt-0.5 shrink-0" variant="secondary">
                3
              </Badge>
              <div>
                <p className="font-medium">{t("sections.how.detectionOrder.step3Title")}</p>
                <p className="text-sm text-muted-foreground">
                  {t("sections.how.detectionOrder.step3Desc")}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </section>

      {/* Supported Frameworks */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.frameworks.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.frameworks.description")}
        </p>

        <Card className="py-0">
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="w-[200px] px-4 py-2 font-semibold">
                    {t("sections.frameworks.table.language")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.frameworks.table.frameworks")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {SUPPORTED_FRAMEWORKS.map(({ frameworks, language }) => (
                  <TableRow key={language}>
                    <TableCell className="px-4 py-2.5 font-medium">{language}</TableCell>
                    <TableCell className="px-4 py-2.5 text-muted-foreground">
                      {frameworks}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </section>

      {/* File Patterns */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.filePatterns.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.filePatterns.description")}
        </p>

        <Card className="py-0">
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="w-[200px] px-4 py-2 font-semibold">
                    {t("sections.filePatterns.table.language")}
                  </TableHead>
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.filePatterns.table.patterns")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {FILE_PATTERNS.map(({ language, patterns }) => (
                  <TableRow key={language}>
                    <TableCell className="px-4 py-2.5 font-medium">{language}</TableCell>
                    <TableCell className="px-4 py-2.5">
                      <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-sm">
                        {patterns}
                      </code>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </section>

      {/* Why Counts May Differ */}
      <section>
        <h2 className="mb-4 text-2xl font-semibold tracking-tight">
          {t("sections.countDifference.title")}
        </h2>
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.countDifference.description")}
        </p>

        <Card className="mb-4 py-0">
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="px-4 py-2 font-semibold">
                    {t("sections.countDifference.table.pattern")}
                  </TableHead>
                  <TableHead className="w-[120px] px-4 py-2 text-center font-semibold">
                    {t("sections.countDifference.table.specvital")}
                  </TableHead>
                  <TableHead className="w-[80px] px-4 py-2 text-center font-semibold">
                    {t("sections.countDifference.table.cli")}
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {DYNAMIC_TEST_EXAMPLES.map(({ cli, parser, pattern }) => (
                  <TableRow key={pattern}>
                    <TableCell className="px-4 py-2.5">
                      <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-sm">
                        {pattern}
                      </code>
                    </TableCell>
                    <TableCell className="px-4 py-2.5 text-center font-medium">{parser}</TableCell>
                    <TableCell className="px-4 py-2.5 text-center text-muted-foreground">
                      {cli}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <Alert className="border-amber-200 bg-amber-50 dark:border-amber-900 dark:bg-amber-950/50">
          <Lightbulb className="size-4 text-amber-600 dark:text-amber-400" />
          <AlertTitle className="text-amber-900 dark:text-amber-100">
            {t("sections.countDifference.tip.title")}
          </AlertTitle>
          <AlertDescription className="text-amber-800 dark:text-amber-200">
            {t("sections.countDifference.tip.description")}
          </AlertDescription>
        </Alert>
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
        <p className="mb-4 leading-7 text-muted-foreground">
          {t("sections.limitations.description")}
        </p>

        <Alert className="mb-4 border-orange-200 bg-orange-50 dark:border-orange-900 dark:bg-orange-950/50">
          <AlertCircle className="size-4 text-orange-600 dark:text-orange-400" />
          <AlertTitle className="text-orange-900 dark:text-orange-100">
            {t("sections.limitations.indirect.title")}
          </AlertTitle>
          <AlertDescription className="text-orange-800 dark:text-orange-200">
            {t("sections.limitations.indirect.description")}
          </AlertDescription>
        </Alert>

        <Card>
          <CardContent className="pt-6">
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li className="flex items-start gap-2">
                <FileCode2 className="mt-0.5 size-4 shrink-0" />
                <span>{t("sections.limitations.items.singleFile")}</span>
              </li>
              <li className="flex items-start gap-2">
                <FileCode2 className="mt-0.5 size-4 shrink-0" />
                <span>{t("sections.limitations.items.dynamicImports")}</span>
              </li>
              <li className="flex items-start gap-2">
                <FileCode2 className="mt-0.5 size-4 shrink-0" />
                <span>{t("sections.limitations.items.macros")}</span>
              </li>
            </ul>
          </CardContent>
        </Card>
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

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-base">{t("sections.faq.q2.question")}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {t("sections.faq.q2.answer")}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-base">{t("sections.faq.q3.question")}</CardTitle>
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
