import type { components } from "@/lib/api/generated-types";

type AnalysisResult = components["schemas"]["AnalysisResult"];
type TestSuite = components["schemas"]["TestSuite"];
type TestCase = components["schemas"]["TestCase"];

type ExportOptions = {
  includeLineNumbers?: boolean;
  includeTimestamp?: boolean;
};

const formatTestCase = (test: TestCase, options: ExportOptions): string => {
  const lineInfo = options.includeLineNumbers ? ` (L${test.line})` : "";

  switch (test.status) {
    case "active":
      return `- [x] ${test.name}${lineInfo}`;
    case "focused":
      return `- [!] ${test.name} *(focused)*${lineInfo}`;
    case "skipped":
      return `- [ ] ~~${test.name}~~${lineInfo}`;
    case "todo":
      return `- [ ] ${test.name} *(todo)*${lineInfo}`;
    case "xfail":
      return `- [?] ${test.name} *(expected fail)*${lineInfo}`;
  }
};

const formatTestSuite = (suite: TestSuite, options: ExportOptions): string => {
  const lines: string[] = [];

  lines.push(`### ${suite.suiteName}`);
  lines.push("");
  lines.push(`📁 \`${suite.filePath}\``);
  lines.push(`🔧 ${suite.framework}`);
  lines.push("");

  for (const test of suite.tests) {
    lines.push(formatTestCase(test, options));
  }

  return lines.join("\n");
};

const groupSuitesByFramework = (suites: TestSuite[]): Map<string, TestSuite[]> => {
  const grouped = new Map<string, TestSuite[]>();

  for (const suite of suites) {
    const existing = grouped.get(suite.framework) ?? [];
    existing.push(suite);
    grouped.set(suite.framework, existing);
  }

  return grouped;
};

export const exportToMarkdown = (data: AnalysisResult, options: ExportOptions = {}): string => {
  const { includeLineNumbers = true, includeTimestamp = true } = options;
  const lines: string[] = [];

  lines.push(`# ${data.owner}/${data.repo} - Test Specification`);
  lines.push("");

  if (includeTimestamp) {
    lines.push(`> Generated: ${new Date().toISOString()}`);
    lines.push(`> Analyzed: ${data.analyzedAt}`);
    if (data.branchName) {
      lines.push(`> Branch: ${data.branchName}`);
    }
    lines.push(`> Commit: ${data.commitSha}`);
    lines.push("");
  }

  lines.push("---");
  lines.push("");

  lines.push("## 📊 Summary");
  lines.push("");
  lines.push(`| Status | Count |`);
  lines.push(`|--------|-------|`);
  lines.push(`| Total | ${data.summary.total} |`);
  lines.push(`| Active | ${data.summary.active} |`);
  lines.push(`| Skipped | ${data.summary.skipped} |`);
  lines.push(`| Todo | ${data.summary.todo} |`);
  if (data.summary.focused > 0) {
    lines.push(`| Focused | ${data.summary.focused} |`);
  }
  if (data.summary.xfail > 0) {
    lines.push(`| Expected Fail | ${data.summary.xfail} |`);
  }
  lines.push("");

  lines.push("### By Framework");
  lines.push("");
  lines.push(`| Framework | Total | Active | Skipped | Todo |`);
  lines.push(`|-----------|-------|--------|---------|------|`);
  for (const fw of data.summary.frameworks) {
    lines.push(`| ${fw.framework} | ${fw.total} | ${fw.active} | ${fw.skipped} | ${fw.todo} |`);
  }
  lines.push("");

  lines.push("---");
  lines.push("");

  lines.push("## 📋 Test Suites");
  lines.push("");

  const groupedSuites = groupSuitesByFramework(data.suites);

  for (const [framework, suites] of groupedSuites) {
    lines.push(`## ${framework}`);
    lines.push("");

    for (const suite of suites) {
      lines.push(formatTestSuite(suite, { includeLineNumbers }));
      lines.push("");
    }
  }

  lines.push("---");
  lines.push("");
  lines.push(
    `*Exported from [SpecVital](https://specvital.com/analyze/${data.owner}/${data.repo})*`
  );

  return lines.join("\n");
};

export const downloadMarkdown = (content: string, filename: string): void => {
  const blob = new Blob([content], { type: "text/markdown;charset=utf-8" });
  const url = URL.createObjectURL(blob);

  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  link.click();

  URL.revokeObjectURL(url);
};

export const copyToClipboard = async (content: string): Promise<boolean> => {
  try {
    await navigator.clipboard.writeText(content);
    return true;
  } catch {
    return false;
  }
};
