import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { isValidGitHubUrl } from "@/lib/github-url";

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

  return (
    <main className="container mx-auto px-4 py-8">
      <div className="space-y-6">
        <header className="space-y-2">
          <h1 className="text-2xl font-bold">
            {owner}/{repo}
          </h1>
          <p className="text-muted-foreground">Test Specification Analysis</p>
        </header>

        <div className="rounded-lg border bg-card p-6">
          <p className="text-muted-foreground">
            Analysis results will be displayed here in Commit 4.
          </p>
        </div>
      </div>
    </main>
  );
};

export default AnalyzePage;
