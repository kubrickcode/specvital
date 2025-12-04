import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { ExternalLink, GitCommit } from "lucide-react";
import type { AnalysisResult } from "@/lib/api";
import { fetchAnalysis } from "@/lib/api";
import { isValidGitHubUrl } from "@/lib/github-url";
import { formatAnalysisDate, SHORT_SHA_LENGTH } from "@/lib/utils";
import { StatsCard } from "@/components/stats-card";
import { TestList } from "@/components/test-list";

export const dynamic = "force-dynamic";

type AnalyzePageProps = {
  params: Promise<{
    owner: string;
    repo: string;
  }>;
};

export const generateMetadata = async ({ params }: AnalyzePageProps): Promise<Metadata> => {
  const { owner, repo } = await params;
  return {
    title: `${owner}/${repo} - Test Analysis | SpecVital`,
    description: `Test specification analysis for ${owner}/${repo} repository`,
  };
};

const AnalyzePage = async ({ params }: AnalyzePageProps) => {
  const { owner, repo } = await params;

  // Validate params using the same validation as URL input
  const mockUrl = `https://github.com/${owner}/${repo}`;
  if (!isValidGitHubUrl(mockUrl)) {
    notFound();
  }

  let result: AnalysisResult;

  try {
    result = await fetchAnalysis(owner, repo);
  } catch (error) {
    // Re-throw to let error.tsx handle it with proper UI
    throw error;
  }

  return (
    <main className="container mx-auto px-4 py-8">
      <div className="space-y-6">
        <header className="space-y-2">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold">
              {result.owner}/{result.repo}
            </h1>
            <a
              href={`https://github.com/${result.owner}/${result.repo}`}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              View on GitHub
              <ExternalLink className="h-4 w-4" />
            </a>
          </div>
          <div className="flex items-center gap-4 text-sm text-muted-foreground">
            <span className="flex items-center gap-1">
              <GitCommit className="h-4 w-4" />
              {result.commitSha.slice(0, SHORT_SHA_LENGTH)}
            </span>
            <span>Analyzed: {formatAnalysisDate(result.analyzedAt)}</span>
          </div>
        </header>

        <StatsCard
          active={result.summary.active}
          label="Test Summary"
          skipped={result.summary.skipped}
          todo={result.summary.todo}
          total={result.summary.total}
        />

        <section className="space-y-4">
          <h2 className="text-xl font-semibold">Test Suites</h2>
          <TestList suites={result.suites} />
        </section>
      </div>
    </main>
  );
};

export default AnalyzePage;
