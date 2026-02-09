import type { TestStatus } from "@/lib/api/types";

import type {
  RepoSpecDocument,
  SpecBehavior,
  SpecDocument,
  SpecDomain,
  SpecFeature,
  SpecLanguage,
} from "../types";

type ExportableDocument = SpecDocument | RepoSpecDocument;

type SpecExportOptions = {
  includeSourceInfo?: boolean;
  includeTimestamp?: boolean;
};

type SpecExportMetadata = {
  owner: string;
  repo: string;
};

const STATUS_ICONS: Record<TestStatus, string> = {
  active: "✅",
  focused: "🎯",
  skipped: "⏭️",
  todo: "📝",
  xfail: "❌",
};

const getCommitSha = (document: ExportableDocument): string | undefined =>
  "commitSha" in document ? document.commitSha : undefined;

const formatBehavior = (behavior: SpecBehavior, options: SpecExportOptions): string => {
  const lines: string[] = [];
  const sourceInfo = behavior.sourceInfo;
  const statusIcon = sourceInfo ? (STATUS_ICONS[sourceInfo.status] ?? "") : "";
  const isSkipped = sourceInfo?.status === "skipped";

  const description = isSkipped
    ? `~~${behavior.convertedDescription}~~`
    : behavior.convertedDescription;

  lines.push(`- ${statusIcon} ${description}`);

  if (options.includeSourceInfo && sourceInfo) {
    lines.push(
      `  - Source: \`${sourceInfo.filePath}:${sourceInfo.lineNumber}\` (${sourceInfo.framework})`
    );
  }

  return lines.join("\n");
};

const formatFeature = (feature: SpecFeature, options: SpecExportOptions): string => {
  const lines: string[] = [];

  lines.push(`### ${feature.name}`);
  if (feature.description) {
    lines.push("");
    lines.push(`> ${feature.description}`);
  }
  lines.push("");
  lines.push("#### Behaviors");
  lines.push("");

  for (const behavior of feature.behaviors) {
    lines.push(formatBehavior(behavior, options));
  }

  return lines.join("\n");
};

const formatDomain = (domain: SpecDomain, options: SpecExportOptions): string => {
  const lines: string[] = [];

  lines.push(`## ${domain.name}`);
  if (domain.description) {
    lines.push("");
    lines.push(`> ${domain.description}`);
  }
  if (domain.classificationConfidence !== undefined) {
    lines.push("");
    lines.push(
      `> Classification Confidence: ${(domain.classificationConfidence * 100).toFixed(0)}%`
    );
  }
  lines.push("");

  for (const feature of domain.features) {
    lines.push(formatFeature(feature, options));
    lines.push("");
  }

  return lines.join("\n");
};

const calculateStats = (
  domains: SpecDomain[]
): { behaviorCount: number; domainCount: number; featureCount: number } => {
  let behaviorCount = 0;
  let featureCount = 0;

  for (const domain of domains) {
    featureCount += domain.features.length;
    for (const feature of domain.features) {
      behaviorCount += feature.behaviors.length;
    }
  }

  return { behaviorCount, domainCount: domains.length, featureCount };
};

export const exportSpecToMarkdown = (
  document: ExportableDocument,
  metadata: SpecExportMetadata,
  options: SpecExportOptions = {}
): string => {
  const { includeSourceInfo = true, includeTimestamp = true } = options;
  const lines: string[] = [];

  lines.push(`# ${metadata.owner}/${metadata.repo} - Specification Document`);
  lines.push("");

  if (includeTimestamp) {
    lines.push(`> Language: ${document.language}`);
    lines.push(`> Version: ${document.version}`);
    lines.push(`> Generated: ${document.createdAt}`);
    if (document.modelId) {
      lines.push(`> Model: ${document.modelId}`);
    }
    const commitSha = getCommitSha(document);
    if (commitSha) {
      lines.push(`> Commit: ${commitSha}`);
    }
    lines.push("");
  }

  if (document.executiveSummary) {
    lines.push("## Executive Summary");
    lines.push("");
    lines.push(document.executiveSummary);
    lines.push("");
  }

  lines.push("---");
  lines.push("");

  const stats = calculateStats(document.domains);
  lines.push("## 📊 Overview");
  lines.push("");
  lines.push("| Domains | Features | Behaviors |");
  lines.push("|---------|----------|-----------|");
  lines.push(`| ${stats.domainCount} | ${stats.featureCount} | ${stats.behaviorCount} |`);
  lines.push("");

  lines.push("---");
  lines.push("");

  for (const domain of document.domains) {
    lines.push(formatDomain(domain, { includeSourceInfo }));
  }

  lines.push("---");
  lines.push("");
  lines.push(
    `*Exported from [SpecVital](https://specvital.com/${metadata.owner}/${metadata.repo})*`
  );

  return lines.join("\n");
};

export const generateSpecFilename = (owner: string, repo: string, language: SpecLanguage): string =>
  `${owner}-${repo}-spec-${language}.md`;

export const downloadSpecMarkdown = (content: string, filename: string): void => {
  const blob = new Blob([content], { type: "text/markdown;charset=utf-8" });
  const url = URL.createObjectURL(blob);

  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  link.click();

  URL.revokeObjectURL(url);
};

export const copySpecToClipboard = async (content: string): Promise<boolean> => {
  try {
    await navigator.clipboard.writeText(content);
    return true;
  } catch {
    return false;
  }
};
